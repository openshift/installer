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

func resourceAlicloudCenTransitRouterRouteTablePropagation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterRouteTablePropagationCreate,
		Read:   resourceAlicloudCenTransitRouterRouteTablePropagationRead,
		Update: resourceAlicloudCenTransitRouterRouteTablePropagationUpdate,
		Delete: resourceAlicloudCenTransitRouterRouteTablePropagationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterRouteTablePropagationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "EnableTransitRouterRouteTablePropagation"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["TransitRouterAttachmentId"] = d.Get("transit_router_attachment_id")
	request["TransitRouterRouteTableId"] = d.Get("transit_router_route_table_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_route_table_propagation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TransitRouterAttachmentId"], ":", request["TransitRouterRouteTableId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterRouteTablePropagationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterRouteTablePropagationRead(d, meta)
}
func resourceAlicloudCenTransitRouterRouteTablePropagationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterRouteTablePropagation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_route_table_propagation cbnService.DescribeCenTransitRouterRouteTablePropagation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("transit_router_attachment_id", parts[0])
	d.Set("transit_router_route_table_id", parts[1])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudCenTransitRouterRouteTablePropagationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudCenTransitRouterRouteTablePropagationRead(d, meta)
}
func resourceAlicloudCenTransitRouterRouteTablePropagationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	cbnService := CbnService{client}
	action := "DisableTransitRouterRouteTablePropagation"
	var response map[string]interface{}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[0],
		"TransitRouterRouteTableId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterRouteTablePropagationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
