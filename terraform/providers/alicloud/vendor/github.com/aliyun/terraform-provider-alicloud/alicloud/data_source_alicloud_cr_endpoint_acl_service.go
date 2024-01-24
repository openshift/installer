package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCrEndpointAclService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrEndpointAclServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Registry", "Chart"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCrEndpointAclServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "UpdateInstanceEndpointStatus"
	request := make(map[string]interface{})
	request["Enable"] = d.Get("enable")
	request["EndpointType"] = d.Get("endpoint_type")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}

	var response map[string]interface{}
	conn, err := client.NewAcrClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cr_instance_endpoint_acl_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["EndpointType"]))
	d.Set("endpoint_type", request["EndpointType"])
	d.Set("instance_id", request["InstanceId"])
	d.Set("enable", request["Enable"])

	if v, ok := d.GetOk("enable"); !ok || !v.(bool) {
		d.Set("status", "")
		return nil
	}

	crService := CrService{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, crService.CrEndpointAclServiceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.Set("status", "RUNNING")
	return nil
}
