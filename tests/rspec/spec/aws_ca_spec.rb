# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-custom-ca' do
  include_examples(
    'withPlannedCluster',
    '../smoke/aws/vars/aws-ca.tfvars.json'
  )
end
