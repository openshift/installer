// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIBMSchematicsResourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSchematicsResourceQueryCreate,
		ReadContext:   resourceIBMSchematicsResourceQueryRead,
		UpdateContext: resourceIBMSchematicsResourceQueryUpdate,
		DeleteContext: resourceIBMSchematicsResourceQueryDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_resource_query", "type"),
				Description:  "Resource type (cluster, vsi, icd, vpc).",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource query name.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"queries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of the query(workspaces).",
						},
						"query_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the resource query param.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value of the resource query param.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of resource query param variable.",
									},
								},
							},
						},
						"query_select": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of query selection parameters.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource query creation time.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who created the Resource query.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource query updation time.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who updated the Resource query.",
			},
		},
	}
}

func ResourceIBMSchematicsResourceQueryValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "vsi",
		})

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_schematics_resource_query", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSchematicsResourceQueryCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_resource_query", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}
	createResourceQueryOptions := &schematicsv1.CreateResourceQueryOptions{}

	if _, ok := d.GetOk("type"); ok {
		createResourceQueryOptions.SetType(d.Get("type").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createResourceQueryOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("queries"); ok {
		var queries []schematicsv1.ResourceQuery
		for _, e := range d.Get("queries").([]interface{}) {
			value := e.(map[string]interface{})
			queriesItem := resourceIBMSchematicsResourceQueryMapToResourceQuery(value)
			queries = append(queries, queriesItem)
		}
		createResourceQueryOptions.SetQueries(queries)
	}

	resourceQueryRecord, response, err := schematicsClient.CreateResourceQueryWithContext(context, createResourceQueryOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryCreate CreateResourceQueryWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_resource_query", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*resourceQueryRecord.ID)

	return resourceIBMSchematicsResourceQueryRead(context, d, meta)
}

func resourceIBMSchematicsResourceQueryMapToResourceQuery(resourceQueryMap map[string]interface{}) schematicsv1.ResourceQuery {
	resourceQuery := schematicsv1.ResourceQuery{}

	if resourceQueryMap["query_type"] != nil {
		resourceQuery.QueryType = core.StringPtr(resourceQueryMap["query_type"].(string))
	}
	if resourceQueryMap["query_condition"] != nil {
		queryCondition := []schematicsv1.ResourceQueryParam{}
		for _, queryConditionItem := range resourceQueryMap["query_condition"].([]interface{}) {
			queryConditionItemModel := resourceIBMSchematicsResourceQueryMapToResourceQueryParam(queryConditionItem.(map[string]interface{}))
			queryCondition = append(queryCondition, queryConditionItemModel)
		}
		resourceQuery.QueryCondition = queryCondition
	}
	if resourceQueryMap["query_select"] != nil {
		querySelect := []string{}
		for _, querySelectItem := range resourceQueryMap["query_select"].([]interface{}) {
			querySelect = append(querySelect, querySelectItem.(string))
		}
		resourceQuery.QuerySelect = querySelect
	}

	return resourceQuery
}

func resourceIBMSchematicsResourceQueryMapToResourceQueryParam(resourceQueryParamMap map[string]interface{}) schematicsv1.ResourceQueryParam {
	resourceQueryParam := schematicsv1.ResourceQueryParam{}

	if resourceQueryParamMap["name"] != nil {
		resourceQueryParam.Name = core.StringPtr(resourceQueryParamMap["name"].(string))
	}
	if resourceQueryParamMap["value"] != nil {
		resourceQueryParam.Value = core.StringPtr(resourceQueryParamMap["value"].(string))
	}
	if resourceQueryParamMap["description"] != nil {
		resourceQueryParam.Description = core.StringPtr(resourceQueryParamMap["description"].(string))
	}

	return resourceQueryParam
}

func resourceIBMSchematicsResourceQueryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	getResourcesQueryOptions := &schematicsv1.GetResourcesQueryOptions{}

	getResourcesQueryOptions.SetQueryID(d.Id())

	resourceQueryRecord, response, err := schematicsClient.GetResourcesQueryWithContext(context, getResourcesQueryOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead GetResourcesQueryWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("type", resourceQueryRecord.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", resourceQueryRecord.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if resourceQueryRecord.Queries != nil {
		queries := []map[string]interface{}{}
		for _, queriesItem := range resourceQueryRecord.Queries {
			queriesItemMap := resourceIBMSchematicsResourceQueryResourceQueryToMap(queriesItem)
			queries = append(queries, queriesItemMap)
		}
		if err = d.Set("queries", queries); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(resourceQueryRecord.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", resourceQueryRecord.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(resourceQueryRecord.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_by", resourceQueryRecord.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryRead failed with error: %s", err), "ibm_schematics_resource_query", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIBMSchematicsResourceQueryResourceQueryToMap(resourceQuery schematicsv1.ResourceQuery) map[string]interface{} {
	resourceQueryMap := map[string]interface{}{}

	if resourceQuery.QueryType != nil {
		resourceQueryMap["query_type"] = resourceQuery.QueryType
	}
	if resourceQuery.QueryCondition != nil {
		queryCondition := []map[string]interface{}{}
		for _, queryConditionItem := range resourceQuery.QueryCondition {
			queryConditionItemMap := resourceIBMSchematicsResourceQueryResourceQueryParamToMap(queryConditionItem)
			queryCondition = append(queryCondition, queryConditionItemMap)
			// TODO: handle QueryCondition of type TypeList -- list of non-primitive, not model items
		}
		resourceQueryMap["query_condition"] = queryCondition
	}
	if resourceQuery.QuerySelect != nil {
		resourceQueryMap["query_select"] = resourceQuery.QuerySelect
	}

	return resourceQueryMap
}

func resourceIBMSchematicsResourceQueryResourceQueryParamToMap(resourceQueryParam schematicsv1.ResourceQueryParam) map[string]interface{} {
	resourceQueryParamMap := map[string]interface{}{}

	if resourceQueryParam.Name != nil {
		resourceQueryParamMap["name"] = resourceQueryParam.Name
	}
	if resourceQueryParam.Value != nil {
		resourceQueryParamMap["value"] = resourceQueryParam.Value
	}
	if resourceQueryParam.Description != nil {
		resourceQueryParamMap["description"] = resourceQueryParam.Description
	}

	return resourceQueryParamMap
}

func resourceIBMSchematicsResourceQueryUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_resource_query", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	replaceResourcesQueryOptions := &schematicsv1.ReplaceResourcesQueryOptions{}

	replaceResourcesQueryOptions.SetQueryID(d.Id())
	if _, ok := d.GetOk("type"); ok {
		replaceResourcesQueryOptions.SetType(d.Get("type").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		replaceResourcesQueryOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("queries"); ok {
		var queries []schematicsv1.ResourceQuery
		for _, e := range d.Get("queries").([]interface{}) {
			value := e.(map[string]interface{})
			queriesItem := resourceIBMSchematicsResourceQueryMapToResourceQuery(value)
			queries = append(queries, queriesItem)
		}
		replaceResourcesQueryOptions.SetQueries(queries)
	}

	_, response, err := schematicsClient.ReplaceResourcesQueryWithContext(context, replaceResourcesQueryOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryUpdate ReplaceResourcesQueryWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_resource_query", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMSchematicsResourceQueryRead(context, d, meta)
}

func resourceIBMSchematicsResourceQueryDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryDelete schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_resource_query", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	deleteResourcesQueryOptions := &schematicsv1.DeleteResourcesQueryOptions{}

	deleteResourcesQueryOptions.SetQueryID(d.Id())

	response, err := schematicsClient.DeleteResourcesQueryWithContext(context, deleteResourcesQueryOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsResourceQueryDelete DeleteResourcesQueryWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_resource_query", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
