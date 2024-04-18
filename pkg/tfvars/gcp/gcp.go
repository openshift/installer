package gcp

import (
	"encoding/json"
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
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
	Region                    string            `json:"gcp_region,omitempty"`
	BootstrapInstanceType     string            `json:"gcp_bootstrap_instance_type,omitempty"`
	CreateFirewallRules       bool              `json:"gcp_create_firewall_rules"`
	MasterInstanceType        string            `json:"gcp_master_instance_type,omitempty"`
	MasterAvailabilityZones   []string          `json:"gcp_master_availability_zones"`
	Image                     string            `json:"gcp_image,omitempty"`
	InstanceServiceAccount    string            `json:"gcp_instance_service_account,omitempty"`
	VolumeType                string            `json:"gcp_master_root_volume_type"`
	VolumeSize                int64             `json:"gcp_master_root_volume_size"`
	VolumeKMSKeyLink          string            `json:"gcp_root_volume_kms_key_link"`
	PublicZoneName            string            `json:"gcp_public_zone_name,omitempty"`
	PrivateZoneName           string            `json:"gcp_private_zone_name,omitempty"`
	PublishStrategy           string            `json:"gcp_publish_strategy,omitempty"`
	PreexistingNetwork        bool              `json:"gcp_preexisting_network,omitempty"`
	ClusterNetwork            string            `json:"gcp_cluster_network,omitempty"`
	ControlPlaneSubnet        string            `json:"gcp_control_plane_subnet,omitempty"`
	ComputeSubnet             string            `json:"gcp_compute_subnet,omitempty"`
	ControlPlaneTags          []string          `json:"gcp_control_plane_tags,omitempty"`
	SecureBoot                string            `json:"gcp_master_secure_boot,omitempty"`
	OnHostMaintenance         string            `json:"gcp_master_on_host_maintenance,omitempty"`
	EnableConfidentialCompute string            `json:"gcp_master_confidential_compute,omitempty"`
	ExtraLabels               map[string]string `json:"gcp_extra_labels,omitempty"`
	UserProvisionedDNS        bool              `json:"gcp_user_provisioned_dns,omitempty"`
	ExtraTags                 map[string]string `json:"gcp_extra_tags,omitempty"`
	IgnitionShim              string            `json:"gcp_ignition_shim,omitempty"`
	PresignedURL              string            `json:"gcp_signed_url"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                Auth
	CreateFirewallRules bool
	MasterConfigs       []*machineapi.GCPMachineProviderSpec
	WorkerConfigs       []*machineapi.GCPMachineProviderSpec
	PublicZoneName      string
	PrivateZoneName     string
	PublishStrategy     types.PublishingStrategy
	PreexistingNetwork  bool
	InfrastructureName  string
	UserProvisionedDNS  bool
	UserTags            map[string]string
	IgnitionShim        string
	PresignedURL        string
}

// TFVars generates gcp-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	workerConfig := sources.WorkerConfigs[0]
	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Zone
	}

	labels := make(map[string]string, len(masterConfig.Labels)+1)
	// add OCP default label
	labels[fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, sources.InfrastructureName)] = "owned"
	for k, v := range masterConfig.Labels {
		labels[k] = v
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
		Image:                     masterConfig.Disks[0].Image,
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
		ExtraLabels:               labels,
		UserProvisionedDNS:        sources.UserProvisionedDNS,
		ExtraTags:                 sources.UserTags,
		IgnitionShim:              sources.IgnitionShim,
		PresignedURL:              sources.PresignedURL,
	}

	if masterConfig.Disks[0].EncryptionKey != nil {
		cfg.VolumeKMSKeyLink = generateDiskEncryptionKeyLink(masterConfig.Disks[0].EncryptionKey, masterConfig.ProjectID)
	}

	instanceServiceAccount := ""
	// Service Account for masters set for xpn installs
	if len(cfg.Auth.NetworkProjectID) > 0 {
		if len(masterConfig.ServiceAccounts) > 0 {
			instanceServiceAccount = masterConfig.ServiceAccounts[0].Email
		}
	}
	cfg.InstanceServiceAccount = instanceServiceAccount

	return json.MarshalIndent(cfg, "", "  ")
}

func generateDiskEncryptionKeyLink(keyRef *machineapi.GCPEncryptionKeyReference, projectID string) string {
	if keyRef.KMSKey.ProjectID != "" {
		projectID = keyRef.KMSKey.ProjectID
	}

	return fmt.Sprintf(kmsKeyNameFmt, projectID, keyRef.KMSKey.Location, keyRef.KMSKey.KeyRing, keyRef.KMSKey.Name)
}
