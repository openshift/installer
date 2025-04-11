// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccReportRule() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportRuleRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of a rule in a report.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule type.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule description.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule version.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule account ID.",
			},
			"created_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user who created the rule.",
			},
			"updated_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was updated.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user who updated the rule.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rule labels.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}

func dataSourceIbmSccReportRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportRuleOptions := &securityandcompliancecenterapiv3.GetReportRuleOptions{}

	getReportRuleOptions.SetReportID(d.Get("report_id").(string))
	getReportRuleOptions.SetRuleID(d.Get("rule_id").(string))
	getReportRuleOptions.SetInstanceID(d.Get("instance_id").(string))

	ruleInfo, response, err := resultsClient.GetReportRuleWithContext(context, getReportRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(*ruleInfo.ID)

	if err = d.Set("id", ruleInfo.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting id: %s", err))
	}

	if err = d.Set("type", ruleInfo.Type); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting type: %s", err))
	}

	if err = d.Set("description", ruleInfo.Description); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting description: %s", err))
	}

	if err = d.Set("version", ruleInfo.Version); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting version: %s", err))
	}

	if err = d.Set("account_id", ruleInfo.AccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting account_id: %s", err))
	}

	if err = d.Set("created_on", flex.DateTimeToString(ruleInfo.CreatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
	}

	if err = d.Set("created_by", ruleInfo.CreatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_on", flex.DateTimeToString(ruleInfo.UpdatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
	}

	if err = d.Set("updated_by", ruleInfo.UpdatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
	}

	return nil
}
