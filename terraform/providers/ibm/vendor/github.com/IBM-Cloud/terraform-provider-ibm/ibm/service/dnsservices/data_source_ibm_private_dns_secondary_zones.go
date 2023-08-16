// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	pdnsSecondaryZones = "secondary_zones"
	pdnsSZResolverID   = "resolver_id"
	pdnsSZId           = "secondary_zone_id"
	pdnsSZDescription  = "description"
	pdnsSZZone         = "zone"
	pdnsSZEnabled      = "enabled"
	pdnsSZTransferFrom = "transfer_from"
	pdnsSZCreatedOn    = "created_on"
	pdnsSZModifiedOn   = "modified_on"
	pdnsSZOffset       = "offset"
	pdnsSZLimit        = "limit"
)

func DataSourceIBMPrivateDNSSecondaryZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMDNSSecondaryZonesRead,

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the DNS Services instance.",
			},
			pdnsSZResolverID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a custom resolver.",
			},
			pdnsSecondaryZones: {
				Type:        schema.TypeList,
				Description: "List of Secondary Zones",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsSZId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the Secondary Zone",
						},
						pdnsSZDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the secondary zone.",
						},
						pdnsSZZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the zone.",
						},
						pdnsSZEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable/Disable the secondary zone.",
						},
						pdnsSZTransferFrom: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The addresses of DNS servers where the secondary zone data is transferred from.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						pdnsSZCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when a secondary zone is created",
						},

						pdnsSZModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The recent time when a secondary zone is modified",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDNSSecondaryZonesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsSZResolverID).(string)

	opt := sess.NewListSecondaryZonesOptions(instanceID, resolverID)

	result, resp, err := sess.ListSecondaryZonesWithContext(context, opt)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error listing the Secondary Zones %s:%s", err, resp))
	}

	secondaryZones := make([]interface{}, 0)
	for _, instance := range result.SecondaryZones {
		secondaryZone := map[string]interface{}{}
		secondaryZone[pdnsSZId] = *instance.ID
		secondaryZone[pdnsSZDescription] = *instance.Description
		secondaryZone[pdnsSZZone] = *instance.Zone
		secondaryZone[pdnsSZEnabled] = *instance.Enabled
		secondaryZone[pdnsSZTransferFrom] = instance.TransferFrom
		secondaryZone[pdnsSZCreatedOn] = (*instance.CreatedOn).String()
		secondaryZone[pdnsSZModifiedOn] = (*instance.ModifiedOn).String()

		secondaryZones = append(secondaryZones, secondaryZone)
	}
	d.SetId(dataSourceIBMPrivateDNSSecondaryZoneID(d))
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsSZResolverID, resolverID)
	d.Set(pdnsSecondaryZones, secondaryZones)
	return nil
}

func dataSourceIBMPrivateDNSSecondaryZoneID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
