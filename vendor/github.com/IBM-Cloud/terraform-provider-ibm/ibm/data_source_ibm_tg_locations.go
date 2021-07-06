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
	tgLocations           = "locations"
	tgLocationsName       = "name"
	tgLocationsBillingLoc = "billing_location"
	tgLocationsType       = "type"
)

func dataSourceIBMTransitGatewaysLocations() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewaysLocationsRead,
		Schema: map[string]*schema.Schema{
			tgLocations: {
				Type:        schema.TypeList,
				Description: "Collection of Transit Gateway locations",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgLocationsName: {
							Type:        schema.TypeString,
							Computed:    true,
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
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewaysLocationsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	listTransitGatewayLocationsOptionsModel := &transitgatewayapisv1.ListGatewayLocationsOptions{}
	listTransitGatewayLocations, response, err := client.ListGatewayLocations(listTransitGatewayLocationsOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while fetching transit gateways locations: %s\n%s", err, response)
	}

	tgLocationsCol := make([]map[string]interface{}, 0)
	for _, instance := range listTransitGatewayLocations.Locations {

		transitgatewayLoc := map[string]interface{}{}
		transitgatewayLoc[tgLocationsName] = instance.Name
		transitgatewayLoc[tgLocationsType] = instance.Type
		transitgatewayLoc[tgLocationsBillingLoc] = instance.BillingLocation

		tgLocationsCol = append(tgLocationsCol, transitgatewayLoc)
	}
	d.Set(tgLocations, tgLocationsCol)
	d.SetId(dataSourceIBMTransitGatewaysID(d))
	return nil
}

// dataSourceIBMTransitGatewaysLocationsID returns a reasonable ID for a transit gateways locations list.
func dataSourceIBMTransitGatewaysLocationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
