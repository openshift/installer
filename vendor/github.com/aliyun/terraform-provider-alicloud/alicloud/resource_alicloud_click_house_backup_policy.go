package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudClickHouseBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudClickHouseBackupPolicyCreate,
		Read:   resourceAlicloudClickHouseBackupPolicyRead,
		Update: resourceAlicloudClickHouseBackupPolicyUpdate,
		Delete: resourceAlicloudClickHouseBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(7, 730),
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"preferred_backup_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudClickHouseBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("backup_retention_period"); ok {
		request["BackupRetentionPeriod"] = v
	}
	request["DBClusterId"] = d.Get("db_cluster_id")
	if v, ok := d.Get("preferred_backup_period").(*schema.Set); ok {
		periodList := expandStringList(v.List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_backup_policy", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DBClusterId"]))
	clickhouseService := ClickhouseService{client}
	stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseBackupPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudClickHouseBackupPolicyRead(d, meta)
}
func resourceAlicloudClickHouseBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	object, err := clickhouseService.DescribeClickHouseBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_backup_policy clickhouseService.DescribeClickHouseBackupPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	if v, ok := object["BackupRetentionPeriod"]; ok && fmt.Sprint(v) != "0" {
		d.Set("backup_retention_period", formatInt(v))
	}
	d.Set("preferred_backup_period", strings.Split(object["PreferredBackupPeriod"].(string), ","))
	d.Set("preferred_backup_time", object["PreferredBackupTime"])
	d.Set("status", fmt.Sprint(object["Switch"]))
	return nil
}
func resourceAlicloudClickHouseBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if d.HasChange("preferred_backup_period") {
		update = true
	}
	if v, ok := d.Get("preferred_backup_period").(*schema.Set); ok {
		periodList := expandStringList(v.List())
		request["PreferredBackupPeriod"] = strings.Join(periodList[:], COMMA_SEPARATED)
	}
	if d.HasChange("preferred_backup_time") {
		update = true
	}
	request["PreferredBackupTime"] = d.Get("preferred_backup_time")
	if d.HasChange("backup_retention_period") {
		update = true
		if v, ok := d.GetOk("backup_retention_period"); ok {
			request["BackupRetentionPeriod"] = v
		}
	}
	if update {
		action := "ModifyBackupPolicy"
		conn, err := client.NewClickhouseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudClickHouseBackupPolicyRead(d, meta)
}
func resourceAlicloudClickHouseBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudClickHouseBackupPolicy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
