// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// GenerateMachines returns manifests and runtime objects to provision control plane nodes using CAPI.
func GenerateMachines(installConfig *installconfig.InstallConfig, infraID string, pool *types.MachinePool, imageName string) ([]*asset.RuntimeFile, error) {
	var result []*asset.RuntimeFile
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcptypes.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.GCP

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	// Create GCP and CAPI machines for all master replicas in pool
	for idx := int64(0); idx < total; idx++ {
		name := fmt.Sprintf("%s-%s-%d", infraID, pool.Name, idx)
		gcpMachine := createGCPMachine(name, installConfig, infraID, mpool, imageName)

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", gcpMachine.Name)},
			Object: gcpMachine,
		})

		dataSecret := fmt.Sprintf("%s-master", infraID)
		capiMachine := createCAPIMachine(gcpMachine.Name, dataSecret, infraID)

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", capiMachine.Name)},
			Object: capiMachine,
		})
	}
	return result, nil
}

// GenerateBootstrapMachines returns a manifest and runtime object for a bootstrap node using CAPI.
func GenerateBootstrapMachines(name string, installConfig *installconfig.InstallConfig, infraID string, pool *types.MachinePool, imageName string) ([]*asset.RuntimeFile, error) {
	var result []*asset.RuntimeFile
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcptypes.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.GCP

	// Create one GCP and CAPI machine for bootstrap
	bootstrapGCPMachine := createGCPMachine(name, installConfig, infraID, mpool, imageName)

	// Identify this as a bootstrap machine
	bootstrapGCPMachine.Labels["install.openshift.io/bootstrap"] = ""

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", bootstrapGCPMachine.Name)},
		Object: bootstrapGCPMachine,
	})

	dataSecret := fmt.Sprintf("%s-%s", infraID, "bootstrap")
	bootstrapCapiMachine := createCAPIMachine(bootstrapGCPMachine.Name, dataSecret, infraID)

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", bootstrapCapiMachine.Name)},
		Object: bootstrapCapiMachine,
	})
	return result, nil
}

// Create a CAPG-specific machine.
func createGCPMachine(name string, installConfig *installconfig.InstallConfig, infraID string, mpool *gcptypes.MachinePool, imageName string) *capg.GCPMachine {
	// Use the rhcosImage as image name if not defined
	var osImage string
	if mpool.OSImage == nil {
		osImage = imageName
	} else {
		osImage = mpool.OSImage.Name
	}

	// TODO tags aren't currently being set in GCPMachine which only has
	// AdditionalNetworkTags []string

	masterSubnet := installConfig.Config.Platform.GCP.ControlPlaneSubnet
	if masterSubnet == "" {
		masterSubnet = gcptypes.DefaultSubnetName(infraID, "master")
	}

	gcpMachine := &capg.GCPMachine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "GCPMachine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capg.GCPMachineSpec{
			InstanceType:     mpool.InstanceType,
			Subnet:           ptr.To(masterSubnet),
			AdditionalLabels: getLabelsFromInstallConfig(installConfig, infraID),
			Image:            ptr.To(osImage),
			RootDeviceType:   ptr.To(capg.DiskType(mpool.OSDisk.DiskType)),
			RootDeviceSize:   mpool.OSDisk.DiskSizeGB,
		},
	}
	// Set optional values from machinepool
	if mpool.OnHostMaintenance != "" {
		gcpMachine.Spec.OnHostMaintenance = ptr.To(capg.HostMaintenancePolicy(mpool.OnHostMaintenance))
	}
	if mpool.ConfidentialCompute != "" {
		gcpMachine.Spec.ConfidentialCompute = ptr.To(capg.ConfidentialComputePolicy(mpool.ConfidentialCompute))
	}
	if mpool.SecureBoot != "" {
		shieldedInstanceConfig := capg.GCPShieldedInstanceConfig{}
		shieldedInstanceConfig.SecureBoot = capg.SecureBootPolicyEnabled
		gcpMachine.Spec.ShieldedInstanceConfig = ptr.To(shieldedInstanceConfig)
	}
	if mpool.ServiceAccount != "" {
		serviceAccount := &capg.ServiceAccount{
			Email: mpool.ServiceAccount,
			// Set scopes to value defined at
			// https://cloud.google.com/compute/docs/access/service-accounts#scopes_best_practice
			Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
		}
		gcpMachine.Spec.ServiceAccount = serviceAccount
	}

	return gcpMachine
}

// Create a CAPI machine based on the CAPG machine.
func createCAPIMachine(name string, dataSecret string, infraID string) *capi.Machine {
	machine := &capi.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capi.MachineSpec{
			ClusterName: infraID,
			// Leave empty until ignition support is added
			// Bootstrap: capi.Bootstrap{
			//	DataSecretName: ptr.To(dataSecret),
			// },
			InfrastructureRef: v1.ObjectReference{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
				Kind:       "GCPMachine",
				Name:       name,
			},
		},
	}

	return machine
}

func getLabelsFromInstallConfig(installConfig *installconfig.InstallConfig, infraID string) map[string]string {
	ic := installConfig.Config

	userLabels := map[string]string{}
	for _, label := range ic.Platform.GCP.UserLabels {
		userLabels[label.Key] = label.Value
	}
	// add OCP default label
	userLabels[fmt.Sprintf("kubernetes-io-cluster-%s", infraID)] = "owned"

	return userLabels
}
