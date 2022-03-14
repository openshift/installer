package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudForwardEntryCreate,
		Read:   resourceAlicloudForwardEntryRead,
		Update: resourceAlicloudForwardEntryUpdate,
		Delete: resourceAlicloudForwardEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"external_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forward_entry_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.119.1. New field 'forward_entry_name' instead.",
				ConflictsWith: []string{"forward_entry_name"},
			},
			"forward_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"any", "tcp", "udp"}, false),
			},
			"port_break": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateForwardEntry"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["ExternalIp"] = d.Get("external_ip")
	request["ExternalPort"] = d.Get("external_port")
	if v, ok := d.GetOk("forward_entry_name"); ok {
		request["ForwardEntryName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ForwardEntryName"] = v
	}

	request["ForwardTableId"] = d.Get("forward_table_id")
	request["InternalIp"] = d.Get("internal_ip")
	request["InternalPort"] = d.Get("internal_port")
	request["IpProtocol"] = d.Get("ip_protocol")
	if v, ok := d.GetOkExists("port_break"); ok {
		request["PortBreak"] = v
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidIp.NotInNatgw"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_forward_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ForwardTableId"], ":", response["ForwardEntryId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.ForwardEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudForwardEntryRead(d, meta)
}
func resourceAlicloudForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	object, err := vpcService.DescribeForwardEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_forward_entry vpcService.DescribeForwardEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("forward_entry_id", parts[1])
	d.Set("forward_table_id", parts[0])
	d.Set("external_ip", object["ExternalIp"])
	d.Set("external_port", object["ExternalPort"])
	d.Set("forward_entry_name", object["ForwardEntryName"])
	d.Set("name", object["ForwardEntryName"])
	d.Set("internal_ip", object["InternalIp"])
	d.Set("internal_port", object["InternalPort"])
	d.Set("ip_protocol", object["IpProtocol"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ForwardEntryId": parts[1],
		"ForwardTableId": parts[0],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("external_ip") {
		update = true
		request["ExternalIp"] = d.Get("external_ip")
	}
	if d.HasChange("external_port") {
		update = true
		request["ExternalPort"] = d.Get("external_port")
	}
	if d.HasChange("forward_entry_name") {
		update = true
		request["ForwardEntryName"] = d.Get("forward_entry_name")
	}
	if d.HasChange("name") {
		update = true
		request["ForwardEntryName"] = d.Get("name")
	}
	if d.HasChange("internal_ip") {
		update = true
		request["InternalIp"] = d.Get("internal_ip")
	}
	if d.HasChange("internal_port") {
		update = true
		request["InternalPort"] = d.Get("internal_port")
	}
	if d.HasChange("ip_protocol") {
		update = true
		request["IpProtocol"] = d.Get("ip_protocol")
	}
	if update {
		if _, ok := d.GetOkExists("port_break"); ok {
			request["PortBreak"] = d.Get("port_break")
		}
		action := "ModifyForwardEntry"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.ForwardEntryStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudForwardEntryRead(d, meta)
}
func resourceAlicloudForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	vpcService := VpcService{client}
	action := "DeleteForwardEntry"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ForwardEntryId": parts[1],
		"ForwardTableId": parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"UnknownError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound", "InvalidRegionId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.ForwardEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
