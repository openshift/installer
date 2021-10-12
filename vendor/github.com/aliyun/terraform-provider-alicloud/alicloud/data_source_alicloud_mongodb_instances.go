package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMongoDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongoDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"sharding", "replicate"}, false),
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
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
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication": {
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
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mongos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"shards": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongoDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if v, ok := d.GetOk("instance_type"); ok {
		request.DBInstanceType = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var reqTags []dds.DescribeDBInstancesTag
		for key, value := range v.(map[string]interface{}) {
			reqTags = append(reqTags, dds.DescribeDBInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var instClass string
	if v, ok := d.GetOk("instance_class"); ok {
		instClass = strings.ToLower(v.(string))
	}

	var az string
	if v, ok := d.GetOk("availability_zone"); ok {
		az = strings.ToLower(v.(string))
	}

	var dbi []dds.DBInstance
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mongodb_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*dds.DescribeDBInstancesResponse)
		if len(response.DBInstances.DBInstance) < 1 {
			break
		}

		for _, item := range response.DBInstances.DBInstance {
			if nameRegex != nil && !nameRegex.MatchString(item.DBInstanceDescription) {
				continue
			}
			if len(instClass) > 0 && instClass != strings.ToLower(string(item.DBInstanceClass)) {
				continue
			}
			if len(az) > 0 && az != strings.ToLower(string(item.ZoneId)) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DBInstanceId]; !ok {
					continue
				}
			}
			dbi = append(dbi, item)
		}

		if len(response.DBInstances.DBInstance) < PageSizeLarge {
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
			"id":                item.DBInstanceId,
			"name":              item.DBInstanceDescription,
			"charge_type":       item.ChargeType,
			"instance_type":     item.DBInstanceType,
			"region_id":         item.RegionId,
			"creation_time":     item.CreationTime,
			"expiration_time":   item.ExpireTime,
			"status":            item.DBInstanceStatus,
			"engine":            item.Engine,
			"engine_version":    item.EngineVersion,
			"network_type":      item.NetworkType,
			"lock_mode":         item.LockMode,
			"availability_zone": item.ZoneId,
			"instance_class":    item.DBInstanceClass,
			"storage":           item.DBInstanceStorage,
			"replication":       item.ReplicationFactor,
			"tags":              ddsService.tagsToMap(item.Tags.Tag),
		}
		mongoList := []map[string]interface{}{}
		for _, v := range item.MongosList.MongosAttribute {
			mongo := map[string]interface{}{
				"description": v.NodeDescription,
				"node_id":     v.NodeId,
				"class":       v.NodeClass,
			}
			mongoList = append(mongoList, mongo)
		}
		shardList := []map[string]interface{}{}
		for _, v := range item.ShardList.ShardAttribute {
			shard := map[string]interface{}{
				"description": v.NodeDescription,
				"node_id":     v.NodeId,
				"class":       v.NodeClass,
				"storage":     v.NodeStorage,
			}
			shardList = append(shardList, shard)
		}
		mapping["mongos"] = mongoList
		mapping["shards"] = shardList
		ids = append(ids, item.DBInstanceId)
		names = append(names, item.DBInstanceDescription)
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
