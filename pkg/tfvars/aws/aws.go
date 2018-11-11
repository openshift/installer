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
)

// AWS converts AWS related config.
type AWS struct {
	EC2AMIOverride string    `json:"aws_ec2_ami_override,omitempty"`
	Endpoints      Endpoints `json:"aws_endpoints,omitempty"`
	External       `json:",inline"`
	ExtraTags      map[string]string `json:"aws_extra_tags,omitempty"`
	InstallerRole  string            `json:"aws_installer_role,omitempty"`
	Master         `json:",inline"`
	Region         string `json:"aws_region,omitempty"`
	VPCCIDRBlock   string `json:"aws_vpc_cidr_block"`
	Worker         `json:",inline"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"aws_external_master_subnet_ids,omitempty"`
	PrivateZone     string   `json:"aws_external_private_zone,omitempty"`
	WorkerSubnetIDs []string `json:"aws_external_worker_subnet_ids,omitempty"`
}

// Master converts master related config.
type Master struct {
	CustomSubnets    map[string]string `json:"aws_master_custom_subnets,omitempty"`
	EC2Type          string            `json:"aws_master_ec2_type,omitempty"`
	ExtraSGIDs       []string          `json:"aws_master_extra_sg_ids,omitempty"`
	IAMRoleName      string            `json:"aws_master_iam_role_name,omitempty"`
	MasterRootVolume `json:",inline"`
}

// MasterRootVolume converts master rool volume related config.
type MasterRootVolume struct {
	IOPS int    `json:"aws_master_root_volume_iops,omitempty"`
	Size int    `json:"aws_master_root_volume_size,omitempty"`
	Type string `json:"aws_master_root_volume_type,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	CustomSubnets map[string]string `json:"aws_worker_custom_subnets,omitempty"`
	IAMRoleName   string            `json:"aws_worker_iam_role_name,omitempty"`
}
