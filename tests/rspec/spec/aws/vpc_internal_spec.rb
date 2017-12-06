# frozen_string_literal: true

require 'shared_examples/k8s'
require 'aws_vpc'
require 'aws_region'
require 'jenkins'
require 'aws_iam'

RSpec.describe 'aws-vpc' do
  include_examples('withBuildFolderSetup', '../smoke/aws/vars/aws-vpc-internal.tfvars.json')

  before(:all) do
    export_random_region_if_not_defined
    # AWSIAM.assume_role if Jenkins.environment?
    @vpc = AWSVPC.new('test-vpc')
    @vpc.create
  end

  context 'with a cluster' do
    include_examples('withRunningClusterExistingBuildFolder', true)
  end

  after(:all) do
    @vpc.destroy
  end
end
