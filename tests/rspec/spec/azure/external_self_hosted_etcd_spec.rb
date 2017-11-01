# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-external-self-hosted-etcd' do
  include_examples('withRunningCluster', '../smoke/azure/vars/external-self-hosted-etcd.tfvars')
end
