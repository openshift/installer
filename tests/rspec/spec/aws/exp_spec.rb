# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'aws-exp' do
  include_examples(
    'withRunningCluster',
    File.join(ENV['RSPEC_PATH'], '../smoke/aws/vars/aws-exp.tfvars.json')
  )
end
