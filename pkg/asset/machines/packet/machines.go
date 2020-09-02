// Package packet generates Machine objects for packet.
package packet

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/packet"

	packetprovider "github.com/packethost/cluster-api-provider-packet/pkg/apis/packetprovider/v1alpha1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != packet.Name {
		return nil, fmt.Errorf("non-packet configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != packet.Name {
		return nil, fmt.Errorf("non-packet machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Packet

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(platform, pool, userDataSecret, osImage)
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

func provider(platform *packet.Platform, pool *types.MachinePool, userDataSecret string, osImage string) *packetprovider.PacketMachineProviderSpec {
	spec := packetprovider.PacketMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "packetprovider.k8s.io/v1alpha1",
			Kind:       "PacketMachineProviderSpec",
		},
		// TODO(displague) which role?
		Roles:     []packetprovider.MachineRole{packetprovider.MasterRole, packetprovider.NodeRole},
		Facility:  []string{platform.FacilityCode},
		OS:        "custom_ipxe",
		ProjectID: platform.ProjectID,
		// TODO(displague) IPXE
		// IPXEScriptURL: osImage / platform.BootstrapOSImage / platform.ClusterOSImage?
		BillingCycle: "hourly",
		MachineType:  "", // TODO(displague) must provide a type
		SshKeys:      []string{},
		// UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		// CredentialsSecret: &corev1.LocalObjectReference{Name: "packet-credentials"},
	}
	return &spec
}
