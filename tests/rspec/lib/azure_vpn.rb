# frozen_string_literal: true

require 'json'
require 'securerandom'
require 'net/http'
require 'uri'
require 'openssl'
require 'time'
require 'net/ssh'
require 'azure_support'

AZURE_VPN_TEMPLATES = '../smoke/azure/fixtures/private-cluster/*.tf'
MAX_RETRIES = 3
TIMEOUT_RETRIES = 30
TIMEOUT_RETRY_DELAY = 3

# Creates a VNET with an OpenVPN gateway instance
# and exposes the connection details for consumption as an external VNET.
#
class AzureVpn
  attr_reader :vnet_id
  attr_reader :rsg_id
  attr_reader :env_config

  extend AzureSupport

  def initialize(tf_vars_file)
    raise 'Invalid tfvars file' if tf_vars_file.nil?
    @vnet_cidr_block = tf_vars_file.tectonic_azure_vnet_cidr_block
    raise 'No CIDR specified for VNet.' if @vnet_cidr_block.empty?
    @ovpn_login = {
      username: "v#{SecureRandom.hex(8)}", password: SecureRandom.urlsafe_base64
    }
    @env_config = {
      'TF_VAR_admin_username' => @ovpn_login[:username],
      'TF_VAR_admin_password' => @ovpn_login[:password],
      'TF_VAR_tectonic_cluster_name' => ENV['CLUSTER'],
      'TF_VAR_tectonic_azure_location' => AzureSupport.random_location_unless_defined
    }
    @build_dir = Dir.mktmpdir('tectonic-test-')
    Dir.glob(AZURE_VPN_TEMPLATES) { |tf| FileUtils.cp(tf, @build_dir) }
    FileUtils.cp(tf_vars_file.path, "#{@build_dir}/terraform.tfvars")
    File.write(File.join(@build_dir, 'vpn_credentials'),
               "#{@ovpn_login[:username]}\n#{@ovpn_login[:password]}\n")
  end

  def start
    create_resources
    parse_terraform_output
    wait_for_vpn_access_server
    query_vpn_config
    connect_to_vpn
    wait_for_vpn_connection
    wait_for_ssh
    import_vnet_resolver
  end

  def create_resources
    commands = ['terraform init', 'terraform apply -auto-approve']

    Dir.chdir(@build_dir) do
      commands.each do |command|
        MAX_RETRIES.times do |count|
          break if system(env_config, command)
          raise "Private vnet: #{command} failed #{MAX_RETRIES} times" if MAX_RETRIES == count + 1
        end
      end
    end
  end

  def parse_terraform_output
    Dir.chdir(@build_dir) do
      tf_out = JSON.parse(`terraform output -json`)
      @vpn_gw_endpoint = tf_out['vpn_gw_endpoint']['value']
      @vpn_gw_private_ip = tf_out['vpn_gw_private_ip']['value']
      @vpn_gw_dns_servers = tf_out['vpn_gw_dns_servers']['value']
      @vnet_id = tf_out['tectonic_azure_external_vnet_id']['value']
      @rsg_id = tf_out['tectonic_azure_external_resource_group']['value']
    end
    raise 'No private IP address for gateway.' if @vpn_gw_private_ip.empty?
  end

  def wait_for_vpn_access_server
    180.times do
      puts 'Waiting for VPN access server to boot...'
      return if system("curl -ksL https://#{@vpn_gw_endpoint}:443/ >/dev/null")
      sleep(1)
    end
    raise 'waiting for vpn access server timed out'
  end

  def query_vpn_config
    return '' if @vpn_gw_endpoint.empty?
    uri = URI.parse("https://#{@vpn_gw_endpoint}:443/rest/GetUserlogin")
    response = Net::HTTP.start(
      uri.hostname, uri.port,
      use_ssl: uri.scheme == 'https', verify_mode: OpenSSL::SSL::VERIFY_NONE
    ) do |http|
      request = Net::HTTP::Get.new(uri)
      request.basic_auth(@ovpn_login[:username], @ovpn_login[:password])
      http.request(request)
    end
    raise 'Unable to get VPN config.' unless response.is_a?(Net::HTTPSuccess)
    Dir.chdir(@build_dir) do
      File.open('config.ovpn', 'w') do |file|
        len = file.write(
          response.body.sub('auth-user-pass', 'auth-user-pass vpn_credentials')
        )
        raise 'Did not receive OpenVPN client config.' unless len.positive?
        puts 'Received OpenVPN client config.'
      end
    end
  end

  def connect_to_vpn
    Dir.chdir(@build_dir) do
      raise 'Cannot start OpenVPN client.' unless system('openvpn --daemon --config config.ovpn --log vpnclient.log')
    end
    puts 'Started OpenVPN client.'
  end

  def wait_for_vpn_connection
    from = Time.now
    loop do
      puts 'Waiting for VPN to connect...'
      break if system("ping -c 3 #{@vpn_gw_private_ip}")
      raise 'Timeout waiting for vpn connection' if Time.now - from > 120
      sleep 3
    end
    puts 'VPN is now connected!'
    from = Time.now
    loop do
      puts 'Waiting for route setup...'
      routes = `ip route show`
      break if routes.include?('dev tun')
      raise 'Timeout waiting for openvpn routes' if Time.now - from > 300
      sleep 3
    end
    puts 'Routes are now set up!'
  end

  def import_vnet_resolver
    Net::SSH.start(@vpn_gw_private_ip, @ovpn_login[:username], timeout: 20, verify_host_key: false) do |ssh_conn|
      res_conf = ssh_conn.exec!('cat /etc/resolv.conf')
      raise "Couldn't get resolv.conf from VPN gateway" if res_conf.empty? || !res_conf.exitstatus.zero?
      resolv_lines = res_conf.split("\n")
      last_ns = resolv_lines.rindex { |l| l =~ /^nameserver/ }
      resolv_lines.insert(last_ns + 1, 'nameserver 8.8.8.8')
      resolv_lines.insert(last_ns + 2, 'nameserver 8.8.4.4')
      File.write('/etc/resolv.conf', resolv_lines.join("\n"))
    end
  rescue Net::SSH::ConnectionTimeout => ex
    puts "Retrying to SSH into #{@vpn_gw_private_ip}..."
    sleep TIMEOUT_RETRY_DELAY
    retries ||= TIMEOUT_RETRIES
    retry if (retries -= 1).positive?
    raise ex
  end

  def wait_for_ssh
    Net::SSH.start(@vpn_gw_private_ip, @ovpn_login[:username], timeout: 20, verify_host_key: false) do |ssh_conn|
      sleep 1 until ssh_conn.exec!('uname').exitstatus.zero?
    end
  rescue Net::SSH::ConnectionTimeout => ex
    puts 'Waiting for SSH to become available on gateway...'
    sleep TIMEOUT_RETRY_DELAY
    retries ||= TIMEOUT_RETRIES
    retry if (retries -= 1).positive?
    raise ex
  end

  def stop
    if ENV.key?('TECTONIC_TESTS_DONT_CLEAN_UP')
      print 'Cleanup inhibiting flag set. Stopping here.'
      return
    end
    destroy_resources
    system("rm -rf #{@build_dir}")
  end

  def destroy_resources
    Dir.chdir(@build_dir) do
      MAX_RETRIES.times do |count|
        raise 'Private vnet: destroy failed too many times' unless MAX_RETRIES > count
        break if system(env_config, 'terraform destroy -force')
      end
    end
  end
end
