package ignition

import (
        "text/template"
        "bytes"
        "github.com/pkg/errors"
)

// This creates a script which would create a node network interface with IP address with the same last 4 bits as the ens3.4094 interface IP

var networkScriptTmpl = template.Must(template.New("user-data").Parse(`#!/bin/bash

#Disable MCO Validation Check
touch /run/machine-config-daemon-force

# These are rendered through Go
KUBE_API_VLAN={{.vlan}}
DEFAULT_GATEWAY={{.defGateway}}
MTU_VALUE={{.mtu}}

IFC_4094=ens3.4094
KUBE_API_VLAN_DEVICE="ens3.${KUBE_API_VLAN}"
FILE_PATH="/etc/sysconfig/network-scripts/ifcfg-api-conn"

# Check if the kube_api interface is present
ifconfig $KUBE_API_VLAN_DEVICE > /dev/null
result=$?
if [[ $result == 0 ]]; then
  echo "Interface created, nothing to do"
  exit 0
else
  echo "Interface not created, continue"
fi

# Find IP address allocated to 4094
ip_string=$(ifconfig $IFC_4094| awk -F ' *|:' '/inet /{print $3}')
if [[ $ip_string == "" ]]; then
  echo "No IP"
else
  echo "Found IP $ip_string"
fi

# Parse out the last value of the address
parse=$(echo $ip_string | tr "." "\n")
last_bits=""
i=0
for addr in $parse
do
    if [ $i -eq 3 ]; then
      last_bits="$addr"
    fi
    ((i=i+1))
done
echo "Last 4 bits are $last_bits"

# Create the node IP out of def gateway and existing bits
i=0
new_ip=""
parse1=$(echo $DEFAULT_GATEWAY | tr "." "\n")
for addr in $parse1
do
    if [ $i -eq 0 ]; then
      new_ip=$addr
    elif [ $i -eq 3 ]; then
      new_ip="${new_ip}.${last_bits}"
    else
      new_ip="${new_ip}.${addr}"
    fi
    ((i=i+1))
done
echo Generated IP is $new_ip

# Create the network script
echo Writing network script file to $FILE_PATH
/bin/cat <<EOM >$FILE_PATH
VLAN=yes
MTU=$MTU_VALUE
TYPE=Vlan
PHYSDEV=ens3
VLAN_ID=$KUBE_API_VLAN
REORDER_HDR=yes
GVRP=no
MVRP=no
PROXY_METHOD=none
BROWSER_ONLY=no
BOOTPROTO=none
IPADDR=$new_ip
PREFIX=24
DEFROUTE=yes
GATEWAY=$DEFAULT_GATEWAY
PEERDNS=no
IPV4_FAILURE_FATAL=no
IPV6INIT=no
NAME=api-conn
DEVICE=$KUBE_API_VLAN_DEVICE
ONBOOT=yes
METRIC=90
EOM

# Change permissions and owner of the network script
chmod 420 $FILE_PATH

# Restarting Network Manager
systemctl restart NetworkManager

`))

func NetworkScript(vlan string, defGateway string, mtu string) ([]byte, error) {
        buf := &bytes.Buffer{}
        data := map[string]string{
		"vlan":          vlan,
                "defGateway":    defGateway,
                "mtu":           mtu,
	}
	if err := networkScriptTmpl.Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute user-data template")
	}
        return buf.Bytes(), nil
}
