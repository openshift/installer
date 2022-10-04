package manifests

import (
	"context"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	vsphereinfra "github.com/openshift/installer/pkg/asset/manifests/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	infraCfgFilename           = filepath.Join(manifestDir, "cluster-infrastructure-02-config.yml")
	cloudControllerUIDFilename = filepath.Join(manifestDir, "cloud-controller-uid-config.yml")
)

// Infrastructure generates the cluster-infrastructure-*.yml files.
type Infrastructure struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Infrastructure)(nil)

// Name returns a human friendly name for the asset.
func (*Infrastructure) Name() string {
	return "Infrastructure Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Infrastructure) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		&CloudProviderConfig{},
		&AdditionalTrustBundleConfig{},
	}
}

// Generate generates the Infrastructure config and its CRD.
func (i *Infrastructure) Generate(dependencies asset.Parents) error {
	cloudProviderConfigMapKey := cloudProviderConfigDataKey
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	cloudproviderconfig := &CloudProviderConfig{}
	trustbundleconfig := &AdditionalTrustBundleConfig{}
	dependencies.Get(clusterID, installConfig, cloudproviderconfig, trustbundleconfig)

	config := &configv1.Infrastructure{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Infrastructure",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.InfrastructureSpec{
			PlatformSpec: configv1.PlatformSpec{},
		},
		Status: configv1.InfrastructureStatus{
			InfrastructureName:   clusterID.InfraID,
			APIServerURL:         getAPIServerURL(installConfig.Config),
			APIServerInternalURL: getInternalAPIServerURL(installConfig.Config),
			PlatformStatus:       &configv1.PlatformStatus{},
		},
	}

	controlPlaneTopology, infrastructureTopology := determineTopologies(installConfig.Config)

	config.Status.InfrastructureTopology = infrastructureTopology
	config.Status.ControlPlaneTopology = controlPlaneTopology

	switch installConfig.Config.Platform.Name() {
	case aws.Name:
		config.Spec.PlatformSpec.Type = configv1.AWSPlatformType
		config.Spec.PlatformSpec.AWS = &configv1.AWSPlatformSpec{}

		var resourceTags []configv1.AWSResourceTag
		if installConfig.Config.AWS.PropagateUserTag {
			resourceTags = make([]configv1.AWSResourceTag, 0, len(installConfig.Config.AWS.UserTags))
			for k, v := range installConfig.Config.AWS.UserTags {
				resourceTags = append(resourceTags, configv1.AWSResourceTag{Key: k, Value: v})
			}
		}
		config.Status.PlatformStatus.AWS = &configv1.AWSPlatformStatus{
			Region:       installConfig.Config.Platform.AWS.Region,
			ResourceTags: resourceTags,
		}

		for _, service := range installConfig.Config.Platform.AWS.ServiceEndpoints {
			config.Spec.PlatformSpec.AWS.ServiceEndpoints = append(config.Spec.PlatformSpec.AWS.ServiceEndpoints, configv1.AWSServiceEndpoint{
				Name: service.Name,
				URL:  service.URL,
			})
			config.Status.PlatformStatus.AWS.ServiceEndpoints = append(config.Status.PlatformStatus.AWS.ServiceEndpoints, configv1.AWSServiceEndpoint{
				Name: service.Name,
				URL:  service.URL,
			})
			sort.Slice(config.Status.PlatformStatus.AWS.ServiceEndpoints, func(i, j int) bool {
				return config.Status.PlatformStatus.AWS.ServiceEndpoints[i].Name <
					config.Status.PlatformStatus.AWS.ServiceEndpoints[j].Name
			})
		}
	case azure.Name:
		config.Spec.PlatformSpec.Type = configv1.AzurePlatformType

		rg := installConfig.Config.Azure.ClusterResourceGroupName(clusterID.InfraID)
		config.Status.PlatformStatus.Azure = &configv1.AzurePlatformStatus{
			ResourceGroupName:        rg,
			NetworkResourceGroupName: rg,
			CloudName:                configv1.AzureCloudEnvironment(installConfig.Config.Platform.Azure.CloudName),
		}
		if nrg := installConfig.Config.Platform.Azure.NetworkResourceGroupName; nrg != "" {
			config.Status.PlatformStatus.Azure.NetworkResourceGroupName = nrg
		}
		if installConfig.Config.Platform.Azure.CloudName == azure.StackCloud {
			config.Status.PlatformStatus.Azure.ARMEndpoint = installConfig.Config.Platform.Azure.ARMEndpoint
		}
	case alibabacloud.Name:
		config.Spec.PlatformSpec.Type = configv1.AlibabaCloudPlatformType
		config.Status.PlatformStatus.AlibabaCloud = &configv1.AlibabaCloudPlatformStatus{
			Region:          installConfig.Config.Platform.AlibabaCloud.Region,
			ResourceGroupID: installConfig.Config.Platform.AlibabaCloud.ResourceGroupID,
		}
	case baremetal.Name:
		config.Spec.PlatformSpec.Type = configv1.BareMetalPlatformType
		config.Status.PlatformStatus.BareMetal = &configv1.BareMetalPlatformStatus{
			APIServerInternalIP:  installConfig.Config.Platform.BareMetal.APIVIPs[0],
			IngressIP:            installConfig.Config.Platform.BareMetal.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.Platform.BareMetal.APIVIPs,
			IngressIPs:           installConfig.Config.Platform.BareMetal.IngressVIPs,
		}
	case gcp.Name:
		config.Spec.PlatformSpec.Type = configv1.GCPPlatformType
		config.Status.PlatformStatus.GCP = &configv1.GCPPlatformStatus{
			ProjectID: installConfig.Config.Platform.GCP.ProjectID,
			Region:    installConfig.Config.Platform.GCP.Region,
		}
		uidConfigMap := gcpmanifests.CloudControllerUID(clusterID.InfraID)
		content, err := yaml.Marshal(uidConfigMap)
		if err != nil {
			return errors.Wrapf(err, "cannot marshal GCP cloud controller UID config map")
		}
		i.FileList = append(i.FileList, &asset.File{
			Filename: cloudControllerUIDFilename,
			Data:     content,
		})
	case ibmcloud.Name:
		config.Spec.PlatformSpec.Type = configv1.IBMCloudPlatformType
		var cisInstanceCRN, dnsInstanceCRN string
		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			dnsInstance, err := installConfig.IBMCloud.DNSInstance(context.TODO())
			if err != nil {
				return errors.Wrap(err, "cannot retrieve IBM DNS Services instance CRN")
			}
			dnsInstanceCRN = dnsInstance.CRN
		} else {
			crn, err := installConfig.IBMCloud.CISInstanceCRN(context.TODO())
			if err != nil {
				return errors.Wrap(err, "cannot retrieve IBM Cloud Internet Services instance CRN")
			}
			cisInstanceCRN = crn
		}
		config.Status.PlatformStatus.IBMCloud = &configv1.IBMCloudPlatformStatus{
			Location:          installConfig.Config.Platform.IBMCloud.Region,
			ResourceGroupName: installConfig.Config.Platform.IBMCloud.ClusterResourceGroupName(clusterID.InfraID),
			CISInstanceCRN:    cisInstanceCRN,
			DNSInstanceCRN:    dnsInstanceCRN,
			ProviderType:      configv1.IBMCloudProviderTypeVPC,
		}
	case libvirt.Name:
		config.Spec.PlatformSpec.Type = configv1.LibvirtPlatformType
	case none.Name:
		config.Spec.PlatformSpec.Type = configv1.NonePlatformType
	case openstack.Name:
		config.Spec.PlatformSpec.Type = configv1.OpenStackPlatformType
		config.Status.PlatformStatus.OpenStack = &configv1.OpenStackPlatformStatus{
			APIServerInternalIP:  installConfig.Config.OpenStack.APIVIPs[0],
			IngressIP:            installConfig.Config.OpenStack.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.OpenStack.APIVIPs,
			IngressIPs:           installConfig.Config.OpenStack.IngressVIPs,
		}
	case vsphere.Name:
		config.Spec.PlatformSpec.Type = configv1.VSpherePlatformType
		if len(installConfig.Config.VSphere.APIVIPs) > 0 {
			config.Status.PlatformStatus.VSphere = &configv1.VSpherePlatformStatus{
				APIServerInternalIP:  installConfig.Config.VSphere.APIVIPs[0],
				IngressIP:            installConfig.Config.VSphere.IngressVIPs[0],
				APIServerInternalIPs: installConfig.Config.VSphere.APIVIPs,
				IngressIPs:           installConfig.Config.VSphere.IngressVIPs,
			}
		}
		config.Spec.PlatformSpec.VSphere = vsphereinfra.GetInfraPlatformSpec(installConfig)

		if _, exists := cloudproviderconfig.ConfigMap.Data["vsphere.conf"]; exists {
			cloudProviderConfigMapKey = "vsphere.conf"
		}

	case ovirt.Name:
		config.Spec.PlatformSpec.Type = configv1.OvirtPlatformType
		config.Status.PlatformStatus.Ovirt = &configv1.OvirtPlatformStatus{
			APIServerInternalIP:  installConfig.Config.Ovirt.APIVIPs[0],
			IngressIP:            installConfig.Config.Ovirt.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.Ovirt.APIVIPs,
			IngressIPs:           installConfig.Config.Ovirt.IngressVIPs,
		}
	case powervs.Name:
		config.Spec.PlatformSpec.Type = configv1.PowerVSPlatformType
		var cisInstanceCRN, dnsInstanceCRN string
		var err error
		switch installConfig.Config.Publish {
		case types.InternalPublishingStrategy:
			dnsInstanceCRN, err = installConfig.PowerVS.DNSInstanceCRN(context.TODO())
			if err != nil {
				return errors.Wrapf(err, "failed to get instance CRN")
			}
		case types.ExternalPublishingStrategy:
			cisInstanceCRN, err = installConfig.PowerVS.CISInstanceCRN(context.TODO())
			if err != nil {
				return errors.Wrapf(err, "failed to get instance CRN")
			}
		default:
			return errors.New("unknown publishing strategy")
		}
		config.Status.PlatformStatus.PowerVS = &configv1.PowerVSPlatformStatus{
			Region:         installConfig.Config.Platform.PowerVS.Region,
			Zone:           installConfig.Config.Platform.PowerVS.Zone,
			CISInstanceCRN: cisInstanceCRN,
			DNSInstanceCRN: dnsInstanceCRN,
		}
	case nutanix.Name:
		nutanixPlatform := installConfig.Config.Nutanix

		// Retrieve the prism element name
		var peName string
		if len(nutanixPlatform.PrismElements[0].Name) == 0 {
			nc, err := nutanix.CreateNutanixClient(context.Background(),
				nutanixPlatform.PrismCentral.Endpoint.Address,
				strconv.Itoa(int(nutanixPlatform.PrismCentral.Endpoint.Port)),
				nutanixPlatform.PrismCentral.Username,
				nutanixPlatform.PrismCentral.Password)
			if err != nil {
				return errors.Wrapf(err, "unable to connect to Prism Central %s", nutanixPlatform.PrismCentral.Endpoint.Address)
			}
			pe, err := nc.V3.GetCluster(nutanixPlatform.PrismElements[0].UUID)
			if err != nil {
				return errors.Wrapf(err, "fail to find the Prism Element (cluster) with uuid %s", nutanixPlatform.PrismElements[0].UUID)
			}
			peName = *pe.Spec.Name
		} else {
			peName = nutanixPlatform.PrismElements[0].Name
		}

		config.Spec.PlatformSpec.Type = configv1.NutanixPlatformType
		config.Spec.PlatformSpec.Nutanix = &configv1.NutanixPlatformSpec{
			PrismCentral: configv1.NutanixPrismEndpoint{
				Address: nutanixPlatform.PrismCentral.Endpoint.Address,
				Port:    nutanixPlatform.PrismCentral.Endpoint.Port,
			},
			PrismElements: []configv1.NutanixPrismElementEndpoint{{
				Name: peName,
				Endpoint: configv1.NutanixPrismEndpoint{
					Address: nutanixPlatform.PrismElements[0].Endpoint.Address,
					Port:    nutanixPlatform.PrismElements[0].Endpoint.Port,
				},
			}},
		}

		if len(installConfig.Config.Nutanix.APIVIPs) > 0 {
			config.Status.PlatformStatus.Nutanix = &configv1.NutanixPlatformStatus{
				APIServerInternalIP:  installConfig.Config.Nutanix.APIVIPs[0],
				IngressIP:            installConfig.Config.Nutanix.IngressVIPs[0],
				APIServerInternalIPs: installConfig.Config.Nutanix.APIVIPs,
				IngressIPs:           installConfig.Config.Nutanix.IngressVIPs,
			}
		}
	default:
		config.Spec.PlatformSpec.Type = configv1.NonePlatformType
	}
	config.Status.Platform = config.Spec.PlatformSpec.Type
	config.Status.PlatformStatus.Type = config.Spec.PlatformSpec.Type

	if cloudproviderconfig.ConfigMap != nil {
		// set the configmap reference.
		config.Spec.CloudConfig = configv1.ConfigMapFileReference{Name: cloudproviderconfig.ConfigMap.Name, Key: cloudProviderConfigMapKey}
		i.FileList = append(i.FileList, cloudproviderconfig.File)
	}

	if trustbundleconfig.ConfigMap != nil {
		i.FileList = append(i.FileList, trustbundleconfig.Files()...)
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal config: %#v", config)
	}
	i.FileList = append(i.FileList, &asset.File{
		Filename: infraCfgFilename,
		Data:     configData,
	})
	return nil
}

// Files returns the files generated by the asset.
func (i *Infrastructure) Files() []*asset.File {
	return i.FileList
}

// Load returns false since this asset is not written to disk by the installer.
func (i *Infrastructure) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
