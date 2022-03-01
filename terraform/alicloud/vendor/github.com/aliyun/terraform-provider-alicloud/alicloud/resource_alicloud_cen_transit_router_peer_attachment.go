package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterPeerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterPeerAttachmentCreate,
		Read:   resourceAlicloudCenTransitRouterPeerAttachmentRead,
		Update: resourceAlicloudCenTransitRouterPeerAttachmentUpdate,
		Delete: resourceAlicloudCenTransitRouterPeerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_publish_route_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cen_bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_transit_router_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "TR",
			},
			"route_table_association_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"route_table_propagation_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_attachment_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"transit_router_attachment_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterPeerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterPeerAttachment"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOk("cen_bandwidth_package_id"); ok {
		request["CenBandwidthPackageId"] = v
	}
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["PeerTransitRouterId"] = d.Get("peer_transit_router_id")
	request["PeerTransitRouterRegionId"] = d.Get("peer_transit_router_region_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}

	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateTransitRouterPeerAttachment")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_peer_attachment", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterPeerAttachmentRead(d, meta)
}
func resourceAlicloudCenTransitRouterPeerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterPeerAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_peer_attachment cbnService.DescribeCenTransitRouterPeerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	d.Set("cen_id", parts[0])
	d.Set("transit_router_attachment_id", parts[1])
	d.Set("auto_publish_route_enabled", object["AutoPublishRouteEnabled"])
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("cen_bandwidth_package_id", object["CenBandwidthPackageId"])
	d.Set("peer_transit_router_id", object["PeerTransitRouterId"])
	d.Set("peer_transit_router_region_id", object["PeerTransitRouterRegionId"])
	d.Set("resource_type", object["ResourceType"])
	d.Set("status", object["Status"])
	d.Set("transit_router_attachment_description", object["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])
	d.Set("transit_router_id", object["TransitRouterId"])
	return nil
}
func resourceAlicloudCenTransitRouterPeerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	update := false
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[1],
	}
	if d.HasChange("auto_publish_route_enabled") || d.IsNewResource() {
		update = true
		request["AutoPublishRouteEnabled"] = d.Get("auto_publish_route_enabled")
	}
	if d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}
	if d.HasChange("cen_bandwidth_package_id") {
		update = true
		request["CenBandwidthPackageId"] = d.Get("cen_bandwidth_package_id")
	}
	if d.HasChange("transit_router_attachment_description") {
		update = true
		request["TransitRouterAttachmentDescription"] = d.Get("transit_router_attachment_description")
	}
	if d.HasChange("transit_router_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_attachment_name")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateTransitRouterPeerAttachmentAttribute"
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCenTransitRouterPeerAttachmentRead(d, meta)
}
func resourceAlicloudCenTransitRouterPeerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterPeerAttachment"
	var response map[string]interface{}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
