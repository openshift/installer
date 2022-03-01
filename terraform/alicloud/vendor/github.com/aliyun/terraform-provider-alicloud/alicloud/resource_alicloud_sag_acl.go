package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSagAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagAclCreate,
		Read:   resourceAlicloudSagAclRead,
		Update: resourceAlicloudSagAclUpdate,
		Delete: resourceAlicloudSagAclDelete,
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

func resourceAlicloudSagAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateACLRequest()

	request.Name = d.Get("name").(string)

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.CreateACL(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_acl", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.CreateACLResponse)
	d.SetId(response.AclId)
	return resourceAlicloudSagAclRead(d, meta)
}

func resourceAlicloudSagAclRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	return nil
}

func resourceAlicloudSagAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := smartag.CreateModifyACLRequest()
		request.AclId = d.Id()
		request.Name = d.Get("name").(string)

		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.ModifyACL(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagAclRead(d, meta)
}

func resourceAlicloudSagAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteACLRequest()
	request.AclId = d.Id()

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.DeleteACL(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagACLId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
