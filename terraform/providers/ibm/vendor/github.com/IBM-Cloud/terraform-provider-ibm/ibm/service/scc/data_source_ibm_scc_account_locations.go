// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
)

func DataSourceIBMSccAccountLocations() *schema.Resource {
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
	adminServiceApiClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
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

	locations_lst := []map[string]interface{}{}
	if locations.Locations != nil {
		for _, modelItem := range locations.Locations {
			modelMap, err := dataSourceIbmSccAccountLocationsLocationToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			locations_lst = append(locations_lst, modelMap)
		}
	}
	if err = d.Set("locations", locations_lst); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting locations %s", err))
	}

	return nil
}

// dataSourceIbmSccAccountLocationsID returns a reasonable ID for the list.
func dataSourceIbmSccAccountLocationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccAccountLocationsLocationToMap(model *adminserviceapiv1.Location) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.MainEndpointURL != nil {
		modelMap["main_endpoint_url"] = *model.MainEndpointURL
	}
	if model.GovernanceEndpointURL != nil {
		modelMap["governance_endpoint_url"] = *model.GovernanceEndpointURL
	}
	if model.ResultsEndpointURL != nil {
		modelMap["results_endpoint_url"] = *model.ResultsEndpointURL
	}
	if model.ComplianceEndpointURL != nil {
		modelMap["compliance_endpoint_url"] = *model.ComplianceEndpointURL
	}
	if model.AnalyticsEndpointURL != nil {
		modelMap["analytics_endpoint_url"] = *model.AnalyticsEndpointURL
	}
	if model.SiEndpointURL != nil {
		modelMap["si_endpoint_url"] = *model.SiEndpointURL
	}
	if model.Regions != nil {
		regions := []map[string]interface{}{}
		for _, regionsItem := range model.Regions {
			regionsItemMap, err := dataSourceIbmSccAccountLocationsRegionToMap(&regionsItem)
			if err != nil {
				return modelMap, err
			}
			regions = append(regions, regionsItemMap)
		}
		modelMap["regions"] = regions
	}
	return modelMap, nil
}

func dataSourceIbmSccAccountLocationsRegionToMap(model *adminserviceapiv1.Region) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}
