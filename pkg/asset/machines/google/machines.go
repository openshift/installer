// Package google generates Machine objects for google.
package google

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/google"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]clusterapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != google.Name {
		return nil, fmt.Errorf("non-GCP configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != google.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	// platform := config.Platform.GCP
	// mpool := pool.Platform.GCP
	// azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []clusterapi.Machine
	for idx := int64(0); idx < total; idx++ {
		// azIndex := int(idx) % len(azs)
		// provider, err := provider(clusterID, clustername, platform, mpool, osImage, azIndex, role, userDataSecret)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "failed to create provider")
		// }
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
				ProviderSpec: clusterapi.ProviderSpec{},
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

// func provider(clusterID, clusterName string, platform *google.Platform, mpool *google.MachinePool, osImage string, azIdx int, role, userDataSecret string) (*awsprovider.AWSMachineProviderConfig, error) {
// 	az := mpool.Zones[azIdx]
// 	amiID := osImage
// 	tags, err := tagsFromUserTags(clusterID, clusterName, platform.UserTags)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to create awsprovider.TagSpecifications from UserTags")
// 	}
// 	return &awsprovider.AWSMachineProviderConfig{
// 		TypeMeta: metav1.TypeMeta{
// 			APIVersion: "awsproviderconfig.k8s.io/v1alpha1",
// 			Kind:       "AWSMachineProviderConfig",
// 		},
// 		InstanceType: mpool.InstanceType,
// 		BlockDevices: []awsprovider.BlockDeviceMappingSpec{
// 			{
// 				EBS: &awsprovider.EBSBlockDeviceSpec{
// 					VolumeType: pointer.StringPtr(mpool.Type),
// 					VolumeSize: pointer.Int64Ptr(int64(mpool.Size)),
// 				},
// 			},
// 		},
// 		AMI:                awsprovider.AWSResourceReference{ID: &amiID},
// 		Tags:               tags,
// 		IAMInstanceProfile: &awsprovider.AWSResourceReference{ID: pointer.StringPtr(fmt.Sprintf("%s-%s-profile", clusterName, role))},
// 		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
// 		Subnet: awsprovider.AWSResourceReference{
// 			Filters: []awsprovider.Filter{{
// 				Name:   "tag:Name",
// 				Values: []string{fmt.Sprintf("%s-%s-%s", clusterName, role, az)},
// 			}},
// 		},
// 		Placement: awsprovider.Placement{Region: platform.Region, AvailabilityZone: az},
// 		SecurityGroups: []awsprovider.AWSResourceReference{{
// 			Filters: []awsprovider.Filter{{
// 				Name:   "tag:Name",
// 				Values: []string{fmt.Sprintf("%s_%s_sg", clusterName, role)},
// 			}},
// 		}},
// 	}, nil
// }

// func tagsFromUserTags(clusterID, clusterName string, usertags map[string]string) ([]awsprovider.TagSpecification, error) {
// 	tags := []awsprovider.TagSpecification{
// 		{Name: "openshiftClusterID", Value: clusterID},
// 		{Name: fmt.Sprintf("kubernetes.io/cluster/%s", clusterName), Value: "owned"},
// 	}
// 	forbiddenTags := sets.NewString()
// 	for idx := range tags {
// 		forbiddenTags.Insert(tags[idx].Name)
// 	}
// 	for k, v := range usertags {
// 		if forbiddenTags.Has(k) {
// 			return nil, fmt.Errorf("user tags may not clobber %s", k)
// 		}
// 		tags = append(tags, awsprovider.TagSpecification{Name: k, Value: v})
// 	}
// 	return tags, nil
// }

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []clusterapi.Machine, clusterName string) {
	// for _, machine := range machines {
	// 	providerSpec := machine.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
	// 	providerSpec.PublicIP = pointer.BoolPtr(true)
	// 	providerSpec.LoadBalancers = []awsprovider.LoadBalancerReference{
	// 		{
	// 			Name: fmt.Sprintf("%s-ext", clusterName),
	// 			Type: awsprovider.NetworkLoadBalancerType,
	// 		},
	// 		{
	// 			Name: fmt.Sprintf("%s-int", clusterName),
	// 			Type: awsprovider.NetworkLoadBalancerType,
	// 		},
	// 	}
	// }
}
