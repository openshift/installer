// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"
	"maps"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta2"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/utils"
)

// ClusterAPIMachineSets returns CAPI MachineSet and AWSMachineTemplate resources.
// This mirrors the MAPI MachineSets() function but produces CAPI-native types.
func ClusterAPIMachineSets(in *MachineSetInput) ([]capa.AWSMachineTemplate, []capi.MachineSet, error) {
	if poolPlatform := in.Pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	mpool := in.Pool.Platform.AWS
	azs := mpool.Zones

	total := int64(0)
	if in.Pool.Replicas != nil {
		total = *in.Pool.Replicas
	}
	numOfAZs := int64(len(azs))

	var templates []capa.AWSMachineTemplate
	var machineSets []capi.MachineSet

	imds := capa.HTTPTokensStateOptional
	if mpool.EC2Metadata.Authentication == "Required" {
		imds = capa.HTTPTokensStateRequired
	}

	instanceProfile := mpool.IAMProfile
	if len(instanceProfile) == 0 {
		instanceProfile = fmt.Sprintf("%s-worker-profile", in.ClusterID)
	}

	tags, err := CapaTagsFromUserTags(in.ClusterID, in.InstallConfigPlatformAWS.UserTags)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CAPA tags from user tags: %w", err)
	}

	for idx, az := range mpool.Zones {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		nodeLabels := map[string]string{
			"node-role.kubernetes.io/worker": "",
		}
		instanceType := mpool.InstanceType
		publicSubnet := in.PublicSubnet
		subnetRef := &capa.AWSResourceReference{}

		if len(in.Subnets) > 0 {
			subnet, ok := in.Subnets[az]
			if !ok {
				return nil, nil, fmt.Errorf("no subnet for zone %s", az)
			}
			publicSubnet = subnet.Public
			subnetRef.ID = ptr.To(subnet.ID)
		} else {
			subnetInternetScope := "private"
			if publicSubnet {
				subnetInternetScope = "public"
			}
			subnetRef.Filters = []capa.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-subnet-%s-%s", in.ClusterID, subnetInternetScope, az)},
				},
			}
		}

		// TODO: edge pools do not share same instance type and regular cluster workloads.
		// The instance type is selected based in the offerings for the location.
		// The labels and taints are set to prevent regular workloads.
		// https://github.com/openshift/enhancements/blob/master/enhancements/installer/aws-custom-edge-machineset-local-zones.md
		// FIXME: node taints on Machine/MachineSet is only supported in CAPI v1.12+ with feature gate MachineTaintPropagation.
		// Until we bump the CAPI version, edge machines can only be provisioned via MAPI.

		dedicatedHost := DedicatedHost(in.Hosts, mpool.HostPlacement, az)

		name := fmt.Sprintf("%s-%s-%s", in.ClusterID, in.Pool.Name, az)

		// Build AWSMachineTemplate for this zone
		template := capa.AWSMachineTemplate{
			TypeMeta: metav1.TypeMeta{
				APIVersion: capa.GroupVersion.String(),
				Kind:       "AWSMachineTemplate",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "openshift-cluster-api",
				Labels: map[string]string{
					"cluster.x-k8s.io/cluster-name": in.ClusterID,
				},
			},
			Spec: capa.AWSMachineTemplateSpec{
				Template: capa.AWSMachineTemplateResource{
					Spec: GenerateCAPIMachineSpec(&CAPIMachineSpecInput{
						InstanceType:       instanceType,
						AMI:                mpool.AMIID,
						IAMInstanceProfile: instanceProfile,
						Subnet:             subnetRef,
						PublicIP:           publicSubnet,
						Tags:               tags,
						EC2RootVolume:      mpool.EC2RootVolume,
						KMSKeyARN:          mpool.KMSKeyARN,
						IMDS:               imds,
						SecurityGroups: []capa.AWSResourceReference{
							{
								Filters: []capa.Filter{{
									Name:   "tag:Name",
									Values: []string{fmt.Sprintf("%s-node", in.ClusterID)},
								}},
							},
							{
								Filters: []capa.Filter{{
									Name:   "tag:Name",
									Values: []string{fmt.Sprintf("%s-lb", in.ClusterID)},
								}},
							},
						},
						AdditionalSecurityGroupIDs: mpool.AdditionalSecurityGroupIDs,
						CPUOptions:                 mpool.CPUOptions,
						Ignition: &capa.Ignition{
							Version: "3.2",
							// Worker machines should get ignition from the MCS on the control plane nodes
							StorageType: capa.IgnitionStorageTypeOptionUnencryptedUserData,
						},
						DedicatedHostID: dedicatedHost,
						IPFamily:        in.InstallConfigPlatformAWS.IPFamily,
					}),
				},
			},
		}
		templates = append(templates, template)

		// Build CAPI MachineSet referencing the template
		machineSet := capi.MachineSet{
			TypeMeta: metav1.TypeMeta{
				APIVersion: capi.GroupVersion.String(),
				Kind:       "MachineSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "openshift-cluster-api",
				Labels: map[string]string{
					"cluster.x-k8s.io/cluster-name": in.ClusterID,
				},
			},
			Spec: capi.MachineSetSpec{
				ClusterName: in.ClusterID,
				Replicas:    &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"cluster.x-k8s.io/cluster-name": in.ClusterID,
						"cluster.x-k8s.io/set-name":     name,
					},
				},
				Template: capi.MachineTemplateSpec{
					ObjectMeta: capi.ObjectMeta{
						Labels: map[string]string{
							"cluster.x-k8s.io/cluster-name": in.ClusterID,
							"cluster.x-k8s.io/set-name":     name,
						},
					},
					Spec: capi.MachineSpec{
						ClusterName: in.ClusterID,
						Bootstrap: capi.Bootstrap{
							DataSecretName: ptr.To(in.UserDataSecret),
						},
						InfrastructureRef: capi.ContractVersionedObjectReference{
							APIGroup: capa.GroupVersion.Group,
							Kind:     "AWSMachineTemplate",
							Name:     name,
						},
					},
				},
			},
		}
		// Machine labels will be synced from the Machine to the corresponding Node
		maps.Copy(machineSet.Spec.Template.ObjectMeta.Labels, nodeLabels)
		utils.SetCAPIMachineSetOSStreamLabels(&machineSet, in.Config)

		machineSets = append(machineSets, machineSet)
	}
	return templates, machineSets, nil
}
