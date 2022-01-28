package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEssScalingConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssScalingConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_in": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"credit_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"system_disk_performance_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disks": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"category": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"device": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"delete_with_instance": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							Computed: true,
						},
						"lifecycle_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_price_limit": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"price_limit": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}

	var allScalingConfigurations []ess.ScalingConfiguration

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_scalingconfigurations", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeScalingConfigurationsResponse)
		if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		allScalingConfigurations = append(allScalingConfigurations, response.ScalingConfigurations.ScalingConfiguration...)
		if len(response.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredScalingConfigurations = make([]ess.ScalingConfiguration, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, configuration := range allScalingConfigurations {
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(configuration.ScalingConfigurationName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[configuration.ScalingConfigurationId]; !ok {
					continue
				}
			}
			filteredScalingConfigurations = append(filteredScalingConfigurations, configuration)
		}
	} else {
		filteredScalingConfigurations = allScalingConfigurations
	}

	return scalingConfigurationsDescriptionAttribute(d, filteredScalingConfigurations, meta)
}

func scalingConfigurationsDescriptionAttribute(d *schema.ResourceData, scalingConfigurations []ess.ScalingConfiguration, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	for _, scalingConfiguration := range scalingConfigurations {
		mapping := map[string]interface{}{
			"id":                            scalingConfiguration.ScalingConfigurationId,
			"name":                          scalingConfiguration.ScalingConfigurationName,
			"scaling_group_id":              scalingConfiguration.ScalingGroupId,
			"image_id":                      scalingConfiguration.ImageId,
			"instance_type":                 scalingConfiguration.InstanceType,
			"security_group_id":             scalingConfiguration.SecurityGroupId,
			"internet_charge_type":          scalingConfiguration.InternetChargeType,
			"internet_max_bandwidth_in":     scalingConfiguration.InternetMaxBandwidthIn,
			"internet_max_bandwidth_out":    scalingConfiguration.InternetMaxBandwidthOut,
			"credit_specification":          scalingConfiguration.CreditSpecification,
			"system_disk_category":          scalingConfiguration.SystemDiskCategory,
			"system_disk_size":              scalingConfiguration.SystemDiskSize,
			"system_disk_performance_level": scalingConfiguration.SystemDiskPerformanceLevel,
			"data_disks":                    essService.flattenDataDiskMappings(scalingConfiguration.DataDisks.DataDisk),
			"lifecycle_state":               scalingConfiguration.LifecycleState,
			"creation_time":                 scalingConfiguration.CreationTime,
			"instance_name":                 scalingConfiguration.InstanceName,
			"host_name":                     scalingConfiguration.HostName,
			"spot_strategy":                 scalingConfiguration.SpotStrategy,
			"spot_price_limit":              essService.flattenSpotPriceLimitMappings(scalingConfiguration.SpotPriceLimit.SpotPriceModel),
		}
		ids = append(ids, scalingConfiguration.ScalingConfigurationId)
		names = append(names, scalingConfiguration.ScalingConfigurationName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
