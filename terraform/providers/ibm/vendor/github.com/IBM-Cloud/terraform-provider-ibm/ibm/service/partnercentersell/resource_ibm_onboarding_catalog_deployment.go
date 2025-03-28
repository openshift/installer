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

func ResourceIbmOnboardingCatalogDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingCatalogDeploymentCreate,
		ReadContext:   resourceIbmOnboardingCatalogDeploymentRead,
		UpdateContext: resourceIbmOnboardingCatalogDeploymentUpdate,
		DeleteContext: resourceIbmOnboardingCatalogDeploymentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "product_id"),
				Description:  "The unique ID of the product.",
			},
			"catalog_product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "catalog_product_id"),
				Description:  "The unique ID of this global catalog product.",
			},
			"catalog_plan_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "catalog_plan_id"),
				Description:  "The unique ID of this global catalog plan.",
			},
			"env": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "env"),
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
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "name"),
				Description:  "The programmatic name of this deployment.",
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
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_catalog_deployment", "kind"),
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
				Description: "Global catalog deployment metadata.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rc_compatible": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the object is compatible with the resource controller service.",
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
										Computed:    true,
										Description: "Deprecated. Controls the Connections tab on the Resource Details page.",
									},
									"plan_updateable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
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
									"parameters": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"displayname": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The display name for custom service parameters.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key of the parameter.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The type of custom service parameters.",
												},
												"options": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"displayname": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The display name for custom service parameters.",
															},
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The value for custom service parameters.",
															},
															"i18n": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The description for the object.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"en": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"de": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"es": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"fr": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"it": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"ja": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"ko": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"pt_br": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"zh_tw": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
																					},
																				},
																			},
																		},
																		"zh_cn": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name and description.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"displayname": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter display name.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The translations for custom service parameter description.",
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
												"value": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"layout": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the layout of check box or radio input types. When unspecified, the default layout is horizontal.",
												},
												"associations": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "A JSON structure to describe the interactions with pricing plans and/or other custom parameters.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"validation_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The validation URL for custom service parameters.",
												},
												"options_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The options URL for custom service parameters.",
												},
												"invalidmessage": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The message that appears when the content of the text box is invalid.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The description of the parameter that is displayed to help users with the value of the parameter.",
												},
												"required": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "A boolean value that indicates whether the parameter must be entered in the IBM Cloud user interface.",
												},
												"pattern": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "A regular expression that the value is checked against.",
												},
												"placeholder": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The placeholder text for custom parameters.",
												},
												"readonly": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "A boolean value that indicates whether the value of the parameter is displayed only and cannot be changed by users. The default value is false.",
												},
												"hidden": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Indicates whether the custom parameters is hidden required or not.",
												},
												"i18n": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The description for the object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"en": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"de": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"es": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"fr": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"it": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"ja": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"ko": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"pt_br": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"zh_tw": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
																		},
																	},
																},
															},
															"zh_cn": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The translations for custom service parameter display name and description.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"displayname": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter display name.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The translations for custom service parameter description.",
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
						"deployment": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The global catalog metadata of the deployment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"broker": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The global catalog metadata of the deployment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the resource broker.",
												},
												"guid": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Crn or guid of the resource broker.",
												},
											},
										},
									},
									"location": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The global catalog deployment location.",
									},
									"location_url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The global catalog deployment URL of location.",
									},
									"target_crn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region crn.",
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
			"catalog_deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a global catalog object.",
			},
		},
	}
}

func ResourceIbmOnboardingCatalogDeploymentValidator() *validate.ResourceValidator {
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
			Identifier:                 "catalog_plan_id",
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
			AllowedValues:              "deployment",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_catalog_deployment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingCatalogDeploymentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createCatalogDeploymentOptions := &partnercentersellv1.CreateCatalogDeploymentOptions{}

	createCatalogDeploymentOptions.SetProductID(d.Get("product_id").(string))
	createCatalogDeploymentOptions.SetCatalogProductID(d.Get("catalog_product_id").(string))
	createCatalogDeploymentOptions.SetCatalogPlanID(d.Get("catalog_plan_id").(string))
	createCatalogDeploymentOptions.SetName(d.Get("name").(string))
	createCatalogDeploymentOptions.SetActive(d.Get("active").(bool))
	createCatalogDeploymentOptions.SetDisabled(d.Get("disabled").(bool))
	createCatalogDeploymentOptions.SetKind(d.Get("kind").(string))
	if _, ok := d.GetOk("env"); ok {
		createCatalogDeploymentOptions.SetEnv(d.Get("env").(string))
	}
	var tags []string
	for _, v := range d.Get("tags").([]interface{}) {
		tagsItem := v.(string)
		tags = append(tags, tagsItem)
	}
	createCatalogDeploymentOptions.SetTags(tags)
	objectProviderModel, err := ResourceIbmOnboardingCatalogDeploymentMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "create", "parse-object_provider").GetDiag()
	}
	createCatalogDeploymentOptions.SetObjectProvider(objectProviderModel)
	if _, ok := d.GetOk("object_id"); ok {
		createCatalogDeploymentOptions.SetObjectID(d.Get("object_id").(string))
	}
	if _, ok := d.GetOk("overview_ui"); ok {
		overviewUiModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "create", "parse-overview_ui").GetDiag()
		}
		createCatalogDeploymentOptions.SetOverviewUi(overviewUiModel)
	}
	if _, ok := d.GetOk("metadata"); ok {
		metadataModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogDeploymentMetadataPrototypePatch(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "create", "parse-metadata").GetDiag()
		}
		createCatalogDeploymentOptions.SetMetadata(metadataModel)
	}
	if _, ok := d.GetOk("env"); ok {
		createCatalogDeploymentOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogDeployment, _, err := partnerCenterSellClient.CreateCatalogDeploymentWithContext(context, createCatalogDeploymentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateCatalogDeploymentWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_deployment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", *createCatalogDeploymentOptions.ProductID, *createCatalogDeploymentOptions.CatalogProductID, *createCatalogDeploymentOptions.CatalogPlanID, *globalCatalogDeployment.ID))

	return resourceIbmOnboardingCatalogDeploymentRead(context, d, meta)
}

func resourceIbmOnboardingCatalogDeploymentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getCatalogDeploymentOptions := &partnercentersellv1.GetCatalogDeploymentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "sep-id-parts").GetDiag()
	}

	getCatalogDeploymentOptions.SetProductID(parts[0])
	getCatalogDeploymentOptions.SetCatalogProductID(parts[1])
	getCatalogDeploymentOptions.SetCatalogPlanID(parts[2])
	getCatalogDeploymentOptions.SetCatalogDeploymentID(parts[3])
	if _, ok := d.GetOk("env"); ok {
		getCatalogDeploymentOptions.SetEnv(d.Get("env").(string))
	}

	globalCatalogDeployment, response, err := partnerCenterSellClient.GetCatalogDeploymentWithContext(context, getCatalogDeploymentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetCatalogDeploymentWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_deployment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(globalCatalogDeployment.ObjectID) {
		if err = d.Set("object_id", globalCatalogDeployment.ObjectID); err != nil {
			err = fmt.Errorf("Error setting object_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-object_id").GetDiag()
		}
	}
	if err = d.Set("name", globalCatalogDeployment.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-name").GetDiag()
	}
	if err = d.Set("active", globalCatalogDeployment.Active); err != nil {
		err = fmt.Errorf("Error setting active: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-active").GetDiag()
	}
	if err = d.Set("disabled", globalCatalogDeployment.Disabled); err != nil {
		err = fmt.Errorf("Error setting disabled: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-disabled").GetDiag()
	}
	if err = d.Set("kind", globalCatalogDeployment.Kind); err != nil {
		err = fmt.Errorf("Error setting kind: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-kind").GetDiag()
	}
	if !core.IsNil(globalCatalogDeployment.OverviewUi) {
		overviewUiMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUIToMap(globalCatalogDeployment.OverviewUi)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "overview_ui-to-map").GetDiag()
		}
		if err = d.Set("overview_ui", []map[string]interface{}{overviewUiMap}); err != nil {
			err = fmt.Errorf("Error setting overview_ui: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-overview_ui").GetDiag()
		}
	}
	if err = d.Set("tags", globalCatalogDeployment.Tags); err != nil {
		err = fmt.Errorf("Error setting tags: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-tags").GetDiag()
	}
	objectProviderMap, err := ResourceIbmOnboardingCatalogDeploymentCatalogProductProviderToMap(globalCatalogDeployment.ObjectProvider)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "object_provider-to-map").GetDiag()
	}
	if err = d.Set("object_provider", []map[string]interface{}{objectProviderMap}); err != nil {
		err = fmt.Errorf("Error setting object_provider: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-object_provider").GetDiag()
	}
	if !core.IsNil(globalCatalogDeployment.Metadata) {
		metadataMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataToMap(globalCatalogDeployment.Metadata)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "metadata-to-map").GetDiag()
		}
		if err = d.Set("metadata", []map[string]interface{}{metadataMap}); err != nil {
			err = fmt.Errorf("Error setting metadata: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-metadata").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogDeployment.GeoTags) {
		if err = d.Set("geo_tags", globalCatalogDeployment.GeoTags); err != nil {
			err = fmt.Errorf("Error setting geo_tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-geo_tags").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogDeployment.URL) {
		if err = d.Set("url", globalCatalogDeployment.URL); err != nil {
			err = fmt.Errorf("Error setting url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-url").GetDiag()
		}
	}
	if parts[0] != "" {
		if err = d.Set("product_id", parts[0]); err != nil {
			err = fmt.Errorf("Error setting product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-product_id").GetDiag()
		}
	}
	if parts[1] != "" {
		if err = d.Set("catalog_product_id", parts[1]); err != nil {
			err = fmt.Errorf("Error setting catalog_product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-catalog_product_id").GetDiag()
		}
	}
	if parts[2] != "" {
		if err = d.Set("catalog_plan_id", parts[2]); err != nil {
			err = fmt.Errorf("Error setting catalog_plan_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-catalog_plan_id").GetDiag()
		}
	}
	if !core.IsNil(globalCatalogDeployment.ID) {
		if err = d.Set("catalog_deployment_id", globalCatalogDeployment.ID); err != nil {
			err = fmt.Errorf("Error setting catalog_deployment_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "read", "set-catalog_deployment_id").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingCatalogDeploymentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateCatalogDeploymentOptions := &partnercentersellv1.UpdateCatalogDeploymentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "update", "sep-id-parts").GetDiag()
	}

	updateCatalogDeploymentOptions.SetProductID(parts[0])
	updateCatalogDeploymentOptions.SetCatalogProductID(parts[1])
	updateCatalogDeploymentOptions.SetCatalogPlanID(parts[2])
	updateCatalogDeploymentOptions.SetCatalogDeploymentID(parts[3])
	if _, ok := d.GetOk("env"); ok {
		updateCatalogDeploymentOptions.SetEnv(d.Get("env").(string))
	}

	hasChange := false

	patchVals := &partnercentersellv1.GlobalCatalogDeploymentPatch{}
	if d.HasChange("product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_deployment", "update", "product_id-forces-new").GetDiag()
	}
	if d.HasChange("catalog_product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "catalog_product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_deployment", "update", "catalog_product_id-forces-new").GetDiag()
	}
	if d.HasChange("catalog_plan_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "catalog_plan_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_catalog_deployment", "update", "catalog_plan_id-forces-new").GetDiag()
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
		overviewUi, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogOverviewUI(d.Get("overview_ui.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "update", "parse-overview_ui").GetDiag()
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
	if d.HasChange("object_provider") {
		objectProvider, err := ResourceIbmOnboardingCatalogDeploymentMapToCatalogProductProvider(d.Get("object_provider.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "update", "parse-object_provider").GetDiag()
		}
		patchVals.ObjectProvider = objectProvider
		hasChange = true
	}
	if d.HasChange("metadata") {
		metadata, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogDeploymentMetadataPrototypePatch(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "update", "parse-metadata").GetDiag()
		}
		patchVals.Metadata = metadata
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateCatalogDeploymentOptions.GlobalCatalogDeploymentPatch = ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateCatalogDeploymentWithContext(context, updateCatalogDeploymentOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateCatalogDeploymentWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_deployment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingCatalogDeploymentRead(context, d, meta)
}

func resourceIbmOnboardingCatalogDeploymentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteCatalogDeploymentOptions := &partnercentersellv1.DeleteCatalogDeploymentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_catalog_deployment", "delete", "sep-id-parts").GetDiag()
	}

	deleteCatalogDeploymentOptions.SetProductID(parts[0])
	deleteCatalogDeploymentOptions.SetCatalogProductID(parts[1])
	deleteCatalogDeploymentOptions.SetCatalogPlanID(parts[2])
	deleteCatalogDeploymentOptions.SetCatalogDeploymentID(parts[3])
	if _, ok := d.GetOk("env"); ok {
		deleteCatalogDeploymentOptions.SetEnv(d.Get("env").(string))
	}

	_, err = partnerCenterSellClient.DeleteCatalogDeploymentWithContext(context, deleteCatalogDeploymentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteCatalogDeploymentWithContext failed: %s", err.Error()), "ibm_onboarding_catalog_deployment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToCatalogProductProvider(modelMap map[string]interface{}) (*partnercentersellv1.CatalogProductProvider, error) {
	model := &partnercentersellv1.CatalogProductProvider{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["email"] != nil && modelMap["email"].(string) != "" {
		model.Email = core.StringPtr(modelMap["email"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogOverviewUI(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUI, error) {
	model := &partnercentersellv1.GlobalCatalogOverviewUI{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 && modelMap["en"].([]interface{})[0] != nil {
		EnModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogOverviewUITranslatedContent(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogOverviewUITranslatedContent(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogOverviewUITranslatedContent, error) {
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

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogDeploymentMetadataPrototypePatch(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogDeploymentMetadataPrototypePatch, error) {
	model := &partnercentersellv1.GlobalCatalogDeploymentMetadataPrototypePatch{}
	if modelMap["rc_compatible"] != nil {
		model.RcCompatible = core.BoolPtr(modelMap["rc_compatible"].(bool))
	}
	if modelMap["service"] != nil && len(modelMap["service"].([]interface{})) > 0 && modelMap["service"].([]interface{})[0] != nil {
		ServiceModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogDeploymentMetadataServicePrototypePatch(modelMap["service"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Service = ServiceModel
	}
	if modelMap["deployment"] != nil && len(modelMap["deployment"].([]interface{})) > 0 && modelMap["deployment"].([]interface{})[0] != nil {
		DeploymentModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataDeployment(modelMap["deployment"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Deployment = DeploymentModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogDeploymentMetadataServicePrototypePatch(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogDeploymentMetadataServicePrototypePatch, error) {
	model := &partnercentersellv1.GlobalCatalogDeploymentMetadataServicePrototypePatch{}
	if modelMap["rc_provisionable"] != nil {
		model.RcProvisionable = core.BoolPtr(modelMap["rc_provisionable"].(bool))
	}
	if modelMap["iam_compatible"] != nil {
		model.IamCompatible = core.BoolPtr(modelMap["iam_compatible"].(bool))
	}
	if modelMap["service_key_supported"] != nil {
		model.ServiceKeySupported = core.BoolPtr(modelMap["service_key_supported"].(bool))
	}
	if modelMap["parameters"] != nil {
		parameters := []partnercentersellv1.GlobalCatalogMetadataServiceCustomParameters{}
		for _, parametersItem := range modelMap["parameters"].([]interface{}) {
			parametersItemModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParameters(parametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			parameters = append(parameters, *parametersItemModel)
		}
		model.Parameters = parameters
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParameters(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataServiceCustomParameters, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataServiceCustomParameters{}
	if modelMap["displayname"] != nil && modelMap["displayname"].(string) != "" {
		model.Displayname = core.StringPtr(modelMap["displayname"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["options"] != nil {
		options := []partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersOptions{}
		for _, optionsItem := range modelMap["options"].([]interface{}) {
			optionsItemModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersOptions(optionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			options = append(options, *optionsItemModel)
		}
		model.Options = options
	}
	if modelMap["value"] != nil {
		value := []string{}
		for _, valueItem := range modelMap["value"].([]interface{}) {
			value = append(value, valueItem.(string))
		}
		model.Value = value
	}
	if modelMap["layout"] != nil && modelMap["layout"].(string) != "" {
		model.Layout = core.StringPtr(modelMap["layout"].(string))
	}
	if modelMap["associations"] != nil {
		model.Associations = modelMap["associations"].(map[string]interface{})
	}
	if modelMap["validation_url"] != nil && modelMap["validation_url"].(string) != "" {
		model.ValidationURL = core.StringPtr(modelMap["validation_url"].(string))
	}
	if modelMap["options_url"] != nil && modelMap["options_url"].(string) != "" {
		model.OptionsURL = core.StringPtr(modelMap["options_url"].(string))
	}
	if modelMap["invalidmessage"] != nil && modelMap["invalidmessage"].(string) != "" {
		model.Invalidmessage = core.StringPtr(modelMap["invalidmessage"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["required"] != nil {
		model.Required = core.BoolPtr(modelMap["required"].(bool))
	}
	if modelMap["pattern"] != nil && modelMap["pattern"].(string) != "" {
		model.Pattern = core.StringPtr(modelMap["pattern"].(string))
	}
	if modelMap["placeholder"] != nil && modelMap["placeholder"].(string) != "" {
		model.Placeholder = core.StringPtr(modelMap["placeholder"].(string))
	}
	if modelMap["readonly"] != nil {
		model.Readonly = core.BoolPtr(modelMap["readonly"].(bool))
	}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	if modelMap["i18n"] != nil && len(modelMap["i18n"].([]interface{})) > 0 {
		I18nModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18n(modelMap["i18n"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.I18n = I18nModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersOptions(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersOptions, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersOptions{}
	if modelMap["displayname"] != nil && modelMap["displayname"].(string) != "" {
		model.Displayname = core.StringPtr(modelMap["displayname"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["i18n"] != nil && len(modelMap["i18n"].([]interface{})) > 0 {
		I18nModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18n(modelMap["i18n"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.I18n = I18nModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18n(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18n, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18n{}
	if modelMap["en"] != nil && len(modelMap["en"].([]interface{})) > 0 {
		EnModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["en"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.En = EnModel
	}
	if modelMap["de"] != nil && len(modelMap["de"].([]interface{})) > 0 {
		DeModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["de"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.De = DeModel
	}
	if modelMap["es"] != nil && len(modelMap["es"].([]interface{})) > 0 {
		EsModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["es"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Es = EsModel
	}
	if modelMap["fr"] != nil && len(modelMap["fr"].([]interface{})) > 0 {
		FrModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["fr"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Fr = FrModel
	}
	if modelMap["it"] != nil && len(modelMap["it"].([]interface{})) > 0 {
		ItModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["it"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.It = ItModel
	}
	if modelMap["ja"] != nil && len(modelMap["ja"].([]interface{})) > 0 {
		JaModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["ja"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ja = JaModel
	}
	if modelMap["ko"] != nil && len(modelMap["ko"].([]interface{})) > 0 {
		KoModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["ko"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ko = KoModel
	}
	if modelMap["pt_br"] != nil && len(modelMap["pt_br"].([]interface{})) > 0 {
		PtBrModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["pt_br"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PtBr = PtBrModel
	}
	if modelMap["zh_tw"] != nil && len(modelMap["zh_tw"].([]interface{})) > 0 {
		ZhTwModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["zh_tw"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ZhTw = ZhTwModel
	}
	if modelMap["zh_cn"] != nil && len(modelMap["zh_cn"].([]interface{})) > 0 {
		ZhCnModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap["zh_cn"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ZhCn = ZhCnModel
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataServiceCustomParametersI18nFields(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18nFields, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18nFields{}
	if modelMap["displayname"] != nil && modelMap["displayname"].(string) != "" {
		model.Displayname = core.StringPtr(modelMap["displayname"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataDeployment(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataDeployment, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataDeployment{}
	if modelMap["broker"] != nil && len(modelMap["broker"].([]interface{})) > 0 && modelMap["broker"].([]interface{})[0] != nil {
		BrokerModel, err := ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataDeploymentBroker(modelMap["broker"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Broker = BrokerModel
	}
	if modelMap["location"] != nil && modelMap["location"].(string) != "" {
		model.Location = core.StringPtr(modelMap["location"].(string))
	}
	if modelMap["location_url"] != nil && modelMap["location_url"].(string) != "" {
		model.LocationURL = core.StringPtr(modelMap["location_url"].(string))
	}
	if modelMap["target_crn"] != nil && modelMap["target_crn"].(string) != "" {
		model.TargetCrn = core.StringPtr(modelMap["target_crn"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentMapToGlobalCatalogMetadataDeploymentBroker(modelMap map[string]interface{}) (*partnercentersellv1.GlobalCatalogMetadataDeploymentBroker, error) {
	model := &partnercentersellv1.GlobalCatalogMetadataDeploymentBroker{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["guid"] != nil && modelMap["guid"].(string) != "" {
		model.Guid = core.StringPtr(modelMap["guid"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUIToMap(model *partnercentersellv1.GlobalCatalogOverviewUI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUITranslatedContentToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUITranslatedContentToMap(model *partnercentersellv1.GlobalCatalogOverviewUITranslatedContent) (map[string]interface{}, error) {
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

func ResourceIbmOnboardingCatalogDeploymentCatalogProductProviderToMap(model *partnercentersellv1.CatalogProductProvider) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataToMap(model *partnercentersellv1.GlobalCatalogDeploymentMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RcCompatible != nil {
		modelMap["rc_compatible"] = *model.RcCompatible
	}
	if model.Service != nil {
		serviceMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataServiceToMap(model.Service)
		if err != nil {
			return modelMap, err
		}
		modelMap["service"] = []map[string]interface{}{serviceMap}
	}
	if model.Deployment != nil {
		deploymentMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentToMap(model.Deployment)
		if err != nil {
			return modelMap, err
		}
		modelMap["deployment"] = []map[string]interface{}{deploymentMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataServiceToMap(model *partnercentersellv1.GlobalCatalogDeploymentMetadataService) (map[string]interface{}, error) {
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
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersToMap(&parametersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersToMap(model *partnercentersellv1.GlobalCatalogMetadataServiceCustomParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Displayname != nil {
		modelMap["displayname"] = *model.Displayname
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Options != nil {
		options := []map[string]interface{}{}
		for _, optionsItem := range model.Options {
			optionsItemMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersOptionsToMap(&optionsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			options = append(options, optionsItemMap)
		}
		modelMap["options"] = options
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Layout != nil {
		modelMap["layout"] = *model.Layout
	}
	if model.Associations != nil {
		associations := make(map[string]interface{})
		for k, v := range model.Associations {
			associations[k] = flex.Stringify(v)
		}
		modelMap["associations"] = associations
	}
	if model.ValidationURL != nil {
		modelMap["validation_url"] = *model.ValidationURL
	}
	if model.OptionsURL != nil {
		modelMap["options_url"] = *model.OptionsURL
	}
	if model.Invalidmessage != nil {
		modelMap["invalidmessage"] = *model.Invalidmessage
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Required != nil {
		modelMap["required"] = *model.Required
	}
	if model.Pattern != nil {
		modelMap["pattern"] = *model.Pattern
	}
	if model.Placeholder != nil {
		modelMap["placeholder"] = *model.Placeholder
	}
	if model.Readonly != nil {
		modelMap["readonly"] = *model.Readonly
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.I18n != nil {
		i18nMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nToMap(model.I18n)
		if err != nil {
			return modelMap, err
		}
		modelMap["i18n"] = []map[string]interface{}{i18nMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersOptionsToMap(model *partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Displayname != nil {
		modelMap["displayname"] = *model.Displayname
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.I18n != nil {
		i18nMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nToMap(model.I18n)
		if err != nil {
			return modelMap, err
		}
		modelMap["i18n"] = []map[string]interface{}{i18nMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nToMap(model *partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18n) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.En != nil {
		enMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.En)
		if err != nil {
			return modelMap, err
		}
		modelMap["en"] = []map[string]interface{}{enMap}
	}
	if model.De != nil {
		deMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.De)
		if err != nil {
			return modelMap, err
		}
		modelMap["de"] = []map[string]interface{}{deMap}
	}
	if model.Es != nil {
		esMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.Es)
		if err != nil {
			return modelMap, err
		}
		modelMap["es"] = []map[string]interface{}{esMap}
	}
	if model.Fr != nil {
		frMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.Fr)
		if err != nil {
			return modelMap, err
		}
		modelMap["fr"] = []map[string]interface{}{frMap}
	}
	if model.It != nil {
		itMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.It)
		if err != nil {
			return modelMap, err
		}
		modelMap["it"] = []map[string]interface{}{itMap}
	}
	if model.Ja != nil {
		jaMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.Ja)
		if err != nil {
			return modelMap, err
		}
		modelMap["ja"] = []map[string]interface{}{jaMap}
	}
	if model.Ko != nil {
		koMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.Ko)
		if err != nil {
			return modelMap, err
		}
		modelMap["ko"] = []map[string]interface{}{koMap}
	}
	if model.PtBr != nil {
		ptBrMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.PtBr)
		if err != nil {
			return modelMap, err
		}
		modelMap["pt_br"] = []map[string]interface{}{ptBrMap}
	}
	if model.ZhTw != nil {
		zhTwMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.ZhTw)
		if err != nil {
			return modelMap, err
		}
		modelMap["zh_tw"] = []map[string]interface{}{zhTwMap}
	}
	if model.ZhCn != nil {
		zhCnMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model.ZhCn)
		if err != nil {
			return modelMap, err
		}
		modelMap["zh_cn"] = []map[string]interface{}{zhCnMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsToMap(model *partnercentersellv1.GlobalCatalogMetadataServiceCustomParametersI18nFields) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Displayname != nil {
		modelMap["displayname"] = *model.Displayname
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentToMap(model *partnercentersellv1.GlobalCatalogMetadataDeployment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Broker != nil {
		brokerMap, err := ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentBrokerToMap(model.Broker)
		if err != nil {
			return modelMap, err
		}
		modelMap["broker"] = []map[string]interface{}{brokerMap}
	}
	if model.Location != nil {
		modelMap["location"] = *model.Location
	}
	if model.LocationURL != nil {
		modelMap["location_url"] = *model.LocationURL
	}
	if model.TargetCrn != nil {
		modelMap["target_crn"] = *model.TargetCrn
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentBrokerToMap(model *partnercentersellv1.GlobalCatalogMetadataDeploymentBroker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Guid != nil {
		modelMap["guid"] = *model.Guid
	}
	return modelMap, nil
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentPatchAsPatch(patchVals *partnercentersellv1.GlobalCatalogDeploymentPatch, d *schema.ResourceData) map[string]interface{} {
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
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUIAsPatch(patch["overview_ui"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "overview_ui")
	}
	path = "tags"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tags"] = nil
	} else if !exists {
		delete(patch, "tags")
	}
	path = "object_provider"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["object_provider"] = nil
	} else if exists && patch["object_provider"] != nil {
		ResourceIbmOnboardingCatalogDeploymentCatalogProductProviderAsPatch(patch["object_provider"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "object_provider")
	}
	path = "metadata"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["metadata"] = nil
	} else if exists && patch["metadata"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataPrototypePatchAsPatch(patch["metadata"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "metadata")
	}

	return patch
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataPrototypePatchAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".rc_compatible"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["rc_compatible"] = nil
	} else if !exists {
		delete(patch, "rc_compatible")
	}
	path = rootPath + ".service"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service"] = nil
	} else if exists && patch["service"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataServicePrototypePatchAsPatch(patch["service"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "service")
	}
	path = rootPath + ".deployment"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["deployment"] = nil
	} else if exists && patch["deployment"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentAsPatch(patch["deployment"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "deployment")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".broker"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["broker"] = nil
	} else if exists && patch["broker"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentBrokerAsPatch(patch["broker"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "broker")
	}
	path = rootPath + ".location"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["location"] = nil
	} else if !exists {
		delete(patch, "location")
	}
	path = rootPath + ".location_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["location_url"] = nil
	} else if !exists {
		delete(patch, "location_url")
	}
	path = rootPath + ".target_crn"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["target_crn"] = nil
	} else if !exists {
		delete(patch, "target_crn")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataDeploymentBrokerAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = rootPath + ".guid"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["guid"] = nil
	} else if !exists {
		delete(patch, "guid")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogDeploymentMetadataServicePrototypePatchAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
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
	path = rootPath + ".service_key_supported"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service_key_supported"] = nil
	} else if !exists {
		delete(patch, "service_key_supported")
	}
	path = rootPath + ".parameters"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["parameters"] = nil
	} else if exists && patch["parameters"] != nil {
		parametersList := patch["parameters"].([]map[string]interface{})
		for i, parametersItem := range parametersList {
			ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersAsPatch(parametersItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "parameters")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".displayname"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["displayname"] = nil
	} else if !exists {
		delete(patch, "displayname")
	}
	path = rootPath + ".name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = rootPath + ".type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	} else if !exists {
		delete(patch, "type")
	}
	path = rootPath + ".options"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options"] = nil
	} else if exists && patch["options"] != nil {
		optionsList := patch["options"].([]map[string]interface{})
		for i, optionsItem := range optionsList {
			ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersOptionsAsPatch(optionsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "options")
	}
	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if !exists {
		delete(patch, "value")
	}
	path = rootPath + ".layout"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["layout"] = nil
	} else if !exists {
		delete(patch, "layout")
	}
	path = rootPath + ".associations"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["associations"] = nil
	} else if !exists {
		delete(patch, "associations")
	}
	path = rootPath + ".validation_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["validation_url"] = nil
	} else if !exists {
		delete(patch, "validation_url")
	}
	path = rootPath + ".options_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options_url"] = nil
	} else if !exists {
		delete(patch, "options_url")
	}
	path = rootPath + ".invalidmessage"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["invalidmessage"] = nil
	} else if !exists {
		delete(patch, "invalidmessage")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".required"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["required"] = nil
	} else if !exists {
		delete(patch, "required")
	}
	path = rootPath + ".pattern"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["pattern"] = nil
	} else if !exists {
		delete(patch, "pattern")
	}
	path = rootPath + ".placeholder"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["placeholder"] = nil
	} else if !exists {
		delete(patch, "placeholder")
	}
	path = rootPath + ".readonly"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["readonly"] = nil
	} else if !exists {
		delete(patch, "readonly")
	}
	path = rootPath + ".hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	} else if !exists {
		delete(patch, "hidden")
	}
	path = rootPath + ".i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["i18n"] = nil
	} else if !exists {
		delete(patch, "i18n")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersOptionsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".displayname"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["displayname"] = nil
	} else if !exists {
		delete(patch, "displayname")
	}
	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if !exists {
		delete(patch, "value")
	}
	path = rootPath + ".i18n"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["i18n"] = nil
	} else if exists && patch["i18n"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nAsPatch(patch["i18n"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "i18n")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsAsPatch(patch["en"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "en")
	}
	path = rootPath + ".de"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["de"] = nil
	} else if !exists {
		delete(patch, "de")
	}
	path = rootPath + ".es"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["es"] = nil
	} else if !exists {
		delete(patch, "es")
	}
	path = rootPath + ".fr"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["fr"] = nil
	} else if !exists {
		delete(patch, "fr")
	}
	path = rootPath + ".it"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["it"] = nil
	} else if !exists {
		delete(patch, "it")
	}
	path = rootPath + ".ja"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ja"] = nil
	} else if !exists {
		delete(patch, "ja")
	}
	path = rootPath + ".ko"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ko"] = nil
	} else if !exists {
		delete(patch, "ko")
	}
	path = rootPath + ".pt_br"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["pt_br"] = nil
	} else if !exists {
		delete(patch, "pt_br")
	}
	path = rootPath + ".zh_tw"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["zh_tw"] = nil
	} else if !exists {
		delete(patch, "zh_tw")
	}
	path = rootPath + ".zh_cn"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["zh_cn"] = nil
	} else if !exists {
		delete(patch, "zh_cn")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogMetadataServiceCustomParametersI18nFieldsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".displayname"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["displayname"] = nil
	} else if !exists {
		delete(patch, "displayname")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if !exists {
		delete(patch, "description")
	}
}

func ResourceIbmOnboardingCatalogDeploymentCatalogProductProviderAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
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

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUIAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
	} else if exists && patch["en"] != nil {
		ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUITranslatedContentAsPatch(patch["en"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "en")
	}
}

func ResourceIbmOnboardingCatalogDeploymentGlobalCatalogOverviewUITranslatedContentAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
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
