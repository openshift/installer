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

func resourceAlicloudDtsSynchronizationInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDtsSynchronizationInstanceCreate,
		Read:   resourceAlicloudDtsSynchronizationInstanceRead,
		Update: resourceAlicloudDtsSynchronizationInstanceUpdate,
		Delete: resourceAlicloudDtsSynchronizationInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"compute_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sync_architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"bidirectional", "oneway"}, false),
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"destination_endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"xxlarge", "xlarge", "large", "medium", "small"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"payment_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDtsSynchronizationInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDtsInstance"
	request := make(map[string]interface{})
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}

	request["AutoPay"] = false
	request["AutoStart"] = true
	request["InstanceClass"] = "small"
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("compute_unit"); ok {
		request["ComputeUnit"] = v
	}
	if v, ok := d.GetOk("database_count"); ok {
		request["DatabaseCount"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request["DestinationEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request["DestinationRegion"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertDtsSyncPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("quantity"); ok {
		request["Quantity"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request["SourceEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceRegion"] = v
	}
	if v, ok := d.GetOk("sync_architecture"); ok {
		request["SyncArchitecture"] = v
	}
	request["Type"] = "SYNC"
	if v, ok := d.GetOk("payment_duration"); ok {
		request["UsedTime"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_synchronization_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	return resourceAlicloudDtsSynchronizationInstanceRead(d, meta)
}
func resourceAlicloudDtsSynchronizationInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSynchronizationInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_synchronization_instance dtsService.DescribeDtsSynchronizationInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_class", object["DtsJobClass"])
	d.Set("payment_type", convertDtsSyncPaymentTypeResponse(object["PayType"]))
	if v, ok := object["SourceEndpoint"].(map[string]interface{}); ok {
		d.Set("source_endpoint_engine_name", v["EngineName"])
		d.Set("source_endpoint_region", v["Region"])
	}
	if v, ok := object["DestinationEndpoint"].(map[string]interface{}); ok {
		d.Set("destination_endpoint_engine_name", v["EngineName"])
		d.Set("destination_endpoint_region", v["Region"])
	}
	return nil
}
func resourceAlicloudDtsSynchronizationInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudDtsSynchronizationInstanceRead(d, meta)
}
func resourceAlicloudDtsSynchronizationInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v.(string) == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource: alicloud_dts_synchronization_job because it's s. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSynchronizationJob"
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SynchronizationJobId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidJobId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertDtsSyncPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
func convertDtsSyncPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
