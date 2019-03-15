#!/bin/bash
# cidr_to_ip - 
#  https://www.terraform.io/docs/providers/external/data_source.html
#  Based on info from here: https://gist.github.com/irvingpop/968464132ded25a206ced835d50afa6b
#  This script takes the CIDR address and cycles through looking for the first available address

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
  eval "$(jq -r '@sh "export CIDR=\(.cidr) master_count=\(.master_count)  worker_count=\(.worker_count)"')"
  if [[ -z "${CIDR}" ]]; then export CIDR=none; fi
  if [[ -z "${master_count}" ]]; then export master_count=none; fi
  if [[ -z "${worker_count}" ]]; then export worker_count=none; fi
}

function produce_output() {

  cidr=$CIDR

  # range is bounded by network (-n) & broadcast (-b) addresses.
  lo=$(ipcalc -n $cidr | cut -f2 -d=)
  hi=$(ipcalc -b $cidr | cut -f2 -d=)

  read a b c d <<< $(echo $lo | tr . ' ')
  read e f g h <<< $(echo $hi | tr . ' ')
  IP_RANGE=$(eval echo {$a..$e}.{$b..$f}.{$c..$g}.{$d..$h})
   
  bs_count=0
  m_count=0
  w_count=0
  for IPADDR in ${IP_RANGE}
  do

	if [ $IPADDR != $(ipcalc -n $cidr | cut -f2 -d=) ] && [ $IPADDR != $(ipcalc -b $cidr | cut -f2 -d=) ] 
        then
  	  ping -c1 -w1 $IPADDR > /dev/null 2>&1
  	  ping_rc=$?
	
	  if [[ $ping_rc -eq 1 ]] && [[ $bs_count -ne 1 ]]
	  then
  		bootstrap_ip+="$IPADDR"
		bs_count=$((bs_count+1))
	  elif [[ $ping_rc -eq 1 ]] && [[ $m_count -ne $master_count ]]
	  then
  		master_ips+="$IPADDR "
		m_count=$((m_count+1))
	  elif [[ $ping_rc -eq 1 ]] && [[ $w_count -ne $worker_count ]]
	  then
  		worker_ips+="$IPADDR "
		w_count=$((w_count+1))
	  elif [[ $bs_count -eq 1 ]] && [[ $m_count -eq $master_count ]] && [[ $w_count -eq $worker_count ]]
	  then
  		jq -n \
    		--arg bootstrap_ip "$bootstrap_ip" \
    		--arg master_ips "$master_ips" \
    		--arg worker_ips "$worker_ips" \
    		'{"bootstrap_ip":$bootstrap_ip,"master_ips":$master_ips,"worker_ips":$worker_ips}'
	        exit 0
	   fi
      	 fi
  done
}

# main()
check_deps
parse_input
produce_output
