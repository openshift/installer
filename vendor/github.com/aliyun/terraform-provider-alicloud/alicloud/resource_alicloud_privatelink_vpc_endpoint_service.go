package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPrivatelinkVpcEndpointService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPrivatelinkVpcEndpointServiceCreate,
		Read:   resourceAlicloudPrivatelinkVpcEndpointServiceRead,
		Update: resourceAlicloudPrivatelinkVpcEndpointServiceUpdate,
		Delete: resourceAlicloudPrivatelinkVpcEndpointServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_accept_connection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connect_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payer": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Endpoint",
				ValidateFunc: validation.StringInSlice([]string{"Endpoint", "EndpointService"}, false),
			},
			"service_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPrivatelinkVpcEndpointServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	var response map[string]interface{}
	action := "CreateVpcEndpointService"
	request := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("auto_accept_connection"); ok {
		request["AutoAcceptEnabled"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if v, ok := d.GetOk("payer"); ok {
		request["Payer"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("service_description"); ok {
		request["ServiceDescription"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointServiceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPrivatelinkVpcEndpointServiceUpdate(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	object, err := privatelinkService.DescribePrivatelinkVpcEndpointService(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service privatelinkService.DescribePrivatelinkVpcEndpointService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auto_accept_connection", object["AutoAcceptEnabled"])
	d.Set("connect_bandwidth", formatInt(object["ConnectBandwidth"]))
	d.Set("service_business_status", object["ServiceBusinessStatus"])
	d.Set("service_description", object["ServiceDescription"])
	d.Set("service_domain", object["ServiceDomain"])
	d.Set("status", object["ServiceStatus"])
	return nil
}
func resourceAlicloudPrivatelinkVpcEndpointServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ServiceId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("auto_accept_connection") {
		update = true
		request["AutoAcceptEnabled"] = d.Get("auto_accept_connection")
	}
	if d.HasChange("connect_bandwidth") {
		update = true
		request["ConnectBandwidth"] = d.Get("connect_bandwidth")
	}
	if !d.IsNewResource() && d.HasChange("service_description") {
		update = true
		request["ServiceDescription"] = d.Get("service_description")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateVpcEndpointServiceAttribute"
		conn, err := client.NewPrivatelinkClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointServiceLocked", "EndpointServiceOperationDenied"}) || NeedRetry(err) {
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
	return resourceAlicloudPrivatelinkVpcEndpointServiceRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcEndpointService"
	var response map[string]interface{}
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServiceId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceConnectionDependence"}) || NeedRetry(err) {
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
