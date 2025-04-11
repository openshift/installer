// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"

	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpmanifests "github.com/openshift/installer/pkg/asset/manifests/gcp"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	masterRole = "master"

	kmsKeyNameFmt = "projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s"
)

func generateDiskEncryptionKeyLink(kmsKey *gcptypes.KMSKeyReference, projectID string) string {
	if kmsKey.ProjectID != "" {
		projectID = kmsKey.ProjectID
	}

	return fmt.Sprintf(kmsKeyNameFmt, projectID, kmsKey.Location, kmsKey.KeyRing, kmsKey.Name)
}

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
		gcpMachine, err := createGCPMachine(name, installConfig, infraID, mpool, imageName)
		if err != nil {
			return nil, fmt.Errorf("failed to create control plane (%d): %w", idx, err)
		}

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", gcpMachine.Name)},
			Object: gcpMachine,
		})

		dataSecret := fmt.Sprintf("%s-%s", infraID, masterRole)
		capiMachine := createCAPIMachine(gcpMachine.Name, dataSecret, infraID)

		if len(mpool.Zones) > 0 {
			// When there are fewer zones than the number of control plane instances,
			// cycle through the zones where the instances will reside.
			zone := mpool.Zones[int(idx)%len(mpool.Zones)]
			capiMachine.Spec.FailureDomain = ptr.To(zone)
		}

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
	bootstrapGCPMachine, err := createGCPMachine(name, installConfig, infraID, mpool, imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to create bootstrap machine: %w", err)
	}

	// Identify this as a bootstrap machine
	bootstrapGCPMachine.Labels["install.openshift.io/bootstrap"] = ""

	bootstrapMachineIsPublic := installConfig.Config.Publish == types.ExternalPublishingStrategy
	bootstrapGCPMachine.Spec.PublicIP = ptr.To(bootstrapMachineIsPublic)

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
func createGCPMachine(name string, installConfig *installconfig.InstallConfig, infraID string, mpool *gcptypes.MachinePool, imageName string) (*capg.GCPMachine, error) {
	// Use the rhcosImage as image name if not defined
	osImage := imageName
	if mpool.OSImage != nil {
		osImage = fmt.Sprintf("projects/%s/global/images/%s", mpool.OSImage.Project, mpool.OSImage.Name)
		logrus.Debugf("overriding gcp machine image: %s", osImage)
	}

	masterSubnet := installConfig.Config.Platform.GCP.ControlPlaneSubnet
	if masterSubnet == "" {
		masterSubnet = gcptypes.DefaultSubnetName(infraID, masterRole)
	}

	gcpMachine := &capg.GCPMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capg.GCPMachineSpec{
			InstanceType:          mpool.InstanceType,
			Subnet:                ptr.To(masterSubnet),
			AdditionalLabels:      getLabelsFromInstallConfig(installConfig, infraID),
			Image:                 ptr.To(osImage),
			RootDeviceType:        ptr.To(capg.DiskType(mpool.OSDisk.DiskType)),
			RootDeviceSize:        mpool.OSDisk.DiskSizeGB,
			AdditionalNetworkTags: mpool.Tags,
			ResourceManagerTags:   gcpmanifests.GetTagsFromInstallConfig(installConfig),
			IPForwarding:          ptr.To(capg.IPForwardingDisabled),
		},
	}
	gcpMachine.SetGroupVersionKind(capg.GroupVersion.WithKind("GCPMachine"))
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

	serviceAccountEmail := gcptypes.GetConfiguredServiceAccount(installConfig.Config.Platform.GCP, mpool)
	if serviceAccountEmail == "" {
		serviceAccountEmail = gcptypes.GetDefaultServiceAccount(installConfig.Config.Platform.GCP, infraID, masterRole[0:1])
	}
	serviceAccount := &capg.ServiceAccount{
		Email: serviceAccountEmail,
		// Set scopes to value defined at
		// https://cloud.google.com/compute/docs/access/service-accounts#scopes_best_practice
		Scopes: []string{compute.CloudPlatformScope},
	}

	gcpMachine.Spec.ServiceAccount = serviceAccount

	if mpool.OSDisk.EncryptionKey != nil {
		encryptionKey := &capg.CustomerEncryptionKey{
			KeyType: capg.CustomerManagedKey,
			ManagedKey: &capg.ManagedKey{
				KMSKeyName: generateDiskEncryptionKeyLink(mpool.OSDisk.EncryptionKey.KMSKey, installConfig.Config.GCP.ProjectID),
			},
		}
		if mpool.OSDisk.EncryptionKey.KMSKeyServiceAccount != "" {
			encryptionKey.KMSKeyServiceAccount = ptr.To(mpool.OSDisk.EncryptionKey.KMSKeyServiceAccount)
		}
		gcpMachine.Spec.RootDiskEncryptionKey = encryptionKey
	}

	return gcpMachine, nil
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
			Bootstrap: capi.Bootstrap{
				DataSecretName: ptr.To(dataSecret),
			},
			InfrastructureRef: v1.ObjectReference{
				APIVersion: capg.GroupVersion.String(),
				Kind:       "GCPMachine",
				Name:       name,
			},
		},
	}
	machine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

	return machine
}

func getLabelsFromInstallConfig(installConfig *installconfig.InstallConfig, infraID string) map[string]string {
	ic := installConfig.Config

	userLabels := map[string]string{}
	for _, label := range ic.Platform.GCP.UserLabels {
		userLabels[label.Key] = label.Value
	}
	// add OCP default label
	userLabels[fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, infraID)] = "owned"

	return userLabels
}
