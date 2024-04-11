package nutanix

import (
	"context"
	"fmt"
	"strings"
	"time"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

const (
	emptyFilter                = ""
	expectedCategoryKeyFormat  = "kubernetes-io-cluster-%s"
	expectedCategoryValueOwned = "owned"
)

// clusterUninstaller holds the various options for the cluster we want to delete.
type clusterUninstaller struct {
	clusterID string
	infraID   string
	v3Client  *nutanixclientv3.Client
	logger    logrus.FieldLogger
}

// New returns an Nutanix destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	v3Client, err := nutanixtypes.CreateNutanixClient(ctx,
		metadata.ClusterPlatformMetadata.Nutanix.PrismCentral,
		metadata.ClusterPlatformMetadata.Nutanix.Port,
		metadata.ClusterPlatformMetadata.Nutanix.Username,
		metadata.ClusterPlatformMetadata.Nutanix.Password,
	)
	if err != nil {
		return nil, err
	}

	return &clusterUninstaller{
		clusterID: metadata.ClusterID,
		infraID:   metadata.InfraID,
		v3Client:  v3Client,
		logger:    logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *clusterUninstaller) Run() (*installertypes.ClusterQuota, error) {
	o.logger.Infof("Starting deletion of Nutanix infrastructure for Openshift cluster %q", o.infraID)
	err := wait.PollImmediateInfinite(time.Second*30, o.destroyCluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}

func (o *clusterUninstaller) destroyCluster() (bool, error) {
	cleanupFuncs := []struct {
		name    string
		execute func(*clusterUninstaller) error
	}{
		{name: "VMs", execute: cleanupVMs},
		{name: "Images", execute: cleanupImages},
		{name: "Categories", execute: cleanupCategories},
	}

	done := true
	for _, cleanupFunc := range cleanupFuncs {
		if err := cleanupFunc.execute(o); err != nil {
			o.logger.Debugf("%s: %v", cleanupFunc.name, err)
			done = false
		}
	}
	return done, nil
}

func cleanupVMs(o *clusterUninstaller) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	matchedVirtualMachineList := make([]*nutanixclientv3.VMIntentResource, 0)
	allVMs, err := o.v3Client.V3.ListAllVM(ctx, emptyFilter)
	if err != nil {
		return err
	}

	for _, v := range allVMs.Entities {
		if hasCategoryOwned(v.Metadata, expectedCategoryKey(o.infraID)) {
			matchedVirtualMachineList = append(matchedVirtualMachineList, v)
		}
	}

	if len(matchedVirtualMachineList) == 0 {
		o.logger.Infof("No VMs found that require deletion for cluster %q", o.clusterID)
	} else {
		logToBeDeletedVMs(matchedVirtualMachineList, o.logger)
		err := deleteVMs(o.v3Client.V3, matchedVirtualMachineList, o.logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupImages(o *clusterUninstaller) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 120*time.Second)
	defer cancel()

	allImages, err := o.v3Client.V3.ListAllImage(ctx, emptyFilter)
	if err != nil {
		return err
	}

	imageDeletionFailed := false
	for _, image := range allImages.Entities {
		if hasCategoryOwned(image.Metadata, expectedCategoryKey(o.infraID)) {
			imageName := *image.Spec.Name
			imageUUID := *image.Metadata.UUID
			o.logger.Infof("Deleting image %q with UUID %q", imageName, imageUUID)
			response, err := o.v3Client.V3.DeleteImage(ctx, imageUUID)
			if err != nil {
				o.logger.Errorf("Failed to delete image %q: %v", imageUUID, err)
				imageDeletionFailed = true
				continue
			}

			if err := nutanixtypes.WaitForTask(o.v3Client.V3, response.Status.ExecutionContext.TaskUUID.(string)); err != nil {
				o.logger.Errorf("Failed to confirm image deletion %q: %v", imageUUID, err)
				imageDeletionFailed = true
			}
		}
	}

	if imageDeletionFailed {
		return fmt.Errorf("failed to cleanup images")
	}

	return nil
}

func cleanupCategories(o *clusterUninstaller) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	expCatKey := expectedCategoryKey(o.infraID)
	key, err := o.v3Client.V3.GetCategoryKey(ctx, expCatKey)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			//Already deleted
			return nil
		}
		return err
	}

	values, err := o.v3Client.V3.ListCategoryValues(context.TODO(), *key.Name, &nutanixclientv3.CategoryListMetadata{})
	if err != nil {
		return err
	}

	categoryDeletionFailed := false
	for _, value := range values.Entities {
		o.logger.Infof("Deleting category value : %s", *value.Value)
		err := o.v3Client.V3.DeleteCategoryValue(ctx, expCatKey, *value.Value)
		if err != nil {
			o.logger.Errorf("Failed to delete category value %q: %v", *value.Value, err)
			categoryDeletionFailed = true
		}
	}

	o.logger.Infof("Deleting category key : %s", expCatKey)
	err = o.v3Client.V3.DeleteCategoryKey(ctx, expCatKey)
	if err != nil {
		o.logger.Errorf("Failed to delete category key %q: %v", expCatKey, err)
		categoryDeletionFailed = true
	}

	if categoryDeletionFailed {
		return fmt.Errorf("failed to delete category")
	}

	return nil
}

func deleteVMs(clientV3 nutanixclientv3.Service, vms []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 120*time.Second)
	defer cancel()

	taskUUIDs := make([]string, 0)
	vmDeletionFailed := false
	for _, vm := range vms {
		l.Infof("Deleting VM %s with ID %s", *vm.Spec.Name, *vm.Metadata.UUID)
		response, err := clientV3.DeleteVM(ctx, *vm.Metadata.UUID)
		if err != nil {
			l.Errorf("Failed to delete VM %q: %v", *vm.Metadata.UUID, err)
			vmDeletionFailed = true
			continue
		}

		taskUUIDs = append(taskUUIDs, response.Status.ExecutionContext.TaskUUID.(string))
	}

	err := nutanixtypes.WaitForTasks(clientV3, taskUUIDs)
	if err != nil {
		l.Errorf("Failed to confirm deletion of VMs: %v", err)
		vmDeletionFailed = true
	}

	if vmDeletionFailed {
		return fmt.Errorf("failed to delete VMs")
	}

	return nil
}

func logToBeDeletedVMs(vms []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) {
	l.Info("Virtual machines scheduled to be deleted: ")
	for _, vm := range vms {
		l.Infof("- %s", *vm.Spec.Name)
	}
}

func expectedCategoryKey(infraID string) string {
	return fmt.Sprintf(expectedCategoryKeyFormat, infraID)
}

func hasCategoryOwned(metadata *nutanixclientv3.Metadata, expectedCategoryKey string) bool {
	value, keyExists := metadata.Categories[expectedCategoryKey]
	return keyExists && value == expectedCategoryValueOwned
}
