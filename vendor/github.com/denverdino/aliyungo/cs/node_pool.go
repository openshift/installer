package cs

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"net/http"
	"time"
)

type SpotPrice struct {
	InstanceType string `json:"instance_type"`
	PriceLimit   string `json:"price_limit"`
}

type NodePoolInfo struct {
	NodePoolId      string        `json:"nodepool_id"`
	RegionId        common.Region `json:"region_id"`
	Name            string        `json:"name"`
	Created         time.Time     `json:"created"`
	Updated         time.Time     `json:"updated"`
	IsDefault       bool          `json:"is_default"`
	NodePoolType    string        `json:"type"`
	ResourceGroupId string        `json:"resource_group_id"`
}

type ScalingGroup struct {
	VpcId                      string             `json:"vpc_id"`
	VswitchIds                 []string           `json:"vswitch_ids"`
	InstanceTypes              []string           `json:"instance_types"`
	LoginPassword              string             `json:"login_password"`
	KeyPair                    string             `json:"key_pair"`
	SecurityGroupId            string             `json:"security_group_id"`
	SystemDiskCategory         ecs.DiskCategory   `json:"system_disk_category"`
	SystemDiskSize             int64              `json:"system_disk_size"`
	SystemDiskPerformanceLevel string             `json:"system_disk_performance_level"`
	DataDisks                  []NodePoolDataDisk `json:"data_disks"` //支持多个数据盘
	Tags                       []Tag              `json:"tags"`
	ImageId                    string             `json:"image_id"`
	Platform                   string             `json:"platform"`
	// 支持包年包月
	InstanceChargeType string `json:"instance_charge_type"`
	Period             int    `json:"period"`
	PeriodUnit         string `json:"period_unit"`
	AutoRenew          bool   `json:"auto_renew"`
	AutoRenewPeriod    int    `json:"auto_renew_period"`
	// spot实例
	SpotStrategy   string      `json:"spot_strategy"`
	SpotPriceLimit []SpotPrice `json:"spot_price_limit"`

	RdsInstances   []string `json:"rds_instances"`
	ScalingPolicy  string   `json:"scaling_policy"`
	ScalingGroupId string   `json:"scaling_group_id"`

	WorkerSnapshotPolicyId string `json:"worker_system_disk_snapshot_policy_id"`
	// 公网ip
	InternetChargeType      string `json:"internet_charge_type"`
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"`
}

type AutoScaling struct {
	Enable       bool   `json:"enable"`
	MaxInstances int64  `json:"max_instances"`
	MinInstances int64  `json:"min_instances"`
	Type         string `json:"type"`
	// eip
	IsBindEip *bool `json:"is_bond_eip"`
	// value: PayByBandwidth / PayByTraffic
	EipInternetChargeType string `json:"eip_internet_charge_type"`
	// default 5
	EipBandWidth int64 `json:"eip_bandwidth"`
}

type KubernetesConfig struct {
	NodeNameMode string  `json:"node_name_mode"`
	Taints       []Taint `json:"taints"`
	Labels       []Label `json:"labels"`
	CpuPolicy    string  `json:"cpu_policy"`
	UserData     string  `json:"user_data"`

	Runtime           string `json:"runtime,omitempty"`
	RuntimeVersion    string `json:"runtime_version"`
	CmsEnabled        bool   `json:"cms_enabled"`
	OverwriteHostname bool   `json:"overwrite_hostname"`
	Unschedulable     bool   `json:"unschedulable"`
}

// 加密计算节点池
type TEEConfig struct {
	TEEType   string `json:"tee_type"`
	TEEEnable bool   `json:"tee_enable"`
}

// 托管节点池配置
type Management struct {
	Enable      bool        `json:"enable"`
	AutoRepair  bool        `json:"auto_repair"`
	UpgradeConf UpgradeConf `json:"upgrade_config"`
}

type UpgradeConf struct {
	AutoUpgrade       bool  `json:"auto_upgrade"`
	Surge             int64 `json:"surge"`
	SurgePercentage   int64 `json:"surge_percentage,omitempty"`
	MaxUnavailable    int64 `json:"max_unavailable"`
	KeepSurgeOnFailed bool  `json:"keep_surge_on_failed"`
}

type CreateNodePoolRequest struct {
	RegionId         common.Region `json:"region_id"`
	Count            int64         `json:"count"`
	NodePoolInfo     `json:"nodepool_info"`
	ScalingGroup     `json:"scaling_group"`
	KubernetesConfig `json:"kubernetes_config"`
	AutoScaling      `json:"auto_scaling"`
	TEEConfig        `json:"tee_config"`
	Management       `json:"management"`
}

type BasicNodePool struct {
	NodePoolInfo   `json:"nodepool_info"`
	NodePoolStatus `json:"status"`
}

type NodePoolStatus struct {
	TotalNodes   int `json:"total_nodes"`
	OfflineNodes int `json:"offline_nodes"`
	ServingNodes int `json:"serving_nodes"`
	//DesiredNodes int  `json:"desired_nodes"`
	RemovingNodes int    `json:"removing_nodes"`
	FailedNodes   int    `json:"failed_nodes"`
	InitialNodes  int    `json:"initial_nodes"`
	HealthyNodes  int    `json:"healthy_nodes"`
	State         string `json:"state"`
}

type NodePoolDetail struct {
	BasicNodePool
	KubernetesConfig `json:"kubernetes_config"`
	ScalingGroup     `json:"scaling_group"`
	AutoScaling      `json:"auto_scaling"`
	TEEConfig        `json:"tee_config"`
	Management       `json:"management"`
}

type CreateNodePoolResponse struct {
	Response
	NodePoolID string `json:"nodepool_id"`
	Message    string `json:"Message"`
	TaskID     string `json:"task_id"`
}

type UpdateNodePoolRequest struct {
	RegionId         common.Region `json:"region_id"`
	Count            int64         `json:"count"`
	NodePoolInfo     `json:"nodepool_info"`
	ScalingGroup     `json:"scaling_group"`
	KubernetesConfig `json:"kubernetes_config"`
	AutoScaling      `json:"auto_scaling"`
	Management       `json:"management"`
}

type NodePoolsDetail struct {
	Response
	NodePools []NodePoolDetail `json:"nodepools"`
}

func (client *Client) CreateNodePool(request *CreateNodePoolRequest, clusterId string) (*CreateNodePoolResponse, error) {
	response := &CreateNodePoolResponse{}
	err := client.Invoke(request.RegionId, http.MethodPost, fmt.Sprintf("/clusters/%s/nodepools", clusterId), nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) DescribeNodePoolDetail(clusterId, nodePoolId string) (*NodePoolDetail, error) {
	nodePool := &NodePoolDetail{}
	err := client.Invoke("", http.MethodGet, fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, nodePoolId), nil, nil, nodePool)
	if err != nil {
		return nil, err
	}
	return nodePool, nil
}

func (client *Client) DescribeClusterNodePools(clusterId string) (*[]NodePoolDetail, error) {
	nodePools := &NodePoolsDetail{}
	err := client.Invoke("", http.MethodGet, fmt.Sprintf("/clusters/%s/nodepools", clusterId), nil, nil, nodePools)
	if err != nil {
		return nil, err
	}
	return &nodePools.NodePools, nil
}

func (client *Client) UpdateNodePool(clusterId string, nodePoolId string, request *UpdateNodePoolRequest) (*Response, error) {
	response := &Response{}
	err := client.Invoke(request.RegionId, http.MethodPut, fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, nodePoolId), nil, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) DeleteNodePool(clusterId, nodePoolId string) error {
	return client.Invoke("", http.MethodDelete, fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, nodePoolId), nil, nil, nil)
}
