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
	// DefaultProfile is the default AWS credentials profile to use.
	DefaultProfile = "default"
	// DefaultRegion is the default AWS region for the cluster.
	DefaultRegion = "us-east-1"
)

// AWS converts AWS related config.
type AWS struct {
	EC2AMIOverride string    `json:"tectonic_aws_ec2_ami_override,omitempty" yaml:"ec2AMIOverride,omitempty"`
	Endpoints      Endpoints `json:"tectonic_aws_endpoints,omitempty" yaml:"endpoints,omitempty"`
	External       `json:",inline" yaml:"external,omitempty"`
	ExtraTags      map[string]string `json:"tectonic_aws_extra_tags,omitempty" yaml:"extraTags,omitempty"`
	InstallerRole  string            `json:"tectonic_aws_installer_role,omitempty" yaml:"installerRole,omitempty"`
	Master         `json:",inline" yaml:"master,omitempty"`
	Profile        string `json:"tectonic_aws_profile,omitempty" yaml:"profile,omitempty"`
	Region         string `json:"tectonic_aws_region,omitempty" yaml:"region,omitempty"`
	VPCCIDRBlock   string `json:"tectonic_aws_vpc_cidr_block,omitempty" yaml:"vpcCIDRBlock,omitempty"`
	Worker         `json:",inline" yaml:"worker,omitempty"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"tectonic_aws_external_master_subnet_ids,omitempty" yaml:"masterSubnetIDs,omitempty"`
	PrivateZone     string   `json:"tectonic_aws_external_private_zone,omitempty" yaml:"privateZone,omitempty"`
	VPCID           string   `json:"tectonic_aws_external_vpc_id,omitempty" yaml:"vpcID,omitempty"`
	WorkerSubnetIDs []string `json:"tectonic_aws_external_worker_subnet_ids,omitempty" yaml:"workerSubnetIDs,omitempty"`
}

// Master converts master related config.
type Master struct {
	CustomSubnets    map[string]string `json:"tectonic_aws_master_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
	EC2Type          string            `json:"tectonic_aws_master_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
	ExtraSGIDs       []string          `json:"tectonic_aws_master_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
	IAMRoleName      string            `json:"tectonic_aws_master_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
	MasterRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}

// MasterRootVolume converts master rool volume related config.
type MasterRootVolume struct {
	IOPS int    `json:"tectonic_aws_master_root_volume_iops,omitempty" yaml:"iops,omitempty"`
	Size int    `json:"tectonic_aws_master_root_volume_size,omitempty" yaml:"size,omitempty"`
	Type string `json:"tectonic_aws_master_root_volume_type,omitempty" yaml:"type,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	CustomSubnets    map[string]string `json:"tectonic_aws_worker_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
	EC2Type          string            `json:"tectonic_aws_worker_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
	ExtraSGIDs       []string          `json:"tectonic_aws_worker_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
	IAMRoleName      string            `json:"tectonic_aws_worker_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
	LoadBalancers    []string          `json:"tectonic_aws_worker_load_balancers,omitempty" yaml:"loadBalancers,omitempty"`
	WorkerRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}

// WorkerRootVolume converts worker rool volume related config.
type WorkerRootVolume struct {
	IOPS int    `json:"tectonic_aws_worker_root_volume_iops,omitempty" yaml:"iops,omitempty"`
	Size int    `json:"tectonic_aws_worker_root_volume_size,omitempty" yaml:"size,omitempty"`
	Type string `json:"tectonic_aws_worker_root_volume_type,omitempty" yaml:"type,omitempty"`
}
