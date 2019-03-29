#!/bin/bash
# cidr_to_ip - 
#  https://www.terraform.io/docs/providers/external/data_source.html
#  Based on info from here: https://gist.github.com/irvingpop/968464132ded25a206ced835d50afa6b
#  This script takes the CIDR address and cycles through looking for the first available address
#  echo '{"cidr": "139.178.89.192/26", "control_plane_count": "3", "compute_count": "3", "cluster_domain": "dphillip.devcluster.openshift.com", "ipam": "ipam_address", "ipam_token", "api_token" }' | ./cidr_to_ip.sh
function error_exit() {
  echo "$1" 1>&2
  exit 1
}

function check_deps() {
  test -f $(which jq) || error_exit "jq command not detected in path, please install it"
  test -f $(which ipcalc) || error_exit "ipcalc command not detected in path, please install it"

}

function parse_input() {
  # jq reads from stdin so we don't have to set up any inputs, but let's validate the outputs
  eval "$(jq -r '@sh "export CIDR=\(.cidr) control_plane_count=\(.control_plane_count)  compute_count=\(.compute_count) cluster_domain=\(.cluster_domain) ipam=\(.ipam) ipam_token=\(.ipam_token)"')"
  if [[ -z "${CIDR}" ]]; then export CIDR=none; fi
  if [[ -z "${control_plane_count}" ]]; then export control_plane_count=none; fi
  if [[ -z "${compute_count}" ]]; then export compute_count=none; fi
  if [[ -z "${cluster_domain}" ]]; then export cluster_domain=none; fi
  if [[ -z "${ipam}" ]]; then export ipam=none; fi
  if [[ -z "${ipam_token}" ]]; then export ipam_token=none; fi
}

function produce_output() {
  cidr=$CIDR

  # Build the curl and run it
  lo=$(ipcalc -n $cidr | cut -f2 -d=)
   
  bs_count=0
  cp_count=0
  c_count=0

  
  if [[ $bs_count -ne 1 ]]
  then
    query_api=$(curl "http://$ipam/api/getFreeIP.php?apiapp=address&apitoken=$ipam_token&subnet=$lo&host=bootstrap-0.${cluster_domain}")
    bootstrap_ip=$query_api
    bs_count=1
  fi

  # check cluster_domain DNS first
  for CP in $(seq 0 $((control_plane_count-1)))
  do
     query_api=$(curl "http://$ipam/api/getFreeIP.php?apiapp=address&apitoken=$ipam_token&subnet=$lo&host=control-plane-$CP.${cluster_domain}")
     control_plane_ips+="$query_api "
     cp_count=$((cp_count+1))
	
  done
  for C in $(seq 0 $((compute_count-1)))
  do 
    query_api=$(curl "http://$ipam/api/getFreeIP.php?apiapp=address&apitoken=$ipam_token&subnet=$lo&host=compute-$C.${cluster_domain}") 
    compute_ips+="$query_api "
    c_count=$((c_count+1))
  done
  
  if [[ $bs_count -eq 1 ]] && [[ $cp_count -eq $control_plane_count ]] && [[ $c_count -eq $compute_count ]]
  then
  	jq -n \
    	--arg bootstrap_ip "$bootstrap_ip" \
    	--arg control_plane_ips "$control_plane_ips" \
    	--arg compute_ips "$compute_ips" \
    	'{"bootstrap_ip":$bootstrap_ip,"control_plane_ips":$control_plane_ips,"compute_ips":$compute_ips}'
	exit 0
  fi
}

# main()
check_deps
parse_input
produce_output
