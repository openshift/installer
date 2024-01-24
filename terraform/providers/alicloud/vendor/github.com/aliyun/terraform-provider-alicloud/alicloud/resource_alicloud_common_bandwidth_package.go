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

func resourceAlicloudCommonBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCommonBandwidthPackageCreate,
		Read:   resourceAlicloudCommonBandwidthPackageRead,
		Update: resourceAlicloudCommonBandwidthPackageUpdate,
		Delete: resourceAlicloudCommonBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_package_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.120.0. New field 'bandwidth_package_name' instead.",
				ConflictsWith: []string{"bandwidth_package_name"},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"force": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BGP", "BGP_PRO"}, false),
				Default:      "BGP",
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayBy95", "PayByBandwidth", "PayByTraffic"}, false),
				Default:      "PayByTraffic",
			},
			"ratio": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(10, 100),
				Default:      100,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCommonBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateCommonBandwidthPackage"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("isp"); ok {
		request["ISP"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := d.GetOk("ratio"); ok {
		request["Ratio"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("zone"); ok {
		request["Zone"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateCommonBandwidthPackage")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"BandwidthPackageOperation.conflict", "Throttling"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_common_bandwidth_package", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BandwidthPackageId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.CommonBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCommonBandwidthPackageRead(d, meta)
}
func resourceAlicloudCommonBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeCommonBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_common_bandwidth_package vpcService.DescribeCommonBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("bandwidth_package_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("isp", object["ISP"])
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("ratio", formatInt(object["Ratio"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("deletion_protection", object["DeletionProtection"])

	return nil
}
func resourceAlicloudCommonBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("deletion_protection") {
		var response map[string]interface{}
		action := "DeletionProtection"
		request := map[string]interface{}{
			"RegionId":         client.RegionId,
			"InstanceId":       d.Id(),
			"ProtectionEnable": d.Get("deletion_protection"),
			"Type":             "CBWP",
		}
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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

		d.SetPartial("deletion_protection")
	}

	update := false
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyCommonBandwidthPackageSpec"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
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
		d.SetPartial("bandwidth")
	}
	update = false
	modifyCommonBandwidthPackageAttributeReq := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
	}
	modifyCommonBandwidthPackageAttributeReq["RegionId"] = client.RegionId
	if d.HasChange("bandwidth_package_name") {
		update = true
		modifyCommonBandwidthPackageAttributeReq["Name"] = d.Get("bandwidth_package_name")
	}
	if d.HasChange("name") {
		update = true
		modifyCommonBandwidthPackageAttributeReq["Name"] = d.Get("name")
	}
	if d.HasChange("description") {
		update = true
		modifyCommonBandwidthPackageAttributeReq["Description"] = d.Get("description")
	}
	if update {
		action := "ModifyCommonBandwidthPackageAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, modifyCommonBandwidthPackageAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyCommonBandwidthPackageAttributeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("bandwidth_package_name")
		d.SetPartial("description")
	}
	update = false
	moveResourceGroupReq := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	moveResourceGroupReq["RegionId"] = client.RegionId
	if d.HasChange("resource_group_id") {
		update = true
	}
	moveResourceGroupReq["NewResourceGroupId"] = d.Get("resource_group_id")
	moveResourceGroupReq["ResourceType"] = "bandwidthpackage"
	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, moveResourceGroupReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, moveResourceGroupReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	d.Partial(false)
	return resourceAlicloudCommonBandwidthPackageRead(d, meta)
}
func resourceAlicloudCommonBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteCommonBandwidthPackage"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
	}

	if v, ok := d.GetOk("force"); ok {
		request["Force"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"InvalidBandwidthPackageId.NotFound", "InvalidRegionId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.CommonBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
