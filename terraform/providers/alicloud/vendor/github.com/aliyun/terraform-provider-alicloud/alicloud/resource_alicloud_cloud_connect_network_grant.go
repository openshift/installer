package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudConnectNetworkGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkGrantCreate,
		Read:   resourceAlicloudCloudConnectNetworkGrantRead,
		Delete: resourceAlicloudCloudConnectNetworkGrantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkGrantCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateGrantInstanceToCbnRequest()

	request.RegionId = client.RegionId
	request.CcnInstanceId = d.Get("ccn_id").(string)
	request.CenInstanceId = d.Get("cen_id").(string)
	request.CenUid = d.Get("cen_uid").(requests.Integer)
	var err error
	var raw interface{}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.GrantInstanceToCbn(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_connect_network_grant", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(fmt.Sprintf("%s%s%s", request.CcnInstanceId, COLON_SEPARATED, request.CenInstanceId))

	return resourceAlicloudCloudConnectNetworkGrantRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkGrantRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeCloudConnectNetworkGrant(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ccn_id", object.CcnInstanceId)
	d.Set("cen_id", object.CenInstanceId)
	d.Set("cen_uid", object.CenUid)

	return nil
}

func resourceAlicloudCloudConnectNetworkGrantDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := smartag.CreateRevokeInstanceFromCbnRequest()
	request.RegionId = client.RegionId
	request.CcnInstanceId = parts[0]
	request.CenInstanceId = parts[1]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.RevokeInstanceFromCbn(request)
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
	return WrapError(sagService.WaitForCloudConnectNetworkGrant(d.Id(), Deleted, DefaultTimeoutMedium))
}
