package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudWafDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafDomainCreate,
		Read:   resourceAlicloudWafDomainRead,
		Update: resourceAlicloudWafDomainUpdate,
		Delete: resourceAlicloudWafDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PhysicalCluster", "VirtualCluster"}, false),
				Default:      "PhysicalCluster",
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"domain_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"domain"},
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'domain' has been deprecated from version 1.94.0. Use 'domain_name' instead.",
				ConflictsWith: []string{"domain_name"},
			},
			"http2_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"http_to_user_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
				Default:      "Off",
			},
			"https_port": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"https_redirect": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
				Default:      "Off",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_access_product": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Off", "On"}, false),
			},
			"load_balancing": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"IpHash", "RoundRobin"}, false),
				Default:      "IpHash",
			},
			"log_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"read_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"write_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},
		},
	}
}

func resourceAlicloudWafDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	var response map[string]interface{}
	action := "CreateDomain"
	request := make(map[string]interface{})
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		request["ClusterType"] = convertClusterTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("connection_time"); ok {
		request["ConnectionTime"] = v
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request["Domain"] = v
	} else if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "domain" or "domain_name" must be set one!`))
	}

	if v, ok := d.GetOk("http2_port"); ok {
		request["Http2Port"] = convertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("http_port"); ok {
		request["HttpPort"] = convertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("http_to_user_ip"); ok {
		request["HttpToUserIp"] = convertHttpToUserIpRequest(v.(string))
	}

	if v, ok := d.GetOk("https_port"); ok {
		request["HttpsPort"] = convertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("https_redirect"); ok {
		request["HttpsRedirect"] = convertHttpsRedirectRequest(v.(string))
	}

	request["InstanceId"] = d.Get("instance_id")

	request["IsAccessProduct"] = convertIsAccessProductRequest(d.Get("is_access_product").(string))
	if v, ok := d.GetOk("load_balancing"); ok {
		request["LoadBalancing"] = convertLoadBalancingRequest(v.(string))
	}

	if v, ok := d.GetOk("log_headers"); ok {
		logHeaders, err := waf_openapiService.convertLogHeadersToString(v.(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
		request["LogHeaders"] = logHeaders
	}

	if v, ok := d.GetOk("read_time"); ok {
		request["ReadTime"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("source_ips"); ok {
		request["SourceIps"] = convertListToJsonString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("write_time"); ok {
		request["WriteTime"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_domain", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"]))

	return resourceAlicloudWafDomainRead(d, meta)
}
func resourceAlicloudWafDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	object, err := waf_openapiService.DescribeWafDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_domain waf_openapiService.DescribeWafDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("domain_name", parts[1])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("cluster_type", convertClusterTypeResponse(formatInt(object["ClusterType"])))
	d.Set("cname", object["Cname"])
	d.Set("connection_time", formatInt(object["ConnectionTime"]))
	d.Set("http2_port", convertJsonStringToStringList(object["Http2Port"]))
	d.Set("http_port", convertJsonStringToStringList(object["HttpPort"]))
	d.Set("http_to_user_ip", convertHttpToUserIpResponse(formatInt(object["HttpToUserIp"])))
	d.Set("https_port", convertJsonStringToStringList(object["HttpsPort"]))
	d.Set("https_redirect", convertHttpsRedirectResponse(formatInt(object["HttpsRedirect"])))
	d.Set("is_access_product", convertIsAccessProductResponse(formatInt(object["IsAccessProduct"])))
	d.Set("load_balancing", convertLoadBalancingResponse(formatInt(object["LoadBalancing"])))
	if v, ok := object["LogHeaders"].([]interface{}); ok {
		logHeaders := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			logHeaders = append(logHeaders, map[string]interface{}{
				"key":   item["k"].(string),
				"value": item["v"].(string),
			})
		}
		if err := d.Set("log_headers", logHeaders); err != nil {
			return WrapError(err)
		}
	}
	d.Set("read_time", formatInt(object["ReadTime"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("source_ips", object["SourceIps"])
	d.Set("write_time", formatInt(object["WriteTime"]))
	return nil
}
func resourceAlicloudWafDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("cluster_type") {
		request := map[string]interface{}{
			"Domain":     parts[1],
			"InstanceId": parts[0],
		}
		request["ClusterType"] = convertClusterTypeRequest(d.Get("cluster_type").(string))
		action := "ModifyDomainClusterType"
		conn, err := client.NewWafClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cluster_type")
	}
	update := false
	request := map[string]interface{}{
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}
	if d.HasChange("is_access_product") {
		update = true
	}
	request["IsAccessProduct"] = convertIsAccessProductRequest(d.Get("is_access_product").(string))
	if d.HasChange("connection_time") {
		update = true
	}
	request["ConnectionTime"] = d.Get("connection_time")
	if d.HasChange("http2_port") {
		update = true
	}
	request["Http2Port"] = convertListToJsonString(d.Get("http2_port").(*schema.Set).List())
	if d.HasChange("http_port") {
		update = true
	}
	request["HttpPort"] = convertListToJsonString(d.Get("http_port").(*schema.Set).List())
	if d.HasChange("http_to_user_ip") {
		update = true
	}
	request["HttpToUserIp"] = convertHttpToUserIpRequest(d.Get("http_to_user_ip").(string))
	if d.HasChange("https_port") {
		update = true
	}
	request["HttpsPort"] = convertListToJsonString(d.Get("https_port").(*schema.Set).List())
	if d.HasChange("https_redirect") {
		update = true
	}
	request["HttpsRedirect"] = convertHttpsRedirectRequest(d.Get("https_redirect").(string))
	if d.HasChange("load_balancing") {
		update = true
	}
	request["LoadBalancing"] = convertLoadBalancingRequest(d.Get("load_balancing").(string))
	if d.HasChange("log_headers") {
		update = true
	}
	logHeaders, err := waf_openapiService.convertLogHeadersToString(d.Get("log_headers").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request["LogHeaders"] = logHeaders
	if d.HasChange("read_time") {
		update = true
	}
	request["ReadTime"] = d.Get("read_time")
	if d.HasChange("source_ips") {
		update = true
	}
	request["SourceIps"] = convertListToJsonString(d.Get("source_ips").(*schema.Set).List())
	if d.HasChange("write_time") {
		update = true
	}
	request["WriteTime"] = d.Get("write_time")
	if update {
		action := "ModifyDomain"
		conn, err := client.NewWafClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("is_access_product")
		d.SetPartial("connection_time")
		d.SetPartial("http2_port")
		d.SetPartial("http_port")
		d.SetPartial("http_to_user_ip")
		d.SetPartial("https_port")
		d.SetPartial("https_redirect")
		d.SetPartial("load_balancing")
		d.SetPartial("log_headers")
		d.SetPartial("read_time")
		d.SetPartial("source_ips")
		d.SetPartial("write_time")
	}
	d.Partial(false)
	return resourceAlicloudWafDomainRead(d, meta)
}
func resourceAlicloudWafDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDomain"
	var response map[string]interface{}
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertClusterTypeRequest(source string) int {
	switch source {
	case "PhysicalCluster":
		return 0
	case "VirtualCluster":
		return 1
	}
	return 0
}
func convertHttpToUserIpRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertHttpsRedirectRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertIsAccessProductRequest(source string) int {
	switch source {
	case "Off":
		return 0
	case "On":
		return 1
	}
	return 0
}
func convertLoadBalancingRequest(source string) int {
	switch source {
	case "IpHash":
		return 0
	case "RoundRobin":
		return 1
	}
	return 0
}
func convertClusterTypeResponse(source int) string {
	switch source {
	case 0:
		return "PhysicalCluster"
	case 1:
		return "VirtualCluster"
	}
	return ""
}
func convertHttpToUserIpResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertHttpsRedirectResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertIsAccessProductResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}
func convertLoadBalancingResponse(source int) string {
	switch source {
	case 0:
		return "IpHash"
	case 1:
		return "RoundRobin"
	}
	return ""
}
