// Package powervs generates Machine objects for powerVS.
package powervs

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	powervsprovider "github.com/openshift/cluster-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string, userTags map[string]string) ([]machineapi.Machine, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, fmt.Errorf("non-PowerVS machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.PowerVS
	mpool := pool.Platform.PowerVS
	//@TODO: clearly this needs to be set in the install config or through the rhcos pkg :D
	mpool.ImageID = "rhcos-48-05132021-tier1"
	mpool.NetworkIDs = []string{"pvs-ipi-net"}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider, err := provider(clusterID, platform, mpool, userDataSecret, userTags)
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
			},
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

func provider(clusterID string, platform *powervs.Platform, mpool *powervs.MachinePool, userDataSecret string, userTags map[string]string) (*powervsprovider.PowerVSMachineProviderConfig, error) {

	//Setting only the mandatory parameters
	config := &powervsprovider.PowerVSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PowerVSMachineProviderConfig",
			APIVersion: powervsprovider.GroupVersion.String(),
		},
		ObjectMeta:        metav1.ObjectMeta{},
		ServiceInstanceID: mpool.ServiceInstance,
		ImageID:           mpool.ImageID,
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "powervs-credentials-secret"},
		SysType:           mpool.SysType,
		ProcType:          mpool.ProcType,
		Processors:        fmt.Sprintf("%f", mpool.Processors),
		Memory:            fmt.Sprintf("%d", mpool.Memory),
		NetworkIDs:        mpool.NetworkIDs,
		KeyPairName:       &mpool.KeyPairName,
	}
	return config, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	//TODO: Revisit this later if required, At the moment we don't know how to handle the ingress data.
}
