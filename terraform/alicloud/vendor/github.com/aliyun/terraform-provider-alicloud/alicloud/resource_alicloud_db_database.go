package alicloud

import (
	"fmt"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBDatabaseCreate,
		Read:   resourceAlicloudDBDatabaseRead,
		Update: resourceAlicloudDBDatabaseUpdate,
		Delete: resourceAlicloudDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_-]*[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and middle lines, and must begin with letters and end with letters or numbers"),
			},

			"character_set": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "CreateDatabase"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"DBInstanceId":     d.Get("instance_id"),
		"DBName":           d.Get("name"),
		"CharacterSetName": d.Get("character_set"),
		"SourceIp":         client.SourceIp,
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request["DBDescription"] = v
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v%s%v", request["DBInstanceId"], COLON_SEPARATED, request["DBName"]))

	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	object, err := rsdService.DescribeDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["DBInstanceId"])
	d.Set("name", object["DBName"])
	d.Set("character_set", object["CharacterSetName"])
	d.Set("description", object["DBDescription"])

	return nil
}

func resourceAlicloudDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("description") && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		action := "ModifyDBDescription"
		request := map[string]interface{}{
			"RegionId":      client.RegionId,
			"DBInstanceId":  parts[0],
			"DBName":        parts[1],
			"DBDescription": d.Get("description"),
			"SourceIp":      client.SourceIp,
		}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
	}
	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDatabase"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": parts[0],
		"DBName":       parts[1],
		"SourceIp":     client.SourceIp,
	}
	// wait instance status is running before deleting database
	if err := rdsService.WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
