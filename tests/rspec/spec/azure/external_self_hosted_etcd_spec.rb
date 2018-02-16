# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-external-self-hosted-etcd' do
  include_examples(
    'withRunningCluster',
    File.join(ENV['RSPEC_PATH'], '../smoke/azure/vars/external-self-hosted-etcd.tfvars')
  )
end
