package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDRDSInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDRDSInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				Deprecated:   "Field 'name_regex' is deprecated and will be removed in a future release. Please use 'description_regex' instead.",
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudDRDSInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := drds.CreateDescribeDrdsInstancesRequest()
	request.RegionId = client.RegionId
	var dbi []drds.Instance
	var regexString *regexp.Regexp
	nameRegex, nameRegexGot := d.GetOk("name_regex")
	descriptionRegex, descriptionRegexGot := d.GetOk("description_regex")
	if nameRegexGot {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			regexString = r
		}
	} else if descriptionRegexGot {
		if r, err := regexp.Compile(descriptionRegex.(string)); err == nil {
			regexString = r
		}
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_drds_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	for _, item := range response.Instances.Instance {
		if regexString != nil {
			if !regexString.MatchString(item.Description) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[item.DrdsInstanceId]; !ok {
				continue
			}
		}

		dbi = append(dbi, item)
	}
	return drdsInstancesDescription(d, dbi)
}
func drdsInstancesDescription(d *schema.ResourceData, dbi []drds.Instance) error {
	var ids []string
	var descriptions []string
	var s []map[string]interface{}
	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":           item.DrdsInstanceId,
			"description":  item.Description,
			"type":         item.Type,
			"create_time":  item.CreateTime,
			"status":       item.Status,
			"network_type": item.NetworkType,
			"zone_id":      item.ZoneId,
			"version":      item.Version,
		}
		ids = append(ids, item.DrdsInstanceId)
		descriptions = append(descriptions, item.Description)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
