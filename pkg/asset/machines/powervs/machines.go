// Package powervs generates Machine objects for powerVS.
package powervs

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, *machinev1.ControlPlaneMachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != powervs.Name {
		return nil, nil, fmt.Errorf("non-PowerVS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, nil, fmt.Errorf("non-PowerVS machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.PowerVS
	mpool := pool.Platform.PowerVS

	// Only the service instance is guaranteed to exist and be passed via the install config
	// The other two, we should standardize a name including the cluster id.
	image := fmt.Sprintf("rhcos-%s", clusterID)
	var network string

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	machineProvider, err := provider(clusterID, platform, mpool, userDataSecret, image, network)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create provider: %w", err)
	}
	for idx := int64(0); idx < total; idx++ {
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
					Value: &runtime.RawExtension{Object: machineProvider},
				},
			},
		}
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
					Spec: machineapi.MachineSpec{
						ProviderSpec: machineapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: machineProvider},
						},
					},
				},
			},
		},
	}
	return machines, controlPlaneMachineSet, nil
}

func provider(clusterID string, platform *powervs.Platform, mpool *powervs.MachinePool, userDataSecret string, image string, network string) (*machinev1.PowerVSMachineProviderConfig, error) {

	if clusterID == "" || platform == nil || mpool == nil || userDataSecret == "" || image == "" {
		return nil, fmt.Errorf("invalid value passed to provider")
	}

	dhcpNetRegex := fmt.Sprintf("^DHCPSERVER.*%s.*_Private$", clusterID)

	var config *machinev1.PowerVSMachineProviderConfig

	// If a service instance GUID was not passed in the install-config.yaml file, then
	// we tell the machine provider to use a specific name via XXXTypeName.  Otherwide,
	// we tell the machine provider the given GUID via XXXTypeID.
	if platform.ServiceInstanceGUID == "" {
		serviceName := fmt.Sprintf("%s-power-iaas", clusterID)

		// Setting only the mandatory parameters
		config = &machinev1.PowerVSMachineProviderConfig{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PowerVSMachineProviderConfig",
				APIVersion: machinev1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{},
			ServiceInstance: machinev1.PowerVSResource{
				Type: machinev1.PowerVSResourceTypeName,
				Name: &serviceName,
			},
			Image: machinev1.PowerVSResource{
				Type: machinev1.PowerVSResourceTypeName,
				Name: &image,
			},
			UserDataSecret: &machinev1.PowerVSSecretReference{
				Name: userDataSecret,
			},
			CredentialsSecret: &machinev1.PowerVSSecretReference{
				Name: "powervs-credentials",
			},
			SystemType:    mpool.SysType,
			ProcessorType: mpool.ProcType,
			Processors:    mpool.Processors,
			MemoryGiB:     mpool.MemoryGiB,
			KeyPairName:   fmt.Sprintf("%s-key", clusterID),
		}
	} else {
		// Setting only the mandatory parameters
		config = &machinev1.PowerVSMachineProviderConfig{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PowerVSMachineProviderConfig",
				APIVersion: machinev1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{},
			ServiceInstance: machinev1.PowerVSResource{
				Type: machinev1.PowerVSResourceTypeID,
				ID:   &platform.ServiceInstanceGUID,
			},
			Image: machinev1.PowerVSResource{
				Type: machinev1.PowerVSResourceTypeName,
				Name: &image,
			},
			UserDataSecret: &machinev1.PowerVSSecretReference{
				Name: userDataSecret,
			},
			CredentialsSecret: &machinev1.PowerVSSecretReference{
				Name: "powervs-credentials",
			},
			SystemType:    mpool.SysType,
			ProcessorType: mpool.ProcType,
			Processors:    mpool.Processors,
			MemoryGiB:     mpool.MemoryGiB,
			KeyPairName:   fmt.Sprintf("%s-key", clusterID),
		}
	}
	if network != "" {
		config.Network = machinev1.PowerVSResource{
			Type: machinev1.PowerVSResourceTypeName,
			Name: &network,
		}
	} else {
		config.Network = machinev1.PowerVSResource{
			Type:  machinev1.PowerVSResourceTypeRegEx,
			RegEx: &dhcpNetRegex,
		}
	}
	return config, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines.
func ConfigMasters(machines []machineapi.Machine, controlPlane *machinev1.ControlPlaneMachineSet, infraID string, publish types.PublishingStrategy) error {
	lbrefs := []machinev1.LoadBalancerReference{{
		Name: fmt.Sprintf("%s-loadbalancer-int", infraID),
		Type: machinev1.ApplicationLoadBalancerType,
	}}

	if publish == types.ExternalPublishingStrategy {
		lbrefs = append(lbrefs, machinev1.LoadBalancerReference{
			Name: fmt.Sprintf("%s-loadbalancer", infraID),
			Type: machinev1.ApplicationLoadBalancerType,
		})
	}

	for _, machine := range machines {
		providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig)
		if !ok {
			return errors.New("unable to set load balancers to control plane machine set")
		}
		providerSpec.LoadBalancers = lbrefs
	}

	providerSpec, ok := controlPlane.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig)
	if !ok {
		return errors.New("unable to set load balancers to control plane machine set")
	}
	providerSpec.LoadBalancers = lbrefs
	return nil
}
