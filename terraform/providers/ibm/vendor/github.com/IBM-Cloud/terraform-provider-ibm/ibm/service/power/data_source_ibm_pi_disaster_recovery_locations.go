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
)

func DataSourceIBMPIDisasterRecoveryLocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDisasterRecoveryLocations,
		Schema: map[string]*schema.Schema{
			// Attributes
			Attr_DisasterRecoveryLocations: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Location: {
							Computed:    true,
							Description: "The region zone of a site.",
							Type:        schema.TypeString,
						},
						Attr_ReplicationSites: {
							Computed:    true,
							Description: "List of Replication Sites.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_IsActive: {
										Computed:    true,
										Description: "Indicates the location is active or not, true if location is active, otherwise it is false.",
										Type:        schema.TypeBool,
									},
									Attr_Location: {
										Computed:    true,
										Description: "The region zone of the location.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
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
						Attr_IsActive: j.IsActive,
						Attr_Location: j.Location,
					}
					replicationSites = append(replicationSites, r)
				}
			}
			l := map[string]interface{}{
				Attr_Location:         i.Location,
				Attr_ReplicationSites: replicationSites,
			}
			results = append(results, l)
		}
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_DisasterRecoveryLocations, results)

	return nil
}
