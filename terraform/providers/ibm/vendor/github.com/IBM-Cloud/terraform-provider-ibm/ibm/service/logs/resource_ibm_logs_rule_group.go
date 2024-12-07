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

func ResourceIbmLogsRuleGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsRuleGroupCreate,
		ReadContext:   resourceIbmLogsRuleGroupRead,
		UpdateContext: resourceIbmLogsRuleGroupUpdate,
		DeleteContext: resourceIbmLogsRuleGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_rule_group", "name"),
				Description:  "The name of the rule group.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_rule_group", "description"),
				Description:  "A description for the rule group, should express what is the rule group purpose.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the rule is enabled.",
			},
			"rule_matchers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "// Optional rule matchers which if matched will make the rule go through the rule group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_name": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "ApplicationName constraint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Only logs with this ApplicationName value will match.",
									},
								},
							},
						},
						"subsystem_name": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "SubsystemName constraint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Only logs with this SubsystemName value will match.",
									},
								},
							},
						},
						"severity": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Severity constraint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Only logs with this severity value will match.",
									},
								},
							},
						},
					},
				},
			},
			"rule_subgroups": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule subgroups. Will try to execute the first rule subgroup, and if not matched will try to match the next one in order.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the rule subgroup.",
						},
						"rules": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Rules to run on the log.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Unique identifier of the rule.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of the rule.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the rule.",
									},
									"source_field": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "A field on which value to execute the rule.",
									},
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Parameters for a rule which specifies how it should run.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"extract_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for text extraction rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"rule": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Regex which will parse the source field and extract the json keys from it while retaining the original log.",
															},
														},
													},
												},
												"json_extract_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for json extract rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_field": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "In which metadata field to store the extracted value.",
															},
														},
													},
												},
												"replace_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for replace rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "In which field to put the modified text.",
															},
															"replace_new_val": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The value to replace the matched text with.",
															},
															"rule": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Regex which will match parts in the text to replace.",
															},
														},
													},
												},
												"parse_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for parse rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "In which field to put the parsed text.",
															},
															"rule": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Regex which will parse the source field and extract the json keys from it while removing the source field.",
															},
														},
													},
												},
												"allow_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for allow rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_blocked_logs": &schema.Schema{
																Type:        schema.TypeBool,
																Required:    true,
																Description: "If true matched logs will be blocked, otherwise matched logs will be kept.",
															},
															"rule": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Regex which will match the source field and decide if the rule will apply.",
															},
														},
													},
												},
												"block_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for block rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_blocked_logs": &schema.Schema{
																Type:        schema.TypeBool,
																Required:    true,
																Description: "If true matched logs will be kept, otherwise matched logs will be blocked.",
															},
															"rule": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Regex which will match the source field and decide if the rule will apply.",
															},
														},
													},
												},
												"extract_timestamp_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for extract timestamp rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"standard": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "What time format to use on the extracted time.",
															},
															"format": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "What time format the the source field to extract from has.",
															},
														},
													},
												},
												"remove_fields_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for remove fields rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"fields": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Json field paths to drop from the log.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"json_stringify_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for json stringify rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Destination field in which to put the json stringified content.",
															},
															"delete_source": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Whether or not to delete the source field after running this rule.",
															},
														},
													},
												},
												"json_parse_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Parameters for json parse rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Destination field under which to put the json object.",
															},
															"delete_source": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Whether or not to delete the source field after running this rule.",
															},
															"override_dest": &schema.Schema{
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Destination field in which to put the json stringified content.",
															},
														},
													},
												},
											},
										},
									},
									"enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether or not to execute the rule.",
									},
									"order": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The ordering of the rule subgroup. Lower order will run first. 0 is considered as no value.",
									},
								},
							},
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the rule subgroup is enabled.",
						},
						"order": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The ordering of the rule subgroup. Lower order will run first. 0 is considered as no value.",
						},
					},
				},
			},
			"order": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "// The order in which the rule group will be evaluated. The lower the order, the more priority the group will have. Not providing the order will by default create a group with the last order.",
			},
			"rule_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule Group Id.",
			},
		},
	}
}

func ResourceIbmLogsRuleGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             255,
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_rule_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsRuleGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_rule_group", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createRuleGroupOptions := &logsv0.CreateRuleGroupOptions{}

	createRuleGroupOptions.SetName(d.Get("name").(string))
	var ruleSubgroups []logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroup
	for _, v := range d.Get("rule_subgroups").([]interface{}) {
		value := v.(map[string]interface{})
		ruleSubgroupsItem, err := ResourceIbmLogsRuleGroupMapToRulesV1CreateRuleGroupRequestCreateRuleSubgroup(value)
		if err != nil {
			return diag.FromErr(err)
		}
		ruleSubgroups = append(ruleSubgroups, *ruleSubgroupsItem)
	}
	createRuleGroupOptions.SetRuleSubgroups(ruleSubgroups)
	if _, ok := d.GetOk("description"); ok {
		createRuleGroupOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOkExists("enabled"); ok {
		createRuleGroupOptions.SetEnabled(d.Get("enabled").(bool))
	}
	if _, ok := d.GetOk("rule_matchers"); ok {
		var ruleMatchers []logsv0.RulesV1RuleMatcherIntf
		for _, v := range d.Get("rule_matchers").([]interface{}) {
			value := v.(map[string]interface{})
			ruleMatchersItem, err := ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcher(value)
			if err != nil {
				return diag.FromErr(err)
			}
			ruleMatchers = append(ruleMatchers, ruleMatchersItem)
		}
		createRuleGroupOptions.SetRuleMatchers(ruleMatchers)
	}
	if _, ok := d.GetOk("order"); ok {
		createRuleGroupOptions.SetOrder(int64(d.Get("order").(int)))
	}

	ruleGroup, _, err := logsClient.CreateRuleGroupWithContext(context, createRuleGroupOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateRuleGroupWithContext failed: %s", err.Error()), "ibm_logs_rule_group", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	ruleGroupId := fmt.Sprintf("%s/%s/%s", region, instanceId, *ruleGroup.ID)
	d.SetId(ruleGroupId)

	return resourceIbmLogsRuleGroupRead(context, d, meta)
}

func resourceIbmLogsRuleGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_rule_group", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, ruleGroupId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleGroupOptions := &logsv0.GetRuleGroupOptions{}

	getRuleGroupOptions.SetGroupID(core.UUIDPtr(strfmt.UUID(ruleGroupId)))

	ruleGroup, response, err := logsClient.GetRuleGroupWithContext(context, getRuleGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRuleGroupWithContext failed: %s", err.Error()), "ibm_logs_rule_group", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("rule_group_id", ruleGroupId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rule_group_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("name", ruleGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(ruleGroup.Description) {
		if err = d.Set("description", ruleGroup.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(ruleGroup.Enabled) {
		if err = d.Set("enabled", ruleGroup.Enabled); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
		}
	}
	if !core.IsNil(ruleGroup.RuleMatchers) {
		ruleMatchers := []map[string]interface{}{}
		for _, ruleMatchersItem := range ruleGroup.RuleMatchers {
			ruleMatchersItemMap, err := ResourceIbmLogsRuleGroupRulesV1RuleMatcherToMap(ruleMatchersItem)
			if err != nil {
				return diag.FromErr(err)
			}
			ruleMatchers = append(ruleMatchers, ruleMatchersItemMap)
		}
		if err = d.Set("rule_matchers", ruleMatchers); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rule_matchers: %s", err))
		}
	}
	ruleSubgroups := []map[string]interface{}{}
	for _, ruleSubgroupsItem := range ruleGroup.RuleSubgroups {
		ruleSubgroupsItemMap, err := ResourceIbmLogsRuleGroupRulesV1RuleSubgroupToMap(&ruleSubgroupsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		ruleSubgroups = append(ruleSubgroups, ruleSubgroupsItemMap)
	}
	if err = d.Set("rule_subgroups", ruleSubgroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rule_subgroups: %s", err))
	}
	if !core.IsNil(ruleGroup.Order) {
		if err = d.Set("order", flex.IntValue(ruleGroup.Order)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting order: %s", err))
		}
	}

	return nil
}

func resourceIbmLogsRuleGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_rule_group", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, ruleGroupId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRuleGroupOptions := &logsv0.UpdateRuleGroupOptions{}

	updateRuleGroupOptions.SetGroupID(core.UUIDPtr(strfmt.UUID(ruleGroupId)))

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("rule_subgroups") ||
		d.HasChange("description") ||
		d.HasChange("enabled") ||
		d.HasChange("rule_matchers") ||
		d.HasChange("order") {
		updateRuleGroupOptions.SetName(d.Get("name").(string))
		var ruleSubgroups []logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroup
		for _, v := range d.Get("rule_subgroups").([]interface{}) {
			value := v.(map[string]interface{})
			ruleSubgroupsItem, err := ResourceIbmLogsRuleGroupMapToRulesV1CreateRuleGroupRequestCreateRuleSubgroup(value)
			if err != nil {
				return diag.FromErr(err)
			}
			ruleSubgroups = append(ruleSubgroups, *ruleSubgroupsItem)
		}
		updateRuleGroupOptions.SetRuleSubgroups(ruleSubgroups)
		if _, ok := d.GetOk("description"); ok {
			updateRuleGroupOptions.SetDescription(d.Get("description").(string))
		}
		if _, ok := d.GetOkExists("enabled"); ok {
			updateRuleGroupOptions.SetEnabled(d.Get("enabled").(bool))
		}
		if _, ok := d.GetOk("rule_matchers"); ok {
			var ruleMatchers []logsv0.RulesV1RuleMatcherIntf
			for _, v := range d.Get("rule_matchers").([]interface{}) {
				value := v.(map[string]interface{})
				ruleMatchersItem, err := ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcher(value)
				if err != nil {
					return diag.FromErr(err)
				}
				ruleMatchers = append(ruleMatchers, ruleMatchersItem)
			}
			updateRuleGroupOptions.SetRuleMatchers(ruleMatchers)
		}
		if _, ok := d.GetOk("order"); ok {
			updateRuleGroupOptions.SetOrder(int64(d.Get("order").(int)))
		}
		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdateRuleGroupWithContext(context, updateRuleGroupOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateRuleGroupWithContext failed: %s", err.Error()), "ibm_logs_rule_group", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsRuleGroupRead(context, d, meta)
}

func resourceIbmLogsRuleGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_rule_group", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, ruleGroupId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	deleteRuleGroupOptions := &logsv0.DeleteRuleGroupOptions{}

	deleteRuleGroupOptions.SetGroupID(core.UUIDPtr(strfmt.UUID(ruleGroupId)))

	_, err = logsClient.DeleteRuleGroupWithContext(context, deleteRuleGroupOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteRuleGroupWithContext failed: %s", err.Error()), "ibm_logs_rule_group", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1CreateRuleGroupRequestCreateRuleSubgroup(modelMap map[string]interface{}) (*logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroup, error) {
	model := &logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroup{}
	rules := []logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroupCreateRule{}
	for _, rulesItem := range modelMap["rules"].([]interface{}) {
		rulesItemModel, err := ResourceIbmLogsRuleGroupMapToRulesV1CreateRuleGroupRequestCreateRuleSubgroupCreateRule(rulesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		rules = append(rules, *rulesItemModel)
	}
	model.Rules = rules
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	model.Order = core.Int64Ptr(int64(modelMap["order"].(int)))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1CreateRuleGroupRequestCreateRuleSubgroupCreateRule(modelMap map[string]interface{}) (*logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroupCreateRule, error) {
	model := &logsv0.RulesV1CreateRuleGroupRequestCreateRuleSubgroupCreateRule{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.SourceField = core.StringPtr(modelMap["source_field"].(string))
	ParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1RuleParameters(modelMap["parameters"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Parameters = ParametersModel
	model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	model.Order = core.Int64Ptr(int64(modelMap["order"].(int)))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParameters(modelMap map[string]interface{}) (logsv0.RulesV1RuleParametersIntf, error) {
	model := &logsv0.RulesV1RuleParameters{}
	if modelMap["extract_parameters"] != nil && len(modelMap["extract_parameters"].([]interface{})) > 0 {
		ExtractParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ExtractParameters(modelMap["extract_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExtractParameters = ExtractParametersModel
	}
	if modelMap["json_extract_parameters"] != nil && len(modelMap["json_extract_parameters"].([]interface{})) > 0 {
		JSONExtractParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONExtractParameters(modelMap["json_extract_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONExtractParameters = JSONExtractParametersModel
	}
	if modelMap["replace_parameters"] != nil && len(modelMap["replace_parameters"].([]interface{})) > 0 {
		ReplaceParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ReplaceParameters(modelMap["replace_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ReplaceParameters = ReplaceParametersModel
	}
	if modelMap["parse_parameters"] != nil && len(modelMap["parse_parameters"].([]interface{})) > 0 {
		ParseParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ParseParameters(modelMap["parse_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ParseParameters = ParseParametersModel
	}
	if modelMap["allow_parameters"] != nil && len(modelMap["allow_parameters"].([]interface{})) > 0 {
		AllowParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1AllowParameters(modelMap["allow_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AllowParameters = AllowParametersModel
	}
	if modelMap["block_parameters"] != nil && len(modelMap["block_parameters"].([]interface{})) > 0 {
		BlockParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1BlockParameters(modelMap["block_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.BlockParameters = BlockParametersModel
	}
	if modelMap["extract_timestamp_parameters"] != nil && len(modelMap["extract_timestamp_parameters"].([]interface{})) > 0 {
		ExtractTimestampParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ExtractTimestampParameters(modelMap["extract_timestamp_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExtractTimestampParameters = ExtractTimestampParametersModel
	}
	if modelMap["remove_fields_parameters"] != nil && len(modelMap["remove_fields_parameters"].([]interface{})) > 0 {
		RemoveFieldsParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1RemoveFieldsParameters(modelMap["remove_fields_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RemoveFieldsParameters = RemoveFieldsParametersModel
	}
	if modelMap["json_stringify_parameters"] != nil && len(modelMap["json_stringify_parameters"].([]interface{})) > 0 {
		JSONStringifyParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONStringifyParameters(modelMap["json_stringify_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONStringifyParameters = JSONStringifyParametersModel
	}
	if modelMap["json_parse_parameters"] != nil && len(modelMap["json_parse_parameters"].([]interface{})) > 0 {
		JSONParseParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONParseParameters(modelMap["json_parse_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONParseParameters = JSONParseParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1ExtractParameters(modelMap map[string]interface{}) (*logsv0.RulesV1ExtractParameters, error) {
	model := &logsv0.RulesV1ExtractParameters{}
	model.Rule = core.StringPtr(modelMap["rule"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1JSONExtractParameters(modelMap map[string]interface{}) (*logsv0.RulesV1JSONExtractParameters, error) {
	model := &logsv0.RulesV1JSONExtractParameters{}
	if modelMap["destination_field"] != nil && modelMap["destination_field"].(string) != "" {
		model.DestinationField = core.StringPtr(modelMap["destination_field"].(string))
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1ReplaceParameters(modelMap map[string]interface{}) (*logsv0.RulesV1ReplaceParameters, error) {
	model := &logsv0.RulesV1ReplaceParameters{}
	model.DestinationField = core.StringPtr(modelMap["destination_field"].(string))
	model.ReplaceNewVal = core.StringPtr(modelMap["replace_new_val"].(string))
	model.Rule = core.StringPtr(modelMap["rule"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1ParseParameters(modelMap map[string]interface{}) (*logsv0.RulesV1ParseParameters, error) {
	model := &logsv0.RulesV1ParseParameters{}
	model.DestinationField = core.StringPtr(modelMap["destination_field"].(string))
	model.Rule = core.StringPtr(modelMap["rule"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1AllowParameters(modelMap map[string]interface{}) (*logsv0.RulesV1AllowParameters, error) {
	model := &logsv0.RulesV1AllowParameters{}
	model.KeepBlockedLogs = core.BoolPtr(modelMap["keep_blocked_logs"].(bool))
	model.Rule = core.StringPtr(modelMap["rule"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1BlockParameters(modelMap map[string]interface{}) (*logsv0.RulesV1BlockParameters, error) {
	model := &logsv0.RulesV1BlockParameters{}
	model.KeepBlockedLogs = core.BoolPtr(modelMap["keep_blocked_logs"].(bool))
	model.Rule = core.StringPtr(modelMap["rule"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1ExtractTimestampParameters(modelMap map[string]interface{}) (*logsv0.RulesV1ExtractTimestampParameters, error) {
	model := &logsv0.RulesV1ExtractTimestampParameters{}
	model.Standard = core.StringPtr(modelMap["standard"].(string))
	model.Format = core.StringPtr(modelMap["format"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RemoveFieldsParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RemoveFieldsParameters, error) {
	model := &logsv0.RulesV1RemoveFieldsParameters{}
	fields := []string{}
	for _, fieldsItem := range modelMap["fields"].([]interface{}) {
		fields = append(fields, fieldsItem.(string))
	}
	model.Fields = fields
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1JSONStringifyParameters(modelMap map[string]interface{}) (*logsv0.RulesV1JSONStringifyParameters, error) {
	model := &logsv0.RulesV1JSONStringifyParameters{}
	model.DestinationField = core.StringPtr(modelMap["destination_field"].(string))
	if modelMap["delete_source"] != nil {
		model.DeleteSource = core.BoolPtr(modelMap["delete_source"].(bool))
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1JSONParseParameters(modelMap map[string]interface{}) (*logsv0.RulesV1JSONParseParameters, error) {
	model := &logsv0.RulesV1JSONParseParameters{}
	model.DestinationField = core.StringPtr(modelMap["destination_field"].(string))
	if modelMap["delete_source"] != nil {
		model.DeleteSource = core.BoolPtr(modelMap["delete_source"].(bool))
	}
	model.OverrideDest = core.BoolPtr(modelMap["override_dest"].(bool))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersExtractParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersExtractParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersExtractParameters{}
	if modelMap["extract_parameters"] != nil && len(modelMap["extract_parameters"].([]interface{})) > 0 {
		ExtractParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ExtractParameters(modelMap["extract_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExtractParameters = ExtractParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersJSONExtractParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters{}
	if modelMap["json_extract_parameters"] != nil && len(modelMap["json_extract_parameters"].([]interface{})) > 0 {
		JSONExtractParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONExtractParameters(modelMap["json_extract_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONExtractParameters = JSONExtractParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersReplaceParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersReplaceParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersReplaceParameters{}
	if modelMap["replace_parameters"] != nil && len(modelMap["replace_parameters"].([]interface{})) > 0 {
		ReplaceParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ReplaceParameters(modelMap["replace_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ReplaceParameters = ReplaceParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersParseParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersParseParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersParseParameters{}
	if modelMap["parse_parameters"] != nil && len(modelMap["parse_parameters"].([]interface{})) > 0 {
		ParseParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ParseParameters(modelMap["parse_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ParseParameters = ParseParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersAllowParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersAllowParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersAllowParameters{}
	if modelMap["allow_parameters"] != nil && len(modelMap["allow_parameters"].([]interface{})) > 0 {
		AllowParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1AllowParameters(modelMap["allow_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AllowParameters = AllowParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersBlockParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersBlockParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersBlockParameters{}
	if modelMap["block_parameters"] != nil && len(modelMap["block_parameters"].([]interface{})) > 0 {
		BlockParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1BlockParameters(modelMap["block_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.BlockParameters = BlockParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersExtractTimestampParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters{}
	if modelMap["extract_timestamp_parameters"] != nil && len(modelMap["extract_timestamp_parameters"].([]interface{})) > 0 {
		ExtractTimestampParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ExtractTimestampParameters(modelMap["extract_timestamp_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExtractTimestampParameters = ExtractTimestampParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersRemoveFieldsParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters{}
	if modelMap["remove_fields_parameters"] != nil && len(modelMap["remove_fields_parameters"].([]interface{})) > 0 {
		RemoveFieldsParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1RemoveFieldsParameters(modelMap["remove_fields_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RemoveFieldsParameters = RemoveFieldsParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersJSONStringifyParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters{}
	if modelMap["json_stringify_parameters"] != nil && len(modelMap["json_stringify_parameters"].([]interface{})) > 0 {
		JSONStringifyParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONStringifyParameters(modelMap["json_stringify_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONStringifyParameters = JSONStringifyParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleParametersRuleParametersJSONParseParameters(modelMap map[string]interface{}) (*logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters, error) {
	model := &logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters{}
	if modelMap["json_parse_parameters"] != nil && len(modelMap["json_parse_parameters"].([]interface{})) > 0 {
		JSONParseParametersModel, err := ResourceIbmLogsRuleGroupMapToRulesV1JSONParseParameters(modelMap["json_parse_parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.JSONParseParameters = JSONParseParametersModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcher(modelMap map[string]interface{}) (logsv0.RulesV1RuleMatcherIntf, error) {
	model := &logsv0.RulesV1RuleMatcher{}
	if modelMap["application_name"] != nil && len(modelMap["application_name"].([]interface{})) > 0 {
		ApplicationNameModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ApplicationNameConstraint(modelMap["application_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ApplicationName = ApplicationNameModel
	}
	if modelMap["subsystem_name"] != nil && len(modelMap["subsystem_name"].([]interface{})) > 0 {
		SubsystemNameModel, err := ResourceIbmLogsRuleGroupMapToRulesV1SubsystemNameConstraint(modelMap["subsystem_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SubsystemName = SubsystemNameModel
	}
	if modelMap["severity"] != nil && len(modelMap["severity"].([]interface{})) > 0 {
		SeverityModel, err := ResourceIbmLogsRuleGroupMapToRulesV1SeverityConstraint(modelMap["severity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Severity = SeverityModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1ApplicationNameConstraint(modelMap map[string]interface{}) (*logsv0.RulesV1ApplicationNameConstraint, error) {
	model := &logsv0.RulesV1ApplicationNameConstraint{}
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1SubsystemNameConstraint(modelMap map[string]interface{}) (*logsv0.RulesV1SubsystemNameConstraint, error) {
	model := &logsv0.RulesV1SubsystemNameConstraint{}
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1SeverityConstraint(modelMap map[string]interface{}) (*logsv0.RulesV1SeverityConstraint, error) {
	model := &logsv0.RulesV1SeverityConstraint{}
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcherConstraintApplicationName(modelMap map[string]interface{}) (*logsv0.RulesV1RuleMatcherConstraintApplicationName, error) {
	model := &logsv0.RulesV1RuleMatcherConstraintApplicationName{}
	if modelMap["application_name"] != nil && len(modelMap["application_name"].([]interface{})) > 0 {
		ApplicationNameModel, err := ResourceIbmLogsRuleGroupMapToRulesV1ApplicationNameConstraint(modelMap["application_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ApplicationName = ApplicationNameModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcherConstraintSubsystemName(modelMap map[string]interface{}) (*logsv0.RulesV1RuleMatcherConstraintSubsystemName, error) {
	model := &logsv0.RulesV1RuleMatcherConstraintSubsystemName{}
	if modelMap["subsystem_name"] != nil && len(modelMap["subsystem_name"].([]interface{})) > 0 {
		SubsystemNameModel, err := ResourceIbmLogsRuleGroupMapToRulesV1SubsystemNameConstraint(modelMap["subsystem_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SubsystemName = SubsystemNameModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupMapToRulesV1RuleMatcherConstraintSeverity(modelMap map[string]interface{}) (*logsv0.RulesV1RuleMatcherConstraintSeverity, error) {
	model := &logsv0.RulesV1RuleMatcherConstraintSeverity{}
	if modelMap["severity"] != nil && len(modelMap["severity"].([]interface{})) > 0 {
		SeverityModel, err := ResourceIbmLogsRuleGroupMapToRulesV1SeverityConstraint(modelMap["severity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Severity = SeverityModel
	}
	return model, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleMatcherToMap(model logsv0.RulesV1RuleMatcherIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintApplicationName); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintApplicationNameToMap(model.(*logsv0.RulesV1RuleMatcherConstraintApplicationName))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintSubsystemName); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintSubsystemNameToMap(model.(*logsv0.RulesV1RuleMatcherConstraintSubsystemName))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintSeverity); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintSeverityToMap(model.(*logsv0.RulesV1RuleMatcherConstraintSeverity))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcher); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.RulesV1RuleMatcher)
		if model.ApplicationName != nil {
			applicationNameMap, err := ResourceIbmLogsRuleGroupRulesV1ApplicationNameConstraintToMap(model.ApplicationName)
			if err != nil {
				return modelMap, err
			}
			modelMap["application_name"] = []map[string]interface{}{applicationNameMap}
		}
		if model.SubsystemName != nil {
			subsystemNameMap, err := ResourceIbmLogsRuleGroupRulesV1SubsystemNameConstraintToMap(model.SubsystemName)
			if err != nil {
				return modelMap, err
			}
			modelMap["subsystem_name"] = []map[string]interface{}{subsystemNameMap}
		}
		if model.Severity != nil {
			severityMap, err := ResourceIbmLogsRuleGroupRulesV1SeverityConstraintToMap(model.Severity)
			if err != nil {
				return modelMap, err
			}
			modelMap["severity"] = []map[string]interface{}{severityMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.RulesV1RuleMatcherIntf subtype encountered")
	}
}

func ResourceIbmLogsRuleGroupRulesV1ApplicationNameConstraintToMap(model *logsv0.RulesV1ApplicationNameConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1SubsystemNameConstraintToMap(model *logsv0.RulesV1SubsystemNameConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1SeverityConstraintToMap(model *logsv0.RulesV1SeverityConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintApplicationNameToMap(model *logsv0.RulesV1RuleMatcherConstraintApplicationName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApplicationName != nil {
		applicationNameMap, err := ResourceIbmLogsRuleGroupRulesV1ApplicationNameConstraintToMap(model.ApplicationName)
		if err != nil {
			return modelMap, err
		}
		modelMap["application_name"] = []map[string]interface{}{applicationNameMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintSubsystemNameToMap(model *logsv0.RulesV1RuleMatcherConstraintSubsystemName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SubsystemName != nil {
		subsystemNameMap, err := ResourceIbmLogsRuleGroupRulesV1SubsystemNameConstraintToMap(model.SubsystemName)
		if err != nil {
			return modelMap, err
		}
		modelMap["subsystem_name"] = []map[string]interface{}{subsystemNameMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleMatcherConstraintSeverityToMap(model *logsv0.RulesV1RuleMatcherConstraintSeverity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severity != nil {
		severityMap, err := ResourceIbmLogsRuleGroupRulesV1SeverityConstraintToMap(model.Severity)
		if err != nil {
			return modelMap, err
		}
		modelMap["severity"] = []map[string]interface{}{severityMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleSubgroupToMap(model *logsv0.RulesV1RuleSubgroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	rules := []map[string]interface{}{}
	for _, rulesItem := range model.Rules {
		rulesItemMap, err := ResourceIbmLogsRuleGroupRulesV1RuleToMap(&rulesItem)
		if err != nil {
			return modelMap, err
		}
		rules = append(rules, rulesItemMap)
	}
	modelMap["rules"] = rules
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	modelMap["order"] = flex.IntValue(model.Order)
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleToMap(model *logsv0.RulesV1Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	modelMap["source_field"] = *model.SourceField
	parametersMap, err := ResourceIbmLogsRuleGroupRulesV1RuleParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	modelMap["enabled"] = *model.Enabled
	modelMap["order"] = flex.IntValue(model.Order)
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersToMap(model logsv0.RulesV1RuleParametersIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersExtractParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersExtractParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersExtractParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONExtractParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersReplaceParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersReplaceParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersReplaceParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersParseParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersParseParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersParseParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersAllowParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersAllowParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersAllowParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersBlockParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersBlockParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersBlockParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersExtractTimestampParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersRemoveFieldsParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONStringifyParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters); ok {
		return ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONParseParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParameters); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.RulesV1RuleParameters)
		if model.ExtractParameters != nil {
			extractParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ExtractParametersToMap(model.ExtractParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["extract_parameters"] = []map[string]interface{}{extractParametersMap}
		}
		if model.JSONExtractParameters != nil {
			jsonExtractParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONExtractParametersToMap(model.JSONExtractParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["json_extract_parameters"] = []map[string]interface{}{jsonExtractParametersMap}
		}
		if model.ReplaceParameters != nil {
			replaceParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ReplaceParametersToMap(model.ReplaceParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["replace_parameters"] = []map[string]interface{}{replaceParametersMap}
		}
		if model.ParseParameters != nil {
			parseParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ParseParametersToMap(model.ParseParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["parse_parameters"] = []map[string]interface{}{parseParametersMap}
		}
		if model.AllowParameters != nil {
			allowParametersMap, err := ResourceIbmLogsRuleGroupRulesV1AllowParametersToMap(model.AllowParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["allow_parameters"] = []map[string]interface{}{allowParametersMap}
		}
		if model.BlockParameters != nil {
			blockParametersMap, err := ResourceIbmLogsRuleGroupRulesV1BlockParametersToMap(model.BlockParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["block_parameters"] = []map[string]interface{}{blockParametersMap}
		}
		if model.ExtractTimestampParameters != nil {
			extractTimestampParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ExtractTimestampParametersToMap(model.ExtractTimestampParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["extract_timestamp_parameters"] = []map[string]interface{}{extractTimestampParametersMap}
		}
		if model.RemoveFieldsParameters != nil {
			removeFieldsParametersMap, err := ResourceIbmLogsRuleGroupRulesV1RemoveFieldsParametersToMap(model.RemoveFieldsParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["remove_fields_parameters"] = []map[string]interface{}{removeFieldsParametersMap}
		}
		if model.JSONStringifyParameters != nil {
			jsonStringifyParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONStringifyParametersToMap(model.JSONStringifyParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["json_stringify_parameters"] = []map[string]interface{}{jsonStringifyParametersMap}
		}
		if model.JSONParseParameters != nil {
			jsonParseParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONParseParametersToMap(model.JSONParseParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["json_parse_parameters"] = []map[string]interface{}{jsonParseParametersMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.RulesV1RuleParametersIntf subtype encountered")
	}
}

func ResourceIbmLogsRuleGroupRulesV1ExtractParametersToMap(model *logsv0.RulesV1ExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1JSONExtractParametersToMap(model *logsv0.RulesV1JSONExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DestinationField != nil {
		modelMap["destination_field"] = *model.DestinationField
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1ReplaceParametersToMap(model *logsv0.RulesV1ReplaceParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	modelMap["replace_new_val"] = *model.ReplaceNewVal
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1ParseParametersToMap(model *logsv0.RulesV1ParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1AllowParametersToMap(model *logsv0.RulesV1AllowParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["keep_blocked_logs"] = *model.KeepBlockedLogs
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1BlockParametersToMap(model *logsv0.RulesV1BlockParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["keep_blocked_logs"] = *model.KeepBlockedLogs
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1ExtractTimestampParametersToMap(model *logsv0.RulesV1ExtractTimestampParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["standard"] = *model.Standard
	modelMap["format"] = *model.Format
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RemoveFieldsParametersToMap(model *logsv0.RulesV1RemoveFieldsParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["fields"] = model.Fields
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1JSONStringifyParametersToMap(model *logsv0.RulesV1JSONStringifyParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	if model.DeleteSource != nil {
		modelMap["delete_source"] = *model.DeleteSource
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1JSONParseParametersToMap(model *logsv0.RulesV1JSONParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	if model.DeleteSource != nil {
		modelMap["delete_source"] = *model.DeleteSource
	}
	modelMap["override_dest"] = *model.OverrideDest
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersExtractParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExtractParameters != nil {
		extractParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ExtractParametersToMap(model.ExtractParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["extract_parameters"] = []map[string]interface{}{extractParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONExtractParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONExtractParameters != nil {
		jsonExtractParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONExtractParametersToMap(model.JSONExtractParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_extract_parameters"] = []map[string]interface{}{jsonExtractParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersReplaceParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersReplaceParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplaceParameters != nil {
		replaceParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ReplaceParametersToMap(model.ReplaceParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["replace_parameters"] = []map[string]interface{}{replaceParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersParseParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ParseParameters != nil {
		parseParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ParseParametersToMap(model.ParseParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parse_parameters"] = []map[string]interface{}{parseParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersAllowParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersAllowParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AllowParameters != nil {
		allowParametersMap, err := ResourceIbmLogsRuleGroupRulesV1AllowParametersToMap(model.AllowParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["allow_parameters"] = []map[string]interface{}{allowParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersBlockParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersBlockParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BlockParameters != nil {
		blockParametersMap, err := ResourceIbmLogsRuleGroupRulesV1BlockParametersToMap(model.BlockParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["block_parameters"] = []map[string]interface{}{blockParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersExtractTimestampParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExtractTimestampParameters != nil {
		extractTimestampParametersMap, err := ResourceIbmLogsRuleGroupRulesV1ExtractTimestampParametersToMap(model.ExtractTimestampParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["extract_timestamp_parameters"] = []map[string]interface{}{extractTimestampParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersRemoveFieldsParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RemoveFieldsParameters != nil {
		removeFieldsParametersMap, err := ResourceIbmLogsRuleGroupRulesV1RemoveFieldsParametersToMap(model.RemoveFieldsParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["remove_fields_parameters"] = []map[string]interface{}{removeFieldsParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONStringifyParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONStringifyParameters != nil {
		jsonStringifyParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONStringifyParametersToMap(model.JSONStringifyParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_stringify_parameters"] = []map[string]interface{}{jsonStringifyParametersMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsRuleGroupRulesV1RuleParametersRuleParametersJSONParseParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONParseParameters != nil {
		jsonParseParametersMap, err := ResourceIbmLogsRuleGroupRulesV1JSONParseParametersToMap(model.JSONParseParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_parse_parameters"] = []map[string]interface{}{jsonParseParametersMap}
	}
	return modelMap, nil
}
