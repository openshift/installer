package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenRouteServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenRouteServicesRead,
		Schema: map[string]*schema.Schema{
			"access_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenRouteServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeRouteServicesInCenRequest()
	if v, ok := d.GetOk("access_region_id"); ok {
		request.AccessRegionId = v.(string)
	}
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("host"); ok {
		request.Host = v.(string)
	}
	if v, ok := d.GetOk("host_region_id"); ok {
		request.HostRegionId = v.(string)
	}
	if v, ok := d.GetOk("host_vpc_id"); ok {
		request.HostVpcId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.RouteServiceEntry
	status, statusOk := d.GetOk("status")
	var response *cbn.DescribeRouteServicesInCenResponse
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeRouteServicesInCen(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_route_services", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*cbn.DescribeRouteServicesInCenResponse)

		for _, item := range response.RouteServiceEntries.RouteServiceEntry {
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.RouteServiceEntries.RouteServiceEntry) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprintf("%v:%v:%v:%v", object.CenId, object.HostRegionId, object.Host, object.AccessRegionId),
			"access_region_id": object.AccessRegionId,
			"cen_id":           object.CenId,
			"cidrs":            object.Cidrs.Cidr,
			"description":      object.Description,
			"host":             object.Host,
			"host_region_id":   object.HostRegionId,
			"host_vpc_id":      object.HostVpcId,
			"status":           object.Status,
			"update_interval":  object.UpdateInterval,
		}
		ids = append(ids, fmt.Sprintf("%v:%v:%v:%v", object.CenId, object.HostRegionId, object.Host, object.AccessRegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("services", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
