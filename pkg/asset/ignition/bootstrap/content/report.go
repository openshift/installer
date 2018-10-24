package content

const (
	// ReportSystemdContents is a service that reports the bootstrap progress
	// via a Kubernetes Event.
	ReportSystemdContents = `
[Unit]
Description=Report the completion of the cluster bootstrap process
# Workaround for https://github.com/systemd/systemd/issues/1312
Wants=bootkube.service tectonic.service
After=bootkube.service tectonic.service

[Service]
# Workaround for https://github.com/systemd/systemd/issues/1312 and https://github.com/opencontainers/runc/pull/1807
ExecStartPre=/usr/bin/test -f /opt/tectonic/.bootkube.done
ExecStartPre=/usr/bin/test -f /opt/tectonic/.tectonic.done
ExecStart=/usr/local/bin/report-progress.sh /opt/tectonic/auth/kubeconfig bootstrap-complete "cluster bootstrapping has completed"

Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
`

	// ReportShFileContents is a script for reporting the bootstrap progress.
	ReportShFileContents = `#!/usr/bin/env bash
set -e

KUBECONFIG="${1}"
NAME="${2}"
MESSAGE="${3}"
TIMESTAMP="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

echo "Reporting install progress..."

oc --config="$KUBECONFIG" create -f - <<EOF
apiVersion: v1
kind: Event
metadata:
  name: "${NAME}"
  namespace: kube-system
involvedObject:
  namespace: kube-system
message: "${MESSAGE}"
firstTimestamp: "${TIMESTAMP}"
lastTimestamp: "${TIMESTAMP}"
count: 1
source:
  component: cluster
  host: $(hostname)
EOF
`
)
