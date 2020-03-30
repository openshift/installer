// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/utils/openstack/clientconfig"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
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
	platform := config.Platform.OpenStack

	az := ""
	trunk := platform.TrunkSupport

	provider := generateProvider(clusterID, platform, pool.Platform.OpenStack, osImage, az, role, userDataSecret, trunk)

	if role == "master" {
		sg, err := getOrCreateServerGroup(platform.Cloud, clusterID+"-"+role, "soft-anti-affinity")
		if err != nil {
			return nil, err
		}

		provider.ServerGroupID = sg.ID
	}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	machines := make([]machineapi.Machine, 0, total)
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
				// we don't need to set Versions, because we control those via operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func generateProvider(clusterID string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, az string, role, userDataSecret string, trunk string) *openstackprovider.OpenstackProviderSpec {
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
		AvailabilityZone: az,
		SecurityGroups:   securityGroups,
		Trunk:            trunkSupportBoolean(trunk),
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
	return &spec
}

func trunkSupportBoolean(trunkSupport string) (result bool) {
	if trunkSupport == "1" {
		result = true
	} else {
		result = false
	}
	return
}

// getOrCreateServerGroup gets a Nova server group with the given name and
// policy or creates a new one if it doesn't exist.
//
// https://docs.openstack.org/api-ref/compute/?expanded=create-server-group-detail#server-groups-os-server-groups
func getOrCreateServerGroup(cloud, serverGroupName, policy string) (*servergroups.ServerGroup, error) {
	conn, err := clientconfig.NewServiceClient(
		"compute",
		&clientconfig.ClientOpts{
			Cloud: cloud,
		},
	)
	if err != nil {
		return nil, err
	}

	// Microversion "2.15" is the first that supports "soft"-anti-affinity.
	// Note that microversions starting from "2.64" use a new field
	// accepting policies as a string instead of an array.
	conn.Microversion = "2.15"

	allPages, err := servergroups.List(conn).AllPages()
	if err != nil {
		return nil, err
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		return nil, err
	}

	for _, serverGroup := range allServerGroups {
		// Reuse the server group if it already exists and is empty
		if serverGroup.Name == serverGroupName && len(serverGroup.Members) == 0 {
			return &serverGroup, nil
		}
	}

	// Create a new group otherwise
	return servergroups.Create(conn, &servergroups.CreateOpts{
		Name:     serverGroupName,
		Policies: []string{policy},
	}).Extract()
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	/*for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}*/
}
