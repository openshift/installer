package manifests

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	externalinfra "github.com/openshift/installer/pkg/asset/manifests/external"
	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	nutanixinfra "github.com/openshift/installer/pkg/asset/manifests/nutanix"
	vsphereinfra "github.com/openshift/installer/pkg/asset/manifests/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
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
//
//nolint:gocyclo
func (i *Infrastructure) Generate(ctx context.Context, dependencies asset.Parents) error {
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
	config.Status.CPUPartitioning = determineCPUPartitioning(installConfig.Config)

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
		// If the user has requested the use of a DNS provisioned by them, then OpenShift needs to
		// start an in-cluster DNS for the installation to succeed. The user can then configure their
		// DNS post-install.
		config.Status.PlatformStatus.AWS.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{
			DNSType: configv1.PlatformDefaultDNSType,
		}
		if installConfig.Config.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled {
			config.Status.PlatformStatus.AWS.CloudLoadBalancerConfig.DNSType = configv1.ClusterHostedDNSType
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
		if len(installConfig.Config.Azure.UserTags) > 0 {
			resourceTags := make([]configv1.AzureResourceTag, 0, len(installConfig.Config.Azure.UserTags))
			for k, v := range installConfig.Config.Azure.UserTags {
				resourceTags = append(resourceTags, configv1.AzureResourceTag{Key: k, Value: v})
			}
			config.Status.PlatformStatus.Azure.ResourceTags = resourceTags
		}
	case baremetal.Name:
		config.Spec.PlatformSpec.Type = configv1.BareMetalPlatformType
		config.Spec.PlatformSpec.BareMetal = &configv1.BareMetalPlatformSpec{}
		config.Status.PlatformStatus.BareMetal = &configv1.BareMetalPlatformStatus{
			APIServerInternalIP:  installConfig.Config.Platform.BareMetal.APIVIPs[0],
			IngressIP:            installConfig.Config.Platform.BareMetal.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.Platform.BareMetal.APIVIPs,
			IngressIPs:           installConfig.Config.Platform.BareMetal.IngressVIPs,
			LoadBalancer:         installConfig.Config.Platform.BareMetal.LoadBalancer,
		}
		config.Spec.PlatformSpec.BareMetal.APIServerInternalIPs = types.StringsToIPs(installConfig.Config.Platform.BareMetal.APIVIPs)
		config.Spec.PlatformSpec.BareMetal.IngressIPs = types.StringsToIPs(installConfig.Config.Platform.BareMetal.IngressVIPs)
		config.Spec.PlatformSpec.BareMetal.MachineNetworks = types.MachineNetworksToCIDRs(installConfig.Config.MachineNetwork)
		config.Status.PlatformStatus.BareMetal.MachineNetworks = types.MachineNetworksToCIDRs(installConfig.Config.MachineNetwork)
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
		if len(installConfig.Config.GCP.UserLabels) > 0 {
			resourceLabels := make([]configv1.GCPResourceLabel, len(installConfig.Config.GCP.UserLabels))
			for i, label := range installConfig.Config.GCP.UserLabels {
				resourceLabels[i] = configv1.GCPResourceLabel{Key: label.Key, Value: label.Value}
			}
			config.Status.PlatformStatus.GCP.ResourceLabels = resourceLabels
		}
		if len(installConfig.Config.GCP.UserTags) > 0 {
			resourceTags := make([]configv1.GCPResourceTag, len(installConfig.Config.GCP.UserTags))
			for i, tag := range installConfig.Config.GCP.UserTags {
				resourceTags[i] = configv1.GCPResourceTag{ParentID: tag.ParentID, Key: tag.Key, Value: tag.Value}
			}
			config.Status.PlatformStatus.GCP.ResourceTags = resourceTags
		}
		// If the user has requested the use of a DNS provisioned by them, then OpenShift needs to
		// start an in-cluster DNS for the installation to succeed. The user can then configure their
		// DNS post-install.
		config.Status.PlatformStatus.GCP.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{
			DNSType: configv1.PlatformDefaultDNSType,
		}
		if installConfig.Config.GCP.UserProvisionedDNS == dns.UserProvisionedDNSEnabled {
			config.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType = configv1.ClusterHostedDNSType
		}
	case ibmcloud.Name:
		config.Spec.PlatformSpec.Type = configv1.IBMCloudPlatformType
		var cisInstanceCRN, dnsInstanceCRN string
		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			dnsInstance, err := installConfig.IBMCloud.DNSInstance(ctx)
			if err != nil {
				return errors.Wrap(err, "cannot retrieve IBM DNS Services instance CRN")
			}
			dnsInstanceCRN = dnsInstance.CRN
		} else {
			crn, err := installConfig.IBMCloud.CISInstanceCRN(ctx)
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
			ServiceEndpoints:  installConfig.Config.Platform.IBMCloud.ServiceEndpoints,
		}
	case external.Name:
		config.Spec.PlatformSpec.Type = configv1.ExternalPlatformType
		config.Spec.PlatformSpec.External = externalinfra.GetInfraPlatformSpec(installConfig)
		config.Status.PlatformStatus.External = externalinfra.GetInfraPlatformStatus(installConfig)
	case none.Name:
		config.Spec.PlatformSpec.Type = configv1.NonePlatformType
	case openstack.Name:
		config.Spec.PlatformSpec.Type = configv1.OpenStackPlatformType
		config.Spec.PlatformSpec.OpenStack = &configv1.OpenStackPlatformSpec{}
		config.Status.PlatformStatus.OpenStack = &configv1.OpenStackPlatformStatus{
			APIServerInternalIP:  installConfig.Config.OpenStack.APIVIPs[0],
			IngressIP:            installConfig.Config.OpenStack.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.OpenStack.APIVIPs,
			IngressIPs:           installConfig.Config.OpenStack.IngressVIPs,
			LoadBalancer:         installConfig.Config.OpenStack.LoadBalancer,
		}
		config.Spec.PlatformSpec.OpenStack.APIServerInternalIPs = types.StringsToIPs(installConfig.Config.Platform.OpenStack.APIVIPs)
		config.Spec.PlatformSpec.OpenStack.IngressIPs = types.StringsToIPs(installConfig.Config.Platform.OpenStack.IngressVIPs)
		config.Spec.PlatformSpec.OpenStack.MachineNetworks = types.MachineNetworksToCIDRs(installConfig.Config.MachineNetwork)
		config.Status.PlatformStatus.OpenStack.MachineNetworks = types.MachineNetworksToCIDRs(installConfig.Config.MachineNetwork)
	case vsphere.Name:
		config.Spec.PlatformSpec.Type = configv1.VSpherePlatformType
		config.Spec.PlatformSpec.VSphere = &configv1.VSpherePlatformSpec{}
		if len(installConfig.Config.VSphere.APIVIPs) > 0 {
			config.Status.PlatformStatus.VSphere = &configv1.VSpherePlatformStatus{
				APIServerInternalIP:  installConfig.Config.VSphere.APIVIPs[0],
				IngressIP:            installConfig.Config.VSphere.IngressVIPs[0],
				APIServerInternalIPs: installConfig.Config.VSphere.APIVIPs,
				IngressIPs:           installConfig.Config.VSphere.IngressVIPs,
				LoadBalancer:         installConfig.Config.VSphere.LoadBalancer,
			}
		} else {
			config.Status.PlatformStatus.VSphere = &configv1.VSpherePlatformStatus{}
		}

		config.Spec.PlatformSpec.VSphere = vsphereinfra.GetInfraPlatformSpec(installConfig, clusterID.InfraID)
		if _, exists := cloudproviderconfig.ConfigMap.Data["vsphere.conf"]; exists {
			cloudProviderConfigMapKey = "vsphere.conf"
		}

		config.Status.PlatformStatus.VSphere.MachineNetworks = types.MachineNetworksToCIDRs(installConfig.Config.MachineNetwork)
	case ovirt.Name:
		config.Spec.PlatformSpec.Type = configv1.OvirtPlatformType
		config.Status.PlatformStatus.Ovirt = &configv1.OvirtPlatformStatus{
			APIServerInternalIP:  installConfig.Config.Ovirt.APIVIPs[0],
			IngressIP:            installConfig.Config.Ovirt.IngressVIPs[0],
			APIServerInternalIPs: installConfig.Config.Ovirt.APIVIPs,
			IngressIPs:           installConfig.Config.Ovirt.IngressVIPs,
			LoadBalancer:         installConfig.Config.Ovirt.LoadBalancer,
		}
	case powervs.Name:
		config.Spec.PlatformSpec.Type = configv1.PowerVSPlatformType
		config.Spec.PlatformSpec.PowerVS = &configv1.PowerVSPlatformSpec{}
		var cisInstanceCRN, dnsInstanceCRN string
		var err error
		switch installConfig.Config.Publish {
		case types.InternalPublishingStrategy:
			dnsInstanceCRN, err = installConfig.PowerVS.DNSInstanceCRN(ctx)
			if err != nil {
				return errors.Wrapf(err, "failed to get instance CRN")
			}
		case types.ExternalPublishingStrategy:
			cisInstanceCRN, err = installConfig.PowerVS.CISInstanceCRN(ctx)
			if err != nil {
				return errors.Wrapf(err, "failed to get instance CRN")
			}
		default:
			return errors.New("unknown publishing strategy")
		}
		for _, service := range installConfig.Config.Platform.PowerVS.ServiceEndpoints {
			config.Spec.PlatformSpec.PowerVS.ServiceEndpoints = append(config.Spec.PlatformSpec.PowerVS.ServiceEndpoints, configv1.PowerVSServiceEndpoint{
				Name: service.Name,
				URL:  service.URL,
			})
		}
		if installConfig.Config.Publish == types.InternalPublishingStrategy &&
			(len(installConfig.Config.ImageDigestSources) > 0 || len(installConfig.Config.DeprecatedImageContentSources) > 0) {
			piRegion := installConfig.Config.PowerVS.Region
			vpcRegion, err := powervs.VPCRegionForPowerVSRegion(piRegion)
			if err != nil {
				return fmt.Errorf("failed to determine VPC region: %w", err)
			}
			cosRegion, err := powervs.COSRegionForPowerVSRegion(piRegion)
			if err != nil {
				return fmt.Errorf("failed to determine COS region: %w", err)
			}
			config.Spec.PlatformSpec.PowerVS.ServiceEndpoints = installConfig.PowerVS.SetDefaultPrivateServiceEndpoints(ctx, config.Spec.PlatformSpec.PowerVS.ServiceEndpoints, cosRegion, vpcRegion)
		}
		config.Status.PlatformStatus.PowerVS = &configv1.PowerVSPlatformStatus{
			Region:           installConfig.Config.Platform.PowerVS.Region,
			Zone:             installConfig.Config.Platform.PowerVS.Zone,
			ResourceGroup:    installConfig.Config.Platform.PowerVS.PowerVSResourceGroup,
			CISInstanceCRN:   cisInstanceCRN,
			DNSInstanceCRN:   dnsInstanceCRN,
			ServiceEndpoints: config.Spec.PlatformSpec.PowerVS.ServiceEndpoints,
		}
	case nutanix.Name:
		config.Spec.PlatformSpec.Type = configv1.NutanixPlatformType

		platformSpec, err := nutanixinfra.GetInfrastructureNutanixPlatformSpec(installConfig)
		if err != nil {
			return fmt.Errorf("failed to create Infrastructure manifest Nutanix platformSpec: %w", err)
		}
		config.Spec.PlatformSpec.Nutanix = platformSpec

		if len(installConfig.Config.Nutanix.APIVIPs) > 0 {
			config.Status.PlatformStatus.Nutanix = &configv1.NutanixPlatformStatus{
				APIServerInternalIP:  installConfig.Config.Nutanix.APIVIPs[0],
				IngressIP:            installConfig.Config.Nutanix.IngressVIPs[0],
				APIServerInternalIPs: installConfig.Config.Nutanix.APIVIPs,
				IngressIPs:           installConfig.Config.Nutanix.IngressVIPs,
				LoadBalancer:         installConfig.Config.Nutanix.LoadBalancer,
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
