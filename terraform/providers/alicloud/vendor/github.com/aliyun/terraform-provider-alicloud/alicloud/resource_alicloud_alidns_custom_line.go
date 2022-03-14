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

func resourceAlicloudAlidnsCustomLine() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsCustomLineCreate,
		Read:   resourceAlicloudAlidnsCustomLineRead,
		Update: resourceAlicloudAlidnsCustomLineUpdate,
		Delete: resourceAlicloudAlidnsCustomLineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"custom_line_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_segment_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudAlidnsCustomLineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddCustomLine"
	request := make(map[string]interface{})
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	for index, ipSegment := range d.Get("ip_segment_list").(*schema.Set).List() {
		ipSegmentArg := ipSegment.(map[string]interface{})
		request[fmt.Sprintf("IpSegment.%d.EndIp", index+1)] = ipSegmentArg["end_ip"]
		request[fmt.Sprintf("IpSegment.%d.StartIp", index+1)] = ipSegmentArg["start_ip"]
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	request["DomainName"] = d.Get("domain_name")
	request["LineName"] = d.Get("custom_line_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_custom_line", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LineId"]))

	return resourceAlicloudAlidnsCustomLineRead(d, meta)
}
func resourceAlicloudAlidnsCustomLineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsCustomLine(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_custom_line alidnsService.DescribeAlidnsCustomLine Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("custom_line_name", object["Name"])
	d.Set("domain_name", object["DomainName"])
	if ipSegmentList, ok := object["IpSegmentList"]; ok {
		ipSegments := make([]map[string]interface{}, 0)
		for _, ipSegmentListItem := range ipSegmentList.([]interface{}) {
			ipSegmentMap := ipSegmentListItem.(map[string]interface{})
			ipSegmentArg := make(map[string]interface{}, 0)
			ipSegmentArg["end_ip"] = ipSegmentMap["EndIp"]
			ipSegmentArg["start_ip"] = ipSegmentMap["StartIp"]
			ipSegments = append(ipSegments, ipSegmentArg)
		}
		d.Set("ip_segment_list", ipSegments)
	}
	return nil
}
func resourceAlicloudAlidnsCustomLineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"LineId": d.Id(),
	}
	if d.HasChange("custom_line_name") {
		update = true
		if v, ok := d.GetOk("custom_line_name"); ok {
			request["LineName"] = v
		}
	}

	if d.HasChange("ip_segment_list") {
		update = true
		for index, ipSegment := range d.Get("ip_segment_list").(*schema.Set).List() {
			ipSegmentArg := ipSegment.(map[string]interface{})
			request[fmt.Sprintf("IpSegment.%d.EndIp", index+1)] = ipSegmentArg["end_ip"]
			request[fmt.Sprintf("IpSegment.%d.StartIp", index+1)] = ipSegmentArg["start_ip"]
		}
	}

	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		action := "UpdateCustomLine"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudAlidnsCustomLineRead(d, meta)
}
func resourceAlicloudAlidnsCustomLineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCustomLines"
	var response map[string]interface{}
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"LineIds": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
