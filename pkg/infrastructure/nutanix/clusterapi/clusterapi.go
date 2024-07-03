package clusterapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"

	infracapi "github.com/openshift/installer/pkg/infrastructure/clusterapi"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

// Provider is the Nutanix implementation of the clusterapi InfraProvider.
type Provider struct{}

var _ infracapi.PreProvider = Provider{}
var _ infracapi.IgnitionProvider = Provider{}

// Name returns the Nutanix provider name.
func (p Provider) Name() string {
	return nutanixtypes.Name
}

// BootstrapHasPublicIP indicates that an ExternalIP is not
// required in the machine ready checks.
func (Provider) BootstrapHasPublicIP() bool { return false }

// PreProvision creates the resources required prior to running capi nutanix controller.
func (p Provider) PreProvision(ctx context.Context, in infracapi.PreProvisionInput) error {
	// create categories with name "kubernetes-io-cluster-<cluster_id>" and values ["owned", "shared"].
	// load the rhcos image to prism_central.

	ic := in.InstallConfig.Config
	nutanixCl, err := nutanixtypes.CreateNutanixClientFromPlatform(ic.Platform.Nutanix)
	if err != nil {
		return fmt.Errorf("failed to create nutanix client: %w", err)
	}

	// create the category key.
	categoryKey := nutanixtypes.CategoryKey(in.InfraID)
	keyBody := &nutanixclientv3.CategoryKey{
		APIVersion:  ptr.To("3.1.0"),
		Description: ptr.To("Openshift Cluster Category Key"),
		Name:        ptr.To(categoryKey),
	}
	respk, err := nutanixCl.V3.CreateOrUpdateCategoryKey(ctx, keyBody)
	if err != nil {
		return fmt.Errorf("failed to create the category key %q: %w", categoryKey, err)
	}
	logrus.Infof("created the category key %q", *respk.Name)

	// create the category value "owned".
	valBody := &nutanixclientv3.CategoryValue{
		Description: ptr.To("Openshift Cluster Category Value: resources owned by the cluster"),
		Value:       ptr.To(nutanixtypes.CategoryValueOwned),
	}
	respv, err := nutanixCl.V3.CreateOrUpdateCategoryValue(ctx, categoryKey, valBody)
	if err != nil {
		return fmt.Errorf("failed to create the category value %q with category key %q: %w", nutanixtypes.CategoryValueOwned, categoryKey, err)
	}
	logrus.Infof("created the category value %q with name %q", *respv.Value, *respv.Name)

	// create the category value "shared".
	valBody = &nutanixclientv3.CategoryValue{
		Description: ptr.To("Openshift Cluster Category Value: resources used but not owned by the cluster"),
		Value:       ptr.To(nutanixtypes.CategoryValueShared),
	}
	respv, err = nutanixCl.V3.CreateOrUpdateCategoryValue(ctx, categoryKey, valBody)
	if err != nil {
		return fmt.Errorf("failed to create the category value %q with category key %q: %w", nutanixtypes.CategoryValueShared, categoryKey, err)
	}
	logrus.Infof("created the category value %q with name %q", *respv.Value, *respv.Name)

	// upload the rhcos image.
	imgName := nutanixtypes.RHCOSImageName(in.InfraID)
	imgURI := in.RhcosImage.ControlPlane
	imgReq := &nutanixclientv3.ImageIntentInput{}
	imgSpec := &nutanixclientv3.Image{
		Name:        &imgName,
		Description: ptr.To("Created By OpenShift Installer"),
		Resources: &nutanixclientv3.ImageResources{
			ImageType: ptr.To("DISK_IMAGE"),
			SourceURI: &imgURI,
		},
	}
	imgReq.Spec = imgSpec
	imgMeta := &nutanixclientv3.Metadata{
		Kind:       ptr.To("image"),
		Categories: map[string]string{categoryKey: nutanixtypes.CategoryValueOwned},
	}
	imgReq.Metadata = imgMeta
	respi, err := nutanixCl.V3.CreateImage(ctx, imgReq)
	if err != nil {
		return fmt.Errorf("failed to create the rhcos image %q: %w", imgName, err)
	}
	imgUUID := *respi.Metadata.UUID
	logrus.Infof("creating the rhcos image %s (uuid: %s).", imgName, imgUUID)

	if taskUUID, ok := respi.Status.ExecutionContext.TaskUUID.(string); ok {
		logrus.Infof("waiting the image data uploading from %s, taskUUID: %s.", imgURI, taskUUID)

		// Wait till the image creation task is successed.
		if err = nutanixtypes.WaitForTask(nutanixCl.V3, taskUUID); err != nil {
			e1 := fmt.Errorf("failed to create the rhcos image %q: %w", imgName, err)
			logrus.Error(e1)
			return e1
		}
		logrus.Infof("created and uploaded the rhcos image data %s (uuid: %s)", imgName, imgUUID)
	} else {
		err = fmt.Errorf("failed to convert the task UUID %v to string", respi.Status.ExecutionContext.TaskUUID)
		logrus.Errorf(err.Error())
		return err
	}

	return nil
}

// Ignition handles preconditions for bootstrap ignition and
// generates ignition data for the CAPI bootstrap ignition secret.
// Load the ignition iso image to prism_central.
func (p Provider) Ignition(ctx context.Context, in infracapi.IgnitionInput) ([]byte, error) {
	ic := in.InstallConfig.Config
	nutanixCl, err := nutanixtypes.CreateNutanixClientFromPlatform(ic.Platform.Nutanix)
	if err != nil {
		return nil, fmt.Errorf("failed to create nutanix client: %w", err)
	}

	imgName := nutanixtypes.BootISOImageName(in.InfraID)
	imgPath, err := nutanixtypes.CreateBootstrapISO(in.InfraID, string(in.BootstrapIgnData))
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap ignition iso: %w", err)
	}

	// upload the bootstrap image.
	imgReq := &nutanixclientv3.ImageIntentInput{}
	imgSpec := &nutanixclientv3.Image{
		Name:        &imgName,
		Description: ptr.To("Created By OpenShift Installer"),
		Resources: &nutanixclientv3.ImageResources{
			ImageType: ptr.To("ISO_IMAGE"),
		},
	}
	imgReq.Spec = imgSpec
	categoryKey := nutanixtypes.CategoryKey(in.InfraID)
	imgMeta := &nutanixclientv3.Metadata{
		Kind:       ptr.To("image"),
		Categories: map[string]string{categoryKey: nutanixtypes.CategoryValueOwned},
	}
	imgReq.Metadata = imgMeta

	// Wait for successful creation of the bootstrap image object and upload the image data to it in PC.
	// Put both createImage() and uploadImageData() in the same wait.ExponentialBackoffWithContext().
	// Because if createImage() succeeds but uploadImageData() fails, we need to delete the image object
	// and retry to call both createImage() and uploadImageData() again. The old-version prism-api server sometimes
	// returns error for the uploadImage call and does not allow to retry the uploadImage call to the same image object.
	timeout := 20 * time.Minute
	if err = wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
		Duration: time.Minute * 4,
		Factor:   float64(1.0),
		Steps:    5,
		Cap:      timeout,
	}, func(ctx context.Context) (bool, error) {
		// create the bootstrap image object in PC
		imgUUID, err1 := createImage(ctx, nutanixCl, imgReq, imgName)
		if err1 != nil {
			logrus.Errorf("failed to create the bootstrap image object %s in PC: %v", imgName, err1)
			// no need to retry if the error code is 401 or 403
			if strings.Contains(err1.Error(), `"code": 401`) || strings.Contains(err1.Error(), `"code": 403`) {
				return false, err1
			}

			// delete the image object if uuid is not empty
			if imgUUID != "" {
				if e2 := deleteImage(ctx, nutanixCl, imgUUID); e2 != nil {
					logrus.Errorf("failed to delete image object %s (uuid: %s): %v", imgName, imgUUID, e2)
				}
			}
			return false, nil
		}

		// upload the image data to the bootstrap image object in PC
		err2 := uploadImageData(ctx, nutanixCl, imgName, imgUUID, imgPath)
		if err2 != nil {
			logrus.Errorf("failed to upload the bootstrap image %s data: %v", imgName, err2)
			// no need to retry if the error code is 401 or 403
			if strings.Contains(err2.Error(), `"code": 401`) || strings.Contains(err2.Error(), `"code": 403`) {
				return false, err2
			}

			// delete the image object
			if e2 := deleteImage(ctx, nutanixCl, imgUUID); e2 != nil {
				logrus.Errorf("failed to delete image object %s (uuid: %s): %v", imgName, imgUUID, e2)
			}
			return false, nil
		}

		return true, nil
	}); err != nil {
		if wait.Interrupted(err) {
			err = fmt.Errorf("timeout/interrupt to create/upload the bootstrap image object %s in PC within %v: %w", imgName, timeout, err)
		} else {
			err = fmt.Errorf("failed to create/upload the bootstrap image object %s in PC: %w", imgName, err)
		}

		return in.BootstrapIgnData, err
	}
	logrus.Infof("Successfully created the bootstrap image object %s and uploaded its image data", imgName)

	return in.BootstrapIgnData, nil
}

// createImage creates the image object in PC, with the provided request input.
// Returns the imageUUID if the image is created.
func createImage(ctx context.Context, nutanixCl *nutanixclientv3.Client, imgReq *nutanixclientv3.ImageIntentInput, imgName string) (string, error) {
	t1 := time.Now()

	// create the image object.
	respi, err := nutanixCl.V3.CreateImage(ctx, imgReq)
	if err != nil {
		return "", fmt.Errorf("failed to create the image %q: %w", imgName, err)
	}
	imgUUID := *respi.Metadata.UUID

	if taskUUID, ok := respi.Status.ExecutionContext.TaskUUID.(string); ok {
		logrus.Infof("creating the image %s (uuid: %s), taskUUID: %s", imgName, imgUUID, taskUUID)

		// Wait for the image creation task
		if err = nutanixtypes.WaitForTask(nutanixCl.V3, taskUUID); err != nil {
			err = fmt.Errorf("failed to create the image %s (uuid: %s), taskUUID: %s: %w", imgName, imgUUID, taskUUID, err)
		} else {
			logrus.Infof("created the image %s (uuid: %s). used_time %v", imgName, imgUUID, time.Since(t1))
		}
	} else {
		err = fmt.Errorf("failed to convert the task UUID %v to string", respi.Status.ExecutionContext.TaskUUID)
	}

	return imgUUID, err
}

// uploadImageData upload the image data from the specified file path to the image object in PC.
func uploadImageData(ctx context.Context, nutanixCl *nutanixclientv3.Client, imgName, imgUUID, imgPath string) error {
	// upload the image data.
	logrus.Infof("preparing to upload the image %s (uuid: %s) data from file %s", imgName, imgUUID, imgPath)
	t1 := time.Now()
	err := nutanixCl.V3.UploadImage(ctx, imgUUID, imgPath)
	if err != nil {
		return fmt.Errorf("failed to upload the image data %q from filepath %s: %w  used_time %v", imgName, imgPath, err, time.Since(t1))
	}
	logrus.Infof("uploading the image %s data. used_time %v", imgName, time.Since(t1))

	// wait for the image data uploading task to complete.
	respb, err := nutanixCl.V3.GetImage(ctx, imgUUID)
	if err != nil {
		return fmt.Errorf("failed to get the image %q. %w", imgName, err)
	}

	if taskUUIDs, ok := respb.Status.ExecutionContext.TaskUUID.([]interface{}); ok {
		tUUIDs := []string{}
		for _, tUUID := range taskUUIDs {
			if tUUIDstr, ok := tUUID.(string); ok {
				tUUIDs = append(tUUIDs, tUUIDstr)
			}
		}
		logrus.Infof("waiting for the image data uploading task to complete,  taskUUIDs: %v", tUUIDs)
		if err = nutanixtypes.WaitForTasks(nutanixCl.V3, tUUIDs); err != nil {
			return fmt.Errorf("failed to upload the bootstrap image data %q from filepath %s: %w", imgName, imgPath, err)
		}
	} else {
		return fmt.Errorf("failed to convert the taskUUIDs %v to array", respb.Status.ExecutionContext.TaskUUID)
	}

	return nil
}

// deleteImage deletes the image object with the given uuid in PC.
func deleteImage(ctx context.Context, nutanixCl *nutanixclientv3.Client, imgUUID string) error {
	logrus.Infof("preparing to delete the image with uuid %s", imgUUID)

	respd, err := nutanixCl.V3.DeleteImage(ctx, imgUUID)
	if err != nil {
		return fmt.Errorf("failed to delete the image with uuid %s: %w", imgUUID, err)
	}

	if taskUUID, ok := respd.Status.ExecutionContext.TaskUUID.(string); ok {
		logrus.Infof("deleting the image with uuid %s, taskUUID: %s", imgUUID, taskUUID)

		// Wait till the image deletion task is successed.
		if err = nutanixtypes.WaitForTask(nutanixCl.V3, taskUUID); err != nil {
			return fmt.Errorf("failed to delete the image with uuid: %s, taskUUID: %s: %w", imgUUID, taskUUID, err)
		}
	} else {
		return fmt.Errorf("failed to convert the task UUID %v to string", respd.Status.ExecutionContext.TaskUUID)
	}

	return nil
}
