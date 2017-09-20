# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-external-experimental' do
  include_examples('withRunningCluster', '../smoke/azure/vars/external-experimental.tfvars')
end
