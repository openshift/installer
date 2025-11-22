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

// MachineSetInput holds the input arguments required to MachineSets for a machinepool.
type MachineSetInput struct {
	ClusterID                string
	InstallConfigPlatformAWS *aws.Platform
	PublicSubnet             bool
	Subnets                  icaws.SubnetsByZone
	Zones                    icaws.Zones
	Pool                     *types.MachinePool
	Role                     string
	UserDataSecret           string
	Hosts                    map[string]icaws.Host
}

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(in *MachineSetInput) ([]*machineapi.MachineSet, error) {
	if poolPlatform := in.Pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := in.Pool.Platform.AWS
	azs := mpool.Zones

	total := int64(0)
	if in.Pool.Replicas != nil {
		total = *in.Pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machinesets []*machineapi.MachineSet
	for idx, az := range mpool.Zones {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		nodeLabels := make(map[string]string, 3)
		nodeTaints := []corev1.Taint{}
		instanceType := mpool.InstanceType
		publicSubnet := in.PublicSubnet
		subnetID := ""
		if len(in.Subnets) > 0 {
			subnet, ok := in.Subnets[az]
			if !ok {
				return nil, errors.Errorf("no subnet for zone %s", az)
			}
			publicSubnet = subnet.Public
			subnetID = subnet.ID
		}

		if in.Pool.Name == types.MachinePoolEdgeRoleName {
			// edge pools not share same instance type and regular cluster workloads.
			// The instance type is selected based in the offerings for the location.
			// The labels and taints are set to prevent regular workloads.
			// https://github.com/openshift/enhancements/blob/master/enhancements/installer/aws-custom-edge-machineset-local-zones.md
			zone := in.Zones[az]
			if zone.PreferredInstanceType != "" {
				instanceType = zone.PreferredInstanceType
			}
			nodeLabels = map[string]string{
				"node-role.kubernetes.io/edge":          "",
				"machine.openshift.io/zone-type":        zone.Type,
				"machine.openshift.io/zone-group":       zone.GroupName,
				"machine.openshift.io/parent-zone-name": zone.ParentZoneName,
			}
			nodeTaints = append(nodeTaints, corev1.Taint{
				Key:    "node-role.kubernetes.io/edge",
				Effect: "NoSchedule",
			})
		}

		instanceProfile := mpool.IAMProfile
		if len(instanceProfile) == 0 {
			instanceProfile = fmt.Sprintf("%s-worker-profile", in.ClusterID)
		}

		dedicatedHost := DedicatedHost(in.Hosts, mpool.HostPlacement, az)

		provider, err := provider(&machineProviderInput{
			clusterID:        in.ClusterID,
			region:           in.InstallConfigPlatformAWS.Region,
			subnet:           subnetID,
			instanceType:     instanceType,
			osImage:          mpool.AMIID,
			zone:             az,
			role:             "worker",
			userDataSecret:   in.UserDataSecret,
			instanceProfile:  instanceProfile,
			root:             &mpool.EC2RootVolume,
			imds:             mpool.EC2Metadata,
			userTags:         in.InstallConfigPlatformAWS.UserTags,
			publicSubnet:     publicSubnet,
			securityGroupIDs: in.Pool.Platform.AWS.AdditionalSecurityGroupIDs,
			cpuOptions:       mpool.CPUOptions,
			dedicatedHost:    dedicatedHost,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}

		// If we are using any feature that is only available via CAPI, we must set the authoritativeAPI = ClusterAPI
		authoritativeAPI := machineapi.MachineAuthorityMachineAPI
		if isAuthoritativeClusterAPIRequired(provider) {
			authoritativeAPI = machineapi.MachineAuthorityClusterAPI
		}

		name := fmt.Sprintf("%s-%s-%s", in.ClusterID, in.Pool.Name, az)
		spec := machineapi.MachineSpec{
			AuthoritativeAPI: authoritativeAPI,
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
					"machine.openshift.io/cluster-api-cluster": in.ClusterID,
				},
			},
			Spec: machineapi.MachineSetSpec{
				AuthoritativeAPI: authoritativeAPI,
				Replicas:         &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    in.ClusterID,
					},
				},
				Template: machineapi.MachineTemplateSpec{
					ObjectMeta: machineapi.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      in.ClusterID,
							"machine.openshift.io/cluster-api-machine-role": in.Role,
							"machine.openshift.io/cluster-api-machine-type": in.Role,
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

// isAuthoritativeClusterAPIRequired is called to determine if the machine spec should have the AuthoritativeAPI set to ClusterAPI.
func isAuthoritativeClusterAPIRequired(provider *machineapi.AWSMachineProviderConfig) bool {
	if provider.HostPlacement != nil && *provider.HostPlacement.Affinity != machineapi.HostAffinityAnyAvailable {
		return true
	}
	return false
}
