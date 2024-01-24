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

func resourceAlicloudExpressConnectPhysicalConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectPhysicalConnectionCreate,
		Read:   resourceAlicloudExpressConnectPhysicalConnectionRead,
		Update: resourceAlicloudExpressConnectPhysicalConnectionUpdate,
		Delete: resourceAlicloudExpressConnectPhysicalConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"circuit_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line_operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peer_location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"1000Base-LX", "1000Base-T", "100Base-T", "10GBase-LR", "10GBase-T"}, false),
			},
			"redundant_physical_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Canceled", "Enabled", "Terminated"}, false),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudExpressConnectPhysicalConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePhysicalConnection"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccessPointId"] = d.Get("access_point_id")
	if v, ok := d.GetOk("bandwidth"); ok {
		request["bandwidth"] = v
	}
	if v, ok := d.GetOk("circuit_code"); ok {
		request["CircuitCode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["LineOperator"] = d.Get("line_operator")
	request["PeerLocation"] = d.Get("peer_location")
	if v, ok := d.GetOk("physical_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("port_type"); ok {
		request["PortType"] = v
	}
	if v, ok := d.GetOk("redundant_physical_connection_id"); ok {
		request["RedundantPhysicalConnectionId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	request["ClientToken"] = buildClientToken("CreatePhysicalConnection")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_physical_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PhysicalConnectionId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Allocated"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, vpcService.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{"Allocation Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudExpressConnectPhysicalConnectionUpdate(d, meta)
}
func resourceAlicloudExpressConnectPhysicalConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_physical_connection vpcService.DescribeExpressConnectPhysicalConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("access_point_id", object["AccessPointId"])
	d.Set("bandwidth", fmt.Sprint(formatInt(object["Bandwidth"])))
	d.Set("circuit_code", object["CircuitCode"])
	d.Set("description", object["Description"])
	d.Set("line_operator", object["LineOperator"])
	d.Set("peer_location", object["PeerLocation"])
	d.Set("physical_connection_name", object["Name"])
	d.Set("port_type", object["PortType"])
	d.Set("redundant_physical_connection_id", object["RedundantPhysicalConnectionId"])
	d.Set("status", object["Status"])
	d.Set("type", object["Type"])
	return nil
}
func resourceAlicloudExpressConnectPhysicalConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PhysicalConnectionId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		if v, ok := d.GetOk("bandwidth"); ok {
			request["bandwidth"] = v
		}
	}
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
	if !d.IsNewResource() && d.HasChange("line_operator") {
		update = true
		request["LineOperator"] = d.Get("line_operator")
	}
	if !d.IsNewResource() && d.HasChange("peer_location") {
		update = true
		request["PeerLocation"] = d.Get("peer_location")
	}
	if !d.IsNewResource() && d.HasChange("physical_connection_name") {
		update = true
		if v, ok := d.GetOk("physical_connection_name"); ok {
			request["Name"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("port_type") {
		update = true
		if v, ok := d.GetOk("port_type"); ok {
			request["PortType"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("redundant_physical_connection_id") {
		update = true
		if v, ok := d.GetOk("redundant_physical_connection_id"); ok {
			request["RedundantPhysicalConnectionId"] = v
		}
	}
	if update {
		action := "ModifyPhysicalConnectionAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyPhysicalConnectionAttribute")
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
		d.SetPartial("bandwidth")
		d.SetPartial("circuit_code")
		d.SetPartial("description")
		d.SetPartial("line_operator")
		d.SetPartial("peer_location")
		d.SetPartial("physical_connection_name")
		d.SetPartial("port_type")
		d.SetPartial("redundant_physical_connection_id")
	}
	if d.HasChange("status") {
		object, err := vpcService.DescribeExpressConnectPhysicalConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Canceled" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "CancelPhysicalConnection"
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("CancelPhysicalConnection")
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
			}
			if target == "Enabled" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "EnablePhysicalConnection"
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("EnablePhysicalConnection")
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
			}
			if target == "Terminated" {
				request := map[string]interface{}{
					"PhysicalConnectionId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "TerminatePhysicalConnection"
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("TerminatePhysicalConnection")
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
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudExpressConnectPhysicalConnectionRead(d, meta)
}
func resourceAlicloudExpressConnectPhysicalConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePhysicalConnection"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PhysicalConnectionId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeletePhysicalConnection")
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
