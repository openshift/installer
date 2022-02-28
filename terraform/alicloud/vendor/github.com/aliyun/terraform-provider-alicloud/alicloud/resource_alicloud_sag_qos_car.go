package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSagQosCar() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagQosCarCreate,
		Read:   resourceAlicloudSagQosCarRead,
		Update: resourceAlicloudSagQosCarUpdate,
		Delete: resourceAlicloudSagQosCarDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"qos_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Absolute", "Percent"}, false),
			},
			"min_bandwidth_abs": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_bandwidth_abs": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_bandwidth_percent": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_bandwidth_percent": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"percent_source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"InternetUpBandwidth", "CcnBandwidth"}, false),
				Default:      "InternetUpBandwidth",
			},
		},
	}
}

func resourceAlicloudSagQosCarCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateQosCarRequest()

	request.QosId = d.Get("qos_id").(string)
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	request.LimitType = d.Get("limit_type").(string)
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("limit_type"); ok && v.(string) == "Absolute" {
		request.MinBandwidthAbs = requests.NewInteger(d.Get("min_bandwidth_abs").(int))
		request.MaxBandwidthAbs = requests.NewInteger(d.Get("max_bandwidth_abs").(int))
	}
	if v, ok := d.GetOk("limit_type"); ok && v.(string) == "Percent" {
		request.MinBandwidthPercent = requests.NewInteger(d.Get("min_bandwidth_percent").(int))
		request.MaxBandwidthPercent = requests.NewInteger(d.Get("max_bandwidth_percent").(int))
		request.PercentSourceType = d.Get("percent_source_type").(string)
	}
	var response *smartag.CreateQosCarResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.CreateQosCar(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceInOperating"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*smartag.CreateQosCarResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_qos_car", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", response.QosId, COLON_SEPARATED, response.QosCarId))

	return resourceAlicloudSagQosCarRead(d, meta)
}

func resourceAlicloudSagQosCarRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagQosCar(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("qos_id", object.QosId)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("priority", object.Priority)
	d.Set("limit_type", object.LimitType)
	d.Set("min_bandwidth_abs", object.MinBandwidthAbs)
	d.Set("max_bandwidth_abs", object.MaxBandwidthAbs)
	d.Set("min_bandwidth_percent", object.MinBandwidthPercent)
	d.Set("max_bandwidth_percent", object.MaxBandwidthPercent)
	d.Set("percent_source_type", object.PercentSourceType)

	return nil
}

func resourceAlicloudSagQosCarUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartag.CreateModifyQosCarRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.QosId = parts[0]
	request.QosCarId = parts[1]

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("priority") {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
		update = true
	}
	if d.HasChange("limit_type") {
		request.LimitType = d.Get("limit_type").(string)
		update = true
	}
	if d.HasChange("min_bandwidth_abs") {
		request.MinBandwidthAbs = requests.NewInteger(d.Get("min_bandwidth_abs").(int))
		update = true
	}
	if d.HasChange("max_bandwidth_abs") {
		request.MaxBandwidthAbs = requests.NewInteger(d.Get("max_bandwidth_abs").(int))
		update = true
	}
	if d.HasChange("min_bandwidth_percent") {
		request.MinBandwidthPercent = requests.NewInteger(d.Get("min_bandwidth_percent").(int))
		update = true
	}
	if d.HasChange("max_bandwidth_percent") {
		request.MaxBandwidthPercent = requests.NewInteger(d.Get("max_bandwidth_percent").(int))
		update = true
	}
	if d.HasChange("percent_source_type") {
		request.PercentSourceType = d.Get("percent_source_type").(string)
		update = true
	}

	if update {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.ModifyQosCar(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagQosCarRead(d, meta)
}

func resourceAlicloudSagQosCarDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteQosCarRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.QosId = parts[0]
	request.QosCarId = parts[1]

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DeleteQosCar(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceInOperating"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		if IsExpectedErrors(err, []string{"ParameterSagQosCarId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(sagService.WaitForSagQosCar(d.Id(), Deleted, DefaultTimeoutMedium))
}
