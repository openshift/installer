#!/bin/bash
set -e
set -u
set -o pipefail
set -x

sudo dnf install -y libvirt libvirt-devel libvirt-client git libvirt-daemon-kvm bind-utils jq gcc-c++
# https://github.com/ironcladlou/openshift4-libvirt-gcp/issues/29
sudo dnf install -y qemu-kvm-2.12.0-88.module+el8.1.0+5708+85d8e057.3

# Install golang
curl -L https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz -o go1.13.8.linux-amd64.tar.gz
tar -xvf go1.13.8.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin

# Install yq to manipulate manifest file created by installer.
if [[ ! -e /usr/local/bin/yq ]]; then
    curl -L https://github.com/mikefarah/yq/releases/download/2.2.1/yq_linux_amd64 -o yq
    chmod +x yq
    sudo mv yq /usr/local/bin/yq
fi

# Enable IP forwarding
# https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md#enable-ip-forwarding
sudo sysctl net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward = 1" | sudo tee /etc/sysctl.d/99-ipforward.conf
sudo sysctl -p /etc/sysctl.d/99-ipforward.conf

# Enable non-root access to libvirt stuff
# https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md#make-sure-you-have-permissions-for-qemusystem
sudo bash -c 'cat > /etc/polkit-1/rules.d/80-libvirt.rules' << EOF
polkit.addRule(function(action, subject) {
  if (action.id == "org.libvirt.unix.manage" subject.isInGroup("google-sudoers")) {
      return polkit.Result.YES;
  }
});
EOF

# Configure libvirt to accept TCP connections
# https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md#configure-libvirt-to-accept-tcp-connections
sudo bash -c 'cat >> /etc/libvirt/libvirtd.conf' << EOF
listen_tls = 0
listen_tcp = 1
auth_tcp="none"
tcp_port = "16509"
EOF
sudo bash -c 'cat >> /etc/sysconfig/libvirtd' << EOF
LIBVIRTD_ARGS="--listen"
EOF
sudo bash -c 'cat >> /etc/modprobe.d/kvm.conf' << EOF
options kvm_intel nested=1
EOF
# Ensure nesting is enabled in the kernel
# TODO: verify this is still necessary
sudo modprobe -r kvm_intel
sudo modprobe kvm_intel nested=1
sudo systemctl restart libvirtd
# Set up iptables and firewalld
# TODO: discover the ports
sudo firewall-cmd --add-rich-rule='rule family=ipv4 source address=192.168.126.0/24 destination address=192.168.122.1 port port=16509 protocol=tcp accept' --permanent --zone=libvirt
sudo firewall-cmd --zone=libvirt --add-service=libvirt --permanent

# Enable NetworkManager DNS overlay
# https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md#set-up-networkmanager-dns-overlay
echo -e "[main]\\ndns=dnsmasq" | sudo tee /etc/NetworkManager/conf.d/openshift.conf
echo server=/openshift.testing/192.168.126.1 | sudo tee /etc/NetworkManager/dnsmasq.d/openshift.conf
# Create new domain for ingress to make sure it able to resolve auth route URL
echo address=/.apps.openshift.testing/192.168.126.51 | sudo tee -a /etc/NetworkManager/dnsmasq.d/openshift.conf
sudo systemctl restart NetworkManager

# Configure the default libvirt storage pool
# https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md#configure-default-libvirt-storage-pool
sudo virsh pool-define /dev/stdin <<EOF
<pool type='dir'>
  <name>default</name>
  <target>
    <path>/var/lib/libvirt/images</path>
  </target>
</pool>
EOF
sudo virsh pool-start default
sudo virsh pool-autostart default

echo "Installing oc client"
curl -OL https://mirror.openshift.com/pub/openshift-v4/clients/oc/latest/linux/oc.tar.gz
tar -zxf oc.tar.gz
rm -fr oc.tar.gz
sudo mv oc /usr/local/bin
sudo ln -s /usr/local/bin/oc /usr/local/bin/kubectl

sudo bash -c 'cat >> /etc/bashrc' << EOF
export PATH=$PATH:/usr/local/go/bin
EOF
