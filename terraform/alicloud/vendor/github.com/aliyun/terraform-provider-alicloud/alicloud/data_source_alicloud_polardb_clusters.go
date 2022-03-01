package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBClustersRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clusters": {
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
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
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
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
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
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_lock": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_node_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_iops": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_node_role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_connections": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_node_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"create_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPolarDBClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := polardb.CreateDescribeDBClustersRequest()

	request.RegionId = client.RegionId
	request.DBClusterStatus = d.Get("status").(string)
	request.DBType = d.Get("db_type").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var dbi []polardb.DBCluster

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
	}

	// ids
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
		var reqTags []polardb.DescribeDBClustersTag
		for key, value := range v.(map[string]interface{}) {
			reqTags = append(reqTags, polardb.DescribeDBClustersTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}
	for {
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusters(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_clusters", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*polardb.DescribeDBClustersResponse)
		if len(response.Items.DBCluster) < 1 {
			break
		}

		for _, item := range response.Items.DBCluster {

			if descriptionRegex != nil {
				if !descriptionRegex.MatchString(item.DBClusterDescription) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DBClusterId]; !ok {
					continue
				}
			}

			dbi = append(dbi, item)
		}

		if len(response.Items.DBCluster) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	return polarDBClustersDescription(d, dbi)
}

func polarDBClustersDescription(d *schema.ResourceData, dbi []polardb.DBCluster) error {
	var ids []string
	var descriptions []string
	var s []map[string]interface{}

	for _, item := range dbi {
		var nodes []map[string]interface{}
		for _, node := range item.DBNodes.DBNode {
			nodeMap := map[string]interface{}{
				"db_node_class":   node.DBNodeClass,
				"max_iops":        node.MaxIOPS,
				"db_node_role":    node.DBNodeRole,
				"region_id":       node.RegionId,
				"max_connections": node.MaxConnections,
				"zone_id":         node.ZoneId,
				"db_node_status":  node.DBNodeStatus,
				"db_node_id":      node.DBNodeId,
				"create_time":     node.CreationTime,
			}
			nodes = append(nodes, nodeMap)
		}
		mapping := map[string]interface{}{
			"id":             item.DBClusterId,
			"description":    item.DBClusterDescription,
			"charge_type":    getChargeType(item.PayType),
			"network_type":   item.DBClusterNetworkType,
			"region_id":      item.RegionId,
			"zone_id":        item.ZoneId,
			"expire_time":    item.ExpireTime,
			"expired":        item.Expired,
			"status":         item.DBClusterStatus,
			"engine":         item.Engine,
			"db_type":        item.DBType,
			"db_version":     item.DBVersion,
			"lock_mode":      item.LockMode,
			"delete_lock":    item.DeletionLock,
			"create_time":    item.CreateTime,
			"vpc_id":         item.VpcId,
			"db_node_number": item.DBNodeNumber,
			"db_node_class":  item.DBNodeClass,
			"storage_used":   item.StorageUsed,
			"db_nodes":       nodes,
		}

		ids = append(ids, item.DBClusterId)
		descriptions = append(descriptions, item.DBClusterDescription)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
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
