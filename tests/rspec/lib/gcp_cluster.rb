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
    EnvVar.contains_any(credential_vars)
  end
end
