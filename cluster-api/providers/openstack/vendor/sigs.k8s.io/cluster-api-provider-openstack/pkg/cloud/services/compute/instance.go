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
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/hash"
)

const (
	retryIntervalInstanceStatus = 10 * time.Second
	timeoutInstanceCreate       = 5
	timeoutInstanceDelete       = 5 * time.Minute
)

func (s *Service) CreateInstance(eventObject runtime.Object, instanceSpec *InstanceSpec, portIDs []string) (*InstanceStatus, error) {
	return s.createInstanceImpl(eventObject, instanceSpec, retryIntervalInstanceStatus, portIDs)
}

func (s *Service) getAndValidateFlavor(flavorName string) (*flavors.Flavor, error) {
	f, err := s.getComputeClient().GetFlavorFromName(flavorName)
	if err != nil {
		return nil, fmt.Errorf("error getting flavor from flavor name %s: %v", flavorName, err)
	}

	return f, nil
}

func (s *Service) createInstanceImpl(eventObject runtime.Object, instanceSpec *InstanceSpec, retryInterval time.Duration, portIDs []string) (*InstanceStatus, error) {
	var server *clients.ServerExt
	portList := []servers.Network{}

	flavor, err := s.getAndValidateFlavor(instanceSpec.Flavor)
	if err != nil {
		return nil, err
	}

	if len(portIDs) == 0 {
		return nil, fmt.Errorf("portIDs cannot be empty")
	}

	for _, portID := range portIDs {
		portList = append(portList, servers.Network{
			Port: portID,
		})
	}

	instanceCreateTimeout := getTimeout("CLUSTER_API_OPENSTACK_INSTANCE_CREATE_TIMEOUT", timeoutInstanceCreate)
	instanceCreateTimeout *= time.Minute

	// Don't set ImageRef on the server if we're booting from volume
	var serverImageRef string
	if !hasRootVolume(instanceSpec) {
		serverImageRef = instanceSpec.ImageID
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

	blockDevices, err := s.getBlockDevices(eventObject, instanceSpec, instanceSpec.ImageID, instanceCreateTimeout, retryInterval)
	if err != nil {
		return nil, err
	}
	if len(blockDevices) > 0 {
		serverCreateOpts = bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: serverCreateOpts,
			BlockDevice:       blockDevices,
		}
	}

	serverCreateOpts = applyServerGroupID(serverCreateOpts, instanceSpec.ServerGroupID)

	server, err = s.getComputeClient().CreateServer(keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           instanceSpec.SSHKeyName,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Openstack instance: %v", err)
	}

	record.Eventf(eventObject, "SuccessfulCreateServer", "Created server %s with id %s", server.Name, server.ID)
	return &InstanceStatus{server, s.scope.Logger()}, nil
}

func volumeName(instanceName string, nameSuffix string) string {
	return fmt.Sprintf("%s-%s", instanceName, nameSuffix)
}

func hasRootVolume(instanceSpec *InstanceSpec) bool {
	return instanceSpec.RootVolume != nil && instanceSpec.RootVolume.SizeGiB > 0
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

// getOrCreateVolume gets or creates a volume with the given options. It returns the volume that already exists or the
// newly created one. It returns an error if the volume creation failed or if the expected volume size is different from
// the one that already exists.
func (s *Service) getOrCreateVolume(eventObject runtime.Object, opts volumes.CreateOpts) (*volumes.Volume, error) {
	existingVolume, err := s.getVolumeByName(opts.Name)
	if err != nil {
		return nil, err
	}
	if existingVolume != nil {
		// TODO(emilien): Improve the checks here, there is an ongoing discussion in the community about how to do this
		// which would involve adding metadata to the volume.
		if existingVolume.Size != opts.Size {
			return nil, fmt.Errorf("expected to find volume %s with size %d; found size %d", opts.Name, opts.Size, existingVolume.Size)
		}

		s.scope.Logger().V(3).Info("Using existing volume", "name", opts.Name, "id", existingVolume.ID)
		return existingVolume, nil
	}

	createdVolume, err := s.getVolumeClient().CreateVolume(opts)
	if err != nil {
		record.Eventf(eventObject, "FailedCreateVolume", "Failed to create volume; name=%s size=%d err=%v", opts.Name, opts.Size, err)
		return nil, err
	}
	record.Eventf(eventObject, "SuccessfulCreateVolume", "Created volume; id=%s", createdVolume.ID)
	return createdVolume, err
}

func (s *Service) waitForVolume(volumeID string, timeout time.Duration, retryInterval time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), retryInterval, timeout, true, func(_ context.Context) (bool, error) {
		volume, err := s.getVolumeClient().GetVolume(volumeID)
		if err != nil {
			if capoerrors.IsRetryable(err) {
				return false, nil
			}
			return false, err
		}

		switch volume.Status {
		case "available":
			return true, nil
		case "error":
			return false, fmt.Errorf("volume %s is in error state", volumeID)
		default:
			return false, nil
		}
	})
}

// getOrCreateVolumeBuilder gets or creates a volume with the given options. It returns the volume that already exists or the newly created one.
// It returns an error if the volume creation failed or if the expected volume is different from the one that already exists.
func (s *Service) getOrCreateVolumeBuilder(eventObject runtime.Object, instanceSpec *InstanceSpec, blockDeviceSpec *infrav1.AdditionalBlockDevice, imageID string, description string) (*volumes.Volume, error) {
	availabilityZone, volType := resolveVolumeOpts(instanceSpec, blockDeviceSpec.Storage.Volume)

	createOpts := volumes.CreateOpts{
		Name:             volumeName(instanceSpec.Name, blockDeviceSpec.Name),
		Description:      description,
		Size:             blockDeviceSpec.SizeGiB,
		ImageID:          imageID,
		Multiattach:      false,
		AvailabilityZone: availabilityZone,
		VolumeType:       volType,
	}

	return s.getOrCreateVolume(eventObject, createOpts)
}

func resolveVolumeOpts(instanceSpec *InstanceSpec, volumeOpts *infrav1.BlockDeviceVolume) (az, volType string) {
	if volumeOpts == nil {
		return
	}

	volType = volumeOpts.Type

	volumeAZ := volumeOpts.AvailabilityZone
	if volumeAZ == nil {
		return
	}

	switch volumeAZ.From {
	case "", infrav1.VolumeAZFromName:
		// volumeAZ.Name is nil case should have been caught by validation
		if volumeAZ.Name != nil {
			az = string(*volumeAZ.Name)
		}
	case infrav1.VolumeAZFromMachine:
		az = instanceSpec.FailureDomain
	}
	return
}

// getBlockDevices returns a list of block devices that were created and attached to the instance. It returns an error
// if the root volume or any of the additional block devices could not be created.
func (s *Service) getBlockDevices(eventObject runtime.Object, instanceSpec *InstanceSpec, imageID string, timeout time.Duration, retryInterval time.Duration) ([]bootfromvolume.BlockDevice, error) {
	blockDevices := []bootfromvolume.BlockDevice{}

	if hasRootVolume(instanceSpec) {
		rootVolumeToBlockDevice := infrav1.AdditionalBlockDevice{
			Name:    "root",
			SizeGiB: instanceSpec.RootVolume.SizeGiB,
			Storage: infrav1.BlockDeviceStorage{
				Type:   infrav1.VolumeBlockDevice,
				Volume: &instanceSpec.RootVolume.BlockDeviceVolume,
			},
		}
		rootVolume, err := s.getOrCreateVolumeBuilder(eventObject, instanceSpec, &rootVolumeToBlockDevice, imageID, fmt.Sprintf("Root volume for %s", instanceSpec.Name))
		if err != nil {
			return nil, err
		}
		blockDevices = append(blockDevices, bootfromvolume.BlockDevice{
			SourceType:          bootfromvolume.SourceVolume,
			DestinationType:     bootfromvolume.DestinationVolume,
			UUID:                rootVolume.ID,
			BootIndex:           0,
			DeleteOnTermination: true,
		})
	} else {
		blockDevices = append(blockDevices, bootfromvolume.BlockDevice{
			SourceType:          bootfromvolume.SourceImage,
			DestinationType:     bootfromvolume.DestinationLocal,
			UUID:                imageID,
			BootIndex:           0,
			DeleteOnTermination: true,
		})
	}

	for i := range instanceSpec.AdditionalBlockDevices {
		blockDeviceSpec := instanceSpec.AdditionalBlockDevices[i]

		var bdUUID string
		var localDiskSizeGiB int
		var sourceType bootfromvolume.SourceType
		var destinationType bootfromvolume.DestinationType

		// There is also a validation in the openstackmachine webhook.
		if blockDeviceSpec.Name == "root" {
			return nil, fmt.Errorf("block device name 'root' is reserved")
		}

		if blockDeviceSpec.Storage.Type == infrav1.VolumeBlockDevice {
			blockDevice, err := s.getOrCreateVolumeBuilder(eventObject, instanceSpec, &blockDeviceSpec, "", fmt.Sprintf("Additional block device for %s", instanceSpec.Name))
			if err != nil {
				return nil, err
			}
			bdUUID = blockDevice.ID
			sourceType = bootfromvolume.SourceVolume
			destinationType = bootfromvolume.DestinationVolume
		} else if blockDeviceSpec.Storage.Type == infrav1.LocalBlockDevice {
			sourceType = bootfromvolume.SourceBlank
			destinationType = bootfromvolume.DestinationLocal
			localDiskSizeGiB = blockDeviceSpec.SizeGiB
		} else {
			return nil, fmt.Errorf("invalid block device type %s", blockDeviceSpec.Storage.Type)
		}

		blockDevices = append(blockDevices, bootfromvolume.BlockDevice{
			SourceType:          sourceType,
			DestinationType:     destinationType,
			UUID:                bdUUID,
			BootIndex:           -1,
			DeleteOnTermination: true,
			VolumeSize:          localDiskSizeGiB,
			Tag:                 blockDeviceSpec.Name,
		})
	}

	// Wait for any volumes in the block devices to become available
	if len(blockDevices) > 0 {
		for _, bd := range blockDevices {
			if bd.SourceType == bootfromvolume.SourceVolume {
				if err := s.waitForVolume(bd.UUID, timeout, retryInterval); err != nil {
					return nil, fmt.Errorf("volume %s did not become available: %w", bd.UUID, err)
				}
			}
		}
	}

	return blockDevices, nil
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

// Helper function for getting image ID from name, ID, or tags.
func (s *Service) GetImageID(image infrav1.ImageParam) (string, error) {
	if image.ID != nil {
		return *image.ID, nil
	}

	if image.Filter == nil {
		// Should have been caught by validation
		return "", errors.New("image id and filter are both nil")
	}

	listOpts := filterconvert.ImageFilterToListOpts(image.Filter)
	allImages, err := s.getImageClient().ListImages(listOpts)
	if err != nil {
		return "", err
	}

	switch len(allImages) {
	case 0:
		return "", fmt.Errorf("no images were found with the given image filter: %v", image)
	case 1:
		return allImages[0].ID, nil
	default:
		// this should never happen
		return "", fmt.Errorf("too many images were found with the given image filter: %v", image)
	}
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

func (s *Service) DeleteInstance(eventObject runtime.Object, instanceStatus *InstanceStatus) error {
	instance := instanceStatus.InstanceIdentifier()

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

// DeleteVolumes deletes any cinder volumes which were created for the instance.
// Note that this need only be called when the server was not successfully
// created. If the server was created the volume will have been added with
// DeleteOnTermination=true, and will be automatically cleaned up with the
// server.
// We don't pass InstanceSpec here because we only require instance name,
// rootVolume, and additionalBlockDevices, and resolving the whole InstanceSpec
// introduces unnecessary failure modes.
func (s *Service) DeleteVolumes(instanceName string, rootVolume *infrav1.RootVolume, additionalBlockDevices []infrav1.AdditionalBlockDevice) error {
	/*
		Attaching volumes to an instance is a two-step process:

		  1. Create the volume
		  2. Create the instance with the created volumes in RootVolume and AdditionalBlockDevices fields with DeleteOnTermination=true

		This has a possible failure mode where creating the volume succeeds but creating the instance
		fails. In this case, we want to make sure that the dangling volumes are cleaned up.

		To handle this safely, we ensure that we never remove a machine finalizer until all resources
		associated with the instance, including volumes, have been deleted. To achieve this:

		  * We always call DeleteInstance when reconciling a delete, even if the instance does not exist
		  * If the instance was already deleted we check that the volumes are also gone

		Note that we don't need to separately delete the volumes when deleting the instance because
		DeleteOnTermination will ensure it is deleted in that case.
	*/

	if rootVolume != nil && rootVolume.SizeGiB > 0 {
		if err := s.deleteVolume(instanceName, "root"); err != nil {
			return err
		}
	}
	for _, volumeSpec := range additionalBlockDevices {
		if err := s.deleteVolume(instanceName, volumeSpec.Name); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) deleteVolume(instanceName string, nameSuffix string) error {
	volumeName := volumeName(instanceName, nameSuffix)
	volume, err := s.getVolumeByName(volumeName)
	if err != nil {
		return err
	}
	if volume == nil {
		return nil
	}

	s.scope.Logger().V(2).Info("Deleting dangling volume", "name", volume.Name, "ID", volume.ID)
	return s.getVolumeClient().DeleteVolume(volume.ID, volumes.DeleteOpts{})
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

func HashInstanceSpec(computeInstance *InstanceSpec) (string, error) {
	instanceHash, err := hash.ComputeSpewHash(computeInstance)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(instanceHash)), nil
}
