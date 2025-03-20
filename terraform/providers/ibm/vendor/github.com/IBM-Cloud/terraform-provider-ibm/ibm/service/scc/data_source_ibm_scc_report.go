// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccReport() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the report.",
			},
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.",
			},
			"created_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the report was created.",
			},
			"scan_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the scan was run.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the scan.",
			},
			"cos_object": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Object Storage object that is associated with the report.",
			},
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The account that is associated with a report.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account name.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account type.",
						},
					},
				},
			},
			"profile": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile name.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile version.",
						},
					},
				},
			},
			"attachment": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attachment that is associated with a report.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attachment ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the attachment.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the attachment.",
						},
						"schedule": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attachment schedule.",
						},
						"scope": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The scope of the attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this scope.",
									},
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The environment that relates to this scope.",
									},
									"properties": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The properties that are supported for scoping by this environment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property name.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property value.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccReportRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportOptions := &securityandcompliancecenterapiv3.GetReportOptions{}

	getReportOptions.SetReportID(d.Get("report_id").(string))
	getReportOptions.SetInstanceID(d.Get("instance_id").(string))

	report, response, err := resultsClient.GetReportWithContext(context, getReportOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportWithContext failed %s\n%s", err, response))
	}

	d.SetId(*report.ID)

	if err = d.Set("id", report.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting id: %s", err))
	}

	if err = d.Set("group_id", report.GroupID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting group_id: %s", err))
	}

	if err = d.Set("created_on", flex.DateTimeToString(report.CreatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
	}

	if err = d.Set("scan_time", report.ScanTime); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting scan_time: %s", err))
	}

	if err = d.Set("type", report.Type); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting type: %s", err))
	}

	if err = d.Set("cos_object", report.CosObject); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting cos_object: %s", err))
	}

	if err = d.Set("instance_id", report.InstanceID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}

	account := []map[string]interface{}{}
	if report.Account != nil {
		modelMap, err := dataSourceIbmSccReportAccountToMap(report.Account)
		if err != nil {
			return diag.FromErr(err)
		}
		account = append(account, modelMap)
	}
	if err = d.Set("account", account); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting account %s", err))
	}

	profile := []map[string]interface{}{}
	if report.Profile != nil {
		modelMap, err := dataSourceIbmSccReportProfileInfoToMap(report.Profile)
		if err != nil {
			return diag.FromErr(err)
		}
		profile = append(profile, modelMap)
	}
	if err = d.Set("profile", profile); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile %s", err))
	}

	attachment := []map[string]interface{}{}
	if report.Attachment != nil {
		modelMap, err := dataSourceIbmSccReportAttachmentToMap(report.Attachment)
		if err != nil {
			return diag.FromErr(err)
		}
		attachment = append(attachment, modelMap)
	}
	if err = d.Set("attachment", attachment); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attachment %s", err))
	}

	return nil
}

func dataSourceIbmSccReportAccountToMap(model *securityandcompliancecenterapiv3.Account) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func dataSourceIbmSccReportProfileInfoToMap(model *securityandcompliancecenterapiv3.ProfileInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	return modelMap, nil
}

func dataSourceIbmSccReportAttachmentToMap(model *securityandcompliancecenterapiv3.Attachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Schedule != nil {
		modelMap["schedule"] = model.Schedule
	}
	if model.Scopes != nil {
		scope := []map[string]interface{}{}
		for _, scopeItem := range model.Scopes {
			scopeItemMap, err := dataSourceIbmSccReportAttachmentScopeToMap(&scopeItem)
			if err != nil {
				return modelMap, err
			}
			scope = append(scope, scopeItemMap)
		}
		modelMap["scope"] = scope
	}
	return modelMap, nil
}

func dataSourceIbmSccReportAttachmentScopeToMap(model *securityandcompliancecenterapiv3.Scope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := dataSourceIbmSccReportScopePropertyToMap(propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	return modelMap, nil
}

func dataSourceIbmSccReportScopePropertyToMap(model securityandcompliancecenterapiv3.ScopePropertyIntf) (map[string]interface{}, error) {
	return scopePropertiesToMap(model)
}
