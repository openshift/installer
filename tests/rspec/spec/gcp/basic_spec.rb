# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'gcp-standard' do
  include_examples('withRunningCluster', '../smoke/gcp/vars/gcp.tfvars.json')
end
