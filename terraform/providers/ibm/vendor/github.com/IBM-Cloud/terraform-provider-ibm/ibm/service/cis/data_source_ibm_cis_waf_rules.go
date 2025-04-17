// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cisWAFRules = "waf_rules"

func DataSourceIBMCISWAFRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISWAFRuleRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_waf_rules",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "CISzone - Domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFRulePackageID: {
				Type:             schema.TypeString,
				Description:      "WAF rule package id",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFRules: {
				Type:        schema.TypeList,
				Description: "Collection of WAF Rules",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule id",
						},
						cisWAFRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule id",
						},
						cisWAFRulePackageID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Package id",
						},
						cisWAFRuleMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIS WAF Rule mode",
						},
						cisWAFRuleDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIS WAF Rule descriptions",
						},
						cisWAFRulePriority: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIS WAF Rule Priority",
						},
						cisWAFRuleGroup: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CIS WAF Rule group",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisWAFRuleGroupID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "waf rule group id",
									},
									cisWAFRuleGroupName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "waf rule group name",
									},
								},
							},
						},
						cisWAFRuleAllowedModes: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CIS WAF Rule allowed modes",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISWAFRulesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISWAFRulesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_waf_rules",
		Schema:       validateSchema}
	return &iBMCISWAFRulesValidator
}
func dataSourceIBMCISWAFRuleRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFRuleClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	packageID, _, _, _ := flex.ConvertTfToCisThreeVar(d.Get(cisWAFRulePackageID).(string))

	opt := cisClient.NewListWafRulesOptions(packageID)
	opt.SetPage(1)
	opt.SetPerPage(1000)
	result, response, err := cisClient.ListWafRules(opt)
	if err != nil {
		log.Printf("List waf rules failed %s\n", response)
		return err
	}
	rules := []interface{}{}
	for _, i := range result.Result {

		groups := []interface{}{}
		group := map[string]interface{}{}
		group[cisWAFRuleGroupID] = *i.Group.ID
		group[cisWAFRuleGroupName] = *i.Group.Name
		groups = append(groups, group)

		rule := map[string]interface{}{}
		rule["id"] = flex.ConvertCisToTfFourVar(*i.ID, *i.PackageID, zoneID, crn)
		rule[cisWAFRuleID] = *i.ID
		rule[cisWAFRulePackageID] = *i.PackageID
		rule[cisWAFRuleMode] = *i.Mode
		rule[cisWAFRuleDesc] = *i.Description
		rule[cisWAFRulePriority] = *i.Priority
		rule[cisWAFRuleGroup] = groups
		rule[cisWAFRuleAllowedModes] = flex.FlattenStringList(i.AllowedModes)

		rules = append(rules, rule)
	}
	d.SetId(dataSourceIBMCISWAFRulesID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFRulePackageID, packageID)
	d.Set(cisWAFRules, rules)
	return nil
}

func dataSourceIBMCISWAFRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
