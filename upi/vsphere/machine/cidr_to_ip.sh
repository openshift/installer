#!/bin/bash
# cidr_to_ip - 
#  https://www.terraform.io/docs/providers/external/data_source.html
#  Based on info from here: https://gist.github.com/irvingpop/968464132ded25a206ced835d50afa6b
#  This script takes requests an IP address from an IPAM server
#  echo '{"cidr": "139.178.89.192/26", "hostname": "control-plane-0.dphillip.devcluster.openshift.com", "ipam": "ipam_address", "ipam_token", "api_token" }' | ./cidr_to_ip.sh
function error_exit() {
  echo "$1" 1>&2
  exit 1
}

function check_deps() {
  test -f "$(command -v jq)" || error_exit "jq command not detected in path, please install it"
  test -f "$(command -v ipcalc)" || error_exit "ipcalc command not detected in path, please install it"

}

function parse_input() {
  # jq reads from stdin so we don't have to set up any inputs, but let's validate the outputs
  eval "$(jq -r '@sh "export CIDR=\(.cidr) hostname=\(.hostname) ipam=\(.ipam) ipam_token=\(.ipam_token)"')"
  if [[ -z "${CIDR}" ]]; then export CIDR=none; fi
  if [[ -z "${hostname}" ]]; then export hostname=none; fi
  if [[ -z "${ipam}" ]]; then export ipam=none; fi
  if [[ -z "${ipam_token}" ]]; then export ipam_token=none; fi
}

function produce_output() {
  cidr=$CIDR

  # Build the curl and run it
  lo=$(ipcalc -n $cidr | cut -f2 -d=)
  
  # Request an IP address. Verify that the IP address reserved matches the IP
  # address returned. Loop until the reservation matches the address returned.
  # The verification and looping is a crude way of overcoming the lack of
  # currency safety in the IPAM server.
  while true
  do 
    ip_address=$(curl -s "http://$ipam/api/getFreeIP.php?apiapp=address&apitoken=$ipam_token&subnet=$lo&host=${hostname}")

    if ! [[ $ip_address =~ ^[0-9]{1,3}(\.[0-9]{1,3}){3}$ ]]; then error_exit "could not reserve an IP address: ${ip_address}"; fi
	
    reserved_ip=$(curl -s "http://$ipam/api/getIPs.php?apiapp=address&apitoken=$ipam_token&domain=${hostname}" | \
      jq -r ".\"${hostname}\"")

    if [[ "$ip_address" == "$reserved_ip" ]]
	then
      jq -n \
        --arg ip_address "$ip_address" \
        '{"ip_address":$ip_address}'
      exit 0
	fi
  done
}

# main()
check_deps
parse_input
produce_output
