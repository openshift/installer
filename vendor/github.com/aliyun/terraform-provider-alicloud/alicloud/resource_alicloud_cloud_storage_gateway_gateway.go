package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCloudStorageGatewayGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayGatewayCreate,
		Read:   resourceAlicloudCloudStorageGatewayGatewayRead,
		Update: resourceAlicloudCloudStorageGatewayGatewayUpdate,
		Delete: resourceAlicloudCloudStorageGatewayGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Cloud", "On_Premise"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"storage_bundle_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"File", "Iscsi"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_class": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Standard", "Enhanced", "Advanced"}, false),
				Optional:     true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				Default:      "PayAsYouGo",
			},

			"public_network_bandwidth": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(5, 200),
				Computed:     true,
				Optional:     true,
			},
			"reason_detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release_after_expiration": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGateway"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("gateway_class"); ok {
		request["GatewayClass"] = v
	}
	request["Location"] = d.Get("location")
	request["Name"] = d.Get("gateway_name")
	if v, ok := d.GetOk("payment_type"); ok {
		request["PostPaid"] = convertCsgGatewayPaymentTypeReq(v.(string))
	}

	if v, ok := d.GetOk("public_network_bandwidth"); ok {
		request["PublicNetworkBandwidth"] = v
	}
	if v, ok := d.GetOkExists("release_after_expiration"); ok {
		request["ReleaseAfterExpiration"] = v
	}
	request["StorageBundleId"] = d.Get("storage_bundle_id")
	request["Type"] = d.Get("type")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"BadRequest"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	d.SetId(fmt.Sprint(response["GatewayId"]))

	if d.Get("location").(string) == "Cloud" {
		action = "DeployGateway"
		request = map[string]interface{}{
			"GatewayId": d.Id(),
		}
		if id, ok := response["GatewayId"]; ok {
			request["GatewayId"] = fmt.Sprint(id)
		}
		if v, ok := d.GetOk("gateway_class"); ok {
			request["GatewayClass"] = v
		}
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"BadRequest"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	err = gatewayDescribeTasks(d, meta, d.Id())
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudCloudStorageGatewayGatewayRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	var object map[string]interface{}
	var err error
	object, err = sgwService.DescribeCloudStorageGatewayGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway sgwService.DescribeCloudStorageGatewayGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("payment_type", convertCsgGatewayPaymentTypeResp(object["IsPostPaid"].(bool)))

	d.Set("description", object["Description"])
	d.Set("gateway_class", object["GatewayClass"])
	d.Set("gateway_name", object["Name"])
	d.Set("location", object["Location"])
	if v, ok := object["PublicNetworkBandwidth"]; ok && fmt.Sprint(v) != "0" {
		d.Set("public_network_bandwidth", formatInt(v))
	}
	d.Set("status", object["Status"])
	d.Set("storage_bundle_id", object["StorageBundleId"])
	d.Set("type", object["Type"])
	d.Set("vswitch_id", object["VSwitchId"])
	return nil
}
func resourceAlicloudCloudStorageGatewayGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewHcsSgwClient()
	var response map[string]interface{}
	d.Partial(true)
	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}
	if d.HasChange("description") || d.HasChange("gateway_name") {
		action := "ModifyGateway"
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		if v, ok := d.GetOk("gateway_name"); ok {
			request["Name"] = v
		}

		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("description")
		d.SetPartial("gateway_name")
	}

	if d.HasChange("public_network_bandwidth") {
		err = gatewayDescribeTasks(d, meta, d.Id())
		if err != nil {
			return WrapError(err)
		}
		action := "ExpandGatewayNetworkBandwidth"
		request = map[string]interface{}{
			"GatewayId": d.Id(),
		}
		if v, ok := d.GetOk("public_network_bandwidth"); ok {
			request["NewNetworkBandwidth"] = v
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		err = gatewayDescribeTasks(d, meta, d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("public_network_bandwidth")
	}

	if d.HasChange("gateway_class") {
		err = gatewayDescribeTasks(d, meta, d.Id())
		if err != nil {
			return WrapError(err)
		}

		action := "ModifyGatewayClass"
		request = map[string]interface{}{
			"GatewayId": d.Id(),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		if v, ok := d.GetOk("gateway_class"); ok {
			request["GatewayClass"] = v
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			d.SetPartial("gateway_class")

			err = gatewayDescribeTasks(d, meta, d.Id())
			if err != nil {
				return WrapError(err)
			}
		}

	}
	d.Partial(false)
	return resourceAlicloudCloudStorageGatewayGatewayRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewHcsSgwClient()
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)

	err = gatewayDescribeTasks(d, meta, d.Id())
	if err != nil {
		return WrapError(err)
	}

	action := "DeleteGateway"
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}
	if v, ok := d.GetOk("reason_detail"); ok {
		request["ReasonDetail"] = v
	}
	if v, ok := d.GetOk("reason_type"); ok {
		request["ReasonType"] = v
	}
	wait = incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	addDebug(action, response, request)
	err = gatewayDescribeTasks(d, meta, d.Id())
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func gatewayDescribeTasks(d *schema.ResourceData, meta interface{}, TargetId string) error {
	action := "DescribeTasks"
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewHcsSgwClient()

	var response map[string]interface{}

	request := map[string]interface{}{
		"GatewayId": d.Id(),
	}
	request["TargetId"] = TargetId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		tasks := response["Tasks"].(map[string]interface{})
		for _, val := range tasks {
			for _, task := range val.([]interface{}) {
				if state, exist := task.(map[string]interface{})["StateCode"]; exist {
					if state != "task.state.completed" {
						return resource.RetryableError(fmt.Errorf("There are still tasks left"))
					}
					continue
				}
			}
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func convertCsgGatewayPaymentTypeReq(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return true
	}
	return source
}

func convertCsgGatewayPaymentTypeResp(source interface{}) interface{} {
	switch source {
	case true:
		return "PayAsYouGo"
	case false:
		return "Subscription"
	}
	return source
}
