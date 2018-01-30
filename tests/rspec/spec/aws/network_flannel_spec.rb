# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-network-flannel' do
  include_examples(
    'withRunningCluster',
    '../smoke/aws/vars/aws-net-flannel.tfvars.json'
  )
end
