# frozen_string_literal: true

require 'json'

PLATFORMS = %w[aws azure metal vmware gcp].freeze

# TFVarsFile represents a Terraform configuration file describing a Tectonic
# cluster configuration
class TFVarsFile
  attr_reader :path, :data
  def initialize(file_path)
    @path = file_path
    raise "file #{file_path} does not exist" unless file_exists?
    @data = JSON.parse(File.read(path))
  end

  def experimental?
    data['tectonic_experimental'] == 'true'
  end

  def calico?
    data['tectonic_calico_network_policy'] == 'true'
  end

  def node_count
    master_count + worker_count
  end

  # TODO: Randomize region on AWS
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
    super
  end

  def respond_to_missing?(method_name, include_private = false)
    data.keys.any? { |k| k == method_name.to_s } || super
  end

  def master_count
    data['tectonic_master_count'].to_i
  end

  def worker_count
    data['tectonic_worker_count'].to_i
  end

  def file_exists?
    File.exist? path
  end
end
