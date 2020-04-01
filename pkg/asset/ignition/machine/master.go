package machine

import (
	"encoding/json"
	"github.com/coreos/ignition/config/util"
	"github.com/openshift/installer/pkg/asset/ignition"
        "net"
	"os"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/pkg/errors"
        "github.com/apparentlymart/go-cidr/cidr"
        "github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
        ign "github.com/openshift/installer/pkg/asset/ignition"
)

const (
	masterIgnFilename = "master.ign"
)

// Master is an asset that generates the ignition config for master nodes.
type Master struct {
	Config *igntypes.Config
	File   *asset.File
}

var _ asset.WritableAsset = (*Master)(nil)

// Dependencies returns the assets on which the Master asset depends.
func (a *Master) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.RootCA{},
	}
}

// Generate generates the ignition config for the Master asset.
func (a *Master) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	rootCA := &tls.RootCA{}
	dependencies.Get(installConfig, rootCA)

	a.Config = pointerIgnitionConfig(installConfig.Config, rootCA.Cert(), "master")

        // Create network Script
        machineCIDR := &installConfig.Config.Networking.DeprecatedMachineCIDR.IPNet
        defaultGateway, _ := cidr.Host(machineCIDR, 1)
        kube_api_vlan := installConfig.Config.Platform.OpenStack.AciNetExt.KubeApiVLAN
        infra_vlan := installConfig.Config.Platform.OpenStack.AciNetExt.InfraVLAN
        mtu_value := installConfig.Config.Platform.OpenStack.AciNetExt.Mtu
        networkScriptString, _ := ign.NetworkScript(kube_api_vlan, defaultGateway.String(), mtu_value)

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
               DEVICE=ens3.` + infra_vlan + `
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

        logrus.Info("Editing Master.........")

        a.Config.Storage.Files = append(a.Config.Storage.Files,
                        ignition.FileFromString("/usr/local/bin/kube-api-interface.sh", "root", 0555, string(networkScriptString)),
			            ignition.FileFromString("/etc/sysconfig/network-scripts/ifcfg-ens3.4094", "root", 0420, ifcfg_ens3_string),
                        ignition.FileFromString("/etc/sysconfig/network-scripts/ifcfg-opflex-conn", "root", 0420, ifcfg_opflex_conn_string),
                        ignition.FileFromString("/etc/sysconfig/network-scripts/ifcfg-uplink-conn", "root", 0420, ifcfg_uplink_conn_string),
                        ignition.FileFromString("/etc/sysconfig/network-scripts/route-opflex-conn", "root", 0420, route_opflex_conn_string),
                        ignition.FileFromString("/etc/sysconfig/network-scripts/route-ens3.4094", "root", 0420, route_ens3_string))

	unit := igntypes.Unit{
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
	a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)

	unit = igntypes.Unit{
		Name:    "machine-config-daemon-force.path",
		Enabled: util.BoolToPtr(true),
		Contents: `[Unit]
Description=Path File for Disabling Machine-Config Validation Check
[Path]
PathChanged=/run/machine-config-daemon-force
Unit=machine-config-daemon-force.service
[Install]
WantedBy=multi-user.target`}

	a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)

	unit = igntypes.Unit{
		Name:    "machine-config-daemon-force.service",
		Enabled: util.BoolToPtr(true),
		Contents: `[Unit]
Description=Disabling Machine-Config Validation Check
[Service]
Type=simple
ExecStart=touch /run/machine-config-daemon-force
[Install]
WantedBy=multi-user.target`}

	a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)

	data, err := json.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal Ignition config")
	}
	a.File = &asset.File{
		Filename: masterIgnFilename,
		Data:     data,
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Master) Name() string {
	return "Master Ignition Config"
}

// Files returns the files generated by the asset.
func (a *Master) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Load returns the master ignitions from disk.
func (a *Master) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(masterIgnFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", masterIgnFilename)
	}

	a.File, a.Config = file, config
	return true, nil
}
