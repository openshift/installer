package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/denverdino/aliyungo/common"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEipAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEipAddressCreate,
		Read:   resourceAlicloudEipAddressRead,
		Update: resourceAlicloudEipAddressUpdate,
		Delete: resourceAlicloudEipAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(9 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"activity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"address_name"},
				Deprecated:    "Field 'name' has been deprecated from provider version 1.126.0 and it will be remove in the future version. Please use the new attribute 'address_name' instead.",
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic", "PayByDominantTraffic"}, false),
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BGP", "BGP_PRO"}, false),
			},
			"netmode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"public"}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc:     validation.Any(validation.IntBetween(1, 9), validation.IntInSlice([]int{12, 24, 36})),
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
				ValidateFunc:  validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"payment_type"},
				Deprecated:    "Field 'instance_charge_type' has been deprecated from provider version 1.126.0 and it will be remove in the future version. Please use the new attribute 'payment_type' instead.",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudEipAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "AllocateEipAddress"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("activity_id"); ok {
		request["ActivityId"] = v
	}
	if v, ok := d.GetOk("address_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("isp"); ok {
		request["ISP"] = v
	}
	if v, ok := d.GetOk("netmode"); ok {
		request["Netmode"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertEipAddressPaymentTypeRequest(v)
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := request["InstanceChargeType"]; ok && v.(string) == "PrePaid" {
		period := 1
		if v, ok := d.GetOk("period"); ok {
			period = v.(int)
		}
		request["Period"] = period
		request["PricingCycle"] = string(Month)
		if period > 9 {
			request["Period"] = period / 12
			request["PricingCycle"] = string(Year)
		}
		autoPay := true
		if v, ok := d.GetOkExists("auto_pay"); ok {
			autoPay = v.(bool)
		}
		request["AutoPay"] = autoPay
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ClientToken"] = buildClientToken("AllocateEipAddress")
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
		if IsExpectedErrors(err, []string{"COMMODITY.INVALID_COMPONENT"}) && request["InternetChargeType"] == string(PayByBandwidth) {
			return WrapErrorf(err, "Your account is international and it can only create '%s' elastic IP. Please change it and try again. %s", PayByTraffic, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AllocationId"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.EipAddressStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEipAddressUpdate(d, meta)
}
func resourceAlicloudEipAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeEipAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eip_address vpcService.DescribeEipAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("address_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("deletion_protection", object["DeletionProtection"])
	d.Set("description", object["Descritpion"])
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("isp", object["ISP"])
	d.Set("payment_type", convertEipAddressPaymentTypeResponse(object["ChargeType"]))
	d.Set("instance_charge_type", object["ChargeType"])
	d.Set("ip_address", object["IpAddress"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])

	listTagResourcesObject, err := vpcService.ListTagResources(d.Id(), "EIP")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceAlicloudEipAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "eip"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("deletion_protection") {
		update = true
	}
	request["ProtectionEnable"] = d.Get("deletion_protection")
	request["RegionId"] = client.RegionId
	request["Type"] = "EIP"
	if update {
		action := "DeletionProtection"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("DeletionProtection")
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
		d.SetPartial("deletion_protection")
	}
	update = false
	moveResourceGroupReq := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	moveResourceGroupReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	moveResourceGroupReq["NewResourceGroupId"] = d.Get("resource_group_id")
	moveResourceGroupReq["ResourceType"] = "eip"
	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, moveResourceGroupReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, moveResourceGroupReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	modifyEipAddressAttributeReq := map[string]interface{}{
		"AllocationId": d.Id(),
	}
	if !d.IsNewResource() && (d.HasChange("address_name") || d.HasChange("name")) {
		update = true
		if _, ok := d.GetOk("address_name"); ok {
			modifyEipAddressAttributeReq["Name"] = d.Get("address_name")
		} else {
			modifyEipAddressAttributeReq["Name"] = d.Get("name")
		}
	}
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		modifyEipAddressAttributeReq["Bandwidth"] = d.Get("bandwidth")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyEipAddressAttributeReq["Description"] = d.Get("description")
	}
	modifyEipAddressAttributeReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyEipAddressAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, modifyEipAddressAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyEipAddressAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("address_name")
		d.SetPartial("bandwidth")
		d.SetPartial("description")
	}
	d.Partial(false)
	return resourceAlicloudEipAddressRead(d, meta)
}
func resourceAlicloudEipAddressDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" || d.Get("instance_charge_type").(string) == "Prepaid" {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_eip_address. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "ReleaseEipAddress"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AllocationId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectEipStatus", "TaskConflict.AssociateGlobalAccelerationInstance"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"InvalidAllocationId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.EipAddressStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertEipAddressPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
func convertEipAddressPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
