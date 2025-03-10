// Package azure generates Machine objects for azure.
package azure

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

const (
	cloudsSecret          = "azure-cloud-credentials"
	cloudsSecretNamespace = "openshift-machine-api"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string, capabilities map[string]string, useImageGallery bool) ([]machineapi.Machine, *machinev1.ControlPlaneMachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != azure.Name {
		return nil, nil, fmt.Errorf("non-Azure configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, nil, fmt.Errorf("non-Azure machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Azure
	mpool := pool.Platform.Azure

	if len(mpool.Zones) == 0 {
		// if no azs are given we set to []string{""} for convenience over later operations.
		// It means no-zoned for the machine API
		mpool.Zones = []string{""}
	}
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	machineSetProvider := &machineapi.AzureMachineProviderSpec{}
	for idx := int64(0); idx < total; idx++ {
		var azIndex int
		if len(azs) > 0 {
			azIndex = int(idx) % len(azs)
		}
		provider, err := provider(platform, mpool, osImage, userDataSecret, clusterID, role, &azIndex, capabilities, useImageGallery)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}
		*machineSetProvider = *provider
		machines = append(machines, machine)
	}

	replicas := int32(total)
	failureDomains := []machinev1.AzureFailureDomain{}

	if len(mpool.Zones) > 1 {
		for _, zone := range mpool.Zones {
			domain := machinev1.AzureFailureDomain{
				Zone: zone,
			}

			failureDomains = append(failureDomains, domain)
		}
		machineSetProvider.Zone = ""
	} else if len(mpool.Zones) == 1 {
		machineSetProvider.Zone = mpool.Zones[0]
	}

	controlPlaneMachineSet := &machinev1.ControlPlaneMachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1",
			Kind:       "ControlPlaneMachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      "cluster",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1.ControlPlaneMachineSetSpec{
			Replicas: &replicas,
			State:    machinev1.ControlPlaneMachineSetStateActive,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
					"machine.openshift.io/cluster-api-cluster":      clusterID,
				},
			},
			Template: machinev1.ControlPlaneMachineSetTemplate{
				MachineType: machinev1.OpenShiftMachineV1Beta1MachineType,
				OpenShiftMachineV1Beta1Machine: &machinev1.OpenShiftMachineV1Beta1MachineTemplate{
					ObjectMeta: machinev1.ControlPlaneMachineSetTemplateObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: machineapi.MachineSpec{
						ProviderSpec: machineapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: machineSetProvider},
						},
					},
				},
			},
		},
	}

	if len(failureDomains) > 0 {
		controlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains = &machinev1.FailureDomains{
			Platform: v1.AzurePlatformType,
			Azure:    &failureDomains,
		}
	}

	return machines, controlPlaneMachineSet, nil
}

func provider(platform *azure.Platform, mpool *azure.MachinePool, osImage string, userDataSecret string, clusterID string, role string, azIdx *int, capabilities map[string]string, useImageGallery bool) (*machineapi.AzureMachineProviderSpec, error) {
	var az string
	if len(mpool.Zones) > 0 && azIdx != nil {
		az = mpool.Zones[*azIdx]
	}

	hyperVGen, err := icazure.GetHyperVGenerationVersion(capabilities, "")
	if err != nil {
		return nil, err
	}

	if mpool.VMNetworkingType == "" {
		acceleratedNetworking := icazure.GetVMNetworkingCapability(capabilities)
		if acceleratedNetworking {
			mpool.VMNetworkingType = string(azure.VMnetworkingTypeAccelerated)
		} else {
			logrus.Infof("Instance type %s does not support Accelerated Networking. Using Basic Networking instead.", mpool.InstanceType)
		}
	}
	rg := platform.ClusterResourceGroupName(clusterID)

	var image machineapi.Image
	if mpool.OSImage.Publisher != "" {
		image.Type = machineapi.AzureImageTypeMarketplaceWithPlan
		if mpool.OSImage.Plan == azure.ImageNoPurchasePlan {
			image.Type = machineapi.AzureImageTypeMarketplaceNoPlan
		}
		image.Publisher = mpool.OSImage.Publisher
		image.Offer = mpool.OSImage.Offer
		image.SKU = mpool.OSImage.SKU
		image.Version = mpool.OSImage.Version
	} else if useImageGallery {
		// image gallery names cannot have dashes
		galleryName := strings.ReplaceAll(clusterID, "-", "_")
		id := clusterID
		if hyperVGen == "V2" {
			id += "-gen2"
		}
		imageID := fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/galleries/gallery_%s/images/%s/versions/latest", rg, galleryName, id)
		image.ResourceID = imageID
	} else {
		imageID := fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/images/%s", rg, clusterID)
		if hyperVGen == "V2" && platform.CloudName != azure.StackCloud {
			imageID += "-gen2"
		}
		image.ResourceID = imageID
	}

	networkResourceGroup, virtualNetwork, subnet, err := getNetworkInfo(platform, clusterID, role)
	if err != nil {
		return nil, err
	}

	if mpool.OSDisk.DiskType == "" {
		mpool.OSDisk.DiskType = "Premium_LRS"
	}

	publicLB := clusterID
	if platform.OutboundType == azure.UserDefinedRoutingOutboundType {
		publicLB = ""
	}

	managedIdentity := ""
	if len(mpool.Identity.UserAssignedIdentities) > 0 {
		managedIdentity = mpool.Identity.UserAssignedIdentities[0].ProviderID()
	}

	var diskEncryptionSet *machineapi.DiskEncryptionSetParameters
	if mpool.OSDisk.DiskEncryptionSet != nil {
		diskEncryptionSet = &machineapi.DiskEncryptionSetParameters{
			ID: mpool.OSDisk.DiskEncryptionSet.ToID(),
		}
	}

	var diskSecurityProfile machineapi.VMDiskSecurityProfile
	if mpool.OSDisk.SecurityProfile != nil && mpool.OSDisk.SecurityProfile.SecurityEncryptionType != "" {
		diskSecurityProfile = machineapi.VMDiskSecurityProfile{
			SecurityEncryptionType: machineapi.SecurityEncryptionTypes(mpool.OSDisk.SecurityProfile.SecurityEncryptionType),
		}

		if mpool.OSDisk.SecurityProfile.DiskEncryptionSet != nil {
			diskSecurityProfile.DiskEncryptionSet = machineapi.DiskEncryptionSetParameters{
				ID: mpool.OSDisk.SecurityProfile.DiskEncryptionSet.ToID(),
			}
		}
	}

	securityProfile := generateSecurityProfile(mpool)

	ultraSSDCapability := machineapi.AzureUltraSSDCapabilityState(mpool.UltraSSDCapability)

	spec := &machineapi.AzureMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "AzureMachineProviderSpec",
		},
		UserDataSecret:    &corev1.SecretReference{Name: userDataSecret},
		CredentialsSecret: &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
		Location:          platform.Region,
		VMSize:            mpool.InstanceType,
		Image:             image,
		OSDisk: machineapi.OSDisk{
			OSType:     "Linux",
			DiskSizeGB: mpool.OSDisk.DiskSizeGB,
			ManagedDisk: machineapi.OSDiskManagedDiskParameters{
				StorageAccountType: mpool.OSDisk.DiskType,
				DiskEncryptionSet:  diskEncryptionSet,
				SecurityProfile:    diskSecurityProfile,
			},
		},
		SecurityProfile:       securityProfile,
		UltraSSDCapability:    ultraSSDCapability,
		Zone:                  az,
		Subnet:                subnet,
		ManagedIdentity:       managedIdentity,
		Vnet:                  virtualNetwork,
		ResourceGroup:         rg,
		NetworkResourceGroup:  networkResourceGroup,
		PublicLoadBalancer:    publicLB,
		AcceleratedNetworking: getVMNetworkingType(mpool.VMNetworkingType),
		Tags:                  platform.UserTags,
	}

	if platform.CloudName == azure.StackCloud {
		spec.AvailabilitySet = fmt.Sprintf("%s-cluster", clusterID)
	}

	return spec, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, controlPlane *machinev1.ControlPlaneMachineSet, clusterID string) error {
	internalLB := fmt.Sprintf("%s-internal", clusterID)

	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*machineapi.AzureMachineProviderSpec)
		providerSpec.InternalLoadBalancer = internalLB
	}
	providerSpec, ok := controlPlane.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value.Object.(*machineapi.AzureMachineProviderSpec)
	if !ok {
		return errors.New("Unable to set internal load balancers to control plane machine set")
	}
	providerSpec.InternalLoadBalancer = internalLB
	return nil
}

func getNetworkInfo(platform *azure.Platform, clusterID, role string) (string, string, string, error) {
	networkResourceGroupName := platform.NetworkResourceGroupName
	if platform.VirtualNetwork == "" {
		networkResourceGroupName = platform.ClusterResourceGroupName(clusterID)
	}

	switch role {
	case "worker":
		return networkResourceGroupName, platform.VirtualNetworkName(clusterID), platform.ComputeSubnetName(clusterID), nil
	case "master":
		return networkResourceGroupName, platform.VirtualNetworkName(clusterID), platform.ControlPlaneSubnetName(clusterID), nil
	default:
		return "", "", "", fmt.Errorf("unrecognized machine role %s", role)
	}
}

// getVMNetworkingType should set the correct capability for instance type
func getVMNetworkingType(value string) bool {
	return value == string(azure.VMnetworkingTypeAccelerated)
}

func generateSecurityProfile(mpool *azure.MachinePool) *machineapi.SecurityProfile {
	securityProfile := &machineapi.SecurityProfile{}

	if mpool.EncryptionAtHost {
		securityProfile.EncryptionAtHost = &mpool.EncryptionAtHost
	}

	if mpool.Settings != nil && mpool.Settings.SecurityType != "" {
		securityProfile.Settings = machineapi.SecuritySettings{
			SecurityType: machineapi.SecurityTypes(mpool.Settings.SecurityType),
		}

		var uefiSettings machineapi.UEFISettings
		if securityProfile.Settings.SecurityType == machineapi.SecurityTypesTrustedLaunch {
			if mpool.Settings.TrustedLaunch != nil && mpool.Settings.TrustedLaunch.UEFISettings != nil {
				if sb := mpool.Settings.TrustedLaunch.UEFISettings.SecureBoot; sb != nil {
					uefiSettings.SecureBoot = machineapi.SecureBootPolicy(*sb)
				}
				if vtpm := mpool.Settings.TrustedLaunch.UEFISettings.VirtualizedTrustedPlatformModule; vtpm != nil {
					uefiSettings.VirtualizedTrustedPlatformModule = machineapi.VirtualizedTrustedPlatformModulePolicy(*vtpm)
				}

				securityProfile.Settings.TrustedLaunch = &machineapi.TrustedLaunch{
					UEFISettings: uefiSettings,
				}
			}
		} else if securityProfile.Settings.SecurityType == machineapi.SecurityTypesConfidentialVM {
			if mpool.Settings.ConfidentialVM != nil && mpool.Settings.ConfidentialVM.UEFISettings != nil {
				if sb := mpool.Settings.ConfidentialVM.UEFISettings.SecureBoot; sb != nil {
					uefiSettings.SecureBoot = machineapi.SecureBootPolicy(*sb)
				}
				if vtpm := mpool.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule; vtpm != nil {
					uefiSettings.VirtualizedTrustedPlatformModule = machineapi.VirtualizedTrustedPlatformModulePolicy(*vtpm)
				}

				securityProfile.Settings.ConfidentialVM = &machineapi.ConfidentialVM{
					UEFISettings: uefiSettings,
				}
			}
		}
	}

	return securityProfile
}
