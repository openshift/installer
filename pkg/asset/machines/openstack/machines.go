// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	netext "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/utils/openstack/clientconfig"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

const (
	// TODO(flaper87): We're choosing to hardcode these values to make
	// the environment more predictable. We expect there to a secret
	// named `openstack-credentials` and a cloud named `openstack` in
	// the clouds file stored in this secret.
	cloudsSecret          = "openstack-cloud-credentials"
	cloudsSecretNamespace = "openshift-machine-api"

	// CloudName is a constant containing the name of the cloud used in the internal cloudsSecret
	CloudName = "openstack"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}

	mpool := pool.Platform.OpenStack
	platform := config.Platform.OpenStack
	trunkSupport, err := checkNetworkExtensionAvailability(platform.Cloud, "trunk", nil)
	if err != nil {
		return nil, err
	}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	machines := make([]machineapi.Machine, 0, total)
	providerConfigs := map[string]*openstackprovider.OpenstackProviderSpec{}
	for idx := int64(0); idx < total; idx++ {
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		var provider *openstackprovider.OpenstackProviderSpec

		if _, ok := providerConfigs[zone]; !ok {
			provider, err = generateProvider(
				clusterID,
				platform,
				mpool,
				osImage,
				zone,
				role,
				userDataSecret,
				trunkSupport,
			)
			if err != nil {
				return nil, err
			}
			providerConfigs[zone] = provider
		}

		provider = providerConfigs[zone]
		if role == "master" {
			provider.ServerGroupName = clusterID + "-master"
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

func generateProvider(clusterID string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, az string, role, userDataSecret string, trunkSupport bool) (*openstackprovider.OpenstackProviderSpec, error) {
	var networks []openstackprovider.NetworkParam
	if platform.MachinesSubnet != "" {
		networks = []openstackprovider.NetworkParam{{
			Subnets: []openstackprovider.SubnetParam{{
				UUID: platform.MachinesSubnet,
			}}},
		}
	} else {
		networks = []openstackprovider.NetworkParam{{
			Subnets: []openstackprovider.SubnetParam{{
				Filter: openstackprovider.SubnetFilter{
					Name: fmt.Sprintf("%s-nodes", clusterID),
					Tags: fmt.Sprintf("%s=%s", "openshiftClusterID", clusterID),
				}},
			}},
		}
	}
	for _, networkID := range mpool.AdditionalNetworkIDs {
		networks = append(networks, openstackprovider.NetworkParam{
			UUID:                  networkID,
			NoAllowedAddressPairs: true,
		})
	}

	securityGroups := []openstackprovider.SecurityGroupParam{
		{
			Name: fmt.Sprintf("%s-%s", clusterID, role),
		},
	}
	for _, sg := range mpool.AdditionalSecurityGroupIDs {
		securityGroups = append(securityGroups, openstackprovider.SecurityGroupParam{
			UUID: sg,
		})
	}

	spec := openstackprovider.OpenstackProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: openstackprovider.SchemeGroupVersion.String(),
			Kind:       "OpenstackProviderSpec",
		},
		Flavor:           mpool.FlavorName,
		CloudName:        CloudName,
		CloudsSecret:     &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
		UserDataSecret:   &corev1.SecretReference{Name: userDataSecret},
		Networks:         networks,
		PrimarySubnet:    platform.MachinesSubnet,
		AvailabilityZone: az,
		SecurityGroups:   securityGroups,
		Trunk:            trunkSupport,
		Tags: []string{
			fmt.Sprintf("openshiftClusterID=%s", clusterID),
		},
		ServerMetadata: map[string]string{
			"Name":               fmt.Sprintf("%s-%s", clusterID, role),
			"openshiftClusterID": clusterID,
		},
	}
	if mpool.RootVolume != nil {
		spec.RootVolume = &openstackprovider.RootVolume{
			Size:       mpool.RootVolume.Size,
			SourceType: "image",
			SourceUUID: osImage,
			VolumeType: mpool.RootVolume.Type,
		}
	} else {
		spec.Image = osImage
	}
	return &spec, nil
}

func checkNetworkExtensionAvailability(cloud, alias string, opts *clientconfig.ClientOpts) (bool, error) {
	if opts == nil {
		opts = openstackdefaults.DefaultClientOpts(cloud)
	}
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return false, err
	}

	res := netext.Get(conn, alias)
	if res.Err != nil {
		if _, ok := res.Err.(gophercloud.ErrDefault404); ok {
			return false, nil
		}
		return false, res.Err
	}

	return true, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	/*for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}*/
}
