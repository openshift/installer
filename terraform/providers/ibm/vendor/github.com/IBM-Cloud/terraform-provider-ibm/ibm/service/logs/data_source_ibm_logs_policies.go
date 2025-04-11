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

func DataSourceIbmLogsPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsPoliciesRead,

		Schema: map[string]*schema.Schema{
			"enabled_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optionally filter only enabled policies.",
			},
			"source_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Source type to filter policies by.",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Company policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID.",
						},
						"company_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Company ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of policy.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of policy.",
						},
						"priority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data pipeline sources that match the policy rules will go through.",
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
						"application_rule": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule for matching with application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier of the rule.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value of the rule.",
									},
								},
							},
						},
						"subsystem_rule": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule for matching with application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier of the rule.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value of the rule.",
									},
								},
							},
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
						"archive_retention": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Archive retention definition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "References archive retention definition.",
									},
								},
							},
						},
						"log_rules": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Log rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severities": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Source severities to match with.",
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
		},
	}
}

func dataSourceIbmLogsPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getCompanyPoliciesOptions := &logsv0.GetCompanyPoliciesOptions{}

	if _, ok := d.GetOk("enabled_only"); ok {
		getCompanyPoliciesOptions.SetEnabledOnly(d.Get("enabled_only").(bool))
	}
	if _, ok := d.GetOk("source_type"); ok {
		getCompanyPoliciesOptions.SetSourceType(d.Get("source_type").(string))
	}

	policyCollection, _, err := logsClient.GetCompanyPoliciesWithContext(context, getCompanyPoliciesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCompanyPoliciesWithContext failed: %s", err.Error()), "(Data) ibm_logs_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsPoliciesID(d))

	policies := []map[string]interface{}{}
	if policyCollection.Policies != nil {
		for _, modelItem := range policyCollection.Policies {
			modelMap, err := DataSourceIbmLogsPoliciesPolicyToMap(modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_policies", "read")
				return tfErr.GetDiag()
			}
			policies = append(policies, modelMap)
		}
	}
	if err = d.Set("policies", policies); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting policies: %s", err), "(Data) ibm_logs_policies", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsPoliciesID returns a reasonable ID for the list.
func dataSourceIbmLogsPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsPoliciesPolicyToMap(model logsv0.PolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.PolicyQuotaV1PolicySourceTypeRulesLogRules); ok {
		return DataSourceIbmLogsPoliciesPolicyQuotaV1PolicySourceTypeRulesLogRulesToMap(model.(*logsv0.PolicyQuotaV1PolicySourceTypeRulesLogRules))
	} else if _, ok := model.(*logsv0.Policy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.Policy)
		modelMap["id"] = model.ID.String()
		modelMap["company_id"] = flex.IntValue(model.CompanyID)
		modelMap["name"] = *model.Name
		modelMap["description"] = *model.Description
		if model.Priority != nil {
			modelMap["priority"] = *model.Priority
		}
		if model.Deleted != nil {
			modelMap["deleted"] = *model.Deleted
		}
		if model.Enabled != nil {
			modelMap["enabled"] = *model.Enabled
		}
		modelMap["order"] = flex.IntValue(model.Order)
		if model.ApplicationRule != nil {
			applicationRuleMap, err := DataSourceIbmLogsPoliciesQuotaV1RuleToMap(model.ApplicationRule)
			if err != nil {
				return modelMap, err
			}
			modelMap["application_rule"] = []map[string]interface{}{applicationRuleMap}
		}
		if model.SubsystemRule != nil {
			subsystemRuleMap, err := DataSourceIbmLogsPoliciesQuotaV1RuleToMap(model.SubsystemRule)
			if err != nil {
				return modelMap, err
			}
			modelMap["subsystem_rule"] = []map[string]interface{}{subsystemRuleMap}
		}
		modelMap["created_at"] = *model.CreatedAt
		modelMap["updated_at"] = *model.UpdatedAt
		if model.ArchiveRetention != nil {
			archiveRetentionMap, err := DataSourceIbmLogsPoliciesQuotaV1ArchiveRetentionToMap(model.ArchiveRetention)
			if err != nil {
				return modelMap, err
			}
			modelMap["archive_retention"] = []map[string]interface{}{archiveRetentionMap}
		}
		if model.LogRules != nil {
			logRulesMap, err := DataSourceIbmLogsPoliciesQuotaV1LogRulesToMap(model.LogRules)
			if err != nil {
				return modelMap, err
			}
			modelMap["log_rules"] = []map[string]interface{}{logRulesMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.PolicyIntf subtype encountered")
	}
}

func DataSourceIbmLogsPoliciesQuotaV1RuleToMap(model *logsv0.QuotaV1Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["rule_type_id"] = *model.RuleTypeID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIbmLogsPoliciesQuotaV1ArchiveRetentionToMap(model *logsv0.QuotaV1ArchiveRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	return modelMap, nil
}

func DataSourceIbmLogsPoliciesQuotaV1LogRulesToMap(model *logsv0.QuotaV1LogRules) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	return modelMap, nil
}

func DataSourceIbmLogsPoliciesPolicyQuotaV1PolicySourceTypeRulesLogRulesToMap(model *logsv0.PolicyQuotaV1PolicySourceTypeRulesLogRules) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["company_id"] = flex.IntValue(model.CompanyID)
	modelMap["name"] = *model.Name
	modelMap["description"] = *model.Description
	if model.Priority != nil {
		modelMap["priority"] = *model.Priority
	}
	if model.Deleted != nil {
		modelMap["deleted"] = *model.Deleted
	}
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	modelMap["order"] = flex.IntValue(model.Order)
	if model.ApplicationRule != nil {
		applicationRuleMap, err := DataSourceIbmLogsPoliciesQuotaV1RuleToMap(model.ApplicationRule)
		if err != nil {
			return modelMap, err
		}
		modelMap["application_rule"] = []map[string]interface{}{applicationRuleMap}
	}
	if model.SubsystemRule != nil {
		subsystemRuleMap, err := DataSourceIbmLogsPoliciesQuotaV1RuleToMap(model.SubsystemRule)
		if err != nil {
			return modelMap, err
		}
		modelMap["subsystem_rule"] = []map[string]interface{}{subsystemRuleMap}
	}
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.ArchiveRetention != nil {
		archiveRetentionMap, err := DataSourceIbmLogsPoliciesQuotaV1ArchiveRetentionToMap(model.ArchiveRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["archive_retention"] = []map[string]interface{}{archiveRetentionMap}
	}
	if model.LogRules != nil {
		logRulesMap, err := DataSourceIbmLogsPoliciesQuotaV1LogRulesToMap(model.LogRules)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_rules"] = []map[string]interface{}{logRulesMap}
	}
	return modelMap, nil
}
