package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCddcDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCddcDedicatedHostCreate,
		Read:   resourceAlicloudCddcDedicatedHostRead,
		Update: resourceAlicloudCddcDedicatedHostUpdate,
		Delete: resourceAlicloudCddcDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allocatable", "Suspended"}, false),
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA_Z][a-zA-Z0-9_-]{1,63}`), "The name must be `1` to `64` characters in length and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"image_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"WindowsWithMssqlEntAlwaysonLicense", "WindowsWithMssqlStdLicense", "WindowsWithMssqlEntLicense", "WindowsWithMssqlWebLicense", "AliLinux"}, false),
			},
			"os_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year", "Week"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"used_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCddcDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDedicatedHost"
	request := make(map[string]interface{})
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["DedicatedHostGroupId"] = d.Get("dedicated_host_group_id")
	request["HostClass"] = d.Get("host_class")
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if v, ok := d.GetOk("image_category"); ok {
		request["ImageCategory"] = v
	}
	if v, ok := d.GetOk("os_password"); ok {
		request["OsPassword"] = v
	}
	request["PayType"] = convertCddcDedicatedPaymentTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = fmt.Sprint(v)
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateDedicatedHost")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cddc_dedicated_host", action, AlibabaCloudSdkGoERROR)
	}
	responseDedicateHostList := response["DedicateHostList"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DedicatedHostGroupId"], ":", responseDedicateHostList["DedicateHostList"].([]interface{})[0].(map[string]interface{})["DedicatedHostId"]))
	cddcService := CddcService{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cddcService.CddcDedicatedHostStateRefreshFunc(d.Id(), []string{"5"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCddcDedicatedHostUpdate(d, meta)
}
func resourceAlicloudCddcDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := CddcService{client}
	object, err := cddcService.DescribeCddcDedicatedHost(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host cddcService.DescribeCddcDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("dedicated_host_group_id", object["DedicatedHostGroupId"])
	d.Set("host_class", object["HostClass"])
	d.Set("host_name", object["HostName"])
	d.Set("status", object["HostStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("allocation_status", convertCddcAllocationStatusResponse(object["AllocationStatus"].(string)))
	d.Set("dedicated_host_id", object["DedicatedHostId"])
	listTagResourcesObject, err := cddcService.ListTagResources(d.Id(), "DEDICATEDHOST")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceAlicloudCddcDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := CddcService{client}
	var response map[string]interface{}
	d.Partial(true)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("tags") {
		if err := cddcService.SetResourceTags(d, "DEDICATEDHOST"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("host_name") || d.HasChange("allocation_status") {
		request := map[string]interface{}{
			"DedicatedHostId": parts[1],
		}
		if v, ok := d.GetOk("host_name"); ok {
			request["HostName"] = v
		}
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("allocation_status"); ok {
			request["AllocationStatus"] = convertCddcAllocationStatusRequest(v.(string))
		}
		action := "ModifyDedicatedHostAttribute"
		conn, err := client.NewCddcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cddcService.CddcDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("host_name")
	}

	if !d.IsNewResource() && d.HasChange("host_class") {
		request := map[string]interface{}{
			"DedicatedHostId": parts[1],
		}
		request["TargetClassCode"] = d.Get("host_class")
		request["RegionId"] = client.RegionId
		action := "ModifyDedicatedHostClass"
		conn, err := client.NewCddcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cddcService.CddcDedicatedHostStateRefreshFunc(d.Id(), []string{"5"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("host_class")
	}
	d.Partial(false)
	return resourceAlicloudCddcDedicatedHostRead(d, meta)
}
func resourceAlicloudCddcDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AlicloudResourceCddcDedicatedHost. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertCddcDedicatedPaymentTypeRequest(source string) string {
	switch source {
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertCddcAllocationStatusRequest(source string) string {
	switch source {
	case "Allocatable":
		return "1"
	case "Suspended":
		return "0"
	}
	return source
}
func convertCddcAllocationStatusResponse(source string) string {
	switch source {
	case "1":
		return "Allocatable"
	case "0":
		return "Suspended"
	}
	return source
}
