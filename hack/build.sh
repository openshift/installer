#!/bin/sh

set -e

cd "$(dirname "$0")/.."

archive=$(tar cz steps modules config.tf | base64 -w 0)
LDFLAGS="-X github.com/openshift/installer/pkg/asset/cluster.TemplateArchive=${archive}"

CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ./bin/openshift-install ./cmd/openshift-install
