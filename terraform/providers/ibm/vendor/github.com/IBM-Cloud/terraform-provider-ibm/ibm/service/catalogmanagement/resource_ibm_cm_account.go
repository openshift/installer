// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.99.1-daeb6e46-20250131-173156
 */

package catalogmanagement

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmAccountCreate,
		ReadContext:   resourceIBMCmAccountRead,
		UpdateContext: resourceIBMCmAccountUpdate,
		DeleteContext: resourceIBMCmAccountDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"hide_ibm_cloud_catalog": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Hide the public catalog in this account.",
			},
			"account_filters": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Filters for account and catalog filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_all": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "-> true - Include all of the public catalog when filtering. Further settings will specifically exclude some offerings. false - Exclude all of the public catalog when filtering. Further settings will specifically include some offerings.",
						},
						"category_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "Filter against offering categories with dynamic keys.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of this category",
									},
									"include": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										Description: "Whether to include the category in the catalog filter.",
									},
									"filter": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Optional:    true,
										Description: "Filter terms related to the category.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Optional:    true,
													Description: "List of filter terms for the category.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"id_filters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "List of values to match against. If include is true, then if the offering has one of the values then the offering is included. If include is false, then if the offering has one of the values then the offering is excluded.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"exclude": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "List of values to match against. If include is true, then if the offering has one of the values then the offering is included. If include is false, then if the offering has one of the values then the offering is excluded.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"region_filter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region filter string.",
			},
		},
	}
}

func resourceIBMCmAccountCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogAccountOptions := &catalogmanagementv1.GetCatalogAccountOptions{}

	account, _, err := catalogManagementClient.GetCatalogAccountWithContext(context, getCatalogAccountOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogAccountWithContext failed: %s", err.Error()), "ibm_cm_account", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*account.ID)

	// Call the Update function to ensure that the resource is in sync with the configuration
	return resourceIBMCmAccountUpdate(context, d, meta)
}

func resourceIBMCmAccountRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogAccountOptions := &catalogmanagementv1.GetCatalogAccountOptions{}

	account, response, err := catalogManagementClient.GetCatalogAccountWithContext(context, getCatalogAccountOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogAccountWithContext failed: %s", err.Error()), "ibm_cm_account", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(account.Rev) {
		if err = d.Set("rev", account.Rev); err != nil {
			err = fmt.Errorf("error setting rev: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "set-rev").GetDiag()
		}
	}
	if !core.IsNil(account.HideIBMCloudCatalog) {
		if err = d.Set("hide_ibm_cloud_catalog", account.HideIBMCloudCatalog); err != nil {
			err = fmt.Errorf("error setting hide_ibm_cloud_catalog: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "set-hide_ibm_cloud_catalog").GetDiag()
		}
	}
	if !core.IsNil(account.AccountFilters) {
		accountFiltersMap, err := ResourceIBMCmAccountFiltersToMap(account.AccountFilters)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "account_filters-to-map").GetDiag()
		}
		if err = d.Set("account_filters", []map[string]interface{}{accountFiltersMap}); err != nil {
			err = fmt.Errorf("error setting account_filters: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "set-account_filters").GetDiag()
		}
	}
	if !core.IsNil(account.RegionFilter) {
		if err = d.Set("region_filter", account.RegionFilter); err != nil {
			err = fmt.Errorf("error setting region_filter: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_account", "read", "set-region_filter").GetDiag()
		}
	}

	return nil
}

func resourceIBMCmAccountUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_version", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogAccountOptions := &catalogmanagementv1.GetCatalogAccountOptions{}

	account, response, err := catalogManagementClient.GetCatalogAccountWithContext(context, getCatalogAccountOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogAccountWithContext failed: %s", err.Error()), "ibm_cm_account", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateCatalogAccountOptions := &catalogmanagementv1.UpdateCatalogAccountOptions{}
	updateCatalogAccountOptions.SetID(*account.ID)
	updateCatalogAccountOptions.SetRev(*account.Rev)

	if d.HasChange("hide_IBM_cloud_catalog") {
		if v, ok := d.GetOk("hide_IBM_cloud_catalog"); ok {
			updateCatalogAccountOptions.SetHideIBMCloudCatalog(v.(bool))
		}
	} else if account.HideIBMCloudCatalog != nil {
		updateCatalogAccountOptions.SetHideIBMCloudCatalog(*account.HideIBMCloudCatalog)
	}

	if d.HasChange("region_filter") {
		updateCatalogAccountOptions.SetRegionFilter(d.Get("region_filter").(string))
	} else if account.RegionFilter != nil {
		updateCatalogAccountOptions.SetRegionFilter(*account.RegionFilter)
	}

	if d.HasChange("account_filters") {
		if v, ok := d.GetOk("account_filters"); ok {
			accountFilters, err := resourceIBMCmAccountFiltersMapToFilters(v.([]interface{})[0].(map[string]interface{}))
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_account", "update")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			updateCatalogAccountOptions.SetAccountFilters(accountFilters)
		}
	} else if account.AccountFilters != nil {
		updateCatalogAccountOptions.SetAccountFilters(account.AccountFilters)
	}

	_, response, err = catalogManagementClient.UpdateCatalogAccountWithContext(context, updateCatalogAccountOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateCatalogAccountWithContext failed %s\n%s", err, response), "ibm_cm_object", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMCmAccountRead(context, d, meta)
}

func resourceIBMCmAccountDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
	d.SetId("")
	return nil
}

func ResourceIBMCmAccountFiltersToMap(model *catalogmanagementv1.Filters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IncludeAll != nil {
		modelMap["include_all"] = *model.IncludeAll
	}
	if model.CategoryFilters != nil {
		var categoryFiltersList []map[string]interface{}
		for k, category := range model.CategoryFilters {
			categoryFilterMap, err := ResourceIBMCmAccountCategoryFilterToMap(k, &category)
			if err != nil {
				return modelMap, err
			}
			categoryFiltersList = append(categoryFiltersList, categoryFilterMap)
		}
		modelMap["category_filters"] = categoryFiltersList
	}
	if model.IDFilters != nil {
		idFiltersMap, err := ResourceIBMCmAccountIDFilterToMap(model.IDFilters)
		if err != nil {
			return modelMap, err
		}
		modelMap["id_filters"] = []map[string]interface{}{idFiltersMap}
	}
	return modelMap, nil
}

func ResourceIBMCmAccountCategoryFilterToMap(key string, model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if key != "" {
		modelMap["category_name"] = key
	}
	if model.Include != nil {
		modelMap["include"] = *model.Include
	}
	if model.Filter != nil {
		filterMap, err := ResourceIBMCmAccountFilterTermsToMap(model.Filter)
		if err != nil {
			return modelMap, err
		}
		modelMap["filter"] = []map[string]interface{}{filterMap}
	}
	return modelMap, nil
}

func ResourceIBMCmAccountFilterTermsToMap(model *catalogmanagementv1.FilterTerms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FilterTerms != nil {
		modelMap["filter_terms"] = model.FilterTerms
	}
	return modelMap, nil
}

func ResourceIBMCmAccountIDFilterToMap(model *catalogmanagementv1.IDFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Include != nil {
		includeMap, err := ResourceIBMCmAccountFilterTermsToMap(model.Include)
		if err != nil {
			return modelMap, err
		}
		modelMap["include"] = []map[string]interface{}{includeMap}
	}
	if model.Exclude != nil {
		excludeMap, err := ResourceIBMCmAccountFilterTermsToMap(model.Exclude)
		if err != nil {
			return modelMap, err
		}
		modelMap["exclude"] = []map[string]interface{}{excludeMap}
	}
	return modelMap, nil
}

func resourceIBMCmAccountFiltersMapToFilters(modelMap map[string]interface{}) (*catalogmanagementv1.Filters, error) {
	model := &catalogmanagementv1.Filters{}
	if modelMap["include_all"] != nil {
		model.IncludeAll = core.BoolPtr(modelMap["include_all"].(bool))
	}
	if modelMap["id_filters"] != nil && len(modelMap["id_filters"].([]interface{})) > 0 {
		var IDFiltersModel *catalogmanagementv1.IDFilter
		var err error
		if modelMap["id_filters"].([]interface{})[0] != nil {
			IDFiltersModel, err = resourceIBMCmCatalogMapToIDFilter(modelMap["id_filters"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
		}
		model.IDFilters = IDFiltersModel
	}
	if modelMap["category_filters"] != nil {
		categoryFiltersList := modelMap["category_filters"].([]interface{})
		categoryFilters := make(map[string]catalogmanagementv1.CategoryFilter) // Initialize the map for category filters

		for _, item := range categoryFiltersList {
			categoryFilterMap := item.(map[string]interface{})
			categoryName := categoryFilterMap["category_name"].(string) // Extract category_name as the map key

			// Convert the category filter to the appropriate struct
			categoryFilter, err := resourceIBMCmAccountMapToCategoryFilter(categoryFilterMap)
			if err != nil {
				return model, err
			}

			// Add the category filter to the map using category_name as the key
			categoryFilters[categoryName] = *categoryFilter
		}

		// Assign the map to the model
		model.CategoryFilters = categoryFilters
	}
	return model, nil
}

func resourceIBMCmAccountMapToCategoryFilter(modelMap map[string]interface{}) (*catalogmanagementv1.CategoryFilter, error) {
	model := &catalogmanagementv1.CategoryFilter{}
	if modelMap["include"] != nil {
		model.Include = core.BoolPtr(modelMap["include"].(bool))
	}
	if modelMap["filter"] != nil && len(modelMap["filter"].([]interface{})) > 0 {
		FilterModel, err := resourceIBMCmCatalogMapToFilterTerms(modelMap["filter"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Filter = FilterModel
	}
	return model, nil
}
