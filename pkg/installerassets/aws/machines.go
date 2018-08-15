package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func masterMachinesRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/aws/99_openshift-cluster-api_master-machines.yaml",
		RebuildHelper: masterMachinesRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"aws/ami",
		"aws/instance-type",
		"aws/region",
		"aws/user-tags",
		"aws/zones",
		"cluster-id",
		"cluster-name",
		"machines/master-count",
	)
	if err != nil {
		return nil, err
	}

	ami := string(parents["aws/ami"].Data)
	clusterID := string(parents["cluster-id"].Data)
	clusterName := string(parents["cluster-name"].Data)
	instanceType := string(parents["aws/instance-type"].Data)
	region := string(parents["aws/region"].Data)
	var userTags map[string]string
	err = yaml.Unmarshal(parents["aws/user-tags"].Data, &userTags)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal user tags")
	}

	masterCount, err := strconv.ParseUint(string(parents["machines/master-count"].Data), 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "parse master count")
	}

	var zones []string
	err = yaml.Unmarshal(parents["aws/zones"].Data, &zones)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal zones")
	}

	role := "master"
	userDataSecret := fmt.Sprintf("%s-user-data", role)
	poolName := role // FIXME: knob to control this?
	total := int64(masterCount)

	var machines []runtime.RawExtension
	for idx := int64(0); idx < total; idx++ {
		zone := zones[int(idx)%len(zones)]
		provider, err := provider(clusterID, clusterName, region, zone, instanceType, ami, userTags, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "create provider")
		}

		provider.PublicIP = pointer.BoolPtr(true)
		provider.LoadBalancers = []awsprovider.LoadBalancerReference{
			{
				Name: fmt.Sprintf("%s-ext", clusterName),
				Type: awsprovider.NetworkLoadBalancerType,
			},
			{
				Name: fmt.Sprintf("%s-int", clusterName),
				Type: awsprovider.NetworkLoadBalancerType,
			},
		}

		machine := clusterapi.Machine{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Machine",
				APIVersion: "cluster.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s-%d", clusterName, poolName, idx),
				Namespace: "openshift-cluster-api",
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clusterName,
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

		machines = append(machines, runtime.RawExtension{Object: &machine})
	}

	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		Items: machines,
	}

	asset.Data, err = yaml.Marshal(list)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func provider(clusterID, clusterName, region, zone, instanceType, ami string, userTags map[string]string, role, userDataSecret string) (*awsprovider.AWSMachineProviderConfig, error) {
	tags, err := tagsFromUserTags(clusterID, clusterName, userTags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create awsprovider.TagSpecifications from UserTags")
	}
	return &awsprovider.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "aws.cluster.k8s.io/v1alpha1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType:       instanceType,
		AMI:                awsprovider.AWSResourceReference{ID: &ami},
		Tags:               tags,
		IAMInstanceProfile: &awsprovider.AWSResourceReference{ID: pointer.StringPtr(fmt.Sprintf("%s-%s-profile", clusterName, role))},
		UserDataSecret:     &corev1.LocalObjectReference{Name: userDataSecret},
		Subnet: awsprovider.AWSResourceReference{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s-%s-%s", clusterName, role, zone)},
			}},
		},
		Placement: awsprovider.Placement{Region: region, AvailabilityZone: zone},
		SecurityGroups: []awsprovider.AWSResourceReference{{
			Filters: []awsprovider.Filter{{
				Name:   "tag:Name",
				Values: []string{fmt.Sprintf("%s_%s_sg", clusterName, role)},
			}},
		}},
	}, nil
}

func init() {
	installerassets.Rebuilders["manifests/aws/99_openshift-cluster-api_master-machines.yaml"] = masterMachinesRebuilder
	installerassets.Defaults["aws/instance-type"] = installerassets.ConstantDefault([]byte("m4.large"))
}
