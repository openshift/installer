// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.92.0-af5c89a5-20240617-153232
 */

package configurationaggregator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/configuration-aggregator-go-sdk/configurationaggregatorv1"
)

func DataSourceIbmConfigAggregatorSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmConfigAggregatorSettingsRead,

		Schema: map[string]*schema.Schema{
			"resource_collection_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The field to check if the resource collection is enabled.",
			},
			"trusted_profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trusted profile ID that provides access to App Configuration instance to retrieve resource metadata.",
			},
			"last_updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the settings was last updated.",
			},
			"regions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Regions for which the resource collection is enabled.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"additional_scope": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The additional scope that enables resource collection for Enterprise acccounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of scope. Currently allowed value is Enterprise.",
						},
						"enterprise_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Enterprise ID.",
						},
						"profile_template": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Profile Template details applied on the enterprise account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Profile Template ID created in the enterprise account that provides access to App Configuration instance for resource collection.",
									},
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trusted profile ID that provides access to App Configuration instance to retrieve template information.",
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

func dataSourceIbmConfigAggregatorSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationAggregatorClient, err := meta.(conns.ClientSession).ConfigurationAggregatorV1()
	region := getConfigurationInstanceRegion(configurationAggregatorClient, d)
	instanceId := d.Get("instance_id").(string)
	configurationAggregatorClient = getClientWithConfigurationInstanceEndpoint(configurationAggregatorClient, instanceId, region)
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_config_aggregator_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSettingsOptions := &configurationaggregatorv1.GetSettingsOptions{}

	settingsResponse, _, err := configurationAggregatorClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSettingsWithContext failed: %s", err.Error()), "(Data) ibm_config_aggregator_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmConfigAggregatorSettingsID(d))

	if err = d.Set("resource_collection_enabled", settingsResponse.ResourceCollectionEnabled); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_collection_enabled: %s", err), "(Data) ibm_config_aggregator_settings", "read", "set-resource_collection_enabled").GetDiag()
	}

	if err = d.Set("trusted_profile_id", settingsResponse.TrustedProfileID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting trusted_profile_id: %s", err), "(Data) ibm_config_aggregator_settings", "read", "set-trusted_profile_id").GetDiag()
	}

	if err = d.Set("last_updated", flex.DateTimeToString(settingsResponse.LastUpdated)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_updated: %s", err), "(Data) ibm_config_aggregator_settings", "read", "set-last_updated").GetDiag()
	}

	additionalScope := []map[string]interface{}{}
	if settingsResponse.AdditionalScope != nil {
		for _, modelItem := range settingsResponse.AdditionalScope {
			modelMap, err := DataSourceIbmConfigAggregatorSettingsAdditionalScopeToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_config_aggregator_settings", "read", "additional_scope-to-map").GetDiag()
			}
			additionalScope = append(additionalScope, modelMap)
		}
	}
	if err = d.Set("additional_scope", additionalScope); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting additional_scope: %s", err), "(Data) ibm_config_aggregator_settings", "read", "set-additional_scope").GetDiag()
	}

	return nil
}

// dataSourceIbmConfigAggregatorSettingsID returns a reasonable ID for the list.
func dataSourceIbmConfigAggregatorSettingsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmConfigAggregatorSettingsAdditionalScopeToMap(model *configurationaggregatorv1.AdditionalScope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.EnterpriseID != nil {
		modelMap["enterprise_id"] = *model.EnterpriseID
	}
	if model.ProfileTemplate != nil {
		profileTemplateMap, err := DataSourceIbmConfigAggregatorSettingsProfileTemplateToMap(model.ProfileTemplate)
		if err != nil {
			return modelMap, err
		}
		modelMap["profile_template"] = []map[string]interface{}{profileTemplateMap}
	}
	return modelMap, nil
}

func DataSourceIbmConfigAggregatorSettingsProfileTemplateToMap(model *configurationaggregatorv1.ProfileTemplate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = *model.TrustedProfileID
	}
	return modelMap, nil
}
