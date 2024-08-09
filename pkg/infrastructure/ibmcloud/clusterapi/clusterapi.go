package clusterapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos/cache"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements IBM Cloud CAPI installation.
type Provider struct{}

// Name returns the IBM Cloud provider name.
func (p Provider) Name() string {
	return ibmcloudtypes.Name
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status when declaring machines ready.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.InternalIP }

// PreProvision creates the IBM Cloud objects required prior to running capibmcloud.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	// Before Provisioning IBM Cloud Infrastructure for the Cluster, we must perform the following.
	// 1. Create the Resource Group to house cluster resources, if necessary (BYO RG).
	// 2. Create a COS Instance and Bucket to host the RHCOS Custom Image file.
	// 3. Upload the RHCOS image to the COS Bucket.
	// 4. Add IAM Authorization for VPC Image Service to access the COS Object/Bucket/Instance.

	// Setup IBM Cloud Client.
	metadata := ibmcloudic.NewMetadata(in.InstallConfig.Config)
	client, err := metadata.Client()
	if err != nil {
		return fmt.Errorf("failed creating IBM Cloud client: %w", err)
	}
	region := in.InstallConfig.Config.Platform.IBMCloud.Region

	// Create cluster's Resource Group, if necessary (BYO RG is supported).
	resourceGroupName := in.InfraID
	if in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName != "" {
		resourceGroupName = in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName
	}

	logrus.Debugf("checking for existing resource group: %s", resourceGroupName)
	// Check whether the Resource Group already exists.
	resourceGroup, err := client.GetResourceGroup(ctx, resourceGroupName)
	if err != nil {
		// If Resource Group cannot be found, but it was provided in install-config (use existing RG), raise an error.
		// We could create the Resource Group, defined by user, but that might make resource cleanup more difficult.
		if in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName != "" {
			return fmt.Errorf("provided resource group not found: %w", err)
		}
	}

	// Create Resource Group if it wasn't found (and isn't expected to be an existing RG).
	if resourceGroup == nil {
		logrus.Debugf("creating resource group: %s", resourceGroupName)
		if err := client.CreateResourceGroup(ctx, resourceGroupName); err != nil {
			return fmt.Errorf("failed creating new resource group: %w", err)
		}
		// Retrieve the newly created resource group
		resourceGroup, err = client.GetResourceGroup(ctx, resourceGroupName)
		if err != nil {
			return fmt.Errorf("failed retrieving new resource group: %w", err)
		}
		logrus.Debugf("created resource group: %s", resourceGroupName)
	}

	// Create a COS Instance and Bucket to host the RHCOS image file.
	// NOTE(cjschaef): Support to use an existing COS Object (RHCO image file) or VPC Custom Image could be added to skip this step.
	cosInstanceName := fmt.Sprintf("%s-cos", in.InfraID)
	logrus.Debugf("checking for existing cos instance: %s", cosInstanceName)
	cosInstance, err := client.GetCOSInstanceByName(ctx, cosInstanceName)
	if err != nil {
		logrus.Debugf("creating cos instance: %s", cosInstanceName)
		cosInstance, err = client.CreateCOSInstance(ctx, cosInstanceName, *resourceGroup.ID)
		if err != nil {
			return fmt.Errorf("failed creating RHCOS image COS instance: %w", err)
		}
		logrus.Debugf("created cos instance: %s", cosInstanceName)
	}
	bucketName := fmt.Sprintf("%s-vsi-image", in.InfraID)
	logrus.Debugf("checking for existing cos bucket: %s", bucketName)
	_, err = client.GetCOSBucketByName(ctx, *cosInstance.ID, bucketName, region)
	if err != nil {
		logrus.Debugf("creating cos bucket: %s", bucketName)
		err = client.CreateCOSBucket(ctx, *cosInstance.ID, bucketName, region)
		if err != nil {
			return fmt.Errorf("failed creating RHCOS image COS bucket: %w", err)
		}
		logrus.Debugf("created cos bucket: %s", bucketName)
	}

	// Upload the RHCOS image to the COS Bucket.
	logrus.Debugf("retreiving rhcos image for upload to cos")
	cachedImage, err := cache.DownloadImageFile(in.RhcosImage.ControlPlane, cache.InstallerApplicationName)
	if err != nil {
		return fmt.Errorf("failed to use cached ibmcloud image: %w", err)
	}
	imageData, err := os.ReadFile(cachedImage)
	if err != nil {
		return fmt.Errorf("failed reading RHCOS image data: %w", err)
	}
	objectName := filepath.Base(cachedImage)
	logrus.Debugf("uploading rhcos image to cos: %s", objectName)
	err = client.CreateCOSObject(ctx, imageData, objectName, *cosInstance.ID, bucketName, region)
	if err != nil {
		return fmt.Errorf("failed uploading RHCOS image: %w", err)
	}
	logrus.Debugf("rhcos image uploaded to cos: %s", objectName)

	// Create IAM authorization for VPC to COS access for Custom Image Creation
	logrus.Debugf("creating iam authorization for vpc to cos access")
	err = client.CreateIAMAuthorizationPolicy(ctx, "is", "image", "cloud-object-storage", *cosInstance.ID, []string{"crn:v1:bluemix:public:iam::::serviceRole:Reader"})
	if err != nil {
		return fmt.Errorf("failed creating vpc-cos IAM authorization policy: %w", err)
	}
	logrus.Debugf("created iam authorization for vpc to cos access")

	return nil
}
