// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	// TODO(flaper87): We're choosing to hardcode these values to make
	// the environment more predictable. We expect there to a secret
	// named `openstack-credentials` and a cloud named `openstack` in
	// the clouds file stored in this secret.
	cloudsSecret          = "openstack-credentials"
	cloudsSecretNamespace = "kube-system"

	// CloudName is a constant containing the name of the cloud used in the internal cloudsSecret
	CloudName = "openstack"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage string, role types.MachineRole, userDataSecret string) ([]machineapi.Machine, error) {
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
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		az := ""
		provider, err := provider(clusterID, clustername, platform, mpool, osImage, az, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-cluster-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clustername,
					"sigs.k8s.io/cluster-api-machine-role": role.ClusterAPIMachineRole(),
					"sigs.k8s.io/cluster-api-machine-type": role.ClusterAPIMachineRole(),
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

func provider(clusterID, clusterName string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, az string, role types.MachineRole, userDataSecret string) (*openstackprovider.OpenstackProviderSpec, error) {
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
		CloudName:      CloudName,
		CloudsSecret:   &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
		UserDataSecret: &corev1.SecretReference{Name: userDataSecret},
		Networks: []openstackprovider.NetworkParam{
			{
				Filter: openstackprovider.Filter{
					Tags: fmt.Sprintf("%s=%s", "openshiftClusterID", clusterID),
				},
			},
		},
		AvailabilityZone: az,
		SecurityGroups:   []string{string(role)},
		// TODO(flaper87): Trunk support missing. Need to add it back
	}, nil
}

// ConfigControlPlane sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigControlPlane(machines []machineapi.Machine, clusterName string) {
	/*for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}*/
}
