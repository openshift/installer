// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"encoding/json"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/rulesetsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	CISRulesetsListOutput                              = "rulesets_list"
	CISRulesetsObjectOutput                            = "rulesets"
	CISRulesetsDescription                             = "description"
	CISRulesetsKind                                    = "kind"
	CISRulesetsName                                    = "name"
	CISRulesetsPhase                                   = "phase"
	CISRulesetsLastUpdatedAt                           = "last_updated"
	CISRulesetsVersion                                 = "version"
	CISRulesetsRule                                    = "rule"
	CISRulesetsRules                                   = "rules"
	CISRulesetsRuleId                                  = "id"
	CISRulesetsRuleVersion                             = "version"
	CISRulesetsRuleAction                              = "action"
	CISRulesetsRuleActionParameters                    = "action_parameters"
	CISRulesetsRuleActionParametersResponse            = "response"
	CISRulesetsRuleActionParametersResponseContent     = "content"
	CISRulesetsRuleActionParametersResponseContentType = "content_type"
	CISRulesetsRuleActionParametersResponseStatusCode  = "status_code"
	CISRulesetsRuleExpression                          = "expression"
	CISRulesetsRuleRef                                 = "ref"
	CISRulesetsRuleLogging                             = "logging"
	CISRulesetsRuleLoggingEnabled                      = "enabled"
	CISRulesetsRuleLastUpdatedAt                       = "last_updated_at"
	CISRulesetsId                                      = "ruleset_id"
	CISRuleset                                         = "ruleset"
	CISRulesetList                                     = "rulesets"
	CISRulesetOverrides                                = "overrides"
	CISRulesetOverridesAction                          = "action"
	CISRulesetOverridesEnabled                         = "enabled"
	CISRulesetOverridesSensitivityLevel                = "sensitivity_level"
	CISRulesetOverridesCategories                      = "categories"
	CISRulesetOverridesCategoriesCategory              = "category"
	CISRulesetOverridesRules                           = "override_rules"
	CISRulesetsRuleActionCategories                    = "categories"
	CISRulesetsRuleActionEnabled                       = "enabled"
	CISRulesetsRuleActionDescription                   = "description"
	CISRulesetsRulePosition                            = "position"
	CISRulesetsRulePositionAfter                       = "after"
	CISRulesetsRulePositionBefore                      = "before"
	CISRulesetsRulePositionIndex                       = "index"
	CISRulesetRuleId                                   = "rule_id"
)

var CISResponseObject = &schema.Resource{
	Schema: map[string]*schema.Schema{
		CISRulesetsDescription: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of the Rulesets",
		},
		CISRulesetsId: {
			Type:        schema.TypeString,
			Description: "Associated Ruleset ID",
			Computed:    true,
		},
		CISRulesetsKind: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kind of the Rulesets",
		},
		CISRulesetsLastUpdatedAt: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Rulesets Last Updated At",
		},
		CISRulesetsName: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Rulesets",
		},
		CISRulesetsPhase: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Phase of the Rulesets",
		},
		CISRulesetsVersion: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Version of the Rulesets",
		},
		CISRulesetsRules: {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Rules of the Rulesets",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					CISRulesetsRuleId: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Id of the Rulesets Rule",
					},
					CISRulesetsRuleVersion: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Version of the Rulesets Rule",
					},
					CISRulesetsRuleAction: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Action of the Rulesets Rule",
					},
					CISRulesetsRuleActionParameters: {
						Type:        schema.TypeSet,
						Computed:    true,
						Description: "Action parameters of the Rulesets Rule",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								CISRulesetsRuleId: {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Id of the Rulesets Rule",
								},
								CISRulesetOverrides: {
									Type:        schema.TypeSet,
									Computed:    true,
									Description: "Override options",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											CISRulesetOverridesAction: {
												Type:        schema.TypeString,
												Computed:    true,
												Description: "Action to perform",
											},
											CISRulesetOverridesEnabled: {
												Type:        schema.TypeBool,
												Computed:    true,
												Description: "Enable Disable Rule",
											},
											// CISRulesetOverridesSensitivityLevel: {
											// 	Type:        schema.TypeString,
											// 	Computed:    true,
											// 	Description: "Sensitivity Level",
											// },
											CISRulesetOverridesRules: {
												Type:        schema.TypeList,
												Computed:    true,
												Description: "Rules",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														CISRulesetRuleId: {
															Type:        schema.TypeString,
															Computed:    true,
															Description: "Id of the Ruleset",
														},
														CISRulesetOverridesEnabled: {
															Type:        schema.TypeBool,
															Computed:    true,
															Description: "Enable Disable Rule",
														},
														CISRulesetOverridesAction: {
															Type:        schema.TypeString,
															Computed:    true,
															Description: "Action to perform",
														},
														// CISRulesetOverridesSensitivityLevel: {
														// 	Type:        schema.TypeString,
														// 	Computed:    true,
														// 	Description: "Sensitivity Level",
														// },
													},
												},
											},
											CISRulesetOverridesCategories: {
												Type:        schema.TypeList,
												Computed:    true,
												Description: "Categories",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														CISRulesetOverridesCategoriesCategory: {
															Type:        schema.TypeString,
															Computed:    true,
															Description: "Category",
														},
														CISRulesetOverridesEnabled: {
															Type:        schema.TypeBool,
															Computed:    true,
															Description: "Enable Disable Rule",
														},
														CISRulesetOverridesAction: {
															Type:        schema.TypeString,
															Computed:    true,
															Description: "Action to perform",
														},
													},
												},
											},
										},
									},
								},
								CISRulesetsVersion: {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Version of the Ruleset",
								},
								CISRuleset: {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "Ruleset ID of the ruleset to apply action to.",
								},
								CISRulesetList: {
									Type:        schema.TypeList,
									Computed:    true,
									Description: "List of Ruleset IDs of the ruleset to apply action to.",
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								CISRulesetsRuleActionParametersResponse: {
									Type:        schema.TypeSet,
									Computed:    true,
									Description: "Action parameters response of the Rulesets Rule",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											CISRulesetsRuleActionParametersResponseContent: {
												Type:        schema.TypeString,
												Computed:    true,
												Description: "Action parameters response content of the Rulesets Rule",
											},
											CISRulesetsRuleActionParametersResponseContentType: {
												Type:        schema.TypeString,
												Computed:    true,
												Description: "Action parameters response type of the Rulesets Rule",
											},
											CISRulesetsRuleActionParametersResponseStatusCode: {
												Type:        schema.TypeInt,
												Computed:    true,
												Description: "Action parameters response status code of the Rulesets Rule",
											},
										},
									},
								},
							},
						},
					},
					CISRulesetsRuleActionCategories: {
						Type:        schema.TypeList,
						Computed:    true,
						Description: "Categories of the Rulesets Rule",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					CISRulesetsRuleActionEnabled: {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "Enable/Disable Ruleset Rule",
					},
					CISRulesetsRuleActionDescription: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Description of the Rulesets Rule",
					},
					CISRulesetsRuleExpression: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Experession of the Rulesets Rule",
					},
					CISRulesetsRuleRef: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Reference of the Rulesets Rule",
					},
					CISRulesetsRuleLogging: {
						Type:        schema.TypeMap,
						Computed:    true,
						Description: "Logging of the Rulesets Rule",
						Elem:        &schema.Schema{Type: schema.TypeBool},
					},
					CISRulesetsRuleLastUpdatedAt: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Rulesets Rule Last Updated At",
					},
				},
			},
		},
	},
}

func DataSourceIBMCISRulesets() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISRulesetsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_rulesets",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Optional:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			CISRulesetsId: {
				Type:        schema.TypeString,
				Description: "Associated Ruleset ID",
				Optional:    true,
			},
			CISRulesetsListOutput: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem:        CISResponseObject,
			},
			CISRulesetsObjectOutput: {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Container for response information.",
				Elem:        CISResponseObject,
			},
		},
	}
}

func DataSourceIBMCISRulesetsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISRulesetsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_rulesets",
		Schema:       validateSchema}
	return &iBMCISRulesetsValidator
}

func dataIBMCISRulesetsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	zoneId := d.Get(cisDomainID).(string)
	rulesetId := d.Get(CISRulesetsId).(string)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)

		if rulesetId != "" {
			opt := sess.NewGetZoneRulesetOptions(rulesetId)
			result, resp, err := sess.GetZoneRuleset(opt)
			if err != nil {
				log.Printf("[WARN] Get Instance ruleset failed: %v\n", resp)
				return err
			}
			rulesetObj := flattenCISRulesets(*result.Result)

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsObjectOutput, rulesetObj)
			d.Set(cisID, crn)

		} else {
			opt := sess.NewGetZoneRulesetsOptions()
			result, resp, err := sess.GetZoneRulesets(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetList := make([]map[string]interface{}, 0)
			for _, rulesetObj := range result.Result {
				rulesetOutput := map[string]interface{}{}
				rulesetOutput[CISRulesetsDescription] = *rulesetObj.Description
				rulesetOutput[CISRulesetsKind] = *rulesetObj.Kind
				rulesetOutput[CISRulesetsName] = *rulesetObj.Name
				rulesetOutput[CISRulesetsPhase] = *rulesetObj.Phase
				rulesetOutput[CISRulesetsLastUpdatedAt] = *rulesetObj.LastUpdated
				rulesetOutput[CISRulesetsVersion] = *rulesetObj.Version
				rulesetOutput[CISRulesetsId] = *rulesetObj.ID

				if rulesetOutput[CISRulesetsPhase] == "http_request_firewall_managed" {
					rulesetList = append(rulesetList, rulesetOutput)
				}
			}

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsListOutput, rulesetList)
			d.Set(cisID, crn)
		}

	} else {

		if rulesetId != "" {
			opt := sess.NewGetInstanceRulesetOptions(rulesetId)
			result, resp, err := sess.GetInstanceRuleset(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetObj := flattenCISRulesets(*result.Result)

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsListOutput, rulesetObj)
			d.Set(cisID, crn)

		} else {
			opt := sess.NewGetInstanceRulesetsOptions()
			result, resp, err := sess.GetInstanceRulesets(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetList := make([]map[string]interface{}, 0)
			for _, rulesetObj := range result.Result {
				rulesetOutput := map[string]interface{}{}
				rulesetOutput[CISRulesetsDescription] = *rulesetObj.Description
				rulesetOutput[CISRulesetsKind] = *rulesetObj.Kind
				rulesetOutput[CISRulesetsName] = *rulesetObj.Name
				rulesetOutput[CISRulesetsPhase] = *rulesetObj.Phase
				rulesetOutput[CISRulesetsLastUpdatedAt] = *rulesetObj.LastUpdated
				rulesetOutput[CISRulesetsVersion] = *rulesetObj.Version
				rulesetOutput[CISRulesetsId] = *rulesetObj.ID

				if rulesetOutput[CISRulesetsPhase] == "http_request_firewall_managed" {
					rulesetList = append(rulesetList, rulesetOutput)
				}
			}

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsListOutput, rulesetList)
			d.Set(cisID, crn)
		}
	}

	return nil
}

func flattenCISRulesets(rulesetObj rulesetsv1.RulesetDetails) interface{} {

	finalrulesetObj := make([]interface{}, 0)

	rulesetOutput := map[string]interface{}{}
	rulesetOutput[CISRulesetsDescription] = *rulesetObj.Description
	rulesetOutput[CISRulesetsKind] = *rulesetObj.Kind
	rulesetOutput[CISRulesetsName] = *rulesetObj.Name
	rulesetOutput[CISRulesetsPhase] = *rulesetObj.Phase
	rulesetOutput[CISRulesetsLastUpdatedAt] = *rulesetObj.LastUpdated
	rulesetOutput[CISRulesetsVersion] = *rulesetObj.Version
	rulesetOutput[CISRulesetsId] = *&rulesetObj.ID

	ruleDetailsList := make([]map[string]interface{}, 0)
	for _, ruleDetailsObj := range rulesetObj.Rules {
		ruleDetails := map[string]interface{}{}
		ruleDetails[CISRulesetsRuleId] = ruleDetailsObj.ID
		ruleDetails[CISRulesetsRuleVersion] = ruleDetailsObj.Version
		ruleDetails[CISRulesetsRuleAction] = ruleDetailsObj.Action
		ruleDetails[CISRulesetsRuleExpression] = ruleDetailsObj.Expression
		ruleDetails[CISRulesetsRuleRef] = ruleDetailsObj.Ref
		ruleDetails[CISRulesetsRuleLastUpdatedAt] = ruleDetailsObj.LastUpdated
		ruleDetails[CISRulesetsRuleActionCategories] = ruleDetailsObj.Categories
		ruleDetails[CISRulesetsRuleActionEnabled] = ruleDetailsObj.Enabled
		ruleDetails[CISRulesetsRuleActionDescription] = ruleDetailsObj.Description

		// Not Applicable for now
		ruleDetails[CISRulesetsRuleLogging] = ruleDetailsObj.Logging

		flattenedActionParameter := flattenCISRulesetsRuleActionParameters(ruleDetailsObj.ActionParameters)

		if len(flattenedActionParameter) != 0 {
			ruleDetails[CISRulesetsRuleActionParameters] = []map[string]interface{}{flattenedActionParameter}
		}

		ruleDetailsList = append(ruleDetailsList, ruleDetails)
	}

	rulesetOutput[CISRulesetsRules] = ruleDetailsList

	finalrulesetObj = append(finalrulesetObj, rulesetOutput)

	return finalrulesetObj
}

func flattenCISRulesetsRuleActionParameters(rulesetsRuleActionParameterObj *rulesetsv1.ActionParameters) map[string]interface{} {
	actionParametersOutput := map[string]interface{}{}
	resultOutput := map[string]interface{}{}

	res, _ := json.Marshal(rulesetsRuleActionParameterObj)
	json.Unmarshal(res, &actionParametersOutput)

	if val, ok := actionParametersOutput["id"]; ok {
		resultOutput[CISRulesetsRuleId] = val.(string)
	}
	if val, ok := actionParametersOutput["ruleset"]; ok {
		resultOutput[CISRuleset] = val.(string)
	}
	if val, ok := actionParametersOutput["version"]; ok {
		resultOutput[CISRulesetsVersion] = val.(string)
	}
	if _, ok := actionParametersOutput["rulesets"]; ok {
		resultOutput[CISRulesetList] = rulesetsRuleActionParameterObj.Rulesets
	}
	if val, ok := actionParametersOutput["response"]; ok {
		response := map[string]interface{}{}

		res, _ := json.Marshal(val)
		json.Unmarshal(res, &response)

		resultOutput[CISRulesetsRuleActionParametersResponse] = response
	}

	if _, ok := actionParametersOutput["overrides"]; ok {
		flattenCISRulesetsRuleActionParameterOverrides := flattenCISRulesetsRuleActionParameterOverrides(rulesetsRuleActionParameterObj.Overrides)
		resultOutput[CISRulesetOverrides] = []map[string]interface{}{flattenCISRulesetsRuleActionParameterOverrides}
	}

	return resultOutput
}

func flattenCISRulesetsRuleActionParameterOverrides(rulesetsRuleActionParameterOverridesObj *rulesetsv1.Overrides) map[string]interface{} {
	actionParameterOverridesOutput := map[string]interface{}{}
	resultOutput := map[string]interface{}{}

	res, _ := json.Marshal(rulesetsRuleActionParameterOverridesObj)
	json.Unmarshal(res, &actionParameterOverridesOutput)

	if val, ok := actionParameterOverridesOutput["action"]; ok {
		resultOutput[CISRulesetOverridesAction] = val.(string)
	}
	if val, ok := actionParameterOverridesOutput["enabled"]; ok {
		resultOutput[CISRulesetOverridesEnabled] = val.(bool)
	}

	if _, ok := actionParameterOverridesOutput["categories"]; ok {
		categoriesList := make([]map[string]interface{}, 0)
		for _, obj := range rulesetsRuleActionParameterOverridesObj.Categories {
			categoriesObj := map[string]interface{}{}
			categoriesObj[CISRulesetOverridesCategoriesCategory] = obj.Category
			categoriesObj[CISRulesetOverridesEnabled] = obj.Enabled
			categoriesObj[CISRulesetOverridesAction] = obj.Action

			categoriesList = append(categoriesList, categoriesObj)
		}

		resultOutput[CISRulesetOverridesCategories] = categoriesList
	}
	if _, ok := actionParameterOverridesOutput["rules"]; ok {

		overrideRulesList := make([]map[string]interface{}, 0)
		for _, obj := range rulesetsRuleActionParameterOverridesObj.Rules {
			overrideRulesObj := map[string]interface{}{}
			overrideRulesObj[CISRulesetRuleId] = obj.ID
			overrideRulesObj[CISRulesetOverridesEnabled] = obj.Enabled
			overrideRulesObj[CISRulesetOverridesAction] = obj.Action

			overrideRulesList = append(overrideRulesList, overrideRulesObj)
		}

		resultOutput[CISRulesetOverridesRules] = overrideRulesList
	}

	return resultOutput
}

func dataSourceCISRulesetsCheckID(d *schema.ResourceData) string {
	return d.Get(CISRulesetsId).(string) + ":" + d.Get(cisDomainID).(string) + ":" + d.Get(cisID).(string)
}
