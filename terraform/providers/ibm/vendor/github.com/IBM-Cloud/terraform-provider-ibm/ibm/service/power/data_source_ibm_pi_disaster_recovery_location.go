// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIDisasterRecoveryLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDisasterRecoveryLocation,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			PIDRLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RegionZone of a site",
			},
			"replication_sites": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						PIDRLocation: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIDisasterRecoveryLocation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	drClient := instance.NewIBMPIDisasterRecoveryLocationClient(ctx, sess, cloudInstanceID)
	drLocationSite, err := drClient.Get()
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(drLocationSite.ReplicationSites))
	for _, i := range drLocationSite.ReplicationSites {
		if i != nil {
			l := map[string]interface{}{
				"is_active":  i.IsActive,
				PIDRLocation: i.Location,
			}
			result = append(result, l)
		}
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(PIDRLocation, drLocationSite.Location)
	d.Set("replication_sites", result)

	return nil
}
