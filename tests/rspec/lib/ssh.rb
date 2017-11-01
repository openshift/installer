# frozen_string_literal: true

require 'net/ssh'
require 'with_retries'

def check_prerequisites
  return if ssh_agent_has_key?
  raise 'No ssh key registered in ssh-agent. Run `ssh-add' \
        '<path-to-private-key>`to add a key.'
end

def ssh_agent_has_key?
  system('ssh-add -l')
end

def ssh_exec(ip_address, command, max_retries = 5)
  status = {}
  stdout = String.new('')
  stderr = String.new('')
  Retriable.with_retries(Errno::ECONNREFUSED, Errno::ECONNRESET, Errno::ETIMEDOUT,
                         Net::SSH::ConnectionTimeout, Net::SSH::Disconnect,
                         IOError, limit: max_retries, sleep: 10) do
    Net::SSH.start(
      ip_address, 'core', forward_agent: true, use_agent: true, verify_host_key: Net::SSH::Verifiers::Null.new
    ) do |ssh|
      ssh.exec! command, status: status do |_ch, stream, data|
        case stream
        when :stdout
          stdout << data
        when :stderr
          stderr << data
        end
      end
    end
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
