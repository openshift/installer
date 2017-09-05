# frozen_string_literal: true

def check_prerequisites
  return if ssh_agent_has_key?
  raise 'No ssh key registered in ssh-agent. Run `ssh-add' \
        '<path-to-private-key>`to add a key.'
end

def ssh_agent_has_key?
  system('ssh-add -l')
end
