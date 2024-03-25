package clusterapi

import (
	"context"
	"fmt"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
)

const (
	// Category Key format: "kubernetes-io-cluster-<cluster-id>"
	CategoryKeyPrefix   = "kubernetes-io-cluster-"
	categoryValueOwned  = "owned"
	categoryValueShared = "shared"
)

// Provider is the Nutanix implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.PreProvider = Provider{}
var _ clusterapi.IgnitionProvider = Provider{}

// Name returns the Nutanix provider name.
func (p Provider) Name() string {
	return nutanixtypes.Name
}

// PreProvision creates the resources required prior to running capi nutanix controller.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	// create categories with name "kubernetes-io-cluster-<cluster_id>" and values ["owned", "shared"]
	// load the rhcos image to prism_central

	ic := in.InstallConfig.Config
	nutanixCl, err := nutanixtypes.CreateNutanixClientFromPlatform(ic.Platform.Nutanix)
	if err != nil {
		return fmt.Errorf("fail to create nutanix client. %w", err)
	}

	// create the category key
	categoryKey := fmt.Sprintf("%s%s", CategoryKeyPrefix, in.InfraID)
	keyBody := &nutanixclientv3.CategoryKey{
		APIVersion:  ptr.To("3.1.0"),
		Description: ptr.To("Openshift Cluster Category Key"),
		Name:        ptr.To(categoryKey),
	}
	respk, err := nutanixCl.V3.CreateOrUpdateCategoryKey(ctx, keyBody)
	if err != nil {
		return fmt.Errorf("failed to create the category key %q. %w", categoryKey, err)
	}
	logrus.Infof("Created the category key %q", *respk.Name)

	// create the category value "owned"
	valBody := &nutanixclientv3.CategoryValue{
		//APIVersion: ptr.To("3.1.0"),
		Description: ptr.To("Openshift Cluster Category Value: resources owned by the cluster"),
		Value:       ptr.To(categoryValueOwned),
	}
	respv, err := nutanixCl.V3.CreateOrUpdateCategoryValue(ctx, categoryKey, valBody)
	if err != nil {
		return fmt.Errorf("failed to create the category value %q with category key %q. %w", categoryValueOwned, categoryKey, err)
	}
	logrus.Infof("Created the category value %q with name %q", *respv.Value, *respv.Name)

	// create the category value "shared"
	valBody = &nutanixclientv3.CategoryValue{
		//APIVersion: ptr.To("3.1.0"),
		Description: ptr.To("Openshift Cluster Category Value: resources used but not owned by the cluster"),
		Value:       ptr.To(categoryValueShared),
	}
	respv, err = nutanixCl.V3.CreateOrUpdateCategoryValue(ctx, categoryKey, valBody)
	if err != nil {
		return fmt.Errorf("failed to create the category value %q with category key %q. %w", categoryValueShared, categoryKey, err)
	}
	logrus.Infof("Created the category value %q with name %q", *respv.Value, *respv.Name)

	// upload the rhcos image
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
		Categories: map[string]string{categoryKey: categoryValueOwned},
	}
	imgReq.Metadata = imgMeta
	respi, err := nutanixCl.V3.CreateImage(ctx, imgReq)
	if err != nil {
		return fmt.Errorf("failed to create the rhcos image %q. %w", imgName, err)
	}
	imgUUID := *respi.Metadata.UUID
	logrus.Debugf("Created the rhcos image %s (uuid: %s), waiting the image data uploading ...", imgName, imgUUID)

	// wait until the image loading completes
	err = nutanixtypes.WaitForImageStateComplete(nutanixCl, imgUUID)
	if err != nil {
		// FIXME: in case of error, delete the image with the uuid
		return fmt.Errorf("failed to upload the rhcos image %q to prism_central. %w", imgName, err)
	}
	logrus.Debugf("Completed uploading the rhcos image %q", imgName)

	return nil
}

// Ignition handles preconditions for bootstrap ignition and
// generates ignition data for the CAPI bootstrap ignition secret.
// Load the ignition iso image to prism_central
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	ic := in.InstallConfig.Config
	nutanixCl, err := nutanixtypes.CreateNutanixClientFromPlatform(ic.Platform.Nutanix)
	if err != nil {
		return nil, fmt.Errorf("fail to create nutanix client. %w", err)
	}

	imgName := nutanixtypes.BootISOImageName(in.InfraID)
	imgPath, err := nutanixtypes.CreateBootstrapISO(in.InfraID, string(in.BootstrapIgnData))
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap ignition iso. %w", err)
	}

	// upload the bootstrap image
	imgReq := &nutanixclientv3.ImageIntentInput{}
	imgSpec := &nutanixclientv3.Image{
		Name:        &imgName,
		Description: ptr.To("Created By OpenShift Installer"),
		Resources: &nutanixclientv3.ImageResources{
			ImageType: ptr.To("ISO_IMAGE"),
		},
	}
	imgReq.Spec = imgSpec
	categoryKey := fmt.Sprintf("%s%s", CategoryKeyPrefix, in.InfraID)
	imgMeta := &nutanixclientv3.Metadata{
		Kind:       ptr.To("image"),
		Categories: map[string]string{categoryKey: categoryValueOwned},
	}
	imgReq.Metadata = imgMeta

	// create the image object
	respi, err := nutanixCl.V3.CreateImage(ctx, imgReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create the bootstrap image %q. %w", imgName, err)
	}
	imgUUID := *respi.Metadata.UUID
	logrus.Debugf("Created the bootstrap image %s (uuid: %s), waiting the image data uploading ...", imgName, imgUUID)

	// upload the image data
	err = nutanixCl.V3.UploadImage(context.TODO(), imgUUID, imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload the bootstrap image data %q. %w", imgName, err)
	}

	// wait until the image loading completes
	err = nutanixtypes.WaitForImageStateComplete(nutanixCl, imgUUID)
	if err != nil {
		// FIXME: in case of error, delete the image with the uuid
		return nil, fmt.Errorf("failed to upload the bootstrap image %q to prism_central. %w", imgName, err)
	}
	logrus.Debugf("Completed uploading the bootstrap image %q", imgName)

	return in.BootstrapIgnData, nil
}
