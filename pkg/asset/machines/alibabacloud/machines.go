// Package alibabacloud generates Machine objects for alibabacloud.
package alibabacloud

import (
	"fmt"

	alibabacloudprovider "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string, resourceTags map[string]string) ([]machineapi.Machine, error) {
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
		provider, err := provider(clusterID, platform, mpool, azIndex, role, userDataSecret, resourceTags)
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
	resourceTags map[string]string,
) (*alibabacloudprovider.AlibabaCloudMachineProviderConfig, error) {
	az := mpool.Zones[azIdx]

	var resourceGroup string
	if platform.ResourceGroupID != "" {
		resourceGroup = platform.ResourceGroupID
	} else {
		return nil, errors.Errorf("Parameter 'ResourceGroup' is empty")
	}

	tags, err := tagsFromResourceTags(clusterID, resourceTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create alibabacloudprovider.Tag from Tags")
	}

	return &alibabacloudprovider.AlibabaCloudMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "alibabacloudmachineproviderconfig.openshift.io/v1beta1",
			Kind:       "AlibabaCloudMachineProviderConfig",
		},
		ImageID:            mpool.ImageID,
		InstanceType:       mpool.InstanceType,
		SystemDiskCategory: string(mpool.SystemDiskCategory),
		SystemDiskSize:     mpool.SystemDiskSize,
		RegionID:           platform.Region,
		ResourceGroupID:    resourceGroup,
		ZoneID:             az,
		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret:  &corev1.LocalObjectReference{Name: "alibabacloud-credentials"},
		Tags:               tags,
	}, nil
}

func tagsFromResourceTags(clusterID string, resourceTags map[string]string) ([]alibabacloudprovider.Tag, error) {
	tags := []alibabacloudprovider.Tag{
		{Key: fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), Value: "owned"},
		{Key: "OCP", Value: "ISV Integration"},
	}
	forbiddenTags := sets.NewString()
	for idx := range tags {
		forbiddenTags.Insert(tags[idx].Key)
	}
	for k, v := range resourceTags {
		if forbiddenTags.Has(k) {
			return nil, fmt.Errorf("user tags may not clobber %s", k)
		}
		tags = append(tags, alibabacloudprovider.Tag{Key: k, Value: v})
	}
	return tags, nil
}
