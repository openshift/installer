// Package gcp generates Machine objects for gcp.
package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, *machinev1.ControlPlaneMachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != gcp.Name {
		return nil, nil, fmt.Errorf("non-GCP configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcp.Name {
		return nil, nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.GCP
	mpool := pool.Platform.GCP
	azs := mpool.Zones

	credentialsMode := config.CredentialsMode

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	machineSetProvider := &machineapi.GCPMachineProviderSpec{}
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, osImage, azIndex, role, userDataSecret, credentialsMode)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create provider")
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
		*machineSetProvider = *provider
		machines = append(machines, machine)
	}
	replicas := int32(total)
	failureDomains := []machinev1.GCPFailureDomain{}
	sort.Strings(mpool.Zones)
	for _, zone := range mpool.Zones {
		domain := machinev1.GCPFailureDomain{
			Zone: zone,
		}
		failureDomains = append(failureDomains, domain)
	}
	machineSetProvider.Zone = ""
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
					FailureDomains: machinev1.FailureDomains{
						Platform: v1.GCPPlatformType,
						GCP:      &failureDomains,
					},
					ObjectMeta: machinev1.ControlPlaneMachineSetTemplateObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: machineapi.MachineSpec{
						ProviderSpec: machineapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: machineSetProvider},
						},
					},
				},
			},
		},
	}

	return machines, controlPlaneMachineSet, nil
}

func provider(clusterID string, platform *gcp.Platform, mpool *gcp.MachinePool, osImage string, azIdx int, role, userDataSecret string, credentialsMode types.CredentialsMode) (*machineapi.GCPMachineProviderSpec, error) {
	az := mpool.Zones[azIdx]
	if len(platform.Licenses) > 0 {
		osImage = fmt.Sprintf("%s-rhcos-image", clusterID)
	}
	network, subnetwork, err := getNetworks(platform, clusterID, role)
	if err != nil {
		return nil, err
	}

	var encryptionKey *machineapi.GCPEncryptionKeyReference

	if mpool.OSDisk.EncryptionKey != nil {
		encryptionKey = &machineapi.GCPEncryptionKeyReference{
			KMSKey: &machineapi.GCPKMSKeyReference{
				Name:      mpool.OSDisk.EncryptionKey.KMSKey.Name,
				KeyRing:   mpool.OSDisk.EncryptionKey.KMSKey.KeyRing,
				ProjectID: mpool.OSDisk.EncryptionKey.KMSKey.ProjectID,
				Location:  mpool.OSDisk.EncryptionKey.KMSKey.Location,
			},
			KMSKeyServiceAccount: mpool.OSDisk.EncryptionKey.KMSKeyServiceAccount,
		}
	}

	instanceServiceAccount := fmt.Sprintf("%s-%s@%s.iam.gserviceaccount.com", clusterID, role[0:1], platform.ProjectID)
	// The installer will create a service account for compute nodes with the above naming convention.
	// The same service account will be used for control plane nodes during a vanilla installation. During a
	// xpn installation, the installer will attempt to use an existing service account either through the
	// credentials or through a user supplied value from the install-config.
	if role == "master" && len(platform.NetworkProjectID) > 0 {
		instanceServiceAccount = mpool.ServiceAccount

		if instanceServiceAccount == "" {
			sess, err := gcpconfig.GetSession(context.TODO())
			if err != nil {
				return nil, err
			}

			var found bool
			serviceAccount := make(map[string]interface{})
			err = json.Unmarshal(sess.Credentials.JSON, &serviceAccount)
			if err != nil {
				return nil, err
			}
			instanceServiceAccount, found = serviceAccount["client_email"].(string)
			if !found {
				return nil, errors.New("could not find google service account")
			}
		}
	}

	shieldedInstanceConfig := machineapi.GCPShieldedInstanceConfig{}
	if mpool.SecureBoot == string(machineapi.SecureBootPolicyEnabled) {
		shieldedInstanceConfig.SecureBoot = machineapi.SecureBootPolicyEnabled
	}
	return &machineapi.GCPMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "GCPMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "gcp-cloud-credentials"},
		Disks: []*machineapi.GCPDisk{{
			AutoDelete:    true,
			Boot:          true,
			SizeGB:        mpool.OSDisk.DiskSizeGB,
			Type:          mpool.OSDisk.DiskType,
			Image:         osImage,
			EncryptionKey: encryptionKey,
		}},
		NetworkInterfaces: []*machineapi.GCPNetworkInterface{{
			Network:    network,
			ProjectID:  platform.NetworkProjectID,
			Subnetwork: subnetwork,
		}},
		ServiceAccounts: []machineapi.GCPServiceAccount{{
			Email:  instanceServiceAccount,
			Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
		}},
		Tags:                   append(mpool.Tags, []string{fmt.Sprintf("%s-%s", clusterID, role)}...),
		MachineType:            mpool.InstanceType,
		Region:                 platform.Region,
		Zone:                   az,
		ProjectID:              platform.ProjectID,
		ShieldedInstanceConfig: shieldedInstanceConfig,
		ConfidentialCompute:    machineapi.ConfidentialComputePolicy(mpool.ConfidentialCompute),
		OnHostMaintenance:      machineapi.GCPHostMaintenanceType(mpool.OnHostMaintenance),
	}, nil
}

// ConfigMasters assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, controlPlane *machinev1.ControlPlaneMachineSet, clusterID string, publish types.PublishingStrategy) error {
	var targetPools []string
	if publish == types.ExternalPublishingStrategy {
		targetPools = append(targetPools, fmt.Sprintf("%s-api", clusterID))
	}

	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
		providerSpec.TargetPools = targetPools
	}

	providerSpec, ok := controlPlane.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
	if !ok {
		return errors.New("Unable to set target pools to control plane machine set")
	}
	providerSpec.TargetPools = targetPools
	return nil
}
func getNetworks(platform *gcp.Platform, clusterID, role string) (string, string, error) {
	if platform.Network == "" {
		return fmt.Sprintf("%s-network", clusterID), fmt.Sprintf("%s-%s-subnet", clusterID, role), nil
	}

	switch role {
	case "worker":
		return platform.Network, platform.ComputeSubnet, nil
	case "master":
		return platform.Network, platform.ControlPlaneSubnet, nil
	default:
		return "", "", fmt.Errorf("unrecognized machine role %s", role)
	}
}
