// Package ovirt supply utilities to extract information from terraform state
package ovirt

import (
	"fmt"
	"net"
	"time"

	"github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/terraform"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const bootstrapSSHPort = "22"

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
	reportedDeviceSlice, _ := reportedDeviceResp.ReportedDevice()

	if len(reportedDeviceSlice.Slice()) == 0 {
		return nil, fmt.Errorf("cannot find IPs for vmId: %s", vmID)
	}
	return reportedDeviceSlice, nil
}

func findVirtualMachineIP(c *ovirtsdk4.Connection, moRefValue string) (string, error) {

	reportedDeviceSlice, err := getReportedDevices(c, moRefValue)
	if err != nil {
		return "", err
	}

	for _, reportedDevice := range reportedDeviceSlice.Slice() {
		ips, hasIps := reportedDevice.Ips()
		if hasIps {
			for _, ip := range ips.Slice() {
				ipres, hasAddress := ip.Address()
				if hasAddress {
					if checkPortIsOpen(ipres, bootstrapSSHPort) {
						logrus.Debugf("ovirt vm id: %s , found usable bootstrap IP %s", moRefValue, ipres)
						return ipres, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("could not find usable bootstrap IP address for vm id: %s", moRefValue)
}

// BootstrapIP returns the ip address for bootstrap host.
// still unsupported, because qemu-ga is not available - see https://bugzilla.redhat.com/show_bug.cgi?id=1764804
func BootstrapIP(tfs *terraform.State) (string, error) {

	client, err := ovirt.NewConnection()
	if err != nil {
		return "", fmt.Errorf("failed to initialize connection to ovirt-engine's %s", err)
	}
	defer client.Close()

	br, err := terraform.LookupResource(tfs, "module.bootstrap", "ovirt_vm", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}

	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}

	vmid, found, err := unstructured.NestedString(br.Instances[0].Attributes, "id")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap managed object reference")
	}
	if !found {
		return "", errors.Errorf("failed to lookup bootstrap by vmID: %s", vmid)
	}

	ip, err := findVirtualMachineIP(client, vmid)
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap ipv4 address")
	}
	return ip, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
// still unsupported, because qemu-ga is not available  - see https://bugzilla.redhat.com/show_bug.cgi?id=1764804
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	return []string{""}, nil
}
