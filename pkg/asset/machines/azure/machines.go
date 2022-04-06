// Package azure generates Machine objects for azure.
package azure

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

const (
	cloudsSecret          = "azure-cloud-credentials"
	cloudsSecretNamespace = "openshift-machine-api"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string, hyperVGen string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != azure.Name {
		return nil, fmt.Errorf("non-Azure configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, fmt.Errorf("non-Azure machine-pool: %q", poolPlatform)
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
	for idx := int64(0); idx < total; idx++ {
		var azIndex int
		if len(azs) > 0 {
			azIndex = int(idx) % len(azs)
		}
		provider, err := provider(platform, mpool, osImage, userDataSecret, clusterID, role, &azIndex, hyperVGen)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
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

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(platform *azure.Platform, mpool *azure.MachinePool, osImage string, userDataSecret string, clusterID string, role string, azIdx *int, hyperVGen string) (*machineapi.AzureMachineProviderSpec, error) {
	var az *string
	if len(mpool.Zones) > 0 && azIdx != nil {
		az = &mpool.Zones[*azIdx]
	}

	rg := platform.ClusterResourceGroupName(clusterID)

	var image machineapi.Image
	if mpool.OSImage.Publisher != "" {
		image.Type = machineapi.AzureImageTypeMarketplaceWithPlan
		image.Publisher = mpool.OSImage.Publisher
		image.Offer = mpool.OSImage.Offer
		image.SKU = mpool.OSImage.SKU
		image.Version = mpool.OSImage.Version
	} else {
		imageID := fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/images/%s", rg, clusterID)
		if hyperVGen == "V2" {
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

	managedIdentity := fmt.Sprintf("%s-identity", clusterID)
	if platform.IsARO() || platform.CloudName == azure.StackCloud {
		managedIdentity = ""
	}

	var diskEncryptionSet *machineapi.DiskEncryptionSetParameters
	if mpool.OSDisk.DiskEncryptionSet != nil {
		diskEncryptionSet = &machineapi.DiskEncryptionSetParameters{
			ID: mpool.OSDisk.DiskEncryptionSet.ToID(),
		}
	}

	var securityProfile *machineapi.SecurityProfile
	if mpool.EncryptionAtHost {
		securityProfile = &machineapi.SecurityProfile{
			EncryptionAtHost: &mpool.EncryptionAtHost,
		}
	}

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
	}

	if platform.CloudName == azure.StackCloud {
		spec.AvailabilitySet = fmt.Sprintf("%s-cluster", clusterID)
	}

	return spec, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	//TODO
}

func getNetworkInfo(platform *azure.Platform, clusterID, role string) (string, string, string, error) {
	if platform.VirtualNetwork == "" {
		return platform.ClusterResourceGroupName(clusterID), fmt.Sprintf("%s-vnet", clusterID), fmt.Sprintf("%s-%s-subnet", clusterID, role), nil
	}

	switch role {
	case "worker":
		return platform.NetworkResourceGroupName, platform.VirtualNetwork, platform.ComputeSubnet, nil
	case "master":
		return platform.NetworkResourceGroupName, platform.VirtualNetwork, platform.ControlPlaneSubnet, nil
	default:
		return "", "", "", fmt.Errorf("unrecognized machine role %s", role)
	}
}

func getVMNetworkingType(value string) bool {
	return value == string(azure.VMnetworkingTypeAccelerated)
}
