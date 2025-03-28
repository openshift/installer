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

func DataSourceIbmSccControlLibraries() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccControlLibrariesRead,

		Schema: map[string]*schema.Schema{
			"control_library_type": {
				Type:         schema.TypeString,
				Description:  "The type of control library to be found.",
				ValidateFunc: validate.InvokeValidator("ibm_scc_control_library", "control_library_type"),
				Optional:     true,
			},
			"control_libraries": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of control libraries found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the control library.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of associated with the control library.",
						},
						// "instance_id": {
						// 	Type:        schema.TypeString,
						// 	Computed:    true,
						// 	Description: "The profile description.",
						// },
						"control_library_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the control library.",
						},
						"control_library_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the control library.",
						},
						"control_library_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the control library.",
						},
						"version_group_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version group label of the control library.",
						},
						"control_library_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the control library.",
						},
						"latest": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The latest version of the control library.",
						},
						// "hierarchy_enabled": {
						// 	Type:        schema.TypeBool,
						// 	Computed:    true,
						// 	Description: "The indication of whether hierarchy is enabled for the control library.",
						// },
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the control library.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the control library was created.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who updated the control library.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the control library was updated.",
						},
						"controls_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of controls for the control library.",
						},
						// "control_parents_count": {
						// 	Type:        schema.TypeInt,
						// 	Computed:    true,
						// 	Description: "The number of parent controls for the control library.",
						// },
					},
				},
			},
		},
	})
}

func dataSourceIbmSccControlLibrariesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listControlLibrariesOptions := &securityandcompliancecenterapiv3.ListControlLibrariesOptions{}
	listControlLibrariesOptions.SetInstanceID(d.Get("instance_id").(string))
	if val, ok := d.GetOk("control_library_type"); ok && val != nil {
		listControlLibrariesOptions.SetControlLibraryType(val.(string))
	}

	pager, err := securityandcompliancecenterapiClient.NewControlLibrariesPager(listControlLibrariesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListControlLibrarysWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListControlLibrarysWithContext failed %s", err))
	}
	controlLibraryList, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ListControlLibrarysWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListControlLibrarysWithContext failed %s", err))
	}
	d.SetId(fmt.Sprintf("%s/control_libraries", d.Get("instance_id").(string)))
	if err = d.Set("instance_id", d.Get("instance_id")); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id %s", err))
	}
	controlLibraries := []map[string]interface{}{}
	for _, cl := range controlLibraryList {
		modelMap, err := dataSourceIbmSccControlLibraryToMap(&cl)
		if err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting control library:%v\n%s", cl, err))
		}
		controlLibraries = append(controlLibraries, modelMap)
	}
	if err = d.Set("control_libraries", controlLibraries); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_libraries: %s", err))
	}
	return nil
}

func dataSourceIbmSccControlLibraryToMap(controlLibrary *securityandcompliancecenterapiv3.ControlLibrary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if controlLibrary.ID != nil {
		modelMap["id"] = controlLibrary.ID
	}
	if controlLibrary.AccountID != nil {
		modelMap["account_id"] = controlLibrary.AccountID
	}
	// if controlLibrary.InstanceID != nil {
	// 	modelMap["instance_id"] = controlLibrary.InstanceID
	// }
	if controlLibrary.ControlLibraryName != nil {
		modelMap["control_library_name"] = controlLibrary.ControlLibraryName
	}
	if controlLibrary.ControlLibraryDescription != nil {
		modelMap["control_library_description"] = controlLibrary.ControlLibraryDescription
	}
	if controlLibrary.ControlLibraryType != nil {
		modelMap["control_library_type"] = controlLibrary.ControlLibraryType
	}
	if controlLibrary.VersionGroupLabel != nil {
		modelMap["version_group_label"] = controlLibrary.VersionGroupLabel
	}
	if controlLibrary.ControlLibraryVersion != nil {
		modelMap["control_library_version"] = controlLibrary.ControlLibraryVersion
	}
	if controlLibrary.Latest != nil {
		modelMap["latest"] = controlLibrary.Latest
	}
	// if controlLibrary.HierarchyEnabled != nil {
	// 	modelMap["hierarchy_enabled"] = controlLibrary.HierarchyEnabled
	// }
	if controlLibrary.CreatedBy != nil {
		modelMap["created_by"] = controlLibrary.CreatedBy
	}
	if controlLibrary.CreatedOn != nil {
		modelMap["created_on"] = controlLibrary.CreatedOn.String()
	}
	if controlLibrary.UpdatedBy != nil {
		modelMap["updated_by"] = controlLibrary.UpdatedBy
	}
	if controlLibrary.UpdatedOn != nil {
		modelMap["updated_on"] = controlLibrary.UpdatedOn.String()
	}
	if controlLibrary.ControlsCount != nil {
		modelMap["controls_count"] = controlLibrary.ControlsCount
	}
	// if controlLibrary.ControlParentCount != nil {
	// 	modelMap["controls_parents_count"] = controlLibrary.ControlParentsCount
	// }
	return modelMap, nil
}
