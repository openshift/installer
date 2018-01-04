# frozen_string_literal: true

require 'fileutils'
require 'json'
require 'securerandom'
require 'net/http'
require 'uri'
require 'openssl'
require 'time'
require 'net/ssh'
require 'open3'
require 'English'

# Versions for Bare metal machine
MATCHBOX_VERSION = 'v0.6.1'
TERRAFORM_VERSION = '0.11.1'
KUBECTL_VERSION = 'v1.8.4'

KUBECTL_URL = "https://storage.googleapis.com/kubernetes-release/release/#{KUBECTL_VERSION}/bin/linux/amd64/kubectl"
TERRAFORM_URL =
  "https://releases.hashicorp.com/terraform/#{TERRAFORM_VERSION}/terraform_#{TERRAFORM_VERSION}_linux_amd64.zip"

# Bare Metal support functions
#
module MetalSupport
  def self.install_base_software
    root = root_path
    execute_command("sudo curl -L -o /usr/local/bin/kubectl #{KUBECTL_URL}")
    execute_command('sudo chmod +x /usr/local/bin/kubectl')
    execute_command("curl #{TERRAFORM_URL} | funzip > /tmp/terraform")
    execute_command('sudo mv /tmp/terraform /usr/local/bin/')
    execute_command('sudo chmod +x /usr/local/bin/terraform')
    execute_command("cd #{root}/ && rm -rf matchbox && git clone https://github.com/coreos/matchbox")
    execute_command("cd #{root}/matchbox && git checkout #{MATCHBOX_VERSION}")
    puts 'Finished initial setup'
  end

  def self.setup_bare(varfile)
    # Copy the certificates to matchbox folder
    tectonic_folder = File.expand_path('../', Dir.pwd)
    root = root_path
    certs = Dir["#{tectonic_folder}/smoke/bare-metal/fake-creds/{ca.crt,server.crt,server.key}"]
    certs.each do |cert|
      filename_dest = cert.split('/')[-1]
      dest_folder = "#{root}/matchbox/examples/etc/matchbox/#{filename_dest}"
      FileUtils.cp(cert, dest_folder)
    end

    # Download CoreOS images
    cl_version = varfile.tectonic_container_linux_version
    system(env_variables_setup, "#{root}/matchbox/scripts/get-coreos stable #{cl_version} ${ASSETS_DIR}")

    # Configuring ssh-agent
    execute_command('eval "$(ssh-agent -s)"')
    execute_command("sudo chmod 600 #{root}/matchbox/tests/smoke/fake_rsa")
    execute_command("ssh-add #{root}/matchbox/tests/smoke/fake_rsa")

    # Setting up the metal0 bridge
    execute_command('sudo mkdir -p /etc/rkt/net.d')
    execute_command("sudo cp #{root}/tests/rspec/utils/20-metal.conf /etc/rkt/net.d/")
    execute_command('cat /etc/rkt/net.d/20-metal.conf')

    # Setting up auth to download images from quay.io
    execute_command('sudo mkdir -p /etc/rkt/auth.d')
    rkt_auth_file = File.read("#{root}/tests/rspec/utils/rkt-auth.json")
    data_hash = JSON.parse(rkt_auth_file)
    data_hash['credentials']['user'] = ENV['QUAY_ROBOT_USERNAME']
    data_hash['credentials']['password'] = ENV['QUAY_ROBOT_SECRET']
    File.open('/tmp/rkt-auth.json', 'w') do |f|
      f.write(data_hash.to_json)
    end
    execute_command('sudo mv /tmp/rkt-auth.json /etc/rkt/auth.d/')

    # Setting up DNS
    `grep -q "172.18.0.3" /etc/resolv.conf`
    return if $CHILD_STATUS.exitstatus.zero?
    execute_command('sudo cp /etc/resolv.conf /etc/resolv.conf.bak')
    execute_command('echo "nameserver 172.18.0.3" | cat - /etc/resolv.conf | sudo tee /etc/resolv.conf >/dev/null')
    execute_command('cat /etc/resolv.conf')
  end

  def self.start_matchbox(varfile)
    root = root_path
    matchbox_dir = "#{root}/matchbox"
    Dir.chdir(matchbox_dir) do
      system(env_variables_setup, 'sudo -E ./scripts/devnet create')
      wait_for_matchbox
      wait_for_terraform(varfile)
    end
  end

  def self.destroy
    root = root_path
    matchbox_dir = "#{root}/matchbox"
    Dir.chdir(matchbox_dir) do
      system(env_variables_setup, 'sudo -E ./scripts/devnet destroy')
      system(env_variables_setup, 'sudo -E ./scripts/libvirt destroy')
    end
    # Restore resolv.conf
    execute_command('sudo cp /etc/resolv.conf.bak /etc/resolv.conf')
    execute_command('sudo cp /dev/null ${HOME}/.ssh/known_hosts')
    %x(for p in \`sudo rkt list | tail -n +2 | awk '{print $1}'\`; do sudo rkt stop --force $p; done)
    execute_command('sudo rkt gc --grace-period=0s')
    %x(for ns in \`ip netns l | grep -o -E '^[[:alnum:]]+'\`; do sudo ip netns del $ns; done; sudo ip l del metal0)
    %x(for veth in \`ip l show | grep -oE 'veth[^@]+'\`; do sudo ip l del $veth; done)
    execute_command('sudo rm -Rf /var/lib/cni/networks/*')
    execute_command('sudo rm -Rf /var/lib/rkt/*')
    execute_command('sudo rm -f /etc/rkt/net.d/20-metal.conf')
    execute_command('sudo rm -rf /tmp/matchbox')
    execute_command('sudo systemctl reset-failed')
  end

  def self.wait_for_matchbox
    from = Time.now
    loop do
      `curl --silent --fail -k http://matchbox.example.com:8080 > /dev/null`
      if $CHILD_STATUS.exitstatus != 0
        puts 'Waiting for matchbox...'
        elapsed = Time.now - from
        raise 'matchbox never returned with successful error code' if elapsed > 1200 # 20 mins timeout
        check_service('dev-matchbox')
        check_service('dev-dnsmasq')
        sleep 5
      else
        puts 'Matchbox up and running...'
        break
      end
    end
  end

  # We fork this to run independently from the current execution because we need to wait for
  # terraform apply to reach some state that provide the ipxe before we boot the nodes.
  def self.wait_for_terraform(varfile)
    fork do
      node_mac = varfile.tectonic_metal_controller_macs[0].tr(':', '-')
      Retriable.with_retries(limit: 24, sleep: 10) do
        succeeded = system('curl --silent --fail -k ' \
                    "\"http://matchbox.example.com:8080/ipxe?uuid=&mac=#{node_mac}&domain=&hostname=&serial=\"" \
                    ' > /dev/null')
        raise 'IPXE is not available yet' unless succeeded
      end
      puts 'Terraform is ready'
      puts 'Starting QEMU/KVM nodes'
      succeeded = system(env_variables_setup, 'sudo -E ./scripts/libvirt create')
      raise 'Failed to start QEMU/KVM nodes' unless succeeded
      puts 'Done with QEMU/KVM nodes'
      exit 0
    end
  end

  def self.check_service(service)
    return true unless `sudo systemctl is-failed #{service}`.chomp.include?('failed')
    print_service_logs(service)
    raise "#{service} failed to load"
  end

  def self.print_service_logs(service)
    cmd = "sudo journalctl --no-pager -u #{service}"
    stdout, stderr, exit_status = Open3.capture3(cmd)
    output = "Journal of #{service} service (exitcode #{exit_status})"
    output += "\nStandard output: \n#{stdout}"
    output += "\nStandard error: \n#{stderr}"
    output += "\nEnd of journal of #{service} service"

    puts output
    save_to_file(service, output)
  end

  def self.save_to_file(service_name, output)
    logs_path = "#{root_path}build/#{ENV['CLUSTER']}/logs/systemd"
    save_file = File.open("#{logs_path}/#{service_name}.log", 'w+')
    save_file << output
    save_file.close
  end

  def self.root_path
    File.expand_path('../../', Dir.pwd)
  end

  def self.env_variables_setup
    {
      'VM_DISK' => '20',
      'VM_MEMORY' => '2048',
      'ASSETS_DIR' => '/tmp/matchbox/assets'
    }
  end

  def self.execute_command(cmd)
    Open3.popen3(cmd) do |_stdin, stdout, stderr, wait_thr|
      exit_status = wait_thr.value
      unless exit_status.success?
        while (line = stdout.gets)
          puts line
        end
        while (line = stderr.gets)
          puts line
        end
        raise "Command execution FAILED! #{cmd}"
      end
    end
  end
end
