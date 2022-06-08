// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	PIPlacementGroups = "placement_groups"
)

func DataSourceIBMPIPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIPlacementGroupsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "PI cloud instance ID",
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PIPlacementGroups: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						PIPlacementGroupMembers: {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
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

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)
	groups, err := client.GetAll()
	if err != nil {
		log.Printf("[ERROR] get all placement groups failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(groups.PlacementGroups))
	for _, placementGroup := range groups.PlacementGroups {
		key := map[string]interface{}{
			"id":                    placementGroup.ID,
			"name":                  placementGroup.Name,
			PIPlacementGroupMembers: placementGroup.Members,
			"policy":                placementGroup.Policy,
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PIPlacementGroups, result)

	return nil
}
