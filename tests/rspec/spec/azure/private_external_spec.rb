# frozen_string_literal: true

require 'shared_examples/k8s'
require 'azure_vpn'

RSpec.describe 'azure-private-external' do
  include_examples('withBuildFolderSetup', '../smoke/azure/vars/private-cluster.tfvars')

  before(:context) do |_context|
    # Save environment, to restore it once the test it done
    # (since we alter it further down)
    @curent_env = ENV.clone
    @azure_ssh_key_path = AzureSupport.set_ssh_key_path
    @vpn_vnet = AzureVpn.new(@tfvars_file)
    @vpn_vnet.start
    # Pick up the VNET and resource group created for the VPN
    # and pass it to the installer as external resources
    ENV['TF_VAR_tectonic_azure_external_vnet_id'] = @vpn_vnet.vnet_id
    ENV['TF_VAR_tectonic_azure_external_resource_group'] = @vpn_vnet.rsg_id
    ENV['TF_VAR_tectonic_azure_location'] = @vpn_vnet.env_config['TF_VAR_tectonic_azure_location']
  end

  after(:context) do |_context|
    @vpn_vnet.stop
    # Restore environment to the state we found before the test.
    ENV.clear
    ENV = @curent_env.clone
  end

  context 'private cluster' do
    include_examples('withRunningClusterExistingBuildFolder', true)

    it 'should have a vpn' do
    end
  end
end
