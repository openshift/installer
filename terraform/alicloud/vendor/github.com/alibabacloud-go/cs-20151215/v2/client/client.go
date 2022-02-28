// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package client

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	endpointutil "github.com/alibabacloud-go/endpoint-util/service"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Runtime struct {
	// 容器运行时名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 容器运行时版本
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s Runtime) String() string {
	return tea.Prettify(s)
}

func (s Runtime) GoString() string {
	return s.String()
}

func (s *Runtime) SetName(v string) *Runtime {
	s.Name = &v
	return s
}

func (s *Runtime) SetVersion(v string) *Runtime {
	s.Version = &v
	return s
}

type Taint struct {
	// key值。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// value值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
	// 污点生效策略。
	Effect *string `json:"effect,omitempty" xml:"effect,omitempty"`
}

func (s Taint) String() string {
	return tea.Prettify(s)
}

func (s Taint) GoString() string {
	return s.String()
}

func (s *Taint) SetKey(v string) *Taint {
	s.Key = &v
	return s
}

func (s *Taint) SetValue(v string) *Taint {
	s.Value = &v
	return s
}

func (s *Taint) SetEffect(v string) *Taint {
	s.Effect = &v
	return s
}

type DataDisk struct {
	// 数据盘类型
	Category *string `json:"category,omitempty" xml:"category,omitempty"`
	// 数据盘大小，取值范围：40～32767
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// 是否对数据盘加密。
	Encrypted *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	// 开启云盘备份时的自动备份策略。
	AutoSnapshotPolicyId *string `json:"auto_snapshot_policy_id,omitempty" xml:"auto_snapshot_policy_id,omitempty"`
}

func (s DataDisk) String() string {
	return tea.Prettify(s)
}

func (s DataDisk) GoString() string {
	return s.String()
}

func (s *DataDisk) SetCategory(v string) *DataDisk {
	s.Category = &v
	return s
}

func (s *DataDisk) SetSize(v int64) *DataDisk {
	s.Size = &v
	return s
}

func (s *DataDisk) SetEncrypted(v string) *DataDisk {
	s.Encrypted = &v
	return s
}

func (s *DataDisk) SetAutoSnapshotPolicyId(v string) *DataDisk {
	s.AutoSnapshotPolicyId = &v
	return s
}

type Tag struct {
	// key值。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// value值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s Tag) String() string {
	return tea.Prettify(s)
}

func (s Tag) GoString() string {
	return s.String()
}

func (s *Tag) SetKey(v string) *Tag {
	s.Key = &v
	return s
}

func (s *Tag) SetValue(v string) *Tag {
	s.Value = &v
	return s
}

type Addon struct {
	// 插件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 插件配置参数。
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
	// 是否禁止默认安装。true | false。
	Disabled *bool `json:"disabled,omitempty" xml:"disabled,omitempty"`
}

func (s Addon) String() string {
	return tea.Prettify(s)
}

func (s Addon) GoString() string {
	return s.String()
}

func (s *Addon) SetName(v string) *Addon {
	s.Name = &v
	return s
}

func (s *Addon) SetConfig(v string) *Addon {
	s.Config = &v
	return s
}

func (s *Addon) SetDisabled(v bool) *Addon {
	s.Disabled = &v
	return s
}

type MaintenanceWindow struct {
	// 是否开启维护窗口。默认值false。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 维护起始时间。Golang标准时间格式"15:04:05Z"。
	MaintenanceTime *string `json:"maintenance_time,omitempty" xml:"maintenance_time,omitempty"`
	// 维护时长。取值范围1～24，单位为小时。 默认值：3h。
	Duration *string `json:"duration,omitempty" xml:"duration,omitempty"`
	// 维护周期。取值范围为:Monday~Sunday，多个值用逗号分隔。 默认值：Thursday。
	WeeklyPeriod *string `json:"weekly_period,omitempty" xml:"weekly_period,omitempty"`
}

func (s MaintenanceWindow) String() string {
	return tea.Prettify(s)
}

func (s MaintenanceWindow) GoString() string {
	return s.String()
}

func (s *MaintenanceWindow) SetEnable(v bool) *MaintenanceWindow {
	s.Enable = &v
	return s
}

func (s *MaintenanceWindow) SetMaintenanceTime(v string) *MaintenanceWindow {
	s.MaintenanceTime = &v
	return s
}

func (s *MaintenanceWindow) SetDuration(v string) *MaintenanceWindow {
	s.Duration = &v
	return s
}

func (s *MaintenanceWindow) SetWeeklyPeriod(v string) *MaintenanceWindow {
	s.WeeklyPeriod = &v
	return s
}

type ListTagResourcesRequest struct {
	// 集群ID列表。
	ResourceIds []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	// 资源类型，只支持Cluster
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 地域ID
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 按标签查找。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 下一次查询Token。
	NextToken *string `json:"next_token,omitempty" xml:"next_token,omitempty"`
}

func (s ListTagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesRequest) GoString() string {
	return s.String()
}

func (s *ListTagResourcesRequest) SetResourceIds(v []*string) *ListTagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *ListTagResourcesRequest) SetResourceType(v string) *ListTagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesRequest) SetRegionId(v string) *ListTagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *ListTagResourcesRequest) SetTags(v []*Tag) *ListTagResourcesRequest {
	s.Tags = v
	return s
}

func (s *ListTagResourcesRequest) SetNextToken(v string) *ListTagResourcesRequest {
	s.NextToken = &v
	return s
}

type ListTagResourcesShrinkRequest struct {
	// 集群ID列表。
	ResourceIdsShrink *string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty"`
	// 资源类型，只支持Cluster
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 地域ID
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 按标签查找。
	TagsShrink *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// 下一次查询Token。
	NextToken *string `json:"next_token,omitempty" xml:"next_token,omitempty"`
}

func (s ListTagResourcesShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesShrinkRequest) GoString() string {
	return s.String()
}

func (s *ListTagResourcesShrinkRequest) SetResourceIdsShrink(v string) *ListTagResourcesShrinkRequest {
	s.ResourceIdsShrink = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetResourceType(v string) *ListTagResourcesShrinkRequest {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetRegionId(v string) *ListTagResourcesShrinkRequest {
	s.RegionId = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetTagsShrink(v string) *ListTagResourcesShrinkRequest {
	s.TagsShrink = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetNextToken(v string) *ListTagResourcesShrinkRequest {
	s.NextToken = &v
	return s
}

type ListTagResourcesResponseBody struct {
	// 下一个查询开始Token，为空说明没有下一个
	NextToken *string `json:"next_token,omitempty" xml:"next_token,omitempty"`
	// 请求ID。
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 资源标签列表。
	TagResources *ListTagResourcesResponseBodyTagResources `json:"tag_resources,omitempty" xml:"tag_resources,omitempty" type:"Struct"`
}

func (s ListTagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBody) SetNextToken(v string) *ListTagResourcesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetRequestId(v string) *ListTagResourcesResponseBody {
	s.RequestId = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetTagResources(v *ListTagResourcesResponseBodyTagResources) *ListTagResourcesResponseBody {
	s.TagResources = v
	return s
}

type ListTagResourcesResponseBodyTagResources struct {
	// 资源标签。
	TagResource []*ListTagResourcesResponseBodyTagResourcesTagResource `json:"tag_resource,omitempty" xml:"tag_resource,omitempty" type:"Repeated"`
}

func (s ListTagResourcesResponseBodyTagResources) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResources) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResources) SetTagResource(v []*ListTagResourcesResponseBodyTagResourcesTagResource) *ListTagResourcesResponseBodyTagResources {
	s.TagResource = v
	return s
}

type ListTagResourcesResponseBodyTagResourcesTagResource struct {
	// 标签key。
	TagKey *string `json:"tag_key,omitempty" xml:"tag_key,omitempty"`
	// 标签值。
	TagValue *string `json:"tag_value,omitempty" xml:"tag_value,omitempty"`
	// 资源ID。
	ResourceId *string `json:"resource_id,omitempty" xml:"resource_id,omitempty"`
	// 资源类型。
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagKey(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagKey = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagValue(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagValue = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceId(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceId = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceType(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceType = &v
	return s
}

type ListTagResourcesResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ListTagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponse) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponse) SetHeaders(v map[string]*string) *ListTagResourcesResponse {
	s.Headers = v
	return s
}

func (s *ListTagResourcesResponse) SetBody(v *ListTagResourcesResponseBody) *ListTagResourcesResponse {
	s.Body = v
	return s
}

type UntagResourcesRequest struct {
	// 资源所属的地域ID
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 资源ID。数组长度取值范围为：1~50
	ResourceIds []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	// 资源类型定义。取值范围： 只支持CLUSTER这一种资源类型
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 资源的标签键。N的取值范围：1~20
	TagKeys []*string `json:"tag_keys,omitempty" xml:"tag_keys,omitempty" type:"Repeated"`
}

func (s UntagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesRequest) GoString() string {
	return s.String()
}

func (s *UntagResourcesRequest) SetRegionId(v string) *UntagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *UntagResourcesRequest) SetResourceIds(v []*string) *UntagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *UntagResourcesRequest) SetResourceType(v string) *UntagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *UntagResourcesRequest) SetTagKeys(v []*string) *UntagResourcesRequest {
	s.TagKeys = v
	return s
}

type UntagResourcesResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UntagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesResponse) GoString() string {
	return s.String()
}

func (s *UntagResourcesResponse) SetHeaders(v map[string]*string) *UntagResourcesResponse {
	s.Headers = v
	return s
}

type ModifyClusterRequest struct {
	// 集群是否绑定EIP，用于公网访问API Server。 true | false
	ApiServerEip *bool `json:"api_server_eip,omitempty" xml:"api_server_eip,omitempty"`
	// 集群API Server 公网连接端点。
	ApiServerEipId *string `json:"api_server_eip_id,omitempty" xml:"api_server_eip_id,omitempty"`
	// 集群是否开启删除保护。默认值false。
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// 实例删除保护，防止通过控制台或API误删除释放节点。
	InstanceDeletionProtection *bool `json:"instance_deletion_protection,omitempty" xml:"instance_deletion_protection,omitempty"`
	// 域名是否重新绑定到Ingress的SLB地址。
	IngressDomainRebinding *string `json:"ingress_domain_rebinding,omitempty" xml:"ingress_domain_rebinding,omitempty"`
	// 集群的Ingress SLB的ID。
	IngressLoadbalancerId *string `json:"ingress_loadbalancer_id,omitempty" xml:"ingress_loadbalancer_id,omitempty"`
	// 集群资源组ID。
	ResourceGroupId   *string            `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
}

func (s ModifyClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterRequest) SetApiServerEip(v bool) *ModifyClusterRequest {
	s.ApiServerEip = &v
	return s
}

func (s *ModifyClusterRequest) SetApiServerEipId(v string) *ModifyClusterRequest {
	s.ApiServerEipId = &v
	return s
}

func (s *ModifyClusterRequest) SetDeletionProtection(v bool) *ModifyClusterRequest {
	s.DeletionProtection = &v
	return s
}

func (s *ModifyClusterRequest) SetInstanceDeletionProtection(v bool) *ModifyClusterRequest {
	s.InstanceDeletionProtection = &v
	return s
}

func (s *ModifyClusterRequest) SetIngressDomainRebinding(v string) *ModifyClusterRequest {
	s.IngressDomainRebinding = &v
	return s
}

func (s *ModifyClusterRequest) SetIngressLoadbalancerId(v string) *ModifyClusterRequest {
	s.IngressLoadbalancerId = &v
	return s
}

func (s *ModifyClusterRequest) SetResourceGroupId(v string) *ModifyClusterRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *ModifyClusterRequest) SetMaintenanceWindow(v *MaintenanceWindow) *ModifyClusterRequest {
	s.MaintenanceWindow = v
	return s
}

type ModifyClusterResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 请求ID。
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ModifyClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyClusterResponseBody) SetClusterId(v string) *ModifyClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ModifyClusterResponseBody) SetRequestId(v string) *ModifyClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ModifyClusterResponseBody) SetTaskId(v string) *ModifyClusterResponseBody {
	s.TaskId = &v
	return s
}

type ModifyClusterResponse struct {
	Headers map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ModifyClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterResponse) SetHeaders(v map[string]*string) *ModifyClusterResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterResponse) SetBody(v *ModifyClusterResponseBody) *ModifyClusterResponse {
	s.Body = v
	return s
}

type DescribeClusterAttachScriptsRequest struct {
	// 节点池ID。将节点加入指定节点池。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 数据盘挂载
	FormatDisk *bool `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	// 保留实例名称
	KeepInstanceName *bool `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	// RDS白名单
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 节点CPU架构,支持amd64、arm、arm64。边缘托管集群专有字段。
	Arch *string `json:"arch,omitempty" xml:"arch,omitempty"`
	// 边缘托管版集群节点的接入配置。
	Options *string `json:"options,omitempty" xml:"options,omitempty"`
}

func (s DescribeClusterAttachScriptsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAttachScriptsRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAttachScriptsRequest) SetNodepoolId(v string) *DescribeClusterAttachScriptsRequest {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetFormatDisk(v bool) *DescribeClusterAttachScriptsRequest {
	s.FormatDisk = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetKeepInstanceName(v bool) *DescribeClusterAttachScriptsRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetRdsInstances(v []*string) *DescribeClusterAttachScriptsRequest {
	s.RdsInstances = v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetArch(v string) *DescribeClusterAttachScriptsRequest {
	s.Arch = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetOptions(v string) *DescribeClusterAttachScriptsRequest {
	s.Options = &v
	return s
}

type DescribeClusterAttachScriptsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *string            `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAttachScriptsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAttachScriptsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAttachScriptsResponse) SetHeaders(v map[string]*string) *DescribeClusterAttachScriptsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAttachScriptsResponse) SetBody(v string) *DescribeClusterAttachScriptsResponse {
	s.Body = &v
	return s
}

type RemoveClusterNodesRequest struct {
	// 是否排空节点上的Pod。
	DrainNode *bool `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	// 要移除的Node列表。
	Nodes []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	// 是否同时释放ECS。
	ReleaseNode *bool `json:"release_node,omitempty" xml:"release_node,omitempty"`
}

func (s RemoveClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s RemoveClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *RemoveClusterNodesRequest) SetDrainNode(v bool) *RemoveClusterNodesRequest {
	s.DrainNode = &v
	return s
}

func (s *RemoveClusterNodesRequest) SetNodes(v []*string) *RemoveClusterNodesRequest {
	s.Nodes = v
	return s
}

func (s *RemoveClusterNodesRequest) SetReleaseNode(v bool) *RemoveClusterNodesRequest {
	s.ReleaseNode = &v
	return s
}

type RemoveClusterNodesResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s RemoveClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *RemoveClusterNodesResponse) SetHeaders(v map[string]*string) *RemoveClusterNodesResponse {
	s.Headers = v
	return s
}

type DescribeKubernetesVersionMetadataRequest struct {
	// 地域ID。
	Region *string `json:"Region,omitempty" xml:"Region,omitempty"`
	// 集群类型。
	ClusterType *string `json:"ClusterType,omitempty" xml:"ClusterType,omitempty"`
	// 要查询的版本，如果为空则查所有版本。
	KubernetesVersion *string `json:"KubernetesVersion,omitempty" xml:"KubernetesVersion,omitempty"`
	// 边缘集群标识，用于区分边缘集群，取值：Default或Edge。
	Profile *string `json:"Profile,omitempty" xml:"Profile,omitempty"`
}

func (s DescribeKubernetesVersionMetadataRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataRequest) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataRequest) SetRegion(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Region = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetClusterType(v string) *DescribeKubernetesVersionMetadataRequest {
	s.ClusterType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetKubernetesVersion(v string) *DescribeKubernetesVersionMetadataRequest {
	s.KubernetesVersion = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetProfile(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Profile = &v
	return s
}

type DescribeKubernetesVersionMetadataResponse struct {
	Headers map[string]*string                               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeKubernetesVersionMetadataResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeKubernetesVersionMetadataResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponse) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponse) SetHeaders(v map[string]*string) *DescribeKubernetesVersionMetadataResponse {
	s.Headers = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponse) SetBody(v []*DescribeKubernetesVersionMetadataResponseBody) *DescribeKubernetesVersionMetadataResponse {
	s.Body = v
	return s
}

type DescribeKubernetesVersionMetadataResponseBody struct {
	// Kubernetes版本特性。
	Capabilities map[string]interface{} `json:"capabilities,omitempty" xml:"capabilities,omitempty"`
	// ECS系统镜像列表。
	Images []*DescribeKubernetesVersionMetadataResponseBodyImages `json:"images,omitempty" xml:"images,omitempty" type:"Repeated"`
	// Kubernetes版本元数据信息。
	MetaData map[string]interface{} `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// 容器运行时详情。
	Runtimes []*Runtime `json:"runtimes,omitempty" xml:"runtimes,omitempty" type:"Repeated"`
	// Kubernetes版本。
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
	// 是否为多可用区。
	MultiAz *string `json:"multi_az,omitempty" xml:"multi_az,omitempty"`
}

func (s DescribeKubernetesVersionMetadataResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetCapabilities(v map[string]interface{}) *DescribeKubernetesVersionMetadataResponseBody {
	s.Capabilities = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetImages(v []*DescribeKubernetesVersionMetadataResponseBodyImages) *DescribeKubernetesVersionMetadataResponseBody {
	s.Images = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetMetaData(v map[string]interface{}) *DescribeKubernetesVersionMetadataResponseBody {
	s.MetaData = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetRuntimes(v []*Runtime) *DescribeKubernetesVersionMetadataResponseBody {
	s.Runtimes = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetVersion(v string) *DescribeKubernetesVersionMetadataResponseBody {
	s.Version = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetMultiAz(v string) *DescribeKubernetesVersionMetadataResponseBody {
	s.MultiAz = &v
	return s
}

type DescribeKubernetesVersionMetadataResponseBodyImages struct {
	// 镜像ID。
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 镜像名称。
	ImageName *string `json:"image_name,omitempty" xml:"image_name,omitempty"`
	// 操作系统发行版。取值范围： CentOS,AliyunLinux,Windows,WindowsCore。
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// 镜像版本。
	OsVersion *string `json:"os_version,omitempty" xml:"os_version,omitempty"`
	// 镜像类型。
	ImageType *string `json:"image_type,omitempty" xml:"image_type,omitempty"`
	// 操作系统发行版本号。
	OsType *string `json:"os_type,omitempty" xml:"os_type,omitempty"`
	// 镜像分类
	ImageCategory *string `json:"image_category,omitempty" xml:"image_category,omitempty"`
}

func (s DescribeKubernetesVersionMetadataResponseBodyImages) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponseBodyImages) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageId(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageId = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageName(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageName = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetPlatform(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.Platform = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetOsVersion(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.OsVersion = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageType(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetOsType(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.OsType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageCategory(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageCategory = &v
	return s
}

type DescribeClusterLogsResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeClusterLogsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClusterLogsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterLogsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterLogsResponse) SetHeaders(v map[string]*string) *DescribeClusterLogsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterLogsResponse) SetBody(v []*DescribeClusterLogsResponseBody) *DescribeClusterLogsResponse {
	s.Body = v
	return s
}

type DescribeClusterLogsResponseBody struct {
	// 日志ID。
	ID *int64 `json:"ID,omitempty" xml:"ID,omitempty"`
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 日志内容。
	ClusterLog *string `json:"cluster_log,omitempty" xml:"cluster_log,omitempty"`
	// 日志等级。
	LogLevel *string `json:"log_level,omitempty" xml:"log_level,omitempty"`
	// 日志创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 日志更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterLogsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterLogsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterLogsResponseBody) SetID(v int64) *DescribeClusterLogsResponseBody {
	s.ID = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetClusterId(v string) *DescribeClusterLogsResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetClusterLog(v string) *DescribeClusterLogsResponseBody {
	s.ClusterLog = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetLogLevel(v string) *DescribeClusterLogsResponseBody {
	s.LogLevel = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetCreated(v string) *DescribeClusterLogsResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetUpdated(v string) *DescribeClusterLogsResponseBody {
	s.Updated = &v
	return s
}

type CreateKubernetesTriggerRequest struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 项目名称。
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	// 触发器行为
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// 触发器类型。默认deployment。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateKubernetesTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerRequest) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerRequest) SetClusterId(v string) *CreateKubernetesTriggerRequest {
	s.ClusterId = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetProjectId(v string) *CreateKubernetesTriggerRequest {
	s.ProjectId = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetAction(v string) *CreateKubernetesTriggerRequest {
	s.Action = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetType(v string) *CreateKubernetesTriggerRequest {
	s.Type = &v
	return s
}

type CreateKubernetesTriggerResponseBody struct {
	// 触发器ID。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 触发器项目名称。
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	// 触发器类型。默认值为 deployment 。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 触发器行为。
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s CreateKubernetesTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerResponseBody) SetId(v string) *CreateKubernetesTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetClusterId(v string) *CreateKubernetesTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetProjectId(v string) *CreateKubernetesTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetType(v string) *CreateKubernetesTriggerResponseBody {
	s.Type = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetAction(v string) *CreateKubernetesTriggerResponseBody {
	s.Action = &v
	return s
}

type CreateKubernetesTriggerResponse struct {
	Headers map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CreateKubernetesTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerResponse) SetHeaders(v map[string]*string) *CreateKubernetesTriggerResponse {
	s.Headers = v
	return s
}

func (s *CreateKubernetesTriggerResponse) SetBody(v *CreateKubernetesTriggerResponseBody) *CreateKubernetesTriggerResponse {
	s.Body = v
	return s
}

type GrantPermissionsRequest struct {
	// 请求体参数
	Body []*GrantPermissionsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s GrantPermissionsRequest) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsRequest) GoString() string {
	return s.String()
}

func (s *GrantPermissionsRequest) SetBody(v []*GrantPermissionsRequestBody) *GrantPermissionsRequest {
	s.Body = v
	return s
}

type GrantPermissionsRequestBody struct {
	// 授权目标集群id
	Cluster *string `json:"cluster,omitempty" xml:"cluster,omitempty"`
	// 该授权是否是自定义授权
	IsCustom *bool `json:"is_custom,omitempty" xml:"is_custom,omitempty"`
	// 预置的角色名称
	RoleName *string `json:"role_name,omitempty" xml:"role_name,omitempty"`
	// 授权类型
	RoleType *string `json:"role_type,omitempty" xml:"role_type,omitempty"`
	// 命名空间名称
	Namespace *string `json:"namespace,omitempty" xml:"namespace,omitempty"`
	// 是否是 RAM 角色授权
	IsRamRole *bool `json:"is_ram_role,omitempty" xml:"is_ram_role,omitempty"`
}

func (s GrantPermissionsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsRequestBody) GoString() string {
	return s.String()
}

func (s *GrantPermissionsRequestBody) SetCluster(v string) *GrantPermissionsRequestBody {
	s.Cluster = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetIsCustom(v bool) *GrantPermissionsRequestBody {
	s.IsCustom = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetRoleName(v string) *GrantPermissionsRequestBody {
	s.RoleName = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetRoleType(v string) *GrantPermissionsRequestBody {
	s.RoleType = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetNamespace(v string) *GrantPermissionsRequestBody {
	s.Namespace = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetIsRamRole(v bool) *GrantPermissionsRequestBody {
	s.IsRamRole = &v
	return s
}

type GrantPermissionsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s GrantPermissionsResponse) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsResponse) GoString() string {
	return s.String()
}

func (s *GrantPermissionsResponse) SetHeaders(v map[string]*string) *GrantPermissionsResponse {
	s.Headers = v
	return s
}

type DescribeClusterDetailResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 集群类型。
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// 集群创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 集群初始化版本。
	InitVersion *string `json:"init_version,omitempty" xml:"init_version,omitempty"`
	// 集群当前版本。
	CurrentVersion *string `json:"current_version,omitempty" xml:"current_version,omitempty"`
	// 集群可升级版本。
	NextVersion *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// 集群是否开启删除保护。
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// 集群内Docker版本。
	DockerVersion *string `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	// 集群Ingress LB实例ID。
	ExternalLoadbalancerId *string `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	// 集群元数据。
	MetaData *string `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// 集群名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群采用的网络类型，例如：VPC网络。
	NetworkMode *string `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	// 集群所在地域ID。
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 集群资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// 集群安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 集群节点数量。
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// 集群运行状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 集群标签。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 集群更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// 集群使用的VPC ID。
	VpcId *string `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	// 集群节点使用的虚拟交换机列表。
	VswitchId *string `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	// Pod网络地址段，必须是有效的私有网段，即以下网段及其子网：10.0.0.0/8，172.16-31.0.0/12-16，192.168.0.0/16。不能与 VPC 及VPC 内已有 Kubernetes 集群使用的网段重复，创建成功后不能修改。  有关集群网络规划，请参见：[VPC下 Kubernetes 的网络地址段规划](https://help.aliyun.com/document_detail/～～86500～～)。
	SubnetCidr *string `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	// 集群所在地域内的可用区ID。
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
	// 集群访问地址。
	MasterUrl *string `json:"master_url,omitempty" xml:"master_url,omitempty"`
	// 集群是否启用用PrivateZone。  true：启用 false：不启用 默认值：false。
	PrivateZone *bool `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	// 面向场景时的集群类型。  Default：非边缘场景集群。 Edge：边缘场景集群。
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// 托管版集群类型，面向托管集群。  ack.pro.small：专业托管集群。 ack.standard ：标准托管集群。
	ClusterSpec *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	// Worker节点RAM角色名称。
	WorkerRamRoleName *string            `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
}

func (s DescribeClusterDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterDetailResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterDetailResponseBody) SetClusterId(v string) *DescribeClusterDetailResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetClusterType(v string) *DescribeClusterDetailResponseBody {
	s.ClusterType = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetCreated(v string) *DescribeClusterDetailResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetInitVersion(v string) *DescribeClusterDetailResponseBody {
	s.InitVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetCurrentVersion(v string) *DescribeClusterDetailResponseBody {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetNextVersion(v string) *DescribeClusterDetailResponseBody {
	s.NextVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetDeletionProtection(v bool) *DescribeClusterDetailResponseBody {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetDockerVersion(v string) *DescribeClusterDetailResponseBody {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetExternalLoadbalancerId(v string) *DescribeClusterDetailResponseBody {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMetaData(v string) *DescribeClusterDetailResponseBody {
	s.MetaData = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetName(v string) *DescribeClusterDetailResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetNetworkMode(v string) *DescribeClusterDetailResponseBody {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetRegionId(v string) *DescribeClusterDetailResponseBody {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetResourceGroupId(v string) *DescribeClusterDetailResponseBody {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSecurityGroupId(v string) *DescribeClusterDetailResponseBody {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSize(v int64) *DescribeClusterDetailResponseBody {
	s.Size = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetState(v string) *DescribeClusterDetailResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetTags(v []*Tag) *DescribeClusterDetailResponseBody {
	s.Tags = v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetUpdated(v string) *DescribeClusterDetailResponseBody {
	s.Updated = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetVpcId(v string) *DescribeClusterDetailResponseBody {
	s.VpcId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetVswitchId(v string) *DescribeClusterDetailResponseBody {
	s.VswitchId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSubnetCidr(v string) *DescribeClusterDetailResponseBody {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetZoneId(v string) *DescribeClusterDetailResponseBody {
	s.ZoneId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMasterUrl(v string) *DescribeClusterDetailResponseBody {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetPrivateZone(v bool) *DescribeClusterDetailResponseBody {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetProfile(v string) *DescribeClusterDetailResponseBody {
	s.Profile = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetClusterSpec(v string) *DescribeClusterDetailResponseBody {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetWorkerRamRoleName(v string) *DescribeClusterDetailResponseBody {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMaintenanceWindow(v *MaintenanceWindow) *DescribeClusterDetailResponseBody {
	s.MaintenanceWindow = v
	return s
}

type DescribeClusterDetailResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterDetailResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterDetailResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterDetailResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterDetailResponse) SetHeaders(v map[string]*string) *DescribeClusterDetailResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterDetailResponse) SetBody(v *DescribeClusterDetailResponseBody) *DescribeClusterDetailResponse {
	s.Body = v
	return s
}

type PauseComponentUpgradeResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s PauseComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s PauseComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *PauseComponentUpgradeResponse) SetHeaders(v map[string]*string) *PauseComponentUpgradeResponse {
	s.Headers = v
	return s
}

type DescribeClustersRequest struct {
	// 集群名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群类型。
	ClusterType *string `json:"clusterType,omitempty" xml:"clusterType,omitempty"`
}

func (s DescribeClustersRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersRequest) GoString() string {
	return s.String()
}

func (s *DescribeClustersRequest) SetName(v string) *DescribeClustersRequest {
	s.Name = &v
	return s
}

func (s *DescribeClustersRequest) SetClusterType(v string) *DescribeClustersRequest {
	s.ClusterType = &v
	return s
}

type DescribeClustersResponse struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeClustersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClustersResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponse) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponse) SetHeaders(v map[string]*string) *DescribeClustersResponse {
	s.Headers = v
	return s
}

func (s *DescribeClustersResponse) SetBody(v []*DescribeClustersResponseBody) *DescribeClustersResponse {
	s.Body = v
	return s
}

type DescribeClustersResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 集群类型。
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// 集群创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 集群当前版本。
	CurrentVersion *string `json:"current_version,omitempty" xml:"current_version,omitempty"`
	// 节点系统盘类型。
	DataDiskCategory *string `json:"data_disk_category,omitempty" xml:"data_disk_category,omitempty"`
	// 节点系统盘大小。
	DataDiskSize *int64 `json:"data_disk_size,omitempty" xml:"data_disk_size,omitempty"`
	// 集群是否开启删除保护。
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// 容器运行时版本。
	DockerVersion *string `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	// 集群Ingerss SLB实例的ID。
	ExternalLoadbalancerId *string `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	// 集群创建时版本。
	InitVersion *string `json:"init_version,omitempty" xml:"init_version,omitempty"`
	// 集群的endpoint地址。
	MasterUrl *string `json:"master_url,omitempty" xml:"master_url,omitempty"`
	// 集群元数据。
	MetaData *string `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// 集群名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群使用的网络类型。
	NetworkMode *string `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	// 集群是否开启Private Zone，默认false。
	PrivateZone *bool `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	// 集群标识，区分是否为边缘托管版。
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// 集群所在地域ID。
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 集群资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// 集群安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 集群内实例数量。
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// 集群运行状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// POD网络。
	SubnetCidr *string `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	// 集群标签。
	Tags []*DescribeClustersResponseBodyTags `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 集群更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// 集群使用的VPC ID。
	VpcId *string `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	// 虚拟交换机网络ID。
	VswitchCidr *string `json:"vswitch_cidr,omitempty" xml:"vswitch_cidr,omitempty"`
	// 节点使用的Vswitch ID。
	VswitchId *string `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	// 集群Worker节点RAM角色名称。
	WorkerRamRoleName *string `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	// 集群所在Region内的区域ID。
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
}

func (s DescribeClustersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponseBody) SetClusterId(v string) *DescribeClustersResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetClusterType(v string) *DescribeClustersResponseBody {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersResponseBody) SetCreated(v string) *DescribeClustersResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClustersResponseBody) SetCurrentVersion(v string) *DescribeClustersResponseBody {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDataDiskCategory(v string) *DescribeClustersResponseBody {
	s.DataDiskCategory = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDataDiskSize(v int64) *DescribeClustersResponseBody {
	s.DataDiskSize = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDeletionProtection(v bool) *DescribeClustersResponseBody {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDockerVersion(v string) *DescribeClustersResponseBody {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetExternalLoadbalancerId(v string) *DescribeClustersResponseBody {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetInitVersion(v string) *DescribeClustersResponseBody {
	s.InitVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetMasterUrl(v string) *DescribeClustersResponseBody {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClustersResponseBody) SetMetaData(v string) *DescribeClustersResponseBody {
	s.MetaData = &v
	return s
}

func (s *DescribeClustersResponseBody) SetName(v string) *DescribeClustersResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClustersResponseBody) SetNetworkMode(v string) *DescribeClustersResponseBody {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClustersResponseBody) SetPrivateZone(v bool) *DescribeClustersResponseBody {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClustersResponseBody) SetProfile(v string) *DescribeClustersResponseBody {
	s.Profile = &v
	return s
}

func (s *DescribeClustersResponseBody) SetRegionId(v string) *DescribeClustersResponseBody {
	s.RegionId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetResourceGroupId(v string) *DescribeClustersResponseBody {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSecurityGroupId(v string) *DescribeClustersResponseBody {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSize(v int64) *DescribeClustersResponseBody {
	s.Size = &v
	return s
}

func (s *DescribeClustersResponseBody) SetState(v string) *DescribeClustersResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSubnetCidr(v string) *DescribeClustersResponseBody {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClustersResponseBody) SetTags(v []*DescribeClustersResponseBodyTags) *DescribeClustersResponseBody {
	s.Tags = v
	return s
}

func (s *DescribeClustersResponseBody) SetUpdated(v string) *DescribeClustersResponseBody {
	s.Updated = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVpcId(v string) *DescribeClustersResponseBody {
	s.VpcId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVswitchCidr(v string) *DescribeClustersResponseBody {
	s.VswitchCidr = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVswitchId(v string) *DescribeClustersResponseBody {
	s.VswitchId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetWorkerRamRoleName(v string) *DescribeClustersResponseBody {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClustersResponseBody) SetZoneId(v string) *DescribeClustersResponseBody {
	s.ZoneId = &v
	return s
}

type DescribeClustersResponseBodyTags struct {
	// 标签名。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// 标签值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s DescribeClustersResponseBodyTags) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponseBodyTags) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponseBodyTags) SetKey(v string) *DescribeClustersResponseBodyTags {
	s.Key = &v
	return s
}

func (s *DescribeClustersResponseBodyTags) SetValue(v string) *DescribeClustersResponseBodyTags {
	s.Value = &v
	return s
}

type DescribeUserPermissionResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeUserPermissionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeUserPermissionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserPermissionResponse) GoString() string {
	return s.String()
}

func (s *DescribeUserPermissionResponse) SetHeaders(v map[string]*string) *DescribeUserPermissionResponse {
	s.Headers = v
	return s
}

func (s *DescribeUserPermissionResponse) SetBody(v []*DescribeUserPermissionResponseBody) *DescribeUserPermissionResponse {
	s.Body = v
	return s
}

type DescribeUserPermissionResponseBody struct {
	// 集群访问配置
	ResourceId *string `json:"resource_id,omitempty" xml:"resource_id,omitempty"`
	// 授权类型
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 自定义角色名称
	RoleName *string `json:"role_name,omitempty" xml:"role_name,omitempty"`
	// 预置的角色类型
	RoleType *string `json:"role_type,omitempty" xml:"role_type,omitempty"`
	// 是否为集群 owner 的授权
	IsOwner *int64 `json:"is_owner,omitempty" xml:"is_owner,omitempty"`
	// 是否为ram 角色授权
	IsRamRole *int64 `json:"is_ram_role,omitempty" xml:"is_ram_role,omitempty"`
}

func (s DescribeUserPermissionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserPermissionResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeUserPermissionResponseBody) SetResourceId(v string) *DescribeUserPermissionResponseBody {
	s.ResourceId = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetResourceType(v string) *DescribeUserPermissionResponseBody {
	s.ResourceType = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetRoleName(v string) *DescribeUserPermissionResponseBody {
	s.RoleName = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetRoleType(v string) *DescribeUserPermissionResponseBody {
	s.RoleType = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetIsOwner(v int64) *DescribeUserPermissionResponseBody {
	s.IsOwner = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetIsRamRole(v int64) *DescribeUserPermissionResponseBody {
	s.IsRamRole = &v
	return s
}

type ModifyClusterNodePoolRequest struct {
	// 自动伸缩节点池配置。
	AutoScaling *ModifyClusterNodePoolRequestAutoScaling `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	// 集群配置。
	KubernetesConfig *ModifyClusterNodePoolRequestKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// 节点池配置。
	NodepoolInfo *ModifyClusterNodePoolRequestNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// 扩容组配置。
	ScalingGroup *ModifyClusterNodePoolRequestScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	// 加密计算配置。
	TeeConfig *ModifyClusterNodePoolRequestTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
	// 托管版节点池配置。
	Management *ModifyClusterNodePoolRequestManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	// 是否同步更新节点标签及污点。
	UpdateNodes *bool `json:"update_nodes,omitempty" xml:"update_nodes,omitempty"`
}

func (s ModifyClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequest) SetAutoScaling(v *ModifyClusterNodePoolRequestAutoScaling) *ModifyClusterNodePoolRequest {
	s.AutoScaling = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetKubernetesConfig(v *ModifyClusterNodePoolRequestKubernetesConfig) *ModifyClusterNodePoolRequest {
	s.KubernetesConfig = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetNodepoolInfo(v *ModifyClusterNodePoolRequestNodepoolInfo) *ModifyClusterNodePoolRequest {
	s.NodepoolInfo = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetScalingGroup(v *ModifyClusterNodePoolRequestScalingGroup) *ModifyClusterNodePoolRequest {
	s.ScalingGroup = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetTeeConfig(v *ModifyClusterNodePoolRequestTeeConfig) *ModifyClusterNodePoolRequest {
	s.TeeConfig = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetManagement(v *ModifyClusterNodePoolRequestManagement) *ModifyClusterNodePoolRequest {
	s.Management = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetUpdateNodes(v bool) *ModifyClusterNodePoolRequest {
	s.UpdateNodes = &v
	return s
}

type ModifyClusterNodePoolRequestAutoScaling struct {
	// 带宽峰值。
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// EIP计费类型。
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// 是否开启自动伸缩。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 是否绑定EIP。
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// 最大实例数。
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// 最小实例数。
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// 自动伸缩节点类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s ModifyClusterNodePoolRequestAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestAutoScaling) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEipBandwidth(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEipInternetChargeType(v string) *ModifyClusterNodePoolRequestAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEnable(v bool) *ModifyClusterNodePoolRequestAutoScaling {
	s.Enable = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetIsBondEip(v bool) *ModifyClusterNodePoolRequestAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetMaxInstances(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetMinInstances(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetType(v string) *ModifyClusterNodePoolRequestAutoScaling {
	s.Type = &v
	return s
}

type ModifyClusterNodePoolRequestKubernetesConfig struct {
	// 是否开启云监控。
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// CPU管理策略。
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 节点标签。
	Labels []*Tag `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	// 容器运行时。
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 容器运行时版本。
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// 污点配置。
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 实例自定义数据。
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s ModifyClusterNodePoolRequestKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestKubernetesConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetCmsEnabled(v bool) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetCpuPolicy(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetLabels(v []*Tag) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Labels = v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetRuntime(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetRuntimeVersion(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetTaints(v []*Taint) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Taints = v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetUserData(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.UserData = &v
	return s
}

type ModifyClusterNodePoolRequestNodepoolInfo struct {
	// 节点池名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
}

func (s ModifyClusterNodePoolRequestNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestNodepoolInfo) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestNodepoolInfo) SetName(v string) *ModifyClusterNodePoolRequestNodepoolInfo {
	s.Name = &v
	return s
}

func (s *ModifyClusterNodePoolRequestNodepoolInfo) SetResourceGroupId(v string) *ModifyClusterNodePoolRequestNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

type ModifyClusterNodePoolRequestScalingGroup struct {
	// 数据盘配置。
	DataDisks []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	// 节点付费类型。
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// 包年包月时长
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// 付费周期
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// 节点池节点是启用自动续费
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// 节点池节点自动续费周期
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// 操作系统发行版。
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// 自定义镜像
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 抢占式实例类型
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// 抢占实例价格上限配置
	SpotPriceLimit []*ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// 节点实例规格。
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// 密钥对名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// SSH登录密码，和key_pari二选一。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// RDS实例列表。
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 扩容策略。
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// 节点系统盘类型。
	SystemDiskCategory *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	// 节点系统盘大小。
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// ECS标签。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 节点使用的虚拟交换机ID。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 多可用区伸缩组ECS实例扩缩容策略
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// 伸缩组所需要按量实例个数的最小值，取值范围：0~1000。当按量实例个数少于该值时，将优先创建按量实例。
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// 伸缩组满足最小按量实例数（OnDemandBaseCapacity）要求后，超出的实例中按量实例应占的比例，取值范围：0～100。
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// 指定可用实例规格的个数，伸缩组将按成本最低的多个规格均衡创建抢占式实例。取值范围：1~10。
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// 是否开启补齐抢占式实例。开启后，当收到抢占式实例将被回收的系统消息时，伸缩组将尝试创建新的实例，替换掉将被回收的抢占式实例。
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// 当MultiAZPolicy取值为COST_OPTIMIZED时，如果因价格、库存等原因无法创建足够的抢占式实例，是否允许自动尝试创建按量实例满足ECS实例数量要求。取值范围：true：允许。false：不允许。默认值：true
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// 节点公网IP网络计费类型
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// 节点公网IP出带宽最大值，单位为Mbps（Mega bit per second），取值范围：1~100
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
}

func (s ModifyClusterNodePoolRequestScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestScalingGroup) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetDataDisks(v []*DataDisk) *ModifyClusterNodePoolRequestScalingGroup {
	s.DataDisks = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInstanceChargeType(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPeriod(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.Period = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPeriodUnit(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetAutoRenew(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetAutoRenewPeriod(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPlatform(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.Platform = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetImageId(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.ImageId = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotStrategy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotPriceLimit(v []*ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInstanceTypes(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetKeyPair(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetLoginPassword(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetRdsInstances(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetScalingPolicy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSystemDiskCategory(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSystemDiskSize(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetTags(v []*Tag) *ModifyClusterNodePoolRequestScalingGroup {
	s.Tags = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetVswitchIds(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.VswitchIds = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetMultiAzPolicy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetOnDemandBaseCapacity(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotInstancePools(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotInstanceRemedy(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetCompensateWithOnDemand(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInternetChargeType(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInternetMaxBandwidthOut(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

type ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit struct {
	// 抢占式实例规格
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// 单台实例上限价格，单位：元/小时。
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) SetInstanceType(v string) *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) SetPriceLimit(v string) *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type ModifyClusterNodePoolRequestTeeConfig struct {
	// 是否为加密计算节点池。
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s ModifyClusterNodePoolRequestTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestTeeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestTeeConfig) SetTeeEnable(v bool) *ModifyClusterNodePoolRequestTeeConfig {
	s.TeeEnable = &v
	return s
}

type ModifyClusterNodePoolRequestManagement struct {
	// 是否启用托管节点池。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 是否开启自动修复。
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// 自动升级配置。
	UpgradeConfig *ModifyClusterNodePoolRequestManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s ModifyClusterNodePoolRequestManagement) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestManagement) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestManagement) SetEnable(v bool) *ModifyClusterNodePoolRequestManagement {
	s.Enable = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagement) SetAutoRepair(v bool) *ModifyClusterNodePoolRequestManagement {
	s.AutoRepair = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagement) SetUpgradeConfig(v *ModifyClusterNodePoolRequestManagementUpgradeConfig) *ModifyClusterNodePoolRequestManagement {
	s.UpgradeConfig = v
	return s
}

type ModifyClusterNodePoolRequestManagementUpgradeConfig struct {
	// 是否启用自动升级，自修复。
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// 额外节点数量。
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// 额外节点比例， 和surge 二选一。
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
	// 最大不可用节点数量。
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
}

func (s ModifyClusterNodePoolRequestManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetAutoUpgrade(v bool) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetSurge(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetSurgePercentage(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetMaxUnavailable(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

type ModifyClusterNodePoolResponseBody struct {
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
	// 节点池ID。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
}

func (s ModifyClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolResponseBody) SetTaskId(v string) *ModifyClusterNodePoolResponseBody {
	s.TaskId = &v
	return s
}

func (s *ModifyClusterNodePoolResponseBody) SetNodepoolId(v string) *ModifyClusterNodePoolResponseBody {
	s.NodepoolId = &v
	return s
}

type ModifyClusterNodePoolResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ModifyClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolResponse) SetHeaders(v map[string]*string) *ModifyClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterNodePoolResponse) SetBody(v *ModifyClusterNodePoolResponseBody) *ModifyClusterNodePoolResponse {
	s.Body = v
	return s
}

type ResumeUpgradeClusterResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s ResumeUpgradeClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ResumeUpgradeClusterResponse) GoString() string {
	return s.String()
}

func (s *ResumeUpgradeClusterResponse) SetHeaders(v map[string]*string) *ResumeUpgradeClusterResponse {
	s.Headers = v
	return s
}

type OpenAckServiceRequest struct {
	// 要开通的服务类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s OpenAckServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceRequest) GoString() string {
	return s.String()
}

func (s *OpenAckServiceRequest) SetType(v string) *OpenAckServiceRequest {
	s.Type = &v
	return s
}

type OpenAckServiceResponseBody struct {
	// 请求ID
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 开通服务的订单号。
	OrderId *string `json:"order_id,omitempty" xml:"order_id,omitempty"`
}

func (s OpenAckServiceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceResponseBody) GoString() string {
	return s.String()
}

func (s *OpenAckServiceResponseBody) SetRequestId(v string) *OpenAckServiceResponseBody {
	s.RequestId = &v
	return s
}

func (s *OpenAckServiceResponseBody) SetOrderId(v string) *OpenAckServiceResponseBody {
	s.OrderId = &v
	return s
}

type OpenAckServiceResponse struct {
	Headers map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *OpenAckServiceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s OpenAckServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceResponse) GoString() string {
	return s.String()
}

func (s *OpenAckServiceResponse) SetHeaders(v map[string]*string) *OpenAckServiceResponse {
	s.Headers = v
	return s
}

func (s *OpenAckServiceResponse) SetBody(v *OpenAckServiceResponseBody) *OpenAckServiceResponse {
	s.Body = v
	return s
}

type ScaleClusterNodePoolRequest struct {
	// 扩容节点数量
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
}

func (s ScaleClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolRequest) SetCount(v int64) *ScaleClusterNodePoolRequest {
	s.Count = &v
	return s
}

type ScaleClusterNodePoolResponseBody struct {
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolResponseBody) SetTaskId(v string) *ScaleClusterNodePoolResponseBody {
	s.TaskId = &v
	return s
}

type ScaleClusterNodePoolResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ScaleClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolResponse) SetHeaders(v map[string]*string) *ScaleClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *ScaleClusterNodePoolResponse) SetBody(v *ScaleClusterNodePoolResponseBody) *ScaleClusterNodePoolResponse {
	s.Body = v
	return s
}

type DescribeClusterNodePoolDetailResponseBody struct {
	// 节点池自动伸缩信息。
	AutoScaling *DescribeClusterNodePoolDetailResponseBodyAutoScaling `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	// 节点池所属集群配置。
	KubernetesConfig *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// 节点池详情。
	NodepoolInfo *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// 节点池扩容组信息。
	ScalingGroup *DescribeClusterNodePoolDetailResponseBodyScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	// 节点池状态。
	Status *DescribeClusterNodePoolDetailResponseBodyStatus `json:"status,omitempty" xml:"status,omitempty" type:"Struct"`
	// 加密计算节点池信息。
	TeeConfig *DescribeClusterNodePoolDetailResponseBodyTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
	// 托管版节点池配置。
	Management *DescribeClusterNodePoolDetailResponseBodyManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetAutoScaling(v *DescribeClusterNodePoolDetailResponseBodyAutoScaling) *DescribeClusterNodePoolDetailResponseBody {
	s.AutoScaling = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetKubernetesConfig(v *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) *DescribeClusterNodePoolDetailResponseBody {
	s.KubernetesConfig = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetNodepoolInfo(v *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) *DescribeClusterNodePoolDetailResponseBody {
	s.NodepoolInfo = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetScalingGroup(v *DescribeClusterNodePoolDetailResponseBodyScalingGroup) *DescribeClusterNodePoolDetailResponseBody {
	s.ScalingGroup = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetStatus(v *DescribeClusterNodePoolDetailResponseBodyStatus) *DescribeClusterNodePoolDetailResponseBody {
	s.Status = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetTeeConfig(v *DescribeClusterNodePoolDetailResponseBodyTeeConfig) *DescribeClusterNodePoolDetailResponseBody {
	s.TeeConfig = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetManagement(v *DescribeClusterNodePoolDetailResponseBodyManagement) *DescribeClusterNodePoolDetailResponseBody {
	s.Management = v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyAutoScaling struct {
	// EIP带宽峰值。
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// EIP实例付费类型。
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// 是否启用自动伸缩。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 是否绑定EIP。
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// 最大实例数。
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// 最小实例数。
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// 扩容组类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyAutoScaling) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEipBandwidth(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEipInternetChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetIsBondEip(v bool) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetMaxInstances(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetMinInstances(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetType(v string) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.Type = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyKubernetesConfig struct {
	// 是否开启云监控
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// CPU管理策略
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 节点标签。
	Labels []*Tag `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	// 容器运行时
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 容器运行时版本。
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// 污点配置。
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 节点自定义数据
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetCmsEnabled(v bool) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetCpuPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetLabels(v []*Tag) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Labels = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetRuntime(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetRuntimeVersion(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetTaints(v []*Taint) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Taints = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetUserData(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.UserData = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyNodepoolInfo struct {
	// 节点池创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 是否为默认节点池。
	IsDefault *bool `json:"is_default,omitempty" xml:"is_default,omitempty"`
	// 节点池名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 节点池ID。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 节点池所属地域ID。
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 节点池所属资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// 节点池类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 节点池更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetCreated(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Created = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetIsDefault(v bool) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.IsDefault = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetName(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Name = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetNodepoolId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetRegionId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetResourceGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetType(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Type = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetUpdated(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Updated = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyScalingGroup struct {
	// 节点是否开启自动续费。
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// 节点自动续费周期。
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// 数据盘配置。
	DataDisks []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	// 自定义镜像ID。
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 节点付费类型。
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// 节点ECS规格类型。
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// 多可用区伸缩组ECS实例扩缩容策略
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// 伸缩组所需要按量实例个数的最小值，取值范围：0~1000。当按量实例个数少于该值时，将优先创建按量实例。
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// 伸缩组满足最小按量实例数（OnDemandBaseCapacity）要求后，超出的实例中按量实例应占的比例，取值范围：0～100。
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// 指定可用实例规格的个数，伸缩组将按成本最低的多个规格均衡创建抢占式实例。取值范围：1~10。
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// 是否开启补齐抢占式实例。开启后，当收到抢占式实例将被回收的系统消息时，伸缩组将尝试创建新的实例，替换掉将被回收的抢占式实例。
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// 当MultiAZPolicy取值为COST_OPTIMIZED时，如果因价格、库存等原因无法创建足够的抢占式实例，是否允许自动尝试创建按量实例满足ECS实例数量要求。取值范围：true：允许。false：不允许。默认值：true
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// 节点包年包月时长。
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// 节点付费周期。
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// 操作系统发行版。取值： CentOS，AliyunLinux，Windows，WindowsCore。
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// 节点RAM 角色名称。
	RamPolicy *string `json:"ram_policy,omitempty" xml:"ram_policy,omitempty"`
	// 抢占式实例类型
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// 抢占式实例价格上限配置。
	SpotPriceLimit []*DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// RDS实例列表。
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 扩容组ID。
	ScalingGroupId *string `json:"scaling_group_id,omitempty" xml:"scaling_group_id,omitempty"`
	// 扩容策略。
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// 节点所属安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 系统盘类型
	SystemDiskCategory *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	// 系统盘大小
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// ECS标签
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 虚拟交换机ID。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 登录密码
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 密钥对名称
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// 节点公网IP网络计费类型
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// 节点公网IP出带宽最大值，单位为Mbps（Mega bit per second），取值范围：1~100
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroup) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetAutoRenew(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetAutoRenewPeriod(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetDataDisks(v []*DataDisk) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.DataDisks = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetImageId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInstanceChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInstanceTypes(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetMultiAzPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetOnDemandBaseCapacity(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotInstancePools(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotInstanceRemedy(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetCompensateWithOnDemand(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPeriod(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Period = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPeriodUnit(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPlatform(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Platform = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetRamPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.RamPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotStrategy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotPriceLimit(v []*DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetRdsInstances(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetScalingGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ScalingGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetScalingPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSecurityGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSystemDiskCategory(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSystemDiskSize(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetTags(v []*Tag) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Tags = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetVswitchIds(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.VswitchIds = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetLoginPassword(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetKeyPair(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInternetChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInternetMaxBandwidthOut(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit struct {
	// 抢占式实例规格。
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// 单台实例上限价格，单位：元/小时。
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) SetInstanceType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) SetPriceLimit(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyStatus struct {
	// 失败节点数。
	FailedNodes *int64 `json:"failed_nodes,omitempty" xml:"failed_nodes,omitempty"`
	// 处于健康状态节点数。
	HealthyNodes *int64 `json:"healthy_nodes,omitempty" xml:"healthy_nodes,omitempty"`
	// 正在初始化节点数。
	InitialNodes *int64 `json:"initial_nodes,omitempty" xml:"initial_nodes,omitempty"`
	// 离线节点数量。
	OfflineNodes *int64 `json:"offline_nodes,omitempty" xml:"offline_nodes,omitempty"`
	// 正在被移除节点数。
	RemovingNodes *int64 `json:"removing_nodes,omitempty" xml:"removing_nodes,omitempty"`
	// 工作节点数量。
	ServingNodes *int64 `json:"serving_nodes,omitempty" xml:"serving_nodes,omitempty"`
	// 节点池状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 总节点数。
	TotalNodes *int64 `json:"total_nodes,omitempty" xml:"total_nodes,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyStatus) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyStatus) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetFailedNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.FailedNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetHealthyNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.HealthyNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetInitialNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.InitialNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetOfflineNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.OfflineNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetRemovingNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.RemovingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetServingNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.ServingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetState(v string) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.State = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetTotalNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.TotalNodes = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyTeeConfig struct {
	// 是否为加密计算节点池。
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyTeeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyTeeConfig) SetTeeEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyTeeConfig {
	s.TeeEnable = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyManagement struct {
	// 是否开启托管版节点池。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 自动修复。
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// 自动升级配置。
	UpgradeConfig *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolDetailResponseBodyManagement) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyManagement) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetAutoRepair(v bool) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.AutoRepair = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetUpgradeConfig(v *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.UpgradeConfig = v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig struct {
	// 是否启用自动升级，自修复。
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// 额外节点数量。
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// 额外节点比例， 和surge 二选一。
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
	// 最大不可用节点数量。
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetAutoUpgrade(v bool) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetSurge(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetSurgePercentage(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetMaxUnavailable(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

type DescribeClusterNodePoolDetailResponse struct {
	Headers map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterNodePoolDetailResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodePoolDetailResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponse) SetHeaders(v map[string]*string) *DescribeClusterNodePoolDetailResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponse) SetBody(v *DescribeClusterNodePoolDetailResponseBody) *DescribeClusterNodePoolDetailResponse {
	s.Body = v
	return s
}

type CreateClusterNodePoolRequest struct {
	// 自动伸缩节点池配置。
	AutoScaling *CreateClusterNodePoolRequestAutoScaling `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	// 集群配置
	KubernetesConfig *CreateClusterNodePoolRequestKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// 节点池配置
	NodepoolInfo *CreateClusterNodePoolRequestNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// 伸缩组配置
	ScalingGroup *CreateClusterNodePoolRequestScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	// 加密计算节点池配置。
	TeeConfig *CreateClusterNodePoolRequestTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
	// 托管节点池配置。
	Management *CreateClusterNodePoolRequestManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	// 节点数量。
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
}

func (s CreateClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequest) SetAutoScaling(v *CreateClusterNodePoolRequestAutoScaling) *CreateClusterNodePoolRequest {
	s.AutoScaling = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetKubernetesConfig(v *CreateClusterNodePoolRequestKubernetesConfig) *CreateClusterNodePoolRequest {
	s.KubernetesConfig = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetNodepoolInfo(v *CreateClusterNodePoolRequestNodepoolInfo) *CreateClusterNodePoolRequest {
	s.NodepoolInfo = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetScalingGroup(v *CreateClusterNodePoolRequestScalingGroup) *CreateClusterNodePoolRequest {
	s.ScalingGroup = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetTeeConfig(v *CreateClusterNodePoolRequestTeeConfig) *CreateClusterNodePoolRequest {
	s.TeeConfig = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetManagement(v *CreateClusterNodePoolRequestManagement) *CreateClusterNodePoolRequest {
	s.Management = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetCount(v int64) *CreateClusterNodePoolRequest {
	s.Count = &v
	return s
}

type CreateClusterNodePoolRequestAutoScaling struct {
	// 是否开启自动伸缩。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 最大实例数。
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// 最小实例数。
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// 扩容节点类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 是否绑定EIP。
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// EIP实例规格。
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// 带宽峰值。
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
}

func (s CreateClusterNodePoolRequestAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestAutoScaling) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEnable(v bool) *CreateClusterNodePoolRequestAutoScaling {
	s.Enable = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetMaxInstances(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetMinInstances(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetType(v string) *CreateClusterNodePoolRequestAutoScaling {
	s.Type = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetIsBondEip(v bool) *CreateClusterNodePoolRequestAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEipInternetChargeType(v string) *CreateClusterNodePoolRequestAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEipBandwidth(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.EipBandwidth = &v
	return s
}

type CreateClusterNodePoolRequestKubernetesConfig struct {
	// 是否开启云监控。
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// CPU管理策略。
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 节点标签。
	Labels []*Tag `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	// 容器运行时。
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 容器运行时版本。
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// 污点信息。
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 节点自定义数据。
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s CreateClusterNodePoolRequestKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestKubernetesConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetCmsEnabled(v bool) *CreateClusterNodePoolRequestKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetCpuPolicy(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetLabels(v []*Tag) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Labels = v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetRuntime(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetRuntimeVersion(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetTaints(v []*Taint) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Taints = v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetUserData(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.UserData = &v
	return s
}

type CreateClusterNodePoolRequestNodepoolInfo struct {
	// 节点池名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
}

func (s CreateClusterNodePoolRequestNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestNodepoolInfo) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestNodepoolInfo) SetName(v string) *CreateClusterNodePoolRequestNodepoolInfo {
	s.Name = &v
	return s
}

func (s *CreateClusterNodePoolRequestNodepoolInfo) SetResourceGroupId(v string) *CreateClusterNodePoolRequestNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroup struct {
	// 节点是否开启自动续费
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// 节点自动续费周期
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// 数据盘配置。
	DataDisks []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	// 自定义镜像。
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 节点付费类型
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// 实例规格。
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// 密钥对名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// SSH登录密码。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 节点包年包月时长。
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// 节点包年包月周期。
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// 操作系统发行版
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// RDS实例列表。
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 抢占式实例类型
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// 抢占实例价格上限配置。
	SpotPriceLimit []*CreateClusterNodePoolRequestScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// 自动伸缩。
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// 安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 节点系统盘类型。
	SystemDiskCategory *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	// 节点系统盘大小。
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// ECS标签
	Tags []*CreateClusterNodePoolRequestScalingGroupTags `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 虚拟交换机ID。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 多可用区伸缩组ECS实例扩缩容策略
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// 伸缩组所需要按量实例个数的最小值，取值范围：0~1000。当按量实例个数少于该值时，将优先创建按量实例。
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// 伸缩组满足最小按量实例数（OnDemandBaseCapacity）要求后，超出的实例中按量实例应占的比例，取值范围：0～100。
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// 指定可用实例规格的个数，伸缩组将按成本最低的多个规格均衡创建抢占式实例。取值范围：1~10。
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// 是否开启补齐抢占式实例。开启后，当收到抢占式实例将被回收的系统消息时，伸缩组将尝试创建新的实例，替换掉将被回收的抢占式实例。
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// 当MultiAZPolicy取值为COST_OPTIMIZED时，如果因价格、库存等原因无法创建足够的抢占式实例，是否允许自动尝试创建按量实例满足ECS实例数量要求。取值范围：true：允许。false：不允许。默认值：true
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// 节点公网IP网络计费类型
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// 节点公网IP出带宽最大值，单位为Mbps（Mega bit per second），取值范围：1~100
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroup) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetAutoRenew(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetAutoRenewPeriod(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetDataDisks(v []*DataDisk) *CreateClusterNodePoolRequestScalingGroup {
	s.DataDisks = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetImageId(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.ImageId = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInstanceChargeType(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInstanceTypes(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetKeyPair(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetLoginPassword(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPeriod(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.Period = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPeriodUnit(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPlatform(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.Platform = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetRdsInstances(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotStrategy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotPriceLimit(v []*CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetScalingPolicy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSecurityGroupId(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskCategory(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskSize(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetTags(v []*CreateClusterNodePoolRequestScalingGroupTags) *CreateClusterNodePoolRequestScalingGroup {
	s.Tags = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetVswitchIds(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.VswitchIds = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetMultiAzPolicy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetOnDemandBaseCapacity(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotInstancePools(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotInstanceRemedy(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetCompensateWithOnDemand(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInternetChargeType(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInternetMaxBandwidthOut(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroupSpotPriceLimit struct {
	// 抢占实例规格。
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// 抢占实例单价。
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) SetInstanceType(v string) *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) SetPriceLimit(v string) *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroupTags struct {
	// key
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// value
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroupTags) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroupTags) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroupTags) SetKey(v string) *CreateClusterNodePoolRequestScalingGroupTags {
	s.Key = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroupTags) SetValue(v string) *CreateClusterNodePoolRequestScalingGroupTags {
	s.Value = &v
	return s
}

type CreateClusterNodePoolRequestTeeConfig struct {
	// 是否为加密计算节点池。
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s CreateClusterNodePoolRequestTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestTeeConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestTeeConfig) SetTeeEnable(v bool) *CreateClusterNodePoolRequestTeeConfig {
	s.TeeEnable = &v
	return s
}

type CreateClusterNodePoolRequestManagement struct {
	// 是否启用托管节点池。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 是否启用自动修复。
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// 自动升级配置。
	UpgradeConfig *CreateClusterNodePoolRequestManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s CreateClusterNodePoolRequestManagement) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestManagement) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestManagement) SetEnable(v bool) *CreateClusterNodePoolRequestManagement {
	s.Enable = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagement) SetAutoRepair(v bool) *CreateClusterNodePoolRequestManagement {
	s.AutoRepair = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagement) SetUpgradeConfig(v *CreateClusterNodePoolRequestManagementUpgradeConfig) *CreateClusterNodePoolRequestManagement {
	s.UpgradeConfig = v
	return s
}

type CreateClusterNodePoolRequestManagementUpgradeConfig struct {
	// 是否启用自动升级
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// 额外节点数量。
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// 额外节点比例。和surge二选一。
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
	// 最大不可用节点数量。
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
}

func (s CreateClusterNodePoolRequestManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetAutoUpgrade(v bool) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetSurge(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetSurgePercentage(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetMaxUnavailable(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

type CreateClusterNodePoolResponseBody struct {
	// 节点池ID
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
}

func (s CreateClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolResponseBody) SetNodepoolId(v string) *CreateClusterNodePoolResponseBody {
	s.NodepoolId = &v
	return s
}

type CreateClusterNodePoolResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CreateClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolResponse) SetHeaders(v map[string]*string) *CreateClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *CreateClusterNodePoolResponse) SetBody(v *CreateClusterNodePoolResponseBody) *CreateClusterNodePoolResponse {
	s.Body = v
	return s
}

type DescribeClusterUserKubeconfigRequest struct {
	// ApiServer是否为内网地址。
	PrivateIpAddress *bool `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
	// 临时kubeconfig有效期，单位：分钟。  最小值：15（15分钟）  最大值：4320（3天）。
	TemporaryDurationMinutes *int64 `json:"TemporaryDurationMinutes,omitempty" xml:"TemporaryDurationMinutes,omitempty"`
}

func (s DescribeClusterUserKubeconfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigRequest) SetPrivateIpAddress(v bool) *DescribeClusterUserKubeconfigRequest {
	s.PrivateIpAddress = &v
	return s
}

func (s *DescribeClusterUserKubeconfigRequest) SetTemporaryDurationMinutes(v int64) *DescribeClusterUserKubeconfigRequest {
	s.TemporaryDurationMinutes = &v
	return s
}

type DescribeClusterUserKubeconfigResponseBody struct {
	// kubeconfig内容。
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
	// kubeconfig过期时间。格式：RFC3339 格式的 UTC 时间。
	Expiration *string `json:"expiration,omitempty" xml:"expiration,omitempty"`
}

func (s DescribeClusterUserKubeconfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigResponseBody) SetConfig(v string) *DescribeClusterUserKubeconfigResponseBody {
	s.Config = &v
	return s
}

func (s *DescribeClusterUserKubeconfigResponseBody) SetExpiration(v string) *DescribeClusterUserKubeconfigResponseBody {
	s.Expiration = &v
	return s
}

type DescribeClusterUserKubeconfigResponse struct {
	Headers map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterUserKubeconfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterUserKubeconfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigResponse) SetHeaders(v map[string]*string) *DescribeClusterUserKubeconfigResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterUserKubeconfigResponse) SetBody(v *DescribeClusterUserKubeconfigResponseBody) *DescribeClusterUserKubeconfigResponse {
	s.Body = v
	return s
}

type ScaleClusterRequest struct {
	// 节点是否安装云监控插件。
	CloudMonitorFlags *bool `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	// 扩容节点数。
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
	// 节点CPU策略。
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 失败是否回滚。
	DisableRollback *bool `json:"disable_rollback,omitempty" xml:"disable_rollback,omitempty"`
	// keypair名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// SSH登录密码。和keypair二选一。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 集群标签。
	Tags []*ScaleClusterRequestTags `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 节点污点标记。
	Taints []*ScaleClusterRequestTaints `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 节点交换机ID列表。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 节点是否开启Worker节点自动续费。
	WorkerAutoRenew *bool `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	// 自动续费周期。
	WorkerAutoRenewPeriod *int64 `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	// 是否挂载数据盘。
	WorkerDataDisk *bool `json:"worker_data_disk,omitempty" xml:"worker_data_disk,omitempty"`
	// Worker数据盘类型、大小等配置的组合。
	WorkerDataDisks []*ScaleClusterRequestWorkerDataDisks `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	// 节点付费类型。
	WorkerInstanceChargeType *string `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	// Worker节点ECS规格类型。
	WorkerInstanceTypes []*string `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	// 节点包年包月时长。
	WorkerPeriod *int64 `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	// 当指定为PrePaid的时候需要指定周期。
	WorkerPeriodUnit *string `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	// 节点系统盘类型。
	WorkerSystemDiskCategory *string `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	// 节点系统盘大小
	WorkerSystemDiskSize *int64 `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
}

func (s ScaleClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequest) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequest) SetCloudMonitorFlags(v bool) *ScaleClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *ScaleClusterRequest) SetCount(v int64) *ScaleClusterRequest {
	s.Count = &v
	return s
}

func (s *ScaleClusterRequest) SetCpuPolicy(v string) *ScaleClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *ScaleClusterRequest) SetDisableRollback(v bool) *ScaleClusterRequest {
	s.DisableRollback = &v
	return s
}

func (s *ScaleClusterRequest) SetKeyPair(v string) *ScaleClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *ScaleClusterRequest) SetLoginPassword(v string) *ScaleClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *ScaleClusterRequest) SetTags(v []*ScaleClusterRequestTags) *ScaleClusterRequest {
	s.Tags = v
	return s
}

func (s *ScaleClusterRequest) SetTaints(v []*ScaleClusterRequestTaints) *ScaleClusterRequest {
	s.Taints = v
	return s
}

func (s *ScaleClusterRequest) SetVswitchIds(v []*string) *ScaleClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerAutoRenew(v bool) *ScaleClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerAutoRenewPeriod(v int64) *ScaleClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerDataDisk(v bool) *ScaleClusterRequest {
	s.WorkerDataDisk = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerDataDisks(v []*ScaleClusterRequestWorkerDataDisks) *ScaleClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerInstanceChargeType(v string) *ScaleClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerInstanceTypes(v []*string) *ScaleClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerPeriod(v int64) *ScaleClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerPeriodUnit(v string) *ScaleClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerSystemDiskCategory(v string) *ScaleClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerSystemDiskSize(v int64) *ScaleClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

type ScaleClusterRequestTags struct {
	// 标签值。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
}

func (s ScaleClusterRequestTags) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestTags) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestTags) SetKey(v string) *ScaleClusterRequestTags {
	s.Key = &v
	return s
}

type ScaleClusterRequestTaints struct {
	// 污点生效策略。
	Effect *string `json:"effect,omitempty" xml:"effect,omitempty"`
	// 污点键。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// 污点值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ScaleClusterRequestTaints) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestTaints) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestTaints) SetEffect(v string) *ScaleClusterRequestTaints {
	s.Effect = &v
	return s
}

func (s *ScaleClusterRequestTaints) SetKey(v string) *ScaleClusterRequestTaints {
	s.Key = &v
	return s
}

func (s *ScaleClusterRequestTaints) SetValue(v string) *ScaleClusterRequestTaints {
	s.Value = &v
	return s
}

type ScaleClusterRequestWorkerDataDisks struct {
	// 数据盘类型。
	Category *string `json:"category,omitempty" xml:"category,omitempty"`
	// 是否对数据盘加密。
	Encrypted *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	// 数据盘大小。
	Size *string `json:"size,omitempty" xml:"size,omitempty"`
}

func (s ScaleClusterRequestWorkerDataDisks) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestWorkerDataDisks) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestWorkerDataDisks) SetCategory(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Category = &v
	return s
}

func (s *ScaleClusterRequestWorkerDataDisks) SetEncrypted(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Encrypted = &v
	return s
}

func (s *ScaleClusterRequestWorkerDataDisks) SetSize(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Size = &v
	return s
}

type ScaleClusterResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 请求ID。
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleClusterResponseBody) SetClusterId(v string) *ScaleClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ScaleClusterResponseBody) SetRequestId(v string) *ScaleClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ScaleClusterResponseBody) SetTaskId(v string) *ScaleClusterResponseBody {
	s.TaskId = &v
	return s
}

type ScaleClusterResponse struct {
	Headers map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ScaleClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterResponse) GoString() string {
	return s.String()
}

func (s *ScaleClusterResponse) SetHeaders(v map[string]*string) *ScaleClusterResponse {
	s.Headers = v
	return s
}

func (s *ScaleClusterResponse) SetBody(v *ScaleClusterResponseBody) *ScaleClusterResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonUpgradeStatusResponse struct {
	Headers map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonUpgradeStatusResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonUpgradeStatusResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonUpgradeStatusResponse {
	s.Body = v
	return s
}

type DescribeAddonsRequest struct {
	// 地域ID。
	Region *string `json:"region,omitempty" xml:"region,omitempty"`
	// 集群类型。  - Kubernetes: 专有版集群。 - ManagedKubernetes：托管版集群。 - Ask：Serverless 集群。 - ExternalKubernetes：注册到ACK的外部集群。
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
}

func (s DescribeAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsRequest) GoString() string {
	return s.String()
}

func (s *DescribeAddonsRequest) SetRegion(v string) *DescribeAddonsRequest {
	s.Region = &v
	return s
}

func (s *DescribeAddonsRequest) SetClusterType(v string) *DescribeAddonsRequest {
	s.ClusterType = &v
	return s
}

type DescribeAddonsResponseBody struct {
	// 组件分组信息，例如：存储类组件，网络组件等。
	ComponentGroups []*DescribeAddonsResponseBodyComponentGroups `json:"ComponentGroups,omitempty" xml:"ComponentGroups,omitempty" type:"Repeated"`
	// 标准组件信息，包含各个组件的描述信息。
	StandardComponents map[string]*StandardComponentsValue `json:"StandardComponents,omitempty" xml:"StandardComponents,omitempty"`
}

func (s DescribeAddonsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBody) SetComponentGroups(v []*DescribeAddonsResponseBodyComponentGroups) *DescribeAddonsResponseBody {
	s.ComponentGroups = v
	return s
}

func (s *DescribeAddonsResponseBody) SetStandardComponents(v map[string]*StandardComponentsValue) *DescribeAddonsResponseBody {
	s.StandardComponents = v
	return s
}

type DescribeAddonsResponseBodyComponentGroups struct {
	// 组件组名称。
	GroupName *string `json:"group_name,omitempty" xml:"group_name,omitempty"`
	// 组件列表
	Items []*DescribeAddonsResponseBodyComponentGroupsItems `json:"items,omitempty" xml:"items,omitempty" type:"Repeated"`
}

func (s DescribeAddonsResponseBodyComponentGroups) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBodyComponentGroups) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBodyComponentGroups) SetGroupName(v string) *DescribeAddonsResponseBodyComponentGroups {
	s.GroupName = &v
	return s
}

func (s *DescribeAddonsResponseBodyComponentGroups) SetItems(v []*DescribeAddonsResponseBodyComponentGroupsItems) *DescribeAddonsResponseBodyComponentGroups {
	s.Items = v
	return s
}

type DescribeAddonsResponseBodyComponentGroupsItems struct {
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s DescribeAddonsResponseBodyComponentGroupsItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBodyComponentGroupsItems) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBodyComponentGroupsItems) SetName(v string) *DescribeAddonsResponseBodyComponentGroupsItems {
	s.Name = &v
	return s
}

type DescribeAddonsResponse struct {
	Headers map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeAddonsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponse) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponse) SetHeaders(v map[string]*string) *DescribeAddonsResponse {
	s.Headers = v
	return s
}

func (s *DescribeAddonsResponse) SetBody(v *DescribeAddonsResponseBody) *DescribeAddonsResponse {
	s.Body = v
	return s
}

type CreateAutoscalingConfigRequest struct {
	// 静默时间，扩容出的节点，在静默时间过后，方可进入缩容判断
	CoolDownDuration *string `json:"cool_down_duration,omitempty" xml:"cool_down_duration,omitempty"`
	// 缩容触发时延，节点缩容时需要连续满足触发时延所设定的时间，方可进行缩容
	UnneededDuration *string `json:"unneeded_duration,omitempty" xml:"unneeded_duration,omitempty"`
	// 缩容阈值，节点上 Request 的资源与总资源量的比值
	UtilizationThreshold *string `json:"utilization_threshold,omitempty" xml:"utilization_threshold,omitempty"`
	// GPU缩容阈值，节点上 Request 的资源与总资源量的比值
	GpuUtilizationThreshold *string `json:"gpu_utilization_threshold,omitempty" xml:"gpu_utilization_threshold,omitempty"`
	// 弹性灵敏度，判断伸缩的间隔时间
	ScanInterval *string `json:"scan_interval,omitempty" xml:"scan_interval,omitempty"`
}

func (s CreateAutoscalingConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAutoscalingConfigRequest) GoString() string {
	return s.String()
}

func (s *CreateAutoscalingConfigRequest) SetCoolDownDuration(v string) *CreateAutoscalingConfigRequest {
	s.CoolDownDuration = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetUnneededDuration(v string) *CreateAutoscalingConfigRequest {
	s.UnneededDuration = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetUtilizationThreshold(v string) *CreateAutoscalingConfigRequest {
	s.UtilizationThreshold = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetGpuUtilizationThreshold(v string) *CreateAutoscalingConfigRequest {
	s.GpuUtilizationThreshold = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetScanInterval(v string) *CreateAutoscalingConfigRequest {
	s.ScanInterval = &v
	return s
}

type CreateAutoscalingConfigResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s CreateAutoscalingConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAutoscalingConfigResponse) GoString() string {
	return s.String()
}

func (s *CreateAutoscalingConfigResponse) SetHeaders(v map[string]*string) *CreateAutoscalingConfigResponse {
	s.Headers = v
	return s
}

type CreateClusterRequest struct {
	// 集群名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 地域ID
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 集群类型
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// 托管版集群类型
	ClusterSpec *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	// 集群版本
	KubernetesVersion *string  `json:"kubernetes_version,omitempty" xml:"kubernetes_version,omitempty"`
	Runtime           *Runtime `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 集群专有网络ID。
	Vpcid *string `json:"vpcid,omitempty" xml:"vpcid,omitempty"`
	// 创建Terway网络类型的集群时，需要为pod指定单独的虚拟交换机
	PodVswitchIds []*string `json:"pod_vswitch_ids,omitempty" xml:"pod_vswitch_ids,omitempty" type:"Repeated"`
	// POD网络网段
	ContainerCidr *string `json:"container_cidr,omitempty" xml:"container_cidr,omitempty"`
	// 服务网络网段
	ServiceCidr *string `json:"service_cidr,omitempty" xml:"service_cidr,omitempty"`
	// 安全组ID，和is_enterprise_security_group二选一
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 自动创建企业安全组
	IsEnterpriseSecurityGroup *bool `json:"is_enterprise_security_group,omitempty" xml:"is_enterprise_security_group,omitempty"`
	// 若您集群内的节点、应用等需要访问公网，勾选该项后我们将为您创建 NAT 网关并自动配置 SNAT 规则
	SnatEntry *bool `json:"snat_entry,omitempty" xml:"snat_entry,omitempty"`
	// 使用EIP暴露apiServer
	EndpointPublicAccess *bool `json:"endpoint_public_access,omitempty" xml:"endpoint_public_access,omitempty"`
	// 允许公网ssh登录
	SshFlags *bool `json:"ssh_flags,omitempty" xml:"ssh_flags,omitempty"`
	// 时区
	Timezone *string `json:"timezone,omitempty" xml:"timezone,omitempty"`
	// 节点IP数量
	NodeCidrMask *string `json:"node_cidr_mask,omitempty" xml:"node_cidr_mask,omitempty"`
	// 自定义集群CA
	UserCa *string `json:"user_ca,omitempty" xml:"user_ca,omitempty"`
	// 节点自定义数据
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
	// 集群本地域名
	ClusterDomain *string `json:"cluster_domain,omitempty" xml:"cluster_domain,omitempty"`
	// 自定义节点名称
	NodeNameMode *string `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	// 自定义证书SAN
	CustomSan *string `json:"custom_san,omitempty" xml:"custom_san,omitempty"`
	// Secret落盘加密
	EncryptionProviderKey *string `json:"encryption_provider_key,omitempty" xml:"encryption_provider_key,omitempty"`
	// serviceaccount token中的签发身份，即token payload中的iss字段。
	ServiceAccountIssuer *string `json:"service_account_issuer,omitempty" xml:"service_account_issuer,omitempty"`
	// 合法的请求token身份，用于apiserver服务端认证请求token是否合法。
	ApiAudiences *string `json:"api_audiences,omitempty" xml:"api_audiences,omitempty"`
	// 自定义镜像
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// RDS白名单
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 集群标签
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 集群组件配置
	Addons []*Addon `json:"addons,omitempty" xml:"addons,omitempty" type:"Repeated"`
	// 节点污点信息
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 为ECS安装云监控
	CloudMonitorFlags *bool `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	// 操作系统发行版
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// 操作系统平台类型
	OsType *string `json:"os_type,omitempty" xml:"os_type,omitempty"`
	// CPU策略
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// Proxy代理模式，ipvs|iptables
	ProxyMode *string `json:"proxy_mode,omitempty" xml:"proxy_mode,omitempty"`
	// 节点服务端口范围
	NodePortRange *string `json:"node_port_range,omitempty" xml:"node_port_range,omitempty"`
	// 密钥对名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// SSH登录密码。密码规则为8~30 个字符，且至少同时包含三项（大小写字母、数字和特殊符号），和key_pair二选一。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 集群Master节点数量
	MasterCount *int64 `json:"master_count,omitempty" xml:"master_count,omitempty"`
	// 集群Master节点使用的虚拟交换机
	MasterVswitchIds []*string `json:"master_vswitch_ids,omitempty" xml:"master_vswitch_ids,omitempty" type:"Repeated"`
	// 集群Master节点类型
	MasterInstanceTypes []*string `json:"master_instance_types,omitempty" xml:"master_instance_types,omitempty" type:"Repeated"`
	// 集群Master节点系统盘类型
	MasterSystemDiskCategory *string `json:"master_system_disk_category,omitempty" xml:"master_system_disk_category,omitempty"`
	// 集群Master节点系统盘大小，至少40
	MasterSystemDiskSize *int64 `json:"master_system_disk_size,omitempty" xml:"master_system_disk_size,omitempty"`
	// 集群Master节点自动快照备份策略
	MasterSystemDiskSnapshotPolicyId *string `json:"master_system_disk_snapshot_policy_id,omitempty" xml:"master_system_disk_snapshot_policy_id,omitempty"`
	// 集群Master节点付费类型
	MasterInstanceChargeType *string `json:"master_instance_charge_type,omitempty" xml:"master_instance_charge_type,omitempty"`
	// 集群Master节点包年包月周期
	MasterPeriodUnit *string `json:"master_period_unit,omitempty" xml:"master_period_unit,omitempty"`
	// 集群Master节点包年包月时长
	MasterPeriod *int64 `json:"master_period,omitempty" xml:"master_period,omitempty"`
	// 集群Master节点是否自动续费
	MasterAutoRenew *bool `json:"master_auto_renew,omitempty" xml:"master_auto_renew,omitempty"`
	// 集群Master节点自动续费时长
	MasterAutoRenewPeriod *int64 `json:"master_auto_renew_period,omitempty" xml:"master_auto_renew_period,omitempty"`
	// 集群Worker节点数量
	NumOfNodes *int64 `json:"num_of_nodes,omitempty" xml:"num_of_nodes,omitempty"`
	// 集群节点所在虚拟交换机。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 集群Worker节点所在虚拟交换机
	WorkerVswitchIds []*string `json:"worker_vswitch_ids,omitempty" xml:"worker_vswitch_ids,omitempty" type:"Repeated"`
	// 集群Worker节点类型
	WorkerInstanceTypes []*string `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	// 集群Worker节点系统盘类型
	WorkerSystemDiskCategory *string `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	// 集群Worker节点系统盘大小
	WorkerSystemDiskSize *int64 `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
	// 集群Worker节点系统盘快照备份策略
	WorkerSystemDiskSnapshotPolicyId *string `json:"worker_system_disk_snapshot_policy_id,omitempty" xml:"worker_system_disk_snapshot_policy_id,omitempty"`
	// 集群Worker节点数据盘配置
	WorkerDataDisks []*DataDisk `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	// 集群Worker节点付费类型
	WorkerInstanceChargeType *string `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	// 集群Worker节点包年包月周期
	WorkerPeriodUnit *string `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	// 集群Worker节点包年包月时长
	WorkerPeriod *int64 `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	// 集群Worker节点到期是否自动续费
	WorkerAutoRenew *bool `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	// 集群Worker节点自动续费时长
	WorkerAutoRenewPeriod *int64 `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	// 使用已有节点创建集群时，已有实例列表
	Instances []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
	// 使用已有节点创建集群时，是否格式化已有实例的磁盘
	FormatDisk *bool `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	// 使用已有节点创建集群时，是否保留实例名称。
	KeepInstanceName *bool `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	// 创建ASK集群时，服务发现类型
	ServiceDiscoveryTypes []*string `json:"service_discovery_types,omitempty" xml:"service_discovery_types,omitempty" type:"Repeated"`
	// 使用自动创建专有网络创建ASK集群时，是否在vpc中创建Nat网关并配置SNAT规则。
	NatGateway *bool `json:"nat_gateway,omitempty" xml:"nat_gateway,omitempty"`
	// 使用自动创建专有网络创建ASK集群时，需要指定专有网络的可用区
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
	// 面向场景时的集群类型。  Default：非边缘场景集群。 Edge：边缘场景集群。
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// 集群删除保护
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// 失败回滚
	DisableRollback *bool `json:"disable_rollback,omitempty" xml:"disable_rollback,omitempty"`
	// 集群创建超时时间
	TimeoutMins *int64 `json:"timeout_mins,omitempty" xml:"timeout_mins,omitempty"`
}

func (s CreateClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterRequest) GoString() string {
	return s.String()
}

func (s *CreateClusterRequest) SetName(v string) *CreateClusterRequest {
	s.Name = &v
	return s
}

func (s *CreateClusterRequest) SetRegionId(v string) *CreateClusterRequest {
	s.RegionId = &v
	return s
}

func (s *CreateClusterRequest) SetClusterType(v string) *CreateClusterRequest {
	s.ClusterType = &v
	return s
}

func (s *CreateClusterRequest) SetClusterSpec(v string) *CreateClusterRequest {
	s.ClusterSpec = &v
	return s
}

func (s *CreateClusterRequest) SetKubernetesVersion(v string) *CreateClusterRequest {
	s.KubernetesVersion = &v
	return s
}

func (s *CreateClusterRequest) SetRuntime(v *Runtime) *CreateClusterRequest {
	s.Runtime = v
	return s
}

func (s *CreateClusterRequest) SetVpcid(v string) *CreateClusterRequest {
	s.Vpcid = &v
	return s
}

func (s *CreateClusterRequest) SetPodVswitchIds(v []*string) *CreateClusterRequest {
	s.PodVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetContainerCidr(v string) *CreateClusterRequest {
	s.ContainerCidr = &v
	return s
}

func (s *CreateClusterRequest) SetServiceCidr(v string) *CreateClusterRequest {
	s.ServiceCidr = &v
	return s
}

func (s *CreateClusterRequest) SetSecurityGroupId(v string) *CreateClusterRequest {
	s.SecurityGroupId = &v
	return s
}

func (s *CreateClusterRequest) SetIsEnterpriseSecurityGroup(v bool) *CreateClusterRequest {
	s.IsEnterpriseSecurityGroup = &v
	return s
}

func (s *CreateClusterRequest) SetSnatEntry(v bool) *CreateClusterRequest {
	s.SnatEntry = &v
	return s
}

func (s *CreateClusterRequest) SetEndpointPublicAccess(v bool) *CreateClusterRequest {
	s.EndpointPublicAccess = &v
	return s
}

func (s *CreateClusterRequest) SetSshFlags(v bool) *CreateClusterRequest {
	s.SshFlags = &v
	return s
}

func (s *CreateClusterRequest) SetTimezone(v string) *CreateClusterRequest {
	s.Timezone = &v
	return s
}

func (s *CreateClusterRequest) SetNodeCidrMask(v string) *CreateClusterRequest {
	s.NodeCidrMask = &v
	return s
}

func (s *CreateClusterRequest) SetUserCa(v string) *CreateClusterRequest {
	s.UserCa = &v
	return s
}

func (s *CreateClusterRequest) SetUserData(v string) *CreateClusterRequest {
	s.UserData = &v
	return s
}

func (s *CreateClusterRequest) SetClusterDomain(v string) *CreateClusterRequest {
	s.ClusterDomain = &v
	return s
}

func (s *CreateClusterRequest) SetNodeNameMode(v string) *CreateClusterRequest {
	s.NodeNameMode = &v
	return s
}

func (s *CreateClusterRequest) SetCustomSan(v string) *CreateClusterRequest {
	s.CustomSan = &v
	return s
}

func (s *CreateClusterRequest) SetEncryptionProviderKey(v string) *CreateClusterRequest {
	s.EncryptionProviderKey = &v
	return s
}

func (s *CreateClusterRequest) SetServiceAccountIssuer(v string) *CreateClusterRequest {
	s.ServiceAccountIssuer = &v
	return s
}

func (s *CreateClusterRequest) SetApiAudiences(v string) *CreateClusterRequest {
	s.ApiAudiences = &v
	return s
}

func (s *CreateClusterRequest) SetImageId(v string) *CreateClusterRequest {
	s.ImageId = &v
	return s
}

func (s *CreateClusterRequest) SetRdsInstances(v []*string) *CreateClusterRequest {
	s.RdsInstances = v
	return s
}

func (s *CreateClusterRequest) SetTags(v []*Tag) *CreateClusterRequest {
	s.Tags = v
	return s
}

func (s *CreateClusterRequest) SetAddons(v []*Addon) *CreateClusterRequest {
	s.Addons = v
	return s
}

func (s *CreateClusterRequest) SetTaints(v []*Taint) *CreateClusterRequest {
	s.Taints = v
	return s
}

func (s *CreateClusterRequest) SetCloudMonitorFlags(v bool) *CreateClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *CreateClusterRequest) SetPlatform(v string) *CreateClusterRequest {
	s.Platform = &v
	return s
}

func (s *CreateClusterRequest) SetOsType(v string) *CreateClusterRequest {
	s.OsType = &v
	return s
}

func (s *CreateClusterRequest) SetCpuPolicy(v string) *CreateClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *CreateClusterRequest) SetProxyMode(v string) *CreateClusterRequest {
	s.ProxyMode = &v
	return s
}

func (s *CreateClusterRequest) SetNodePortRange(v string) *CreateClusterRequest {
	s.NodePortRange = &v
	return s
}

func (s *CreateClusterRequest) SetKeyPair(v string) *CreateClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *CreateClusterRequest) SetLoginPassword(v string) *CreateClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *CreateClusterRequest) SetMasterCount(v int64) *CreateClusterRequest {
	s.MasterCount = &v
	return s
}

func (s *CreateClusterRequest) SetMasterVswitchIds(v []*string) *CreateClusterRequest {
	s.MasterVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetMasterInstanceTypes(v []*string) *CreateClusterRequest {
	s.MasterInstanceTypes = v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskCategory(v string) *CreateClusterRequest {
	s.MasterSystemDiskCategory = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskSize(v int64) *CreateClusterRequest {
	s.MasterSystemDiskSize = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskSnapshotPolicyId(v string) *CreateClusterRequest {
	s.MasterSystemDiskSnapshotPolicyId = &v
	return s
}

func (s *CreateClusterRequest) SetMasterInstanceChargeType(v string) *CreateClusterRequest {
	s.MasterInstanceChargeType = &v
	return s
}

func (s *CreateClusterRequest) SetMasterPeriodUnit(v string) *CreateClusterRequest {
	s.MasterPeriodUnit = &v
	return s
}

func (s *CreateClusterRequest) SetMasterPeriod(v int64) *CreateClusterRequest {
	s.MasterPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetMasterAutoRenew(v bool) *CreateClusterRequest {
	s.MasterAutoRenew = &v
	return s
}

func (s *CreateClusterRequest) SetMasterAutoRenewPeriod(v int64) *CreateClusterRequest {
	s.MasterAutoRenewPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetNumOfNodes(v int64) *CreateClusterRequest {
	s.NumOfNodes = &v
	return s
}

func (s *CreateClusterRequest) SetVswitchIds(v []*string) *CreateClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetWorkerVswitchIds(v []*string) *CreateClusterRequest {
	s.WorkerVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetWorkerInstanceTypes(v []*string) *CreateClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskCategory(v string) *CreateClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskSize(v int64) *CreateClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskSnapshotPolicyId(v string) *CreateClusterRequest {
	s.WorkerSystemDiskSnapshotPolicyId = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerDataDisks(v []*DataDisk) *CreateClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *CreateClusterRequest) SetWorkerInstanceChargeType(v string) *CreateClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerPeriodUnit(v string) *CreateClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerPeriod(v int64) *CreateClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerAutoRenew(v bool) *CreateClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerAutoRenewPeriod(v int64) *CreateClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetInstances(v []*string) *CreateClusterRequest {
	s.Instances = v
	return s
}

func (s *CreateClusterRequest) SetFormatDisk(v bool) *CreateClusterRequest {
	s.FormatDisk = &v
	return s
}

func (s *CreateClusterRequest) SetKeepInstanceName(v bool) *CreateClusterRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *CreateClusterRequest) SetServiceDiscoveryTypes(v []*string) *CreateClusterRequest {
	s.ServiceDiscoveryTypes = v
	return s
}

func (s *CreateClusterRequest) SetNatGateway(v bool) *CreateClusterRequest {
	s.NatGateway = &v
	return s
}

func (s *CreateClusterRequest) SetZoneId(v string) *CreateClusterRequest {
	s.ZoneId = &v
	return s
}

func (s *CreateClusterRequest) SetProfile(v string) *CreateClusterRequest {
	s.Profile = &v
	return s
}

func (s *CreateClusterRequest) SetDeletionProtection(v bool) *CreateClusterRequest {
	s.DeletionProtection = &v
	return s
}

func (s *CreateClusterRequest) SetDisableRollback(v bool) *CreateClusterRequest {
	s.DisableRollback = &v
	return s
}

func (s *CreateClusterRequest) SetTimeoutMins(v int64) *CreateClusterRequest {
	s.TimeoutMins = &v
	return s
}

type CreateClusterResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 请求ID。
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s CreateClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterResponseBody) GoString() string {
	return s.String()
}

func (s *CreateClusterResponseBody) SetClusterId(v string) *CreateClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *CreateClusterResponseBody) SetRequestId(v string) *CreateClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *CreateClusterResponseBody) SetTaskId(v string) *CreateClusterResponseBody {
	s.TaskId = &v
	return s
}

type CreateClusterResponse struct {
	Headers map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CreateClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterResponse) GoString() string {
	return s.String()
}

func (s *CreateClusterResponse) SetHeaders(v map[string]*string) *CreateClusterResponse {
	s.Headers = v
	return s
}

func (s *CreateClusterResponse) SetBody(v *CreateClusterResponseBody) *CreateClusterResponse {
	s.Body = v
	return s
}

type UpgradeClusterRequest struct {
	// 组件名称，集群升级时取值"k8s"。
	ComponentName *string `json:"component_name,omitempty" xml:"component_name,omitempty"`
	// 目标版本。
	NextVersion *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// 当前版本。
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s UpgradeClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterRequest) GoString() string {
	return s.String()
}

func (s *UpgradeClusterRequest) SetComponentName(v string) *UpgradeClusterRequest {
	s.ComponentName = &v
	return s
}

func (s *UpgradeClusterRequest) SetNextVersion(v string) *UpgradeClusterRequest {
	s.NextVersion = &v
	return s
}

func (s *UpgradeClusterRequest) SetVersion(v string) *UpgradeClusterRequest {
	s.Version = &v
	return s
}

type UpgradeClusterResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UpgradeClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterResponse) GoString() string {
	return s.String()
}

func (s *UpgradeClusterResponse) SetHeaders(v map[string]*string) *UpgradeClusterResponse {
	s.Headers = v
	return s
}

type CancelWorkflowRequest struct {
	// 执行的操作，目前只支持cancel。
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s CancelWorkflowRequest) String() string {
	return tea.Prettify(s)
}

func (s CancelWorkflowRequest) GoString() string {
	return s.String()
}

func (s *CancelWorkflowRequest) SetAction(v string) *CancelWorkflowRequest {
	s.Action = &v
	return s
}

type CancelWorkflowResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s CancelWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelWorkflowResponse) GoString() string {
	return s.String()
}

func (s *CancelWorkflowResponse) SetHeaders(v map[string]*string) *CancelWorkflowResponse {
	s.Headers = v
	return s
}

type AttachInstancesRequest struct {
	// 实例列表。
	Instances []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
	// key_pair名称，与login_password二选一
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// password，与key_pair二选一。
	Password *string `json:"password,omitempty" xml:"password,omitempty"`
	// 是否格式化数据盘。
	FormatDisk *bool `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	// 是否保留实例名称。
	KeepInstanceName *bool `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	// 是否为边缘节点。
	IsEdgeWorker *bool `json:"is_edge_worker,omitempty" xml:"is_edge_worker,omitempty"`
	// 节点池ID，欲将节点添加到哪个节点池中。。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 自定义镜像ID。
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// CPU亲和策略。
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 节点自定义数据。
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
	// RDS实例列表。
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	Runtime      *Runtime  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 节点标签。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
}

func (s AttachInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesRequest) GoString() string {
	return s.String()
}

func (s *AttachInstancesRequest) SetInstances(v []*string) *AttachInstancesRequest {
	s.Instances = v
	return s
}

func (s *AttachInstancesRequest) SetKeyPair(v string) *AttachInstancesRequest {
	s.KeyPair = &v
	return s
}

func (s *AttachInstancesRequest) SetPassword(v string) *AttachInstancesRequest {
	s.Password = &v
	return s
}

func (s *AttachInstancesRequest) SetFormatDisk(v bool) *AttachInstancesRequest {
	s.FormatDisk = &v
	return s
}

func (s *AttachInstancesRequest) SetKeepInstanceName(v bool) *AttachInstancesRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *AttachInstancesRequest) SetIsEdgeWorker(v bool) *AttachInstancesRequest {
	s.IsEdgeWorker = &v
	return s
}

func (s *AttachInstancesRequest) SetNodepoolId(v string) *AttachInstancesRequest {
	s.NodepoolId = &v
	return s
}

func (s *AttachInstancesRequest) SetImageId(v string) *AttachInstancesRequest {
	s.ImageId = &v
	return s
}

func (s *AttachInstancesRequest) SetCpuPolicy(v string) *AttachInstancesRequest {
	s.CpuPolicy = &v
	return s
}

func (s *AttachInstancesRequest) SetUserData(v string) *AttachInstancesRequest {
	s.UserData = &v
	return s
}

func (s *AttachInstancesRequest) SetRdsInstances(v []*string) *AttachInstancesRequest {
	s.RdsInstances = v
	return s
}

func (s *AttachInstancesRequest) SetRuntime(v *Runtime) *AttachInstancesRequest {
	s.Runtime = v
	return s
}

func (s *AttachInstancesRequest) SetTags(v []*Tag) *AttachInstancesRequest {
	s.Tags = v
	return s
}

type AttachInstancesResponseBody struct {
	// 节点信息列表。
	List []*AttachInstancesResponseBodyList `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s AttachInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponseBody) SetList(v []*AttachInstancesResponseBodyList) *AttachInstancesResponseBody {
	s.List = v
	return s
}

func (s *AttachInstancesResponseBody) SetTaskId(v string) *AttachInstancesResponseBody {
	s.TaskId = &v
	return s
}

type AttachInstancesResponseBodyList struct {
	// 状态码。
	Code *string `json:"code,omitempty" xml:"code,omitempty"`
	// 实例ID。
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	// 添加结果描述。
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s AttachInstancesResponseBodyList) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponseBodyList) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponseBodyList) SetCode(v string) *AttachInstancesResponseBodyList {
	s.Code = &v
	return s
}

func (s *AttachInstancesResponseBodyList) SetInstanceId(v string) *AttachInstancesResponseBodyList {
	s.InstanceId = &v
	return s
}

func (s *AttachInstancesResponseBodyList) SetMessage(v string) *AttachInstancesResponseBodyList {
	s.Message = &v
	return s
}

type AttachInstancesResponse struct {
	Headers map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *AttachInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AttachInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponse) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponse) SetHeaders(v map[string]*string) *AttachInstancesResponse {
	s.Headers = v
	return s
}

func (s *AttachInstancesResponse) SetBody(v *AttachInstancesResponseBody) *AttachInstancesResponse {
	s.Body = v
	return s
}

type DescribeTemplatesRequest struct {
	// 模板类型，部署模板类型，目前一共有2种类型，取值为：kubernetes或compose。
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
	// 对查询结果进行分页处理，指定返回第几页的数据。  默认值为 1
	PageNum *int64 `json:"page_num,omitempty" xml:"page_num,omitempty"`
	// 对查询结果进行分页处理，指定每页包含的数据条数。  默认值为 10
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
}

func (s DescribeTemplatesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesRequest) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesRequest) SetTemplateType(v string) *DescribeTemplatesRequest {
	s.TemplateType = &v
	return s
}

func (s *DescribeTemplatesRequest) SetPageNum(v int64) *DescribeTemplatesRequest {
	s.PageNum = &v
	return s
}

func (s *DescribeTemplatesRequest) SetPageSize(v int64) *DescribeTemplatesRequest {
	s.PageSize = &v
	return s
}

type DescribeTemplatesResponseBody struct {
	// 模板列表。
	Templates []*DescribeTemplatesResponseBodyTemplates `json:"templates,omitempty" xml:"templates,omitempty" type:"Repeated"`
	// 分页信息。
	PageInfo *DescribeTemplatesResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeTemplatesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBody) SetTemplates(v []*DescribeTemplatesResponseBodyTemplates) *DescribeTemplatesResponseBody {
	s.Templates = v
	return s
}

func (s *DescribeTemplatesResponseBody) SetPageInfo(v *DescribeTemplatesResponseBodyPageInfo) *DescribeTemplatesResponseBody {
	s.PageInfo = v
	return s
}

type DescribeTemplatesResponseBodyTemplates struct {
	// 模板访问权限，取值为：private、pubilc或shared。。
	Acl *string `json:"acl,omitempty" xml:"acl,omitempty"`
	// 模板ID。会模板随着更新而变化。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 模板名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 模板描述信息。
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 模板标签，如果不显式指定了，默认为模板名称。
	Tags *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// 模板详情。
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
	// 部署模板类型。
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
	// 模板创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 模板修改时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// 模板唯一ID。
	TemplateWithHistId *string `json:"template_with_hist_id,omitempty" xml:"template_with_hist_id,omitempty"`
}

func (s DescribeTemplatesResponseBodyTemplates) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBodyTemplates) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBodyTemplates) SetAcl(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Acl = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetId(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Id = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetName(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Name = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetDescription(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Description = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTags(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Tags = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplate(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Template = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplateType(v string) *DescribeTemplatesResponseBodyTemplates {
	s.TemplateType = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetCreated(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Created = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetUpdated(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Updated = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplateWithHistId(v string) *DescribeTemplatesResponseBodyTemplates {
	s.TemplateWithHistId = &v
	return s
}

type DescribeTemplatesResponseBodyPageInfo struct {
	// 当前页数。
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// 单页最大数据条数。
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 结果总数。
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeTemplatesResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetPageNumber(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetPageSize(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetTotalCount(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeTemplatesResponse struct {
	Headers map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeTemplatesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeTemplatesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponse) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponse) SetHeaders(v map[string]*string) *DescribeTemplatesResponse {
	s.Headers = v
	return s
}

func (s *DescribeTemplatesResponse) SetBody(v *DescribeTemplatesResponseBody) *DescribeTemplatesResponse {
	s.Body = v
	return s
}

type PauseClusterUpgradeResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s PauseClusterUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s PauseClusterUpgradeResponse) GoString() string {
	return s.String()
}

func (s *PauseClusterUpgradeResponse) SetHeaders(v map[string]*string) *PauseClusterUpgradeResponse {
	s.Headers = v
	return s
}

type DeleteTemplateResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s DeleteTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteTemplateResponse) GoString() string {
	return s.String()
}

func (s *DeleteTemplateResponse) SetHeaders(v map[string]*string) *DeleteTemplateResponse {
	s.Headers = v
	return s
}

type DescribeTemplateAttributeRequest struct {
	// 模板类型，值为创建部署模板时指定的模板类型。
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s DescribeTemplateAttributeRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeRequest) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeRequest) SetTemplateType(v string) *DescribeTemplateAttributeRequest {
	s.TemplateType = &v
	return s
}

type DescribeTemplateAttributeResponse struct {
	Headers map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeTemplateAttributeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeTemplateAttributeResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeResponse) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeResponse) SetHeaders(v map[string]*string) *DescribeTemplateAttributeResponse {
	s.Headers = v
	return s
}

func (s *DescribeTemplateAttributeResponse) SetBody(v []*DescribeTemplateAttributeResponseBody) *DescribeTemplateAttributeResponse {
	s.Body = v
	return s
}

type DescribeTemplateAttributeResponseBody struct {
	// 编排模板ID，模板每次修改，这个ID都会改变。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 编排模板权限。取值：private，public，shared。
	Acl *string `json:"acl,omitempty" xml:"acl,omitempty"`
	// 编排模板名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 编排模板内容。
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
	// 编排模板类型
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
	// 编排模板描述。
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 部署模板的标签。
	Tags *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// 编排模板ID，该ID唯一不随更新而改变。
	TemplateWithHistId *string `json:"template_with_hist_id,omitempty" xml:"template_with_hist_id,omitempty"`
	// 编排模板创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 编排模板修改时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeTemplateAttributeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeResponseBody) SetId(v string) *DescribeTemplateAttributeResponseBody {
	s.Id = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetAcl(v string) *DescribeTemplateAttributeResponseBody {
	s.Acl = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetName(v string) *DescribeTemplateAttributeResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplate(v string) *DescribeTemplateAttributeResponseBody {
	s.Template = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplateType(v string) *DescribeTemplateAttributeResponseBody {
	s.TemplateType = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetDescription(v string) *DescribeTemplateAttributeResponseBody {
	s.Description = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTags(v string) *DescribeTemplateAttributeResponseBody {
	s.Tags = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplateWithHistId(v string) *DescribeTemplateAttributeResponseBody {
	s.TemplateWithHistId = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetCreated(v string) *DescribeTemplateAttributeResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetUpdated(v string) *DescribeTemplateAttributeResponseBody {
	s.Updated = &v
	return s
}

type CreateTemplateRequest struct {
	// 模板名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// YAML格式的模板内容。
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
	// 模板标签。
	Tags *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// 模板描述。
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 模板类型。默认值：kubernetes
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s CreateTemplateRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateRequest) GoString() string {
	return s.String()
}

func (s *CreateTemplateRequest) SetName(v string) *CreateTemplateRequest {
	s.Name = &v
	return s
}

func (s *CreateTemplateRequest) SetTemplate(v string) *CreateTemplateRequest {
	s.Template = &v
	return s
}

func (s *CreateTemplateRequest) SetTags(v string) *CreateTemplateRequest {
	s.Tags = &v
	return s
}

func (s *CreateTemplateRequest) SetDescription(v string) *CreateTemplateRequest {
	s.Description = &v
	return s
}

func (s *CreateTemplateRequest) SetTemplateType(v string) *CreateTemplateRequest {
	s.TemplateType = &v
	return s
}

type CreateTemplateResponseBody struct {
	// 模板ID。
	TemplateId *string `json:"template_id,omitempty" xml:"template_id,omitempty"`
}

func (s CreateTemplateResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateResponseBody) GoString() string {
	return s.String()
}

func (s *CreateTemplateResponseBody) SetTemplateId(v string) *CreateTemplateResponseBody {
	s.TemplateId = &v
	return s
}

type CreateTemplateResponse struct {
	Headers map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *CreateTemplateResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateResponse) GoString() string {
	return s.String()
}

func (s *CreateTemplateResponse) SetHeaders(v map[string]*string) *CreateTemplateResponse {
	s.Headers = v
	return s
}

func (s *CreateTemplateResponse) SetBody(v *CreateTemplateResponseBody) *CreateTemplateResponse {
	s.Body = v
	return s
}

type DescribeClusterNodesRequest struct {
	// 节点实例ID，按照实例ID进行过滤。  节点池ID不为空时会忽略此字段。多节点用逗号分割
	InstanceIds *string `json:"instanceIds,omitempty" xml:"instanceIds,omitempty"`
	// 节点池ID。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 节点状态。默认值：all。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 每页大小。
	PageSize *string `json:"pageSize,omitempty" xml:"pageSize,omitempty"`
	// 分页数量
	PageNumber *string `json:"pageNumber,omitempty" xml:"pageNumber,omitempty"`
}

func (s DescribeClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesRequest) SetInstanceIds(v string) *DescribeClusterNodesRequest {
	s.InstanceIds = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetNodepoolId(v string) *DescribeClusterNodesRequest {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetState(v string) *DescribeClusterNodesRequest {
	s.State = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetPageSize(v string) *DescribeClusterNodesRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetPageNumber(v string) *DescribeClusterNodesRequest {
	s.PageNumber = &v
	return s
}

type DescribeClusterNodesResponseBody struct {
	// 节点信息列表。
	Nodes []*DescribeClusterNodesResponseBodyNodes `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	// 分页信息。
	Page *DescribeClusterNodesResponseBodyPage `json:"page,omitempty" xml:"page,omitempty" type:"Struct"`
}

func (s DescribeClusterNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBody) SetNodes(v []*DescribeClusterNodesResponseBodyNodes) *DescribeClusterNodesResponseBody {
	s.Nodes = v
	return s
}

func (s *DescribeClusterNodesResponseBody) SetPage(v *DescribeClusterNodesResponseBodyPage) *DescribeClusterNodesResponseBody {
	s.Page = v
	return s
}

type DescribeClusterNodesResponseBodyNodes struct {
	// 节点创建时间。
	CreationTime *string `json:"creation_time,omitempty" xml:"creation_time,omitempty"`
	// 错误信息说明。
	ErrorMessage *string `json:"error_message,omitempty" xml:"error_message,omitempty"`
	// 节点过期时间。
	ExpiredTime *string `json:"expired_time,omitempty" xml:"expired_time,omitempty"`
	// 节点主机名。
	HostName *string `json:"host_name,omitempty" xml:"host_name,omitempty"`
	// 节点使用的镜像ID。
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 节点付费类型。
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// 节点实例ID。
	InstanceId *string `json:"instance_id,omitempty" xml:"instance_id,omitempty"`
	// 节点名称。
	InstanceName *string `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
	// 节点实例角色类型，Master或Worker。
	InstanceRole *string `json:"instance_role,omitempty" xml:"instance_role,omitempty"`
	// 节点实例状态，
	InstanceStatus *string `json:"instance_status,omitempty" xml:"instance_status,omitempty"`
	// 节点实例类型。
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// 节点实例所属ECS实例簇名称。
	InstanceTypeFamily *string `json:"instance_type_family,omitempty" xml:"instance_type_family,omitempty"`
	// 节点IP地址。
	IpAddress []*string `json:"ip_address,omitempty" xml:"ip_address,omitempty" type:"Repeated"`
	// 节点是否为aliyun实例。
	IsAliyunNode *bool `json:"is_aliyun_node,omitempty" xml:"is_aliyun_node,omitempty"`
	// 节点名称，该名称是k8s专用名称。
	NodeName *string `json:"node_name,omitempty" xml:"node_name,omitempty"`
	// 节点状态，是否Ready。
	NodeStatus *string `json:"node_status,omitempty" xml:"node_status,omitempty"`
	// 节点池ID。
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 节点通过什么方式创建出来的，例如：ROS。
	Source *string `json:"source,omitempty" xml:"source,omitempty"`
	// ECS运行状态，例如：Running。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 抢占时实例类型
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
}

func (s DescribeClusterNodesResponseBodyNodes) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBodyNodes) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBodyNodes) SetCreationTime(v string) *DescribeClusterNodesResponseBodyNodes {
	s.CreationTime = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetErrorMessage(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ErrorMessage = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetExpiredTime(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ExpiredTime = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetHostName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.HostName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetImageId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceChargeType(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceRole(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceRole = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceStatus(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceStatus = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceType(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceTypeFamily(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceTypeFamily = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetIpAddress(v []*string) *DescribeClusterNodesResponseBodyNodes {
	s.IpAddress = v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetIsAliyunNode(v bool) *DescribeClusterNodesResponseBodyNodes {
	s.IsAliyunNode = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodeName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodeName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodeStatus(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodeStatus = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodepoolId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetSource(v string) *DescribeClusterNodesResponseBodyNodes {
	s.Source = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetState(v string) *DescribeClusterNodesResponseBodyNodes {
	s.State = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetSpotStrategy(v string) *DescribeClusterNodesResponseBodyNodes {
	s.SpotStrategy = &v
	return s
}

type DescribeClusterNodesResponseBodyPage struct {
	// 总页数。
	PageNumber *int32 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// 单页展示结果数量。
	PageSize *int32 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 结果总条数。
	TotalCount *int32 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClusterNodesResponseBodyPage) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBodyPage) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBodyPage) SetPageNumber(v int32) *DescribeClusterNodesResponseBodyPage {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyPage) SetPageSize(v int32) *DescribeClusterNodesResponseBodyPage {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyPage) SetTotalCount(v int32) *DescribeClusterNodesResponseBodyPage {
	s.TotalCount = &v
	return s
}

type DescribeClusterNodesResponse struct {
	Headers map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponse) SetHeaders(v map[string]*string) *DescribeClusterNodesResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodesResponse) SetBody(v *DescribeClusterNodesResponseBody) *DescribeClusterNodesResponse {
	s.Body = v
	return s
}

type DeleteClusterRequest struct {
	// 是否保留所有资源,如果设置了该值，将会忽略retain_resources。  true：保留 false：不保留 默认值：fase。
	RetainAllResources *bool `json:"retain_all_resources,omitempty" xml:"retain_all_resources,omitempty"`
	// 是否保留SLB。  true：保留 false：不保留 默认值：false。
	KeepSlb *bool `json:"keep_slb,omitempty" xml:"keep_slb,omitempty"`
	// 要保留的资源列表。
	RetainResources []*string `json:"retain_resources,omitempty" xml:"retain_resources,omitempty" type:"Repeated"`
}

func (s DeleteClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterRequest) SetRetainAllResources(v bool) *DeleteClusterRequest {
	s.RetainAllResources = &v
	return s
}

func (s *DeleteClusterRequest) SetKeepSlb(v bool) *DeleteClusterRequest {
	s.KeepSlb = &v
	return s
}

func (s *DeleteClusterRequest) SetRetainResources(v []*string) *DeleteClusterRequest {
	s.RetainResources = v
	return s
}

type DeleteClusterShrinkRequest struct {
	// 是否保留所有资源,如果设置了该值，将会忽略retain_resources。  true：保留 false：不保留 默认值：fase。
	RetainAllResources *bool `json:"retain_all_resources,omitempty" xml:"retain_all_resources,omitempty"`
	// 是否保留SLB。  true：保留 false：不保留 默认值：false。
	KeepSlb *bool `json:"keep_slb,omitempty" xml:"keep_slb,omitempty"`
	// 要保留的资源列表。
	RetainResourcesShrink *string `json:"retain_resources,omitempty" xml:"retain_resources,omitempty"`
}

func (s DeleteClusterShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterShrinkRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterShrinkRequest) SetRetainAllResources(v bool) *DeleteClusterShrinkRequest {
	s.RetainAllResources = &v
	return s
}

func (s *DeleteClusterShrinkRequest) SetKeepSlb(v bool) *DeleteClusterShrinkRequest {
	s.KeepSlb = &v
	return s
}

func (s *DeleteClusterShrinkRequest) SetRetainResourcesShrink(v string) *DeleteClusterShrinkRequest {
	s.RetainResourcesShrink = &v
	return s
}

type DeleteClusterResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s DeleteClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterResponse) SetHeaders(v map[string]*string) *DeleteClusterResponse {
	s.Headers = v
	return s
}

type CancelComponentUpgradeResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s CancelComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *CancelComponentUpgradeResponse) SetHeaders(v map[string]*string) *CancelComponentUpgradeResponse {
	s.Headers = v
	return s
}

type MigrateClusterResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s MigrateClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s MigrateClusterResponse) GoString() string {
	return s.String()
}

func (s *MigrateClusterResponse) SetHeaders(v map[string]*string) *MigrateClusterResponse {
	s.Headers = v
	return s
}

type DescribeClusterAddonsVersionResponse struct {
	Headers map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonsVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsVersionResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsVersionResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonsVersionResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonsVersionResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonsVersionResponse {
	s.Body = v
	return s
}

type DescribeExternalAgentRequest struct {
	// 是否获取内网访问凭据。  true：获取内网连接凭据 false：获取公网连接凭据 默认值：false。
	PrivateIpAddress *string `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
}

func (s DescribeExternalAgentRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentRequest) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentRequest) SetPrivateIpAddress(v string) *DescribeExternalAgentRequest {
	s.PrivateIpAddress = &v
	return s
}

type DescribeExternalAgentResponseBody struct {
	// 代理配置。
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s DescribeExternalAgentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentResponseBody) SetConfig(v string) *DescribeExternalAgentResponseBody {
	s.Config = &v
	return s
}

type DescribeExternalAgentResponse struct {
	Headers map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeExternalAgentResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeExternalAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentResponse) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentResponse) SetHeaders(v map[string]*string) *DescribeExternalAgentResponse {
	s.Headers = v
	return s
}

func (s *DescribeExternalAgentResponse) SetBody(v *DescribeExternalAgentResponseBody) *DescribeExternalAgentResponse {
	s.Body = v
	return s
}

type UnInstallClusterAddonsRequest struct {
	// 卸载组件列表。
	Addons []*UnInstallClusterAddonsRequestAddons `json:"addons,omitempty" xml:"addons,omitempty" type:"Repeated"`
}

func (s UnInstallClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsRequest) SetAddons(v []*UnInstallClusterAddonsRequestAddons) *UnInstallClusterAddonsRequest {
	s.Addons = v
	return s
}

type UnInstallClusterAddonsRequestAddons struct {
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s UnInstallClusterAddonsRequestAddons) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsRequestAddons) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsRequestAddons) SetName(v string) *UnInstallClusterAddonsRequestAddons {
	s.Name = &v
	return s
}

type UnInstallClusterAddonsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UnInstallClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsResponse) SetHeaders(v map[string]*string) *UnInstallClusterAddonsResponse {
	s.Headers = v
	return s
}

type ResumeComponentUpgradeResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s ResumeComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s ResumeComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *ResumeComponentUpgradeResponse) SetHeaders(v map[string]*string) *ResumeComponentUpgradeResponse {
	s.Headers = v
	return s
}

type DescribeClustersV1Request struct {
	// 通过集群名称进行模糊查询。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群类型。  Kubernetes: 专有版集群。 ManagedKubernetes：托管版集群。 Ask：Serverless集群。 ExternalKubernetes：注册集群。 ServiceMesh：ASM集群。
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// 单页大小。
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 分页数。
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
}

func (s DescribeClustersV1Request) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1Request) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1Request) SetName(v string) *DescribeClustersV1Request {
	s.Name = &v
	return s
}

func (s *DescribeClustersV1Request) SetClusterType(v string) *DescribeClustersV1Request {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersV1Request) SetPageSize(v int64) *DescribeClustersV1Request {
	s.PageSize = &v
	return s
}

func (s *DescribeClustersV1Request) SetPageNumber(v int64) *DescribeClustersV1Request {
	s.PageNumber = &v
	return s
}

type DescribeClustersV1ResponseBody struct {
	// 集群详情列表。
	Clusters []*DescribeClustersV1ResponseBodyClusters `json:"clusters,omitempty" xml:"clusters,omitempty" type:"Repeated"`
	// 分页信息。
	PageInfo *DescribeClustersV1ResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeClustersV1ResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBody) SetClusters(v []*DescribeClustersV1ResponseBodyClusters) *DescribeClustersV1ResponseBody {
	s.Clusters = v
	return s
}

func (s *DescribeClustersV1ResponseBody) SetPageInfo(v *DescribeClustersV1ResponseBodyPageInfo) *DescribeClustersV1ResponseBody {
	s.PageInfo = v
	return s
}

type DescribeClustersV1ResponseBodyClusters struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 集群类型。
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// 集群初始化时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 集群初始化版本。
	InitVersion *string `json:"init_version,omitempty" xml:"init_version,omitempty"`
	// 集群当前版本。
	CurrentVersion *string `json:"current_version,omitempty" xml:"current_version,omitempty"`
	// 集群可升级版本。
	NextVersion *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// 集群是否开启删除保护。
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// 集群使用的Docker版本。
	DockerVersion *string `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	// 集群负载均衡服务的ID。
	ExternalLoadbalancerId *string `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	// 集群访问地址列表。
	MasterUrl *string `json:"master_url,omitempty" xml:"master_url,omitempty"`
	// 集群元数据信息。
	MetaData *string `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// 集群名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群使用的网络类型，例如：VPC网络。
	NetworkMode *string `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	// 集群是否开启Private Zone。
	PrivateZone *bool `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	// 边缘集群表示，用于区分边缘托管版集群。
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// 地域ID。
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 集群资源组ID。
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// 集群安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 集群节点数。
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// 集群运行状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// POD网段地址。
	SubnetCidr *string `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	// 集群标签。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 集群更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// 集群所在的VPC ID。
	VpcId *string `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	// 集群使用的虚拟交换ID。
	VswitchId *string `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	// 集群Worker RAM角色。
	WorkerRamRoleName *string `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	// 可用区ID。
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
	// 托管版集群类型，面向托管集群。 • ack.pro.small：专业托管集群。 • ack.standard ：标准托管集群。
	ClusterSpec       *string            `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
}

func (s DescribeClustersV1ResponseBodyClusters) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBodyClusters) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterType(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetCreated(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Created = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetInitVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.InitVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetCurrentVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetNextVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.NextVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetDeletionProtection(v bool) *DescribeClustersV1ResponseBodyClusters {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetDockerVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetExternalLoadbalancerId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMasterUrl(v string) *DescribeClustersV1ResponseBodyClusters {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMetaData(v string) *DescribeClustersV1ResponseBodyClusters {
	s.MetaData = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetName(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Name = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetNetworkMode(v string) *DescribeClustersV1ResponseBodyClusters {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetPrivateZone(v bool) *DescribeClustersV1ResponseBodyClusters {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetProfile(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Profile = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetRegionId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.RegionId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetResourceGroupId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSecurityGroupId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSize(v int64) *DescribeClustersV1ResponseBodyClusters {
	s.Size = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetState(v string) *DescribeClustersV1ResponseBodyClusters {
	s.State = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSubnetCidr(v string) *DescribeClustersV1ResponseBodyClusters {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetTags(v []*Tag) *DescribeClustersV1ResponseBodyClusters {
	s.Tags = v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetUpdated(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Updated = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetVpcId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.VpcId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetVswitchId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.VswitchId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetWorkerRamRoleName(v string) *DescribeClustersV1ResponseBodyClusters {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetZoneId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ZoneId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterSpec(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMaintenanceWindow(v *MaintenanceWindow) *DescribeClustersV1ResponseBodyClusters {
	s.MaintenanceWindow = v
	return s
}

type DescribeClustersV1ResponseBodyPageInfo struct {
	// 分页数。
	PageNumber *int32 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// 单页大小。
	PageSize *int32 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 结果总数。
	TotalCount *int32 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClustersV1ResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetPageNumber(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetPageSize(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetTotalCount(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeClustersV1Response struct {
	Headers map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClustersV1ResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClustersV1Response) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1Response) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1Response) SetHeaders(v map[string]*string) *DescribeClustersV1Response {
	s.Headers = v
	return s
}

func (s *DescribeClustersV1Response) SetBody(v *DescribeClustersV1ResponseBody) *DescribeClustersV1Response {
	s.Body = v
	return s
}

type ModifyClusterConfigurationRequest struct {
	// 自定义配置。
	CustomizeConfig []*ModifyClusterConfigurationRequestCustomizeConfig `json:"customize_config,omitempty" xml:"customize_config,omitempty" type:"Repeated"`
}

func (s ModifyClusterConfigurationRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequest) SetCustomizeConfig(v []*ModifyClusterConfigurationRequestCustomizeConfig) *ModifyClusterConfigurationRequest {
	s.CustomizeConfig = v
	return s
}

type ModifyClusterConfigurationRequestCustomizeConfig struct {
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 组件配置。
	Configs []*ModifyClusterConfigurationRequestCustomizeConfigConfigs `json:"configs,omitempty" xml:"configs,omitempty" type:"Repeated"`
}

func (s ModifyClusterConfigurationRequestCustomizeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequestCustomizeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequestCustomizeConfig) SetName(v string) *ModifyClusterConfigurationRequestCustomizeConfig {
	s.Name = &v
	return s
}

func (s *ModifyClusterConfigurationRequestCustomizeConfig) SetConfigs(v []*ModifyClusterConfigurationRequestCustomizeConfigConfigs) *ModifyClusterConfigurationRequestCustomizeConfig {
	s.Configs = v
	return s
}

type ModifyClusterConfigurationRequestCustomizeConfigConfigs struct {
	// key值。
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
	// value值。
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ModifyClusterConfigurationRequestCustomizeConfigConfigs) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequestCustomizeConfigConfigs) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequestCustomizeConfigConfigs) SetKey(v string) *ModifyClusterConfigurationRequestCustomizeConfigConfigs {
	s.Key = &v
	return s
}

func (s *ModifyClusterConfigurationRequestCustomizeConfigConfigs) SetValue(v string) *ModifyClusterConfigurationRequestCustomizeConfigConfigs {
	s.Value = &v
	return s
}

type ModifyClusterConfigurationResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s ModifyClusterConfigurationResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationResponse) SetHeaders(v map[string]*string) *ModifyClusterConfigurationResponse {
	s.Headers = v
	return s
}

type DescribeTaskInfoResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
	// 任务创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 任务更新时间。
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// 任务当前状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 当前任务类型。
	TaskType *string `json:"task_type,omitempty" xml:"task_type,omitempty"`
	// 任务执行详情。
	TaskResult []*DescribeTaskInfoResponseBodyTaskResult `json:"task_result,omitempty" xml:"task_result,omitempty" type:"Repeated"`
}

func (s DescribeTaskInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBody) SetClusterId(v string) *DescribeTaskInfoResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskId(v string) *DescribeTaskInfoResponseBody {
	s.TaskId = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetCreated(v string) *DescribeTaskInfoResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetUpdated(v string) *DescribeTaskInfoResponseBody {
	s.Updated = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetState(v string) *DescribeTaskInfoResponseBody {
	s.State = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskType(v string) *DescribeTaskInfoResponseBody {
	s.TaskType = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskResult(v []*DescribeTaskInfoResponseBodyTaskResult) *DescribeTaskInfoResponseBody {
	s.TaskResult = v
	return s
}

type DescribeTaskInfoResponseBodyTaskResult struct {
	// 操作的资源，例如：实例ID。
	Data *string `json:"data,omitempty" xml:"data,omitempty"`
	// 资源的状态。
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s DescribeTaskInfoResponseBodyTaskResult) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyTaskResult) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyTaskResult) SetData(v string) *DescribeTaskInfoResponseBodyTaskResult {
	s.Data = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyTaskResult) SetStatus(v string) *DescribeTaskInfoResponseBodyTaskResult {
	s.Status = &v
	return s
}

type DescribeTaskInfoResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeTaskInfoResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeTaskInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponse) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponse) SetHeaders(v map[string]*string) *DescribeTaskInfoResponse {
	s.Headers = v
	return s
}

func (s *DescribeTaskInfoResponse) SetBody(v *DescribeTaskInfoResponseBody) *DescribeTaskInfoResponse {
	s.Body = v
	return s
}

type DescirbeWorkflowResponseBody struct {
	// 工作流创建时间。
	CreateTime *string `json:"create_time,omitempty" xml:"create_time,omitempty"`
	// 工作流经过时长。
	Duration *string `json:"duration,omitempty" xml:"duration,omitempty"`
	// 任务结束时间。
	FinishTime *string `json:"finish_time,omitempty" xml:"finish_time,omitempty"`
	// 输入数据大小。
	InputDataSize *string `json:"input_data_size,omitempty" xml:"input_data_size,omitempty"`
	// 工作流名称。
	JobName *string `json:"job_name,omitempty" xml:"job_name,omitempty"`
	// 工作流所在命名空间。
	JobNamespace *string `json:"job_namespace,omitempty" xml:"job_namespace,omitempty"`
	// 输出数据大小。
	OutputDataSize *string `json:"output_data_size,omitempty" xml:"output_data_size,omitempty"`
	// 工作流当前状态。
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 碱基对个数。
	TotalBases *string `json:"total_bases,omitempty" xml:"total_bases,omitempty"`
	// Reads个数。
	TotalReads *string `json:"total_reads,omitempty" xml:"total_reads,omitempty"`
	// 用户输入参数。
	UserInputData *string `json:"user_input_data,omitempty" xml:"user_input_data,omitempty"`
}

func (s DescirbeWorkflowResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescirbeWorkflowResponseBody) GoString() string {
	return s.String()
}

func (s *DescirbeWorkflowResponseBody) SetCreateTime(v string) *DescirbeWorkflowResponseBody {
	s.CreateTime = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetDuration(v string) *DescirbeWorkflowResponseBody {
	s.Duration = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetFinishTime(v string) *DescirbeWorkflowResponseBody {
	s.FinishTime = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetInputDataSize(v string) *DescirbeWorkflowResponseBody {
	s.InputDataSize = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetJobName(v string) *DescirbeWorkflowResponseBody {
	s.JobName = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetJobNamespace(v string) *DescirbeWorkflowResponseBody {
	s.JobNamespace = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetOutputDataSize(v string) *DescirbeWorkflowResponseBody {
	s.OutputDataSize = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetStatus(v string) *DescirbeWorkflowResponseBody {
	s.Status = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetTotalBases(v string) *DescirbeWorkflowResponseBody {
	s.TotalBases = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetTotalReads(v string) *DescirbeWorkflowResponseBody {
	s.TotalReads = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetUserInputData(v string) *DescirbeWorkflowResponseBody {
	s.UserInputData = &v
	return s
}

type DescirbeWorkflowResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescirbeWorkflowResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescirbeWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s DescirbeWorkflowResponse) GoString() string {
	return s.String()
}

func (s *DescirbeWorkflowResponse) SetHeaders(v map[string]*string) *DescirbeWorkflowResponse {
	s.Headers = v
	return s
}

func (s *DescirbeWorkflowResponse) SetBody(v *DescirbeWorkflowResponseBody) *DescirbeWorkflowResponse {
	s.Body = v
	return s
}

type CancelClusterUpgradeResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s CancelClusterUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelClusterUpgradeResponse) GoString() string {
	return s.String()
}

func (s *CancelClusterUpgradeResponse) SetHeaders(v map[string]*string) *CancelClusterUpgradeResponse {
	s.Headers = v
	return s
}

type RemoveWorkflowResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s RemoveWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveWorkflowResponse) GoString() string {
	return s.String()
}

func (s *RemoveWorkflowResponse) SetHeaders(v map[string]*string) *RemoveWorkflowResponse {
	s.Headers = v
	return s
}

type UpdateTemplateRequest struct {
	// 部署模板描述信息。
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 部署模板名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 部署模板标签
	Tags *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// 部署模板yaml。
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
	// 部署模板类型。
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s UpdateTemplateRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateTemplateRequest) GoString() string {
	return s.String()
}

func (s *UpdateTemplateRequest) SetDescription(v string) *UpdateTemplateRequest {
	s.Description = &v
	return s
}

func (s *UpdateTemplateRequest) SetName(v string) *UpdateTemplateRequest {
	s.Name = &v
	return s
}

func (s *UpdateTemplateRequest) SetTags(v string) *UpdateTemplateRequest {
	s.Tags = &v
	return s
}

func (s *UpdateTemplateRequest) SetTemplate(v string) *UpdateTemplateRequest {
	s.Template = &v
	return s
}

func (s *UpdateTemplateRequest) SetTemplateType(v string) *UpdateTemplateRequest {
	s.TemplateType = &v
	return s
}

type UpdateTemplateResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UpdateTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateTemplateResponse) GoString() string {
	return s.String()
}

func (s *UpdateTemplateResponse) SetHeaders(v map[string]*string) *UpdateTemplateResponse {
	s.Headers = v
	return s
}

type UpgradeClusterAddonsRequest struct {
	// Request body，类型是对象数组。
	Body []*UpgradeClusterAddonsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s UpgradeClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsRequest) SetBody(v []*UpgradeClusterAddonsRequestBody) *UpgradeClusterAddonsRequest {
	s.Body = v
	return s
}

type UpgradeClusterAddonsRequestBody struct {
	// 组件名称。
	ComponentName *string `json:"component_name,omitempty" xml:"component_name,omitempty"`
	// 可升级版本。
	NextVersion *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// 当前版本。
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s UpgradeClusterAddonsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsRequestBody) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsRequestBody) SetComponentName(v string) *UpgradeClusterAddonsRequestBody {
	s.ComponentName = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetNextVersion(v string) *UpgradeClusterAddonsRequestBody {
	s.NextVersion = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetVersion(v string) *UpgradeClusterAddonsRequestBody {
	s.Version = &v
	return s
}

type UpgradeClusterAddonsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UpgradeClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsResponse) SetHeaders(v map[string]*string) *UpgradeClusterAddonsResponse {
	s.Headers = v
	return s
}

type DescribeClusterNamespacesResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*string          `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClusterNamespacesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNamespacesResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNamespacesResponse) SetHeaders(v map[string]*string) *DescribeClusterNamespacesResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNamespacesResponse) SetBody(v []*string) *DescribeClusterNamespacesResponse {
	s.Body = v
	return s
}

type DeleteKubernetesTriggerResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s DeleteKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *DeleteKubernetesTriggerResponse) SetHeaders(v map[string]*string) *DeleteKubernetesTriggerResponse {
	s.Headers = v
	return s
}

type DescribeUserQuotaResponseBody struct {
	// 托管版集群配额。
	AmkClusterQuota *int64 `json:"amk_cluster_quota,omitempty" xml:"amk_cluster_quota,omitempty"`
	// Serverless集群配额。
	AskClusterQuota *int64 `json:"ask_cluster_quota,omitempty" xml:"ask_cluster_quota,omitempty"`
	// 集群节点池配额。
	ClusterNodepoolQuota *int64 `json:"cluster_nodepool_quota,omitempty" xml:"cluster_nodepool_quota,omitempty"`
	// 专有版集群托管版集群的总配额。
	ClusterQuota *int64 `json:"cluster_quota,omitempty" xml:"cluster_quota,omitempty"`
	// 单集群的节点配额。
	NodeQuota *int64 `json:"node_quota,omitempty" xml:"node_quota,omitempty"`
}

func (s DescribeUserQuotaResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserQuotaResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeUserQuotaResponseBody) SetAmkClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.AmkClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetAskClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.AskClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetClusterNodepoolQuota(v int64) *DescribeUserQuotaResponseBody {
	s.ClusterNodepoolQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.ClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetNodeQuota(v int64) *DescribeUserQuotaResponseBody {
	s.NodeQuota = &v
	return s
}

type DescribeUserQuotaResponse struct {
	Headers map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeUserQuotaResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeUserQuotaResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserQuotaResponse) GoString() string {
	return s.String()
}

func (s *DescribeUserQuotaResponse) SetHeaders(v map[string]*string) *DescribeUserQuotaResponse {
	s.Headers = v
	return s
}

func (s *DescribeUserQuotaResponse) SetBody(v *DescribeUserQuotaResponseBody) *DescribeUserQuotaResponse {
	s.Body = v
	return s
}

type DeleteClusterNodepoolResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s DeleteClusterNodepoolResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodepoolResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodepoolResponse) SetHeaders(v map[string]*string) *DeleteClusterNodepoolResponse {
	s.Headers = v
	return s
}

type DescribeClusterAddonsUpgradeStatusRequest struct {
	// 组件名称列表。
	ComponentIds []*string `json:"componentIds,omitempty" xml:"componentIds,omitempty" type:"Repeated"`
}

func (s DescribeClusterAddonsUpgradeStatusRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusRequest) SetComponentIds(v []*string) *DescribeClusterAddonsUpgradeStatusRequest {
	s.ComponentIds = v
	return s
}

type DescribeClusterAddonsUpgradeStatusShrinkRequest struct {
	// 组件名称列表。
	ComponentIdsShrink *string `json:"componentIds,omitempty" xml:"componentIds,omitempty"`
}

func (s DescribeClusterAddonsUpgradeStatusShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusShrinkRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusShrinkRequest) SetComponentIdsShrink(v string) *DescribeClusterAddonsUpgradeStatusShrinkRequest {
	s.ComponentIdsShrink = &v
	return s
}

type DescribeClusterAddonsUpgradeStatusResponse struct {
	Headers map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonsUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonsUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonsUpgradeStatusResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonsUpgradeStatusResponse {
	s.Body = v
	return s
}

type DescribeWorkflowsResponseBody struct {
	// job信息
	Jobs []*DescribeWorkflowsResponseBodyJobs `json:"jobs,omitempty" xml:"jobs,omitempty" type:"Repeated"`
}

func (s DescribeWorkflowsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponseBody) SetJobs(v []*DescribeWorkflowsResponseBodyJobs) *DescribeWorkflowsResponseBody {
	s.Jobs = v
	return s
}

type DescribeWorkflowsResponseBodyJobs struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 工作流名称。
	JobName *string `json:"job_name,omitempty" xml:"job_name,omitempty"`
	// 工作流创建时间。
	CreateTime *string `json:"create_time,omitempty" xml:"create_time,omitempty"`
}

func (s DescribeWorkflowsResponseBodyJobs) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponseBodyJobs) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponseBodyJobs) SetClusterId(v string) *DescribeWorkflowsResponseBodyJobs {
	s.ClusterId = &v
	return s
}

func (s *DescribeWorkflowsResponseBodyJobs) SetJobName(v string) *DescribeWorkflowsResponseBodyJobs {
	s.JobName = &v
	return s
}

func (s *DescribeWorkflowsResponseBodyJobs) SetCreateTime(v string) *DescribeWorkflowsResponseBodyJobs {
	s.CreateTime = &v
	return s
}

type DescribeWorkflowsResponse struct {
	Headers map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeWorkflowsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeWorkflowsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponse) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponse) SetHeaders(v map[string]*string) *DescribeWorkflowsResponse {
	s.Headers = v
	return s
}

func (s *DescribeWorkflowsResponse) SetBody(v *DescribeWorkflowsResponseBody) *DescribeWorkflowsResponse {
	s.Body = v
	return s
}

type InstallClusterAddonsRequest struct {
	// Addon列表。
	Body []*InstallClusterAddonsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s InstallClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsRequest) SetBody(v []*InstallClusterAddonsRequestBody) *InstallClusterAddonsRequest {
	s.Body = v
	return s
}

type InstallClusterAddonsRequestBody struct {
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 组件版本号。
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
	// 组件配置信息。
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s InstallClusterAddonsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsRequestBody) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsRequestBody) SetName(v string) *InstallClusterAddonsRequestBody {
	s.Name = &v
	return s
}

func (s *InstallClusterAddonsRequestBody) SetVersion(v string) *InstallClusterAddonsRequestBody {
	s.Version = &v
	return s
}

func (s *InstallClusterAddonsRequestBody) SetConfig(v string) *InstallClusterAddonsRequestBody {
	s.Config = &v
	return s
}

type InstallClusterAddonsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s InstallClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsResponse) SetHeaders(v map[string]*string) *InstallClusterAddonsResponse {
	s.Headers = v
	return s
}

type DescribeClusterNodePoolsResponseBody struct {
	// 节点池列表。
	Nodepools []*DescribeClusterNodePoolsResponseBodyNodepools `json:"nodepools,omitempty" xml:"nodepools,omitempty" type:"Repeated"`
}

func (s DescribeClusterNodePoolsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBody) SetNodepools(v []*DescribeClusterNodePoolsResponseBodyNodepools) *DescribeClusterNodePoolsResponseBody {
	s.Nodepools = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepools struct {
	// 自动伸缩配置详情。
	AutoScaling *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	// 集群配置信息。
	KubernetesConfig *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// 节点池配置详情。
	NodepoolInfo *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// 扩容组配置详情。
	ScalingGroup *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	// 节点池状态详情。
	Status *DescribeClusterNodePoolsResponseBodyNodepoolsStatus `json:"status,omitempty" xml:"status,omitempty" type:"Struct"`
	// 加密计算配置详情。
	TeeConfig *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
	// 托管节点池配置。
	Management *DescribeClusterNodePoolsResponseBodyNodepoolsManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepools) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepools) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetAutoScaling(v *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.AutoScaling = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetKubernetesConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.KubernetesConfig = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetNodepoolInfo(v *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.NodepoolInfo = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetScalingGroup(v *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.ScalingGroup = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetStatus(v *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.Status = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetTeeConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.TeeConfig = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetManagement(v *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.Management = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling struct {
	// EIP带宽峰值
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// 是否绑定EIP
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// EIP实例计费方式
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// 自动伸缩。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 最大节点数
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// 最小节点数
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// 扩容组类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEipBandwidth(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetIsBondEip(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEipInternetChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetMaxInstances(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetMinInstances(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.Type = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig struct {
	// 是否开启云监控
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// CPU管理策略
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// ECS标签。
	Labels []*Tag `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	// 容器运行时
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// 容器运行时版本
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// 污点配置
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// 节点自定义数据
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetCmsEnabled(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetCpuPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetLabels(v []*Tag) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Labels = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetRuntime(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetRuntimeVersion(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetTaints(v []*Taint) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Taints = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetUserData(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.UserData = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo struct {
	// 节点池创建时间
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 是否为默认节点池
	IsDefault *bool `json:"is_default,omitempty" xml:"is_default,omitempty"`
	// 节点池名称
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 节点池ID
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// 节点池所在地域ID。
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 资源组ID
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// 节点池类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 节点池更新时间
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetCreated(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Created = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetIsDefault(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.IsDefault = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetName(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Name = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetNodepoolId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetRegionId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetResourceGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Type = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetUpdated(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Updated = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup struct {
	// 自动续费
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// 自动付费时长
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// 数据盘配置
	DataDisks []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	// 镜像ID
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 节点付费类型
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// 节点类型
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// 多可用区伸缩组ECS实例扩缩容策略
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// 伸缩组所需要按量实例个数的最小值，取值范围：0~1000。当按量实例个数少于该值时，将优先创建按量实例。
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// 伸缩组满足最小按量实例数（OnDemandBaseCapacity）要求后，超出的实例中按量实例应占的比例，取值范围：0～100。
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// 指定可用实例规格的个数，伸缩组将按成本最低的多个规格均衡创建抢占式实例。取值范围：1~10。
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// 是否开启补齐抢占式实例。开启后，当收到抢占式实例将被回收的系统消息时，伸缩组将尝试创建新的实例，替换掉将被回收的抢占式实例。
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// 当MultiAZPolicy取值为COST_OPTIMIZED时，如果因价格、库存等原因无法创建足够的抢占式实例，是否允许自动尝试创建按量实例满足ECS实例数量要求。取值范围：true：允许。false：不允许。默认值：true
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// 包年包月时长
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// 自动付费周期
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// 操作系统发行版。取值： CentOS，AliyunLinux，Windows，WindowsCore。
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// RAM 角色名称
	RamPolicy *string `json:"ram_policy,omitempty" xml:"ram_policy,omitempty"`
	// 抢占式实例类型
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// 抢占实例价格上限配置。
	SpotPriceLimit []*DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// RDS列表
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 扩容组ID
	ScalingGroupId *string `json:"scaling_group_id,omitempty" xml:"scaling_group_id,omitempty"`
	// 扩容节点策略
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// 安全组ID。
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// 系统盘类型。
	SystemDiskCategory *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	// 系统盘大小
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// 节点标签
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 虚拟交换机ID。
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// 登录密码，返回结果是加密的。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 密钥对名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// 节点公网IP网络计费类型
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// 节点公网IP出带宽最大值，单位为Mbps（Mega bit per second），取值范围：1~100
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetAutoRenew(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetAutoRenewPeriod(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetDataDisks(v []*DataDisk) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.DataDisks = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetImageId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInstanceChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInstanceTypes(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetMultiAzPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetOnDemandBaseCapacity(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotInstancePools(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotInstanceRemedy(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetCompensateWithOnDemand(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPeriod(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Period = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPeriodUnit(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPlatform(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Platform = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetRamPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.RamPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotStrategy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotPriceLimit(v []*DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetRdsInstances(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetScalingGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ScalingGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetScalingPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSecurityGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSystemDiskCategory(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSystemDiskSize(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetTags(v []*Tag) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Tags = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetVswitchIds(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.VswitchIds = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetLoginPassword(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetKeyPair(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInternetChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInternetMaxBandwidthOut(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit struct {
	// 抢占式实例规格。
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// 单台实例上限价格，单位：元/小时。
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) SetInstanceType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) SetPriceLimit(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsStatus struct {
	// 失败的节点数
	FailedNodes *int64 `json:"failed_nodes,omitempty" xml:"failed_nodes,omitempty"`
	// 处于健康状态的节点数
	HealthyNodes *int64 `json:"healthy_nodes,omitempty" xml:"healthy_nodes,omitempty"`
	// 正在创建的节点数
	InitialNodes *int64 `json:"initial_nodes,omitempty" xml:"initial_nodes,omitempty"`
	// 离线节点数
	OfflineNodes *int64 `json:"offline_nodes,omitempty" xml:"offline_nodes,omitempty"`
	// 真在被移除的节点数。
	RemovingNodes *int64 `json:"removing_nodes,omitempty" xml:"removing_nodes,omitempty"`
	// 正在工作节点数
	ServingNodes *int64 `json:"serving_nodes,omitempty" xml:"serving_nodes,omitempty"`
	// 节点池状态
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// 节点总数
	TotalNodes *int64 `json:"total_nodes,omitempty" xml:"total_nodes,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsStatus) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsStatus) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetFailedNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.FailedNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetHealthyNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.HealthyNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetInitialNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.InitialNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetOfflineNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.OfflineNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetRemovingNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.RemovingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetServingNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.ServingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetState(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.State = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetTotalNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.TotalNodes = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig struct {
	// 是否为加密计算节点池
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) SetTeeEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig {
	s.TeeEnable = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsManagement struct {
	// 是否开启托管版节点池。
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// 是否启用自动修复。
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// 是否启用自动修复。
	UpgradeConfig *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagement) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagement) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetAutoRepair(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.AutoRepair = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetUpgradeConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.UpgradeConfig = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig struct {
	// 是否启用自动升级，自修复。
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// 额外节点数量。
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// 额外节点比例， 和surge 二选一。
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
	// 最大不可用节点数量。
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetAutoUpgrade(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetSurge(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetSurgePercentage(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetMaxUnavailable(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

type DescribeClusterNodePoolsResponse struct {
	Headers map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterNodePoolsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodePoolsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponse) SetHeaders(v map[string]*string) *DescribeClusterNodePoolsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodePoolsResponse) SetBody(v *DescribeClusterNodePoolsResponseBody) *DescribeClusterNodePoolsResponse {
	s.Body = v
	return s
}

type DescribeClusterV2UserKubeconfigRequest struct {
	// 是否为内网访问。
	PrivateIpAddress *bool `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
}

func (s DescribeClusterV2UserKubeconfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigRequest) SetPrivateIpAddress(v bool) *DescribeClusterV2UserKubeconfigRequest {
	s.PrivateIpAddress = &v
	return s
}

type DescribeClusterV2UserKubeconfigResponseBody struct {
	// kubeconfig内容。
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s DescribeClusterV2UserKubeconfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigResponseBody) SetConfig(v string) *DescribeClusterV2UserKubeconfigResponseBody {
	s.Config = &v
	return s
}

type DescribeClusterV2UserKubeconfigResponse struct {
	Headers map[string]*string                           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeClusterV2UserKubeconfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterV2UserKubeconfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigResponse) SetHeaders(v map[string]*string) *DescribeClusterV2UserKubeconfigResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterV2UserKubeconfigResponse) SetBody(v *DescribeClusterV2UserKubeconfigResponseBody) *DescribeClusterV2UserKubeconfigResponse {
	s.Body = v
	return s
}

type StartWorkflowRequest struct {
	// 工作流类型，可选值：wgs或mapping。
	WorkflowType *string `json:"workflow_type,omitempty" xml:"workflow_type,omitempty"`
	// SLA类型，可选值：s、g、p。 白银级（s）：超过90 Gbp的部分，按1.5 Gbp/min计算增加的时间。 黄金级（g）：超过90 Gbp的部分，按2 Gbp/min计算增加的时间。 铂金级（p）：超过90 Gbp的部分，按3 Gbp/min计算增加的时间。
	Service *string `json:"service,omitempty" xml:"service,omitempty"`
	// mapping oss数据的存放region。
	MappingOssRegion *string `json:"mapping_oss_region,omitempty" xml:"mapping_oss_region,omitempty"`
	// mapping的第一个fastq文件名。
	MappingFastqFirstFilename *string `json:"mapping_fastq_first_filename,omitempty" xml:"mapping_fastq_first_filename,omitempty"`
	// mapping的第二个fastq文件名。
	MappingFastqSecondFilename *string `json:"mapping_fastq_second_filename,omitempty" xml:"mapping_fastq_second_filename,omitempty"`
	// 存放mapping的bucket名称。
	MappingBucketName *string `json:"mapping_bucket_name,omitempty" xml:"mapping_bucket_name,omitempty"`
	// mapping的fastq文件路径。
	MappingFastqPath *string `json:"mapping_fastq_path,omitempty" xml:"mapping_fastq_path,omitempty"`
	// mapping的reference文件位置。
	MappingReferencePath *string `json:"mapping_reference_path,omitempty" xml:"mapping_reference_path,omitempty"`
	// 是否进行dup。
	MappingIsMarkDup *string `json:"mapping_is_mark_dup,omitempty" xml:"mapping_is_mark_dup,omitempty"`
	// bam文件输出路径。
	MappingBamOutPath *string `json:"mapping_bam_out_path,omitempty" xml:"mapping_bam_out_path,omitempty"`
	// bam文件输出名称。
	MappingBamOutFilename *string `json:"mapping_bam_out_filename,omitempty" xml:"mapping_bam_out_filename,omitempty"`
	// wgs oss数据的存放region。
	WgsOssRegion *string `json:"wgs_oss_region,omitempty" xml:"wgs_oss_region,omitempty"`
	// wgs的第一个fastq文件名。
	WgsFastqFirstFilename *string `json:"wgs_fastq_first_filename,omitempty" xml:"wgs_fastq_first_filename,omitempty"`
	// wgs的第二个fastq文件名。
	WgsFastqSecondFilename *string `json:"wgs_fastq_second_filename,omitempty" xml:"wgs_fastq_second_filename,omitempty"`
	// 存放wgs的bucket名称。
	WgsBucketName *string `json:"wgs_bucket_name,omitempty" xml:"wgs_bucket_name,omitempty"`
	// wgs的fastq文件路径。
	WgsFastqPath *string `json:"wgs_fastq_path,omitempty" xml:"wgs_fastq_path,omitempty"`
	// wgs的reference文件路径。
	WgsReferencePath *string `json:"wgs_reference_path,omitempty" xml:"wgs_reference_path,omitempty"`
	// wgs的vcf输出路径。
	WgsVcfOutPath *string `json:"wgs_vcf_out_path,omitempty" xml:"wgs_vcf_out_path,omitempty"`
	// wgs的vcf输出文件名称。
	WgsVcfOutFilename *string `json:"wgs_vcf_out_filename,omitempty" xml:"wgs_vcf_out_filename,omitempty"`
}

func (s StartWorkflowRequest) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowRequest) GoString() string {
	return s.String()
}

func (s *StartWorkflowRequest) SetWorkflowType(v string) *StartWorkflowRequest {
	s.WorkflowType = &v
	return s
}

func (s *StartWorkflowRequest) SetService(v string) *StartWorkflowRequest {
	s.Service = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingOssRegion(v string) *StartWorkflowRequest {
	s.MappingOssRegion = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqFirstFilename(v string) *StartWorkflowRequest {
	s.MappingFastqFirstFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqSecondFilename(v string) *StartWorkflowRequest {
	s.MappingFastqSecondFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingBucketName(v string) *StartWorkflowRequest {
	s.MappingBucketName = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqPath(v string) *StartWorkflowRequest {
	s.MappingFastqPath = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingReferencePath(v string) *StartWorkflowRequest {
	s.MappingReferencePath = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingIsMarkDup(v string) *StartWorkflowRequest {
	s.MappingIsMarkDup = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingBamOutPath(v string) *StartWorkflowRequest {
	s.MappingBamOutPath = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingBamOutFilename(v string) *StartWorkflowRequest {
	s.MappingBamOutFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsOssRegion(v string) *StartWorkflowRequest {
	s.WgsOssRegion = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqFirstFilename(v string) *StartWorkflowRequest {
	s.WgsFastqFirstFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqSecondFilename(v string) *StartWorkflowRequest {
	s.WgsFastqSecondFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsBucketName(v string) *StartWorkflowRequest {
	s.WgsBucketName = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqPath(v string) *StartWorkflowRequest {
	s.WgsFastqPath = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsReferencePath(v string) *StartWorkflowRequest {
	s.WgsReferencePath = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsVcfOutPath(v string) *StartWorkflowRequest {
	s.WgsVcfOutPath = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsVcfOutFilename(v string) *StartWorkflowRequest {
	s.WgsVcfOutFilename = &v
	return s
}

type StartWorkflowResponseBody struct {
	// 工作流名称
	JobName *string `json:"JobName,omitempty" xml:"JobName,omitempty"`
}

func (s StartWorkflowResponseBody) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowResponseBody) GoString() string {
	return s.String()
}

func (s *StartWorkflowResponseBody) SetJobName(v string) *StartWorkflowResponseBody {
	s.JobName = &v
	return s
}

type StartWorkflowResponse struct {
	Headers map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *StartWorkflowResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s StartWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowResponse) GoString() string {
	return s.String()
}

func (s *StartWorkflowResponse) SetHeaders(v map[string]*string) *StartWorkflowResponse {
	s.Headers = v
	return s
}

func (s *StartWorkflowResponse) SetBody(v *StartWorkflowResponseBody) *StartWorkflowResponse {
	s.Body = v
	return s
}

type ScaleOutClusterRequest struct {
	// 扩容节点数
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
	// keypair名称，和login_password二选一。
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// SSH登录密码，和key_pair二选一。
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// 虚拟交换机
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// Worker节点付费类型
	WorkerInstanceChargeType *string `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	// Worker节点包年包月时长
	WorkerPeriod *int64 `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	// Worker节点包年包月周期
	WorkerPeriodUnit *string `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	// Worker节点到期是否自动续费
	WorkerAutoRenew *bool `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	// Worker节点自动续费时长
	WorkerAutoRenewPeriod *int64 `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	// Worker节点实例规格
	WorkerInstanceTypes []*string `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	// Worker节点系统盘类型
	WorkerSystemDiskCategory *string `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	// Worker节点系统盘大小
	WorkerSystemDiskSize *int64 `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
	// Worker节点数据盘配置
	WorkerDataDisks []*ScaleOutClusterRequestWorkerDataDisks `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	// 在节点上安装云监控
	CloudMonitorFlags *bool `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	// CPU亲和性策略
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// 自定义镜像
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// 节点自定义数据
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
	// RDS白名单
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// 节点标签
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// 节点污点信息
	Taints  []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	Runtime *Runtime `json:"runtime,omitempty" xml:"runtime,omitempty"`
}

func (s ScaleOutClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterRequest) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterRequest) SetCount(v int64) *ScaleOutClusterRequest {
	s.Count = &v
	return s
}

func (s *ScaleOutClusterRequest) SetKeyPair(v string) *ScaleOutClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *ScaleOutClusterRequest) SetLoginPassword(v string) *ScaleOutClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *ScaleOutClusterRequest) SetVswitchIds(v []*string) *ScaleOutClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerInstanceChargeType(v string) *ScaleOutClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerPeriod(v int64) *ScaleOutClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerPeriodUnit(v string) *ScaleOutClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerAutoRenew(v bool) *ScaleOutClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerAutoRenewPeriod(v int64) *ScaleOutClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerInstanceTypes(v []*string) *ScaleOutClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerSystemDiskCategory(v string) *ScaleOutClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerSystemDiskSize(v int64) *ScaleOutClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerDataDisks(v []*ScaleOutClusterRequestWorkerDataDisks) *ScaleOutClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *ScaleOutClusterRequest) SetCloudMonitorFlags(v bool) *ScaleOutClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *ScaleOutClusterRequest) SetCpuPolicy(v string) *ScaleOutClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *ScaleOutClusterRequest) SetImageId(v string) *ScaleOutClusterRequest {
	s.ImageId = &v
	return s
}

func (s *ScaleOutClusterRequest) SetUserData(v string) *ScaleOutClusterRequest {
	s.UserData = &v
	return s
}

func (s *ScaleOutClusterRequest) SetRdsInstances(v []*string) *ScaleOutClusterRequest {
	s.RdsInstances = v
	return s
}

func (s *ScaleOutClusterRequest) SetTags(v []*Tag) *ScaleOutClusterRequest {
	s.Tags = v
	return s
}

func (s *ScaleOutClusterRequest) SetTaints(v []*Taint) *ScaleOutClusterRequest {
	s.Taints = v
	return s
}

func (s *ScaleOutClusterRequest) SetRuntime(v *Runtime) *ScaleOutClusterRequest {
	s.Runtime = v
	return s
}

type ScaleOutClusterRequestWorkerDataDisks struct {
	// 数据盘类型,默认值：cloud_efficiency
	Category *string `json:"category,omitempty" xml:"category,omitempty"`
	// 数据盘大小，单位为GiB。  取值范围：[40,32768]
	Size *string `json:"size,omitempty" xml:"size,omitempty"`
	// 是否对数据盘加密
	Encrypted *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	// 自动快照策略ID，云盘会按照快照策略自动备份。
	AutoSnapshotPolicyId *string `json:"auto_snapshot_policy_id,omitempty" xml:"auto_snapshot_policy_id,omitempty"`
}

func (s ScaleOutClusterRequestWorkerDataDisks) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterRequestWorkerDataDisks) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetCategory(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Category = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetSize(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Size = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetEncrypted(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Encrypted = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetAutoSnapshotPolicyId(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.AutoSnapshotPolicyId = &v
	return s
}

type ScaleOutClusterResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 请求ID。
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleOutClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterResponseBody) SetClusterId(v string) *ScaleOutClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ScaleOutClusterResponseBody) SetRequestId(v string) *ScaleOutClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ScaleOutClusterResponseBody) SetTaskId(v string) *ScaleOutClusterResponseBody {
	s.TaskId = &v
	return s
}

type ScaleOutClusterResponse struct {
	Headers map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *ScaleOutClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleOutClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterResponse) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterResponse) SetHeaders(v map[string]*string) *ScaleOutClusterResponse {
	s.Headers = v
	return s
}

func (s *ScaleOutClusterResponse) SetBody(v *ScaleOutClusterResponseBody) *ScaleOutClusterResponse {
	s.Body = v
	return s
}

type DescribeEventsRequest struct {
	// 集群ID
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 事件类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 页数
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 没页记录数量
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
}

func (s DescribeEventsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsRequest) GoString() string {
	return s.String()
}

func (s *DescribeEventsRequest) SetClusterId(v string) *DescribeEventsRequest {
	s.ClusterId = &v
	return s
}

func (s *DescribeEventsRequest) SetType(v string) *DescribeEventsRequest {
	s.Type = &v
	return s
}

func (s *DescribeEventsRequest) SetPageSize(v int64) *DescribeEventsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeEventsRequest) SetPageNumber(v int64) *DescribeEventsRequest {
	s.PageNumber = &v
	return s
}

type DescribeEventsResponseBody struct {
	Events   []*DescribeEventsResponseBodyEvents `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	PageInfo *DescribeEventsResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeEventsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBody) SetEvents(v []*DescribeEventsResponseBodyEvents) *DescribeEventsResponseBody {
	s.Events = v
	return s
}

func (s *DescribeEventsResponseBody) SetPageInfo(v *DescribeEventsResponseBodyPageInfo) *DescribeEventsResponseBody {
	s.PageInfo = v
	return s
}

type DescribeEventsResponseBodyEvents struct {
	// 事件ID
	EventId *string `json:"event_id,omitempty" xml:"event_id,omitempty"`
	// 事件类型
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 事件源
	Source *string `json:"source,omitempty" xml:"source,omitempty"`
	// 事件
	Subject *string `json:"subject,omitempty" xml:"subject,omitempty"`
	// 事件开始事件
	Time *string `json:"time,omitempty" xml:"time,omitempty"`
	// 集群ID
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 事件描述
	Data *DescribeEventsResponseBodyEventsData `json:"data,omitempty" xml:"data,omitempty" type:"Struct"`
}

func (s DescribeEventsResponseBodyEvents) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyEvents) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyEvents) SetEventId(v string) *DescribeEventsResponseBodyEvents {
	s.EventId = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetType(v string) *DescribeEventsResponseBodyEvents {
	s.Type = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetSource(v string) *DescribeEventsResponseBodyEvents {
	s.Source = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetSubject(v string) *DescribeEventsResponseBodyEvents {
	s.Subject = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetTime(v string) *DescribeEventsResponseBodyEvents {
	s.Time = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetClusterId(v string) *DescribeEventsResponseBodyEvents {
	s.ClusterId = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetData(v *DescribeEventsResponseBodyEventsData) *DescribeEventsResponseBodyEvents {
	s.Data = v
	return s
}

type DescribeEventsResponseBodyEventsData struct {
	// 事件级别
	Level *string `json:"level,omitempty" xml:"level,omitempty"`
	// 事件状态
	Reason *string `json:"reason,omitempty" xml:"reason,omitempty"`
	// 事件详情
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s DescribeEventsResponseBodyEventsData) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyEventsData) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyEventsData) SetLevel(v string) *DescribeEventsResponseBodyEventsData {
	s.Level = &v
	return s
}

func (s *DescribeEventsResponseBodyEventsData) SetReason(v string) *DescribeEventsResponseBodyEventsData {
	s.Reason = &v
	return s
}

func (s *DescribeEventsResponseBodyEventsData) SetMessage(v string) *DescribeEventsResponseBodyEventsData {
	s.Message = &v
	return s
}

type DescribeEventsResponseBodyPageInfo struct {
	// 页数
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// 每页记录数量
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// 结果总数
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeEventsResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyPageInfo) SetPageSize(v int64) *DescribeEventsResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeEventsResponseBodyPageInfo) SetPageNumber(v int64) *DescribeEventsResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeEventsResponseBodyPageInfo) SetTotalCount(v int64) *DescribeEventsResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeEventsResponse struct {
	Headers map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *DescribeEventsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEventsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponse) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponse) SetHeaders(v map[string]*string) *DescribeEventsResponse {
	s.Headers = v
	return s
}

func (s *DescribeEventsResponse) SetBody(v *DescribeEventsResponseBody) *DescribeEventsResponse {
	s.Body = v
	return s
}

type UpdateK8sClusterUserConfigExpireResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s UpdateK8sClusterUserConfigExpireResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateK8sClusterUserConfigExpireResponse) GoString() string {
	return s.String()
}

func (s *UpdateK8sClusterUserConfigExpireResponse) SetHeaders(v map[string]*string) *UpdateK8sClusterUserConfigExpireResponse {
	s.Headers = v
	return s
}

type TagResourcesRequest struct {
	// 资源ID列表
	ResourceIds []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	// 资源类型定义。取值范围：  只支持CLUSTER这一种资源类型
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 资源所属的地域ID
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// 资源的标签键值对。数组长度范围：1~20。一旦传值，则不允许为空字符串。最多支持128个字符，不能以aliyun和acs:开头，不能包含http://或者https://。
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
}

func (s TagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesRequest) GoString() string {
	return s.String()
}

func (s *TagResourcesRequest) SetResourceIds(v []*string) *TagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *TagResourcesRequest) SetResourceType(v string) *TagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *TagResourcesRequest) SetRegionId(v string) *TagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *TagResourcesRequest) SetTags(v []*Tag) *TagResourcesRequest {
	s.Tags = v
	return s
}

type TagResourcesResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s TagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesResponse) GoString() string {
	return s.String()
}

func (s *TagResourcesResponse) SetHeaders(v map[string]*string) *TagResourcesResponse {
	s.Headers = v
	return s
}

type ModifyClusterTagsRequest struct {
	// 集群标签列表。
	Body []*Tag `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s ModifyClusterTagsRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterTagsRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterTagsRequest) SetBody(v []*Tag) *ModifyClusterTagsRequest {
	s.Body = v
	return s
}

type ModifyClusterTagsResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s ModifyClusterTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterTagsResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterTagsResponse) SetHeaders(v map[string]*string) *ModifyClusterTagsResponse {
	s.Headers = v
	return s
}

type GetKubernetesTriggerRequest struct {
	// 应用所属命名空间。
	Namespace *string `json:"Namespace,omitempty" xml:"Namespace,omitempty"`
	// 应用类型。
	Type *string `json:"Type,omitempty" xml:"Type,omitempty"`
	// 应用名称。
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// 触发器行为。
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s GetKubernetesTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerRequest) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerRequest) SetNamespace(v string) *GetKubernetesTriggerRequest {
	s.Namespace = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetType(v string) *GetKubernetesTriggerRequest {
	s.Type = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetName(v string) *GetKubernetesTriggerRequest {
	s.Name = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetAction(v string) *GetKubernetesTriggerRequest {
	s.Action = &v
	return s
}

type GetKubernetesTriggerResponse struct {
	Headers map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*GetKubernetesTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s GetKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerResponse) SetHeaders(v map[string]*string) *GetKubernetesTriggerResponse {
	s.Headers = v
	return s
}

func (s *GetKubernetesTriggerResponse) SetBody(v []*GetKubernetesTriggerResponseBody) *GetKubernetesTriggerResponse {
	s.Body = v
	return s
}

type GetKubernetesTriggerResponseBody struct {
	// 触发器ID。
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// 触发器名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 集群ID
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 触发器项目名称
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	// 触发器类型。
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// 触发器行为
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// Token
	Token *string `json:"token,omitempty" xml:"token,omitempty"`
}

func (s GetKubernetesTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerResponseBody) SetId(v string) *GetKubernetesTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetName(v string) *GetKubernetesTriggerResponseBody {
	s.Name = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetClusterId(v string) *GetKubernetesTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetProjectId(v string) *GetKubernetesTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetType(v string) *GetKubernetesTriggerResponseBody {
	s.Type = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetAction(v string) *GetKubernetesTriggerResponseBody {
	s.Action = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetToken(v string) *GetKubernetesTriggerResponseBody {
	s.Token = &v
	return s
}

type GetUpgradeStatusResponseBody struct {
	// 错误信息描述。
	ErrorMessage *string `json:"error_message,omitempty" xml:"error_message,omitempty"`
	// 预检查返回ID。
	PrecheckReportId *string `json:"precheck_report_id,omitempty" xml:"precheck_report_id,omitempty"`
	// 升级状态。
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 升级任务执行到哪一步了。
	UpgradeStep *string `json:"upgrade_step,omitempty" xml:"upgrade_step,omitempty"`
	// 升级任务详情。
	UpgradeTask *GetUpgradeStatusResponseBodyUpgradeTask `json:"upgrade_task,omitempty" xml:"upgrade_task,omitempty" type:"Struct"`
}

func (s GetUpgradeStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponseBody) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponseBody) SetErrorMessage(v string) *GetUpgradeStatusResponseBody {
	s.ErrorMessage = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetPrecheckReportId(v string) *GetUpgradeStatusResponseBody {
	s.PrecheckReportId = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetStatus(v string) *GetUpgradeStatusResponseBody {
	s.Status = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetUpgradeStep(v string) *GetUpgradeStatusResponseBody {
	s.UpgradeStep = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetUpgradeTask(v *GetUpgradeStatusResponseBodyUpgradeTask) *GetUpgradeStatusResponseBody {
	s.UpgradeTask = v
	return s
}

type GetUpgradeStatusResponseBodyUpgradeTask struct {
	// 任务状态：  emptry、running、success、failed
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// 任务描述信息。
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s GetUpgradeStatusResponseBodyUpgradeTask) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponseBodyUpgradeTask) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponseBodyUpgradeTask) SetStatus(v string) *GetUpgradeStatusResponseBodyUpgradeTask {
	s.Status = &v
	return s
}

func (s *GetUpgradeStatusResponseBodyUpgradeTask) SetMessage(v string) *GetUpgradeStatusResponseBodyUpgradeTask {
	s.Message = &v
	return s
}

type GetUpgradeStatusResponse struct {
	Headers map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    *GetUpgradeStatusResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponse) SetHeaders(v map[string]*string) *GetUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *GetUpgradeStatusResponse) SetBody(v *GetUpgradeStatusResponseBody) *GetUpgradeStatusResponse {
	s.Body = v
	return s
}

type DescribeClusterResourcesResponse struct {
	Headers map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	Body    []*DescribeClusterResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClusterResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterResourcesResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterResourcesResponse) SetHeaders(v map[string]*string) *DescribeClusterResourcesResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterResourcesResponse) SetBody(v []*DescribeClusterResourcesResponseBody) *DescribeClusterResourcesResponse {
	s.Body = v
	return s
}

type DescribeClusterResourcesResponseBody struct {
	// 集群ID。
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// 资源创建时间。
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// 资源实例ID。
	InstanceId *string `json:"instance_id,omitempty" xml:"instance_id,omitempty"`
	// 资源元信息。
	ResourceInfo *string `json:"resource_info,omitempty" xml:"resource_info,omitempty"`
	// 资源类型。
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// 资源状态。
	State *string `json:"state,omitempty" xml:"state,omitempty"`
}

func (s DescribeClusterResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterResourcesResponseBody) SetClusterId(v string) *DescribeClusterResourcesResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetCreated(v string) *DescribeClusterResourcesResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetInstanceId(v string) *DescribeClusterResourcesResponseBody {
	s.InstanceId = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetResourceInfo(v string) *DescribeClusterResourcesResponseBody {
	s.ResourceInfo = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetResourceType(v string) *DescribeClusterResourcesResponseBody {
	s.ResourceType = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetState(v string) *DescribeClusterResourcesResponseBody {
	s.State = &v
	return s
}

type DeleteClusterNodesRequest struct {
	// 是否自动排空节点上的Pod。
	DrainNode *bool `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	// 是否同时释放 ECS
	ReleaseNode *bool `json:"release_node,omitempty" xml:"release_node,omitempty"`
	// 移除节点列表。
	Nodes []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
}

func (s DeleteClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodesRequest) SetDrainNode(v bool) *DeleteClusterNodesRequest {
	s.DrainNode = &v
	return s
}

func (s *DeleteClusterNodesRequest) SetReleaseNode(v bool) *DeleteClusterNodesRequest {
	s.ReleaseNode = &v
	return s
}

func (s *DeleteClusterNodesRequest) SetNodes(v []*string) *DeleteClusterNodesRequest {
	s.Nodes = v
	return s
}

type DeleteClusterNodesResponse struct {
	Headers map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
}

func (s DeleteClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodesResponse) SetHeaders(v map[string]*string) *DeleteClusterNodesResponse {
	s.Headers = v
	return s
}

type StandardComponentsValue struct {
	// 组件名称。
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// 组件版本。
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
	// 组件描述信息。
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// 是否为必需组件。
	Required *string `json:"required,omitempty" xml:"required,omitempty"`
	// 是否禁止默认安装。
	Disabled *bool `json:"disabled,omitempty" xml:"disabled,omitempty"`
}

func (s StandardComponentsValue) String() string {
	return tea.Prettify(s)
}

func (s StandardComponentsValue) GoString() string {
	return s.String()
}

func (s *StandardComponentsValue) SetName(v string) *StandardComponentsValue {
	s.Name = &v
	return s
}

func (s *StandardComponentsValue) SetVersion(v string) *StandardComponentsValue {
	s.Version = &v
	return s
}

func (s *StandardComponentsValue) SetDescription(v string) *StandardComponentsValue {
	s.Description = &v
	return s
}

func (s *StandardComponentsValue) SetRequired(v string) *StandardComponentsValue {
	s.Required = &v
	return s
}

func (s *StandardComponentsValue) SetDisabled(v bool) *StandardComponentsValue {
	s.Disabled = &v
	return s
}

type Client struct {
	openapi.Client
}

func NewClient(config *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	client.EndpointRule = tea.String("regional")
	client.EndpointMap = map[string]*string{
		"ap-northeast-2-pop":          tea.String("cs.aliyuncs.com"),
		"cn-beijing-finance-1":        tea.String("cs.aliyuncs.com"),
		"cn-beijing-finance-pop":      tea.String("cs.aliyuncs.com"),
		"cn-beijing-gov-1":            tea.String("cs.aliyuncs.com"),
		"cn-beijing-nu16-b01":         tea.String("cs.aliyuncs.com"),
		"cn-edge-1":                   tea.String("cs.aliyuncs.com"),
		"cn-fujian":                   tea.String("cs.aliyuncs.com"),
		"cn-haidian-cm12-c01":         tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-bj-b01":          tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-finance":         tea.String("cs-vpc.cn-hangzhou-finance.aliyuncs.com"),
		"cn-hangzhou-internal-prod-1": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-1": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-2": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-3": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-test-306":        tea.String("cs.aliyuncs.com"),
		"cn-hongkong-finance-pop":     tea.String("cs.aliyuncs.com"),
		"cn-huhehaote-nebula-1":       tea.String("cs.aliyuncs.com"),
		"cn-qingdao-nebula":           tea.String("cs.aliyuncs.com"),
		"cn-shanghai-et15-b01":        tea.String("cs.aliyuncs.com"),
		"cn-shanghai-et2-b01":         tea.String("cs.aliyuncs.com"),
		"cn-shanghai-finance-1":       tea.String("cs-vpc.cn-shanghai-finance-1.aliyuncs.com"),
		"cn-shanghai-inner":           tea.String("cs.aliyuncs.com"),
		"cn-shanghai-internal-test-1": tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-finance-1":       tea.String("cs-vpc.cn-shenzhen-finance-1.aliyuncs.com"),
		"cn-shenzhen-inner":           tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-st4-d01":         tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-su18-b01":        tea.String("cs.aliyuncs.com"),
		"cn-wuhan":                    tea.String("cs.aliyuncs.com"),
		"cn-yushanfang":               tea.String("cs.aliyuncs.com"),
		"cn-zhangbei":                 tea.String("cs.aliyuncs.com"),
		"cn-zhangbei-na61-b01":        tea.String("cs.aliyuncs.com"),
		"cn-zhangjiakou-na62-a01":     tea.String("cs.aliyuncs.com"),
		"cn-zhengzhou-nebula-1":       tea.String("cs.aliyuncs.com"),
		"eu-west-1-oxs":               tea.String("cs.aliyuncs.com"),
		"rus-west-1-pop":              tea.String("cs.aliyuncs.com"),
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}
	client.Endpoint, _err = client.GetEndpoint(tea.String("cs"), client.RegionId, client.EndpointRule, client.Network, client.Suffix, client.EndpointMap, client.Endpoint)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *Client) GetEndpoint(productId *string, regionId *string, endpointRule *string, network *string, suffix *string, endpointMap map[string]*string, endpoint *string) (_result *string, _err error) {
	if !tea.BoolValue(util.Empty(endpoint)) {
		_result = endpoint
		return _result, _err
	}

	if !tea.BoolValue(util.IsUnset(endpointMap)) && !tea.BoolValue(util.Empty(endpointMap[tea.StringValue(regionId)])) {
		_result = endpointMap[tea.StringValue(regionId)]
		return _result, _err
	}

	_body, _err := endpointutil.GetEndpointRules(productId, regionId, endpointRule, network, suffix)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTagResources(request *ListTagResourcesRequest) (_result *ListTagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ListTagResourcesResponse{}
	_body, _err := client.ListTagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTagResourcesWithOptions(tmpReq *ListTagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ListTagResourcesResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &ListTagResourcesShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.ResourceIds)) {
		request.ResourceIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.ResourceIds, tea.String("resource_ids"), tea.String("json"))
	}

	if !tea.BoolValue(util.IsUnset(tmpReq.Tags)) {
		request.TagsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.Tags, tea.String("tags"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ResourceIdsShrink)) {
		query["resource_ids"] = request.ResourceIdsShrink
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.TagsShrink)) {
		query["tags"] = request.TagsShrink
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["next_token"] = request.NextToken
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &ListTagResourcesResponse{}
	_body, _err := client.DoROARequest(tea.String("ListTagResources"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/tags"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UntagResources(request *UntagResourcesRequest) (_result *UntagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UntagResourcesResponse{}
	_body, _err := client.UntagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UntagResourcesWithOptions(request *UntagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UntagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceIds)) {
		query["resource_ids"] = request.ResourceIds
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.TagKeys)) {
		query["tag_keys"] = request.TagKeys
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &UntagResourcesResponse{}
	_body, _err := client.DoROARequest(tea.String("UntagResources"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/tags"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyCluster(ClusterId *string, request *ModifyClusterRequest) (_result *ModifyClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterResponse{}
	_body, _err := client.ModifyClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterWithOptions(ClusterId *string, request *ModifyClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ApiServerEip)) {
		body["api_server_eip"] = request.ApiServerEip
	}

	if !tea.BoolValue(util.IsUnset(request.ApiServerEipId)) {
		body["api_server_eip_id"] = request.ApiServerEipId
	}

	if !tea.BoolValue(util.IsUnset(request.DeletionProtection)) {
		body["deletion_protection"] = request.DeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceDeletionProtection)) {
		body["instance_deletion_protection"] = request.InstanceDeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.IngressDomainRebinding)) {
		body["ingress_domain_rebinding"] = request.IngressDomainRebinding
	}

	if !tea.BoolValue(util.IsUnset(request.IngressLoadbalancerId)) {
		body["ingress_loadbalancer_id"] = request.IngressLoadbalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		body["resource_group_id"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.MaintenanceWindow))) {
		body["maintenance_window"] = request.MaintenanceWindow
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ModifyClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("ModifyCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAttachScripts(ClusterId *string, request *DescribeClusterAttachScriptsRequest) (_result *DescribeClusterAttachScriptsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAttachScriptsResponse{}
	_body, _err := client.DescribeClusterAttachScriptsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAttachScriptsWithOptions(ClusterId *string, request *DescribeClusterAttachScriptsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAttachScriptsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		body["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.Arch)) {
		body["arch"] = request.Arch
	}

	if !tea.BoolValue(util.IsUnset(request.Options)) {
		body["options"] = request.Options
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &DescribeClusterAttachScriptsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterAttachScripts"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/attachscript"), tea.String("string"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RemoveClusterNodes(ClusterId *string, request *RemoveClusterNodesRequest) (_result *RemoveClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RemoveClusterNodesResponse{}
	_body, _err := client.RemoveClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RemoveClusterNodesWithOptions(ClusterId *string, request *RemoveClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RemoveClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DrainNode)) {
		body["drain_node"] = request.DrainNode
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.ReleaseNode)) {
		body["release_node"] = request.ReleaseNode
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &RemoveClusterNodesResponse{}
	_body, _err := client.DoROARequest(tea.String("RemoveClusterNodes"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/nodes/remove"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeKubernetesVersionMetadata(request *DescribeKubernetesVersionMetadataRequest) (_result *DescribeKubernetesVersionMetadataResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeKubernetesVersionMetadataResponse{}
	_body, _err := client.DescribeKubernetesVersionMetadataWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeKubernetesVersionMetadataWithOptions(request *DescribeKubernetesVersionMetadataRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeKubernetesVersionMetadataResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Region)) {
		query["Region"] = request.Region
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["ClusterType"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesVersion)) {
		query["KubernetesVersion"] = request.KubernetesVersion
	}

	if !tea.BoolValue(util.IsUnset(request.Profile)) {
		query["Profile"] = request.Profile
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeKubernetesVersionMetadataResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeKubernetesVersionMetadata"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/api/v1/metadata/versions"), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterLogs(ClusterId *string) (_result *DescribeClusterLogsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterLogsResponse{}
	_body, _err := client.DescribeClusterLogsWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterLogsWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterLogsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterLogsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterLogs"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/logs"), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateKubernetesTrigger(request *CreateKubernetesTriggerRequest) (_result *CreateKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateKubernetesTriggerResponse{}
	_body, _err := client.CreateKubernetesTriggerWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateKubernetesTriggerWithOptions(request *CreateKubernetesTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateKubernetesTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterId)) {
		body["cluster_id"] = request.ClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.ProjectId)) {
		body["project_id"] = request.ProjectId
	}

	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		body["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateKubernetesTriggerResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateKubernetesTrigger"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/triggers"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GrantPermissions(uid *string, request *GrantPermissionsRequest) (_result *GrantPermissionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GrantPermissionsResponse{}
	_body, _err := client.GrantPermissionsWithOptions(uid, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GrantPermissionsWithOptions(uid *string, request *GrantPermissionsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GrantPermissionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	_result = &GrantPermissionsResponse{}
	_body, _err := client.DoROARequest(tea.String("GrantPermissions"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/permissions/users/"+tea.StringValue(uid)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterDetail(ClusterId *string) (_result *DescribeClusterDetailResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterDetailResponse{}
	_body, _err := client.DescribeClusterDetailWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterDetailWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterDetailResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterDetailResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterDetail"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PauseComponentUpgrade(clusterid *string, componentid *string) (_result *PauseComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &PauseComponentUpgradeResponse{}
	_body, _err := client.PauseComponentUpgradeWithOptions(clusterid, componentid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PauseComponentUpgradeWithOptions(clusterid *string, componentid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *PauseComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &PauseComponentUpgradeResponse{}
	_body, _err := client.DoROARequest(tea.String("PauseComponentUpgrade"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(clusterid)+"/components/"+tea.StringValue(componentid)+"/pause"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusters(request *DescribeClustersRequest) (_result *DescribeClustersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClustersResponse{}
	_body, _err := client.DescribeClustersWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClustersWithOptions(request *DescribeClustersRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClustersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["clusterType"] = request.ClusterType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClustersResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusters"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters"), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeUserPermission(uid *string) (_result *DescribeUserPermissionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeUserPermissionResponse{}
	_body, _err := client.DescribeUserPermissionWithOptions(uid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeUserPermissionWithOptions(uid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeUserPermissionResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeUserPermissionResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeUserPermission"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/permissions/users/"+tea.StringValue(uid)), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterNodePool(ClusterId *string, NodepoolId *string, request *ModifyClusterNodePoolRequest) (_result *ModifyClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterNodePoolResponse{}
	_body, _err := client.ModifyClusterNodePoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterNodePoolWithOptions(ClusterId *string, NodepoolId *string, request *ModifyClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.AutoScaling))) {
		body["auto_scaling"] = request.AutoScaling
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.KubernetesConfig))) {
		body["kubernetes_config"] = request.KubernetesConfig
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.NodepoolInfo))) {
		body["nodepool_info"] = request.NodepoolInfo
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.ScalingGroup))) {
		body["scaling_group"] = request.ScalingGroup
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.TeeConfig))) {
		body["tee_config"] = request.TeeConfig
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.Management))) {
		body["management"] = request.Management
	}

	if !tea.BoolValue(util.IsUnset(request.UpdateNodes)) {
		body["update_nodes"] = request.UpdateNodes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ModifyClusterNodePoolResponse{}
	_body, _err := client.DoROARequest(tea.String("ModifyClusterNodePool"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools/"+tea.StringValue(NodepoolId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResumeUpgradeCluster(ClusterId *string) (_result *ResumeUpgradeClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ResumeUpgradeClusterResponse{}
	_body, _err := client.ResumeUpgradeClusterWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResumeUpgradeClusterWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ResumeUpgradeClusterResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &ResumeUpgradeClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("ResumeUpgradeCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/upgrade/resume"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) OpenAckService(request *OpenAckServiceRequest) (_result *OpenAckServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &OpenAckServiceResponse{}
	_body, _err := client.OpenAckServiceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) OpenAckServiceWithOptions(request *OpenAckServiceRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *OpenAckServiceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &OpenAckServiceResponse{}
	_body, _err := client.DoROARequest(tea.String("OpenAckService"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/service/open"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ScaleClusterNodePool(ClusterId *string, NodepoolId *string, request *ScaleClusterNodePoolRequest) (_result *ScaleClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleClusterNodePoolResponse{}
	_body, _err := client.ScaleClusterNodePoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ScaleClusterNodePoolWithOptions(ClusterId *string, NodepoolId *string, request *ScaleClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ScaleClusterNodePoolResponse{}
	_body, _err := client.DoROARequest(tea.String("ScaleClusterNodePool"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools/"+tea.StringValue(NodepoolId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolDetail(ClusterId *string, NodepoolId *string) (_result *DescribeClusterNodePoolDetailResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodePoolDetailResponse{}
	_body, _err := client.DescribeClusterNodePoolDetailWithOptions(ClusterId, NodepoolId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolDetailWithOptions(ClusterId *string, NodepoolId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodePoolDetailResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterNodePoolDetailResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterNodePoolDetail"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools/"+tea.StringValue(NodepoolId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateClusterNodePool(ClusterId *string, request *CreateClusterNodePoolRequest) (_result *CreateClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateClusterNodePoolResponse{}
	_body, _err := client.CreateClusterNodePoolWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateClusterNodePoolWithOptions(ClusterId *string, request *CreateClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.AutoScaling))) {
		body["auto_scaling"] = request.AutoScaling
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.KubernetesConfig))) {
		body["kubernetes_config"] = request.KubernetesConfig
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.NodepoolInfo))) {
		body["nodepool_info"] = request.NodepoolInfo
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.ScalingGroup))) {
		body["scaling_group"] = request.ScalingGroup
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.TeeConfig))) {
		body["tee_config"] = request.TeeConfig
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.Management))) {
		body["management"] = request.Management
	}

	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateClusterNodePoolResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateClusterNodePool"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterUserKubeconfig(ClusterId *string, request *DescribeClusterUserKubeconfigRequest) (_result *DescribeClusterUserKubeconfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterUserKubeconfigResponse{}
	_body, _err := client.DescribeClusterUserKubeconfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterUserKubeconfigWithOptions(ClusterId *string, request *DescribeClusterUserKubeconfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterUserKubeconfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	if !tea.BoolValue(util.IsUnset(request.TemporaryDurationMinutes)) {
		query["TemporaryDurationMinutes"] = request.TemporaryDurationMinutes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClusterUserKubeconfigResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterUserKubeconfig"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/k8s/"+tea.StringValue(ClusterId)+"/user_config"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ScaleCluster(ClusterId *string, request *ScaleClusterRequest) (_result *ScaleClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleClusterResponse{}
	_body, _err := client.ScaleClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ScaleClusterWithOptions(ClusterId *string, request *ScaleClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.DisableRollback)) {
		body["disable_rollback"] = request.DisableRollback
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisk)) {
		body["worker_data_disk"] = request.WorkerDataDisk
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ScaleClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("ScaleCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonUpgradeStatus(ClusterId *string, ComponentId *string) (_result *DescribeClusterAddonUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonUpgradeStatusResponse{}
	_body, _err := client.DescribeClusterAddonUpgradeStatusWithOptions(ClusterId, ComponentId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonUpgradeStatusWithOptions(ClusterId *string, ComponentId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonUpgradeStatusResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterAddonUpgradeStatusResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterAddonUpgradeStatus"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/"+tea.StringValue(ComponentId)+"/upgradestatus"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeAddons(request *DescribeAddonsRequest) (_result *DescribeAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeAddonsResponse{}
	_body, _err := client.DescribeAddonsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeAddonsWithOptions(request *DescribeAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Region)) {
		query["region"] = request.Region
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["cluster_type"] = request.ClusterType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeAddonsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeAddons"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/components/metadata"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateAutoscalingConfig(ClusterId *string, request *CreateAutoscalingConfigRequest) (_result *CreateAutoscalingConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateAutoscalingConfigResponse{}
	_body, _err := client.CreateAutoscalingConfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateAutoscalingConfigWithOptions(ClusterId *string, request *CreateAutoscalingConfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateAutoscalingConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CoolDownDuration)) {
		body["cool_down_duration"] = request.CoolDownDuration
	}

	if !tea.BoolValue(util.IsUnset(request.UnneededDuration)) {
		body["unneeded_duration"] = request.UnneededDuration
	}

	if !tea.BoolValue(util.IsUnset(request.UtilizationThreshold)) {
		body["utilization_threshold"] = request.UtilizationThreshold
	}

	if !tea.BoolValue(util.IsUnset(request.GpuUtilizationThreshold)) {
		body["gpu_utilization_threshold"] = request.GpuUtilizationThreshold
	}

	if !tea.BoolValue(util.IsUnset(request.ScanInterval)) {
		body["scan_interval"] = request.ScanInterval
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateAutoscalingConfigResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateAutoscalingConfig"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/cluster/"+tea.StringValue(ClusterId)+"/autoscale/config/"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateCluster(request *CreateClusterRequest) (_result *CreateClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateClusterResponse{}
	_body, _err := client.CreateClusterWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateClusterWithOptions(request *CreateClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		body["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		body["cluster_type"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterSpec)) {
		body["cluster_spec"] = request.ClusterSpec
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesVersion)) {
		body["kubernetes_version"] = request.KubernetesVersion
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.Runtime))) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Vpcid)) {
		body["vpcid"] = request.Vpcid
	}

	if !tea.BoolValue(util.IsUnset(request.PodVswitchIds)) {
		body["pod_vswitch_ids"] = request.PodVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.ContainerCidr)) {
		body["container_cidr"] = request.ContainerCidr
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceCidr)) {
		body["service_cidr"] = request.ServiceCidr
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityGroupId)) {
		body["security_group_id"] = request.SecurityGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.IsEnterpriseSecurityGroup)) {
		body["is_enterprise_security_group"] = request.IsEnterpriseSecurityGroup
	}

	if !tea.BoolValue(util.IsUnset(request.SnatEntry)) {
		body["snat_entry"] = request.SnatEntry
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointPublicAccess)) {
		body["endpoint_public_access"] = request.EndpointPublicAccess
	}

	if !tea.BoolValue(util.IsUnset(request.SshFlags)) {
		body["ssh_flags"] = request.SshFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Timezone)) {
		body["timezone"] = request.Timezone
	}

	if !tea.BoolValue(util.IsUnset(request.NodeCidrMask)) {
		body["node_cidr_mask"] = request.NodeCidrMask
	}

	if !tea.BoolValue(util.IsUnset(request.UserCa)) {
		body["user_ca"] = request.UserCa
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterDomain)) {
		body["cluster_domain"] = request.ClusterDomain
	}

	if !tea.BoolValue(util.IsUnset(request.NodeNameMode)) {
		body["node_name_mode"] = request.NodeNameMode
	}

	if !tea.BoolValue(util.IsUnset(request.CustomSan)) {
		body["custom_san"] = request.CustomSan
	}

	if !tea.BoolValue(util.IsUnset(request.EncryptionProviderKey)) {
		body["encryption_provider_key"] = request.EncryptionProviderKey
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceAccountIssuer)) {
		body["service_account_issuer"] = request.ServiceAccountIssuer
	}

	if !tea.BoolValue(util.IsUnset(request.ApiAudiences)) {
		body["api_audiences"] = request.ApiAudiences
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Addons)) {
		body["addons"] = request.Addons
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Platform)) {
		body["platform"] = request.Platform
	}

	if !tea.BoolValue(util.IsUnset(request.OsType)) {
		body["os_type"] = request.OsType
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.ProxyMode)) {
		body["proxy_mode"] = request.ProxyMode
	}

	if !tea.BoolValue(util.IsUnset(request.NodePortRange)) {
		body["node_port_range"] = request.NodePortRange
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.MasterCount)) {
		body["master_count"] = request.MasterCount
	}

	if !tea.BoolValue(util.IsUnset(request.MasterVswitchIds)) {
		body["master_vswitch_ids"] = request.MasterVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.MasterInstanceTypes)) {
		body["master_instance_types"] = request.MasterInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskCategory)) {
		body["master_system_disk_category"] = request.MasterSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskSize)) {
		body["master_system_disk_size"] = request.MasterSystemDiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskSnapshotPolicyId)) {
		body["master_system_disk_snapshot_policy_id"] = request.MasterSystemDiskSnapshotPolicyId
	}

	if !tea.BoolValue(util.IsUnset(request.MasterInstanceChargeType)) {
		body["master_instance_charge_type"] = request.MasterInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.MasterPeriodUnit)) {
		body["master_period_unit"] = request.MasterPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.MasterPeriod)) {
		body["master_period"] = request.MasterPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.MasterAutoRenew)) {
		body["master_auto_renew"] = request.MasterAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.MasterAutoRenewPeriod)) {
		body["master_auto_renew_period"] = request.MasterAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.NumOfNodes)) {
		body["num_of_nodes"] = request.NumOfNodes
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerVswitchIds)) {
		body["worker_vswitch_ids"] = request.WorkerVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSnapshotPolicyId)) {
		body["worker_system_disk_snapshot_policy_id"] = request.WorkerSystemDiskSnapshotPolicyId
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.Instances)) {
		body["instances"] = request.Instances
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceDiscoveryTypes)) {
		body["service_discovery_types"] = request.ServiceDiscoveryTypes
	}

	if !tea.BoolValue(util.IsUnset(request.NatGateway)) {
		body["nat_gateway"] = request.NatGateway
	}

	if !tea.BoolValue(util.IsUnset(request.ZoneId)) {
		body["zone_id"] = request.ZoneId
	}

	if !tea.BoolValue(util.IsUnset(request.Profile)) {
		body["profile"] = request.Profile
	}

	if !tea.BoolValue(util.IsUnset(request.DeletionProtection)) {
		body["deletion_protection"] = request.DeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.DisableRollback)) {
		body["disable_rollback"] = request.DisableRollback
	}

	if !tea.BoolValue(util.IsUnset(request.TimeoutMins)) {
		body["timeout_mins"] = request.TimeoutMins
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeCluster(ClusterId *string, request *UpgradeClusterRequest) (_result *UpgradeClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpgradeClusterResponse{}
	_body, _err := client.UpgradeClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeClusterWithOptions(ClusterId *string, request *UpgradeClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpgradeClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ComponentName)) {
		body["component_name"] = request.ComponentName
	}

	if !tea.BoolValue(util.IsUnset(request.NextVersion)) {
		body["next_version"] = request.NextVersion
	}

	if !tea.BoolValue(util.IsUnset(request.Version)) {
		body["version"] = request.Version
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &UpgradeClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("UpgradeCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/upgrade"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelWorkflow(workflowName *string, request *CancelWorkflowRequest) (_result *CancelWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelWorkflowResponse{}
	_body, _err := client.CancelWorkflowWithOptions(workflowName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelWorkflowWithOptions(workflowName *string, request *CancelWorkflowRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelWorkflowResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CancelWorkflowResponse{}
	_body, _err := client.DoROARequest(tea.String("CancelWorkflow"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/gs/workflow/"+tea.StringValue(workflowName)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) AttachInstances(ClusterId *string, request *AttachInstancesRequest) (_result *AttachInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &AttachInstancesResponse{}
	_body, _err := client.AttachInstancesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) AttachInstancesWithOptions(ClusterId *string, request *AttachInstancesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *AttachInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Instances)) {
		body["instances"] = request.Instances
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.Password)) {
		body["password"] = request.Password
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.IsEdgeWorker)) {
		body["is_edge_worker"] = request.IsEdgeWorker
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		body["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.Runtime))) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &AttachInstancesResponse{}
	_body, _err := client.DoROARequest(tea.String("AttachInstances"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/attach"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTemplates(request *DescribeTemplatesRequest) (_result *DescribeTemplatesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTemplatesResponse{}
	_body, _err := client.DescribeTemplatesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTemplatesWithOptions(request *DescribeTemplatesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTemplatesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		query["template_type"] = request.TemplateType
	}

	if !tea.BoolValue(util.IsUnset(request.PageNum)) {
		query["page_num"] = request.PageNum
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeTemplatesResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeTemplates"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/templates"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PauseClusterUpgrade(ClusterId *string) (_result *PauseClusterUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &PauseClusterUpgradeResponse{}
	_body, _err := client.PauseClusterUpgradeWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PauseClusterUpgradeWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *PauseClusterUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &PauseClusterUpgradeResponse{}
	_body, _err := client.DoROARequest(tea.String("PauseClusterUpgrade"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/upgrade/pause"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteTemplate(TemplateId *string) (_result *DeleteTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteTemplateResponse{}
	_body, _err := client.DeleteTemplateWithOptions(TemplateId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteTemplateWithOptions(TemplateId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteTemplateResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DeleteTemplateResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteTemplate"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/templates/"+tea.StringValue(TemplateId)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTemplateAttribute(TemplateId *string, request *DescribeTemplateAttributeRequest) (_result *DescribeTemplateAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTemplateAttributeResponse{}
	_body, _err := client.DescribeTemplateAttributeWithOptions(TemplateId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTemplateAttributeWithOptions(TemplateId *string, request *DescribeTemplateAttributeRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTemplateAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		query["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeTemplateAttributeResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeTemplateAttribute"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/templates/"+tea.StringValue(TemplateId)), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateTemplate(request *CreateTemplateRequest) (_result *CreateTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateTemplateResponse{}
	_body, _err := client.CreateTemplateWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateTemplateWithOptions(request *CreateTemplateRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateTemplateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Template)) {
		body["template"] = request.Template
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		body["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &CreateTemplateResponse{}
	_body, _err := client.DoROARequest(tea.String("CreateTemplate"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/templates"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodes(ClusterId *string, request *DescribeClusterNodesRequest) (_result *DescribeClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodesResponse{}
	_body, _err := client.DescribeClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodesWithOptions(ClusterId *string, request *DescribeClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceIds)) {
		query["instanceIds"] = request.InstanceIds
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		query["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.State)) {
		query["state"] = request.State
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["pageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["pageNumber"] = request.PageNumber
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClusterNodesResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterNodes"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodes"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteCluster(ClusterId *string, request *DeleteClusterRequest) (_result *DeleteClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterResponse{}
	_body, _err := client.DeleteClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteClusterWithOptions(ClusterId *string, tmpReq *DeleteClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &DeleteClusterShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.RetainResources)) {
		request.RetainResourcesShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.RetainResources, tea.String("retain_resources"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.RetainAllResources)) {
		query["retain_all_resources"] = request.RetainAllResources
	}

	if !tea.BoolValue(util.IsUnset(request.KeepSlb)) {
		query["keep_slb"] = request.KeepSlb
	}

	if !tea.BoolValue(util.IsUnset(request.RetainResourcesShrink)) {
		query["retain_resources"] = request.RetainResourcesShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DeleteClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelComponentUpgrade(clusterId *string, componentId *string) (_result *CancelComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelComponentUpgradeResponse{}
	_body, _err := client.CancelComponentUpgradeWithOptions(clusterId, componentId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelComponentUpgradeWithOptions(clusterId *string, componentId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &CancelComponentUpgradeResponse{}
	_body, _err := client.DoROARequest(tea.String("CancelComponentUpgrade"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(clusterId)+"/components/"+tea.StringValue(componentId)+"/cancel"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) MigrateCluster(clusterId *string) (_result *MigrateClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &MigrateClusterResponse{}
	_body, _err := client.MigrateClusterWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) MigrateClusterWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *MigrateClusterResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &MigrateClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("MigrateCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(clusterId)+"/migrate"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonsVersion(ClusterId *string) (_result *DescribeClusterAddonsVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonsVersionResponse{}
	_body, _err := client.DescribeClusterAddonsVersionWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonsVersionWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonsVersionResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterAddonsVersionResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterAddonsVersion"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/version"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeExternalAgent(ClusterId *string, request *DescribeExternalAgentRequest) (_result *DescribeExternalAgentResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeExternalAgentResponse{}
	_body, _err := client.DescribeExternalAgentWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeExternalAgentWithOptions(ClusterId *string, request *DescribeExternalAgentRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeExternalAgentResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeExternalAgentResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeExternalAgent"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/k8s/"+tea.StringValue(ClusterId)+"/external/agent/deployment"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UnInstallClusterAddons(ClusterId *string, request *UnInstallClusterAddonsRequest) (_result *UnInstallClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UnInstallClusterAddonsResponse{}
	_body, _err := client.UnInstallClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UnInstallClusterAddonsWithOptions(ClusterId *string, request *UnInstallClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UnInstallClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Addons),
	}
	_result = &UnInstallClusterAddonsResponse{}
	_body, _err := client.DoROARequest(tea.String("UnInstallClusterAddons"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/uninstall"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResumeComponentUpgrade(clusterid *string, componentid *string) (_result *ResumeComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ResumeComponentUpgradeResponse{}
	_body, _err := client.ResumeComponentUpgradeWithOptions(clusterid, componentid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResumeComponentUpgradeWithOptions(clusterid *string, componentid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ResumeComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &ResumeComponentUpgradeResponse{}
	_body, _err := client.DoROARequest(tea.String("ResumeComponentUpgrade"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(clusterid)+"/components/"+tea.StringValue(componentid)+"/resume"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClustersV1(request *DescribeClustersV1Request) (_result *DescribeClustersV1Response, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClustersV1Response{}
	_body, _err := client.DescribeClustersV1WithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClustersV1WithOptions(request *DescribeClustersV1Request, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClustersV1Response, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["cluster_type"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClustersV1Response{}
	_body, _err := client.DoROARequest(tea.String("DescribeClustersV1"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/api/v1/clusters"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterConfiguration(ClusterId *string, request *ModifyClusterConfigurationRequest) (_result *ModifyClusterConfigurationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterConfigurationResponse{}
	_body, _err := client.ModifyClusterConfigurationWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterConfigurationWithOptions(ClusterId *string, request *ModifyClusterConfigurationRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterConfigurationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CustomizeConfig)) {
		body["customize_config"] = request.CustomizeConfig
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ModifyClusterConfigurationResponse{}
	_body, _err := client.DoROARequest(tea.String("ModifyClusterConfiguration"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/configuration"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTaskInfo(taskId *string) (_result *DescribeTaskInfoResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTaskInfoResponse{}
	_body, _err := client.DescribeTaskInfoWithOptions(taskId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTaskInfoWithOptions(taskId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTaskInfoResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeTaskInfoResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeTaskInfo"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/tasks/"+tea.StringValue(taskId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescirbeWorkflow(workflowName *string) (_result *DescirbeWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescirbeWorkflowResponse{}
	_body, _err := client.DescirbeWorkflowWithOptions(workflowName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescirbeWorkflowWithOptions(workflowName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescirbeWorkflowResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescirbeWorkflowResponse{}
	_body, _err := client.DoROARequest(tea.String("DescirbeWorkflow"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/gs/workflow/"+tea.StringValue(workflowName)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelClusterUpgrade(ClusterId *string) (_result *CancelClusterUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelClusterUpgradeResponse{}
	_body, _err := client.CancelClusterUpgradeWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelClusterUpgradeWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelClusterUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &CancelClusterUpgradeResponse{}
	_body, _err := client.DoROARequest(tea.String("CancelClusterUpgrade"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/upgrade/cancel"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RemoveWorkflow(workflowName *string) (_result *RemoveWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RemoveWorkflowResponse{}
	_body, _err := client.RemoveWorkflowWithOptions(workflowName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RemoveWorkflowWithOptions(workflowName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RemoveWorkflowResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &RemoveWorkflowResponse{}
	_body, _err := client.DoROARequest(tea.String("RemoveWorkflow"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/gs/workflow/"+tea.StringValue(workflowName)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateTemplate(TemplateId *string, request *UpdateTemplateRequest) (_result *UpdateTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateTemplateResponse{}
	_body, _err := client.UpdateTemplateWithOptions(TemplateId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateTemplateWithOptions(TemplateId *string, request *UpdateTemplateRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateTemplateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Template)) {
		body["template"] = request.Template
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		body["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &UpdateTemplateResponse{}
	_body, _err := client.DoROARequest(tea.String("UpdateTemplate"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/templates/"+tea.StringValue(TemplateId)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeClusterAddons(ClusterId *string, request *UpgradeClusterAddonsRequest) (_result *UpgradeClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpgradeClusterAddonsResponse{}
	_body, _err := client.UpgradeClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeClusterAddonsWithOptions(ClusterId *string, request *UpgradeClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpgradeClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	_result = &UpgradeClusterAddonsResponse{}
	_body, _err := client.DoROARequest(tea.String("UpgradeClusterAddons"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/upgrade"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNamespaces(ClusterId *string) (_result *DescribeClusterNamespacesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNamespacesResponse{}
	_body, _err := client.DescribeClusterNamespacesWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNamespacesWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNamespacesResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterNamespacesResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterNamespaces"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/k8s/"+tea.StringValue(ClusterId)+"/namespaces"), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteKubernetesTrigger(Id *string) (_result *DeleteKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteKubernetesTriggerResponse{}
	_body, _err := client.DeleteKubernetesTriggerWithOptions(Id, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteKubernetesTriggerWithOptions(Id *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteKubernetesTriggerResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DeleteKubernetesTriggerResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteKubernetesTrigger"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/triggers/revoke/"+tea.StringValue(Id)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeUserQuota() (_result *DescribeUserQuotaResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeUserQuotaResponse{}
	_body, _err := client.DescribeUserQuotaWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeUserQuotaWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeUserQuotaResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeUserQuotaResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeUserQuota"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/quota"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteClusterNodepool(ClusterId *string, NodepoolId *string) (_result *DeleteClusterNodepoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterNodepoolResponse{}
	_body, _err := client.DeleteClusterNodepoolWithOptions(ClusterId, NodepoolId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteClusterNodepoolWithOptions(ClusterId *string, NodepoolId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterNodepoolResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DeleteClusterNodepoolResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteClusterNodepool"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("DELETE"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools/"+tea.StringValue(NodepoolId)), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonsUpgradeStatus(ClusterId *string, request *DescribeClusterAddonsUpgradeStatusRequest) (_result *DescribeClusterAddonsUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonsUpgradeStatusResponse{}
	_body, _err := client.DescribeClusterAddonsUpgradeStatusWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonsUpgradeStatusWithOptions(ClusterId *string, tmpReq *DescribeClusterAddonsUpgradeStatusRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonsUpgradeStatusResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &DescribeClusterAddonsUpgradeStatusShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.ComponentIds)) {
		request.ComponentIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.ComponentIds, tea.String("componentIds"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ComponentIdsShrink)) {
		query["componentIds"] = request.ComponentIdsShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClusterAddonsUpgradeStatusResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterAddonsUpgradeStatus"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/upgradestatus"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeWorkflows() (_result *DescribeWorkflowsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeWorkflowsResponse{}
	_body, _err := client.DescribeWorkflowsWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeWorkflowsWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeWorkflowsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeWorkflowsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeWorkflows"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/gs/workflows"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) InstallClusterAddons(ClusterId *string, request *InstallClusterAddonsRequest) (_result *InstallClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &InstallClusterAddonsResponse{}
	_body, _err := client.InstallClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InstallClusterAddonsWithOptions(ClusterId *string, request *InstallClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *InstallClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	_result = &InstallClusterAddonsResponse{}
	_body, _err := client.DoROARequest(tea.String("InstallClusterAddons"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/components/install"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodePools(ClusterId *string) (_result *DescribeClusterNodePoolsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodePoolsResponse{}
	_body, _err := client.DescribeClusterNodePoolsWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolsWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodePoolsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterNodePoolsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterNodePools"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodepools"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterV2UserKubeconfig(ClusterId *string, request *DescribeClusterV2UserKubeconfigRequest) (_result *DescribeClusterV2UserKubeconfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterV2UserKubeconfigResponse{}
	_body, _err := client.DescribeClusterV2UserKubeconfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterV2UserKubeconfigWithOptions(ClusterId *string, request *DescribeClusterV2UserKubeconfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterV2UserKubeconfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeClusterV2UserKubeconfigResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterV2UserKubeconfig"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/api/v2/k8s/"+tea.StringValue(ClusterId)+"/user_config"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) StartWorkflow(request *StartWorkflowRequest) (_result *StartWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &StartWorkflowResponse{}
	_body, _err := client.StartWorkflowWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) StartWorkflowWithOptions(request *StartWorkflowRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *StartWorkflowResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.WorkflowType)) {
		body["workflow_type"] = request.WorkflowType
	}

	if !tea.BoolValue(util.IsUnset(request.Service)) {
		body["service"] = request.Service
	}

	if !tea.BoolValue(util.IsUnset(request.MappingOssRegion)) {
		body["mapping_oss_region"] = request.MappingOssRegion
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqFirstFilename)) {
		body["mapping_fastq_first_filename"] = request.MappingFastqFirstFilename
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqSecondFilename)) {
		body["mapping_fastq_second_filename"] = request.MappingFastqSecondFilename
	}

	if !tea.BoolValue(util.IsUnset(request.MappingBucketName)) {
		body["mapping_bucket_name"] = request.MappingBucketName
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqPath)) {
		body["mapping_fastq_path"] = request.MappingFastqPath
	}

	if !tea.BoolValue(util.IsUnset(request.MappingReferencePath)) {
		body["mapping_reference_path"] = request.MappingReferencePath
	}

	if !tea.BoolValue(util.IsUnset(request.MappingIsMarkDup)) {
		body["mapping_is_mark_dup"] = request.MappingIsMarkDup
	}

	if !tea.BoolValue(util.IsUnset(request.MappingBamOutPath)) {
		body["mapping_bam_out_path"] = request.MappingBamOutPath
	}

	if !tea.BoolValue(util.IsUnset(request.MappingBamOutFilename)) {
		body["mapping_bam_out_filename"] = request.MappingBamOutFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsOssRegion)) {
		body["wgs_oss_region"] = request.WgsOssRegion
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqFirstFilename)) {
		body["wgs_fastq_first_filename"] = request.WgsFastqFirstFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqSecondFilename)) {
		body["wgs_fastq_second_filename"] = request.WgsFastqSecondFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsBucketName)) {
		body["wgs_bucket_name"] = request.WgsBucketName
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqPath)) {
		body["wgs_fastq_path"] = request.WgsFastqPath
	}

	if !tea.BoolValue(util.IsUnset(request.WgsReferencePath)) {
		body["wgs_reference_path"] = request.WgsReferencePath
	}

	if !tea.BoolValue(util.IsUnset(request.WgsVcfOutPath)) {
		body["wgs_vcf_out_path"] = request.WgsVcfOutPath
	}

	if !tea.BoolValue(util.IsUnset(request.WgsVcfOutFilename)) {
		body["wgs_vcf_out_filename"] = request.WgsVcfOutFilename
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &StartWorkflowResponse{}
	_body, _err := client.DoROARequest(tea.String("StartWorkflow"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/gs/workflow"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ScaleOutCluster(ClusterId *string, request *ScaleOutClusterRequest) (_result *ScaleOutClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleOutClusterResponse{}
	_body, _err := client.ScaleOutClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ScaleOutClusterWithOptions(ClusterId *string, request *ScaleOutClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleOutClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(tea.ToMap(request.Runtime))) {
		body["runtime"] = request.Runtime
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &ScaleOutClusterResponse{}
	_body, _err := client.DoROARequest(tea.String("ScaleOutCluster"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEvents(request *DescribeEventsRequest) (_result *DescribeEventsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEventsResponse{}
	_body, _err := client.DescribeEventsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEventsWithOptions(request *DescribeEventsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEventsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterId)) {
		query["cluster_id"] = request.ClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["type"] = request.Type
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &DescribeEventsResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeEvents"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/events"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateK8sClusterUserConfigExpire(ClusterId *string) (_result *UpdateK8sClusterUserConfigExpireResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateK8sClusterUserConfigExpireResponse{}
	_body, _err := client.UpdateK8sClusterUserConfigExpireWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateK8sClusterUserConfigExpireWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateK8sClusterUserConfigExpireResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &UpdateK8sClusterUserConfigExpireResponse{}
	_body, _err := client.DoROARequest(tea.String("UpdateK8sClusterUserConfigExpire"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/k8s/"+tea.StringValue(ClusterId)+"/user_config/expire"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TagResources(request *TagResourcesRequest) (_result *TagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &TagResourcesResponse{}
	_body, _err := client.TagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TagResourcesWithOptions(request *TagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *TagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ResourceIds)) {
		body["resource_ids"] = request.ResourceIds
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		body["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		body["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &TagResourcesResponse{}
	_body, _err := client.DoROARequest(tea.String("TagResources"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("PUT"), tea.String("AK"), tea.String("/tags"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterTags(ClusterId *string, request *ModifyClusterTagsRequest) (_result *ModifyClusterTagsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterTagsResponse{}
	_body, _err := client.ModifyClusterTagsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterTagsWithOptions(ClusterId *string, request *ModifyClusterTagsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterTagsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	_result = &ModifyClusterTagsResponse{}
	_body, _err := client.DoROARequest(tea.String("ModifyClusterTags"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/tags"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetKubernetesTrigger(ClusterId *string, request *GetKubernetesTriggerRequest) (_result *GetKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GetKubernetesTriggerResponse{}
	_body, _err := client.GetKubernetesTriggerWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetKubernetesTriggerWithOptions(ClusterId *string, request *GetKubernetesTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GetKubernetesTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Namespace)) {
		query["Namespace"] = request.Namespace
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["Type"] = request.Type
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["Name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Action)) {
		query["action"] = request.Action
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	_result = &GetKubernetesTriggerResponse{}
	_body, _err := client.DoROARequest(tea.String("GetKubernetesTrigger"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/triggers/"+tea.StringValue(ClusterId)), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetUpgradeStatus(ClusterId *string) (_result *GetUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GetUpgradeStatusResponse{}
	_body, _err := client.GetUpgradeStatusWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetUpgradeStatusWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GetUpgradeStatusResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &GetUpgradeStatusResponse{}
	_body, _err := client.DoROARequest(tea.String("GetUpgradeStatus"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/api/v2/clusters/"+tea.StringValue(ClusterId)+"/upgrade/status"), tea.String("json"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterResources(ClusterId *string) (_result *DescribeClusterResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterResourcesResponse{}
	_body, _err := client.DescribeClusterResourcesWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterResourcesWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterResourcesResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	_result = &DescribeClusterResourcesResponse{}
	_body, _err := client.DoROARequest(tea.String("DescribeClusterResources"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("GET"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/resources"), tea.String("array"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteClusterNodes(ClusterId *string, request *DeleteClusterNodesRequest) (_result *DeleteClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterNodesResponse{}
	_body, _err := client.DeleteClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteClusterNodesWithOptions(ClusterId *string, request *DeleteClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DrainNode)) {
		body["drain_node"] = request.DrainNode
	}

	if !tea.BoolValue(util.IsUnset(request.ReleaseNode)) {
		body["release_node"] = request.ReleaseNode
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	_result = &DeleteClusterNodesResponse{}
	_body, _err := client.DoROARequest(tea.String("DeleteClusterNodes"), tea.String("2015-12-15"), tea.String("HTTPS"), tea.String("POST"), tea.String("AK"), tea.String("/clusters/"+tea.StringValue(ClusterId)+"/nodes"), tea.String("none"), req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}
