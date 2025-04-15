// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISPPPlacementGroup() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPISPPPlacementGroupRead,
		Schema: map[string]*schema.Schema{
			Arg_SPPPlacementGroupID: {
				Type:     schema.TypeString,
				Required: true,
			},

			Attr_SPPPlacementGroupName: {
				Type:     schema.TypeString,
				Computed: true,
			},

			Attr_SPPPlacementGroupPolicy: {
				Type:     schema.TypeString,
				Computed: true,
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			Attr_SPPPlacementGroupMembers: {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPISPPPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	placementGroupID := d.Get(Arg_SPPPlacementGroupID).(string)
	client := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(placementGroupID)
	if err != nil || response == nil {
		return diag.Errorf("error fetching the spp placement group: %v", err)
	}

	d.SetId(*response.ID)
	d.Set(Attr_SPPPlacementGroupName, response.Name)
	d.Set(Attr_SPPPlacementGroupPolicy, response.Policy)
	d.Set(Attr_SPPPlacementGroupMembers, response.MemberSharedProcessorPools)

	return nil
}
