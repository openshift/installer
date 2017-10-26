# frozen_string_literal: true

require 'cluster'
require 'env_var'

# GCPCluster represents a k8s cluster on CGP cloud provider
class GcpCluster < Cluster
  def env_variables
    variables = super
    variables['PLATFORM'] = 'gcp'
    variables
  end

  def check_prerequisites
    raise 'GCP credentials not defined' unless credentials_defined?

    super
  end

  def credentials_defined?
    credential_vars = %w[
      GOOGLE_CREDENTIALS
      GOOGLE_CLOUD_KEYFILE_JSON
      GCLOUD_KEYFILE_JSON
      GOOGLE_APPLICATION_CREDENTIALS
    ]
    EnvVar.contains_any?(credential_vars)
  end

  def master_ip_addresses
    ip_addresses = []
    Dir.chdir(@build_path) do
      ip_address = `echo module.network.master_ip | terraform console ../../platforms/gcp`.chomp
      if ip_address.empty?
        raise 'should get the master_ip_address to use in the tests.'
      end
      ip_addresses.push(ip_address)
      ip_addresses
    end
  end

  def master_ip_address
    master_ip_addresses[0]
  end

  def tectonic_console_url
    Dir.chdir(@build_path) do
      console_url = `echo module.dns.kube_ingress_fqdn | terraform console ../../platforms/gcp`.chomp
      if console_url.empty?
        raise 'should get the console url to use in the UI tests.'
      end
      console_url
    end
  end
end
