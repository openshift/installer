# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-exp' do
  include_examples(
    'withRunningCluster',
    '../smoke/aws/vars/aws-exp.tfvars.json'
  )
end
