// Package azure generates Machine objects for azure.
package azure

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
)

const (
	cloudsSecret          = "azure-cloud-credentials"
	cloudsSecretNamespace = "openshift-machine-api"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != azure.Name {
		return nil, fmt.Errorf("non-Azure configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, fmt.Errorf("non-Azure machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Azure
	mpool := pool.Platform.Azure

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider, err := provider(platform, mpool, osImage, userDataSecret, clusterID, role)
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

func provider(platform *azure.Platform, mpool *azure.MachinePool, osImage string, userDataSecret string, clusterID string, role string) (*azureprovider.AzureMachineProviderSpec, error) {
	return &azureprovider.AzureMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "azureprovider.k8s.io/v1alpha1",
			Kind:       "AzureMachineProviderSpec",
		},
		UserDataSecret:    &corev1.SecretReference{Name: userDataSecret},
		CredentialsSecret: &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
		Location:          platform.Region,
		VMSize:            mpool.InstanceType,
		Image: azureprovider.Image{
			ResourceID: osImage,
		},
		OSDisk: azureprovider.OSDisk{
			OSType:     "Linux",
			DiskSizeGB: mpool.OSDisk.DiskSizeGB,
			ManagedDisk: azureprovider.ManagedDisk{
				StorageAccountType: "Premium_LRS",
			},
		},
		Subnet:          fmt.Sprintf("%s-%s-subnet", clusterID, role),
		ManagedIdentity: fmt.Sprintf("%s-identity", clusterID),
		Vnet:            fmt.Sprintf("%s-vnet", clusterID),
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	//TODO
}
