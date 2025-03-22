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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/ibmcloudshellv1"
)

func DataSourceIBMCloudShellAccountSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCloudShellAccountSettingsRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account ID in which the account settings belong to.",
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique revision number for the settings object.",
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
			"default_enable_new_features": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "You can choose which Cloud Shell features are available in the account and whether any new features are enabled as they become available. The feature settings apply only to the enabled Cloud Shell locations.",
			},
			"default_enable_new_regions": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set whether Cloud Shell is enabled in a specific location for the account. The location determines where user and session data are stored. By default, users are routed to the nearest available location.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When enabled, Cloud Shell is available to all users in the account.",
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Cloud Shell features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "State of the feature.",
						},
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the feature.",
						},
					},
				},
			},
			"regions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Cloud Shell region settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "State of the region.",
						},
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the region.",
						},
					},
				},
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

func dataSourceIBMCloudShellAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudShellClient, err := meta.(conns.ClientSession).IBMCloudShellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cloud_shell_account_settings", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAccountSettingsOptions := &ibmcloudshellv1.GetAccountSettingsOptions{}

	getAccountSettingsOptions.SetAccountID(d.Get("account_id").(string))

	accountSettings, _, err := ibmCloudShellClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAccountSettingsWithContext failed: %s", err.Error()), "(Data) ibm_cloud_shell_account_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getAccountSettingsOptions.AccountID)

	if !core.IsNil(accountSettings.Rev) {
		if err = d.Set("rev", accountSettings.Rev); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-rev").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.CreatedAt) {
		if err = d.Set("created_at", flex.IntValue(accountSettings.CreatedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-created_at").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.CreatedBy) {
		if err = d.Set("created_by", accountSettings.CreatedBy); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_by: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-created_by").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.DefaultEnableNewFeatures) {
		if err = d.Set("default_enable_new_features", accountSettings.DefaultEnableNewFeatures); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting default_enable_new_features: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-default_enable_new_features").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.DefaultEnableNewRegions) {
		if err = d.Set("default_enable_new_regions", accountSettings.DefaultEnableNewRegions); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting default_enable_new_regions: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-default_enable_new_regions").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.Enabled) {
		if err = d.Set("enabled", accountSettings.Enabled); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enabled: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-enabled").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.Features) {
		features := []map[string]interface{}{}
		for _, featuresItem := range accountSettings.Features {
			featuresItemMap, err := DataSourceIBMCloudShellAccountSettingsFeatureToMap(&featuresItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cloud_shell_account_settings", "read", "features-to-map").GetDiag()
			}
			features = append(features, featuresItemMap)
		}
		if err = d.Set("features", features); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting features: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-features").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.Regions) {
		regions := []map[string]interface{}{}
		for _, regionsItem := range accountSettings.Regions {
			regionsItemMap, err := DataSourceIBMCloudShellAccountSettingsRegionSettingToMap(&regionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cloud_shell_account_settings", "read", "regions-to-map").GetDiag()
			}
			regions = append(regions, regionsItemMap)
		}
		if err = d.Set("regions", regions); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting regions: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-regions").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.Type) {
		if err = d.Set("type", accountSettings.Type); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-type").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.UpdatedAt) {
		if err = d.Set("updated_at", flex.IntValue(accountSettings.UpdatedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-updated_at").GetDiag()
		}
	}

	if !core.IsNil(accountSettings.UpdatedBy) {
		if err = d.Set("updated_by", accountSettings.UpdatedBy); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_by: %s", err), "(Data) ibm_cloud_shell_account_settings", "read", "set-updated_by").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMCloudShellAccountSettingsFeatureToMap(model *ibmcloudshellv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	return modelMap, nil
}

func DataSourceIBMCloudShellAccountSettingsRegionSettingToMap(model *ibmcloudshellv1.RegionSetting) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	return modelMap, nil
}
