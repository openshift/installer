// Package nutanix generates Machine objects for nutanix.
package nutanix

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, *machinev1.ControlPlaneMachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != nutanix.Name {
		return nil, nil, fmt.Errorf("non nutanix configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != nutanix.Name {
		return nil, nil, fmt.Errorf("non-nutanix machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Nutanix
	mpool := pool.Platform.Nutanix

	failureDomains := make([]*nutanix.FailureDomain, 0, len(mpool.FailureDomains))
	for _, fdName := range mpool.FailureDomains {
		fd, err := platform.GetFailureDomainByName(fdName)
		if err != nil {
			return nil, nil, err
		}
		failureDomains = append(failureDomains, fd)
	}

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	var machineSetProvider *machinev1.NutanixMachineProviderConfig
	for idx := int64(0); idx < total; idx++ {
		var failureDomain *nutanix.FailureDomain
		if len(failureDomains) > 0 {
			failureDomain = failureDomains[idx%int64(len(failureDomains))]
		}
		provider, err := provider(clusterID, platform, mpool, osImage, userDataSecret, failureDomain)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create provider: %w", err)
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
		machineSetProvider = provider.DeepCopy()
		machines = append(machines, machine)
	}

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
			Replicas: &replicas,
			State:    machinev1.ControlPlaneMachineSetStateActive,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
					"machine.openshift.io/cluster-api-cluster":      clusterID,
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
				},
			},
		},
	}

	if len(failureDomains) > 0 {
		fdRefs := make([]machinev1.NutanixFailureDomainReference, 0, len(failureDomains))
		for _, fd := range failureDomains {
			fdRefs = append(fdRefs, machinev1.NutanixFailureDomainReference{Name: fd.Name})
		}
		controlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains = &machinev1.FailureDomains{
			Platform: configv1.NutanixPlatformType,
			Nutanix:  fdRefs,
		}

		// Reset the providerSpec fields related to the failure domain
		machineSetProvider.Cluster = machinev1.NutanixResourceIdentifier{}
		machineSetProvider.Subnets = []machinev1.NutanixResourceIdentifier{}
		machineSetProvider.FailureDomain = nil
	}

	controlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec = machineapi.MachineSpec{
		ProviderSpec: machineapi.ProviderSpec{
			Value: &runtime.RawExtension{Object: machineSetProvider},
		},
	}

	return machines, controlPlaneMachineSet, nil
}

func provider(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string, failureDomain *nutanix.FailureDomain) (*machinev1.NutanixMachineProviderConfig, error) {
	// cluster
	peUUID := platform.PrismElements[0].UUID
	if failureDomain != nil {
		peUUID = failureDomain.PrismElement.UUID
	}

	// subnets
	subnets := []machinev1.NutanixResourceIdentifier{}
	subnetUUIDs := platform.SubnetUUIDs
	if failureDomain != nil {
		subnetUUIDs = failureDomain.SubnetUUIDs
	}
	for _, subnetUUID := range subnetUUIDs {
		subnet := machinev1.NutanixResourceIdentifier{
			Type: machinev1.NutanixIdentifierUUID,
			UUID: &subnetUUID,
		}
		subnets = append(subnets, subnet)
	}

	providerCfg := &machinev1.NutanixMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machinev1.GroupVersion.String(),
			Kind:       "NutanixMachineProviderConfig",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: nutanix.CredentialsSecretName},
		Image: machinev1.NutanixResourceIdentifier{
			Type: machinev1.NutanixIdentifierName,
			Name: &osImage,
		},
		Subnets:        subnets,
		VCPUsPerSocket: int32(mpool.NumCoresPerSocket),
		VCPUSockets:    int32(mpool.NumCPUs),
		MemorySize:     resource.MustParse(fmt.Sprintf("%dMi", mpool.MemoryMiB)),
		Cluster: machinev1.NutanixResourceIdentifier{
			Type: machinev1.NutanixIdentifierUUID,
			UUID: &peUUID,
		},
		SystemDiskSize: resource.MustParse(fmt.Sprintf("%dGi", mpool.OSDisk.DiskSizeGiB)),
	}

	// FailureDomain
	if failureDomain != nil {
		providerCfg.FailureDomain = &machinev1.NutanixFailureDomainReference{Name: failureDomain.Name}
	}

	if len(mpool.BootType) != 0 {
		providerCfg.BootType = mpool.BootType
	}

	if mpool.Project != nil && mpool.Project.Type == machinev1.NutanixIdentifierUUID {
		providerCfg.Project = machinev1.NutanixResourceIdentifier{
			Type: machinev1.NutanixIdentifierUUID,
			UUID: mpool.Project.UUID,
		}
	}

	if len(mpool.Categories) > 0 {
		providerCfg.Categories = mpool.Categories
	}

	return providerCfg, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
}
