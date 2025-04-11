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

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsDataAccessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsDataAccessRuleCreate,
		ReadContext:   resourceIbmLogsDataAccessRuleRead,
		UpdateContext: resourceIbmLogsDataAccessRuleUpdate,
		DeleteContext: resourceIbmLogsDataAccessRuleDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"display_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_data_access_rule", "display_name"),
				Description:  "Data Access Rule Display Name.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_data_access_rule", "description"),
				Description:  "Optional Data Access Rule Description.",
			},
			"filters": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of filters that the Data Access Rule is composed of.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter's Entity Type.",
						},
						"expression": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter's Expression.",
						},
					},
				},
			},
			"default_expression": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_data_access_rule", "default_expression"),
				Description:  "Default expression to use when no filter matches the query.",
			},
			"access_rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data access rule ID.",
			},
		},
	}
}

func ResourceIbmLogsDataAccessRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "display_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9_\-\s]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "default_expression",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Za-z0-9_\.,\-"{}()\[\]=!:#\/$|'<> ]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_data_access_rule", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsDataAccessRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createDataAccessRuleOptions := &logsv0.CreateDataAccessRuleOptions{}

	createDataAccessRuleOptions.SetDisplayName(d.Get("display_name").(string))
	var filters []logsv0.DataAccessRuleFilter
	for _, v := range d.Get("filters").([]interface{}) {
		value := v.(map[string]interface{})
		filtersItem, err := ResourceIbmLogsDataAccessRuleMapToDataAccessRuleFilter(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "create", "parse-filters").GetDiag()
		}
		filters = append(filters, *filtersItem)
	}
	createDataAccessRuleOptions.SetFilters(filters)
	createDataAccessRuleOptions.SetDefaultExpression(d.Get("default_expression").(string))
	if _, ok := d.GetOk("description"); ok {
		createDataAccessRuleOptions.SetDescription(d.Get("description").(string))
	}

	dataAccessRule, _, err := logsClient.CreateDataAccessRuleWithContext(context, createDataAccessRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDataAccessRuleWithContext failed: %s", err.Error()), "ibm_logs_data_access_rule", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	dataAccessRuleId := fmt.Sprintf("%s/%s/%s", region, instanceId, *dataAccessRule.ID)
	d.SetId(dataAccessRuleId)

	return resourceIbmLogsDataAccessRuleRead(context, d, meta)
}

func resourceIbmLogsDataAccessRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, accessRuleId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	listDataAccessRulesOptions := &logsv0.ListDataAccessRulesOptions{}

	ruleID := core.UUIDPtr(strfmt.UUID(accessRuleId))
	listDataAccessRulesOptions.ID = []strfmt.UUID{*ruleID}

	dataAccessRuleCollection, response, err := logsClient.ListDataAccessRulesWithContext(context, listDataAccessRulesOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListDataAccessRulesWithContext failed: %s", err.Error()), "ibm_logs_data_access_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("access_rule_id", accessRuleId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting access_rule_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if dataAccessRuleCollection != nil && len(dataAccessRuleCollection.DataAccessRules) > 0 && &dataAccessRuleCollection.DataAccessRules[0] != nil {

		dataAccessRule := dataAccessRuleCollection.DataAccessRules[0]

		if err = d.Set("display_name", dataAccessRule.DisplayName); err != nil {
			err = fmt.Errorf("Error setting display_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read", "set-display_name").GetDiag()
		}
		if !core.IsNil(dataAccessRule.Description) {
			if err = d.Set("description", dataAccessRule.Description); err != nil {
				err = fmt.Errorf("Error setting description: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read", "set-description").GetDiag()
			}
		}
		filters := []map[string]interface{}{}
		for _, filtersItem := range dataAccessRule.Filters {
			filtersItemMap, err := ResourceIbmLogsDataAccessRuleDataAccessRuleFilterToMap(&filtersItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read", "filters-to-map").GetDiag()
			}
			filters = append(filters, filtersItemMap)
		}
		if err = d.Set("filters", filters); err != nil {
			err = fmt.Errorf("Error setting filters: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read", "set-filters").GetDiag()
		}
		if err = d.Set("default_expression", dataAccessRule.DefaultExpression); err != nil {
			err = fmt.Errorf("Error setting default_expression: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "read", "set-default_expression").GetDiag()
		}
	}

	return nil
}

func resourceIbmLogsDataAccessRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, accessRuleId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateDataAccessRuleOptions := &logsv0.UpdateDataAccessRuleOptions{}

	updateDataAccessRuleOptions.SetID(core.UUIDPtr(strfmt.UUID(accessRuleId)))

	hasChange := false

	if d.HasChange("display_name") ||
		d.HasChange("filters") ||
		d.HasChange("default_expression") ||
		d.HasChange("description") {

		updateDataAccessRuleOptions.SetDisplayName(d.Get("display_name").(string))
		var filters []logsv0.DataAccessRuleFilter
		for _, v := range d.Get("filters").([]interface{}) {
			value := v.(map[string]interface{})
			filtersItem, err := ResourceIbmLogsDataAccessRuleMapToDataAccessRuleFilter(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "update", "parse-filters").GetDiag()
			}
			filters = append(filters, *filtersItem)
		}
		updateDataAccessRuleOptions.SetFilters(filters)
		updateDataAccessRuleOptions.SetDefaultExpression(d.Get("default_expression").(string))

		updateDataAccessRuleOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdateDataAccessRuleWithContext(context, updateDataAccessRuleOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateDataAccessRuleWithContext failed: %s", err.Error()), "ibm_logs_data_access_rule", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsDataAccessRuleRead(context, d, meta)
}

func resourceIbmLogsDataAccessRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_access_rule", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, accessRuleId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDataAccessRuleOptions := &logsv0.DeleteDataAccessRuleOptions{}
	deleteDataAccessRuleOptions.SetID(core.UUIDPtr(strfmt.UUID(accessRuleId)))

	_, err = logsClient.DeleteDataAccessRuleWithContext(context, deleteDataAccessRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDataAccessRuleWithContext failed: %s", err.Error()), "ibm_logs_data_access_rule", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsDataAccessRuleMapToDataAccessRuleFilter(modelMap map[string]interface{}) (*logsv0.DataAccessRuleFilter, error) {
	model := &logsv0.DataAccessRuleFilter{}
	model.EntityType = core.StringPtr(modelMap["entity_type"].(string))
	model.Expression = core.StringPtr(modelMap["expression"].(string))
	return model, nil
}

func ResourceIbmLogsDataAccessRuleDataAccessRuleFilterToMap(model *logsv0.DataAccessRuleFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["entity_type"] = *model.EntityType
	modelMap["expression"] = *model.Expression
	return modelMap, nil
}
