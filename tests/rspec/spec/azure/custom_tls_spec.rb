# frozen_string_literal: true

require 'shared_examples/k8s'
require 'azure_support'

RSpec.describe 'azure-custom-tls' do
  context 'with a cluster' do
    azure_region = ENV['TF_VAR_tectonic_azure_location'] || AzureSupport.random_location_unless_defined
    ENV['TF_VAR_tectonic_azure_location'] = azure_region
    @domain = "#{azure_region}.cloudapp.azure.com"

    include_examples('withRunningClusterWithCustomTLS', '../smoke/azure/vars/basic.tfvars', @domain, false)
  end
end
