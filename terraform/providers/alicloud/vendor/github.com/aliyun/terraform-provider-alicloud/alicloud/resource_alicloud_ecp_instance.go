package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcpInstanceCreate,
		Read:   resourceAlicloudEcpInstanceRead,
		Update: resourceAlicloudEcpInstanceUpdate,
		Delete: resourceAlicloudEcpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"eip_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 128), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5}),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Year", "Month"}, false),
			},
			"resolution": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnc_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
		},
	}
}

func resourceAlicloudEcpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RunInstances"
	request := make(map[string]interface{})
	conn, err := client.NewCloudphoneClient()
	if err != nil {
		return WrapError(err)
	}
	request["Amount"] = 1
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEcpSyncPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("eip_bandwidth"); ok {
		request["EipBandwidth"] = v
	}
	request["ImageId"] = d.Get("image_id")
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resolution"); ok {
		request["Resolution"] = v
	}
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["ClientToken"] = buildClientToken("RunInstances")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecp_instance", action, AlibabaCloudSdkGoERROR)
	}
	responseInstanceIds := response["InstanceIds"].(map[string]interface{})
	cloudphoneService := CloudphoneService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudphoneService.EcpInstanceStateRefreshFunc(fmt.Sprint(responseInstanceIds["InstanceId"].([]interface{})[0]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId(fmt.Sprint(responseInstanceIds["InstanceId"].([]interface{})[0]))
	return resourceAlicloudEcpInstanceUpdate(d, meta)
}
func resourceAlicloudEcpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudphoneService := CloudphoneService{client}
	object, err := cloudphoneService.DescribeEcpInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecp_instance cloudphoneService.DescribeEcpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("image_id", object["ImageId"])
	d.Set("instance_name", object["InstanceName"])
	d.Set("instance_type", object["InstanceType"])
	d.Set("payment_type", convertEcpSyncPaymentTypeResponse(object["ChargeType"]))
	d.Set("key_pair_name", object["KeyPairName"])
	d.Set("resolution", object["Resolution"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["VpcAttributes"].(map[string]interface{})["VSwitchId"])
	return nil
}
func resourceAlicloudEcpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudphoneService := CloudphoneService{client}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		if v, ok := d.GetOk("instance_name"); ok {
			request["InstanceName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("key_pair_name") {
		update = true
		if v, ok := d.GetOk("key_pair_name"); ok {
			request["KeyPairName"] = v
		}
	}
	if d.HasChange("vnc_password") {
		update = true
		if v, ok := d.GetOk("vnc_password"); ok {
			request["VncPassword"] = v
		}
	}
	if update {
		action := "UpdateInstanceAttribute"
		conn, err := client.NewCloudphoneClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

		d.SetPartial("description")
		d.SetPartial("instance_name")
		d.SetPartial("key_pair_name")
		d.SetPartial("vnc_password")
	}

	if d.HasChange("status") {
		object, err := cloudphoneService.DescribeEcpInstance(d.Id())

		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != "Running" {
			if target == "Running" {
				request := map[string]interface{}{
					"InstanceId": []string{d.Id()},
				}
				request["RegionId"] = client.RegionId
				action := "StartInstances"
				conn, err := client.NewCloudphoneClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudphoneService.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

		}
		if object["Status"].(string) != "Stopped" {
			if target == "Stopped" {
				request := map[string]interface{}{
					"InstanceId": []string{d.Id()},
				}
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOkExists("force"); ok {
					request["Force"] = v
				}
				action := "StopInstances"
				conn, err := client.NewCloudphoneClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudphoneService.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudEcpInstanceRead(d, meta)
}
func resourceAlicloudEcpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstances"
	var response map[string]interface{}
	conn, err := client.NewCloudphoneClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": []string{d.Id()},
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	request["RegionId"] = client.RegionId

	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"CloudPhoneInstances.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertEcpSyncPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
func convertEcpSyncPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
