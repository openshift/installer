# frozen_string_literal: true

require 'shared_examples/k8s'

RSpec.describe 'azure-dns' do
  include_examples('withRunningCluster', '../smoke/azure/vars/dns.tfvars')
end
