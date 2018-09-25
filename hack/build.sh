#!/bin/sh

set -e

cd "$(dirname "$0")/.."

CGO_ENABLED=0 go build -o ./bin/openshift-install ./cmd/openshift-install
