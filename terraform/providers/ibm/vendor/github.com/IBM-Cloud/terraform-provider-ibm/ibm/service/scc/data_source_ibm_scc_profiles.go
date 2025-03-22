// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProfiles() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProfilesRead,

		Schema: map[string]*schema.Schema{
			"profile_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of profiles found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile ID.",
						},
						"profile_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile name.",
						},
						"profile_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile description.",
						},
						"profile_type": {
							Type:         schema.TypeString,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_scc_profile", "profile_type"),
							Description:  "The profile type, such as custom or predefined.",
						},
						"profile_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version status of the profile.",
						},
						"version_group_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version group label of the profile.",
						},
						"latest": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The latest version of the profile.",
						},
						"hierarchy_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The indication of whether hierarchy is enabled for the profile.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the profile.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the profile was created.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who updated the profile.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the profile was updated.",
						},
						"controls_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of controls for the profile.",
						},
						"control_parents_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of parent controls for the profile.",
						},
						"attachments_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of attachments related to this profile.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listProfilesOptions := &securityandcompliancecenterapiv3.ListProfilesOptions{}
	listProfilesOptions.SetInstanceID(d.Get("instance_id").(string))
	if val, ok := d.GetOk("profile_type"); ok && val != nil {
		listProfilesOptions.SetProfileType(d.Get("profile_type").(string))
	}

	pager, err := securityandcompliancecenterapiClient.NewProfilesPager(listProfilesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListProfilesWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListProfilesWithContext failed %s", err))
	}
	profileList, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ListProfilesWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListProfilesWithContext failed %s", err))
	}
	d.SetId(fmt.Sprintf("%s/profiles", d.Get("instance_id").(string)))
	if err = d.Set("instance_id", d.Get("instance_id")); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id %s", err))
	}
	profiles := []map[string]interface{}{}
	for _, profile := range profileList {
		modelMap, err := dataSourceIbmSccProfileToMap(&profile)
		if err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting profile:%v\n%s", profile, err))
		}
		profiles = append(profiles, modelMap)
	}
	if err = d.Set("profiles", profiles); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profiles: %s", err))
	}
	return nil
}

func dataSourceIbmSccProfileToMap(profile *securityandcompliancecenterapiv3.Profile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if profile.ID != nil {
		modelMap["id"] = profile.ID
	}
	if profile.ProfileName != nil {
		modelMap["profile_name"] = profile.ProfileName
	}
	if profile.ProfileDescription != nil {
		modelMap["profile_description"] = profile.ProfileDescription
	}
	if profile.ProfileType != nil {
		modelMap["profile_type"] = profile.ProfileType
	}
	if profile.ProfileVersion != nil {
		modelMap["profile_version"] = profile.ProfileVersion
	}
	if profile.VersionGroupLabel != nil {
		modelMap["version_group_label"] = profile.VersionGroupLabel
	}
	if profile.Latest != nil {
		modelMap["latest"] = profile.Latest
	}
	if profile.CreatedBy != nil {
		modelMap["created_by"] = profile.CreatedBy
	}
	if profile.CreatedOn != nil {
		modelMap["created_on"] = profile.CreatedOn.String()
	}
	if profile.UpdatedBy != nil {
		modelMap["updated_by"] = profile.UpdatedBy
	}
	if profile.UpdatedOn != nil {
		modelMap["updated_on"] = profile.UpdatedOn.String()
	}
	if profile.ControlsCount != nil {
		modelMap["controls_count"] = profile.ControlsCount
	}
	if profile.AttachmentsCount != nil {
		modelMap["attachments_count"] = profile.AttachmentsCount
	}
	return modelMap, nil
}
