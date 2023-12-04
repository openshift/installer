// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
)

func DataSourceIBMIAMAccessGroupTemplateAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIAMAccessGroupTemplateAssignmentRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enterprise account ID.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by Template Id.",
			},
			"template_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by Template Version.",
			},
			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by the assignment target.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by the assignment status.",
			},
			"transaction_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional transaction id for the request.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of items returned in the response.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Index of the first item returned in the response.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of items matching the query.",
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A string containing the link’s URL.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A string containing the link’s URL.",
						},
					},
				},
			},
			"assignments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of template assignments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the assignment.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the account that the assignment belongs to.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the template that the assignment is based on.",
						},
						"template_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the template that the assignment is based on.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the entity that the assignment applies to.",
						},
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the entity that the assignment applies to.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation that the assignment applies to (e.g. 'assign', 'update', 'remove').",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the assignment (e.g. 'accepted', 'in_progress', 'succeeded', 'failed', 'superseded').",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the assignment resource.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time when the assignment was created.",
						},
						"created_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user or system that created the assignment.",
						},
						"last_modified_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time when the assignment was last updated.",
						},
						"last_modified_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user or system that last updated the assignment.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMAccessGroupTemplateAssignmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listAssignmentsOptions := &iamaccessgroupsv2.ListAssignmentsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	listAssignmentsOptions.SetAccountID(userDetails.UserAccount)
	if _, ok := d.GetOk("template_id"); ok {
		listAssignmentsOptions.SetTemplateID(d.Get("template_id").(string))
	}
	if _, ok := d.GetOk("template_version"); ok {
		listAssignmentsOptions.SetTemplateVersion(d.Get("template_version").(string))
	}
	if _, ok := d.GetOk("target"); ok {
		listAssignmentsOptions.SetTarget(d.Get("target").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		listAssignmentsOptions.SetStatus(d.Get("status").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		listAssignmentsOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	listTemplateAssignmentResponse, response, err := iamAccessGroupsClient.ListAssignmentsWithContext(context, listAssignmentsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListAssignmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListAssignmentsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIAMAccessGroupTemplateAssignmentID(d))

	if err = d.Set("limit", flex.IntValue(listTemplateAssignmentResponse.Limit)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting limit: %s", err))
	}

	if err = d.Set("offset", flex.IntValue(listTemplateAssignmentResponse.Offset)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting offset: %s", err))
	}

	if err = d.Set("total_count", flex.IntValue(listTemplateAssignmentResponse.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	first := []map[string]interface{}{}
	if listTemplateAssignmentResponse.First != nil {
		modelMap, err := dataSourceIBMIAMAccessGroupTemplateAssignmentHrefStructToMap(listTemplateAssignmentResponse.First)
		if err != nil {
			return diag.FromErr(err)
		}
		first = append(first, modelMap)
	}
	if err = d.Set("first", first); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting first %s", err))
	}

	last := []map[string]interface{}{}
	if listTemplateAssignmentResponse.Last != nil {
		modelMap, err := dataSourceIBMIAMAccessGroupTemplateAssignmentHrefStructToMap(listTemplateAssignmentResponse.Last)
		if err != nil {
			return diag.FromErr(err)
		}
		last = append(last, modelMap)
	}
	if err = d.Set("last", last); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last %s", err))
	}

	assignments := []map[string]interface{}{}
	if listTemplateAssignmentResponse.Assignments != nil {
		for _, modelItem := range listTemplateAssignmentResponse.Assignments {
			modelMap, err := dataSourceIBMIAMAccessGroupTemplateAssignmentTemplateAssignmentResponseToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			assignments = append(assignments, modelMap)
		}
	}
	if err = d.Set("assignments", assignments); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting assignments %s", err))
	}

	return nil
}

// dataSourceIBMIAMAccessGroupTemplateAssignmentID returns a reasonable ID for the list.
func dataSourceIBMIAMAccessGroupTemplateAssignmentID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIAMAccessGroupTemplateAssignmentHrefStructToMap(model *iamaccessgroupsv2.HrefStruct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateAssignmentTemplateAssignmentResponseToMap(model *iamaccessgroupsv2.TemplateAssignmentResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["account_id"] = model.AccountID
	modelMap["template_id"] = model.TemplateID
	modelMap["template_version"] = model.TemplateVersion
	modelMap["target_type"] = model.TargetType
	modelMap["target"] = model.Target
	modelMap["operation"] = model.Operation
	modelMap["status"] = model.Status
	modelMap["href"] = model.Href
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["created_by_id"] = model.CreatedByID
	modelMap["last_modified_at"] = model.LastModifiedAt.String()
	modelMap["last_modified_by_id"] = model.LastModifiedByID
	return modelMap, nil
}
