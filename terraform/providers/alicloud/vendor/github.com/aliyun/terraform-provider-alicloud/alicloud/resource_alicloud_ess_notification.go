package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEssNotificationCreate,
		Read:   resourceAlicloudEssNotificationRead,
		Update: resourceAlicloudEssNotificationUpdate,
		Delete: resourceAlicloudEssNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"notification_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notification_types": {
				Required: true,
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEssNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ess.CreateCreateNotificationConfigurationRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	request.NotificationArn = d.Get("notification_arn").(string)
	if v, ok := d.GetOk("notification_types"); ok {
		notificationTypes := make([]string, 0)
		notificationTypeList := v.(*schema.Set).List()
		if len(notificationTypeList) > 0 {
			for _, n := range notificationTypeList {
				notificationTypes = append(notificationTypes, n.(string))
			}
		}
		if len(notificationTypes) > 0 {
			request.NotificationType = &notificationTypes
		}
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateNotificationConfiguration(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_notification", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(fmt.Sprintf("%s:%s", request.ScalingGroupId, request.NotificationArn))
	return resourceAlicloudEssNotificationRead(d, meta)
}

func resourceAlicloudEssNotificationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssNotification(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("notification_arn", object.NotificationArn)
	d.Set("notification_types", object.NotificationTypes.NotificationType)
	return nil
}

func resourceAlicloudEssNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyNotificationConfigurationRequest()
	request.RegionId = client.RegionId
	parts := strings.SplitN(d.Id(), ":", 2)
	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]
	if d.HasChange("notification_types") {
		v := d.Get("notification_types")
		notificationTypes := make([]string, 0)
		notificationTypeList := v.(*schema.Set).List()
		if len(notificationTypeList) > 0 {
			for _, n := range notificationTypeList {
				notificationTypes = append(notificationTypes, n.(string))
			}
		}
		if len(notificationTypes) > 0 {
			request.NotificationType = &notificationTypes
		}
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyNotificationConfiguration(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlicloudEssNotificationRead(d, meta)
}

func resourceAlicloudEssNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	request := ess.CreateDeleteNotificationConfigurationRequest()
	request.RegionId = client.RegionId
	parts := strings.SplitN(d.Id(), ":", 2)

	request.ScalingGroupId = parts[0]
	request.NotificationArn = parts[1]

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteNotificationConfiguration(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssNotification(d.Id(), Deleted, DefaultTimeout))
}
