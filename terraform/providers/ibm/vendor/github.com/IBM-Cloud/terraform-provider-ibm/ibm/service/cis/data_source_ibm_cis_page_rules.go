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

const (
	cisPageRules = "cis_page_rules"
)

func DataSourceIBMCISPageRules() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISPageRulesRead,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Zone CRN",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_page_rules",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "DNS Zone ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisPageRules: {
				Type:        schema.TypeList,
				Description: "Collection of page rules detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Page rule identifier",
						},
						cisPageRuleID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisPageRulePriority: {
							Type:        schema.TypeInt,
							Description: "Page rule priority",
							Computed:    true,
						},
						cisPageRuleStatus: {
							Type:        schema.TypeString,
							Description: "Page Rule status",
							Computed:    true,
						},
						cisPageRuleTargets: {
							Type:        schema.TypeList,
							Description: "Page rule targets",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisPageRuleTargetsTarget: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page rule target url",
									},
									cisPageRuleTargetsConstraint: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Page rule constraint",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												cisPageRuleTargetsConstraintOperator: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Constraint operator",
												},
												cisPageRuleTargetsConstraintValue: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Constraint value",
												},
											},
										},
									},
								},
							},
						},
						cisPageRuleActions: {
							Type:        schema.TypeList,
							Description: "Page rule actions",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisPageRuleActionsID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page rule target url",
									},
									cisPageRuleActionsValue: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page rule target url",
									},
									cisPageRuleActionsValueURL: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page rule actions value url",
									},
									cisPageRuleActionsValueStatusCode: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Page rule actions status code",
									},
									cisPageRuleActionsMinifyCSS: {
										Type:        schema.TypeString,
										Description: "Minify CSS value",
										Computed:    true,
									},
									cisPageRuleActionsMinifyHTML: {
										Type:        schema.TypeString,
										Description: "Minify HTML value",
										Computed:    true,
									},
									cisPageRuleActionsMinifyJS: {
										Type:        schema.TypeString,
										Description: "Minify JS value",
										Computed:    true,
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
func DataSourceIBMCISPageRulesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISPageRulesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_page_rules",
		Schema:       validateSchema}
	return &iBMCISPageRulesValidator
}
func dataSourceIBMCISPageRulesRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisPageRuleClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	opt := sess.NewListPageRulesOptions()

	result, resp, err := sess.ListPageRules(opt)
	if err != nil {
		log.Printf("Error listing page rules detail: %s", resp)
		return err
	}

	pageRules := make([]map[string]interface{}, 0)
	for _, instance := range result.Result {
		pageRule := map[string]interface{}{}
		pageRule["id"] = flex.ConvertCisToTfThreeVar(*instance.ID, zoneID, crn)
		pageRule[cisPageRuleID] = *instance.ID
		pageRule[cisPageRulePriority] = *instance.Priority
		pageRule[cisPageRuleStatus] = *instance.Status
		pageRule[cisPageRuleTargets] = flattenCISPageRuleTargets(instance.Targets)
		pageRule[cisPageRuleActions] = flattenCISPageRuleActions(instance.Actions)
		pageRules = append(pageRules, pageRule)
	}
	d.SetId(dataSourceIBMCISPageRulesID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisPageRules, pageRules)
	return nil
}

func dataSourceIBMCISPageRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
