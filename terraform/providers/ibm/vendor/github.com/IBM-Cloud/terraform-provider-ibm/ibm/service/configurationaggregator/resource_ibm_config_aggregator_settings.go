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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/configuration-aggregator-go-sdk/configurationaggregatorv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmConfigAggregatorSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmConfigAggregatorSettingsCreate,
		ReadContext:   resourceIbmConfigAggregatorSettingsRead,
		DeleteContext: resourceIbmConfigAggregatorSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"resource_collection_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "The field denoting if the resource collection is enabled.",
			},
			"trusted_profile_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_config_aggregator_settings", "trusted_profile_id"),
				Description:  "The trusted profile id that provides Reader access to the App Configuration instance to collect resource metadata.",
			},
			"resource_collection_regions": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of regions across which the resource collection is enabled.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"additional_scope": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The additional scope that enables resource collection for Enterprise acccounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of scope. Currently allowed value is Enterprise.",
						},
						"enterprise_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Enterprise ID.",
						},
						"profile_template": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The Profile Template details applied on the enterprise account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Profile Template ID created in the enterprise account that provides access to App Configuration instance for resource collection.",
									},
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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

func ResourceIbmConfigAggregatorSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "trusted_profile_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-]*$`,
			MinValueLength:             44,
			MaxValueLength:             44,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_config_aggregator_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmConfigAggregatorSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationAggregatorClient, err := meta.(conns.ClientSession).ConfigurationAggregatorV1()
	region := getConfigurationInstanceRegion(configurationAggregatorClient, d)
	instanceId := d.Get("instance_id").(string)
	configurationAggregatorClient = getClientWithConfigurationInstanceEndpoint(configurationAggregatorClient, instanceId, region)
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	replaceSettingsOptions := &configurationaggregatorv1.ReplaceSettingsOptions{}
	replaceSettingsOptions.SetResourceCollectionEnabled(d.Get("resource_collection_enabled").(bool))
	if _, ok := d.GetOk("trusted_profile_id"); ok {
		replaceSettingsOptions.SetTrustedProfileID(d.Get("trusted_profile_id").(string))
	}
	if _, ok := d.GetOk("resource_collection_regions"); ok {
		var regions []string
		for _, v := range d.Get("resource_collection_regions").([]interface{}) {
			regionsItem := v.(string)
			regions = append(regions, regionsItem)
		}
		replaceSettingsOptions.SetRegions(regions)
	}
	if _, ok := d.GetOk("additional_scope"); ok {
		var additionalScope []configurationaggregatorv1.AdditionalScope
		for _, v := range d.Get("additional_scope").([]interface{}) {
			value := v.(map[string]interface{})
			additionalScopeItem, err := ResourceIbmConfigAggregatorSettingsMapToAdditionalScope(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "create", "parse-additional_scope").GetDiag()
			}
			additionalScope = append(additionalScope, *additionalScopeItem)
		}
		replaceSettingsOptions.SetAdditionalScope(additionalScope)
	}
	_, _, err = configurationAggregatorClient.ReplaceSettings(replaceSettingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceSettingsWithContext failed: %s", err.Error()), "ibm_config_aggregator_settings", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	aggregatorID := fmt.Sprintf("%s/%s", region, instanceId)
	d.SetId(aggregatorID)

	return resourceIbmConfigAggregatorSettingsRead(context, d, meta)
}

func resourceIbmConfigAggregatorSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationAggregatorClient, err := meta.(conns.ClientSession).ConfigurationAggregatorV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSettingsOptions := &configurationaggregatorv1.GetSettingsOptions{}
	var region string
	var instanceId string
	configurationAggregatorClient, region, instanceId, err = updateClientURLWithInstanceEndpoint(d.Id(), configurationAggregatorClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	settingsResponse, response, err := configurationAggregatorClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSettingsWithContext failed: %s", err.Error()), "ibm_config_aggregator_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if !core.IsNil(settingsResponse.ResourceCollectionEnabled) {
		if err = d.Set("resource_collection_enabled", settingsResponse.ResourceCollectionEnabled); err != nil {
			err = fmt.Errorf("Error setting resource_collection_enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read", "set-resource_collection_enabled").GetDiag()
		}
	}
	if !core.IsNil(settingsResponse.TrustedProfileID) {
		if err = d.Set("trusted_profile_id", settingsResponse.TrustedProfileID); err != nil {
			err = fmt.Errorf("Error setting trusted_profile_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read", "set-trusted_profile_id").GetDiag()
		}
	}
	if !core.IsNil(settingsResponse.Regions) {
		if err = d.Set("resource_collection_regions", settingsResponse.Regions); err != nil {
			err = fmt.Errorf("Error setting regions: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read", "set-regions").GetDiag()
		}
	}
	if !core.IsNil(settingsResponse.AdditionalScope) {
		additionalScope := []map[string]interface{}{}
		for _, additionalScopeItem := range settingsResponse.AdditionalScope {
			additionalScopeItemMap, err := ResourceIbmConfigAggregatorSettingsAdditionalScopeToMap(&additionalScopeItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read", "additional_scope-to-map").GetDiag()
			}
			additionalScope = append(additionalScope, additionalScopeItemMap)
		}
		if err = d.Set("additional_scope", additionalScope); err != nil {
			err = fmt.Errorf("Error setting additional_scope: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_config_aggregator_settings", "read", "set-additional_scope").GetDiag()
		}
	}

	return nil
}

func resourceIbmConfigAggregatorSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
	d.SetId("")
	return nil
}

func ResourceIbmConfigAggregatorSettingsMapToAdditionalScope(modelMap map[string]interface{}) (*configurationaggregatorv1.AdditionalScope, error) {
	model := &configurationaggregatorv1.AdditionalScope{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["enterprise_id"] != nil && modelMap["enterprise_id"].(string) != "" {
		model.EnterpriseID = core.StringPtr(modelMap["enterprise_id"].(string))
	}
	if modelMap["profile_template"] != nil && len(modelMap["profile_template"].([]interface{})) > 0 {
		ProfileTemplateModel, err := ResourceIbmConfigAggregatorSettingsMapToProfileTemplate(modelMap["profile_template"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ProfileTemplate = ProfileTemplateModel
	}
	return model, nil
}

func ResourceIbmConfigAggregatorSettingsMapToProfileTemplate(modelMap map[string]interface{}) (*configurationaggregatorv1.ProfileTemplate, error) {
	model := &configurationaggregatorv1.ProfileTemplate{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["trusted_profile_id"] != nil && modelMap["trusted_profile_id"].(string) != "" {
		model.TrustedProfileID = core.StringPtr(modelMap["trusted_profile_id"].(string))
	}
	return model, nil
}

func ResourceIbmConfigAggregatorSettingsAdditionalScopeToMap(model *configurationaggregatorv1.AdditionalScope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.EnterpriseID != nil {
		modelMap["enterprise_id"] = *model.EnterpriseID
	}
	if model.ProfileTemplate != nil {
		profileTemplateMap, err := ResourceIbmConfigAggregatorSettingsProfileTemplateToMap(model.ProfileTemplate)
		if err != nil {
			return modelMap, err
		}
		modelMap["profile_template"] = []map[string]interface{}{profileTemplateMap}
	}
	return modelMap, nil
}

func ResourceIbmConfigAggregatorSettingsProfileTemplateToMap(model *configurationaggregatorv1.ProfileTemplate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = *model.TrustedProfileID
	}
	return modelMap, nil
}
