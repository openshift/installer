package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudHBaseInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHBaseInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
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
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"core_node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"core_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pay_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHBaseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	request := hbase.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var az string
	if v, ok := d.GetOk("availability_zone"); ok {
		az = strings.ToLower(v.(string))
	}

	var dbi []hbase.Instance
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		var reqTags []hbase.DescribeInstancesTag
		for key, value := range v.(map[string]interface{}) {
			reqTags = append(reqTags, hbase.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	for {
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DescribeInstances(request)
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbase_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*hbase.DescribeInstancesResponse)
		if len(response.Instances.Instance) < 1 {
			break
		}

		for _, item := range response.Instances.Instance {
			if nameRegex != nil && !nameRegex.MatchString(item.InstanceName) {
				continue
			}
			if len(az) > 0 && az != strings.ToLower(item.ZoneId) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.InstanceId]; !ok {
					continue
				}
			}
			dbi = append(dbi, item)
		}

		if len(response.Instances.Instance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range dbi {
		mapping := map[string]interface{}{
			"id":                   item.InstanceId,
			"name":                 item.InstanceName,
			"region_id":            item.RegionId,
			"zone_id":              item.ZoneId,
			"engine":               item.Engine,
			"engine_version":       item.MajorVersion,
			"network_type":         item.NetworkType,
			"master_instance_type": item.MasterInstanceType,
			"master_node_count":    item.MasterNodeCount,
			"core_instance_type":   item.CoreInstanceType,
			"core_node_count":      item.CoreNodeCount,
			"core_disk_type":       item.CoreDiskType,
			"core_disk_size":       item.CoreDiskSize,
			"vpc_id":               item.VpcId,
			"vswitch_id":           item.VswitchId,
			"pay_type":             item.PayType,
			"status":               item.Status,
			"backup_status":        item.BackupStatus,
			"created_time":         item.CreatedTimeUTC,
			"expire_time":          item.ExpireTimeUTC,
			"deletion_protection":  item.IsDeletionProtection,
			"tags":                 hbaseService.tagsToMap(item.Tags.Tag),
		}
		ids = append(ids, item.InstanceId)
		names = append(names, item.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("instances", s); err != nil {
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
		err := writeToFile(output.(string), s)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
