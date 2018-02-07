# frozen_string_literal: true

require 'shared_examples/k8s'

DOMAIN = 'gcp.tectonic-ci.de'

RSpec.describe 'gcp-custom-tls' do
  context 'with a cluster' do
    include_examples(
      'withRunningClusterWithCustomTLS',
      File.join(ENV['RSPEC_PATH'], '../smoke/gcp/vars/gcp.tfvars.json'),
      DOMAIN,
      false
    )
  end
end
