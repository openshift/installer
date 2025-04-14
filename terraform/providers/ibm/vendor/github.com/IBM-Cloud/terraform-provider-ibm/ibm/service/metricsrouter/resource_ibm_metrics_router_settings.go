// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package metricsrouter

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/metricsrouterv3"
)

func ResourceIBMMetricsRouterSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMMetricsRouterSettingsCreate,
		ReadContext:   resourceIBMMetricsRouterSettingsRead,
		UpdateContext: resourceIBMMetricsRouterSettingsUpdate,
		DeleteContext: resourceIBMMetricsRouterSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"default_targets": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of default target references.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The target uuid for a pre-defined metrics router target.",
							ValidateFunc: validate.InvokeValidator("ibm_metrics_router_settings", "id"),
						},
					},
				},
			},
			"permitted_target_regions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If present then only these regions may be used to define a target.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"primary_metadata_region": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_metrics_router_settings", "primary_metadata_region"),
				Description:  "To store all your meta data in a single region.",
			},
			"backup_metadata_region": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_metrics_router_settings", "backup_metadata_region"),
				Description:  "To backup all your meta data in a different region.",
			},
			"private_api_endpoint_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If you set this true then you cannot access api through public network.",
			},
		},
	}
}

func ResourceIBMMetricsRouterSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "primary_metadata_region",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-_]+$`,
			MinValueLength:             3,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "backup_metadata_region",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-_]+$`,
			MinValueLength:             3,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-._:]+$`,
			MinValueLength:             3,
			MaxValueLength:             1000,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_metrics_router_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMMetricsRouterSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSettingsOptions := &metricsrouterv3.UpdateSettingsOptions{}

	if _, ok := d.GetOk("default_targets"); ok {
		var defaultTargets []metricsrouterv3.TargetIdentity
		for _, e := range d.Get("default_targets").([]interface{}) {
			value := e.(map[string]interface{})
			defaultTargetsItem, err := resourceIBMMetricsRouterSettingsMapToTargetIdentity(value)
			if err != nil {
				return diag.FromErr(err)
			}
			defaultTargets = append(defaultTargets, *defaultTargetsItem)
		}
		updateSettingsOptions.SetDefaultTargets(defaultTargets)
	}

	if _, ok := d.GetOk("permitted_target_regions"); ok {
		updateSettingsOptions.SetPermittedTargetRegions(resourceInterfaceToStringArray(d.Get("permitted_target_regions").([]interface{})))
	}
	if _, ok := d.GetOk("primary_metadata_region"); ok {
		updateSettingsOptions.SetPrimaryMetadataRegion(d.Get("primary_metadata_region").(string))
	}
	if _, ok := d.GetOk("backup_metadata_region"); ok {
		updateSettingsOptions.SetBackupMetadataRegion(d.Get("backup_metadata_region").(string))
	}
	if _, ok := d.GetOk("private_api_endpoint_only"); ok {
		updateSettingsOptions.SetPrivateAPIEndpointOnly(d.Get("private_api_endpoint_only").(bool))
	}

	setting, response, err := metricsRouterClient.UpdateSettingsWithContext(context, updateSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*setting.PrimaryMetadataRegion)

	return resourceIBMMetricsRouterSettingsRead(context, d, meta)
}

func resourceIBMMetricsRouterSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &metricsrouterv3.GetSettingsOptions{}

	setting, response, err := metricsRouterClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	defaultTargets := []map[string]interface{}{}
	if setting.DefaultTargets != nil {
		for _, defaultTargetsItem := range setting.DefaultTargets {
			tId := &metricsrouterv3.TargetIdentity{
				ID: defaultTargetsItem.ID,
			}
			defaultTargetsItemMap, err := resourceIBMMetricsRouterSettingsTargetIdentityToMap(tId)
			if err != nil {
				return diag.FromErr(err)
			}
			defaultTargets = append(defaultTargets, defaultTargetsItemMap)
		}
	}
	if err = d.Set("default_targets", defaultTargets); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting default_targets: %s", err))
	}
	if setting.PermittedTargetRegions != nil {
		if err = d.Set("permitted_target_regions", setting.PermittedTargetRegions); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting permitted_target_regions: %s", err))
		}
	}
	if err = d.Set("primary_metadata_region", setting.PrimaryMetadataRegion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_metadata_region: %s", err))
	}
	if err = d.Set("backup_metadata_region", setting.BackupMetadataRegion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting backup_metadata_region: %s", err))
	}
	if err = d.Set("private_api_endpoint_only", setting.PrivateAPIEndpointOnly); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_api_endpoint_only: %s", err))
	}

	return nil
}

func resourceIBMMetricsRouterSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSettingsOptions := &metricsrouterv3.UpdateSettingsOptions{}

	updateSettingsOptions.SetPrimaryMetadataRegion(d.Id())

	hasChange := false

	if d.HasChange("default_targets") {
		if _, ok := d.GetOk("default_targets"); ok {
			var defaultTargets []metricsrouterv3.TargetIdentity
			for _, e := range d.Get("default_targets").([]interface{}) {
				value := e.(map[string]interface{})
				defaultTargetsItem, err := resourceIBMMetricsRouterSettingsMapToTargetIdentity(value)
				if err != nil {
					return diag.FromErr(err)
				}
				defaultTargets = append(defaultTargets, *defaultTargetsItem)
			}
			updateSettingsOptions.SetDefaultTargets(defaultTargets)
		} else {
			// In case, need to remove all the default_targets
			updateSettingsOptions.SetDefaultTargets([]metricsrouterv3.TargetIdentity{})
		}
		hasChange = true
	}
	if d.HasChange("permitted_target_regions") {
		if _, ok := d.GetOk("permitted_target_regions"); ok {
			updateSettingsOptions.SetPermittedTargetRegions(resourceInterfaceToStringArray(d.Get("permitted_target_regions").([]interface{})))
		} else {
			// In case, need to remove all the permitted_target_regions
			updateSettingsOptions.SetPermittedTargetRegions([]string{})
		}
		hasChange = true
	}
	if d.HasChange("primary_metadata_region") {
		updateSettingsOptions.SetPrimaryMetadataRegion(d.Get("primary_metadata_region").(string))
		hasChange = true
	}
	if d.HasChange("backup_metadata_region") {
		updateSettingsOptions.SetBackupMetadataRegion(d.Get("backup_metadata_region").(string))
		hasChange = true
	}
	if d.HasChange("private_api_endpoint_only") {
		updateSettingsOptions.SetPrivateAPIEndpointOnly(d.Get("private_api_endpoint_only").(bool))
		hasChange = true
	}

	if hasChange {
		_, response, err := metricsRouterClient.UpdateSettingsWithContext(context, updateSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSettingsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateSettingsWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMMetricsRouterSettingsRead(context, d, meta)
}

func resourceIBMMetricsRouterSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve old settings and put them for required fields.  Remove all other fields
	settings, response, err := metricsRouterClient.GetSettingsWithContext(context, &metricsrouterv3.GetSettingsOptions{})
	if err != nil {
		log.Printf("[DEBUG] UpdateSettingsWithContext with GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("with GetSettingsWithContext failed %s\n%s", err, response))
	}

	updateSettingsOptions := &metricsrouterv3.UpdateSettingsOptions{}

	updateSettingsOptions.PrimaryMetadataRegion = settings.PrimaryMetadataRegion
	updateSettingsOptions.BackupMetadataRegion = settings.BackupMetadataRegion
	updateSettingsOptions.PrivateAPIEndpointOnly = settings.PrivateAPIEndpointOnly
	updateSettingsOptions.PermittedTargetRegions = []string{}
	updateSettingsOptions.DefaultTargets = []metricsrouterv3.TargetIdentity{}

	_, res, err := metricsRouterClient.UpdateSettingsWithContext(context, updateSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateSettingsWithContext failed %s\n%s", err, res)
		return diag.FromErr(fmt.Errorf("UpdateSettingsWithContext failed %s\n%s", err, res))
	}

	d.SetId("")

	return nil
}

func resourceIBMMetricsRouterSettingsMapToTargetIdentity(modelMap map[string]interface{}) (*metricsrouterv3.TargetIdentity, error) {
	model := &metricsrouterv3.TargetIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMMetricsRouterSettingsMapToTargetRefernce(modelMap map[string]interface{}) (*metricsrouterv3.TargetReference, error) {
	model := &metricsrouterv3.TargetReference{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	return model, nil
}

func resourceIBMMetricsRouterSettingsTargetReferanceToMap(model *metricsrouterv3.TargetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["crn"] = model.CRN
	modelMap["name"] = model.Name
	modelMap["target_type"] = model.TargetType
	return modelMap, nil
}

func resourceIBMMetricsRouterSettingsTargetIdentityToMap(model *metricsrouterv3.TargetIdentity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceInterfaceToStringArray(resources []interface{}) (result []string) {
	result = make([]string, 0)
	for _, item := range resources {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return result
}
