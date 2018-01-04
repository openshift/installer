#!/usr/bin/env bash
# shellcheck disable=SC2002,SC2015,SC2086,SC2091
# This scripts brings up a bare-metal Tectonic cluster using VMs and
# containerized matchbox/dnsmasq servers.
#
# The following environment variables are expected:
# - BRANCH_NAME
# - BUILD_ID
# The script setups the environment and calls `make apply` from the repository
# root.
#
# Due to the assumptions made by this script, it is *not* safe to run multiple
# instances of it on a single host and the Terraform configuration must be
# matching the infrastructure. Notably:
# - matchbox is expected on 172.18.0.2,
# - three nodes are expected on 172.18.0.21 (master), 172.18.0.22 (worker), 172.18.0.23 (worker).
#
# This script requires the following packages on the host:
# - qemu-kvm
# - libvirt-bin
# - virt-manager
# - curl
# - unzip
# - policycoreutils.
set -xe

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."
BIN_DIR="$ROOT/bin_test"

MATCHBOX_VERSION=v0.6.1
KUBECTL_VERSION=v1.6.4
TERRAFORM_VERSION=0.11.1

KUBECTL_URL="https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl"
TERRAFORM_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"

export VM_DISK='20'
export VM_MEMORY='2048'
export ASSETS_DIR="${ASSETS_DIR:-/tmp/matchbox/assets}"

main() {
  if [ -z "${BRANCH_NAME}" ] || [ -z "${BUILD_ID}" ]; then
    echo "\$BRANCH_NAME, \$BUILD_ID are required"
    return 1
  fi
  if [ -z "$1" ]; then
    echo "$0 <tfvars file's relative path>"
    return 1
  fi

  echo "Installing required binaries"
  install

  echo "Cleanup testing environment"
  cleanup &>/dev/null
  trap kill_terraform_and_cleanup EXIT

  echo "Setting up configuration and environment"
  configure "$1"
  setup

  echo "Starting matchbox"
  (cd "${ROOT}"/matchbox && sudo -S -E ./scripts/devnet create)
  echo "Waiting for matchbox..."
  until $(curl --silent --fail -k http://matchbox.example.com:8080 > /dev/null); do
    echo "Waiting for matchbox..."
    sleep 5

    if sudo -E systemctl is-failed dev-matchbox; then
      sudo -E journalctl -u dev-matchbox
      exit 1
    fi

    if sudo -E systemctl is-failed dev-dnsmasq; then
      sudo -E journalctl -u dev-dnsmasq
      exit 1
    fi
  done

  echo "Starting Terraform"
  (cd "${ROOT}" && make apply || kill $$) & # Self-destruct and trigger trap on failure
  TERRAFORM_PID=$!
  sleep 15

  echo "Starting QEMU/KVM nodes"
  (cd "${ROOT}"/matchbox && sudo -E ./scripts/libvirt create)

  echo "Waiting for Kubernetes/Tectonic cluster to be up and running:"
  cluster_up

  echo "Running Go smoke tests"
  test_cluster

  echo "SUCCESS: Tectonic bare-metal cluster came up!"
  cleanup
}

install() {
  mkdir -p $BIN_DIR
  export PATH=$BIN_DIR:$PATH
    
  echo "Installing kubectl"
  curl -L -o ${BIN_DIR}/kubectl ${KUBECTL_URL}
  chmod +x ${BIN_DIR}/kubectl

  echo "Installing Terraform"
  curl ${TERRAFORM_URL} | funzip > $BIN_DIR/terraform
  sudo chmod +x $BIN_DIR/terraform

  echo "Installing matchbox"
  (cd ${ROOT}/ && rm -rf matchbox && git clone https://github.com/coreos/matchbox)
  (cd ${ROOT}/matchbox && git checkout $MATCHBOX_VERSION)
}

setup() {
  echo "Copying matchbook test credentials"
  cp ${DIR}/fake-creds/{ca.crt,server.crt,server.key} ${ROOT}/matchbox/examples/etc/matchbox

  if [ ! -d $ASSETS_DIR/coreos/$COREOS_VERSION ]; then
    echo "Downloading CoreOS image"
    ${ROOT}/matchbox/scripts/get-coreos $COREOS_CHANNEL $COREOS_VERSION $ASSETS_DIR
  fi

  echo "Configuring ssh-agent"
  eval "$(ssh-agent -s)"
  chmod 600 ${ROOT}/matchbox/tests/smoke/fake_rsa
  ssh-add ${ROOT}/matchbox/tests/smoke/fake_rsa

  echo "Setting up the metal0 bridge"
  sudo mkdir -p /etc/rkt/net.d
  sudo bash -c 'cat > /etc/rkt/net.d/20-metal.conf << EOF
  {
    "name": "metal0",
    "type": "bridge",
    "bridge": "metal0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
      "type": "host-local",
      "subnet": "172.18.0.0/24",
      "routes" : [ { "dst" : "0.0.0.0/0" } ]
     }
  }
EOF'

  echo "Setting up DNS"
  if ! grep -q "172.18.0.3" /etc/resolv.conf; then
    echo "nameserver 172.18.0.3" | cat - /etc/resolv.conf | sudo tee /etc/resolv.conf >/dev/null
  fi
}

configure() {
  export PLATFORM=metal
  export CLUSTER="tf-${PLATFORM}-${BRANCH_NAME}-${BUILD_ID}"
  export TF_VAR_tectonic_cluster_name
  TF_VAR_tectonic_cluster_name=$(echo ${CLUSTER} | awk '{print tolower($0)}')

  CONFIG=${DIR}/$1
  make localconfig
  ln -sf ${CONFIG} ${ROOT}/build/${CLUSTER}/terraform.tfvars

  COREOS_CHANNEL=$(awk -F "=" '/^tectonic_container_linux_channel/ {gsub(/[ \t"]/, "", $2); print $2}' ${CONFIG})
  COREOS_VERSION=$(awk -F "=" '/^tectonic_container_linux_version/ {gsub(/[ \t"]/, "", $2); print $2}' ${CONFIG})

  export SMOKE_KUBECONFIG=${ROOT}/build/${CLUSTER}/generated/auth/kubeconfig
}

cleanup() {
  set +e

  # Kill any remaining VMs.
  (cd ${ROOT}/matchbox && sudo ./scripts/libvirt destroy)

  # shellcheck disable=SC2006
  # Reset rkt pods and CNI entirely, to avoid IP conflict due to leakage bug.
  for p in `sudo rkt list | tail -n +2 | awk '{print $1}'`; do sudo rkt stop --force $p; done
  sudo rkt gc --grace-period=0s

  # shellcheck disable=SC2006
  for ns in `ip netns l | grep -o -E '^[[:alnum:]]+'`; do sudo ip netns del $ns; done
  sudo ip l del metal0
  # shellcheck disable=SC2006
  for veth in `ip l show | grep -oE 'veth[^@]+'`; do sudo ip l del $veth; done

  sudo rm -Rf /var/lib/cni/networks/*
  sudo rm -Rf /var/lib/rkt/*
  sudo rm -f /etc/rkt/net.d/20-metal.conf

  # Reset DNS.
  cat /etc/resolv.conf | grep -v 172.18.0.3 | sudo tee /etc/resolv.conf

  # Reset failed units (i.e. matchbox, dnsmasq which we just killed).
  sudo systemctl reset-failed

  set -e
}

kill_terraform_and_cleanup() {
  echo "Killing Terraform"
  kill ${TERRAFORM_PID} || true

  echo "WARNING: Cleanup is temporarily disabled on failure for debugging purposes. Next job will clean at startup."
  #echo "Cleanup testing environment"
  #cleanup
}

kubelet_up() {
  ssh -q -i ${ROOT}/matchbox/tests/smoke/fake_rsa \
   -o StrictHostKeyChecking=no \
   -o UserKnownHostsFile=/dev/null \
   -o PreferredAuthentications=publickey \
   core@$1 /usr/bin/systemctl status k8s-node-bootstrap kubelet
  curl --silent --fail -m 1 "http://$1:10255/healthz" > /dev/null
}

cluster_up() {
  echo "Waiting for Kubelets to start..."
  until kubelet_up "node1.example.com" \
    && kubelet_up "node2.example.com" \
    && kubelet_up "node3.example.com"
  do
    sleep 15
    echo "Waiting for Kubelets to start..."
  done

  echo "$(readyNodes) of 3 nodes are Ready..."
  until [[ "$(readyNodes)" == "3" ]]; do
    sleep 5
    echo "$(readyNodes) of 3 nodes are Ready..."
  done

  echo "List of nodes:"
  k8s get nodes

  sleep 5
  until [[ "$(readyPods)" == "$(podCount)" && "$(readyPods)" -gt "0" ]]; do
    sleep 15
    echo "$(readyPods) / $(podCount) pods are Running..."
    k8s get pods --all-namespaces || true
  done

  echo "List of pods:"
  k8s get pods --all-namespaces || true

  echo "Waiting for Tectonic Console..."
  until $(curl --silent --fail -k https://tectonic.example.com > /dev/null); do
    echo "Waiting for Tectonic Console..."
    k8s get pods --all-namespaces || true
    sleep 15
  done
}

k8s() {
  ${BIN_DIR}/kubectl --kubeconfig=${SMOKE_KUBECONFIG} "$@"
}

# ready nodes returns the number of Ready Kubernetes nodes
readyNodes() {
  # shellcheck disable=SC2126
  k8s get nodes -o template --template='{{range .items}}{{range .status.conditions}}{{if eq .type "Ready"}}{{.}}{{end}}{{end}}{{end}}' | grep -o -E True | wc -l
}

# ready pods returns the number of Running pods
readyPods() {
  # shellcheck disable=SC2126
  k8s get pods --all-namespaces -o template --template='{{range .items}}{{range .status.conditions}}{{if eq .type "Ready"}}{{.}}{{end}}{{end}}{{end}}' | grep -o -E True | wc -l
}

# podCount returns the number of pods
podCount() {
  # shellcheck disable=SC2126
  k8s get pods --all-namespaces -o template --template='{{range .items}}{{range .status.conditions}}{{if eq .type "Ready"}}{{.}}{{end}}{{end}}{{end}}' | grep -o -E status | wc -l
}

test_cluster() {
  MASTER_COUNT=$(grep tectonic_master_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
  WORKER_COUNT=$(grep tectonic_worker_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
  export SMOKE_NODE_COUNT=$(( MASTER_COUNT + WORKER_COUNT ))
  export SMOKE_MANIFEST_PATHS=${ROOT}/build/${CLUSTER}/generated/
  # shellcheck disable=SC2155
  export SMOKE_NETWORKING=$(grep tectonic_networking "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}' | tr -d ' ')
  bin/smoke -test.v -test.parallel=1 --cluster
}

main "$@"
