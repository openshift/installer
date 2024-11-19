// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/rulesetsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var CISRulesetsRulesObject = &schema.Resource{
	Schema: map[string]*schema.Schema{
		CISRulesetsRuleId: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the Rulesets Rule",
			Computed:    true,
		},
		CISRulesetsRuleVersion: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Version of the Rulesets Rule",
		},
		CISRulesetsRuleAction: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Action of the Rulesets Rule",
		},
		CISRulesetsRuleActionParameters: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Action parameters of the Rulesets Rule",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					CISRulesetsRuleId: {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Id of the Rulesets Rule",
						Computed:    true,
					},
					CISRulesetOverrides: {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "Override options",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								CISRulesetOverridesAction: {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Action to perform",
								},
								CISRulesetOverridesEnabled: {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Enable Disable Rule",
								},
								// CISRulesetOverridesSensitivityLevel: {
								// 	Type:        schema.TypeString,
								// 	Optional:    true,
								// 	Description: "Sensitivity Level",
								// },
								CISRulesetOverridesRules: {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Rules",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											CISRulesetRuleId: {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Id of the Ruleset",
											},
											CISRulesetOverridesEnabled: {
												Type:        schema.TypeBool,
												Optional:    true,
												Description: "Enable Disable Rule",
											},
											CISRulesetOverridesAction: {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Action to perform",
											},
											CISRulesetOverridesSensitivityLevel: {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Sensitivity Level",
											},
										},
									},
								},
								CISRulesetOverridesCategories: {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Categories",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											CISRulesetOverridesCategoriesCategory: {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Category",
											},
											CISRulesetOverridesEnabled: {
												Type:        schema.TypeBool,
												Optional:    true,
												Description: "Enable Disable Rule",
											},
											CISRulesetOverridesAction: {
												Type:        schema.TypeString,
												Optional:    true,
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
						Optional:    true,
						Description: "Version of the Ruleset",
					},
					CISRuleset: {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Ruleset ID of the ruleset to apply action to.",
					},
					CISRulesetList: {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of Ruleset IDs of the ruleset to apply action to.",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					CISRulesetsRuleActionParametersResponse: {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "Action parameters response of the Rulesets Rule",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								CISRulesetsRuleActionParametersResponseContent: {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Action parameters response content of the Rulesets Rule",
								},
								CISRulesetsRuleActionParametersResponseContentType: {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Action parameters response type of the Rulesets Rule",
								},
								CISRulesetsRuleActionParametersResponseStatusCode: {
									Type:        schema.TypeInt,
									Optional:    true,
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
			Optional:    true,
			Description: "Categories of the Rulesets Rule",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		CISRulesetsRuleActionEnabled: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable/Disable Ruleset Rule",
		},
		CISRulesetsRuleActionDescription: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the Rulesets Rule",
		},
		CISRulesetsRuleExpression: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Experession of the Rulesets Rule",
		},
		CISRulesetsRuleRef: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Reference of the Rulesets Rule",
		},
		CISRulesetsRuleLogging: {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Logging of the Rulesets Rule",
			Elem:        &schema.Schema{Type: schema.TypeBool},
		},
		CISRulesetsRulePosition: {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Position of Rulesets Rule",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					CISRulesetsRulePositionAfter: {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Ruleset before Position of Rulesets Rule",
					},
					CISRulesetsRulePositionBefore: {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Ruleset after Position of Rulesets Rule",
					},
					CISRulesetsRulePositionIndex: {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Index of Rulesets Rule",
					},
				},
			},
		},
		CISRulesetsRuleLastUpdatedAt: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Rulesets Rule Last Updated At",
		},
	},
}

func ResourceIBMCISRulesetRule() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISRulesetRuleCreate,
		Read:     ResourceIBMCISRulesetRuleRead,
		Update:   ResourceIBMCISRulesetRuleUpdate,
		Delete:   ResourceIBMCISRulesetRuleDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_ruleset_rule",
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
				Required:    true,
			},
			CISRulesetsRule: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rules of the Rulesets",
				Elem:        CISRulesetsRulesObject,
			},
		},
	}
}
func ResourceIBMCISRulesetRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISRulesetValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_ruleset_rule",
		Schema:       validateSchema}
	return &ibmCISRulesetValidator
}

func ResourceIBMCISRulesetRuleCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}
	crn := d.Get(cisID).(string)
	zoneId := d.Get(cisDomainID).(string)
	rulesetId := d.Get(CISRulesetsId).(string)

	sess.Crn = core.StringPtr(crn)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)
		opt := sess.NewCreateZoneRulesetRuleOptions(rulesetId)

		rulesObject := d.Get(CISRulesetsRule).([]interface{})[0].(map[string]interface{})

		opt.SetRulesetID(rulesetId)
		opt.SetExpression(rulesObject[CISRulesetsRuleExpression].(string))
		opt.SetAction(rulesObject[CISRulesetsRuleAction].(string))
		opt.SetDescription(rulesObject[CISRulesetsRuleActionDescription].(string))
		opt.SetEnabled(rulesObject[CISRulesetsRuleActionEnabled].(bool))
		opt.SetRef(rulesObject[CISRulesetsRuleRef].(string))

		position := rulesetsv1.Position{}
		if reflect.ValueOf(rulesObject[CISRulesetsRulePosition]).IsNil() {
			position = expandCISRulesetsRulesPositions(rulesObject[CISRulesetsRulePosition])
		}
		opt.SetPosition(&position)

		actionParameterObj := rulesetsv1.ActionParameters{}
		if len(rulesObject[CISRulesetsRuleActionParameters].(*schema.Set).List()) != 0 {
			actionParameterObj = expandCISRulesetsRulesActionParameters(rulesObject[CISRulesetsRuleActionParameters])
		}
		opt.SetActionParameters(&actionParameterObj)

		result, resp, err := sess.CreateZoneRulesetRule(opt)

		if err != nil {
			return fmt.Errorf("[ERROR] Error while creating the zone Rule %s", resp)
		}
		len_rules := len(result.Result.Rules)
		opt.SetID(*result.Result.Rules[len_rules-1].ID)

		d.SetId(dataSourceCISRulesetsRuleCheckID(d, *result.Result.Rules[len_rules-1].ID))

	} else {
		opt := sess.NewCreateInstanceRulesetRuleOptions(rulesetId)

		rulesObject := d.Get(CISRulesetsRule).([]interface{})[0].(map[string]interface{})

		opt.SetRulesetID(rulesetId)
		opt.SetExpression(rulesObject[CISRulesetsRuleExpression].(string))
		opt.SetAction(rulesObject[CISRulesetsRuleAction].(string))
		opt.SetDescription(rulesObject[CISRulesetsRuleActionDescription].(string))
		opt.SetEnabled(rulesObject[CISRulesetsRuleActionEnabled].(bool))
		opt.SetRef(rulesObject[CISRulesetsRuleRef].(string))

		position := rulesetsv1.Position{}
		if reflect.ValueOf(rulesObject[CISRulesetsRulePosition]).IsNil() {
			position = expandCISRulesetsRulesPositions(rulesObject[CISRulesetsRulePosition])
		}
		opt.SetPosition(&position)

		actionParameterObj := rulesetsv1.ActionParameters{}
		if len(rulesObject[CISRulesetsRuleActionParameters].(*schema.Set).List()) != 0 {
			actionParameterObj = expandCISRulesetsRulesActionParameters(rulesObject[CISRulesetsRuleActionParameters])
		}
		opt.SetActionParameters(&actionParameterObj)

		result, resp, err := sess.CreateInstanceRulesetRule(opt)

		if err != nil {
			return fmt.Errorf("[ERROR] Error while creating the instance Rule %s", resp)
		}

		len_rules := len(result.Result.Rules)
		opt.SetID(*result.Result.Rules[len_rules-1].ID)

		d.SetId(dataSourceCISRulesetsRuleCheckID(d, *result.Result.Rules[len_rules-1].ID))
	}
	return nil
}

func ResourceIBMCISRulesetRuleRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func ResourceIBMCISRulesetRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}

	ruleId, rulesetId, zoneId, crn, err := flex.ConvertTfToCisFourVar(d.Id())
	sess.Crn = core.StringPtr(crn)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)

		opt := sess.NewUpdateZoneRulesetRuleOptions(rulesetId, ruleId)

		rulesetsRuleObject := d.Get(CISRulesetsRule).([]interface{})[0].(map[string]interface{})
		opt.SetDescription(rulesetsRuleObject[CISRulesetsDescription].(string))
		opt.SetAction(rulesetsRuleObject[CISRulesetsRuleAction].(string))
		actionParameters := expandCISRulesetsRulesActionParameters(rulesetsRuleObject[CISRulesetsRuleActionParameters])
		opt.SetActionParameters(&actionParameters)
		opt.SetEnabled(rulesetsRuleObject[CISRulesetsRuleActionEnabled].(bool))
		opt.SetExpression(rulesetsRuleObject[CISRulesetsRuleExpression].(string))
		opt.SetRef(rulesetsRuleObject[CISRulesetsRuleRef].(string))
		position := expandCISRulesetsRulesPositions(rulesetsRuleObject[CISRulesetsRulePosition])
		opt.SetPosition(&position)

		opt.SetRulesetID(rulesetId)
		opt.SetRuleID(ruleId)
		opt.SetID(ruleId)

		_, _, err := sess.UpdateZoneRulesetRule(opt)

		if err != nil {
			return fmt.Errorf("[ERROR] Error while updating the zone Ruleset %s", err)
		}

		d.SetId(dataSourceCISRulesetsRuleCheckID(d, ruleId))

	} else {
		opt := sess.NewUpdateInstanceRulesetRuleOptions(rulesetId, ruleId)

		rulesetsRuleObject := d.Get(CISRulesetsRule).([]interface{})[0].(map[string]interface{})
		opt.SetDescription(rulesetsRuleObject[CISRulesetsDescription].(string))
		opt.SetAction(rulesetsRuleObject[CISRulesetsRuleAction].(string))
		actionParameters := expandCISRulesetsRulesActionParameters(rulesetsRuleObject[CISRulesetsRuleActionParameters])
		opt.SetActionParameters(&actionParameters)
		opt.SetEnabled(rulesetsRuleObject[CISRulesetsRuleActionEnabled].(bool))
		opt.SetExpression(rulesetsRuleObject[CISRulesetsRuleExpression].(string))
		opt.SetRef(rulesetsRuleObject[CISRulesetsRuleAction].(string))
		position := expandCISRulesetsRulesPositions(rulesetsRuleObject[CISRulesetsRulePosition])
		opt.SetPosition(&position)

		opt.SetRulesetID(rulesetId)
		opt.SetRuleID(ruleId)
		opt.SetID(ruleId)

		_, _, err := sess.UpdateInstanceRulesetRule(opt)

		if err != nil {
			return fmt.Errorf("[ERROR] Error while updating the zone Ruleset %s", err)
		}

		d.SetId(dataSourceCISRulesetsRuleCheckID(d, ruleId))
	}
	return nil
}

func ResourceIBMCISRulesetRuleDelete(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}

	ruleId, rulesetId, zoneId, crn, err := flex.ConvertTfToCisFourVar(d.Id())
	sess.Crn = core.StringPtr(crn)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)
		opt := sess.NewDeleteZoneRulesetRuleOptions(rulesetId, ruleId)
		_, res, err := sess.DeleteZoneRulesetRule(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error deleting the zone ruleset rule %s:%s", err, res)
		}
	} else {
		opt := sess.NewDeleteInstanceRulesetRuleOptions(rulesetId, ruleId)
		_, res, err := sess.DeleteInstanceRulesetRule(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error deleting the Instance ruleset rule %s:%s", err, res)
		}
	}

	d.SetId("")
	return nil
}

func dataSourceCISRulesetsRuleCheckID(d *schema.ResourceData, ruleId string) string {
	return ruleId + ":" + d.Get(CISRulesetsId).(string) + ":" + d.Get(cisDomainID).(string) + ":" + d.Get(cisID).(string)
}
