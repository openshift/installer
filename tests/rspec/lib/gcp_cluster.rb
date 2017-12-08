# frozen_string_literal: true

require 'cluster'
require 'env_var'
require 'gcloud_helper'

# GCPCluster represents a k8s cluster on GCP cloud provider
class GcpCluster < Cluster
  def initialize(tfvars_file)
    super(tfvars_file)
    @gcloud = GcloudHelper.new
  end

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
    project_vars = %w[
      GOOGLE_PROJECT
    ]
    EnvVar.contains_any?(credential_vars) && EnvVar.set?(project_vars)
  end

  def master_ip_addresses
    gcloud_command = "compute instances list \
--format='value(networkInterfaces[0].accessConfigs[0].natIP)' \
--filter='name~#{@name}.*master.*'"
    ip_addresses = @gcloud.run(gcloud_command).split("\n")
    ip_addresses
  end

  def master_ip_address
    master_ip_addresses[0]
  end

  def worker_ip_addresses
    gcloud_command = "compute instances list \
--format='value(networkInterfaces[0].accessConfigs[0].natIP)' \
--filter='name~#{@name}.*worker.*'"
    ip_addresses = @gcloud.run(gcloud_command).split("\n")
    ip_addresses
  end

  def etcd_ip_addresses
    @tfstate_file.output('etcd', 'etcd_ip_addresses')
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
