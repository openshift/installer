#!/usr/bin/env bash
set -x

chmod 700 /root/.ssh
chmod 600 /root/.ssh/*

# shellcheck disable=SC1091
source /root/.bash_profile

# DNSMASQ setup
cat <<EOF > /etc/dnsmasq.conf
bind-interfaces
interface=lo
strict-order
user=root
domain-needed
bogus-priv
filterwin2k
localise-queries
no-negcache
no-resolv
$(grep -oE 'nameserver.*' /etc/resolv.conf | sed -E 's/^nameserver (.*)/server=\1/')
# server=$(ip route get 1.1.1.1 | grep -oE 'via ([^ ]+)' | sed -E 's/via //')
 server=/tt.testing/192.168.126.1
EOF

cp /etc/resolv.conf{,.bkp}
cat <<EOF > /etc/resolv.conf
nameserver 127.0.0.1
EOF

dnsmasq

# Start LIBVIRT
libvirtd -d --listen -f /etc/libvirt/libvirtd.conf
virtlockd -d
virtlogd -d

mkdir -p "/opt/app-root/src/github.com/openshift"
cd "/opt/app-root/src/github.com/openshift" || exit 1
git clone "https://github.com/${REPO_OWNER}/installer.git" || exit 1

cd "/opt/app-root/src/github.com/openshift/installer" || exit 1
git checkout "$BRANCH" || exit 1
./hack/build.sh
bash -i
