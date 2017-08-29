require 'cluster'
require 'aws_region'

# AWSCluster represents a k8s cluster on AWS cloud provider
class AWSCluster < Cluster
  def initialize(prefix, tf_vars_path)
    export_random_region_if_not_defined

    super(prefix, tf_vars_path)
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
