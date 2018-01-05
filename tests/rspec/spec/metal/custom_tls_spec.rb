# frozen_string_literal: true

require 'shared_examples/k8s'
require 'tls_certs'
require 'metal_support'

DOMAIN = 'example.com'

RSpec.describe 'metal-custom-tls' do
  include_examples('withBuildFolderSetup', '../smoke/bare-metal/vars/metal.tfvars.json')
  include_examples('withTLSSetup', DOMAIN)

  before(:all) do
    MetalSupport.install_base_software
    MetalSupport.setup_bare(@tfvars_file)
    MetalSupport.start_matchbox(@tfvars_file)
  end

  after(:context) do |_context|
    MetalSupport.destroy
  end

  context 'with a cluster' do
    include_examples('withRunningClusterExistingBuildFolder')
  end
end
