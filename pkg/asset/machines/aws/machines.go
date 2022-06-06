// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, region string, subnets map[string]string, pool *types.MachinePool, role, userDataSecret string, userTags map[string]string) ([]machineapi.Machine, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.AWS

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		subnet, ok := subnets[zone]
		if len(subnets) > 0 && !ok {
			return nil, errors.Errorf("no subnet for zone %s", zone)
		}
		provider, err := provider(
			clusterID,
			region,
			subnet,
			mpool.InstanceType,
			&mpool.EC2RootVolume,
			mpool.EC2Metadata,
			mpool.AMIID,
			zone,
			role,
			userDataSecret,
			userTags,
		)
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

func provider(clusterID string, region string, subnet string, instanceType string, root *aws.EC2RootVolume, imds aws.EC2Metadata, osImage string, zone, role, userDataSecret string, userTags map[string]string) (*machineapi.AWSMachineProviderConfig, error) {
	tags, err := tagsFromUserTags(clusterID, userTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create machineapi.TagSpecifications from UserTags")
	}
	config := &machineapi.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType: instanceType,
		BlockDevices: []machineapi.BlockDeviceMappingSpec{
			{
				EBS: &machineapi.EBSBlockDeviceSpec{
					VolumeType: pointer.StringPtr(root.Type),
					VolumeSize: pointer.Int64Ptr(int64(root.Size)),
					Iops:       pointer.Int64Ptr(int64(root.IOPS)),
					Encrypted:  pointer.BoolPtr(true),
					KMSKey:     machineapi.AWSResourceReference{ARN: pointer.StringPtr(root.KMSKeyARN)},
				},
			},
		},
		Tags:               tags,
		IAMInstanceProfile: &machineapi.AWSResourceReference{ID: pointer.StringPtr(fmt.Sprintf("%s-%s-profile", clusterID, role))},
		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret:  &corev1.LocalObjectReference{Name: "aws-cloud-credentials"},
		Placement:          machineapi.Placement{Region: region, AvailabilityZone: zone},
		SecurityGroups: []machineapi.AWSResourceReference{{
			Filters: []machineapi.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s-%s-sg", clusterID, role)},
			}},
		}},
	}

	if subnet == "" {
		config.Subnet.Filters = []machineapi.Filter{{
			Name:   "tag:Name",
			Values: []string{fmt.Sprintf("%s-private-%s", clusterID, zone)},
		}}
	} else {
		config.Subnet.ID = pointer.StringPtr(subnet)
	}

	if osImage == "" {
		config.AMI.Filters = []machineapi.Filter{{
			Name:   "tag:Name",
			Values: []string{fmt.Sprintf("%s-ami-%s", clusterID, region)},
		}}
	} else {
		config.AMI.ID = pointer.StringPtr(osImage)
	}

	if imds.Authentication != "" {
		config.MetadataServiceOptions.Authentication = machineapi.MetadataServiceAuthentication(imds.Authentication)
	}

	return config, nil
}

func tagsFromUserTags(clusterID string, usertags map[string]string) ([]machineapi.TagSpecification, error) {
	tags := []machineapi.TagSpecification{
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
		tags = append(tags, machineapi.TagSpecification{Name: k, Value: v})
	}
	return tags, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string, publish types.PublishingStrategy) {
	lbrefs := []machineapi.LoadBalancerReference{{
		Name: fmt.Sprintf("%s-int", clusterID),
		Type: machineapi.NetworkLoadBalancerType,
	}}

	if publish == types.ExternalPublishingStrategy {
		lbrefs = append(lbrefs, machineapi.LoadBalancerReference{
			Name: fmt.Sprintf("%s-ext", clusterID),
			Type: machineapi.NetworkLoadBalancerType,
		})
	}

	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*machineapi.AWSMachineProviderConfig)
		providerSpec.LoadBalancers = lbrefs
	}
}
