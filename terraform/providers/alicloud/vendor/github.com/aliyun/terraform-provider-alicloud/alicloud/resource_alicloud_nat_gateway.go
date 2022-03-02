package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNatGatewayCreate,
		Read:   resourceAlicloudNatGatewayRead,
		Update: resourceAlicloudNatGatewayUpdate,
		Delete: resourceAlicloudNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"forward_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByLcu", "PayBySpec"}, false),
			},
			"nat_gateway_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nat_gateway_name"},
			},
			"nat_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enhanced", "Normal"}, false),
				Computed:     true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ConflictsWith: []string{"instance_charge_type"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"payment_type"},
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},
			"bandwidth_packages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"public_ip_addresses": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 4,
				Optional: true,
				Removed:  "Field 'bandwidth_packages' has been removed from provider version 1.121.0.",
			},
			"bandwidth_package_ids": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'bandwidth_package_ids' has been removed from provider version 1.121.0.",
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'spec' has been removed from provider version 1.121.0, replace by 'specification'.",
			},
			"snat_table_ids": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Large", "Middle", "Small", "XLarge.1"}, false),
				Computed:     true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("internet_charge_type").(string) == "PayByLcu"
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("nat_type").(string) != "Enhanced"
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
		},
	}
}

func resourceAlicloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateNatGateway"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := d.GetOk("nat_gateway_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	request["NatType"] = d.Get("nat_type")
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertNatGatewayPaymentTypeRequest(v.(string))
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}

	if v, ok := request["InstanceChargeType"]; ok && v.(string) == "PrePaid" {
		period := d.Get("period").(int)
		request["Duration"] = strconv.Itoa(period)
		request["PricingCycle"] = "Month"
		if period > 9 {
			request["Duration"] = strconv.Itoa(period / 12)
			request["PricingCycle"] = string(Year)
		}
		request["AutoPay"] = true
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("specification"); ok {
		request["Spec"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateNatGateway")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "VswitchStatusError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nat_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["NatGatewayId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNatGatewayUpdate(d, meta)
}
func resourceAlicloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeNatGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nat_gateway vpcService.DescribeNatGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	if v, ok := object["ForwardTableIds"].(map[string]interface{})["ForwardTableId"].([]interface{}); ok {
		ids := []string{}
		for _, id := range v {
			ids = append(ids, id.(string))
		}
		d.Set("forward_table_ids", strings.Join(ids, ","))
	}
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("nat_gateway_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("nat_type", object["NatType"])
	d.Set("payment_type", convertNatGatewayPaymentTypeResponse(object["InstanceChargeType"].(string)))
	d.Set("instance_charge_type", object["InstanceChargeType"])
	d.Set("network_type", object["NetworkType"])
	//if object["InstanceChargeType"] == "PrePaid" {
	//	period, err := computePeriodByUnit(object["CreationTime"], object["ExpiredTime"], d.Get("period").(int), "Month")
	//	if err != nil {
	//		return WrapError(err)
	//	}
	//	d.Set("period", period)
	//}
	if v, ok := object["SnatTableIds"].(map[string]interface{})["SnatTableId"].([]interface{}); ok {
		ids := []string{}
		for _, id := range v {
			ids = append(ids, id.(string))
		}
		d.Set("snat_table_ids", strings.Join(ids, ","))
	}
	d.Set("specification", object["Spec"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["NatGatewayPrivateInfo"].(map[string]interface{})["VswitchId"])
	d.Set("vpc_id", object["VpcId"])

	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "NATGATEWAY")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	d.Set("deletion_protection", object["DeletionProtection"])
	return nil
}
func resourceAlicloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "NATGATEWAY"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("deletion_protection") {
		var response map[string]interface{}
		action := "DeletionProtection"
		request := map[string]interface{}{
			"RegionId":         client.RegionId,
			"InstanceId":       d.Id(),
			"ProtectionEnable": d.Get("deletion_protection"),
			"Type":             "NATGW",
		}
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("deletion_protection")
	}
	update := false
	request := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("nat_gateway_name") {
		update = true
		request["Name"] = d.Get("nat_gateway_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if update {
		action := "ModifyNatGatewayAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("nat_gateway_name")
	}
	if !d.IsNewResource() && d.HasChange("specification") {
		request := map[string]interface{}{
			"NatGatewayId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["Spec"] = d.Get("specification")
		action := "ModifyNatGatewaySpec"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("specification")
	}
	update = false
	updateNatGatewayNatTypeReq := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("nat_type") {
		update = true
	}
	updateNatGatewayNatTypeReq["NatType"] = d.Get("nat_type")
	updateNatGatewayNatTypeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
	}
	updateNatGatewayNatTypeReq["VSwitchId"] = d.Get("vswitch_id")
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			updateNatGatewayNatTypeReq["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateNatGatewayNatType"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		updateNatGatewayNatTypeReq["ClientToken"] = buildClientToken("UpdateNatGatewayNatType")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, updateNatGatewayNatTypeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationFailed.NatGwRouteInMiddleStatus", "TaskConflict", "Throttling", "UnknownError"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, updateNatGatewayNatTypeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("nat_type")
		d.SetPartial("vswitch_id")
		d.SetPartial("dry_run")
	}
	d.Partial(false)
	return resourceAlicloudNatGatewayRead(d, meta)
}
func resourceAlicloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteNatGateway"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"NatGatewayId": d.Id(),
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.BandwidthPackages"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXISTS", "IncorrectStatus.NatGateway", "InvalidNatGatewayId.NotFound", "InvalidRegionId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.NatGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertNatGatewayPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func convertNatGatewayPaymentTypeResponse(source string) string {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
