package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDatahubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubSubscriptionCreate,
		Read:   resourceAliyunDatahubSubscriptionRead,
		Update: resourceAliyunDatahubSubscriptionUpdate,
		Delete: resourceAliyunDatahubSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "subscription added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"sub_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	var requestInfo *datahub.DataHub

	raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
		requestInfo = dataHubClient.(*datahub.DataHub)
		return dataHubClient.CreateSubscription(projectName, topicName, subComment)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_datahub_subscription", "CreateSubscription", AliyunDatahubSdkGo)
	}
	subscriptionResult := raw.(*datahub.CreateSubscriptionResult)
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubComment"] = subComment
		addDebug("CreateSubscription", raw, requestInfo, requestMap)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(topicName), COLON_SEPARATED, subscriptionResult.SubId))
	return resourceAliyunDatahubSubscriptionRead(d, meta)
}

func resourceAliyunDatahubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	projectName := parts[0]
	object, err := datahubService.DescribeDatahubSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(object.TopicName), COLON_SEPARATED, object.SubId))

	d.Set("project_name", projectName)
	d.Set("topic_name", object.TopicName)
	d.Set("sub_id", object.SubId)
	d.Set("comment", object.Comment)
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))
	return nil
}

func resourceAliyunDatahubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("comment") {
		subComment := d.Get("comment").(string)

		var requestInfo *datahub.DataHub

		raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			requestInfo = dataHubClient.(*datahub.DataHub)
			return dataHubClient.UpdateSubscription(projectName, topicName, subId, subComment)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSubscription", AliyunDatahubSdkGo)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			requestMap["SubId"] = subId
			requestMap["SubComment"] = subComment
			addDebug("UpdateSubscription", raw, requestInfo, requestMap)
		}
	}

	return resourceAliyunDatahubSubscriptionRead(d, meta)
}

func resourceAliyunDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	var requestInfo *datahub.DataHub

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			requestInfo = dataHubClient.(*datahub.DataHub)
			return dataHubClient.DeleteSubscription(projectName, topicName, subId)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			requestMap["SubId"] = subId
			addDebug("DeleteSubscription", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteSubscription", AliyunDatahubSdkGo)
	}
	return WrapError(datahubService.WaitForDatahubSubscription(d.Id(), Deleted, DefaultTimeout))
}
