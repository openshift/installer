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

func resourceAlicloudDtsSynchronizationJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDtsSynchronizationJobCreate,
		Read:   resourceAlicloudDtsSynchronizationJobRead,
		Update: resourceAlicloudDtsSynchronizationJobUpdate,
		Delete: resourceAlicloudDtsSynchronizationJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"structure_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"synchronization_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Forward", "Reverse"}, false),
			},
			"db_list": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserve": {
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
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AS400", "DB2", "DMSPOLARDB", "HBASE", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "POSTGRESQL", "TERADATA"}, false),
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"destination_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ads", "CEN", "DATAHUB", "DG", "ECS", "EXPRESS", "GREENPLUM", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS"}, false),
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADB20", "ADB30", "AS400", "DATAHUB", "DB2", "GREENPLUM", "KAFKA", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "PostgreSQL"}, false),
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
			"delay_notice": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delay_phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delay_rule_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"error_notice": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"error_phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Synchronizing", "Suspending"}, false),
			},
		},
	}
}

func resourceAlicloudDtsSynchronizationJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ConfigureDtsJob"
	request := make(map[string]interface{})
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("dts_instance_id"); ok {
		request["DtsInstanceId"] = v
	}
	request["DtsJobName"] = d.Get("dts_job_name")
	if v, ok := d.GetOk("checkpoint"); ok {
		request["Checkpoint"] = v
	}
	request["DataInitialization"] = d.Get("data_initialization")
	request["DataSynchronization"] = d.Get("data_synchronization")
	request["StructureInitialization"] = d.Get("structure_initialization")
	request["SynchronizationDirection"] = d.Get("synchronization_direction")
	request["DbList"] = d.Get("db_list")
	if v, ok := d.GetOkExists("delay_notice"); ok {
		request["DelayNotice"] = v
	}
	if v, ok := d.GetOk("delay_phone"); ok {
		request["DelayPhone"] = v
	}
	if v, ok := d.GetOk("delay_rule_time"); ok {
		request["DelayRuleTime"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_database_name"); ok {
		request["DestinationEndpointDataBaseName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request["DestinationEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_ip"); ok {
		request["DestinationEndpointIP"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_instance_id"); ok {
		request["DestinationEndpointInstanceID"] = v
	}
	request["DestinationEndpointInstanceType"] = d.Get("destination_endpoint_instance_type")
	if v, ok := d.GetOk("destination_endpoint_oracle_sid"); ok {
		request["DestinationEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_password"); ok {
		request["DestinationEndpointPassword"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_port"); ok {
		request["DestinationEndpointPort"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request["DestinationEndpointRegion"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
		request["DestinationEndpointUserName"] = v
	}
	if v, ok := d.GetOkExists("error_notice"); ok {
		request["ErrorNotice"] = v
	}
	if v, ok := d.GetOk("error_phone"); ok {
		request["ErrorPhone"] = v
	}
	request["JobType"] = "SYNC"
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("reserve"); ok {
		request["Reserve"] = v
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		request["SourceEndpointDatabaseName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request["SourceEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		request["SourceEndpointIP"] = v
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		request["SourceEndpointInstanceID"] = v
	}
	request["SourceEndpointInstanceType"] = d.Get("source_endpoint_instance_type")
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		request["SourceEndpointOracleSID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		request["SourceEndpointOwnerID"] = v
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		request["SourceEndpointPassword"] = v
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		request["SourceEndpointPort"] = v
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request["SourceEndpointRegion"] = v
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		request["SourceEndpointRole"] = v
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		request["SourceEndpointUserName"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_synchronization_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DtsJobId"]))
	d.Set("dts_instance_id", response["DtsInstanceId"])
	dtsService := DtsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Synchronizing"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{"InitializeFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDtsSynchronizationJobUpdate(d, meta)
}
func resourceAlicloudDtsSynchronizationJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSynchronizationJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_synchronization_job dtsService.DescribeDtsSynchronizationJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	migrationModeObj := object["MigrationMode"].(map[string]interface{})
	destinationEndpointObj := object["DestinationEndpoint"].(map[string]interface{})
	sourceEndpointObj := object["SourceEndpoint"].(map[string]interface{})
	d.Set("checkpoint", fmt.Sprint(formatInt(object["Checkpoint"])))
	d.Set("data_initialization", migrationModeObj["DataInitialization"])
	d.Set("data_synchronization", migrationModeObj["DataSynchronization"])
	d.Set("db_list", object["DbObject"])
	d.Set("destination_endpoint_database_name", destinationEndpointObj["DatabaseName"])
	d.Set("destination_endpoint_engine_name", destinationEndpointObj["EngineName"])
	d.Set("destination_endpoint_ip", destinationEndpointObj["Ip"])
	d.Set("destination_endpoint_instance_id", destinationEndpointObj["InstanceID"])
	d.Set("destination_endpoint_instance_type", destinationEndpointObj["InstanceType"])
	d.Set("destination_endpoint_oracle_sid", destinationEndpointObj["OracleSID"])
	d.Set("destination_endpoint_port", destinationEndpointObj["Port"])
	d.Set("destination_endpoint_region", destinationEndpointObj["Region"])
	d.Set("destination_endpoint_user_name", destinationEndpointObj["UserName"])
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	d.Set("source_endpoint_database_name", sourceEndpointObj["DatabaseName"])
	d.Set("source_endpoint_engine_name", sourceEndpointObj["EngineName"])
	d.Set("source_endpoint_ip", sourceEndpointObj["Ip"])
	d.Set("source_endpoint_instance_id", sourceEndpointObj["InstanceID"])
	d.Set("source_endpoint_instance_type", sourceEndpointObj["InstanceType"])
	d.Set("source_endpoint_oracle_sid", sourceEndpointObj["OracleSID"])
	d.Set("source_endpoint_owner_id", sourceEndpointObj["AliyunUid"])
	d.Set("source_endpoint_port", sourceEndpointObj["Port"])
	d.Set("source_endpoint_region", sourceEndpointObj["Region"])
	d.Set("source_endpoint_role", sourceEndpointObj["RoleName"])
	d.Set("source_endpoint_user_name", sourceEndpointObj["UserName"])
	d.Set("status", object["Status"])
	d.Set("structure_initialization", migrationModeObj["StructureInitialization"])
	d.Set("synchronization_direction", object["SynchronizationDirection"])

	return nil
}
func resourceAlicloudDtsSynchronizationJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("dts_job_name") {
		update = true
		request["DtsJobName"] = d.Get("dts_job_name")
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyDtsJobName"
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("dts_job_name")
	}

	modifyDtsJobPasswordReq := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	modifyDtsJobPasswordReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("source_endpoint_password") {

		modifyDtsJobPasswordReq["Endpoint"] = "src"
		if v, ok := d.GetOk("source_endpoint_password"); ok {
			modifyDtsJobPasswordReq["Password"] = v
		}
		if v, ok := d.GetOk("source_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq["UserName"] = v
		}

		action := "ModifyDtsJobPassword"
		conn, err := client.NewDtsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, modifyDtsJobPasswordReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDtsJobPasswordReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("source_endpoint_password")
		d.SetPartial("source_endpoint_user_name")

		target := d.Get("status").(string)
		err = resourceAlicloudDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	if !d.IsNewResource() && d.HasChange("destination_endpoint_password") {

		modifyDtsJobPasswordReq["Endpoint"] = "src"
		if v, ok := d.GetOk("destination_endpoint_password"); ok {
			modifyDtsJobPasswordReq["Password"] = v
		}
		if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq["UserName"] = v
		}

		action := "ModifyDtsJobPassword"
		conn, err := client.NewDtsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, modifyDtsJobPasswordReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDtsJobPasswordReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("destination_endpoint_password")
		d.SetPartial("destination_endpoint_user_name")

		target := d.Get("status").(string)
		err = resourceAlicloudDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	update = false
	request = map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("instance_class") {
		if v, ok := d.GetOk("instance_class"); ok {
			request["InstanceClass"] = v
		}
		update = true
	}
	request["RegionId"] = client.RegionId
	request["OrderType"] = "UPGRADE"

	if update {
		action := "TransferInstanceClass"
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	if !d.IsNewResource() && d.HasChange("status") {
		target := d.Get("status").(string)
		err := resourceAlicloudDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	d.Partial(false)
	return resourceAlicloudDtsSynchronizationJobRead(d, meta)
}
func resourceAlicloudDtsSynchronizationJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ResetDtsJob"
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}

	if v, ok := d.GetOk("dts_instance_id"); ok {
		request["DtsInstanceId"] = v
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
		if IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func resourceAlicloudDtsSynchronizationJobStatusFlow(d *schema.ResourceData, meta interface{}, target string) error {

	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	var response map[string]interface{}
	object, err := dtsService.DescribeDtsSynchronizationJob(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if object["Status"].(string) != target {
		if target == "Synchronizing" || target == "Suspending" {
			request := map[string]interface{}{
				"DtsJobId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request["SynchronizationDirection"] = v
			}
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
			stateConf := BuildStateConf([]string{}, []string{"Synchronizing"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{"InitializeFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if target == "Suspending" {
			request := map[string]interface{}{
				"DtsJobId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request["SynchronizationDirection"] = v
			}
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
			stateConf := BuildStateConf([]string{}, []string{"Suspending"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("status")
	}

	return nil
}
