package manifests

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	infraCrdFilename           = filepath.Join(manifestDir, "cluster-infrastructure-01-crd.yaml")
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
		Status: configv1.InfrastructureStatus{
			InfrastructureName:   clusterID.InfraID,
			APIServerURL:         getAPIServerURL(installConfig.Config),
			APIServerInternalURL: getInternalAPIServerURL(installConfig.Config),
			EtcdDiscoveryDomain:  getEtcdDiscoveryDomain(installConfig.Config),
			PlatformStatus:       &configv1.PlatformStatus{},
		},
	}

	switch installConfig.Config.Platform.Name() {
	case aws.Name:
		config.Status.PlatformStatus.Type = configv1.AWSPlatformType
		config.Status.PlatformStatus.AWS = &configv1.AWSPlatformStatus{
			Region: installConfig.Config.Platform.AWS.Region,
		}
	case azure.Name:
		config.Status.PlatformStatus.Type = configv1.AzurePlatformType
		config.Status.PlatformStatus.Azure = &configv1.AzurePlatformStatus{
			ResourceGroupName: fmt.Sprintf("%s-rg", clusterID.InfraID),
		}
	case baremetal.Name:
		config.Status.PlatformStatus.Type = configv1.BareMetalPlatformType
		config.Status.PlatformStatus.BareMetal = &configv1.BareMetalPlatformStatus{
			APIServerInternalIP: installConfig.Config.Platform.BareMetal.APIVIP,
			NodeDNSIP:           installConfig.Config.Platform.BareMetal.DNSVIP,
			IngressIP:           installConfig.Config.Platform.BareMetal.IngressVIP,
		}
	case gcp.Name:
		config.Status.PlatformStatus.Type = configv1.GCPPlatformType
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
		config.Status.PlatformStatus.Type = configv1.LibvirtPlatformType
	case none.Name:
		config.Status.PlatformStatus.Type = configv1.NonePlatformType
	case openstack.Name:
		config.Status.PlatformStatus.Type = configv1.OpenStackPlatformType
		config.Status.PlatformStatus.OpenStack = &configv1.OpenStackPlatformStatus{
			APIServerInternalIP: installConfig.Config.Platform.OpenStack.APIVIP,
			NodeDNSIP:           installConfig.Config.Platform.OpenStack.DNSVIP,
			IngressIP:           installConfig.Config.Platform.OpenStack.IngressVIP,
		}
	case vsphere.Name:
		config.Status.PlatformStatus.Type = configv1.VSpherePlatformType
	default:
		config.Status.PlatformStatus.Type = configv1.NonePlatformType
	}
	config.Status.Platform = config.Status.PlatformStatus.Type

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
