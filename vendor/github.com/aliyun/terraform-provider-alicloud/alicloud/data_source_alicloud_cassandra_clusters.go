package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCassandraClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCassandraClustersRead,
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
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_center_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"major_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"minor_version": {
							Type:     schema.TypeString,
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
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCassandraClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}

	request := cassandra.CreateDescribeClustersRequest()
	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = v.(string)
	}
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cassandra.Cluster
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
	if v, ok := d.GetOk("tags"); ok {
		var reqTags []cassandra.DescribeClustersTag
		for key, value := range v.(map[string]interface{}) {
			reqTags = append(reqTags, cassandra.DescribeClustersTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}
	for {
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.DescribeClusters(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cassandra_clusters", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cassandra.DescribeClustersResponse)

		for _, item := range response.Clusters.Cluster {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.ClusterName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.ClusterId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Clusters.Cluster) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"id":                object.ClusterId,
			"cluster_id":        object.ClusterId,
			"cluster_name":      object.ClusterName,
			"created_time":      object.CreatedTime,
			"data_center_count": object.DataCenterCount,
			"expire_time":       object.ExpireTime,
			"lock_mode":         object.LockMode,
			"major_version":     object.MajorVersion,
			"minor_version":     object.MinorVersion,
			"pay_type":          object.PayType,
			"status":            object.Status,
			"tags":              cassandraService.tagsToMap(object.Tags.Tag),
		}
		ids[i] = object.ClusterId
		names[i] = object.ClusterName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
