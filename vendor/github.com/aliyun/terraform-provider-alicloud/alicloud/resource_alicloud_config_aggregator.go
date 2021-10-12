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

func resourceAlicloudConfigAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigAggregatorCreate,
		Read:   resourceAlicloudConfigAggregatorRead,
		Update: resourceAlicloudConfigAggregatorUpdate,
		Delete: resourceAlicloudConfigAggregatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregator_accounts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"aggregator_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregator_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM", "RD"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudConfigAggregatorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	action := "CreateAggregator"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	aggregatorAccountsMaps := make([]map[string]interface{}, 0)
	for _, aggregatorAccounts := range d.Get("aggregator_accounts").(*schema.Set).List() {
		aggregatorAccountsArg := aggregatorAccounts.(map[string]interface{})
		aggregatorAccountsMap := map[string]interface{}{
			"AccountId":   aggregatorAccountsArg["account_id"],
			"AccountName": aggregatorAccountsArg["account_name"],
			"AccountType": aggregatorAccountsArg["account_type"],
		}
		aggregatorAccountsMaps = append(aggregatorAccountsMaps, aggregatorAccountsMap)
	}
	if v, err := convertArrayObjectToJsonString(aggregatorAccountsMaps); err == nil {
		request["AggregatorAccounts"] = v
	} else {
		return WrapError(err)
	}
	request["AggregatorName"] = d.Get("aggregator_name")
	if v, ok := d.GetOk("aggregator_type"); ok {
		request["AggregatorType"] = v
	}
	request["Description"] = d.Get("description")
	request["ClientToken"] = buildClientToken("CreateAggregator")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AggregatorId"]))
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigAggregatorStateRefreshFunc(d.Id(), []string{"0"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigAggregatorRead(d, meta)
}
func resourceAlicloudConfigAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigAggregator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregator configService.DescribeConfigAggregator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	aggregatorAccounts := make([]map[string]interface{}, 0)
	if aggregatorAccountsList, ok := object["AggregatorAccounts"].([]interface{}); ok {
		for _, v := range aggregatorAccountsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"account_id":   fmt.Sprint(m1["AccountId"]),
					"account_name": m1["AccountName"],
					"account_type": m1["AccountType"],
				}
				aggregatorAccounts = append(aggregatorAccounts, temp1)

			}
		}
	}
	if err := d.Set("aggregator_accounts", aggregatorAccounts); err != nil {
		return WrapError(err)
	}
	d.Set("aggregator_name", object["AggregatorName"])
	d.Set("aggregator_type", object["AggregatorType"])
	d.Set("description", object["Description"])
	d.Set("status", convertConfigAggregatorStatusResponse(formatInt(object["AggregatorStatus"])))
	return nil
}
func resourceAlicloudConfigAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AggregatorId": d.Id(),
	}
	if d.HasChange("aggregator_accounts") {
		update = true
	}
	aggregatorAccountsMaps := make([]map[string]interface{}, 0)
	for _, aggregatorAccounts := range d.Get("aggregator_accounts").(*schema.Set).List() {
		aggregatorAccountsArg := aggregatorAccounts.(map[string]interface{})
		aggregatorAccountsMap := map[string]interface{}{
			"AccountId":   aggregatorAccountsArg["account_id"],
			"AccountName": aggregatorAccountsArg["account_name"],
			"AccountType": aggregatorAccountsArg["account_type"],
		}
		aggregatorAccountsMaps = append(aggregatorAccountsMaps, aggregatorAccountsMap)
	}
	if v, err := convertArrayObjectToJsonString(aggregatorAccountsMaps); err == nil {
		request["AggregatorAccounts"] = v
	} else {
		return WrapError(err)
	}
	if d.HasChange("aggregator_name") {
		update = true
	}
	request["AggregatorName"] = d.Get("aggregator_name")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if update {
		action := "UpdateAggregator"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAggregator")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
	return resourceAlicloudConfigAggregatorRead(d, meta)
}
func resourceAlicloudConfigAggregatorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAggregators"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AggregatorIds": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteAggregators")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"AccountNotExisted", "Invalid.AggregatorIds.Empty"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertConfigAggregatorStatusResponse(source interface{}) interface{} {
	switch source {
	case 0:
		return "Creating"
	case 2:
		return "Deleting"
	case 1:
		return "Normal"
	}
	return ""
}
