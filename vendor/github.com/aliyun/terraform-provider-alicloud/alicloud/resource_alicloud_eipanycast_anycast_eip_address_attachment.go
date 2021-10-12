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

func resourceAlicloudEipanycastAnycastEipAddressAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEipanycastAnycastEipAddressAttachmentCreate,
		Read:   resourceAlicloudEipanycastAnycastEipAddressAttachmentRead,
		Delete: resourceAlicloudEipanycastAnycastEipAddressAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"anycast_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SlbInstance"}, false),
			},
			"bind_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEipanycastAnycastEipAddressAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := EipanycastService{client}
	var response map[string]interface{}
	action := "AssociateAnycastEipAddress"
	request := make(map[string]interface{})
	conn, err := client.NewEipanycastClient()
	if err != nil {
		return WrapError(err)
	}
	request["AnycastId"] = d.Get("anycast_id")
	request["BindInstanceId"] = d.Get("bind_instance_id")
	request["BindInstanceRegionId"] = d.Get("bind_instance_region_id")
	request["BindInstanceType"] = d.Get("bind_instance_type")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("AssociateAnycastEipAddress")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eipanycast_anycast_eip_address_attachment", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(request["AnycastId"], ":", request["BindInstanceId"], ":", request["BindInstanceRegionId"], ":", request["BindInstanceType"]))
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEipanycastAnycastEipAddressAttachmentRead(d, meta)
}
func resourceAlicloudEipanycastAnycastEipAddressAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := EipanycastService{client}
	object, err := eipanycastService.DescribeEipanycastAnycastEipAddressAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eipanycast_anycast_eip_address_attachment eipanycastService.DescribeEipanycastAnycastEipAddressAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("anycast_id", parts[0])
	d.Set("bind_instance_id", parts[1])
	d.Set("bind_instance_region_id", parts[2])
	d.Set("bind_instance_type", parts[3])
	if v, ok := object["AnycastEipBindInfoList"].([]interface{})[0].(map[string]interface{}); ok {
		d.Set("bind_time", v["BindTime"])
	}
	return nil
}

func resourceAlicloudEipanycastAnycastEipAddressAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := EipanycastService{client}
	action := "UnassociateAnycastEipAddress"
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	conn, err := client.NewEipanycastClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AnycastId":            parts[0],
		"BindInstanceId":       parts[1],
		"BindInstanceRegionId": parts[2],
		"BindInstanceType":     parts[3],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
