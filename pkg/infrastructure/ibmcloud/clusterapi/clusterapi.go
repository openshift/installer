package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api/util/patch"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/openshift/api/config/v1"
	ibmcloudbootstrap "github.com/openshift/installer/pkg/asset/ignition/bootstrap/ibmcloud"
	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

var _ clusterapi.BootstrapDestroyer = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)
var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)
var _ clusterapi.Provider = (*Provider)(nil)
var _ clusterapi.Timeouts = (*Provider)(nil)

// Provider implements IBM Cloud CAPI installation.
type Provider struct{}

// Name returns the IBM Cloud provider name.
func (p Provider) Name() string {
	return ibmcloudtypes.Name
}

// NetworkTimeout allows platform provider to override the timeout
// when waiting for the network infrastructure to become ready.
func (p Provider) NetworkTimeout() time.Duration {
	// IBM Cloud requires additional time for VPC Custom Image creation and Load Balancer reconciliation.
	return 30 * time.Minute
}

// ProvisionTimeout allows platform provider to override the timeout
// when waiting for the machines to provision.
func (p Provider) ProvisionTimeout() time.Duration {
	return 25 * time.Minute
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
		// Retrieve the newly created resource group.
		// Use retry logic to wait for the new resource group if necessary.
		backoff := wait.Backoff{
			Duration: 10 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32,
		}

		var lastErr error
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			resourceGroup, lastErr = client.GetResourceGroup(ctx, resourceGroupName)
			if lastErr == nil {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			if lastErr != nil {
				err = lastErr
			}
			return fmt.Errorf("failed retrieving new resource group: %w", err)
		}
		logrus.Debugf("created resource group: %s", resourceGroupName)
	}

	// Create a COS Instance and Bucket to host the RHCOS image file.
	// NOTE(cjschaef): Support to use an existing COS Object (RHCO image file) or VPC Custom Image could be added to skip this step.
	cosInstanceName := ibmcloudic.COSInstanceName(in.InfraID)
	logrus.Debugf("checking for existing cos instance: %s", cosInstanceName)
	var cosInstanceNotFoundError *ibmcloudic.COSResourceNotFoundError
	cosInstance, err := client.GetCOSInstanceByName(ctx, cosInstanceName)
	if err != nil {
		if errors.As(err, &cosInstanceNotFoundError) {
			// Attempt to create the COS Instance, since it was not found.
			logrus.Debugf("creating cos instance: %s", cosInstanceName)
			cosInstance, err = client.CreateCOSInstance(ctx, cosInstanceName, *resourceGroup.ID)
			if err != nil {
				return fmt.Errorf("failed creating RHCOS image COS instance: %w", err)
			}
			logrus.Debugf("created cos instance: %s", cosInstanceName)
		} else {
			return fmt.Errorf("failed checking for cos instance %s: %w", cosInstanceName, err)
		}
	}
	bucketName := ibmcloudic.VSIImageCOSBucketName(in.InfraID)
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
	err = client.CreateIAMAuthorizationPolicy(ctx, "is", "image", "cloud-object-storage", *cosInstance.GUID, []string{"crn:v1:bluemix:public:iam::::serviceRole:Reader"})
	if err != nil {
		return fmt.Errorf("failed creating vpc-cos IAM authorization policy: %w", err)
	}
	logrus.Debugf("created iam authorization for vpc to cos access")

	return nil
}

// InfraReady is called once cluster.Status.InfrastructureReady is true.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	// 1. Collect necessary details from Cluster resource.
	// 2. Create DNS Records for the Control Plane's Load Balancers (one for public LB - 'api' and one for private LB - 'api-int'). For Public/External clusters, records in the CIS instance is created. For Private/Internal cluster, the records are created in the DNS Services instance.
	// 3. For Private/Internal cluster, add the VPC to the DNS Services Zone's Permitted Networks, if not already there.

	// Setup IBM Cloud Client.
	metadata := ibmcloudic.NewMetadata(in.InstallConfig.Config)
	client, err := metadata.Client()
	if err != nil {
		return fmt.Errorf("failed creating IBM Cloud client in InfraReady: %w", err)
	}

	ibmcloudCluster := &capibmcloud.IBMVPCCluster{}

	// Get the cluster from the provider.
	key := crclient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	logrus.Debugf("InfraReady: clusterKey = %+v", key)
	if err = in.Client.Get(ctx, key, ibmcloudCluster); err != nil {
		return fmt.Errorf("failed to get ibmcloud cluster in InfraReady: %w", err)
	}
	logrus.Debugf("InfraReady: ibmcloudCluster = %+v", ibmcloudCluster)
	logrus.Debugf("InfraReady: ibmcloudCluster.Status = %+v", ibmcloudCluster.Status)
	if ibmcloudCluster.Status.Network == nil || ibmcloudCluster.Status.Network.VPC == nil || ibmcloudCluster.Status.Network.VPC.ID == "" {
		return fmt.Errorf("vpc missing from ibmcloudCluster.Status in InfraReady")
	}
	logrus.Debugf("InfraReady: ibmcloudCluster.Status.Network.VPC.ID = %s", ibmcloudCluster.Status.Network.VPC.ID)

	// Collect the Load Balancer's and their hostnames to use for DNS Record creation.
	if len(ibmcloudCluster.Status.Network.LoadBalancers) == 0 {
		return fmt.Errorf("load balancers missing from ibmcloudCluster.Status in InfraReady")
	}

	domain := in.InstallConfig.Config.ClusterDomain()

	// For now, we expect one of two LB configurations,
	// 1. One Public LB for Public APIServer traffic and one Private LB for Private APIServer traffic (Public/External clusters).
	// 2. One Private LB for Private APIServer traffic and "Public" traffic only accessible within VPC (Private/Internal clusters).
	for lbID, lb := range ibmcloudCluster.Status.Network.LoadBalancers {
		// Verify that the Load Balancer is ready (active).
		if lb.State != capibmcloud.VPCLoadBalancerStateActive {
			return fmt.Errorf("load balancer %s is not ready, infrastructure not ready: %s", lbID, lb.State)
		}
		// Lookup Load Balancer details to use during DNS Record creation.
		lbDetails, err := client.GetLoadBalancer(ctx, lbID)
		if err != nil {
			return fmt.Errorf("failed retrieving load balancer for dns record creation: %w", err)
		} else if lbDetails == nil {
			return fmt.Errorf("failed to find load balancer for dns record creation by id: %s", lbID)
		}
		switch in.InstallConfig.Config.Publish {
		case types.ExternalPublishingStrategy:
			var recordName string
			// Build the record name based on the LB name/type, ignore LB's not named and configured for Kube API Server traffic.
			switch {
			case strings.HasSuffix(*lbDetails.Name, ibmcloudic.KubernetesAPIPublicSuffix):
				recordName = fmt.Sprintf("%s%s", ibmcloudic.PublicHostPrefix, domain)
			case strings.HasSuffix(*lbDetails.Name, ibmcloudic.KubernetesAPIPrivateSuffix):
				recordName = fmt.Sprintf("%s%s", ibmcloudic.PrivateHostPrefix, domain)
			default:
				logrus.Debug("ignoring unexpected load balancer for external cluster", "lbName", *lbDetails.Name)
				continue
			}
			err = metadata.CreateDNSRecord(ctx, recordName, lbDetails)
		case types.InternalPublishingStrategy:
			// Create both DNS Records for the expected single Private LB, ignore the LB if not named and configured for Kube API Server traffic.
			if !strings.HasSuffix(*lbDetails.Name, ibmcloudic.KubernetesAPIPrivateSuffix) {
				logrus.Debug("ignoring unexpected load balancer for internal cluster", "lbName", *lbDetails.Name)
				continue
			}
			err = metadata.CreateDNSRecord(ctx, fmt.Sprintf("%s%s", ibmcloudic.PublicHostPrefix, domain), lbDetails)
			if err != nil {
				return fmt.Errorf("failed to create public dns record for private load balancer: %w", err)
			}
			logrus.Debug("public dns record created for private load balancer", "hostName", *lbDetails.Hostname)
			err = metadata.CreateDNSRecord(ctx, fmt.Sprintf("%s%s", ibmcloudic.PrivateHostPrefix, domain), lbDetails)
		default:
			return fmt.Errorf("failed to create dns record, invalid publish strategy: %s", in.InstallConfig.Config.Publish)
		}
		if err != nil {
			return fmt.Errorf("failed to create dns record for load balancer: %w", err)
		}
		logrus.Debug("dns record created for load balancer", "hostName", *lbDetails.Hostname)
	}

	logrus.Debug("checking cluster publishing strategy", "publish", in.InstallConfig.Config.Publish)
	// For Private/Internal cluster, check DNS Services Zone's Permitted Network.
	if in.InstallConfig.Config.Publish == types.InternalPublishingStrategy {
		logrus.Debug("checking dns services permitted network for vpc", "vpcName", in.InstallConfig.Config.IBMCloud.VPCName)
		// Determine whether the VPC is already a Permitted Network, if not, add the VPC to the DNS Services Zone's Permitted Networks.
		// Since this check is based on the value provided for vpcName in the InstallConfig, pass that value to help shortcut the check (when vpcName is empty or not provided, assume a new VPC was created and thus needs to be added to Permitted Networks).
		permitted, err := metadata.IsVPCPermittedNetwork(ctx, in.InstallConfig.Config.IBMCloud.VPCName)
		if err != nil {
			return fmt.Errorf("failed to check whether vpc is a permitted network: %w", err)
		}
		// If VPC is already a PermittedNetwork, no further action necessary.
		if permitted {
			logrus.Debug("vpc already a permitted network", "vpcName", in.InstallConfig.Config.IBMCloud.VPCName)
			return nil
		}
		// If not, attempt to add the VPC to PermittedNetworks.
		if err = metadata.AddVPCToPermittedNetworks(ctx, ibmcloudCluster.Status.Network.VPC.ID); err != nil {
			return fmt.Errorf("failed to add vpc %s to dns services permitted networks: %w", ibmcloudCluster.Status.Network.VPC.ID, err)
		}
	}

	return nil
}

// PostProvision is called once Infrastructure provisioning is completed, including machines.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	// If the cluster is Public/External, attach a floating IP to the bootstrap node for debugging purposes. It will need to be cleaned up during bootstrap destroy however.

	if in.InstallConfig.Config.Publish == types.InternalPublishingStrategy {
		logrus.Debugf("no post provisioning required for internal clusters")
		return nil
	}

	metadata := ibmcloudic.NewMetadata(in.InstallConfig.Config)
	client, err := metadata.Client()
	if err != nil {
		return fmt.Errorf("failed creating IBM Cloud client in PostProvision: %w", err)
	}

	// Collect the bootstrap machine from the provider.
	ibmcloudMachine := &capibmcloud.IBMVPCMachine{}
	key := crclient.ObjectKey{
		Name:      capiutils.GenerateBoostrapMachineName(in.InfraID),
		Namespace: capiutils.Namespace,
	}
	logrus.Debugf("PostProvision: machineKey = %+v", key)
	if err = in.Client.Get(ctx, key, ibmcloudMachine); err != nil {
		return fmt.Errorf("failed to get ibmcloud bootstrap machine in PostProvision: %w", err)
	}

	resourceGroupName := in.InfraID
	if in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName != "" {
		resourceGroupName = in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName
	}
	logrus.Debugf("collected resource group name %s for floating ip", resourceGroupName)

	// Attach a new VPC Floating IP resource to the bootstrap machine.
	logrus.Debugf("creating floating ip for bootstrap machine")
	floatingIPDetails, err := client.AttachFloatingIP(ctx, ibmcloudMachine.Name, ibmcloudMachine.Status.InstanceID, in.InstallConfig.Config.Platform.IBMCloud.Region, resourceGroupName)
	if err != nil {
		return fmt.Errorf("failed to attaching floating ip to bootstrap machine: %w", err)
	}

	// Create patch helper to update the bootstrap machine addresses.
	patchHelper, err := patch.NewHelper(ibmcloudMachine, in.Client)
	if err != nil {
		return fmt.Errorf("failed creating patch helper to update machine addresses")
	}

	// Update the bootstrap machine's external address with the new floating IP address.
	ibmcloudMachine.Status.Addresses = append(ibmcloudMachine.Status.Addresses, corev1.NodeAddress{
		Address: *floatingIPDetails.Address,
		Type:    corev1.NodeExternalIP,
	})

	// Patch the bootstrap machine.
	if err = patchHelper.Patch(ctx, ibmcloudMachine); err != nil {
		return fmt.Errorf("failed to patch bootstrap machine with floating ip address: %w", err)
	}

	logrus.Debugf("floating ip attached to bootstrap machine at %s", *floatingIPDetails.Address)

	return nil
}

func leftInContext(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return math.MaxInt64
	}
	return time.Until(deadline)
}

// Ignition provisions the IBM Cloud COS Bucket and Object containing the Ignition based configuration.
// The Bootstrap ignition data is too large to be passed as userdata to the IBM Cloud VPC VSI, so instead it is pulled from COS.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]*corev1.Secret, error) {
	// Setup IBM Cloud Client.
	metadata := ibmcloudic.NewMetadata(in.InstallConfig.Config)
	client, err := metadata.Client()
	if err != nil {
		return nil, fmt.Errorf("failed creating IBM Cloud client in Ignition: %w", err)
	}
	region := in.InstallConfig.Config.Platform.IBMCloud.Region

	// Get the Resource Group name, which should already exist, and lookup the Resource Group ID.
	resourceGroupName := in.InfraID
	if in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName != "" {
		resourceGroupName = in.InstallConfig.Config.Platform.IBMCloud.ResourceGroupName
	}
	logrus.Debugf("retrieving resource group id for: %s", resourceGroupName)
	resourceGroup, err := client.GetResourceGroup(ctx, resourceGroupName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource group %s: %w", resourceGroupName, err)
	}
	logrus.Debugf("retrieved resource group id: %s", *resourceGroup.ID)

	// Get the COS Instance, possibly created for RHCOS image, or create the COS Instance.
	cosInstanceName := ibmcloudic.COSInstanceName(in.InfraID)
	var cosInstanceNotFoundError *ibmcloudic.COSResourceNotFoundError
	cosInstance, err := client.GetCOSInstanceByName(ctx, cosInstanceName)
	if err != nil {
		if errors.As(err, &cosInstanceNotFoundError) {
			// Attempt to create the COS Instance, since it was not found.
			logrus.Debugf("creating cos instance: %s", cosInstanceName)
			cosInstance, err = client.CreateCOSInstance(ctx, cosInstanceName, *resourceGroup.ID)
			if err != nil {
				return nil, fmt.Errorf("failed creating ignition COS instance: %w", err)
			}
			logrus.Debugf("created cos instance: %s", cosInstanceName)
		} else {
			return nil, fmt.Errorf("failed retrieving cos instance %s for ignition: %w", cosInstanceName, err)
		}
	}

	// Create new bucket for bootstrap's temporary Ignition Config.
	logrus.Debugf("fetching cos instance for cluster: %s", cosInstanceName)
	bucketName := ibmcloudbootstrap.GetIgnitionBucketName((in.InfraID))
	logrus.Debugf("creating cos bucket for bootstrap ignition config: %s", bucketName)
	err = client.CreateCOSBucket(ctx, *cosInstance.ID, bucketName, region)
	if err != nil {
		return nil, fmt.Errorf("failed creating ignition COS bucket: %w", err)
	}
	logrus.Infof("created cos bucket for bootstrap ignition config: %s/%s", cosInstanceName, bucketName)

	// Default to using the direct regional COS endpoint.
	cosEndpoint := fmt.Sprintf("s3.direct.%s.cloud-object-storage.appdomain.cloud", region)
	// Check whether an endpoint override was provided for COS.
	if endpointURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCOS, in.InstallConfig.Config.IBMCloud.ServiceEndpoints); endpointURL != "" {
		cosEndpoint = endpointURL
	}

	// Upload Ignition Config to COS bucket.
	logrus.Debugf("uploading bootstrap ignition config to bucket: %s", bucketName)
	ignitionFile := ibmcloudbootstrap.GetIgnitionFileName()
	err = client.CreateCOSObject(ctx, in.BootstrapIgnData, ignitionFile, *cosInstance.ID, bucketName, region)
	if err != nil {
		return nil, fmt.Errorf("failed uploading ignition data: %w", err)
	}
	logrus.Debugf("bootstrap ignition config upload complete to %s/%s/%s", cosInstanceName, bucketName, ignitionFile)

	// Build the URL for the ignition config.
	cosURL, err := url.Parse(cosEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ibmcloud cos url from %s: %w", cosEndpoint, err)
	}
	// Make sure the COS URL has an https scheme, if one isn't already set.
	if cosURL.Scheme == "" {
		cosURL.Scheme = "https"
	}
	ignitionURL := cosURL.JoinPath(bucketName, ignitionFile)

	// Build Ignition Config for Secret to direct bootstrap to consume COS Ignition Config.
	logrus.Debugf("building ignition config data for bootstrap secret")
	// Get IAM token for bootstrap node to access the Ignition config in COS.
	iamToken, err := metadata.GetIAMToken(client.GetAPIKey())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve iam token for ignition: %w", err)
	}

	logrus.Debugf("building bootstrap ignition config for bootstrap secret")
	// NOTE(cjschaef): Replace the reliance on using the IAM token with a Service ID credential, when working with the COS Instance.
	ignShim, err := ibmcloudbootstrap.GenerateIgnitionShimWithCredentials(ignitionURL.String(), *iamToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition shim: %w", err)
	}
	logrus.Debugf("bootstrap ignition config built for bootstrap secret")

	ignSecrets := []*corev1.Secret{
		clusterapi.IgnitionSecret(ignShim, in.InfraID, "bootstrap"),
		clusterapi.IgnitionSecret(in.MasterIgnData, in.InfraID, "master"),
	}

	return ignSecrets, nil
}

// DestroyBootstrap destroys the remaining temporary bootstrap resources.
func (p Provider) DestroyBootstrap(ctx context.Context, in clusterapi.BootstrapDestroyInput) error {
	// Currently expect the Bootstrap Node will be deleted via CAPI (along with its LB Pool Members).
	// Remaining IBM Cloud Resources that need to be cleaned up:
	// 1. COS Bucket for temporary bootstrap Ignition Config.
	// 2. Potentially the COS Instance, if no other Buckets exist (specifically the Bucket for VSI Image).
	// 3. Bootstrap Security Group.
	// 4. Bootstrap Floating IP (if one was created, for External/Public clusters).
	// 5. TBD - Bootstrap Machine's LoadBalancerPoolMembers (check whether this cleanup is required here or elsewhere).

	logrus.Infof("bootstrap destroy ibmcloud resource cleanup")
	// Setup IBM Cloud Client.
	client, err := ibmcloudic.NewClient(in.Metadata.IBMCloud.ServiceEndpoints)
	if err != nil {
		return fmt.Errorf("failed creating IBM Cloud client during destroy bootstrap: %w", err)
	}
	infraID := in.Metadata.InfraID
	region := in.Metadata.IBMCloud.Region
	bootstrapMachineName := capiutils.GenerateBoostrapMachineName(infraID)
	logrus.Debugf("bootstrap destroy cleanup resources for %s", infraID)

	// Cleanup COS Bucket for bootstrap Ignition Config
	cosInstanceName := ibmcloudic.COSInstanceName(infraID)
	bucketName := ibmcloudbootstrap.GetIgnitionBucketName(infraID)
	// Get COS Instance ID
	logrus.Debugf("bootstrap destroy check cos resources for %s/%s", cosInstanceName, bucketName)
	cosInstanceDetails, err := client.GetCOSInstanceByName(ctx, cosInstanceName)
	switch {
	case err != nil:
		return fmt.Errorf("failed retrieving cos instance for destroy bootstrap: %w", err)
	case cosInstanceDetails != nil:
		err = cleanupIgnitionCOSBucket(ctx, client, *cosInstanceDetails.ID, bucketName, region)
		if err != nil {
			return fmt.Errorf("failed to cleanup bootstrap ignition cos bucket %s: %w", bucketName, err)
		}
		logrus.Debugf("cleaned up bootstrap ignition cos bucket: %s/%s", *cosInstanceDetails.ID, bucketName)
	default:
		logrus.Debugf("no cos instance found for the cluster as %s, skipping ignition bucket cleanup", cosInstanceName)
	}

	// If there are no additional Buckets within the COS Instance (specifically the Bucket used for the VSI Image), cleanup the COS Instance as well.
	logrus.Debugf("checking whether cos instance should be cleaned up as well: %s", *cosInstanceDetails.ID)
	cosBucketsOutput, err := client.ListCOSBuckets(ctx, *cosInstanceDetails.ID, region)
	switch {
	case err != nil:
		return fmt.Errorf("failed listing cos buckets in cos instance %s: %w", *cosInstanceDetails.ID, err)
	case cosBucketsOutput == nil || len(cosBucketsOutput.Buckets) == 0:
		logrus.Debugf("no remaining buckets in cos instance %s, attempting to cleanup instance", *cosInstanceDetails.ID)
		err := client.DeleteCOSInstance(ctx, *cosInstanceDetails.ID)
		if err != nil {
			return fmt.Errorf("failed to delete empty cos instance %s: %w", *cosInstanceDetails.ID, err)
		}
	case len(cosBucketsOutput.Buckets) == 1 && cosBucketsOutput.Buckets[0].Name != nil && *cosBucketsOutput.Buckets[0].Name == bucketName:
		logrus.Debugf("bootstrap ignition cos bucket %s still listed in cos instance, proceeding to cleanup cos instance: %s", bucketName, *cosInstanceDetails.ID)
		err := client.DeleteCOSInstance(ctx, *cosInstanceDetails.ID)
		if err != nil {
			return fmt.Errorf("failed to delete cos instance %s, with single bucket %s: %w", *cosInstanceDetails.ID, bucketName, err)
		}
	default:
		logrus.Debugf("cos instance contains additional buckets, skipping cos instance %s cleanup", *cosInstanceDetails.ID)
	}

	// Cleanup the Security Group for the Bootstrap node.
	vpcName := in.Metadata.IBMCloud.VPC
	// If VPC name is empty, it was auto-generated from the InfraID (not a pre-existing VPC for BYON support).
	if vpcName == "" {
		vpcName = ibmcloudic.VPCName(infraID)
	}

	logrus.Debugf("bootstrap destroy checking for vpc %s", vpcName)
	var vpcResourceNotFoundError *ibmcloudic.VPCResourceNotFoundError
	vpcDetails, err := client.GetVPCByName(ctx, vpcName)
	if err != nil {
		if errors.As(err, &vpcResourceNotFoundError) {
			return fmt.Errorf("no vpc found for the cluster as %s, unable to complete cleanup", vpcName)
		}
		return fmt.Errorf("failed retrieving vpc for destroy bootstrap: %w", err)
	}

	bootstrapSGName := ibmcloudic.BootstrapSecurityGroupName(infraID)
	logrus.Debugf("cleaning up bootstrap security group %s", bootstrapSGName)
	if err = cleanupBootstrapSecurityGroup(ctx, client, *vpcDetails.ID, bootstrapSGName, region); err != nil {
		return fmt.Errorf("failed to cleanup boostrap security group %s: %w", bootstrapSGName, err)
	}

	// Cleanup the Floating IP for the Bootstrap node, if this is an External/Public cluster.
	logrus.Debugf("checking for cleanup of floating ip %s", bootstrapMachineName)
	floatingIPDetails, err := client.GetFloatingIPByName(ctx, bootstrapMachineName, region)
	switch {
	case err != nil:
		// If this isn't a Not Found error (it is okay if it doesn't exist), return the unexpected error.
		if !errors.As(err, &vpcResourceNotFoundError) {
			return fmt.Errorf("failed checking for existing bootstrap floating ip: %w", err)
		}
		logrus.Debugf("no bootstrap floating ip found for cleanup, skipping")
	case floatingIPDetails != nil:
		logrus.Debugf("found bootstrap floating ip for cleanup: %s=%s", *floatingIPDetails.Name, *floatingIPDetails.ID)
		err = client.DeleteFloatingIP(ctx, *floatingIPDetails.ID, region)
		if err != nil {
			return fmt.Errorf("failed deleting bootstrap floating ip: %w", err)
		}
	default:
		logrus.Debugf("no bootstrap floating ip found for cleanup, skipping")
	}

	logrus.Infof("destroy bootstrap cleanup completed")
	return nil
}
