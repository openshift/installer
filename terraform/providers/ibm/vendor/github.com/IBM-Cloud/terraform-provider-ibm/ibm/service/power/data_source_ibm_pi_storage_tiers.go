// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
)

func DataSourceIBMPIStorageTiers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStorageTiersRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_RegionStorageTiers: {
				Computed:    true,
				Description: "An array of of storage tiers supported in a region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Computed:    true,
							Description: "Description of the storage tier label.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the storage tier.",
							Type:        schema.TypeString,
						},
						Attr_State: {
							Computed:    true,
							Description: "State of the storage tier (active or inactive).",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIStorageTiersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPIStorageTierClient(ctx, sess, cloudInstanceID)
	rst, err := client.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)

	regionStorageTiers := []map[string]interface{}{}
	if len(rst) > 0 {
		for _, storageTier := range rst {
			regionStorageTier := storageTierToMap(storageTier)
			regionStorageTiers = append(regionStorageTiers, regionStorageTier)
		}
	}
	d.Set(Attr_RegionStorageTiers, regionStorageTiers)

	return nil
}

func storageTierToMap(storageTier *models.StorageTier) map[string]interface{} {
	storageTierMap := make(map[string]interface{})
	if storageTier.Description != "" {
		storageTierMap[Attr_Description] = storageTier.Description
	}
	if storageTier.Name != "" {
		storageTierMap[Attr_Name] = storageTier.Name
	}
	if storageTier.State != nil {
		storageTierMap[Attr_State] = storageTier.State
	}
	return storageTierMap
}
