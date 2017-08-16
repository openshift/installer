require 'cluster'

# AWSCluster represents a k8s cluster on AWS cloud provider
class AWSCluster < Cluster
  def env_variables
    variables = super
    variables['PLATFORM'] = 'aws'
    variables
  end

  def check_prerequisites
    raise 'AWS credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_aws_ssh_key is not defined' unless ssh_key_defined?

    super
  end

  def credentials_defined?
    credential_names = %w[AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY]
    EnvVar.set?(credential_names)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_ssh_key])
  end
end
