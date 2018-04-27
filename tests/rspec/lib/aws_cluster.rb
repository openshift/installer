# frozen_string_literal: true

require 'cluster'
require 'aws_region'
require 'json'
require 'jenkins'
require 'grafiti'
require 'env_var'
require 'aws_iam'
require 'aws_support'
require 'tfstate_file'
require 'fileutils'
require 'with_retries'

# AWSCluster represents a k8s cluster on AWS cloud provider
class AwsCluster < Cluster
  TIMEOUT_IN_SECONDS = (30 * 60).freeze # 30 minutes

  attr_reader :config_file, :kubeconfig, :manifest_path, :build_path,
              :tectonic_admin_email, :tectonic_admin_password, :tfstate_file

  def initialize(config_file)
    @config_file = config_file
    export_random_region_if_not_defined if Jenkins.environment?
    @aws_region = ENV['TF_VAR_tectonic_aws_region']

    @name = @config_file.cluster_name
    @build_path = File.join(File.dirname(ENV['RELEASE_TARBALL_PATH']), "tectonic-dev/#{@name}")
    @manifest_path = File.join(@build_path, 'generated')
    @kubeconfig = File.join(@build_path, 'generated/auth/kubeconfig')

    @role_credentials = nil
    @role_credentials = AWSIAM.assume_role(@aws_region) if ENV.key?('TECTONIC_INSTALLER_ROLE')

    unless ssh_key_defined?
      ENV['TF_VAR_tectonic_aws_ssh_key'] = AwsSupport.create_aws_key_pairs(@aws_region, @role_credentials)
    end

    @config_file.change_aws_region(@aws_region)
    @config_file.change_license(ENV['TF_VAR_tectonic_license_path'])
    @config_file.change_pull_secret(ENV['TF_VAR_tectonic_pull_secret_path'])
    @config_file.change_ssh_key(@config_file.platform, ENV['TF_VAR_tectonic_aws_ssh_key'])

    @tectonic_admin_email = NameGenerator.generate_fake_email if @config_file.admin_credentials[0].nil?
    @tectonic_admin_password = PasswordGenerator.generate_password if @config_file.admin_credentials[1].nil?
    @config_file.change_admin_credentials(@tectonic_admin_email, @tectonic_admin_password)

    @tfstate_file = TFStateFile.new(@build_path, 'bootstrap.tfstate')
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'aws'
    variables['TF_VAR_tectonic_cluster_name'] = @config_file.cluster_name
    variables['CLUSTER'] = @config_file.cluster_name

    # Unless base domain is provided by the user:
    unless ENV.key?('TF_VAR_tectonic_base_domain')
      variables['TF_VAR_tectonic_base_domain'] = 'tectonic-ci.de'
      @config_file.change_base_domain('tectonic-ci.de')
    end

    variables
  end

  def stop
    if ENV['TF_VAR_tectonic_aws_ssh_key'].include?('rspec-')
      AwsSupport.delete_aws_key_pairs(ENV['TF_VAR_tectonic_aws_ssh_key'], @aws_region, @role_credentials)
    end

    super
  end

  def machine_boot_console_logs
    instances_id = retrieve_instances_ids('module.masters.aws_autoscaling_group.masters')
    # Return the log output in a hash {ip => log}
    hash_log_ip = instances_id.map do |instance_id|
      {
        instance_id_to_ip_address(instance_id) =>
        AwsSupport.collect_ec2_console_logs(instance_id, @aws_region, @role_credentials)
      }
    end
    # convert the array to hash [{k1=>v1},{k2=>v2}] to {k1=>v1,k2=>v2}
    hash_log_ip.reduce({}, :update)
  end

  def retrieve_instances_ids(auto_scaling_groups)
    aws_autoscaling_group_master = @tfstate_file.value(auto_scaling_groups, 'id')
    AwsSupport.sorted_auto_scaling_instances(aws_autoscaling_group_master, @aws_region, @role_credentials)
  end

  def instance_id_to_ip_address(instance_id)
    AwsSupport.instance_ip_address(instance_id, @aws_region, @role_credentials)
  end

  def master_ip_addresses
    instances_id = retrieve_instances_ids('module.masters.aws_autoscaling_group.masters')
    instances_id.map { |instance_id| AwsSupport.instance_ip_address(instance_id, @aws_region, @role_credentials) }
  end

  def master_ip_address
    master_ip_addresses[0]
  end

  def worker_ip_addresses
    instances_id = retrieve_instances_ids('module.workers.aws_autoscaling_group.workers')
    instances_id.map { |instance_id| AwsSupport.instance_ip_address(instance_id, @aws_region, @role_credentials) }
  end

  def etcd_ip_addresses
    @tfstate_file.output('etcd', 'ip_addresses')
  end

  def check_prerequisites
    raise 'AWS credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_aws_ssh_key is not defined' unless ssh_key_defined?
    raise 'TF_VAR_tectonic_aws_region is not defined' unless region_defined?

    super
  end

  def region_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_region])
  end

  def credentials_defined?
    credential_names = %w[AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY]
    profile_name = %w[AWS_PROFILE]
    session_token = %w[
      AWS_ACCESS_KEY_ID
      AWS_SECRET_ACCESS_KEY
      AWS_SESSION_TOKEN
    ]
    EnvVar.set?(credential_names) ||
      EnvVar.set?(profile_name) ||
      EnvVar.set?(session_token)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_aws_ssh_key])
  end

  def recover_from_failed_destroy
    Grafiti.new(@build_path, ENV['TF_VAR_tectonic_aws_region']).clean
    super
  end

  def tectonic_console_url
    Dir.chdir(@build_path) do
      ingress_ext = @tfstate_file.output('dns', 'ingress_external_fqdn')
      ingress_int = @tfstate_file.output('dns', 'ingress_internal_fqdn')
      if ingress_ext.empty?
        if ingress_int.empty?
          raise 'failed to get the console url to use in the UI tests.'
        end
        return ingress_int
      end
      ingress_ext
    end
  end

  # TODO: Remove once other platforms caught up

  def init
    env = env_variables
    env['TF_INIT_OPTIONS'] = '-no-color'

    run_tectonic_cli(env, 'init', '--config=config.yaml')
  end

  def apply
    Retriable.with_retries(limit: 3) do
      env = env_variables
      env['TF_APPLY_OPTIONS'] = '-no-color'
      env['TF_INIT_OPTIONS'] = '-no-color'

      run_tectonic_cli(env, 'install', "--dir=#{@name}")
    end
  end

  def destroy
    describe_network_interfaces
    Retriable.with_retries(limit: 3) do
      env = env_variables
      env['TF_DESTROY_OPTIONS'] = '-no-color'
      env['TF_INIT_OPTIONS'] = '-no-color'
      run_tectonic_cli(env, 'destroy', "--dir=#{@name}")
    end

    recover_from_failed_destroy
    raise 'Destroying cluster failed'
  rescue => e
    recover_from_failed_destroy
    raise e
  end

  def run_tectonic_cli(env, cmd, flags = '')
    tectonic_binary = File.join(
      File.dirname(ENV['RELEASE_TARBALL_PATH']),
      'tectonic-dev/installer/tectonic'
    )

    tectonic_logs = File.join(
      File.dirname(ENV['RELEASE_TARBALL_PATH']),
      "tectonic-dev/#{@name}/logs/tectonic-#{cmd}.log"
    )

    ::Timeout.timeout(TIMEOUT_IN_SECONDS) do
      command = "#{tectonic_binary} #{cmd} #{flags}"
      output = ''
      Open3.popen3(env, "bash -coxe pipefail '#{command}'") do |_stdin, stdout, stderr, wait_thr|
        puts "Only printing tectonic logs to stdout/stderr on failure for command: #{command}.\nLogs are preserved via log files."

        while (line = stdout.gets)
          output += line
        end
        while (line = stderr.gets)
          output += line
        end

        save_terraform_logs(tectonic_logs, output)
        unless wait_thr.value.success?
          puts output
          raise "failed to execute command: #{command}"
        end
      end
    end
  rescue Timeout::Error
    save_terraform_logs(tectonic_logs, output)
    forensic(false)
    raise 'Applying cluster failed'
  end

  def wait_nodes_ready
    from = Time.now
    loop do
      puts 'Waiting for nodes become in ready state after an update'
      Retriable.with_retries(KubeCTL::KubeCTLCmdError, limit: 5, sleep: 10) do
        nodes = describe_nodes
        nodes_ready = Array.new(@config_file.node_count, false)
        nodes['items'].each_with_index do |item, index|
          item['status']['conditions'].each do |condition|
            if condition['type'] == 'Ready' && condition['status'] == 'True'
              nodes_ready[index] = true
            end
          end
        end
        if nodes_ready.uniq.length == 1 && nodes_ready.uniq.include?(true)
          puts '**All nodes are Ready!**'
          return true
        end
        puts "One or more nodes are not ready yet or missing nodes. Waiting...\n" \
             "# of returned nodes #{nodes['items'].size}. Expected #{@config_file.node_count}"
        elapsed = Time.now - from
        raise 'waiting for all nodes to become ready timed out' if elapsed > 1200 # 20 mins timeout
        sleep 20
      end
    end
  end

  # TODO: (carlos) remove this
  def tf_var(v)
    tf_value "var.#{v}"
  end

  # TODO: (carlos) remove this
  def tf_value(v)
    Dir.chdir(@build_path) do
      `echo '#{v}' | terraform console ../steps/bootstrap`.chomp
    end
  end

  private

  # def destroy
  #   # For debugging purposes (see: https://github.com/terraform-providers/terraform-provider-aws/pull/1051)
  #   describe_network_interfaces

  #   super

  #   # For debugging purposes (see: https://github.com/terraform-providers/terraform-provider-aws/pull/1051)
  #   describe_network_interfaces
  # end

  def describe_network_interfaces
    puts 'describing network interfaces for debugging purposes'
    vpc_id = @tfstate_file.value('module.vpc.aws_vpc.cluster_vpc', 'id')
    filter = "--filters=Name=vpc-id,Values=#{vpc_id}"
    region = "--region #{@aws_region}"

    # TODO: use aws sdk instead of command line
    success = system("aws ec2 describe-network-interfaces #{filter}  #{region}")
    raise 'failed to describe network interfaces by vpc' unless success

  # Do not fail build. This is only for debugging purposes
  rescue => e
    puts e
  end

  def save_terraform_logs(tectonic_logs, output)
    # Save output in logs/
    FileUtils.mkdir_p(File.dirname(tectonic_logs))
    save_to_file = File.open(tectonic_logs, 'a')
    save_to_file << output
    save_to_file.close
  end
end
