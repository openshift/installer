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
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
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

func (s *Service) CreatePort(eventObject runtime.Object, portSpec *infrav1.ResolvedPortSpec) (*ports.Port, error) {
	var addressPairs []ports.AddressPair
	if !ptr.Deref(portSpec.DisablePortSecurity, false) {
		for _, ap := range portSpec.AllowedAddressPairs {
			addressPairs = append(addressPairs, ports.AddressPair{
				IPAddress:  ap.IPAddress,
				MACAddress: ptr.Deref(ap.MACAddress, ""),
			})
		}
	}

	var fixedIPs []ports.IP
	if len(portSpec.FixedIPs) > 0 {
		fixedIPs = make([]ports.IP, len(portSpec.FixedIPs))
		for i, fixedIP := range portSpec.FixedIPs {
			fixedIPs[i] = ports.IP{
				SubnetID:  ptr.Deref(fixedIP.SubnetID, ""),
				IPAddress: ptr.Deref(fixedIP.IPAddress, ""),
			}
		}
	}

	var valueSpecs *map[string]string
	if len(portSpec.ValueSpecs) > 0 {
		vs := make(map[string]string, len(portSpec.ValueSpecs))
		for _, valueSpec := range portSpec.ValueSpecs {
			vs[valueSpec.Key] = valueSpec.Value
		}
		valueSpecs = &vs
	}

	var builder ports.CreateOptsBuilder
	createOpts := ports.CreateOpts{
		Name:                  portSpec.Name,
		NetworkID:             portSpec.NetworkID,
		Description:           portSpec.Description,
		AdminStateUp:          portSpec.AdminStateUp,
		MACAddress:            ptr.Deref(portSpec.MACAddress, ""),
		AllowedAddressPairs:   addressPairs,
		ValueSpecs:            valueSpecs,
		PropagateUplinkStatus: portSpec.PropagateUplinkStatus,
	}
	if fixedIPs != nil {
		createOpts.FixedIPs = fixedIPs
	}
	if portSpec.SecurityGroups != nil {
		createOpts.SecurityGroups = &portSpec.SecurityGroups
	}
	builder = createOpts

	if portSpec.DisablePortSecurity != nil {
		portSecurity := !*portSpec.DisablePortSecurity
		portSecurityOpts := portsecurity.PortCreateOptsExt{
			CreateOptsBuilder:   builder,
			PortSecurityEnabled: &portSecurity,
		}
		builder = portSecurityOpts
	}

	portsBindingOpts := portsbinding.CreateOptsExt{
		CreateOptsBuilder: builder,
		HostID:            ptr.Deref(portSpec.HostID, ""),
		VNICType:          ptr.Deref(portSpec.VNICType, ""),
		Profile:           getPortProfile(portSpec.Profile),
	}
	builder = portsBindingOpts

	port, err := s.client.CreatePort(builder)
	if err != nil {
		record.Warnf(eventObject, "FailedCreatePort", "Failed to create port %s: %v", portSpec.Name, err)
		return nil, err
	}

	if len(portSpec.Tags) > 0 {
		if err = s.replaceAllAttributesTags(eventObject, portResource, port.ID, portSpec.Tags); err != nil {
			record.Warnf(eventObject, "FailedReplaceTags", "Failed to replace port tags %s: %v", portSpec.Name, err)
			return nil, err
		}
	}
	record.Eventf(eventObject, "SuccessfulCreatePort", "Created port %s with id %s", port.Name, port.ID)
	if ptr.Deref(portSpec.Trunk, false) {
		trunk, err := s.getOrCreateTrunkForPort(eventObject, port)
		if err != nil {
			record.Warnf(eventObject, "FailedCreateTrunk", "Failed to create trunk for port %s: %v", port.Name, err)
			return nil, err
		}
		if err = s.replaceAllAttributesTags(eventObject, trunkResource, trunk.ID, portSpec.Tags); err != nil {
			record.Warnf(eventObject, "FailedReplaceTags", "Failed to replace trunk tags %s: %v", port.Name, err)
			return nil, err
		}
	}

	return port, nil
}

func getPortProfile(p *infrav1.BindingProfile) map[string]interface{} {
	if p == nil {
		return nil
	}

	portProfile := make(map[string]interface{})

	// if p.OVSHWOffload is true, we need to set the profile
	// to enable hardware offload for the port
	if ptr.Deref(p.OVSHWOffload, false) {
		portProfile["capabilities"] = []string{"switchdev"}
	}
	if ptr.Deref(p.TrustedVF, false) {
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

// DeleteTrunk deletes the Neutron trunk and port with the given ID.
func (s *Service) DeleteInstanceTrunkAndPort(eventObject runtime.Object, port infrav1.PortStatus, trunkSupported bool) error {
	if trunkSupported {
		if err := s.DeleteTrunk(eventObject, port.ID); err != nil {
			return fmt.Errorf("error deleting trunk of port %s: %v", port.ID, err)
		}
	}
	if err := s.DeletePort(eventObject, port.ID); err != nil {
		return fmt.Errorf("error deleting port %s: %v", port.ID, err)
	}

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
		if strings.HasPrefix(port.Name, openStackCluster.Name) {
			if err := s.DeletePort(openStackCluster, port.ID); err != nil {
				return fmt.Errorf("error deleting port %s: %v", port.ID, err)
			}
		}
	}

	return nil
}

// getPortName appends a suffix to an instance name in order to try and get a unique name per port.
func getPortName(baseName string, portSpec *infrav1.PortOpts, netIndex int) string {
	if portSpec != nil && portSpec.NameSuffix != nil {
		return fmt.Sprintf("%s-%s", baseName, *portSpec.NameSuffix)
	}
	return fmt.Sprintf("%s-%d", baseName, netIndex)
}

func (s *Service) CreatePorts(eventObject runtime.Object, desiredPorts []infrav1.ResolvedPortSpec, resources *infrav1.MachineResources) error {
	for i := range desiredPorts {
		// Skip creation of ports which already exist
		if i < len(resources.Ports) {
			continue
		}

		portSpec := &desiredPorts[i]
		// Events are recorded in CreatePort
		port, err := s.CreatePort(eventObject, portSpec)
		if err != nil {
			return err
		}

		resources.Ports = append(resources.Ports, infrav1.PortStatus{
			ID: port.ID,
		})
	}

	return nil
}

// ConstructPorts builds an array of ports from the machine spec.
// If no ports are in the spec, returns a single port for a network connection to the default cluster network.
func (s *Service) ConstructPorts(spec *infrav1.OpenStackMachineSpec, clusterResourceName, baseName string, defaultNetwork *infrav1.NetworkStatusWithSubnets, managedSecurityGroup *string, baseTags []string) ([]infrav1.ResolvedPortSpec, error) {
	ports := spec.Ports

	defaultSecurityGroupIDs, err := s.GetSecurityGroups(spec.SecurityGroups)
	if err != nil {
		return nil, fmt.Errorf("error getting security groups: %v", err)
	}
	if managedSecurityGroup != nil {
		defaultSecurityGroupIDs = append(defaultSecurityGroupIDs, *managedSecurityGroup)
	}

	// Ensure user-specified ports have all required fields
	resolvedPorts, err := s.normalizePorts(ports, clusterResourceName, baseName, spec.Trunk, defaultSecurityGroupIDs, defaultNetwork, baseTags)
	if err != nil {
		return nil, err
	}

	// no networks or ports found in the spec, so create a port on the cluster network
	if len(resolvedPorts) == 0 {
		resolvedPorts = make([]infrav1.ResolvedPortSpec, 1)
		resolvedPort := &resolvedPorts[0]
		resolvedPort.Name = getPortName(baseName, nil, 0)
		resolvedPort.Description = names.GetDescription(clusterResourceName)
		if len(baseTags) > 0 {
			resolvedPort.Tags = baseTags
		}
		if spec.Trunk {
			resolvedPort.Trunk = &spec.Trunk
		}
		resolvedPort.SecurityGroups = defaultSecurityGroupIDs
		resolvedPort.NetworkID, resolvedPort.FixedIPs, _ = defaultNetworkTarget(defaultNetwork)
	}

	// trunk support is required if any port has trunk enabled
	portUsesTrunk := func() bool {
		for _, port := range resolvedPorts {
			if ptr.Deref(port.Trunk, false) {
				return true
			}
		}
		return false
	}
	if portUsesTrunk() {
		trunkSupported, err := s.IsTrunkExtSupported()
		if err != nil {
			return nil, err
		}

		if !trunkSupported {
			return nil, fmt.Errorf("there is no trunk support. please ensure that the trunk extension is enabled in your OpenStack deployment")
		}
	}

	return resolvedPorts, nil
}

// normalizePorts ensures that a user-specified PortOpts has all required fields set. Specifically it:
// - sets the Trunk field to the instance spec default if not specified
// - sets the Network ID field if not specified.
func (s *Service) normalizePorts(ports []infrav1.PortOpts, clusterResourceName, baseName string, trunkEnabled bool, defaultSecurityGroupIDs []string, defaultNetwork *infrav1.NetworkStatusWithSubnets, baseTags []string) ([]infrav1.ResolvedPortSpec, error) {
	normalizedPorts := make([]infrav1.ResolvedPortSpec, len(ports))
	for i := range ports {
		port := &ports[i]
		normalizedPort := &normalizedPorts[i]

		// Copy fields which don't need to be resolved
		normalizedPort.ResolvedPortSpecFields = port.ResolvedPortSpecFields

		// Generate a standardised name
		normalizedPort.Name = getPortName(baseName, port, i)

		// Generate a description if none is provided
		if port.Description != nil {
			normalizedPort.Description = *port.Description
		} else {
			normalizedPort.Description = names.GetDescription(clusterResourceName)
		}

		// Tags are inherited base tags plus any port-specific tags
		normalizedPort.Tags = slices.Concat(baseTags, port.Tags)

		// No Trunk field specified for the port, inherit the machine default
		if port.Trunk == nil {
			if trunkEnabled {
				normalizedPort.Trunk = &trunkEnabled
			}
		} else {
			normalizedPort.Trunk = port.Trunk
		}

		// Resolve network ID and fixed IPs
		var err error
		normalizedPort.NetworkID, normalizedPort.FixedIPs, err = s.normalizePortTarget(port, defaultNetwork, i)
		if err != nil {
			return nil, err
		}

		// Resolve security groups
		if len(port.SecurityGroups) == 0 {
			normalizedPort.SecurityGroups = defaultSecurityGroupIDs
		} else {
			normalizedPort.SecurityGroups, err = s.GetSecurityGroups(port.SecurityGroups)
			if err != nil {
				return nil, fmt.Errorf("error getting security groups: %v", err)
			}
		}
	}
	return normalizedPorts, nil
}

func defaultNetworkTarget(network *infrav1.NetworkStatusWithSubnets) (string, []infrav1.ResolvedFixedIP, error) {
	networkID := network.ID
	fixedIPs := make([]infrav1.ResolvedFixedIP, len(network.Subnets))
	for i := range network.Subnets {
		subnet := &network.Subnets[i]
		fixedIPs[i].SubnetID = &subnet.ID
	}
	return networkID, fixedIPs, nil
}

// normalizePortTarget ensures that the port has a network ID.
func (s *Service) normalizePortTarget(port *infrav1.PortOpts, defaultNetwork *infrav1.NetworkStatusWithSubnets, portIdx int) (string, []infrav1.ResolvedFixedIP, error) {
	// No network or subnets defined: use cluster defaults
	if port.Network == nil && len(port.FixedIPs) == 0 {
		return defaultNetworkTarget(defaultNetwork)
	}

	var networkID string
	var resolvedFixedIPs []infrav1.ResolvedFixedIP
	if len(port.FixedIPs) > 0 {
		resolvedFixedIPs = make([]infrav1.ResolvedFixedIP, len(port.FixedIPs))
	}

	switch {
	case port.Network != nil:
		var err error
		networkID, err = s.GetNetworkIDByParam(port.Network)
		if err != nil {
			return "", nil, err
		}

	// No network, but fixed IPs are defined(we handled the no fixed
	// IPs case above): try to infer network from a subnet
	case len(port.FixedIPs) > 0:
		s.scope.Logger().V(4).Info("No network defined for port, attempting to infer from subnet", "port", portIdx)

		// Look for a unique subnet defined in FixedIPs.  If we find one
		// we can use it to infer the network ID. We don't need to worry
		// here about the case where different FixedIPs have different
		// networks because that will cause an error later.
		var err error
		networkID, err = func() (string, error) {
			for i, fixedIP := range port.FixedIPs {
				resolvedFixedIP := &resolvedFixedIPs[i]

				if fixedIP.Subnet == nil {
					continue
				}

				subnet, err := s.GetSubnetByParam(fixedIP.Subnet)
				if err != nil {
					// Multiple matches might be ok later when we restrict matches to a single network
					if errors.Is(err, ErrMultipleMatches) {
						s.scope.Logger().V(4).Info("Couldn't infer network from subnet", "subnetIndex", i, "err", err)
						continue
					}

					return "", err
				}

				// Cache the known subnet ID in the FixedIP so we don't fetch it again later
				resolvedFixedIP.SubnetID = &subnet.ID
				return subnet.NetworkID, nil
			}

			// TODO: This is a spec error: it should set the machine to failed
			return "", fmt.Errorf("port %d has no network and unable to infer from fixed IPs", portIdx)
		}()
		if err != nil {
			return "", nil, err
		}

	default:
		// TODO: This is a spec errors: it should set the machine to failed
		return "", nil, fmt.Errorf("unable to determine network for port %d", portIdx)
	}

	// Network ID is now known. Resolve all FixedIPs
	for i, fixedIP := range port.FixedIPs {
		resolvedFixedIP := &resolvedFixedIPs[i]
		resolvedFixedIP.IPAddress = fixedIP.IPAddress
		if fixedIP.Subnet != nil && resolvedFixedIP.SubnetID == nil {
			subnet, err := s.GetNetworkSubnetByParam(networkID, fixedIP.Subnet)
			if err != nil {
				return "", nil, err
			}
			resolvedFixedIP.SubnetID = &subnet.ID
		}
	}

	return networkID, resolvedFixedIPs, nil
}

// IsTrunkExtSupported verifies trunk setup on the OpenStack deployment.
func (s *Service) IsTrunkExtSupported() (trunknSupported bool, err error) {
	trunkSupport, err := s.GetTrunkSupport()
	if err != nil {
		return false, fmt.Errorf("there was an issue verifying whether trunk support is available, Please try again later: %v", err)
	}
	if !trunkSupport {
		return false, nil
	}
	return true, nil
}

// AdoptPorts looks for ports in desiredPorts which were previously created, and adds them to resources.Ports.
// A port matches if it has the same name and network ID as the desired port.
func (s *Service) AdoptPorts(scope *scope.WithLogger, desiredPorts []infrav1.ResolvedPortSpec, resources *infrav1.MachineResources) error {
	// We can skip adoption if the ports are already in the status
	if len(desiredPorts) == len(resources.Ports) {
		return nil
	}

	scope.Logger().V(5).Info("Adopting ports")

	// We create ports in order and adopt them in order in PortsStatus.
	// This means that if port N doesn't exist we know that ports >N don't exist.
	// We can therefore stop searching for ports once we find one that doesn't exist.
	for i := range desiredPorts {
		// check if the port is in status first and if it is, skip it
		if i < len(resources.Ports) {
			scope.Logger().V(5).Info("Port already in status, skipping it", "port index", i)
			continue
		}

		portSpec := &desiredPorts[i]
		ports, err := s.client.ListPort(ports.ListOpts{
			Name:      portSpec.Name,
			NetworkID: portSpec.NetworkID,
		})
		if err != nil {
			return fmt.Errorf("searching for existing port %s in network %s: %v", portSpec.Name, portSpec.NetworkID, err)
		}
		// if the port is not found, we stop the adoption of ports since the rest of the ports will not be found either
		// and will be created after the adoption
		if len(ports) == 0 {
			scope.Logger().V(5).Info("Port not found, stopping the adoption of ports", "port index", i)
			return nil
		}
		if len(ports) > 1 {
			return fmt.Errorf("found multiple ports with name %s", portSpec.Name)
		}

		// The desired port was found, so we add it to the status
		portID := ports[0].ID
		scope.Logger().Info("Adopted previously created port which was not in status", "port index", i, "portID", portID)
		resources.Ports = append(resources.Ports, infrav1.PortStatus{ID: portID})
	}

	return nil
}
