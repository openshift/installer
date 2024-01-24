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

func resourceAlicloudLindormInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLindormInstanceCreate,
		Read:   resourceAlicloudLindormInstanceRead,
		Update: resourceAlicloudLindormInstanceUpdate,
		Delete: resourceAlicloudLindormInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cold_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntInSlice([]int{0}), validation.IntBetween(800, 100000)),
			},
			"core_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"core_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deletion_proection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disk_category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"capacity_cloud_storage", "cloud_efficiency", "cloud_essd", "cloud_ssd"}, false),
			},
			"duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 9),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("payment_type").(string) != "Subscription"
				},
			},
			"file_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"file_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.c.xlarge"}, false),
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_storage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_white_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lts_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"lts_node_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.g.2xlarge", "lindorm.g.xlarge"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"phoenix_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"phoenix_node_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.c.2xlarge", "lindorm.c.4xlarge", "lindorm.c.8xlarge", "lindorm.c.xlarge", "lindorm.g.2xlarge", "lindorm.g.4xlarge", "lindorm.g.8xlarge", "lindorm.g.xlarge"}, false),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("payment_type").(string) != "Subscription"
				},
			},
			"search_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"search_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.g.2xlarge", "lindorm.g.4xlarge", "lindorm.g.8xlarge", "lindorm.g.xlarge"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"table_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.c.2xlarge", "lindorm.c.4xlarge", "lindorm.c.8xlarge", "lindorm.c.xlarge", "lindorm.g.2xlarge", "lindorm.g.4xlarge", "lindorm.g.8xlarge", "lindorm.g.xlarge"}, false),
			},
			"time_series_engine_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"time_serires_engine_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm.g.2xlarge", "lindorm.g.4xlarge", "lindorm.g.8xlarge", "lindorm.g.xlarge"}, false),
			},
			"upgrade_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"open-bds-transfer", "open-bds-transfer-only", "open-lindorm-engine", "open-phoenix-engine", "open-search-engine", "open-tsdb-engine", "upgrade-bds-core-num", "upgrade-bds-transfer", "upgrade-cold-storage", "upgrade-disk-size", "upgrade-file-core-num", "upgrade-file-engine", "upgrade-lindorm-core-num", "upgrade-lindorm-engine", "upgrade-phoenix-core-num", "upgrade-phoenix-engine", "upgrade-search-core-num", "upgrade-search-engine", "upgrade-tsdb-core-num", "upgrade-tsdb-engine"}, false),
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

func resourceAlicloudLindormInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLindormInstance"
	request := make(map[string]interface{})
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("cold_storage"); ok {
		request["ColdStorage"] = v
	}
	if v, ok := d.GetOk("core_spec"); ok {
		request["CoreSpec"] = v
	}
	request["DiskCategory"] = d.Get("disk_category")
	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("file_engine_node_count"); ok {
		request["FilestoreNum"] = v
	}
	if v, ok := d.GetOk("file_engine_specification"); ok {
		request["FilestoreSpec"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceAlias"] = v
	}
	if v, ok := d.GetOk("instance_storage"); ok {
		request["InstanceStorage"] = v
	}
	request["PayType"] = convertLindormInstancePaymentTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("search_engine_node_count"); ok {
		request["SolrNum"] = v
	}
	if v, ok := d.GetOk("search_engine_specification"); ok {
		request["SolrSpec"] = v
	}
	if v, ok := d.GetOk("table_engine_node_count"); ok {
		request["LindormNum"] = v
	}
	if v, ok := d.GetOk("table_engine_specification"); ok {
		request["LindormSpec"] = v
	}
	if v, ok := d.GetOk("time_series_engine_node_count"); ok {
		request["TsdbNum"] = v
	}
	if v, ok := d.GetOk("time_serires_engine_specification"); ok {
		request["TsdbSpec"] = v
	}
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_lindorm_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))
	hitsdbService := HitsdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, hitsdbService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudLindormInstanceRead(d, meta)
}
func resourceAlicloudLindormInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	object, err := hitsdbService.DescribeLindormInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_lindorm_instance hitsdbService.DescribeLindormInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["ColdStorage"]; ok {
		d.Set("cold_storage", formatInt(v))
	}
	d.Set("deletion_proection", object["DeletionProtection"])
	d.Set("disk_category", object["DiskCategory"])
	d.Set("instance_name", object["InstanceAlias"])
	d.Set("instance_storage", object["InstanceStorage"])
	d.Set("payment_type", convertLindormInstancePaymentTypeResponse(object["PayType"]))
	d.Set("status", object["InstanceStatus"])
	d.Set("vswitch_id", object["VswitchId"])
	d.Set("zone_id", object["ZoneId"])
	ipWhite, err := hitsdbService.GetInstanceIpWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ip_white_list", ipWhite)
	getLindormInstanceEngineInfoObject, err := hitsdbService.GetLindormInstanceEngineInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := getLindormInstanceEngineInfoObject["FileEngineNodeCount"]; ok {
		d.Set("file_engine_node_count", formatInt(v))
	}
	d.Set("file_engine_specification", getLindormInstanceEngineInfoObject["FileEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["LtsNodeCount"]; ok {
		d.Set("lts_node_count", formatInt(v))
	}
	d.Set("lts_node_specification", getLindormInstanceEngineInfoObject["LtsNodeSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["PhoenixNodeCount"]; ok {
		d.Set("phoenix_node_count", formatInt(v))
	}
	d.Set("phoenix_node_specification", getLindormInstanceEngineInfoObject["PhoenixNodeSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["SearchEngineNodeCount"]; ok {
		d.Set("search_engine_node_count", formatInt(v))
	}
	d.Set("search_engine_specification", getLindormInstanceEngineInfoObject["SearchEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["TableEngineNodeCount"]; ok {
		d.Set("table_engine_node_count", formatInt(v))
	}
	d.Set("table_engine_specification", getLindormInstanceEngineInfoObject["TableEngineSpecification"])
	if v, ok := getLindormInstanceEngineInfoObject["TimeSeriesNodeCount"]; ok {
		d.Set("time_series_engine_node_count", formatInt(v))
	}
	d.Set("time_serires_engine_specification", getLindormInstanceEngineInfoObject["TimeSeriesSpecification"])
	return nil
}
func resourceAlicloudLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hitsdbService := HitsdbService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("ip_white_list") {
		update = true
		if v, ok := d.GetOk("ip_white_list"); ok && v != nil {
			request["SecurityIpList"] = convertListToCommaSeparate(v.(*schema.Set).List())
		}
	}
	if update {
		if v, ok := d.GetOk("group_name"); ok {
			request["GroupName"] = v
		}
		action := "UpdateInstanceIpWhiteList"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.NotActive", "Lindorm.Errorcode.ParameterInvaild"}) || NeedRetry(err) {
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
		d.SetPartial("ip_white_list")
	}
	update = false
	updateLindormInstanceAttributeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		updateLindormInstanceAttributeReq["InstanceAlias"] = v
	}
	if d.HasChange("deletion_proection") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("deletion_proection"); ok {
			updateLindormInstanceAttributeReq["DeletionProtection"] = v
		}
	}
	if update {
		action := "UpdateLindormInstanceAttribute"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, updateLindormInstanceAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.NotActive", "Lindorm.Errorcode.ParameterInvaild"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLindormInstanceAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_name")
		d.SetPartial("deletion_proection")
	}
	update = false
	upgradeLindormInstanceReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("cold_storage") {
		update = true
	}
	if v, ok := d.GetOk("cold_storage"); ok {
		upgradeLindormInstanceReq["ColdStorage"] = v
	}
	if d.HasChange("instance_storage") {
		update = true
	}
	if v, ok := d.GetOk("instance_storage"); ok {
		upgradeLindormInstanceReq["ClusterStorage"] = v
	}
	upgradeLindormInstanceReq["RegionId"] = client.RegionId
	if d.HasChange("table_engine_node_count") {
		update = true
	}
	if v, ok := d.GetOk("table_engine_node_count"); ok {
		upgradeLindormInstanceReq["LindormNum"] = v
	}

	upgradeLindormInstanceReq["ZoneId"] = d.Get("zone_id")
	if d.HasChange("file_engine_node_count") {
		update = true
		if v, ok := d.GetOk("file_engine_node_count"); ok {
			upgradeLindormInstanceReq["FilestoreSpec"] = v
		}
	}
	if d.HasChange("file_engine_specification") {
		update = true
		if v, ok := d.GetOk("file_engine_specification"); ok {
			upgradeLindormInstanceReq["FilestoreNum"] = v
		}
	}
	if d.HasChange("lts_node_count") {
		update = true
		if v, ok := d.GetOk("lts_node_count"); ok {
			upgradeLindormInstanceReq["LtsCoreNum"] = v
		}
	}
	if d.HasChange("lts_node_specification") {
		update = true
		if v, ok := d.GetOk("lts_node_specification"); ok {
			upgradeLindormInstanceReq["LtsCoreSpec"] = v
		}
	}
	if d.HasChange("phoenix_node_count") {
		update = true
		if v, ok := d.GetOk("phoenix_node_count"); ok {
			upgradeLindormInstanceReq["PhoenixCoreNum"] = v
		}
	}
	if d.HasChange("phoenix_node_specification") {
		update = true
		if v, ok := d.GetOk("phoenix_node_specification"); ok {
			upgradeLindormInstanceReq["PhoenixCoreSpec"] = v
		}
	}
	if d.HasChange("search_engine_node_count") {
		update = true
		if v, ok := d.GetOk("search_engine_node_count"); ok {
			upgradeLindormInstanceReq["SolrNum"] = v
		}
	}
	if d.HasChange("search_engine_specification") {
		update = true
		if v, ok := d.GetOk("search_engine_specification"); ok {
			upgradeLindormInstanceReq["SolrSpec"] = v
		}
	}
	if d.HasChange("table_engine_specification") {
		update = true
		if v, ok := d.GetOk("table_engine_specification"); ok {
			upgradeLindormInstanceReq["LindormSpec"] = v
		}
	}
	if d.HasChange("time_series_engine_node_count") {
		update = true
		if v, ok := d.GetOk("time_series_engine_node_count"); ok {
			upgradeLindormInstanceReq["TsdbNum"] = v
		}
	}
	if d.HasChange("time_serires_engine_specification") {
		update = true
		if v, ok := d.GetOk("time_serires_engine_specification"); ok {
			upgradeLindormInstanceReq["TsdbSpec"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("upgrade_type"); ok {
			upgradeLindormInstanceReq["UpgradeType"] = v
		}
		if v, ok := d.GetOk("core_num"); ok {
			upgradeLindormInstanceReq["CoreNum"] = v
		}
		if v, ok := d.GetOk("core_spec"); ok {
			upgradeLindormInstanceReq["CoreSpec"] = v
		}
		action := "UpgradeLindormInstance"
		conn, err := client.NewHitsdbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 30*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, upgradeLindormInstanceReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.NotActive", "Lindorm.Errorcode.ParameterInvaild"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, upgradeLindormInstanceReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, hitsdbService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("cold_storage")
		d.SetPartial("instance_storage")
		d.SetPartial("table_engine_node_count")
		d.SetPartial("zone_id")
		d.SetPartial("file_engine_node_count")
		d.SetPartial("file_engine_specification")
		d.SetPartial("lts_node_count")
		d.SetPartial("lts_node_specification")
		d.SetPartial("phoenix_node_count")
		d.SetPartial("phoenix_node_specification")
		d.SetPartial("search_engine_node_count")
		d.SetPartial("search_engine_specification")
		d.SetPartial("table_engine_specification")
		d.SetPartial("time_series_engine_node_count")
		d.SetPartial("time_serires_engine_specification")
	}
	d.Partial(false)
	return resourceAlicloudLindormInstanceRead(d, meta)
}
func resourceAlicloudLindormInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseLindormInstance"
	var response map[string]interface{}
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"Lindorm.Errorcode.InstanceNotFound", "Instance.IsDeleted"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertLindormInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}
func convertLindormInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
