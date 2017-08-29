require 'shared_examples/k8s'
require 'aws_vpc'
require 'aws_region'

RSpec.describe 'aws-vpc' do
  before(:all) do
    export_random_region_if_not_defined
    @vpc = AWSVPC.new('test-vpc')
    @vpc.create
  end

  context 'with a cluster' do
    include_examples(
      'withCluster',
      '../smoke/aws/vars/aws-vpc-internal.tfvars.json'
    )
  end

  after(:all) do
    @vpc.destroy
  end
end
