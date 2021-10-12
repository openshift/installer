package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDedicatedHostCreate,
		Read:   resourceAlicloudEcsDedicatedHostRead,
		Update: resourceAlicloudEcsDedicatedHostUpdate,
		Delete: resourceAlicloudEcsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_on_maintenance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stop", "Migrate"}, false),
				Default:      "Stop",
			},
			"auto_placement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "on",
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_over_commit_ratio": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"dedicated_host_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail_fee": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"network_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"slb_udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Default:      "PostPaid",
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sale_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AllocateDedicatedHosts"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("action_on_maintenance"); ok {
		request["ActionOnMaintenance"] = v
	}

	if v, ok := d.GetOk("auto_placement"); ok {
		request["AutoPlacement"] = v
	}

	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	if v, ok := d.GetOk("cpu_over_commit_ratio"); ok {
		request["CpuOverCommitRatio"] = v
	}

	if v, ok := d.GetOk("dedicated_host_cluster_id"); ok {
		request["DedicatedHostClusterId"] = v
	}

	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request["DedicatedHostName"] = v
	}

	request["DedicatedHostType"] = d.Get("dedicated_host_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("expired_time"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("min_quantity"); ok {
		request["MinQuantity"] = v
	}

	if v, ok := d.GetOk("network_attributes"); ok {
		networkAttributesMap := make(map[string]interface{})
		for _, networkAttributes := range v.(*schema.Set).List() {
			networkAttributesArg := networkAttributes.(map[string]interface{})
			networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
			networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
		}
		request["NetworkAttributes"] = networkAttributesMap

	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}

	request["Quantity"] = 1
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("sale_cycle"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_dedicated_host", action, AlibabaCloudSdkGoERROR)
	}
	responseDedicatedHostIdSets := response["DedicatedHostIdSets"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseDedicatedHostIdSets["DedicatedHostId"].([]interface{})[0]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{"PermanentFailure"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsDedicatedHostRead(d, meta)
}
func resourceAlicloudEcsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDedicatedHost(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_dedicated_host ecsService.DescribeEcsDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("action_on_maintenance", object["ActionOnMaintenance"])
	d.Set("auto_placement", object["AutoPlacement"])
	d.Set("auto_release_time", object["AutoReleaseTime"])
	d.Set("cpu_over_commit_ratio", object["CpuOverCommitRatio"])
	d.Set("dedicated_host_name", object["DedicatedHostName"])
	d.Set("dedicated_host_type", object["DedicatedHostType"])
	d.Set("description", object["Description"])
	d.Set("expired_time", object["ExpiredTime"])

	networkAttributesSli := make([]map[string]interface{}, 0)
	if len(object["NetworkAttributes"].(map[string]interface{})) > 0 {
		networkAttributes := object["NetworkAttributes"]
		networkAttributesMap := make(map[string]interface{})
		networkAttributesMap["slb_udp_timeout"] = networkAttributes.(map[string]interface{})["SlbUdpTimeout"]
		networkAttributesMap["udp_timeout"] = networkAttributes.(map[string]interface{})["UdpTimeout"]
		networkAttributesSli = append(networkAttributesSli, networkAttributesMap)
	}
	d.Set("network_attributes", networkAttributesSli)
	d.Set("payment_type", object["ChargeType"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("sale_cycle", object["SaleCycle"])
	d.Set("status", object["Status"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("zone_id", object["ZoneId"])
	return nil
}
func resourceAlicloudEcsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "ddh"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if !d.IsNewResource() && d.HasChange("auto_release_time") {
		request := map[string]interface{}{
			"DedicatedHostId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["AutoReleaseTime"] = d.Get("auto_release_time")
		action := "ModifyDedicatedHostAutoReleaseTime"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("auto_release_time")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"ResourceId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = d.Get("resource_group_id")
		request["ResourceType"] = "ddh"
		action := "JoinResourceGroup"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update := false
	request := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	if !d.IsNewResource() && d.HasChange("expired_time") {
		update = true
		request["Period"] = d.Get("expired_time")
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("sale_cycle") {
		update = true
		request["PeriodUnit"] = d.Get("sale_cycle")
	}
	if update {
		action := "RenewDedicatedHosts"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("expired_time")
		d.SetPartial("sale_cycle")
	}
	update = false
	modifyDedicatedHostsChargeTypeReq := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	modifyDedicatedHostsChargeTypeReq["RegionId"] = client.RegionId
	modifyDedicatedHostsChargeTypeReq["AutoPay"] = true
	modifyDedicatedHostsChargeTypeReq["Period"] = d.Get("expired_time")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		modifyDedicatedHostsChargeTypeReq["DedicatedHostChargeType"] = d.Get("payment_type")
	}
	modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = d.Get("sale_cycle")
	if update {
		if _, ok := d.GetOkExists("detail_fee"); ok {
			modifyDedicatedHostsChargeTypeReq["DetailFee"] = d.Get("detail_fee")
		}
		if _, ok := d.GetOkExists("dry_run"); ok {
			modifyDedicatedHostsChargeTypeReq["DryRun"] = d.Get("dry_run")
		}
		action := "ModifyDedicatedHostsChargeType"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, modifyDedicatedHostsChargeTypeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyDedicatedHostsChargeTypeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("detail_fee")
		d.SetPartial("dry_run")
		d.SetPartial("expired_time")
		d.SetPartial("payment_type")
		d.SetPartial("sale_cycle")
	}
	update = false
	modifyDedicatedHostAttributeReq := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	modifyDedicatedHostAttributeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("action_on_maintenance") {
		update = true
		modifyDedicatedHostAttributeReq["ActionOnMaintenance"] = d.Get("action_on_maintenance")
	}
	if !d.IsNewResource() && d.HasChange("auto_placement") {
		update = true
		modifyDedicatedHostAttributeReq["AutoPlacement"] = d.Get("auto_placement")
	}
	if !d.IsNewResource() && d.HasChange("cpu_over_commit_ratio") {
		update = true
		modifyDedicatedHostAttributeReq["CpuOverCommitRatio"] = d.Get("cpu_over_commit_ratio")
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_name") {
		update = true
		modifyDedicatedHostAttributeReq["DedicatedHostName"] = d.Get("dedicated_host_name")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyDedicatedHostAttributeReq["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("network_attributes") {
		update = true
		if d.Get("network_attributes") != nil {
			networkAttributesMap := make(map[string]interface{})
			for _, networkAttributes := range d.Get("network_attributes").(*schema.Set).List() {
				networkAttributesArg := networkAttributes.(map[string]interface{})
				networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
				networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
			}
			modifyDedicatedHostAttributeReq["NetworkAttributes"] = networkAttributesMap
		}
	}
	if update {
		if _, ok := d.GetOk("dedicated_host_cluster_id"); ok {
			modifyDedicatedHostAttributeReq["DedicatedHostClusterId"] = d.Get("dedicated_host_cluster_id")
		}
		action := "ModifyDedicatedHostAttribute"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, modifyDedicatedHostAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyDedicatedHostAttributeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("action_on_maintenance")
		d.SetPartial("auto_placement")
		d.SetPartial("cpu_over_commit_ratio")
		d.SetPartial("dedicated_host_cluster_id")
		d.SetPartial("dedicated_host_name")
		d.SetPartial("description")
		d.SetPartial("network_attributes")
	}
	d.Partial(false)
	return resourceAlicloudEcsDedicatedHostRead(d, meta)
}
func resourceAlicloudEcsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseDedicatedHost"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectHostStatus.Initializing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDedicatedHostId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
