# frozen_string_literal: true

require 'shared_examples/k8s'
require 'govcloud_vpc'
require 'aws_region'
require 'jenkins'
require 'aws_iam'

RSpec.describe 'govcloud-vpc' do
  include_examples(
    'withBuildFolderSetup',
    File.join(ENV['RSPEC_PATH'], '../smoke/govcloud/vars/govcloud-vpc-internal.tfvars.json')
  )

  before(:all) do
    @role_credentials = nil
    @role_credentials = AWSIAM.assume_role('us-gov-west-1') if ENV.key?('TECTONIC_INSTALLER_ROLE')
    @ssh_key = ENV['TF_VAR_tectonic_govcloud_ssh_key'] || AwsSupport.create_aws_key_pairs('us-gov-west-1',
                                                                                          @role_credentials)
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
