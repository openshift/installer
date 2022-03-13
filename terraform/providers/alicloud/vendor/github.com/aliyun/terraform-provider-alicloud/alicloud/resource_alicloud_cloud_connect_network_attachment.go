package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudConnectNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkAttachmentCreate,
		Read:   resourceAlicloudCloudConnectNetworkAttachmentRead,
		Delete: resourceAlicloudCloudConnectNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateBindSmartAccessGatewayRequest()

	request.RegionId = client.RegionId
	request.CcnId = d.Get("ccn_id").(string)
	request.SmartAGId = d.Get("sag_id").(string)
	var err error
	var raw interface{}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.BindSmartAccessGateway(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "UnknownError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_connect_network_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s%s%s", request.CcnId, COLON_SEPARATED, request.SmartAGId))

	return resourceAlicloudCloudConnectNetworkAttachmentRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeCloudConnectNetworkAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ccn_id", object.AssociatedCcnId)
	d.Set("sag_id", object.SmartAGId)

	return nil
}

func resourceAlicloudCloudConnectNetworkAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	ccnId := parts[0]
	sagId := parts[1]

	request := smartag.CreateUnbindSmartAccessGatewayRequest()
	request.RegionId = client.RegionId
	request.CcnId = ccnId
	request.SmartAGId = sagId

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.UnbindSmartAccessGateway(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus", "TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForCloudConnectNetworkAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
