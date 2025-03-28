// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBackupPolicyCreate,
		ReadContext:   resourceIBMIsBackupPolicyRead,
		UpdateContext: resourceIBMIsBackupPolicyUpdate,
		DeleteContext: resourceIBMIsBackupPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"match_resource_types": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Set:           schema.HashString,
				Deprecated:    "match_resource_types is being deprecated. Use match_resource_type instead",
				Description:   "A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.",
				ConflictsWith: []string{"match_resource_type"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"match_resource_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "volume",
				ForceNew:      true,
				ConflictsWith: []string{"match_resource_types"},
				ValidateFunc:  validate.InvokeValidator("ibm_is_backup_policy", "match_resource_type"),
				Description:   "A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.",
			},
			"match_user_tags": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"included_content": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Set:         schema.HashString,
				Description: "The included content for backups created using this policy",
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_backup_policy", "included_content")},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_backup_policy", "name"),
				Description:  "The user-defined name for this backup policy. Names must be unique within the region this backup policy resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this backup policy.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy.",
			},
			"last_job_completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the most recent job for this backup policy completed.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the backup policy.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current health_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},
			"health_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource",
			},
			"scope": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope for this backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CRN for this enterprise.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this enterprise or account.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMIsBackupPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "match_user_tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "included_content",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "boot_volume, data_volumes",
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "match_resource_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "instance, volume, share",
		},
	)
	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_backup_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBackupPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createBackupPolicyOptions := &vpcv1.CreateBackupPolicyOptions{}
	backupPolicyPrototype := &vpcv1.BackupPolicyPrototype{}

	if matchResourceType, ok := d.GetOk("match_resource_type"); ok {
		backupPolicyPrototype.MatchResourceType = core.StringPtr(matchResourceType.(string))
	} else if matchResourceTypes, ok := d.GetOk("match_resource_types"); ok {
		matchResourceTypeList := flex.ExpandStringList((matchResourceTypes.(*schema.Set)).List())
		backupPolicyPrototype.MatchResourceType = core.StringPtr(matchResourceTypeList[0])
	}
	if _, ok := d.GetOk("included_content"); ok {
		backupPolicyPrototype.IncludedContent = flex.ExpandStringList((d.Get("included_content").(*schema.Set)).List())
	}

	if _, ok := d.GetOk("match_user_tags"); ok {
		backupPolicyPrototype.MatchUserTags = flex.ExpandStringList((d.Get("match_user_tags").(*schema.Set)).List())
	}
	if _, ok := d.GetOk("name"); ok {
		backupPolicyPrototype.Name = core.StringPtr(d.Get("name").(string))
	}
	if resGroup, ok := d.GetOk("resource_group"); ok {
		resourceGroupStr := resGroup.(string)
		resourceGroup := vpcv1.ResourceGroupIdentity{
			ID: &resourceGroupStr,
		}
		backupPolicyPrototype.ResourceGroup = &resourceGroup
	}

	if _, ok := d.GetOk("scope"); ok {
		bkpPolicyScopePrototypeMap := d.Get("scope.0").(map[string]interface{})
		bkpPolicyScopePrototype := vpcv1.BackupPolicyScopePrototype{}
		if bkpPolicyScopePrototypeMap["crn"] != nil {
			if crnStr := bkpPolicyScopePrototypeMap["crn"].(string); crnStr != "" {
				bkpPolicyScopePrototype.CRN = core.StringPtr(crnStr)
			}
		}
		backupPolicyPrototype.Scope = &bkpPolicyScopePrototype
	}
	createBackupPolicyOptions.SetBackupPolicyPrototype(backupPolicyPrototype)
	backupPolicyIntf, response, err := vpcClient.CreateBackupPolicyWithContext(context, createBackupPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] CreateBackupPolicyWithContext failed %s\n%s", err, response))
	}

	backupPolicy := backupPolicyIntf.(*vpcv1.BackupPolicy)
	d.SetId(*backupPolicy.ID)

	return resourceIBMIsBackupPolicyRead(context, d, meta)
}

func resourceIBMIsBackupPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyOptions := &vpcv1.GetBackupPolicyOptions{}

	getBackupPolicyOptions.SetID(d.Id())

	backupPolicyIntf, response, err := vpcClient.GetBackupPolicyWithContext(context, getBackupPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetBackupPolicyWithContext failed %s\n%s", err, response))
	}
	backupPolicy := backupPolicyIntf.(*vpcv1.BackupPolicy)

	if backupPolicy.MatchResourceType != nil {
		matchResourceTypes := *backupPolicy.MatchResourceType
		matchResourceTypesList := []string{matchResourceTypes}
		if err = d.Set("match_resource_types", matchResourceTypesList); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting match_resource_types: %s", err))
		}
		if err = d.Set("match_resource_type", backupPolicy.MatchResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting match_resource_type: %s", err))
		}
	}
	if backupPolicy.IncludedContent != nil {
		if err = d.Set("included_content", backupPolicy.IncludedContent); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting included_content: %s", err))
		}
	}

	if backupPolicy.MatchUserTags != nil {
		if err = d.Set("match_user_tags", backupPolicy.MatchUserTags); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting match_user_tags: %s", err))
		}
	}
	if backupPolicy.Name != nil {
		if err = d.Set("name", backupPolicy.Name); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
		}
	}
	if backupPolicy.ResourceGroup != nil {
		resourceGroupID := *backupPolicy.ResourceGroup.ID
		if err = d.Set("resource_group", resourceGroupID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
		}
	}
	if backupPolicy.CreatedAt != nil {
		if err = d.Set("created_at", flex.DateTimeToString(backupPolicy.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
	}

	if backupPolicy.CRN != nil {
		if err = d.Set("crn", backupPolicy.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
		}
	}

	if backupPolicy.Href != nil {
		if err = d.Set("href", backupPolicy.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
	}

	if backupPolicy.LastJobCompletedAt != nil {
		if err = d.Set("last_job_completed_at", flex.DateTimeToString(backupPolicy.LastJobCompletedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_job_completed_at: %s", err))
		}
	}

	if backupPolicy.LifecycleState != nil {
		if err = d.Set("lifecycle_state", backupPolicy.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
	}

	if backupPolicy.ResourceType != nil {
		if err = d.Set("resource_type", backupPolicy.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}
	}

	if backupPolicy.HealthReasons != nil {
		healthReasonsList := make([]map[string]interface{}, 0)
		for _, sr := range backupPolicy.HealthReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR["code"] = *sr.Code
				currentSR["message"] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR["more_info"] = *sr.Message
				}
				healthReasonsList = append(healthReasonsList, currentSR)
			}
		}
		d.Set("health_reasons", healthReasonsList)
	}
	if err = d.Set("health_state", backupPolicy.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_state: %s", err))
	}

	if backupPolicy.Scope != nil {
		scope := []map[string]interface{}{}
		scopeMap := resourceIbmIsBackupPolicyScopeToMap(*backupPolicy.Scope.(*vpcv1.BackupPolicyScope))
		scope = append(scope, scopeMap)

		if err = d.Set("scope", scope); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting scope: %s", err))
		}
	}

	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting version: %s", err))
	}

	return nil
}

func resourceIbmIsBackupPolicyScopeToMap(scope vpcv1.BackupPolicyScope) map[string]interface{} {
	scopeMap := map[string]interface{}{}

	scopeMap["crn"] = scope.CRN
	scopeMap["id"] = scope.ID
	scopeMap["resource_type"] = scope.ResourceType

	return scopeMap
}

func resourceIBMIsBackupPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	updateBackupPolicyOptions := &vpcv1.UpdateBackupPolicyOptions{}
	updateBackupPolicyOptions.SetID(d.Id())
	hasChange := false
	patchVals := &vpcv1.BackupPolicyPatch{}
	if d.HasChange("match_user_tags") {
		patchVals.MatchUserTags = (flex.ExpandStringList((d.Get("match_user_tags").(*schema.Set)).List()))
		hasChange = true
	}
	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("included_content") {
		patchVals.IncludedContent = (flex.ExpandStringList((d.Get("included_content").(*schema.Set)).List()))
		hasChange = true
	}
	updateBackupPolicyOptions.SetIfMatch(d.Get("version").(string))
	if hasChange {
		updateBackupPolicyOptions.BackupPolicyPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateBackupPolicyWithContext(context, updateBackupPolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBackupPolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] UpdateBackupPolicyWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsBackupPolicyRead(context, d, meta)
}

func resourceIBMIsBackupPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	deleteBackupPolicyOptions := &vpcv1.DeleteBackupPolicyOptions{}
	deleteBackupPolicyOptions.SetID(d.Id())
	deleteBackupPolicyOptions.SetIfMatch(d.Get("version").(string))
	_, response, err := vpcClient.DeleteBackupPolicyWithContext(context, deleteBackupPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBackupPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] DeleteBackupPolicyWithContext failed %s\n%s", err, response))
	}
	d.SetId("")
	return nil
}
