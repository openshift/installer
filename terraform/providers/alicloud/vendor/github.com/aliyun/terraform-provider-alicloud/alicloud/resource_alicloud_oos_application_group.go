package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOosApplicationGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosApplicationGroupCreate,
		Read:   resourceAlicloudOosApplicationGroupRead,
		Delete: resourceAlicloudOosApplicationGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"import_tag_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"import_tag_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudOosApplicationGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateApplicationGroup"
	request := make(map[string]interface{})
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request["Name"] = d.Get("application_group_name")
	request["ApplicationName"] = d.Get("application_name")
	request["DeployRegionId"] = d.Get("deploy_region_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("import_tag_key"); ok {
		if _, ok := d.GetOk("import_tag_value"); !ok {
			return WrapError(fmt.Errorf("The tag key must be passed in at the same time as the tag value (import_tag_value) or none, not just one. "))
		}
		request["ImportTagKey"] = v
	}
	if v, ok := d.GetOk("import_tag_value"); ok {
		if _, ok := d.GetOk("import_tag_key"); !ok {
			return WrapError(fmt.Errorf("The tag value must be passed in at the same time as the tag key (import_tag_key) or none, not just one. "))
		}
		request["ImportTagValue"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateApplicationGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_application_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ApplicationName"], ":", request["Name"]))

	return resourceAlicloudOosApplicationGroupRead(d, meta)
}
func resourceAlicloudOosApplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosApplicationGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_application_group oosService.DescribeOosApplicationGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("application_group_name", parts[1])
	d.Set("application_name", parts[0])
	d.Set("deploy_region_id", object["DeployRegionId"])
	d.Set("description", object["Description"])
	d.Set("import_tag_key", object["ImportTagKey"])
	d.Set("import_tag_value", object["ImportTagValue"])
	return nil
}
func resourceAlicloudOosApplicationGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteApplicationGroup"
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Name":            parts[1],
		"ApplicationName": parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
