// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudshell

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/ibmcloudshellv1"
)

func DataSourceIBMCloudShellAccountSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCloudShellAccountSettingsRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account ID in which the account settings belong to.",
			},
			"rev": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique revision number for the settings object.",
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
			"default_enable_new_features": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "You can choose which Cloud Shell features are available in the account and whether any new features are enabled as they become available. The feature settings apply only to the enabled Cloud Shell locations.",
			},
			"default_enable_new_regions": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and session data are stored. By default, users are routed to the nearest available location.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When enabled, Cloud Shell is available to all users in the account.",
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Cloud Shell features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "State of the feature.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the feature.",
						},
					},
				},
			},
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Cloud Shell region settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "State of the region.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the region.",
						},
					},
				},
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

func dataSourceIBMCloudShellAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &ibmcloudshellv1.GetAccountSettingsOptions{}

	getAccountSettingsOptions.SetAccountID(d.Get("account_id").(string))

	accountSettings, response, err := ibmCloudShellClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil || accountSettings == nil {
		log.Printf("[DEBUG] GetAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAccountSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*accountSettings.AccountID)
	if err = d.Set("rev", accountSettings.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting rev: %s", err))
	}
	if err = d.Set("created_at", flex.IntValue(accountSettings.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", accountSettings.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
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
		err = d.Set("features", dataSourceAccountSettingsFlattenFeatures(accountSettings.Features))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting features %s", err))
		}
	}

	if accountSettings.Regions != nil {
		err = d.Set("regions", dataSourceAccountSettingsFlattenRegions(accountSettings.Regions))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting regions %s", err))
		}
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

func dataSourceAccountSettingsFlattenFeatures(result []ibmcloudshellv1.Feature) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, dataSourceAccountSettingsFeaturesToMap(featuresItem))
	}

	return features
}

func dataSourceAccountSettingsFeaturesToMap(featuresItem ibmcloudshellv1.Feature) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.Enabled != nil {
		featuresMap["enabled"] = featuresItem.Enabled
	}
	if featuresItem.Key != nil {
		featuresMap["key"] = featuresItem.Key
	}

	return featuresMap
}

func dataSourceAccountSettingsFlattenRegions(result []ibmcloudshellv1.RegionSetting) (regions []map[string]interface{}) {
	for _, regionsItem := range result {
		regions = append(regions, dataSourceAccountSettingsRegionsToMap(regionsItem))
	}

	return regions
}

func dataSourceAccountSettingsRegionsToMap(regionsItem ibmcloudshellv1.RegionSetting) (regionsMap map[string]interface{}) {
	regionsMap = map[string]interface{}{}

	if regionsItem.Enabled != nil {
		regionsMap["enabled"] = regionsItem.Enabled
	}
	if regionsItem.Key != nil {
		regionsMap["key"] = regionsItem.Key
	}

	return regionsMap
}
