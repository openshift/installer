package ovirt

import (
	"fmt"
	"net"
	"strconv"
	"time"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
)

const bootstrapSSHPort = 22

var bootstrapSSHPortAsString = strconv.Itoa(22)

// PlatformStages are the stages to run to provision the infrastructure in oVirt.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		ovirttypes.Name,
		"image",
		[]providers.Provider{providers.OVirt},
		stages.WithNormalBootstrapDestroy(),
	),
	stages.NewStage(
		ovirttypes.Name,
		"cluster",
		[]providers.Provider{providers.OVirt},
		stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
	),
	stages.NewStage(
		ovirttypes.Name,
		"bootstrap",
		[]providers.Provider{providers.OVirt},
		stages.WithNormalBootstrapDestroy(),
		stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
	),
}

func extractOutputHostAddresses(s stages.SplitStage, directory string, ic *types.InstallConfig) (bootstrapIP string, sshPort int, controlPlaneIPs []string, returnErr error) {
	sshPort = bootstrapSSHPort

	outputs, err := stages.GetTerraformOutputs(s, directory)
	if err != nil {
		returnErr = err
		return
	}

	client, err := ovirt.NewConnection()
	if err != nil {
		returnErr = errors.Wrap(err, "failed to initialize connection to ovirt-engine")
		return
	}
	defer client.Close()

	if vmIDRaw, ok := outputs["bootstrap_vm_id"]; ok {
		vmID, ok := vmIDRaw.(string)
		if !ok {
			returnErr = errors.New("could not read bootstrap VM ID from terraform outputs")
			return
		}
		ip, err := findVirtualMachineIP(vmID, client)
		if err != nil {
			returnErr = errors.Wrapf(err, "could not find IP address for bootstrap instance %q", vmID)
			return
		}
		bootstrapIP = ip
	}

	if vmIDsRaw, ok := outputs["control_plane_vm_ids"]; ok {
		vmIDs, ok := vmIDsRaw.([]interface{})
		if !ok {
			returnErr = errors.New("could not read control plane VM IDs from terraform outputs")
			return
		}
		controlPlaneIPs = make([]string, len(vmIDs))
		for i, vmIDRaw := range vmIDs {
			vmID, ok := vmIDRaw.(string)
			if !ok {
				returnErr = errors.New("could not read control plane VM ID from terraform outputs")
				return
			}
			ip, err := findVirtualMachineIP(vmID, client)
			if err != nil {
				returnErr = errors.Wrapf(err, "could not find IP address for bootstrap instance %q", vmID)
				return
			}
			controlPlaneIPs[i] = ip
		}
	}

	return
}

func checkPortIsOpen(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		logrus.Debugf("connection error: %v", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
	}
	return conn != nil
}

func getReportedDevices(c *ovirtsdk4.Connection, vmID string) (*ovirtsdk4.ReportedDeviceSlice, error) {
	vmsService := c.SystemService().VmsService()
	// Look up the vm by id:
	vmResp, err := vmsService.VmService(vmID).Get().Send()
	if err != nil {
		return nil, fmt.Errorf("failed to find VM, by id %v, reason: %v", vmID, err)
	}
	vm := vmResp.MustVm()

	// Get the reported-devices service for this vm:
	reportedDevicesService := vmsService.VmService(vm.MustId()).ReportedDevicesService()

	// Get the guest reported devices
	reportedDeviceResp, err := reportedDevicesService.List().Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get reported devices list, reason: %v", err)
	}
	reportedDeviceSlice, hasIps := reportedDeviceResp.ReportedDevice()

	if !hasIps {
		return nil, fmt.Errorf("cannot find IPs for vmId: %s", vmID)
	}
	return reportedDeviceSlice, nil
}

func findVirtualMachineIP(instanceID string, client *ovirtsdk4.Connection) (string, error) {
	reportedDeviceSlice, err := getReportedDevices(client, instanceID)
	if err == nil {
		return "", errors.Wrapf(err, "could not Find IP Address for vm id: %s", instanceID)
	}

	for _, reportedDevice := range reportedDeviceSlice.Slice() {
		ips, hasIps := reportedDevice.Ips()
		if hasIps {
			for _, ip := range ips.Slice() {
				ipres, hasAddress := ip.Address()
				if hasAddress {
					if checkPortIsOpen(ipres, bootstrapSSHPortAsString) {
						logrus.Debugf("ovirt vm id: %s , found usable IP Address: %s", instanceID, ipres)
						return ipres, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("could not find usable IP address for vm id: %s", instanceID)
}
