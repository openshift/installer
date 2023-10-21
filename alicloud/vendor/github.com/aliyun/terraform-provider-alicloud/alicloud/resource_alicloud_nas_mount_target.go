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

func resourceAlicloudNasMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasMountTargetCreate,
		Read:   resourceAlicloudNasMountTargetRead,
		Update: resourceAlicloudNasMountTargetUpdate,
		Delete: resourceAlicloudNasMountTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudNasMountTargetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	var response map[string]interface{}
	action := "CreateMountTarget"
	request := make(map[string]interface{})
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	}

	request["FileSystemId"] = d.Get("file_system_id")
	request["NetworkType"] = string(Classic)
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["NetworkType"] = string(Vpc)
		request["VpcId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_mount_target", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"], ":", response["MountTargetDomain"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, nasService.NasMountTargetStateRefreshFunc(d.Id(), []string{"Inactive"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNasMountTargetUpdate(d, meta)
}
func resourceAlicloudNasMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	object, err := nasService.DescribeNasMountTarget(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_mount_target nasService.DescribeNasMountTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("file_system_id", parts[0])
	d.Set("access_group_name", object["AccessGroup"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["VswId"])
	return nil
}
func resourceAlicloudNasMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FileSystemId":      parts[0],
		"MountTargetDomain": parts[1],
	}
	if !d.IsNewResource() && d.HasChange("access_group_name") {
		update = true
		request["AccessGroupName"] = d.Get("access_group_name")
	}
	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}
	if update {
		action := "ModifyMountTarget"
		conn, err := client.NewNasClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
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
	}
	return resourceAlicloudNasMountTargetRead(d, meta)
}
func resourceAlicloudNasMountTargetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", strings.Split(d.Id(), "-")[0], d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteMountTarget"
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FileSystemId":      parts[0],
		"MountTargetDomain": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
