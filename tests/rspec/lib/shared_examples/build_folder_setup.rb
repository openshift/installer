# frozen_string_literal: true

require 'fileutils'
require 'name_generator'

RSpec.shared_examples 'withBuildFolderSetup' do |tf_vars_path|
  before(:all) do
    temp_tfvars_file = TFVarsFile.new(tf_vars_path)
    @name = ENV['CLUSTER'] || NameGenerator.generate(temp_tfvars_file.prefix)
    ENV['CLUSTER'] = @name
    file_path = "../../build/#{@name}"
    FileUtils.mkdir_p file_path

    FileUtils.cp(
      tf_vars_path,
      Dir.pwd + "/#{file_path}/terraform.tfvars"
    )
    @tfvars_file = TFVarsFile.new("#{file_path}/terraform.tfvars")
  end
end
