package ibmcloud

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != ibmcloud.Name {
		return nil, fmt.Errorf("non-IBMCloud configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != ibmcloud.Name {
		return nil, fmt.Errorf("non-IBMCloud machine-pool: %q", poolPlatform)
	}
	// platform := config.Platform.IBMCloud
	// mpool := pool.Platform.IBMCloud
	// azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		// azIndex := int(idx) % len(azs)
		// provider, err := provider(clusterID, platform, mpool, osImage, azIndex, role, userDataSecret)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "failed to create provider")
		// }
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
					// TODO: IBM: uncomment and use provider when it's implemented.
					// Value: &runtime.RawExtension{Object: provider},
				},
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

//nolint
func provider(clusterID string, platform *ibmcloud.Platform, mpool *ibmcloud.MachinePool, osImage string, azIdx int, role, userDataSecret string) (*interface{}, error) {
	// TODO: IBM: return cluster api provider, error
	return nil, nil
}
