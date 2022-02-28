package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAdbBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbBackupPolicyCreate,
		Read:   resourceAlicloudAdbBackupPolicyRead,
		Update: resourceAlicloudAdbBackupPolicyUpdate,
		Delete: resourceAlicloudAdbBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"preferred_backup_period": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Required: true,
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Required:     true,
			},
			"backup_retention_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAdbBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("db_cluster_id").(string))

	return resourceAlicloudAdbBackupPolicyUpdate(d, meta)
}

func resourceAlicloudAdbBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	d.Set("backup_retention_period", object.BackupRetentionPeriod)
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("preferred_backup_time", object.PreferredBackupTime)

	return nil
}

func resourceAlicloudAdbBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	if d.HasChange("preferred_backup_period") || d.HasChange("preferred_backup_time") {
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		preferredBackupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		preferredBackupTime := d.Get("preferred_backup_time").(string)

		// wait instance running before modifying
		if err := adbService.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := adbService.ModifyAdbBackupPolicy(d.Id(), preferredBackupTime, preferredBackupPeriod); err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudAdbBackupPolicyRead(d, meta)
}

func resourceAlicloudAdbBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// In case of a delete we are resetting to default values which is Tuesday,Friday each 1am-2am
	client := meta.(*connectivity.AliyunClient)
	request := adb.CreateModifyBackupPolicyRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Id()

	request.PreferredBackupTime = "01:00Z-02:00Z"
	request.PreferredBackupPeriod = "Tuesday,Friday"

	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
