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

def ssh_exec(ip_address, command)
  status = {}
  stdout = ''
  stderr = ''
  Net::SSH.start(ip_address, 'core', forward_agent: true, use_agent: true) do |ssh|
    ssh.exec! command, status: status do |_ch, stream, data|
      if stream == :stdout
        stdout = data
      else
        stderr = data
      end
    end
  end
  [stdout, stderr, status[:exit_code]]
end
