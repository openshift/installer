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

package net

import (
	"context"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// NetworkStatus provides information about one of a VM's networks.
type NetworkStatus struct {
	// Connected is a flag that indicates whether this network is currently
	// connected to the VM.
	Connected bool `json:"connected,omitempty"`

	// IPAddrs is one or more IP addresses reported by vm-tools.
	// +optional
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// MACAddr is the MAC address of the network device.
	MACAddr string `json:"macAddr"`

	// NetworkName is the name of the network.
	// +optional
	NetworkName string `json:"networkName,omitempty"`
}

// GetNetworkStatus returns the network information for the specified VM.
func GetNetworkStatus(
	ctx context.Context,
	client *vim25.Client,
	moRef types.ManagedObjectReference) ([]NetworkStatus, error) {
	var (
		obj mo.VirtualMachine

		pc    = property.DefaultCollector(client)
		props = []string{
			"config.hardware.device",
			"guest.net",
		}
	)

	if err := pc.RetrieveOne(ctx, moRef, props, &obj); err != nil {
		return nil, errors.Wrapf(err, "unable to fetch props %v for vm %v", props, moRef)
	}
	if obj.Config == nil {
		return nil, errors.New("config.hardware.device is nil")
	}

	var allNetStatus []NetworkStatus

	for _, device := range obj.Config.Hardware.Device {
		if dev, ok := device.(types.BaseVirtualEthernetCard); ok {
			nic := dev.GetVirtualEthernetCard()
			netStatus := NetworkStatus{
				MACAddr: nic.MacAddress,
			}
			if obj.Guest != nil {
				for _, i := range obj.Guest.Net {
					if strings.EqualFold(nic.MacAddress, i.MacAddress) {
						netStatus.IPAddrs = i.IpAddress
						netStatus.NetworkName = i.Network
						netStatus.Connected = i.Connected
					}
				}
			}
			allNetStatus = append(allNetStatus, netStatus)
		}
	}

	return allNetStatus, nil
}

// ErrOnLocalOnlyIPAddr returns an error if the provided IP address is
// accessible only on the VM's guest OS.
func ErrOnLocalOnlyIPAddr(addr string) error {
	var reason string
	a := net.ParseIP(addr)
	switch {
	case len(a) == 0:
		reason = "invalid"
	case a.IsUnspecified():
		reason = "unspecified"
	case a.IsLinkLocalMulticast():
		reason = "link-local-mutlicast"
	case a.IsLinkLocalUnicast():
		reason = "link-local-unicast"
	case a.IsLoopback():
		reason = "loopback"
	}
	if reason != "" {
		return errors.Errorf("failed to validate ip addr=%v: %s", addr, reason)
	}
	return nil
}
