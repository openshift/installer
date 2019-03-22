// Package azure generates Machine objects for azure.
package azure

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
)

const (
	// TODO(flaper87): We're choosing to hardcode these values to make
	// the environment more predictable. We expect there to a secret
	// named `azure-credentials` and a cloud named `AzurePublicCloud` in
	// the clouds file stored in this secret.
	cloudsSecret          = "azure-credentials"
	cloudsSecretNamespace = "kube-system"

	// CloudName is a constant containing the name of the cloud used in the internal cloudsSecret
	// TODO: parameterized this to support non-public cloud regions and azure stack
	CloudName = "AzurePublicCloud"
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
		az := ""
		provider, err := provider(clusterID, platform, mpool, osImage, az, role, userDataSecret)
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

func provider(clusterID string, platform *azure.Platform, mpool *azure.MachinePool, osImage string, az string, role, userDataSecret string) (*azureprovider.AzureMachineProviderSpec, error) {

	return &azureprovider.AzureMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "azureprovider.k8s.io/v1alpha1",
			Kind:       "AzureMachineProviderSpec",
		},
		Roles:    []azureprovider.MachineRole{azureprovider.MachineRole(azureprovider.Node)},
		Location: platform.Region,
		VMSize:   mpool.InstanceType,
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	//TODO
}
