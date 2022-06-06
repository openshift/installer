// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudshell

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The account ID in which the account settings belong to.",
			},
			"rev": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unique revision number for the settings object.",
			},
			"default_enable_new_features": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "You can choose which Cloud Shell features are available in the account and whether any new features are enabled as they become available. The feature settings apply only to the enabled Cloud Shell locations.",
			},
			"default_enable_new_regions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and session data are stored. By default, users are routed to the nearest available location.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When enabled, Cloud Shell is available to all users in the account.",
			},
			"features": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of Cloud Shell features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "State of the feature.",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the feature.",
						},
					},
				},
			},
			"regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of Cloud Shell region settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "State of the region.",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the region.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation timestamp in Unix epoch time.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of creator.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of api response object.",
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of last update in Unix epoch time.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of last updater.",
			},
		},
	}
}

func resourceIBMCloudShellAccountSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		return diag.FromErr(err)
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
		for _, e := range d.Get("features").([]interface{}) {
			value := e.(map[string]interface{})
			featuresItem := resourceIBMCloudShellAccountSettingsMapToFeature(value)
			features = append(features, featuresItem)
		}
		updateAccountSettingsOptions.SetFeatures(features)
	}
	if _, ok := d.GetOk("regions"); ok {
		var regions []ibmcloudshellv1.RegionSetting
		for _, e := range d.Get("regions").([]interface{}) {
			value := e.(map[string]interface{})
			regionsItem := resourceIBMCloudShellAccountSettingsMapToRegionSetting(value)
			regions = append(regions, regionsItem)
		}
		updateAccountSettingsOptions.SetRegions(regions)
	}

	accountSettings, response, err := ibmCloudShellClient.UpdateAccountSettingsWithContext(context, updateAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateAccountSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*accountSettings.ID)

	return resourceIBMCloudShellAccountSettingsRead(context, d, meta)
}

func resourceIBMCloudShellAccountSettingsMapToFeature(featureMap map[string]interface{}) ibmcloudshellv1.Feature {
	feature := ibmcloudshellv1.Feature{}

	if featureMap["enabled"] != nil {
		feature.Enabled = core.BoolPtr(featureMap["enabled"].(bool))
	}
	if featureMap["key"] != nil {
		feature.Key = core.StringPtr(featureMap["key"].(string))
	}

	return feature
}

func resourceIBMCloudShellAccountSettingsMapToRegionSetting(regionSettingMap map[string]interface{}) ibmcloudshellv1.RegionSetting {
	regionSetting := ibmcloudshellv1.RegionSetting{}

	if regionSettingMap["enabled"] != nil {
		regionSetting.Enabled = core.BoolPtr(regionSettingMap["enabled"].(bool))
	}
	if regionSettingMap["key"] != nil {
		regionSetting.Key = core.StringPtr(regionSettingMap["key"].(string))
	}

	return regionSetting
}

func resourceIBMCloudShellAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &ibmcloudshellv1.GetAccountSettingsOptions{}

	getAccountSettingsOptions.SetAccountID(strings.TrimPrefix(d.Id(), "ac-"))

	accountSettings, response, err := ibmCloudShellClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAccountSettingsWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("account_id", accountSettings.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}
	if err = d.Set("rev", accountSettings.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting rev: %s", err))
	}
	if err = d.Set("default_enable_new_features", accountSettings.DefaultEnableNewFeatures); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting default_enable_new_features: %s", err))
	}
	if err = d.Set("default_enable_new_regions", accountSettings.DefaultEnableNewRegions); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting default_enable_new_regions: %s", err))
	}
	if err = d.Set("enabled", accountSettings.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enabled: %s", err))
	}
	if accountSettings.Features != nil {
		features := []map[string]interface{}{}
		for _, featuresItem := range accountSettings.Features {
			featuresItemMap := resourceIBMCloudShellAccountSettingsFeatureToMap(featuresItem)
			features = append(features, featuresItemMap)
		}
		if err = d.Set("features", features); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting features: %s", err))
		}
	}
	if accountSettings.Regions != nil {
		regions := []map[string]interface{}{}
		for _, regionsItem := range accountSettings.Regions {
			regionsItemMap := resourceIBMCloudShellAccountSettingsRegionSettingToMap(regionsItem)
			regions = append(regions, regionsItemMap)
		}
		if err = d.Set("regions", regions); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting regions: %s", err))
		}
	}
	if err = d.Set("created_at", flex.IntValue(accountSettings.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", accountSettings.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}
	if err = d.Set("type", accountSettings.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}
	if err = d.Set("updated_at", flex.IntValue(accountSettings.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", accountSettings.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_by: %s", err))
	}

	return nil
}

func resourceIBMCloudShellAccountSettingsFeatureToMap(feature ibmcloudshellv1.Feature) map[string]interface{} {
	featureMap := map[string]interface{}{}

	if feature.Enabled != nil {
		featureMap["enabled"] = feature.Enabled
	}
	if feature.Key != nil {
		featureMap["key"] = feature.Key
	}

	return featureMap
}

func resourceIBMCloudShellAccountSettingsRegionSettingToMap(regionSetting ibmcloudshellv1.RegionSetting) map[string]interface{} {
	regionSettingMap := map[string]interface{}{}

	if regionSetting.Enabled != nil {
		regionSettingMap["enabled"] = regionSetting.Enabled
	}
	if regionSetting.Key != nil {
		regionSettingMap["key"] = regionSetting.Key
	}

	return regionSettingMap
}

func resourceIBMCloudShellAccountSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountSettingsOptions := &ibmcloudshellv1.UpdateAccountSettingsOptions{}

	updateAccountSettingsOptions.SetAccountID(strings.TrimPrefix(d.Id(), "ac-"))
	hasChange := false
	updateAccountSettingsOptions.SetRev(d.Get("rev").(string))
	if d.HasChange("default_enable_new_features") {
		updateAccountSettingsOptions.SetDefaultEnableNewFeatures(d.Get("default_enable_new_features").(bool))
		hasChange = true
	}
	if d.HasChange("default_enable_new_regions") {
		updateAccountSettingsOptions.SetDefaultEnableNewRegions(d.Get("default_enable_new_regions").(bool))
		hasChange = true
	}
	if d.HasChange("enabled") {
		updateAccountSettingsOptions.SetEnabled(d.Get("enabled").(bool))
		hasChange = true
	}
	if d.HasChange("features") {
		var features []ibmcloudshellv1.Feature
		for _, e := range d.Get("features").([]interface{}) {
			value := e.(map[string]interface{})
			featuresItem := resourceIBMCloudShellAccountSettingsMapToFeature(value)
			features = append(features, featuresItem)
		}
		updateAccountSettingsOptions.SetFeatures(features)
		hasChange = true
	}
	if d.HasChange("regions") {
		var regions []ibmcloudshellv1.RegionSetting
		for _, e := range d.Get("regions").([]interface{}) {
			value := e.(map[string]interface{})
			regionsItem := resourceIBMCloudShellAccountSettingsMapToRegionSetting(value)
			regions = append(regions, regionsItem)
		}
		updateAccountSettingsOptions.SetRegions(regions)
		hasChange = true
	}

	if hasChange {
		_, response, err := ibmCloudShellClient.UpdateAccountSettingsWithContext(context, updateAccountSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountSettingsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateAccountSettingsWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCloudShellAccountSettingsRead(context, d, meta)
}

func resourceIBMCloudShellAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Cloud Shell does not support delete of account settings subsequently delete is a no-op.
	return nil
}
