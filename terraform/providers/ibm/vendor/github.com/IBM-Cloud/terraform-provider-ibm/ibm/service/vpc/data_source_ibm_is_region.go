// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isRegionEndpoint = "endpoint"
	isRegionName     = "name"
	isRegionStatus   = "status"
)

func DataSourceIBMISRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISRegionRead,

		Schema: map[string]*schema.Schema{

			isRegionEndpoint: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isRegionName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			isRegionStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMISRegionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	if name == "" {
		bmxSess, err := meta.(conns.ClientSession).BluemixSession()
		if err != nil {
			tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_region", "read", "session-initialize-client")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		name = bmxSess.Config.Region
	}
	return regionGet(context, d, meta, name)
}

func regionGet(context context.Context, d *schema.ResourceData, meta interface{}, name string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_region", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getRegionOptions := &vpcv1.GetRegionOptions{
		Name: &name,
	}
	region, _, err := sess.GetRegionWithContext(context, getRegionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRegionWithContext failed: %s", err.Error()), "(Data) ibm_is_region", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	// For lack of anything better, compose our id from region name.
	d.SetId(*region.Name)
	if err = d.Set("endpoint", region.Endpoint); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting endpoint: %s", err), "(Data) ibm_is_region", "read", "set-endpoint").GetDiag()
	}
	if err = d.Set("name", region.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_region", "read", "set-name").GetDiag()
	}
	if err = d.Set("status", region.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_region", "read", "set-status").GetDiag()
	}
	return nil
}
