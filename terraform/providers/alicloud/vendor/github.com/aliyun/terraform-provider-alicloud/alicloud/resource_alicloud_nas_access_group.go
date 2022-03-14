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

func resourceAlicloudNasAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAccessGroupCreate,
		Read:   resourceAlicloudNasAccessGroupRead,
		Update: resourceAlicloudNasAccessGroupUpdate,
		Delete: resourceAlicloudNasAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"access_group_name"},
			},
			"access_group_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"type"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Classic", "Vpc"}, false),
				ConflictsWith: []string{"access_group_type"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"extreme", "standard"}, false),
				Default:      "standard",
			},
		},
	}
}

func resourceAlicloudNasAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessGroup"
	request := make(map[string]interface{})
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["AccessGroupName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "access_group_name" must be set one!`))
	}

	if v, ok := d.GetOk("access_group_type"); ok {
		request["AccessGroupType"] = v
	} else if v, ok := d.GetOk("type"); ok {
		request["AccessGroupType"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "type" or "access_group_type" must be set one!`))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "ServiceTimeout", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AccessGroupName"], ":", request["FileSystemType"]))

	return resourceAlicloudNasAccessGroupRead(d, meta)
}
func resourceAlicloudNasAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	object, err := nasService.DescribeNasAccessGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_access_group nasService.DescribeNasAccessGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("access_group_name", parts[0])
	d.Set("name", parts[0])
	d.Set("file_system_type", parts[1])
	d.Set("access_group_type", object["AccessGroupType"])
	d.Set("type", object["AccessGroupType"])
	d.Set("description", object["Description"])
	return nil
}
func resourceAlicloudNasAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("description") {
		request := map[string]interface{}{
			"AccessGroupName": parts[0],
			"FileSystemType":  parts[1],
		}
		request["Description"] = d.Get("description")
		action := "ModifyAccessGroup"
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
	return resourceAlicloudNasAccessGroupRead(d, meta)
}
func resourceAlicloudNasAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 2 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), "standard"))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAccessGroup"
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AccessGroupName": parts[0],
		"FileSystemType":  parts[1],
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
