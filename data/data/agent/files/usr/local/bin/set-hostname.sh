#!/bin/bash

cat <<EOF >/etc/motd
The primary service is agent.service. To watch its status, run:

  journalctl -u agent.service

To view the agent log, run:

  journalctl TAG=agent
EOF
echo "Waiting for network to determine if this is the rendezvous host." > /etc/motd.d/60-rendezvous-host

#
# The hostnames defined in agent-config.yaml are written out
# to files at /etc/assisted/hostnames/<MAC-address>.
#
# If a host has multiple interfaces, the host's first network
# interface's MAC address is used.
#
# This script compares the MAC addresses on the current host
# with the addresses in /etc/assisted/hostnames/.
#
# If a match is found, then the hostname in the file is set
# as this host's hostname.
#

# This is a workaround to handle the case when NetworkManager
# cannot get the hostname from the IP via reverse DNS
function lookup_address() {

   ips=$(hostname -I)
   while [[ "$ips" == "" ]]
   do
      echo "Waiting for IPs to appear for reverse DNS" 1>&2
      sleep 5
      ips=$(hostname -I)
   done

   for ip in ${ips}
   do
      echo "Requesting hostname lookup for $ip" 1>&2
      hostname=$(dig -x "$ip" +short)
      # If hostname returned, configure first entry
      for hn in ${hostname}; do
         if [[ "$hn" != "" ]]; then
            hostnamectl set-hostname "${hn}"
            echo "Setting hostname to ${hn}" 1>&2
	    break
         fi
      done
   done
}

hostname_set=0
HOSTNAMES_PATH=/etc/assisted/hostnames
FILES=$(ls $HOSTNAMES_PATH)
for filename in ${FILES}
do
    MATCHED_MAC_ADDRESS_WITH_HOST=$(ip address | grep "${filename}")
    if [ "$MATCHED_MAC_ADDRESS_WITH_HOST" != "" ]; then
        HOSTNAME="$(cat "${HOSTNAMES_PATH}/${filename}")"
        echo "Host has matching MAC address: ${filename}" 1>&2
        echo "Setting hostname to ${HOSTNAME}" 1>&2
        hostnamectl set-hostname "${HOSTNAME}"
	hostname_set=1
    else
        echo "MAC address, ${filename}, does not exist on this host" 1>&2
    fi
done

if [[ $hostname_set -eq 0 ]]; then
   lookup_address
fi
