# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-example' do
  include_examples('withRunningCluster', '../smoke/azure/vars/example.tfvars')
end
