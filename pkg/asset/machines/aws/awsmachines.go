// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, region string, subnets map[string]string, pool *types.MachinePool, role string, userTags map[string]string) ([]*asset.RuntimeFile, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.AWS

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	tags, err := CapaTagsFromUserTags(clusterID, userTags)
	if err != nil {
		return nil, fmt.Errorf("failed to create machineapi.TagSpecifications from UserTags: %w", err)
	}

	var result []*asset.RuntimeFile

	for idx := int64(0); idx < total; idx++ {
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		subnetID, ok := subnets[zone]
		if len(subnets) > 0 && !ok {
			return nil, fmt.Errorf("no subnet for zone %s", zone)
		}
		subnet := &capa.AWSResourceReference{}
		if subnetID == "" {
			subnet.Filters = []capa.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-subnet-private-%s", clusterID, zone)},
				},
			}
		} else {
			subnet.ID = ptr.To(subnetID)
		}

		awsMachine := &capa.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capa.AWSMachineSpec{
				Ignition:             &capa.Ignition{Version: "3.2"},
				UncompressedUserData: ptr.To(true),
				InstanceType:         mpool.InstanceType,
				AMI:                  capa.AMIReference{ID: ptr.To(mpool.AMIID)},
				SSHKeyName:           ptr.To(""),
				IAMInstanceProfile:   fmt.Sprintf("%s-master-profile", clusterID),
				Subnet:               subnet,
				AdditionalTags:       tags,
				RootVolume: &capa.Volume{
					Size:          int64(mpool.EC2RootVolume.Size),
					Type:          capa.VolumeType(mpool.EC2RootVolume.Type),
					IOPS:          int64(mpool.EC2RootVolume.IOPS),
					Encrypted:     ptr.To(true),
					EncryptionKey: mpool.KMSKeyARN,
				},
			},
		}
		awsMachine.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSMachine"))

		// Handle additional security groups.
		for _, sg := range mpool.AdditionalSecurityGroupIDs {
			awsMachine.Spec.AdditionalSecurityGroups = append(
				awsMachine.Spec.AdditionalSecurityGroups,
				capa.AWSResourceReference{ID: ptr.To(sg)},
			)
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", awsMachine.Name)},
			Object: awsMachine,
		})

		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: awsMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capa.GroupVersion.String(),
					Kind:       "AWSMachine",
					Name:       awsMachine.Name,
				},
			},
		}
		machine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
			Object: machine,
		})
	}
	return result, nil
}

// CapaTagsFromUserTags converts a map of user tags to a map of capa.Tags.
func CapaTagsFromUserTags(clusterID string, usertags map[string]string) (capa.Tags, error) {
	tags := capa.Tags{}
	tags[fmt.Sprintf("kubernetes.io/cluster/%s", clusterID)] = "owned"

	forbiddenTags := sets.New[string]()
	for key := range tags {
		forbiddenTags.Insert(key)
	}

	userTagKeys := sets.New[string]()
	for key := range usertags {
		userTagKeys.Insert(key)
	}

	if clobberedTags := userTagKeys.Intersection(forbiddenTags); clobberedTags.Len() > 0 {
		return nil, fmt.Errorf("user tag keys %v are not allowed", sets.List(clobberedTags))
	}

	for _, k := range sets.List(userTagKeys) {
		tags[k] = usertags[k]
	}
	return tags, nil
}
