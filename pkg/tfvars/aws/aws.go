package aws

// Endpoints is the type of the AWS endpoints.
type Endpoints string

const (
	// EndpointsAll represents the configuration for using both private and public endpoints.
	EndpointsAll Endpoints = "all"
	// EndpointsPrivate represents the configuration for using only private endpoints.
	EndpointsPrivate Endpoints = "private"
	// EndpointsPublic represents the configuration for using only public endpoints.
	EndpointsPublic Endpoints = "public"
	// DefaultVPCCIDRBlock is the default CIDR range for an AWS VPC.
	DefaultVPCCIDRBlock = "10.0.0.0/16"
	// DefaultRegion is the default AWS region for the cluster.
	DefaultRegion = "us-east-1"
)

// AWS converts AWS related config.
type AWS struct {
	EC2AMIOverride string    `json:"tectonic_aws_ec2_ami_override,omitempty"`
	Endpoints      Endpoints `json:"tectonic_aws_endpoints,omitempty"`
	External       `json:",inline"`
	ExtraTags      map[string]string `json:"tectonic_aws_extra_tags,omitempty"`
	InstallerRole  string            `json:"tectonic_aws_installer_role,omitempty"`
	Master         `json:",inline"`
	Region         string `json:"tectonic_aws_region,omitempty"`
	VPCCIDRBlock   string `json:"tectonic_aws_vpc_cidr_block,omitempty"`
	Worker         `json:",inline"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"tectonic_aws_external_master_subnet_ids,omitempty"`
	PrivateZone     string   `json:"tectonic_aws_external_private_zone,omitempty"`
	VPCID           string   `json:"tectonic_aws_external_vpc_id,omitempty"`
	WorkerSubnetIDs []string `json:"tectonic_aws_external_worker_subnet_ids,omitempty"`
}

// Master converts master related config.
type Master struct {
	CustomSubnets    map[string]string `json:"tectonic_aws_master_custom_subnets,omitempty"`
	EC2Type          string            `json:"tectonic_aws_master_ec2_type,omitempty"`
	ExtraSGIDs       []string          `json:"tectonic_aws_master_extra_sg_ids,omitempty"`
	IAMRoleName      string            `json:"tectonic_aws_master_iam_role_name,omitempty"`
	MasterRootVolume `json:",inline"`
}

// MasterRootVolume converts master rool volume related config.
type MasterRootVolume struct {
	IOPS int    `json:"tectonic_aws_master_root_volume_iops,omitempty"`
	Size int    `json:"tectonic_aws_master_root_volume_size,omitempty"`
	Type string `json:"tectonic_aws_master_root_volume_type,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	CustomSubnets map[string]string `json:"tectonic_aws_worker_custom_subnets,omitempty"`
	IAMRoleName   string            `json:"tectonic_aws_worker_iam_role_name,omitempty"`
}
