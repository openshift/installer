package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCassandraDataCenters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCassandraDataCentersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"centers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commodity_instance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_center_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_center_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pay_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCassandraDataCentersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cassandra.CreateDescribeDataCentersRequest()
	request.RegionId = client.RegionId
	request.ClusterId = d.Get("cluster_id").(string)
	var objects []cassandra.DataCenter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeDataCenters(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cassandra_data_centers", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cassandra.DescribeDataCentersResponse)

	for _, item := range response.DataCenters.DataCenter {
		if nameRegex != nil {
			if !nameRegex.MatchString(item.DataCenterName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%v:%v", item.DataCenterId, item.ClusterId)]; !ok {
				continue
			}
		}
		if item.ClusterId == item.CommodityInstance {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"cluster_id":         object.ClusterId,
			"commodity_instance": object.CommodityInstance,
			"created_time":       object.CreatedTime,
			"data_center_id":     object.DataCenterId,
			"data_center_name":   object.DataCenterName,
			"disk_size":          object.DiskSize,
			"disk_type":          object.DiskType,
			"expire_time":        object.ExpireTime,
			"instance_type":      object.InstanceType,
			"lock_mode":          object.LockMode,
			"node_count":         object.NodeCount,
			"pay_type":           object.PayType,
			"status":             object.Status,
			"vpc_id":             object.VpcId,
			"vswitch_id":         object.VswitchId,
			"zone_id":            object.ZoneId,
		}
		ids[i] = fmt.Sprintf("%v:%v", object.DataCenterId, object.ClusterId)
		names[i] = object.DataCenterName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("centers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
