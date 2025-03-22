// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccProfileAttachment() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccProfileAttachmentCreate,
		ReadContext:   resourceIbmSccProfileAttachmentRead,
		UpdateContext: resourceIbmSccProfileAttachmentUpdate,
		DeleteContext: resourceIbmSccProfileAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile attachment ID.",
			},
			"profile_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_profile_attachment", "profile_id"),
				Description:  "The ID of the profile that is specified in the attachment.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID that is associated to the attachment.",
			},
			"scope": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The scope/scopes to link the profile attachment.",
				Elem:        schemaAttachmentScopes(),
				Set:         attachmentHashSchemaSetFunc("id"),
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
				Required:    true,
				Description: "The status of an attachment evaluation.",
			},
			"schedule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The schedule of an attachment evaluation.",
			},
			"notifications": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The request payload of the attachment notifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "enabled notifications.",
							DefaultFunc: func() (any, error) {
								return false, nil
							},
						},
						"controls": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Optional:    true,
							Description: "The failed controls.",
							DefaultFunc: func() (any, error) {
								return []map[string]interface{}{
									{
										"threshold_limit":    15,
										"failed_control_ids": []string{},
									},
								}, nil
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The threshold limit.",
										DefaultFunc: func() (any, error) {
											return 15, nil
										},
									},
									"failed_control_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The failed control IDs.",
										Elem:        &schema.Schema{Type: schema.TypeString},
										DefaultFunc: func() (any, error) {
											return []string{}, nil
										},
									},
								},
							},
						},
					},
				},
			},
			"attachment_parameters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The profile parameters for the attachment.",
				Elem:        schemaAttachmentParameters(),
				Set:         attachmentHashSchemaSetFunc("assessment_id", "parameter_name", "parameter_display_name", "parameter_type", "parameter_value"),
			},
			"last_scan": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of the last scan of an attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the last scan of an attachment.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The status of the last scan of an attachment.",
						},
						"time": {
							Type:        schema.TypeString,
							Optional:    true,
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
				Required:    true,
				Description: "The name of the attachment.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description for the attachment.",
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the attachment.",
			},
		},
	})
}

func ResourceIbmSccProfileAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_profile_attachment", Schema: validateSchema}
	return &resourceValidator
}

// hashAttachmentParameters will determine how to hash the AttachmentParameters schema.Resource
// It uses the 'assessment_id' in order to determine the difference.
func attachmentHashSchemaSetFunc(keys ...string) schema.SchemaSetFunc {
	return func(v interface{}) int {
		var str strings.Builder

		if m, ok := v.(map[string]interface{}); ok {
			for _, key := range keys {
				if v, ok := m[key]; ok {
					switch v := v.(type) {
					case bool:
						str.WriteRune('-')
						str.WriteString(strconv.FormatBool(v))
					case int:
						str.WriteRune('-')
						str.WriteString(strconv.Itoa(v))
					case string:
						str.WriteRune('-')
						str.WriteString(v)
					case []interface{}:
						str.WriteRune('-')
						s := make([]string, len(v))
						for i, v := range v {
							s[i] = fmt.Sprint(v)
						}
						str.WriteString(fmt.Sprintf("[%s]", strings.Join(s, ",")))
					}
				}
			}
		}
		log.Printf("[DEBUG] attachmentParameterSchemaSet value %s\n", str.String())

		return stringHashcode(str.String())
	}
}

// stringHashcode will return an integer to use for unique keys in the Schema.Set
func stringHashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// schemaAttachmentParameters returns a *schema.Resource for AttachmentParameters
func schemaAttachmentParameters() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"assessment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The implementation ID of the parameter.",
			},
			"assessment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "automated",
				Description: "The type of the implementation.",
			},
			"parameter_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter name.",
			},
			"parameter_display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter display name.",
			},
			"parameter_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter type.",
			},
			"parameter_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the parameter.",
			},
		},
	}
}

// schemaAttachmentScopes returns a *schema.Resource for AttachmentScopes
func schemaAttachmentScopes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The environment that relates to this scope.",
			},
			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The properties supported for scoping by this environment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the property.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the property.",
						},
					},
				},
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The scope id to target.",
			},
		},
	}
}

func resourceIbmSccProfileAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	bodyModelMap := map[string]interface{}{}
	createAttachmentOptions := &securityandcompliancecenterapiv3.CreateProfileAttachmentOptions{}
	instance_id := d.Get("instance_id").(string)
	bodyModelMap["instance_id"] = instance_id

	if _, ok := d.GetOk("profile_id"); ok {
		bodyModelMap["profile_id"] = d.Get("profile_id")
	}
	if _, ok := d.GetOk("description"); ok {
		bodyModelMap["description"] = d.Get("description")
	}
	if _, ok := d.GetOk("scope"); ok {
		bodyModelMap["scope"] = d.Get("scope")
	}

	// manual change
	if _, ok := d.GetOk("attachment_parameters"); ok {
		bodyModelMap["attachment_parameters"] = d.Get("attachment_parameters")
	}
	if _, ok := d.GetOk("notifications"); ok {
		bodyModelMap["notifications"] = d.Get("notifications")
	}
	// end manual change
	if _, ok := d.GetOk("status"); ok {
		bodyModelMap["status"] = d.Get("status")
	}
	if _, ok := d.GetOk("schedule"); ok {
		bodyModelMap["schedule"] = d.Get("schedule")
	}
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	convertedModel, err := resourceIbmSccProfileAttachmentMapToAttachmentPrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createAttachmentOptions = convertedModel

	attachmentPrototype, response, err := securityandcompliancecenterapiClient.CreateProfileAttachmentWithContext(context, createAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("CreateAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instance_id, *createAttachmentOptions.ProfileID, *attachmentPrototype.Attachments[0].ID))

	return resourceIbmSccProfileAttachmentRead(context, d, meta)
}

func resourceIbmSccProfileAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileAttachmentOptions := &securityandcompliancecenterapiv3.GetProfileAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileAttachmentOptions.SetInstanceID(parts[0])
	getProfileAttachmentOptions.SetProfileID(parts[1])
	getProfileAttachmentOptions.SetAttachmentID(parts[2])

	attachmentItem, response, err := securityandcompliancecenterapiClient.GetProfileAttachmentWithContext(context, getProfileAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProfileAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProfileAttachmentWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_id", parts[0]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}
	if !core.IsNil(attachmentItem.ID) {
		if err = d.Set("profile_attachment_id", attachmentItem.ID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting profile_id: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.ProfileID) {
		if err = d.Set("profile_id", attachmentItem.ProfileID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting profile_id: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.AccountID) {
		if err = d.Set("account_id", attachmentItem.AccountID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting account_id: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Scope) {
		// scope := []map[string]interface{}{}
		// for _, scopeItem := range attachmentItem.Scope {
		// 	scopeItemMap, err := resourceIbmSccProfileAttachmentMultiCloudScopeToMap(&scopeItem)
		// 	if err != nil {
		// 		return diag.FromErr(err)
		// 	}
		// 	scope = append(scope, scopeItemMap)
		// }
		attachmentScopes := &schema.Set{
			F: attachmentHashSchemaSetFunc("id"),
		}
		for _, scopeItem := range attachmentItem.Scope {
			scopeItemMap, err := resourceIbmSccProfileAttachmentMultiCloudScopeToMap(&scopeItem)
			if err != nil {
				return diag.FromErr(err)
			}
			attachmentScopes.Add(scopeItemMap)
		}
		if err = d.Set("scope", attachmentScopes); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting scope: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.CreatedOn) {
		if err = d.Set("created_on", flex.DateTimeToString(attachmentItem.CreatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.CreatedBy) {
		if err = d.Set("created_by", attachmentItem.CreatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.UpdatedOn) {
		if err = d.Set("updated_on", flex.DateTimeToString(attachmentItem.UpdatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.UpdatedBy) {
		if err = d.Set("updated_by", attachmentItem.UpdatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Status) {
		if err = d.Set("status", attachmentItem.Status); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting status: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Schedule) {
		if err = d.Set("schedule", attachmentItem.Schedule); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting schedule: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Notifications) {
		notificationsMap, err := resourceIbmSccProfileAttachmentAttachmentsNotificationsPrototypeToMap(attachmentItem.Notifications)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("notifications", []map[string]interface{}{notificationsMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting notifications: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.AttachmentParameters) {
		attachmentParameters := &schema.Set{
			F: attachmentHashSchemaSetFunc("assessment_id", "parameter_name", "parameter_display_name", "parameter_type", "parameter_value"),
		}
		for _, attachmentParametersItem := range attachmentItem.AttachmentParameters {
			attachmentParametersItemMap, err := resourceIbmSccProfileAttachmentAttachmentParameterPrototypeToMap(&attachmentParametersItem)
			if err != nil {
				return diag.FromErr(err)
			}
			attachmentParameters.Add(attachmentParametersItemMap)
		}
		if err = d.Set("attachment_parameters", attachmentParameters); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting attachment_parameters: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.LastScan) {
		lastScanMap, err := resourceIbmSccProfileAttachmentLastScanToMap(attachmentItem.LastScan)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("last_scan", []map[string]interface{}{lastScanMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting last_scan: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.NextScanTime) {
		if err = d.Set("next_scan_time", flex.DateTimeToString(attachmentItem.NextScanTime)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting next_scan_time: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Name) {
		if err = d.Set("name", attachmentItem.Name); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.Description) {
		if err = d.Set("description", attachmentItem.Description); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(attachmentItem.ID) {
		if err = d.Set("attachment_id", attachmentItem.ID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting attachment_id: %s", err))
		}
	}

	return nil
}

func resourceIbmSccProfileAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceProfileAttachmentOptions := &securityandcompliancecenterapiv3.ReplaceProfileAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	replaceProfileAttachmentOptions.SetInstanceID(parts[0])
	replaceProfileAttachmentOptions.SetProfileID(parts[1])
	replaceProfileAttachmentOptions.SetAttachmentID(parts[2])

	hasChange := false

	if d.HasChange("profile_id") {
		return diag.FromErr(flex.FmtErrorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "profile_id"))
	}

	if d.HasChange("schedule") {
		replaceProfileAttachmentOptions.SetSchedule(d.Get("schedule").(string))
		hasChange = true
	}

	if d.HasChange("name") {
		replaceProfileAttachmentOptions.SetName(d.Get("name").(string))
		hasChange = true
	}

	if d.HasChange("description") {
		replaceProfileAttachmentOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	if d.HasChange("attachment_parameters") {
		attachmentItems := d.Get("attachment_parameters")
		attachmentParameters := []securityandcompliancecenterapiv3.Parameter{}
		for _, attachmentParametersItem := range attachmentItems.(*schema.Set).List() {
			attachmentParametersItemModel, err := resourceIbmSccProfileAttachmentMapToAttachmentParameterPrototype(attachmentParametersItem.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			attachmentParameters = append(attachmentParameters, *attachmentParametersItemModel)
		}
		replaceProfileAttachmentOptions.SetAttachmentParameters(attachmentParameters)
		hasChange = true
	}

	if d.HasChange("scope") {
		hasChange = true
	}

	if d.HasChange("notifications") {
		notificationsItem := d.Get("notifications.0").(map[string]interface{})
		updateNotifications, err := resourceIbmSccProfileAttachmentMapToAttachmentsNotificationsPrototype(notificationsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		replaceProfileAttachmentOptions.SetNotifications(updateNotifications)
		hasChange = true
	}

	if hasChange {
		if replaceProfileAttachmentOptions.Name == nil {
			replaceProfileAttachmentOptions.SetName(d.Get("name").(string))
		}
		if replaceProfileAttachmentOptions.Status == nil {
			replaceProfileAttachmentOptions.SetStatus(d.Get("status").(string))
		}
		if replaceProfileAttachmentOptions.Schedule == nil {
			replaceProfileAttachmentOptions.SetSchedule(d.Get("schedule").(string))
		}
		if replaceProfileAttachmentOptions.Description == nil {
			replaceProfileAttachmentOptions.SetDescription(d.Get("description").(string))
		}
		if replaceProfileAttachmentOptions.Notifications == nil {
			notificationsItem := d.Get("notifications.0").(map[string]interface{})
			updateNotifications, err := resourceIbmSccProfileAttachmentMapToAttachmentsNotificationsPrototype(notificationsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			replaceProfileAttachmentOptions.SetNotifications(updateNotifications)
		}
		if replaceProfileAttachmentOptions.AttachmentParameters == nil || d.Get("attachment_parameters") != nil {
			attachmentItems := d.Get("attachment_parameters")
			attachmentParameters := []securityandcompliancecenterapiv3.Parameter{}
			for _, attachmentParametersItem := range attachmentItems.(*schema.Set).List() {
				attachmentParametersItemModel, err := resourceIbmSccProfileAttachmentMapToAttachmentParameterPrototype(attachmentParametersItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				attachmentParameters = append(attachmentParameters, *attachmentParametersItemModel)
			}
			replaceProfileAttachmentOptions.SetAttachmentParameters(attachmentParameters)
		}
		if replaceProfileAttachmentOptions.Status == nil {
			replaceProfileAttachmentOptions.SetSchedule(d.Get("status").(string))
		}
		if replaceProfileAttachmentOptions.Scope == nil || d.Get("scope") != nil {
			scopeItems := d.Get("scope")
			scopes := []securityandcompliancecenterapiv3.MultiCloudScopePayload{}
			for _, scopeItem := range scopeItems.(*schema.Set).List() {
				scopeItemModel, err := resourceIbmSccProfileAttachmentMapToMultiCloudScope(scopeItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				scopes = append(scopes, *scopeItemModel)
			}
			replaceProfileAttachmentOptions.SetScope(scopes)
		}
		_, response, err := securityandcompliancecenterapiClient.ReplaceProfileAttachmentWithContext(context, replaceProfileAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceProfileAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(flex.FmtErrorf("ReplaceProfileAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSccProfileAttachmentRead(context, d, meta)
}

func resourceIbmSccProfileAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProfileAttachmentOptions := &securityandcompliancecenterapiv3.DeleteProfileAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProfileAttachmentOptions.SetInstanceID(parts[0])
	deleteProfileAttachmentOptions.SetProfileID(parts[1])
	deleteProfileAttachmentOptions.SetAttachmentID(parts[2])

	_, response, err := securityandcompliancecenterapiClient.DeleteProfileAttachmentWithContext(context, deleteProfileAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProfileAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("DeleteProfileAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSccProfileAttachmentMapToAttachmentsPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ProfileAttachmentBase, error) {
	model := &securityandcompliancecenterapiv3.ProfileAttachmentBase{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	scope := []securityandcompliancecenterapiv3.MultiCloudScopePayload{}
	for _, scopeItem := range modelMap["scope"].(*schema.Set).List() {
		scopeItemModel, err := resourceIbmSccProfileAttachmentMapToMultiCloudScope(scopeItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		scope = append(scope, *scopeItemModel)
	}
	model.Scope = scope
	model.Status = core.StringPtr(modelMap["status"].(string))
	model.Schedule = core.StringPtr(modelMap["schedule"].(string))
	if modelMap["notifications"] != nil && len(modelMap["notifications"].([]interface{})) > 0 {
		NotificationsModel, err := resourceIbmSccProfileAttachmentMapToAttachmentsNotificationsPrototype(modelMap["notifications"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Notifications = NotificationsModel
	}
	attachmentParameters := []securityandcompliancecenterapiv3.Parameter{}
	if modelMap["attachment_parameters"] != nil {
		for _, attachmentParametersItem := range modelMap["attachment_parameters"].(*schema.Set).List() {
			if attachmentParametersItem != nil {
				attachmentParametersItemModel, err := resourceIbmSccProfileAttachmentMapToAttachmentParameterPrototype(attachmentParametersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				attachmentParameters = append(attachmentParameters, *attachmentParametersItemModel)
			}
		}
	}
	model.AttachmentParameters = attachmentParameters
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToMultiCloudScope(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.MultiCloudScopePayload, error) {
	model := &securityandcompliancecenterapiv3.MultiCloudScopePayload{}
	properties := []securityandcompliancecenterapiv3.ScopePropertyIntf{}
	// Manual Change: If the ID is found from a previous call, return only the ID
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
		return model, nil
	}
	for _, propertiesItem := range modelMap["properties"].([]interface{}) {
		propertiesItemModel, err := resourceIbmSccProfileAttachmentMapToPropertyItem(propertiesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		properties = append(properties, propertiesItemModel)
	}
	model.Properties = properties
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToPropertyItem(modelMap map[string]interface{}) (securityandcompliancecenterapiv3.ScopePropertyIntf, error) {
	model := &securityandcompliancecenterapiv3.ScopeProperty{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToAttachmentsNotificationsPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.AttachmentNotifications, error) {
	model := &securityandcompliancecenterapiv3.AttachmentNotifications{}
	model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	ControlsModel, err := resourceIbmSccProfileAttachmentMapToFailedControls(modelMap["controls"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Controls = ControlsModel
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToFailedControls(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.AttachmentNotificationsControls, error) {
	model := &securityandcompliancecenterapiv3.AttachmentNotificationsControls{}
	if modelMap["threshold_limit"] != nil {
		model.ThresholdLimit = core.Int64Ptr(int64(modelMap["threshold_limit"].(int)))
	}
	if modelMap["failed_control_ids"] != nil {
		failedControlIds := []string{}
		for _, failedControlIdsItem := range modelMap["failed_control_ids"].([]interface{}) {
			failedControlIds = append(failedControlIds, failedControlIdsItem.(string))
		}
		model.FailedControlIds = failedControlIds
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToAttachmentParameterPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Parameter, error) {
	model := &securityandcompliancecenterapiv3.Parameter{}
	if modelMap["assessment_type"] != nil && modelMap["assessment_type"].(string) != "" {
		model.AssessmentType = core.StringPtr(modelMap["assessment_type"].(string))
	}
	if modelMap["assessment_id"] != nil && modelMap["assessment_id"].(string) != "" {
		model.AssessmentID = core.StringPtr(modelMap["assessment_id"].(string))
	}
	if modelMap["parameter_name"] != nil && modelMap["parameter_name"].(string) != "" {
		model.ParameterName = core.StringPtr(modelMap["parameter_name"].(string))
	}
	if modelMap["parameter_value"] != nil && modelMap["parameter_value"].(string) != "" {
		model.ParameterValue = core.StringPtr(modelMap["parameter_value"].(string))
	}
	if modelMap["parameter_display_name"] != nil && modelMap["parameter_display_name"].(string) != "" {
		model.ParameterDisplayName = core.StringPtr(modelMap["parameter_display_name"].(string))
	}
	if modelMap["parameter_type"] != nil && modelMap["parameter_type"].(string) != "" {
		model.ParameterType = core.StringPtr(modelMap["parameter_type"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToAttachmentItem(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ProfileAttachment, error) {
	model := &securityandcompliancecenterapiv3.ProfileAttachment{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["profile_id"] != nil && modelMap["profile_id"].(string) != "" {
		model.ProfileID = core.StringPtr(modelMap["profile_id"].(string))
	}
	if modelMap["account_id"] != nil && modelMap["account_id"].(string) != "" {
		model.AccountID = core.StringPtr(modelMap["account_id"].(string))
	}
	if modelMap["instance_id"] != nil && modelMap["instance_id"].(string) != "" {
		model.InstanceID = core.StringPtr(modelMap["instance_id"].(string))
	}
	if modelMap["scope"] != nil {
		scope := []securityandcompliancecenterapiv3.MultiCloudScopePayload{}
		for _, scopeItem := range modelMap["scope"].([]interface{}) {
			scopeItemModel, err := resourceIbmSccProfileAttachmentMapToMultiCloudScope(scopeItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			scope = append(scope, *scopeItemModel)
		}
		model.Scope = scope
	}
	if modelMap["created_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["created_on"].(string))
		if err != nil {
			return model, err
		}
		model.CreatedOn = &dateTime
	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["updated_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["updated_on"].(string))
		if err != nil {
			return model, err
		}
		model.UpdatedOn = &dateTime
	}
	if modelMap["updated_by"] != nil && modelMap["updated_by"].(string) != "" {
		model.UpdatedBy = core.StringPtr(modelMap["updated_by"].(string))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["schedule"] != nil && modelMap["schedule"].(string) != "" {
		model.Schedule = core.StringPtr(modelMap["schedule"].(string))
	}
	if modelMap["notifications"] != nil && len(modelMap["notifications"].([]interface{})) > 0 {
		NotificationsModel, err := resourceIbmSccProfileAttachmentMapToAttachmentsNotificationsPrototype(modelMap["notifications"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Notifications = NotificationsModel
	}
	if modelMap["attachment_parameters"] != nil {
		attachmentParameters := []securityandcompliancecenterapiv3.Parameter{}
		for _, attachmentParametersItem := range modelMap["attachment_parameters"].([]interface{}) {
			attachmentParametersItemModel, err := resourceIbmSccProfileAttachmentMapToAttachmentParameterPrototype(attachmentParametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			attachmentParameters = append(attachmentParameters, *attachmentParametersItemModel)
		}
		model.AttachmentParameters = attachmentParameters
	}
	if modelMap["last_scan"] != nil && len(modelMap["last_scan"].([]interface{})) > 0 {
		LastScanModel, err := resourceIbmSccProfileAttachmentMapToLastScan(modelMap["last_scan"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LastScan = LastScanModel
	}
	if modelMap["next_scan_time"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["next_scan_time"].(string))
		if err != nil {
			return model, err
		}
		model.NextScanTime = &dateTime
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToLastScan(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.LastScan, error) {
	model := &securityandcompliancecenterapiv3.LastScan{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["time"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["time"].(string))
		if err != nil {
			return model, err
		}
		model.Time = &dateTime
	}
	return model, nil
}

func resourceIbmSccProfileAttachmentMapToAttachmentPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.CreateProfileAttachmentOptions, error) {
	model := &securityandcompliancecenterapiv3.CreateProfileAttachmentOptions{}
	model.SetInstanceID(modelMap["instance_id"].(string))
	if modelMap["profile_id"] != nil && modelMap["profile_id"].(string) != "" {
		model.ProfileID = core.StringPtr(modelMap["profile_id"].(string))
	}
	attachments := []securityandcompliancecenterapiv3.ProfileAttachmentBase{}
	attachmentsItemModel, err := resourceIbmSccProfileAttachmentMapToAttachmentsPrototype(modelMap)
	if err != nil {
		return model, err
	}
	attachments = append(attachments, *attachmentsItemModel)
	model.Attachments = attachments
	return model, nil
}

func resourceIbmSccProfileAttachmentMultiCloudScopeToMap(model *securityandcompliancecenterapiv3.MultiCloudScopePayload) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	properties := []map[string]interface{}{}
	for _, propertiesItem := range model.Properties {
		propertiesItemMap, err := scopePropertiesToMap(propertiesItem)
		if err != nil {
			return modelMap, err
		}
		properties = append(properties, propertiesItemMap)
	}
	modelMap["properties"] = properties
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func resourceIbmSccProfileAttachmentAttachmentsNotificationsPrototypeToMap(model *securityandcompliancecenterapiv3.AttachmentNotifications) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["enabled"] = model.Enabled
	controlsMap, err := resourceIbmSccProfileAttachmentFailedControlsToMap(model.Controls)
	if err != nil {
		return modelMap, err
	}
	modelMap["controls"] = []map[string]interface{}{controlsMap}
	return modelMap, nil
}

func resourceIbmSccProfileAttachmentFailedControlsToMap(model *securityandcompliancecenterapiv3.AttachmentNotificationsControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ThresholdLimit != nil {
		modelMap["threshold_limit"] = flex.IntValue(model.ThresholdLimit)
	}
	if model.FailedControlIds != nil {
		modelMap["failed_control_ids"] = model.FailedControlIds
	}
	return modelMap, nil
}

func resourceIbmSccProfileAttachmentAttachmentParameterPrototypeToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AssessmentType != nil {
		modelMap["assessment_type"] = flex.StringValue(model.AssessmentType)
	}
	if model.AssessmentID != nil {
		modelMap["assessment_id"] = flex.StringValue(model.AssessmentID)
	}
	if model.ParameterName != nil {
		modelMap["parameter_name"] = flex.StringValue(model.ParameterName)
	}
	if model.ParameterDisplayName != nil {
		modelMap["parameter_display_name"] = flex.StringValue(model.ParameterDisplayName)
	}
	if model.ParameterType != nil {
		modelMap["parameter_type"] = flex.StringValue(model.ParameterType)
	}
	if model.ParameterValue != nil {
		val, ok := model.ParameterValue.(string)
		if ok {
			modelMap["parameter_value"] = val
		}
	}
	return modelMap, nil
}

func resourceIbmSccProfileAttachmentLastScanToMap(model *securityandcompliancecenterapiv3.LastScan) (map[string]interface{}, error) {
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
