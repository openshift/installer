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

func DataSourceIbmSccReportResources() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportResourcesRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the resource.",
			},
			"resource_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the resource.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the account owning a resource.",
			},
			"component_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of component.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The compliance status value.",
			},
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field sorts resources by using a valid sort field. To learn more, see [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).",
			},
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The page reference.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for the first and next page.",
						},
					},
				},
			},
			"home_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the home account.",
			},
			"resources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of resource evaluation summaries that are on the page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the report.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource CRN.",
						},
						"resource_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource name.",
						},
						"component_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the component.",
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment.",
						},
						"account": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The account that is associated with a report.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account name.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account type.",
									},
								},
							},
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allowed values of an aggregated status for controls, specifications, assessments, and resources.",
						},
						"total_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of evaluations.",
						},
						"pass_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of passed evaluations.",
						},
						"failure_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of failed evaluations.",
						},
						"error_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of evaluations that started, but did not finish, and ended with errors.",
						},
						"completed_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of completed evaluations.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccReportResourcesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listReportResourcesOptions := &securityandcompliancecenterapiv3.ListReportResourcesOptions{}

	listReportResourcesOptions.SetReportID(d.Get("report_id").(string))
	listReportResourcesOptions.SetInstanceID(d.Get("instance_id").(string))
	if _, ok := d.GetOk("id"); ok {
		listReportResourcesOptions.SetID(d.Get("id").(string))
	}
	if _, ok := d.GetOk("resource_name"); ok {
		listReportResourcesOptions.SetResourceName(d.Get("resource_name").(string))
	}
	if _, ok := d.GetOk("account_id"); ok {
		listReportResourcesOptions.SetAccountID(d.Get("account_id").(string))
	}
	if _, ok := d.GetOk("component_id"); ok {
		listReportResourcesOptions.SetComponentID(d.Get("component_id").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		listReportResourcesOptions.SetStatus(d.Get("status").(string))
	}
	if _, ok := d.GetOk("sort"); ok {
		listReportResourcesOptions.SetSort(d.Get("sort").(string))
	}

	var pager *securityandcompliancecenterapiv3.ReportResourcesPager
	pager, err = resultsClient.NewReportResourcesPager(listReportResourcesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ReportResourcesPager.GetAll() failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ReportResourcesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIbmSccReportResourcesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmSccReportResourcesResourceToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("resources", mapSlice); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting resources %s", err))
	}

	return nil
}

// dataSourceIbmSccReportResourcesID returns a reasonable ID for the list.
func dataSourceIbmSccReportResourcesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportResourcesResourceToMap(model *securityandcompliancecenterapiv3.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReportID != nil {
		modelMap["report_id"] = model.ReportID
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.ResourceName != nil {
		modelMap["resource_name"] = model.ResourceName
	}
	if model.ComponentID != nil {
		modelMap["component_id"] = model.ComponentID
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Account != nil {
		accountMap, err := dataSourceIbmSccReportResourcesAccountToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.TotalCount != nil {
		modelMap["total_count"] = flex.IntValue(model.TotalCount)
	}
	if model.PassCount != nil {
		modelMap["pass_count"] = flex.IntValue(model.PassCount)
	}
	if model.FailureCount != nil {
		modelMap["failure_count"] = flex.IntValue(model.FailureCount)
	}
	if model.ErrorCount != nil {
		modelMap["error_count"] = flex.IntValue(model.ErrorCount)
	}
	if model.CompletedCount != nil {
		modelMap["completed_count"] = flex.IntValue(model.CompletedCount)
	}
	return modelMap, nil
}

func dataSourceIbmSccReportResourcesAccountToMap(model *securityandcompliancecenterapiv3.Account) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}
