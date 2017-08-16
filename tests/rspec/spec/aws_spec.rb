require 'shared_examples/k8s'

RSpec.describe 'aws-standard' do
  include_examples('withCluster', '../smoke/aws/vars/aws.tfvars.json')
end
