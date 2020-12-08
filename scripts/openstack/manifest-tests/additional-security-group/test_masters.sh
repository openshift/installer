#!/usr/bin/env bash

set -Eeuo pipefail

declare -r assets_dir="$1"
declare result='PASS'

declare -a machines=(
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-0.yaml"
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-1.yaml"
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-2.yaml"
)

for machine in "${machines[@]}"; do
	if ! [ -f "$machine" ]; then
		>&2 echo "Machine resource $machine not found"
		result='FAIL'
	fi

	if ! >/dev/null yq -e '.spec.providerSpec.value.securityGroups[] | select(.uuid=="aaaaaaaa-bbbb-4ccc-dddd-111111111111")' "$machine"; then
		>&2 echo "Security group UUID not found in Machine $machine"
		>&2 echo
		>&2 echo 'The file was:'
		>&2 cat "$machine"
		>&2 echo
		result='FAIL'
	fi
done

if [ "$result" != 'PASS' ]; then
	exit 1
fi
