# frozen_string_literal: true

require 'cluster'
require 'json'
require 'jenkins'
require 'env_var'
require 'metal_support'

# BareMetalCluster represents a k8s cluster on AWS cloud provider
class MetalCluster < Cluster
  def env_variables
    variables = super
    variables['PLATFORM'] = 'metal'
    variables
  end

  def master_ip_addresses
    master_ip_address = []
    master_ip_address.push(tf_value('var.tectonic_metal_controller_domains[0]'))
  end

  def master_ip_address
    master_ip_addresses[0]
  end

  def tectonic_console_url
    console_url = tf_value('var.tectonic_metal_ingress_domain')
    if console_url.empty?
      raise 'failed to get the console url to use in the UI tests.'
    end
    console_url
  end
end
