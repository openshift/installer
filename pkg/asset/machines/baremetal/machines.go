// Package baremetal generates Machine objects for bare metal.
package baremetal

import (
	"fmt"

	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.BareMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider, err := provider(platform, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
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
				// we don't need to set Versions, because we control those via cluster operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(platform *baremetal.Platform, userDataSecret string) (*baremetalprovider.BareMetalMachineProviderSpec, error) {
	config := &baremetalprovider.BareMetalMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "baremetal.cluster.k8s.io/v1alpha1",
			Kind:       "BareMetalMachineProviderSpec",
		},
		CustomDeploy: &baremetalprovider.CustomDeploy{
			Method: "install_coreos",
		},
		UserData: &corev1.SecretReference{Name: userDataSecret},
	}
	return config, nil
}
