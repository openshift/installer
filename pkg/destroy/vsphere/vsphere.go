package vsphere

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	ClusterID string
	InfraID   string

	Client     *vim25.Client
	RestClient *rest.Client

	Logger logrus.FieldLogger
}

// New returns an VSphere destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {

	vim25Client, restClient, err := vspheretypes.CreateVSphereClients(context.TODO(),
		metadata.ClusterPlatformMetadata.VSphere.VCenter,
		metadata.ClusterPlatformMetadata.VSphere.Username,
		metadata.ClusterPlatformMetadata.VSphere.Password)

	if err != nil {
		return nil, err
	}

	return &ClusterUninstaller{
		ClusterID:  metadata.ClusterID,
		InfraID:    metadata.InfraID,
		Client:     vim25Client,
		RestClient: restClient,
		Logger:     logger,
	}, nil
}

func deleteVirtualMachines(ctx context.Context, client *vim25.Client, virtualMachineMoList []mo.VirtualMachine, logger logrus.FieldLogger) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	if len(virtualMachineMoList) != 0 {
		for _, vmMO := range virtualMachineMoList {
			virtualMachineLogger := logger.WithField("VirtualMachine", vmMO.Name)
			vm := object.NewVirtualMachine(client, vmMO.Reference())
			if vmMO.Summary.Runtime.PowerState == "poweredOn" {
				task, err := vm.PowerOff(ctx)
				if err != nil {
					return err
				}
				task.Wait(ctx)
				virtualMachineLogger.Debug("Powered off")
			}

			task, err := vm.Destroy(ctx)
			if err != nil {
				return err
			}
			task.Wait(ctx)
			virtualMachineLogger.Info("Destroyed")
		}
	}
	return nil
}
func deleteFolder(ctx context.Context, client *vim25.Client, folderMoList []mo.Folder, logger logrus.FieldLogger) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	// If there are no children in the folder go ahead an remove it
	if len(folderMoList[0].ChildEntity) == 0 {
		folderLogger := logger.WithField("Folder", folderMoList[0].Name)

		folder := object.NewFolder(client, folderMoList[0].Reference())
		task, err := folder.Destroy(ctx)
		if err != nil {
			return err
		}
		task.Wait(ctx)
		folderLogger.Info("Destroyed")
	} else {
		return errors.Errorf("Expected Folder %s to be empty", folderMoList[0].Name)
	}

	return nil
}
func getFolderManagedObjects(ctx context.Context, client *vim25.Client, moRef []types.ManagedObjectReference) ([]mo.Folder, error) {
	var folderMoList []mo.Folder
	pc := property.DefaultCollector(client)
	err := pc.Retrieve(ctx, moRef, nil, &folderMoList)
	if err != nil {
		return nil, err
	}
	return folderMoList, nil
}
func getVirtualMachineManagedObjects(ctx context.Context, client *vim25.Client, moRef []types.ManagedObjectReference) ([]mo.VirtualMachine, error) {
	var virtualMachineMoList []mo.VirtualMachine

	pc := property.DefaultCollector(client)
	err := pc.Retrieve(ctx, moRef, nil, &virtualMachineMoList)
	if err != nil {
		return nil, err
	}
	return virtualMachineMoList, nil
}

func getAttachedObjectsOnTag(ctx context.Context, client *rest.Client, tagName string) ([]tags.AttachedObjects, error) {
	tagManager := tags.NewManager(client)
	attached, err := tagManager.GetAttachedObjectsOnTags(ctx, []string{tagName})
	if err != nil {
		return nil, err
	}

	return attached, nil
}

func deleteTag(ctx context.Context, client *rest.Client, tagID string) error {
	tagManager := tags.NewManager(client)
	tag, err := tagManager.GetTag(ctx, tagID)
	if err == nil {
		err = tagManager.DeleteTag(ctx, tag)
	}
	return err
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	var folderList []types.ManagedObjectReference
	var virtualMachineList []types.ManagedObjectReference

	o.Logger.Debug("find attached objects on tag")
	tagAttachedObjects, err := getAttachedObjectsOnTag(context.TODO(), o.RestClient, o.InfraID)
	if err != nil {
		return err
	}

	// Seperate the objects attached to the tag based on type
	// We only need Folder and VirtualMachine
	for _, attachedObject := range tagAttachedObjects {
		for _, ref := range attachedObject.ObjectIDs {
			if ref.Reference().Type == "Folder" {
				folderList = append(folderList, ref.Reference())
			}
			if ref.Reference().Type == "VirtualMachine" {
				virtualMachineList = append(virtualMachineList, ref.Reference())
			}
		}
	}

	// The installer should create at most one parent,
	// the parent to the VirtualMachines.
	// If there are more or less fail with error message.
	if len(folderList) > 1 {
		return errors.Errorf("Expected 1 Folder per tag but got %d", len(folderList))
	}

	o.Logger.Debug("find VirtualMachine objects")
	virtualMachineMoList, err := getVirtualMachineManagedObjects(context.TODO(), o.Client, virtualMachineList)
	if err != nil {
		return err
	}
	o.Logger.Debug("delete VirtualMachines")
	err = deleteVirtualMachines(context.TODO(), o.Client, virtualMachineMoList, o.Logger)
	if err != nil {
		return err
	}

	// In this case, folder was user-provided
	// and should not be deleted so we are done.
	if len(folderList) == 0 {
		return nil
	}

	o.Logger.Debug("find Folder objects")
	folderMoList, err := getFolderManagedObjects(context.TODO(), o.Client, folderList)
	if err != nil {
		o.Logger.Errorln(err)
		return err
	}

	o.Logger.Debug("delete Folder")
	err = deleteFolder(context.TODO(), o.Client, folderMoList, o.Logger)
	if err != nil {
		o.Logger.Errorln(err)
		return err
	}

	o.Logger.Debug("delete tag")
	if err = deleteTag(context.TODO(), o.RestClient, o.InfraID); err != nil {
		o.Logger.Errorln(err)
		return err
	}

	return nil
}
