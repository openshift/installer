# frozen_string_literal: true

require 'json'

PLATFORMS = %w[govcloud aws azure metal vmware gcp].freeze

# TFVarsFile represents a Terraform configuration file describing a Tectonic
# cluster configuration
class TFVarsFile
  attr_reader :path, :data
  def initialize(file_path)
    @path = file_path
    raise "file #{file_path} does not exist" unless file_exists?
    @data = JSON.parse(File.read(path))
  end

  def networking
    data['tectonic_networking']
  end

  def node_count
    master_count + worker_count
  end

  def master_count
    count = if platform.eql?('metal')
              data['tectonic_metal_controller_names'].count
            else
              data['tectonic_master_count'].to_i
            end
    count
  end

  def worker_count
    count = if platform.eql?('metal')
              data['tectonic_metal_worker_names'].count
            else
              data['tectonic_worker_count'].to_i
            end
    count
  end

  def etcd_count
    data['tectonic_etcd_count'].to_i
  end

  def add_worker_node(node_count)
    data['tectonic_worker_count'] = node_count.to_s
    save
  end

  def change_cluster_name(cluster_name)
    data['tectonic_cluster_name'] = cluster_name
    save
  end

  def change_dns_name(dns_name)
    data['tectonic_dns_name'] = dns_name
    save
  end

  def region
    data['tectonic_aws_region']
  end

  def prefix
    prefix = File.basename(path).split('.').first
    raise 'could not extract prefix from tfvars file name' if prefix == ''
    prefix
  end

  def platform
    PLATFORMS.each do |platform|
      return platform if data.keys.any? do |key|
        key.start_with?("tectonic_#{platform}")
      end
    end
  end

  private

  def method_missing(method_name, *arguments, &block)
    data.fetch(method_name.to_s)
  rescue KeyError
    return ENV["TF_VAR_#{method_name}"] if ENV.key?("TF_VAR_#{method_name}")
    super
  end

  def respond_to_missing?(method_name, include_private = false)
    data.keys.any? { |k| k == method_name.to_s } || super
  end

  def file_exists?
    File.exist? path
  end

  def save
    File.open(path, 'w+') do |f|
      f << data.to_json
    end
  end
end
