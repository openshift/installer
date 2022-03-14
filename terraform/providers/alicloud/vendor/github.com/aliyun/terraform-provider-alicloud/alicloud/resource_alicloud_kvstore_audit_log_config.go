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

func resourceAlicloudKvstoreAuditLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKvstoreAuditLogConfigCreate,
		Read:   resourceAlicloudKvstoreAuditLogConfigRead,
		Update: resourceAlicloudKvstoreAuditLogConfigUpdate,
		Delete: resourceAlicloudKvstoreAuditLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_audit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention": {
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

func resourceAlicloudKvstoreAuditLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyAuditLogConfig"
	request := make(map[string]interface{})
	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("db_audit"); ok {
		request["DbAudit"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	} else if v, ok := d.GetOkExists("db_audit"); ok && v.(bool) {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "retention", "db_audit", d.Get("db_audit")))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_audit_log_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))
	rKvstoreService := RKvstoreService{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rKvstoreService.KvstoreAuditLogConfigStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudKvstoreAuditLogConfigRead(d, meta)
}
func resourceAlicloudKvstoreAuditLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rKvstoreService := RKvstoreService{client}
	object, err := rKvstoreService.DescribeKvstoreAuditLogConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kvstore_audit_log_config rKvstoreService.DescribeKvstoreAuditLogConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("db_audit", object["DbAudit"])
	if v, ok := object["Retention"]; ok && v.(string) != "" {
		d.Set("retention", formatInt(v))
	}
	describeInstanceAttributeObject, err := rKvstoreService.DescribeInstanceAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("create_time", describeInstanceAttributeObject["CreateTime"])
	d.Set("status", describeInstanceAttributeObject["InstanceStatus"])
	return nil
}
func resourceAlicloudKvstoreAuditLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rKvstoreService := RKvstoreService{client}
	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("db_audit") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("db_audit"); ok {
			request["DbAudit"] = v
		}
	}
	if d.HasChange("retention") {
		update = true
		if v, ok := d.GetOk("retention"); ok {
			request["Retention"] = v
		} else if v, ok := d.GetOkExists("db_audit"); ok && v.(bool) {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "retention", "db_audit", d.Get("db_audit")))
		}
	}
	if update {
		action := "ModifyAuditLogConfig"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rKvstoreService.KvstoreAuditLogConfigStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudKvstoreAuditLogConfigRead(d, meta)
}
func resourceAlicloudKvstoreAuditLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudKvstoreAuditLogConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
