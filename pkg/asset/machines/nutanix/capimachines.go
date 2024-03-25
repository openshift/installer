// Package aws generates capi Machine objects for nutanix.
package nutanix

import (
	"fmt"

	capnv1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	masterRole = "master"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage string, role string) ([]*asset.RuntimeFile, error) {
	machines, _, err := Machines(clusterID, config, pool, osImage, role, "")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve machines: %w", err)
	}

	ntxMachines := make([]*capnv1.NutanixMachine, 0, len(machines))
	result := make([]*asset.RuntimeFile, 0, 2*len(machines))

	for _, machine := range machines {
		providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig)
		if !ok {
			return nil, fmt.Errorf("unable to convert ProviderSpec to NutanixMachineProviderConfig")
		}

		// create the NutanixMachine object
		ntxMachine := generateNutanixMachine(machine.Name, providerSpec)
		ntxMachines = append(ntxMachines, ntxMachine)
		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", ntxMachine.Name)},
			Object: ntxMachine,
		})

		// create the capi Machine object
		capiMachine := &capv1.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.x-k8s.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiutils.Namespace,
				Name:      ntxMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capv1.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capv1.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:       "NutanixMachine",
					Name:       ntxMachine.Name,
				},
			},
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", capiMachine.Name)},
			Object: capiMachine,
		})
	}

	// as part of provisioning control plane nodes, we need to create a bootstrap node as well
	if role == masterRole {
		bootstrapSpec := ntxMachines[0].Spec.DeepCopy()
		bootstrapSpec.VCPUsPerSocket = 4
		bootstrapSpec.VCPUSockets = 1
		bootstrapImgName := nutanixtypes.BootISOImageName(clusterID)
		bootstrapSpec.Image.Name = &bootstrapImgName
		bootstrapNtxMachine := &capnv1.NutanixMachine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
				Kind:       "NutanixMachine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiutils.Namespace,
				Name:      capiutils.GenerateBoostrapMachineName(clusterID),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
					"install.openshift.io/bootstrap": "",
				},
			},
			Spec: *bootstrapSpec,
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", bootstrapNtxMachine.Name)},
			Object: bootstrapNtxMachine,
		})

		bootstrapCapiMachine := &capv1.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.x-k8s.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiutils.Namespace,
				Name:      bootstrapNtxMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capv1.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capv1.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-bootstrap", clusterID)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:       "NutanixMachine",
					Name:       bootstrapNtxMachine.Name,
				},
			},
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", bootstrapCapiMachine.Name)},
			Object: bootstrapCapiMachine,
		})
	}

	return result, nil
}

func generateNutanixMachine(machineName string, providerSpec *machinev1.NutanixMachineProviderConfig) *capnv1.NutanixMachine {
	ntxMachine := &capnv1.NutanixMachine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "NutanixMachine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: capiutils.Namespace,
			Name:      machineName,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capnv1.NutanixMachineSpec{
			VCPUsPerSocket: providerSpec.VCPUsPerSocket,
			VCPUSockets:    providerSpec.VCPUSockets,
			MemorySize:     providerSpec.MemorySize,
			SystemDiskSize: providerSpec.SystemDiskSize,
			Image: capnv1.NutanixResourceIdentifier{
				Type: capnv1.NutanixIdentifierType(providerSpec.Image.Type),
				Name: providerSpec.Image.Name,
				UUID: providerSpec.Image.UUID,
			},
			Cluster: capnv1.NutanixResourceIdentifier{
				Type: capnv1.NutanixIdentifierType(providerSpec.Cluster.Type),
				Name: providerSpec.Cluster.Name,
				UUID: providerSpec.Cluster.UUID,
			},
			Subnets:              []capnv1.NutanixResourceIdentifier{},
			AdditionalCategories: []capnv1.NutanixCategoryIdentifier{},
			BootType:             capnv1.NutanixBootType(providerSpec.BootType),
		},
	}

	for _, subnet := range providerSpec.Subnets {
		ntxMachine.Spec.Subnets = append(ntxMachine.Spec.Subnets, capnv1.NutanixResourceIdentifier{
			Type: capnv1.NutanixIdentifierType(subnet.Type),
			Name: subnet.Name,
			UUID: subnet.UUID,
		})
	}

	for _, category := range providerSpec.Categories {
		ntxMachine.Spec.AdditionalCategories = append(ntxMachine.Spec.AdditionalCategories,
			capnv1.NutanixCategoryIdentifier{
				Key:   category.Key,
				Value: category.Value,
			})
	}

	if providerSpec.BootType != "" {
		ntxMachine.Spec.BootType = capnv1.NutanixBootType(providerSpec.BootType)
	}

	if providerSpec.Project.Type != "" {
		ntxMachine.Spec.Project = &capnv1.NutanixResourceIdentifier{
			Type: capnv1.NutanixIdentifierType(providerSpec.Project.Type),
			Name: providerSpec.Project.Name,
			UUID: providerSpec.Project.UUID,
		}
	}

	return ntxMachine
}
