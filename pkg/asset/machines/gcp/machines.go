// Package gcp generates Machine objects for gcp.
package gcp

import (
	"context"
	"fmt"

	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	googleoauth "golang.org/x/oauth2/google"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.GCP
	mpool := pool.Platform.GCP
	azs := mpool.Zones

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		azIndex := int(idx) % len(azs)
		provider, err := provider(clusterID, platform, mpool, osImage, azIndex, config.CredentialsMode, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
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

		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterID string, platform *gcp.Platform, mpool *gcp.MachinePool, osImage string, azIdx int, cm types.CredentialsMode, role, userDataSecret string) (*gcpprovider.GCPMachineProviderSpec, error) {
	az := mpool.Zones[azIdx]
	if len(platform.Licenses) > 0 {
		osImage = fmt.Sprintf("%s-rhcos-image", clusterID)
	}
	network, subnetwork, err := getNetworks(platform, clusterID, role)
	if err != nil {
		return nil, err
	}

	var encryptionKey *gcpprovider.GCPEncryptionKeyReference

	if mpool.OSDisk.EncryptionKey != nil {
		encryptionKey = &gcpprovider.GCPEncryptionKeyReference{
			KMSKey: &gcpprovider.GCPKMSKeyReference{
				Name:      mpool.OSDisk.EncryptionKey.KMSKey.Name,
				KeyRing:   mpool.OSDisk.EncryptionKey.KMSKey.KeyRing,
				ProjectID: mpool.OSDisk.EncryptionKey.KMSKey.ProjectID,
				Location:  mpool.OSDisk.EncryptionKey.KMSKey.Location,
			},
			KMSKeyServiceAccount: mpool.OSDisk.EncryptionKey.KMSKeyServiceAccount,
		}
	}

	sa, err := getServiceAccounts(cm, clusterID, role, platform.ProjectID)
	if err != nil {
		return nil, err
	}

	return &gcpprovider.GCPMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "gcpprovider.openshift.io/v1beta1",
			Kind:       "GCPMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "gcp-cloud-credentials"},
		Disks: []*gcpprovider.GCPDisk{{
			AutoDelete:    true,
			Boot:          true,
			SizeGb:        mpool.OSDisk.DiskSizeGB,
			Type:          mpool.OSDisk.DiskType,
			Image:         osImage,
			EncryptionKey: encryptionKey,
		}},
		NetworkInterfaces: []*gcpprovider.GCPNetworkInterface{{
			Network:    network,
			Subnetwork: subnetwork,
		}},
		ServiceAccounts: sa,
		Tags:            []string{fmt.Sprintf("%s-%s", clusterID, role)},
		MachineType:     mpool.InstanceType,
		Region:          platform.Region,
		Zone:            az,
		ProjectID:       platform.ProjectID,
	}, nil
}

// ConfigMasters assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string, publish types.PublishingStrategy) {
	var targetPools []string
	if publish == types.ExternalPublishingStrategy {
		targetPools = append(targetPools, fmt.Sprintf("%s-api", clusterID))
	}

	for _, machine := range machines {
		providerSpec := machine.Spec.ProviderSpec.Value.Object.(*gcpprovider.GCPMachineProviderSpec)
		providerSpec.TargetPools = targetPools
	}
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

// getServiceAccounts retrieves the appropriate service account value based on the credentials mode:
// Mint: generates a name to match the SA created in Terraform
// Passthrough: reads the email address from the session credentials
// Manual: leaves email blank, expecting users to manually edit manifests
func getServiceAccounts(cm types.CredentialsMode, clusterID, role, projectID string) ([]gcpprovider.GCPServiceAccount, error) {
	sa := gcpprovider.GCPServiceAccount{
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	}
	switch cm {
	case types.PassthroughCredentialsMode:
		var err error
		if sa.Email, err = getServiceAccountEmail(); err != nil {
			return nil, errors.Wrap(err, "failed to get service account email for passthrough mode")
		}
	case types.ManualCredentialsMode:
		// sa.Email is empty
	default:
		// types.MintCredentialsMode is the default behavior
		sa.Email = fmt.Sprintf("%s-%s@%s.iam.gserviceaccount.com", clusterID, role[0:1], projectID)
	}
	return []gcpprovider.GCPServiceAccount{sa}, nil
}

// getServiceAccountEmail retrieves the SA email address from the session credentials
func getServiceAccountEmail() (string, error) {
	ssn, err := ic.GetSession(context.TODO())
	if err != nil {
		return "", errors.Wrap(err, "failed to get session")
	}
	cfg, err := googleoauth.JWTConfigFromJSON(ssn.Credentials.JSON, "")
	if err != nil {
		return "", errors.Wrap(err, "failed to parse service account from credentials")
	}
	if cfg.Email == "" {
		return "", errors.New("CredentialsPassthroughMode requires GCP session credentials to have a valid client_email")
	}
	return cfg.Email, nil
}
