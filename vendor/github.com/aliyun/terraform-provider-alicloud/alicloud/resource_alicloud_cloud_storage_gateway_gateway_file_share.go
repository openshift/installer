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

func resourceAlicloudCloudStorageGatewayGatewayFileShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayGatewayFileShareCreate,
		Read:   resourceAlicloudCloudStorageGatewayGatewayFileShareRead,
		Update: resourceAlicloudCloudStorageGatewayGatewayFileShareUpdate,
		Delete: resourceAlicloudCloudStorageGatewayGatewayFileShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_based_enumeration": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "SMB" {
						if _, ok := d.GetOkExists("windows_acl"); ok {
							return false
						}
					}
					return true
				},
			},
			"backend_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1280),
			},
			"browsable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "SMB" {
						return false
					}
					return true
				},
			},
			"cache_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Cache", "Sync"}, false),
			},
			"direct_io": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"download_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1280),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cache_mode"); ok && v.(string) == "Sync" {
						return false
					}
					return true
				},
			},
			"fast_reclaim": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"fe_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1280),
			},
			"gateway_file_share_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][A-Za-z0-9\\_\\.\\-]{1,254}$`), "Length from `1` to `255` characters can contain lowercase letters, digits, (.), (_) Or (-), at the same time, must start with a lowercase letter."),
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ignore_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"in_place": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"index_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lag_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(5, 120),
			},
			"local_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nfs_v4_optimization": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "NFS" {
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
			"partial_sync_paths": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"path_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"polling_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 36000),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NFS", "SMB"}, false),
			},
			"remote_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"remote_sync_download": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cache_mode"); ok && v.(string) == "Sync" {

						return false
					}
					return true
				},
			},
			"ro_client_list": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "NFS" {
						return false
					}
					return true
				},
			},
			"ro_user_list": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "SMB" {
						return false
					}
					return true
				},
			},
			"rw_client_list": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "NFS" {
						return false
					}
					return true
				},
			},
			"rw_user_list": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "SMB" {
						return false
					}
					return true
				},
			},
			"squash": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"all_anonymous", "all_squash", "none", "root_squash"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "NFS" {
						return false
					}
					return true
				},
			},
			"support_archive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"transfer_acceleration": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"windows_acl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("protocol"); ok && v.(string) == "SMB" {
						return false
					}
					return true
				},
			},
			"bypass_cache_read": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayGatewayFileShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	var response map[string]interface{}
	action := "CreateGatewayFileShare"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("access_based_enumeration"); ok {
		request["AccessBasedEnumeration"] = v
	}
	if v, ok := d.GetOk("backend_limit"); ok {
		request["BackendLimit"] = v
	}
	if v, ok := d.GetOkExists("browsable"); ok {
		request["Browsable"] = v
	}
	if v, ok := d.GetOk("cache_mode"); ok {
		request["CacheMode"] = v
	}
	if v, ok := d.GetOkExists("direct_io"); ok {
		request["DirectIO"] = v
	}
	if v, ok := d.GetOk("download_limit"); ok {
		request["DownloadLimit"] = v
	}
	if v, ok := d.GetOkExists("fast_reclaim"); ok {
		request["FastReclaim"] = v
	}
	if v, ok := d.GetOk("fe_limit"); ok {
		request["FrontendLimit"] = v
	}
	if v, ok := d.GetOkExists("bypass_cache_read"); ok {
		request["BypassCacheRead"] = v
	}
	request["Name"] = d.Get("gateway_file_share_name")
	request["GatewayId"] = d.Get("gateway_id")
	if v, ok := d.GetOkExists("ignore_delete"); ok {
		request["IgnoreDelete"] = v
	}
	if v, ok := d.GetOkExists("in_place"); ok {
		request["InPlace"] = v
	}
	if v, ok := d.GetOk("lag_period"); ok {
		request["LagPeriod"] = v
	}
	request["LocalFilePath"] = d.Get("local_path")
	if v, ok := d.GetOkExists("nfs_v4_optimization"); ok {
		request["NfsV4Optimization"] = v
	}
	request["OssBucketName"] = d.Get("oss_bucket_name")
	if v, ok := d.GetOkExists("oss_bucket_ssl"); ok {
		request["OssBucketSsl"] = v
	}
	request["OssEndpoint"] = d.Get("oss_endpoint")
	if v, ok := d.GetOk("partial_sync_paths"); ok {
		request["PartialSyncPaths"] = v
	}
	if v, ok := d.GetOk("path_prefix"); ok {
		request["PathPrefix"] = v
	}
	if v, ok := d.GetOk("polling_interval"); ok {
		request["PollingInterval"] = v
	}
	request["ShareProtocol"] = d.Get("protocol")
	if v, ok := d.GetOkExists("remote_sync"); ok {
		request["RemoteSync"] = v
	}
	if v, ok := d.GetOkExists("remote_sync_download"); ok {
		request["RemoteSyncDownload"] = v
	}
	if v, ok := d.GetOk("ro_client_list"); ok {
		request["ReadOnlyClientList"] = v
	}
	if v, ok := d.GetOk("ro_user_list"); ok {
		request["ReadOnlyUserList"] = v
	}
	if v, ok := d.GetOk("rw_client_list"); ok {
		request["ReadWriteClientList"] = v
	}
	if v, ok := d.GetOk("rw_user_list"); ok {
		request["ReadWriteUserList"] = v
	}
	if v, ok := d.GetOk("squash"); ok {
		request["Squash"] = v
	}
	if v, ok := d.GetOkExists("support_archive"); ok {
		request["SupportArchive"] = v
	}
	if v, ok := d.GetOkExists("transfer_acceleration"); ok {
		request["TransferAcceleration"] = v
	}
	if v, ok := d.GetOkExists("windows_acl"); ok {
		request["WindowsAcl"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_file_share", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	object, err := sgwService.DescribeTasks(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		return nil
	}

	d.SetId(fmt.Sprint(request["GatewayId"], ":", object["RelatedResourceId"]))

	return resourceAlicloudCloudStorageGatewayGatewayFileShareRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayFileShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	object, err := sgwService.DescribeCloudStorageGatewayGatewayFileShare(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_file_share sgwService.DescribeCloudStorageGatewayGatewayFileShare Failed!!! %s", err)
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
	d.Set("access_based_enumeration", object["AccessBasedEnumeration"])
	d.Set("backend_limit", formatInt(object["BeLimit"]))
	d.Set("browsable", object["Browsable"])
	d.Set("cache_mode", object["CacheMode"])
	d.Set("direct_io", object["DirectIO"])
	d.Set("download_limit", formatInt(object["DownloadLimit"]))
	d.Set("fast_reclaim", object["FastReclaim"])
	d.Set("fe_limit", formatInt(object["FeLimit"]))
	d.Set("gateway_file_share_name", object["Name"])
	d.Set("ignore_delete", object["IgnoreDelete"])
	d.Set("in_place", object["InPlace"])
	d.Set("lag_period", formatInt(object["LagPeriod"]))
	d.Set("local_path", object["LocalPath"])
	d.Set("nfs_v4_optimization", object["NfsV4Optimization"])
	d.Set("oss_bucket_name", object["OssBucketName"])
	d.Set("oss_bucket_ssl", object["OssBucketSsl"])
	d.Set("oss_endpoint", object["OssEndpoint"])
	d.Set("partial_sync_paths", object["PartialSyncPaths"])
	d.Set("path_prefix", object["PathPrefix"])
	if v, ok := object["PollingInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("polling_interval", formatInt(v))
	}
	d.Set("protocol", object["Protocol"])
	d.Set("bypass_cache_read", object["BypassCacheRead"])
	d.Set("remote_sync", object["RemoteSync"])
	d.Set("remote_sync_download", object["RemoteSyncDownload"])
	d.Set("ro_client_list", object["RoClientList"])
	d.Set("ro_user_list", object["RoUserList"])
	d.Set("rw_client_list", object["RwClientList"])
	d.Set("rw_user_list", object["RwUserList"])
	d.Set("squash", object["Squash"])
	d.Set("support_archive", object["SupportArchive"])
	d.Set("transfer_acceleration", object["TransferAcceleration"])
	d.Set("windows_acl", object["WindowsAcl"])
	return nil
}
func resourceAlicloudCloudStorageGatewayGatewayFileShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
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
	request["Name"] = d.Get("gateway_file_share_name")
	if d.HasChange("access_based_enumeration") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("access_based_enumeration"); ok {
		request["AccessBasedEnumeration"] = v
	}
	if v, ok := d.GetOkExists("backend_limit"); ok {
		request["BackendLimit"] = v
	}
	if d.HasChange("backend_limit") {
		update = true
	}
	if d.HasChange("browsable") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("browsable"); ok {
		request["Browsable"] = v
	}
	if d.HasChange("bypass_cache_read") {
		update = true
	}
	if v, ok := d.GetOkExists("bypass_cache_read"); ok {
		request["BypassCacheRead"] = v
	}
	if v, ok := d.GetOk("cache_mode"); ok {
		request["CacheMode"] = v
	}
	if d.HasChange("download_limit") {
		update = true
	}
	if v, ok := d.GetOkExists("download_limit"); ok {
		request["DownloadLimit"] = v
	}
	if d.HasChange("fast_reclaim") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("fast_reclaim"); ok {
		request["FastReclaim"] = v
	}
	if d.HasChange("fe_limit") {
		update = true
	}
	if v, ok := d.GetOkExists("fe_limit"); ok {
		request["FrontendLimit"] = v
	}
	if d.HasChange("ignore_delete") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("ignore_delete"); ok {
		request["IgnoreDelete"] = v
	}
	if v, ok := d.GetOkExists("in_place"); ok {
		request["InPlace"] = v
	}
	if d.HasChange("lag_period") {
		update = true
	}
	if v, ok := d.GetOk("lag_period"); ok {
		request["LagPeriod"] = v
	}
	if d.HasChange("nfs_v4_optimization") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("nfs_v4_optimization"); ok {
		request["NfsV4Optimization"] = v
	}
	if d.HasChange("polling_interval") {
		update = true
	}
	if v, ok := d.GetOkExists("polling_interval"); ok {
		request["PollingInterval"] = v
	}
	if d.HasChange("remote_sync") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("remote_sync"); ok {
		request["RemoteSync"] = v
	}
	if d.HasChange("remote_sync_download") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("remote_sync_download"); ok {
		request["RemoteSyncDownload"] = v
	}
	if d.HasChange("ro_client_list") {
		update = true
	}
	if v, ok := d.GetOk("ro_client_list"); ok {
		request["ReadOnlyClientList"] = v
	}
	if d.HasChange("ro_user_list") {
		update = true
	}
	if v, ok := d.GetOk("ro_user_list"); ok {
		request["ReadOnlyUserList"] = v
	}
	if d.HasChange("rw_client_list") {
		update = true
	}
	if v, ok := d.GetOk("rw_client_list"); ok {
		request["ReadWriteClientList"] = v
	}
	if d.HasChange("rw_user_list") {
		update = true
	}
	if v, ok := d.GetOk("rw_user_list"); ok {
		request["ReadWriteUserList"] = v
	}
	if d.HasChange("squash") {
		update = true
	}
	if v, ok := d.GetOk("squash"); ok {
		request["Squash"] = v
	}
	if d.HasChange("transfer_acceleration") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("transfer_acceleration"); ok {
		request["TransferAcceleration"] = v
	}
	if d.HasChange("windows_acl") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("windows_acl"); ok {
		request["WindowsAcl"] = v
	}
	if update {
		action := "UpdateGatewayFileShare"
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
		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudCloudStorageGatewayGatewayFileShareRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayFileShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteGatewayFileShares"
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GatewayId": parts[0],
		"IndexId":   parts[1],
		"Force":     true,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"GatewayDeletionError"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
