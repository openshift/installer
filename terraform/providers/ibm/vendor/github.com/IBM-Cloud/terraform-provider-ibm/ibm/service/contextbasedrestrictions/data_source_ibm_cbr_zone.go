// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func DataSourceIBMCbrZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCbrZoneRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of a zone.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique ID of the zone.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone CRN.",
			},
			"address_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of addresses in the zone.",
			},
			"excluded_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of excluded addresses in the zone.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the zone.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the account owning this zone.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the zone.",
			},
			"addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of addresses in the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of address.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"ref": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A service reference value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the account owning the service.",
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service type.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name.",
									},
									"service_instance": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service instance.",
									},
								},
							},
						},
					},
				},
			},
			"excluded": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of excluded addresses in the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of address.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"ref": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A service reference value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the account owning the service.",
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service type.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service name.",
									},
									"service_instance": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service instance.",
									},
								},
							},
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href link to the resource.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the resource was created.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the resource.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the resource was modified.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which modified the resource.",
			},
		},
	}
}

func dataSourceIBMCbrZoneRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getZoneOptions := &contextbasedrestrictionsv1.GetZoneOptions{}

	getZoneOptions.SetZoneID(d.Get("zone_id").(string))

	zone, response, err := contextBasedRestrictionsClient.GetZoneWithContext(context, getZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] GetZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId(*getZoneOptions.ZoneID)
	if err = d.Set("id", zone.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting id: %s", err))
	}
	if err = d.Set("crn", zone.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("address_count", flex.IntValue(zone.AddressCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting address_count: %s", err))
	}
	if err = d.Set("excluded_count", flex.IntValue(zone.ExcludedCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting excluded_count: %s", err))
	}
	if err = d.Set("name", zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("account_id", zone.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}
	if err = d.Set("description", zone.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}

	if zone.Addresses != nil {
		err = d.Set("addresses", dataSourceZoneFlattenAddresses(zone.Addresses))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting addresses %s", err))
		}
	}

	if zone.Excluded != nil {
		err = d.Set("excluded", dataSourceZoneFlattenExcluded(zone.Excluded))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting excluded %s", err))
		}
	}
	if err = d.Set("href", zone.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(zone.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", flex.DateTimeToString(zone.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_modified_at: %s", err))
	}
	if err = d.Set("last_modified_by_id", zone.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_modified_by_id: %s", err))
	}

	return nil
}

func dataSourceZoneFlattenAddresses(result []contextbasedrestrictionsv1.AddressIntf) (addresses []map[string]interface{}) {
	for _, addressesItem := range result {
		addresses = append(addresses, dataSourceZoneAddressesToMap(addressesItem))
	}

	return addresses
}

func dataSourceZoneAddressesToMap(addressesItem contextbasedrestrictionsv1.AddressIntf) (addressesMap map[string]interface{}) {

	buf, err := json.Marshal(addressesItem)

	if err == nil {
		err = json.Unmarshal(buf, &addressesMap)
	}

	if err != nil {
		panic(err)
	}

	return addressesMap
}

func dataSourceZoneAddressesRefToMap(refItem contextbasedrestrictionsv1.ServiceRefValue) (refMap map[string]interface{}) {
	refMap = map[string]interface{}{}

	if refItem.AccountID != nil {
		refMap["account_id"] = refItem.AccountID
	}
	if refItem.ServiceType != nil {
		refMap["service_type"] = refItem.ServiceType
	}
	if refItem.ServiceName != nil {
		refMap["service_name"] = refItem.ServiceName
	}
	if refItem.ServiceInstance != nil {
		refMap["service_instance"] = refItem.ServiceInstance
	}

	return refMap
}

func dataSourceZoneFlattenExcluded(result []contextbasedrestrictionsv1.AddressIntf) (excluded []map[string]interface{}) {
	for _, excludedItem := range result {
		excluded = append(excluded, dataSourceZoneExcludedToMap(excludedItem))
	}

	return excluded
}

func dataSourceZoneExcludedToMap(excludedItem contextbasedrestrictionsv1.AddressIntf) (excludedMap map[string]interface{}) {

	buf, err := json.Marshal(excludedItem)

	if err == nil {
		err = json.Unmarshal(buf, &excludedMap)
	}

	if err != nil {
		panic(err)
	}

	return excludedMap
}

func dataSourceZoneExcludedRefToMap(refItem contextbasedrestrictionsv1.ServiceRefValue) (refMap map[string]interface{}) {
	refMap = map[string]interface{}{}

	if refItem.AccountID != nil {
		refMap["account_id"] = refItem.AccountID
	}
	if refItem.ServiceType != nil {
		refMap["service_type"] = refItem.ServiceType
	}
	if refItem.ServiceName != nil {
		refMap["service_name"] = refItem.ServiceName
	}
	if refItem.ServiceInstance != nil {
		refMap["service_instance"] = refItem.ServiceInstance
	}

	return refMap
}
