# frozen_string_literal: true

require 'ssh'

SSH_CMD_CONTAINER_LINUX_VERSION = 'if sudo [ -f /var/lib/update_engine/prefs/aleph-version ]; then \
  sudo cat /var/lib/update_engine/prefs/aleph-version; \
else \
  source /usr/share/coreos/release && echo "$COREOS_RELEASE_VERSION"; \
fi'
SSH_CMD_CONTAINER_LINUX_CHANNEL = 'for conf in /usr/share/coreos/update.conf /etc/coreos/update.conf; do \
  [ -f "$conf" ] && source "$conf"; \
done; \
echo "$GROUP"'

# ContainerLinux provides helpers to find OS-level properties for a cluster.
module ContainerLinux
  def self.version(cluster)
    v, err, = ssh_exec(cluster.master_ip_address, SSH_CMD_CONTAINER_LINUX_VERSION)
    raise "failed to get Container Linux version for #{cluster.master_ip_address}" if err != ''
    v.chomp
  end

  def self.channel(cluster)
    c, err, = ssh_exec(cluster.master_ip_address, SSH_CMD_CONTAINER_LINUX_CHANNEL)
    raise "failed to get Container Linux channel for #{cluster.master_ip_address}" if err != ''
    c.chomp
  end
end
