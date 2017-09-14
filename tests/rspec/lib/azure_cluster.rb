# frozen_string_literal: true

require 'cluster'
require 'azure_support'

# AzureCluster represents a k8s cluster on Azure cloud provider
class AzureCluster < Cluster
  extend AzureSupport

  def env_variables
    variables = super
    variables['PLATFORM'] = 'azure'
    variables['TF_VAR_tectonic_azure_location'] = AzureSupport.random_location_unless_defined
    variables
  end

  def check_prerequisites
    raise 'Azure credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_azure_ssh_key is not defined' unless ssh_key_defined?

    super
  end

  def credentials_defined?
    credential_names = %w[
      ARM_SUBSCRIPTION_ID ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_TENANT_ID
    ]
    EnvVar.set?(credential_names)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_azure_ssh_key])
  end
end
