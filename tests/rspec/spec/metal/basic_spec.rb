# frozen_string_literal: true

require 'shared_examples/k8s'
require 'metal_support'

RSpec.describe 'bare-metal-standard' do
  include_examples('withBuildFolderSetup', '../smoke/bare-metal/vars/metal.tfvars.json')

  before(:context) do |_context|
    MetalSupport.install_base_software
    MetalSupport.setup_bare(@tfvars_file)
    MetalSupport.start_matchbox(@tfvars_file)
  end

  after(:context) do |_context|
    MetalSupport.destroy
  end

  context 'running Tectonic' do
    include_examples('withRunningClusterExistingBuildFolder')
  end
end
