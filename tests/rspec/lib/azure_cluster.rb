# frozen_string_literal: true

require 'cluster'
require 'azure_support'
require 'net/ssh'

SSH_CMD_BOOTKUBE_DONE = "journalctl --no-pager -u bootkube | grep -q 'Started Bootstrap a Kubernetes cluster.'"
SSH_CMD_TECTONIC_DONE = "journalctl --no-pager -u tectonic | grep -q 'Started Bootstrap a Tectonic cluster.'"

# AzureCluster represents a k8s cluster on Azure cloud provider
#
class AzureCluster < Cluster
  extend AzureSupport

  def start
    super
    # Wait for bootstrapping to complete
    wait_for_bootstrapping
  end

  def wait_for_bootstrapping
    ssh_ip = master_ip_address
    from = Time.now
    Net::SSH.start(ssh_ip, 'core') do |ssh|
      loop do
        puts 'Waiting for bootstrapping to complete...'
        raise 'timeout waiting for bootstrapping' if Time.now - from > 1200 # 20 mins timeout
        bootkube_done = ssh.exec!(SSH_CMD_BOOTKUBE_DONE).exitstatus.zero?
        tectonic_done = ssh.exec!(SSH_CMD_TECTONIC_DONE).exitstatus.zero?
        break if bootkube_done && tectonic_done
        sleep(5)
      end
    end
    puts 'HOORAY! The cluster is up'
  end

  def master_ip_address
    Dir.chdir(@build_path) do
      `echo 'module.vnet.api_ip_addresses[0]' | terraform console ../../platforms/azure`.chomp
    end
  end

  def env_variables
    variables = super
    variables['PLATFORM'] = 'azure'
    variables['TF_VAR_tectonic_azure_location'] = AzureSupport.random_location_unless_defined
    variables
  end

  def check_prerequisites
    raise 'Azure credentials not defined' unless credentials_defined?
    raise 'TF_VAR_tectonic_azure_ssh_key is not defined' unless ssh_key_defined?

    super
  end

  def credentials_defined?
    credential_names = %w[
      ARM_SUBSCRIPTION_ID ARM_CLIENT_ID ARM_CLIENT_SECRET ARM_TENANT_ID
    ]
    EnvVar.set?(credential_names)
  end

  def ssh_key_defined?
    EnvVar.set?(%w[TF_VAR_tectonic_azure_ssh_key])
  end
end
