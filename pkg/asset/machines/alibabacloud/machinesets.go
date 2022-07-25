package alibabacloud

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string, resourceTags map[string]string, vswitchMaps map[string]string) ([]*machinev1beta1.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != alibabacloud.Name {
		return nil, fmt.Errorf("non-AlibabaCloud configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != alibabacloud.Name {
		return nil, fmt.Errorf("non-AlibabaCloud machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.AlibabaCloud
	mpool := pool.Platform.AlibabaCloud
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machinesets []*machinev1beta1.MachineSet
	for idx, az := range azs {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		vswitchID, ok := vswitchMaps[az]
		if len(vswitchMaps) > 0 && !ok {
			return nil, errors.Errorf("no VSwitch for zone %s", az)
		}
		provider, err := provider(clusterID, platform, mpool, az, role, userDataSecret, resourceTags, vswitchID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		name := fmt.Sprintf("%s-%s-%s", clusterID, pool.Name, strings.TrimPrefix(az, fmt.Sprintf("%s-", platform.Region)))
		mset := &machinev1beta1.MachineSet{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "MachineSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      name,
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster": clusterID,
				},
			},
			Spec: machinev1beta1.MachineSetSpec{
				Replicas: &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Template: machinev1beta1.MachineTemplateSpec{
					ObjectMeta: machinev1beta1.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: machinev1beta1.MachineSpec{
						ProviderSpec: machinev1beta1.ProviderSpec{
							Value: &runtime.RawExtension{Object: provider},
						},
					},
				},
			},
		}
		machinesets = append(machinesets, mset)
	}
	return machinesets, nil
}
