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

func resourceAlicloudExpressConnectVirtualBorderRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectVirtualBorderRouterCreate,
		Read:   resourceAlicloudExpressConnectVirtualBorderRouterRead,
		Update: resourceAlicloudExpressConnectVirtualBorderRouterUpdate,
		Delete: resourceAlicloudExpressConnectVirtualBorderRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associated_physical_connections": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detect_multiplier": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 10),
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"local_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"min_rx_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(200, 1000),
			},
			"min_tx_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(200, 1000),
			},
			"peer_gateway_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_ipv6_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_ipv6_subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peering_subnet_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "deleting", "recovering", "terminated", "terminating", "unconfirmed"}, false),
			},
			"vbr_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_border_router_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 2999),
			},
		},
	}
}

func resourceAlicloudExpressConnectVirtualBorderRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVirtualBorderRouter"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("circuit_code"); ok {
		request["CircuitCode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}
	request["LocalGatewayIp"] = d.Get("local_gateway_ip")
	if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
		request["LocalIpv6GatewayIp"] = v
	}
	request["PeerGatewayIp"] = d.Get("peer_gateway_ip")
	if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
		request["PeerIpv6GatewayIp"] = v
	}
	if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
		request["PeeringIpv6SubnetMask"] = v
	}
	request["PeeringSubnetMask"] = d.Get("peering_subnet_mask")
	request["PhysicalConnectionId"] = d.Get("physical_connection_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vbr_owner_id"); ok {
		request["VbrOwnerId"] = v
	}
	if v, ok := d.GetOk("virtual_border_router_name"); ok {
		request["Name"] = v
	}
	request["VlanId"] = d.Get("vlan_id")
	request["ClientToken"] = buildClientToken("CreateVirtualBorderRouter")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_virtual_border_router", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VbrId"]))

	return resourceAlicloudExpressConnectVirtualBorderRouterUpdate(d, meta)
}
func resourceAlicloudExpressConnectVirtualBorderRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeExpressConnectVirtualBorderRouter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_virtual_border_router vpcService.DescribeExpressConnectVirtualBorderRouter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("description", object["Description"])
	if v, ok := object["DetectMultiplier"]; ok && fmt.Sprint(v) != "0" {
		d.Set("detect_multiplier", formatInt(v))
	}
	d.Set("enable_ipv6", object["EnableIpv6"])
	d.Set("local_gateway_ip", object["LocalGatewayIp"])
	d.Set("local_ipv6_gateway_ip", object["LocalIpv6GatewayIp"])
	if v, ok := object["MinRxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_rx_interval", formatInt(v))
	}
	if v, ok := object["MinTxInterval"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_tx_interval", formatInt(v))
	}
	d.Set("peer_gateway_ip", object["PeerGatewayIp"])
	d.Set("peer_ipv6_gateway_ip", object["PeerIpv6GatewayIp"])
	d.Set("peering_ipv6_subnet_mask", object["PeeringIpv6SubnetMask"])
	d.Set("peering_subnet_mask", object["PeeringSubnetMask"])
	d.Set("physical_connection_id", object["PhysicalConnectionId"])
	d.Set("status", object["Status"])
	d.Set("virtual_border_router_name", object["Name"])
	if v, ok := object["VlanId"]; ok && fmt.Sprint(v) != "0" {
		d.Set("vlan_id", formatInt(v))
	}
	return nil
}
func resourceAlicloudExpressConnectVirtualBorderRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"VbrId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("circuit_code") {
		update = true
		if v, ok := d.GetOk("circuit_code"); ok {
			request["CircuitCode"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("detect_multiplier") {
		update = true
		if v, ok := d.GetOk("detect_multiplier"); ok {
			request["DetectMultiplier"] = v
		} else if v, ok := d.GetOk("min_rx_interval"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_tx_interval"); ok && fmt.Sprint(v) != "" {
				return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "detect_multiplier", "min_rx_interval", d.Get("min_rx_interval"), "min_tx_interval", d.Get("min_tx_interval")))
			}
		}
		request["MinRxInterval"] = d.Get("min_rx_interval")
		request["MinTxInterval"] = d.Get("min_tx_interval")
	}
	if !d.IsNewResource() && d.HasChange("enable_ipv6") {
		update = true
		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request["EnableIpv6"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("local_gateway_ip") {
		update = true
		request["LocalGatewayIp"] = d.Get("local_gateway_ip")
	}
	if !d.IsNewResource() && d.HasChange("local_ipv6_gateway_ip") {
		update = true
		if v, ok := d.GetOk("local_ipv6_gateway_ip"); ok {
			request["LocalIpv6GatewayIp"] = v
		}
	}
	if d.HasChange("min_rx_interval") {
		update = true
		if v, ok := d.GetOk("min_rx_interval"); ok {
			request["MinRxInterval"] = v
		} else if v, ok := d.GetOk("detect_multiplier"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_tx_interval"); ok && fmt.Sprint(v) != "" {
				return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "min_rx_interval", "detect_multiplier", d.Get("detect_multiplier"), "min_tx_interval", d.Get("min_tx_interval")))
			}
		}
		request["DetectMultiplier"] = d.Get("detect_multiplier")
		request["MinTxInterval"] = d.Get("min_tx_interval")
	}
	if d.HasChange("min_tx_interval") {
		update = true
		if v, ok := d.GetOk("min_tx_interval"); ok {
			request["MinTxInterval"] = v
		} else if v, ok := d.GetOk("detect_multiplier"); ok && fmt.Sprint(v) != "" {
			if v, ok := d.GetOk("min_rx_interval"); ok && fmt.Sprint(v) != "" {
				return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "min_tx_interval", "detect_multiplier", d.Get("detect_multiplier"), "min_rx_interval", d.Get("min_rx_interval")))
			}
		}
		request["DetectMultiplier"] = d.Get("detect_multiplier")
		request["MinRxInterval"] = d.Get("min_rx_interval")
	}
	if !d.IsNewResource() && d.HasChange("peer_gateway_ip") {
		update = true
		request["PeerGatewayIp"] = d.Get("peer_gateway_ip")
	}
	if !d.IsNewResource() && d.HasChange("peer_ipv6_gateway_ip") {
		update = true
		if v, ok := d.GetOk("peer_ipv6_gateway_ip"); ok {
			request["PeerIpv6GatewayIp"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("peering_ipv6_subnet_mask") {
		update = true
		if v, ok := d.GetOk("peering_ipv6_subnet_mask"); ok {
			request["PeeringIpv6SubnetMask"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("peering_subnet_mask") {
		update = true
		request["PeeringSubnetMask"] = d.Get("peering_subnet_mask")
	}
	if !d.IsNewResource() && d.HasChange("virtual_border_router_name") {
		update = true
		if v, ok := d.GetOk("virtual_border_router_name"); ok {
			request["Name"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("vlan_id") {
		update = true
		request["VlanId"] = d.Get("vlan_id")
	}
	if update {
		if v, ok := d.GetOk("associated_physical_connections"); ok {
			request["AssociatedPhysicalConnections"] = v
		}
		if v, ok := d.GetOk("bandwidth"); ok {
			request["Bandwidth"] = v
		}
		action := "ModifyVirtualBorderRouterAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyVirtualBorderRouterAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		d.SetPartial("circuit_code")
		d.SetPartial("description")
		d.SetPartial("detect_multiplier")
		d.SetPartial("enable_ipv6")
		d.SetPartial("local_gateway_ip")
		d.SetPartial("local_ipv6_gateway_ip")
		d.SetPartial("min_rx_interval")
		d.SetPartial("min_tx_interval")
		d.SetPartial("peer_gateway_ip")
		d.SetPartial("peer_ipv6_gateway_ip")
		d.SetPartial("peering_ipv6_subnet_mask")
		d.SetPartial("peering_subnet_mask")
		d.SetPartial("virtual_border_router_name")
		d.SetPartial("vlan_id")
	}
	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectVirtualBorderRouter(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "active" {
				request := map[string]interface{}{
					"VbrId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "RecoverVirtualBorderRouter"
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("RecoverVirtualBorderRouter")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "terminated" {
				request := map[string]interface{}{
					"VbrId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "TerminateVirtualBorderRouter"
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("TerminateVirtualBorderRouter")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
				stateConf := BuildStateConf([]string{}, []string{"terminated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.ExpressConnectVirtualBorderRouterStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudExpressConnectVirtualBorderRouterRead(d, meta)
}
func resourceAlicloudExpressConnectVirtualBorderRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVirtualBorderRouter"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"VbrId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteVirtualBorderRouter")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
