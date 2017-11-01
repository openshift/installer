# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-self-hosted-etcd' do
  include_examples('withRunningCluster', '../smoke/azure/vars/self-hosted-etcd.tfvars')
end
