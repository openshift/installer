package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPvtzRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzRuleAttachmentCreate,
		Read:   resourceAlicloudPvtzRuleAttachmentRead,
		Update: resourceAlicloudPvtzRuleAttachmentUpdate,
		Delete: resourceAlicloudPvtzRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpcs": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudPvtzRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "BindResolverRuleVpc"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request["Lang"] = "en"
	request["RuleId"] = d.Get("rule_id")

	for k, vpcConfig := range d.Get("vpcs").(*schema.Set).List() {
		vpcConfigArg := vpcConfig.(map[string]interface{})
		request[fmt.Sprintf("Vpc.%d.VpcId", k+1)] = vpcConfigArg["vpc_id"]
		request[fmt.Sprintf("Vpc.%d.RegionId", k+1)] = vpcConfigArg["region_id"]
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_rule_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleId"]))

	return resourceAlicloudPvtzRuleAttachmentRead(d, meta)
}
func resourceAlicloudPvtzRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzRuleAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_rule_attachment pvtzService.DescribePvtzRuleAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	vpcsSli := make([]map[string]interface{}, 0)
	if vpcs, ok := object["BindVpcs"].([]interface{}); ok {
		for _, vpcConfigArgs := range vpcs {
			vpcConfigArg := vpcConfigArgs.(map[string]interface{})
			vpcsMap := make(map[string]interface{})
			vpcsMap["region_id"] = vpcConfigArg["RegionId"]
			vpcsMap["vpc_id"] = vpcConfigArg["VpcId"]
			vpcsSli = append(vpcsSli, vpcsMap)
		}
	}
	d.Set("vpcs", vpcsSli)
	d.Set("rule_id", d.Id())
	return nil
}
func resourceAlicloudPvtzRuleAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("vpcs") {
		update = true
		for k, vpcConfig := range d.Get("vpcs").(*schema.Set).List() {
			vpcConfigArg := vpcConfig.(map[string]interface{})
			request[fmt.Sprintf("Vpc.%d.VpcId", k+1)] = vpcConfigArg["vpc_id"]
			request[fmt.Sprintf("Vpc.%d.RegionId", k+1)] = vpcConfigArg["region_id"]
		}
	}

	if update {
		request["Lang"] = "en"
		action := "BindResolverRuleVpc"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudPvtzRuleAttachmentRead(d, meta)
}
func resourceAlicloudPvtzRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "BindResolverRuleVpc"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}

	request["Lang"] = "en"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
