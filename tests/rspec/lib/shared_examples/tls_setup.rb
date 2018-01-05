# frozen_string_literal: true

require 'fileutils'
require 'tls_certs'

RSpec.shared_examples 'withTLSSetup' do |domain|
  before(:all) do
    test_folder = File.expand_path('..', Dir.pwd)
    generate_tls("#{test_folder}/smoke/user_provided_tls/certs/", @name, domain, @tfvars_file.etcd_count)

    root_folder = File.expand_path('../..', Dir.pwd)
    custom_tls_tf = "#{test_folder}/smoke/user_provided_tls/tls.tf"
    dest_folder = "#{root_folder}/platforms/#{@tfvars_file.platform}"
    original_tls_tf = "#{dest_folder}/tls.tf"

    FileUtils.mv(original_tls_tf, "#{dest_folder}/tls.tf.original")
    FileUtils.cp(custom_tls_tf, dest_folder)
  end
end
