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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmCatalog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmCatalogCreate,
		ReadContext:   resourceIBMCmCatalogRead,
		UpdateContext: resourceIBMCmCatalogUpdate,
		DeleteContext: resourceIBMCmCatalogDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display Name in the requested language.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"catalog_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL for an icon associated with this catalog.",
			},
			"catalog_banner_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL for a banner image for this catalog.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of tags associated with this catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of features associated with this catalog.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Heading.",
						},
						"title_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Feature description.",
						},
						"description_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether a catalog is disabled.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
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
							Optional:    true,
							Description: "-> true - Include all of the public catalog when filtering. Further settings will specifically exclude some offerings. false - Exclude all of the public catalog when filtering. Further settings will specifically include some offerings.",
						},
						"category_filters": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Filter against offering properties.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"id_filters": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of values to match against. If include is true, then if the offering has one of the values then the offering is included. If include is false, then if the offering has one of the values then the offering is excluded.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"exclude": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
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
			"syndication_settings": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Feature information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remove_related_components": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Remove related components.",
						},
						"clusters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Syndication clusters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cluster region.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cluster ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cluster name.",
									},
									"resource_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource group ID.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Syndication type.",
									},
									"namespaces": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Syndicated namespaces.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"all_namespaces": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Syndicated to all namespaces on cluster.",
									},
								},
							},
						},
						"history": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Feature information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespaces": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Array of syndicated namespaces.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"clusters": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Array of syndicated namespaces.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Cluster region.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Cluster ID.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Cluster name.",
												},
												"resource_group_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource group ID.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Syndication type.",
												},
												"namespaces": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Syndicated namespaces.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"all_namespaces": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Syndicated to all namespaces on cluster.",
												},
											},
										},
									},
									"last_run": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Date and time last syndicated.",
									},
								},
							},
						},
						"authorization": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Feature information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Array of syndicated namespaces.",
									},
									"last_run": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				Optional:    true,
				Description: "Kind of catalog. Supported kinds are offering and vpe.",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Catalog specific metadata.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
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
			"target_account_contexts": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of target accounts contexts on this catalog.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "API key of the target account.",
							Sensitive:   true,
						},
						"trusted_profile": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Trusted profile information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Target service ID.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique identifier/name for this target account context.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label for this target account context.",
						},
						"project_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Project ID.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMCmCatalogCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createCatalogOptions := &catalogmanagementv1.CreateCatalogOptions{}

	if _, ok := d.GetOk("label"); ok {
		createCatalogOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("short_description"); ok {
		createCatalogOptions.SetShortDescription(d.Get("short_description").(string))
	}
	if _, ok := d.GetOk("catalog_icon_url"); ok {
		createCatalogOptions.SetCatalogIconURL(d.Get("catalog_icon_url").(string))
	}
	if _, ok := d.GetOk("catalog_banner_url"); ok {
		createCatalogOptions.SetCatalogBannerURL(d.Get("catalog_banner_url").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createCatalogOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("features"); ok {
		var features []catalogmanagementv1.Feature
		for _, e := range d.Get("features").([]interface{}) {
			value := e.(map[string]interface{})
			featuresItem, err := resourceIBMCmCatalogMapToFeature(value)
			if err != nil {
				return diag.FromErr(err)
			}
			features = append(features, *featuresItem)
		}
		createCatalogOptions.SetFeatures(features)
	}
	if _, ok := d.GetOk("disabled"); ok {
		createCatalogOptions.SetDisabled(d.Get("disabled").(bool))
	}
	if _, ok := d.GetOk("resource_group_id"); ok {
		createCatalogOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("owning_account"); ok {
		createCatalogOptions.SetOwningAccount(d.Get("owning_account").(string))
	}
	if _, ok := d.GetOk("catalog_filters"); ok {
		catalogFiltersModel, err := resourceIBMCmCatalogMapToFilters(d.Get("catalog_filters.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createCatalogOptions.SetCatalogFilters(catalogFiltersModel)
	}
	if _, ok := d.GetOk("syndication_settings"); ok {
		syndicationSettingsModel, err := resourceIBMCmCatalogMapToSyndicationResource(d.Get("syndication_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createCatalogOptions.SetSyndicationSettings(syndicationSettingsModel)
	}
	if _, ok := d.GetOk("kind"); ok {
		createCatalogOptions.SetKind(d.Get("kind").(string))
	}
	if _, ok := d.GetOk("target_account_contexts"); ok {
		var target_account_contexts []catalogmanagementv1.TargetAccountContext
		for _, e := range d.Get("target_account_contexts").([]interface{}) {
			value := e.(map[string]interface{})
			targetAccountContextsItem, err := resourceIBMCmCatalogMapToTargetAccountContexts(nil, value)
			if err != nil {
				return diag.FromErr(err)
			}
			target_account_contexts = append(target_account_contexts, *targetAccountContextsItem)
		}
		createCatalogOptions.SetTargetAccountContexts(target_account_contexts)
	}

	catalog, response, err := catalogManagementClient.CreateCatalogWithContext(context, createCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateCatalogWithContext failed %s\n%s", err, response))
	}

	d.SetId(*catalog.ID)

	return resourceIBMCmCatalogRead(context, d, meta)
}

func resourceIBMCmCatalogRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}
	getCatalogOptions.SetCatalogIdentifier(d.Id())

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCatalogWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("rev", catalog.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rev: %s", err))
	}
	if err = d.Set("label", catalog.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}
	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting short_description: %s", err))
	}
	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_icon_url: %s", err))
	}
	if err = d.Set("catalog_banner_url", catalog.CatalogBannerURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_banner_url: %s", err))
	}
	if catalog.Tags != nil {
		if err = d.Set("tags", catalog.Tags); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tags: %s", err))
		}
	}
	features := []map[string]interface{}{}
	if catalog.Features != nil {
		for _, featuresItem := range catalog.Features {
			featuresItemMap, err := resourceIBMCmCatalogFeatureToMap(&featuresItem)
			if err != nil {
				return diag.FromErr(err)
			}
			features = append(features, featuresItemMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting features: %s", err))
	}
	if err = d.Set("disabled", catalog.Disabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting disabled: %s", err))
	}
	if err = d.Set("resource_group_id", catalog.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("owning_account", catalog.OwningAccount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting owning_account: %s", err))
	}
	if catalog.CatalogFilters != nil {
		catalogFiltersMap, err := resourceIBMCmCatalogFiltersToMap(catalog.CatalogFilters)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("catalog_filters", []map[string]interface{}{catalogFiltersMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting catalog_filters: %s", err))
		}
	}
	if catalog.SyndicationSettings != nil {
		syndicationSettingsMap, err := resourceIBMCmCatalogSyndicationResourceToMap(catalog.SyndicationSettings)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("syndication_settings", []map[string]interface{}{syndicationSettingsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting syndication_settings: %s", err))
		}
	}
	if err = d.Set("kind", catalog.Kind); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind: %s", err))
	}
	if err = d.Set("rev", catalog.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rev: %s", err))
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
	if err = d.Set("created", flex.DateTimeToString(catalog.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created: %s", err))
	}
	if err = d.Set("updated", flex.DateTimeToString(catalog.Updated)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated: %s", err))
	}
	targetAccountContexts := []map[string]interface{}{}
	if catalog.TargetAccountContexts != nil {
		for _, tacItem := range catalog.TargetAccountContexts {
			tacItemMap, err := resourceIBMCmCatalogTargetAccountContextToMap(&tacItem)
			if err != nil {
				return diag.FromErr(err)
			}
			targetAccountContexts = append(targetAccountContexts, tacItemMap)
		}
	}
	if err = d.Set("target_account_contexts", targetAccountContexts); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target_account_contexts: %s", err))
	}

	return nil
}

func resourceIBMCmCatalogUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}
	getCatalogOptions.SetCatalogIdentifier(d.Id())

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCatalogWithContext failed %s\n%s", err, response))
	}

	replaceCatalogOptions := &catalogmanagementv1.ReplaceCatalogOptions{}
	replaceCatalogOptions.SetCatalogIdentifier(*catalog.ID)
	replaceCatalogOptions.SetID(*catalog.ID)
	replaceCatalogOptions.SetRev(*catalog.Rev)

	if _, ok := d.GetOk("label"); ok {
		replaceCatalogOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("short_description"); ok {
		replaceCatalogOptions.SetShortDescription(d.Get("short_description").(string))
	}
	if _, ok := d.GetOk("catalog_icon_url"); ok {
		replaceCatalogOptions.SetCatalogIconURL(d.Get("catalog_icon_url").(string))
	}
	if _, ok := d.GetOk("catalog_banner_url"); ok {
		replaceCatalogOptions.SetCatalogBannerURL(d.Get("catalog_banner_url").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		replaceCatalogOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("features"); ok {
		var features []catalogmanagementv1.Feature
		for _, e := range d.Get("features").([]interface{}) {
			value := e.(map[string]interface{})
			featuresItem, err := resourceIBMCmCatalogMapToFeature(value)
			if err != nil {
				return diag.FromErr(err)
			}
			features = append(features, *featuresItem)
		}
		replaceCatalogOptions.SetFeatures(features)
	}
	if _, ok := d.GetOk("disabled"); ok {
		replaceCatalogOptions.SetDisabled(d.Get("disabled").(bool))
	}
	if _, ok := d.GetOk("resource_group_id"); ok {
		replaceCatalogOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("owning_account"); ok {
		replaceCatalogOptions.SetOwningAccount(d.Get("owning_account").(string))
	}
	if _, ok := d.GetOk("catalog_filters"); ok {
		catalogFilters, err := resourceIBMCmCatalogMapToFilters(d.Get("catalog_filters.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceCatalogOptions.SetCatalogFilters(catalogFilters)
	}
	if _, ok := d.GetOk("syndication_settings"); ok {
		syndicationSettings, err := resourceIBMCmCatalogMapToSyndicationResource(d.Get("syndication_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceCatalogOptions.SetSyndicationSettings(syndicationSettings)
	}
	if _, ok := d.GetOk("kind"); ok {
		replaceCatalogOptions.SetKind(d.Get("kind").(string))
	}
	if _, ok := d.GetOk("target_account_contexts"); ok {
		var target_account_contexts []catalogmanagementv1.TargetAccountContext
		for _, e := range d.Get("target_account_contexts").([]interface{}) {
			value := e.(map[string]interface{})
			targetAccountContextsItem, err := resourceIBMCmCatalogMapToTargetAccountContexts(catalog, value)
			if err != nil {
				return diag.FromErr(err)
			}
			target_account_contexts = append(target_account_contexts, *targetAccountContextsItem)
		}
		replaceCatalogOptions.SetTargetAccountContexts(target_account_contexts)
	}

	_, response, err = catalogManagementClient.ReplaceCatalogWithContext(context, replaceCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ReplaceCatalogWithContext failed %s\n%s", err, response))
	}

	return resourceIBMCmCatalogRead(context, d, meta)
}

func resourceIBMCmCatalogDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteCatalogOptions := &catalogmanagementv1.DeleteCatalogOptions{}

	deleteCatalogOptions.SetCatalogIdentifier(d.Id())

	response, err := catalogManagementClient.DeleteCatalogWithContext(context, deleteCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCatalogWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteCatalogWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCmCatalogMapToFeature(modelMap map[string]interface{}) (*catalogmanagementv1.Feature, error) {
	model := &catalogmanagementv1.Feature{}
	if modelMap["title"] != nil && modelMap["title"].(string) != "" {
		model.Title = core.StringPtr(modelMap["title"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmCatalogMapToTargetAccountContexts(catalog *catalogmanagementv1.Catalog, modelMap map[string]interface{}) (*catalogmanagementv1.TargetAccountContext, error) {
	model := &catalogmanagementv1.TargetAccountContext{}
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.APIKey = core.StringPtr(modelMap["api_key"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["project_id"] != nil && modelMap["project_id"].(string) != "" {
		model.ProjectID = core.StringPtr(modelMap["project_id"].(string))
	}
	if modelMap["trusted_profile"] != nil && len(modelMap["trusted_profile"].([]interface{})) > 0 {
		var trustedProfile *catalogmanagementv1.TrustedProfileInfo
		var err error
		if modelMap["trusted_profile"].([]interface{})[0] != nil {
			trustedProfile, err = resourceIBMCmCatalogMapToTrustedProfile(catalog, modelMap["trusted_profile"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
		}
		model.TrustedProfile = trustedProfile
	}
	return model, nil
}

func resourceIBMCmCatalogMapToTrustedProfile(catalog *catalogmanagementv1.Catalog, modelMap map[string]interface{}) (*catalogmanagementv1.TrustedProfileInfo, error) {
	model := &catalogmanagementv1.TrustedProfileInfo{}
	if catalog != nil {
		model.CatalogCRN = catalog.CRN
		model.CatalogName = catalog.Label
	} else {
		if modelMap["catalog_crn"] != nil && modelMap["catalog_crn"].(string) != "" {
			model.CatalogCRN = core.StringPtr(modelMap["catalog_crn"].(string))
		}
		if modelMap["catalog_name"] != nil && modelMap["catalog_name"].(string) != "" {
			model.CatalogName = core.StringPtr(modelMap["catalog_name"].(string))
		}
	}
	if modelMap["trusted_profile_id"] != nil && modelMap["trusted_profile_id"].(string) != "" {
		model.TrustedProfileID = core.StringPtr(modelMap["trusted_profile_id"].(string))
	}
	if modelMap["target_service_id"] != nil && modelMap["target_service_id"].(string) != "" {
		model.TargetServiceID = core.StringPtr(modelMap["target_service_id"].(string))
	}
	return model, nil
}

func resourceIBMCmCatalogMapToFilters(modelMap map[string]interface{}) (*catalogmanagementv1.Filters, error) {
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
	return model, nil
}

func resourceIBMCmCatalogMapToCategoryFilter(modelMap map[string]interface{}) (*catalogmanagementv1.CategoryFilter, error) {
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

func resourceIBMCmCatalogMapToFilterTerms(modelMap map[string]interface{}) (*catalogmanagementv1.FilterTerms, error) {
	model := &catalogmanagementv1.FilterTerms{}
	if modelMap["filter_terms"] != nil {
		filterTerms := []string{}
		for _, filterTermsItem := range modelMap["filter_terms"].([]interface{}) {
			filterTerms = append(filterTerms, filterTermsItem.(string))
		}
		model.FilterTerms = filterTerms
	}
	return model, nil
}

func resourceIBMCmCatalogMapToIDFilter(modelMap map[string]interface{}) (*catalogmanagementv1.IDFilter, error) {
	model := &catalogmanagementv1.IDFilter{}
	if modelMap["include"] != nil && len(modelMap["include"].([]interface{})) > 0 {
		IncludeModel, err := resourceIBMCmCatalogMapToFilterTerms(modelMap["include"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Include = IncludeModel
	}
	if modelMap["exclude"] != nil && len(modelMap["exclude"].([]interface{})) > 0 {
		ExcludeModel, err := resourceIBMCmCatalogMapToFilterTerms(modelMap["exclude"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Exclude = ExcludeModel
	}
	return model, nil
}

func resourceIBMCmCatalogMapToSyndicationResource(modelMap map[string]interface{}) (*catalogmanagementv1.SyndicationResource, error) {
	model := &catalogmanagementv1.SyndicationResource{}
	if modelMap["remove_related_components"] != nil {
		model.RemoveRelatedComponents = core.BoolPtr(modelMap["remove_related_components"].(bool))
	}
	if modelMap["clusters"] != nil {
		clusters := []catalogmanagementv1.SyndicationCluster{}
		for _, clustersItem := range modelMap["clusters"].([]interface{}) {
			clustersItemModel, err := resourceIBMCmCatalogMapToSyndicationCluster(clustersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			clusters = append(clusters, *clustersItemModel)
		}
		model.Clusters = clusters
	}
	if modelMap["history"] != nil && len(modelMap["history"].([]interface{})) > 0 {
		var HistoryModel *catalogmanagementv1.SyndicationHistory
		var err error
		if modelMap["history"].([]interface{})[0] != nil {
			HistoryModel, err = resourceIBMCmCatalogMapToSyndicationHistory(modelMap["history"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
		}
		model.History = HistoryModel
	}
	if modelMap["authorization"] != nil && len(modelMap["authorization"].([]interface{})) > 0 {
		AuthorizationModel, err := resourceIBMCmCatalogMapToSyndicationAuthorization(modelMap["authorization"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorization = AuthorizationModel
	}
	return model, nil
}

func resourceIBMCmCatalogMapToSyndicationCluster(modelMap map[string]interface{}) (*catalogmanagementv1.SyndicationCluster, error) {
	model := &catalogmanagementv1.SyndicationCluster{}
	if modelMap["region"] != nil && modelMap["region"].(string) != "" {
		model.Region = core.StringPtr(modelMap["region"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["resource_group_name"] != nil && modelMap["resource_group_name"].(string) != "" {
		model.ResourceGroupName = core.StringPtr(modelMap["resource_group_name"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["namespaces"] != nil {
		namespaces := []string{}
		for _, namespacesItem := range modelMap["namespaces"].([]interface{}) {
			namespaces = append(namespaces, namespacesItem.(string))
		}
		model.Namespaces = namespaces
	}
	if modelMap["all_namespaces"] != nil {
		model.AllNamespaces = core.BoolPtr(modelMap["all_namespaces"].(bool))
	}
	return model, nil
}

func resourceIBMCmCatalogMapToSyndicationHistory(modelMap map[string]interface{}) (*catalogmanagementv1.SyndicationHistory, error) {
	model := &catalogmanagementv1.SyndicationHistory{}
	if modelMap["namespaces"] != nil {
		namespaces := []string{}
		for _, namespacesItem := range modelMap["namespaces"].([]interface{}) {
			namespaces = append(namespaces, namespacesItem.(string))
		}
		model.Namespaces = namespaces
	}
	if modelMap["clusters"] != nil {
		clusters := []catalogmanagementv1.SyndicationCluster{}
		for _, clustersItem := range modelMap["clusters"].([]interface{}) {
			clustersItemModel, err := resourceIBMCmCatalogMapToSyndicationCluster(clustersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			clusters = append(clusters, *clustersItemModel)
		}
		model.Clusters = clusters
	}
	if modelMap["last_run"] != nil {

	}
	return model, nil
}

func resourceIBMCmCatalogMapToSyndicationAuthorization(modelMap map[string]interface{}) (*catalogmanagementv1.SyndicationAuthorization, error) {
	model := &catalogmanagementv1.SyndicationAuthorization{}
	if modelMap["token"] != nil && modelMap["token"].(string) != "" {
		model.Token = core.StringPtr(modelMap["token"].(string))
	}
	if modelMap["last_run"] != nil {

	}
	return model, nil
}

func resourceIBMCmCatalogFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Title != nil {
		modelMap["title"] = model.Title
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmCatalogTargetAccountContextToMap(model *catalogmanagementv1.TargetAccountContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.APIKey != nil {
		modelMap["api_key"] = model.APIKey
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Label != nil {
		modelMap["label"] = model.Label
	}
	if model.ProjectID != nil {
		modelMap["project_id"] = model.ProjectID
	}
	if model.TrustedProfile != nil {
		trustedProfileMap, err := resourceIBMCmCatalogTrustedProfileToMap(model.TrustedProfile)
		if err != nil {
			return modelMap, err
		}
		modelMap["trusted_profile"] = []map[string]interface{}{trustedProfileMap}
	}
	return modelMap, nil
}

func resourceIBMCmCatalogTrustedProfileToMap(model *catalogmanagementv1.TrustedProfileInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = model.TrustedProfileID
	}
	if model.CatalogCRN != nil {
		modelMap["catalog_crn"] = model.CatalogCRN
	}
	if model.CatalogName != nil {
		modelMap["catalog_name"] = model.CatalogName
	}
	if model.TargetServiceID != nil {
		modelMap["target_service_id"] = model.TargetServiceID
	}
	return modelMap, nil
}

func resourceIBMCmCatalogFiltersToMap(model *catalogmanagementv1.Filters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IncludeAll != nil {
		modelMap["include_all"] = model.IncludeAll
	}
	if model.IDFilters != nil {
		idFiltersMap, err := resourceIBMCmCatalogIDFilterToMap(model.IDFilters)
		if err != nil {
			return modelMap, err
		}
		modelMap["id_filters"] = []map[string]interface{}{idFiltersMap}
	}
	return modelMap, nil
}

func resourceIBMCmCatalogCategoryFilterToMap(model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Include != nil {
		modelMap["include"] = model.Include
	}
	if model.Filter != nil {
		filterMap, err := resourceIBMCmCatalogFilterTermsToMap(model.Filter)
		if err != nil {
			return modelMap, err
		}
		modelMap["filter"] = []map[string]interface{}{filterMap}
	}
	return modelMap, nil
}

func resourceIBMCmCatalogFilterTermsToMap(model *catalogmanagementv1.FilterTerms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FilterTerms != nil {
		modelMap["filter_terms"] = model.FilterTerms
	}
	return modelMap, nil
}

func resourceIBMCmCatalogIDFilterToMap(model *catalogmanagementv1.IDFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Include != nil {
		includeMap, err := resourceIBMCmCatalogFilterTermsToMap(model.Include)
		if err != nil {
			return modelMap, err
		}
		modelMap["include"] = []map[string]interface{}{includeMap}
	}
	if model.Exclude != nil {
		excludeMap, err := resourceIBMCmCatalogFilterTermsToMap(model.Exclude)
		if err != nil {
			return modelMap, err
		}
		modelMap["exclude"] = []map[string]interface{}{excludeMap}
	}
	return modelMap, nil
}

func resourceIBMCmCatalogSyndicationResourceToMap(model *catalogmanagementv1.SyndicationResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RemoveRelatedComponents != nil {
		modelMap["remove_related_components"] = model.RemoveRelatedComponents
	}
	if model.Clusters != nil {
		clusters := []map[string]interface{}{}
		for _, clustersItem := range model.Clusters {
			clustersItemMap, err := resourceIBMCmCatalogSyndicationClusterToMap(&clustersItem)
			if err != nil {
				return modelMap, err
			}
			clusters = append(clusters, clustersItemMap)
		}
		modelMap["clusters"] = clusters
	}
	if model.History != nil {
		historyMap, err := resourceIBMCmCatalogSyndicationHistoryToMap(model.History)
		if err != nil {
			return modelMap, err
		}
		modelMap["history"] = []map[string]interface{}{historyMap}
	}
	if model.Authorization != nil {
		authorizationMap, err := resourceIBMCmCatalogSyndicationAuthorizationToMap(model.Authorization)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorization"] = []map[string]interface{}{authorizationMap}
	}
	return modelMap, nil
}

func resourceIBMCmCatalogSyndicationClusterToMap(model *catalogmanagementv1.SyndicationCluster) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		modelMap["region"] = model.Region
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.ResourceGroupName != nil {
		modelMap["resource_group_name"] = model.ResourceGroupName
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Namespaces != nil {
		modelMap["namespaces"] = model.Namespaces
	}
	if model.AllNamespaces != nil {
		modelMap["all_namespaces"] = model.AllNamespaces
	}
	return modelMap, nil
}

func resourceIBMCmCatalogSyndicationHistoryToMap(model *catalogmanagementv1.SyndicationHistory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Namespaces != nil {
		modelMap["namespaces"] = model.Namespaces
	}
	if model.Clusters != nil {
		clusters := []map[string]interface{}{}
		for _, clustersItem := range model.Clusters {
			clustersItemMap, err := resourceIBMCmCatalogSyndicationClusterToMap(&clustersItem)
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

func resourceIBMCmCatalogSyndicationAuthorizationToMap(model *catalogmanagementv1.SyndicationAuthorization) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Token != nil {
		modelMap["token"] = model.Token
	}
	if model.LastRun != nil {
		modelMap["last_run"] = model.LastRun.String()
	}
	return modelMap, nil
}
