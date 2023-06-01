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

	Logger logrus.FieldLogger
	client API
}

// New returns an VSphere destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {
	client, err := NewClient(metadata.VSphere.VCenter, metadata.VSphere.Username, metadata.VSphere.Password)
	if err != nil {
		return nil, err
	}
	return newWithClient(logger, metadata, client), nil
}

func newWithClient(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata, client API) *ClusterUninstaller {
	return &ClusterUninstaller{
		ClusterID:         metadata.ClusterID,
		InfraID:           metadata.InfraID,
		terraformPlatform: metadata.VSphere.TerraformPlatform,

		Logger: logger,
		client: client,
	}
}

func (o *ClusterUninstaller) deleteFolder(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	o.Logger.Debug("Delete Folder")

	folderMoList, err := o.client.ListFolders(ctx, o.InfraID)
	if err != nil {
		return err
	}

	if len(folderMoList) == 0 {
		o.Logger.Debug("All folders deleted")
		return nil
	}

	// If there are no children in the folder, go ahead and remove it

	for _, f := range folderMoList {
		folderLogger := o.Logger.WithField("Folder", f.Name)
		if numChildren := len(f.ChildEntity); numChildren > 0 {
			entities := make([]string, 0, numChildren)
			for _, child := range f.ChildEntity {
				entities = append(entities, fmt.Sprintf("%s:%s", child.Type, child.Value))
			}
			folderLogger.Errorf("Folder should be empty but contains %d objects: %s. The installer will retry removing \"virtualmachine\" objects, but any other type will need to be removed manually before the deprovision can proceed", numChildren, strings.Join(entities, ", "))
			return errors.Errorf("Expected Folder %s to be empty", f.Name)
		}
		err = o.client.DeleteFolder(ctx, f)
		if err != nil {
			folderLogger.Debug(err)
			return err
		}
		folderLogger.Info("Destroyed")
	}

	return nil
}

func (o *ClusterUninstaller) deleteStoragePolicy(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	policyName := fmt.Sprintf("openshift-storage-policy-%s", o.InfraID)
	policyLogger := o.Logger.WithField("StoragePolicy", policyName)
	policyLogger.Debug("Delete")
	err := o.client.DeleteStoragePolicy(ctx, policyName)
	if err != nil {
		policyLogger.Debug(err)
		return err
	}
	policyLogger.Info("Destroyed")

	return nil
}

func (o *ClusterUninstaller) deleteTag(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	tagLogger := o.Logger.WithField("Tag", o.InfraID)
	tagLogger.Debug("Delete")
	err := o.client.DeleteTag(ctx, o.InfraID)
	if err != nil {
		tagLogger.Debug(err)
		return err
	}
	tagLogger.Info("Deleted")

	return nil
}

func (o *ClusterUninstaller) deleteTagCategory(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	categoryID := "openshift-" + o.InfraID
	tcLogger := o.Logger.WithField("TagCategory", categoryID)
	tcLogger.Debug("Delete")
	err := o.client.DeleteTagCategory(ctx, categoryID)
	if err != nil {
		tcLogger.Errorln(err)
		return err
	}
	tcLogger.Info("Deleted")

	return nil
}

func (o *ClusterUninstaller) stopVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name)
	err := o.client.StopVirtualMachine(ctx, vmMO)
	if err != nil {
		virtualMachineLogger.Debug(err)
		return err
	}
	virtualMachineLogger.Debug("Powered off")

	return nil
}

func (o *ClusterUninstaller) stopVirtualMachines(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	o.Logger.Debug("Power Off Virtual Machines")
	found, err := o.client.ListVirtualMachines(ctx, o.InfraID)
	if err != nil {
		o.Logger.Debug(err)
		return err
	}

	var errs []error
	for _, vmMO := range found {
		if !isPoweredOff(vmMO) {
			if err := o.stopVirtualMachine(ctx, vmMO); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) deleteVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name)
	err := o.client.DeleteVirtualMachine(ctx, vmMO)
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

	o.Logger.Debug("Delete Virtual Machines")
	found, err := o.client.ListVirtualMachines(ctx, o.InfraID)
	if err != nil {
		o.Logger.Debug(err)
		return err
	}

	var errs []error
	for _, vmMO := range found {
		if err := o.deleteVirtualMachine(ctx, vmMO); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) destroyCluster(ctx context.Context) (bool, error) {
	stagedFuncs := [][]struct {
		name    string
		execute func(context.Context) error
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
	defer o.client.Logout()

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
