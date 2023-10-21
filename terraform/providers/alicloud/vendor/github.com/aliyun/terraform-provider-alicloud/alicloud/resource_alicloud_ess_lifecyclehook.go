package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssLifeCycleHookCreate,
		Read:   resourceAliyunEssLifeCycleHookRead,
		Update: resourceAliyunEssLifeCycleHookUpdate,
		Delete: resourceAliyunEssLifeCycleHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lifecycle_transition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SCALE_IN", "SCALE_OUT"}, false),
			},
			"heartbeat_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validation.IntBetween(30, 21600),
			},
			"default_result": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CONTINUE",
				ValidateFunc: validation.StringInSlice([]string{"CONTINUE", "ABANDON"}, false),
			},
			"notification_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notification_metadata": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEssLifeCycleHookCreate(d *schema.ResourceData, meta interface{}) error {

	request := buildAlicloudEssLifeCycleHookArgs(d)
	client := meta.(*connectivity.AliyunClient)
	request.RegionId = client.RegionId
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateLifecycleHook(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateLifecycleHookResponse)
		d.SetId(response.LifecycleHookId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_lifecyclehook", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunEssLifeCycleHookRead(d, meta)
}

func resourceAliyunEssLifeCycleHookRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssLifecycleHook(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("name", object.LifecycleHookName)
	d.Set("lifecycle_transition", object.LifecycleTransition)
	d.Set("heartbeat_timeout", object.HeartbeatTimeout)
	d.Set("default_result", object.DefaultResult)
	d.Set("notification_arn", object.NotificationArn)
	d.Set("notification_metadata", object.NotificationMetadata)

	return nil
}

func resourceAliyunEssLifeCycleHookUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyLifecycleHookRequest()
	request.LifecycleHookId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("lifecycle_transition") {
		request.LifecycleTransition = d.Get("lifecycle_transition").(string)
	}

	if d.HasChange("heartbeat_timeout") {
		request.HeartbeatTimeout = requests.NewInteger(d.Get("heartbeat_timeout").(int))
	}

	if d.HasChange("default_result") {
		request.DefaultResult = d.Get("default_result").(string)
	}

	if d.HasChange("notification_arn") {
		request.NotificationArn = d.Get("notification_arn").(string)
	}

	if d.HasChange("notification_metadata") {
		request.NotificationMetadata = d.Get("notification_metadata").(string)
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyLifecycleHook(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAliyunEssLifeCycleHookRead(d, meta)
}

func resourceAliyunEssLifeCycleHookDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	request := ess.CreateDeleteLifecycleHookRequest()
	request.LifecycleHookId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteLifecycleHook(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLifecycleHookId.NotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(essService.WaitForEssLifecycleHook(d.Id(), Deleted, DefaultTimeout))

}

func buildAlicloudEssLifeCycleHookArgs(d *schema.ResourceData) *ess.CreateLifecycleHookRequest {
	request := ess.CreateCreateLifecycleHookRequest()

	request.ScalingGroupId = d.Get("scaling_group_id").(string)

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.LifecycleHookName = v.(string)
	}

	if transition := d.Get("lifecycle_transition").(string); transition != "" {
		request.LifecycleTransition = transition
	}

	if timeout, ok := d.GetOk("heartbeat_timeout"); ok && timeout.(int) > 0 {
		request.HeartbeatTimeout = requests.NewInteger(timeout.(int))
	}

	if v, ok := d.GetOk("default_result"); ok && v.(string) != "" {
		request.DefaultResult = v.(string)
	}

	if v, ok := d.GetOk("notification_arn"); ok && v.(string) != "" {
		request.NotificationArn = v.(string)
	}

	if v, ok := d.GetOk("notification_metadata"); ok && v.(string) != "" {
		request.NotificationMetadata = v.(string)
	}

	return request
}
