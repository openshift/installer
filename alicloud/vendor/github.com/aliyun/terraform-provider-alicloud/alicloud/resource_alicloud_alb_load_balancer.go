package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbLoadBalancerCreate,
		Read:   resourceAlicloudAlbLoadBalancerRead,
		Update: resourceAlicloudAlbLoadBalancerUpdate,
		Delete: resourceAlicloudAlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_log_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_store": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"address_allocated_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Dynamic", "Fixed"}, false),
			},
			"address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Intranet", "Internet"}, false),
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"load_balancer_billing_config": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pay_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
						},
					},
				},
				ForceNew: true,
			},
			"load_balancer_edition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Standard"}, false),
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modification_protection_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-.]{1,127}$`), "The reason must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."),
							DiffSuppressFunc: modificationProtectionConfigDiffSuppressFunc,
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"ConsoleProtection", "NonProtection"}, false),
						},
					},
				},
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudAlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLoadBalancer"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("address_allocated_mode"); ok {
		request["AddressAllocatedMode"] = v
	}
	request["AddressType"] = d.Get("address_type")
	if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
		request["DeletionProtectionEnabled"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	request["LoadBalancerName"] = d.Get("load_balancer_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	for _, zoneMappings := range d.Get("zone_mappings").(*schema.Set).List() {
		zoneMappingsArg := zoneMappings.(map[string]interface{})
		zoneMappingsMap := map[string]interface{}{}
		zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
		zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]
		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}
	request["ZoneMappings"] = zoneMappingsMaps
	loadBalancerBillingConfigMap := map[string]interface{}{}
	for _, loadBalancerBillingConfigs := range d.Get("load_balancer_billing_config").(*schema.Set).List() {
		loadBalancerBillingConfigArg := loadBalancerBillingConfigs.(map[string]interface{})
		loadBalancerBillingConfigMap["PayType"] = convertAlbLoadBalancerPaymentTypeRequest(loadBalancerBillingConfigArg["pay_type"].(string))
	}
	request["LoadBalancerBillingConfig"] = loadBalancerBillingConfigMap
	modificationProtectionConfigMap := map[string]interface{}{}
	for _, modificationProtectionConfigs := range d.Get("modification_protection_config").(*schema.Set).List() {
		modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
		modificationProtectionConfigMap["Reason"] = modificationProtectionConfigArg["reason"]
		modificationProtectionConfigMap["Status"] = modificationProtectionConfigArg["status"]
	}
	request["ModificationProtectionConfig"] = modificationProtectionConfigMap
	request["ClientToken"] = buildClientToken("CreateLoadBalancer")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))
	albService := AlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlbLoadBalancerUpdate(d, meta)
}
func resourceAlicloudAlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbLoadBalancer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer albService.DescribeAlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	accessLogConfigSli := make([]map[string]interface{}, 0)
	if object["AccessLogConfig"] != nil && len(object["AccessLogConfig"].(map[string]interface{})) > 0 {
		accessLogConfig := object["AccessLogConfig"].(map[string]interface{})
		accessLogConfigMap := make(map[string]interface{})
		accessLogConfigMap["log_project"] = accessLogConfig["LogProject"]
		accessLogConfigMap["log_store"] = accessLogConfig["LogStore"]
		accessLogConfigSli = append(accessLogConfigSli, accessLogConfigMap)
	}
	d.Set("access_log_config", accessLogConfigSli)
	d.Set("address_allocated_mode", object["AddressAllocatedMode"])
	d.Set("address_type", object["AddressType"])

	loadBalancerBillingConfigSli := make([]map[string]interface{}, 0)
	if object["LoadBalancerBillingConfig"] != nil && len(object["LoadBalancerBillingConfig"].(map[string]interface{})) > 0 {
		loadBalancerBillingConfig := object["LoadBalancerBillingConfig"]
		loadBalancerBillingConfigMap := make(map[string]interface{})
		loadBalancerBillingConfigMap["pay_type"] = convertAlbLoadBalancerPaymentTypeResponse(loadBalancerBillingConfig.(map[string]interface{})["PayType"].(string))
		loadBalancerBillingConfigSli = append(loadBalancerBillingConfigSli, loadBalancerBillingConfigMap)
	}
	d.Set("load_balancer_billing_config", loadBalancerBillingConfigSli)
	d.Set("load_balancer_edition", object["LoadBalancerEdition"])
	d.Set("load_balancer_name", object["LoadBalancerName"])

	modificationProtectionConfigSli := make([]map[string]interface{}, 0)
	if object["ModificationProtectionConfig"] != nil && len(object["ModificationProtectionConfig"].(map[string]interface{})) > 0 {
		modificationProtectionConfig := object["ModificationProtectionConfig"].(map[string]interface{})
		modificationProtectionConfigMap := make(map[string]interface{})
		modificationProtectionConfigMap["reason"] = modificationProtectionConfig["Reason"]
		modificationProtectionConfigMap["status"] = modificationProtectionConfig["Status"]
		modificationProtectionConfigSli = append(modificationProtectionConfigSli, modificationProtectionConfigMap)
	}
	d.Set("modification_protection_config", modificationProtectionConfigSli)
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["LoadBalancerStatus"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("vpc_id", object["VpcId"])
	if zoneMappingsList, ok := object["ZoneMappings"]; ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, zoneMappingsListItem := range zoneMappingsList.([]interface{}) {
			if zoneMappingsListItemMap, ok := zoneMappingsListItem.(map[string]interface{}); ok {
				zoneMappingsArg := map[string]interface{}{}
				zoneMappingsArg["vswitch_id"] = zoneMappingsListItemMap["VSwitchId"]
				zoneMappingsArg["zone_id"] = zoneMappingsListItemMap["ZoneId"]
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsArg)
			}
		}
		d.Set("zone_mappings", zoneMappingsMaps)
	}

	return nil
}
func resourceAlicloudAlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}
	if update {
		request["ResourceType"] = "loadbalancer"
		action := "MoveResourceGroup"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}

	update = false
	updateLoadBalancerEditionReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_edition") {
		update = true
		updateLoadBalancerEditionReq["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateLoadBalancerEdition"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateLoadBalancerEdition")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerEditionReq, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_edition")
	}

	if d.HasChange("access_log_config") {
		oraw, _ := d.GetChange("access_log_config")

		if oraw != nil && oraw.(*schema.Set).Len() > 1 {
			disableLoadBalancerAccessLogReq := map[string]interface{}{
				"LoadBalancerId": d.Id(),
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "DisableLoadBalancerAccessLog"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("DisableLoadBalancerAccessLog")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, disableLoadBalancerAccessLogReq, &runtime)
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

		} else {

			enableLoadBalancerAccessLogReq := map[string]interface{}{
				"LoadBalancerId": d.Id(),
			}

			if v, ok := d.GetOk("access_log_config"); ok {
				for _, enableLoadBalancerAccessLogs := range v.(*schema.Set).List() {
					enableLoadBalancerAccessArg := enableLoadBalancerAccessLogs.(map[string]interface{})
					enableLoadBalancerAccessLogReq["LogProject"] = enableLoadBalancerAccessArg["log_project"]
					enableLoadBalancerAccessLogReq["LogStore"] = enableLoadBalancerAccessArg["log_store"]
				}
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "EnableLoadBalancerAccessLog"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("EnableLoadBalancerAccessLog")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, enableLoadBalancerAccessLogReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IdempotenceProcessing", "SystemBusy", "Throttling"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("access_log_config")
	}

	update = false
	updateLoadBalancerAttributeReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		if v, ok := d.GetOk("load_balancer_name"); ok {
			updateLoadBalancerAttributeReq["LoadBalancerName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
		modificationProtectionConfigMap := map[string]interface{}{}
		if v, ok := d.GetOk("modification_protection_config"); ok {
			for _, modificationProtectionConfigs := range v.(*schema.Set).List() {
				modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
				modificationProtectionConfigMap["Reason"] = modificationProtectionConfigArg["reason"]
				modificationProtectionConfigMap["Status"] = modificationProtectionConfigArg["status"]
			}
		}
		updateLoadBalancerAttributeReq["ModificationProtectionConfig"] = modificationProtectionConfigMap
	}

	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}

		action := "UpdateLoadBalancerAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		updateLoadBalancerAttributeReq["ClientToken"] = buildClientToken("UpdateLoadBalancerAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerAttributeReq, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_name")
		d.SetPartial("modification_protection_config")
	}

	if d.HasChange("deletion_protection_enabled") {
		target := strconv.FormatBool(d.Get("deletion_protection_enabled").(bool))
		if target == "false" {
			request := map[string]interface{}{
				"ResourceId": d.Id(),
			}
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "DisableDeletionProtection"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("DisableDeletionProtection")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IdempotenceProcessing", "SystemBusy", "Throttling"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if target == "true" {
			request := map[string]interface{}{
				"ResourceId": d.Id(),
			}
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "EnableDeletionProtection"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("EnableDeletionProtection")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"IdempotenceProcessing", "SystemBusy", "Throttling"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbLoadBalancerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("deletion_protection_enabled")
	}

	d.Partial(false)
	return resourceAlicloudAlbLoadBalancerRead(d, meta)
}
func resourceAlicloudAlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteLoadBalancer")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "ResourceNotFound.LoadBalancer", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertAlbLoadBalancerPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}
	return source
}
func convertAlbLoadBalancerPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}
	return source
}
func modificationProtectionConfigDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {

	if v, ok := d.GetOk("modification_protection_config"); ok {
		val := v.(*schema.Set).List()
		if len(val) > 2 {
			// modification_protection_config 为 Object 类型
			return true
		}
		for _, modificationProtectionConfigs := range val {
			modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
			return fmt.Sprintf(modificationProtectionConfigArg["status"].(string)) != "ConsoleProtection"
		}
	}

	return true
}
