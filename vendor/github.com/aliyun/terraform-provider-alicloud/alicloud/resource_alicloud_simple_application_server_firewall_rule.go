package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSimpleApplicationServerFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSimpleApplicationServerFirewallRuleCreate,
		Read:   resourceAlicloudSimpleApplicationServerFirewallRuleRead,
		Delete: resourceAlicloudSimpleApplicationServerFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"firewall_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rule_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Tcp", "TcpAndUdp", "Udp"}, false),
			},
		},
	}
}

func resourceAlicloudSimpleApplicationServerFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFirewallRule"
	request := make(map[string]interface{})
	conn, err := client.NewSwasClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceId"] = d.Get("instance_id")
	request["Port"] = d.Get("port")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}
	request["RuleProtocol"] = d.Get("rule_protocol")
	request["ClientToken"] = buildClientToken("CreateFirewallRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_simple_application_server_firewall_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["FirewallId"]))

	return resourceAlicloudSimpleApplicationServerFirewallRuleRead(d, meta)
}
func resourceAlicloudSimpleApplicationServerFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := SwasOpenService{client}
	object, err := swasOpenService.DescribeSimpleApplicationServerFirewallRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_firewall_rule swasOpenService.DescribeSimpleApplicationServerFirewallRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("firewall_rule_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("port", object["Port"])
	d.Set("remark", object["Remark"])
	d.Set("rule_protocol", convertSimpleApplicationServerFirewallRuleRuleProtocolResponse(object["RuleProtocol"].(string)))
	return nil
}
func resourceAlicloudSimpleApplicationServerFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteFirewallRule"
	var response map[string]interface{}
	conn, err := client.NewSwasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RuleId":     parts[1],
		"InstanceId": parts[0],
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteFirewallRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertSimpleApplicationServerFirewallRuleRuleProtocolResponse(source string) string {
	switch source {
	case "TCP":
		return "Tcp"
	case "UDP":
		return "Udp"
	case "TCP+UDP":
		return "TcpAndUdp"
	}
	return source
}
