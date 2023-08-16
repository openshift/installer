// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIDisasterRecoveryLocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDisasterRecoveryLocations,
		Schema: map[string]*schema.Schema{

			// Computed Attributes
			"disaster_recovery_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceIBMPIDisasterRecoveryLocations(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	drClient := instance.NewIBMPIDisasterRecoveryLocationClient(ctx, sess, "")
	drLocationSites, err := drClient.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	results := make([]map[string]interface{}, 0, len(drLocationSites.DisasterRecoveryLocations))
	for _, i := range drLocationSites.DisasterRecoveryLocations {
		if i != nil {
			replicationSites := make([]map[string]interface{}, 0, len(i.ReplicationSites))
			for _, j := range i.ReplicationSites {
				if j != nil {
					r := map[string]interface{}{
						"is_active":  j.IsActive,
						PIDRLocation: j.Location,
					}
					replicationSites = append(replicationSites, r)
				}
			}
			l := map[string]interface{}{
				"location":          i.Location,
				"replication_sites": replicationSites,
			}
			results = append(results, l)
		}
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("disaster_recovery_locations", results)

	return nil
}
