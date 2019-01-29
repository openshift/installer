// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

const cloudsSecret = "openstack-credentials"

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]clusterapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.OpenStack
	mpool := pool.Platform.OpenStack

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []clusterapi.Machine
	for idx := int64(0); idx < total; idx++ {
		az := ""
		provider, err := provider(clusterID, clustername, platform, mpool, osImage, az, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := clusterapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.k8s.io/v1alpha1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-cluster-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clustername,
					"sigs.k8s.io/cluster-api-machine-role": role,
					"sigs.k8s.io/cluster-api-machine-type": role,
				},
			},
			Spec: clusterapi.MachineSpec{
				ProviderSpec: clusterapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID, clusterName string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, az string, role, userDataSecret string) (*openstackprovider.OpenstackProviderSpec, error) {
	return &openstackprovider.OpenstackProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "openstackproviderconfig.k8s.io/v1alpha1",
			Kind:       "OpenstackProviderSpec",
		},
		Flavor: mpool.FlavorName,
		/*RootVolume: openstackprovider.RootVolume{
			VolumeType: pointer.StringPtr(mpool.Type),
			Size:       pointer.Int64Ptr(int64(mpool.Size)),
		},*/
		Image:          osImage,
		CloudName:      platform.Cloud,
		CloudsSecret:   &corev1.SecretReference{Name: cloudsSecret},
		UserDataSecret: &corev1.SecretReference{Name: userDataSecret},
		Networks: []openstackprovider.NetworkParam{
			{
				Filter: openstackprovider.Filter{
					Tags: fmt.Sprintf("%s=%s", "openshiftClusterID", clusterID),
				},
			},
		},
		AvailabilityZone: az,
		SecurityGroups:   []string{role},
		// TODO(flaper87): Trunk support missing. Need to add it back
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []clusterapi.Machine, clusterName string) {
	/*for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}*/
}
