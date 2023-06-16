// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/api/machine/v1beta1"
	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, region string, subnets icaws.Subnets, pool *types.MachinePool, role, userDataSecret string, userTags map[string]string) ([]*machineapi.MachineSet, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.AWS
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machinesets []*machineapi.MachineSet
	for idx, az := range mpool.Zones {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}
		subnet, ok := subnets[az]
		if len(subnets) > 0 && !ok {
			return nil, errors.Errorf("no subnet for zone %s", az)
		}

		publicSubnet := subnet.Public
		instanceType := mpool.InstanceType
		nodeLabels := make(map[string]string, 3)
		nodeTaints := []corev1.Taint{}

		if pool.Name == types.MachinePoolEdgeRoleName {
			// edge pools typically do not receive the same workloads between
			// different zoneGroups, thus the installer will discover preferred
			// instance based on the installer's preferred instance lookup.
			if subnet.PreferredEdgeInstanceType != "" {
				instanceType = subnet.PreferredEdgeInstanceType
			}
			nodeLabels = map[string]string{
				"node-role.kubernetes.io/edge":    "",
				"machine.openshift.io/zone-type":  subnet.ZoneType,
				"machine.openshift.io/zone-group": subnet.ZoneGroupName,
			}
			nodeTaints = append(nodeTaints, corev1.Taint{
				Key:    "node-role.kubernetes.io/edge",
				Effect: "NoSchedule",
			})
		}

		provider, err := provider(&machineProviderInput{
			clusterID:        clusterID,
			region:           region,
			subnet:           subnet.ID,
			instanceType:     instanceType,
			osImage:          mpool.AMIID,
			zone:             az,
			role:             "worker",
			userDataSecret:   userDataSecret,
			root:             &mpool.EC2RootVolume,
			imds:             mpool.EC2Metadata,
			userTags:         userTags,
			publicSubnet:     publicSubnet,
			securityGroupIDs: pool.Platform.AWS.AdditionalSecurityGroupIDs,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		name := fmt.Sprintf("%s-%s-%s", clusterID, pool.Name, az)
		spec := machineapi.MachineSpec{
			ProviderSpec: machineapi.ProviderSpec{
				Value: &runtime.RawExtension{Object: provider},
			},
			ObjectMeta: machineapi.ObjectMeta{
				Labels: nodeLabels,
			},
			Taints: nodeTaints,
		}

		mset := &machineapi.MachineSet{
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
			Spec: machineapi.MachineSetSpec{
				Replicas: &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Template: machineapi.MachineTemplateSpec{
					ObjectMeta: machineapi.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: spec,
					// we don't need to set Versions, because we control those via cluster operators.
				},
			},
		}
		machinesets = append(machinesets, mset)
	}

	return machinesets, nil
}
