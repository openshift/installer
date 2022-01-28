package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudStorageGatewayGatewayBlockVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate,
		Read:   resourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead,
		Update: resourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate,
		Delete: resourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cache_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Cache", "WriteThrough"}, false),
			},
			"chap_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"chap_in_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(12, 16),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("chap_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"chap_in_user": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]{0,32}$`), "The chap_in_user must be 1 to 32 characters in length, and can contain letters and digits."),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("chap_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"chunk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{131072, 16384, 32768, 65536, 8192}),
			},
			"gateway_block_volume_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]{0,31}$`), "The name must be 1 to 32 characters in length, and can contain lower case letters and digits."),
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"index_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_source_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"local_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cache_mode"); ok && v.(string) == "Cache" {
						return false
					}
					return true
				},
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"oss_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 262144),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGatewayBlockVolume"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("cache_mode"); ok {
		request["CacheMode"] = v
	}
	if v, ok := d.GetOkExists("chap_enabled"); ok {
		request["ChapEnabled"] = v
	}
	if v, ok := d.GetOk("chap_in_user"); ok {
		request["ChapInUser"] = v
	}
	if v, ok := d.GetOk("chap_in_password"); ok {
		request["ChapInPassword"] = v
	}
	if v, ok := d.GetOk("chunk_size"); ok {
		request["ChunkSize"] = v
	}
	request["Name"] = d.Get("gateway_block_volume_name")
	request["GatewayId"] = d.Get("gateway_id")
	if v, ok := d.GetOk("local_path"); ok {
		request["LocalFilePath"] = v
	}
	request["OssBucketName"] = d.Get("oss_bucket_name")
	if v, ok := d.GetOkExists("oss_bucket_ssl"); ok {
		request["OssBucketSsl"] = v
	}
	request["OssEndpoint"] = d.Get("oss_endpoint")
	request["VolumeProtocol"] = d.Get("protocol")
	if v, ok := d.GetOkExists("recovery"); ok {
		request["Recovery"] = v
	}
	if v, ok := d.GetOk("size"); ok {
		request["Size"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_block_volume", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	sgwService := SgwService{client}
	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response, err = sgwService.DescribeTasks(request["GatewayId"].(string), response["TaskId"].(string))
	if err != nil {
		d.SetId("")
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(request["GatewayId"], ":", response["RelatedResourceId"]))
	return resourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	object, err := sgwService.DescribeCloudStorageGatewayGatewayBlockVolume(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_block_volume sgwService.DescribeCloudStorageGatewayGatewayBlockVolume Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("gateway_id", parts[0])
	d.Set("index_id", parts[1])
	d.Set("cache_mode", object["CacheMode"])
	d.Set("chap_enabled", object["ChapEnabled"])
	d.Set("chap_in_user", object["ChapInUser"])
	if v, ok := object["ChunkSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("chunk_size", formatInt(v))
	}
	d.Set("gateway_block_volume_name", object["Name"])
	d.Set("local_path", object["LocalPath"])
	d.Set("oss_bucket_name", object["OssBucketName"])
	d.Set("oss_bucket_ssl", object["OssBucketSsl"])
	d.Set("oss_endpoint", object["OssEndpoint"])
	d.Set("protocol", object["Protocol"])
	d.Set("status", fmt.Sprint(formatInt(object["VolumeState"])))
	return nil
}
func resourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"GatewayId": parts[0],
		"IndexId":   parts[1],
	}
	if v, ok := d.GetOkExists("chap_enabled"); ok {
		request["ChapEnabled"] = v
	}
	if d.HasChange("chap_enabled") {
		update = true
	}
	if v, ok := d.GetOk("chap_in_password"); ok {
		request["ChapInPassword"] = v
	}
	if d.HasChange("chap_in_password") {
		update = true
	}
	if v, ok := d.GetOk("chap_in_user"); ok {
		request["ChapInUser"] = v
	}
	if d.HasChange("chap_in_user") {
		update = true
	}
	if d.HasChange("size") {
		update = true
		if v, ok := d.GetOk("size"); ok {
			request["Size"] = v
		}
	}
	if update {
		action := "UpdateGatewayBlockVolume"
		conn, err := client.NewHcsSgwClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

		sgwService := SgwService{client}
		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteGatewayBlockVolumes"
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GatewayId": parts[0],
		"IndexId":   parts[1],
	}

	if v, ok := d.GetOkExists("is_source_deletion"); ok {
		request["IsSourceDeletion"] = v
	} else {
		request["IsSourceDeletion"] = true
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	sgwService := SgwService{client}
	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
