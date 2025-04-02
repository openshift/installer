// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_zone_addresses", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	zoneAddressesId := d.Get("zone_addresses_id").(string)
	zoneId, addressesId := decomposeZoneAddressesId(zoneAddressesId)

	if zoneId == "" || addressesId == "" {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("zone_addresses_id %s not found", zoneAddressesId), "(Data) ibm_cbr_zone_addresses", "read", "zone_addresses_id-not-found").GetDiag()
	}

	var zone *contextbasedrestrictionsv1.Zone
	var found bool
	zone, _, found, err = getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("%s", err.Error()), "(Data) ibm_cbr_zone_addresses", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if !found {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("zone not found by zone_id %s", zoneId), "(Data) ibm_cbr_zone_addresses", "read", "zone-not-found-from-getZone").GetDiag()
	}

	var addresses []map[string]interface{}
	addresses, err = DataSourceDecodeAddressList(zone.Addresses, addressesId)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_zone_addresses", "read", "DataSourceDecodeAddressList").GetDiag()
	}
	if len(addresses) == 0 {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("zone_addresses_id %s not found", zoneAddressesId), "(Data) ibm_cbr_zone_addresses", "read", "zone_addresses_id-not-found").GetDiag()
	}

	d.SetId(zoneAddressesId)

	if err = d.Set("zone_id", zoneId); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting zone_id: %s", err), "(Data) ibm_cbr_zone_addresses", "read", "set-zone_id").GetDiag()
	}

	if err = d.Set("addresses", addresses); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting addresses %s", err), "(Data) ibm_cbr_zone_addresses", "read", "set-addresses").GetDiag()
	}

	return nil
}
