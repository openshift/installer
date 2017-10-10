# frozen_string_literal: true

require 'kubectl_helpers'
require 'securerandom'
require 'jenkins'
require 'tfvars_file'
require 'fileutils'
require 'name_generator'
require 'password_generator'
require 'ssh'

SSH_CMD_BOOTKUBE_DONE = 'systemctl is-active bootkube'
SSH_CMD_TECTONIC_DONE = 'systemctl is-active tectonic'

# Cluster represents a k8s cluster
class Cluster
  attr_reader :tfvars_file, :kubeconfig, :manifest_path, :build_path,
              :tectonic_admin_email, :tectonic_admin_password

  def initialize(tfvars_file)
    @tfvars_file = tfvars_file

    # Enable local testers to specify a static cluster name
    # S3 buckets can only handle lower case names
    @name = ENV['CLUSTER'] || NameGenerator.generate(tfvars_file.prefix)
    @tectonic_admin_email = ENV['TF_VAR_tectonic_admin_email'] || NameGenerator.generate_fake_email
    @tectonic_admin_password = ENV['TF_VAR_tectonic_admin_password'] || PasswordGenerator.generate_password

    @build_path = File.join(File.realpath('../../'), "build/#{@name}")
    @manifest_path = File.join(@build_path, 'generated')
    @kubeconfig = File.join(manifest_path, 'auth/kubeconfig')
    check_prerequisites
    localconfig
    prepare_assets
  end

  def plan
    succeeded = system(env_variables, 'make -C ../.. plan')
    raise 'Planning cluster failed' unless succeeded
  end

  def start
    apply
    wait_til_ready
  end

  def stop
    if ENV.key?('TECTONIC_TESTS_DONT_CLEAN_UP')
      puts "*** Cleanup inhibiting flag set. Stopping here. ***\n"
      puts '*** Your email/password to use in the tectonic console is:'\
           "#{@tectonic_admin_email} / #{@tectonic_admin_password} ***\n"
      return
    end
    destroy
    clean if Jenkins.environment?
  end

  def check_prerequisites
    return if license_and_pull_secret_defined?
    raise 'Tectonic license and pull secret are not defined as environment'\
          'variables.'
  end

  def env_variables
    {
      'CLUSTER' => @name,
      'TF_VAR_tectonic_cluster_name' => @name,
      'TF_VAR_tectonic_admin_email' => @tectonic_admin_email,
      'TF_VAR_tectonic_admin_password' => @tectonic_admin_password
    }
  end

  def tf_var(v)
    tf_value "var.#{v}"
  end

  def tf_value(v)
    Dir.chdir(@build_path) do
      `echo "#{v}" | terraform console ../../platforms/#{env_variables['PLATFORM']}`.chomp
    end
  end

  private

  def license_and_pull_secret_defined?
    license_path = 'TF_VAR_tectonic_license_path'
    pull_secret_path = 'TF_VAR_tectonic_pull_secret_path'

    EnvVar.set?([license_path, pull_secret_path])
  end

  def prepare_assets
    FileUtils.cp(
      @tfvars_file.path,
      Dir.pwd + "/../../build/#{@name}/terraform.tfvars"
    )
  end

  def localconfig
    succeeded = system(env_variables, 'make -C ../.. localconfig')
    raise 'Run localconfig failed' unless succeeded
  end

  def apply
    3.times do
      return if system(env_variables, 'make -C ../.. apply')
    end
    raise 'Applying cluster failed'
  end

  def destroy
    3.times do
      return if system(env_variables, 'make -C ../.. destroy')
    end

    recover_from_failed_destroy
    raise 'Destroying cluster failed'
  end

  def recover_from_failed_destroy() end

  def clean
    succeeded = system(env_variables, 'make -C ../.. clean')
    raise 'could not clean build directory' unless succeeded
  end

  def wait_til_ready
    from = Time.now
    loop do
      begin
        KubeCTL.run(@kubeconfig, 'cluster-info')
        break
      rescue KubeCTL::KubeCTLCmdError
        elapsed = Time.now - from
        raise 'kubectl cluster-info never returned with successful error code' if elapsed > 1200 # 20 mins timeout
        sleep 10
      end
    end
    wait_for_bootstrapping
  end

  def wait_for_bootstrapping
    from = Time.now
    loop do
      _, _, bootkube_exitstatus = ssh_exec(master_ip_address, SSH_CMD_BOOTKUBE_DONE)
      _, _, tectonic_exitstatus = ssh_exec(master_ip_address, SSH_CMD_TECTONIC_DONE)
      break if bootkube_exitstatus.zero? && tectonic_exitstatus.zero?
      elapsed = Time.now - from
      puts 'Waiting for bootstrapping to complete...' if (elapsed.round % 5).zero?
      raise 'timeout waiting for bootstrapping' if elapsed > 1200 # 20 mins timeout
      sleep 2
    end
    puts 'HOORAY! The cluster is up'
  end
end
