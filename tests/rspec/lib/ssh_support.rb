# frozen_string_literal: true

# Shared support code for AWS-based operations
#
module SshSupport
  def self.create_ssh_key
    key_file = '${HOME}/.ssh/id_rsa'
    `KEY_FILE=#{key_file} && ssh-keygen -f "${KEY_FILE}" -t rsa -N ''`
    %x(eval `ssh-agent -a /tmp/ssh-agent.sock`; ssh-add ${HOME}/.ssh/id_rsa)
    ENV['SSH_AUTH_SOCK'] = '/tmp/ssh-agent.sock'
  end
end
