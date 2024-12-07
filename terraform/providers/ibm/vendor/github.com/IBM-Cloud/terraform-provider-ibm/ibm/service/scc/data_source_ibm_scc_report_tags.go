// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccReportTags() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportTagsRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of different types of tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The collection of user tags.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"access": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The collection of access tags.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The collection of service tags.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccReportTagsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportTagsOptions := &securityandcompliancecenterapiv3.GetReportTagsOptions{}

	getReportTagsOptions.SetReportID(d.Get("report_id").(string))
	getReportTagsOptions.SetInstanceID(d.Get("instance_id").(string))

	reportTags, response, err := resultsClient.GetReportTagsWithContext(context, getReportTagsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportTagsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportTagsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccReportTagsID(d))

	tags := []map[string]interface{}{}
	if reportTags.Tags != nil {
		modelMap, err := dataSourceIbmSccReportTagsTagsToMap(reportTags.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
		tags = append(tags, modelMap)
	}
	if err = d.Set("tags", tags); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting tags %s", err))
	}

	return nil
}

// dataSourceIbmSccReportTagsID returns a reasonable ID for the list.
func dataSourceIbmSccReportTagsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportTagsTagsToMap(model *securityandcompliancecenterapiv3.Tags) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.User != nil {
		modelMap["user"] = model.User
	}
	if model.Access != nil {
		modelMap["access"] = model.Access
	}
	if model.Service != nil {
		modelMap["service"] = model.Service
	}
	return modelMap, nil
}
