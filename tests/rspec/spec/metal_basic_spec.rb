# frozen_string_literal: true

require 'shared_examples/k8s'
require 'metal_support'

TEST_CLUSTER_CONFIG_FILE = '../smoke/bare-metal/vars/metal.tfvars.json'

RSpec.describe 'bare-metal-standard' do
  before(:context) do |_context|
    @varfile = TFVarsFile.new(TEST_CLUSTER_CONFIG_FILE)
    ENV['CLUSTER'] ||= NameGenerator.generate(@varfile.prefix)
    MetalSupport.install_base_software
    MetalSupport.setup_bare(@varfile)
    MetalSupport.start_matchbox
  end

  after(:context) do |_context|
    MetalSupport.destroy
  end

  context 'running Tectonic' do
    include_examples('withRunningCluster', TEST_CLUSTER_CONFIG_FILE)
  end
end
