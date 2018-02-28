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

  def self.from_config_file(config_file)
    cluster_class = "#{config_file.platform.downcase.capitalize}Cluster"
    Object.const_get(cluster_class).new(config_file)
  end

  def self.from_variable(cluster_type, config_file)
    Object.const_get("#{cluster_type.downcase.capitalize}Cluster").new(config_file)
  end
end
