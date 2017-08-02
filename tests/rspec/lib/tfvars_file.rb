require 'json'

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

  private

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
