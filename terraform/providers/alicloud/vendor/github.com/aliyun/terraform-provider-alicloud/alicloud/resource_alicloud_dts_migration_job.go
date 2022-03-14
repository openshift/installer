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

func resourceAlicloudDtsMigrationJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDtsMigrationJobCreate,
		Read:   resourceAlicloudDtsMigrationJobRead,
		Update: resourceAlicloudDtsMigrationJobUpdate,
		Delete: resourceAlicloudDtsMigrationJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"data_synchronization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"db_list": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"destination_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDS", "PolarDB", "POLARDBX20", "ADS", "MONGODB", "GREENPLUM", "DATAHUB", "OTHER", "ECS", "EXPRESS", "CEN", "DG"}, false),
			},

			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADS", "ADB30", "AS400", "DATAHUB", "DB2", "GREENPLUM", "KAFKA", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "PostgreSQL"}, false),
			},
			"destination_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CEN", "DG", "DISTRIBUTED_DMSLOGICDB", "ECS", "EXPRESS", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS"}, false),
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AS400", "DB2", "DMSPOLARDB", "HBASE", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "POSTGRESQL", "TERADATA"}, false),
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"structure_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"checkpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"xxlarge", "xlarge", "large", "medium", "small"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Suspending", "Migrating"}, false),
			},
		},
	}
}

func resourceAlicloudDtsMigrationJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ConfigureDtsJob"
	request := make(map[string]interface{})
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("dts_job_name"); ok {
		request["DtsJobName"] = v
	}
	if v, ok := d.GetOk("checkpoint"); ok {
		request["Checkpoint"] = v
	}

	request["DtsInstanceId"] = d.Get("dts_instance_id")
	request["DataInitialization"] = d.Get("data_initialization")
	request["DataSynchronization"] = d.Get("data_synchronization")
	request["StructureInitialization"] = d.Get("structure_initialization")
	request["DbList"] = d.Get("db_list")

	request["DestinationEndpointInstanceType"] = d.Get("destination_endpoint_instance_type")
	request["DestinationEndpointEngineName"] = d.Get("destination_endpoint_engine_name")
	if v, ok := d.GetOk("destination_endpoint_instance_id"); ok {
		request["DestinationEndpointInstanceID"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request["DestinationEndpointRegion"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_ip"); ok {
		request["DestinationEndpointIP"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_port"); ok {
		request["DestinationEndpointPort"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_database_name"); ok {
		request["DestinationEndpointDataBaseName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
		request["DestinationEndpointUserName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_password"); ok {
		request["DestinationEndpointPassword"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_oracle_sid"); ok {
		request["DestinationEndpointOracleSID"] = v
	}

	request["SourceEndpointInstanceType"] = d.Get("source_endpoint_instance_type")
	request["SourceEndpointEngineName"] = d.Get("source_endpoint_engine_name")
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		request["SourceEndpointInstanceID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceEndpointRegion"] = v
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		request["SourceEndpointIP"] = v
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		request["SourceEndpointPort"] = v
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		request["SourceEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		request["SourceEndpointDatabaseName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		request["SourceEndpointPassword"] = v
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		request["SourceEndpointUserName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		request["SourceEndpointOwnerID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		request["SourceEndpointRole"] = v
	}

	request["PayType"] = convertDtsMigrationJobPaymentTypeRequest(d.Get("payment_type"))
	request["JobType"] = "MIGRATION"
	request["RegionId"] = client.RegionId
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_migration_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DtsJobId"]))
	dtsService := DtsService{client}
	stateConf := BuildStateConf([]string{}, []string{"PreCheckPass", "Migrating", "PrecheckFailed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dtsService.DtsMigrationJobStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDtsMigrationJobUpdate(d, meta)
}
func resourceAlicloudDtsMigrationJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsMigrationJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_migration_job dtsService.DescribeDtsMigrationJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("checkpoint", fmt.Sprint(object["Checkpoint"]))
	if v, ok := object["MigrationMode"].(map[string]interface{}); ok {
		d.Set("data_initialization", v["DataInitialization"])
		d.Set("data_synchronization", v["DataSynchronization"])
		d.Set("structure_initialization", v["StructureInitialization"])

	}
	d.Set("db_list", object["DbObject"])
	if v, ok := object["DestinationEndpoint"].(map[string]interface{}); ok {
		d.Set("destination_endpoint_database_name", v["DatabaseName"])
		d.Set("destination_endpoint_engine_name", v["EngineName"])
		d.Set("destination_endpoint_ip", v["Ip"])
		d.Set("destination_endpoint_instance_id", v["InstanceID"])
		d.Set("destination_endpoint_instance_type", v["InstanceType"])
		d.Set("destination_endpoint_oracle_sid", v["OracleSID"])
		d.Set("destination_endpoint_port", v["Port"])
		d.Set("destination_endpoint_region", v["Region"])
		d.Set("destination_endpoint_user_name", v["UserName"])
	}
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	d.Set("payment_type", convertDtsMigrationJobPaymentTypeResponse(object["PayType"]))
	if v, ok := object["SourceEndpoint"].(map[string]interface{}); ok {
		d.Set("source_endpoint_database_name", v["DatabaseName"])
		d.Set("source_endpoint_engine_name", v["EngineName"])
		d.Set("source_endpoint_ip", v["Ip"])
		d.Set("source_endpoint_instance_id", v["InstanceID"])
		d.Set("source_endpoint_instance_type", v["InstanceType"])
		d.Set("source_endpoint_oracle_sid", v["OracleSID"])
		d.Set("source_endpoint_owner_id", v["AliyunUid"])
		d.Set("source_endpoint_port", v["Port"])
		d.Set("source_endpoint_region", v["Region"])
		d.Set("source_endpoint_role", v["RoleName"])
		d.Set("source_endpoint_user_name", v["UserName"])
	}
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudDtsMigrationJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	var response map[string]interface{}
	d.Partial(false)

	if d.HasChange("status") {
		object, err := dtsService.DescribeDtsMigrationJob(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Migrating" {
				request := map[string]interface{}{
					"DtsJobId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "StartDtsJob"
				conn, err := client.NewDtsClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "Suspending" {
				request := map[string]interface{}{
					"DtsJobId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "SuspendDtsJob"
				conn, err := client.NewDtsClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("status")
		}
	}
	return resourceAlicloudDtsMigrationJobRead(d, meta)
}
func resourceAlicloudDtsMigrationJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDtsJob"
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
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
		if IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertDtsMigrationJobPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
func convertDtsMigrationJobPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
