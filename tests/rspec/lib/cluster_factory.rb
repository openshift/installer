# frozen_string_literal: true

require 'aws_cluster'
require 'govcloud_cluster'
require 'azure_cluster'
require 'metal_cluster'
require 'gcp_cluster'

# Creates a platform specific Cluster object based on the provided tfvars file.
#
module ClusterFactory
  def self.from_tf_vars(tf_vars_file)
    cluster_class = "#{tf_vars_file.platform.downcase.capitalize}Cluster"
    Object.const_get(cluster_class).new(tf_vars_file)
  end
end
