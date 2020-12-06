#!/usr/bin/env bash

set -Eeuxo pipefail

declare -r assets_dir="$1"

declare -a machines=(
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-0.yaml"
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-1.yaml"
	"${assets_dir}/openshift/99_openshift-cluster-api_master-machines-2.yaml"
)

declare -i exit_code=0

for machine in "${machines[@]}"; do
	if ! [ -f "$machine" ]; then
		>&2 echo "Machine resource $machine not found"
		exit_code=$((exit_code+1))
	fi

	if ! >/dev/null yq -e '.spec.providerSpec.value.securityGroups[] | select(.uuid=="aaaaaaaa-bbbb-4ccc-dddd-111111111111")' "$machine"; then
		>&2 echo "Security group UUID not found in Machine $machine"
		>&2 echo
		>&2 echo 'The file was:'
		>&2 cat "$machine"
		>&2 echo
		exit_code=$((exit_code+1))
	fi
done

exit $exit_code
