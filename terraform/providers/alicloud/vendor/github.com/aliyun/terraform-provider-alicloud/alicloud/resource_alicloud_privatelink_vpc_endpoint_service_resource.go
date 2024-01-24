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

func resourceAlicloudPrivatelinkVpcEndpointServiceResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPrivatelinkVpcEndpointServiceResourceCreate,
		Read:   resourceAlicloudPrivatelinkVpcEndpointServiceResourceRead,
		Update: resourceAlicloudPrivatelinkVpcEndpointServiceResourceUpdate,
		Delete: resourceAlicloudPrivatelinkVpcEndpointServiceResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPrivatelinkVpcEndpointServiceResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachResourceToVpcEndpointService"
	request := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["RegionId"] = client.RegionId
	request["ResourceId"] = d.Get("resource_id")
	request["ResourceType"] = d.Get("resource_type")
	request["ServiceId"] = d.Get("service_id")
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_resource", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ServiceId"], ":", request["ResourceId"]))

	return resourceAlicloudPrivatelinkVpcEndpointServiceResourceRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointServiceResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	object, err := privatelinkService.DescribePrivatelinkVpcEndpointServiceResource(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service_resource privatelinkService.DescribePrivatelinkVpcEndpointServiceResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("resource_id", parts[1])
	d.Set("service_id", parts[0])
	d.Set("resource_type", object["ResourceType"])
	return nil
}
func resourceAlicloudPrivatelinkVpcEndpointServiceResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudPrivatelinkVpcEndpointServiceResourceRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointServiceResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DetachResourceFromVpcEndpointService"
	var response map[string]interface{}
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ResourceId": parts[1],
		"ServiceId":  parts[0],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	request["ResourceType"] = d.Get("resource_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
