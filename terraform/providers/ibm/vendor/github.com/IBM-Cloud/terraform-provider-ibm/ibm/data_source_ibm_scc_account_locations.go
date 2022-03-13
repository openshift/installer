// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/adminserviceapiv1"
)

func dataSourceIbmSccAccountLocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSccAccountLocationsRead,

		Schema: map[string]*schema.Schema{
			"locations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The programatic ID of the location that you want to work in.",
						},
						"main_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The base URL for the service.",
						},
						"governance_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint that is used to call the Configuration Governance APIs.",
						},
						"results_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint that is used to get the results for the Configuration Governance component.",
						},
						"compliance_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint that is used to call the Posture Management APIs.",
						},
						"analytics_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint that is used to generate analytics for the Posture Management component.",
						},
						"si_endpoint_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint that is used to call the Security Insights APIs.",
						},
						"regions": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The programatic ID of the available regions.",
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

func dataSourceIbmSccAccountLocationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceApiClient, err := meta.(ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	listLocationsOptions := &adminserviceapiv1.ListLocationsOptions{}

	locations, response, err := adminServiceApiClient.ListLocationsWithContext(context, listLocationsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLocationsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLocationsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccAccountLocationsID(d))

	if locations.Locations != nil {
		err = d.Set("locations", dataSourceLocationsFlattenLocations(locations.Locations))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting locations %s", err))
		}
	}

	return nil
}

// dataSourceIbmSccAccountLocationsID returns a reasonable ID for the list.
func dataSourceIbmSccAccountLocationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLocationsFlattenLocations(result []adminserviceapiv1.Location) (locations []map[string]interface{}) {
	for _, locationsItem := range result {
		locations = append(locations, dataSourceLocationsLocationsToMap(locationsItem))
	}

	return locations
}

func dataSourceLocationsLocationsToMap(locationsItem adminserviceapiv1.Location) (locationsMap map[string]interface{}) {
	locationsMap = map[string]interface{}{}

	if locationsItem.ID != nil {
		locationsMap["id"] = locationsItem.ID
	}
	if locationsItem.MainEndpointURL != nil {
		locationsMap["main_endpoint_url"] = locationsItem.MainEndpointURL
	}
	if locationsItem.GovernanceEndpointURL != nil {
		locationsMap["governance_endpoint_url"] = locationsItem.GovernanceEndpointURL
	}
	if locationsItem.ResultsEndpointURL != nil {
		locationsMap["results_endpoint_url"] = locationsItem.ResultsEndpointURL
	}
	if locationsItem.ComplianceEndpointURL != nil {
		locationsMap["compliance_endpoint_url"] = locationsItem.ComplianceEndpointURL
	}
	if locationsItem.AnalyticsEndpointURL != nil {
		locationsMap["analytics_endpoint_url"] = locationsItem.AnalyticsEndpointURL
	}
	if locationsItem.SiEndpointURL != nil {
		locationsMap["si_endpoint_url"] = locationsItem.SiEndpointURL
	}
	if locationsItem.Regions != nil {
		regionsList := []map[string]interface{}{}
		for _, regionsItem := range locationsItem.Regions {
			regionsList = append(regionsList, dataSourceLocationsLocationsRegionsToMap(regionsItem))
		}
		locationsMap["regions"] = regionsList
	}

	return locationsMap
}

func dataSourceLocationsLocationsRegionsToMap(regionsItem adminserviceapiv1.Region) (regionsMap map[string]interface{}) {
	regionsMap = map[string]interface{}{}

	if regionsItem.ID != nil {
		regionsMap["id"] = regionsItem.ID
	}

	return regionsMap
}
