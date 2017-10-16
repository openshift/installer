# frozen_string_literal: true

require 'net/ssh'

def check_prerequisites
  return if ssh_agent_has_key?
  raise 'No ssh key registered in ssh-agent. Run `ssh-add' \
        '<path-to-private-key>`to add a key.'
end

def ssh_agent_has_key?
  system('ssh-add -l')
end

def ssh_exec(ip_address, command, max_retries = 5)
  retries = 0
  status = {}
  stdout = ''
  stderr = ''
  begin
    Net::SSH.start(ip_address, 'core', forward_agent: true, use_agent: true) do |ssh|
      ssh.exec! command, status: status do |_ch, stream, data|
        if stream == :stdout
          stdout = data
        else
          stderr = data
        end
      end
    end
  rescue Errno::ECONNREFUSED, Errno::ECONNRESET, Errno::ETIMEDOUT, IOError, Net::SSH::ConnectionTimeout,
         Net::SSH::Disconnect
    raise "failed to exec '#{command}' in #{max_retries} retries" if retries >= max_retries
    retries += 1
    sleep_time = 5 * retries
    puts "failed to exec '#{command}'; retrying in #{sleep_time} seconds"
    sleep sleep_time
    retry
  end
  [stdout, stderr, status[:exit_code]]
end

def create_if_not_exist_and_add_ssh_key
  key_file = '${HOME}/.ssh/id_rsa'
  ssh_file = File.join(Dir.home, '.ssh/id_rsa')
  `KEY_FILE=#{key_file} && ssh-keygen -f "${KEY_FILE}" -t rsa -N ''` unless File.exist?(ssh_file)
  %x(eval `ssh-agent -a /tmp/ssh-agent.sock`; ssh-add ${HOME}/.ssh/id_rsa)
  ENV['SSH_AUTH_SOCK'] = '/tmp/ssh-agent.sock'
end
