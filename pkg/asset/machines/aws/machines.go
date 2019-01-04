// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// Machines returns a list of machines for a machinepool.
func Machines(config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]clusterapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.AWS
	mpool := pool.Platform.AWS
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []clusterapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(config.ClusterID, clustername, platform, mpool, azIndex, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := clusterapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.k8s.io/v1alpha1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-cluster-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clustername,
					"sigs.k8s.io/cluster-api-machine-role": role,
					"sigs.k8s.io/cluster-api-machine-type": role,
				},
			},
			Spec: clusterapi.MachineSpec{
				ProviderConfig: clusterapi.ProviderConfig{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID, clusterName string, platform *aws.Platform, mpool *aws.MachinePool, azIdx int, role, userDataSecret string) (*awsprovider.AWSMachineProviderConfig, error) {
	az := mpool.Zones[azIdx]
	tags, err := tagsFromUserTags(clusterID, clusterName, platform.UserTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create awsprovider.TagSpecifications from UserTags")
	}
	return &awsprovider.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "awsproviderconfig.k8s.io/v1alpha1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType:       mpool.InstanceType,
		AMI:                awsprovider.AWSResourceReference{ID: &mpool.AMIID},
		Tags:               tags,
		IAMInstanceProfile: &awsprovider.AWSResourceReference{ID: pointer.StringPtr(fmt.Sprintf("%s-%s-profile", clusterName, role))},
		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
		Subnet: awsprovider.AWSResourceReference{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s-%s-%s", clusterName, role, az)},
			}},
		},
		Placement: awsprovider.Placement{Region: platform.Region, AvailabilityZone: az},
		SecurityGroups: []awsprovider.AWSResourceReference{{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s_%s_sg", clusterName, role)},
			}},
		}},
	}, nil
}

func tagsFromUserTags(clusterID, clusterName string, usertags map[string]string) ([]awsprovider.TagSpecification, error) {
	tags := []awsprovider.TagSpecification{
		{Name: "openshiftClusterID", Value: clusterID},
		{Name: fmt.Sprintf("kubernetes.io/cluster/%s", clusterName), Value: "owned"},
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
func ConfigMasters(machines []clusterapi.Machine, clusterName string) {
	for _, machine := range machines {
		providerConfig := machine.Spec.ProviderConfig.Value.Object.(*awsprovider.AWSMachineProviderConfig)
		providerConfig.PublicIP = pointer.BoolPtr(true)
		providerConfig.LoadBalancers = []awsprovider.LoadBalancerReference{
			{
				Name: fmt.Sprintf("%s-ext", clusterName),
				Type: awsprovider.NetworkLoadBalancerType,
			},
			{
				Name: fmt.Sprintf("%s-int", clusterName),
				Type: awsprovider.NetworkLoadBalancerType,
			},
		}
	}
}
