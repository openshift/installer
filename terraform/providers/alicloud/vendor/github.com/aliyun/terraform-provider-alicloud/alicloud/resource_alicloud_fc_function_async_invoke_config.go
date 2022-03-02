package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudFCFunctionAsyncInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCAsyncInvokeConfigCreate,
		Read:   resourceAlicloudFCAsyncInvokeConfigRead,
		Update: resourceAlicloudFCAsyncInvokeConfigUpdate,
		Delete: resourceAlicloudFCAsyncInvokeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"function_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"qualifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "LATEST",
				ValidateFunc: validation.NoZeroValues,
			},
			"maximum_event_age_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(60, 21600),
			},
			"maximum_retry_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(0, 2),
			},
			"destination_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_failure": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:     schema.TypeString,
										Required: true,
										// TODO: Implement validateArn
										//ValidateFunc: validateArn,
									},
								},
							},
						},
						"on_success": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:     schema.TypeString,
										Required: true,
										// TODO: Implement validateArn
										//ValidateFunc: validateArn,
									},
								},
							},
						},
					},
				},
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCAsyncInvokeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)
	functionName := d.Get("function_name").(string)
	qualifier := d.Get("qualifier").(string)

	id := fmt.Sprintf("%s:%s:%s", serviceName, functionName, qualifier)

	request := &fc.PutFunctionAsyncInvokeConfigInput{
		ServiceName:  StringPointer(serviceName),
		FunctionName: StringPointer(functionName),
	}

	if v, ok := d.GetOk("maximum_retry_attempts"); ok {
		request.MaxAsyncRetryAttempts = Int64Pointer(int64(v.(int)))
	}
	if v, ok := d.GetOk("destination_config"); ok {
		request.DestinationConfig = parseFCDestinationConfig(v.([]interface{}))
	}

	if qualifier != "" {
		request.Qualifier = StringPointer(qualifier)
	}

	if v, ok := d.GetOk("maximum_event_age_in_seconds"); ok {
		request.MaxAsyncEventAgeInSeconds = Int64Pointer(int64(v.(int)))
	}

	var response *fc.PutFunctionAsyncInvokeConfigOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.PutFunctionAsyncInvokeConfig(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("PutFunctionAsyncInvokeConfig", raw, requestInfo, request)
		response, _ = raw.(*fc.PutFunctionAsyncInvokeConfigOutput)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_function_async_invoke_config", "PutFunctionAsyncInvokeConfig", FcGoSdk)
	}

	if response == nil {
		return WrapError(Error("Putting function async invoke config got an empty response"))
	}

	d.SetId(id)

	return resourceAlicloudFCAsyncInvokeConfigRead(d, meta)
}

func resourceAlicloudFCAsyncInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	object, err := fcService.DescribeFcFunctionAsyncInvokeConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fc_function_async_invoke_config: fcService.DescribeFcFunctionAsyncInvokeConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if err := d.Set("destination_config", flattenFCFunctionEventInvokeConfigDestinationConfig(object.DestinationConfig)); err != nil {
		return WrapError(err)
	}

	d.Set("service_name", object.Service)
	d.Set("function_name", object.Function)
	d.Set("qualifier", object.Qualifier)
	d.Set("maximum_event_age_in_seconds", object.MaxAsyncEventAgeInSeconds)
	d.Set("maximum_retry_attempts", object.MaxAsyncRetryAttempts)
	d.Set("created_time", object.CreatedTime)
	d.Set("last_modified_time", object.LastModifiedTime)

	return nil
}

func resourceAlicloudFCAsyncInvokeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName, functionName, qualifier, err := parseFCDestinationConfigId(d.Id())

	if err != nil {
		return WrapError(err)
	}

	request := &fc.PutFunctionAsyncInvokeConfigInput{
		ServiceName:  StringPointer(serviceName),
		FunctionName: StringPointer(functionName),
	}
	request.MaxAsyncRetryAttempts = Int64Pointer(int64(d.Get("maximum_retry_attempts").(int)))
	request.DestinationConfig = parseFCDestinationConfig(d.Get("destination_config").([]interface{}))
	request.Qualifier = StringPointer(qualifier)

	if v, ok := d.GetOk("maximum_event_age_in_seconds"); ok {
		request.MaxAsyncEventAgeInSeconds = Int64Pointer(int64(v.(int)))
	}

	var response *fc.PutFunctionAsyncInvokeConfigOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.PutFunctionAsyncInvokeConfig(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("PutFunctionAsyncInvokeConfig", raw, requestInfo, request)
		response, _ = raw.(*fc.PutFunctionAsyncInvokeConfigOutput)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_function_async_invoke_config", "PutFunctionAsyncInvokeConfig", FcGoSdk)
	}

	if response == nil {
		return WrapError(Error("Putting function async invoke config got an empty response"))
	}

	return resourceAlicloudFCAsyncInvokeConfigRead(d, meta)
}

func resourceAlicloudFCAsyncInvokeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	serviceName, functionName, qualifier, err := parseFCDestinationConfigId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	request := &fc.DeleteFunctionAsyncInvokeConfigInput{
		ServiceName:  &serviceName,
		FunctionName: &functionName,
		Qualifier:    &qualifier,
	}
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteFunctionAsyncInvokeConfig(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AsyncConfigNotExists"}) {
			return nil
		}
		b, _ := json.Marshal(request)
		return WrapErrorf(err, DefaultErrorMsg, string(b), "DeleteFunctionAsyncInvokeConfig", FcGoSdk)
	}
	addDebug("DeleteFunctionAsyncInvokeConfig", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcFunctionAsyncInvokeConfig(d.Id(), Deleted, DefaultTimeout))
}

func parseFCDestinationConfig(l []interface{}) *fc.DestinationConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := &fc.DestinationConfig{}
	m := l[0].(map[string]interface{})
	if v, ok := m["on_failure"].([]interface{}); ok {
		config.OnFailure = expandFCDestinationConfigOnFailure(v)
	}
	if v, ok := m["on_success"].([]interface{}); ok {
		config.OnSuccess = expandFCDestinationConfigOnSuccess(v)
	}
	return config
}

func expandFCDestinationConfigOnSuccess(l []interface{}) *fc.Destination {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	var onSuccess *fc.Destination
	if v, ok := m["destination"].(string); ok {
		onSuccess = &fc.Destination{
			Destination: StringPointer(v),
		}
	}

	return onSuccess
}

func expandFCDestinationConfigOnFailure(l []interface{}) *fc.Destination {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	var onFailure *fc.Destination
	if v, ok := m["destination"].(string); ok {
		onFailure = &fc.Destination{
			Destination: StringPointer(v),
		}
	}

	return onFailure
}

func parseFCDestinationConfigId(id string) (serviceName string, functionName string, qualifier string, err error) {
	idParts := strings.Split(id, ":")

	if len(idParts) == 2 && idParts[0] != "" && idParts[1] != "" {
		serviceName = idParts[0]
		functionName = idParts[1]
		return
	}

	if len(idParts) == 3 && idParts[0] != "" && idParts[1] != "" && idParts[2] != "" {
		serviceName = idParts[0]
		functionName = idParts[1]
		qualifier = idParts[2]
		return
	}

	return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected service_name:function_name or service_name:function_name:qualifier", id)
}

func flattenFCFunctionEventInvokeConfigDestinationConfig(destinationConfig *fc.DestinationConfig) []interface{} {
	// The API will respond with empty OnFailure and OnSuccess destinations when unconfigured:
	// "DestinationConfig":{"OnFailure":{"Destination":null},"OnSuccess":{"Destination":null}}
	// Return no destination configuration to prevent Terraform state difference

	if destinationConfig == nil {
		return []interface{}{}
	}

	if destinationConfig.OnFailure == nil && destinationConfig.OnSuccess == nil {
		return []interface{}{}
	}

	if (destinationConfig.OnFailure != nil && destinationConfig.OnFailure.Destination == nil) && (destinationConfig.OnSuccess != nil && destinationConfig.OnSuccess.Destination == nil) {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"on_failure": flattenFCFunctionEventInvokeConfigDestinationConfigOnFailure(destinationConfig.OnFailure),
		"on_success": flattenFCFunctionEventInvokeConfigDestinationConfigOnSuccess(destinationConfig.OnSuccess),
	}

	return []interface{}{m}
}

func flattenFCFunctionEventInvokeConfigDestinationConfigOnFailure(onFailure *fc.Destination) []interface{} {
	// The API will respond with empty OnFailure destination when unconfigured:
	// "DestinationConfig":{"OnFailure":{"Destination":null},"OnSuccess":{"Destination":null}}
	// Return no on failure configuration to prevent Terraform state difference

	if onFailure == nil || onFailure.Destination == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"destination": *onFailure.Destination,
	}

	return []interface{}{m}
}

func flattenFCFunctionEventInvokeConfigDestinationConfigOnSuccess(onSuccess *fc.Destination) []interface{} {
	// The API will respond with empty OnSuccess destination when unconfigured:
	// "DestinationConfig":{"OnFailure":{"Destination":null},"OnSuccess":{"Destination":null}}
	// Return no on success configuration to prevent Terraform state difference

	if onSuccess == nil || onSuccess.Destination == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"destination": *onSuccess.Destination,
	}

	return []interface{}{m}
}
