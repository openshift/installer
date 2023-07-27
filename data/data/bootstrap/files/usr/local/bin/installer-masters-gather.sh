#!/usr/bin/env bash

# Get target architecture
arch=$(uname -m)

if test "x${1}" = 'x--id'
then
	GATHER_ID="${2}"
	shift 2
fi

ARTIFACTS="/tmp/artifacts-${GATHER_ID}"
mkdir -p "${ARTIFACTS}"

echo "Gathering master systemd summary ..."
LANG=POSIX systemctl list-units --state=failed >& "${ARTIFACTS}/failed-units.txt"

echo "Gathering master failed systemd unit status ..."
mkdir -p "${ARTIFACTS}/unit-status"
sed -n 's/^\* \([^ ]*\) .*/\1/p' < "${ARTIFACTS}/failed-units.txt" | while read -r UNIT
do
    systemctl status --full "${UNIT}" >& "${ARTIFACTS}/unit-status/${UNIT}.txt"
    journalctl -u "${UNIT}" > "${ARTIFACTS}/unit-status/${UNIT}.log"
done

echo "Gathering master journals ..."
mkdir -p "${ARTIFACTS}/journals"
for service in crio kubelet machine-config-daemon-host machine-config-daemon-firstboot openshift-azure-routes openshift-gcp-routes pivot sssd
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/journals/${service}.log"
done

journalctl -o with-unit --no-pager | gzip > "${ARTIFACTS}/journals/journal.log.gz"

echo "Gathering master networking ..."
mkdir -p "${ARTIFACTS}/network"
ip addr >& "${ARTIFACTS}/network/ip-addr.txt"
ip route >& "${ARTIFACTS}/network/ip-route.txt"
hostname >& "${ARTIFACTS}/network/hostname.txt"
netstat -anp >& "${ARTIFACTS}/network/netstat.txt"
cp -r /etc/resolv.conf "${ARTIFACTS}/network/"

echo "Gathering master containers ..."
mkdir -p "${ARTIFACTS}/containers"
for container in $(crictl ps --all --quiet)
do
    container_name=$(crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")
    crictl logs "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.log"
    crictl inspect "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.inspect"
done

podman ps --all --format "{{ .ID }} {{ .Names }}" | while read -r container_id container_name
do
    podman logs "${container_id}" >& "${ARTIFACTS}/containers/${container_name}-${container_id}.log"
    podman inspect "${container_id}" >& "${ARTIFACTS}/containers/${container_name}-${container_id}.inspect"
done

echo "Gathering master rpm-ostree info ..."
mkdir -p "${ARTIFACTS}/rpm-ostree"
sudo rpm-ostree status >& "${ARTIFACTS}/rpm-ostree/status"
sudo rpm-ostree ex history >& "${ARTIFACTS}/rpm-ostree/history"

# Collect system information specific to IBM Linux Z (s390x) systems. The dbginfo
# script is available by default as part of the s390-utils rpm package
if [ "$arch" == "s390x" ]
then
    echo "Gathering dbginfo for the s390x system"
    mkdir -p "${ARTIFACTS}/node-dbginfo"
    /usr/sbin/dbginfo.sh -d "${ARTIFACTS}/node-dbginfo"
    find "${ARTIFACTS}/node-dbginfo" -print0 | xargs -0 chmod a+r
fi

echo "Waiting for logs ..."
while wait -n; do jobs; done
