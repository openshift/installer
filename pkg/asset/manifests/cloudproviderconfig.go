package manifests

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	ibmcloudmachines "github.com/openshift/installer/pkg/asset/machines/ibmcloud"
	alibabacloudmanifests "github.com/openshift/installer/pkg/asset/manifests/alibabacloud"
	"github.com/openshift/installer/pkg/asset/manifests/azure"
	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	ibmcloudmanifests "github.com/openshift/installer/pkg/asset/manifests/ibmcloud"
	openstackmanifests "github.com/openshift/installer/pkg/asset/manifests/openstack"
	powervsmanifests "github.com/openshift/installer/pkg/asset/manifests/powervs"
	vspheremanifests "github.com/openshift/installer/pkg/asset/manifests/vsphere"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

var (
	cloudProviderConfigFileName = filepath.Join(manifestDir, "cloud-provider-config.yaml")
)

const (
	cloudProviderConfigDataKey         = "config"
	cloudProviderConfigCABundleDataKey = "ca-bundle.pem"
	cloudProviderEndpointsKey          = "endpoints"
)

// CloudProviderConfig generates the cloud-provider-config.yaml files.
type CloudProviderConfig struct {
	ConfigMap *corev1.ConfigMap
	File      *asset.File
}

var _ asset.WritableAsset = (*CloudProviderConfig)(nil)

// Name returns a human friendly name for the asset.
func (*CloudProviderConfig) Name() string {
	return "Cloud Provider Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*CloudProviderConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},

		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
	}
}

// Generate generates the CloudProviderConfig.
func (cpc *CloudProviderConfig) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	dependencies.Get(installConfig, clusterID)

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-config",
			Name:      "cloud-provider-config",
		},
		Data: map[string]string{},
	}

	switch installConfig.Config.Platform.Name() {
	case libvirttypes.Name, nonetypes.Name, baremetaltypes.Name, ovirttypes.Name:
		return nil
	case awstypes.Name:
		// Store the additional trust bundle in the ca-bundle.pem key if the cluster is being installed on a C2S region.
		trustBundle := installConfig.Config.AdditionalTrustBundle
		if trustBundle == "" || !awstypes.IsSecretRegion(installConfig.Config.AWS.Region) {
			return nil
		}
		cm.Data[cloudProviderConfigCABundleDataKey] = trustBundle
		// Include a non-empty kube config to appease components--such as the kube-apiserver--that
		// expect there to be a kube config if the cloud-provider-config ConfigMap exists. See
		// https://bugzilla.redhat.com/show_bug.cgi?id=1926975.
		// Note that the newline is required in order to be valid yaml.
		cm.Data[cloudProviderConfigDataKey] = `[Global]
`
	case alibabacloudtypes.Name:
		alibabacloudConfig, err := alibabacloudmanifests.CloudConfig{
			Global: alibabacloudmanifests.GlobalConfig{
				ClusterID: clusterID.InfraID,
				Region:    installConfig.Config.AlibabaCloud.Region,
			},
		}.JSON()
		if err != nil {
			return errors.Wrap(err, "could not create Alibaba Cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = alibabacloudConfig
	case openstacktypes.Name:
		cloudProviderConfigData, cloudProviderConfigCABundleData, err := openstackmanifests.GenerateCloudProviderConfig(*installConfig.Config)
		if err != nil {
			return errors.Wrap(err, "failed to generate OpenStack provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = cloudProviderConfigData
		if cloudProviderConfigCABundleData != "" {
			cm.Data[cloudProviderConfigCABundleDataKey] = cloudProviderConfigCABundleData
		}

	case azuretypes.Name:
		session, err := installConfig.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "could not get azure session")
		}

		nsg := fmt.Sprintf("%s-nsg", clusterID.InfraID)
		nrg := installConfig.Config.Azure.ClusterResourceGroupName(clusterID.InfraID)
		if installConfig.Config.Azure.NetworkResourceGroupName != "" {
			nrg = installConfig.Config.Azure.NetworkResourceGroupName
		}
		vnet := fmt.Sprintf("%s-vnet", clusterID.InfraID)
		if installConfig.Config.Azure.VirtualNetwork != "" {
			vnet = installConfig.Config.Azure.VirtualNetwork
		}
		subnet := fmt.Sprintf("%s-worker-subnet", clusterID.InfraID)
		if installConfig.Config.Azure.ComputeSubnet != "" {
			subnet = installConfig.Config.Azure.ComputeSubnet
		}
		azureConfig, err := azure.CloudProviderConfig{
			CloudName:                installConfig.Config.Azure.CloudName,
			ResourceGroupName:        installConfig.Config.Azure.ClusterResourceGroupName(clusterID.InfraID),
			GroupLocation:            installConfig.Config.Azure.Region,
			ResourcePrefix:           clusterID.InfraID,
			SubscriptionID:           session.Credentials.SubscriptionID,
			TenantID:                 session.Credentials.TenantID,
			NetworkResourceGroupName: nrg,
			NetworkSecurityGroupName: nsg,
			VirtualNetworkName:       vnet,
			SubnetName:               subnet,
			ResourceManagerEndpoint:  installConfig.Config.Azure.ARMEndpoint,
			ARO:                      installConfig.Config.Azure.IsARO(),
		}.JSON()
		if err != nil {
			return errors.Wrap(err, "could not create cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = azureConfig

		if installConfig.Azure.CloudName == azuretypes.StackCloud {
			b, err := json.Marshal(session.Environment)
			if err != nil {
				return errors.Wrap(err, "could not serialize Azure Stack endpoints")
			}
			cm.Data[cloudProviderEndpointsKey] = string(b)
		}
	case gcptypes.Name:
		subnet := fmt.Sprintf("%s-worker-subnet", clusterID.InfraID)
		if installConfig.Config.GCP.ComputeSubnet != "" {
			subnet = installConfig.Config.GCP.ComputeSubnet
		}
		gcpConfig, err := gcpmanifests.CloudProviderConfig(clusterID.InfraID, installConfig.Config.GCP.ProjectID, subnet)
		if err != nil {
			return errors.Wrap(err, "could not create cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = gcpConfig
	case ibmcloudtypes.Name:
		accountID, err := installConfig.IBMCloud.AccountID(context.TODO())
		if err != nil {
			return err
		}

		controlPlane := &ibmcloudtypes.MachinePool{}
		controlPlane.Set(installConfig.Config.Platform.IBMCloud.DefaultMachinePlatform)
		controlPlane.Set(installConfig.Config.ControlPlane.Platform.IBMCloud)
		compute := &ibmcloudtypes.MachinePool{}
		compute.Set(installConfig.Config.Platform.IBMCloud.DefaultMachinePlatform)
		compute.Set(installConfig.Config.WorkerMachinePool().Platform.IBMCloud)

		if len(controlPlane.Zones) == 0 || len(compute.Zones) == 0 {
			zones, err := ibmcloudmachines.AvailabilityZones(installConfig.Config.IBMCloud.Region)
			if err != nil {
				return errors.Wrapf(err, "could not get availability zones for %s", installConfig.Config.IBMCloud.Region)
			}
			if len(controlPlane.Zones) == 0 {
				controlPlane.Zones = zones
			}
			if len(compute.Zones) == 0 {
				compute.Zones = zones
			}
		}

		resourceGroupName := installConfig.Config.Platform.IBMCloud.ClusterResourceGroupName(clusterID.InfraID)
		ibmcloudConfig, err := ibmcloudmanifests.CloudProviderConfig(clusterID.InfraID, accountID, installConfig.Config.IBMCloud.Region, resourceGroupName, controlPlane.Zones, compute.Zones)
		if err != nil {
			return errors.Wrap(err, "could not create cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = ibmcloudConfig
	case powervstypes.Name:
		var vpcRegion string
		var err error

		accountID, err := installConfig.PowerVS.AccountID(context.TODO())
		if err != nil {
			return err
		}

		vpcRegion, err = powervstypes.VPCRegionForPowerVSRegion(installConfig.Config.PowerVS.Region)
		if err != nil {
			return err
		}

		powervsConfig, err := powervsmanifests.CloudProviderConfig(
			clusterID.InfraID,
			accountID,
			installConfig.Config.PowerVS.VPC,
			vpcRegion,
			installConfig.Config.Platform.PowerVS.PowerVSResourceGroup,
			installConfig.Config.PowerVS.Subnets)
		if err != nil {
			return errors.Wrap(err, "could not create cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = powervsConfig
	case vspheretypes.Name:
		folderPath := installConfig.Config.Platform.VSphere.Folder
		if len(folderPath) == 0 {
			dataCenter := installConfig.Config.Platform.VSphere.Datacenter
			folderPath = fmt.Sprintf("/%s/vm/%s", dataCenter, clusterID.InfraID)
		}
		vsphereConfig, err := vspheremanifests.CloudProviderConfig(
			folderPath,
			installConfig.Config.Platform.VSphere,
		)
		if err != nil {
			return errors.Wrap(err, "could not create cloud provider config")
		}
		cm.Data[cloudProviderConfigDataKey] = vsphereConfig
	default:
		return errors.New("invalid Platform")
	}

	cmData, err := yaml.Marshal(cm)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifest", cpc.Name())
	}
	cpc.ConfigMap = cm
	cpc.File = &asset.File{
		Filename: cloudProviderConfigFileName,
		Data:     cmData,
	}
	return nil
}

// Files returns the files generated by the asset.
func (cpc *CloudProviderConfig) Files() []*asset.File {
	if cpc.File != nil {
		return []*asset.File{cpc.File}
	}
	return []*asset.File{}
}

// Load loads the already-rendered files back from disk.
func (cpc *CloudProviderConfig) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
