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

func resourceAlicloudHbrRestoreJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrRestoreJobCreate,
		Read:   resourceAlicloudHbrRestoreJobRead,
		Update: resourceAlicloudHbrRestoreJobUpdate,
		Delete: resourceAlicloudHbrRestoreJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"exclude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"restore_job_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"restore_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS_FILE", "NAS", "OSS"}, false),
			},
			"snapshot_hash": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS_FILE", "NAS", "OSS"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_container_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_create_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_data_source_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_file_system_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudHbrRestoreJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRestoreJob"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("exclude"); ok {
		request["Exclude"] = v
	}
	if v, ok := d.GetOk("include"); ok {
		request["Include"] = v
	}
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("restore_job_id"); ok {
		request["RestoreId"] = v
	}
	request["RestoreType"] = d.Get("restore_type")
	if v, ok := d.GetOk("snapshot_hash"); ok {
		request["SnapshotHash"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	request["SourceType"] = d.Get("source_type")
	if v, ok := d.GetOk("target_bucket"); ok {
		request["TargetBucket"] = v
	}
	if v, ok := d.GetOk("target_client_id"); ok {
		request["TargetClientId"] = v
	}
	if v, ok := d.GetOk("target_container"); ok {
		request["TargetContainer"] = v
	}
	if v, ok := d.GetOk("target_container_cluster_id"); ok {
		request["TargetContainerClusterId"] = v
	}
	if v, ok := d.GetOk("target_create_time"); ok {
		request["TargetCreateTime"] = ConvertNasFileSystemStringToUnix(v.(string))
	}
	if v, ok := d.GetOk("target_data_source_id"); ok {
		request["TargetDataSourceId"] = v
	}
	if v, ok := d.GetOk("target_file_system_id"); ok {
		request["TargetFileSystemId"] = v
	}
	if v, ok := d.GetOk("target_instance_id"); ok {
		request["TargetInstanceId"] = v
	}
	if v, ok := d.GetOk("target_path"); ok {
		request["TargetPath"] = v
	}
	if v, ok := d.GetOk("target_prefix"); ok {
		request["TargetPrefix"] = v
	}
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateRestoreJob")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_restore_job", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["RestoreId"], ":", request["RestoreType"]))
	hbrService := HbrService{client}
	stateConf := BuildStateConf([]string{"PARTIAL_COMPLETE", "CREATED"}, []string{"RUNNING", "COMPLETE", "FAILED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrRestoreJobStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudHbrRestoreJobRead(d, meta)
}
func resourceAlicloudHbrRestoreJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrRestoreJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_restore_job hbrService.DescribeHbrRestoreJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("restore_job_id", parts[0])
	d.Set("restore_type", parts[1])
	d.Set("options", object["Options"])
	d.Set("snapshot_hash", object["SnapshotHash"])
	d.Set("snapshot_id", object["SnapshotId"])
	d.Set("source_type", object["SourceType"])
	d.Set("status", object["Status"])
	d.Set("target_bucket", object["TargetBucket"])
	d.Set("target_client_id", object["TargetClientId"])
	if object["TargetCreateTime"] != nil {
		t := int64(formatInt(object["TargetCreateTime"]))
		d.Set("target_create_time", ConvertNasFileSystemUnixToString(d.Get("target_create_time").(string), t))
	}
	d.Set("target_data_source_id", object["TargetDataSourceId"])
	d.Set("target_file_system_id", object["TargetFileSystemId"])
	d.Set("target_instance_id", object["TargetInstanceId"])
	d.Set("target_path", object["TargetPath"])
	d.Set("target_prefix", object["TargetPrefix"])
	d.Set("vault_id", object["VaultId"])
	return nil
}
func resourceAlicloudHbrRestoreJobUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudHbrRestoreJobRead(d, meta)
}
func resourceAlicloudHbrRestoreJobDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudHbrRestoreJob. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
