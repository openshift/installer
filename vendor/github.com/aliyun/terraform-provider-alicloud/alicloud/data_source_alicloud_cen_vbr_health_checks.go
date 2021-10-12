package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCenVbrHealthChecks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenVbrHealthChecksRead,
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vbr_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vbr_instance_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"vbr_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_target_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vbr_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vbr_instance_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenVbrHealthChecksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenVbrHealthCheckRequest()
	if v, ok := d.GetOk("cen_id"); ok {
		request.CenId = v.(string)
	}
	if v, ok := d.GetOk("vbr_instance_id"); ok {
		request.VbrInstanceId = v.(string)
	}
	if v, ok := d.GetOk("vbr_instance_owner_id"); ok {
		request.VbrInstanceOwnerId = requests.NewInteger(v.(int))
	}
	request.VbrInstanceRegionId = d.Get("vbr_instance_region_id").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.VbrHealthCheck
	var response *cbn.DescribeCenVbrHealthCheckResponse
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenVbrHealthCheck(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_vbr_health_checks", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*cbn.DescribeCenVbrHealthCheckResponse)

		for _, item := range response.VbrHealthChecks.VbrHealthCheck {
			objects = append(objects, item)
		}
		if len(response.VbrHealthChecks.VbrHealthCheck) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                     fmt.Sprintf("%v:%v", object.VbrInstanceId, object.VbrInstanceRegionId),
			"cen_id":                 object.CenId,
			"health_check_interval":  object.HealthCheckInterval,
			"health_check_source_ip": object.HealthCheckSourceIp,
			"health_check_target_ip": object.HealthCheckTargetIp,
			"healthy_threshold":      object.HealthyThreshold,
			"vbr_instance_id":        object.VbrInstanceId,
			"vbr_instance_region_id": object.VbrInstanceRegionId,
		}
		ids = append(ids, fmt.Sprintf("%v:%v", object.VbrInstanceId, object.VbrInstanceRegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("checks", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
