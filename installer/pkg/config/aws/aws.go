package aws

// AWS converts AWS related config.
type AWS struct {
	AssetsS3BucketName        string `json:"tectonic_aws_assets_s3_bucket_name,omitempty" yaml:"assetsS3BucketName,omitempty"`
	AutoScalingGroupExtraTags string `json:"tectonic_autoscaling_group_extra_tags,omitempty" yaml:"autoScalingGroupExtraTags,omitempty"`
	EC2AMIOverride            string `json:"tectonic_aws_ec2_ami_override,omitempty" yaml:"ec2AMIOverride,omitempty"`
	Etcd                      `json:",inline" yaml:"etcd,omitempty"`
	External                  `json:",inline" yaml:"external,omitempty"`
	ExtraTags                 string `json:"tectonic_aws_extra_tags,omitempty" yaml:"extraTags,omitempty"`
	InstallerRole             string `json:"tectonic_aws_installer_role,omitempty" yaml:"installerRole,omitempty"`
	Master                    `json:",inline" yaml:"master,omitempty"`
	PrivateEndpoints          bool   `json:"tectonic_aws_private_endpoints,omitempty" yaml:"privateEndpoints,omitempty"`
	Profile                   string `json:"tectonic_aws_profile,omitempty" yaml:"profile,omitempty"`
	PublicEndpoints           bool   `json:"tectonic_aws_public_endpoints,omitempty" yaml:"publicEndpoints,omitempty"`
	Region                    string `json:"tectonic_aws_region,omitempty" yaml:"region,omitempty"`
	SSHKey                    string `json:"tectonic_aws_ssh_key,omitempty" yaml:"sshKey,omitempty"`
	VPCCIDRBlock              string `json:"tectonic_aws_vpc_cidr_block,omitempty" yaml:"vpcCIDRBlock,omitempty"`
	Worker                    `json:",inline" yaml:"worker,omitempty"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs string `json:"tectonic_aws_external_master_subnet_ids,omitempty" yaml:"masterSubnetIDs,omitempty"`
	PrivateZone     string `json:"tectonic_aws_external_private_zone,omitempty" yaml:"privateZone,omitempty"`
	VPCID           string `json:"tectonic_aws_external_vpc_id,omitempty" yaml:"vpcID,omitempty"`
	WorkerSubnetIDs string `json:"tectonic_aws_external_worker_subnet_ids,omitempty" yaml:"workerSubnetIDs,omitempty"`
}

// Etcd converts etcd related config.
type Etcd struct {
	EC2Type        string `json:"tectonic_aws_etcd_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
	ExtraSGIDs     string `json:"tectonic_aws_etcd_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
	IAMRoleName    string `json:"tectonic_aws_etcd_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
	EtcdRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}

// EtcdRootVolume converts etcd rool volume related config.
type EtcdRootVolume struct {
	IOPS int    `json:"tectonic_aws_etcd_root_volume_iops,omitempty" yaml:"iops,omitempty"`
	Size int    `json:"tectonic_aws_etcd_root_volume_size,omitempty" yaml:"size,omitempty"`
	Type string `json:"tectonic_aws_etcd_root_volume_type,omitempty" yaml:"type,omitempty"`
}

// Master converts master related config.
type Master struct {
	CustomSubnets    string `json:"tectonic_aws_master_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
	EC2Type          string `json:"tectonic_aws_master_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
	ExtraSGIDs       string `json:"tectonic_aws_master_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
	IAMRoleName      string `json:"tectonic_aws_master_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
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
	CustomSubnets    string `json:"tectonic_aws_worker_custom_subnets,omitempty" yaml:"customSubnets,omitempty"`
	EC2Type          string `json:"tectonic_aws_worker_ec2_type,omitempty" yaml:"ec2Type,omitempty"`
	ExtraSGIDs       string `json:"tectonic_aws_worker_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
	IAMRoleName      string `json:"tectonic_aws_worker_iam_role_name,omitempty" yaml:"iamRoleName,omitempty"`
	LoadBalancers    string `json:"tectonic_aws_worker_load_balancers,omitempty" yaml:"loadBalancers,omitempty"`
	WorkerRootVolume `json:",inline" yaml:"rootVolume,omitempty"`
}

// WorkerRootVolume converts worker rool volume related config.
type WorkerRootVolume struct {
	IOPS int    `json:"tectonic_aws_worker_root_volume_iops,omitempty" yaml:"iops,omitempty"`
	Size int    `json:"tectonic_aws_worker_root_volume_size,omitempty" yaml:"size,omitempty"`
	Type string `json:"tectonic_aws_worker_root_volume_type,omitempty" yaml:"type,omitempty"`
}
