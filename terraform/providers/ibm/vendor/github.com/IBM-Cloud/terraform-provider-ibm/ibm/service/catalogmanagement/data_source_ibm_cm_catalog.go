// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmCatalog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmCatalogRead,

		Schema: map[string]*schema.Schema{
			"catalog_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Catalog identifier.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID.",
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display Name in the requested language.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"catalog_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an icon associated with this catalog.",
			},
			"catalog_banner_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for a banner image for this catalog.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tags associated with this catalog.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url for this specific catalog.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN associated with the catalog.",
			},
			"offerings_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL path to offerings.",
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of features associated with this catalog.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Heading.",
						},
						"title_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Feature description.",
						},
						"description_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes whether a catalog is disabled.",
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date-time this catalog was created.",
			},
			"updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date-time this catalog was last updated.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group id the catalog is owned by.",
			},
			"owning_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account that owns catalog.",
			},
			"catalog_filters": &schema.Schema{
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
							Description: "Filter against offering categories with dynamic keys.",
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
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kind of catalog. Supported kinds are offering and vpe.",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Catalog specific metadata.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_account_contexts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of target accounts contexts on this catalog.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API key of the target account.",
							Sensitive:   true,
						},
						"trusted_profile": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trusted profile information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Trusted profile ID.",
									},
									"catalog_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN of this catalog.",
									},
									"catalog_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of this catalog.",
									},
									"target_service_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target service ID.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier/name for this target account context.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Label for this target account context.",
						},
						"project_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCmCatalogRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listCatalogsOptions := &catalogmanagementv1.ListCatalogsOptions{}
	catalogs, response, err := catalogManagementClient.ListCatalogsWithContext(context, listCatalogsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListCatalogsWithContext failed %s\n%s", err, response), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	var catalogId string
	for _, catalog := range catalogs.Resources {
		if *catalog.ID == d.Get("catalog_identifier").(string) || *catalog.Label == d.Get("label").(string) {
			catalogId = *catalog.ID
		}
	}

	if catalogId == "" {
		tfErr := flex.TerraformErrorf(flex.FmtErrorf("Could not find catalog from provided ID or label"), "Could not find catalog from provided ID or label", "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}

	getCatalogOptions.SetCatalogIdentifier(catalogId)

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogWithContext failed %s\n%s", err, response), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getCatalogOptions.CatalogIdentifier))

	if err = d.Set("id", catalog.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting id: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("rev", catalog.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("label", catalog.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalog.LabelI18n != nil {
		if err = d.Set("label_i18n", catalog.LabelI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label_i18n: %s", err), "(Data) ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalog.ShortDescriptionI18n != nil {
		if err = d.Set("short_description_i18n", catalog.ShortDescriptionI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description_i18n: %s", err), "(Data) ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_icon_url: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_banner_url", catalog.CatalogBannerURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_banner_url: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("url", catalog.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", catalog.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offerings_url", catalog.OfferingsURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offerings_url: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	features := []map[string]interface{}{}
	if catalog.Features != nil {
		for _, modelItem := range catalog.Features {
			modelMap, err := dataSourceIBMCmCatalogFeatureToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting features: %s", err), "(Data) ibm_cm_catalog", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			features = append(features, modelMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting features: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("disabled", catalog.Disabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disabled: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created", flex.DateTimeToString(catalog.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated", flex.DateTimeToString(catalog.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group_id", catalog.ResourceGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("owning_account", catalog.OwningAccount); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting owning_account: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	catalogFilters := []map[string]interface{}{}
	if catalog.CatalogFilters != nil {
		modelMap, err := dataSourceIBMCmCatalogFiltersToMap(catalog.CatalogFilters)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_filters: %s", err), "(Data) ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		catalogFilters = append(catalogFilters, modelMap)
	}
	if err = d.Set("catalog_filters", catalogFilters); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_filters: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("kind", catalog.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kind: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalog.Metadata != nil {
		convertedMap := make(map[string]interface{}, len(catalog.Metadata))
		for k, v := range catalog.Metadata {
			convertedMap[k] = v
		}

		if err = d.Set("metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metadata: %s", err), "(Data) ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	targetAccountContexts := []map[string]interface{}{}
	if catalog.TargetAccountContexts != nil {
		for _, tacItem := range catalog.TargetAccountContexts {
			tacItemMap, err := resourceIBMCmCatalogTargetAccountContextToMap(&tacItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting target_account_contexts: %s", err), "(Data) ibm_cm_catalog", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			targetAccountContexts = append(targetAccountContexts, tacItemMap)
		}
	}
	if err = d.Set("target_account_contexts", targetAccountContexts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting target_account_contexts: %s", err), "(Data) ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIBMCmCatalogFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Title != nil {
		modelMap["title"] = *model.Title
	}
	if model.TitleI18n != nil {
		titleI18nMap := make(map[string]interface{}, len(model.TitleI18n))
		for k, v := range model.TitleI18n {
			titleI18nMap[k] = v
		}
		modelMap["title_i18n"] = flex.Flatten(titleI18nMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.DescriptionI18n != nil {
		descriptionI18nMap := make(map[string]interface{}, len(model.DescriptionI18n))
		for k, v := range model.DescriptionI18n {
			descriptionI18nMap[k] = v
		}
		modelMap["description_i18n"] = flex.Flatten(descriptionI18nMap)
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogFiltersToMap(model *catalogmanagementv1.Filters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IncludeAll != nil {
		modelMap["include_all"] = *model.IncludeAll
	}
	if model.CategoryFilters != nil {
		var categoryFiltersList []map[string]interface{}
		for k, category := range model.CategoryFilters {
			categoryFilterMap, err := dataSourceIBMCmCatalogCategoryFilterToMap(k, &category)
			if err != nil {
				return modelMap, err
			}
			categoryFiltersList = append(categoryFiltersList, categoryFilterMap)
		}
		modelMap["category_filters"] = categoryFiltersList
	}
	if model.IDFilters != nil {
		idFiltersMap, err := dataSourceIBMCmCatalogIDFilterToMap(model.IDFilters)
		if err != nil {
			return modelMap, err
		}
		modelMap["id_filters"] = []map[string]interface{}{idFiltersMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogCategoryFilterToMap(key string, model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if key != "" {
		modelMap["category_name"] = key
	}
	if model.Include != nil {
		modelMap["include"] = *model.Include
	}
	if model.Filter != nil {
		filterMap, err := dataSourceIBMCmCatalogFilterTermsToMap(model.Filter)
		if err != nil {
			return modelMap, err
		}
		modelMap["filter"] = []map[string]interface{}{filterMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogFilterTermsToMap(model *catalogmanagementv1.FilterTerms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FilterTerms != nil {
		modelMap["filter_terms"] = model.FilterTerms
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogIDFilterToMap(model *catalogmanagementv1.IDFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Include != nil {
		includeMap, err := dataSourceIBMCmCatalogFilterTermsToMap(model.Include)
		if err != nil {
			return modelMap, err
		}
		modelMap["include"] = []map[string]interface{}{includeMap}
	}
	if model.Exclude != nil {
		excludeMap, err := dataSourceIBMCmCatalogFilterTermsToMap(model.Exclude)
		if err != nil {
			return modelMap, err
		}
		modelMap["exclude"] = []map[string]interface{}{excludeMap}
	}
	return modelMap, nil
}
