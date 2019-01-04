#!/usr/bin/env bash
set -xe

# LIBVIRT
cat <<EOF >> /etc/polkit-1/rules.d/80-libvirt.rules
polkit.addRule(function(action, subject) {
  if (action.id == "org.libvirt.unix.manage" && subject.local && subject.active && subject.isInGroup("wheel")) {
      return polkit.Result.YES;
  }
});
EOF

sed -i 's/#user = "root"/user = "root"/; s/#group = "root"/group = "root"/' /etc/libvirt/qemu.conf

cat <<EOF >>/etc/libvirt/libvirtd.conf
listen_tls = 0
listen_tcp = 1
auth_tcp="none"
tcp_port = "16509"
log_level = 4
EOF

cat <<EOF >>/etc/sysconfig/libvirtd
LIBVIRTD_ARGS="--listen"
EOF

libvirtd -d
virsh --connect qemu:///system pool-create --file=/opt/app-root/src/libvirt_config/libvirt-storage-pool.xml

# TERRAFORM
cat <<EOF > "${HOME}/.terraformrc"
plugin_cache_dir = "${HOME}/.terraform.d/plugin-cache"
EOF
