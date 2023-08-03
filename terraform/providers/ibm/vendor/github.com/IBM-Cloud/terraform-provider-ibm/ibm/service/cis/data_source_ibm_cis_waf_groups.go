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

const cisWAFGroups = "waf_groups"

func DataSourceIBMCISWAFGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISWAFGroupsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_waf_groups",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFGroupPackageID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "WAF Rule package id",
			},
			cisWAFGroups: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule group id",
						},
						cisWAFGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule group id",
						},
						cisWAFGroupMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule group mode on/off",
						},
						cisWAFGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule group name",
						},
						cisWAFGroupDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF Rule group description",
						},
						cisWAFGroupRulesCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "WAF Rule group rules count",
						},
						cisWAFGroupModifiedRulesCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "WAF Rule group modified rules count",
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISWAFGroupsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISWAFGroupsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_waf_groups",
		Schema:       validateSchema}
	return &iBMCISWAFGroupsValidator
}

func dataSourceIBMCISWAFGroupsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFGroupClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	packageID, _, _, _ := flex.ConvertTfToCisThreeVar(d.Get(cisWAFGroupPackageID).(string))

	opt := cisClient.NewListWafRuleGroupsOptions(packageID)
	opt.SetPage(1)
	opt.SetPerPage(100)
	result, resp, err := cisClient.ListWafRuleGroups(opt)
	if err != nil {
		log.Printf("List waf rule groups failed: %s\n", resp)
		return err
	}
	wafGroups := []interface{}{}
	for _, i := range result.Result {
		waf := map[string]interface{}{}
		waf["id"] = flex.ConvertCisToTfFourVar(*i.ID, packageID, zoneID, crn)
		waf[cisWAFGroupID] = *i.ID
		waf[cisWAFGroupName] = *i.Name
		waf[cisWAFGroupDesc] = *i.Description
		waf[cisWAFGroupMode] = *i.Mode
		waf[cisWAFGroupModifiedRulesCount] = *i.ModifiedRulesCount
		waf[cisWAFGroupRulesCount] = *i.RulesCount
		wafGroups = append(wafGroups, waf)
	}
	d.SetId(dataSourceIBMCISWAFGroupID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFGroupPackageID, packageID)
	d.Set(cisWAFGroups, wafGroups)
	return nil
}

func dataSourceIBMCISWAFGroupID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
