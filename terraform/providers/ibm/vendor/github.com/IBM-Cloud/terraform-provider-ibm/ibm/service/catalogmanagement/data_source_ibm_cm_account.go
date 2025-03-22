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

func DataSourceIBMCmAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmAccountRead,

		Schema: map[string]*schema.Schema{
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"hide_ibm_cloud_catalog": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Hide the public catalog in this account.",
			},
			"account_filters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Filters for account and catalog filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_all": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "-> true - Include all of the public catalog when filtering. Further settings will specifically exclude some offerings. false - Exclude all of the public catalog when filtering. Further settings will specifically include some offerings.",
						},
						"category_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Filter against offering categories.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of this category",
									},
									"include": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to include the category in the catalog filter.",
									},
									"filter": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Filter terms related to the category.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
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
							Computed:    true,
							Description: "Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of values to match against. If include is true, then if the offering has one of the values then the offering is included. If include is false, then if the offering has one of the values then the offering is excluded.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"exclude": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of values to match against. If include is true, then if the offering has one of the values then the offering is included. If include is false, then if the offering has one of the values then the offering is excluded.",
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
					},
				},
			},
			"region_filter": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region filter string.",
			},
		},
	}
}

func dataSourceIBMCmAccountRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cm_account", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogAccountOptions := &catalogmanagementv1.GetCatalogAccountOptions{}

	account, _, err := catalogManagementClient.GetCatalogAccountWithContext(context, getCatalogAccountOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogAccountWithContext failed: %s", err.Error()), "(Data) ibm_cm_account", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*account.ID)

	if !core.IsNil(account.Rev) {
		if err = d.Set("rev", account.Rev); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cm_account", "read", "set-rev").GetDiag()
		}
	}

	if !core.IsNil(account.HideIBMCloudCatalog) {
		if err = d.Set("hide_ibm_cloud_catalog", account.HideIBMCloudCatalog); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting hide_ibm_cloud_catalog: %s", err), "(Data) ibm_cm_account", "read", "set-hide_ibm_cloud_catalog").GetDiag()
		}
	}

	if !core.IsNil(account.AccountFilters) {
		accountFilters := []map[string]interface{}{}
		accountFiltersMap, err := DataSourceIBMCmAccountFiltersToMap(account.AccountFilters)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cm_account", "read", "account_filters-to-map").GetDiag()
		}
		accountFilters = append(accountFilters, accountFiltersMap)
		if err = d.Set("account_filters", accountFilters); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting account_filters: %s", err), "(Data) ibm_cm_account", "read", "set-account_filters").GetDiag()
		}
	}

	if !core.IsNil(account.RegionFilter) {
		if err = d.Set("region_filter", account.RegionFilter); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region_filter: %s", err), "(Data) ibm_cm_account", "read", "set-region_filter").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMCmAccountFiltersToMap(model *catalogmanagementv1.Filters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IncludeAll != nil {
		modelMap["include_all"] = *model.IncludeAll
	}
	if model.CategoryFilters != nil {
		var categoryFiltersList []map[string]interface{}
		for k, category := range model.CategoryFilters {
			categoryFilterMap, err := DataSourceIBMCmAccountCategoryFilterToMap(k, &category)
			if err != nil {
				return modelMap, err
			}
			categoryFiltersList = append(categoryFiltersList, categoryFilterMap)
		}
		modelMap["category_filters"] = categoryFiltersList
	}
	if model.IDFilters != nil {
		idFiltersMap, err := DataSourceIBMCmAccountIDFilterToMap(model.IDFilters)
		if err != nil {
			return modelMap, err
		}
		modelMap["id_filters"] = []map[string]interface{}{idFiltersMap}
	}
	return modelMap, nil
}

func DataSourceIBMCmAccountCategoryFilterToMap(key string, model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if key != "" {
		modelMap["category_name"] = key
	}
	if model.Include != nil {
		modelMap["include"] = *model.Include
	}
	if model.Filter != nil {
		filterMap, err := DataSourceIBMCmAccountFilterTermsToMap(model.Filter)
		if err != nil {
			return modelMap, err
		}
		modelMap["filter"] = []map[string]interface{}{filterMap}
	}
	return modelMap, nil
}

func DataSourceIBMCmAccountFilterTermsToMap(model *catalogmanagementv1.FilterTerms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FilterTerms != nil {
		modelMap["filter_terms"] = model.FilterTerms
	}
	return modelMap, nil
}

func DataSourceIBMCmAccountIDFilterToMap(model *catalogmanagementv1.IDFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Include != nil {
		includeMap, err := DataSourceIBMCmAccountFilterTermsToMap(model.Include)
		if err != nil {
			return modelMap, err
		}
		modelMap["include"] = []map[string]interface{}{includeMap}
	}
	if model.Exclude != nil {
		excludeMap, err := DataSourceIBMCmAccountFilterTermsToMap(model.Exclude)
		if err != nil {
			return modelMap, err
		}
		modelMap["exclude"] = []map[string]interface{}{excludeMap}
	}
	return modelMap, nil
}
