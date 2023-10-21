package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSagQos() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagQosCreate,
		Read:   resourceAlicloudSagQosRead,
		Update: resourceAlicloudSagQosUpdate,
		Delete: resourceAlicloudSagQosDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
		},
	}
}

func resourceAlicloudSagQosCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateQosRequest()

	request.QosName = d.Get("name").(string)

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.CreateQos(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_qos", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.CreateQosResponse)
	d.SetId(response.QosId)
	return resourceAlicloudSagQosRead(d, meta)
}

func resourceAlicloudSagQosRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagQos(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.QosName)

	return nil
}

func resourceAlicloudSagQosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := smartag.CreateModifyQosRequest()
		request.QosId = d.Id()
		request.QosName = d.Get("name").(string)
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.ModifyQos(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagQosRead(d, meta)
}

func resourceAlicloudSagQosDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteQosRequest()
	request.QosId = d.Id()

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.DeleteQos(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagQosId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagQos(d.Id(), Deleted, DefaultTimeoutMedium))
}
