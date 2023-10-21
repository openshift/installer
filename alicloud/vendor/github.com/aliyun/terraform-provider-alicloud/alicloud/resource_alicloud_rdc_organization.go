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

func resourceAlicloudRdcOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRdcOrganizationCreate,
		Read:   resourceAlicloudRdcOrganizationRead,
		Update: resourceAlicloudRdcOrganizationUpdate,
		Delete: resourceAlicloudRdcOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"desired_member_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"real_pk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudRdcOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDevopsOrganization"
	request := make(map[string]interface{})
	conn, err := client.NewDevopsrdcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("desired_member_count"); ok {
		request["DesiredMemberCount"] = v
	}
	request["OrgName"] = d.Get("organization_name")
	if v, ok := d.GetOk("real_pk"); ok {
		request["RealPk"] = v
	}
	request["Source"] = d.Get("source")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rdc_organization", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Object"]))

	return resourceAlicloudRdcOrganizationRead(d, meta)
}
func resourceAlicloudRdcOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	devopsRdcService := DevopsRdcService{client}
	object, err := devopsRdcService.DescribeRdcOrganization(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rdc_organization devopsRdcService.DescribeRdcOrganization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("organization_name", object["Name"])
	return nil
}
func resourceAlicloudRdcOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudRdcOrganizationRead(d, meta)
}
func resourceAlicloudRdcOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDevopsOrganization"
	var response map[string]interface{}
	conn, err := client.NewDevopsrdcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"OrgId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidOrganization.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
