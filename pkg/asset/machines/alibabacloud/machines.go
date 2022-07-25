// Package alibabacloud generates Machine objects for alibabacloud.
package alibabacloud

import (
	"fmt"

	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string, resourceTags map[string]string, vswitchMaps map[string]string) ([]machinev1beta1.Machine, error) {
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

	var machines []machinev1beta1.Machine
	for idx := int64(0); idx < total; idx++ {
		zoneID := azs[int(idx)%len(azs)]
		vswitchID := vswitchMaps[zoneID]
		provider, err := provider(clusterID, platform, mpool, zoneID, role, userDataSecret, resourceTags, vswitchID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machinev1beta1.Machine{
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
			Spec: machinev1beta1.MachineSpec{
				ProviderSpec: machinev1beta1.ProviderSpec{
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
	zoneID string,
	role string,
	userDataSecret string,
	resourceTags map[string]string,
	vswitchID string,
) (*machinev1.AlibabaCloudMachineProviderConfig, error) {
	tags, err := tagsFromResourceTags(clusterID, resourceTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create alibabacloudprovider.Tag from Tags")
	}
	sgTags := append(tags, machinev1.Tag{
		Key:   "Name",
		Value: fmt.Sprintf("%s-sg-%s", clusterID, role),
	})
	sgResourceRef := []machinev1.AlibabaResourceReference{
		{
			Type: machinev1.AlibabaResourceReferenceTypeTags,
			Tags: &sgTags,
		},
	}
	config := &machinev1.AlibabaCloudMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1",
			Kind:       "AlibabaCloudMachineProviderConfig",
		},
		ImageID:      mpool.ImageID,
		InstanceType: mpool.InstanceType,
		SystemDisk: machinev1.SystemDiskProperties{
			Category: string(mpool.SystemDiskCategory),
			Size:     int64(mpool.SystemDiskSize),
		},
		RegionID:          platform.Region,
		ZoneID:            zoneID,
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "alibabacloud-credentials"},
		Tags:              tags,
		RAMRoleName:       fmt.Sprintf("%s-role-%s", clusterID, role),
		SecurityGroups:    sgResourceRef,
	}

	if platform.ResourceGroupID != "" {
		config.ResourceGroup = machinev1.AlibabaResourceReference{
			Type: machinev1.AlibabaResourceReferenceTypeID,
			ID:   &platform.ResourceGroupID,
		}
	} else {
		rgname := platform.ClusterResourceGroupName(clusterID)
		config.ResourceGroup = machinev1.AlibabaResourceReference{
			Type: machinev1.AlibabaResourceReferenceTypeName,
			Name: &rgname,
		}
	}

	if vswitchID != "" {
		config.VSwitch = machinev1.AlibabaResourceReference{
			Type: machinev1.AlibabaResourceReferenceTypeID,
			ID:   &vswitchID,
		}
	} else {
		vstags := append(tags, machinev1.Tag{
			Key:   "Name",
			Value: fmt.Sprintf("%s-vswitch-%s", clusterID, zoneID),
		})
		config.VSwitch = machinev1.AlibabaResourceReference{
			Type: machinev1.AlibabaResourceReferenceTypeTags,
			Tags: &vstags,
		}
	}
	return config, nil
}

func tagsFromResourceTags(clusterID string, resourceTags map[string]string) ([]machinev1.Tag, error) {
	tags := []machinev1.Tag{
		{Key: fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), Value: "owned"},
		{Key: "GISV", Value: "ocp"},
		{Key: "sigs.k8s.io/cloud-provider-alibaba/origin", Value: "ocp"},
	}
	forbiddenTags := sets.NewString()
	for idx := range tags {
		forbiddenTags.Insert(tags[idx].Key)
	}
	for k, v := range resourceTags {
		if forbiddenTags.Has(k) {
			return nil, fmt.Errorf("user tags may not clobber %s", k)
		}
		tags = append(tags, machinev1.Tag{Key: k, Value: v})
	}
	return tags, nil
}
