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

func resourceAlicloudCassandraBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCassandraBackupPlanCreate,
		Read:   resourceAlicloudCassandraBackupPlanRead,
		Update: resourceAlicloudCassandraBackupPlanUpdate,
		Delete: resourceAlicloudCassandraBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"backup_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_center_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
		},
	}
}

func resourceAlicloudCassandraBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	conn, err := client.NewCdsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("active"); ok {
		request["Active"] = v
	}
	if v, ok := d.GetOk("backup_period"); ok {
		request["BackupPeriod"] = v
	}
	request["BackupTime"] = d.Get("backup_time")
	request["ClusterId"] = d.Get("cluster_id")
	request["DataCenterId"] = d.Get("data_center_id")
	if v, ok := d.GetOk("retention_period"); ok {
		request["RetentionPeriod"] = v
	}
	request["ClientToken"] = buildClientToken("CreateBackupPlan")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cassandra_backup_plan", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ClusterId"], ":", request["DataCenterId"]))

	return resourceAlicloudCassandraBackupPlanRead(d, meta)
}
func resourceAlicloudCassandraBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}
	object, err := cassandraService.DescribeCassandraBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cassandra_backup_plan cassandraService.DescribeCassandraBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	d.Set("data_center_id", parts[1])
	d.Set("active", object["Active"])
	d.Set("backup_period", object["BackupPeriod"])
	d.Set("backup_time", object["BackupTime"])
	d.Set("retention_period", formatInt(object["RetentionPeriod"]))
	return nil
}
func resourceAlicloudCassandraBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ClusterId":    parts[0],
		"DataCenterId": parts[1],
	}
	if d.HasChange("backup_time") {
		update = true
	}
	request["BackupTime"] = d.Get("backup_time")

	if d.HasChange("active") || d.IsNewResource() {
		update = true
		request["Active"] = d.Get("active")
	}
	if d.HasChange("backup_period") {
		update = true
		request["BackupPeriod"] = d.Get("backup_period")
	}
	if d.HasChange("retention_period") {
		update = true
		request["RetentionPeriod"] = d.Get("retention_period")
	}
	if update {
		action := "ModifyBackupPlan"
		conn, err := client.NewCdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCassandraBackupPlanRead(d, meta)
}
func resourceAlicloudCassandraBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteBackupPlan"
	var response map[string]interface{}
	conn, err := client.NewCdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId":    parts[0],
		"DataCenterId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
