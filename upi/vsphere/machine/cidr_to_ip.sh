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
  input=$(jq .)
  CIDR=$(echo "$input" | jq -r .cidr)
  hostname=$(echo "$input" | jq -r .hostname)
  ipam=$(echo "$input" | jq -r .ipam)
  ipam_token=$(echo "$input" | jq -r .ipam_token)
}

is_ip_address() {
  if [[ $1 =~ ^[0-9]{1,3}(\.[0-9]{1,3}){3}$ ]]
  then
    echo "true"
  else
    echo "false"
  fi
}

get_reservation() {
  reserved_ip=$(curl -s "http://${ipam}/api/getIPs.php?apiapp=address&apitoken=${ipam_token}&domain=${hostname}" | \
      jq -r ".\"${hostname}\"")
  if [ "$(is_ip_address "${reserved_ip}")" == "false" ]
  then
    reserved_ip=""
  fi
  echo $reserved_ip
}

function produce_output() {
  if [[ "${CIDR}" == "null" ]]
  then
    jq -n \
      --arg ip_address "$(get_reservation)" \
      '{"ip_address":$ip_address}'
	exit 0
  fi

  cidr=$CIDR

  # Build the curl and run it
  lo=$(ipcalc -n "$cidr" | cut -f2 -d=)
  
  # Request an IP address. Verify that the IP address reserved matches the IP
  # address returned. Loop until the reservation matches the address returned.
  # The verification and looping is a crude way of overcoming the lack of
  # currency safety in the IPAM server.
  while true
  do 
    ip_address=$(curl -s "http://$ipam/api/getFreeIP.php?apiapp=address&apitoken=$ipam_token&subnet=$lo&host=${hostname}")

    if [[ "$(is_ip_address "${ip_address}")" != "true" ]]; then error_exit "could not reserve an IP address: ${ip_address}"; fi
	
    if [[ "$ip_address" == "$(get_reservation)" ]]
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
