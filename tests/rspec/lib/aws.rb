require 'env_var'

# AWS contains helper functions for the AWS cloud provider
module AWS
  def self.check_prerequisites
    raise 'AWS credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_aws_ssh_key is not defined' unless ssh_key_defined?
  end

  def self.credentials_defined?
    credential_names = %w[AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY]
    EnvVar.set?(credential_names)
  end

  def self.ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_ssh_key])
  end
end
