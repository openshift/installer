# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-experimental' do
  include_examples('withRunningCluster', '../smoke/azure/vars/experimental.tfvars')
end
