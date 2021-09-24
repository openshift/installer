package alicloud

import (
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSlbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbZonesRead,

		Schema: map[string]*schema.Schema{
			"available_slb_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Vpc", "classic_intranet", "classic_internet"}, false),
			},
			"available_slb_address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_slave_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slaveZones := make(map[string][]string)

	request := slb.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	if ipVersion, ok := d.GetOk("available_slb_address_ip_version"); ok {
		request.AddressIPVersion = ipVersion.(string)
	}
	if addressType, ok := d.GetOk("available_slb_address_type"); ok {
		request.AddressType = addressType.(string)
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeAvailableResourceResponse)
	for _, resource := range response.AvailableResources.AvailableResource {
		slaveIds := slaveZones[resource.MasterZoneId]
		slaveIds = append(slaveIds, resource.SlaveZoneId)
		if len(slaveIds) > 0 {
			sort.Strings(slaveIds)
		}
		slaveZones[resource.MasterZoneId] = slaveIds
	}

	var ids []string
	for v, _ := range slaveZones {
		ids = append(ids, v)
	}
	if len(ids) > 0 {
		sort.Strings(ids)
	}

	var s []map[string]interface{}
	for _, zoneId := range ids {
		mapping := map[string]interface{}{"id": zoneId}
		if len(slaveZones) > 0 {
			mapping["slb_slave_zone_ids"] = slaveZones[zoneId]
		}
		if !d.Get("enable_details").(bool) {
			s = append(s, mapping)
			continue
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
