package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRouteEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRouteEntriesRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudRouteEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.RouteTableId = d.Get("route_table_id").(string)

	var allRouteEntries []vpc.RouteEntry
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			response, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			raw = response
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_route_entries", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) < 1 {
			break
		}

		for _, entries := range response.RouteTables.RouteTable[0].RouteEntrys.RouteEntry {
			if instance_id, ok := d.GetOk("instance_id"); ok && entries.InstanceId != instance_id.(string) {
				continue
			}
			if route_entry_type, ok := d.GetOk("type"); ok && entries.Type != route_entry_type.(string) {
				continue
			}
			if cidr_block, ok := d.GetOk("cidr_block"); ok && entries.DestinationCidrBlock != cidr_block.(string) {
				continue
			}
			allRouteEntries = append(allRouteEntries, entries)
		}

		if len(response.RouteTables.RouteTable) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return RouteEntriesDecriptionAttributes(d, allRouteEntries, meta)
}

func RouteEntriesDecriptionAttributes(d *schema.ResourceData, entries []vpc.RouteEntry, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, entry := range entries {
		mapping := map[string]interface{}{
			"route_table_id": entry.RouteTableId,
			"instance_id":    entry.InstanceId,
			"status":         entry.Status,
			"next_hop_type":  entry.NextHopType,
			"type":           entry.Type,
			"cidr_block":     entry.DestinationCidrBlock,
		}
		ids = append(ids, entry.RouteTableId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
