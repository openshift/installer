// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func DataSourceIBMCbrZoneAddresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCbrZoneAddressesRead,

		Schema: map[string]*schema.Schema{
			"zone_addresses_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of a zone addresses resource.",
			},
			"zone_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the zone.",
			},
			"addresses": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of addresses added to the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of address.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"ref": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A service reference value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the account owning the service.",
									},
									"service_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service type.",
									},
									"service_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name.",
									},
									"service_instance": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service instance.",
									},
									"location": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The location.",
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

func dataSourceIBMCbrZoneAddressesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	zoneAddressesId := d.Get("zone_addresses_id").(string)
	zoneId, addressesId := decomposeZoneAddressesId(zoneAddressesId)

	if zoneId == "" || addressesId == "" {
		return diag.Errorf("zone_addresses_id %s not found", zoneAddressesId)
	}

	var zone *contextbasedrestrictionsv1.Zone
	var found bool
	zone, _, found, err = getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}
	if !found {
		return diag.Errorf("zone_addresses_id %s not found", zoneAddressesId)
	}

	var addresses []map[string]interface{}
	addresses, err = dataSourceDecodeAddressList(zone.Addresses, addressesId)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(addresses) == 0 {
		return diag.Errorf("zone_addresses_id %s not found", zoneAddressesId)
	}

	d.SetId(zoneAddressesId)

	if err = d.Set("zone_id", zoneId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone_id: %s", err))
	}

	if err = d.Set("addresses", addresses); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting addresses %s", err))
	}

	return nil
}
