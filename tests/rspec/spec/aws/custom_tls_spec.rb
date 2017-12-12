# frozen_string_literal: true

require 'shared_examples/k8s'

DOMAIN = 'tectonic-ci.de'

RSpec.describe 'aws-custom-tls' do
  context 'with a cluster' do
    include_examples('withRunningClusterWithCustomTLS', '../smoke/aws/vars/aws.tfvars.json', DOMAIN, false)
  end
end
