package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBEndpointsRead,

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_add_new_nodes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_write_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_items": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"net_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connection_string": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
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
									"ip_address": {
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

func dataSourceAlicloudPolarDBEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := polardb.CreateDescribeDBClusterEndpointsRequest()

	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	if v, ok := d.GetOk("db_endpoint_id"); ok && v.(string) != "" {
		request.DBEndpointId = v.(string)
	}

	var dbi []polardb.DBEndpoint

	raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.DescribeDBClusterEndpoints(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_endpoints", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)

	for _, item := range response.Items {
		dbi = append(dbi, item)
	}

	return polarDBClusterEndpointsDescription(d, dbi)
}

func polarDBClusterEndpointsDescription(d *schema.ResourceData, dbi []polardb.DBEndpoint) error {
	var ids []string
	var s []map[string]interface{}
	for _, item := range dbi {
		var addrs []map[string]interface{}
		for _, addr := range item.AddressItems {
			addrMap := map[string]interface{}{
				"net_type":          addr.NetType,
				"connection_string": addr.ConnectionString,
				"port":              addr.Port,
				"vpc_id":            addr.VPCId,
				"vswitch_id":        addr.VSwitchId,
				"ip_address":        addr.IPAddress,
			}
			addrs = append(addrs, addrMap)
		}
		mapping := map[string]interface{}{
			"db_endpoint_id":     item.DBEndpointId,
			"auto_add_new_nodes": item.AutoAddNewNodes,
			"endpoint_config":    item.EndpointConfig,
			"endpoint_type":      item.EndpointType,
			"nodes":              item.Nodes,
			"read_write_mode":    item.ReadWriteMode,
			"address_items":      addrs,
		}
		ids = append(ids, item.DBEndpointId)
		s = append(s, mapping)
	}
	if err := d.Set("endpoints", s); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
