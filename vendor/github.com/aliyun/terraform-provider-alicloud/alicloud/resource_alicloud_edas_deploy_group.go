package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasDeployGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasDeployGroupCreate,
		Read:   resourceAlicloudEdasDeployGroupRead,
		Delete: resourceAlicloudEdasDeployGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_type": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEdasDeployGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	regionId := client.RegionId
	groupName := d.Get("group_name").(string)

	request := edas.CreateInsertDeployGroupRequest()
	request.RegionId = regionId
	request.AppId = appId
	request.GroupName = groupName

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.InsertDeployGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response := raw.(*edas.InsertDeployGroupResponse)
		deployGroup := response.DeployGroupEntity
		d.SetId(appId + ":" + groupName + ":" + deployGroup.Id)
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_deploy_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudEdasDeployGroupRead(d, meta)
}

func resourceAlicloudEdasDeployGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	strs, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	appId := strs[0]
	groupId := strs[2]

	deployGroup, err := edasService.GetDeployGroup(appId, groupId)
	if err != nil {
		return WrapError(err)
	}
	if deployGroup == nil {
		return nil
	}

	d.Set("group_type", deployGroup.GroupType)
	d.Set("app_id", deployGroup.AppId)
	d.Set("group_name", deployGroup.GroupName)

	return nil
}

func resourceAlicloudEdasDeployGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateDeleteDeployGroupRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Get("app_id").(string)
	request.GroupName = d.Get("group_name").(string)

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteDeployGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}
