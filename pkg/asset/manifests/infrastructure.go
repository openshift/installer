package manifests

import (
	"path/filepath"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
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

	if installConfig.Config.ControlPlane.Replicas != nil && *installConfig.Config.ControlPlane.Replicas < 3 {
		config.Status.ControlPlaneTopology = configv1.SingleReplicaTopologyMode
	} else {
		config.Status.ControlPlaneTopology = configv1.HighlyAvailableTopologyMode
	}

	numOfWorkers := int64(0)
	for _, mp := range installConfig.Config.Compute {
		if mp.Replicas != nil {
			numOfWorkers += *mp.Replicas
		}
	}
	switch numOfWorkers {
	case 0:
		config.Status.InfrastructureTopology = config.Status.ControlPlaneTopology
	case 1:
		config.Status.InfrastructureTopology = configv1.SingleReplicaTopologyMode
	default:
		config.Status.InfrastructureTopology = configv1.HighlyAvailableTopologyMode
	}

	switch installConfig.Config.Platform.Name() {
	case aws.Name:
		config.Spec.PlatformSpec.Type = configv1.AWSPlatformType
		config.Spec.PlatformSpec.AWS = &configv1.AWSPlatformSpec{}

		var resourceTags []configv1.AWSResourceTag
		if installConfig.Config.AWS.ExperimentalPropagateUserTag {
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
	case baremetal.Name:
		config.Spec.PlatformSpec.Type = configv1.BareMetalPlatformType
		config.Status.PlatformStatus.BareMetal = &configv1.BareMetalPlatformStatus{
			APIServerInternalIP: installConfig.Config.Platform.BareMetal.APIVIP,
			IngressIP:           installConfig.Config.Platform.BareMetal.IngressVIP,
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
	case libvirt.Name:
		config.Spec.PlatformSpec.Type = configv1.LibvirtPlatformType
	case none.Name:
		config.Spec.PlatformSpec.Type = configv1.NonePlatformType
	case openstack.Name:
		config.Spec.PlatformSpec.Type = configv1.OpenStackPlatformType
		config.Status.PlatformStatus.OpenStack = &configv1.OpenStackPlatformStatus{
			APIServerInternalIP: installConfig.Config.OpenStack.APIVIP,
			IngressIP:           installConfig.Config.OpenStack.IngressVIP,
		}
	case vsphere.Name:
		config.Spec.PlatformSpec.Type = configv1.VSpherePlatformType
		if installConfig.Config.VSphere.APIVIP != "" {
			config.Status.PlatformStatus.VSphere = &configv1.VSpherePlatformStatus{
				APIServerInternalIP: installConfig.Config.VSphere.APIVIP,
				IngressIP:           installConfig.Config.VSphere.IngressVIP,
			}
		}
	case ovirt.Name:
		config.Spec.PlatformSpec.Type = configv1.OvirtPlatformType
		config.Status.PlatformStatus.Ovirt = &configv1.OvirtPlatformStatus{
			APIServerInternalIP: installConfig.Config.Ovirt.APIVIP,
			IngressIP:           installConfig.Config.Ovirt.IngressVIP,
		}
	case kubevirt.Name:
		config.Spec.PlatformSpec.Type = configv1.KubevirtPlatformType
		config.Status.PlatformStatus.Kubevirt = &configv1.KubevirtPlatformStatus{
			APIServerInternalIP: installConfig.Config.Kubevirt.APIVIP,
			IngressIP:           installConfig.Config.Kubevirt.IngressVIP,
		}
	default:
		config.Spec.PlatformSpec.Type = configv1.NonePlatformType
	}
	config.Status.Platform = config.Spec.PlatformSpec.Type
	config.Status.PlatformStatus.Type = config.Spec.PlatformSpec.Type

	if cloudproviderconfig.ConfigMap != nil {
		// set the configmap reference.
		config.Spec.CloudConfig = configv1.ConfigMapFileReference{Name: cloudproviderconfig.ConfigMap.Name, Key: cloudProviderConfigDataKey}
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
