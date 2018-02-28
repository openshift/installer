package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// GovCloud defines all variables for this platform.
type GovCloud struct {
	AssetsS3BucketName      string `json:"tectonic_govcloud_assets_s3_bucket_name,omitempty"`
	DNSServerIP             string `json:"tectonic_govcloud_dns_server_ip,omitempty"`
	EtcdEC2Type             string `json:"tectonic_govcloud_etcd_ec2_type,omitempty"`
	EtcdExtraSGIDs          string `json:"tectonic_govcloud_etcd_extra_sg_ids,omitempty"`
	EtcdRootVolumeIOPS      int    `json:"tectonic_govcloud_etcd_root_volume_iops,omitempty"`
	EtcdRootVolumeSize      int    `json:"tectonic_govcloud_etcd_root_volume_size,omitempty"`
	EtcdRootVolumeType      string `json:"tectonic_govcloud_etcd_root_volume_type,omitempty"`
	ExternalMasterSubnetIDs string `json:"tectonic_govcloud_external_master_subnet_ids,omitempty"`
	ExternalPrivateZone     string `json:"tectonic_govcloud_external_private_zone,omitempty"`
	ExternalVPCID           string `json:"tectonic_govcloud_external_vpc_id,omitempty"`
	ExternalWorkerSubnetIDs string `json:"tectonic_govcloud_external_worker_subnet_ids,omitempty"`
	ExtraTags               string `json:"tectonic_govcloud_extra_tags,omitempty"`
	MasterCustomSubnets     string `json:"tectonic_govcloud_master_custom_subnets,omitempty"`
	MasterEC2Type           string `json:"tectonic_govcloud_master_ec2_type,omitempty"`
	MasterExtraSGIDs        string `json:"tectonic_govcloud_master_extra_sg_ids,omitempty"`
	MasterIAMRoleName       string `json:"tectonic_govcloud_master_iam_role_name,omitempty"`
	MasterRootVolumeIOPS    int    `json:"tectonic_govcloud_master_root_volume_iops,omitempty"`
	MasterRootVolumeSize    int    `json:"tectonic_govcloud_master_root_volume_size,omitempty"`
	MasterRootVolumeType    string `json:"tectonic_govcloud_master_root_volume_type,omitempty"`
	PrivateEndpoints        bool   `json:"tectonic_govcloud_private_endpoints,omitempty"`
	Profile                 string `json:"tectonic_govcloud_profile,omitempty"`
	PublicEndpoints         bool   `json:"tectonic_govcloud_public_endpoints,omitempty"`
	SSHKey                  string `json:"tectonic_govcloud_ssh_key,omitempty"`
	VPCCIDRBlock            string `json:"tectonic_govcloud_vpc_cidr_block,omitempty"`
	WorkerCustomSubnets     string `json:"tectonic_govcloud_worker_custom_subnets,omitempty"`
	WorkerEC2Type           string `json:"tectonic_govcloud_worker_ec2_type,omitempty"`
	WorkerExtraSGIDs        string `json:"tectonic_govcloud_worker_extra_sg_ids,omitempty"`
	WorkerIAMRoleName       string `json:"tectonic_govcloud_worker_iam_role_name,omitempty"`
	WorkerLoadBalancers     string `json:"tectonic_govcloud_worker_load_balancers,omitempty"`
	WorkerRootVolumeIOPS    int    `json:"tectonic_govcloud_worker_root_volume_iops,omitempty"`
	WorkerRootVolumeSize    int    `json:"tectonic_govcloud_worker_root_volume_size,omitempty"`
	WorkerRootVolumeType    string `json:"tectonic_govcloud_worker_root_volume_type,omitempty"`
}

// NewGovCloud returns the config for GovCloud.
func NewGovCloud(cluster config.Cluster) GovCloud {
	return GovCloud{
		AssetsS3BucketName:      cluster.GovCloud.AWS.AssetsS3BucketName,
		DNSServerIP:             cluster.GovCloud.DNSServerIP,
		EtcdEC2Type:             cluster.GovCloud.AWS.Etcd.EC2Type,
		EtcdExtraSGIDs:          cluster.GovCloud.AWS.Etcd.ExtraSGIDs,
		EtcdRootVolumeIOPS:      cluster.GovCloud.AWS.Etcd.RootVolume.IOPS,
		EtcdRootVolumeSize:      cluster.GovCloud.AWS.Etcd.RootVolume.Size,
		EtcdRootVolumeType:      cluster.GovCloud.AWS.Etcd.RootVolume.Type,
		ExternalMasterSubnetIDs: cluster.GovCloud.AWS.External.MasterSubnetIDs,
		ExternalPrivateZone:     cluster.GovCloud.AWS.External.PrivateZone,
		ExternalVPCID:           cluster.GovCloud.AWS.External.VPCID,
		ExternalWorkerSubnetIDs: cluster.GovCloud.AWS.External.WorkerSubnetIDs,
		ExtraTags:               cluster.GovCloud.AWS.ExtraTags,
		MasterCustomSubnets:     cluster.GovCloud.AWS.Master.CustomSubnets,
		MasterEC2Type:           cluster.GovCloud.AWS.Master.EC2Type,
		MasterExtraSGIDs:        cluster.GovCloud.AWS.Master.ExtraSGIDs,
		MasterIAMRoleName:       cluster.GovCloud.AWS.Master.IAMRoleName,
		MasterRootVolumeIOPS:    cluster.GovCloud.AWS.Master.RootVolume.IOPS,
		MasterRootVolumeSize:    cluster.GovCloud.AWS.Master.RootVolume.Size,
		MasterRootVolumeType:    cluster.GovCloud.AWS.Master.RootVolume.Type,
		PrivateEndpoints:        cluster.GovCloud.AWS.PrivateEndpoints,
		Profile:                 cluster.GovCloud.AWS.Profile,
		PublicEndpoints:         cluster.GovCloud.AWS.PublicEndpoints,
		SSHKey:                  cluster.GovCloud.AWS.SSHKey,
		VPCCIDRBlock:            cluster.GovCloud.AWS.VPCCIDRBlock,
		WorkerCustomSubnets:     cluster.GovCloud.AWS.Worker.CustomSubnets,
		WorkerEC2Type:           cluster.GovCloud.AWS.Worker.EC2Type,
		WorkerExtraSGIDs:        cluster.GovCloud.AWS.Worker.ExtraSGIDs,
		WorkerIAMRoleName:       cluster.GovCloud.AWS.Worker.IAMRoleName,
		WorkerLoadBalancers:     cluster.GovCloud.AWS.Worker.LoadBalancers,
		WorkerRootVolumeIOPS:    cluster.GovCloud.AWS.Worker.RootVolume.IOPS,
		WorkerRootVolumeSize:    cluster.GovCloud.AWS.Worker.RootVolume.Size,
		WorkerRootVolumeType:    cluster.GovCloud.AWS.Worker.RootVolume.Type,
	}
}
