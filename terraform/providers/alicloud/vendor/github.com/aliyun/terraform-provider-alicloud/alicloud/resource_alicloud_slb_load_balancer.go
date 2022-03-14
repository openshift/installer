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

func resourceAlicloudSlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbLoadBalancerCreate,
		Read:   resourceAlicloudSlbLoadBalancerRead,
		Update: resourceAlicloudSlbLoadBalancerUpdate,
		Delete: resourceAlicloudSlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(9 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"internet": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Removed:  "Field 'internet' has been removed from provider version 1.124. Use 'address_type' replaces it.",
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("address_type"); ok && v.(string) == "internet" {
						return true
					}
					return false
				},
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
				Default:      "ipv4",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("address_type"); ok && v.(string) == "intranet" {
						return true
					}
					return false
				},
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 1000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("internet_charge_type").(string) == "PayByTraffic" {
						return true
					}
					return false
				},
			},
			"delete_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "off",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return true
					}
					if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "PrePaid" {
						return true
					}
					return false
				},
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, true),
				Default:      "PayByTraffic",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("address_type"); ok && v.(string) == "intranet" {
						return true
					}
					return false
				},
			},
			"load_balancer_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.123.1. New field 'load_balancer_name' instead",
				ConflictsWith: []string{"load_balancer_name"},
			},
			"load_balancer_spec": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"slb.s1.small", "slb.s2.medium", "slb.s2.small", "slb.s3.large", "slb.s3.medium", "slb.s3.small", "slb.s4.large"}, false),
				ConflictsWith: []string{"specification"},
			},
			"specification": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"slb.s1.small", "slb.s2.medium", "slb.s2.small", "slb.s3.large", "slb.s3.medium", "slb.s3.small", "slb.s4.large"}, false),
				Deprecated:    "Field 'specification' has been deprecated from provider version 1.123.1. New field 'load_balancer_spec' instead",
				ConflictsWith: []string{"load_balancer_spec"},
			},
			"master_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"modification_protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("modification_protection_status"); ok && v.(string) == "NonProtection" {
						return true
					}
					return false
				},
			},
			"modification_protection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ConsoleProtection", "NonProtection"}, false),
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
				ValidateFunc:  validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ConflictsWith: []string{"instance_charge_type"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"payment_type"},
				Deprecated:    "Field 'instance_charge_type' has been deprecated from provider version 1.124. Use 'payment_type' replaces it.",
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"slave_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "inactive"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("address_type"); ok && v.(string) == "internet" {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceAlicloudSlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	var response map[string]interface{}
	action := "CreateLoadBalancer"
	request := make(map[string]interface{})
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("address"); ok {
		request["Address"] = v
	}

	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}

	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = v
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOk("delete_protection"); ok {
		request["DeleteProtection"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = convertSlbLoadBalancerInternetChargeTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["LoadBalancerName"] = v
	}

	if v, ok := d.GetOk("load_balancer_spec"); ok {
		request["LoadBalancerSpec"] = v
	} else if v, ok := d.GetOk("specification"); ok {
		request["LoadBalancerSpec"] = v
	}

	if v, ok := d.GetOk("master_zone_id"); ok {
		request["MasterZoneId"] = v
	}

	if v, ok := d.GetOk("modification_protection_reason"); ok {
		request["ModificationProtectionReason"] = v
	}

	if v, ok := d.GetOk("modification_protection_status"); ok {
		request["ModificationProtectionStatus"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertSlbLoadBalancerPaymentTypeRequest(v.(string))
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["PayType"] = convertSlbLoadBalancerInstanceChargeTypeRequest(v.(string))
	}
	if v, ok := request["PayType"]; ok && v.(string) == "PrePay" {
		period := 1
		if v, ok := d.GetOk("period"); ok {
			period = v.(int)
		}
		request["Duration"] = period
		request["PricingCycle"] = string(Month)
		if period > 9 {
			request["Duration"] = period / 12
			request["PricingCycle"] = string(Year)
		}
		request["AutoPay"] = true
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("slave_zone_id"); ok {
		request["SlaveZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["VpcId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
	}
	request["ClientToken"] = buildClientToken("CreateLoadBalancer")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.TokenIsProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, slbService.SlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudSlbLoadBalancerUpdate(d, meta)
}
func resourceAlicloudSlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbLoadBalancer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_load_balancer slbService.DescribeSlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("address", object["Address"])
	d.Set("address_ip_version", object["AddressIPVersion"])
	d.Set("address_type", object["AddressType"])
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("delete_protection", object["DeleteProtection"])
	d.Set("internet_charge_type", convertSlbLoadBalancerInternetChargeTypeResponse(object["InternetChargeType"]))
	d.Set("load_balancer_name", object["LoadBalancerName"])
	d.Set("name", object["LoadBalancerName"])
	d.Set("load_balancer_spec", object["LoadBalancerSpec"])
	d.Set("specification", object["LoadBalancerSpec"])
	d.Set("master_zone_id", object["MasterZoneId"])
	d.Set("modification_protection_reason", object["ModificationProtectionReason"])
	d.Set("modification_protection_status", object["ModificationProtectionStatus"])
	d.Set("payment_type", convertSlbLoadBalancerPaymentTypeResponse(object["PayType"]))
	d.Set("instance_charge_type", convertSlbLoadBalancerInstanceChargeTypeResponse(object["PayType"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("slave_zone_id", object["SlaveZoneId"])
	d.Set("status", object["LoadBalancerStatus"])
	d.Set("vswitch_id", object["VSwitchId"])

	listTagResourcesObject, err := slbService.ListTagResources(d.Id(), "instance")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceAlicloudSlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := slbService.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("status") {
		request := map[string]interface{}{
			"LoadBalancerId": d.Id(),
		}
		request["LoadBalancerStatus"] = d.Get("status")
		action := "SetLoadBalancerStatus"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("status")
	}
	if !d.IsNewResource() && d.HasChange("delete_protection") {
		request := map[string]interface{}{
			"LoadBalancerId": d.Id(),
		}
		request["DeleteProtection"] = d.Get("delete_protection")
		request["RegionId"] = client.RegionId
		action := "SetLoadBalancerDeleteProtection"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("delete_protection")
	}
	update := false
	request := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	} else if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["LoadBalancerName"] = d.Get("name")
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "SetLoadBalancerName"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("name")
		d.SetPartial("load_balancer_name")
	}
	update = false
	modifyLoadBalancerInstanceSpecReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_spec") {
		update = true
		modifyLoadBalancerInstanceSpecReq["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	} else if !d.IsNewResource() && d.HasChange("specification") {
		update = true
		modifyLoadBalancerInstanceSpecReq["LoadBalancerSpec"] = d.Get("specification")
	}
	modifyLoadBalancerInstanceSpecReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyLoadBalancerInstanceSpec"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, modifyLoadBalancerInstanceSpecReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyLoadBalancerInstanceSpecReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("specification")
		d.SetPartial("load_balancer_spec")
	}
	update = false
	setLoadBalancerModificationProtectionReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("modification_protection_status") {
		update = true
	}
	setLoadBalancerModificationProtectionReq["ModificationProtectionStatus"] = d.Get("modification_protection_status")
	setLoadBalancerModificationProtectionReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("modification_protection_reason") {
		update = true
		setLoadBalancerModificationProtectionReq["ModificationProtectionReason"] = d.Get("modification_protection_reason")
	}
	if update {
		action := "SetLoadBalancerModificationProtection"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, setLoadBalancerModificationProtectionReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, setLoadBalancerModificationProtectionReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("modification_protection_status")
		d.SetPartial("modification_protection_reason")
	}
	update = false
	modifyLoadBalancerInternetSpecReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		modifyLoadBalancerInternetSpecReq["Bandwidth"] = d.Get("bandwidth")
	}
	if !d.IsNewResource() && d.HasChange("internet_charge_type") {
		update = true
		modifyLoadBalancerInternetSpecReq["InternetChargeType"] = d.Get("internet_charge_type")
	}
	modifyLoadBalancerInternetSpecReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyLoadBalancerInternetSpec"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, modifyLoadBalancerInternetSpecReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyLoadBalancerInternetSpecReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("bandwidth")
		d.SetPartial("internet_charge_type")
	}
	update = false
	modifyLoadBalancerPayTypeReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	modifyLoadBalancerPayTypeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		modifyLoadBalancerPayTypeReq["PayType"] = convertSlbLoadBalancerPaymentTypeRequest(d.Get("payment_type").(string))
	}
	if !d.IsNewResource() && d.HasChange("instance_charge_type") {
		update = true
		modifyLoadBalancerPayTypeReq["PayType"] = convertSlbLoadBalancerInstanceChargeTypeRequest(d.Get("instance_charge_type").(string))
	}
	if v, ok := modifyLoadBalancerPayTypeReq["PayType"]; ok && v.(string) == "PrePay" {
		period := 1
		if v, ok := d.GetOk("period"); ok {
			period = v.(int)
		}
		modifyLoadBalancerPayTypeReq["Duration"] = period
		modifyLoadBalancerPayTypeReq["PricingCycle"] = string(Month)
		if period > 9 {
			modifyLoadBalancerPayTypeReq["Duration"] = period / 12
			modifyLoadBalancerPayTypeReq["PricingCycle"] = string(Year)
		}
		modifyLoadBalancerPayTypeReq["AutoPay"] = true
	}
	if update {
		action := "ModifyLoadBalancerPayType"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, modifyLoadBalancerPayTypeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyLoadBalancerPayTypeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_charge_type")
		d.SetPartial("payment_type")
	}
	d.Partial(false)
	return resourceAlicloudSlbLoadBalancerRead(d, meta)
}
func resourceAlicloudSlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" || d.Get("instance_charge_type").(string) == "Prepaid" {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_slb_load_balancer. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	action := "DeleteLoadBalancer"
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, slbService.SlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertSlbLoadBalancerInternetChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayByBandwidth":
		return "paybybandwidth"
	case "PayByTraffic":
		return "paybytraffic"
	}
	return source
}
func convertSlbLoadBalancerPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PayOnDemand"
	case "Subscription":
		return "PrePay"
	}
	return source
}
func convertSlbLoadBalancerInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayOnDemand"
	case "PrePaid":
		return "PrePay"
	}
	return source
}
func convertSlbLoadBalancerInternetChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "paybybandwidth":
		return "PayByBandwidth"
	case "paybytraffic":
		return "PayByTraffic"
	}
	return source
}
func convertSlbLoadBalancerPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PayOnDemand":
		return "PayAsYouGo"
	case "PrePay":
		return "Subscription"
	}
	return source
}
func convertSlbLoadBalancerInstanceChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PayOnDemand":
		return "PostPaid"
	case "PrePay":
		return "PrePaid"
	}
	return source
}
