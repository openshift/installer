/*
Copyright 2021 The Kubernetes Authors.

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

package networking

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
)

const (
	timeoutPortDelete       = 3 * time.Minute
	retryIntervalPortDelete = 5 * time.Second
)

// GetPortFromInstanceIP returns at most one port attached to the instance with given ID
// and with the IP address provided.
func (s *Service) GetPortFromInstanceIP(instanceID string, ip string) ([]ports.Port, error) {
	portOpts := ports.ListOpts{
		DeviceID: instanceID,
		FixedIPs: []ports.FixedIPOpts{
			{
				IPAddress: ip,
			},
		},
		Limit: 1,
	}
	return s.client.ListPort(portOpts)
}

type PortListOpts struct {
	DeviceOwner []string `q:"device_owner"`
	NetworkID   string   `q:"network_id"`
}

func (p *PortListOpts) ToPortListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(p)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func (s *Service) GetPortForExternalNetwork(instanceID string, externalNetworkID string) (*ports.Port, error) {
	instancePortsOpts := ports.ListOpts{
		DeviceID: instanceID,
	}
	instancePorts, err := s.client.ListPort(instancePortsOpts)
	if err != nil {
		return nil, fmt.Errorf("lookup ports for server %s: %w", instanceID, err)
	}

	for _, instancePort := range instancePorts {
		networkPortsOpts := &PortListOpts{
			NetworkID:   instancePort.NetworkID,
			DeviceOwner: []string{"network:router_interface", "network:router_interface_distributed", "network:ha_router_replicated_interface", "network:router_ha_interface"},
		}

		networkPorts, err := s.client.ListPort(networkPortsOpts)
		if err != nil {
			return nil, fmt.Errorf("lookup ports for network %s: %w", instancePort.NetworkID, err)
		}

		for _, networkPort := range networkPorts {
			// Check if the instance port and the network port share a subnet
			matchingSubnet := false
			for _, fixedIP := range instancePort.FixedIPs {
				for _, networkFixedIP := range networkPort.FixedIPs {
					if fixedIP.SubnetID == networkFixedIP.SubnetID {
						matchingSubnet = true
						break
					}
				}
				if matchingSubnet {
					break
				}
			}
			if !matchingSubnet {
				continue
			}

			router, err := s.client.GetRouter(networkPort.DeviceID)
			if err != nil {
				return nil, fmt.Errorf("lookup router %s: %w", networkPort.DeviceID, err)
			}

			if router.GatewayInfo.NetworkID == externalNetworkID {
				return &instancePort, nil
			}
		}
	}
	return nil, nil
}

// DeletePort deletes the Neutron port with the given ID.
func (s *Service) DeletePort(eventObject runtime.Object, portID string) error {
	var err error
	err = wait.PollUntilContextTimeout(context.TODO(), retryIntervalPortDelete, timeoutPortDelete, true, func(_ context.Context) (bool, error) {
		err = s.client.DeletePort(portID)
		if err != nil {
			if capoerrors.IsNotFound(err) {
				record.Eventf(eventObject, "SuccessfulDeletePort", "Port with id %d did not exist", portID)
				// this is success so we return without another try
				return true, nil
			}
			if capoerrors.IsRetryable(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	if err != nil {
		record.Warnf(eventObject, "FailedDeletePort", "Failed to delete port with id %s: %v", portID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDeletePort", "Deleted port with id %s", portID)
	return nil
}

// DeleteClusterPorts deletes all ports created for the cluster.
func (s *Service) DeleteClusterPorts(openStackCluster *infrav1.OpenStackCluster) error {
	// If the network is not ready, do nothing
	if openStackCluster.Status.Network == nil || openStackCluster.Status.Network.ID == "" {
		return nil
	}
	networkID := openStackCluster.Status.Network.ID

	portList, err := s.client.ListPort(ports.ListOpts{
		NetworkID:   networkID,
		DeviceOwner: "",
	})
	s.scope.Logger().Info("Deleting cluster ports", "networkID", networkID, "portList", portList)
	if err != nil {
		if capoerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("list ports of network %q: %v", networkID, err)
	}

	for _, port := range portList {
		// The bastion port in 0.10 was prefixed with the namespace and then the current port name
		// so in order to cleanup the old bastion port we need to check for the old format.
		bastionLegacyPortPrefix := fmt.Sprintf("%s-%s", openStackCluster.Namespace, openStackCluster.Name)
		if strings.HasPrefix(port.Name, openStackCluster.Name) || strings.HasPrefix(port.Name, bastionLegacyPortPrefix) {
			if err := s.DeletePort(openStackCluster, port.ID); err != nil {
				return fmt.Errorf("error deleting port %s: %v", port.ID, err)
			}
		}
	}

	return nil
}

// UpdateAllowedAddressPairs updates the allowedAddressPairs on an existing Neutron port.
func (s *Service) UpdateAllowedAddressPairs(portID string, pairs []infrav1.AddressPair) error {
	addressPairs := make([]ports.AddressPair, len(pairs))
	for i, ap := range pairs {
		addressPairs[i] = ports.AddressPair{
			IPAddress:  ap.IPAddress,
			MACAddress: ptr.Deref(ap.MACAddress, ""),
		}
	}
	_, err := s.client.UpdatePort(portID, ports.UpdateOpts{
		AllowedAddressPairs: &addressPairs,
	})
	return err
}
