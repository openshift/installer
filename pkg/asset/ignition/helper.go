package ignition

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/coreos/ignition/config/util"
	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/stretchr/testify/assert"
)

func CheckIgnitionFiles(t *testing.T, ignConfig *igntypes.Config) {

	actualStringList := [6]string{ GetActualKubeApiInterfaceStr(),
				GetActualIfcfgEns34094Str(),
				GetActualIfcfgOpflexConnStr(),
				GetActualIfcfgUplinkConnStr(),
				GetActualRouteOpflexConnStr(),
				GetActualRouteEns34094Str() }

	for i := 0; i < 6; i++ {
		expected := ignConfig.Storage.Files[i].FileEmbedded1.Contents.Source
		actualStr := actualStringList[i]
		compareScripts(t, expected, actualStr)
	}

}

func compareScripts(t *testing.T, expected string, actualStr string) {
        expectedSplit := strings.Split(expected, ",")[1]
        expectedDecoded, _ := base64.StdEncoding.DecodeString(expectedSplit)
        expectedStr := string(expectedDecoded)
        assert.Equal(t, expectedStr, actualStr, "unexpected ignition file")
}

func CheckSystemdUnitFiles(t *testing.T, actualUnits []igntypes.Unit) {

	expectedUnitData := [3]igntypes.Unit{GetNodeInterfaceService(),
					GetMachineConfigDaemonForcePath(),
					GetMachineConfigDaemonService()}
					
        for i, u := range actualUnits {
		assert.Equal(t, u, expectedUnitData[i], "unexpected " + expectedUnitData[i].Name)
        }

}

func GetActualKubeApiInterfaceStr() string {

	actualKubeApiInterfaceStr := `#!/bin/bash

# These are rendered through Go
KUBE_API_VLAN=1021
DEFAULT_GATEWAY=1.2.3.1
MTU_VALUE=1600

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

# Add iptable rule to accept igmp
iptables -I INPUT 1 -j ACCEPT -p igmp

`
	return actualKubeApiInterfaceStr
}

func GetActualIfcfgEns34094Str() string {

	actualIfcfgEns34094Str := `DEVICE=ens3.4094
               ONBOOT=yes
               BOOTPROTO=dhcp
               MTU=1600
               TYPE=Vlan
               VLAN=yes
               PHYSDEV=ens3
               VLAN_ID=4094
               REORDER_HDR=yes
               GVRP=no
               MVRP=no
               PROXY_METHOD=none
               BROWSER_ONLY=no
               DEFROUTE=no
               IPV4_FAILURE_FATAL=no
               IPV6INIT=no`

	return actualIfcfgEns34094Str
}

func GetActualIfcfgOpflexConnStr() string {

	actualIfcfgOpflexConnStr := `VLAN=yes
               TYPE=Vlan
               PHYSDEV=ens3
               VLAN_ID=4094
               REORDER_HDR=yes
               GVRP=no
               MVRP=no
               PROXY_METHOD=none
               BROWSER_ONLY=no
               BOOTPROTO=dhcp
               DEFROUTE=no
               IPV4_FAILURE_FATAL=no
               IPV6INIT=no
               NAME=opflex-conn
               DEVICE=ens3.4094
               ONBOOT=yes
               MTU=1600`

	return actualIfcfgOpflexConnStr
}

func GetActualIfcfgUplinkConnStr() string { 

        actualIfcfgUplinkConnStr := `TYPE=Ethernet
               PROXY_METHOD=none
               BROWSER_ONLY=no
               DEFROUTE=yes
               IPV4_FAILURE_FATAL=no
               IPV6INIT=no
               NAME=uplink-conn
               DEVICE=ens3
               ONBOOT=yes
               BOOTPROTO=none
               MTU=1600`

	return actualIfcfgUplinkConnStr
}

func GetActualRouteOpflexConnStr() string {

	actualRouteOpflexConnStr := `ADDRESS0=224.0.0.0
               NETMASK0=240.0.0.0
               METRIC0=1000`
	return actualRouteOpflexConnStr
}

func GetActualRouteEns34094Str() string {

	actualRouteEns34094Str := `ADDRESS0=9.10.11.12
               NETMASK0=255.192.0.0
               METRIC0=1000
               GATEWAY0=5.6.7.9`

	return actualRouteEns34094Str
}

func GetNodeInterfaceService() igntypes.Unit {
	nodeService := igntypes.Unit{
		Name:    "node-interface.service",
		Enabled: util.BoolToPtr(true),
		Contents: `[Unit]
		Description=Adding Node Network Interface to MachineSet
		Wants=network-online.target
		After=network-online.target
		[Service]
		Type=simple
		ExecStart=/usr/local/bin/kube-api-interface.sh
		[Install]
		WantedBy=multi-user.target`}
	return nodeService
}

func GetMachineConfigDaemonForcePath() igntypes.Unit {
	machineConfigDaemonPath := igntypes.Unit{
		Name:    "machine-config-daemon-force.path",
		Enabled: util.BoolToPtr(true),
		Contents: `[Unit]
Description=Path File for Disabling Machine-Config Validation Check
[Path]
PathChanged=/run/machine-config-daemon-force
Unit=machine-config-daemon-force.service
[Install]
WantedBy=multi-user.target`}
	return machineConfigDaemonPath
}

func GetMachineConfigDaemonService() igntypes.Unit {
	machineConfigDaemonService := igntypes.Unit{
		Name:    "machine-config-daemon-force.service",
		Enabled: util.BoolToPtr(true),
		Contents: `[Unit]
Description=Disabling Machine-Config Validation Check
[Service]
Type=simple
ExecStart=touch /run/machine-config-daemon-force
[Install]
WantedBy=multi-user.target`}
	return machineConfigDaemonService
}
