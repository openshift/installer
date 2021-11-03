package alicloud

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	cs "github.com/alibabacloud-go/cs-20151215/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const resourceName = "resource_alicloud_cs_autoscaling_config"

func resourceAlicloudCSAutoscalingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSAutoscalingConfigCreate,
		Read:   resourceAlicloudCSAutoscalingConfigRead,
		Update: resourceAlicloudCSAutoscalingConfigUpdate,
		Delete: resourceAlicloudCSAutoscalingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cool_down_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"unneeded_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"gpu_utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"scan_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "30s",
			},
		},
	}
}

func resourceAlicloudCSAutoscalingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudCSAutoscalingConfigUpdate(d, meta)
}

func resourceAlicloudCSAutoscalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudCSAutoscalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// cluster id
	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	// auto scaling config
	updateAutoscalingConfigRequest := &cs.CreateAutoscalingConfigRequest{}
	if v, ok := d.GetOk("cool_down_duration"); ok {
		updateAutoscalingConfigRequest.CoolDownDuration = tea.String(v.(string))
	}
	if v, ok := d.GetOk("unneeded_duration"); ok {
		updateAutoscalingConfigRequest.UnneededDuration = tea.String(v.(string))
	}
	if v, ok := d.GetOk("utilization_threshold"); ok {
		updateAutoscalingConfigRequest.UtilizationThreshold = tea.String(v.(string))
	}
	if v, ok := d.GetOk("gpu_utilization_threshold"); ok {
		updateAutoscalingConfigRequest.GpuUtilizationThreshold = tea.String(v.(string))
	}
	if v, ok := d.GetOk("scan_interval"); ok {
		updateAutoscalingConfigRequest.ScanInterval = tea.String(v.(string))
	}

	_, err = client.CreateAutoscalingConfig(tea.String(clusterId), updateAutoscalingConfigRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "CreateAutoscalingConfig", AliyunTablestoreGoSdk)
	}

	addDebug("CreateAutoscalingConfig", updateAutoscalingConfigRequest, err)
	d.SetId(clusterId)
	d.Partial(false)

	return resourceAlicloudCSAutoscalingConfigRead(d, meta)
}

func resourceAlicloudCSAutoscalingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
