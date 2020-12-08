#!/usr/bin/env bash

set -Eeuo pipefail

declare -r assets_dir="$1"
declare result='PASS'

declare machineset="${assets_dir}/openshift/99_openshift-cluster-api_worker-machineset-0.yaml"

if ! [ -f "$machineset" ]; then
	>&2 echo 'MachineSet not found'
	result='FAIL'
fi

if ! >/dev/null yq -e '.spec.template.spec.providerSpec.value.securityGroups[] | select(.uuid=="aaaaaaaa-bbbb-4ccc-dddd-000000000000")' "$machineset"; then
	>&2 echo 'Security group UUID not found in the MachineSet'
	>&2 echo
	>&2 echo 'The file was:'
	>&2 cat "$machineset"
	>&2 echo
	result='FAIL'
fi

if [ "$result" != 'PASS' ]; then
	exit 1
fi
