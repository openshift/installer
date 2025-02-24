// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISPPPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISPPPlacementGroupsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_SPPPlacementGroups: {
				Computed:    true,
				Description: "List of all the shared processor pool placement groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Members: {
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of shared processor pool IDs that are members of the shared processor pool placement group.",
							Type:        schema.TypeList,
						},
						Attr_Name: {
							Computed:    true,
							Description: "User defined name for the shared processor pool placement group.",
							Type:        schema.TypeString,
						},
						Attr_Policy: {
							Computed:    true,
							Description: "The value of the group's affinity policy. Valid values are affinity and anti-affinity.",
							Type:        schema.TypeString,
						},
						Attr_SPPPlacementGroupID: {
							Computed:    true,
							Description: "The ID of the shared processor pool placement group.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPISPPPlacementGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)
	groups, err := client.GetAll()
	if err != nil || groups == nil {
		return diag.Errorf("error fetching spp placement groups: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(groups.SppPlacementGroups))
	for _, placementGroup := range groups.SppPlacementGroups {
		key := map[string]interface{}{
			Attr_Members:             placementGroup.MemberSharedProcessorPools,
			Attr_Name:                placementGroup.Name,
			Attr_Policy:              placementGroup.Policy,
			Attr_SPPPlacementGroupID: placementGroup.ID,
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_SPPPlacementGroups, result)

	return nil
}
