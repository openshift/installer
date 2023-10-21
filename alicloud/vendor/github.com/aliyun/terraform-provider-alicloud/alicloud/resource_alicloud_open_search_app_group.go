package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOpenSearchAppGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOpenSearchAppGroupCreate,
		Read:   resourceAlicloudOpenSearchAppGroupRead,
		Update: resourceAlicloudOpenSearchAppGroupUpdate,
		Delete: resourceAlicloudOpenSearchAppGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "enhanced"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"quota": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"doc_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"compute_resource": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"qps": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"opensearch.share.junior", "opensearch.share.common", "opensearch.share.compute", "opensearch.share.storage", "opensearch.private.common", "opensearch.private.compute", "opensearch.private.storage"}, false),
						},
					},
				},
			},
			"order": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pricing_cycle": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Year", "Month"}, false),
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"charge_way": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"compute_resource", "qps"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudOpenSearchAppGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/v4/openapi/app-groups"
	body := make(map[string]interface{})
	conn, err := client.NewOpensearchClient()
	if err != nil {
		return WrapError(err)
	}
	body["name"] = d.Get("app_group_name")
	body["type"] = d.Get("type")
	body["chargeType"] = convertOpenSearchAppGroupPaymentTypeRequest(d.Get("payment_type").(string))
	quotaMaps := make(map[string]interface{}, 0)
	for _, quotas := range d.Get("quota").(*schema.Set).List() {
		quota := quotas.(map[string]interface{})
		quotaMaps = map[string]interface{}{
			"docSize":         quota["doc_size"],
			"computeResource": quota["compute_resource"],
			"spec":            quota["spec"],
			"qps":             quota["qps"],
		}
	}
	body["quota"] = quotaMaps

	if _, ok := d.GetOk("order"); ok {
		orderMaps := make(map[string]interface{}, 0)
		for _, quotas := range d.Get("order").(*schema.Set).List() {
			quota := quotas.(map[string]interface{})
			orderMaps = map[string]interface{}{
				"duration":     quota["duration"],
				"pricingCycle": quota["pricing_cycle"],
				"autoRenew":    quota["auto_renew"],
			}
		}
		body["order"] = orderMaps
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("POST "+action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_open_search_app_group", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	d.SetId(fmt.Sprint(response["result"].(map[string]interface{})["name"]))

	return resourceAlicloudOpenSearchAppGroupUpdate(d, meta)
}
func resourceAlicloudOpenSearchAppGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	openSearchService := OpenSearchService{client}
	object, err := openSearchService.DescribeOpenSearchAppGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_open_search_app_group openSearchService.DescribeOpenSearchAppGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("payment_type", convertOpenSearchAppGroupPaymentTypeResponse(object["chargeType"].(string)))
	d.Set("app_group_name", object["name"])
	d.Set("description", object["description"])
	d.Set("type", object["type"])
	d.Set("current_version", object["currentVersion"])
	d.Set("status", object["status"])
	d.Set("charge_way", convertOpenSearchchargingWayResponse(object["chargingWay"].(json.Number).String()))
	quotaSli := make([]map[string]interface{}, 0)
	if _, exist := object["quota"]; exist {
		quotaval := object["quota"].(map[string]interface{})
		quotaSli = append(quotaSli, map[string]interface{}{
			"doc_size":         quotaval["docSize"],
			"compute_resource": quotaval["computeResource"],
			"spec":             quotaval["spec"],
			"qps":              quotaval["qps"],
		})
	}
	d.Set("quota", quotaSli)
	return nil
}

func resourceAlicloudOpenSearchAppGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	openSearchService := OpenSearchService{client}
	var response map[string]interface{}
	update := false
	body := map[string]interface{}{
		"appGroupIdentity": d.Id(),
	}
	if d.HasChange("description") || d.IsNewResource() {
		if v, ok := d.GetOk("description"); ok {
			body["description"] = v
			update = true
		}
	}
	if d.HasChange("current_version") || d.IsNewResource() {
		if v, ok := d.GetOk("current_version"); ok {
			body["currentVersion"] = v
			update = true
		}
	}
	if d.HasChange("charge_way") || d.IsNewResource() {
		if v, ok := d.GetOk("charge_way"); ok {
			body["chargingWay"] = convertOpenSearchchargingWayRequest(v.(string))
			update = true
		}
	}
	if update {
		action := "/v4/openapi/app-groups/" + d.Id()
		conn, err := client.NewOpensearchClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug("PUT "+action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if code, exist := response["code"]; exist && code.(string) != "Success" {
			return WrapError(Error("Update AppGroup failed for " + response["message"].(string)))
		}

	}
	update = false
	if d.HasChange("quota") && !d.IsNewResource() {
		update = true
		if val, exist := d.GetOk("order_type"); exist {
			body["orderType"] = val
		} else {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is updated ", "order_type", "quota"))
		}
		for _, quotas := range d.Get("quota").(*schema.Set).List() {
			quota := quotas.(map[string]interface{})
			body["docSize"] = quota["doc_size"]
			body["computeResource"] = quota["compute_resource"]
			body["spec"] = quota["spec"]
		}
	}
	if update {
		action := fmt.Sprintf("/v4/openapi/app-groups/%s/quota", d.Id())
		conn, err := client.NewOpensearchClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug("PUT "+action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if code, exist := response["code"]; exist && code.(string) != "Success" {
			return WrapError(Error("Update AppGroup failed for " + response["message"].(string)))
		}
	}
	stateConf := BuildStateConf([]string{}, []string{"config_pending", "normal"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, openSearchService.OpenSearchAppStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudOpenSearchAppGroupRead(d, meta)
}

func resourceAlicloudOpenSearchAppGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.Get("payment_type").(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resource: alicloud_open_search_app_group. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	action := "/v4/openapi/app-groups/" + d.Id()
	var response map[string]interface{}
	conn, err := client.NewOpensearchClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"appGroupIdentity": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-12-25"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("DELETE "+action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertOpenSearchAppGroupPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}
func convertOpenSearchAppGroupPaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertOpenSearchchargingWayRequest(source string) interface{} {
	switch source {
	case "compute_resource":
		return 1
	case "qps":
		return 2
	}
	return source
}
func convertOpenSearchchargingWayResponse(source string) interface{} {
	switch source {
	case "1":
		return "compute_resource"
	case "2":
		return "qps"
	}
	return source
}
