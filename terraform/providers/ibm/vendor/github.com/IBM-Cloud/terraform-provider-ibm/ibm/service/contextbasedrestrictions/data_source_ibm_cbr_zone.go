// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMCbrZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCbrZoneRead,

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of a zone.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone CRN.",
			},
			"address_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of addresses in the zone.",
			},
			"excluded_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of excluded addresses in the zone.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the zone.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the account owning this zone.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the zone.",
			},
			"addresses": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of addresses in the zone.",
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
			"excluded": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of excluded addresses in the zone. Only addresses of type `ipAddress`, `ipRange`, and `subnet` can be excluded.",
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
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href link to the resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the resource was created.",
			},
			"created_by_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the resource.",
			},
			"last_modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the resource was modified.",
			},
			"last_modified_by_id": &schema.Schema{
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

	d.SetId(fmt.Sprintf("%s", *getZoneOptions.ZoneID))

	if err = d.Set("crn", zone.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("address_count", flex.IntValue(zone.AddressCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address_count: %s", err))
	}

	if err = d.Set("excluded_count", flex.IntValue(zone.ExcludedCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded_count: %s", err))
	}

	if err = d.Set("name", zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("account_id", zone.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}

	if err = d.Set("description", zone.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	addresses := []map[string]interface{}{}
	if zone.Addresses != nil {
		for _, modelItem := range zone.Addresses {
			modelMap, err := dataSourceIBMCbrZoneAddressToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			addresses = append(addresses, modelMap)
		}
	}
	if err = d.Set("addresses", addresses); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting addresses %s", err))
	}

	excluded := []map[string]interface{}{}
	if zone.Excluded != nil {
		for _, modelItem := range zone.Excluded {
			modelMap, err := dataSourceIBMCbrZoneAddressToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			excluded = append(excluded, modelMap)
		}
	}
	if err = d.Set("excluded", excluded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded %s", err))
	}

	if err = d.Set("href", zone.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(zone.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}

	if err = d.Set("last_modified_at", flex.DateTimeToString(zone.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}

	if err = d.Set("last_modified_by_id", zone.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}

	return nil
}

func dataSourceIBMCbrZoneAddressToMap(model contextbasedrestrictionsv1.AddressIntf) (map[string]interface{}, error) {
	if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddress); ok {
		return dataSourceIBMCbrZoneAddressIPAddressToMap(model.(*contextbasedrestrictionsv1.AddressIPAddress))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddressRange); ok {
		return dataSourceIBMCbrZoneAddressIPAddressRangeToMap(model.(*contextbasedrestrictionsv1.AddressIPAddressRange))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressSubnet); ok {
		return dataSourceIBMCbrZoneAddressSubnetToMap(model.(*contextbasedrestrictionsv1.AddressSubnet))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressVPC); ok {
		return dataSourceIBMCbrZoneAddressVPCToMap(model.(*contextbasedrestrictionsv1.AddressVPC))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressServiceRef); ok {
		return dataSourceIBMCbrZoneAddressServiceRefToMap(model.(*contextbasedrestrictionsv1.AddressServiceRef))
	} else if _, ok := model.(*contextbasedrestrictionsv1.Address); ok {
		modelMap := make(map[string]interface{})
		model := model.(*contextbasedrestrictionsv1.Address)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Value != nil {
			modelMap["value"] = *model.Value
		}
		if model.Ref != nil {
			refMap, err := dataSourceIBMCbrZoneServiceRefValueToMap(model.Ref)
			if err != nil {
				return modelMap, err
			}
			modelMap["ref"] = []map[string]interface{}{refMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized contextbasedrestrictionsv1.AddressIntf subtype encountered")
	}
}

func dataSourceIBMCbrZoneServiceRefValueToMap(model *contextbasedrestrictionsv1.ServiceRefValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccountID != nil {
		modelMap["account_id"] = *model.AccountID
	}
	if model.ServiceType != nil {
		modelMap["service_type"] = *model.ServiceType
	}
	if model.ServiceName != nil {
		modelMap["service_name"] = *model.ServiceName
	}
	if model.ServiceInstance != nil {
		modelMap["service_instance"] = *model.ServiceInstance
	}
	if model.Location != nil {
		modelMap["location"] = *model.Location
	}
	return modelMap, nil
}

func dataSourceIBMCbrZoneAddressIPAddressToMap(model *contextbasedrestrictionsv1.AddressIPAddress) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func dataSourceIBMCbrZoneAddressServiceRefToMap(model *contextbasedrestrictionsv1.AddressServiceRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Ref != nil {
		refMap, err := dataSourceIBMCbrZoneServiceRefValueToMap(model.Ref)
		if err != nil {
			return modelMap, err
		}
		modelMap["ref"] = []map[string]interface{}{refMap}
	}
	return modelMap, nil
}

func dataSourceIBMCbrZoneAddressSubnetToMap(model *contextbasedrestrictionsv1.AddressSubnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func dataSourceIBMCbrZoneAddressIPAddressRangeToMap(model *contextbasedrestrictionsv1.AddressIPAddressRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func dataSourceIBMCbrZoneAddressVPCToMap(model *contextbasedrestrictionsv1.AddressVPC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}
