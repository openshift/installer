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

func DataSourceIBMPIDisasterRecoveryLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDisasterRecoveryLocation,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Location: {
				Computed:    true,
				Description: "The region zone of a site.",
				Type:        schema.TypeString,
			},
			Attr_ReplicationSites: {
				Computed:    true,
				Description: "List of replication sites.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_IsActive: {
							Computed:    true,
							Description: "Indicates the location is active or not, true if location is active , otherwise it is false.",
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
	}
}

func dataSourceIBMPIDisasterRecoveryLocation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	drClient := instance.NewIBMPIDisasterRecoveryLocationClient(ctx, sess, cloudInstanceID)
	drLocationSite, err := drClient.Get()
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(drLocationSite.ReplicationSites))
	for _, i := range drLocationSite.ReplicationSites {
		if i != nil {
			l := map[string]interface{}{
				Attr_IsActive: i.IsActive,
				Attr_Location: i.Location,
			}
			result = append(result, l)
		}
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_Location, drLocationSite.Location)
	d.Set(Attr_ReplicationSites, result)

	return nil
}
