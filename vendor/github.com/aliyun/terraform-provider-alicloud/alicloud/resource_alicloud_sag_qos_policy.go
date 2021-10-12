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

func resourceAlicloudSagQosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSagQosPolicyCreate,
		Read:   resourceAlicloudSagQosPolicyRead,
		Update: resourceAlicloudSagQosPolicyUpdate,
		Delete: resourceAlicloudSagQosPolicyDelete,
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
			"ip_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudSagQosPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateQosPolicyRequest()

	request.QosId = d.Get("qos_id").(string)
	request.Priority = requests.NewInteger(d.Get("priority").(int))
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.SourceCidr = d.Get("source_cidr").(string)
	request.SourcePortRange = d.Get("source_port_range").(string)
	request.DestCidr = d.Get("dest_cidr").(string)
	request.DestPortRange = d.Get("dest_port_range").(string)

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("start_time"); ok && v.(string) != "" {
		request.StartTime = v.(string)
	}
	if v, ok := d.GetOk("end_time"); ok && v.(string) != "" {
		request.EndTime = v.(string)
	}

	var response *smartag.CreateQosPolicyResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.CreateQosPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceInOperating"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*smartag.CreateQosPolicyResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sag_qos_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", response.QosId, COLON_SEPARATED, response.QosPolicyId))

	return resourceAlicloudSagQosPolicyRead(d, meta)
}

func resourceAlicloudSagQosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	sagService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := sagService.DescribeSagQosPolicy(d.Id())
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
	d.Set("ip_protocol", object.IpProtocol)
	d.Set("source_cidr", object.SourceCidr)
	d.Set("source_port_range", object.SourcePortRange)
	d.Set("dest_cidr", object.DestCidr)
	d.Set("dest_port_range", object.DestPortRange)
	d.Set("start_time", object.StartTime)
	d.Set("end_time", object.EndTime)

	return nil
}

func resourceAlicloudSagQosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartag.CreateModifyQosPolicyRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.QosId = parts[0]
	request.QosPolicyId = parts[1]

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
	if d.HasChange("ip_protocol") {
		request.IpProtocol = d.Get("ip_protocol").(string)
		update = true
	}
	if d.HasChange("source_cidr") {
		request.SourceCidr = d.Get("source_cidr").(string)
		update = true
	}
	if d.HasChange("source_port_range") {
		request.SourcePortRange = d.Get("source_port_range").(string)
		update = true
	}
	if d.HasChange("dest_cidr") {
		request.DestCidr = d.Get("dest_cidr").(string)
		update = true
	}
	if d.HasChange("dest_port_range") {
		request.DestPortRange = d.Get("dest_port_range").(string)
		update = true
	}
	if d.HasChange("start_time") {
		request.StartTime = d.Get("start_time").(string)
		update = true
	}
	if d.HasChange("end_time") {
		request.EndTime = d.Get("end_time").(string)
		update = true
	}

	if update {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.ModifyQosPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSagQosPolicyRead(d, meta)
}

func resourceAlicloudSagQosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteQosPolicyRequest()
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request.QosId = parts[0]
	request.QosPolicyId = parts[1]

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DeleteQosPolicy(request)
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
		if IsExpectedErrors(err, []string{"ParameterSagQosPolicyId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(sagService.WaitForSagQosPolicy(d.Id(), Deleted, DefaultTimeoutMedium))
}
