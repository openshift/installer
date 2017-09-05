# frozen_string_literal: true

require 'kubectl_helpers'
require 'securerandom'
require 'jenkins'
require 'tfvars_file'
require 'fileutils'

# Cluster represents a k8s cluster
class Cluster
  MAX_NAME_LENGTH = 28
  RANDOM_HASH_LENGTH = 5

  attr_reader :tfvars_file, :kubeconfig, :manifest_path, :build_path

  def initialize(tfvars_file)
    @tfvars_file = tfvars_file

    # Enable local testers to specify a static cluster name
    # S3 buckets can only handle lower case names
    @name = (ENV['CLUSTER'] || generate_name(tfvars_file.prefix)).downcase

    @build_path = File.join(File.realpath('../../'), "build/#{@name}")
    @manifest_path = File.join(@build_path, 'generated')
    @kubeconfig = File.join(manifest_path, 'auth/kubeconfig')
  end

  def start
    check_prerequisites
    localconfig
    prepare_assets
    plan
    apply
    wait_til_ready
  end

  def stop
    destroy
    clean
  end

  def check_prerequisites
    return if license_and_pull_secret_defined?
    raise 'Tectonic license and pull secret are not defined as environment'\
          'variables.'
  end

  def env_variables
    {
      'CLUSTER' => @name,
      'TF_VAR_tectonic_cluster_name' => @name
    }
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

  def plan
    succeeded = system(env_variables, 'make -C ../.. plan')
    raise 'Planning cluster failed' unless succeeded
  end

  def apply
    succeeded = system(env_variables, 'make -C ../.. apply')
    raise 'Applying cluster failed' unless succeeded
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
    90.times do
      begin
        KubeCTL.run(@kubeconfig, 'cluster-info')
        return
      rescue KubeCTL::KubeCTLCmdError
        sleep 10
      end
    end

    raise 'kubectl cluster-info never returned with successful error code'
  end

  def generate_name(prefix)
    name = prefix

    if Jenkins.environment?
      build_id = ENV['BUILD_ID']
      branch_name = ENV['BRANCH_NAME']
      name = "#{prefix}-#{branch_name}-#{build_id}"
    end

    name = name[0..(MAX_NAME_LENGTH - RANDOM_HASH_LENGTH - 1)]
    name += SecureRandom.hex[0...RANDOM_HASH_LENGTH]
    name
  end
end
