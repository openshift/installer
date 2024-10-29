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

package compute

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	corev1 "k8s.io/api/core/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// InstanceSpec defines the fields which can be set on a new OpenStack instance.
type InstanceSpec struct {
	Name                          string
	ImageID                       string
	FlavorID                      string
	SSHKeyName                    string
	UserData                      string
	Metadata                      map[string]string
	ConfigDrive                   bool
	FailureDomain                 string
	RootVolume                    *infrav1.RootVolume
	AdditionalBlockDevices        []infrav1.AdditionalBlockDevice
	ServerGroupID                 string
	Trunk                         bool
	Tags                          []string
	SchedulerAdditionalProperties []infrav1.SchedulerHintAdditionalProperty
}

// InstanceIdentifier describes an instance which has not necessarily been fetched.
type InstanceIdentifier struct {
	ID   string
	Name string
}

// InstanceStatus represents instance data which has been returned by OpenStack.
type InstanceStatus struct {
	server *servers.Server
	logger logr.Logger
}

func NewInstanceStatusFromServer(server *servers.Server, logger logr.Logger) *InstanceStatus {
	return &InstanceStatus{server, logger}
}

type networkInterface struct {
	Address string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type"`
}

// InstanceNetworkStatus represents the network status of an OpenStack instance
// as used by CAPO. Therefore it may use more context than just data which was
// returned by OpenStack.
type InstanceNetworkStatus struct {
	addresses map[string][]corev1.NodeAddress
}

func (is *InstanceStatus) ID() string {
	return is.server.ID
}

func (is *InstanceStatus) Name() string {
	return is.server.Name
}

func (is *InstanceStatus) State() infrav1.InstanceState {
	return infrav1.InstanceState(is.server.Status)
}

func (is *InstanceStatus) SSHKeyName() string {
	return is.server.KeyName
}

func (is *InstanceStatus) AvailabilityZone() string {
	return is.server.AvailabilityZone
}

// BastionStatus updates BastionStatus in openStackCluster.
func (is *InstanceStatus) UpdateBastionStatus(openStackCluster *infrav1.OpenStackCluster) {
	if openStackCluster.Status.Bastion == nil {
		openStackCluster.Status.Bastion = &infrav1.BastionStatus{}
	}

	openStackCluster.Status.Bastion.ID = is.ID()
	openStackCluster.Status.Bastion.Name = is.Name()
	openStackCluster.Status.Bastion.SSHKeyName = is.SSHKeyName()
	openStackCluster.Status.Bastion.State = is.State()

	ns, err := is.NetworkStatus()
	if err != nil {
		// Bastion IP won't be saved in status, error is not critical
		return
	}

	clusterNetwork := openStackCluster.Status.Network.Name
	openStackCluster.Status.Bastion.IP = ns.IP(clusterNetwork)
}

// InstanceIdentifier returns an InstanceIdentifier object for an InstanceStatus.
func (is *InstanceStatus) InstanceIdentifier() *InstanceIdentifier {
	return &InstanceIdentifier{
		ID:   is.ID(),
		Name: is.Name(),
	}
}

// NetworkStatus returns an InstanceNetworkStatus object for an InstanceStatus.
func (is *InstanceStatus) NetworkStatus() (*InstanceNetworkStatus, error) {
	// Gophercloud doesn't give us a struct for server addresses: we get a
	// map of networkname -> interface{}. That interface{} is a list of
	// addresses as in the example output here:
	// https://docs.openstack.org/api-ref/compute/?expanded=show-server-details-detail#show-server-details
	//
	// Here we convert the interface{} into something more usable by
	// marshalling it to json, then unmarshalling it back into our own
	// struct.
	addressesByNetwork := make(map[string][]corev1.NodeAddress)
	for networkName, b := range is.server.Addresses {
		list, err := json.Marshal(b)
		if err != nil {
			return nil, fmt.Errorf("error marshalling addresses for instance %s: %w", is.ID(), err)
		}
		var interfaceList []networkInterface
		err = json.Unmarshal(list, &interfaceList)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling addresses for instance %s: %w", is.ID(), err)
		}

		var IPv4addresses, IPv6addresses []corev1.NodeAddress
		for i := range interfaceList {
			address := &interfaceList[i]

			var addressType corev1.NodeAddressType
			switch address.Type {
			case "floating":
				addressType = corev1.NodeExternalIP
			case "fixed":
				addressType = corev1.NodeInternalIP
			default:
				is.logger.V(6).Info("Ignoring address with unknown type", "address", address.Address, "type", address.Type)
				continue
			}
			if address.Version == 4 {
				IPv4addresses = append(IPv4addresses, corev1.NodeAddress{
					Type:    addressType,
					Address: address.Address,
				})
			} else {
				IPv6addresses = append(IPv6addresses, corev1.NodeAddress{
					Type:    addressType,
					Address: address.Address,
				})
			}
		}
		// Maintain IPv4 addresses being first ones on Machine's status given there are operations, e.g. reconcile load-balancer member, that use the first address by network type
		addressesByNetwork[networkName] = append(IPv4addresses, IPv6addresses...)
	}

	return &InstanceNetworkStatus{addressesByNetwork}, nil
}

// Addresses returns a list of NodeAddresses containing all addresses which will
// be reported on the OpenStackMachine object.
func (ns *InstanceNetworkStatus) Addresses() []corev1.NodeAddress {
	// We want the returned order of addresses to be deterministic to make
	// it easy to detect changes and avoid unnecessary updates. Iteration
	// over maps is non-deterministic, so we explicitly iterate over the
	// address map in lexical order of network names. This order is
	// arbitrary.
	// Pull out addresses map keys (network names) and sort them lexically
	networks := make([]string, 0, len(ns.addresses))
	for network := range ns.addresses {
		networks = append(networks, network)
	}
	sort.Strings(networks)

	var addresses []corev1.NodeAddress
	for _, network := range networks {
		addressList := ns.addresses[network]
		addresses = append(addresses, addressList...)
	}

	return addresses
}

func (ns *InstanceNetworkStatus) firstAddressByNetworkAndType(networkName string, addressType corev1.NodeAddressType) string {
	if addressList, ok := ns.addresses[networkName]; ok {
		for i := range addressList {
			address := &addressList[i]
			if address.Type == addressType {
				return address.Address
			}
		}
	}
	return ""
}

// IP returns the first listed ip of an instance for the given network name.
func (ns *InstanceNetworkStatus) IP(networkName string) string {
	return ns.firstAddressByNetworkAndType(networkName, corev1.NodeInternalIP)
}

// FloatingIP returns the first listed floating ip of an instance for the given
// network name.
func (ns *InstanceNetworkStatus) FloatingIP(networkName string) string {
	return ns.firstAddressByNetworkAndType(networkName, corev1.NodeExternalIP)
}
