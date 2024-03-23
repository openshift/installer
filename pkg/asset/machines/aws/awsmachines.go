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
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// MachineInput defines the inputs needed to generate a machine asset.
type MachineInput struct {
	Role     string
	Pool     *types.MachinePool
	Subnets  map[string]string
	Tags     capa.Tags
	PublicIP bool
}

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, in *MachineInput) ([]*asset.RuntimeFile, error) {
	if poolPlatform := in.Pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := in.Pool.Platform.AWS

	total := int64(1)
	if in.Pool.Replicas != nil {
		total = *in.Pool.Replicas
	}

	var result []*asset.RuntimeFile

	for idx := int64(0); idx < total; idx++ {
		var subnet *capa.AWSResourceReference
		// By not setting subnets for the machine, we let CAPA choose one for us
		if len(in.Subnets) > 0 {
			zone := mpool.Zones[int(idx)%len(mpool.Zones)]
			subnetID, ok := in.Subnets[zone]
			if len(in.Subnets) > 0 && !ok {
				return nil, fmt.Errorf("no subnet for zone %s", zone)
			}
			subnet = &capa.AWSResourceReference{}
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
		}

		awsMachine := &capa.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, in.Pool.Name, idx),
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
				AdditionalTags:       in.Tags,
				RootVolume: &capa.Volume{
					Size:          int64(mpool.EC2RootVolume.Size),
					Type:          capa.VolumeType(mpool.EC2RootVolume.Type),
					IOPS:          int64(mpool.EC2RootVolume.IOPS),
					Encrypted:     ptr.To(true),
					EncryptionKey: mpool.KMSKeyARN,
				},
				InstanceMetadataOptions: &capa.InstanceMetadataOptions{
					HTTPTokens:   capa.HTTPTokensState(mpool.EC2Metadata.Authentication),
					HTTPEndpoint: capa.InstanceMetadataEndpointStateEnabled,
				},
			},
		}
		awsMachine.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSMachine"))

		if in.Role == "bootstrap" {
			awsMachine.Name = capiutils.GenerateBoostrapMachineName(clusterID)
			awsMachine.Labels["install.openshift.io/bootstrap"] = ""
		}

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
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, in.Role)),
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
