package vsphere

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/client/vsphere"
	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
)

var defaultTimeout = time.Minute * 5

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	ClusterID string
	InfraID   string
	vCenter   string
	username  string
	password  string

	Client     *vim25.Client
	RestClient *rest.Client

	Logger logrus.FieldLogger

	context context.Context
}

// New returns an VSphere destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		ClusterID: metadata.ClusterID,
		InfraID:   metadata.InfraID,
		vCenter:   metadata.VSphere.VCenter,
		username:  metadata.VSphere.Username,
		password:  metadata.VSphere.Password,
		Logger:    logger,
		context:   context.Background(),
	}, nil
}

func isNotFound(err error) bool {
	return err != nil && strings.HasSuffix(err.Error(), http.StatusText(http.StatusNotFound))
}

func (o *ClusterUninstaller) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(o.context, defaultTimeout)
}

func (o *ClusterUninstaller) getAttachedObjectsOnTag(objType string) ([]types.ManagedObjectReference, error) {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Find attached %s on tag", objType)
	tagManager := tags.NewManager(o.RestClient)
	attached, err := tagManager.GetAttachedObjectsOnTags(ctx, []string{o.InfraID})
	if err != nil && !isNotFound(err) {
		return nil, err
	}

	// Separate the objects attached to the tag based on type
	var objectList []types.ManagedObjectReference
	for _, attachedObject := range attached {
		for _, ref := range attachedObject.ObjectIDs {
			if ref.Reference().Type == objType {
				objectList = append(objectList, ref.Reference())
			}
		}
	}

	return objectList, nil
}

func (o *ClusterUninstaller) getFolderManagedObjects(moRef []types.ManagedObjectReference) ([]mo.Folder, error) {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	var folderMoList []mo.Folder
	if len(moRef) > 0 {
		pc := property.DefaultCollector(o.Client)
		err := pc.Retrieve(ctx, moRef, nil, &folderMoList)
		if err != nil {
			return nil, err
		}
	}
	return folderMoList, nil
}

func (o *ClusterUninstaller) listFolders() ([]mo.Folder, error) {
	folderList, err := o.getAttachedObjectsOnTag("Folder")
	if err != nil {
		return nil, err
	}

	return o.getFolderManagedObjects(folderList)
}

func (o *ClusterUninstaller) deleteFolder() error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debug("Delete Folder")

	folderMoList, err := o.listFolders()
	if err != nil {
		return err
	}

	// The installer should create at most one parent,
	// the parent to the VirtualMachines.
	// If there are more or less fail with error message.
	if len(folderMoList) > 1 {
		return errors.Errorf("Expected 1 Folder per tag but got %d", len(folderMoList))
	}

	if len(folderMoList) == 0 {
		o.Logger.Debug("All folders deleted")
		return nil
	}

	// If there are no children in the folder, go ahead and remove it
	if len(folderMoList[0].ChildEntity) == 0 {
		folderLogger := o.Logger.WithField("Folder", folderMoList[0].Name)

		folder := object.NewFolder(o.Client, folderMoList[0].Reference())
		task, err := folder.Destroy(ctx)
		if err == nil {
			err = task.Wait(ctx)
		}
		if err != nil {
			folderLogger.Debug(err)
			return err
		}
		folderLogger.Info("Destroyed")
	} else {
		return errors.Errorf("Expected Folder %s to be empty", folderMoList[0].Name)
	}

	return nil
}

func (o *ClusterUninstaller) deleteStoragePolicy() error {
	ctx, cancel := context.WithTimeout(o.context, time.Minute*30)
	defer cancel()

	o.Logger.Debug("Delete Storage Policy")
	rtype := pbmtypes.PbmProfileResourceType{
		ResourceType: string(pbmtypes.PbmProfileResourceTypeEnumSTORAGE),
	}

	category := pbmtypes.PbmProfileCategoryEnumREQUIREMENT

	pbmClient, err := pbm.NewClient(ctx, o.Client)
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
	policyName := fmt.Sprintf("openshift-storage-policy-%s", o.InfraID)
	policyLogger := o.Logger.WithField("StoragePolicy", policyName)

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

func (o *ClusterUninstaller) deleteTag() error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	tagLogger := o.Logger.WithField("Tag", o.InfraID)
	tagLogger.Debug("Delete")

	tagManager := tags.NewManager(o.RestClient)
	tag, err := tagManager.GetTag(ctx, o.InfraID)
	if err == nil {
		err = tagManager.DeleteTag(ctx, tag)
		if err == nil {
			tagLogger.Info("Deleted")
		}
	}
	if isNotFound(err) {
		return nil
	}
	return err
}

func (o *ClusterUninstaller) deleteTagCategory() error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	categoryID := "openshift-" + o.InfraID
	tcLogger := o.Logger.WithField("TagCategory", categoryID)
	tcLogger.Debug("Delete")

	tagManager := tags.NewManager(o.RestClient)
	ids, err := tagManager.ListCategories(ctx)
	if err != nil {
		tcLogger.Errorln(err)
		return err
	}

	var errs []error
	for _, id := range ids {
		category, err := tagManager.GetCategory(ctx, id)
		if err != nil {
			if !isNotFound(err) {
				errs = append(errs, errors.Wrapf(err, "could not get category %q", id))
			}
			continue
		}
		if category.Name == categoryID {
			if err = tagManager.DeleteCategory(ctx, category); err != nil {
				tcLogger.Errorln(err)
				return err
			}
			tcLogger.Info("Deleted")
			return nil
		}
	}

	if len(errs) == 0 {
		tcLogger.Debug("Not found")
	}
	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) destroyCluster() (bool, error) {
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Stop virtual machines", execute: o.stopVirtualMachines},
	}, {
		{name: "Virtual Machines", execute: o.deleteVirtualMachines},
	}, {
		{name: "Folder", execute: o.deleteFolder},
	}, {
		{name: "Storage Policy", execute: o.deleteStoragePolicy},
		{name: "Tag", execute: o.deleteTag},
		{name: "Tag Category", execute: o.deleteTagCategory},
	}}

	stageFailed := false
	for _, stage := range stagedFuncs {
		if stageFailed {
			break
		}
		for _, f := range stage {
			err := f.execute()
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
	vim25Client, restClient, cleanup, err := vsphere.CreateVSphereClients(context.TODO(),
		o.vCenter,
		o.username,
		o.password)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	o.Client = vim25Client
	o.RestClient = restClient

	err = wait.PollImmediateInfinite(
		time.Second*10,
		o.destroyCluster,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}
