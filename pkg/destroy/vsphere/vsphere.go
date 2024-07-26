package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/mo"
	errorsutil "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	ClusterID         string
	InfraID           string
	terraformPlatform string

	Logger        logrus.FieldLogger
	client        API
	KubeClientset *kubernetes.Clientset

	deleteVolumes bool
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
	clusterUninstaller := &ClusterUninstaller{
		ClusterID:         metadata.ClusterID,
		InfraID:           metadata.InfraID,
		terraformPlatform: metadata.VSphere.TerraformPlatform,

		Logger:        logger,
		client:        client,
		KubeClientset: nil,
		deleteVolumes: metadata.DeleteVolumes,
	}

	if metadata.DeleteVolumes {
		config, err := clientcmd.RESTConfigFromKubeConfig(*metadata.Auth)
		if err == nil {
			clientset, err := kubernetes.NewForConfig(config)

			if clientset != nil && err == nil {
				clusterUninstaller.KubeClientset = clientset
			}
		}
	}
	return clusterUninstaller
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

type jsonPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// isTransientConnectionError checks whether given error is "Connection refused" or
// "Connection reset" error which usually means that apiserver is temporarily
// unavailable.
func isTransientConnectionError(err error) bool {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.ECONNREFUSED) || errors.Is(errno, syscall.ECONNRESET)
	}
	return false
}

func isTransientError(err error) bool {
	if isTransientConnectionError(err) {
		return true
	}

	if t, ok := err.(errorsutil.APIStatus); ok && t.Status().Code >= 500 {
		return true
	}

	return errorsutil.IsTooManyRequests(err)
}

func (o *ClusterUninstaller) deletePersistentVolumes(ctx context.Context) error {
	if o.KubeClientset == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()
	removeFinalizer := []jsonPatch{{
		Op:    "remove",
		Path:  "/metadata/finalizers",
		Value: "kubernetes",
	}}

	// todo: maybe the errors should be checked for cluster not available...

	removeFinalizerBytes, err := json.Marshal(removeFinalizer)
	if err != nil {
		return err
	}

	pvList, err := o.KubeClientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		o.Logger.Warnf("Unable to list persistent volumes: %v", err)
		return nil
	}

	if len(pvList.Items) == 0 {
		return nil
	}

	for _, pv := range pvList.Items {
		o.Logger.Debugf("deleting volume %s", pv.Name)
		err = o.KubeClientset.CoreV1().PersistentVolumes().Delete(ctx, pv.Name, metav1.DeleteOptions{})

		if err != nil {
			o.Logger.Warnf("Unable to delete persistent volumes: %v", err)
			return nil
		}

		_, err := o.KubeClientset.CoreV1().PersistentVolumes().Patch(ctx, pv.Name, types.JSONPatchType, removeFinalizerBytes, metav1.PatchOptions{})
		if err != nil {
			o.Logger.Warnf("Unable to patch persistent volumes: %v", err)
			return nil
		}
	}

	// todo: is there another way to know if a PV once deleted has been reconciled?
	time.Sleep(time.Second * 30)

	for {
		pvList, err := o.KubeClientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
		if err != nil {
			o.Logger.Warnf("Unable to list persistent volumes: %v", err)
			return nil
		}

		o.Logger.Debugf("%d remaining persistent volumes", len(pvList.Items))

		if len(pvList.Items) == 0 {
			break
		}

		time.Sleep(time.Second * 30)
	}

	return nil
}

func (o *ClusterUninstaller) deleteCnsVolumes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()
	o.Logger.Debug("Delete CNS Volumes")

	cnsVolumes, err := o.client.GetCnsVolumes(ctx, o.InfraID)
	if err != nil {
		return err
	}

	for _, cv := range cnsVolumes {
		cnsVolumeLogger := o.Logger.WithField("CNS Volume", cv.VolumeId.Id)
		err := o.client.DeleteCnsVolumes(ctx, cv)
		if err != nil {
			return err
		}
		cnsVolumeLogger.Info("Destroyed")
	}

	return nil
}

type StagedFunctions struct {
	Name    string
	Execute func(ctx context.Context) error
}

func (o *ClusterUninstaller) destroyCluster(ctx context.Context) (bool, error) {
	var stagedFuncs [][]StagedFunctions

	if o.deleteVolumes {
		var deleteVolumeStagedFunctions []StagedFunctions

		if o.KubeClientset != nil {
			deleteVolumeStagedFunctions = append(deleteVolumeStagedFunctions, StagedFunctions{
				Name:    "Delete Persistent Volumes",
				Execute: o.deletePersistentVolumes,
			})
		}
		deleteVolumeStagedFunctions = append(deleteVolumeStagedFunctions, StagedFunctions{
			Name:    "Delete CNS Volumes",
			Execute: o.deleteCnsVolumes,
		})
		stagedFuncs = append(stagedFuncs, deleteVolumeStagedFunctions)
	}

	deleteVirtualMachinesFuncs := []StagedFunctions{
		{
			Name: "Stop virtual machines", Execute: o.stopVirtualMachines,
		},
		{
			Name: "Delete Virtual Machines", Execute: o.deleteVirtualMachines,
		},
	}

	deleteVCenterObjectsFuncs := []StagedFunctions{
		{
			Name: "Folder", Execute: o.deleteFolder,
		},
		{
			Name: "Storage Policy", Execute: o.deleteStoragePolicy,
		},
		{
			Name: "Tag", Execute: o.deleteTag,
		},
		{
			Name: "Tag Category", Execute: o.deleteTagCategory,
		},
	}

	stagedFuncs = append(stagedFuncs, deleteVirtualMachinesFuncs)
	stagedFuncs = append(stagedFuncs, deleteVCenterObjectsFuncs)

	for _, sf := range stagedFuncs {
		for _, f := range sf {
			fmt.Print(f.Name)
		}
	}

	stageFailed := false
	for _, stage := range stagedFuncs {
		if stageFailed {
			break
		}
		for _, f := range stage {
			err := f.Execute(ctx)
			if err != nil {
				o.Logger.Debugf("%s: %v", f.Name, err)
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
