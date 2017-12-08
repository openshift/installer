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

  def worker_ip_addresses
    workers_ip_address = []
    workers = tf_value('var.tectonic_metal_worker_domains')

    # the output is "[\n  node2.example.com,\n  node3.example.com\n]"
    workers = workers.delete("\n[]").split(',')
    workers.each do |value|
      value.strip!
      workers_ip_address.push(value)
    end
    workers_ip_address
  end

  def etcd_ip_addresses
    # for metal etcd are in the same server as the master.
    master_ip_addresses
  end

  def tectonic_console_url
    console_url = tf_value('var.tectonic_metal_ingress_domain')
    if console_url.empty?
      raise 'failed to get the console url to use in the UI tests.'
    end
    console_url
  end
end
