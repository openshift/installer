package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Provider is the vSphere implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.PreProvider = Provider{}

// Name returns the vsphere provider name.
func (p Provider) Name() string {
	return vsphere.Name
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status when declaring machines ready.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.InternalIP }

func initializeFoldersAndTemplates(ctx context.Context, cachedImage string, failureDomain vsphere.FailureDomain, session *session.Session, diskType vsphere.DiskType, clusterID, tagID string) error {
	finder := session.Finder

	dc, err := finder.Datacenter(ctx, failureDomain.Topology.Datacenter)
	if err != nil {
		return err
	}

	// Upstream govmomi bug, workaround
	// https://github.com/vmware/govmomi/issues/3523
	folderPath := path.Join(dc.InventoryPath, "vm", clusterID)

	// we must set the Folder to the infraId somewhere, we will need to remove that.
	// if we are overwriting folderPath it needs to have a slash (path)
	folder := failureDomain.Topology.Folder
	if strings.Contains(folder, "/") {
		folderPath = folder
	}

	var folderObj *object.Folder

	// Only createFolder() and attach the tag if the folder does not exist prior to installing
	if folderObj, err = folderExists(ctx, folderPath, session); folderObj == nil && err == nil {
		folderObj, err = createFolder(ctx, folderPath, session)
		if err != nil {
			return fmt.Errorf("unable to create folder: %w", err)
		}
		// attach tag to folder
		err = session.TagManager.AttachTag(ctx, tagID, folderObj.Reference())
		if err != nil {
			return fmt.Errorf("unable to attach tag to folder: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to get folder: %w", err)
	}

	// if the template is empty, the ova must be imported
	if len(failureDomain.Topology.Template) == 0 {
		if err = importRhcosOva(ctx, session, folderObj,
			cachedImage, clusterID, tagID, string(diskType), failureDomain); err != nil {
			return fmt.Errorf("failed to import ova: %w", err)
		}
	}
	return nil
}

// PreProvision creates the vCenter objects required prior to running capv.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	/*
	 * one locally cached image
	 * one tag and tag category per vcenter
	 * one folder per datacenter
	 * one template per region/zone aka failuredomain
	 */

	installConfig := in.InstallConfig
	clusterID := &installconfig.ClusterID{InfraID: in.InfraID}
	var tagID, cachedImage string
	var err error

	if downloadImage(installConfig) {
		cachedImage, err = cache.DownloadImageFile(in.RhcosImage.ControlPlane, cache.InstallerApplicationName)
		if err != nil {
			return fmt.Errorf("failed to use cached vsphere image: %w", err)
		}
	}

	for _, vcenter := range installConfig.Config.VSphere.VCenters {
		server := vcenter.Server
		vctrSession, err := installConfig.VSphere.Session(context.TODO(), server)

		if err != nil {
			return err
		}

		tagID, err = createClusterTagID(ctx, vctrSession, clusterID.InfraID)
		if err != nil {
			return err
		}

		for i := range in.MachineManifests {
			if vm, ok := in.MachineManifests[i].(*v1beta1.VSphereMachine); ok {
				if vm.Spec.Server == server {
					vm.Spec.TagIDs = append(vm.Spec.TagIDs, tagID)
				}
			}
		}

		for _, failureDomain := range installConfig.Config.VSphere.FailureDomains {
			if failureDomain.Server != server {
				continue
			}

			if err = initializeFoldersAndTemplates(ctx, cachedImage, failureDomain, vctrSession, installConfig.Config.VSphere.DiskType, clusterID.InfraID, tagID); err != nil {
				return fmt.Errorf("unable to initialize folders and templates: %w", err)
			}
		}
	}

	return nil
}

// InfraReady is called once cluster.Status.InfrastructureReady
// is true, typically after load balancers have been provisioned. It can be used
// to create DNS records.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	return nil
}

// PostProvision should be called to add or update and vSphere resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	// We will want to check to see if ControlPlane machines need additional disks
	cpPool := in.InstallConfig.Config.ControlPlane.Platform.VSphere

	// If not set, nothing to do here.
	if cpPool == nil {
		return nil
	}

	// If we have any additional disks defined, we'll need to create them here
	if len(cpPool.AdditionalDisks) > 0 {
		logrus.Info("Adding additional disks to control plane machines")

		for i := range in.MachineManifests {
			if vm, ok := in.MachineManifests[i].(*v1beta1.VSphereMachine); ok {
				if !strings.HasSuffix(vm.Name, "bootstrap") {
					logrus.Infof("Adding additional disks to vm %s", vm.Name)
					server := vm.Spec.Server
					vctrSession, err := in.InstallConfig.VSphere.Session(context.TODO(), server)

					if err != nil {
						return err
					}
					err = addAdditionalDisks(ctx, vm, cpPool, vctrSession)

					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// downloadImage if any failure domains don't have a defined template, this function
// returns true.
func downloadImage(installConfig *installconfig.InstallConfig) bool {
	for _, fd := range installConfig.Config.VSphere.FailureDomains {
		if fd.Topology.Template == "" {
			return true
		}
	}
	return false
}

func addAdditionalDisks(ctx context.Context, machine *v1beta1.VSphereMachine, pool *vsphere.MachinePool, session *session.Session) error {
	logrus.Debugf("Getting vm %v", machine.Name)
	vm, err := session.Finder.VirtualMachine(ctx, fmt.Sprintf("%s/%s", machine.Spec.Folder, machine.Name))
	offset := 1 // For now, we assume only one disk is current attached to VM.

	if err != nil {
		return err
	}

	logrus.Debugf("Getting VM devices")
	devices, err := vm.Device(ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("Getting datastore %v", machine.Spec.Datastore)
	ds, err := session.Finder.Datastore(ctx, machine.Spec.Datastore)
	if err != nil {
		return err
	}

	// For now, we only do the active scsi controller
	logrus.Debug("Getting scsi controller")
	controller, err := devices.FindSCSIController("")
	if err != nil {
		return err
	}

	for diskIndex, newDisk := range pool.AdditionalDisks {
		logrus.Debugf("Attempting to add disk %d with size %dGB", diskIndex, newDisk.DiskSizeGB)

		disk := devices.CreateDisk(controller, ds.Reference(), "")

		existing := devices.SelectByBackingInfo(disk.Backing)

		if len(existing) > 0 {
			logrus.Warningf("Disk already present for index %d", diskIndex)
			return errors.New("disk already present")
		}

		disk.CapacityInKB = int64(newDisk.DiskSizeGB) * 1024 * 1024
		unitNumber := int32(offset + diskIndex)
		disk.VirtualDevice.UnitNumber = &unitNumber

		// Add disk using default profile of VM.
		logrus.Infof("Adding disk device to vm %v", disk)
		err = vm.AddDevice(ctx, disk)
		if err != nil {
			return err
		}
	}

	return nil
}
