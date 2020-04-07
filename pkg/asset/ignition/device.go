package ignition

import (
        "text/template"
	"net"
        "bytes"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/coreos/ignition/config/util"
        "github.com/pkg/errors"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// This creates a script which would create a node network interface with IP address with the same last 4 bits as the ens3.4094 interface IP

var networkScriptTmpl = template.Must(template.New("user-data").Parse(`#!/bin/bash

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

# Add iptable rule to accept igmp
iptables -I INPUT 1 -j ACCEPT -p igmp

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

func IgnitionFiles(installConfig *installconfig.InstallConfig) []igntypes.File {
	//installConfig := &installconfig.InstallConfig{}
        machineCIDR := &installConfig.Config.Networking.DeprecatedMachineCIDR.IPNet
        defaultGateway, _ := cidr.Host(machineCIDR, 1)
        kube_api_vlan := installConfig.Config.Platform.OpenStack.AciNetExt.KubeApiVLAN
        infra_vlan := installConfig.Config.Platform.OpenStack.AciNetExt.InfraVLAN
        mtu_value := installConfig.Config.Platform.OpenStack.AciNetExt.Mtu
        networkScriptString, _ := NetworkScript(kube_api_vlan, defaultGateway.String(), mtu_value)

        neutronCIDR := &installConfig.Config.Platform.OpenStack.AciNetExt.NeutronCIDR.IPNet
        defaultNeutronGateway, _ := cidr.Host(neutronCIDR, 1)
        defaultNeutronGatewayStr := defaultNeutronGateway.String()

        installerHostSubnet := installConfig.Config.Platform.OpenStack.AciNetExt.InstallerHostSubnet
        installerHostIP, installerHostNet, _ := net.ParseCIDR(installerHostSubnet)
        installerNetmask := net.IP(installerHostNet.Mask)
	
	ifcfg_ens3_string := `DEVICE=ens3.4094
               ONBOOT=yes
               BOOTPROTO=dhcp
               MTU=` + mtu_value + `
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

        ifcfg_opflex_conn_string := `VLAN=yes
               TYPE=Vlan
               PHYSDEV=ens3
               VLAN_ID=` + infra_vlan + `
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
               DEVICE=ens3.`+ infra_vlan +`
               ONBOOT=yes
               MTU=` + mtu_value

        ifcfg_uplink_conn_string := `TYPE=Ethernet
               PROXY_METHOD=none
               BROWSER_ONLY=no
               DEFROUTE=yes
               IPV4_FAILURE_FATAL=no
               IPV6INIT=no
               NAME=uplink-conn
               DEVICE=ens3
               ONBOOT=yes
               BOOTPROTO=none
               MTU=` + mtu_value

        route_opflex_conn_string := `ADDRESS0=224.0.0.0
               NETMASK0=240.0.0.0
               METRIC0=1000`

        route_ens3_string := `ADDRESS0=` + installerHostIP.String() + `
               NETMASK0=` + installerNetmask.String() + `
               METRIC0=1000
               GATEWAY0=` + defaultNeutronGatewayStr

	var ignitionFiles []igntypes.File
	ignitionFiles = append(ignitionFiles,
				FileFromString("/usr/local/bin/kube-api-interface.sh", "root", 0555, string(networkScriptString)),
				FileFromString("/etc/sysconfig/network-scripts/ifcfg-ens3.4094", "root", 0420, ifcfg_ens3_string),
				FileFromString("/etc/sysconfig/network-scripts/ifcfg-opflex-conn", "root", 0420, ifcfg_opflex_conn_string),
				FileFromString("/etc/sysconfig/network-scripts/ifcfg-uplink-conn", "root", 0420, ifcfg_uplink_conn_string),
                        	FileFromString("/etc/sysconfig/network-scripts/route-opflex-conn", "root", 0420, route_opflex_conn_string),
                        	FileFromString("/etc/sysconfig/network-scripts/route-ens3.4094", "root", 0420, route_ens3_string))

	return ignitionFiles
}

func SystemdUnitFiles(installConfig *installconfig.InstallConfig) []igntypes.Unit {
	var systemdUnits []igntypes.Unit

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
	systemdUnits = append(systemdUnits, nodeService)

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
	systemdUnits = append(systemdUnits, machineConfigDaemonPath)

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
	systemdUnits = append(systemdUnits, machineConfigDaemonService)

	return systemdUnits
}
