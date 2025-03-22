// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isRegions    = "regions"
	isRegionHref = "href"
)

func DataSourceIBMISRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISRegionsRead,
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

func dataSourceIBMISRegionsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_regions", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	listRegionOptions := &vpcv1.ListRegionsOptions{}
	regioncollection, _, err := sess.ListRegionsWithContext(context, listRegionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListRegionsWithContext failed: %s", err.Error()), "(Data) ibm_is_regions", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
	if err = d.Set("regions", regionInfo); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting regions: %s", err), "(Data) ibm_is_regions", "read", "set-regions").GetDiag()
	}
	return nil
}

// dataSourceIBMISRegionsID returns a reasonable ID for the region list.
func dataSourceIBMISRegionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
