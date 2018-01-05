# frozen_string_literal: true

require 'shared_examples/k8s'
require 'govcloud_vpc'
require 'aws_region'
require 'jenkins'
require 'aws_iam'

RSpec.describe 'govcloud-vpc' do
  include_examples('withBuildFolderSetup', '../smoke/govcloud/vars/govcloud-vpc-internal.tfvars.json')

  before(:all) do
    @ssh_key = ENV['TF_VAR_tectonic_govcloud_ssh_key'] || AwsSupport.create_aws_key_pairs('us-gov-west-1')
    ENV['TF_VAR_tectonic_govcloud_ssh_key'] = @ssh_key
    ENV['TF_VAR_ssh_key'] = @ssh_key

    # AWSIAM.assume_role if Jenkins.environment?
    @vpc = GovcloudVPC.new('test-vpc-govcloud')
    @vpc.create
  end

  context 'with a cluster' do
    include_examples('withRunningClusterExistingBuildFolder')
  end

  after(:all) do
    @vpc.destroy
  end
end
