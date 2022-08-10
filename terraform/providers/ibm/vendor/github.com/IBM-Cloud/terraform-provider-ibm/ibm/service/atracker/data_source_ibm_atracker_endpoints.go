// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"
)

func DataSourceIBMAtrackerEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceIBMAtrackerEndpointsRead,
		DeprecationMessage: "use Settings instead",
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Activity Tracker API endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public URL of Activity Tracker in a region.",
						},
						"public_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether or not the public endpoint is enabled in the account.",
						},
						"private_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private URL of Activity Tracker. This URL cannot be disabled.",
						},
						"private_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The private endpoint is always enabled.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAtrackerEndpointsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, _, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getEndpointsOptions := &atrackerv1.GetEndpointsOptions{}

	endpoints, response, err := atrackerClient.GetEndpointsWithContext(context, getEndpointsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetEndpointsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetEndpointsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMAtrackerEndpointsID(d))

	if endpoints.APIEndpoint != nil {
		err = d.Set("api_endpoint", dataSourceEndpointsFlattenAPIEndpoint(*endpoints.APIEndpoint))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting api_endpoint %s", err))
		}
	}

	return nil
}

// dataSourceIBMAtrackerEndpointsID returns a reasonable ID for the list.
func dataSourceIBMAtrackerEndpointsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceEndpointsFlattenAPIEndpoint(result atrackerv1.APIEndpoint) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceEndpointsAPIEndpointToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceEndpointsAPIEndpointToMap(apiEndpointItem atrackerv1.APIEndpoint) (apiEndpointMap map[string]interface{}) {
	apiEndpointMap = map[string]interface{}{}

	if apiEndpointItem.PublicURL != nil {
		apiEndpointMap["public_url"] = apiEndpointItem.PublicURL
	}
	if apiEndpointItem.PublicEnabled != nil {
		apiEndpointMap["public_enabled"] = apiEndpointItem.PublicEnabled
	}
	if apiEndpointItem.PrivateURL != nil {
		apiEndpointMap["private_url"] = apiEndpointItem.PrivateURL
	}
	if apiEndpointItem.PrivateEnabled != nil {
		apiEndpointMap["private_enabled"] = apiEndpointItem.PrivateEnabled
	}

	return apiEndpointMap
}
