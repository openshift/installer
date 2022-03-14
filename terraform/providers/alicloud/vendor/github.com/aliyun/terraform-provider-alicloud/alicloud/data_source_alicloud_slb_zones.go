package alicloud

import (
	"strings"

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
				ValidateFunc: validation.StringInSlice([]string{"vpc", "classic_intranet", "classic_internet", "Vpc"}, false),
			},
			"available_slb_address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"master_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slave_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "The parameter enable_details has been deprecated from version v1.154.0+",
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
						"master_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_slave_zone_ids": {
							Type:       schema.TypeList,
							Computed:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Deprecated: "the attribute slb_slave_zone_ids has been deprecated from version 1.157.0 and use slave_zone_id instead.",
						},
						"supported_resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"address_ip_version": {
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

func dataSourceAlicloudSlbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slb.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	if ipVersion, ok := d.GetOk("available_slb_address_ip_version"); ok {
		request.AddressIPVersion = ipVersion.(string)
	}
	if addressType, ok := d.GetOk("available_slb_address_type"); ok {
		request.AddressType = strings.ToLower(addressType.(string))
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeAvailableResourceResponse)
	var ids []string
	var s []map[string]interface{}
	for _, r := range response.AvailableResources.AvailableResource {
		if v, ok := d.GetOk("master_zone_id"); ok && v.(string) != r.MasterZoneId {
			continue
		}
		if v, ok := d.GetOk("slave_zone_id"); ok && v.(string) != r.SlaveZoneId {
			continue
		}
		ids = append(ids, r.MasterZoneId)
		mapping := map[string]interface{}{
			"id":                 r.MasterZoneId,
			"master_zone_id":     r.MasterZoneId,
			"slave_zone_id":      r.SlaveZoneId,
			"slb_slave_zone_ids": []string{r.SlaveZoneId},
		}
		supportedResourceList := make([]map[string]interface{}, 0)
		for _, v := range r.SupportResources.SupportResource {
			supportedResourceList = append(supportedResourceList, map[string]interface{}{
				"address_type":       v.AddressType,
				"address_ip_version": v.AddressIPVersion,
			})
		}
		mapping["supported_resources"] = supportedResourceList
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
