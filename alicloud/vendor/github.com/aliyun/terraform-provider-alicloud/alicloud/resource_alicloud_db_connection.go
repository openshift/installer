package alicloud

import (
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const dbConnectionSuffixRegex = "\\.mysql\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const dbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionIdWithSuffixRegexp = regexp.MustCompile(dbConnectionIdWithSuffixRegex)

func resourceAlicloudDBConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBConnectionCreate,
		Read:   resourceAlicloudDBConnectionRead,
		Update: resourceAlicloudDBConnectionUpdate,
		Delete: resourceAlicloudDBConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				Default:      "3306",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}
	action := "AllocateInstancePublicConnection"
	request := map[string]interface{}{
		"RegionId":               client.RegionId,
		"DBInstanceId":           instanceId,
		"ConnectionStringPrefix": prefix,
		"Port":                   d.Get("port"),
		"SourceIp":               client.SourceIp,
	}

	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%v", instanceId, COLON_SEPARATED, request["ConnectionStringPrefix"]))

	if err := rdsService.WaitForDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBConnection(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object["Port"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("ip_address", object["IPAddress"])

	return nil
}

func resourceAlicloudDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") {
		action := "ModifyDBInstanceConnectionString"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": parts[0],
			"SourceIp":     client.SourceIp,
		}
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request["CurrentConnectionString"] = object["ConnectionString"]
		request["ConnectionStringPrefix"] = parts[1]
		request["Port"] = d.Get("port")
		runtime := util.RuntimeOptions{}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	action := "ReleaseInstancePublicConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": split[0],
		"SourceIp":     client.SourceIp,
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request["CurrentConnectionString"] = object["ConnectionString"]
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return rdsService.WaitForDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
