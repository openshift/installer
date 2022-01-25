package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRdsBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdsBackupCreate,
		Read:   resourceAlicloudRdsBackupRead,
		Update: resourceAlicloudRdsBackupUpdate,
		Delete: resourceAlicloudRdsBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remove_from_state": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"store_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRdsBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateBackup"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("backup_method"); ok {
		request["BackupMethod"] = v
	}
	if v, ok := d.GetOk("backup_strategy"); ok {
		request["BackupStrategy"] = v
	}
	if v, ok := d.GetOk("backup_type"); ok {
		request["BackupType"] = v
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	if v, ok := d.GetOk("db_name"); ok {
		request["DBName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_backup", action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, rdsService.RdsBackupStateRefreshFunc(d.Get("db_instance_id").(string), response["BackupJobId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// Wait one minute because the query API(DescribeBackups) has not been synchronized when the backup status is Finished
	time.Sleep(1 * time.Minute)
	object, err := rdsService.DescribeBackupTasks(d.Get("db_instance_id").(string), response["BackupJobId"].(string))
	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", object["BackupId"].(string)))
	return resourceAlicloudRdsBackupRead(d, meta)
}

func resourceAlicloudRdsBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeRdsBackup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_backup rdsService.DescribeRdsBackup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_method", object["BackupMethod"])
	d.Set("backup_type", object["BackupType"])
	d.Set("db_instance_id", object["DBInstanceId"])
	d.Set("backup_id", object["BackupId"])
	d.Set("store_status", object["StoreStatus"])
	return nil
}
func resourceAlicloudRdsBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudRdsBackupRead(d, meta)
}
func resourceAlicloudRdsBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackup"
	parts, err := ParseResourceId(d.Id(), 2)
	if d.Get("store_status").(string) == "Disabled" {
		if !d.Get("remove_from_state").(bool) {
			return WrapError(Error("the resource can not be deleted at this time and you can set remove_from_state to true to remove it."))
		} else {
			return nil
		}
	}
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BackupId":     parts[1],
		"DBInstanceId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
