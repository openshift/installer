package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasInstanceClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasInstanceClusterAttachmentCreate,
		Read:   resourceAlicloudEdasInstanceClusterAttachmentRead,
		Delete: resourceAlicloudEdasInstanceClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			"status_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
				ForceNew: true,
			},
			"ecu_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
			"cluster_member_ids": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasInstanceClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	clusterId := d.Get("cluster_id").(string)
	instanceIds := d.Get("instance_ids").([]interface{})
	aString := make([]string, len(instanceIds))
	for i, v := range instanceIds {
		aString[i] = v.(string)
	}

	request := edas.CreateInstallAgentRequest()
	request.ClusterId = clusterId
	request.RegionId = client.RegionId

	request.InstanceIds = strings.Join(aString, ",")
	request.SetReadTimeout(30 * time.Second)

	if err := edasService.SyncResource("ecs"); err != nil {
		return err
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.InstallAgent(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ := raw.(*edas.InstallAgentResponse)

		if response.Code != 200 {
			return resource.NonRetryableError(Error("insert instances to cluster failed for " + response.Message))
		}

		var instanceIdFailed []string
		for _, result := range response.ExecutionResultList.ExecutionResult {
			if !result.Success {
				instanceIdFailed = append(instanceIdFailed, result.InstanceId)
			}
		}
		if len(instanceIdFailed) > 0 {
			err = Error("instances still import failed, try again")
			request.InstanceIds = strings.Join(instanceIdFailed, ",")
			return resource.RetryableError(err)
		}

		d.SetId(clusterId + ":" + strings.Join(aString, ","))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_cluster_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudEdasInstanceClusterAttachmentRead(d, meta)
}

func resourceAlicloudEdasInstanceClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	clusterId := strs[0]
	regionId := client.RegionId
	instanceIdstr := strs[1]

	request := edas.CreateListClusterMembersRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListClusterMembers(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_cluster_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	statusMap := make(map[string]int)
	ecuMap := make(map[string]string)
	memMap := make(map[string]string)
	response := raw.(*edas.ListClusterMembersResponse)
	for _, member := range response.ClusterMemberPage.ClusterMemberList.ClusterMember {
		if strings.Contains(instanceIdstr, member.EcsId) {
			statusMap[member.EcsId] = member.Status
			ecuMap[member.EcsId] = member.EcuId
			memMap[member.EcsId] = member.ClusterMemberId
		}
	}

	d.Set("status_map", statusMap)
	d.Set("ecu_map", ecuMap)
	d.Set("cluster_member_ids", memMap)

	return nil
}

//有问题 单个实例删除失败会影响整个过程
func resourceAlicloudEdasInstanceClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	memIds := d.Get("cluster_member_ids").(map[string]interface{})
	for instanceId, memberId := range memIds {
		request := edas.CreateDeleteClusterMemberRequest()
		request.RegionId = client.RegionId
		request.ClusterId = d.Get("cluster_id").(string)
		request.ClusterMemberId = memberId.(string)

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteClusterMember(request)

			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RoaRequest, request)
			response, _ := raw.(*edas.DeleteClusterMemberResponse)
			if strings.Contains(response.Message, "there are still applications deployed in this cluster") {
				err = Error("there are still applications deployed in this cluster")
				return resource.RetryableError(err)
			} else if response.Code != 200 {
				return resource.NonRetryableError(Error("delete instance:" + instanceId + " from cluster failed for " + response.Message))
			}

			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_instance_cluster_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}
