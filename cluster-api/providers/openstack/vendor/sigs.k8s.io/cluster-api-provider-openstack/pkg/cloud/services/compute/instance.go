/*
Copyright 2018 The Kubernetes Authors.

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
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/hash"
)

const (
	retryIntervalInstanceStatus = 10 * time.Second
	timeoutInstanceCreate       = 5
	timeoutInstanceDelete       = 5 * time.Minute
)

// normalizePortTarget ensures that the port has a network ID.
func (s *Service) normalizePortTarget(port *infrav1.PortOpts, openStackCluster *infrav1.OpenStackCluster, portIdx int) error {
	// Treat no Network and empty Network the same
	noNetwork := port.Network == nil || (*port.Network == infrav1.NetworkFilter{})

	// No network or subnets defined: use cluster defaults
	if noNetwork && len(port.FixedIPs) == 0 {
		port.Network = &infrav1.NetworkFilter{
			ID: openStackCluster.Status.Network.ID,
		}
		for _, subnet := range openStackCluster.Status.Network.Subnets {
			port.FixedIPs = append(port.FixedIPs, infrav1.FixedIP{
				Subnet: &infrav1.SubnetFilter{
					ID: subnet.ID,
				},
			})
		}

		return nil
	}

	// No network, but fixed IPs are defined(we handled the no fixed
	// IPs case above): try to infer network from a subnet
	if noNetwork {
		s.scope.Logger().V(4).Info("No network defined for port, attempting to infer from subnet", "port", portIdx)

		// Look for a unique subnet defined in FixedIPs.  If we find one
		// we can use it to infer the network ID. We don't need to worry
		// here about the case where different FixedIPs have different
		// networks because that will cause an error later when we try
		// to create the port.
		networkID, err := func() (string, error) {
			networkingService, err := s.getNetworkingService()
			if err != nil {
				return "", err
			}

			for i, fixedIP := range port.FixedIPs {
				if fixedIP.Subnet == nil {
					continue
				}

				subnet, err := networkingService.GetSubnetByFilter(fixedIP.Subnet)
				if err != nil {
					// Multiple matches might be ok later when we restrict matches to a single network
					if errors.Is(err, networking.ErrMultipleMatches) {
						s.scope.Logger().V(4).Info("Couldn't infer network from subnet", "subnetIndex", i, "err", err)
						continue
					}

					return "", err
				}

				// Cache the subnet ID in the FixedIP
				fixedIP.Subnet.ID = subnet.ID
				return subnet.NetworkID, nil
			}

			// TODO: This is a spec error: it should set the machine to failed
			return "", fmt.Errorf("port %d has no network and unable to infer from fixed IPs", portIdx)
		}()
		if err != nil {
			return err
		}

		port.Network = &infrav1.NetworkFilter{
			ID: networkID,
		}

		return nil
	}

	// Nothing to do if network ID is already set
	if port.Network.ID != "" {
		return nil
	}

	// Network is defined by Filter
	networkingService, err := s.getNetworkingService()
	if err != nil {
		return err
	}

	netIDs, err := networkingService.GetNetworkIDsByFilter(port.Network.ToListOpt())
	if err != nil {
		return err
	}

	// TODO: These are spec errors: they should set the machine to failed
	if len(netIDs) > 1 {
		return fmt.Errorf("network filter for port %d returns more than one result", portIdx)
	} else if len(netIDs) == 0 {
		return fmt.Errorf("network filter for port %d returns no networks", portIdx)
	}

	port.Network.ID = netIDs[0]

	return nil
}

// normalizePorts ensures that a user-specified PortOpts has all required fields set. Specifically it:
// - sets the Trunk field to the instance spec default if not specified
// - sets the Network ID field if not specified.
func (s *Service) normalizePorts(ports []infrav1.PortOpts, openStackCluster *infrav1.OpenStackCluster, instanceSpec *InstanceSpec) ([]infrav1.PortOpts, error) {
	normalizedPorts := make([]infrav1.PortOpts, 0, len(ports))
	for i := range ports {
		// Deep copy the port to avoid mutating the original
		port := ports[i].DeepCopy()

		// No Trunk field specified for the port, inherit the machine default
		if port.Trunk == nil {
			port.Trunk = &instanceSpec.Trunk
		}

		if err := s.normalizePortTarget(port, openStackCluster, i); err != nil {
			return nil, err
		}

		normalizedPorts = append(normalizedPorts, *port)
	}
	return normalizedPorts, nil
}

// constructPorts builds an array of ports from the instance spec.
// If no ports are in the spec, returns a single port for a network connection to the default cluster network.
func (s *Service) constructPorts(openStackCluster *infrav1.OpenStackCluster, instanceSpec *InstanceSpec) ([]infrav1.PortOpts, error) {
	// Ensure user-specified ports have all required fields
	ports, err := s.normalizePorts(instanceSpec.Ports, openStackCluster, instanceSpec)
	if err != nil {
		return nil, err
	}

	// no networks or ports found in the spec, so create a port on the cluster network
	if len(ports) == 0 {
		port := infrav1.PortOpts{
			Network: &infrav1.NetworkFilter{
				ID: openStackCluster.Status.Network.ID,
			},
			Trunk: &instanceSpec.Trunk,
		}
		for _, subnet := range openStackCluster.Status.Network.Subnets {
			port.FixedIPs = append(port.FixedIPs, infrav1.FixedIP{
				Subnet: &infrav1.SubnetFilter{
					ID: subnet.ID,
				},
			})
		}
		ports = []infrav1.PortOpts{port}
	}

	// trunk support is required if any port has trunk enabled
	portUsesTrunk := func() bool {
		for _, port := range ports {
			if port.Trunk != nil && *port.Trunk {
				return true
			}
		}
		return false
	}
	if portUsesTrunk() {
		trunkSupported, err := s.isTrunkExtSupported()
		if err != nil {
			return nil, err
		}
		if !trunkSupported {
			return nil, fmt.Errorf("there is no trunk support. please ensure that the trunk extension is enabled in your OpenStack deployment")
		}
	}

	return ports, nil
}

func (s *Service) CreateInstance(eventObject runtime.Object, openStackCluster *infrav1.OpenStackCluster, instanceSpec *InstanceSpec, clusterName string, isBastion bool) (*InstanceStatus, error) {
	return s.createInstanceImpl(eventObject, openStackCluster, instanceSpec, clusterName, isBastion, retryIntervalInstanceStatus)
}

func (s *Service) getAndValidateFlavor(flavorName string, isBastion bool) (*flavors.Flavor, error) {
	f, err := s.getComputeClient().GetFlavorFromName(flavorName)
	if err != nil {
		return nil, fmt.Errorf("error getting flavor from flavor name %s: %v", flavorName, err)
	}
	if !isBastion && f.VCPUs <= 1 {
		return nil, fmt.Errorf("kubeadm requires a minimum of 2 vCPUs, pick a flavor with at least 2 vCPUs")
	}

	return f, nil
}

func (s *Service) createInstanceImpl(eventObject runtime.Object, openStackCluster *infrav1.OpenStackCluster, instanceSpec *InstanceSpec, clusterName string, isBastion bool, retryInterval time.Duration) (*InstanceStatus, error) {
	var server *clients.ServerExt
	portList := []servers.Network{}

	imageID, err := s.getImageID(instanceSpec.ImageUUID, instanceSpec.Image)
	if err != nil {
		return nil, fmt.Errorf("error getting image ID: %v", err)
	}

	flavor, err := s.getAndValidateFlavor(instanceSpec.Flavor, isBastion)
	if err != nil {
		return nil, err
	}

	// Ensure we delete the ports we created if we haven't created the server.
	defer func() {
		if server != nil {
			return
		}

		if err := s.deletePorts(eventObject, portList); err != nil {
			s.scope.Logger().V(4).Error(err, "Failed to clean up ports after failure")
		}
	}()

	ports, err := s.constructPorts(openStackCluster, instanceSpec)
	if err != nil {
		return nil, err
	}

	networkingService, err := s.getNetworkingService()
	if err != nil {
		return nil, err
	}

	securityGroups, err := networkingService.GetSecurityGroups(instanceSpec.SecurityGroups)
	if err != nil {
		return nil, fmt.Errorf("error getting security groups: %v", err)
	}

	for i := range ports {
		portOpts := &ports[i]
		iTags := []string{}
		if len(instanceSpec.Tags) > 0 {
			iTags = instanceSpec.Tags
		}
		portName := networking.GetPortName(instanceSpec.Name, portOpts, i)
		port, err := networkingService.GetOrCreatePort(eventObject, clusterName, portName, portOpts, securityGroups, iTags)
		if err != nil {
			return nil, err
		}

		portList = append(portList, servers.Network{
			Port: port.ID,
		})
	}

	volume, err := s.getOrCreateRootVolume(eventObject, instanceSpec, imageID)
	if err != nil {
		return nil, fmt.Errorf("error in get or create root volume: %w", err)
	}

	instanceCreateTimeout := getTimeout("CLUSTER_API_OPENSTACK_INSTANCE_CREATE_TIMEOUT", timeoutInstanceCreate)
	instanceCreateTimeout *= time.Minute

	// Wait for volume to become available
	if volume != nil {
		err = wait.PollUntilContextTimeout(context.TODO(), retryIntervalInstanceStatus, instanceCreateTimeout, true, func(_ context.Context) (bool, error) {
			createdVolume, err := s.getVolumeClient().GetVolume(volume.ID)
			if err != nil {
				if capoerrors.IsRetryable(err) {
					return false, nil
				}
				return false, err
			}

			switch createdVolume.Status {
			case "available":
				return true, nil
			case "error":
				return false, fmt.Errorf("volume %s is in error state", volume.ID)
			default:
				return false, nil
			}
		})
		if err != nil {
			return nil, fmt.Errorf("volume %s did not become available: %w", volume.ID, err)
		}
	}

	// Don't set ImageRef on the server if we're booting from volume
	var serverImageRef string
	if volume == nil {
		serverImageRef = imageID
	}

	var serverCreateOpts servers.CreateOptsBuilder = servers.CreateOpts{
		Name:             instanceSpec.Name,
		ImageRef:         serverImageRef,
		FlavorRef:        flavor.ID,
		AvailabilityZone: instanceSpec.FailureDomain,
		Networks:         portList,
		UserData:         []byte(instanceSpec.UserData),
		Tags:             instanceSpec.Tags,
		Metadata:         instanceSpec.Metadata,
		ConfigDrive:      &instanceSpec.ConfigDrive,
	}

	serverCreateOpts = applyRootVolume(serverCreateOpts, volume)

	serverCreateOpts = applyServerGroupID(serverCreateOpts, instanceSpec.ServerGroupID)

	server, err = s.getComputeClient().CreateServer(keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           instanceSpec.SSHKeyName,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Openstack instance: %v", err)
	}

	var createdInstance *InstanceStatus
	err = wait.PollUntilContextTimeout(context.TODO(), retryInterval, instanceCreateTimeout, true, func(_ context.Context) (bool, error) {
		createdInstance, err = s.GetInstanceStatus(server.ID)
		if err != nil {
			if capoerrors.IsRetryable(err) {
				return false, nil
			}
			return false, err
		}
		if createdInstance.State() == infrav1.InstanceStateError {
			return false, fmt.Errorf("error creating OpenStack instance %s, status changed to error", createdInstance.ID())
		}
		return createdInstance.State() == infrav1.InstanceStateActive, nil
	})
	if err != nil {
		record.Warnf(eventObject, "FailedCreateServer", "Failed to create server %s: %v", createdInstance.Name(), err)
		return nil, err
	}

	record.Eventf(eventObject, "SuccessfulCreateServer", "Created server %s with id %s", createdInstance.Name(), createdInstance.ID())
	return createdInstance, nil
}

func rootVolumeName(instanceName string) string {
	return fmt.Sprintf("%s-root", instanceName)
}

func hasRootVolume(rootVolume *infrav1.RootVolume) bool {
	return rootVolume != nil && rootVolume.Size > 0
}

func (s *Service) getVolumeByName(name string) (*volumes.Volume, error) {
	listOpts := volumes.ListOpts{
		AllTenants: false,
		Name:       name,
		TenantID:   s.scope.ProjectID(),
	}
	volumeList, err := s.getVolumeClient().ListVolumes(listOpts)
	if err != nil {
		return nil, fmt.Errorf("error listing volumes: %w", err)
	}
	if len(volumeList) > 1 {
		return nil, fmt.Errorf("expected to find a single volume called %s; found %d", name, len(volumeList))
	}
	if len(volumeList) == 0 {
		return nil, nil
	}
	return &volumeList[0], nil
}

func (s *Service) getOrCreateRootVolume(eventObject runtime.Object, instanceSpec *InstanceSpec, imageID string) (*volumes.Volume, error) {
	rootVolume := instanceSpec.RootVolume
	if !hasRootVolume(rootVolume) {
		return nil, nil
	}

	name := rootVolumeName(instanceSpec.Name)
	size := rootVolume.Size

	volume, err := s.getVolumeByName(name)
	if err != nil {
		return nil, err
	}
	if volume != nil {
		if volume.Size != size {
			return nil, fmt.Errorf("exected to find volume %s with size %d; found size %d", name, size, volume.Size)
		}

		s.scope.Logger().V(3).Info("Using existing root volume", "name", name)
		return volume, nil
	}

	availabilityZone := instanceSpec.FailureDomain
	if rootVolume.AvailabilityZone != "" {
		availabilityZone = rootVolume.AvailabilityZone
	}

	createOpts := volumes.CreateOpts{
		Size:             rootVolume.Size,
		Description:      fmt.Sprintf("Root volume for %s", instanceSpec.Name),
		Name:             rootVolumeName(instanceSpec.Name),
		ImageID:          imageID,
		Multiattach:      false,
		AvailabilityZone: availabilityZone,
		VolumeType:       rootVolume.VolumeType,
	}
	volume, err = s.getVolumeClient().CreateVolume(createOpts)
	if err != nil {
		record.Eventf(eventObject, "FailedCreateVolume", "Failed to create root volume; size=%d imageID=%s err=%v", size, imageID, err)
		return nil, err
	}
	record.Eventf(eventObject, "SuccessfulCreateVolume", "Created root volume; id=%s", volume.ID)
	return volume, err
}

// applyRootVolume sets a root volume if the root volume Size is not 0.
func applyRootVolume(opts servers.CreateOptsBuilder, volume *volumes.Volume) servers.CreateOptsBuilder {
	if volume == nil {
		return opts
	}

	block := bootfromvolume.BlockDevice{
		SourceType:          bootfromvolume.SourceVolume,
		BootIndex:           0,
		UUID:                volume.ID,
		DeleteOnTermination: true,
		DestinationType:     bootfromvolume.DestinationVolume,
	}
	return bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: opts,
		BlockDevice:       []bootfromvolume.BlockDevice{block},
	}
}

// applyServerGroupID adds a scheduler hint to the CreateOptsBuilder, if the
// spec contains a server group ID.
func applyServerGroupID(opts servers.CreateOptsBuilder, serverGroupID string) servers.CreateOptsBuilder {
	if serverGroupID != "" {
		return schedulerhints.CreateOptsExt{
			CreateOptsBuilder: opts,
			SchedulerHints: schedulerhints.SchedulerHints{
				Group: serverGroupID,
			},
		}
	}
	return opts
}

// Helper function for getting image id from name.
func (s *Service) getImageIDFromName(imageName string) (string, error) {
	var opts images.ListOpts

	opts.Name = imageName

	allImages, err := s.getImageClient().ListImages(opts)
	if err != nil {
		return "", err
	}

	switch len(allImages) {
	case 0:
		return "", fmt.Errorf("no image with the Name %s could be found", imageName)
	case 1:
		return allImages[0].ID, nil
	default:
		// this should never happen
		return "", fmt.Errorf("too many images with the name, %s, were found", imageName)
	}
}

// Helper function for getting image ID from name or ID.
func (s *Service) getImageID(imageUUID, imageName string) (string, error) {
	if imageUUID != "" {
		// we return imageUUID without check
		return imageUUID, nil
	} else if imageName != "" {
		return s.getImageIDFromName(imageName)
	}

	return "", nil
}

// GetManagementPort returns the port which is used for management and external
// traffic. Cluster floating IPs must be associated with this port.
func (s *Service) GetManagementPort(openStackCluster *infrav1.OpenStackCluster, instanceStatus *InstanceStatus) (*ports.Port, error) {
	ns, err := instanceStatus.NetworkStatus()
	if err != nil {
		return nil, err
	}

	networkingService, err := s.getNetworkingService()
	if err != nil {
		return nil, err
	}

	allPorts, err := networkingService.GetPortFromInstanceIP(instanceStatus.ID(), ns.IP(openStackCluster.Status.Network.Name))
	if err != nil {
		return nil, fmt.Errorf("lookup management port for server %s: %w", instanceStatus.ID(), err)
	}
	if len(allPorts) < 1 {
		return nil, fmt.Errorf("did not find management port for server %s", instanceStatus.ID())
	}
	return &allPorts[0], nil
}

func (s *Service) DeleteInstance(openStackCluster *infrav1.OpenStackCluster, eventObject runtime.Object, instanceStatus *InstanceStatus, instanceSpec *InstanceSpec) error {
	if instanceStatus == nil {
		/*
			We create a boot-from-volume instance in 2 steps:
			1. Create the volume
			2. Create the instance with the created root volume and set DeleteOnTermination

			This introduces a new failure mode which has implications for safely deleting instances: we
			might create the volume, but the instance create fails. This would leave us with a dangling
			volume with no instance.

			To handle this safely, we ensure that we never remove a machine finalizer until all resources
			associated with the instance, including a root volume, have been deleted. To achieve this:
			* We always call DeleteInstance when reconciling a delete, regardless of
			  whether the instance exists or not.
			* If the instance was already deleted we check that the volume is also gone.

			Note that we don't need to separately delete the root volume when deleting the instance because
			DeleteOnTermination will ensure it is deleted in that case.
		*/
		if hasRootVolume(instanceSpec.RootVolume) {
			name := rootVolumeName(instanceSpec.Name)
			volume, err := s.getVolumeByName(name)
			if err != nil {
				return err
			}
			if volume == nil {
				return nil
			}

			s.scope.Logger().V(2).Info("Deleting dangling root volume", "name", volume.Name, "ID", volume.ID)
			return s.getVolumeClient().DeleteVolume(volume.ID, volumes.DeleteOpts{})
		}

		return nil
	}

	instanceInterfaces, err := s.getComputeClient().ListAttachedInterfaces(instanceStatus.ID())
	if err != nil {
		return err
	}

	trunkSupported, err := s.isTrunkExtSupported()
	if err != nil {
		return fmt.Errorf("obtaining network extensions: %v", err)
	}

	networkingService, err := s.getNetworkingService()
	if err != nil {
		return err
	}

	// get and delete trunks
	for _, port := range instanceInterfaces {
		if err = s.deleteAttachInterface(eventObject, instanceStatus.InstanceIdentifier(), port.PortID); err != nil {
			return err
		}

		if trunkSupported {
			if err = networkingService.DeleteTrunk(eventObject, port.PortID); err != nil {
				return err
			}
		}
		if err = networkingService.DeletePort(eventObject, port.PortID); err != nil {
			return err
		}
	}

	// delete port of error instance
	if instanceStatus.State() == infrav1.InstanceStateError {
		portOpts, err := s.constructPorts(openStackCluster, instanceSpec)
		if err != nil {
			return err
		}

		if err := networkingService.GarbageCollectErrorInstancesPort(eventObject, instanceSpec.Name, portOpts); err != nil {
			return err
		}
	}

	return s.deleteInstance(eventObject, instanceStatus.InstanceIdentifier())
}

func (s *Service) deletePorts(eventObject runtime.Object, nets []servers.Network) error {
	trunkSupported, err := s.isTrunkExtSupported()
	if err != nil {
		return err
	}

	for _, n := range nets {
		networkingService, err := s.getNetworkingService()
		if err != nil {
			return err
		}

		if trunkSupported {
			if err = networkingService.DeleteTrunk(eventObject, n.Port); err != nil {
				return err
			}
		}
		if err := networkingService.DeletePort(eventObject, n.Port); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) deleteAttachInterface(eventObject runtime.Object, instance *InstanceIdentifier, portID string) error {
	err := s.getComputeClient().DeleteAttachedInterface(instance.ID, portID)
	if err != nil {
		if capoerrors.IsNotFound(err) {
			record.Eventf(eventObject, "SuccessfulDeleteAttachInterface", "Attach interface did not exist: instance %s, port %s", instance.ID, portID)
			return nil
		}
		if capoerrors.IsConflict(err) {
			// we don't want to block deletion because of Conflict
			// due to instance must be paused/active/shutoff in order to detach interface
			return nil
		}
		record.Warnf(eventObject, "FailedDeleteAttachInterface", "Failed to delete attach interface: instance %s, port %s: %v", instance.ID, portID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDeleteAttachInterface", "Deleted attach interface: instance %s, port %s", instance.ID, portID)
	return nil
}

func (s *Service) deleteInstance(eventObject runtime.Object, instance *InstanceIdentifier) error {
	err := s.getComputeClient().DeleteServer(instance.ID)
	if err != nil {
		if capoerrors.IsNotFound(err) {
			record.Eventf(eventObject, "SuccessfulDeleteServer", "Server %s with id %s did not exist", instance.Name, instance.ID)
			return nil
		}
		record.Warnf(eventObject, "FailedDeleteServer", "Failed to delete server %s with id %s: %v", instance.Name, instance.ID, err)
		return err
	}

	err = wait.PollUntilContextTimeout(context.TODO(), retryIntervalInstanceStatus, timeoutInstanceDelete, true, func(_ context.Context) (bool, error) {
		i, err := s.GetInstanceStatus(instance.ID)
		if err != nil {
			return false, err
		}
		if i != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		record.Warnf(eventObject, "FailedDeleteServer", "Failed to delete server %s with id %s: %v", instance.Name, instance.ID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDeleteServer", "Deleted server %s with id %s", instance.Name, instance.ID)
	return nil
}

func (s *Service) GetInstanceStatus(resourceID string) (instance *InstanceStatus, err error) {
	if resourceID == "" {
		return nil, fmt.Errorf("resourceId should be specified to get detail")
	}

	server, err := s.getComputeClient().GetServer(resourceID)
	if err != nil {
		if capoerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get server %q detail failed: %v", resourceID, err)
	}

	return &InstanceStatus{server, s.scope.Logger()}, nil
}

func (s *Service) GetInstanceStatusByName(eventObject runtime.Object, name string) (instance *InstanceStatus, err error) {
	var listOpts servers.ListOpts
	if name != "" {
		listOpts = servers.ListOpts{
			// The name parameter to /servers is a regular expression. Unless we
			// explicitly specify a whole string match this will be a substring
			// match.
			Name: fmt.Sprintf("^%s$", name),
		}
	} else {
		listOpts = servers.ListOpts{}
	}

	serverList, err := s.getComputeClient().ListServers(listOpts)
	if err != nil {
		return nil, fmt.Errorf("get server list: %v", err)
	}

	if len(serverList) > 1 {
		record.Warnf(eventObject, "DuplicateServerNames", "Found %d servers with name '%s'. This is likely to cause errors.", len(serverList), name)
	}

	// Return the first returned server, if any
	for i := range serverList {
		return &InstanceStatus{&serverList[i], s.scope.Logger()}, nil
	}
	return nil, nil
}

func getTimeout(name string, timeout int) time.Duration {
	if v := os.Getenv(name); v != "" {
		timeout, err := strconv.Atoi(v)
		if err == nil {
			return time.Duration(timeout)
		}
	}
	return time.Duration(timeout)
}

// isTrunkExtSupported verifies trunk setup on the OpenStack deployment.
func (s *Service) isTrunkExtSupported() (trunknSupported bool, err error) {
	networkingService, err := s.getNetworkingService()
	if err != nil {
		return false, err
	}

	trunkSupport, err := networkingService.GetTrunkSupport()
	if err != nil {
		return false, fmt.Errorf("there was an issue verifying whether trunk support is available, Please try again later: %v", err)
	}
	if !trunkSupport {
		return false, nil
	}
	return true, nil
}

func HashInstanceSpec(computeInstance *InstanceSpec) (string, error) {
	instanceHash, err := hash.ComputeSpewHash(computeInstance)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(instanceHash)), nil
}
