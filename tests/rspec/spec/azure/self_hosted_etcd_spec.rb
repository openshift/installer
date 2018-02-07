# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-self-hosted-etcd' do
  include_examples('withRunningCluster', File.join(ENV['RSPEC_PATH'], '../smoke/azure/vars/self-hosted-etcd.tfvars'))
end
