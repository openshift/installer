#!/usr/bin/env bash

# This script is called instead of installer/data/data/bootstrap/files/usr/local/bin/installer-gather.sh for post-pivot
# bootstrap-in-place log gathering.
# Instead of gathering logs from the various masters, it simply runs the installer-masters-gather.sh script on itself.
# Instead of gathering bootstrap logs from "itself", it unpacks the bundle transferred to it from the bootstrap phase
# via the master ignition.
# The script will also gather the bootstrap-in-place-post-reboot service logs.

if test "x${1}" = 'x--id'
then
	GATHER_ID="${2}"
	shift 2
fi

ARTIFACTS="/tmp/bip-artifacts-${GATHER_ID}"
MASTER_ARTIFACTS="/tmp/artifacts-${GATHER_ID}"
mkdir -p "${ARTIFACTS}"

exec &> >(tee "${ARTIFACTS}/gather.log")

mkdir -p "${ARTIFACTS}/control-plane/master"
sudo /usr/local/bin/installer-masters-gather.sh --id "${GATHER_ID}" </dev/null
cp -r "$MASTER_ARTIFACTS"/* "${ARTIFACTS}/control-plane/master/"

BOOTSTRAP_PHASE_LOG_BUNDLE_NAME="log-bundle-bootstrap-in-place-pre-reboot"
BOOTSTRAP_PHASE_LOG_BUNDLE_ARCHIVE_PATH="/var/log/$BOOTSTRAP_PHASE_LOG_BUNDLE_NAME.tar.gz"

if [[ -f $BOOTSTRAP_PHASE_LOG_BUNDLE_ARCHIVE_PATH ]]; then
  echo "Found log bundle from bootstrap phase"
  tar -xzf $BOOTSTRAP_PHASE_LOG_BUNDLE_ARCHIVE_PATH
  cp -r $BOOTSTRAP_PHASE_LOG_BUNDLE_NAME/* "${ARTIFACTS}/"
  rm -rf $BOOTSTRAP_PHASE_LOG_BUNDLE_NAME
else
  echo "Bootstrap phase logs not found, including only master logs in log bundle"
fi

echo "Gathering bootstrap-in-place-post-reboot journal"
mkdir -p "${ARTIFACTS}/bootstrap-in-place/journals"
journalctl --boot --no-pager --output=short --unit=bootstrap-in-place-post-reboot > "${ARTIFACTS}/bootstrap-in-place/journals/bootstrap-in-place-post-reboot.log"

TAR_FILE="${TAR_FILE:-${HOME}/log-bundle-${GATHER_ID}.tar.gz}"
tar cz -C "${ARTIFACTS}" --transform "s?^\\.?log-bundle-${GATHER_ID}?" . > "${TAR_FILE}"
echo "Log bundle written to ${TAR_FILE}"
