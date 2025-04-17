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
				Required:    true,
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
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Filter against offering properties.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
			"syndication_settings": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Feature information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remove_related_components": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Remove related components.",
						},
						"clusters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Syndication clusters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster region.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name.",
									},
									"resource_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource group ID.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Syndication type.",
									},
									"namespaces": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Syndicated namespaces.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"all_namespaces": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Syndicated to all namespaces on cluster.",
									},
								},
							},
						},
						"history": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Feature information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespaces": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Array of syndicated namespaces.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"clusters": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Array of syndicated namespaces.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster region.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster ID.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster name.",
												},
												"resource_group_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource group ID.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Syndication type.",
												},
												"namespaces": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Syndicated namespaces.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"all_namespaces": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Syndicated to all namespaces on cluster.",
												},
											},
										},
									},
									"last_run": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Date and time last syndicated.",
									},
								},
							},
						},
						"authorization": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Feature information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Array of syndicated namespaces.",
									},
									"last_run": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Date and time last updated.",
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
		},
	}
}

func dataSourceIBMCmCatalogRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}

	getCatalogOptions.SetCatalogIdentifier(d.Get("catalog_identifier").(string))

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCatalogWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getCatalogOptions.CatalogIdentifier))

	if err = d.Set("id", catalog.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if err = d.Set("rev", catalog.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rev: %s", err))
	}

	if err = d.Set("label", catalog.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}

	if catalog.LabelI18n != nil {
		if err = d.Set("label_i18n", catalog.LabelI18n); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting label_i18n: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting label_i18n %s", err))
		}
	}

	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting short_description: %s", err))
	}

	if catalog.ShortDescriptionI18n != nil {
		if err = d.Set("short_description_i18n", catalog.ShortDescriptionI18n); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting short_description_i18n: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting short_description_i18n %s", err))
		}
	}

	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_icon_url: %s", err))
	}

	if err = d.Set("catalog_banner_url", catalog.CatalogBannerURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_banner_url: %s", err))
	}

	if err = d.Set("url", catalog.URL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting url: %s", err))
	}

	if err = d.Set("crn", catalog.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("offerings_url", catalog.OfferingsURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting offerings_url: %s", err))
	}

	features := []map[string]interface{}{}
	if catalog.Features != nil {
		for _, modelItem := range catalog.Features {
			modelMap, err := dataSourceIBMCmCatalogFeatureToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			features = append(features, modelMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting features %s", err))
	}

	if err = d.Set("disabled", catalog.Disabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting disabled: %s", err))
	}

	if err = d.Set("created", flex.DateTimeToString(catalog.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created: %s", err))
	}

	if err = d.Set("updated", flex.DateTimeToString(catalog.Updated)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated: %s", err))
	}

	if err = d.Set("resource_group_id", catalog.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}

	if err = d.Set("owning_account", catalog.OwningAccount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting owning_account: %s", err))
	}

	catalogFilters := []map[string]interface{}{}
	if catalog.CatalogFilters != nil {
		modelMap, err := dataSourceIBMCmCatalogFiltersToMap(catalog.CatalogFilters)
		if err != nil {
			return diag.FromErr(err)
		}
		catalogFilters = append(catalogFilters, modelMap)
	}
	if err = d.Set("catalog_filters", catalogFilters); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_filters %s", err))
	}

	syndicationSettings := []map[string]interface{}{}
	if catalog.SyndicationSettings != nil {
		modelMap, err := dataSourceIBMCmCatalogSyndicationResourceToMap(catalog.SyndicationSettings)
		if err != nil {
			return diag.FromErr(err)
		}
		syndicationSettings = append(syndicationSettings, modelMap)
	}
	if err = d.Set("syndication_settings", syndicationSettings); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting syndication_settings %s", err))
	}

	if err = d.Set("kind", catalog.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind: %s", err))
	}

	if catalog.Metadata != nil {
		convertedMap := make(map[string]interface{}, len(catalog.Metadata))
		for k, v := range catalog.Metadata {
			convertedMap[k] = v
		}

		if err = d.Set("metadata", flex.Flatten(convertedMap)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting metadata: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
		}
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
		categoryFiltersMap := make(map[string]interface{}, len(model.CategoryFilters))
		modelMap["category_filters"] = flex.Flatten(categoryFiltersMap)
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

func dataSourceIBMCmCatalogCategoryFilterToMap(model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
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

func dataSourceIBMCmCatalogSyndicationResourceToMap(model *catalogmanagementv1.SyndicationResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RemoveRelatedComponents != nil {
		modelMap["remove_related_components"] = *model.RemoveRelatedComponents
	}
	if model.Clusters != nil {
		clusters := []map[string]interface{}{}
		for _, clustersItem := range model.Clusters {
			clustersItemMap, err := dataSourceIBMCmCatalogSyndicationClusterToMap(&clustersItem)
			if err != nil {
				return modelMap, err
			}
			clusters = append(clusters, clustersItemMap)
		}
		modelMap["clusters"] = clusters
	}
	if model.History != nil {
		historyMap, err := dataSourceIBMCmCatalogSyndicationHistoryToMap(model.History)
		if err != nil {
			return modelMap, err
		}
		modelMap["history"] = []map[string]interface{}{historyMap}
	}
	if model.Authorization != nil {
		authorizationMap, err := dataSourceIBMCmCatalogSyndicationAuthorizationToMap(model.Authorization)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorization"] = []map[string]interface{}{authorizationMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogSyndicationClusterToMap(model *catalogmanagementv1.SyndicationCluster) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		modelMap["region"] = *model.Region
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroupName != nil {
		modelMap["resource_group_name"] = *model.ResourceGroupName
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Namespaces != nil {
		modelMap["namespaces"] = model.Namespaces
	}
	if model.AllNamespaces != nil {
		modelMap["all_namespaces"] = *model.AllNamespaces
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogSyndicationHistoryToMap(model *catalogmanagementv1.SyndicationHistory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Namespaces != nil {
		modelMap["namespaces"] = model.Namespaces
	}
	if model.Clusters != nil {
		clusters := []map[string]interface{}{}
		for _, clustersItem := range model.Clusters {
			clustersItemMap, err := dataSourceIBMCmCatalogSyndicationClusterToMap(&clustersItem)
			if err != nil {
				return modelMap, err
			}
			clusters = append(clusters, clustersItemMap)
		}
		modelMap["clusters"] = clusters
	}
	if model.LastRun != nil {
		modelMap["last_run"] = model.LastRun.String()
	}
	return modelMap, nil
}

func dataSourceIBMCmCatalogSyndicationAuthorizationToMap(model *catalogmanagementv1.SyndicationAuthorization) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Token != nil {
		modelMap["token"] = *model.Token
	}
	if model.LastRun != nil {
		modelMap["last_run"] = model.LastRun.String()
	}
	return modelMap, nil
}
