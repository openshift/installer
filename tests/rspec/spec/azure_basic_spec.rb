# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-basic' do
  include_examples('withRunningCluster', '../smoke/azure/vars/basic.tfvars')
end
