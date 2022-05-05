// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isRegions    = "regions"
	isRegionHref = "href"
)

func DataSourceIBMISRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISRegionsRead,
		Schema: map[string]*schema.Schema{

			isRegions: {
				Type:        schema.TypeList,
				Description: "List of regions",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						isRegionEndpoint: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isRegionHref: {
							Type:     schema.TypeString,
							Computed: true,
						},

						isRegionName: {
							Type:     schema.TypeString,
							Computed: true,
						},

						isRegionStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISRegionsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listRegionOptions := &vpcv1.ListRegionsOptions{}
	regioncollection, _, err := sess.ListRegions(listRegionOptions)
	if err != nil {
		return err
	}
	regions := regioncollection.Regions
	regionInfo := make([]map[string]interface{}, 0)
	for _, region := range regions {
		l := map[string]interface{}{
			isRegionEndpoint: *region.Endpoint,
			isRegionName:     *region.Name,
			isRegionStatus:   *region.Status,
			isRegionHref:     *region.Href,
		}
		regionInfo = append(regionInfo, l)
	}
	d.SetId(dataSourceIBMISRegionsID(d))
	d.Set(isRegions, regionInfo)
	return nil
}

// dataSourceIBMISRegionsID returns a reasonable ID for the region list.
func dataSourceIBMISRegionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
