# frozen_string_literal: true

require 'json'
require 'securerandom'
require 'net/http'
require 'uri'
require 'openssl'
require 'time'
require 'azure_support'

AZURE_VPN_TEMPLATES = '../smoke/azure/fixtures/private-cluster/*.tf'

# Creates a VNET with an OpenVPN gateway instance
# and exposes the connection details for consumption as an external VNET.
#
class AzureVpn
  attr_reader :vpn_gw_endpoint
  attr_reader :ovpn_login
  attr_reader :vnet_id
  attr_reader :rsg_id
  attr_reader :build_dir
  attr_reader :env_config

  extend AzureSupport

  def initialize(tf_vars_file)
    raise 'Invalid tfvars file' if tf_vars_file.nil?

    @ovpn_login = {
      username: "v#{SecureRandom.hex(8)}", password: SecureRandom.urlsafe_base64
    }
    @build_dir = Dir.mktmpdir('tectonic-test-')
    @env_config = {
      'TF_VAR_admin_username' => ovpn_login[:username],
      'TF_VAR_admin_password' => ovpn_login[:password],
      'TF_VAR_tectonic_cluster_name' => ENV['CLUSTER'],
      'TF_VAR_tectonic_azure_location' => AzureSupport.random_location_unless_defined
    }
    Dir.glob(AZURE_VPN_TEMPLATES) { |tf| FileUtils.cp(tf, build_dir) }
    FileUtils.cp(tf_vars_file.path, "#{build_dir}/terraform.tfvars")
    File.write(File.join(build_dir, 'vpn_credentials'),
               "#{ovpn_login[:username]}\n#{ovpn_login[:password]}\n")
  end

  def start
    create_resources
    parse_terraform_output
    wait_for_vpn_access_server
    connect_to_vpn
    wait_for_vpn_connection
    import_vnet_resolver
  end

  def stop
    if ENV.key?('TECTONIC_TESTS_DONT_CLEAN_UP')
      print 'Cleanup inhibiting flag set. Stopping here.'
      return
    end
    Dir.chdir(build_dir) do
      unless system(env_config, 'terraform destroy -force')
        raise 'Private vnet: destroy failed'
      end
    end
    system("rm -rf #{build_dir}")
  end

  def wait_for_vpn_access_server
    90.times do
      puts 'Waiting for VPN access server to boot...'
      return if system("curl -ksL https://#{vpn_gw_endpoint}:443/ >/dev/null")
      sleep(0.5)
    end
    raise 'waiting for vpn access server timed out'
  end

  def wait_for_vpn_connection
    gwip = gw_private_ip_address
    from = Time.now
    loop do
      puts 'Waiting for VPN to connect...'
      break if system("ping -c 1 #{gwip}")
      raise 'waiting for vpn connection timed out' if Time.now - from > 120
    end
  end

  def query_vpn_config
    return '' if @vpn_gw_endpoint.empty?
    uri = URI.parse("https://#{vpn_gw_endpoint}:443/rest/GetUserlogin")
    request = Net::HTTP::Get.new(uri)
    request.basic_auth(@ovpn_login[:username], @ovpn_login[:password])
    req_options = {
      use_ssl: uri.scheme == 'https',
      verify_mode: OpenSSL::SSL::VERIFY_NONE
    }
    Net::HTTP.start(uri.hostname, uri.port, req_options) do |http|
      http.request(request)
    end
  end

  def vpn_config
    response = query_vpn_config
    raise 'Unable to get VPN config.' unless response.is_a?(Net::HTTPSuccess)
    Dir.chdir(build_dir) do
      File.open('config.ovpn', 'w') do |file|
        len = file.write(
          response.body.sub('auth-user-pass', 'auth-user-pass vpn_credentials')
        )
        raise 'Did not receive OpenVPN client config.' unless len.positive?
        puts 'Received OpenVPN client config.'
      end
    end
  end

  def parse_terraform_output
    Dir.chdir(build_dir) do
      tf_out = JSON.parse(`terraform output -json`)
      @vpn_gw_endpoint = tf_out['vpn-gw-endpoint']['value']
      @vnet_id = tf_out['tectonic_azure_external_vnet_id']['value']
      @rsg_id = tf_out['tectonic_azure_external_resource_group']['value']
    end
  end

  def create_resources
    Dir.chdir(build_dir) do
      unless system(env_config, 'terraform init')
        raise 'Private vnet: init failed'
      end
      unless system(env_config, 'terraform apply')
        raise 'Private vnet: apply failed'
      end
    end
  end

  def connect_to_vpn
    vpn_config
    Dir.chdir(build_dir) do
      unless system('openvpn --daemon --config config.ovpn')
        raise 'Couldn\'t start OpenVPN client. Aborting.'
      end
    end
  end

  def import_vnet_resolver
    resolv_conf = `ssh -o 'StrictHostKeyChecking no' #{ovpn_login[:username]}@#{vpn_gw_endpoint} 'cat /etc/resolv.conf'`
    resolv_lines = resolv_conf.split("\n")
    last_ns = resolv_lines.rindex { |l| l =~ /^nameserver/ }
    resolv_lines.insert(last_ns + 1, 'nameserver 8.8.8.8')
    resolv_lines.insert(last_ns + 2, 'nameserver 8.8.4.4')
    File.write('/etc/resolv.conf', resolv_lines.join("\n"))
  end

  def gw_private_ip_address
    Dir.chdir(build_dir) do
      gw_atrributes = `terraform state show azurerm_network_interface.vpn_gw`.split("\n")
      gw_atrributes.select { |attr| attr =~ /^private_ip_address/ }.first.split('=')[1].strip
    end
  end
end
