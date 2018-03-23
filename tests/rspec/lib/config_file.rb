# frozen_string_literal: true

require 'yaml'

PLATFORMS = %w[govcloud aws azure metal vmware gcp].freeze

# ConfigFile represents a Terraform configuration file describing a Tectonic
# cluster configuration
class ConfigFile
  attr_reader :path
  def initialize(file_path)
    @path = file_path
    raise "file #{file_path} does not exist" unless file_exists?
  end

  def data
    YAML.safe_load(File.read(@path))
  end

  def networking
    data['networking']['type']
  end

  def node_count
    master_count + worker_count
  end

  def master_count
    get_node_count(data['master']['nodePools'])
  end

  def worker_count
    get_node_count(data['worker']['nodePools'])
  end

  def etcd_count
    get_node_count(data['etcd']['nodePools'])
  end

  def add_worker_node(node_count)
    new_data = set_node_count(data['worker']['nodePools'][0], node_count)
    save(new_data)
  end

  def change_cluster_name(cluster_name)
    new_data = data
    new_data['name'] = cluster_name
    save(new_data)
  end

  def cluster_name
    data['name']
  end

  def change_aws_region(region)
    new_data = data
    new_data['aws']['region'] = region
    save(new_data)
  end

  def region(platform)
    data[platform]['region']
  end

  def change_license(license_path)
    new_data = data
    new_data['licensePath'] = license_path
    save(new_data)
  end

  def change_pull_secret(pull_secret_path)
    new_data = data
    new_data['pullSecretPath'] = pull_secret_path
    save(new_data)
  end

  def change_base_domain(base_domain)
    new_data = data
    new_data['baseDomain'] = base_domain
    save(new_data)
  end

  def license
    data['licensePath']
  end

  def pull_secret
    data['pullSecretPath']
  end

  def change_admin_credentials(admin_email, admin_passwd)
    new_data = data
    new_data['admin'] ||= {}
    new_data['admin']['email'] = admin_email
    new_data['admin']['password'] = admin_passwd
    save(new_data)
  end

  def admin_credentials
    admin_email = data.dig('admin', 'email')
    admin_passwd = data.dig('admin', 'password')
    [admin_email, admin_passwd]
  end

  def prefix
    prefix = File.basename(path).split('.').first
    raise 'could not extract prefix from tfvars file name' if prefix == ''
    prefix
  end

  def change_ssh_key(platform, ssh_key)
    new_data = data
    new_data[platform]['sshKey'] = ssh_key
    save(new_data)
  end

  def platform
    PLATFORMS.each do |plat|
      return plat if data['platform'].downcase.eql?(plat)
    end
  end

  private

  def file_exists?
    File.exist? path
  end

  def save(data)
    File.open(path, 'w+') do |f|
      f << data.to_yaml
    end
  end

  def get_node_count(names)
    count = 0
    names.each do |name|
      data['nodePools'].each do |n|
        count += n['count'] if n['name'] == name
      end
    end
    count
  end

  def set_node_count(name, count)
    d = data
    d['nodePools'].each do |n|
      if n['name'] == name
        n['count'] = count
        break
      end
    end
    d
  end
end
