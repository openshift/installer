/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"regexp"
	"text/template"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/integer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
)

// GetVSphereMachine gets a vmware.infrastructure.cluster.x-k8s.io.VSphereMachine resource for the given CAPI Machine.
func GetVSphereMachine(
	ctx context.Context,
	controllerClient client.Client,
	namespace, machineName string) (*vmwarev1.VSphereMachine, error) {
	machine := &vmwarev1.VSphereMachine{}
	namespacedName := apitypes.NamespacedName{
		Namespace: namespace,
		Name:      machineName,
	}
	if err := controllerClient.Get(ctx, namespacedName, machine); err != nil {
		return nil, err
	}
	return machine, nil
}

// ErrNoMachineIPAddr indicates that no valid IP addresses were found in a machine context.
var ErrNoMachineIPAddr = errors.New("no IP addresses found for machine")

// GetMachinePreferredIPAddress returns the preferred IP address for a
// VSphereMachine resource.
func GetMachinePreferredIPAddress(machine *infrav1.VSphereMachine) (string, error) {
	var cidr *net.IPNet
	if cidrString := machine.Spec.Network.PreferredAPIServerCIDR; cidrString != "" {
		var err error
		if _, cidr, err = net.ParseCIDR(cidrString); err != nil {
			return "", errors.New("error parsing preferred API server CIDR")
		}
	}

	for _, machineAddr := range machine.Status.Addresses {
		if machineAddr.Type != clusterv1.MachineExternalIP {
			continue
		}
		if cidr == nil {
			return machineAddr.Address, nil
		}
		if cidr.Contains(net.ParseIP(machineAddr.Address)) {
			return machineAddr.Address, nil
		}
	}

	return "", ErrNoMachineIPAddr
}

// IsControlPlaneMachine returns true if the provided resource is
// a member of the control plane.
func IsControlPlaneMachine(machine metav1.Object) bool {
	_, ok := machine.GetLabels()[clusterv1.MachineControlPlaneLabel]
	return ok
}

// GetMachineMetadata the cloud-init metadata as a base-64 encoded
// string for a given VSphereMachine.
// IPAM state includes IP and Gateways that should be added to each device.
func GetMachineMetadata(hostname string, vsphereVM infrav1.VSphereVM, ipamState map[string]infrav1.NetworkDeviceSpec, networkStatuses ...infrav1.NetworkStatus) ([]byte, error) {
	// Create a copy of the devices and add their MAC addresses from a network status.
	devices := make([]infrav1.NetworkDeviceSpec, integer.IntMax(len(vsphereVM.Spec.Network.Devices), len(networkStatuses)))

	var waitForIPv4, waitForIPv6 bool
	for i := range vsphereVM.Spec.Network.Devices {
		vsphereVM.Spec.Network.Devices[i].DeepCopyInto(&devices[i])

		// Add the MAC Address to the network device
		if len(networkStatuses) > i {
			devices[i].MACAddr = networkStatuses[i].MACAddr
		}

		if state, ok := ipamState[devices[i].MACAddr]; ok {
			devices[i].IPAddrs = append(devices[i].IPAddrs, state.IPAddrs...)
			devices[i].Gateway4 = state.Gateway4
			devices[i].Gateway6 = state.Gateway6
		}

		if waitForIPv4 && waitForIPv6 {
			// break early as we already wait for ipv4 and ipv6
			continue
		}
		// check static IPs
		for _, ipStr := range vsphereVM.Spec.Network.Devices[i].IPAddrs {
			ip := net.ParseIP(ipStr)
			// check the IP family
			if ip != nil {
				if ip.To4() == nil {
					waitForIPv6 = true
				} else {
					waitForIPv4 = true
				}
			}
		}
		// check if DHCP is enabled
		if vsphereVM.Spec.Network.Devices[i].DHCP4 {
			waitForIPv4 = true
		}
		if vsphereVM.Spec.Network.Devices[i].DHCP6 {
			waitForIPv6 = true
		}
	}

	// Add the MAC Address to the network device
	// networkStatuses may be longer than devices
	// and we want to add all the networks
	for i, status := range networkStatuses {
		devices[i].MACAddr = status.MACAddr
	}

	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("t").Funcs(
		template.FuncMap{
			"nameservers": func(spec infrav1.NetworkDeviceSpec) bool {
				return len(spec.Nameservers) > 0 || len(spec.SearchDomains) > 0
			},
		}).Parse(metadataFormat))
	if err := tpl.Execute(buf, struct {
		Hostname    string
		Devices     []infrav1.NetworkDeviceSpec
		Routes      []infrav1.NetworkRouteSpec
		WaitForIPv4 bool
		WaitForIPv6 bool
	}{
		Hostname:    hostname, // note that hostname determines the Kubernetes node name
		Devices:     devices,
		Routes:      vsphereVM.Spec.Network.Routes,
		WaitForIPv4: waitForIPv4,
		WaitForIPv6: waitForIPv6,
	}); err != nil {
		return nil, errors.Wrapf(
			err,
			"error getting cloud init metadata for vsphereVM %s/%s",
			vsphereVM.Namespace, vsphereVM.Name)
	}
	return buf.Bytes(), nil
}

// GetOwnerVSphereMachine returns the VSphereMachine owner for the passed object.
func GetOwnerVSphereMachine(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*infrav1.VSphereMachine, error) {
	for _, ref := range obj.OwnerReferences {
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, err
		}
		if ref.Kind == "VSphereMachine" && gv.Group == infrav1.GroupVersion.Group {
			return getVSphereMachineByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

func getVSphereMachineByName(ctx context.Context, c client.Client, namespace, name string) (*infrav1.VSphereMachine, error) {
	m := &infrav1.VSphereMachine{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}

const (
	// ProviderIDPrefix is the string data prefixed to a BIOS UUID in order
	// to build a provider ID.
	ProviderIDPrefix = "vsphere://"

	// ProviderIDPattern is a regex pattern and is used by ConvertProviderIDToUUID
	// to convert a providerID into a UUID string.
	ProviderIDPattern = `(?i)^` + ProviderIDPrefix + `([a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12})$`

	// UUIDPattern is a regex pattern and is used by ConvertUUIDToProviderID
	// to convert a UUID into a providerID string.
	UUIDPattern = `(?i)^[a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12}$`
)

// ConvertProviderIDToUUID transforms a provider ID into a UUID string.
// If providerID is nil, empty, or invalid, then an empty string is returned.
// A valid providerID should adhere to the format specified by
// ProviderIDPattern.
func ConvertProviderIDToUUID(providerID *string) string {
	if providerID == nil || *providerID == "" {
		return ""
	}
	pattern := regexp.MustCompile(ProviderIDPattern)
	matches := pattern.FindStringSubmatch(*providerID)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

// ConvertUUIDToProviderID transforms a UUID string into a provider ID.
// If the supplied UUID is empty or invalid then an empty string is returned.
// A valid UUID should adhere to the format specified by UUIDPattern.
func ConvertUUIDToProviderID(uuid string) string {
	if uuid == "" {
		return ""
	}
	pattern := regexp.MustCompile(UUIDPattern)
	if !pattern.MatchString(uuid) {
		return ""
	}
	return ProviderIDPrefix + uuid
}

// MachinesAsString constructs a string (with correct punctuations) to be
// used in logging and error messages.
func MachinesAsString(machines []*clusterv1.Machine) string {
	var message string
	count := 1
	for _, m := range machines {
		if count == 1 {
			message = fmt.Sprintf("%s/%s", m.Namespace, m.Name)
		} else {
			var format string
			if count > 1 && count != len(machines) {
				format = "%s, %s/%s"
			} else if count == len(machines) {
				format = "%s and %s/%s"
			}
			message = fmt.Sprintf(format, message, m.Namespace, m.Name)
		}
		count++
	}
	return message
}
