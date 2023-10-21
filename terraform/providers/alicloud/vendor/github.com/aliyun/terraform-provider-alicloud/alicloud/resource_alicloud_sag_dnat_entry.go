package alicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSagDnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagDnatCreate,
		Read:   resourceAlicloudSagDnatRead,
		Delete: resourceAlicloudSagDnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Intranet", "Internet"}, false),
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "any"}, false),
			},
			"external_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.SingleIP(),
				DiffSuppressFunc: sagDnatEntryTypeDiffSuppressFunc,
			},
			"external_port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"internal_port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSagDnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateAddDnatEntryRequest()

	request.SagId = d.Get("sag_id").(string)
	request.Type = d.Get("type").(string)
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.ExternalIp = d.Get("external_ip").(string)
	request.InternalIp = d.Get("internal_ip").(string)
	request.ExternalPort = d.Get("external_port").(string)
	request.InternalPort = d.Get("internal_port").(string)

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.AddDnatEntry(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_dnat_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.AddDnatEntryResponse)
	d.SetId(fmt.Sprintf("%s%s%s", request.SagId, COLON_SEPARATED, response.DnatEntryId))

	return resourceAlicloudSagDnatRead(d, meta)
}

func resourceAlicloudSagDnatRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagDnatEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("sag_id", object.SagId)
	d.Set("type", object.Type)
	d.Set("ip_protocol", object.IpProtocol)
	d.Set("external_ip", object.ExternalIp)
	d.Set("external_port", object.ExternalPort)
	d.Set("internal_ip", object.InternalIp)
	d.Set("internal_port", object.InternalPort)

	return nil
}

func resourceAlicloudSagDnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteDnatEntryRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.SagId = parts[0]
	request.DnatEntryId = parts[1]

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.DeleteDnatEntry(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagClientId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagDnatEntry(d.Id(), Deleted, DefaultTimeoutMedium))
}
