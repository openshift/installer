// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsRuleGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsRuleGroupsRead,

		Schema: map[string]*schema.Schema{
			"rulegroups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rule groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the rule group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the rule group.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description for the rule group, should express what is the rule group purpose.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether or not the rule is enabled.",
						},
						"rule_matchers": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "// Optional rule matchers which if matched will make the rule go through the rule group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_name": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ApplicationName constraint.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Only logs with this ApplicationName value will match.",
												},
											},
										},
									},
									"subsystem_name": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SubsystemName constraint.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Only logs with this SubsystemName value will match.",
												},
											},
										},
									},
									"severity": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Severity constraint.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
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
							Computed:    true,
							Description: "Rule subgroups. Will try to execute the first rule subgroup, and if not matched will try to match the next one in order.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the rule subgroup.",
									},
									"rules": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Rules to run on the log.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique identifier of the rule.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the rule.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Description of the rule.",
												},
												"source_field": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A field on which value to execute the rule.",
												},
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Parameters for a rule which specifies how it should run.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"extract_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for text extraction rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"rule": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex which will parse the source field and extract the json keys from it while retaining the original log.",
																		},
																	},
																},
															},
															"json_extract_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for json extract rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "In which metadata field to store the extracted value.",
																		},
																	},
																},
															},
															"replace_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for replace rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "In which field to put the modified text.",
																		},
																		"replace_new_val": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The value to replace the matched text with.",
																		},
																		"rule": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex which will match parts in the text to replace.",
																		},
																	},
																},
															},
															"parse_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for parse rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "In which field to put the parsed text.",
																		},
																		"rule": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex which will parse the source field and extract the json keys from it while removing the source field.",
																		},
																	},
																},
															},
															"allow_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for allow rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"keep_blocked_logs": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true matched logs will be blocked, otherwise matched logs will be kept.",
																		},
																		"rule": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex which will match the source field and decide if the rule will apply.",
																		},
																	},
																},
															},
															"block_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for block rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"keep_blocked_logs": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true matched logs will be kept, otherwise matched logs will be blocked.",
																		},
																		"rule": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex which will match the source field and decide if the rule will apply.",
																		},
																	},
																},
															},
															"extract_timestamp_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for extract timestamp rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"standard": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "What time format to use on the extracted time.",
																		},
																		"format": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "What time format the the source field to extract from has.",
																		},
																	},
																},
															},
															"remove_fields_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for remove fields rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"fields": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Json field paths to drop from the log.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"json_stringify_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for json stringify rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Destination field in which to put the json stringified content.",
																		},
																		"delete_source": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether or not to delete the source field after running this rule.",
																		},
																	},
																},
															},
															"json_parse_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Parameters for json parse rule.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Destination field under which to put the json object.",
																		},
																		"delete_source": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Whether or not to delete the source field after running this rule.",
																		},
																		"override_dest": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
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
													Computed:    true,
													Description: "Whether or not to execute the rule.",
												},
												"order": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ordering of the rule subgroup. Lower order will run first. 0 is considered as no value.",
												},
											},
										},
									},
									"enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether or not the rule subgroup is enabled.",
									},
									"order": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ordering of the rule subgroup. Lower order will run first. 0 is considered as no value.",
									},
								},
							},
						},
						"order": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "// The order in which the rule group will be evaluated. The lower the order, the more priority the group will have. Not providing the order will by default create a group with the last order.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsRuleGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_rule_groups", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	listRuleGroupsOptions := &logsv0.ListRuleGroupsOptions{}

	ruleGroupCollection, _, err := logsClient.ListRuleGroupsWithContext(context, listRuleGroupsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListRuleGroupsWithContext failed: %s", err.Error()), "(Data) ibm_logs_rule_groups", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsRuleGroupsID(d))

	rulegroups := []map[string]interface{}{}
	if ruleGroupCollection.Rulegroups != nil {
		for _, modelItem := range ruleGroupCollection.Rulegroups {
			modelMap, err := DataSourceIbmLogsRuleGroupsRuleGroupToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_rule_groups", "read")
				return tfErr.GetDiag()
			}
			rulegroups = append(rulegroups, modelMap)
		}
	}
	if err = d.Set("rulegroups", rulegroups); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rulegroups: %s", err), "(Data) ibm_logs_rule_groups", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsRuleGroupsID returns a reasonable ID for the list.
func dataSourceIbmLogsRuleGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsRuleGroupsRuleGroupToMap(model *logsv0.RuleGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.RuleMatchers != nil {
		ruleMatchers := []map[string]interface{}{}
		for _, ruleMatchersItem := range model.RuleMatchers {
			ruleMatchersItemMap, err := DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherToMap(ruleMatchersItem)
			if err != nil {
				return modelMap, err
			}
			ruleMatchers = append(ruleMatchers, ruleMatchersItemMap)
		}
		modelMap["rule_matchers"] = ruleMatchers
	}
	ruleSubgroups := []map[string]interface{}{}
	for _, ruleSubgroupsItem := range model.RuleSubgroups {
		ruleSubgroupsItemMap, err := DataSourceIbmLogsRuleGroupsRulesV1RuleSubgroupToMap(&ruleSubgroupsItem)
		if err != nil {
			return modelMap, err
		}
		ruleSubgroups = append(ruleSubgroups, ruleSubgroupsItemMap)
	}
	modelMap["rule_subgroups"] = ruleSubgroups
	if model.Order != nil {
		modelMap["order"] = flex.IntValue(model.Order)
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherToMap(model logsv0.RulesV1RuleMatcherIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintApplicationName); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintApplicationNameToMap(model.(*logsv0.RulesV1RuleMatcherConstraintApplicationName))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintSubsystemName); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintSubsystemNameToMap(model.(*logsv0.RulesV1RuleMatcherConstraintSubsystemName))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcherConstraintSeverity); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintSeverityToMap(model.(*logsv0.RulesV1RuleMatcherConstraintSeverity))
	} else if _, ok := model.(*logsv0.RulesV1RuleMatcher); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.RulesV1RuleMatcher)
		if model.ApplicationName != nil {
			applicationNameMap, err := DataSourceIbmLogsRuleGroupsRulesV1ApplicationNameConstraintToMap(model.ApplicationName)
			if err != nil {
				return modelMap, err
			}
			modelMap["application_name"] = []map[string]interface{}{applicationNameMap}
		}
		if model.SubsystemName != nil {
			subsystemNameMap, err := DataSourceIbmLogsRuleGroupsRulesV1SubsystemNameConstraintToMap(model.SubsystemName)
			if err != nil {
				return modelMap, err
			}
			modelMap["subsystem_name"] = []map[string]interface{}{subsystemNameMap}
		}
		if model.Severity != nil {
			severityMap, err := DataSourceIbmLogsRuleGroupsRulesV1SeverityConstraintToMap(model.Severity)
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

func DataSourceIbmLogsRuleGroupsRulesV1ApplicationNameConstraintToMap(model *logsv0.RulesV1ApplicationNameConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1SubsystemNameConstraintToMap(model *logsv0.RulesV1SubsystemNameConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1SeverityConstraintToMap(model *logsv0.RulesV1SeverityConstraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintApplicationNameToMap(model *logsv0.RulesV1RuleMatcherConstraintApplicationName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApplicationName != nil {
		applicationNameMap, err := DataSourceIbmLogsRuleGroupsRulesV1ApplicationNameConstraintToMap(model.ApplicationName)
		if err != nil {
			return modelMap, err
		}
		modelMap["application_name"] = []map[string]interface{}{applicationNameMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintSubsystemNameToMap(model *logsv0.RulesV1RuleMatcherConstraintSubsystemName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SubsystemName != nil {
		subsystemNameMap, err := DataSourceIbmLogsRuleGroupsRulesV1SubsystemNameConstraintToMap(model.SubsystemName)
		if err != nil {
			return modelMap, err
		}
		modelMap["subsystem_name"] = []map[string]interface{}{subsystemNameMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleMatcherConstraintSeverityToMap(model *logsv0.RulesV1RuleMatcherConstraintSeverity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severity != nil {
		severityMap, err := DataSourceIbmLogsRuleGroupsRulesV1SeverityConstraintToMap(model.Severity)
		if err != nil {
			return modelMap, err
		}
		modelMap["severity"] = []map[string]interface{}{severityMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleSubgroupToMap(model *logsv0.RulesV1RuleSubgroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	rules := []map[string]interface{}{}
	for _, rulesItem := range model.Rules {
		rulesItemMap, err := DataSourceIbmLogsRuleGroupsRulesV1RuleToMap(&rulesItem)
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

func DataSourceIbmLogsRuleGroupsRulesV1RuleToMap(model *logsv0.RulesV1Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	modelMap["source_field"] = *model.SourceField
	parametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1RuleParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	modelMap["enabled"] = *model.Enabled
	modelMap["order"] = flex.IntValue(model.Order)
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersToMap(model logsv0.RulesV1RuleParametersIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersExtractParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersExtractParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersExtractParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONExtractParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersReplaceParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersReplaceParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersReplaceParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersParseParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersParseParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersParseParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersAllowParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersAllowParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersAllowParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersBlockParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersBlockParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersBlockParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersExtractTimestampParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersRemoveFieldsParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONStringifyParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters); ok {
		return DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONParseParametersToMap(model.(*logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters))
	} else if _, ok := model.(*logsv0.RulesV1RuleParameters); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.RulesV1RuleParameters)
		if model.ExtractParameters != nil {
			extractParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ExtractParametersToMap(model.ExtractParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["extract_parameters"] = []map[string]interface{}{extractParametersMap}
		}
		if model.JSONExtractParameters != nil {
			jsonExtractParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONExtractParametersToMap(model.JSONExtractParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["json_extract_parameters"] = []map[string]interface{}{jsonExtractParametersMap}
		}
		if model.ReplaceParameters != nil {
			replaceParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ReplaceParametersToMap(model.ReplaceParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["replace_parameters"] = []map[string]interface{}{replaceParametersMap}
		}
		if model.ParseParameters != nil {
			parseParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ParseParametersToMap(model.ParseParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["parse_parameters"] = []map[string]interface{}{parseParametersMap}
		}
		if model.AllowParameters != nil {
			allowParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1AllowParametersToMap(model.AllowParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["allow_parameters"] = []map[string]interface{}{allowParametersMap}
		}
		if model.BlockParameters != nil {
			blockParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1BlockParametersToMap(model.BlockParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["block_parameters"] = []map[string]interface{}{blockParametersMap}
		}
		if model.ExtractTimestampParameters != nil {
			extractTimestampParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ExtractTimestampParametersToMap(model.ExtractTimestampParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["extract_timestamp_parameters"] = []map[string]interface{}{extractTimestampParametersMap}
		}
		if model.RemoveFieldsParameters != nil {
			removeFieldsParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1RemoveFieldsParametersToMap(model.RemoveFieldsParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["remove_fields_parameters"] = []map[string]interface{}{removeFieldsParametersMap}
		}
		if model.JSONStringifyParameters != nil {
			jsonStringifyParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONStringifyParametersToMap(model.JSONStringifyParameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["json_stringify_parameters"] = []map[string]interface{}{jsonStringifyParametersMap}
		}
		if model.JSONParseParameters != nil {
			jsonParseParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONParseParametersToMap(model.JSONParseParameters)
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

func DataSourceIbmLogsRuleGroupsRulesV1ExtractParametersToMap(model *logsv0.RulesV1ExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1JSONExtractParametersToMap(model *logsv0.RulesV1JSONExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DestinationField != nil {
		modelMap["destination_field"] = *model.DestinationField
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1ReplaceParametersToMap(model *logsv0.RulesV1ReplaceParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	modelMap["replace_new_val"] = *model.ReplaceNewVal
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1ParseParametersToMap(model *logsv0.RulesV1ParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1AllowParametersToMap(model *logsv0.RulesV1AllowParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["keep_blocked_logs"] = *model.KeepBlockedLogs
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1BlockParametersToMap(model *logsv0.RulesV1BlockParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["keep_blocked_logs"] = *model.KeepBlockedLogs
	modelMap["rule"] = *model.Rule
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1ExtractTimestampParametersToMap(model *logsv0.RulesV1ExtractTimestampParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["standard"] = *model.Standard
	modelMap["format"] = *model.Format
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RemoveFieldsParametersToMap(model *logsv0.RulesV1RemoveFieldsParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["fields"] = model.Fields
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1JSONStringifyParametersToMap(model *logsv0.RulesV1JSONStringifyParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	if model.DeleteSource != nil {
		modelMap["delete_source"] = *model.DeleteSource
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1JSONParseParametersToMap(model *logsv0.RulesV1JSONParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["destination_field"] = *model.DestinationField
	if model.DeleteSource != nil {
		modelMap["delete_source"] = *model.DeleteSource
	}
	modelMap["override_dest"] = *model.OverrideDest
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersExtractParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExtractParameters != nil {
		extractParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ExtractParametersToMap(model.ExtractParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["extract_parameters"] = []map[string]interface{}{extractParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONExtractParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONExtractParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONExtractParameters != nil {
		jsonExtractParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONExtractParametersToMap(model.JSONExtractParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_extract_parameters"] = []map[string]interface{}{jsonExtractParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersReplaceParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersReplaceParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplaceParameters != nil {
		replaceParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ReplaceParametersToMap(model.ReplaceParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["replace_parameters"] = []map[string]interface{}{replaceParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersParseParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ParseParameters != nil {
		parseParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ParseParametersToMap(model.ParseParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parse_parameters"] = []map[string]interface{}{parseParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersAllowParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersAllowParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AllowParameters != nil {
		allowParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1AllowParametersToMap(model.AllowParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["allow_parameters"] = []map[string]interface{}{allowParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersBlockParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersBlockParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BlockParameters != nil {
		blockParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1BlockParametersToMap(model.BlockParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["block_parameters"] = []map[string]interface{}{blockParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersExtractTimestampParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersExtractTimestampParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExtractTimestampParameters != nil {
		extractTimestampParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1ExtractTimestampParametersToMap(model.ExtractTimestampParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["extract_timestamp_parameters"] = []map[string]interface{}{extractTimestampParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersRemoveFieldsParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersRemoveFieldsParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RemoveFieldsParameters != nil {
		removeFieldsParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1RemoveFieldsParametersToMap(model.RemoveFieldsParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["remove_fields_parameters"] = []map[string]interface{}{removeFieldsParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONStringifyParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONStringifyParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONStringifyParameters != nil {
		jsonStringifyParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONStringifyParametersToMap(model.JSONStringifyParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_stringify_parameters"] = []map[string]interface{}{jsonStringifyParametersMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsRuleGroupsRulesV1RuleParametersRuleParametersJSONParseParametersToMap(model *logsv0.RulesV1RuleParametersRuleParametersJSONParseParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.JSONParseParameters != nil {
		jsonParseParametersMap, err := DataSourceIbmLogsRuleGroupsRulesV1JSONParseParametersToMap(model.JSONParseParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["json_parse_parameters"] = []map[string]interface{}{jsonParseParametersMap}
	}
	return modelMap, nil
}
