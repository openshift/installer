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
				Optional:    true,
				Description: "Filters for account and catalog filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_all": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
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
							MaxItems:    1,
							Computed:    true,
							Optional:    true,
							Description: "Filter on offering ID's. There is an include filter and an exclule filter. Both can be set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Optional:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
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
										Computed:    true,
										Description: "Offering filter terms.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_terms": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
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
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "create")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
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
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "create")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createCatalogOptions.SetCatalogFilters(catalogFiltersModel)
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
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "create")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			target_account_contexts = append(target_account_contexts, *targetAccountContextsItem)
		}
		createCatalogOptions.SetTargetAccountContexts(target_account_contexts)
	}

	catalog, response, err := catalogManagementClient.CreateCatalogWithContext(context, createCatalogOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateCatalogWithContext failed %s\n%s", err, response), "ibm_cm_catalog", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*catalog.ID)

	return resourceIBMCmCatalogRead(context, d, meta)
}

func resourceIBMCmCatalogRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}
	getCatalogOptions.SetCatalogIdentifier(d.Id())

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogWithContext failed %s\n%s", err, response), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("rev", catalog.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("label", catalog.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_icon_url: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("catalog_banner_url", catalog.CatalogBannerURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_banner_url: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if catalog.Tags != nil {
		if err = d.Set("tags", catalog.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting tags: %s", err), "ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	features := []map[string]interface{}{}
	if catalog.Features != nil {
		for _, featuresItem := range catalog.Features {
			featuresItemMap, err := resourceIBMCmCatalogFeatureToMap(&featuresItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "read")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			features = append(features, featuresItemMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting features: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("disabled", catalog.Disabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disabled: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group_id", catalog.ResourceGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("owning_account", catalog.OwningAccount); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting owning_account: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if catalog.CatalogFilters != nil {
		catalogFiltersMap, err := resourceIBMCmCatalogFiltersToMap(catalog.CatalogFilters)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("catalog_filters", []map[string]interface{}{catalogFiltersMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_filters: %s", err), "ibm_cm_catalog", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("kind", catalog.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kind: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("rev", catalog.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("url", catalog.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", catalog.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("offerings_url", catalog.OfferingsURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offerings_url: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created", flex.DateTimeToString(catalog.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated", flex.DateTimeToString(catalog.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	targetAccountContexts := []map[string]interface{}{}
	if catalog.TargetAccountContexts != nil {
		for _, tacItem := range catalog.TargetAccountContexts {
			tacItemMap, err := resourceIBMCmCatalogTargetAccountContextToMap(&tacItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "read")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			targetAccountContexts = append(targetAccountContexts, tacItemMap)
		}
	}
	if err = d.Set("target_account_contexts", targetAccountContexts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting target_account_contexts: %s", err), "ibm_cm_catalog", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIBMCmCatalogUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "update")
		log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}
	getCatalogOptions.SetCatalogIdentifier(d.Id())

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context, getCatalogOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogWithContext failed %s\n%s", err, response), "ibm_cm_catalog", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "update")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
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
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "update")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		replaceCatalogOptions.SetCatalogFilters(catalogFilters)
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
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "update")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			target_account_contexts = append(target_account_contexts, *targetAccountContextsItem)
		}
		replaceCatalogOptions.SetTargetAccountContexts(target_account_contexts)
	}

	_, response, err = catalogManagementClient.ReplaceCatalogWithContext(context, replaceCatalogOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceCatalogWithContext failed %s\n%s", err, response), "ibm_cm_catalog", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMCmCatalogRead(context, d, meta)
}

func resourceIBMCmCatalogDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_catalog", "delete")
		log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteCatalogOptions := &catalogmanagementv1.DeleteCatalogOptions{}

	deleteCatalogOptions.SetCatalogIdentifier(d.Id())

	response, err := catalogManagementClient.DeleteCatalogWithContext(context, deleteCatalogOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteCatalogWithContext failed %s\n%s", err, response), "ibm_cm_catalog", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
	if modelMap["category_filters"] != nil {
		categoryFiltersList := modelMap["category_filters"].([]interface{})
		categoryFilters := make(map[string]catalogmanagementv1.CategoryFilter) // Initialize the map for category filters

		for _, item := range categoryFiltersList {
			categoryFilterMap := item.(map[string]interface{})
			categoryName := categoryFilterMap["category_name"].(string) // Extract category_name as the map key

			// Convert the category filter to the appropriate struct
			categoryFilter, err := resourceIBMCmCatalogMapToCategoryFilter(categoryFilterMap)
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
		if modelMap["include"].([]interface{})[0] != nil && modelMap["include"].([]interface{})[0].(map[string]interface{}) != nil {
			IncludeModel, err := resourceIBMCmCatalogMapToFilterTerms(modelMap["include"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Include = IncludeModel
		}
	}
	if modelMap["exclude"] != nil && len(modelMap["exclude"].([]interface{})) > 0 {
		if modelMap["exclude"].([]interface{})[0] != nil && modelMap["exclude"].([]interface{})[0].(map[string]interface{}) != nil {
			ExcludeModel, err := resourceIBMCmCatalogMapToFilterTerms(modelMap["exclude"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Exclude = ExcludeModel
		}
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
	if model.CategoryFilters != nil {
		var categoryFiltersList []map[string]interface{}
		for k, category := range model.CategoryFilters {
			categoryFilterMap, err := resourceIBMCmCatalogCategoryFilterToMap(k, &category)
			if err != nil {
				return modelMap, err
			}
			categoryFiltersList = append(categoryFiltersList, categoryFilterMap)
		}
		modelMap["category_filters"] = categoryFiltersList
	}
	return modelMap, nil
}

func resourceIBMCmCatalogCategoryFilterToMap(key string, model *catalogmanagementv1.CategoryFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if key != "" {
		modelMap["category_name"] = key
	}
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
