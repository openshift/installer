# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-standard' do
  include_examples('withRunningCluster', '../smoke/aws/vars/aws.tfvars.json')
end
