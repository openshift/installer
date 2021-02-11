#!/usr/bin/env bash
set -euoE pipefail ## -E option will cause functions to inherit trap

# This script is executed by bootkube.sh when installing single node with bootstrap in place
CLUSTER_BOOTSTRAP_IMAGE=$1


bootkube_podman_run() {
  # we run all commands in the host-network to prevent IP conflicts with
  # end-user infrastructure.
  podman run --quiet --net=host "${@}"
}

if [ ! -f stop-etcd.done ]; then
  echo "Stop etcd static pod by moving the manifest"
  mv /etc/kubernetes/manifests/etcd-member-pod.yaml /etc/kubernetes || echo "already moved etcd-member-pod.yaml"

  until ! crictl ps | grep etcd; do
    echo "Waiting for etcd to go down"
    sleep 10
  done

  touch stop-etcd.done
fi

if [ ! -f master-ignition.done ]; then
  echo "Creating master ignition and writing it to disk"
  # Get the master ignition from MCS
  curl --header 'Accept:application/vnd.coreos.ignition+json;version=3.1.0' \
    http://localhost:22624/config/master -o /opt/openshift/original-master.ign

  GATHER_ID="bootstrap"
  GATHER_TAR_FILE=log-bundle-${GATHER_ID}.tar.gz

  echo "Gathering installer bootstrap logs"
  TAR_FILE=${GATHER_TAR_FILE} /usr/local/bin/installer-gather.sh --id ${GATHER_ID}

  echo "Adding bootstrap control plane and bootstrap installer-gather bundle to master ignition"
  bootkube_podman_run \
    --rm \
    --privileged \
    --volume "$PWD:/assets:z" \
    --volume "/usr/local/bin/:/assets/bin" \
    --volume "/var/lib/etcd/:/assets/etcd-data" \
    --volume "/etc/kubernetes:/assets/kubernetes" \
    "${CLUSTER_BOOTSTRAP_IMAGE}" \
    bootstrap-in-place \
    --asset-dir /assets \
    --input /assets/bootstrap-in-place/master-update.fcc \
    --output /assets/master.ign

  touch master-ignition.done
fi
