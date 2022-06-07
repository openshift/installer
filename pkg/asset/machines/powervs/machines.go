// Package powervs generates Machine objects for powerVS.
package powervs

import (
	"fmt"

	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != powervs.Name {
		return nil, fmt.Errorf("non-PowerVS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, fmt.Errorf("non-PowerVS machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.PowerVS
	mpool := pool.Platform.PowerVS

	// Only the service instance is guaranteed to exist and be passed via the install config
	// The other two, we should standardize a name including the cluster id.
	image := fmt.Sprintf("rhcos-%s", clusterID)
	var network string
	if platform.ClusterOSImage != "" {
		image = platform.ClusterOSImage
	}
	if platform.PVSNetworkName != "" {
		network = platform.PVSNetworkName
	}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider, err := provider(clusterID, platform, mpool, userDataSecret, image, network)
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

func provider(clusterID string, platform *powervs.Platform, mpool *powervs.MachinePool, userDataSecret string, image string, network string) (*machinev1.PowerVSMachineProviderConfig, error) {

	if clusterID == "" || platform == nil || mpool == nil || userDataSecret == "" || image == "" {
		return nil, fmt.Errorf("invalid value passed to provider")
	}

	dhcpNetRegex := "^DHCPSERVER[0-9a-z]{32}_Private$"

	//Setting only the mandatory parameters
	config := &machinev1.PowerVSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PowerVSMachineProviderConfig",
			APIVersion: machinev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{},
		ServiceInstance: machinev1.PowerVSResource{
			Type: machinev1.PowerVSResourceTypeID,
			ID:   &platform.ServiceInstanceID,
		},
		Image: machinev1.PowerVSResource{
			Type: machinev1.PowerVSResourceTypeName,
			Name: &image,
		},
		UserDataSecret: &machinev1.PowerVSSecretReference{
			Name: userDataSecret,
		},
		CredentialsSecret: &machinev1.PowerVSSecretReference{
			Name: "powervs-credentials",
		},
		SystemType:    mpool.SysType,
		ProcessorType: mpool.ProcType,
		Processors:    mpool.Processors,
		MemoryGiB:     mpool.MemoryGiB,
		KeyPairName:   fmt.Sprintf("%s-key", clusterID),
	}
	if network != "" {
		config.Network = machinev1.PowerVSResource{
			Type: machinev1.PowerVSResourceTypeName,
			Name: &network,
		}
	} else {
		config.Network = machinev1.PowerVSResource{
			Type:  machinev1.PowerVSResourceTypeRegEx,
			RegEx: &dhcpNetRegex,
		}
	}
	return config, nil
}

// ConfigMasters sets the network and boot image IDs
func ConfigMasters(machines []machineapi.Machine, clusterID string) {

}
