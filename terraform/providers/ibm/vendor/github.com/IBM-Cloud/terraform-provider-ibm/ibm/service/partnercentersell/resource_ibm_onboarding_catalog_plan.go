// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.99.1-daeb6e46-20250131-173156
 */

package partnercentersell

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/partnercentersellv1"
)

func ResourceIbmOnboardingCatalogPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingCatalogPlanCreate,
		ReadContext:   resourceIbmOnboardingCatalogPlanRead,
		UpdateContext: resourceIbmOnboardingCatalogPlanUpdate,
		DeleteContext: resourceIbmOnboardingCatalogPlanDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_plan", "product_id"),
				Description:  "The unique ID of the product.",
			},
			"catalog_product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_plan", "catalog_product_id"),
				Description:  "The unique ID of this global catalog product.",
			},
			"env": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_plan", "env"),
				Description:  "The environment to fetch this object from.",
			},
			"object_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The desired ID of the global catalog object.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_plan", "name"),
				Description:  "The programmatic name of this plan.",
			},
			"active": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the service is active.",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines the global visibility for the catalog entry, and its children. If it is not enabled, all plans are disabled.",
			},
			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_plan", "kind"),
				Description:  "The kind of the global catalog object.",
			},
			"overview_ui": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The object that contains the service details from the Overview page in global catalog.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"en": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Translated details about the service, for example, display name, short description, and long description.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"display_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The display name of the product.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The short description of the product that is displayed in your catalog entry.",
									},
									"long_description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The detailed description of your product that is displayed at the beginning of your product page in the catalog. Markdown markup language is supported.",
									},
								},
							},
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of tags that carry information about your product. These tags can be used to find your product in the IBM Cloud catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pricing_tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of tags that carry information about the pricing information of your product.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"object_provider": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The provider or owner of the product.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the provider.",
						},
						"email": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email address of the provider.",
						},
					},
				},
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Global catalog plan metadata.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rc_compatible": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the object is compatible with the resource controller service.",
						},
						"ui": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The UI metadata of this service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"strings": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The data strings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"en": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Translated content of additional information about the service.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bullets": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The list of features that highlights your product's attributes and benefits for users.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The description about the features of the product.",
																		},
																		"title": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The descriptive title for the feature.",
																		},
																	},
																},
															},
															"media": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The list of supporting media for this product.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"caption": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Provide a descriptive caption that indicates what the media illustrates. This caption is displayed in the catalog.",
																		},
																		"thumbnail": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The reduced-size version of your images and videos.",
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The type of the media.",
																		},
																		"url": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL that links to the media that shows off the product.",
																		},
																	},
																},
															},
															"navigation_items": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of custom navigation panel.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Id of custom navigation panel.",
																		},
																		"url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Url for custom navigation panel.",
																		},
																		"label": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Url for custom navigation panel.",
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
									"urls": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Metadata with URLs related to a service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"doc_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The URL for your product's documentation.",
												},
												"apidocs_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The URL for your product's API documentation.",
												},
												"terms_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The URL for your product's end user license agreement.",
												},
												"instructions_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Controls the Getting Started tab on the Resource Details page. Setting it the content is loaded from the specified URL.",
												},
												"catalog_details_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Controls the Provisioning page URL, if set the assumption is that this URL is the provisioning URL for your service.",
												},
												"custom_create_page_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Controls the Provisioning page URL, if set the assumption is that this URL is the provisioning URL for your service.",
												},
												"dashboard": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Controls if your service has a custom dashboard or Resource Detail page.",
												},
											},
										},
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the object is hidden from the consumption catalog.",
									},
									"side_by_side_index": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "When the objects are listed side-by-side, this value controls the ordering.",
									},
								},
							},
						},
						"service": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The global catalog metadata of the service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rc_provisionable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the service is provisionable by the resource controller service.",
									},
									"iam_compatible": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the service is compatible with the IAM service.",
									},
									"bindable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Deprecated. Controls the Connections tab on the Resource Details page.",
									},
									"plan_updateable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates plan update support and controls the Plan tab on the Resource Details page.",
									},
									"service_key_supported": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates service credentials support and controls the Service Credential tab on Resource Details page.",
									},
									"unique_api_key": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Sensitive:   true,
										Description: "Indicates whether the deployment uses a unique API key or not.",
									},
								},
							},
						},
						"pricing": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The pricing metadata of this object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of the pricing plan.",
									},
									"origin": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The source of the pricing information: global_catalog or pricing_catalog.",
									},
									"sales_avenue": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Currently only available for MCSP subscription plans. Shows how the customer can purchase the plan.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"plan": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Metadata controlling Plan related settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_internal_users": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Controls if IBMers are allowed to provision this plan.",
									},
									"bindable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deprecated. Controls the Connections tab on the Resource Details page.",
									},
									"provision_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SaaS services that integrate with MCSP for provisioning and metering should select `mcsp`. If the subscription is captured as a resource subscription, but SSM/MCSP is not involved, select `ibm_cloud`. Otherwise, select `legacy` if the subscription is only captured as a billable subscription, and there's no resource subscription involved.",
									},
									"reservable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "This field is not available for free and lite plans.",
									},
								},
							},
						},
						"other": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Global catalog plan metadata other field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_controller": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The resource controller metadata.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subscription_provider_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The broker ID for the plan. Only needed if the service is MCSP.",
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
			"geo_tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global catalog URL of your product.",
			},
			"catalog_plan_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a global catalog object.",
			},
		},
	}
}

func ResourceIbmOnboardingCatalogPlanValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "product_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9]{32}:o:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
			MinValueLength:             71,
			MaxValueLength:             71,
		},
		validate.ValidateSchema{
			Identifier:                 "catalog_product_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z\-_\d]+$`,
			MinValueLength:             2,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "env",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9\-.]+$`,
		},
		validate.ValidateSchema{
			Identifier:                 "kind",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "plan",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_catalog_plan", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingCatalogPlanCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createCatalogPlanOptions := &partnercentersellv1.CreateCatalogPlanOptions{}

	createCatalogPlanOptions.SetProductID(d.Get("product_id").(string))
	createCatalogPlanOptions.SetCatalogProductID(d.Get("catalog_product_id").(string))
	createCatalogPlanOptions.SetName(d.Get("name").(string))
	createCatalogPlanOptions.SetActive(d.Get("active").(bool))
	createCatalogPlanOptions.SetDisabled(d.Get("disabled").(bool))
	createCatalogPlanOptions.SetKind(d.Get("kind").(string))
	var tags []string
	for _, v := range d.Get("tags").([]interface{}) {
		tagsItem := v.(string)
		tags = append(tags, tagsItem)
	}
	createCatalogPlanOptions.SetTags(tags)
	objectProviderModel, err := ResourceIbmOnboardingCatalogPlanMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "create", "parse-object_provider").GetDiag()
	}
	createCatalogPlanOptions.SetObjectProvider(objectProviderModel)
	if _, ok := d.GetOk("object_id"); ok {
		createCatalogPlanOptions.SetObjectID(d.Get("object_id").(string))
	}
	if _, ok := d.GetOk("overview_ui"); ok {
		overviewUiModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "create", "parse-overview_ui").GetDiag()
		}
		createCatalogPlanOptions.SetOverviewUi(overviewUiModel)
	}
	if _, ok := d.GetOk("pricing_tags"); ok {
		var pricingTags []string
		for _, v := range d.Get("pricing_tags").([]interface{}) {
			pricingTagsItem := v.(string)
			pricingTags = append(pricingTags, pricingTagsItem)
		}
		createCatalogPlanOptions.SetPricingTags(pricingTags)
	}
	if _, ok := d.GetOk("metadata"); ok {
		metadataModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataPrototypePatch(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "create", "parse-metadata").GetDiag()
		}
		createCatalogPlanOptions.SetMetadata(metadataModel)
	}
	if _, ok := d.GetOk("env"); ok {
		createCatalogPlanOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogPlan, _, err := partnerCenterSellClient.CreateCatalogPlanWithContext(context, createCatalogPlanOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateCatalogPlanWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_plan", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createCatalogPlanOptions.ProductID, *createCatalogPlanOptions.CatalogProductID, *globalCatalogPlan.ID))

	return resourceIbmOnboardingCatalogPlanRead(context, d, meta)
}

func resourceIbmOnboardingCatalogPlanRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogPlanOptions := &partnercentersellv1.GetCatalogPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "sep-id-parts").GetDiag()
	}

	getCatalogPlanOptions.SetProductID(parts[0])
	getCatalogPlanOptions.SetCatalogProductID(parts[1])
	getCatalogPlanOptions.SetCatalogPlanID(parts[2])
	if _, ok := d.GetOk("env"); ok {
		getCatalogPlanOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogPlan, response, err := partnerCenterSellClient.GetCatalogPlanWithContext(context, getCatalogPlanOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogPlanWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_plan", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(globalCatalogPlan.ObjectID) {
		if err = d.Set("object_id", globalCatalogPlan.ObjectID); err != nil {
			err = fmt.Errorf("Error setting object_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-object_id").GetDiag()
		}
	}
	if err = d.Set("name", globalCatalogPlan.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-name").GetDiag()
	}
	if err = d.Set("active", globalCatalogPlan.Active); err != nil {
		err = fmt.Errorf("Error setting active: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-active").GetDiag()
	}
	if err = d.Set("disabled", globalCatalogPlan.Disabled); err != nil {
		err = fmt.Errorf("Error setting disabled: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-disabled").GetDiag()
	}
	if err = d.Set("kind", globalCatalogPlan.Kind); err != nil {
		err = fmt.Errorf("Error setting kind: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-kind").GetDiag()
	}
	if !core.IsNil(globalCatalogPlan.OverviewUi) {
		overviewUiMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUIToMap(globalCatalogPlan.OverviewUi)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "overview_ui-to-map").GetDiag()
		}
		if err = d.Set("overview_ui", []map[string]interface{}{overviewUiMap}); err != nil {
			err = fmt.Errorf("Error setting overview_ui: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-overview_ui").GetDiag()
		}
	}
	if err = d.Set("tags", globalCatalogPlan.Tags); err != nil {
		err = fmt.Errorf("Error setting tags: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-tags").GetDiag()
	}
	if !core.IsNil(globalCatalogPlan.PricingTags) {
		if err = d.Set("pricing_tags", globalCatalogPlan.PricingTags); err != nil {
			err = fmt.Errorf("Error setting pricing_tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-pricing_tags").GetDiag()
		}
	}
	objectProviderMap, err := ResourceIbmOnboardingCatalogPlanCatalogProductProviderToMap(globalCatalogPlan.ObjectProvider)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "object_provider-to-map").GetDiag()
	}
	if err = d.Set("object_provider", []map[string]interface{}{objectProviderMap}); err != nil {
		err = fmt.Errorf("Error setting object_provider: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-object_provider").GetDiag()
	}
	if !core.IsNil(globalCatalogPlan.Metadata) {
		metadataMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataToMap(globalCatalogPlan.Metadata)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "metadata-to-map").GetDiag()
		}
		if err = d.Set("metadata", []map[string]interface{}{metadataMap}); err != nil {
			err = fmt.Errorf("Error setting metadata: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-metadata").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogPlan.GeoTags) {
		if err = d.Set("geo_tags", globalCatalogPlan.GeoTags); err != nil {
			err = fmt.Errorf("Error setting geo_tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-geo_tags").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogPlan.URL) {
		if err = d.Set("url", globalCatalogPlan.URL); err != nil {
			err = fmt.Errorf("Error setting url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-url").GetDiag()
		}
	}
	if parts[0] != "" {
		if err = d.Set("product_id", parts[0]); err != nil {
			err = fmt.Errorf("Error setting product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-product_id").GetDiag()
		}
	}
	if parts[1] != "" {
		if err = d.Set("catalog_product_id", parts[1]); err != nil {
			err = fmt.Errorf("Error setting catalog_product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-catalog_product_id").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogPlan.ID) {
		if err = d.Set("catalog_plan_id", globalCatalogPlan.ID); err != nil {
			err = fmt.Errorf("Error setting catalog_plan_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "read", "set-catalog_plan_id").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingCatalogPlanUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateCatalogPlanOptions := &partnercentersellv1.UpdateCatalogPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "update", "sep-id-parts").GetDiag()
	}

	updateCatalogPlanOptions.SetProductID(parts[0])
	updateCatalogPlanOptions.SetCatalogProductID(parts[1])
	updateCatalogPlanOptions.SetCatalogPlanID(parts[2])
	if _, ok := d.GetOk("env"); ok {
		updateCatalogPlanOptions.SetEnv(d.Get("env").(string))
	}

	hasChange := false

	patchVals := &partnercentersellv1.GlobalCatalogPlanPatch{}
	if d.HasChange("product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_plan", "update", "product_id-forces-new").GetDiag()
	}
	if d.HasChange("catalog_product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "catalog_product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_plan", "update", "catalog_product_id-forces-new").GetDiag()
	}
	if d.HasChange("active") {
		newActive := d.Get("active").(bool)
		patchVals.Active = &newActive
		hasChange = true
	}
	if d.HasChange("disabled") {
		newDisabled := d.Get("disabled").(bool)
		patchVals.Disabled = &newDisabled
		hasChange = true
	}
	if d.HasChange("overview_ui") {
		overviewUi, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "update", "parse-overview_ui").GetDiag()
		}
		patchVals.OverviewUi = overviewUi
		hasChange = true
	}
	if d.HasChange("tags") {
		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tagsItem := v.(string)
			tags = append(tags, tagsItem)
		}
		patchVals.Tags = tags
		hasChange = true
	}
	if d.HasChange("pricing_tags") {
		var pricingTags []string
		for _, v := range d.Get("pricing_tags").([]interface{}) {
			pricingTagsItem := v.(string)
			pricingTags = append(pricingTags, pricingTagsItem)
		}
		patchVals.PricingTags = pricingTags
		hasChange = true
	}
	if d.HasChange("object_provider") {
		objectProvider, err := ResourceIbmOnboardingCatalogPlanMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "update", "parse-object_provider").GetDiag()
		}
		patchVals.ObjectProvider = objectProvider
		hasChange = true
	}
	if d.HasChange("metadata") {
		metadata, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataPrototypePatch(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "update", "parse-metadata").GetDiag()
		}
		patchVals.Metadata = metadata
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateCatalogPlanOptions.GlobalCatalogPlanPatch = ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateCatalogPlanWithContext(context, updateCatalogPlanOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateCatalogPlanWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_plan", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingCatalogPlanRead(context, d, meta)
}

func resourceIbmOnboardingCatalogPlanDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteCatalogPlanOptions := &partnercentersellv1.DeleteCatalogPlanOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_plan", "delete", "sep-id-parts").GetDiag()
	}

	deleteCatalogPlanOptions.SetProductID(parts[0])
	deleteCatalogPlanOptions.SetCatalogProductID(parts[1])
	deleteCatalogPlanOptions.SetCatalogPlanID(parts[2])
	if _, ok := d.GetOk("env"); ok {
		deleteCatalogPlanOptions.SetEnv(d.Get("env").(string))
	}

	_, err = partnerCenterSellClient.DeleteCatalogPlanWithContext(context, deleteCatalogPlanOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteCatalogPlanWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_plan", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingCatalogPlanMapToCatalogProductProvider(modelMap map[string]interface{}) (*partnercentersellv1.CatalogProductProvider, error) {
	model := &partnercentersellv1.CatalogProductProvider{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["email"] != nil && modelMap["email"].(string) != "" {
		model.Email = core.StringPtr(modelMap["email"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogOverviewUI(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUI, error) {
	model := &partnercentersellv1.GlobalCatalogOverviewUI{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 && modelMap["en"].([]interface{})[0] != nil {
		EnModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogOverviewUITranslatedContent(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogOverviewUITranslatedContent(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUITranslatedContent, error) {
	model := &partnercentersellv1.GlobalCatalogOverviewUITranslatedContent{}
	if modelMap["display_name"] != nil && modelMap["display_name"].(string) != "" {
		model.DisplayName = core.StringPtr(modelMap["display_name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["long_description"] != nil && modelMap["long_description"].(string) != "" {
		model.LongDescription = core.StringPtr(modelMap["long_description"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataPrototypePatch(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataPrototypePatch, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataPrototypePatch{}
	if modelMap["rc_compatible"] != nil {
		model.RcCompatible = core.BoolPtr(modelMap["rc_compatible"].(bool))
	}
	if modelMap["ui"] != nil && len(modelMap["ui"].([]interface{})) > 0 && modelMap["ui"].([]interface{})[0] != nil {
		UiModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataUI(modelMap["ui"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ui = UiModel
	}
	if modelMap["service"] != nil && len(modelMap["service"].([]interface{})) > 0 && modelMap["service"].([]interface{})[0] != nil {
		ServiceModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataServicePrototypePatch(modelMap["service"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Service = ServiceModel
	}
	if modelMap["pricing"] != nil && len(modelMap["pricing"].([]interface{})) > 0 && modelMap["pricing"].([]interface{})[0] != nil {
		PricingModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataPricing(modelMap["pricing"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Pricing = PricingModel
	}
	if modelMap["plan"] != nil && len(modelMap["plan"].([]interface{})) > 0 && modelMap["plan"].([]interface{})[0] != nil {
		PlanModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataPlan(modelMap["plan"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Plan = PlanModel
	}
	if modelMap["other"] != nil && len(modelMap["other"].([]interface{})) > 0 {
		OtherModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataOther(modelMap["other"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Other = OtherModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataUI(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataUI, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataUI{}
	if modelMap["strings"] != nil && len(modelMap["strings"].([]interface{})) > 0 && modelMap["strings"].([]interface{})[0] != nil {
		StringsModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIStrings(modelMap["strings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Strings = StringsModel
	}
	if modelMap["urls"] != nil && len(modelMap["urls"].([]interface{})) > 0 && modelMap["urls"].([]interface{})[0] != nil {
		UrlsModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIUrls(modelMap["urls"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Urls = UrlsModel
	}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	if modelMap["side_by_side_index"] != nil {
		model.SideBySideIndex = core.Float64Ptr(modelMap["side_by_side_index"].(float64))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIStrings(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIStrings, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUIStrings{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 && modelMap["en"].([]interface{})[0] != nil {
		EnModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIStringsContent(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIStringsContent(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIStringsContent, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUIStringsContent{}
	if modelMap["bullets"] != nil {
		bullets := []partnercentersellv1.CatalogHighlightItem{}
		for _, bulletsItem := range modelMap["bullets"].([]interface{}) {
			bulletsItemModel, err := ResourceIbmOnboardingCatalogPlanMapToCatalogHighlightItem(bulletsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			bullets = append(bullets, *bulletsItemModel)
		}
		model.Bullets = bullets
	}
	if modelMap["media"] != nil {
		media := []partnercentersellv1.CatalogProductMediaItem{}
		for _, mediaItem := range modelMap["media"].([]interface{}) {
			mediaItemModel, err := ResourceIbmOnboardingCatalogPlanMapToCatalogProductMediaItem(mediaItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			media = append(media, *mediaItemModel)
		}
		model.Media = media
	}
	if modelMap["navigation_items"] != nil {
		navigationItems := []partnercentersellv1.GlobalCatalogMetadataUINavigationItem{}
		for _, navigationItemsItem := range modelMap["navigation_items"].([]interface{}) {
			navigationItemsItemModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUINavigationItem(navigationItemsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			navigationItems = append(navigationItems, *navigationItemsItemModel)
		}
		model.NavigationItems = navigationItems
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToCatalogHighlightItem(modelMap map[string]interface{}) (*partnercentersellv1.CatalogHighlightItem, error) {
	model := &partnercentersellv1.CatalogHighlightItem{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["title"] != nil && modelMap["title"].(string) != "" {
		model.Title = core.StringPtr(modelMap["title"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToCatalogProductMediaItem(modelMap map[string]interface{}) (*partnercentersellv1.CatalogProductMediaItem, error) {
	model := &partnercentersellv1.CatalogProductMediaItem{}
	model.Caption = core.StringPtr(modelMap["caption"].(string))
	if modelMap["thumbnail"] != nil && modelMap["thumbnail"].(string) != "" {
		model.Thumbnail = core.StringPtr(modelMap["thumbnail"].(string))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.URL = core.StringPtr(modelMap["url"].(string))
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUINavigationItem(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUINavigationItem, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUINavigationItem{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataUIUrls(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIUrls, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUIUrls{}
	if modelMap["doc_url"] != nil && modelMap["doc_url"].(string) != "" {
		model.DocURL = core.StringPtr(modelMap["doc_url"].(string))
	}
	if modelMap["apidocs_url"] != nil && modelMap["apidocs_url"].(string) != "" {
		model.ApidocsURL = core.StringPtr(modelMap["apidocs_url"].(string))
	}
	if modelMap["terms_url"] != nil && modelMap["terms_url"].(string) != "" {
		model.TermsURL = core.StringPtr(modelMap["terms_url"].(string))
	}
	if modelMap["instructions_url"] != nil && modelMap["instructions_url"].(string) != "" {
		model.InstructionsURL = core.StringPtr(modelMap["instructions_url"].(string))
	}
	if modelMap["catalog_details_url"] != nil && modelMap["catalog_details_url"].(string) != "" {
		model.CatalogDetailsURL = core.StringPtr(modelMap["catalog_details_url"].(string))
	}
	if modelMap["custom_create_page_url"] != nil && modelMap["custom_create_page_url"].(string) != "" {
		model.CustomCreatePageURL = core.StringPtr(modelMap["custom_create_page_url"].(string))
	}
	if modelMap["dashboard"] != nil && modelMap["dashboard"].(string) != "" {
		model.Dashboard = core.StringPtr(modelMap["dashboard"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataServicePrototypePatch(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataServicePrototypePatch, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataServicePrototypePatch{}
	if modelMap["rc_provisionable"] != nil {
		model.RcProvisionable = core.BoolPtr(modelMap["rc_provisionable"].(bool))
	}
	if modelMap["iam_compatible"] != nil {
		model.IamCompatible = core.BoolPtr(modelMap["iam_compatible"].(bool))
	}
	if modelMap["bindable"] != nil {
		model.Bindable = core.BoolPtr(modelMap["bindable"].(bool))
	}
	if modelMap["plan_updateable"] != nil {
		model.PlanUpdateable = core.BoolPtr(modelMap["plan_updateable"].(bool))
	}
	if modelMap["service_key_supported"] != nil {
		model.ServiceKeySupported = core.BoolPtr(modelMap["service_key_supported"].(bool))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogMetadataPricing(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataPricing, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataPricing{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["origin"] != nil && modelMap["origin"].(string) != "" {
		model.Origin = core.StringPtr(modelMap["origin"].(string))
	}
	if modelMap["sales_avenue"] != nil {
		salesAvenue := []string{}
		for _, salesAvenueItem := range modelMap["sales_avenue"].([]interface{}) {
			salesAvenue = append(salesAvenue, salesAvenueItem.(string))
		}
		model.SalesAvenue = salesAvenue
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataPlan(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataPlan, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataPlan{}
	if modelMap["allow_internal_users"] != nil {
		model.AllowInternalUsers = core.BoolPtr(modelMap["allow_internal_users"].(bool))
	}
	if modelMap["bindable"] != nil {
		model.Bindable = core.BoolPtr(modelMap["bindable"].(bool))
	}
	if modelMap["provision_type"] != nil && modelMap["provision_type"].(string) != "" {
		model.ProvisionType = core.StringPtr(modelMap["provision_type"].(string))
	}
	if modelMap["reservable"] != nil {
		model.Reservable = core.BoolPtr(modelMap["reservable"].(bool))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataOther(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataOther, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataOther{}
	if modelMap["resource_controller"] != nil && len(modelMap["resource_controller"].([]interface{})) > 0 {
		ResourceControllerModel, err := ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataOtherResourceController(modelMap["resource_controller"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceController = ResourceControllerModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanMapToGlobalCatalogPlanMetadataOtherResourceController(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogPlanMetadataOtherResourceController, error) {
	model := &partnercentersellv1.GlobalCatalogPlanMetadataOtherResourceController{}
	if modelMap["subscription_provider_id"] != nil && modelMap["subscription_provider_id"].(string) != "" {
		model.SubscriptionProviderID = core.StringPtr(modelMap["subscription_provider_id"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUIToMap(model *partnercentersellv1.GlobalCatalogOverviewUI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUITranslatedContentToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUITranslatedContentToMap(model *partnercentersellv1.GlobalCatalogOverviewUITranslatedContent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DisplayName != nil {
		modelMap["display_name"] = *model.DisplayName
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.LongDescription != nil {
		modelMap["long_description"] = *model.LongDescription
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanCatalogProductProviderToMap(model *partnercentersellv1.CatalogProductProvider) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataToMap(model *partnercentersellv1.GlobalCatalogPlanMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RcCompatible != nil {
		modelMap["rc_compatible"] = *model.RcCompatible
	}
	if model.Ui != nil {
		uiMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataUIToMap(model.Ui)
		if err != nil {
			return modelMap, err
		}
		modelMap["ui"] = []map[string]interface{}{uiMap}
	}
	if model.Service != nil {
		serviceMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataServiceToMap(model.Service)
		if err != nil {
			return modelMap, err
		}
		modelMap["service"] = []map[string]interface{}{serviceMap}
	}
	if model.Pricing != nil {
		pricingMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataPricingToMap(model.Pricing)
		if err != nil {
			return modelMap, err
		}
		modelMap["pricing"] = []map[string]interface{}{pricingMap}
	}
	if model.Plan != nil {
		planMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPlanToMap(model.Plan)
		if err != nil {
			return modelMap, err
		}
		modelMap["plan"] = []map[string]interface{}{planMap}
	}
	if model.Other != nil {
		otherMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherToMap(model.Other)
		if err != nil {
			return modelMap, err
		}
		modelMap["other"] = []map[string]interface{}{otherMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataUIToMap(model *partnercentersellv1.GlobalCatalogPlanMetadataUI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Strings != nil {
		stringsMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsToMap(model.Strings)
		if err != nil {
			return modelMap, err
		}
		modelMap["strings"] = []map[string]interface{}{stringsMap}
	}
	if model.Urls != nil {
		urlsMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIUrlsToMap(model.Urls)
		if err != nil {
			return modelMap, err
		}
		modelMap["urls"] = []map[string]interface{}{urlsMap}
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.SideBySideIndex != nil {
		modelMap["side_by_side_index"] = *model.SideBySideIndex
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsToMap(model *partnercentersellv1.GlobalCatalogMetadataUIStrings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsContentToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsContentToMap(model *partnercentersellv1.GlobalCatalogMetadataUIStringsContent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Bullets != nil {
		bullets := []map[string]interface{}{}
		for _, bulletsItem := range model.Bullets {
			bulletsItemMap, err := ResourceIbmOnboardingCatalogPlanCatalogHighlightItemToMap(&bulletsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			bullets = append(bullets, bulletsItemMap)
		}
		modelMap["bullets"] = bullets
	}
	if model.Media != nil {
		media := []map[string]interface{}{}
		for _, mediaItem := range model.Media {
			mediaItemMap, err := ResourceIbmOnboardingCatalogPlanCatalogProductMediaItemToMap(&mediaItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			media = append(media, mediaItemMap)
		}
		modelMap["media"] = media
	}
	if model.NavigationItems != nil {
		navigationItems := []map[string]interface{}{}
		for _, navigationItemsItem := range model.NavigationItems {
			navigationItemsItemMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUINavigationItemToMap(&navigationItemsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			navigationItems = append(navigationItems, navigationItemsItemMap)
		}
		modelMap["navigation_items"] = navigationItems
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanCatalogHighlightItemToMap(model *partnercentersellv1.CatalogHighlightItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Title != nil {
		modelMap["title"] = *model.Title
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanCatalogProductMediaItemToMap(model *partnercentersellv1.CatalogProductMediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["caption"] = *model.Caption
	if model.Thumbnail != nil {
		modelMap["thumbnail"] = *model.Thumbnail
	}
	modelMap["type"] = *model.Type
	modelMap["url"] = *model.URL
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUINavigationItemToMap(model *partnercentersellv1.GlobalCatalogMetadataUINavigationItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIUrlsToMap(model *partnercentersellv1.GlobalCatalogMetadataUIUrls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DocURL != nil {
		modelMap["doc_url"] = *model.DocURL
	}
	if model.ApidocsURL != nil {
		modelMap["apidocs_url"] = *model.ApidocsURL
	}
	if model.TermsURL != nil {
		modelMap["terms_url"] = *model.TermsURL
	}
	if model.InstructionsURL != nil {
		modelMap["instructions_url"] = *model.InstructionsURL
	}
	if model.CatalogDetailsURL != nil {
		modelMap["catalog_details_url"] = *model.CatalogDetailsURL
	}
	if model.CustomCreatePageURL != nil {
		modelMap["custom_create_page_url"] = *model.CustomCreatePageURL
	}
	if model.Dashboard != nil {
		modelMap["dashboard"] = *model.Dashboard
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataServiceToMap(model *partnercentersellv1.GlobalCatalogPlanMetadataService) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RcProvisionable != nil {
		modelMap["rc_provisionable"] = *model.RcProvisionable
	}
	if model.IamCompatible != nil {
		modelMap["iam_compatible"] = *model.IamCompatible
	}
	if model.Bindable != nil {
		modelMap["bindable"] = *model.Bindable
	}
	if model.PlanUpdateable != nil {
		modelMap["plan_updateable"] = *model.PlanUpdateable
	}
	if model.ServiceKeySupported != nil {
		modelMap["service_key_supported"] = *model.ServiceKeySupported
	}
	if model.UniqueApiKey != nil {
		modelMap["unique_api_key"] = *model.UniqueApiKey
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataPricingToMap(model *partnercentersellv1.GlobalCatalogMetadataPricing) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Origin != nil {
		modelMap["origin"] = *model.Origin
	}
	if model.SalesAvenue != nil {
		modelMap["sales_avenue"] = model.SalesAvenue
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPlanToMap(model *partnercentersellv1.GlobalCatalogPlanMetadataPlan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AllowInternalUsers != nil {
		modelMap["allow_internal_users"] = *model.AllowInternalUsers
	}
	if model.Bindable != nil {
		modelMap["bindable"] = *model.Bindable
	}
	if model.ProvisionType != nil {
		modelMap["provision_type"] = *model.ProvisionType
	}
	if model.Reservable != nil {
		modelMap["reservable"] = *model.Reservable
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherToMap(model *partnercentersellv1.GlobalCatalogPlanMetadataOther) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceController != nil {
		resourceControllerMap, err := ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherResourceControllerToMap(model.ResourceController)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_controller"] = []map[string]interface{}{resourceControllerMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherResourceControllerToMap(model *partnercentersellv1.GlobalCatalogPlanMetadataOtherResourceController) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SubscriptionProviderID != nil {
		modelMap["subscription_provider_id"] = *model.SubscriptionProviderID
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanPatchAsPatch(patchVals *partnercentersellv1.GlobalCatalogPlanPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "active"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["active"] = nil
	} else if !exists {
		delete(patch, "active")
	}
	path = "disabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["disabled"] = nil
	} else if !exists {
		delete(patch, "disabled")
	}
	path = "overview_ui"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["overview_ui"] = nil
	} else if exists && patch["overview_ui"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUIAsPatch(patch["overview_ui"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "overview_ui")
	}
	path = "tags"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tags"] = nil
	} else if !exists {
		delete(patch, "tags")
	}
	path = "pricing_tags"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["pricing_tags"] = nil
	} else if !exists {
		delete(patch, "pricing_tags")
	}
	path = "object_provider"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["object_provider"] = nil
	} else if exists && patch["object_provider"] != nil {
		ResourceIbmOnboardingCatalogPlanCatalogProductProviderAsPatch(patch["object_provider"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "object_provider")
	}
	path = "metadata"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["metadata"] = nil
	} else if exists && patch["metadata"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPrototypePatchAsPatch(patch["metadata"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "metadata")
	}

	return patch
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPrototypePatchAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".rc_compatible"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["rc_compatible"] = nil
	} else if !exists {
		delete(patch, "rc_compatible")
	}
	path = rootPath + ".ui"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ui"] = nil
	} else if exists && patch["ui"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataUIAsPatch(patch["ui"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "ui")
	}
	path = rootPath + ".service"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service"] = nil
	} else if exists && patch["service"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataServicePrototypePatchAsPatch(patch["service"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "service")
	}
	path = rootPath + ".pricing"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["pricing"] = nil
	} else if exists && patch["pricing"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataPricingAsPatch(patch["pricing"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "pricing")
	}
	path = rootPath + ".plan"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["plan"] = nil
	} else if exists && patch["plan"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPlanAsPatch(patch["plan"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "plan")
	}
	path = rootPath + ".other"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["other"] = nil
	} else if exists && patch["other"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherAsPatch(patch["other"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "other")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".resource_controller"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["resource_controller"] = nil
	} else if exists && patch["resource_controller"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherResourceControllerAsPatch(patch["resource_controller"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "resource_controller")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataOtherResourceControllerAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".subscription_provider_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["subscription_provider_id"] = nil
	} else if !exists {
		delete(patch, "subscription_provider_id")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataPlanAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".allow_internal_users"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["allow_internal_users"] = nil
	} else if !exists {
		delete(patch, "allow_internal_users")
	}
	path = rootPath + ".provision_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["provision_type"] = nil
	} else if !exists {
		delete(patch, "provision_type")
	}
	path = rootPath + ".reservable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["reservable"] = nil
	} else if !exists {
		delete(patch, "reservable")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataPricingAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	} else if !exists {
		delete(patch, "type")
	}
	path = rootPath + ".origin"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["origin"] = nil
	} else if !exists {
		delete(patch, "origin")
	}
	path = rootPath + ".sales_avenue"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["sales_avenue"] = nil
	} else if !exists {
		delete(patch, "sales_avenue")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataServicePrototypePatchAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".rc_provisionable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["rc_provisionable"] = nil
	} else if !exists {
		delete(patch, "rc_provisionable")
	}
	path = rootPath + ".iam_compatible"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["iam_compatible"] = nil
	} else if !exists {
		delete(patch, "iam_compatible")
	}
	path = rootPath + ".bindable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["bindable"] = nil
	} else if !exists {
		delete(patch, "bindable")
	}
	path = rootPath + ".plan_updateable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["plan_updateable"] = nil
	} else if !exists {
		delete(patch, "plan_updateable")
	}
	path = rootPath + ".service_key_supported"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service_key_supported"] = nil
	} else if !exists {
		delete(patch, "service_key_supported")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogPlanMetadataUIAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".strings"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["strings"] = nil
	} else if exists && patch["strings"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsAsPatch(patch["strings"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "strings")
	}
	path = rootPath + ".urls"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["urls"] = nil
	} else if exists && patch["urls"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIUrlsAsPatch(patch["urls"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "urls")
	}
	path = rootPath + ".hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	} else if !exists {
		delete(patch, "hidden")
	}
	path = rootPath + ".side_by_side_index"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["side_by_side_index"] = nil
	} else if !exists {
		delete(patch, "side_by_side_index")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIUrlsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".doc_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["doc_url"] = nil
	} else if !exists {
		delete(patch, "doc_url")
	}
	path = rootPath + ".apidocs_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["apidocs_url"] = nil
	} else if !exists {
		delete(patch, "apidocs_url")
	}
	path = rootPath + ".terms_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["terms_url"] = nil
	} else if !exists {
		delete(patch, "terms_url")
	}
	path = rootPath + ".instructions_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["instructions_url"] = nil
	} else if !exists {
		delete(patch, "instructions_url")
	}
	path = rootPath + ".catalog_details_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["catalog_details_url"] = nil
	} else if !exists {
		delete(patch, "catalog_details_url")
	}
	path = rootPath + ".custom_create_page_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["custom_create_page_url"] = nil
	} else if !exists {
		delete(patch, "custom_create_page_url")
	}
	path = rootPath + ".dashboard"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["dashboard"] = nil
	} else if !exists {
		delete(patch, "dashboard")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsContentAsPatch(patch["en"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "en")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUIStringsContentAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".bullets"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["bullets"] = nil
	} else if exists && patch["bullets"] != nil {
		bulletsList := patch["bullets"].([]map[string]interface{})
		for i, bulletsItem := range bulletsList {
			ResourceIbmOnboardingCatalogPlanCatalogHighlightItemAsPatch(bulletsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "bullets")
	}
	path = rootPath + ".media"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["media"] = nil
	} else if exists && patch["media"] != nil {
		mediaList := patch["media"].([]map[string]interface{})
		for i, mediaItem := range mediaList {
			ResourceIbmOnboardingCatalogPlanCatalogProductMediaItemAsPatch(mediaItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "media")
	}
	path = rootPath + ".navigation_items"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["navigation_items"] = nil
	} else if exists && patch["navigation_items"] != nil {
		navigation_itemsList := patch["navigation_items"].([]map[string]interface{})
		for i, navigation_itemsItem := range navigation_itemsList {
			ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUINavigationItemAsPatch(navigation_itemsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "navigation_items")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogMetadataUINavigationItemAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["id"] = nil
	} else if !exists {
		delete(patch, "id")
	}
	path = rootPath + ".url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["url"] = nil
	} else if !exists {
		delete(patch, "url")
	}
	path = rootPath + ".label"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["label"] = nil
	} else if !exists {
		delete(patch, "label")
	}
}

func ResourceIbmOnboardingCatalogPlanCatalogProductMediaItemAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".thumbnail"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["thumbnail"] = nil
	} else if !exists {
		delete(patch, "thumbnail")
	}
}

func ResourceIbmOnboardingCatalogPlanCatalogHighlightItemAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".title"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["title"] = nil
	} else if !exists {
		delete(patch, "title")
	}
}

func ResourceIbmOnboardingCatalogPlanCatalogProductProviderAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = rootPath + ".email"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["email"] = nil
	} else if !exists {
		delete(patch, "email")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUIAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUITranslatedContentAsPatch(patch["en"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "en")
	}
}

func ResourceIbmOnboardingCatalogPlanGlobalCatalogOverviewUITranslatedContentAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if !exists {
		delete(patch, "display_name")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".long_description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["long_description"] = nil
	} else if !exists {
		delete(patch, "long_description")
	}
}
