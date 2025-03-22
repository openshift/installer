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
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProfileAttachment() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProfileAttachmentRead,

		Schema: map[string]*schema.Schema{
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The attachment ID.",
			},
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The profile ID.",
			},
			"attachment_item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the attachment.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID that is associated to the attachment.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID of the account that is associated to the attachment.",
			},
			"scope": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The scope payload for the multi cloud feature.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment that relates to this scope.",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The properties supported for scoping by this environment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the property.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the property.",
									},
								},
							},
						},
					},
				},
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the attachment was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the attachment.",
			},
			"updated_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the attachment was updated.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who updated the attachment.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of an attachment evaluation.",
			},
			"schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The schedule of an attachment evaluation.",
			},
			"notifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The request payload of the attachment notifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "enabled notifications.",
						},
						"controls": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The failed controls.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The threshold limit.",
									},
									"failed_control_ids": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The failed control IDs.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"attachment_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile parameters for the attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assessment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the implementation.",
						},
						"assessment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The implementation ID of the parameter.",
						},
						"parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter name.",
						},
						"parameter_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the parameter.",
						},
						"parameter_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter display name.",
						},
						"parameter_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter type.",
						},
					},
				},
			},
			"last_scan": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of the last scan of an attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the last scan of an attachment.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the last scan of an attachment.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the last scan started.",
						},
					},
				},
			},
			"next_scan_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the next scan.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the attachment.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description for the attachment.",
			},
		},
	})
}

func dataSourceIbmSccProfileAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileAttachmentOptions := &securityandcompliancecenterapiv3.GetProfileAttachmentOptions{}

	getProfileAttachmentOptions.SetAttachmentID(d.Get("attachment_id").(string))
	getProfileAttachmentOptions.SetProfileID(d.Get("profile_id").(string))
	getProfileAttachmentOptions.SetInstanceID(d.Get("instance_id").(string))

	attachmentItem, response, err := securityandcompliancecenterapiClient.GetProfileAttachmentWithContext(context, getProfileAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProfileAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProfileAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getProfileAttachmentOptions.AttachmentID, *getProfileAttachmentOptions.ProfileID))

	if err = d.Set("attachment_item_id", attachmentItem.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attachment_item_id: %s", err))
	}

	if err = d.Set("account_id", attachmentItem.AccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting account_id: %s", err))
	}

	if err = d.Set("instance_id", attachmentItem.InstanceID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}

	scope := []map[string]interface{}{}
	if attachmentItem.Scope != nil {
		for _, modelItem := range attachmentItem.Scope {
			modelMap, err := dataSourceIbmSccProfileAttachmentMultiCloudScopeToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			scope = append(scope, modelMap)
		}
	}
	if err = d.Set("scope", scope); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting scope %s", err))
	}

	if err = d.Set("created_on", flex.DateTimeToString(attachmentItem.CreatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
	}

	if err = d.Set("created_by", attachmentItem.CreatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_on", flex.DateTimeToString(attachmentItem.UpdatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
	}

	if err = d.Set("updated_by", attachmentItem.UpdatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
	}

	if err = d.Set("status", attachmentItem.Status); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting status: %s", err))
	}

	if err = d.Set("schedule", attachmentItem.Schedule); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting schedule: %s", err))
	}

	notifications := []map[string]interface{}{}
	if attachmentItem.Notifications != nil {
		modelMap, err := dataSourceIbmSccProfileAttachmentAttachmentsNotificationsPrototypeToMap(attachmentItem.Notifications)
		if err != nil {
			return diag.FromErr(err)
		}
		notifications = append(notifications, modelMap)
	}
	if err = d.Set("notifications", notifications); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting notifications %s", err))
	}

	attachmentParameters := []map[string]interface{}{}
	if attachmentItem.AttachmentParameters != nil {
		for _, modelItem := range attachmentItem.AttachmentParameters {
			modelMap, err := dataSourceIbmSccProfileAttachmentAttachmentParameterPrototypeToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			attachmentParameters = append(attachmentParameters, modelMap)
		}
	}
	if err = d.Set("attachment_parameters", attachmentParameters); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attachment_parameters %s", err))
	}

	lastScan := []map[string]interface{}{}
	if attachmentItem.LastScan != nil {
		modelMap, err := dataSourceIbmSccProfileAttachmentLastScanToMap(attachmentItem.LastScan)
		if err != nil {
			return diag.FromErr(err)
		}
		lastScan = append(lastScan, modelMap)
	}
	if err = d.Set("last_scan", lastScan); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting last_scan %s", err))
	}

	if err = d.Set("next_scan_time", flex.DateTimeToString(attachmentItem.NextScanTime)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting next_scan_time: %s", err))
	}

	if err = d.Set("name", attachmentItem.Name); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting name: %s", err))
	}

	if err = d.Set("description", attachmentItem.Description); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting description: %s", err))
	}

	return nil
}

func dataSourceIbmSccProfileAttachmentMultiCloudScopeToMap(model securityandcompliancecenterapiv3.MultiCloudScopePayload) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["environment"] = model.Environment
	properties := []map[string]interface{}{}
	for _, propertiesItem := range model.Properties {
		propertiesItemMap, err := dataSourceIbmSccProfileAttachmentPropertyItemToMap(propertiesItem)
		if err != nil {
			return modelMap, err
		}
		properties = append(properties, propertiesItemMap)
	}
	modelMap["properties"] = properties
	return modelMap, nil
}

func dataSourceIbmSccProfileAttachmentPropertyItemToMap(model securityandcompliancecenterapiv3.ScopePropertyIntf) (map[string]interface{}, error) {
	return scopePropertiesToMap(model)
}

func dataSourceIbmSccProfileAttachmentAttachmentsNotificationsPrototypeToMap(model *securityandcompliancecenterapiv3.AttachmentNotifications) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["enabled"] = model.Enabled
	controlsMap, err := dataSourceIbmSccProfileAttachmentFailedControlsToMap(model.Controls)
	if err != nil {
		return modelMap, err
	}
	modelMap["controls"] = []map[string]interface{}{controlsMap}
	return modelMap, nil
}

func dataSourceIbmSccProfileAttachmentFailedControlsToMap(model *securityandcompliancecenterapiv3.AttachmentNotificationsControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ThresholdLimit != nil {
		modelMap["threshold_limit"] = flex.IntValue(model.ThresholdLimit)
	}
	if model.FailedControlIds != nil {
		modelMap["failed_control_ids"] = model.FailedControlIds
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileAttachmentAttachmentParameterPrototypeToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AssessmentType != nil {
		modelMap["assessment_type"] = model.AssessmentType
	}
	if model.AssessmentID != nil {
		modelMap["assessment_id"] = model.AssessmentID
	}
	if model.ParameterName != nil {
		modelMap["parameter_name"] = model.ParameterName
	}
	if model.ParameterValue != nil {
		modelMap["parameter_value"] = model.ParameterValue
	}
	if model.ParameterDisplayName != nil {
		modelMap["parameter_display_name"] = model.ParameterDisplayName
	}
	if model.ParameterType != nil {
		modelMap["parameter_type"] = model.ParameterType
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileAttachmentLastScanToMap(model *securityandcompliancecenterapiv3.LastScan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.Time != nil {
		modelMap["time"] = model.Time.String()
	}
	return modelMap, nil
}
