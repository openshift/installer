package clusterapi

import (
	"context"
	"fmt"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/sirupsen/logrus"
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
	imgURI := string(*in.RhcosImage)
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

	// create the image object.
	respi, err := nutanixCl.V3.CreateImage(ctx, imgReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create the bootstrap image %q: %w", imgName, err)
	}
	imgUUID := *respi.Metadata.UUID

	if taskUUID, ok := respi.Status.ExecutionContext.TaskUUID.(string); ok {
		logrus.Infof("creating the bootstrap image %s (uuid: %s), taskUUID: %s.", imgName, imgUUID, taskUUID)

		// Wait till the image creation task is successed.
		if err = nutanixtypes.WaitForTask(nutanixCl.V3, taskUUID); err != nil {
			err = fmt.Errorf("failed to create the bootstrap image %q: %w", imgName, err)
			logrus.Errorf(err.Error())
			return nil, err
		}
		logrus.Infof("created the bootstrap image %s (uuid: %s).", imgName, imgUUID)
	} else {
		err = fmt.Errorf("failed to convert the task UUID %v to string", respi.Status.ExecutionContext.TaskUUID)
		logrus.Error(err)
		return nil, err
	}

	// upload the image data.
	logrus.Infof("preparing to upload the bootstrap image %s (uuid: %s) data from file %s", imgName, imgUUID, imgPath)
	err = nutanixCl.V3.UploadImage(ctx, imgUUID, imgPath)
	if err != nil {
		e1 := fmt.Errorf("failed to upload the bootstrap image data %q from filepath %s: %w", imgName, imgPath, err)
		logrus.Error(e1)
		return nil, e1
	}
	logrus.Infof("uploading the bootstrap image %s data", imgName)
	// wait for the image data uploading task to complete.
	respb, err := nutanixCl.V3.GetImage(ctx, imgUUID)
	if err != nil {
		e1 := fmt.Errorf("failed to get the bootstrap image %q. %w", imgName, err)
		logrus.Error(e1)
		return nil, e1
	}

	if taskUUIDs, ok := respb.Status.ExecutionContext.TaskUUID.([]interface{}); ok {
		tUUIDs := []string{}
		for _, tUUID := range taskUUIDs {
			if tUUIDstr, ok := tUUID.(string); ok {
				tUUIDs = append(tUUIDs, tUUIDstr)
			}
		}
		logrus.Infof("waiting for the bootstrap image data uploading task to complete,  taskUUIDs: %v", tUUIDs)
		if err = nutanixtypes.WaitForTasks(nutanixCl.V3, tUUIDs); err != nil {
			e1 := fmt.Errorf("failed to upload the bootstrap image data %q from filepath %s: %w", imgName, imgPath, err)
			logrus.Error(e1)
			return nil, e1
		}
		logrus.Infof("completed uploading the bootstrap image data %s (uuid: %s)", imgName, imgUUID)
	} else {
		err = fmt.Errorf("failed to convert the taskUUIDs %v to array", respb.Status.ExecutionContext.TaskUUID)
		logrus.Error(err)
		return nil, err
	}

	return in.BootstrapIgnData, nil
}
