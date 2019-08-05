// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.AWS
	mpool := pool.Platform.AWS
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, osImage, azIndex, role, userDataSecret)
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
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID string, platform *aws.Platform, mpool *aws.MachinePool, osImage string, azIdx int, role, userDataSecret string) (*awsprovider.AWSMachineProviderConfig, error) {
	az := mpool.Zones[azIdx]
	amiID := osImage
	tags, err := tagsFromUserTags(clusterID, platform.UserTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create awsprovider.TagSpecifications from UserTags")
	}
	return &awsprovider.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "awsproviderconfig.openshift.io/v1beta1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType: mpool.InstanceType,
		BlockDevices: []awsprovider.BlockDeviceMappingSpec{
			{
				EBS: &awsprovider.EBSBlockDeviceSpec{
					VolumeType: pointer.StringPtr(mpool.Type),
					VolumeSize: pointer.Int64Ptr(int64(mpool.Size)),
					Iops:       pointer.Int64Ptr(int64(mpool.IOPS)),
					Encrypted:  pointer.BoolPtr(true),
				},
			},
		},
		AMI:                awsprovider.AWSResourceReference{ID: &amiID},
		Tags:               tags,
		IAMInstanceProfile: &awsprovider.AWSResourceReference{ID: pointer.StringPtr(fmt.Sprintf("%s-%s-profile", clusterID, role))},
		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret:  &corev1.LocalObjectReference{Name: "aws-cloud-credentials"},
		Subnet: awsprovider.AWSResourceReference{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s-private-%s", clusterID, az)},
			}},
		},
		Placement: awsprovider.Placement{Region: platform.Region, AvailabilityZone: az},
		SecurityGroups: []awsprovider.AWSResourceReference{{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s-%s-sg", clusterID, role)},
			}},
		}},
	}, nil
}

func tagsFromUserTags(clusterID string, usertags map[string]string) ([]awsprovider.TagSpecification, error) {
	tags := []awsprovider.TagSpecification{
		{Name: fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), Value: "owned"},
	}
	forbiddenTags := sets.NewString()
	for idx := range tags {
		forbiddenTags.Insert(tags[idx].Name)
	}
	for k, v := range usertags {
		if forbiddenTags.Has(k) {
			return nil, fmt.Errorf("user tags may not clobber %s", k)
		}
		tags = append(tags, awsprovider.TagSpecification{Name: k, Value: v})
	}
	return tags, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
		providerSpec.LoadBalancers = []awsprovider.LoadBalancerReference{
			{
				Name: fmt.Sprintf("%s-ext", clusterID),
				Type: awsprovider.NetworkLoadBalancerType,
			},
			{
				Name: fmt.Sprintf("%s-int", clusterID),
				Type: awsprovider.NetworkLoadBalancerType,
			},
		}
	}
}
