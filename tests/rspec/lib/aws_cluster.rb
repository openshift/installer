# frozen_string_literal: true

require 'cluster'
require 'aws_region'

# AWSCluster represents a k8s cluster on AWS cloud provider

class AwsCluster < Cluster
  def initialize(tfvars_file)
    export_random_region_if_not_defined

    super(tfvars_file)
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'aws'
    variables
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
    EnvVar.set?(credential_names)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_ssh_key])
  end
end
