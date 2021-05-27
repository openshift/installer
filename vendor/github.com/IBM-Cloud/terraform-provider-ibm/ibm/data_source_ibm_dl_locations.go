// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlBillingLocation         = "billing_location"
	dlLocDisplayName          = "display_name"
	dlLocationType            = "location_type"
	dlMarket                  = "market"
	dlMarketGeography         = "market_geography"
	dlMzr                     = "mzr"
	dlLocShortName            = "name"
	dlBuildingColocationOwner = "building_colocation_owner"
	dlVpcRegion               = "vpc_region"
	dlLocations               = "locations"
	dlMacsec                  = "macsec_enabled"
	dlProvisionEnabled        = "provision_enabled"
)

func dataSourceIBMDLLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLOfferingLocationsRead,
		Schema: map[string]*schema.Schema{
			dlOfferingType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"dedicated", "connect"}),
				Description:  "The Direct Link offering type. Current supported values (dedicated and connect).",
			},
			dlLocations: {
				Type:        schema.TypeList,
				Description: "Collection of valid locations for the specified Direct Link offering.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlOfferingType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Direct Link offering type. Current supported values (dedicated and connect).",
						},
						dlBillingLocation: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing location. Only present for locations where provisioning is enabled.",
						},
						dlBuildingColocationOwner: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Building colocation owner. Only present for offering_type=dedicated locations where provisioning is enabled.",
						},
						dlLocationType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location type",
						},
						dlLocShortName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location short name",
						}, dlLocDisplayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location long name",
						},
						dlMarket: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location market",
						},
						dlMarketGeography: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location geography. Only present for locations where provisioning is enabled.",
						},
						dlMzr: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is location a multi-zone region (MZR). Only present for locations where provisioning is enabled.",
						},
						dlVpcRegion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Location's VPC region. Only present for locations where provisioning is enabled.",
						},
						dlMacsec: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether location supports MACsec.",
						},
						dlProvisionEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates for the specific offering_type whether this location supports gateway provisioning.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDLOfferingLocationsRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := meta.(ClientSession).DirectlinkV1API()
	if err != nil {
		return err
	}
	listOfferingTypeLocationsOptions := &directlinkv1.ListOfferingTypeLocationsOptions{}
	listOfferingTypeLocationsOptions.SetOfferingType(d.Get(dlOfferingType).(string))
	listLocations, response, err := directLink.ListOfferingTypeLocations(listOfferingTypeLocationsOptions)
	if err != nil {
		return fmt.Errorf("Error while listing directlink gateway's locations %s\n%s", err, response)
	}

	locations := make([]map[string]interface{}, 0)
	for _, instance := range listLocations.Locations {
		location := map[string]interface{}{}
		if instance.BuildingColocationOwner != nil {
			location[dlBuildingColocationOwner] = *instance.BuildingColocationOwner
		}

		if instance.DisplayName != nil {
			location[dlLocDisplayName] = *instance.DisplayName
		}
		if instance.Name != nil {
			location[dlLocShortName] = *instance.Name
		}
		if instance.LocationType != nil {
			location[dlLocationType] = *instance.LocationType
		}
		if instance.OfferingType != nil {
			location[dlOfferingType] = *instance.OfferingType
		}
		if instance.Market != nil {
			location[dlMarket] = *instance.Market
		}

		if instance.MarketGeography != nil {
			location[dlMarketGeography] = *instance.MarketGeography
		}
		if instance.Mzr != nil {
			location[dlMzr] = *instance.Mzr
		}
		if instance.VpcRegion != nil {
			location[dlVpcRegion] = *instance.VpcRegion
		}
		if instance.BillingLocation != nil {
			location[dlBillingLocation] = *instance.BillingLocation
		}
		if instance.MacsecEnabled != nil {
			location[dlMacsec] = *instance.MacsecEnabled
		}
		if instance.ProvisionEnabled != nil {
			location[dlProvisionEnabled] = *instance.ProvisionEnabled
		}
		locations = append(locations, location)
	}

	d.SetId(dataSourceIBMDLLocationsID(d))
	d.Set(dlLocations, locations)
	return nil
}

// dataSourceIBMDLLocationsID returns a reasonable ID for a direct link offering locations list.
func dataSourceIBMDLLocationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
