// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	netext "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/utils/openstack/clientconfig"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

const (
	// TODO(flaper87): We're choosing to hardcode these values to make
	// the environment more predictable. We expect there to a secret
	// named `openstack-credentials` and a cloud named `openstack` in
	// the clouds file stored in this secret.
	cloudsSecret          = "openstack-cloud-credentials"
	cloudsSecretNamespace = "openshift-machine-api"

	// CloudName is a constant containing the name of the cloud used in the internal cloudsSecret
	CloudName = "openstack"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, *machinev1.ControlPlaneMachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}

	mpool := pool.Platform.OpenStack
	platform := config.Platform.OpenStack
	trunkSupport, err := checkNetworkExtensionAvailability(platform.Cloud, "trunk", nil)
	if err != nil {
		return nil, nil, err
	}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	machines := make([]machineapi.Machine, 0, total)
	failureDomains := failureDomainsFromSpec(*mpool)
	for idx := int64(0); idx < total; idx++ {
		failureDomain := failureDomains[uint(idx)%uint(len(failureDomains))]

		providerSpec := generateProviderSpec(
			clusterID,
			platform,
			mpool,
			osImage,
			role,
			userDataSecret,
			trunkSupport,
			failureDomain,
		)

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
					Value: &runtime.RawExtension{Object: providerSpec},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}
		machines = append(machines, machine)
	}

	machineSetProviderSpec := generateProviderSpec(
		clusterID,
		platform,
		mpool,
		osImage,
		role,
		userDataSecret,
		trunkSupport,
		machinev1.OpenStackFailureDomain{RootVolume: &machinev1.RootVolume{}},
	)

	replicas := int32(total)

	controlPlaneMachineSet := &machinev1.ControlPlaneMachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1",
			Kind:       "ControlPlaneMachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      "cluster",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1.ControlPlaneMachineSetSpec{
			State:    machinev1.ControlPlaneMachineSetStateActive,
			Replicas: &replicas,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Template: machinev1.ControlPlaneMachineSetTemplate{
				MachineType: machinev1.OpenShiftMachineV1Beta1MachineType,
				OpenShiftMachineV1Beta1Machine: &machinev1.OpenShiftMachineV1Beta1MachineTemplate{
					ObjectMeta: machinev1.ControlPlaneMachineSetTemplateObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: machineapi.MachineSpec{
						ProviderSpec: machineapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: machineSetProviderSpec},
						},
					},
				},
			},
		},
	}

	if CPMSFailureDomains := pruneFailureDomains(failureDomains); CPMSFailureDomains != nil {
		controlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains = machinev1.FailureDomains{
			Platform:  v1.OpenStackPlatformType,
			OpenStack: CPMSFailureDomains,
		}
	}
	return machines, controlPlaneMachineSet, nil
}

func generateProviderSpec(clusterID string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, role, userDataSecret string, trunkSupport bool, failureDomain machinev1.OpenStackFailureDomain) *machinev1alpha1.OpenstackProviderSpec {
	var controlPlaneNetwork machinev1alpha1.NetworkParam
	additionalNetworks := make([]machinev1alpha1.NetworkParam, 0, len(mpool.AdditionalNetworkIDs))
	primarySubnet := ""

	if platform.ControlPlanePort != nil {
		var subnets []machinev1alpha1.SubnetParam
		controlPlanePort := platform.ControlPlanePort

		for _, fixedIP := range controlPlanePort.FixedIPs {
			subnets = append(subnets, machinev1alpha1.SubnetParam{
				Filter: machinev1alpha1.SubnetFilter{ID: fixedIP.Subnet.ID, Name: fixedIP.Subnet.Name},
			})
		}
		controlPlaneNetwork = machinev1alpha1.NetworkParam{
			Subnets: subnets,
			Filter: machinev1alpha1.Filter{
				Name: controlPlanePort.Network.Name,
				ID:   controlPlanePort.Network.ID,
			},
		}
		primarySubnet = controlPlanePort.FixedIPs[0].Subnet.ID
	} else {
		controlPlaneNetwork = machinev1alpha1.NetworkParam{
			Subnets: []machinev1alpha1.SubnetParam{
				{
					Filter: machinev1alpha1.SubnetFilter{
						Name: fmt.Sprintf("%s-nodes", clusterID),
						Tags: fmt.Sprintf("openshiftClusterID=%s", clusterID),
					},
				},
			},
		}
	}

	for _, networkID := range mpool.AdditionalNetworkIDs {
		additionalNetworks = append(additionalNetworks, machinev1alpha1.NetworkParam{
			UUID:                  networkID,
			NoAllowedAddressPairs: true,
		})
	}

	securityGroups := []machinev1alpha1.SecurityGroupParam{
		{
			Name: fmt.Sprintf("%s-%s", clusterID, role),
		},
	}
	for _, sg := range mpool.AdditionalSecurityGroupIDs {
		securityGroups = append(securityGroups, machinev1alpha1.SecurityGroupParam{
			UUID: sg,
		})
	}

	serverGroupName := clusterID + "-" + role
	// We initially used the AZ name as part of the server group name for the masters
	// but we realized that it was not useful. Whether or not the AZ is specified, the
	// masters will be spread across multiple hosts by default by the Nova scheduler
	// (the policy can be changed via `serverGroupPolicy` in install-config.yaml).
	// For the workers, we still use the AZ name as part of the server group name
	// so the user can control the scheduling policy per AZ and change the MachineSets
	// if needed on a day 2 operation.
	if role == "worker" && failureDomain.AvailabilityZone != "" {
		serverGroupName += "-" + failureDomain.AvailabilityZone
	}

	spec := machinev1alpha1.OpenstackProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machinev1alpha1.GroupVersion.String(),
			Kind:       "OpenstackProviderSpec",
		},
		Flavor:           mpool.FlavorName,
		CloudName:        CloudName,
		CloudsSecret:     &corev1.SecretReference{Name: cloudsSecret, Namespace: cloudsSecretNamespace},
		UserDataSecret:   &corev1.SecretReference{Name: userDataSecret},
		Networks:         append([]machinev1alpha1.NetworkParam{controlPlaneNetwork}, additionalNetworks...),
		PrimarySubnet:    primarySubnet,
		AvailabilityZone: failureDomain.AvailabilityZone,
		SecurityGroups:   securityGroups,
		ServerGroupName:  serverGroupName,
		Trunk:            trunkSupport,
		Tags: []string{
			fmt.Sprintf("openshiftClusterID=%s", clusterID),
		},
		ServerMetadata: map[string]string{
			"Name":               fmt.Sprintf("%s-%s", clusterID, role),
			"openshiftClusterID": clusterID,
		},
	}
	if mpool.RootVolume != nil {
		spec.RootVolume = &machinev1alpha1.RootVolume{
			Size:       mpool.RootVolume.Size,
			SourceUUID: osImage,
			VolumeType: failureDomain.RootVolume.VolumeType,
			Zone:       failureDomain.RootVolume.AvailabilityZone,
		}
	} else {
		spec.Image = osImage
	}
	return &spec
}

// failureDomainIsEmpty returns true if the failure domain only contains nil or
// zero values.
func failureDomainIsEmpty(failureDomain machinev1.OpenStackFailureDomain) bool {
	if failureDomain.AvailabilityZone == "" {
		if failureDomain.RootVolume == nil {
			return true
		}
		if failureDomain.RootVolume.AvailabilityZone == "" && failureDomain.RootVolume.VolumeType == "" {
			return true
		}
	}
	return false
}

// pruneFailureDomains returns nil if the only failure domain in the given
// slice is empty. One empty failure domain is not syntactically valid in CPMS.
func pruneFailureDomains(failureDomains []machinev1.OpenStackFailureDomain) []machinev1.OpenStackFailureDomain {
	if len(failureDomains) == 1 && failureDomainIsEmpty(failureDomains[0]) {
		return nil
	}
	return failureDomains
}

// failureDomainsFromSpec returns as many failure domains as there are zones in
// the given machine-pool. The returned failure domains have nil RootVolume if
// and only if the given machine-pool has nil RootVolume. The returned failure
// domain slice is guaranteed to have at least one element.
func failureDomainsFromSpec(mpool openstack.MachinePool) []machinev1.OpenStackFailureDomain {
	var numberOfFailureDomains int
	if mpool.RootVolume != nil {
		// At this point, after validation, compute availability zones,
		// storage avalability zones and root volume types must all be
		// equal in number. However, we want to accept case where any
		// of them has zero or one value (which means: apply the same
		// value to all failure domains).
		var (
			highestCardinality      int
			highestCardinalityField string
		)
		for field, cardinality := range map[string]int{
			"compute availability zones": len(mpool.Zones),
			"storage availability zones": len(mpool.RootVolume.Zones),
			"root volume types":          len(mpool.RootVolume.Types),
		} {
			if cardinality > 1 {
				if highestCardinality > 1 && cardinality != highestCardinality {
					panic(highestCardinalityField + " and " + field + " should have equal length")
				}
				highestCardinality = cardinality
				highestCardinalityField = field
			}
		}
		numberOfFailureDomains = highestCardinality
	} else {
		numberOfFailureDomains = len(mpool.Zones)
	}

	// No failure domain is exactly like one failure domain with the default values.
	if numberOfFailureDomains < 1 {
		numberOfFailureDomains = 1
	}

	failureDomains := make([]machinev1.OpenStackFailureDomain, numberOfFailureDomains)

	for i := range failureDomains {
		switch len(mpool.Zones) {
		case 0:
			failureDomains[i].AvailabilityZone = openstackdefaults.DefaultComputeAZ()
		case 1:
			failureDomains[i].AvailabilityZone = mpool.Zones[0]
		default:
			failureDomains[i].AvailabilityZone = mpool.Zones[i]
		}

		if mpool.RootVolume != nil {
			switch len(mpool.RootVolume.Zones) {
			case 0:
				failureDomains[i].RootVolume = &machinev1.RootVolume{
					AvailabilityZone: openstackdefaults.DefaultRootVolumeAZ(),
				}
			case 1:
				failureDomains[i].RootVolume = &machinev1.RootVolume{
					AvailabilityZone: mpool.RootVolume.Zones[0],
				}
			default:
				failureDomains[i].RootVolume = &machinev1.RootVolume{
					AvailabilityZone: mpool.RootVolume.Zones[i],
				}
			}

			switch len(mpool.RootVolume.Types) {
			case 0:
				panic("Root volume types should have been validated to have at least one element")
			case 1:
				failureDomains[i].RootVolume.VolumeType = mpool.RootVolume.Types[0]
			default:
				failureDomains[i].RootVolume.VolumeType = mpool.RootVolume.Types[i]
			}
		}
	}
	return failureDomains
}

func checkNetworkExtensionAvailability(cloud, alias string, opts *clientconfig.ClientOpts) (bool, error) {
	if opts == nil {
		opts = openstackdefaults.DefaultClientOpts(cloud)
	}
	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return false, err
	}

	res := netext.Get(conn, alias)
	if res.Err != nil {
		if _, ok := res.Err.(gophercloud.ErrDefault404); ok {
			return false, nil
		}
		return false, res.Err
	}

	return true, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
	/*for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}*/
}
