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

func resourceAlicloudTsdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudTsdbInstanceCreate,
		Read:   resourceAlicloudTsdbInstanceRead,
		Update: resourceAlicloudTsdbInstanceUpdate,
		Delete: resourceAlicloudTsdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(31 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_efficiency", "cloud_essd", "cloud_ssd"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine_type").(string) != "tsdb_influxdb"
				},
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tsdb_influxdb", "tsdb_tsdb"}, false),
			},
			"instance_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"influxdata.n1.16xlarge", "influxdata.n1.16xlarge_ha", "influxdata.n1.2xlarge", "influxdata.n1.2xlarge_ha", "influxdata.n1.4xlarge", "influxdata.n1.4xlarge_ha", "influxdata.n1.8xlarge", "influxdata.n1.8xlarge_ha", "influxdata.n1.mxlarge", "influxdata.n1.mxlarge_ha", "influxdata.n1.xlarge", "influxdata.n1.xlarge_ha", "tsdb.12x.standard", "tsdb.1x.basic", "tsdb.24x.standard", "tsdb.3x.basic", "tsdb.48x.large", "tsdb.4x.basic", "tsdb.96x.large", "tsdb.iot.1x.small"}, false),
			},
			"instance_storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudTsdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	var response map[string]interface{}
	action := "CreateHiTSDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("app_key"); ok {
		request["AppKey"] = v
	}

	if v, ok := d.GetOk("disk_category"); ok {
		request["DiskCategory"] = v
	}

	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}

	if v, ok := d.GetOk("engine_type"); ok {
		request["EngineType"] = v
	}

	if v, ok := d.GetOk("instance_alias"); ok {
		request["InstanceAlias"] = v
	}

	request["InstanceClass"] = d.Get("instance_class")
	request["InstanceStorage"] = d.Get("instance_storage")
	request["PayType"] = convertTsdbInstancePaymentTypeRequest(d.Get("payment_type").(string))
	if request["PayType"].(string) == "PREPAY" {
		request["PricingCycle"] = "Month"
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["VPCId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateHiTSDBInstance")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_tsdb_instance", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["InstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, hitsdbService.TsdbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudTsdbInstanceRead(d, meta)
}
func resourceAlicloudTsdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	object, err := hitsdbService.DescribeTsdbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_tsdb_instance hitsdbService.DescribeTsdbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("disk_category", object["DiskCategory"])
	d.Set("engine_type", object["EngineType"])
	d.Set("instance_alias", object["InstanceAlias"])
	d.Set("instance_class", object["InstanceClass"])
	d.Set("instance_storage", object["InstanceStorage"])
	d.Set("payment_type", convertTsdbInstancePaymentTypeResponse(object["PaymentType"].(string)))
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["VswitchId"])
	d.Set("zone_id", object["ZoneId"])
	return nil
}
func resourceAlicloudTsdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("instance_alias") {
		update = true
	}
	request["InstanceAlias"] = d.Get("instance_alias")
	if update {
		if _, ok := d.GetOk("app_key"); ok {
			request["AppKey"] = d.Get("app_key")
		}
		action := "RenameHiTSDBInstanceAlias"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("instance_alias")
	}
	update = false
	modifyHiTSDBInstanceClassReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("instance_class") {
		update = true
	}
	modifyHiTSDBInstanceClassReq["InstanceClass"] = d.Get("instance_class")
	if d.HasChange("instance_storage") {
		update = true
	}
	modifyHiTSDBInstanceClassReq["InstanceStorage"] = d.Get("instance_storage")
	if update {
		if _, ok := d.GetOk("app_key"); ok {
			modifyHiTSDBInstanceClassReq["AppKey"] = d.Get("app_key")
		}
		action := "ModifyHiTSDBInstanceClass"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 30*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, modifyHiTSDBInstanceClassReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessingError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyHiTSDBInstanceClassReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, hitsdbService.TsdbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_class")
		d.SetPartial("instance_storage")
	}
	d.Partial(false)
	return resourceAlicloudTsdbInstanceRead(d, meta)
}
func resourceAlicloudTsdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHiTSDBInstance"
	var response map[string]interface{}
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("app_key"); ok {
		request["AppKey"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"Instance.IsNotAvailable", "Instance.IsNotPostPay"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertTsdbInstancePaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}
func convertTsdbInstancePaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
