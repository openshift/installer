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

func resourceAlicloudPrivatelinkVpcEndpointConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPrivatelinkVpcEndpointConnectionCreate,
		Read:   resourceAlicloudPrivatelinkVpcEndpointConnectionRead,
		Update: resourceAlicloudPrivatelinkVpcEndpointConnectionUpdate,
		Delete: resourceAlicloudPrivatelinkVpcEndpointConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPrivatelinkVpcEndpointConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	var response map[string]interface{}
	action := "EnableVpcEndpointConnection"
	request := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["EndpointId"] = d.Get("endpoint_id")
	request["RegionId"] = client.RegionId
	request["ServiceId"] = d.Get("service_id")
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointConnectionOperationDenied"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ServiceId"], ":", request["EndpointId"]))
	stateConf := BuildStateConf([]string{}, []string{"Connected"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointConnectionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPrivatelinkVpcEndpointConnectionRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	object, err := privatelinkService.DescribePrivatelinkVpcEndpointConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_connection privatelinkService.DescribePrivatelinkVpcEndpointConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("endpoint_id", parts[1])
	d.Set("service_id", parts[0])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("status", object["ConnectionStatus"])
	return nil
}
func resourceAlicloudPrivatelinkVpcEndpointConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"EndpointId": parts[1],
		"ServiceId":  parts[0],
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateVpcEndpointConnectionAttribute"
		conn, err := client.NewPrivatelinkClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateVpcEndpointConnectionAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointConnectionOperationDenied"}) {
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
	return resourceAlicloudPrivatelinkVpcEndpointConnectionRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	privatelinkService := PrivatelinkService{client}
	action := "DisableVpcEndpointConnection"
	var response map[string]interface{}
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointId": parts[1],
		"ServiceId":  parts[0],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointConnectionOperationDenied"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointConnectionNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Disconnected"}, d.Timeout(schema.TimeoutDelete), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointConnectionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
