package cs

import (
	"encoding/json"
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"net/http"
	"time"

	"github.com/denverdino/aliyungo/common"
)

type TaskState string

const (
	// task status
	Task_Status_Running = "running"
	Task_Status_Success = "Success"
	Task_Status_Failed  = "Failed"

	// upgrade state
	UpgradeStep_NotStart    = "not_start"
	UpgradeStep_Prechecking = "prechecking"
	UpgradeStep_Upgrading   = "upgrading"
	UpgradeStep_Pause       = "pause"
	UpgradeStep_Success     = "success"
)

type WeeklyPeriod string
type MaintenanceTime string

// maintenance window
type MaintenanceWindow struct {
	Enable          bool            `json:"enable"`
	MaintenanceTime MaintenanceTime `json:"maintenance_time"`
	Duration        string          `json:"duration"`
	Recurrence      string          `json:"recurrence,omitempty"`
	WeeklyPeriod    WeeklyPeriod    `json:"weekly_period"`
}

//modify cluster,include DeletionProtection and so on
type ModifyClusterArgs struct {
	DeletionProtection bool              `json:"deletion_protection"`
	ResourceGroupId    string            `json:"resource_group_id"`
	MaintenanceWindow  MaintenanceWindow `json:"maintenance_window"`
}

type UpgradeClusterArgs struct {
	Version string `json:"version"`
}

type UpgradeClusterResult struct {
	Status           TaskState `json:"status"`
	PrecheckReportId string    `json:"precheck_report_id"`
	UpgradeStep      string    `json:"upgrade_step"`
	ErrorMessage     string    `json:"error_message"`
	*UpgradeTask     `json:"upgrade_task,omitempty"`
}

type UpgradeTask struct {
	FieldRetries    int           `json:"retries,omitempty"`
	FieldCreatedAt  time.Time     `json:"created_at"`
	FieldMessage    string        `json:"message,omitempty"`
	FieldStatus     string        `json:"status"` // empty|running|success|failed
	FieldFinishedAt time.Time     `json:"finished_at,omitempty"`
	UpgradeStatus   UpgradeStatus `json:"upgrade_status"`
}

type UpgradeStatus struct {
	State      string  `json:"state"`
	Phase      string  `json:"phase"` // {Master1, Master2, Master3, Nodes}
	Total      int     `json:"total"`
	Succeeded  int     `json:"succeeded"`
	Failed     string  `json:"failed"`
	Events     []Event `json:"events"`
	IsCanceled bool    `json:"is_canceled"`
}

type Event struct {
	Timestamp time.Time
	Type      string
	Reason    string
	Message   string
	Source    string
}

//modify cluster
func (client *Client) ModifyCluster(clusterId string, args *ModifyClusterArgs) error {
	return client.Invoke("", http.MethodPut, "/api/v2/clusters/"+clusterId, nil, args, nil)
}

//upgrade cluster
func (client *Client) UpgradeCluster(clusterId string, args *UpgradeClusterArgs) error {
	return client.Invoke("", http.MethodPost, fmt.Sprintf("/api/v2/clusters/%s/upgrade", clusterId), nil, args, nil)
}

//cancel upgrade cluster
func (client *Client) CancelUpgradeCluster(clusterId string) error {
	return client.Invoke("", http.MethodPost, fmt.Sprintf("/api/v2/clusters/%s/upgrade/cancel", clusterId), nil, nil, nil)
}

func (client *Client) QueryUpgradeClusterResult(clusterId string) (*UpgradeClusterResult, error) {
	cluster := &UpgradeClusterResult{}
	err := client.Invoke("", http.MethodGet, fmt.Sprintf("/api/v2/clusters/%s/upgrade/status", clusterId), nil, nil, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

//Cluster Info
type KubernetesClusterType string

var (
	DelicatedKubernetes  = KubernetesClusterType("Kubernetes")
	ManagedKubernetes    = KubernetesClusterType("ManagedKubernetes")
	ServerlessKubernetes = KubernetesClusterType("Ask")
)

type ProxyMode string

var (
	IPTables = "iptables"
	IPVS     = "ipvs"
)

type Effect string

var (
	TaintNoExecute        = "NoExecute"
	TaintNoSchedule       = "NoSchedule"
	TaintPreferNoSchedule = "PreferNoSchedule"
)

type ClusterArgs struct {
	DisableRollback bool `json:"disable_rollback"`
	TimeoutMins     int  `json:"timeout_mins"`

	Name               string                `json:"name"`
	ClusterType        KubernetesClusterType `json:"cluster_type"`
	Profile            string                `json:"profile"`
	KubernetesVersion  string                `json:"kubernetes_version"`
	DeletionProtection bool                  `json:"deletion_protection"`

	NodeCidrMask string `json:"node_cidr_mask"`
	UserCa       string `json:"user_ca"`

	OsType   string `json:"os_type"`
	Platform string `json:"platform"`

	UserData string `json:"user_data"`

	NodePortRange string `json:"node_port_range"`
	NodeNameMode  string `json:"node_name_mode"`

	//ImageId
	ImageId string `json:"image_id"`

	PodVswitchIds []string `json:"pod_vswitch_ids"` // eni多网卡模式下，需要传额外的vswitchid给addon

	LoginPassword string `json:"login_password"` //和KeyPair 二选一
	KeyPair       string `json:"key_pair"`       ////LoginPassword 二选一

	RegionId      common.Region `json:"region_id"`
	VpcId         string        `json:"vpcid"`
	ContainerCidr string        `json:"container_cidr"`
	ServiceCidr   string        `json:"service_cidr"`

	CloudMonitorFlags bool `json:"cloud_monitor_flags"`

	SecurityGroupId           string    `json:"security_group_id"`
	IsEnterpriseSecurityGroup bool      `json:"is_enterprise_security_group"`
	EndpointPublicAccess      bool      `json:"endpoint_public_access"`
	LoadBalancerSpec          string    `json:"load_balancer_spec"` //api server slb实例规格
	ProxyMode                 ProxyMode `json:"proxy_mode"`
	SnatEntry                 bool      `json:"snat_entry"`
	ResourceGroupId           string    `json:"resource_group_id"`

	Addons []Addon `json:"addons"`
	Tags   []Tag   `json:"tags"`

	Taints []Taint `json:"taints"`

	ApiAudiences          string            `json:"api_audiences,omitempty"`
	ServiceAccountIssuer  string            `json:"service_account_issuer,omitempty"`
	CustomSAN             string            `json:"custom_san,omitempty"`
	ClusterSpec           string            `json:"cluster_spec"`
	Timezone              string            `json:"timezone"`
	ClusterDomain         string            `json:"cluster_domain"`
	RdsInstances          []string          `json:"rds_instances"`
	EncryptionProviderKey string            `json:"encryption_provider_key"`
	MaintenanceWindow     MaintenanceWindow `json:"maintenance_window"`

	//controlplane log parms
	ControlplaneLogProject string   `json:"controlplane_log_project"`
	ControlplaneLogTTL     string   `json:"controlplane_log_ttl"`
	ControlplaneComponents []string `json:"controlplane_log_components"`

	// Operating system hardening
	SocEnabled *bool `json:"soc_enabled"`
	CisEnabled *bool `json:"cis_enabled"`
}

//addon
type Addon struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Config   string `json:"config"`
	Disabled bool   `json:"disabled"` // 是否禁止默认安装
}

//taint
type Taint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect Effect `json:"effect"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MasterArgs struct {
	MasterCount         int      `json:"master_count"`
	MasterVSwitchIds    []string `json:"master_vswitch_ids"`
	MasterInstanceTypes []string `json:"master_instance_types"`

	MasterInstanceChargeType string `json:"master_instance_charge_type"`
	MasterPeriod             int    `json:"master_period"`
	MasterPeriodUnit         string `json:"master_period_unit"`

	MasterAutoRenew       bool `json:"master_auto_renew"`
	MasterAutoRenewPeriod int  `json:"master_auto_renew_period"`

	MasterSystemDiskCategory         ecs.DiskCategory `json:"master_system_disk_category"`
	MasterSystemDiskSize             int64            `json:"master_system_disk_size"`
	MasterSystemDiskPerformanceLevel string           `json:"master_system_disk_performance_level"`

	MasterDataDisks []DataDisk `json:"master_data_disks"` //支持多个数据盘

	//support hpc/scc
	MasterHpcClusterId    string `json:"master_hpc_cluster_id"`
	MasterDeploymentSetId string `json:"master_deploymentset_id"`

	//master node deletion protection
	MasterDeletionProtection *bool `json:"master_deletion_protection"`

	// disk snapshot policy
	MasterSnapshotPolicyId string `json:"master_system_disk_snapshot_policy_id"`
}

type WorkerArgs struct {
	WorkerVSwitchIds    []string `json:"worker_vswitch_ids"`
	WorkerInstanceTypes []string `json:"worker_instance_types"`

	NumOfNodes int64 `json:"num_of_nodes"`

	WorkerInstanceChargeType string `json:"worker_instance_charge_type"`
	WorkerPeriod             int    `json:"worker_period"`
	WorkerPeriodUnit         string `json:"worker_period_unit"`

	WorkerAutoRenew       bool `json:"worker_auto_renew"`
	WorkerAutoRenewPeriod int  `json:"worker_auto_renew_period"`

	WorkerSystemDiskCategory         ecs.DiskCategory `json:"worker_system_disk_category"`
	WorkerSystemDiskSize             int64            `json:"worker_system_disk_size"`
	WorkerSystemDiskPerformanceLevel string           `json:"worker_system_disk_performance_level"`

	WorkerDataDisk  bool       `json:"worker_data_disk"`
	WorkerDataDisks []DataDisk `json:"worker_data_disks"` //支持多个数据盘

	WorkerHpcClusterId    string `json:"worker_hpc_cluster_id"`
	WorkerDeploymentSetId string `json:"worker_deploymentset_id"`

	//worker node deletion protection
	WorkerDeletionProtection *bool `json:"worker_deletion_protection"`

	//Runtime only for worker nodes
	Runtime Runtime `json:"runtime"`

	// disk snapshot policy
	WorkerSnapshotPolicyId string `json:"worker_system_disk_snapshot_policy_id"`
}

type ScaleOutKubernetesClusterRequest struct {
	LoginPassword string `json:"login_password"` //和KeyPair 二选一
	KeyPair       string `json:"key_pair"`       ////LoginPassword 二选一

	WorkerVSwitchIds    []string `json:"worker_vswitch_ids"`
	WorkerInstanceTypes []string `json:"worker_instance_types"`

	WorkerInstanceChargeType string `json:"worker_instance_charge_type"`
	WorkerPeriod             int    `json:"worker_period"`
	WorkerPeriodUnit         string `json:"worker_period_unit"`

	WorkerAutoRenew       bool `json:"worker_auto_renew"`
	WorkerAutoRenewPeriod int  `json:"worker_auto_renew_period"`

	WorkerSystemDiskCategory         ecs.DiskCategory `json:"worker_system_disk_category"`
	WorkerSystemDiskSize             int64            `json:"worker_system_disk_size"`
	WorkerSnapshotPolicyId           string           `json:"worker_system_disk_snapshot_policy_id"`
	WorkerSystemDiskPerformanceLevel string           `json:"worker_system_disk_performance_level"`

	WorkerDataDisk  bool       `json:"worker_data_disk"`
	WorkerDataDisks []DataDisk `json:"worker_data_disks"` //支持多个数据盘

	Tags    []Tag   `json:"tags"`
	Taints  []Taint `json:"taints"`
	ImageId string  `json:"image_id"`

	UserData string `json:"user_data"`

	Count             int64    `json:"count"`
	CpuPolicy         string   `json:"cpu_policy"`
	Runtime           Runtime  `json:"runtime"`
	CloudMonitorFlags bool     `json:"cloud_monitor_flags"`
	RdsInstances      []string ` json:"rds_instances"`
}

//datadiks
type DataDisk struct {
	Category             string `json:"category"`
	KMSKeyId             string `json:"kms_key_id"`
	Encrypted            string `json:"encrypted"` // true|false
	Device               string `json:"device"`    //  could be /dev/xvd[a-z]. If not specification, will use default value.
	Size                 string `json:"size"`
	DiskName             string `json:"name"`
	AutoSnapshotPolicyId string `json:"auto_snapshot_policy_id"`
	PerformanceLevel     string `json:"performance_Level"`
}

type NodePoolDataDisk struct {
	Category             string `json:"category"`
	KMSKeyId             string `json:"kms_key_id"`
	Encrypted            string `json:"encrypted"` // true|false
	Device               string `json:"device"`    //  could be /dev/xvd[a-z]. If not specification, will use default value.
	Size                 int    `json:"size"`
	DiskName             string `json:"name"`
	AutoSnapshotPolicyId string `json:"auto_snapshot_policy_id"`
	PerformanceLevel     string `json:"performance_Level"`
}

//runtime
type Runtime struct {
	Name                       string   `json:"name"`
	Version                    string   `json:"version"`
	RuntimeClass               []string `json:"runtimeClass,omitempty"`
	Exist                      bool     `json:"exist"`
	AvailableNetworkComponents []string `json:"availableNetworkComponents,omitempty"`
}

//DelicatedKubernetes
type DelicatedKubernetesClusterCreationRequest struct {
	ClusterArgs
	MasterArgs
	WorkerArgs
}

//ManagedKubernetes
type ManagedKubernetesClusterCreationRequest struct {
	ClusterArgs
	WorkerArgs
}

//Validate

//Create DelicatedKubernetes Cluster
func (client *Client) CreateDelicatedKubernetesCluster(request *DelicatedKubernetesClusterCreationRequest) (*ClusterCommonResponse, error) {
	response := &ClusterCommonResponse{}
	path := "/clusters"
	if request.ResourceGroupId != "" {
		// 创建集群到指定资源组
		path = fmt.Sprintf("/resource_groups/%s/clusters", request.ResourceGroupId)
	}
	err := client.Invoke(request.RegionId, http.MethodPost, path, nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//Create ManagedKubernetes Cluster
func (client *Client) CreateManagedKubernetesCluster(request *ManagedKubernetesClusterCreationRequest) (*ClusterCommonResponse, error) {
	response := &ClusterCommonResponse{}
	path := "/clusters"
	if request.ResourceGroupId != "" {
		// 创建集群到指定资源组
		path = fmt.Sprintf("/resource_groups/%s/clusters", request.ResourceGroupId)
	}
	err := client.Invoke(request.RegionId, http.MethodPost, path, nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//ScaleKubernetesCluster
func (client *Client) ScaleOutKubernetesCluster(clusterId string, request *ScaleOutKubernetesClusterRequest) (*ClusterCommonResponse, error) {
	response := &ClusterCommonResponse{}
	err := client.Invoke("", http.MethodPost, fmt.Sprintf("/api/v2/clusters/%s", clusterId), nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//DeleteClusterNodes
type DeleteKubernetesClusterNodesRequest struct {
	ReleaseNode bool     `json:"release_node"` //if set to true, the ecs instance will be released
	Nodes       []string `json:"nodes"`        //the format is regionId.instanceId|Ip ,for example  cn-hangzhou.192.168.1.2 or cn-hangzhou.i-abc
	DrainNode   bool     `json:"drain_node"`   //same as Nodes
}

//DeleteClusterNodes
func (client *Client) DeleteKubernetesClusterNodes(clusterId string, request *DeleteKubernetesClusterNodesRequest) (*common.Response, error) {
	response := &common.Response{}
	err := client.Invoke("", http.MethodPost, fmt.Sprintf("/api/v2/clusters/%s/nodes/remove", clusterId), nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//Cluster definition
type KubernetesClusterDetail struct {
	RegionId common.Region `json:"region_id"`

	Name        string                `json:"name"`
	ClusterId   string                `json:"cluster_id"`
	Size        int64                 `json:"size"`
	ClusterType KubernetesClusterType `json:"cluster_type"`
	Profile     string                `json:"profile"`

	VpcId                 string `json:"vpc_id"`
	VSwitchIds            string `json:"vswitch_id"`
	SecurityGroupId       string `json:"security_group_id"`
	IngressLoadbalancerId string `json:"external_loadbalancer_id"`
	ResourceGroupId       string `json:"resource_group_id"`
	NetworkMode           string `json:"network_mode"`
	ContainerCIDR         string `json:"subnet_cidr"`

	Tags  []Tag  `json:"tags"`
	State string `json:"state"`

	InitVersion        string `json:"init_version"`
	CurrentVersion     string `json:"current_version"`
	PrivateZone        bool   `json:"private_zone"`
	DeletionProtection bool   `json:"deletion_protection"`
	MetaData           string `json:"meta_data"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	WorkerRamRoleName string            `json:"worker_ram_role_name"`
	ClusterSpec       string            `json:"cluster_spec"`
	OSType            string            `json:"os_type"`
	MasterURL         string            `json:"master_url"`
	MaintenanceWindow MaintenanceWindow `json:"maintenance_window"`
}

//GetMetaData
func (c *KubernetesClusterDetail) GetMetaData() map[string]interface{} {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(c.MetaData), &m)
	return m
}

//查询集群详情
func (client *Client) DescribeKubernetesClusterDetail(clusterId string) (*KubernetesClusterDetail, error) {
	cluster := &KubernetesClusterDetail{}
	err := client.Invoke("", http.MethodGet, "/clusters/"+clusterId, nil, nil, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

//DeleteKubernetesCluster
func (client *Client) DeleteKubernetesCluster(clusterId string) error {
	return client.Invoke("", http.MethodDelete, "/clusters/"+clusterId, nil, nil, nil)
}
