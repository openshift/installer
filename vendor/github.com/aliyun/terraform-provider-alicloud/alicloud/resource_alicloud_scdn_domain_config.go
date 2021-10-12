package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudScdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudScdnDomainConfigCreate,
		Read:   resourceAlicloudScdnDomainConfigRead,
		Update: resourceAlicloudScdnDomainConfigUpdate,
		Delete: resourceAlicloudScdnDomainConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(5, 67),
			},
			"function_args": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arg_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"arg_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudScdnDomainConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := ScdnService{client}
	var response map[string]interface{}
	action := "BatchSetScdnDomainConfigs"
	request := make(map[string]interface{})
	conn, err := client.NewScdnClient()
	if err != nil {
		return WrapError(err)
	}

	config := make([]map[string]interface{}, 1)
	functionArgs := d.Get("function_args").(*schema.Set).List()
	args := make([]map[string]interface{}, len(functionArgs))
	for key, value := range functionArgs {
		arg := value.(map[string]interface{})
		args[key] = map[string]interface{}{
			"argName":  arg["arg_name"],
			"argValue": arg["arg_value"],
		}
	}
	config[0] = map[string]interface{}{
		"functionArgs": args,
		"functionName": d.Get("function_name").(string),
	}
	bytconfig, _ := json.Marshal(config)
	request["DomainNames"] = d.Get("domain_name").(string)
	request["Functions"] = string(bytconfig)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_scdn_domain", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%s:%s:%s", request["DomainNames"], d.Get("function_name"), ""))

	object, err := scdnService.DescribeScdnDomainConfig(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s:%s:%s", request["DomainNames"], object["FunctionName"], object["ConfigId"]))

	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 2*time.Second, scdnService.ScdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudScdnDomainConfigRead(d, meta)
}

func resourceAlicloudScdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := &ScdnService{client: client}

	object, err := scdnService.DescribeScdnDomainConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	functionArgMaps := make([]map[string]interface{}, 0)

	for _, functionArg := range object["FunctionArgs"].(map[string]interface{})["FunctionArg"].([]interface{}) {
		functionArgItem := functionArg.(map[string]interface{})
		functionArgMap := make(map[string]interface{}, 0)
		functionArgMap["arg_name"] = functionArgItem["ArgName"]
		functionArgMap["arg_value"] = functionArgItem["ArgValue"]
		functionArgMaps = append(functionArgMaps, functionArgMap)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("config_id", parts[2])
	d.Set("domain_name", parts[0])
	d.Set("function_args", functionArgMaps)
	d.Set("function_name", parts[1])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudScdnDomainConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := ScdnService{client}
	var response map[string]interface{}
	if d.HasChange("function_args") {
		action := "BatchSetScdnDomainConfigs"
		request := make(map[string]interface{})
		conn, err := client.NewScdnClient()
		if err != nil {
			return WrapError(err)
		}
		parts, err := ParseResourceId(d.Id(), 3)
		if err != nil {
			return WrapError(err)
		}

		config := make([]map[string]interface{}, 1)
		functionArgs := d.Get("function_args").(*schema.Set).List()
		args := make([]map[string]interface{}, len(functionArgs))
		for key, value := range functionArgs {
			arg := value.(map[string]interface{})
			args[key] = map[string]interface{}{
				"argName":  arg["arg_name"],
				"argValue": arg["arg_value"],
			}
		}
		config[0] = map[string]interface{}{
			"functionArgs": args,
			"functionName": parts[1],
			"configId":     parts[2],
		}
		bytconfig, _ := json.Marshal(config)
		request["DomainNames"] = parts[0]
		request["Functions"] = string(bytconfig)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 2*time.Second, scdnService.ScdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudScdnDomainConfigRead(d, meta)
}

func resourceAlicloudScdnDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	scdnService := ScdnService{client}
	var response map[string]interface{}
	action := "DeleteScdnSpecificConfig"
	conn, err := client.NewScdnClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := make(map[string]interface{})
	request["ConfigId"] = parts[2]
	request["DomainName"] = parts[0]
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Second, scdnService.ScdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
