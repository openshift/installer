# frozen_string_literal: true

require 'cluster'
require 'azure_support'

# AzureCluster represents a k8s cluster on Azure cloud provider
#
class AzureCluster < Cluster
  extend AzureSupport

  def initialize(tfvars_file)
    @random_location = AzureSupport.random_location_unless_defined
    @azure_ssh_key_path = AzureSupport.set_ssh_key_path unless EnvVar.set?(%w[TF_VAR_tectonic_azure_ssh_key])
    super(tfvars_file)
  end

  def master_ip_address
    Dir.chdir(@build_path) do
      `echo 'module.vnet.api_ip_addresses[0]' | terraform console ../../platforms/azure`.chomp
    end
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'azure'
    variables['TF_VAR_tectonic_azure_location'] = @random_location
    variables['TF_VAR_tectonic_azure_client_secret'] = ENV['ARM_CLIENT_SECRET']

    # To use Azure-provided DNS, `tectonic_base_domain` should be set to `""`
    unless ENV.key?('TF_VAR_tectonic_base_domain')
      variables['TF_VAR_tectonic_base_domain'] = ''
    end

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

  def tectonic_console_url
    Dir.chdir(@build_path) do
      console_url = `echo module.vnet.ingress_fqdn | terraform console ../../platforms/azure`.chomp
      if console_url.empty?
        raise 'should get the console url to use in the UI tests.'
      end
      console_url
    end
  end
end
