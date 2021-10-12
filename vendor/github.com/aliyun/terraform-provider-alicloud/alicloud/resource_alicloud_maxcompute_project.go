package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMaxcomputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMaxcomputeProjectCreate,
		Read:   resourceAlicloudMaxcomputeProjectRead,
		Delete: resourceAlicloudMaxcomputeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"order_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"project_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(3, 27),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(3, 27),
				ConflictsWith: []string{"project_name"},
			},
			"specification_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OdpsStandard"}, false),
			},
		},
	}
}

func resourceAlicloudMaxcomputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateProject"
	request := make(map[string]interface{})
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request["OrderType"] = d.Get("order_type")
	if v, ok := d.GetOk("project_name"); ok {
		request["ProjectName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ProjectName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "project_name" must be set one!`))
	}

	request["OdpsRegionId"] = client.RegionId
	request["OdpsSpecificationType"] = d.Get("specification_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("CreateProject failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(request["ProjectName"]))

	return resourceAlicloudMaxcomputeProjectRead(d, meta)
}
func resourceAlicloudMaxcomputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_maxcompute_project maxcomputeService.DescribeMaxcomputeProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project_name", d.Id())
	d.Set("name", d.Id())
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(object["Data"].(string)), &data); err != nil {
		return WrapError(Error("%v", object))
	}
	d.Set("order_type", data["orderType"].(string))
	d.Set("name", data["projectName"].(string))
	return nil
}
func resourceAlicloudMaxcomputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteProject"
	var response map[string]interface{}
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ProjectName": d.Id(),
	}

	request["RegionIdName"] = client.RegionId
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
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
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"102", "403"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("DeleteProject failed for " + response["Message"].(string)))
	}
	return nil
}
