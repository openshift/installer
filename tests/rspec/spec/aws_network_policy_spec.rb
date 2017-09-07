# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-network-policy' do
  include_examples(
    'withRunningCluster',
    '../smoke/aws/vars/aws-net-policy.tfvars.json'
  )
end
