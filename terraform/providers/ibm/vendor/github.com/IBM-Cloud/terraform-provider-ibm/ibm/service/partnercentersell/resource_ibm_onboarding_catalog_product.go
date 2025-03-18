// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
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

func ResourceIbmOnboardingCatalogProduct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingCatalogProductCreate,
		ReadContext:   resourceIbmOnboardingCatalogProductRead,
		UpdateContext: resourceIbmOnboardingCatalogProductUpdate,
		DeleteContext: resourceIbmOnboardingCatalogProductDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_product", "product_id"),
				Description:  "The unique ID of the product.",
			},
			"env": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_product", "env"),
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
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_product", "name"),
				Description:  "The programmatic name of this product.",
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
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_product", "kind"),
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
			"images": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Images from the global catalog entry that help illustrate the service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL for your product logo.",
						},
					},
				},
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
				Description: "The global catalog service metadata object.",
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
																		"description_i18n": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "The description about the features of the product in translation.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"title": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The descriptive title for the feature.",
																		},
																		"title_i18n": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "The descriptive title for the feature in translation.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
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
																		"caption_i18n": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "The brief explanation for your images and videos in translation.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
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
															"embeddable_dashboard": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "On a service kind record this controls if your service has a custom dashboard or Resource Detail page.",
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
								},
							},
						},
						"other": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The additional metadata of the service in global catalog.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pc": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The metadata of the service owned and managed by Partner Center - Sell.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"support": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The support metadata of the service.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The support site URL where the support for your service is available.",
															},
															"status_url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The URL where the status of your service is available.",
															},
															"locations": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The countries in which your support is available. Provide a list of country codes.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"languages": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The languages in which support is available.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"process": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The description of your support process.",
															},
															"process_i18n": &schema.Schema{
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "The description of your support process.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"support_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The type of support provided.",
															},
															"support_escalation": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The details of the support escalation process.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"contact": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The support contact information of the escalation team.",
																		},
																		"escalation_wait_time": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The time interval of providing support in units and values.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"value": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Optional:    true,
																						Description: "The number of time units.",
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The unit of the time.",
																					},
																				},
																			},
																		},
																		"response_wait_time": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The time interval of providing support in units and values.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"value": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Optional:    true,
																						Description: "The number of time units.",
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The unit of the time.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"support_details": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "The support options for the service.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The type of support for this support channel.",
																		},
																		"contact": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The contact information for this support channel.",
																		},
																		"response_wait_time": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The time interval of providing support in units and values.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"value": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Optional:    true,
																						Description: "The number of time units.",
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The unit of the time.",
																					},
																				},
																			},
																		},
																		"availability": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The time period during which support is available for the service.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"times": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "The support hours available for the service.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"day": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Optional:    true,
																									Description: "The number of days in a week when support is available for the service.",
																								},
																								"start_time": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The time in the day when support starts for the service.",
																								},
																								"end_time": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The time in the day when support ends for the service.",
																								},
																							},
																						},
																					},
																					"timezone": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The timezones in which support is available. Only relevant if `always_available` is set to false.",
																					},
																					"always_available": &schema.Schema{
																						Type:        schema.TypeBool,
																						Optional:    true,
																						Description: "Whether the support for the service is always available.",
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
										},
									},
									"composite": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Optional metadata of the service defining it as a composite.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"composite_kind": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The type of the composite service.",
												},
												"composite_tag": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The tag used for the composite parent and its children.",
												},
												"children": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The type of the composite child.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The name of the composite child.",
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
				},
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global catalog URL of your product.",
			},
			"catalog_product_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a global catalog object.",
			},
		},
	}
}

func ResourceIbmOnboardingCatalogProductValidator() *validate.ResourceValidator {
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
			Regexp:                     `^[a-z0-9\-.]+$`,
		},
		validate.ValidateSchema{
			Identifier:                 "kind",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "composite, platform_service, service",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_catalog_product", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingCatalogProductCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createCatalogProductOptions := &partnercentersellv1.CreateCatalogProductOptions{}

	createCatalogProductOptions.SetProductID(d.Get("product_id").(string))
	createCatalogProductOptions.SetName(d.Get("name").(string))
	createCatalogProductOptions.SetActive(d.Get("active").(bool))
	createCatalogProductOptions.SetDisabled(d.Get("disabled").(bool))
	createCatalogProductOptions.SetKind(d.Get("kind").(string))
	var tags []string
	for _, v := range d.Get("tags").([]interface{}) {
		tagsItem := v.(string)
		tags = append(tags, tagsItem)
	}
	createCatalogProductOptions.SetTags(tags)
	objectProviderModel, err := ResourceIbmOnboardingCatalogProductMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "create", "parse-object_provider").GetDiag()
	}
	createCatalogProductOptions.SetObjectProvider(objectProviderModel)
	if _, ok := d.GetOk("object_id"); ok {
		createCatalogProductOptions.SetObjectID(d.Get("object_id").(string))
	}
	if _, ok := d.GetOk("overview_ui"); ok {
		overviewUiModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "create", "parse-overview_ui").GetDiag()
		}
		createCatalogProductOptions.SetOverviewUi(overviewUiModel)
	}
	if _, ok := d.GetOk("images"); ok {
		imagesModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductImages(d.Get("images.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "create", "parse-images").GetDiag()
		}
		createCatalogProductOptions.SetImages(imagesModel)
	}
	if _, ok := d.GetOk("metadata"); ok {
		metadataModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadata(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "create", "parse-metadata").GetDiag()
		}
		createCatalogProductOptions.SetMetadata(metadataModel)
	}
	if _, ok := d.GetOk("env"); ok {
		createCatalogProductOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogProduct, _, err := partnerCenterSellClient.CreateCatalogProductWithContext(context, createCatalogProductOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateCatalogProductWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_product", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createCatalogProductOptions.ProductID, *globalCatalogProduct.ID))

	return resourceIbmOnboardingCatalogProductRead(context, d, meta)
}

func resourceIbmOnboardingCatalogProductRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogProductOptions := &partnercentersellv1.GetCatalogProductOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "sep-id-parts").GetDiag()
	}

	getCatalogProductOptions.SetProductID(parts[0])
	getCatalogProductOptions.SetCatalogProductID(parts[1])
	if _, ok := d.GetOk("env"); ok {
		getCatalogProductOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogProduct, response, err := partnerCenterSellClient.GetCatalogProductWithContext(context, getCatalogProductOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogProductWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_product", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(globalCatalogProduct.ObjectID) {
		if err = d.Set("object_id", globalCatalogProduct.ObjectID); err != nil {
			err = fmt.Errorf("Error setting object_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-object_id").GetDiag()
		}
	}
	if err = d.Set("name", globalCatalogProduct.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-name").GetDiag()
	}
	if err = d.Set("active", globalCatalogProduct.Active); err != nil {
		err = fmt.Errorf("Error setting active: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-active").GetDiag()
	}
	if err = d.Set("disabled", globalCatalogProduct.Disabled); err != nil {
		err = fmt.Errorf("Error setting disabled: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-disabled").GetDiag()
	}
	if err = d.Set("kind", globalCatalogProduct.Kind); err != nil {
		err = fmt.Errorf("Error setting kind: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-kind").GetDiag()
	}
	if !core.IsNil(globalCatalogProduct.OverviewUi) {
		overviewUiMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUIToMap(globalCatalogProduct.OverviewUi)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "overview_ui-to-map").GetDiag()
		}
		if err = d.Set("overview_ui", []map[string]interface{}{overviewUiMap}); err != nil {
			err = fmt.Errorf("Error setting overview_ui: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-overview_ui").GetDiag()
		}
	}
	if err = d.Set("tags", globalCatalogProduct.Tags); err != nil {
		err = fmt.Errorf("Error setting tags: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-tags").GetDiag()
	}
	if !core.IsNil(globalCatalogProduct.Images) {
		imagesMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductImagesToMap(globalCatalogProduct.Images)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "images-to-map").GetDiag()
		}
		if err = d.Set("images", []map[string]interface{}{imagesMap}); err != nil {
			err = fmt.Errorf("Error setting images: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-images").GetDiag()
		}
	}
	objectProviderMap, err := ResourceIbmOnboardingCatalogProductCatalogProductProviderToMap(globalCatalogProduct.ObjectProvider)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "object_provider-to-map").GetDiag()
	}
	if err = d.Set("object_provider", []map[string]interface{}{objectProviderMap}); err != nil {
		err = fmt.Errorf("Error setting object_provider: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-object_provider").GetDiag()
	}
	if !core.IsNil(globalCatalogProduct.Metadata) {
		metadataMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataToMap(globalCatalogProduct.Metadata)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "metadata-to-map").GetDiag()
		}
		if err = d.Set("metadata", []map[string]interface{}{metadataMap}); err != nil {
			err = fmt.Errorf("Error setting metadata: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-metadata").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogProduct.URL) {
		if err = d.Set("url", globalCatalogProduct.URL); err != nil {
			err = fmt.Errorf("Error setting url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-url").GetDiag()
		}
	}
	if parts[0] != "" {
		if err = d.Set("product_id", parts[0]); err != nil {
			err = fmt.Errorf("Error setting product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-product_id").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogProduct.ID) {
		if err = d.Set("catalog_product_id", globalCatalogProduct.ID); err != nil {
			err = fmt.Errorf("Error setting catalog_product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "read", "set-catalog_product_id").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingCatalogProductUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateCatalogProductOptions := &partnercentersellv1.UpdateCatalogProductOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "sep-id-parts").GetDiag()
	}

	updateCatalogProductOptions.SetProductID(parts[0])
	updateCatalogProductOptions.SetCatalogProductID(parts[1])
	if _, ok := d.GetOk("env"); ok {
		updateCatalogProductOptions.SetEnv(d.Get("env").(string))
	}

	hasChange := false

	patchVals := &partnercentersellv1.GlobalCatalogProductPatch{}
	if d.HasChange("product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_product", "update", "product_id-forces-new").GetDiag()
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
		overviewUi, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "parse-overview_ui").GetDiag()
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
	if d.HasChange("images") {
		images, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductImages(d.Get("images.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "parse-images").GetDiag()
		}
		patchVals.Images = images
		hasChange = true
	}
	if d.HasChange("object_provider") {
		objectProvider, err := ResourceIbmOnboardingCatalogProductMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "parse-object_provider").GetDiag()
		}
		patchVals.ObjectProvider = objectProvider
		hasChange = true
	}
	if d.HasChange("metadata") {
		metadata, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadata(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "update", "parse-metadata").GetDiag()
		}
		patchVals.Metadata = metadata
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateCatalogProductOptions.GlobalCatalogProductPatch = ResourceIbmOnboardingCatalogProductGlobalCatalogProductPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateCatalogProductWithContext(context, updateCatalogProductOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateCatalogProductWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_product", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingCatalogProductRead(context, d, meta)
}

func resourceIbmOnboardingCatalogProductDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteCatalogProductOptions := &partnercentersellv1.DeleteCatalogProductOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_product", "delete", "sep-id-parts").GetDiag()
	}

	deleteCatalogProductOptions.SetProductID(parts[0])
	deleteCatalogProductOptions.SetCatalogProductID(parts[1])
	if _, ok := d.GetOk("env"); ok {
		deleteCatalogProductOptions.SetEnv(d.Get("env").(string))
	}

	_, err = partnerCenterSellClient.DeleteCatalogProductWithContext(context, deleteCatalogProductOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteCatalogProductWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_product", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingCatalogProductMapToCatalogProductProvider(modelMap map[string]interface{}) (*partnercentersellv1.CatalogProductProvider, error) {
	model := &partnercentersellv1.CatalogProductProvider{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["email"] != nil && modelMap["email"].(string) != "" {
		model.Email = core.StringPtr(modelMap["email"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogOverviewUI(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUI, error) {
	model := &partnercentersellv1.GlobalCatalogOverviewUI{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 && modelMap["en"].([]interface{})[0] != nil {
		EnModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogOverviewUITranslatedContent(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogOverviewUITranslatedContent(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUITranslatedContent, error) {
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

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductImages(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductImages, error) {
	model := &partnercentersellv1.GlobalCatalogProductImages{}
	if modelMap["image"] != nil && modelMap["image"].(string) != "" {
		model.Image = core.StringPtr(modelMap["image"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadata(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadata, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadata{}
	if modelMap["rc_compatible"] != nil {
		model.RcCompatible = core.BoolPtr(modelMap["rc_compatible"].(bool))
	}
	if modelMap["ui"] != nil && len(modelMap["ui"].([]interface{})) > 0 && modelMap["ui"].([]interface{})[0] != nil {
		UiModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUI(modelMap["ui"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ui = UiModel
	}
	if modelMap["service"] != nil && len(modelMap["service"].([]interface{})) > 0 && modelMap["service"].([]interface{})[0] != nil {
		ServiceModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataService(modelMap["service"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Service = ServiceModel
	}
	if modelMap["other"] != nil && len(modelMap["other"].([]interface{})) > 0 && modelMap["other"].([]interface{})[0] != nil {
		OtherModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOther(modelMap["other"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Other = OtherModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUI(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUI, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUI{}
	if modelMap["strings"] != nil && len(modelMap["strings"].([]interface{})) > 0 && modelMap["strings"].([]interface{})[0] != nil {
		StringsModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIStrings(modelMap["strings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Strings = StringsModel
	}
	if modelMap["urls"] != nil && len(modelMap["urls"].([]interface{})) > 0 && modelMap["urls"].([]interface{})[0] != nil {
		UrlsModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIUrls(modelMap["urls"].([]interface{})[0].(map[string]interface{}))
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

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIStrings(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIStrings, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUIStrings{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 && modelMap["en"].([]interface{})[0] != nil {
		EnModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIStringsContent(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIStringsContent(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIStringsContent, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataUIStringsContent{}
	if modelMap["bullets"] != nil {
		bullets := []partnercentersellv1.CatalogHighlightItem{}
		for _, bulletsItem := range modelMap["bullets"].([]interface{}) {
			bulletsItemModel, err := ResourceIbmOnboardingCatalogProductMapToCatalogHighlightItem(bulletsItem.(map[string]interface{}))
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
			mediaItemModel, err := ResourceIbmOnboardingCatalogProductMapToCatalogProductMediaItem(mediaItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			media = append(media, *mediaItemModel)
		}
		model.Media = media
	}
	if modelMap["embeddable_dashboard"] != nil && modelMap["embeddable_dashboard"].(string) != "" {
		model.EmbeddableDashboard = core.StringPtr(modelMap["embeddable_dashboard"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToCatalogHighlightItem(modelMap map[string]interface{}) (*partnercentersellv1.CatalogHighlightItem, error) {
	model := &partnercentersellv1.CatalogHighlightItem{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["description_i18n"] != nil {
		model.DescriptionI18n = make(map[string]string)
		for key, value := range modelMap["description_i18n"].(map[string]interface{}) {
			if str, ok := value.(string); ok {
				model.DescriptionI18n[key] = str
			}
		}
	}
	if modelMap["title"] != nil && modelMap["title"].(string) != "" {
		model.Title = core.StringPtr(modelMap["title"].(string))
	}
	if modelMap["title_i18n"] != nil {
		model.TitleI18n = make(map[string]string)
		for key, value := range modelMap["title_i18n"].(map[string]interface{}) {
			if str, ok := value.(string); ok {
				model.TitleI18n[key] = str
			}
		}
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToCatalogProductMediaItem(modelMap map[string]interface{}) (*partnercentersellv1.CatalogProductMediaItem, error) {
	model := &partnercentersellv1.CatalogProductMediaItem{}
	model.Caption = core.StringPtr(modelMap["caption"].(string))
	if modelMap["caption_i18n"] != nil {
		model.CaptionI18n = make(map[string]string)
		for key, value := range modelMap["caption_i18n"].(map[string]interface{}) {
			if str, ok := value.(string); ok {
				model.CaptionI18n[key] = str
			}
		}
	}
	if modelMap["thumbnail"] != nil && modelMap["thumbnail"].(string) != "" {
		model.Thumbnail = core.StringPtr(modelMap["thumbnail"].(string))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.URL = core.StringPtr(modelMap["url"].(string))
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataUIUrls(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataUIUrls, error) {
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

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogMetadataService(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataService, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataService{}
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

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOther(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadataOther, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadataOther{}
	if modelMap["pc"] != nil && len(modelMap["pc"].([]interface{})) > 0 && modelMap["pc"].([]interface{})[0] != nil {
		PCModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherPC(modelMap["pc"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PC = PCModel
	}
	if modelMap["composite"] != nil && len(modelMap["composite"].([]interface{})) > 0 {
		CompositeModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherComposite(modelMap["composite"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Composite = CompositeModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherPC(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadataOtherPC, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadataOtherPC{}
	if modelMap["support"] != nil && len(modelMap["support"].([]interface{})) > 0 && modelMap["support"].([]interface{})[0] != nil {
		SupportModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherPCSupport(modelMap["support"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Support = SupportModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherPCSupport(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadataOtherPCSupport, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadataOtherPCSupport{}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["status_url"] != nil && modelMap["status_url"].(string) != "" {
		model.StatusURL = core.StringPtr(modelMap["status_url"].(string))
	}
	if modelMap["locations"] != nil {
		locations := []string{}
		for _, locationsItem := range modelMap["locations"].([]interface{}) {
			locations = append(locations, locationsItem.(string))
		}
		model.Locations = locations
	}
	if modelMap["languages"] != nil {
		languages := []string{}
		for _, languagesItem := range modelMap["languages"].([]interface{}) {
			languages = append(languages, languagesItem.(string))
		}
		model.Languages = languages
	}
	if modelMap["process"] != nil && modelMap["process"].(string) != "" {
		model.Process = core.StringPtr(modelMap["process"].(string))
	}
	if modelMap["process_i18n"] != nil {
		model.ProcessI18n = make(map[string]string)
		for key, value := range modelMap["process_i18n"].(map[string]interface{}) {
			if str, ok := value.(string); ok {
				model.ProcessI18n[key] = str
			}
		}
	}
	if modelMap["support_type"] != nil && modelMap["support_type"].(string) != "" {
		model.SupportType = core.StringPtr(modelMap["support_type"].(string))
	}
	if modelMap["support_escalation"] != nil && len(modelMap["support_escalation"].([]interface{})) > 0 && modelMap["support_escalation"].([]interface{})[0] != nil {
		SupportEscalationModel, err := ResourceIbmOnboardingCatalogProductMapToSupportEscalation(modelMap["support_escalation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SupportEscalation = SupportEscalationModel
	}
	if modelMap["support_details"] != nil {
		supportDetails := []partnercentersellv1.SupportDetailsItem{}
		for _, supportDetailsItem := range modelMap["support_details"].([]interface{}) {
			supportDetailsItemModel, err := ResourceIbmOnboardingCatalogProductMapToSupportDetailsItem(supportDetailsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			supportDetails = append(supportDetails, *supportDetailsItemModel)
		}
		model.SupportDetails = supportDetails
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToSupportEscalation(modelMap map[string]interface{}) (*partnercentersellv1.SupportEscalation, error) {
	model := &partnercentersellv1.SupportEscalation{}
	if modelMap["contact"] != nil && modelMap["contact"].(string) != "" {
		model.Contact = core.StringPtr(modelMap["contact"].(string))
	}
	if modelMap["escalation_wait_time"] != nil && len(modelMap["escalation_wait_time"].([]interface{})) > 0 && modelMap["escalation_wait_time"].([]interface{})[0] != nil {
		EscalationWaitTimeModel, err := ResourceIbmOnboardingCatalogProductMapToSupportTimeInterval(modelMap["escalation_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.EscalationWaitTime = EscalationWaitTimeModel
	}
	if modelMap["response_wait_time"] != nil && len(modelMap["response_wait_time"].([]interface{})) > 0 && modelMap["response_wait_time"].([]interface{})[0] != nil {
		ResponseWaitTimeModel, err := ResourceIbmOnboardingCatalogProductMapToSupportTimeInterval(modelMap["response_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResponseWaitTime = ResponseWaitTimeModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToSupportTimeInterval(modelMap map[string]interface{}) (*partnercentersellv1.SupportTimeInterval, error) {
	model := &partnercentersellv1.SupportTimeInterval{}
	if modelMap["value"] != nil {
		model.Value = core.Float64Ptr(modelMap["value"].(float64))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToSupportDetailsItem(modelMap map[string]interface{}) (*partnercentersellv1.SupportDetailsItem, error) {
	model := &partnercentersellv1.SupportDetailsItem{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["contact"] != nil && modelMap["contact"].(string) != "" {
		model.Contact = core.StringPtr(modelMap["contact"].(string))
	}
	if modelMap["response_wait_time"] != nil && len(modelMap["response_wait_time"].([]interface{})) > 0 && modelMap["response_wait_time"].([]interface{})[0] != nil {
		ResponseWaitTimeModel, err := ResourceIbmOnboardingCatalogProductMapToSupportTimeInterval(modelMap["response_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResponseWaitTime = ResponseWaitTimeModel
	}
	if modelMap["availability"] != nil && len(modelMap["availability"].([]interface{})) > 0 && modelMap["availability"].([]interface{})[0] != nil {
		AvailabilityModel, err := ResourceIbmOnboardingCatalogProductMapToSupportDetailsItemAvailability(modelMap["availability"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Availability = AvailabilityModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToSupportDetailsItemAvailability(modelMap map[string]interface{}) (*partnercentersellv1.SupportDetailsItemAvailability, error) {
	model := &partnercentersellv1.SupportDetailsItemAvailability{}
	if modelMap["times"] != nil {
		times := []partnercentersellv1.SupportDetailsItemAvailabilityTime{}
		for _, timesItem := range modelMap["times"].([]interface{}) {
			timesItemModel, err := ResourceIbmOnboardingCatalogProductMapToSupportDetailsItemAvailabilityTime(timesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			times = append(times, *timesItemModel)
		}
		model.Times = times
	}
	if modelMap["timezone"] != nil && modelMap["timezone"].(string) != "" {
		model.Timezone = core.StringPtr(modelMap["timezone"].(string))
	}
	if modelMap["always_available"] != nil {
		model.AlwaysAvailable = core.BoolPtr(modelMap["always_available"].(bool))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToSupportDetailsItemAvailabilityTime(modelMap map[string]interface{}) (*partnercentersellv1.SupportDetailsItemAvailabilityTime, error) {
	model := &partnercentersellv1.SupportDetailsItemAvailabilityTime{}
	if modelMap["day"] != nil {
		model.Day = core.Float64Ptr(modelMap["day"].(float64))
	}
	if modelMap["start_time"] != nil && modelMap["start_time"].(string) != "" {
		model.StartTime = core.StringPtr(modelMap["start_time"].(string))
	}
	if modelMap["end_time"] != nil && modelMap["end_time"].(string) != "" {
		model.EndTime = core.StringPtr(modelMap["end_time"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherComposite(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadataOtherComposite, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadataOtherComposite{}
	if modelMap["composite_kind"] != nil && modelMap["composite_kind"].(string) != "" {
		model.CompositeKind = core.StringPtr(modelMap["composite_kind"].(string))
	}
	if modelMap["composite_tag"] != nil && modelMap["composite_tag"].(string) != "" {
		model.CompositeTag = core.StringPtr(modelMap["composite_tag"].(string))
	}
	if modelMap["children"] != nil {
		children := []partnercentersellv1.GlobalCatalogProductMetadataOtherCompositeChild{}
		for _, childrenItem := range modelMap["children"].([]interface{}) {
			childrenItemModel, err := ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherCompositeChild(childrenItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			children = append(children, *childrenItemModel)
		}
		model.Children = children
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductMapToGlobalCatalogProductMetadataOtherCompositeChild(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogProductMetadataOtherCompositeChild, error) {
	model := &partnercentersellv1.GlobalCatalogProductMetadataOtherCompositeChild{}
	if modelMap["kind"] != nil && modelMap["kind"].(string) != "" {
		model.Kind = core.StringPtr(modelMap["kind"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUIToMap(model *partnercentersellv1.GlobalCatalogOverviewUI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUITranslatedContentToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUITranslatedContentToMap(model *partnercentersellv1.GlobalCatalogOverviewUITranslatedContent) (map[string]interface{}, error) {
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

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductImagesToMap(model *partnercentersellv1.GlobalCatalogProductImages) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Image != nil {
		modelMap["image"] = *model.Image
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductCatalogProductProviderToMap(model *partnercentersellv1.CatalogProductProvider) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataToMap(model *partnercentersellv1.GlobalCatalogProductMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RcCompatible != nil {
		modelMap["rc_compatible"] = *model.RcCompatible
	}
	if model.Ui != nil {
		uiMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIToMap(model.Ui)
		if err != nil {
			return modelMap, err
		}
		modelMap["ui"] = []map[string]interface{}{uiMap}
	}
	if model.Service != nil {
		serviceMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataServiceToMap(model.Service)
		if err != nil {
			return modelMap, err
		}
		modelMap["service"] = []map[string]interface{}{serviceMap}
	}
	if model.Other != nil {
		otherMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherToMap(model.Other)
		if err != nil {
			return modelMap, err
		}
		modelMap["other"] = []map[string]interface{}{otherMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIToMap(model *partnercentersellv1.GlobalCatalogMetadataUI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Strings != nil {
		stringsMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsToMap(model.Strings)
		if err != nil {
			return modelMap, err
		}
		modelMap["strings"] = []map[string]interface{}{stringsMap}
	}
	if model.Urls != nil {
		urlsMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIUrlsToMap(model.Urls)
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

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsToMap(model *partnercentersellv1.GlobalCatalogMetadataUIStrings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsContentToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsContentToMap(model *partnercentersellv1.GlobalCatalogMetadataUIStringsContent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Bullets != nil {
		bullets := []map[string]interface{}{}
		for _, bulletsItem := range model.Bullets {
			bulletsItemMap, err := ResourceIbmOnboardingCatalogProductCatalogHighlightItemToMap(&bulletsItem) // #nosec G601
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
			mediaItemMap, err := ResourceIbmOnboardingCatalogProductCatalogProductMediaItemToMap(&mediaItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			media = append(media, mediaItemMap)
		}
		modelMap["media"] = media
	}
	if model.EmbeddableDashboard != nil {
		modelMap["embeddable_dashboard"] = *model.EmbeddableDashboard
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductCatalogHighlightItemToMap(model *partnercentersellv1.CatalogHighlightItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.DescriptionI18n != nil {
		descriptionI18n := make(map[string]interface{})
		for k, v := range model.DescriptionI18n {
			descriptionI18n[k] = flex.Stringify(v)
		}
		modelMap["description_i18n"] = descriptionI18n
	}
	if model.Title != nil {
		modelMap["title"] = *model.Title
	}
	if model.TitleI18n != nil {
		titleI18n := make(map[string]interface{})
		for k, v := range model.TitleI18n {
			titleI18n[k] = flex.Stringify(v)
		}
		modelMap["title_i18n"] = titleI18n
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductCatalogProductMediaItemToMap(model *partnercentersellv1.CatalogProductMediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["caption"] = *model.Caption
	if model.CaptionI18n != nil {
		captionI18n := make(map[string]interface{})
		for k, v := range model.CaptionI18n {
			captionI18n[k] = flex.Stringify(v)
		}
		modelMap["caption_i18n"] = captionI18n
	}
	if model.Thumbnail != nil {
		modelMap["thumbnail"] = *model.Thumbnail
	}
	modelMap["type"] = *model.Type
	modelMap["url"] = *model.URL
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIUrlsToMap(model *partnercentersellv1.GlobalCatalogMetadataUIUrls) (map[string]interface{}, error) {
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

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataServiceToMap(model *partnercentersellv1.GlobalCatalogMetadataService) (map[string]interface{}, error) {
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
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherToMap(model *partnercentersellv1.GlobalCatalogProductMetadataOther) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PC != nil {
		pcMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCToMap(model.PC)
		if err != nil {
			return modelMap, err
		}
		modelMap["pc"] = []map[string]interface{}{pcMap}
	}
	if model.Composite != nil {
		compositeMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeToMap(model.Composite)
		if err != nil {
			return modelMap, err
		}
		modelMap["composite"] = []map[string]interface{}{compositeMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCToMap(model *partnercentersellv1.GlobalCatalogProductMetadataOtherPC) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Support != nil {
		supportMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCSupportToMap(model.Support)
		if err != nil {
			return modelMap, err
		}
		modelMap["support"] = []map[string]interface{}{supportMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCSupportToMap(model *partnercentersellv1.GlobalCatalogProductMetadataOtherPCSupport) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.StatusURL != nil {
		modelMap["status_url"] = *model.StatusURL
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	if model.Languages != nil {
		modelMap["languages"] = model.Languages
	}
	if model.Process != nil {
		modelMap["process"] = *model.Process
	}
	if model.ProcessI18n != nil {
		processI18n := make(map[string]interface{})
		for k, v := range model.ProcessI18n {
			processI18n[k] = flex.Stringify(v)
		}
		modelMap["process_i18n"] = processI18n
	}
	if model.SupportType != nil {
		modelMap["support_type"] = *model.SupportType
	}
	if model.SupportEscalation != nil {
		supportEscalationMap, err := ResourceIbmOnboardingCatalogProductSupportEscalationToMap(model.SupportEscalation)
		if err != nil {
			return modelMap, err
		}
		modelMap["support_escalation"] = []map[string]interface{}{supportEscalationMap}
	}
	if model.SupportDetails != nil {
		supportDetails := []map[string]interface{}{}
		for _, supportDetailsItem := range model.SupportDetails {
			supportDetailsItemMap, err := ResourceIbmOnboardingCatalogProductSupportDetailsItemToMap(&supportDetailsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			supportDetails = append(supportDetails, supportDetailsItemMap)
		}
		modelMap["support_details"] = supportDetails
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductSupportEscalationToMap(model *partnercentersellv1.SupportEscalation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Contact != nil {
		modelMap["contact"] = *model.Contact
	}
	if model.EscalationWaitTime != nil {
		escalationWaitTimeMap, err := ResourceIbmOnboardingCatalogProductSupportTimeIntervalToMap(model.EscalationWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["escalation_wait_time"] = []map[string]interface{}{escalationWaitTimeMap}
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := ResourceIbmOnboardingCatalogProductSupportTimeIntervalToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductSupportTimeIntervalToMap(model *partnercentersellv1.SupportTimeInterval) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemToMap(model *partnercentersellv1.SupportDetailsItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Contact != nil {
		modelMap["contact"] = *model.Contact
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := ResourceIbmOnboardingCatalogProductSupportTimeIntervalToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	if model.Availability != nil {
		availabilityMap, err := ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityToMap(model.Availability)
		if err != nil {
			return modelMap, err
		}
		modelMap["availability"] = []map[string]interface{}{availabilityMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityToMap(model *partnercentersellv1.SupportDetailsItemAvailability) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Times != nil {
		times := []map[string]interface{}{}
		for _, timesItem := range model.Times {
			timesItemMap, err := ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityTimeToMap(&timesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			times = append(times, timesItemMap)
		}
		modelMap["times"] = times
	}
	if model.Timezone != nil {
		modelMap["timezone"] = *model.Timezone
	}
	if model.AlwaysAvailable != nil {
		modelMap["always_available"] = *model.AlwaysAvailable
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityTimeToMap(model *partnercentersellv1.SupportDetailsItemAvailabilityTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Day != nil {
		modelMap["day"] = *model.Day
	}
	if model.StartTime != nil {
		modelMap["start_time"] = *model.StartTime
	}
	if model.EndTime != nil {
		modelMap["end_time"] = *model.EndTime
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeToMap(model *partnercentersellv1.GlobalCatalogProductMetadataOtherComposite) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CompositeKind != nil {
		modelMap["composite_kind"] = *model.CompositeKind
	}
	if model.CompositeTag != nil {
		modelMap["composite_tag"] = *model.CompositeTag
	}
	if model.Children != nil {
		children := []map[string]interface{}{}
		for _, childrenItem := range model.Children {
			childrenItemMap, err := ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeChildToMap(&childrenItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			children = append(children, childrenItemMap)
		}
		modelMap["children"] = children
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeChildToMap(model *partnercentersellv1.GlobalCatalogProductMetadataOtherCompositeChild) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Kind != nil {
		modelMap["kind"] = *model.Kind
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductPatchAsPatch(patchVals *partnercentersellv1.GlobalCatalogProductPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "active"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["active"] = nil
	}
	path = "disabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["disabled"] = nil
	}
	path = "overview_ui"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["overview_ui"] = nil
	} else if exists && patch["overview_ui"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUIAsPatch(patch["overview_ui"].(map[string]interface{}), d)
	}
	path = "tags"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tags"] = nil
	}
	path = "images"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["images"] = nil
	} else if exists && patch["images"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductImagesAsPatch(patch["images"].(map[string]interface{}), d)
	}
	path = "object_provider"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["object_provider"] = nil
	} else if exists && patch["object_provider"] != nil {
		ResourceIbmOnboardingCatalogProductCatalogProductProviderAsPatch(patch["object_provider"].(map[string]interface{}), d)
	}
	path = "metadata"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["metadata"] = nil
	} else if exists && patch["metadata"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataAsPatch(patch["metadata"].(map[string]interface{}), d)
	}

	return patch
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.rc_compatible"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["rc_compatible"] = nil
	}
	path = "metadata.0.ui"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ui"] = nil
	} else if exists && patch["ui"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIAsPatch(patch["ui"].(map[string]interface{}), d)
	}
	path = "metadata.0.service"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service"] = nil
	} else if exists && patch["service"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataServiceAsPatch(patch["service"].(map[string]interface{}), d)
	}
	path = "metadata.0.other"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["other"] = nil
	} else if exists && patch["other"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherAsPatch(patch["other"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["pc"] = nil
	} else if exists && patch["pc"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCAsPatch(patch["pc"].(map[string]interface{}), d)
	}
	path = "metadata.0.other.0.composite"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["composite"] = nil
	} else if exists && patch["composite"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeAsPatch(patch["composite"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.composite.0.composite_kind"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["composite_kind"] = nil
	}
	path = "metadata.0.other.0.composite.0.composite_tag"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["composite_tag"] = nil
	}
	path = "metadata.0.other.0.composite.0.children"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["children"] = nil
	} else if exists && patch["children"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeChildAsPatch(patch["children"].([]interface{})[0].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherCompositeChildAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.composite.0.children.0.kind"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["kind"] = nil
	}
	path = "metadata.0.other.0.composite.0.children.0.name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["support"] = nil
	} else if exists && patch["support"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCSupportAsPatch(patch["support"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductMetadataOtherPCSupportAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["url"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.status_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["status_url"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.locations"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["locations"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.languages"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["languages"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.process"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["process"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.process_i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["process_i18n"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["support_type"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_escalation"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["support_escalation"] = nil
	} else if exists && patch["support_escalation"] != nil {
		ResourceIbmOnboardingCatalogProductSupportEscalationAsPatch(patch["support_escalation"].(map[string]interface{}), d)
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["support_details"] = nil
	} else if exists && patch["support_details"] != nil {
		ResourceIbmOnboardingCatalogProductSupportDetailsItemAsPatch(patch["support_details"].([]interface{})[0].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.support_details.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.contact"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["contact"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.response_wait_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["response_wait_time"] = nil
	} else if exists && patch["response_wait_time"] != nil {
		ResourceIbmOnboardingCatalogProductSupportTimeIntervalAsPatch(patch["response_wait_time"].(map[string]interface{}), d)
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["availability"] = nil
	} else if exists && patch["availability"] != nil {
		ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityAsPatch(patch["availability"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.times"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["times"] = nil
	} else if exists && patch["times"] != nil {
		ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityTimeAsPatch(patch["times"].([]interface{})[0].(map[string]interface{}), d)
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.timezone"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["timezone"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.always_available"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["always_available"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductSupportDetailsItemAvailabilityTimeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.times.0.day"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["day"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.times.0.start_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["start_time"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.availability.0.times.0.end_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["end_time"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductSupportEscalationAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.support_escalation.0.contact"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["contact"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_escalation.0.escalation_wait_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["escalation_wait_time"] = nil
	} else if exists && patch["escalation_wait_time"] != nil {
		ResourceIbmOnboardingCatalogProductSupportTimeIntervalAsPatch(patch["escalation_wait_time"].(map[string]interface{}), d)
	}
	path = "metadata.0.other.0.pc.0.support.0.support_escalation.0.response_wait_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["response_wait_time"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductSupportTimeIntervalAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.other.0.pc.0.support.0.support_details.0.response_wait_time.0.value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	}
	path = "metadata.0.other.0.pc.0.support.0.support_details.0.response_wait_time.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataServiceAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.service.0.rc_provisionable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["rc_provisionable"] = nil
	}
	path = "metadata.0.service.0.iam_compatible"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["iam_compatible"] = nil
	}
	path = "metadata.0.service.0.bindable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["bindable"] = nil
	}
	path = "metadata.0.service.0.plan_updateable"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["plan_updateable"] = nil
	}
	path = "metadata.0.service.0.service_key_supported"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service_key_supported"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.strings"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["strings"] = nil
	} else if exists && patch["strings"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsAsPatch(patch["strings"].(map[string]interface{}), d)
	}
	path = "metadata.0.ui.0.urls"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["urls"] = nil
	} else if exists && patch["urls"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIUrlsAsPatch(patch["urls"].(map[string]interface{}), d)
	}
	path = "metadata.0.ui.0.hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	}
	path = "metadata.0.ui.0.side_by_side_index"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["side_by_side_index"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIUrlsAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.urls.0.doc_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["doc_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.apidocs_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["apidocs_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.terms_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["terms_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.instructions_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["instructions_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.catalog_details_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["catalog_details_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.custom_create_page_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["custom_create_page_url"] = nil
	}
	path = "metadata.0.ui.0.urls.0.dashboard"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["dashboard"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.strings.0.en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsContentAsPatch(patch["en"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogMetadataUIStringsContentAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.strings.0.en.0.bullets"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["bullets"] = nil
	} else if exists && patch["bullets"] != nil {
		ResourceIbmOnboardingCatalogProductCatalogHighlightItemAsPatch(patch["bullets"].([]interface{})[0].(map[string]interface{}), d)
	}
	path = "metadata.0.ui.0.strings.0.en.0.media"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["media"] = nil
	} else if exists && patch["media"] != nil {
		ResourceIbmOnboardingCatalogProductCatalogProductMediaItemAsPatch(patch["media"].([]interface{})[0].(map[string]interface{}), d)
	}
	path = "metadata.0.ui.0.strings.0.en.0.embeddable_dashboard"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["embeddable_dashboard"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductCatalogProductMediaItemAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.strings.0.en.0.media.0.caption_i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["caption_i18n"] = nil
	}
	path = "metadata.0.ui.0.strings.0.en.0.media.0.thumbnail"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["thumbnail"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductCatalogHighlightItemAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "metadata.0.ui.0.strings.0.en.0.bullets.0.description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	}
	path = "metadata.0.ui.0.strings.0.en.0.bullets.0.description_i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description_i18n"] = nil
	}
	path = "metadata.0.ui.0.strings.0.en.0.bullets.0.title"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["title"] = nil
	}
	path = "metadata.0.ui.0.strings.0.en.0.bullets.0.title_i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["title_i18n"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductCatalogProductProviderAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "object_provider.0.name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
	path = "object_provider.0.email"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["email"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogProductImagesAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "images.0.image"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image"] = nil
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUIAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "overview_ui.0.en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUITranslatedContentAsPatch(patch["en"].(map[string]interface{}), d)
	}
}

func ResourceIbmOnboardingCatalogProductGlobalCatalogOverviewUITranslatedContentAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "overview_ui.0.en.0.display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	}
	path = "overview_ui.0.en.0.description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	}
	path = "overview_ui.0.en.0.long_description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["long_description"] = nil
	}
}
