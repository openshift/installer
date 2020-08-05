// Package vsphere generates Machine objects for vsphere.
package vsphere

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	vsphereapis "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != vsphere.Name {
		return nil, fmt.Errorf("non vsphere configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != vsphere.Name {
		return nil, fmt.Errorf("non-VSphere machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider, err := provider(clusterID, platform, mpool, osImage, userDataSecret)
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

func provider(clusterID string, platform *vsphere.Platform, mpool *vsphere.MachinePool, osImage string, userDataSecret string) (*vsphereapis.VSphereMachineProviderSpec, error) {
	folder := fmt.Sprintf("/%s/vm/%s", platform.Datacenter, clusterID)
	resourcePool := fmt.Sprintf("/%s/host/%s/Resources", platform.Datacenter, platform.Cluster)
	if platform.Folder != "" {
		folder = platform.Folder
	}

	return &vsphereapis.VSphereMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vsphereapis.SchemeGroupVersion.String(),
			Kind:       "VSphereMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "vsphere-cloud-credentials"},
		Template:          osImage,
		Network: vsphereapis.NetworkSpec{
			Devices: []vsphereapis.NetworkDeviceSpec{
				{
					NetworkName: platform.Network,
				},
			},
		},
		Workspace: &vsphereapis.Workspace{
			Server:       platform.VCenter,
			Datacenter:   platform.Datacenter,
			Datastore:    platform.DefaultDatastore,
			Folder:       folder,
			ResourcePool: resourcePool,
		},
		NumCPUs:           mpool.NumCPUs,
		NumCoresPerSocket: mpool.NumCoresPerSocket,
		MemoryMiB:         mpool.MemoryMiB,
		DiskGiB:           mpool.OSDisk.DiskSizeGB,
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
}
