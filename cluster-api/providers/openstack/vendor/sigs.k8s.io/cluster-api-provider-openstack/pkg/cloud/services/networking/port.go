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

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
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

func (s *Service) GetOrCreatePort(eventObject runtime.Object, clusterName string, portName string, portOpts *infrav1.PortOpts, instanceSecurityGroups []string, instanceTags []string) (*ports.Port, error) {
	networkID := portOpts.Network.ID

	existingPorts, err := s.client.ListPort(ports.ListOpts{
		Name:      portName,
		NetworkID: networkID,
	})
	if err != nil {
		return nil, fmt.Errorf("searching for existing port for server: %v", err)
	}

	if len(existingPorts) == 1 {
		return &existingPorts[0], nil
	}

	if len(existingPorts) > 1 {
		return nil, fmt.Errorf("multiple ports found with name \"%s\"", portName)
	}

	description := portOpts.Description
	if description == "" {
		description = names.GetDescription(clusterName)
	}

	var securityGroups []string
	addressPairs := []ports.AddressPair{}
	if portOpts.DisablePortSecurity == nil || !*portOpts.DisablePortSecurity {
		for _, ap := range portOpts.AllowedAddressPairs {
			addressPairs = append(addressPairs, ports.AddressPair{
				IPAddress:  ap.IPAddress,
				MACAddress: ap.MACAddress,
			})
		}
		if portOpts.SecurityGroupFilters != nil {
			securityGroups, err = s.GetSecurityGroups(portOpts.SecurityGroupFilters)
			if err != nil {
				return nil, fmt.Errorf("error getting security groups: %v", err)
			}
		}
		// inherit port security groups from the instance if not explicitly specified
		if len(securityGroups) == 0 {
			securityGroups = instanceSecurityGroups
		}
	}

	var fixedIPs interface{}
	if len(portOpts.FixedIPs) > 0 {
		fips := make([]ports.IP, 0, len(portOpts.FixedIPs)+1)
		for _, fixedIP := range portOpts.FixedIPs {
			subnetID, err := s.getSubnetIDForFixedIP(fixedIP.Subnet, networkID)
			if err != nil {
				return nil, err
			}
			fips = append(fips, ports.IP{
				SubnetID:  subnetID,
				IPAddress: fixedIP.IPAddress,
			})
		}
		fixedIPs = fips
	}

	var valueSpecs *map[string]string
	if len(portOpts.ValueSpecs) > 0 {
		vs := make(map[string]string, len(portOpts.ValueSpecs))
		for _, valueSpec := range portOpts.ValueSpecs {
			vs[valueSpec.Key] = valueSpec.Value
		}
		valueSpecs = &vs
	}

	var createOpts ports.CreateOptsBuilder

	// Gophercloud expects a *[]string. We translate a nil slice to a nil pointer.
	var securityGroupsPtr *[]string
	if securityGroups != nil {
		securityGroupsPtr = &securityGroups
	}

	createOpts = ports.CreateOpts{
		Name:                  portName,
		NetworkID:             networkID,
		Description:           description,
		AdminStateUp:          portOpts.AdminStateUp,
		MACAddress:            portOpts.MACAddress,
		SecurityGroups:        securityGroupsPtr,
		AllowedAddressPairs:   addressPairs,
		FixedIPs:              fixedIPs,
		ValueSpecs:            valueSpecs,
		PropagateUplinkStatus: portOpts.PropagateUplinkStatus,
	}

	if portOpts.DisablePortSecurity != nil {
		portSecurity := !*portOpts.DisablePortSecurity
		createOpts = portsecurity.PortCreateOptsExt{
			CreateOptsBuilder:   createOpts,
			PortSecurityEnabled: &portSecurity,
		}
	}

	createOpts = portsbinding.CreateOptsExt{
		CreateOptsBuilder: createOpts,
		HostID:            portOpts.HostID,
		VNICType:          portOpts.VNICType,
		Profile:           getPortProfile(portOpts.Profile),
	}

	port, err := s.client.CreatePort(createOpts)
	if err != nil {
		record.Warnf(eventObject, "FailedCreatePort", "Failed to create port %s: %v", portName, err)
		return nil, err
	}

	var tags []string
	tags = append(tags, instanceTags...)
	tags = append(tags, portOpts.Tags...)
	if len(tags) > 0 {
		if err = s.replaceAllAttributesTags(eventObject, portResource, port.ID, tags); err != nil {
			record.Warnf(eventObject, "FailedReplaceTags", "Failed to replace port tags %s: %v", portName, err)
			return nil, err
		}
	}
	record.Eventf(eventObject, "SuccessfulCreatePort", "Created port %s with id %s", port.Name, port.ID)
	if portOpts.Trunk != nil && *portOpts.Trunk {
		trunk, err := s.getOrCreateTrunk(eventObject, clusterName, port.Name, port.ID)
		if err != nil {
			record.Warnf(eventObject, "FailedCreateTrunk", "Failed to create trunk for port %s: %v", portName, err)
			return nil, err
		}
		if err = s.replaceAllAttributesTags(eventObject, trunkResource, trunk.ID, tags); err != nil {
			record.Warnf(eventObject, "FailedReplaceTags", "Failed to replace trunk tags %s: %v", portName, err)
			return nil, err
		}
	}

	return port, nil
}

func (s *Service) getSubnetIDForFixedIP(subnet *infrav1.SubnetFilter, networkID string) (string, error) {
	if subnet == nil {
		return "", nil
	}
	// Do not query for subnets if UUID is already provided
	if subnet.ID != "" {
		return subnet.ID, nil
	}

	opts := subnet.ToListOpt()
	opts.NetworkID = networkID
	subnets, err := s.client.ListSubnet(opts)
	if err != nil {
		return "", err
	}

	switch len(subnets) {
	case 0:
		return "", fmt.Errorf("subnet query %v, returns no subnets", *subnet)
	case 1:
		return subnets[0].ID, nil
	default:
		return "", fmt.Errorf("subnet query %v, returns too many subnets: %v", *subnet, subnets)
	}
}

func getPortProfile(p infrav1.BindingProfile) map[string]interface{} {
	portProfile := make(map[string]interface{})

	// if p.OVSHWOffload is true, we need to set the profile
	// to enable hardware offload for the port
	if p.OVSHWOffload {
		portProfile["capabilities"] = []string{"switchdev"}
	}
	if p.TrustedVF {
		portProfile["trusted"] = true
	}

	// We need return nil if there is no profiles
	// to have backward compatible defaults.
	// To set profiles, your tenant needs this permission:
	// rule:create_port and rule:create_port:binding:profile
	if len(portProfile) == 0 {
		return nil
	}
	return portProfile
}

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

func (s *Service) DeletePorts(openStackCluster *infrav1.OpenStackCluster) error {
	// If the network is not ready, do nothing
	if openStackCluster.Status.Network == nil || openStackCluster.Status.Network.ID == "" {
		return nil
	}
	networkID := openStackCluster.Status.Network.ID

	portList, err := s.client.ListPort(ports.ListOpts{
		NetworkID:   networkID,
		DeviceOwner: "",
	})
	if err != nil {
		if capoerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("list ports of network %q: %v", networkID, err)
	}

	for _, port := range portList {
		if strings.HasPrefix(port.Name, openStackCluster.Name) {
			err := s.DeletePort(openStackCluster, port.ID)
			if err != nil {
				return fmt.Errorf("delete port %s of network %q failed : %v", port.ID, networkID, err)
			}
		}
	}

	return nil
}

func (s *Service) GarbageCollectErrorInstancesPort(eventObject runtime.Object, instanceName string, portOpts []infrav1.PortOpts) error {
	for i := range portOpts {
		portOpt := &portOpts[i]

		portName := GetPortName(instanceName, portOpt, i)

		// TODO: whould be nice if gophercloud could be persuaded to accept multiple
		// names as is allowed by the API in order to reduce API traffic.
		portList, err := s.client.ListPort(ports.ListOpts{Name: portName})
		if err != nil {
			return err
		}

		// NOTE: https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/1476
		// It is up to the end user to specify a UNIQUE cluster name when provisioning in the
		// same project, otherwise things will alias and we could delete more than we should.
		if len(portList) > 1 {
			return fmt.Errorf("garbage collection of port %s failed, found %d ports with the same name", portName, len(portList))
		}

		if err := s.DeletePort(eventObject, portList[0].ID); err != nil {
			return err
		}
	}

	return nil
}

// GetPortName appends a suffix to an instance name in order to try and get a unique name per port.
func GetPortName(instanceName string, opts *infrav1.PortOpts, netIndex int) string {
	if opts != nil && opts.NameSuffix != "" {
		return fmt.Sprintf("%s-%s", instanceName, opts.NameSuffix)
	}
	return fmt.Sprintf("%s-%d", instanceName, netIndex)
}
