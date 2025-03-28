// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package cloudshell

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/ibmcloudshellv1"
)

func ResourceIBMCloudShellAccountSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCloudShellAccountSettingsCreate,
		ReadContext:   resourceIBMCloudShellAccountSettingsRead,
		UpdateContext: resourceIBMCloudShellAccountSettingsUpdate,
		DeleteContext: resourceIBMCloudShellAccountSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cloud_shell_account_settings", "account_id"),
				Description:  "The account ID in which the account settings belong to.",
			},
			"rev": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cloud_shell_account_settings", "rev"),
				Description:  "Unique revision number for the settings object.",
			},
			"default_enable_new_features": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "You can choose which Cloud Shell features are available in the account and whether any new features are enabled as they become available. The feature settings apply only to the enabled Cloud Shell locations.",
			},
			"default_enable_new_regions": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and session data are stored. By default, users are routed to the nearest available location.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When enabled, Cloud Shell is available to all users in the account.",
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Cloud Shell features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "State of the feature.",
						},
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the feature.",
						},
					},
				},
			},
			"regions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Cloud Shell region settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "State of the region.",
						},
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the region.",
						},
					},
				},
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique id of the settings object.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation timestamp in Unix epoch time.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of creator.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of api response object.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of last update in Unix epoch time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of last updater.",
			},
		},
	}
}

func ResourceIBMCloudShellAccountSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "account_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-]*$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
		validate.ValidateSchema{
			Identifier:                 "rev",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9-]*$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cloud_shell_account_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCloudShellAccountSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateAccountSettingsOptions := &ibmcloudshellv1.UpdateAccountSettingsOptions{}

	updateAccountSettingsOptions.SetAccountID(d.Get("account_id").(string))
	if _, ok := d.GetOk("rev"); ok {
		updateAccountSettingsOptions.SetRev(d.Get("rev").(string))
	}
	if _, ok := d.GetOk("default_enable_new_features"); ok {
		updateAccountSettingsOptions.SetDefaultEnableNewFeatures(d.Get("default_enable_new_features").(bool))
	}
	if _, ok := d.GetOk("default_enable_new_regions"); ok {
		updateAccountSettingsOptions.SetDefaultEnableNewRegions(d.Get("default_enable_new_regions").(bool))
	}
	if _, ok := d.GetOk("enabled"); ok {
		updateAccountSettingsOptions.SetEnabled(d.Get("enabled").(bool))
	}
	if _, ok := d.GetOk("features"); ok {
		var features []ibmcloudshellv1.Feature
		for _, v := range d.Get("features").([]interface{}) {
			value := v.(map[string]interface{})
			featuresItem, err := ResourceIBMCloudShellAccountSettingsMapToFeature(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "parse-features").GetDiag()
			}
			features = append(features, *featuresItem)
		}
		updateAccountSettingsOptions.SetFeatures(features)
	}
	if _, ok := d.GetOk("regions"); ok {
		var regions []ibmcloudshellv1.RegionSetting
		for _, v := range d.Get("regions").([]interface{}) {
			value := v.(map[string]interface{})
			regionsItem, err := ResourceIBMCloudShellAccountSettingsMapToRegionSetting(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "parse-regions").GetDiag()
			}
			regions = append(regions, *regionsItem)
		}
		updateAccountSettingsOptions.SetRegions(regions)
	}

	accountSettings, _, err := ibmCloudShellClient.UpdateAccountSettingsWithContext(context, updateAccountSettingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateAccountSettingsWithContext failed: %s", err.Error()), "ibm_cloud_shell_account_settings", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*accountSettings.ID)

	return resourceIBMCloudShellAccountSettingsRead(context, d, meta)
}

func resourceIBMCloudShellAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAccountSettingsOptions := &ibmcloudshellv1.GetAccountSettingsOptions{}

	getAccountSettingsOptions.SetAccountID(strings.TrimPrefix(d.Id(), "ac-"))

	accountSettings, response, err := ibmCloudShellClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAccountSettingsWithContext failed: %s", err.Error()), "ibm_cloud_shell_account_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("account_id", accountSettings.AccountID); err != nil {
		err = fmt.Errorf("Error setting account_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-account_id").GetDiag()
	}
	if !core.IsNil(accountSettings.Rev) {
		if err = d.Set("rev", accountSettings.Rev); err != nil {
			err = fmt.Errorf("Error setting rev: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-rev").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.DefaultEnableNewFeatures) {
		if err = d.Set("default_enable_new_features", accountSettings.DefaultEnableNewFeatures); err != nil {
			err = fmt.Errorf("Error setting default_enable_new_features: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-default_enable_new_features").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.DefaultEnableNewRegions) {
		if err = d.Set("default_enable_new_regions", accountSettings.DefaultEnableNewRegions); err != nil {
			err = fmt.Errorf("Error setting default_enable_new_regions: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-default_enable_new_regions").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.Enabled) {
		if err = d.Set("enabled", accountSettings.Enabled); err != nil {
			err = fmt.Errorf("Error setting enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-enabled").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.Features) {
		features := []map[string]interface{}{}
		for _, featuresItem := range accountSettings.Features {
			featuresItemMap, err := ResourceIBMCloudShellAccountSettingsFeatureToMap(&featuresItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "features-to-map").GetDiag()
			}
			features = append(features, featuresItemMap)
		}
		if err = d.Set("features", features); err != nil {
			err = fmt.Errorf("Error setting features: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-features").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.Regions) {
		regions := []map[string]interface{}{}
		for _, regionsItem := range accountSettings.Regions {
			regionsItemMap, err := ResourceIBMCloudShellAccountSettingsRegionSettingToMap(&regionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "regions-to-map").GetDiag()
			}
			regions = append(regions, regionsItemMap)
		}
		if err = d.Set("regions", regions); err != nil {
			err = fmt.Errorf("Error setting regions: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-regions").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.ID) {
		if err = d.Set("id", accountSettings.ID); err != nil {
			err = fmt.Errorf("Error setting id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-id").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.CreatedAt) {
		if err = d.Set("created_at", flex.IntValue(accountSettings.CreatedAt)); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.CreatedBy) {
		if err = d.Set("created_by", accountSettings.CreatedBy); err != nil {
			err = fmt.Errorf("Error setting created_by: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-created_by").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.Type) {
		if err = d.Set("type", accountSettings.Type); err != nil {
			err = fmt.Errorf("Error setting type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-type").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.UpdatedAt) {
		if err = d.Set("updated_at", flex.IntValue(accountSettings.UpdatedAt)); err != nil {
			err = fmt.Errorf("Error setting updated_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-updated_at").GetDiag()
		}
	}
	if !core.IsNil(accountSettings.UpdatedBy) {
		if err = d.Set("updated_by", accountSettings.UpdatedBy); err != nil {
			err = fmt.Errorf("Error setting updated_by: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "read", "set-updated_by").GetDiag()
		}
	}

	return nil
}

func resourceIBMCloudShellAccountSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateAccountSettingsOptions := &ibmcloudshellv1.UpdateAccountSettingsOptions{}

	updateAccountSettingsOptions.SetAccountID(strings.TrimPrefix(d.Id(), "ac-"))

	hasChange := false

	updateAccountSettingsOptions.SetRev(d.Get("rev").(string))
	updateAccountSettingsOptions.SetDefaultEnableNewFeatures(d.Get("default_enable_new_features").(bool))
	updateAccountSettingsOptions.SetDefaultEnableNewRegions(d.Get("default_enable_new_regions").(bool))
	updateAccountSettingsOptions.SetEnabled(d.Get("enabled").(bool))

	var features []ibmcloudshellv1.Feature
	for _, v := range d.Get("features").([]interface{}) {
		value := v.(map[string]interface{})
		featuresItem, err := ResourceIBMCloudShellAccountSettingsMapToFeature(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "parse-features").GetDiag()
		}
		features = append(features, *featuresItem)
	}
	updateAccountSettingsOptions.SetFeatures(features)

	var regions []ibmcloudshellv1.RegionSetting
	for _, v := range d.Get("regions").([]interface{}) {
		value := v.(map[string]interface{})
		regionsItem, err := ResourceIBMCloudShellAccountSettingsMapToRegionSetting(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cloud_shell_account_settings", "delete", "parse-regions").GetDiag()
		}
		regions = append(regions, *regionsItem)
	}
	updateAccountSettingsOptions.SetRegions(regions)

	if d.HasChange("account_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "account_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cloud_shell_account_settings", "delete", "account_id-forces-new").GetDiag()
	}

	if d.HasChange("default_enable_new_features") || d.HasChange("default_enable_new_regions") || d.HasChange("enabled") || d.HasChange("features") || d.HasChange("regions") {
		hasChange = true
	}

	if hasChange {
		_, _, err = ibmCloudShellClient.UpdateAccountSettingsWithContext(context, updateAccountSettingsOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateAccountSettingsWithContext failed: %s", err.Error()), "ibm_cloud_shell_account_settings", "delete")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCloudShellAccountSettingsRead(context, d, meta)
}

func resourceIBMCloudShellAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Cloud Shell does not support delete of account settings subsequently delete is a no-op.
	return nil
}

func ResourceIBMCloudShellAccountSettingsMapToFeature(modelMap map[string]interface{}) (*ibmcloudshellv1.Feature, error) {
	model := &ibmcloudshellv1.Feature{}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	return model, nil
}

func ResourceIBMCloudShellAccountSettingsMapToRegionSetting(modelMap map[string]interface{}) (*ibmcloudshellv1.RegionSetting, error) {
	model := &ibmcloudshellv1.RegionSetting{}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	return model, nil
}

func ResourceIBMCloudShellAccountSettingsFeatureToMap(model *ibmcloudshellv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	return modelMap, nil
}

func ResourceIBMCloudShellAccountSettingsRegionSettingToMap(model *ibmcloudshellv1.RegionSetting) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	return modelMap, nil
}
