package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEdasClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEdasClustersRead,

		Schema: map[string]*schema.Schema{
			"logical_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEdasClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	logicalRegionId := d.Get("logical_region_id").(string)
	request := edas.CreateListClusterRequest()
	request.LogicalRegionId = logicalRegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			if id == nil {
				continue
			}
			idsMap[Trim(id.(string))] = Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_clusters", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListClusterResponse)
	if response.Code != 200 {
		return WrapError(Error(response.Message))
	}

	var filteredClusters []edas.Cluster
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
		}
		for _, cluster := range response.ClusterList.Cluster {
			if r != nil && !r.MatchString(cluster.ClusterName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[cluster.ClusterId]; !ok {
					continue
				}
			}
			filteredClusters = append(filteredClusters, cluster)
		}
	} else {
		filteredClusters = response.ClusterList.Cluster
	}

	return edasClusterDescriptionAttributes(d, filteredClusters)
}

func edasClusterDescriptionAttributes(d *schema.ResourceData, clusters []edas.Cluster) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, cluster := range clusters {
		mapping := map[string]interface{}{
			"cluster_id":   cluster.ClusterId,
			"cluster_name": cluster.ClusterName,
			"cluster_type": cluster.ClusterType,
			"create_time":  cluster.CreateTime,
			"update_time":  cluster.UpdateTime,
			"cpu":          cluster.Cpu,
			"cpu_used":     cluster.CpuUsed,
			"mem":          cluster.Mem,
			"mem_used":     cluster.MemUsed,
			"network_mode": cluster.NetworkMode,
			"node_num":     cluster.NodeNum,
			"vpc_id":       cluster.VpcId,
			"region_id":    cluster.RegionId,
		}
		ids = append(ids, cluster.ClusterId)
		s = append(s, mapping)
		names = append(names, cluster.ClusterName)
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
