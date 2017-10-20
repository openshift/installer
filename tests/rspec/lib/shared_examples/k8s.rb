# frozen_string_literal: true

require 'smoke_test'
require 'forensic'
require 'cluster_factory'
require 'container_linux'
require 'operators'
require 'pages/login_page'
require 'name_generator'
require 'password_generator'
require 'webdriver_helpers'

RSpec.shared_examples 'withRunningCluster' do |tf_vars_path|
  before(:all) do
    @cluster = ClusterFactory.from_tf_vars(TFVarsFile.new(tf_vars_path))
    begin
      @cluster.start
    rescue => e
      Forensic.run(@cluster)
      raise "Aborting execution due startup failed. Error: #{e}"
    end
  end

  after(:each) do |example|
    Forensic.run(@cluster) if example.exception
  end

  after(:all) do
    @cluster.stop
  end

  it 'generates operator manifests' do
    expect { Operators.manifests_generated?(@cluster.manifest_path) }
      .to_not raise_error
  end

  it 'succeeds with the golang test suit' do
    expect { SmokeTest.run(@cluster) }.to_not raise_error
  end

  it 'installs the correct Container Linux version' do
    version = @cluster.tf_var('tectonic_container_linux_version')
    version = @cluster.tf_value('module.container_linux.version') if version == 'latest'
    expect(ContainerLinux.version(@cluster)).to eq(version)
  end

  it 'installs the correct Container Linux channel' do
    expect(ContainerLinux.channel(@cluster)).to eq(@cluster.tf_var('tectonic_container_linux_channel'))
  end

  describe 'Interact with tectonic console' do
    before(:all) do
      @driver = WebdriverHelpers.start_webdriver
      @login = Login.new(@driver)
      @console_url = @cluster.tectonic_console_url
    end

    after(:all) do
      WebdriverHelpers.stop_webdriver(@driver)
    end

    it 'can login in the tectonic console', :ui, retry: 3, retry_wait: 10 do
      @login.login_page "https://#{@console_url}"
      @login.with(@cluster.tectonic_admin_email, @cluster.tectonic_admin_password)
      expect(@login.success_login?).to be_truthy
      @login.logout
    end

    it 'fail to login with wrong credentials', :ui, retry: 3, retry_wait: 10 do
      @login.login_page "https://#{@console_url}"
      @login.with(NameGenerator.generate_fake_email, PasswordGenerator.generate_password)
      expect(@login.fail_to_login?).to be_truthy
    end
  end
end

RSpec.shared_examples 'withPlannedCluster' do |tf_vars_path|
  before(:all) do
    @cluster = ClusterFactory.from_tf_vars(TFVarsFile.new(tf_vars_path))
  end

  it 'terraform plan succeeds' do
    @cluster.plan
  end

  after(:all) do
    @cluster.stop
  end
end
