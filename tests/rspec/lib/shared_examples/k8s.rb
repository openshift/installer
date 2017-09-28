# frozen_string_literal: true

require 'smoke_test'
require 'forensic'
require 'cluster_factory'
require 'operators'
require 'pages/login_page'
require 'name_generator'
require 'password_generator'
require 'webdriver_helpers'

RSpec.shared_examples 'withRunningCluster' do |tf_vars_path|
  before(:all) do
    @cluster = ClusterFactory.from_tf_vars(TFVarsFile.new(tf_vars_path))
    @cluster.start
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

  describe 'Interact with tectonic console' do
    before(:all) do
      @driver = WebdriverHelpers.start_webdriver
      @login = Login.new(@driver)
      console_url = @cluster.tectonic_console_url
      @login.login_page "https://#{console_url}"
    end

    after(:all) do
      WebdriverHelpers.stop_webdriver(@driver)
    end

    it 'can login in the tectonic console' do
      @login.with(@cluster.tectonic_admin_email, @cluster.tectonic_admin_password)
      expect(@login.success_login?).to be_truthy
      @login.logout
    end

    it 'fail to login with wrong credentials' do
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
