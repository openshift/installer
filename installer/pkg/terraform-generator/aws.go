package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// AWS defines all variables for this platform.
type AWS struct {
	AssetsS3BucketName      string `json:"tectonic_aws_assets_s3_bucket_name,omitempty"`
	EtcdEC2Type             string `json:"tectonic_aws_etcd_ec2_type,omitempty"`
	EtcdExtraSGIDs          string `json:"tectonic_aws_etcd_extra_sg_ids,omitempty"`
	EtcdIAMRoleName         string `json:"tectonic_aws_etcd_iam_role_name,omitempty"`
	EtcdRootVolumeIOPS      string `json:"tectonic_aws_etcd_root_volume_iops,omitempty"`
	EtcdRootVolumeSize      string `json:"tectonic_aws_etcd_root_volume_size,omitempty"`
	EtcdRootVolumeType      string `json:"tectonic_aws_etcd_root_volume_type,omitempty"`
	ExternalMasterSubnetIDs string `json:"tectonic_aws_external_master_subnet_ids,omitempty"`
	ExternalPrivateZone     string `json:"tectonic_aws_external_private_zone,omitempty"`
	ExternalVPCID           string `json:"tectonic_aws_external_vpc_id,omitempty"`
	ExternalWorkerSubnetIDs string `json:"tectonic_aws_external_worker_subnet_ids,omitempty"`
	ExtraTags               string `json:"tectonic_aws_extra_tags,omitempty"`
	MasterCustomSubnets     string `json:"tectonic_aws_master_custom_subnets,omitempty"`
	MasterEC2Type           string `json:"tectonic_aws_master_ec2_type,omitempty"`
	MasterExtraSGIDs        string `json:"tectonic_aws_master_extra_sg_ids,omitempty"`
	MasterIAMRoleName       string `json:"tectonic_aws_master_iam_role_name,omitempty"`
	MasterRootVolumeIOPS    string `json:"tectonic_aws_master_root_volume_iops,omitempty"`
	MasterRootVolumeSize    string `json:"tectonic_aws_master_root_volume_size,omitempty"`
	MasterRootVolumeType    string `json:"tectonic_aws_master_root_volume_type,omitempty"`
	PrivateEndpoints        string `json:"tectonic_aws_private_endpoints,omitempty"`
	Profile                 string `json:"tectonic_aws_profile,omitempty"`
	PublicEndpoints         string `json:"tectonic_aws_public_endpoints,omitempty"`
	Region                  string `json:"tectonic_aws_region,omitempty"`
	SSHKey                  string `json:"tectonic_aws_ssh_key,omitempty"`
	VPCCIDRBlock            string `json:"tectonic_aws_vpc_cidr_block,omitempty"`
	WorkerCustomSubnets     string `json:"tectonic_aws_worker_custom_subnets,omitempty"`
	WorkerEC2Type           string `json:"tectonic_aws_worker_ec2_type,omitempty"`
	WorkerExtraSGIDs        string `json:"tectonic_aws_worker_extra_sg_ids,omitempty"`
	WorkerIAMRoleName       string `json:"tectonic_aws_worker_iam_role_name,omitempty"`
	WorkerLoadBalancers     string `json:"tectonic_aws_worker_load_balancers,omitempty"`
	WorkerRootVolumeIOPS    string `json:"tectonic_aws_worker_root_volume_iops,omitempty"`
	WorkerRootVolumeSize    string `json:"tectonic_aws_worker_root_volume_size,omitempty"`
	WorkerRootVolumeType    string `json:"tectonic_aws_worker_root_volume_type,omitempty"`
}

// NewAWS returns the config for AWS.
func NewAWS(cluster config.Cluster) AWS {
	return AWS{
		AssetsS3BucketName:      cluster.AWS.AssetsS3BucketName,
		EtcdEC2Type:             cluster.AWS.Etcd.EC2Type,
		EtcdExtraSGIDs:          cluster.AWS.Etcd.ExtraSGIDs,
		EtcdIAMRoleName:         cluster.AWS.Etcd.IAMRoleName,
		EtcdRootVolumeIOPS:      cluster.AWS.Etcd.RootVolume.IOPS,
		EtcdRootVolumeSize:      cluster.AWS.Etcd.RootVolume.Size,
		EtcdRootVolumeType:      cluster.AWS.Etcd.RootVolume.Type,
		ExternalMasterSubnetIDs: cluster.AWS.External.MasterSubnetIDs,
		ExternalPrivateZone:     cluster.AWS.External.PrivateZone,
		ExternalVPCID:           cluster.AWS.External.VPCID,
		ExternalWorkerSubnetIDs: cluster.AWS.External.WorkerSubnetIDs,
		ExtraTags:               cluster.AWS.ExtraTags,
		MasterCustomSubnets:     cluster.AWS.Master.CustomSubnets,
		MasterEC2Type:           cluster.AWS.Master.EC2Type,
		MasterExtraSGIDs:        cluster.AWS.Master.ExtraSGIDs,
		MasterIAMRoleName:       cluster.AWS.Master.IAMRoleName,
		MasterRootVolumeIOPS:    cluster.AWS.Master.RootVolume.IOPS,
		MasterRootVolumeSize:    cluster.AWS.Master.RootVolume.Size,
		MasterRootVolumeType:    cluster.AWS.Master.RootVolume.Type,
		PrivateEndpoints:        cluster.AWS.PrivateEndpoints,
		Profile:                 cluster.AWS.Profile,
		PublicEndpoints:         cluster.AWS.PublicEndpoints,
		Region:                  cluster.AWS.Region,
		SSHKey:                  cluster.AWS.SSHKey,
		VPCCIDRBlock:            cluster.AWS.VPCCIDRBlock,
		WorkerCustomSubnets:     cluster.AWS.Worker.CustomSubnets,
		WorkerEC2Type:           cluster.AWS.Worker.EC2Type,
		WorkerExtraSGIDs:        cluster.AWS.Worker.ExtraSGIDs,
		WorkerIAMRoleName:       cluster.AWS.Worker.IAMRoleName,
		WorkerLoadBalancers:     cluster.AWS.Worker.LoadBalancers,
		WorkerRootVolumeIOPS:    cluster.AWS.Worker.RootVolume.IOPS,
		WorkerRootVolumeSize:    cluster.AWS.Worker.RootVolume.Size,
		WorkerRootVolumeType:    cluster.AWS.Worker.RootVolume.Type,
	}
}
