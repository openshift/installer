// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIPlacementGroupsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_PlacementGroups: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ID: {
							Computed:    true,
							Description: "The ID of the placement group.",
							Type:        schema.TypeString,
						},
						Attr_Members: {
							Computed:    true,
							Description: "List of server instances IDs that are members of the placement group.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
						Attr_Name: {
							Computed:    true,
							Description: "User defined name for the placement group.",
							Type:        schema.TypeString,
						},
						Attr_Policy: {
							Computed:    true,
							Description: "The value of the group's affinity policy. Valid values are affinity and anti-affinity.",
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIPlacementGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)
	groups, err := client.GetAll()
	if err != nil {
		log.Printf("[ERROR] get all placement groups failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(groups.PlacementGroups))
	for _, placementGroup := range groups.PlacementGroups {
		key := map[string]interface{}{
			Attr_ID:      placementGroup.ID,
			Attr_Members: placementGroup.Members,
			Attr_Name:    placementGroup.Name,
			Attr_Policy:  placementGroup.Policy,
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_PlacementGroups, result)

	return nil
}
