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

func resourceAlicloudEcdPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdPolicyGroupCreate,
		Read:   resourceAlicloudEcdPolicyGroupRead,
		Update: resourceAlicloudEcdPolicyGroupUpdate,
		Delete: resourceAlicloudEcdPolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"authorize_access_policy_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"authorize_security_policy_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"clipboard": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "read", "readwrite"}, false),
			},
			"domain_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"html_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"html_file_transfer": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "download", "off", "upload"}, false),
			},
			"local_drive": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"readwrite", "off", "read"}, false),
			},
			"policy_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usb_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"visual_quality": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "lossless", "low", "medium"}, false),
			},
			"watermark": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
			},
			"watermark_transparency": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DARK", "LIGHT", "MIDDLE"}, false),
			},
			"watermark_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"EndUserId", "HostName"}, false),
			},
		},
	}
}

func resourceAlicloudEcdPolicyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicyGroup"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("authorize_access_policy_rules"); ok {
		for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range v.(*schema.Set).List() {
			authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
			request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
			request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
		}
	}
	if v, ok := d.GetOk("authorize_security_policy_rules"); ok {
		for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range v.(*schema.Set).List() {
			authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
			request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
		}
	}
	if v, ok := d.GetOk("clipboard"); ok {
		request["Clipboard"] = v
	}
	if v, ok := d.GetOk("domain_list"); ok {
		request["DomainList"] = v
	}
	if v, ok := d.GetOk("html_access"); ok {
		request["Html5Access"] = v
	}
	if v, ok := d.GetOk("html_file_transfer"); ok {
		request["Html5FileTransfer"] = v
	}
	if v, ok := d.GetOk("local_drive"); ok {
		request["LocalDrive"] = v
	}
	if v, ok := d.GetOk("policy_group_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("usb_redirect"); ok {
		request["UsbRedirect"] = v
	}
	if v, ok := d.GetOk("visual_quality"); ok {
		request["VisualQuality"] = v
	}
	if v, ok := d.GetOk("watermark"); ok {
		request["Watermark"] = v
	}
	if v, ok := d.GetOk("watermark_transparency"); ok {
		request["WatermarkTransparency"] = v
	}
	if v, ok := d.GetOk("watermark_type"); ok {
		request["WatermarkType"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_policy_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyGroupId"]))

	return resourceAlicloudEcdPolicyGroupRead(d, meta)
}
func resourceAlicloudEcdPolicyGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdPolicyGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_policy_group ecdService.DescribeEcdPolicyGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if authorizeAccessPolicyRulesList, ok := object["AuthorizeAccessPolicyRules"]; ok && authorizeAccessPolicyRulesList != nil {
		authorizeAccessPolicyRulesMaps := make([]map[string]interface{}, 0)
		for _, authorizeAccessPolicyRulesListItem := range authorizeAccessPolicyRulesList.([]interface{}) {
			if authorizeAccessPolicyRulesListItemMap, ok := authorizeAccessPolicyRulesListItem.(map[string]interface{}); ok {
				authorizeAccessPolicyRulesMap := make(map[string]interface{})
				authorizeAccessPolicyRulesMap["cidr_ip"] = authorizeAccessPolicyRulesListItemMap["CidrIp"]
				authorizeAccessPolicyRulesMap["description"] = authorizeAccessPolicyRulesListItemMap["Description"]
				authorizeAccessPolicyRulesMaps = append(authorizeAccessPolicyRulesMaps, authorizeAccessPolicyRulesMap)
			}
		}
		d.Set("authorize_access_policy_rules", authorizeAccessPolicyRulesMaps)
	}
	if v, ok := object["AuthorizeSecurityPolicyRules"].([]interface{}); ok {
		authorizeSecurityPolicyRules := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"cidr_ip":     item["CidrIp"],
				"description": item["Description"],
				"ip_protocol": item["IpProtocol"],
				"policy":      item["Policy"],
				"port_range":  item["PortRange"],
				"priority":    item["Priority"],
				"type":        item["Type"],
			}

			authorizeSecurityPolicyRules = append(authorizeSecurityPolicyRules, temp)
		}
		if err := d.Set("authorize_security_policy_rules", authorizeSecurityPolicyRules); err != nil {
			return WrapError(err)
		}
	}
	d.Set("clipboard", object["Clipboard"])
	d.Set("domain_list", object["DomainList"])
	d.Set("html_access", object["Html5Access"])
	d.Set("html_file_transfer", object["Html5FileTransfer"])
	d.Set("local_drive", object["LocalDrive"])
	d.Set("policy_group_name", object["Name"])
	d.Set("status", object["PolicyStatus"])
	d.Set("usb_redirect", object["UsbRedirect"])
	d.Set("visual_quality", object["VisualQuality"])
	d.Set("watermark", object["Watermark"])
	d.Set("watermark_transparency", object["WatermarkTransparency"])
	d.Set("watermark_type", object["WatermarkType"])
	return nil
}
func resourceAlicloudEcdPolicyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PolicyGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("authorize_access_policy_rules") {
		update = true
		err := ecdService.setAuthAccessPolicyRules(d, request, "authorize_access_policy_rules")
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("authorize_access_policy_rules")
	}
	if d.HasChange("authorize_security_policy_rules") {
		update = true
		err := ecdService.setAuthSecurityPolicyRules(d, request, "authorize_security_policy_rules")
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("authorize_security_policy_rules")
	}
	if d.HasChange("cidr_ip") {
		update = true
		if v, ok := d.GetOk("cidr_ip"); ok {
			request["AuthorizeSecurityPolicyRule.*.CidrIp"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["AuthorizeSecurityPolicyRule.*.Description"] = v
		}
	}
	if d.HasChange("ip_protocol") {
		update = true
		if v, ok := d.GetOk("ip_protocol"); ok {
			request["AuthorizeSecurityPolicyRule.*.IpProtocol"] = v
		}
	}
	if d.HasChange("policy") {
		update = true
		if v, ok := d.GetOk("policy"); ok {
			request["AuthorizeSecurityPolicyRule.*.Policy"] = v
		}
	}
	if d.HasChange("port_range") {
		update = true
		if v, ok := d.GetOk("port_range"); ok {
			request["AuthorizeSecurityPolicyRule.*.PortRange"] = v
		}
	}
	if d.HasChange("priority") {
		update = true
		if v, ok := d.GetOk("priority"); ok {
			request["AuthorizeSecurityPolicyRule.*.Priority"] = v
		}
	}
	if d.HasChange("type") {
		update = true
		if v, ok := d.GetOk("type"); ok {
			request["AuthorizeSecurityPolicyRule.*.Type"] = v
		}
	}
	if d.HasChange("clipboard") {
		update = true
		if v, ok := d.GetOk("clipboard"); ok {
			request["Clipboard"] = v
		}
	}
	if d.HasChange("domain_list") {
		update = true
		if v, ok := d.GetOk("domain_list"); ok {
			request["DomainList"] = v
		}
	}
	if d.HasChange("html_access") {
		update = true
		if v, ok := d.GetOk("html_access"); ok {
			request["Html5Access"] = v
		}
	}
	if d.HasChange("html_file_transfer") {
		update = true
		if v, ok := d.GetOk("html_file_transfer"); ok {
			request["Html5FileTransfer"] = v
		}
	}
	if d.HasChange("local_drive") {
		update = true
		if v, ok := d.GetOk("local_drive"); ok {
			request["LocalDrive"] = v
		}
	}
	if d.HasChange("policy_group_name") {
		update = true
		if v, ok := d.GetOk("policy_group_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("usb_redirect") {
		update = true
		if v, ok := d.GetOk("usb_redirect"); ok {
			request["UsbRedirect"] = v
		}
	}
	if d.HasChange("visual_quality") {
		update = true
	}
	if v, ok := d.GetOk("visual_quality"); ok {
		request["VisualQuality"] = v
	}
	if d.HasChange("watermark") {
		update = true
	}
	if v, ok := d.GetOk("watermark"); ok {
		request["Watermark"] = v
	}
	if d.HasChange("watermark_transparency") {
		update = true
	}
	if v, ok := d.GetOk("watermark_transparency"); ok {
		request["WatermarkTransparency"] = v
	}
	if d.HasChange("watermark_type") {
		update = true
	}
	if v, ok := d.GetOk("watermark_type"); ok {
		request["WatermarkType"] = v
	}
	if update {
		action := "ModifyPolicyGroup"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudEcdPolicyGroupRead(d, meta)
}
func resourceAlicloudEcdPolicyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicyGroups"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PolicyGroupId": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
