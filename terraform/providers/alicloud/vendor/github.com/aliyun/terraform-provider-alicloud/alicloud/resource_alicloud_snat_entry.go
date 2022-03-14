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
)

func resourceAlicloudSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSnatEntryCreate,
		Read:   resourceAlicloudSnatEntryRead,
		Update: resourceAlicloudSnatEntryUpdate,
		Delete: resourceAlicloudSnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"snat_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snat_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snat_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: strings.Fields("source_vswitch_id"),
			},
			"source_vswitch_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: strings.Fields("source_cidr"),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateSnatEntry"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("snat_entry_name"); ok {
		request["SnatEntryName"] = v
	}
	request["SnatIp"] = d.Get("snat_ip")
	request["SnatTableId"] = d.Get("snat_table_id")
	if v, ok := d.GetOk("source_cidr"); ok {
		request["SourceCIDR"] = v
	}
	if v, ok := d.GetOk("source_vswitch_id"); ok {
		request["SourceVSwitchId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateSnatEntry")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"EIP_NOT_IN_GATEWAY", "InternalError", "OperationFailed.Throttling", "OperationUnsupported.EipInBinding", "OperationUnsupported.EipNatBWPCheck", "OperationConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_snat_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["SnatTableId"], ":", response["SnatEntryId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.SnatEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudSnatEntryUpdate(d, meta)
}
func resourceAlicloudSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	object, err := vpcService.DescribeSnatEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_snat_entry vpcService.DescribeSnatEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("snat_entry_id", parts[1])
	d.Set("snat_table_id", parts[0])
	d.Set("snat_entry_name", object["SnatEntryName"])
	d.Set("snat_ip", object["SnatIp"])
	d.Set("source_cidr", object["SourceCIDR"])
	d.Set("source_vswitch_id", object["SourceVSwitchId"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudSnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if !d.IsNewResource() && d.HasChange("snat_entry_name") {
		request := map[string]interface{}{
			"SnatEntryId": parts[1],
			"SnatTableId": parts[0],
		}
		request["RegionId"] = client.RegionId
		request["SnatEntryName"] = d.Get("snat_entry_name")
		action := "ModifySnatEntry"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifySnatEntry")
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.SnatEntryStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudSnatEntryRead(d, meta)
}
func resourceAlicloudSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	vpcService := VpcService{client}
	action := "DeleteSnatEntry"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SnatEntryId": parts[1],
		"SnatTableId": parts[0],
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteSnatEntry")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorretSnatEntryStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnatEntryId.NotFound", "InvalidSnatTableId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.SnatEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
