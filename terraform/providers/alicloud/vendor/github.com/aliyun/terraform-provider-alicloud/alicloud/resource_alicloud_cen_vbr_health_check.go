package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenVbrHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenVbrHealthCheckCreate,
		Read:   resourceAlicloudCenVbrHealthCheckRead,
		Update: resourceAlicloudCenVbrHealthCheckUpdate,
		Delete: resourceAlicloudCenVbrHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"health_check_source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_target_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},
			"vbr_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vbr_instance_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vbr_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenVbrHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateEnableCenVbrHealthCheckRequest()
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("health_check_interval"); ok {
		request.HealthCheckInterval = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("health_check_source_ip"); ok {
		request.HealthCheckSourceIp = v.(string)
	}

	request.HealthCheckTargetIp = d.Get("health_check_target_ip").(string)
	if v, ok := d.GetOk("healthy_threshold"); ok {
		request.HealthyThreshold = requests.NewInteger(v.(int))
	}

	request.VbrInstanceId = d.Get("vbr_instance_id").(string)
	if v, ok := d.GetOk("vbr_instance_owner_id"); ok {
		request.VbrInstanceOwnerId = requests.NewInteger(v.(int))
	}

	request.VbrInstanceRegionId = d.Get("vbr_instance_region_id").(string)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.EnableCenVbrHealthCheck(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InstanceStatus.NotSupport", "Operation.Blocking", "InvalidOperation.CenInstanceStatus", "InvalidOperation.VbrNotAttachedToCen", "OperationFailed.StatusNotSupport"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(fmt.Sprintf("%v:%v", request.VbrInstanceId, request.VbrInstanceRegionId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_vbr_health_check", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudCenVbrHealthCheckRead(d, meta)
}
func resourceAlicloudCenVbrHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenVbrHealthCheck(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_vbr_health_check cbnService.DescribeCenVbrHealthCheck Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vbr_instance_id", parts[0])
	d.Set("vbr_instance_region_id", parts[1])
	d.Set("cen_id", object.CenId)
	d.Set("health_check_interval", object.HealthCheckInterval)
	d.Set("health_check_source_ip", object.HealthCheckSourceIp)
	d.Set("health_check_target_ip", object.HealthCheckTargetIp)
	d.Set("healthy_threshold", object.HealthyThreshold)
	return nil
}
func resourceAlicloudCenVbrHealthCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := cbn.CreateEnableCenVbrHealthCheckRequest()
	request.VbrInstanceId = parts[0]
	request.VbrInstanceRegionId = parts[1]
	request.CenId = d.Get("cen_id").(string)
	if d.HasChange("health_check_target_ip") {
		update = true
	}
	request.HealthCheckTargetIp = d.Get("health_check_target_ip").(string)
	if d.HasChange("health_check_interval") {
		update = true
		request.HealthCheckInterval = requests.NewInteger(d.Get("health_check_interval").(int))
	}
	if d.HasChange("health_check_source_ip") {
		update = true
		request.HealthCheckSourceIp = d.Get("health_check_source_ip").(string)
	}
	if d.HasChange("healthy_threshold") {
		update = true
		request.HealthyThreshold = requests.NewInteger(d.Get("healthy_threshold").(int))
	}
	if d.HasChange("vbr_instance_owner_id") {
		update = true
		request.VbrInstanceOwnerId = requests.NewInteger(d.Get("vbr_instance_owner_id").(int))
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.EnableCenVbrHealthCheck(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.CenInstanceStatus", "InstanceStatus.NotSupport"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCenVbrHealthCheckRead(d, meta)
}
func resourceAlicloudCenVbrHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := cbn.CreateDisableCenVbrHealthCheckRequest()
	request.VbrInstanceId = parts[0]
	request.VbrInstanceRegionId = parts[1]
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("vbr_instance_owner_id"); ok {
		request.VbrInstanceOwnerId = requests.NewInteger(v.(int))
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DisableCenVbrHealthCheck(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InstanceStatus.NotSupport"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
