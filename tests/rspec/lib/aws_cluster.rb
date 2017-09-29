# frozen_string_literal: true

require 'cluster'
require 'aws_region'
require 'json'
require 'jenkins'
require 'grafiti'
require 'env_var'
require 'aws_iam'
require 'aws_support'

# AWSCluster represents a k8s cluster on AWS cloud provider
class AwsCluster < Cluster
  def initialize(tfvars_file)
    if Jenkins.environment?
      export_random_region_if_not_defined
      # AWSIAM.assume_role
    end
    @aws_region = tfvars_file.tectonic_aws_region
    @aws_ssh_key = ENV['TF_VAR_tectonic_aws_ssh_key'] = AwsSupport.create_aws_key_pairs(@aws_region)
    super(tfvars_file)
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'aws'
    variables
  end

  def stop
    AwsSupport.delete_aws_key_pairs(@aws_ssh_key, @aws_region)

    super
  end

  def master_ip_address
    ssh_master_ip = nil
    Dir.chdir(@build_path) do
      terraform_state = `terraform state show module.masters.aws_autoscaling_group.masters`.chomp.split("\n")
      terraform_state.each do |value|
        attributes = value.split('=')
        next unless attributes[0].strip.eql?('id')
        instances_id = AwsSupport.sorted_auto_scaling_instances(attributes[1].strip.chomp, @aws_region)
        ssh_master_ip = AwsSupport.preferred_instance_ip_address(instances_id[0], @aws_region)
        break
      end
    end
    ssh_master_ip
  end

  def check_prerequisites
    raise 'AWS credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_aws_ssh_key is not defined' unless ssh_key_defined?
    raise 'TF_VAR_tectonic_aws_region is not defined' unless region_defined?

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
    EnvVar.set?(%w[TF_VAR_tectonic_aws_ssh_key])
  end

  def recover_from_failed_destroy
    Grafiti.new(@build_path, ENV['TF_VAR_tectonic_aws_region']).clean
    super
  end

  def tectonic_console_url
    Dir.chdir(@build_path) do
      ingress_ext = `echo module.masters.ingress_external_fqdn | terraform console ../../platforms/aws`.chomp
      ingress_int = `echo module.masters.ingress_internal_fqdn | terraform console ../../platforms/aws`.chomp
      if ingress_ext.empty?
        if ingress_int.empty?
          raise 'should get the console url to use in the UI tests.'
        end
        return ingress_int
      end
      ingress_ext
    end
  end
end
