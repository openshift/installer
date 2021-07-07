// Package alibabacloud generates Machine objects for alibabacloud.
package alibabacloud

import (
	"fmt"

	alibabacloudprovider "github.com/openshift/installer/pkg/tfvars/alibabacloud"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != alibabacloud.Name {
		return nil, fmt.Errorf("non-AlibabaCloud configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != alibabacloud.Name {
		return nil, fmt.Errorf("non-AlibabaCloud machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.AlibabaCloud
	mpool := pool.Platform.AlibabaCloud
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, azIndex, role, userDataSecret)
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

func provider(clusterID string,
	platform *alibabacloud.Platform,
	mpool *alibabacloud.MachinePool,
	azIdx int,
	role string,
	userDataSecret string,
) (*alibabacloudprovider.MachineProviderSpec, error) {
	// az := mpool.Zones[azIdx]

	var resourceGroup string
	if platform.ResourceGroupName != "" {
		resourceGroup = platform.ResourceGroupName
	} else {
		return nil, errors.Errorf("Parameter 'ResourceGroup' is empty")
	}

	return &alibabacloudprovider.MachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "alibabacloudproviderconfig.openshift.io/v1beta1",
			Kind:       "MachineProviderSpec",
		},
		// VPC:           vpc,
		// Tags:          []alibabacloudprovider.TagSpecs{},
		// Image:         fmt.Sprintf("%s-rhcos", clusterID),
		// Profile:       mpool.InstanceType,
		// Region:        platform.Region,
		ResourceGroupID: resourceGroup,
		// Zone:          az,
		// PrimaryNetworkInterface: alibabacloudprovider.NetworkInterface{
		// 	SecurityGroups: []string{securityGroup},
		// },
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "alibabacloud-credentials"},
		// TODO: AlibabaCloud:
	}, nil
}
