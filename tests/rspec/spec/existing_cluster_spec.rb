# frozen_string_literal: true

require 'shared_examples/k8s'
require 'env_var'

RSpec.describe 'existing_cluster' do
  before(:all) do
    raise 'Missing existing platform. Please set PLATFORM environment variable' unless EnvVar.set?(%w[PLATFORM])
    raise 'Missing existing platform. Please set CLUSTER environment variable' unless EnvVar.set?(%w[CLUSTER])
  end

  context 'with a cluster' do
    vars_file_path = "../../build/#{ENV['CLUSTER']}/terraform.tfvars"
    raise 'Missing tfvars. Aborting...' unless File.exist?(vars_file_path)
    existing_vars_file = TFVarsFile.new(vars_file_path)

    cred_file = "../../build/#{ENV['CLUSTER']}/utils/console_creds.txt"
    if File.exist?(cred_file)
      creds = File.read(cred_file)
      creds_hash = JSON.parse(creds)
      tectonic_admin_email = creds_hash['tectonic_admin_email']
      tectonic_admin_password = creds_hash['tectonic_admin_password']
    elsif existing_vars_file.tectonic_admin_email && existing_vars_file.tectonic_admin_password
      tectonic_admin_email = existing_vars_file.tectonic_admin_email
      tectonic_admin_password = existing_vars_file.tectonic_admin_password
    end

    ENV['TF_VAR_tectonic_admin_email'] = tectonic_admin_email
    ENV['TF_VAR_tectonic_admin_password'] = tectonic_admin_password

    include_examples('withRunningClusterExistingBuildFolder',
                     false,
                     ENV['PLATFORM'],
                     existing_vars_file)
  end
end
