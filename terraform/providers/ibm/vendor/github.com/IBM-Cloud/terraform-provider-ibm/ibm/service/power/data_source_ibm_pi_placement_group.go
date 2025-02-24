// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIPlacementGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIPlacementGroupRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_PlacementGroupName: {
				Description:  "The name of the placement group.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attribute
			Attr_Members: {
				Computed:    true,
				Description: "List of server instances IDs that are members of the placement group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_Policy: {
				Computed:    true,
				Description: "The value of the group's affinity policy. Valid values are affinity and anti-affinity.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	placementGroupName := d.Get(Arg_PlacementGroupName).(string)
	client := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(placementGroupName)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(*response.ID)
	d.Set(Attr_Members, response.Members)
	d.Set(Attr_Policy, response.Policy)

	return nil
}
