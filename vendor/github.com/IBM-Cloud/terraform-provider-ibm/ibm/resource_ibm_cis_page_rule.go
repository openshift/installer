// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"strconv"

	"github.com/IBM/go-sdk-core/v4/core"
	cispagerulev1 "github.com/IBM/networking-go-sdk/pageruleapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISPageRule                       = "ibm_cis_page_rule"
	cisPageRuleID                        = "rule_id"
	cisPageRuleTargets                   = "targets"
	cisPageRuleTargetsConstraint         = "constraint"
	cisPageRuleTargetsConstraintOperator = "operator"
	cisPageRuleTargetsConstraintValue    = "value"
	cisPageRuleTargetsTarget             = "target"
	cisPageRuleActions                   = "actions"
	cisPageRuleActionsID                 = "id"
	cisPageRuleActionsValue              = "value"
	cisPageRuleActionsValueURL           = "url"
	cisPageRuleActionsValueStatusCode    = "status_code"
	cisPageRulePriority                  = "priority"
	cisPageRuleStatus                    = "status"
	cisPageRuleActionsIDForwardingURL    = "forwarding_url"
	cisPageRuleActionsIDEdgeCacheTTL     = "edge_cache_ttl"
	cisPageRuleActionsIDBrowserCacheTTL  = "browser_cache_ttl"
	cisPageRuleActionsIDDisableSecurity  = "disable_security"
	cisPageRuleActionsIDAlwaysUseHTTPS   = "always_use_https"
)

func resourceIBMCISPageRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCISPageRuleCreate,
		Read:     resourceCISPageRuleRead,
		Update:   resourceCISPageRuleUpdate,
		Delete:   resourceCISPageRuleDelete,
		Exists:   resourceCISPageRuleExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisPageRuleID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisPageRulePriority: {
				Type:        schema.TypeInt,
				Description: "Page rule priority",
				Optional:    true,
				Default:     1,
			},
			cisPageRuleStatus: {
				Type:        schema.TypeString,
				Description: "Page Rule status",
				Optional:    true,
				Default:     "disabled",
				ValidateFunc: InvokeValidator(
					ibmCISPageRule, cisPageRuleStatus),
			},
			cisPageRuleTargets: {
				Type:        schema.TypeSet,
				Description: "Page rule targets",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisPageRuleTargetsTarget: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Page rule target url",
						},
						cisPageRuleTargetsConstraint: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Page rule constraint",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisPageRuleTargetsConstraintOperator: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Constraint operator",
									},
									cisPageRuleTargetsConstraintValue: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Constraint value",
									},
								},
							},
						},
					},
				},
			},
			cisPageRuleActions: {
				Type:        schema.TypeSet,
				Description: "Page rule actions",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisPageRuleActionsID: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Page rule target url",
							ValidateFunc: InvokeValidator(
								ibmCISPageRule, cisPageRuleActionsID),
						},
						cisPageRuleActionsValue: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Page rule target url",
						},
						cisPageRuleActionsValueURL: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Page rule actions value url",
						},
						cisPageRuleActionsValueStatusCode: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Page rule actions status code",
						},
					},
				},
			},
		},
	}
}

func resourceCISPageRuleValidator() *ResourceValidator {
	actions := "disable_security, always_use_https, always_online, ssl, browser_cache_ttl, " +
		"security_level, cache_level, edge_cache_ttl, bypass_cache_on_cookie, " +
		"browser_check, server_side_exclude, serve_stale_content, email_obfuscation, " +
		"automatic_https_rewrites, opportunistic_encryption, ip_geolocation, " +
		"explicit_cache_control, cache_deception_armor, waf, forwarding_url, " +
		"host_header_override, resolve_override, cache_on_cookie, disable_apps, " +
		"disable_performance, image_load_optimization, origin_error_page_pass_thru, " +
		"response_buffering, image_size_optimization, script_load_optimization, " +
		"true_client_ip_header, sort_query_string_for_cache, respect_strong_etag"
	status := "active, disabled"
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisPageRuleActionsID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              actions})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisPageRuleStatus,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              status})
	cisPageRuleValidator := ResourceValidator{ResourceName: ibmCISPageRule, Schema: validateSchema}
	return &cisPageRuleValidator
}

func resourceCISPageRuleCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisPageRuleClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	targets := expandCISPageRuleTargets(d.Get(cisPageRuleTargets))
	actions := expandCISPageRuleActions(d.Get(cisPageRuleActions))

	opt := cisClient.NewCreatePageRuleOptions()
	opt.SetTargets(targets)
	opt.SetActions(actions)
	if value, ok := d.GetOk(cisPageRulePriority); ok {
		opt.SetPriority(int64(value.(int)))
	}
	if value, ok := d.GetOk(cisPageRuleStatus); ok {
		opt.SetStatus(value.(string))
	}

	result, response, err := cisClient.CreatePageRule(opt)
	if err != nil {
		log.Printf("Create page rule failed: %v", response)
		return err
	}
	d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return resourceCISPageRuleRead(d, meta)
}
func resourceCISPageRuleRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisPageRuleClientSession()
	if err != nil {
		return err
	}

	ruleID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	opt := cisClient.NewGetPageRuleOptions(ruleID)
	result, response, err := cisClient.GetPageRule(opt)
	if err != nil {
		log.Printf("Get page rule failed: %v", response)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisPageRuleID, result.Result.ID)
	d.Set(cisPageRulePriority, result.Result.Priority)
	d.Set(cisPageRuleStatus, result.Result.Status)
	d.Set(cisPageRuleTargets, flattenCISPageRuleTargets(result.Result.Targets))
	d.Set(cisPageRuleActions, flattenCISPageRuleActions(result.Result.Actions))
	return nil
}
func resourceCISPageRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisPageRuleClientSession()
	if err != nil {
		return err
	}

	ruleID, zoneID, crn, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	if d.HasChange(cisPageRuleTargets) ||
		d.HasChange(cisPageRuleActions) ||
		d.HasChange(cisPageRulePriority) ||
		d.HasChange(cisPageRuleStatus) {

		targets := expandCISPageRuleTargets(d.Get(cisPageRuleTargets))
		actions := expandCISPageRuleActions(d.Get(cisPageRuleActions))

		opt := cisClient.NewUpdatePageRuleOptions(ruleID)
		opt.SetTargets(targets)
		opt.SetActions(actions)
		if value, ok := d.GetOk(cisPageRulePriority); ok {
			opt.SetPriority(int64(value.(int)))
		}
		if value, ok := d.GetOk(cisPageRuleStatus); ok {
			opt.SetStatus(value.(string))
		}

		_, response, err := cisClient.UpdatePageRule(opt)
		if err != nil {
			log.Printf("Update page rule failed: %v", response)
			return err
		}
	}
	return resourceCISPageRuleRead(d, meta)
}

func resourceCISPageRuleDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisPageRuleClientSession()
	if err != nil {
		return err
	}

	ruleID, zoneID, crn, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	opt := cisClient.NewDeletePageRuleOptions(ruleID)
	_, response, err := cisClient.DeletePageRule(opt)
	if err != nil {
		log.Printf("Delete page rule failed: %v", response)
		return err
	}
	return nil
}

func resourceCISPageRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisPageRuleClientSession()
	if err != nil {
		return false, err
	}

	ruleID, zoneID, crn, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	opt := cisClient.NewGetPageRuleOptions(ruleID)
	_, response, err := cisClient.GetPageRule(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Page rule does not exist.")
			return false, nil
		}
		log.Printf("Get page rule failed: %v", response)
		return false, err
	}
	return true, nil
}

func expandCISPageRuleTargets(targets interface{}) []cispagerulev1.TargetsItem {
	targetsInput := targets.(*schema.Set).List()
	targetsOuptut := make([]cispagerulev1.TargetsItem, 0)
	for _, instance := range targetsInput {
		targetsItem := instance.(map[string]interface{})
		targetsTarget := targetsItem[cisPageRuleTargetsTarget].(string)
		targetsConstraint := targetsItem[cisPageRuleTargetsConstraint].([]interface{})[0].(map[string]interface{})
		targetsConstraintOperator := targetsConstraint[cisPageRuleTargetsConstraintOperator].(string)
		targetsConstraintValue := targetsConstraint[cisPageRuleTargetsConstraintValue].(string)
		targetsConstraintOpt := cispagerulev1.TargetsItemConstraint{
			Operator: &targetsConstraintOperator,
			Value:    &targetsConstraintValue,
		}
		targetItemOpt := cispagerulev1.TargetsItem{
			Target:     &targetsTarget,
			Constraint: &targetsConstraintOpt,
		}
		targetsOuptut = append(targetsOuptut, targetItemOpt)
	}
	return targetsOuptut
}

func expandCISPageRuleActions(actions interface{}) []cispagerulev1.PageRulesBodyActionsItemIntf {
	actionsInput := actions.(*schema.Set).List()

	actionsOutput := make([]cispagerulev1.PageRulesBodyActionsItemIntf, 0)
	for _, action := range actionsInput {
		instance := action.(map[string]interface{})
		id := instance[cisPageRuleActionsID].(string)
		var value interface{}
		switch id {
		case cisPageRuleActionsIDDisableSecurity,
			cisPageRuleActionsIDAlwaysUseHTTPS:
			actionItem := &cispagerulev1.PageRulesBodyActionsItem{
				ID: &id,
			}
			actionsOutput = append(actionsOutput, actionItem)
			break
		case cisPageRuleActionsIDBrowserCacheTTL,
			cisPageRuleActionsIDEdgeCacheTTL:
			valueStr := instance[cisPageRuleActionsValue].(string)
			value, _ = strconv.ParseInt(valueStr, 10, 64)
			actionItem := &cispagerulev1.PageRulesBodyActionsItem{
				ID:    &id,
				Value: &value,
			}
			actionsOutput = append(actionsOutput, actionItem)
			break
		case cisPageRuleActionsIDForwardingURL:
			forwardingURL := instance[cisPageRuleActionsValueURL].(string)
			statusCode := instance[cisPageRuleActionsValueStatusCode].(int)
			value = cispagerulev1.ActionsForwardingUrlValue{
				URL:        &forwardingURL,
				StatusCode: core.Int64Ptr(int64(statusCode)),
			}
			actionItem := &cispagerulev1.PageRulesBodyActionsItem{
				ID:    &id,
				Value: &value,
			}
			actionsOutput = append(actionsOutput, actionItem)
			break
		default:
			value = instance[cisPageRuleActionsValue]
			actionItem := &cispagerulev1.PageRulesBodyActionsItem{
				ID:    &id,
				Value: &value,
			}
			actionsOutput = append(actionsOutput, actionItem)
		}
	}
	return actionsOutput
}

func flattenCISPageRuleTargets(targets []cispagerulev1.TargetsItem) interface{} {
	targetsOutput := make([]interface{}, 0)

	for _, item := range targets {
		targetItemOutput := map[string]interface{}{}
		constraints := []interface{}{}
		constraint := map[string]interface{}{}
		// flatten constraint
		constraint[cisPageRuleTargetsConstraintOperator] = *item.Constraint.Operator
		constraint[cisPageRuleTargetsConstraintValue] = *item.Constraint.Value
		constraints = append(constraints, constraint)

		// flatten target item
		targetItemOutput[cisPageRuleTargetsConstraint] = constraints
		targetItemOutput[cisPageRuleTargetsTarget] = *item.Target

		targetsOutput = append(targetsOutput, targetItemOutput)
	}
	return targetsOutput
}

func flattenCISPageRuleActions(actions []cispagerulev1.PageRulesBodyActionsItemIntf) interface{} {
	actionsOutput := make([]interface{}, 0)

	for _, instance := range actions {
		actionItemOutput := map[string]interface{}{}
		item := instance.(*cispagerulev1.PageRulesBodyActionsItem)
		actionItemOutput[cisPageRuleActionsID] = *item.ID
		if *item.ID == cisPageRuleActionsIDForwardingURL {
			value := item.Value.(map[string]interface{})
			actionItemOutput[cisPageRuleActionsValueURL] = value[cisPageRuleActionsValueURL]
			actionItemOutput[cisPageRuleActionsValueStatusCode] = value[cisPageRuleActionsValueStatusCode]
		} else {
			actionItemOutput[cisPageRuleActionsValue] = item.Value
		}
		actionsOutput = append(actionsOutput, actionItemOutput)
	}
	return actionsOutput
}
