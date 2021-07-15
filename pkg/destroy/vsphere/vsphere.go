package vsphere

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/pbm"
	pbmtypes "github.com/vmware/govmomi/pbm/types"
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

func deleteStoragePolicy(ctx context.Context, client *vim25.Client, infraID string, logger logrus.FieldLogger) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	rtype := pbmtypes.PbmProfileResourceType{
		ResourceType: string(pbmtypes.PbmProfileResourceTypeEnumSTORAGE),
	}

	category := pbmtypes.PbmProfileCategoryEnumREQUIREMENT

	pbmClient, err := pbm.NewClient(ctx, client)
	if err != nil {
		return err
	}

	ids, err := pbmClient.QueryProfile(ctx, rtype, string(category))
	if err != nil {
		return err
	}

	profiles, err := pbmClient.RetrieveContent(ctx, ids)
	if err != nil {
		return err
	}
	policyName := fmt.Sprintf("openshift-storage-policy-%s", infraID)
	policyLogger := logger.WithField("StoragePolicy", policyName)

	matchingProfileIds := []pbmtypes.PbmProfileId{}
	for _, p := range profiles {
		if p.GetPbmProfile().Name == policyName {
			profileID := p.GetPbmProfile().ProfileId
			matchingProfileIds = append(matchingProfileIds, profileID)
		}
	}
	if len(matchingProfileIds) > 0 {
		_, err = pbmClient.DeleteProfile(ctx, matchingProfileIds)
		if err != nil {
			return err
		}
		policyLogger.Info("Destroyed")

	}
	return nil
}

func deleteTag(ctx context.Context, client *rest.Client, tagID string) error {
	tagManager := tags.NewManager(client)
	tag, err := tagManager.GetTag(ctx, tagID)
	if err == nil {
		err = tagManager.DeleteTag(ctx, tag)
	}
	return err
}

func deleteTagCategory(ctx context.Context, client *rest.Client, categoryID string) error {
	tagManager := tags.NewManager(client)
	category, err := tagManager.GetCategory(ctx, categoryID)
	if err == nil {
		err = tagManager.DeleteCategory(ctx, category)
	}
	return err
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	var folderList []types.ManagedObjectReference
	var virtualMachineList []types.ManagedObjectReference

	o.Logger.Debug("Find attached objects on tag")
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

	if len(virtualMachineList) > 0 {
		o.Logger.Debug("Find VirtualMachine objects")
		virtualMachineMoList, err := getVirtualMachineManagedObjects(context.TODO(), o.Client, virtualMachineList)
		if err != nil {
			return err
		}
		o.Logger.Debug("Delete VirtualMachines")
		err = deleteVirtualMachines(context.TODO(), o.Client, virtualMachineMoList, o.Logger)
		if err != nil {
			return err
		}
	} else {
		o.Logger.Debug("No VirtualMachines found")
	}

	if len(folderList) > 0 {
		o.Logger.Debug("Find Folder objects")
		folderMoList, err := getFolderManagedObjects(context.TODO(), o.Client, folderList)
		if err != nil {
			o.Logger.Errorln(err)
			return err
		}

		o.Logger.Debug("Delete Folder")
		err = deleteFolder(context.TODO(), o.Client, folderMoList, o.Logger)
		if err != nil {
			o.Logger.Errorln(err)
			return err
		}
	} else {
		o.Logger.Debug("No managed Folder found")
	}

	err = deleteStoragePolicy(context.TODO(), o.Client, o.InfraID, o.Logger)
	if err != nil {
		return errors.Errorf("error deleting storage policy: %v", err)
	}

	o.Logger.Debug("Delete tag")
	tagLogger := o.Logger.WithField("Tag", o.InfraID)
	if err = deleteTag(context.TODO(), o.RestClient, o.InfraID); err != nil {
		tagLogger.Errorln(err)
		return err
	}
	tagLogger.Info("Destroyed")

	o.Logger.Debug("Delete tag category")
	tcLogger := o.Logger.WithField("TagCategory", "openshift-"+o.InfraID)
	if err = deleteTagCategory(context.TODO(), o.RestClient, "openshift-"+o.InfraID); err != nil {
		tcLogger.Errorln(err)
		return err
	}
	tcLogger.Info("Destroyed")

	return nil
}
