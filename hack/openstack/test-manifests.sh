#!/usr/bin/env bash

set -Eeuo pipefail

[[ $(type -P yq) ]] || { >&2 echo "Required tool 'yq' not found in PATH" ; exit 1; }

declare \
	tests_dir='./scripts/openstack/manifest-tests' \
	openshift_install='./bin/openshift-install' \
	api_fip='' \
	os_cloud='' \
	external_network='' \
	compute_flavor=''

# Let install-config describe a configuration that is incompatible with the
# target CI infrastructure
export OPENSHFIT_INSTALL_SKIP_PREFLIGHT_VALIDATIONS=1

print_help() {
	set +x

	echo -e "Test the OpenStack manifest generation."
	echo
	echo -e "Required configuration:"
	echo
	echo -e "\\t-c\\tOS_CLOUD"
	echo -e "\\t-e\\tExternal network"
	echo -e "\\t-f\\tA valid flavor"
	echo
	echo -e "Use:"
	echo -e "\\t${0} -c <cloud> -e <external network> -f <flavor> [-a <fip>] [-i <openshift-install>] [-t <test-dir>]"
	echo
	echo -e 'Additional arguments:'
	echo -e "\\t-a\\tapiFloatingIP"
	echo -e "\\t-i\\tpath to openshift-install (defaults to '${openshift_install}')"
	echo -e "\\t-t\\tpath to the tests to be run (defaults to '${tests_dir}')"
}

fill_install_config() {
	declare -r \
		template="$1" \
		pull_secret="'"'{"auths":{"registry.svc.ci.openshift.org":{"auth":"QW4gYWN0dWFsIHB1bGwgc2VjcmV0IGlzIG5vdCBuZWNlc3NhcnkK"}}}'"'"

	sed '
		s|${\?OS_CLOUD}\?|'"${os_cloud}"'|;
		s|${\?EXTERNAL_NETWORK}\?|'"${external_network}"'|;
		s|${\?COMPUTE_FLAVOR}\?|'"${compute_flavor}"'|;
		s|${\?API_FIP}\?|'"${api_fip}"'|;
		s|${\?PULL_SECRET}\?|'"${pull_secret}"'|;
		' "$template"
}

validate_configuration() {
	declare -a required_values=("os_cloud" "external_network" "compute_flavor")
	declare fail=false

	for val in "${required_values[@]}"; do
		declare required=${!val:-}
		if [ -z "${required}" ]; then
			>&2 echo "Missing required argument '${val}'."
			fail=true
		fi
	done

	if [ "$fail" = true ]; then
		print_help
		exit 1
	fi
}

while getopts a:c:e:f:i:t:h o; do
	case "$o" in
		a) api_fip="$OPTARG" ;;
		c) os_cloud="$OPTARG" ;;
		e) external_network="$OPTARG" ;;
		f) compute_flavor="$OPTARG" ;;
		i) openshift_install="$OPTARG" ;;
		t) tests_dir="$OPTARG" ;;
		h) print_help; exit 0 ;;
		*) print_help; exit 1 ;;
	esac
done
readonly api_fip os_cloud external_network compute_flavor openshift_install tests_dir

declare -a temp_dirs
cleanup() {
	for temp_dir in "${temp_dirs[@]}"; do
		rm -rf "$temp_dir"
	done
}
trap cleanup EXIT

validate_configuration

>&2 echo "Running the tests from '${tests_dir}' against the Installer binary '${openshift_install}'."

declare result='PASS'
for testcase in "${tests_dir}"/* ; do
	if [ -d "$testcase" ]; then
		assets_dir="$(mktemp -d)"
		temp_dirs+=("$assets_dir")
		fill_install_config "${testcase}/install-config.yaml" > "${assets_dir}/install-config.yaml"
		"$openshift_install" create manifests --dir "$assets_dir"
		for t in "${testcase}"/test_*; do
			if $t "$assets_dir"; then
				echo "PASS: '$t'"
			else
				result='FAIL'
				echo "FAIL: '$t'"
			fi
		done
	fi
done

if [ "$result" != 'PASS' ]; then
	echo "FAIL"
	exit 1
fi

echo 'PASS'
