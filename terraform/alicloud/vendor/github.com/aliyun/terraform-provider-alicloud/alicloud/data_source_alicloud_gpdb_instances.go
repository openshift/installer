package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudGpdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGpdbInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"names": {
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_group_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGpdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	// name regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		} else {
			return WrapError(err)
		}
	}
	// availability zone
	availabilityZone := d.Get("availability_zone").(string)
	// vSwitchId
	vSwitchId := d.Get("vswitch_id").(string)
	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	request := gpdb.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	if v, ok := d.GetOk("tags"); ok {
		var reqTags []gpdb.DescribeDBInstancesTag
		for key, value := range v.(map[string]interface{}) {
			reqTags = append(reqTags, gpdb.DescribeDBInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	var dbi []gpdb.DBInstanceAttribute
	for {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_gpdb_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*gpdb.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)
		if len(response.Items.DBInstance) < 1 {
			break
		}

		for _, item := range response.Items.DBInstance {
			// filter by description regex
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}
			// filter by instance id
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DBInstanceId]; !ok {
					continue
				}
			}
			// filter by availability zone
			if availabilityZone != "" && availabilityZone != strings.ToLower(string(item.ZoneId)) {
				continue
			}
			// filter by vSwitchId
			if vSwitchId != "" && vSwitchId != string(item.VSwitchId) {
				continue
			}

			// describe instance
			instanceAttribute, err := gpdbService.DescribeGpdbInstance(item.DBInstanceId)
			if err != nil {
				return WrapError(err)
			}
			dbi = append(dbi, instanceAttribute)
		}

		if len(response.Items.DBInstance) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return describeGpdbInstances(d, dbi)
}

func describeGpdbInstances(d *schema.ResourceData, dbi []gpdb.DBInstanceAttribute) error {
	var ids []string
	var names []string
	var instances []map[string]interface{}
	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                    item.DBInstanceId,
			"description":           item.DBInstanceDescription,
			"region_id":             item.RegionId,
			"availability_zone":     item.ZoneId,
			"creation_time":         item.CreationTime,
			"status":                item.DBInstanceStatus,
			"engine":                item.Engine,
			"engine_version":        item.EngineVersion,
			"charge_type":           item.PayType,
			"instance_class":        item.DBInstanceClass,
			"instance_group_count":  item.DBInstanceGroupCount,
			"instance_network_type": item.InstanceNetworkType,
		}
		ids = append(ids, item.DBInstanceId)
		names = append(names, item.DBInstanceDescription)
		instances = append(instances, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", instances); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), instances)
	}
	return nil
}
