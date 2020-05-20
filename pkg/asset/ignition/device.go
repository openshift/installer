package ignition

import (
        "text/template"
	"net"
        "bytes"
	"strconv"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/coreos/ignition/config/util"
        "github.com/pkg/errors"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// This creates a script which would add node labels to kubelet service
var cloudProviderScriptTmpl = template.Must(template.New("user-data").Parse(`#!/bin/bash
# These are rendered through Go
KUBE_API_VLAN={{.vlan}}
KUBE_API_VLAN_DEVICE="ens3.${KUBE_API_VLAN}"
ip=$(/sbin/ip -o -4 addr list $KUBE_API_VLAN_DEVICE | awk '{print $4}' | cut -d/ -f1)
retVal=1
while [ $retVal -ne 0 ]; do
echo "Node IP is ${ip}"
oc get nodes --kubeconfig=/var/lib/kubelet/kubeconfig  -o wide | grep $ip
retVal=$?
done
grep -zo 'cloud-provider=openstack \\' /etc/systemd/system/kubelet.service  || sed -i '/kubelet \\/a\      \--cloud-provider=openstack \\\n      --cloud-config=/etc/kubernetes/cloud.conf \\' /etc/systemd/system/kubelet.service
systemctl daemon-reload
systemctl restart kubelet`))

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
PREFIX={{.machine_cidr_mask}}
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

func CloudProviderScript(vlan string) ([]byte, error) {
	buf := &bytes.Buffer{}
	data := map[string]string{
		"vlan":          vlan,
	}
	if err := cloudProviderScriptTmpl.Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute user-data template")
	}
	return buf.Bytes(), nil
}

func NetworkScript(vlan string, defGateway string, mtu string, machineCIDRMask string) ([]byte, error) {
        buf := &bytes.Buffer{}
        data := map[string]string{
		"vlan":          	vlan,
		"defGateway":    	defGateway,
		"mtu":           	mtu,
                "machine_cidr_mask":	machineCIDRMask,
	}
	if err := networkScriptTmpl.Execute(buf, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute user-data template")
	}
        return buf.Bytes(), nil
}

func IgnitionFiles(installConfig *installconfig.InstallConfig, is_bootstrap bool) []igntypes.File {
	if installConfig.Config.Networking.NetworkType != "CiscoAci" {
		return nil
	}
        machineCIDR := &installConfig.Config.Networking.MachineNetwork[0].CIDR.IPNet
	machineStr := machineCIDR.String()
        machineCIDRMask := strings.Split(machineStr, "/")[1]
        defaultGateway, _ := cidr.Host(machineCIDR, 1)
        kubeApiVLAN := installConfig.Config.Platform.OpenStack.AciNetExt.KubeApiVLAN
        infraVLAN := installConfig.Config.Platform.OpenStack.AciNetExt.InfraVLAN
        mtuString, _ := strconv.Atoi(installConfig.Config.Platform.OpenStack.AciNetExt.Mtu)
        mtuValue := strconv.Itoa(mtuString - 100)
	cloudProviderScriptString, _ := CloudProviderScript(kubeApiVLAN)
        networkScriptString, _ := NetworkScript(kubeApiVLAN, defaultGateway.String(), mtuValue, machineCIDRMask)
        neutronCIDR := &installConfig.Config.Platform.OpenStack.AciNetExt.NeutronCIDR.IPNet
        defaultNeutronGateway, _ := cidr.Host(neutronCIDR, 1)
        defaultNeutronGatewayStr := defaultNeutronGateway.String()

        installerHostSubnet := installConfig.Config.Platform.OpenStack.AciNetExt.InstallerHostSubnet
        installerHostIP, installerHostNet, _ := net.ParseCIDR(installerHostSubnet)
        installerNetmask := net.IP(installerHostNet.Mask)
	
	ifcfg_ens3_string := `DEVICE=ens3.4094
               ONBOOT=yes
               BOOTPROTO=dhcp
               MTU=` + mtuValue + `
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
               VLAN_ID=` + infraVLAN + `
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
               DEVICE=ens3.`+ infraVLAN +`
               ONBOOT=yes
               MTU=` + mtuValue

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
               MTU=` + mtuValue

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

	if !is_bootstrap{
		ignitionFiles = append(ignitionFiles,FileFromString(
			"/usr/local/bin/node-cloud-provider.sh", "root", 0555,
			string(cloudProviderScriptString)))
	}

	return ignitionFiles
}

func SystemdUnitFiles(installConfig *installconfig.InstallConfig, is_bootstrap bool) []igntypes.Unit {
	if installConfig.Config.Networking.NetworkType != "CiscoAci" {
                return nil
        }
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

	if !is_bootstrap {
		nodeCloudProvider := igntypes.Unit{
			Name:    "node-cloud-provider.service",
			Enabled: util.BoolToPtr(true),
			Contents: `[Unit]
		Description=Assigning Cloud Provider Extension to kubelet
		Wants=kubelet.service
		After=kubelet.service
		[Service]
		Type=simple
		ExecStart=/usr/local/bin/node-cloud-provider.sh
		[Install]
		WantedBy=multi-user.target`}

		systemdUnits = append(systemdUnits, nodeCloudProvider)
	}

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
