// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_zone", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	zoneId := d.Get("zone_id").(string)

	var zone *contextbasedrestrictionsv1.Zone
	var found bool
	zone, _, found, err = getZone(contextBasedRestrictionsClient, context, zoneId)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("getZone failed: %s", err.Error()), "(Data) ibm_cbr_zone", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if !found {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("zone_id %s not found", zoneId), "(Data) ibm_cbr_zone", "read", "zone_id_not_found").GetDiag()
	}

	d.SetId(zoneId)

	if err = d.Set("crn", zone.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cbr_zone", "read", "set-crn").GetDiag()
	}

	if err = d.Set("address_count", flex.IntValue(zone.AddressCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting address_count: %s", err), "(Data) ibm_cbr_zone", "read", "set-address_count").GetDiag()
	}

	if err = d.Set("excluded_count", flex.IntValue(zone.ExcludedCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting excluded_count: %s", err), "(Data) ibm_cbr_zone", "read", "set-excluded_count").GetDiag()
	}

	if err = d.Set("name", zone.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cbr_zone", "read", "set-name").GetDiag()
	}

	if err = d.Set("account_id", zone.AccountID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting account_id: %s", err), "(Data) ibm_cbr_zone", "read", "set-account_id").GetDiag()
	}

	if err = d.Set("description", zone.Description); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_cbr_zone", "read", "set-description").GetDiag()
	}

	var addresses []map[string]interface{}
	addresses, err = DataSourceDecodeAddressList(zone.Addresses, cbrZoneAddressIdDefault)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_zone", "read", "DataSourceDecodeAddressList").GetDiag()
	}
	if err = d.Set("addresses", addresses); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting addresses: %s", err), "(Data) ibm_cbr_zone", "read", "set-addresses").GetDiag()
	}

	var excluded []map[string]interface{}
	excluded, err = DataSourceDecodeAddressList(zone.Excluded, cbrZoneAddressIdDefault)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_zone", "read", "DataSourceDecodeAddressList_excluded_address").GetDiag()
	}
	if err = d.Set("excluded", excluded); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting excluded: %s", err), "(Data) ibm_cbr_zone", "read", "set-excluded").GetDiag()
	}

	if err = d.Set("href", zone.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cbr_zone", "read", "set-href").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(zone.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_cbr_zone", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("created_by_id", zone.CreatedByID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_by_id: %s", err), "(Data) ibm_cbr_zone", "read", "set-created_by_id").GetDiag()
	}

	if err = d.Set("last_modified_at", flex.DateTimeToString(zone.LastModifiedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_modified_at: %s", err), "(Data) ibm_cbr_zone", "read", "set-last_modified_at").GetDiag()
	}

	if err = d.Set("last_modified_by_id", zone.LastModifiedByID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_modified_by_id: %s", err), "(Data) ibm_cbr_zone", "read", "set-last_modified_by_id").GetDiag()
	}

	return nil
}

func DataSourceIBMCbrZoneAddressToMap(model contextbasedrestrictionsv1.AddressIntf) (modelMap map[string]interface{}, addressId string, err error) {
	if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddress); ok {
		return DataSourceIBMCbrZoneAddressIPAddressToMap(model.(*contextbasedrestrictionsv1.AddressIPAddress))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressIPAddressRange); ok {
		return DataSourceIBMCbrZoneAddressIPAddressRangeToMap(model.(*contextbasedrestrictionsv1.AddressIPAddressRange))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressSubnet); ok {
		return DataSourceIBMCbrZoneAddressSubnetToMap(model.(*contextbasedrestrictionsv1.AddressSubnet))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressVPC); ok {
		return DataSourceIBMCbrZoneAddressVPCToMap(model.(*contextbasedrestrictionsv1.AddressVPC))
	} else if _, ok := model.(*contextbasedrestrictionsv1.AddressServiceRef); ok {
		return DataSourceIBMCbrZoneAddressServiceRefToMap(model.(*contextbasedrestrictionsv1.AddressServiceRef))
	} else if _, ok := model.(*contextbasedrestrictionsv1.Address); ok {
		modelMap = make(map[string]interface{})
		address := model.(*contextbasedrestrictionsv1.Address)
		if address.Type != nil {
			modelMap["type"] = *address.Type
		}
		if address.Value != nil {
			modelMap["value"] = *address.Value
		}
		if address.Ref != nil {
			var refMap map[string]interface{}
			refMap, err = DataSourceIBMCbrZoneServiceRefValueToMap(address.Ref)
			if err != nil {
				return
			}
			modelMap["ref"] = []map[string]interface{}{refMap}
		}
		if address.ID != nil {
			addressId = *address.ID
		}
	} else {
		err = fmt.Errorf("Unrecognized contextbasedrestrictionsv1.AddressIntf subtype encountered")
	}

	return
}

func DataSourceIBMCbrZoneServiceRefValueToMap(model *contextbasedrestrictionsv1.ServiceRefValue) (map[string]interface{}, error) {
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

func DataSourceIBMCbrZoneAddressIPAddressToMap(model *contextbasedrestrictionsv1.AddressIPAddress) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func DataSourceIBMCbrZoneAddressServiceRefToMap(model *contextbasedrestrictionsv1.AddressServiceRef) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Ref != nil {
		var refMap map[string]interface{}
		refMap, err = DataSourceIBMCbrZoneServiceRefValueToMap(model.Ref)
		if err != nil {
			return
		}
		modelMap["ref"] = []map[string]interface{}{refMap}
	}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func DataSourceIBMCbrZoneAddressSubnetToMap(model *contextbasedrestrictionsv1.AddressSubnet) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func DataSourceIBMCbrZoneAddressIPAddressRangeToMap(model *contextbasedrestrictionsv1.AddressIPAddressRange) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func DataSourceIBMCbrZoneAddressVPCToMap(model *contextbasedrestrictionsv1.AddressVPC) (modelMap map[string]interface{}, addressId string, err error) {
	modelMap = make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.ID != nil {
		addressId = *model.ID
	}
	return
}

func DataSourceDecodeAddressList(addresses []contextbasedrestrictionsv1.AddressIntf, wantAddressId string) (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0, len(addresses))
	for _, addr := range addresses {
		var m map[string]interface{}
		var addressId string
		m, addressId, err = DataSourceIBMCbrZoneAddressToMap(addr)
		if err != nil {
			return
		}
		if addressId == wantAddressId {
			result = append(result, m)
		}
	}
	return
}
