package alicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSagAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagAclRuleCreate,
		Read:   resourceAlicloudSagAclRuleRead,
		Update: resourceAlicloudSagAclRuleUpdate,
		Delete: resourceAlicloudSagAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudSagAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateAddACLRuleRequest()
	request.AclId = d.Get("acl_id").(string)
	request.Policy = d.Get("policy").(string)
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.Direction = d.Get("direction").(string)
	request.SourceCidr = d.Get("source_cidr").(string)
	request.SourcePortRange = d.Get("source_port_range").(string)
	request.DestCidr = d.Get("dest_cidr").(string)
	request.DestPortRange = d.Get("dest_port_range").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("priority"); ok && v.(int) != 0 {
		request.Priority = requests.NewInteger(v.(int))
	}

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.AddACLRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_acl_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.AddACLRuleResponse)
	d.SetId(fmt.Sprintf("%s%s%s", response.AclId, COLON_SEPARATED, response.AcrId))

	return resourceAlicloudSagAclRuleRead(d, meta)
}

func resourceAlicloudSagAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagAclRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("acl_id", object.AclId)
	d.Set("description", object.Description)
	d.Set("policy", object.Policy)
	d.Set("ip_protocol", object.IpProtocol)
	d.Set("direction", object.Direction)
	d.Set("source_cidr", object.SourceCidr)
	d.Set("source_port_range", object.SourcePortRange)
	d.Set("dest_cidr", object.DestCidr)
	d.Set("dest_port_range", object.DestPortRange)
	d.Set("priority", object.Priority)

	return nil
}

func resourceAlicloudSagAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartag.CreateModifyACLRuleRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.AclId = parts[0]
	request.AcrId = parts[1]

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("policy") {
		request.Policy = d.Get("policy").(string)
		update = true
	}
	if d.HasChange("ip_protocol") {
		request.IpProtocol = d.Get("ip_protocol").(string)
		update = true
	}
	if d.HasChange("direction") {
		request.Direction = d.Get("direction").(string)
		update = true
	}
	if d.HasChange("source_cidr") {
		request.SourceCidr = d.Get("source_cidr").(string)
		update = true
	}
	if d.HasChange("source_port_range") {
		request.SourcePortRange = d.Get("source_port_range").(string)
		update = true
	}
	if d.HasChange("dest_cidr") {
		request.DestCidr = d.Get("dest_cidr").(string)
		update = true
	}
	if d.HasChange("dest_port_range") {
		request.DestPortRange = d.Get("dest_port_range").(string)
		update = true
	}
	if d.HasChange("priority") {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
		update = true
	}

	if update {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.ModifyACLRule(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagAclRuleRead(d, meta)
}

func resourceAlicloudSagAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteACLRuleRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.AclId = parts[0]
	request.AcrId = parts[1]

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.DeleteACLRule(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagACLRuleId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForSagAclRule(d.Id(), Deleted, DefaultTimeoutMedium))
}
