package alicloud

import (
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDdoscooSchedulerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooSchedulerRuleCreate,
		Read:   resourceAlicloudDdoscooSchedulerRuleRead,
		Update: resourceAlicloudDdoscooSchedulerRuleUpdate,
		Delete: resourceAlicloudDdoscooSchedulerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 6}),
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"A", "CNAME"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 6}),
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudDdoscooSchedulerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	request := ddoscoo.CreateCreateSchedulerRuleRequest()
	if v, ok := d.GetOk("param"); ok {
		request.Param = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	request.RuleName = d.Get("rule_name").(string)
	request.RuleType = requests.NewInteger(d.Get("rule_type").(int))
	rules, err := ddoscooService.convertRulesToString(d.Get("rules").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request.Rules = rules
	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.CreateSchedulerRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_scheduler_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ddoscoo.CreateSchedulerRuleResponse)
	d.SetId(fmt.Sprintf("%v", response.RuleName))

	return resourceAlicloudDdoscooSchedulerRuleRead(d, meta)
}
func resourceAlicloudDdoscooSchedulerRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdoscooSchedulerRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_name", d.Id())
	d.Set("cname", object.Cname)
	rule_type, _ := strconv.Atoi(object.RuleType)
	d.Set("rule_type", rule_type)
	rules := make([]map[string]interface{}, len(object.Rules))
	for i, v := range object.Rules {
		rules[i] = map[string]interface{}{
			"priority":   v.Priority,
			"region_id":  v.RegionId,
			"status":     v.Status,
			"type":       v.Type,
			"value":      v.Value,
			"value_type": v.ValueType,
		}
	}
	if err := d.Set("rules", rules); err != nil {
		return WrapError(err)
	}
	return nil
}
func resourceAlicloudDdoscooSchedulerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	update := false
	request := ddoscoo.CreateModifySchedulerRuleRequest()
	request.RuleName = d.Id()
	if d.HasChange("rule_type") {
		update = true
	}
	request.RuleType = requests.NewInteger(d.Get("rule_type").(int))
	if d.HasChange("rules") {
		update = true
	}
	rules, err := ddoscooService.convertRulesToString(d.Get("rules").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request.Rules = rules
	if d.HasChange("param") {
		update = true
		request.Param = d.Get("param").(string)
	}
	if d.HasChange("resource_group_id") {
		update = true
		request.ResourceGroupId = d.Get("resource_group_id").(string)
	}
	if update {
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ModifySchedulerRule(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudDdoscooSchedulerRuleRead(d, meta)
}
func resourceAlicloudDdoscooSchedulerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ddoscoo.CreateDeleteSchedulerRuleRequest()
	request.RuleName = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DeleteSchedulerRule(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
