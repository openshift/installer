// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsPolicyCreate,
		ReadContext:   resourceIbmLogsPolicyRead,
		UpdateContext: resourceIbmLogsPolicyUpdate,
		DeleteContext: resourceIbmLogsPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_policy", "name"),
				Description:  "Name of policy.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_policy", "description"),
				Description:  "Description of policy.",
			},
			"priority": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_policy", "priority"),
				Description:  "The data pipeline sources that match the policy rules will go through.",
			},
			"application_rule": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Rule for matching with application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Identifier of the rule.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the rule.",
						},
					},
				},
			},
			"subsystem_rule": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Rule for matching with application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Identifier of the rule.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the rule.",
						},
					},
				},
			},
			"archive_retention": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Archive retention definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "References archive retention definition.",
						},
					},
				},
			},
			"log_rules": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Log rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severities": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Source severities to match with.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"company_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Company ID.",
			},
			"deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Soft deletion flag.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enabled flag.",
			},
			"order": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Order of policy in relation to other policies.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created at date at utc+0.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updated at date at utc+0.",
			},
			"policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy Id.",
			},
		},
	}
}

func ResourceIbmLogsPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9_\-\s]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "priority",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "type_block, type_high, type_low, type_medium, type_unspecified",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	bodyModelMap := map[string]interface{}{}
	createPolicyOptions := &logsv0.CreatePolicyOptions{}

	bodyModelMap["name"] = d.Get("name")
	if _, ok := d.GetOk("description"); ok {
		bodyModelMap["description"] = d.Get("description")
	}
	bodyModelMap["priority"] = d.Get("priority")
	if _, ok := d.GetOk("application_rule"); ok {
		bodyModelMap["application_rule"] = d.Get("application_rule")
	}
	if _, ok := d.GetOk("subsystem_rule"); ok {
		bodyModelMap["subsystem_rule"] = d.Get("subsystem_rule")
	}
	if _, ok := d.GetOk("archive_retention"); ok {
		bodyModelMap["archive_retention"] = d.Get("archive_retention")
	}
	if _, ok := d.GetOk("log_rules"); ok {
		bodyModelMap["log_rules"] = d.Get("log_rules")
	}
	convertedModel, err := ResourceIbmLogsPolicyMapToPolicyPrototype(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "create")
		return tfErr.GetDiag()
	}
	createPolicyOptions.PolicyPrototype = convertedModel

	policyIntf, _, err := logsClient.CreatePolicyWithContext(context, createPolicyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreatePolicyWithContext failed: %s", err.Error()), "ibm_logs_policy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	policy := policyIntf.(*logsv0.Policy)

	policyId := fmt.Sprintf("%s/%s/%s", region, instanceId, *policy.ID)
	d.SetId(policyId)

	return resourceIbmLogsPolicyRead(context, d, meta)
}

func resourceIbmLogsPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, policyId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getPolicyOptions := &logsv0.GetPolicyOptions{}

	getPolicyOptions.SetID(core.UUIDPtr(strfmt.UUID(policyId)))

	policyIntf, response, err := logsClient.GetPolicyWithContext(context, getPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetPolicyWithContext failed: %s", err.Error()), "ibm_logs_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	policy := policyIntf.(*logsv0.Policy)

	if err = d.Set("policy_id", policyId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting policy_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("name", policy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(policy.Description) {
		if err = d.Set("description", policy.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if err = d.Set("priority", policy.Priority); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting priority: %s", err))
	}
	if !core.IsNil(policy.ApplicationRule) {
		applicationRuleMap, err := ResourceIbmLogsPolicyQuotaV1RuleToMap(policy.ApplicationRule)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("application_rule", []map[string]interface{}{applicationRuleMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting application_rule: %s", err))
		}
	}
	if !core.IsNil(policy.SubsystemRule) {
		subsystemRuleMap, err := ResourceIbmLogsPolicyQuotaV1RuleToMap(policy.SubsystemRule)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("subsystem_rule", []map[string]interface{}{subsystemRuleMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subsystem_rule: %s", err))
		}
	}
	if !core.IsNil(policy.ArchiveRetention) {
		archiveRetentionMap, err := ResourceIbmLogsPolicyQuotaV1ArchiveRetentionToMap(policy.ArchiveRetention)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("archive_retention", []map[string]interface{}{archiveRetentionMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting archive_retention: %s", err))
		}
	}
	if !core.IsNil(policy.LogRules) {
		logRulesMap, err := ResourceIbmLogsPolicyQuotaV1LogRulesToMap(policy.LogRules)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("log_rules", []map[string]interface{}{logRulesMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting log_rules: %s", err))
		}
	}
	if !core.IsNil(policy.CompanyID) {
		if err = d.Set("company_id", flex.IntValue(policy.CompanyID)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting company_id: %s", err))
		}
	}
	if !core.IsNil(policy.Deleted) {
		if err = d.Set("deleted", policy.Deleted); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting deleted: %s", err))
		}
	}
	if !core.IsNil(policy.Enabled) {
		if err = d.Set("enabled", policy.Enabled); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
		}
	}
	if !core.IsNil(policy.Order) {
		if err = d.Set("order", flex.IntValue(policy.Order)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting order: %s", err))
		}
	}
	if !core.IsNil(policy.CreatedAt) {
		if err = d.Set("created_at", policy.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if !core.IsNil(policy.UpdatedAt) {
		if err = d.Set("updated_at", policy.UpdatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
	}

	return nil
}

func resourceIbmLogsPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, policyId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePolicyOptions := &logsv0.UpdatePolicyOptions{}

	updatePolicyOptions.SetID(core.UUIDPtr(strfmt.UUID(policyId)))

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("priority") ||
		d.HasChange("application_rule") ||
		d.HasChange("subsystem_rule") ||
		d.HasChange("archive_retention") ||
		d.HasChange("log_rules") {

		bodyModelMap := map[string]interface{}{}

		bodyModelMap["name"] = d.Get("name")

		if _, ok := d.GetOk("description"); ok {
			bodyModelMap["description"] = d.Get("description")
		}
		bodyModelMap["priority"] = d.Get("priority")
		if _, ok := d.GetOk("application_rule"); ok {
			bodyModelMap["application_rule"] = d.Get("application_rule")
		}
		if _, ok := d.GetOk("subsystem_rule"); ok {
			bodyModelMap["subsystem_rule"] = d.Get("subsystem_rule")
		}
		if _, ok := d.GetOk("archive_retention"); ok {
			bodyModelMap["archive_retention"] = d.Get("archive_retention")
		}
		if _, ok := d.GetOk("log_rules"); ok {
			bodyModelMap["log_rules"] = d.Get("log_rules")
		}
		convertedModel, err := ResourceIbmLogsPolicyMapToPolicyPrototype(bodyModelMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "create")
			return tfErr.GetDiag()
		}
		updatePolicyOptions.PolicyPrototype = convertedModel

		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdatePolicyWithContext(context, updatePolicyOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdatePolicyWithContext failed: %s", err.Error()), "ibm_logs_policy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsPolicyRead(context, d, meta)
}

func resourceIbmLogsPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_policy", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, policyId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deletePolicyOptions := &logsv0.DeletePolicyOptions{}

	deletePolicyOptions.SetID(core.UUIDPtr(strfmt.UUID(policyId)))

	_, err = logsClient.DeletePolicyWithContext(context, deletePolicyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeletePolicyWithContext failed: %s", err.Error()), "ibm_logs_policy", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsPolicyMapToQuotaV1Rule(modelMap map[string]interface{}) (*logsv0.QuotaV1Rule, error) {
	model := &logsv0.QuotaV1Rule{}
	model.RuleTypeID = core.StringPtr(modelMap["rule_type_id"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIbmLogsPolicyMapToQuotaV1ArchiveRetention(modelMap map[string]interface{}) (*logsv0.QuotaV1ArchiveRetention, error) {
	model := &logsv0.QuotaV1ArchiveRetention{}
	model.ID = core.UUIDPtr(strfmt.UUID(modelMap["id"].(string)))
	return model, nil
}

func ResourceIbmLogsPolicyMapToQuotaV1LogRules(modelMap map[string]interface{}) (*logsv0.QuotaV1LogRules, error) {
	model := &logsv0.QuotaV1LogRules{}
	if modelMap["severities"] != nil {
		severities := []string{}
		for _, severitiesItem := range modelMap["severities"].([]interface{}) {
			severities = append(severities, severitiesItem.(string))
		}
		model.Severities = severities
	}
	return model, nil
}

func ResourceIbmLogsPolicyMapToPolicyPrototype(modelMap map[string]interface{}) (logsv0.PolicyPrototypeIntf, error) {
	model := &logsv0.PolicyPrototype{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Priority = core.StringPtr(modelMap["priority"].(string))
	if modelMap["application_rule"] != nil && len(modelMap["application_rule"].([]interface{})) > 0 {
		ApplicationRuleModel, err := ResourceIbmLogsPolicyMapToQuotaV1Rule(modelMap["application_rule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ApplicationRule = ApplicationRuleModel
	}
	if modelMap["subsystem_rule"] != nil && len(modelMap["subsystem_rule"].([]interface{})) > 0 {
		SubsystemRuleModel, err := ResourceIbmLogsPolicyMapToQuotaV1Rule(modelMap["subsystem_rule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SubsystemRule = SubsystemRuleModel
	}
	if modelMap["archive_retention"] != nil && len(modelMap["archive_retention"].([]interface{})) > 0 {
		ArchiveRetentionModel, err := ResourceIbmLogsPolicyMapToQuotaV1ArchiveRetention(modelMap["archive_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchiveRetention = ArchiveRetentionModel
	}
	if modelMap["log_rules"] != nil && len(modelMap["log_rules"].([]interface{})) > 0 {
		LogRulesModel, err := ResourceIbmLogsPolicyMapToQuotaV1LogRules(modelMap["log_rules"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRules = LogRulesModel
	}
	return model, nil
}

func ResourceIbmLogsPolicyMapToPolicyPrototypeQuotaV1CreatePolicyRequestSourceTypeRulesLogRules(modelMap map[string]interface{}) (*logsv0.PolicyPrototypeQuotaV1CreatePolicyRequestSourceTypeRulesLogRules, error) {
	model := &logsv0.PolicyPrototypeQuotaV1CreatePolicyRequestSourceTypeRulesLogRules{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Priority = core.StringPtr(modelMap["priority"].(string))
	if modelMap["application_rule"] != nil && len(modelMap["application_rule"].([]interface{})) > 0 {
		ApplicationRuleModel, err := ResourceIbmLogsPolicyMapToQuotaV1Rule(modelMap["application_rule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ApplicationRule = ApplicationRuleModel
	}
	if modelMap["subsystem_rule"] != nil && len(modelMap["subsystem_rule"].([]interface{})) > 0 {
		SubsystemRuleModel, err := ResourceIbmLogsPolicyMapToQuotaV1Rule(modelMap["subsystem_rule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SubsystemRule = SubsystemRuleModel
	}
	if modelMap["archive_retention"] != nil && len(modelMap["archive_retention"].([]interface{})) > 0 {
		ArchiveRetentionModel, err := ResourceIbmLogsPolicyMapToQuotaV1ArchiveRetention(modelMap["archive_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchiveRetention = ArchiveRetentionModel
	}
	if modelMap["log_rules"] != nil && len(modelMap["log_rules"].([]interface{})) > 0 {
		LogRulesModel, err := ResourceIbmLogsPolicyMapToQuotaV1LogRules(modelMap["log_rules"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRules = LogRulesModel
	}
	return model, nil
}

func ResourceIbmLogsPolicyQuotaV1RuleToMap(model *logsv0.QuotaV1Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["rule_type_id"] = *model.RuleTypeID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIbmLogsPolicyQuotaV1ArchiveRetentionToMap(model *logsv0.QuotaV1ArchiveRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	return modelMap, nil
}

func ResourceIbmLogsPolicyQuotaV1LogRulesToMap(model *logsv0.QuotaV1LogRules) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	return modelMap, nil
}
