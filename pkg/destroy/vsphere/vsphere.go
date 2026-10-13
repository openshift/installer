package vsphere

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/mo"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	ClusterID         string
	InfraID           string
	terraformPlatform string
	Logger            logrus.FieldLogger
	clients           []API
}

// New returns an VSphere destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {
	var clients []API

	// We have two ways of processing metadata. Older metadata has only 1 vcenter but configured at root level.  New
	// way is for all vcenter data to be part of the vcenters array.
	if len(metadata.VSphere.VCenters) > 0 {
		for _, vsphere := range metadata.VSphere.VCenters {
			logger.Info(fmt.Sprintf("Creating client for vCenter %v for destroy", vsphere.VCenter))
			client, err := NewClient(vsphere.VCenter, vsphere.Username, vsphere.Password, logger)
			if err != nil {
				return nil, err
			}
			clients = append(clients, client)
		}
	} else {
		client, err := NewClient(metadata.VSphere.VCenter, metadata.VSphere.Username, metadata.VSphere.Password, logger)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return newWithClient(logger, metadata, clients), nil
}

func newWithClient(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata, clients []API) *ClusterUninstaller {
	clusterUninstaller := &ClusterUninstaller{
		ClusterID:         metadata.ClusterID,
		InfraID:           metadata.InfraID,
		terraformPlatform: metadata.VSphere.TerraformPlatform,

		clients: clients,
		Logger:  logger,
	}

	return clusterUninstaller
}

func (o *ClusterUninstaller) deleteFolder(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for _, client := range o.clients {
		folderMoList, err := client.ListFolders(ctx, o.InfraID)
		if err != nil {
			return err
		}

		if len(folderMoList) == 0 {
			o.Logger.Debug("All folders deleted")
			return nil
		}

		// If there are no children in the folder, go ahead and remove it

		for _, f := range folderMoList {
			folderLogger := o.Logger.WithField("Folder", f.Name).WithField("vCenter", client.GetVCenterName())
			if numChildren := len(f.ChildEntity); numChildren > 0 {
				entities := make([]string, 0, numChildren)
				for _, child := range f.ChildEntity {
					entities = append(entities, fmt.Sprintf("%s:%s", child.Type, child.Value))
				}
				folderLogger.Errorf("Folder should be empty but contains %d objects: %s. The installer will retry removing \"virtualmachine\" objects, but any other type will need to be removed manually before the deprovision can proceed", numChildren, strings.Join(entities, ", "))
				return errors.Errorf("Expected Folder %s to be empty", f.Name)
			}
			err = client.DeleteFolder(ctx, f)
			if err != nil {
				folderLogger.Debug(err)
				return err
			}
			folderLogger.Info("Destroyed")
		}
	}

	return nil
}

func (o *ClusterUninstaller) deleteStoragePolicy(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	policyName := fmt.Sprintf("openshift-storage-policy-%s", o.InfraID)
	for _, client := range o.clients {
		policyLogger := o.Logger.WithField("StoragePolicy", policyName).WithField("vCenter", client.GetVCenterName())
		policyLogger.Debug("Destroying")
		err := client.DeleteStoragePolicy(ctx, policyName)
		if err != nil {
			policyLogger.Debug(err)
			return err
		}
		policyLogger.Info("Destroyed")
	}

	return nil
}

func (o *ClusterUninstaller) deleteTag(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	for _, client := range o.clients {
		tagLogger := o.Logger.WithField("Tag", o.InfraID).WithField("vCenter", client.GetVCenterName())
		tagLogger.Debug("Delete")
		err := client.DeleteTag(ctx, o.InfraID)
		if err != nil {
			tagLogger.Debug(err)
			return err
		}
		tagLogger.Info("Deleted")
	}

	return nil
}

func (o *ClusterUninstaller) deleteTagCategory(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	categoryID := "openshift-" + o.InfraID
	for _, client := range o.clients {
		tcLogger := o.Logger.WithField("TagCategory", categoryID).WithField("vCenter", client.GetVCenterName())
		tcLogger.Debug("Delete")
		err := client.DeleteTagCategory(ctx, categoryID)
		if err != nil {
			tcLogger.Errorln(err)
			return err
		}
		tcLogger.Info("Deleted")
	}

	return nil
}

func (o *ClusterUninstaller) stopVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine, client API) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name).WithField("vCenter", client.GetVCenterName())
	virtualMachineLogger.Debug("Powering off")
	err := client.StopVirtualMachine(ctx, vmMO)
	if err != nil {
		virtualMachineLogger.Debug(err)
		return err
	}
	virtualMachineLogger.Info("Powered off")

	return nil
}

func (o *ClusterUninstaller) stopVirtualMachines(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	var errs []error
	for _, client := range o.clients {
		clientLogger := o.Logger.WithField("vCenter", client.GetVCenterName())
		clientLogger.Debug("Powering Off Virtual Machines")
		found, err := client.ListVirtualMachines(ctx, o.InfraID)
		if err != nil {
			o.Logger.Debug(err)
			return err
		}

		// In theory, all failure domains should at least have the template in it.  If we get zero VMs back, this may
		// signify an issue getting VMs by tags from the vCenter.
		if len(found) == 0 {
			clientLogger.Warning("No Virtual Machines Found")
		}

		for _, vmMO := range found {
			stopLogger := o.Logger.WithField("VirtualMachine", vmMO.Name).WithField("IsPoweredOff", isPoweredOff(vmMO))
			stopLogger.Debug("Power State")
			if !isPoweredOff(vmMO) {
				if err := o.stopVirtualMachine(ctx, vmMO, client); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) deleteVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine, client API) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name)
	err := client.DeleteVirtualMachine(ctx, vmMO)
	if err != nil {
		virtualMachineLogger.Debug(err)
		return err
	}
	virtualMachineLogger.Info("Destroyed")

	return nil
}

func (o *ClusterUninstaller) deleteVirtualMachines(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	var errs []error
	for _, client := range o.clients {
		clientLogger := o.Logger.WithField("vCenter", client.GetVCenterName())
		clientLogger.Debug("Delete Virtual Machines")
		found, err := client.ListVirtualMachines(ctx, o.InfraID)
		if err != nil {
			o.Logger.Debug(err)
			return err
		}

		// In theory, all failure domains should at least have the template in it.  If we get zero VMs back, this may
		// signify an issue getting VMs by tags from the vCenter.
		if len(found) == 0 {
			clientLogger.Warning("No Virtual Machines Found")
		}

		for _, vmMO := range found {
			if err := o.deleteVirtualMachine(ctx, vmMO, client); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) deleteHostZoneObjects(ctx context.Context) error {
	var errs []error
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	for _, client := range o.clients {
		o.Logger.WithField("vCenter", client.GetVCenterName()).Debug("Delete Host Zone Objects")
		if err := client.DeleteHostZoneObjects(ctx, o.InfraID); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) deleteCnsVolumes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	for _, client := range o.clients {
		o.Logger.WithField("vCenter", client.GetVCenterName()).Debug("Delete CNS Volumes")
		cnsVolumes, err := client.GetCnsVolumes(ctx, o.InfraID)
		if err != nil {
			return err
		}

		for _, cv := range cnsVolumes {
			cnsVolumeLogger := o.Logger.WithField("CNSVolume", cv.VolumeId.Id).WithField("vCenter", client.GetVCenterName())
			cnsVolumeLogger.Debug("Destroying")
			err := client.DeleteCnsVolumes(ctx, cv)
			if err != nil {
				return err
			}
			cnsVolumeLogger.Info("Destroyed")
		}
	}

	return nil
}

func (o *ClusterUninstaller) destroyCluster(ctx context.Context) (bool, error) {
	o.Logger.Debug("Destroying cluster")
	stagedFuncs := [][]struct {
		name    string
		execute func(context.Context) error
	}{
		{
			{name: "Stop virtual machines", execute: o.stopVirtualMachines},
		},
		{
			{name: "Delete Virtual Machines", execute: o.deleteVirtualMachines},
		},
		{
			{name: "Delete CNS Volumes", execute: o.deleteCnsVolumes},
		},
		{
			{name: "Folder", execute: o.deleteFolder},
		},
		{
			{name: "Storage Policy", execute: o.deleteStoragePolicy},
			{name: "Tag", execute: o.deleteTag},
			{name: "Tag Category", execute: o.deleteTagCategory},
		},
		{
			{name: "VM Groups and VM Host Rules", execute: o.deleteHostZoneObjects},
		},
	}

	stageFailed := false
	for _, stage := range stagedFuncs {
		if stageFailed {
			break
		}
		for _, f := range stage {
			err := f.execute(ctx)
			if err != nil {
				o.Logger.Debugf("%s: %v", f.name, err)
				stageFailed = true
			}
		}
	}

	return !stageFailed, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*installertypes.ClusterQuota, error) {
	for _, client := range o.clients {
		defer client.Logout()
	}

	err := wait.PollUntilContextCancel(
		context.Background(),
		time.Second*10,
		true,
		o.destroyCluster,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}
