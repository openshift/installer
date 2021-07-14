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
	reportedDeviceSlice, hasIps := reportedDeviceResp.ReportedDevice()

	if !hasIps {
		return nil, fmt.Errorf("cannot find IPs for vmId: %s", vmID)
	}
	return reportedDeviceSlice, nil
}

func findVirtualMachineIP(moRefValue string, client *ovirtsdk4.Connection) (string, error) {
	reportedDeviceSlice, err := getReportedDevices(client, moRefValue)
	if err == nil {
		return "", errors.Wrapf(err, "couldnt Find IP Address for vm id: %s", moRefValue)
	}

	for _, reportedDevice := range reportedDeviceSlice.Slice() {
		ips, hasIps := reportedDevice.Ips()
		if hasIps {
			for _, ip := range ips.Slice() {
				ipres, hasAddress := ip.Address()
				if hasAddress {
					if checkPortIsOpen(ipres, bootstrapSSHPort) {
						logrus.Debugf("ovirt vm id: %s , found usable  IP Address: %s", moRefValue, ipres)
						return ipres, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("could not find usable IP address for vm id: %s", moRefValue)
}

func lookupVMResources(tfs *terraform.State, moduleName string, name string) ([]string, error) {

	client, err := ovirt.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection to ovirt-engine's %s", err)
	}
	defer client.Close()

	br, err := terraform.LookupResource(tfs, moduleName, "ovirt_vm", name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to lookup %s VM", name)
	}

	if len(br.Instances) == 0 {
		return nil, errors.New(fmt.Sprintf("no %s instance found", name))
	}

	var ips []string

	for _, instance := range br.Instances {
		vmid, found, err := unstructured.NestedString(instance.Attributes, "id")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to lookup %s managed object reference", name)
		}
		if !found {
			return nil, errors.Errorf("failed to lookup %s by vmID: %s", name, vmid)
		}
		ip, err := findVirtualMachineIP(vmid, client)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to lookup %s ipv4 address", name)
		}
		ips = append(ips, ip)
	}

	return ips, nil
}

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	ips, err := lookupVMResources(tfs, "module.bootstrap", "bootstrap")
	if len(ips) == 0 {
		return "", errors.Wrapf(err, "no ips Found for Bootstrap VM")
	}

	return ips[0], err
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	return lookupVMResources(tfs, "module.masters", "master")
}
