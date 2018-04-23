# frozen_string_literal: true

require 'fileutils'
require 'name_generator'

RSpec.shared_examples 'withBuildFolderSetupWithConfig' do |config_path|
  before(:all) do
    Dir.chdir(File.join(File.dirname(ENV['RELEASE_TARBALL_PATH']), 'tectonic-dev'))
    # TODO: Only ignore on AWS
    temp_config_file = ConfigFile.new(config_path)
    @name = ENV['CLUSTER'] || NameGenerator.generate(temp_config_file.prefix)
    ENV['CLUSTER'] = @name

    FileUtils.cp(
      config_path,
      'config.yaml'
    )
    @config_file = ConfigFile.new(File.expand_path('config.yaml'))
  end
end

RSpec.shared_examples 'withBuildFolderSetup' do |tf_vars_path|
  before(:all) do
    Dir.chdir(File.join(File.dirname(ENV['RELEASE_TARBALL_PATH']), 'tectonic-dev'))
    temp_tfvars_file = TFVarsFile.new(tf_vars_path)
    @name = ENV['CLUSTER'] || NameGenerator.generate(temp_tfvars_file.prefix)
    ENV['CLUSTER'] = @name
    file_path = "build/#{@name}"
    FileUtils.mkdir_p file_path
    Dir.chdir(file_path)

    FileUtils.cp(
      tf_vars_path,
      'terraform.tfvars'
    )
    @tfvars_file = TFVarsFile.new(File.expand_path('terraform.tfvars'))
  end
end
