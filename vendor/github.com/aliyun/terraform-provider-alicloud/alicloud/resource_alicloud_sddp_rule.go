package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSddpRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSddpRuleCreate,
		Read:   resourceAlicloudSddpRuleRead,
		Update: resourceAlicloudSddpRuleUpdate,
		Delete: resourceAlicloudSddpRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 2}),
				ForceNew:     true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content_category": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"custom_type": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
				Default:      "zh",
				Optional:     true,
			},
			"product_code": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ODPS", "OSS", "RDS"}, false),
				Optional:     true,
			},
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "2", "5"}, false),
			},
			"risk_level_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2", "3", "4", "5"}, false),
			},
			"rule_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"stat_express": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"warn_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
		},
	}
}

func resourceAlicloudSddpRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRule"
	request := make(map[string]interface{})
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	request["Category"] = d.Get("category").(int)
	request["Content"] = d.Get("content")
	if v, ok := d.GetOk("content_category"); ok {
		request["ContentCategory"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}
	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
	}
	if v, ok := d.GetOk("risk_level_id"); ok {
		request["RiskLevelId"] = v
	}
	request["Name"] = d.Get("rule_name")
	if v, ok := d.GetOk("rule_type"); ok {
		request["RuleType"] = v
	}
	if v, ok := d.GetOk("stat_express"); ok {
		request["StatExpress"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("target"); ok {
		request["Target"] = v
	}
	if v, ok := d.GetOk("warn_level"); ok {
		request["WarnLevel"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sddp_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudSddpRuleUpdate(d, meta)
}
func resourceAlicloudSddpRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sddpService := SddpService{client}
	object, err := sddpService.DescribeSddpRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sddp_rule sddpService.DescribeSddpRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["Category"]; ok {
		err = d.Set("category", formatInt(v))
		if err != nil {
			return WrapError(err)
		}
	}
	err = d.Set("content", object["Content"])
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("content_category", object["ContentCategory"])
	if err != nil {
		return WrapError(err)
	}
	if v, ok := object["CustomType"]; ok {
		err = d.Set("custom_type", formatInt(v))
		if err != nil {
			return WrapError(err)
		}
	}
	err = d.Set("description", object["Description"])
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("product_code", object["ProductCode"])
	if err != nil {
		return WrapError(err)
	}
	if v, ok := object["ProductId"]; ok {
		v = formatInt(v)
		err = d.Set("product_id", fmt.Sprint(v))
		if err != nil {
			return WrapError(err)
		}
	}
	if v, ok := object["RiskLevelId"]; ok {
		v = formatInt(v)
		err = d.Set("risk_level_id", fmt.Sprint(v))
		if err != nil {
			return WrapError(err)
		}
	}
	if v, ok := object["Name"]; ok {
		err = d.Set("rule_name", v.(string))
		if err != nil {
			return WrapError(err)
		}
	}
	err = d.Set("stat_express", object["StatExpress"])
	if err != nil {
		return WrapError(err)
	}
	if v, ok := object["Status"]; ok {
		d.Set("status", formatInt(v))
	}
	err = d.Set("target", object["Target"])
	if err != nil {
		return WrapError(err)
	}
	if v, ok := object["WarnLevel"]; ok {
		err = d.Set("warn_level", formatInt(v))
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
func resourceAlicloudSddpRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("status") {
		update = true
		_, new := d.GetChange("status")
		if new != nil {
			request["Status"] = new.(int)
		}
	}
	if update {
		action := "ModifyRuleStatus"
		conn, err := client.NewSddpClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("status")
		d.SetPartial("lang")
	}

	update = false
	modifyRuleReq := map[string]interface{}{
		"Id": d.Id(),
	}

	modifyRuleReq["Category"] = d.Get("category")
	modifyRuleReq["Content"] = d.Get("content")

	modifyRuleReq["CustomType"] = 1
	modifyRuleReq["Name"] = d.Get("rule_name")
	if !d.IsNewResource() && d.HasChange("content_category") {
		update = true
	}
	if v, ok := d.GetOk("content_category"); ok {
		modifyRuleReq["ContentCategory"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		modifyRuleReq["Description"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		modifyRuleReq["Lang"] = v
	}
	if !d.IsNewResource() && d.HasChange("product_code") {
		update = true
	}
	if v, ok := d.GetOk("product_code"); ok {
		modifyRuleReq["ProductCode"] = v
	}
	if !d.IsNewResource() && d.HasChange("product_id") {
		update = true
	}
	if v, ok := d.GetOk("product_id"); ok {
		modifyRuleReq["ProductId"] = v
	}
	if !d.IsNewResource() && d.HasChange("risk_level_id") {
		update = true
	}
	if v, ok := d.GetOk("risk_level_id"); ok {
		modifyRuleReq["RiskLevelId"] = v
	}
	if !d.IsNewResource() && d.HasChange("rule_type") {
		update = true
	}
	if v, ok := d.GetOk("rule_type"); ok {
		modifyRuleReq["RuleType"] = v
	}
	if !d.IsNewResource() && d.HasChange("stat_express") {
		update = true
	}
	if v, ok := d.GetOk("stat_express"); ok {
		modifyRuleReq["StatExpress"] = v
	}
	if !d.IsNewResource() && d.HasChange("target") {
		update = true
	}
	if v, ok := d.GetOk("target"); ok {
		modifyRuleReq["Target"] = v
	}
	if !d.IsNewResource() && d.HasChange("warn_level") {
		update = true
	}
	if v, ok := d.GetOk("warn_level"); ok {
		modifyRuleReq["WarnLevel"] = v
	}
	if update {
		action := "ModifyRule"
		conn, err := client.NewSddpClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, modifyRuleReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyRuleReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("category")
		d.SetPartial("content")
		d.SetPartial("custom_type")
		d.SetPartial("rule_name")
		d.SetPartial("content_category")
		d.SetPartial("description")
		d.SetPartial("lang")
		d.SetPartial("product_code")
		d.SetPartial("product_id")
		d.SetPartial("risk_level_id")
		d.SetPartial("rule_type")
		d.SetPartial("stat_express")
		d.SetPartial("target")
		d.SetPartial("warn_level")
	}
	d.Partial(false)
	return resourceAlicloudSddpRuleRead(d, meta)
}
func resourceAlicloudSddpRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRule"
	var response map[string]interface{}
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
