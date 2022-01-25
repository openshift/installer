package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudVpcTrafficMirrorFilterEgressRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcTrafficMirrorFilterEgressRuleCreate,
		Read:   resourceAlicloudVpcTrafficMirrorFilterEgressRuleRead,
		Update: resourceAlicloudVpcTrafficMirrorFilterEgressRuleUpdate,
		Delete: resourceAlicloudVpcTrafficMirrorFilterEgressRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "ICMP" {
						return true
					}
					return false
				},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "ICMP", "TCP", "UDP"}, false),
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},
			"source_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "ICMP" {
						return true
					}
					return false
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_mirror_filter_egress_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcTrafficMirrorFilterEgressRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTrafficMirrorFilterRules"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	requestEgressRules := make(map[string]interface{})
	requestEgressRulesMap := make([]interface{}, 0)
	requestEgressRules["Action"] = d.Get("rule_action")
	requestEgressRules["DestinationCidrBlock"] = d.Get("destination_cidr_block")
	requestEgressRules["Priority"] = d.Get("priority")
	requestEgressRules["Protocol"] = d.Get("protocol")
	requestEgressRules["SourceCidrBlock"] = d.Get("source_cidr_block")
	if fmt.Sprint(d.Get("protocol")) != "ICMP" {
		if v, ok := d.GetOk("source_port_range"); ok {
			requestEgressRules["SourcePortRange"] = v
		}
		if v, ok := d.GetOk("destination_port_range"); ok {
			requestEgressRules["DestinationPortRange"] = v
		}
	}
	requestEgressRulesMap = append(requestEgressRulesMap, requestEgressRules)
	request["EgressRules"] = requestEgressRulesMap

	request["RegionId"] = client.RegionId
	request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateTrafficMirrorFilterRules")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_filter_egress_rule", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.EgressRules", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(request["TrafficMirrorFilterId"], ":", response["InstanceId"]))

	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcTrafficMirrorFilterEgressRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcTrafficMirrorFilterEgressRuleRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorFilterEgressRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcTrafficMirrorFilterEgressRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_filter_egress_rule vpcService.DescribeVpcTrafficMirrorFilterEgressRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("traffic_mirror_filter_id", object["TrafficMirrorFilterId"])
	d.Set("destination_cidr_block", object["DestinationCidrBlock"])
	d.Set("destination_port_range", object["DestinationPortRange"])
	d.Set("priority", object["Priority"])
	d.Set("protocol", object["Protocol"])
	d.Set("rule_action", object["Action"])
	d.Set("source_cidr_block", object["SourceCidrBlock"])
	d.Set("source_port_range", object["SourcePortRange"])
	d.Set("status", object["TrafficMirrorFilterRuleStatus"])
	d.Set("traffic_mirror_filter_egress_rule_id", fmt.Sprint(object["TrafficMirrorFilterRuleId"]))
	return nil
}
func resourceAlicloudVpcTrafficMirrorFilterEgressRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TrafficMirrorFilterRuleId": parts[1],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("destination_cidr_block") {
		update = true
		if v, ok := d.GetOk("destination_cidr_block"); ok {
			request["DestinationCidrBlock"] = v
		}
	}
	if d.HasChange("destination_port_range") {
		update = true
		if v, ok := d.GetOk("destination_port_range"); ok {
			request["DestinationPortRange"] = v
		}
	}
	if d.HasChange("priority") {
		update = true
		if v, ok := d.GetOk("priority"); ok {
			request["Priority"] = v
		}
	}
	if d.HasChange("protocol") {
		update = true
		if v, ok := d.GetOk("protocol"); ok {
			request["Protocol"] = v
		}
	}
	if d.HasChange("source_cidr_block") {
		update = true
		if v, ok := d.GetOk("source_cidr_block"); ok {
			request["SourceCidrBlock"] = v
		}
	}
	if d.HasChange("source_port_range") {
		update = true
		if v, ok := d.GetOk("source_port_range"); ok {
			request["SourcePortRange"] = v
		}
	}
	if d.HasChange("rule_action") {
		update = true
		if v, ok := d.GetOk("rule_action"); ok {
			request["RuleAction"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateTrafficMirrorFilterRuleAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateTrafficMirrorFilterRuleAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcTrafficMirrorFilterEgressRuleStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudVpcTrafficMirrorFilterEgressRuleRead(d, meta)
}
func resourceAlicloudVpcTrafficMirrorFilterEgressRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteTrafficMirrorFilterRules"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TrafficMirrorFilterRuleIds": []string{parts[1]},
		"TrafficMirrorFilterId":      parts[0],
		"RegionId":                   client.RegionId,
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteTrafficMirrorFilterRules")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
