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

func resourceAlicloudGaIpSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaIpSetCreate,
		Read:   resourceAlicloudGaIpSetRead,
		Update: resourceAlicloudGaIpSetUpdate,
		Delete: resourceAlicloudGaIpSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerate_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_address_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
				Default:      "IPv4",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaIpSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateIpSets"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccelerateRegion.1.AccelerateRegionId"] = d.Get("accelerate_region_id")
	request["AcceleratorId"] = d.Get("accelerator_id")
	if v, ok := d.GetOk("bandwidth"); ok {
		request["AccelerateRegion.1.Bandwidth"] = v
	}

	if v, ok := d.GetOk("ip_version"); ok {
		request["AccelerateRegion.1.IpVersion"] = v
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateIpSets")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_ip_set", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.IpSets", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(response["IpSetId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaIpSetRead(d, meta)
}
func resourceAlicloudGaIpSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaIpSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_ip_set gaService.DescribeGaIpSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("accelerate_region_id", object["AccelerateRegionId"])
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("ip_address_list", object["IpAddressList"])
	d.Set("ip_version", object["IpVersion"])
	d.Set("status", object["State"])
	return nil
}
func resourceAlicloudGaIpSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"IpSetId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")
	request["RegionId"] = client.RegionId
	if update {
		action := "UpdateIpSet"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateIpSet")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "GreaterThanGa.IpSetBandwidth"}) {
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
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudGaIpSetRead(d, meta)
}
func resourceAlicloudGaIpSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteIpSet"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"IpSetId": d.Id(),
	}

	request["AcceleratorId"] = d.Get("accelerator_id")
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteIpSet")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
