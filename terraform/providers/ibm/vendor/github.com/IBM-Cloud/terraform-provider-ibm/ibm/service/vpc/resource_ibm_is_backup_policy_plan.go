// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsBackupPolicyPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBackupPolicyPlanCreate,
		ReadContext:   resourceIBMIsBackupPolicyPlanRead,
		UpdateContext: resourceIBMIsBackupPolicyPlanUpdate,
		DeleteContext: resourceIBMIsBackupPolicyPlanDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The backup policy identifier.",
			},
			"backup_policy_plan_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The backup policy identifier.",
			},
			"cron_spec": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_backup_policy_plan", "cron_spec"),
				Description:  "The cron specification for the backup schedule.",
			},
			"active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the plan is active.",
			},
			"attach_user_tags": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         schema.HashString,
				Description: "User tags to attach to each backup (snapshot) created by this plan. If unspecified, no user tags will be attached.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"copy_user_tags": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether to copy the source's user tags to the created backups (snapshots).",
			},
			"deletion_trigger": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_after": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "The maximum number of days to keep each backup after creation.",
						},
						"delete_over_count": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The maximum number of recent backups to keep. If unspecified, there will be no maximum.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_backup_policy_plan", "name"),
				Description:  "The user-defined name for this backup policy plan. Names must be unique within the backup policy this plan resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy plan was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy plan.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this backup policy plan.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the BackupPolicyPlan.",
			},
		},
	}
}

func ResourceIBMIsBackupPolicyPlanValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cron_spec",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^((((\d+,)+\d+|([\d\*]+(\/|-)\d+)|\d+|\*) ?){5,7})$`,
			MinValueLength:             9,
			MaxValueLength:             63,
		},
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_backup_policy_plan", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBackupPolicyPlanCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createBackupPolicyPlanOptions := &vpcv1.CreateBackupPolicyPlanOptions{}

	createBackupPolicyPlanOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))
	createBackupPolicyPlanOptions.SetCronSpec(d.Get("cron_spec").(string))
	if _, ok := d.GetOk("active"); ok {
		createBackupPolicyPlanOptions.SetActive(d.Get("active").(bool))
	}
	if _, ok := d.GetOk("attach_user_tags"); ok {
		createBackupPolicyPlanOptions.SetAttachUserTags((flex.ExpandStringList((d.Get("attach_user_tags").(*schema.Set)).List())))
	}
	if _, ok := d.GetOk("copy_user_tags"); ok {
		createBackupPolicyPlanOptions.SetCopyUserTags(d.Get("copy_user_tags").(bool))
	}
	if _, ok := d.GetOk("deletion_trigger"); ok {
		backupPolicyPlanDeletionTriggerPrototypeMap := d.Get("deletion_trigger.0").(map[string]interface{})
		backupPolicyPlanDeletionTriggerPrototype := vpcv1.BackupPolicyPlanDeletionTriggerPrototype{}

		if backupPolicyPlanDeletionTriggerPrototypeMap["delete_after"] != nil {
			backupPolicyPlanDeletionTriggerPrototype.DeleteAfter = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerPrototypeMap["delete_after"].(int)))
		}
		if backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"] != nil {
			deleteOverCountString := backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"].(string)
			if deleteOverCountString != "" && deleteOverCountString != "null" {
				deleteOverCount, err := strconv.ParseInt(backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"].(string), 10, 64)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error setting delete_over_count: %s", err))
				}
				deleteOverCountint := int64(deleteOverCount)
				if deleteOverCountint >= int64(0) {
					backupPolicyPlanDeletionTriggerPrototype.DeleteOverCount = core.Int64Ptr(deleteOverCountint)
				} else {
					return diag.FromErr(fmt.Errorf("[ERROR] Error setting delete_over_count: Retention count and days cannot be both zero"))
				}
			}
		}
		createBackupPolicyPlanOptions.SetDeletionTrigger(&backupPolicyPlanDeletionTriggerPrototype)
	}
	if _, ok := d.GetOk("name"); ok {
		createBackupPolicyPlanOptions.SetName(d.Get("name").(string))
	}

	backupPolicyPlan, response, err := vpcClient.CreateBackupPolicyPlanWithContext(context, createBackupPolicyPlanOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] CreateBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createBackupPolicyPlanOptions.BackupPolicyID, *backupPolicyPlan.ID))

	return resourceIBMIsBackupPolicyPlanRead(context, d, meta)
}

func resourceIBMIsBackupPolicyPlanRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyPlanOptions := &vpcv1.GetBackupPolicyPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	getBackupPolicyPlanOptions.SetID(parts[1])

	backupPolicyPlan, response, err := vpcClient.GetBackupPolicyPlanWithContext(context, getBackupPolicyPlanOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	if getBackupPolicyPlanOptions.BackupPolicyID != nil {
		if err = d.Set("backup_policy_id", getBackupPolicyPlanOptions.BackupPolicyID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting backup_policy_id: %s", err))
		}
	}

	if getBackupPolicyPlanOptions.ID != nil {
		if err = d.Set("backup_policy_plan_id", getBackupPolicyPlanOptions.ID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting backup_policy_plan_id: %s", err))
		}
	}

	if backupPolicyPlan.CronSpec != nil {
		if err = d.Set("cron_spec", backupPolicyPlan.CronSpec); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting cron_spec: %s", err))
		}
	}

	if backupPolicyPlan.Active != nil {
		if err = d.Set("active", backupPolicyPlan.Active); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting active: %s", err))
		}
	}

	if backupPolicyPlan.AttachUserTags != nil {
		if err = d.Set("attach_user_tags", backupPolicyPlan.AttachUserTags); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting attach_user_tags: %s", err))
		}
	}

	if backupPolicyPlan.CopyUserTags != nil {
		if err = d.Set("copy_user_tags", backupPolicyPlan.CopyUserTags); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting copy_user_tags: %s", err))
		}
	}

	if backupPolicyPlan.DeletionTrigger != nil {
		deletionTriggerMap := resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerPrototypeToMap(*backupPolicyPlan.DeletionTrigger)
		if err = d.Set("deletion_trigger", []map[string]interface{}{deletionTriggerMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting deletion_trigger: %s", err))
		}
	}

	if backupPolicyPlan.Name != nil {
		if err = d.Set("name", backupPolicyPlan.Name); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
		}
	}

	if backupPolicyPlan.CreatedAt != nil {
		if err = d.Set("created_at", flex.DateTimeToString(backupPolicyPlan.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
	}

	if err = d.Set("href", backupPolicyPlan.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", backupPolicyPlan.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", backupPolicyPlan.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting version: %s", err))
	}

	return nil
}

func resourceIBMIsBackupPolicyPlanBackupPolicyPlanDeletionTriggerPrototypeToMap(backupPolicyPlanDeletionTriggerPrototype vpcv1.BackupPolicyPlanDeletionTrigger) map[string]interface{} {
	backupPolicyPlanDeletionTriggerPrototypeMap := map[string]interface{}{}

	if backupPolicyPlanDeletionTriggerPrototype.DeleteAfter != nil {
		backupPolicyPlanDeletionTriggerPrototypeMap["delete_after"] = flex.IntValue(backupPolicyPlanDeletionTriggerPrototype.DeleteAfter)
	}
	if backupPolicyPlanDeletionTriggerPrototype.DeleteOverCount != nil {
		backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"] = strconv.FormatInt(*backupPolicyPlanDeletionTriggerPrototype.DeleteOverCount, 10)
	} else {
		backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"] = "null"
	}

	return backupPolicyPlanDeletionTriggerPrototypeMap
}

func resourceIBMIsBackupPolicyPlanUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBackupPolicyPlanOptions := &vpcv1.UpdateBackupPolicyPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	updateBackupPolicyPlanOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.BackupPolicyPlanPatch{}
	if d.HasChange("cron_spec") {
		patchVals.CronSpec = core.StringPtr(d.Get("cron_spec").(string))
		hasChange = true
	}
	if d.HasChange("active") {
		patchVals.Active = core.BoolPtr(d.Get("active").(bool))
		hasChange = true
	}
	if d.HasChange("attach_user_tags") {
		patchVals.AttachUserTags = (flex.ExpandStringList((d.Get("attach_user_tags").(*schema.Set)).List()))
		hasChange = true
	}
	if d.HasChange("copy_user_tags") {
		patchVals.CopyUserTags = core.BoolPtr(d.Get("copy_user_tags").(bool))
		hasChange = true
	}

	deleteOverCountBool := false
	if d.HasChange("deletion_trigger") {
		backupPolicyPlanDeletionTriggerPrototype := vpcv1.BackupPolicyPlanDeletionTriggerPatch{}
		backupPolicyPlanDeletionTriggerPrototypeMap := d.Get("deletion_trigger.0").(map[string]interface{})
		if backupPolicyPlanDeletionTriggerPrototypeMap["delete_after"] != nil {
			backupPolicyPlanDeletionTriggerPrototype.DeleteAfter = core.Int64Ptr(int64(backupPolicyPlanDeletionTriggerPrototypeMap["delete_after"].(int)))
		}
		if backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"] != nil {
			deleteOverCountString := backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"].(string)
			if deleteOverCountString != "" && deleteOverCountString != "null" {
				deleteOverCount, err := strconv.ParseInt(backupPolicyPlanDeletionTriggerPrototypeMap["delete_over_count"].(string), 10, 64)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error setting delete_over_count: %s", err))
				}
				deleteOverCountint := int64(deleteOverCount)
				if deleteOverCountint >= int64(0) {
					backupPolicyPlanDeletionTriggerPrototype.DeleteOverCount = core.Int64Ptr(deleteOverCountint)
				}
			} else {
				deleteOverCountBool = true
			}
		}
		patchVals.DeletionTrigger = &backupPolicyPlanDeletionTriggerPrototype
		hasChange = true
	}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}
	updateBackupPolicyPlanOptions.SetIfMatch(d.Get("version").(string))

	if hasChange {
		backupPolicyPlanPatch, err := patchVals.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] [ERROR] Error calling asPatch for BackupPolicyPlanPatch: %s", err))
		}

		if deleteOverCountBool {
			backupPolicyPlanDeletionTriggerMap := backupPolicyPlanPatch["deletion_trigger"]
			backupPolicyPlanDeletionTrigger := backupPolicyPlanDeletionTriggerMap.(map[string]interface{})
			backupPolicyPlanDeletionTrigger["delete_over_count"] = nil
			backupPolicyPlanPatch["deletion_trigger"] = backupPolicyPlanDeletionTrigger
		}

		updateBackupPolicyPlanOptions.BackupPolicyPlanPatch = backupPolicyPlanPatch
		_, response, err := vpcClient.UpdateBackupPolicyPlanWithContext(context, updateBackupPolicyPlanOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBackupPolicyPlanWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] UpdateBackupPolicyPlanWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsBackupPolicyPlanRead(context, d, meta)
}

func resourceIBMIsBackupPolicyPlanDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackupPolicyPlanOptions := &vpcv1.DeleteBackupPolicyPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackupPolicyPlanOptions.SetBackupPolicyID(parts[0])
	deleteBackupPolicyPlanOptions.SetID(parts[1])

	deleteBackupPolicyPlanOptions.SetIfMatch(d.Get("version").(string))

	_, response, err := vpcClient.DeleteBackupPolicyPlanWithContext(context, deleteBackupPolicyPlanOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBackupPolicyPlanWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] DeleteBackupPolicyPlanWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
