package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDbfsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbfsSnapshotCreate,
		Read:   resourceAlicloudDbfsSnapshotRead,
		Update: resourceAlicloudDbfsSnapshotUpdate,
		Delete: resourceAlicloudDbfsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65536),
			},
			"snapshot_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9:_-]{1,127}$`), "The name must be `2` to `128` characters in length and can contain digits, colons (:), underscores (_), and hyphens (-)."), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)|(^auto.*)`), "It must cannot start with `https://`, `https://` and `auto`.")),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDbfsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSnapshot"
	request := make(map[string]interface{})
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["FsId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("retention_days"); ok {
		request["RetentionDays"] = v
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	}
	request["ClientToken"] = buildClientToken("CreateSnapshot")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dbfs_snapshot", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotId"]))
	dbfsService := DbfsService{client}
	stateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbfsService.DbfsSnapshotStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDbfsSnapshotRead(d, meta)
}
func resourceAlicloudDbfsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	object, err := dbfsService.DescribeDbfsSnapshot(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_snapshot dbfsService.DescribeDbfsSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	if v, ok := object["RetentionDays"]; ok && fmt.Sprint(v) != "0" {
		d.Set("retention_days", formatInt(v))
	}
	d.Set("snapshot_name", object["SnapshotName"])
	d.Set("status", object["Status"])
	d.Set("instance_id", object["SourceFsId"])
	return nil
}
func resourceAlicloudDbfsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudDbfsSnapshotRead(d, meta)
}
func resourceAlicloudDbfsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	action := "DeleteSnapshot"
	var response map[string]interface{}
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SnapshotId": d.Id(),
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dbfsService.DbfsSnapshotStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
