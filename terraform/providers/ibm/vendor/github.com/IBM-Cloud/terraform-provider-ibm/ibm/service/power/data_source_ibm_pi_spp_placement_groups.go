// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	PISPPPlacementGroups = "spp_placement_groups"
)

func DataSourceIBMPISPPPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISPPPlacementGroupsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "PI cloud instance ID",
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PISPPPlacementGroups: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_SPPPlacementGroupID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						Attr_SPPPlacementGroupName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						Attr_SPPPlacementGroupMembers: {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						Attr_SPPPlacementGroupPolicy: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISPPPlacementGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)
	groups, err := client.GetAll()
	if err != nil || groups == nil {
		return diag.Errorf("error fetching spp placement groups: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(groups.SppPlacementGroups))
	for _, placementGroup := range groups.SppPlacementGroups {
		key := map[string]interface{}{
			Attr_SPPPlacementGroupID:      placementGroup.ID,
			Attr_SPPPlacementGroupName:    placementGroup.Name,
			Attr_SPPPlacementGroupMembers: placementGroup.MemberSharedProcessorPools,
			Attr_SPPPlacementGroupPolicy:  placementGroup.Policy,
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PISPPPlacementGroups, result)

	return nil
}
