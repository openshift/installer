// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package logs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsDataAccessRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsDataAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"logs_data_access_rules_id": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Array of data access rule IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_access_rules": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data Access Rule details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data Access Rule ID.",
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data Access Rule Display Name.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional Data Access Rule Description.",
						},
						"filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of filters that the Data Access Rule is composed of.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entity_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Filter's Entity Type.",
									},
									"expression": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Filter's Expression.",
									},
								},
							},
						},
						"default_expression": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default expression to use when no filter matches the query.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsDataAccessRulesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_data_access_rules", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	listDataAccessRulesOptions := &logsv0.ListDataAccessRulesOptions{}

	if _, ok := d.GetOk("logs_data_access_rules_id"); ok {
		var id []strfmt.UUID
		for _, v := range d.Get("logs_data_access_rules_id").([]interface{}) {
			idItem := strfmt.UUID(v.(string))
			id = append(id, idItem)
		}
		listDataAccessRulesOptions.SetID(id)
	}

	dataAccessRuleCollection, _, err := logsClient.ListDataAccessRulesWithContext(context, listDataAccessRulesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListDataAccessRulesWithContext failed: %s", err.Error()), "(Data) ibm_logs_data_access_rules", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsDataAccessRulesID(d))

	dataAccessRules := []map[string]interface{}{}
	if dataAccessRuleCollection.DataAccessRules != nil {
		for _, modelItem := range dataAccessRuleCollection.DataAccessRules {
			modelMap, err := DataSourceIbmLogsDataAccessRulesDataAccessRuleToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_data_access_rules", "read", "data_access_rules-to-map").GetDiag()
			}
			dataAccessRules = append(dataAccessRules, modelMap)
		}
	}
	if err = d.Set("data_access_rules", dataAccessRules); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting data_access_rules: %s", err), "(Data) ibm_logs_data_access_rules", "read", "set-data_access_rules").GetDiag()
	}

	return nil
}

// dataSourceIbmLogsDataAccessRulesID returns a reasonable ID for the list.
func dataSourceIbmLogsDataAccessRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsDataAccessRulesDataAccessRuleToMap(model *logsv0.DataAccessRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["display_name"] = *model.DisplayName
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDataAccessRulesDataAccessRuleFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	modelMap["default_expression"] = *model.DefaultExpression
	return modelMap, nil
}

func DataSourceIbmLogsDataAccessRulesDataAccessRuleFilterToMap(model *logsv0.DataAccessRuleFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["entity_type"] = *model.EntityType
	modelMap["expression"] = *model.Expression
	return modelMap, nil
}
