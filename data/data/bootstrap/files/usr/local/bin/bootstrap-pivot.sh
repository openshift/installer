#!/usr/bin/env bash
set -euo pipefail
# Update bootstrap node with latest content. Required for OKD as it starts with FCOS image
# and needs to be updated to OKD's layered fedora-coreos image used by masters / nodes.

# shellcheck disable=SC1091
. /usr/local/bin/bootstrap-service-record.sh

# Rebase the FCOS bootstrap node to OKD's fedora-coreos image
if [ ! -f /etc/okd-booted.stamp ]; then
  # shellcheck disable=SC1091
  . /usr/local/bin/release-image.sh

  MACHINE_OS_IMAGE=$(image_for fedora-coreos)
  echo "Pulling ${MACHINE_OS_IMAGE}..."
  while true
  do
      record_service_stage_start "pull-fedora-coreos image"
      if podman pull --quiet "${MACHINE_OS_IMAGE}"
      then
          record_service_stage_success
          break
      else
          record_service_stage_failure
          echo "Pull failed. Retrying ${MACHINE_OS_IMAGE}..."
      fi
  done

  # Rebase to fedora-coreos image
  record_service_stage_start "rebase-to-fedora-coreos-image"
  rpm-ostree rebase --experimental "ostree-unverified-registry:${MACHINE_OS_IMAGE}"

  # Change default kargs
  # We have to do this here for now because it's not supported
  # in the layered image build, yet. See:
  # https://github.com/coreos/rpm-ostree/issues/3738
  rpm-ostree kargs \
    --replace mitigations=auto,nosmt=off \
    --append intel_pstate=disable && \

  record_service_stage_success

  systemctl reboot
fi

# Copy manifests
record_service_stage_start "copy-okd-bootstrap-manifests"
cp /manifests/* /opt/openshift/openshift/ -rvf
record_service_stage_success
