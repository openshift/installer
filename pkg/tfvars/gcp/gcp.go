package gcp

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
)

const (
	kmsKeyNameFmt = "projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	ProjectID        string `json:"gcp_project_id,omitempty"`
	NetworkProjectID string `json:"gcp_network_project_id,omitempty"`
	ServiceAccount   string `json:"gcp_service_account,omitempty"`
}

type config struct {
	Auth                      `json:",inline"`
	Region                    string   `json:"gcp_region,omitempty"`
	BootstrapInstanceType     string   `json:"gcp_bootstrap_instance_type,omitempty"`
	CreateBootstrapSA         bool     `json:"gcp_create_bootstrap_sa"`
	CreateFirewallRules       bool     `json:"gcp_create_firewall_rules"`
	MasterInstanceType        string   `json:"gcp_master_instance_type,omitempty"`
	MasterAvailabilityZones   []string `json:"gcp_master_availability_zones"`
	ImageURI                  string   `json:"gcp_image_uri,omitempty"`
	Image                     string   `json:"gcp_image,omitempty"`
	PreexistingImage          bool     `json:"gcp_preexisting_image"`
	InstanceServiceAccount    string   `json:"gcp_instance_service_account,omitempty"`
	ImageLicenses             []string `json:"gcp_image_licenses,omitempty"`
	VolumeType                string   `json:"gcp_master_root_volume_type"`
	VolumeSize                int64    `json:"gcp_master_root_volume_size"`
	VolumeKMSKeyLink          string   `json:"gcp_root_volume_kms_key_link"`
	PublicZoneName            string   `json:"gcp_public_zone_name,omitempty"`
	PrivateZoneName           string   `json:"gcp_private_zone_name,omitempty"`
	PublishStrategy           string   `json:"gcp_publish_strategy,omitempty"`
	PreexistingNetwork        bool     `json:"gcp_preexisting_network,omitempty"`
	ClusterNetwork            string   `json:"gcp_cluster_network,omitempty"`
	ControlPlaneSubnet        string   `json:"gcp_control_plane_subnet,omitempty"`
	ComputeSubnet             string   `json:"gcp_compute_subnet,omitempty"`
	ControlPlaneTags          []string `json:"gcp_control_plane_tags,omitempty"`
	SecureBoot                string   `json:"gcp_master_secure_boot,omitempty"`
	OnHostMaintenance         string   `json:"gcp_master_on_host_maintenance,omitempty"`
	EnableConfidentialCompute string   `json:"gcp_master_confidential_compute,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                Auth
	CreateFirewallRules bool
	ImageURI            string
	ImageLicenses       []string
	MasterConfigs       []*machineapi.GCPMachineProviderSpec
	WorkerConfigs       []*machineapi.GCPMachineProviderSpec
	PublicZoneName      string
	PrivateZoneName     string
	PublishStrategy     types.PublishingStrategy
	PreexistingNetwork  bool
}

// TFVars generates gcp-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	workerConfig := sources.WorkerConfigs[0]
	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Zone
	}

	cfg := &config{
		Auth:                      sources.Auth,
		Region:                    masterConfig.Region,
		BootstrapInstanceType:     masterConfig.MachineType,
		CreateFirewallRules:       sources.CreateFirewallRules,
		MasterInstanceType:        masterConfig.MachineType,
		MasterAvailabilityZones:   masterAvailabilityZones,
		VolumeType:                masterConfig.Disks[0].Type,
		VolumeSize:                masterConfig.Disks[0].SizeGB,
		ImageURI:                  sources.ImageURI,
		Image:                     masterConfig.Disks[0].Image,
		ImageLicenses:             sources.ImageLicenses,
		PublicZoneName:            sources.PublicZoneName,
		PrivateZoneName:           sources.PrivateZoneName,
		PublishStrategy:           string(sources.PublishStrategy),
		ClusterNetwork:            masterConfig.NetworkInterfaces[0].Network,
		ControlPlaneSubnet:        masterConfig.NetworkInterfaces[0].Subnetwork,
		ComputeSubnet:             workerConfig.NetworkInterfaces[0].Subnetwork,
		PreexistingNetwork:        sources.PreexistingNetwork,
		ControlPlaneTags:          masterConfig.Tags,
		SecureBoot:                string(masterConfig.ShieldedInstanceConfig.SecureBoot),
		EnableConfidentialCompute: string(masterConfig.ConfidentialCompute),
		OnHostMaintenance:         string(masterConfig.OnHostMaintenance),
	}

	cfg.PreexistingImage = true
	if len(sources.ImageLicenses) > 0 {
		cfg.PreexistingImage = false
	}

	if masterConfig.Disks[0].EncryptionKey != nil {
		cfg.VolumeKMSKeyLink = generateDiskEncryptionKeyLink(masterConfig.Disks[0].EncryptionKey, masterConfig.ProjectID)
	}

	serviceAccount := make(map[string]interface{})

	if err := json.Unmarshal([]byte(cfg.Auth.ServiceAccount), &serviceAccount); len(cfg.Auth.ServiceAccount) > 0 && err != nil {
		return nil, errors.Wrapf(err, "unmarshaling service account")
	}

	instanceServiceAccount := ""
	// Passthrough service accounts are only needed for GCP XPN.
	if len(cfg.Auth.NetworkProjectID) > 0 {
		var found bool
		instanceServiceAccount, found = serviceAccount["client_email"].(string)
		if !found {
			return nil, errors.New("could not find google service account")
		}
	}
	cfg.InstanceServiceAccount = instanceServiceAccount

	// A private key is needed to sign the URL for bootstrap ignition.
	// If there is no key in the credentials, we need to generate a new SA.
	_, foundKey := serviceAccount["private_key"]
	cfg.CreateBootstrapSA = !foundKey

	return json.MarshalIndent(cfg, "", "  ")
}

func generateDiskEncryptionKeyLink(keyRef *machineapi.GCPEncryptionKeyReference, projectID string) string {
	if keyRef.KMSKey.ProjectID != "" {
		projectID = keyRef.KMSKey.ProjectID
	}

	return fmt.Sprintf(kmsKeyNameFmt, projectID, keyRef.KMSKey.Location, keyRef.KMSKey.KeyRing, keyRef.KMSKey.Name)
}
