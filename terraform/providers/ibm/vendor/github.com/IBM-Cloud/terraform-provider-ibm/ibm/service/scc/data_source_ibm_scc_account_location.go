// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
)

func DataSourceIBMSccAccountLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSccAccountLocationRead,
		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The programatic ID of the location that you want to work in.",
			},
			"main_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base URL for the service.",
			},
			"governance_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint that is used to call the Configuration Governance APIs.",
			},
			"results_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint that is used to get the results for the Configuration Governance component.",
			},
			"compliance_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint that is used to call the Posture Management APIs.",
			},
			"analytics_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint that is used to generate analytics for the Posture Management component.",
			},
			"si_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint that is used to call the Security Insights APIs.",
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The programatic ID of the available regions.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSccAccountLocationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceApiClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getLocationOptions := &adminserviceapiv1.GetLocationOptions{}

	getLocationOptions.SetLocationID(d.Get("location_id").(string))

	location, response, err := adminServiceApiClient.GetLocationWithContext(context, getLocationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLocationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLocationWithContext failed %s\n%s", err, response))
	}

	d.SetId(*location.ID)
	if err = d.Set("main_endpoint_url", location.MainEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting main_endpoint_url: %s", err))
	}
	if err = d.Set("governance_endpoint_url", location.GovernanceEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting governance_endpoint_url: %s", err))
	}
	if err = d.Set("results_endpoint_url", location.ResultsEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting results_endpoint_url: %s", err))
	}
	if err = d.Set("compliance_endpoint_url", location.ComplianceEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting compliance_endpoint_url: %s", err))
	}
	if err = d.Set("analytics_endpoint_url", location.AnalyticsEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting analytics_endpoint_url: %s", err))
	}
	if err = d.Set("si_endpoint_url", location.SiEndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting si_endpoint_url: %s", err))
	}

	if location.Regions != nil {
		err = d.Set("regions", dataSourceLocationFlattenRegions(location.Regions))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting regions %s", err))
		}
	}

	return nil
}

func dataSourceLocationFlattenRegions(result []adminserviceapiv1.Region) (regions []map[string]interface{}) {
	for _, regionsItem := range result {
		regions = append(regions, dataSourceLocationRegionsToMap(regionsItem))
	}

	return regions
}

func dataSourceLocationRegionsToMap(regionsItem adminserviceapiv1.Region) (regionsMap map[string]interface{}) {
	regionsMap = map[string]interface{}{}

	if regionsItem.ID != nil {
		regionsMap["id"] = regionsItem.ID
	}

	return regionsMap
}
