package gcp

import (
	"encoding/json"
	"fmt"

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
	CreateFirewallRules       bool     `json:"gcp_create_firewall_rules"`
	CreatePublicZoneRecords   bool     `json:"gcp_create_public_zone_records"`
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
	PrivateZoneProject        string   `json:"gcp_private_zone_project,omitempty"`
	PublicZoneName            string   `json:"gcp_public_zone_name,omitempty"`
	PublicZoneProject         string   `json:"gcp_public_zone_project,omitempty"`
	PublishStrategy           string   `json:"gcp_publish_strategy,omitempty"`
	PreexistingNetwork        bool     `json:"gcp_preexisting_network,omitempty"`
	ClusterNetwork            string   `json:"gcp_cluster_network,omitempty"`
	ControlPlaneSubnet        string   `json:"gcp_control_plane_subnet,omitempty"`
	ComputeSubnet             string   `json:"gcp_compute_subnet,omitempty"`
	ControlPlaneTags          []string `json:"gcp_control_plane_tags,omitempty"`
	OnHostMaintenance         string   `json:"gcp_master_on_host_maintenance,omitempty"`
	EnableConfidentialCompute bool     `json:"gcp_master_enable_confidential_compute,omitempty"`
	EnableSecureBoot          bool     `json:"gcp_master_enable_secure_boot,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                    Auth
	CreateFirewallRules     bool
	CreatePublicZoneRecords bool
	ImageURI                string
	ImageLicenses           []string
	InstanceServiceAccount  string
	MasterConfigs           []*machineapi.GCPMachineProviderSpec
	WorkerConfigs           []*machineapi.GCPMachineProviderSpec
	PrivateZoneProject      string
	PublicZoneName          string
	PublicZoneProject       string
	PublishStrategy         types.PublishingStrategy
	PreexistingNetwork      bool
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
		CreatePublicZoneRecords:   sources.CreatePublicZoneRecords,
		MasterInstanceType:        masterConfig.MachineType,
		MasterAvailabilityZones:   masterAvailabilityZones,
		VolumeType:                masterConfig.Disks[0].Type,
		VolumeSize:                masterConfig.Disks[0].SizeGB,
		ImageURI:                  sources.ImageURI,
		Image:                     masterConfig.Disks[0].Image,
		ImageLicenses:             sources.ImageLicenses,
		InstanceServiceAccount:    sources.InstanceServiceAccount,
		PrivateZoneProject:        sources.PrivateZoneProject,
		PublicZoneName:            sources.PublicZoneName,
		PublicZoneProject:         sources.PublicZoneProject,
		PublishStrategy:           string(sources.PublishStrategy),
		ClusterNetwork:            masterConfig.NetworkInterfaces[0].Network,
		ControlPlaneSubnet:        masterConfig.NetworkInterfaces[0].Subnetwork,
		ComputeSubnet:             workerConfig.NetworkInterfaces[0].Subnetwork,
		PreexistingNetwork:        sources.PreexistingNetwork,
		ControlPlaneTags:          masterConfig.Tags,
		EnableConfidentialCompute: masterConfig.EnableConfidentialCompute,
		OnHostMaintenance:         string(masterConfig.OnHostMaintenance),
		EnableSecureBoot:          masterConfig.EnableSecureBoot,
	}

	cfg.PreexistingImage = true
	if len(sources.ImageLicenses) > 0 {
		cfg.PreexistingImage = false
	}

	if masterConfig.Disks[0].EncryptionKey != nil {
		cfg.VolumeKMSKeyLink = generateDiskEncryptionKeyLink(masterConfig.Disks[0].EncryptionKey, masterConfig.ProjectID)
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func generateDiskEncryptionKeyLink(keyRef *machineapi.GCPEncryptionKeyReference, projectID string) string {
	if keyRef.KMSKey.ProjectID != "" {
		projectID = keyRef.KMSKey.ProjectID
	}

	return fmt.Sprintf(kmsKeyNameFmt, projectID, keyRef.KMSKey.Location, keyRef.KMSKey.KeyRing, keyRef.KMSKey.Name)
}
