# frozen_string_literal: true

require 'cluster'
require 'aws_region'
require 'json'
require 'jenkins'
require 'grafiti'
require 'env_var'
require 'aws_iam'
require 'aws_support'
require 'tfstate_file'

# AWSCluster represents a k8s cluster on AWS cloud provider
class GovcloudCluster < Cluster
  def initialize(tfvars_file)
    if Jenkins.environment?
      export_random_region_if_not_defined
      # AWSIAM.assume_role
    end
    @aws_region = 'us-gov-west-1'

    unless ssh_key_defined?
      ENV['TF_VAR_tectonic_govcloud_ssh_key'] = AwsSupport.create_aws_key_pairs(@aws_region)
    end

    super(tfvars_file)
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'govcloud'

    # Unless base domain is provided by the user:
    unless ENV.key?('TF_VAR_tectonic_base_domain')
      variables['TF_VAR_tectonic_base_domain'] = 'tectonic-ci.de'
    end

    variables
  end

  def stop
    if ENV['TF_VAR_tectonic_govcloud_ssh_key'].include?('rspec-')
      AwsSupport.delete_aws_key_pairs(ENV['TF_VAR_tectonic_govcloud_ssh_key'], @aws_region)
    end

    super
  end

  def master_ip_addresses
    ssh_master_ip_addresses = []
    Dir.chdir(@build_path) do
      terraform_state = `terraform state show module.masters.aws_autoscaling_group.masters`.chomp.split("\n")
      terraform_state.each do |value|
        attributes = value.split('=')
        next unless attributes[0].strip.eql?('id')
        instances_id = AwsSupport.sorted_auto_scaling_instances(attributes[1].strip.chomp, @aws_region)
        instances_id.each do |instance_id|
          ssh_master_ip_addresses.push AwsSupport.preferred_instance_ip_address(instance_id, @aws_region)
        end
      end
    end
    ssh_master_ip_addresses
  end

  def master_ip_address
    master_ip_addresses[0]
  end

  def check_prerequisites
    raise 'AWS credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_govcloud_ssh_key is not defined' unless ssh_key_defined?
    super
  end

  def region_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_region])
  end

  def credentials_defined?
    credential_names = %w[AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY]
    profile_name = %w[AWS_PROFILE]
    session_token = %w[
      AWS_ACCESS_KEY_ID
      AWS_SECRET_ACCESS_KEY
      AWS_SESSION_TOKEN
    ]
    EnvVar.set?(credential_names) ||
      EnvVar.set?(profile_name) ||
      EnvVar.set?(session_token)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_govcloud_ssh_key])
  end

  def recover_from_failed_destroy
    Grafiti.new(@build_path, @aws_region).clean
    super
  end

  def tectonic_console_url
    Dir.chdir(@build_path) do
      ingress_ext = `echo module.dns.ingress_external_fqdn | terraform console ../../platforms/govcloud`.chomp
      ingress_int = `echo module.dns.ingress_internal_fqdn | terraform console ../../platforms/govcloud`.chomp
      if ingress_ext.empty?
        if ingress_int.empty?
          raise 'failed to get the console url to use in the UI tests.'
        end
        return ingress_int
      end
      ingress_ext
    end
  end

  private

  def destroy
    # For debugging purposes (see: https://github.com/terraform-providers/terraform-provider-aws/pull/1051)
    describe_network_interfaces

    super

    # For debugging purposes (see: https://github.com/terraform-providers/terraform-provider-aws/pull/1051)
    describe_network_interfaces
  end

  def describe_network_interfaces
    puts 'describing network interfaces for debugging purposes'
    vpc_id = @tfstate_file.value('module.vpc.aws_vpc.cluster_vpc', 'id')
    filter = "--filters=Name=vpc-id,Values=#{vpc_id}"
    region = "--region #{@aws_region}"

    success = system("aws ec2 describe-network-interfaces #{filter}  #{region}")
    raise 'failed to describe network interfaces by vpc' unless success

  # Do not fail build. This is only for debugging purposes
  rescue => e
    puts e
  end
end
