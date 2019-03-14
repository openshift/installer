#!/bin/bash
# cidr_to_ip - 
#   https://www.terraform.io/docs/providers/external/data_source.html
#
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
  eval "$(jq -r '@sh "export CIDR=\(.cidr)"')"
  if [[ -z "${CIDR}" ]]; then export CIDR=none; fi
}

function produce_output() {

  cidr=$CIDR

  # range is bounded by network (-n) & broadcast (-b) addresses.
  lo=$(ipcalc -n $cidr | cut -f2 -d=)
  hi=$(ipcalc -b $cidr | cut -f2 -d=)

  read a b c d <<< $(echo $lo | tr . ' ')
  read e f g h <<< $(echo $hi | tr . ' ')
  IP_RANGE=$(eval echo {$a..$e}.{$b..$f}.{$c..$g}.{$d..$h})
  
  count=0
  for IPADDR in ${IP_RANGE}
  do

	if [ $IPADDR != $(ipcalc -n $cidr | cut -f2 -d=) ] && [ $IPADDR != $(ipcalc -b $cidr | cut -f2 -d=) ] 
        then
  	  ping -c1 -w1 $IPADDR > /dev/null 2>&1
  	  ping_rc=$?
	
	  if [[ $ping_rc -eq 1 ]]
	  then
  		ipaddress+="$IPADDR "
		count=$((count+1))
	  fi
	  if [[ $count -eq 7 ]]
	  then
  		jq -n \
    		--arg ipaddress "$ipaddress" \
    		'{"ipaddress":$ipaddress}'
	        exit 0
	  fi
       fi
  done
}

# main()
check_deps
parse_input
produce_output
