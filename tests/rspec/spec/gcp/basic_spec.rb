# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'gcp-basic' do
  include_examples('withRunningCluster', File.join(ENV['RSPEC_PATH'], '../smoke/gcp/vars/gcp.tfvars.json'))
end
