// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	tgLocalConnectionLocations = "local_connection_locations"
	tgLocationsDisplayName     = "display_name"
)

func dataSourceIBMTransitGatewaysLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMTransitGatewaysLocationRead,
		Schema: map[string]*schema.Schema{
			tgLocationsName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Location.",
			},
			tgLocationsType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the location, determining is this a multi-zone region, a single data center, or a point of presence.",
			},

			tgLocationsBillingLoc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The geographical location of this location, used for billing purposes.",
			},
			tgLocalConnectionLocations: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The set of network locations that are considered local for this Transit Gateway location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgLocationsName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Location.",
						},
						tgLocationsDisplayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A descriptive display name for the location.",
						},
						tgLocationsType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the location, determining is this a multi-zone region, a single data center, or a point of presence.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewaysLocationRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	detailGatewayLocationOptionsModel := &transitgatewayapisv1.GetGatewayLocationOptions{}
	locName := d.Get(tgLocationsName).(string)
	detailGatewayLocationOptionsModel.Name = &locName
	detailTransitGatewayLocation, response, err := client.GetGatewayLocation(detailGatewayLocationOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while fetching transit gateway detailed location: %s\n%s", err, response)
	}

	if detailTransitGatewayLocation != nil {
		d.SetId(dataSourceIBMTransitGatewaysLocationID(d))
		d.Set(tgLocationsType, *detailTransitGatewayLocation.Type)
		d.Set(tgLocationsBillingLoc, *detailTransitGatewayLocation.BillingLocation)
		tgConnLocationsCol := make([]map[string]interface{}, 0)
		for _, instance := range detailTransitGatewayLocation.LocalConnectionLocations {
			tgConnLocation := map[string]interface{}{}
			tgConnLocation[tgLocationsName] = *instance.Name
			tgConnLocation[tgLocationsDisplayName] = *instance.DisplayName
			tgConnLocation[tgLocationsType] = *instance.Type
			tgConnLocationsCol = append(tgConnLocationsCol, tgConnLocation)
		}
		d.Set(tgLocalConnectionLocations, tgConnLocationsCol)
	}
	return nil
}

// dataSourceIBMTransitGatewaysLocationID returns a reasonable ID for a transit gateways location list.
func dataSourceIBMTransitGatewaysLocationID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
