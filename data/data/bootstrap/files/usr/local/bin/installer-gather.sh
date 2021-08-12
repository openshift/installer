#!/usr/bin/env bash

# shellcheck disable=SC1091
. /usr/local/bin/bootstrap-cluster-gather.sh

if test "x${1}" = 'x--id'
then
	GATHER_ID="${2}"
	shift 2
fi

ARTIFACTS="/tmp/artifacts-${GATHER_ID}"
mkdir -p "${ARTIFACTS}"

exec &> >(tee "${ARTIFACTS}/gather.log")

echo "Gathering bootstrap service records ..."
mkdir -p "${ARTIFACTS}/bootstrap/services"
sudo cp -r /var/log/openshift/* "${ARTIFACTS}/bootstrap/services/"

echo "Gathering bootstrap systemd summary ..."
LANG=POSIX systemctl list-units --state=failed >& "${ARTIFACTS}/failed-units.txt"

echo "Gathering bootstrap failed systemd unit status ..."
mkdir -p "${ARTIFACTS}/unit-status"
sed -n 's/^\* \([^ ]*\) .*/\1/p' < "${ARTIFACTS}/failed-units.txt" | while read -r UNIT
do
    systemctl status --full "${UNIT}" >& "${ARTIFACTS}/unit-status/${UNIT}.txt"
    journalctl -u "${UNIT}" > "${ARTIFACTS}/unit-status/${UNIT}.log"
done

echo "Gathering bootstrap journals ..."
mkdir -p "${ARTIFACTS}/bootstrap/journals"
for service in release-image release-image-download crio-configure bootkube kubelet crio approve-csr ironic master-bmh-update
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/bootstrap/journals/${service}.log"
done

echo "Gathering bootstrap networking ..."
mkdir -p "${ARTIFACTS}/bootstrap/network"
ip addr >& "${ARTIFACTS}/bootstrap/network/ip-addr.txt"
ip route >& "${ARTIFACTS}/bootstrap/network/ip-route.txt"
hostname >& "${ARTIFACTS}/bootstrap/network/hostname.txt"
cp -r /etc/resolv.conf "${ARTIFACTS}/bootstrap/network/"

echo "Gathering bootstrap containers ..."
mkdir -p "${ARTIFACTS}/bootstrap/containers"
sudo crictl ps --all --quiet | while read -r container
do
    container_name="$(sudo crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")"
    sudo crictl logs "${container}" >& "${ARTIFACTS}/bootstrap/containers/${container_name}-${container}.log"
    sudo crictl inspect "${container}" >& "${ARTIFACTS}/bootstrap/containers/${container_name}-${container}.inspect"
done
sudo cp -r /var/log/bootstrap-control-plane/ "${ARTIFACTS}/bootstrap/containers"
mkdir -p "${ARTIFACTS}/bootstrap/pods"
sudo podman ps --all --format "{{ .ID }} {{ .Names }}" | while read -r container_id container_name
do
    sudo podman logs "${container_id}" >& "${ARTIFACTS}/bootstrap/pods/${container_name}-${container_id}.log"
    sudo podman inspect "${container_id}" >& "${ARTIFACTS}/bootstrap/pods/${container_name}-${container_id}.inspect"
done

echo "Gathering rendered assets..."
mkdir -p "${ARTIFACTS}/rendered-assets"
sudo cp -r /var/opt/openshift/ "${ARTIFACTS}/rendered-assets"
sudo chown -R "${USER}":"${USER}" "${ARTIFACTS}/rendered-assets"
sudo find "${ARTIFACTS}/rendered-assets" -type d -print0 | xargs -0 sudo chmod u+x
# remove sensitive information
# TODO leave tls.crt inside of secret yaml files
find "${ARTIFACTS}/rendered-assets" -name "*secret*" -print0 | xargs -0 rm -rf
find "${ARTIFACTS}/rendered-assets" -name "*kubeconfig*" -print0 | xargs -0 rm
find "${ARTIFACTS}/rendered-assets" -name "*.key" -print0 | xargs -0 rm
find "${ARTIFACTS}/rendered-assets" -name ".kube" -print0 | xargs -0 rm -rf

# Collect cluster data
GATHER_KUBECONFIG="/opt/openshift/auth/kubeconfig"

ALTERNATIVE_KUBECONFIG="/etc/kubernetes/bootstrap-secrets/kubeconfig"
if [[ -f ${ALTERNATIVE_KUBECONFIG} ]]; then
    GATHER_KUBECONFIG=${ALTERNATIVE_KUBECONFIG}
fi

echo "Using ${GATHER_KUBECONFIG} as KUBECONFIG"

cluster_bootstrap_gather
if [ -d "${API_SERVER_ARTIFACTS_DIR}/resources" ]
then
    cp -r "${API_SERVER_ARTIFACTS_DIR}/resources" "${ARTIFACTS}"
fi
# The existence of the file located in LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME is used
# as indication that a bootstrap process has already previously taken place and the resulting
# bundle already exists in the filesystem. In that case, we include said bundle inside the log
# bundle created by this script.
#
# An example for a scenario where this happens is when we're running inside a single-node
# bootstrap-in-place deployment post-pivot master node, rather than a typical bootstrap node.
# In that case, the bootstrap node collects and bundles logs pre-reboot and transfers that bundle
# via an ignition file to this post-pivot master node.
LOG_BUNDLE_BOOTSTRAP_NAME="log-bundle-bootstrap"
LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME="/var/log/${LOG_BUNDLE_BOOTSTRAP_NAME}.tar.gz"

if [[ -f ${LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME} ]]; then
    echo "Including existing bootstrap bundle ${LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME}"
    tar -xzf ${LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME} --directory "${ARTIFACTS}"
fi
mkdir -p "${ARTIFACTS}/control-plane"
echo "Gather remote logs"
export MASTERS=()
MASTER_GATHER_ID="master-${GATHER_ID}"
if [[ -f ${LOG_BUNDLE_BOOTSTRAP_ARCHIVE_NAME} ]]; then
    # Instead of running installer-masters-gather.sh on remote masters, run it on the current node
    MASTER_ARTIFACTS="/tmp/artifacts-${MASTER_GATHER_ID}"
    mkdir -p "${ARTIFACTS}/control-plane/master"
    sudo /usr/local/bin/installer-masters-gather.sh --id "${MASTER_GATHER_ID}" </dev/null
    cp -r "${MASTER_ARTIFACTS}"/* "${ARTIFACTS}/control-plane/master/"
elif [ "$#" -ne 0 ]; then
    MASTERS=( "$@" )
elif test -s "${ARTIFACTS}/resources/masters.list"; then
    mapfile -t MASTERS < "${ARTIFACTS}/resources/masters.list"
else
    echo "No masters found!"
fi

for master in "${MASTERS[@]}"
do
    echo "Collecting info from ${master}"
    scp -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null -q /usr/local/bin/installer-masters-gather.sh "core@[${master}]:"
    mkdir -p "${ARTIFACTS}/control-plane/${master}"
    ssh -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null "core@${master}" -C "sudo ./installer-masters-gather.sh --id '${MASTER_GATHER_ID}'" </dev/null
    scp -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null -r -q "core@[${master}]:/tmp/artifacts-${MASTER_GATHER_ID}/*" "${ARTIFACTS}/control-plane/${master}/"
done

TAR_FILE="${TAR_FILE:-${HOME}/log-bundle-${GATHER_ID}.tar.gz}"
tar cz -C "${ARTIFACTS}" --transform "s?^\\.?log-bundle-${GATHER_ID}?" . > "${TAR_FILE}"
echo "Log bundle written to ${TAR_FILE}"
