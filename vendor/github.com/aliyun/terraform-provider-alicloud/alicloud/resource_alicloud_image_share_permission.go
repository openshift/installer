package alicloud

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageSharePermissionCreate,
		Read:   resourceAliCloudImageSharePermissionRead,
		Delete: resourceAliCloudImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	imageId := d.Get("image_id").(string)
	accountId := d.Get("account_id").(string)
	request := ecs.CreateModifyImageSharePermissionRequest()
	request.RegionId = client.RegionId
	request.ImageId = imageId
	accountSli := []string{accountId}
	request.AddAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_share_permission", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(imageId + ":" + accountId)
	return resourceAliCloudImageSharePermissionRead(d, meta)
}

func resourceAliCloudImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client: client}
	object, err := ecsService.DescribeImageShareByImageId(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_image_share_permission ecsService.DescribeImageShareByImageId Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	d.Set("image_id", object.ImageId)
	d.Set("account_id", parts[1])
	return WrapError(err)
}

func resourceAliCloudImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateModifyImageSharePermissionRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 2)
	request.ImageId = parts[0]
	accountSli := []string{parts[1]}
	request.RemoveAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_share_permission", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
