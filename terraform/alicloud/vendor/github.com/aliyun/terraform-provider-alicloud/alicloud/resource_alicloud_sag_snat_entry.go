package alicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSagSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagSnatCreate,
		Read:   resourceAlicloudSagSnatRead,
		Delete: resourceAlicloudSagSnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"snat_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
		},
	}
}

func resourceAlicloudSagSnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateAddSnatEntryRequest()

	request.SmartAGId = d.Get("sag_id").(string)
	request.CidrBlock = d.Get("cidr_block").(string)
	request.SnatIp = d.Get("snat_ip").(string)

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.AddSnatEntry(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_snat_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.AddSnatEntryResponse)
	d.SetId(fmt.Sprintf("%s%s%s", request.SmartAGId, COLON_SEPARATED, response.InstanceId))

	return resourceAlicloudSagSnatRead(d, meta)
}

func resourceAlicloudSagSnatRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagSnatEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, _ := ParseResourceId(d.Id(), 2)
	d.Set("sag_id", parts[0])
	d.Set("cidr_block", object.CidrBlock)
	d.Set("snat_ip", object.SnatIp)

	return nil
}

func resourceAlicloudSagSnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteSnatEntryRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.SmartAGId = parts[0]
	request.InstanceId = parts[1]

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.DeleteSnatEntry(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagSnatId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagSnatEntry(d.Id(), Deleted, DefaultTimeoutMedium))
}
