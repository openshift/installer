# frozen_string_literal: true

require 'json'

# AWSVPC represents an AWS virtual private cloud
class AWSVPC
  attr_reader :vpn_url
  attr_reader :ovpn_password
  attr_reader :name
  attr_reader :vpc_dns
  attr_reader :vpc_id
  attr_reader :private_zone_id
  attr_reader :master_subnet_ids
  attr_reader :worker_subnet_ids
  attr_reader :vpn_connection

  def initialize(name)
    @name = name
    @ovpn_password =
      `tr -cd '[:alnum:]' < /dev/urandom | head -c 32 ; echo`.chomp
  end

  def env_variables
    {
      'TF_VAR_vpc_aws_region' => ENV['TF_VAR_tectonic_aws_region'],
      'TF_VAR_vpc_name' => @name,
      'TF_VAR_ovpn_password' => @ovpn_password,
      'TF_VAR_base_domain' => 'tectonic-ci.de'
    }
  end

  def export_tfvars
    vars = {
      'TF_VAR_tectonic_aws_external_private_zone' => @private_zone_id,
      'TF_VAR_tectonic_aws_external_vpc_id' => @vpc_id,
      'TF_VAR_tectonic_aws_external_master_subnet_ids' => @master_subnet_ids,
      'TF_VAR_tectonic_aws_external_worker_subnet_ids' => @worker_subnet_ids
    }
    vars.each do |key, value|
      ENV[key] = value
    end
  end

  def create
    Dir.chdir('../../contrib/internal-cluster') do
      succeeded = system(env_variables, 'terraform init')
      raise 'could not init Terraform to create VPC' unless succeeded
      succeeded = system(env_variables, 'terraform apply -auto-approve')
      raise 'could not create vpc with Terraform' unless succeeded

      parse_terraform_output
      wait_for_vpn_access_server

      @vpn_connection = VPNConnection.new(@ovpn_password, @vpn_url)
      @vpn_connection.start
    end

    set_nameserver
    export_tfvars
  end

  def parse_terraform_output
    tf_out = JSON.parse(`terraform output -json`)
    @vpn_url = tf_out['ovpn_url']['value']
    @vpc_dns = tf_out['vpc_dns']['value']
    @vpc_id = tf_out['vpc_id']['value']
    @private_zone_id = tf_out['private_zone_id']['value']
    parse_subnets(tf_out)
  end

  def parse_subnets(tf_out)
    subnets = tf_out['subnets']['value']
    @master_subnet_ids =
      "[\"#{subnets[0]}\", \"#{subnets[1]}\"]"
    @worker_subnet_ids =
      "[\"#{subnets[2]}\", \"#{subnets[3]}\"]"
  end

  def destroy
    @vpn_connection.stop
  rescue
    raise 'could not disconnect from vpn'
  ensure
    terraform_destroy
    recover_etc_resolv
  end

  def terraform_destroy
    Dir.chdir('../../contrib/internal-cluster') do
      3.times do
        return if system(env_variables, 'terraform destroy -force')
      end
    end

    raise 'could not destroy vpc with Terraform'
  end

  def wait_for_vpn_access_server
    90.times do
      succeeded = system("curl -k -L --silent #{@vpn_url} > /dev/null")
      return if succeeded
      sleep(0.5)
    end
    raise 'waiting for vpn access server timed out'
  end

  def set_nameserver
    # Use AWS VPC DNS rather than host's.
    FileUtils.cp '/etc/resolv.conf', '/etc/resolv.conf.bak'
    IO.write('/etc/resolv.conf', "nameserver #{@vpc_dns}\nnameserver 8.8.8.8\n")
    system('cat /etc/resolv.conf')
  end

  def recover_etc_resolv
    FileUtils.cp '/etc/resolv.conf.bak', '/etc/resolv.conf'
  end
end

# VPNConnection represents a VPN connection via the VPN server in an AWS VPC
class VPNConnection
  attr_reader :vpn_url
  attr_reader :ovpn_password
  attr_reader :vpn_conf

  def initialize(ovpn_password, vpn_url)
    @ovpn_password = ovpn_password
    @vpn_url = vpn_url
  end

  def curl_vpn_config
    cmd = 'curl -k -L ' \
          "-u 'openvpn:#{@ovpn_password}' " \
          '--silent ' \
          '--fail ' \
          "#{@vpn_url}/rest/GetUserlogin"
    @vpn_conf = `#{cmd}`.chomp
  end

  def start
    curl_vpn_config

    IO.write('vpn_credentials', "openvpn\n#{@ovpn_password}\n")
    @vpn_conf = @vpn_conf.sub(
      'auth-user-pass',
      'auth-user-pass vpn_credentials'
    )
    IO.write('vpn.conf', @vpn_conf)

    succeeded = system('openvpn --config vpn.conf --daemon')
    raise 'could not start vpn' unless succeeded

    wait_for_network
  end

  def stop
    system('pkill openvpn || true')
    wait_for_network
  end

  def wait_for_network
    90.times do
      succeeded = system('ping -c 1 8.8.8.8 > /dev/null')
      return if succeeded
    end
    raise 'waiting for network timed out'
  end
end
