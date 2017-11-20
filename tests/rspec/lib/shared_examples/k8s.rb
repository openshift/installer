# frozen_string_literal: true

require 'shared_examples/build_folder_setup'
require 'smoke_test'
require 'forensic'
require 'cluster_factory'
require 'container_linux'
require 'operators'
require 'pages/login_page'
require 'name_generator'
require 'password_generator'
require 'webdriver_helpers'
require 'k8s_conformance_tests'

RSpec.shared_examples 'withRunningCluster' do |tf_vars_path|
  include_examples('withBuildFolderSetup', tf_vars_path)
  include_examples('withRunningClusterExistingBuildFolder')
end

RSpec.shared_examples 'withRunningClusterExistingBuildFolder' do
  before(:all) do
    @cluster = ClusterFactory.from_tf_vars(@tfvars_file)
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

  it 'verifies api checkpoint manifests' do
    @cluster.master_ip_addresses.each do |ip|
      cmd = "sudo sh -c 'cat /etc/kubernetes/inactive-manifests/kube-system-kube-apiserver-*.json'"

      retries = 0
      begin
        stdout, _, fin = ssh_exec(ip, cmd, 20)
        raise unless fin.zero?
        expect { JSON.parse(stdout) }.to_not raise_error
      rescue
        retries += 1
        expect(retries).to be < 20
        puts "failed to exec '#{cmd}'; retrying in 3 seconds"
        sleep 3
        retry
      end
    end
  end

  it 'verifies api secret checkpoints' do
    secrets = @cluster.secret_files('kube-system', 'kube-apiserver')

    @cluster.master_ip_addresses.each do |ip|
      path = '/etc/kubernetes/checkpoint-secrets/kube-system/kube-apiserver-*/kube-apiserver'
      cmd = secrets
            .map { |secret| "test -e #{path}/#{secret}" }
            .join(' && ')
      cmd = "sudo sh -c '#{cmd}'"

      retries = 0
      begin
        _, _, fin = ssh_exec(ip, cmd, 20)
        raise unless fin.zero?
      rescue
        retries += 1
        expect(retries).to be < 20
        puts "failed to exec '#{cmd}'; retrying in 3 seconds"
        sleep 3
        retry
      end
    end
  end

  it 'succeeds with the golang test suit', :smoke_tests do
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
      WebdriverHelpers.stop_webdriver(@driver) if defined? @driver
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

  describe 'scale up worker cluster' do
    before(:all) do
      platform = @cluster.env_variables['PLATFORM']
      # remove platform AZURE when the JIRA https://jira.prod.coreos.systems/browse/INST-619 is fixed
      skip_platform = %w[metal azure]
      skip "This test is not ready to run in #{platform}" if skip_platform.include?(platform)
    end

    it 'can scale up nodes by 1 worker' do
      @cluster.tfvars_file.add_worker_node(@cluster.tfvars_file.worker_count + 1)

      expect { @cluster.update_cluster }.to_not raise_error
    end
  end

  it 'passes the k8s conformance tests', :conformance_tests do
    conformance_test = K8sConformanceTest.new(@cluster.kubeconfig, vpn_tunnel)
    expect { conformance_test.run }.to_not raise_error
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
