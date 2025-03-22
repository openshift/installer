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

func ResourceIbmOnboardingIamRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingIamRegistrationCreate,
		ReadContext:   resourceIbmOnboardingIamRegistrationRead,
		UpdateContext: resourceIbmOnboardingIamRegistrationUpdate,
		DeleteContext: resourceIbmOnboardingIamRegistrationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_iam_registration", "product_id"),
				Description:  "The unique ID of the product.",
			},
			"env": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_iam_registration", "env"),
				Description:  "The environment to fetch this object from.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_iam_registration", "name"),
				Description:  "The IAM registration name, which must be the programmatic name of the product.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the service is enabled or disabled for IAM.",
			},
			"service_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_iam_registration", "service_type"),
				Description:  "The type of the service.",
			},
			"actions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The product access management action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier for the action.",
						},
						"roles": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of roles for the action.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The description for the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The display name of the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"options": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Extra options.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Optional opt-in if action is hidden from customers.",
									},
								},
							},
						},
					},
				},
			},
			"additional_policy_scopes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of additional policy scopes.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The display name of the object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The fallback string for the description object.",
						},
						"en": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "English.",
						},
						"de": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "German.",
						},
						"es": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Spanish.",
						},
						"fr": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "French.",
						},
						"it": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Italian.",
						},
						"ja": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Japanese.",
						},
						"ko": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Korean.",
						},
						"pt_br": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Portuguese (Brazil).",
						},
						"zh_tw": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traditional Chinese.",
						},
						"zh_cn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Simplified Chinese.",
						},
					},
				},
			},
			"parent_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of parent IDs for product access management.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"resource_hierarchy_attribute": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The resource hierarchy key-value pair for composite services.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource hierarchy key.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource hierarchy value.",
						},
					},
				},
			},
			"supported_anonymous_accesses": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of supported anonymous accesses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The attributes for anonymous accesses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "An account id.",
									},
									"service_name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the service.",
									},
									"additional_properties": &schema.Schema{
										Type:        schema.TypeMap,
										Required:    true,
										Description: "Additional properties the key must come from supported attributes.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"roles": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The roles of supported anonymous accesses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"supported_attributes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of supported attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The supported attribute key.",
						},
						"options": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The list of support attribute options.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operators": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The supported attribute operator.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Optional opt-in if attribute is hidden from customers (customer can still use it if they found out themselves).",
									},
									"supported_patterns": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of supported patterns.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"policy_types": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of policy types.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"is_empty_value_supported": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicate whether the empty value is supported.",
									},
									"is_string_exists_false_value_supported": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicate whether the false value is supported for stringExists operator.",
									},
									"key": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of attribute.",
									},
									"resource_hierarchy": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Resource hierarchy options for composite services.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Hierarchy description key.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Key.",
															},
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value.",
															},
														},
													},
												},
												"value": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Hierarchy description value.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Key.",
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
						"display_name": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The display name of the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The description for the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"ui": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The user interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"input_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of the input.",
									},
									"input_details": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The details of the input.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "They type of the input details.",
												},
												"values": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The provided values of input details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The values of input details.",
															},
															"display_name": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The display name of the object.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"default": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The fallback string for the description object.",
																		},
																		"en": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "English.",
																		},
																		"de": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "German.",
																		},
																		"es": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Spanish.",
																		},
																		"fr": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "French.",
																		},
																		"it": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Italian.",
																		},
																		"ja": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Japanese.",
																		},
																		"ko": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Korean.",
																		},
																		"pt_br": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Portuguese (Brazil).",
																		},
																		"zh_tw": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Traditional Chinese.",
																		},
																		"zh_cn": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Simplified Chinese.",
																		},
																	},
																},
															},
														},
													},
												},
												"gst": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Required if type is gst.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"query": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The query to use.",
															},
															"value_property_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The value of the property name.",
															},
															"label_property_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "One of labelPropertyName or inputOptionLabel is required.",
															},
															"input_option_label": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The label for option input.",
															},
														},
													},
												},
												"url": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The URL data for user interface.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"url_endpoint": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The URL of the user interface interface.",
															},
															"input_option_label": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The label options for the user interface URL.",
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
			"supported_authorization_subjects": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of supported authorization subjects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The list of supported authorization subject properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the service.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of the service.",
									},
								},
							},
						},
						"roles": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of roles for authorization.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"supported_roles": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of roles that you can use to assign access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value belonging to the key.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The description for the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The display name of the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fallback string for the description object.",
									},
									"en": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "English.",
									},
									"de": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "German.",
									},
									"es": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Spanish.",
									},
									"fr": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "French.",
									},
									"it": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Italian.",
									},
									"ja": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Japanese.",
									},
									"ko": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Korean.",
									},
									"pt_br": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Portuguese (Brazil).",
									},
									"zh_tw": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Traditional Chinese.",
									},
									"zh_cn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Simplified Chinese.",
									},
								},
							},
						},
						"options": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The supported role options.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_policy": &schema.Schema{
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Optional opt-in to require access control on the role.",
									},
									"policy_type": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Optional opt-in to require checking policy type when applying the role.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"account_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional opt-in to require checking account type when applying the role.",
									},
								},
							},
						},
					},
				},
			},
			"supported_network": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The registration of set of endpoint types that are supported by your service in the `networkType` environment attribute. This constrains the context-based restriction rules specific to the service such that they describe access restrictions on only this set of endpoints.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_attributes": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The environment attribute for support.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the key.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of values that belong to the key.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The list of options for supported networks.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hidden": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the attribute is hidden or not.",
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
	}
}

func ResourceIbmOnboardingIamRegistrationValidator() *validate.ResourceValidator {
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
			Regexp:                     `^[a-zA-Z0-9\-.]+$`,
		},
		validate.ValidateSchema{
			Identifier:                 "service_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "platform_service, service",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_iam_registration", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingIamRegistrationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createIamRegistrationOptions := &partnercentersellv1.CreateIamRegistrationOptions{}

	createIamRegistrationOptions.SetProductID(d.Get("product_id").(string))
	createIamRegistrationOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("enabled"); ok {
		createIamRegistrationOptions.SetEnabled(d.Get("enabled").(bool))
	}
	if _, ok := d.GetOk("service_type"); ok {
		createIamRegistrationOptions.SetServiceType(d.Get("service_type").(string))
	}
	if _, ok := d.GetOk("actions"); ok {
		var actions []partnercentersellv1.IamServiceRegistrationAction
		for _, v := range d.Get("actions").([]interface{}) {
			value := v.(map[string]interface{})
			actionsItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationAction(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-actions").GetDiag()
			}
			actions = append(actions, *actionsItem)
		}
		createIamRegistrationOptions.SetActions(actions)
	}
	if _, ok := d.GetOk("additional_policy_scopes"); ok {
		var additionalPolicyScopes []string
		for _, v := range d.Get("additional_policy_scopes").([]interface{}) {
			additionalPolicyScopesItem := v.(string)
			additionalPolicyScopes = append(additionalPolicyScopes, additionalPolicyScopesItem)
		}
		createIamRegistrationOptions.SetAdditionalPolicyScopes(additionalPolicyScopes)
	}
	if _, ok := d.GetOk("display_name"); ok {
		displayNameModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(d.Get("display_name.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-display_name").GetDiag()
		}
		createIamRegistrationOptions.SetDisplayName(displayNameModel)
	}
	if _, ok := d.GetOk("parent_ids"); ok {
		var parentIds []string
		for _, v := range d.Get("parent_ids").([]interface{}) {
			parentIdsItem := v.(string)
			parentIds = append(parentIds, parentIdsItem)
		}
		createIamRegistrationOptions.SetParentIds(parentIds)
	}
	if _, ok := d.GetOk("resource_hierarchy_attribute"); ok {
		resourceHierarchyAttributeModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationResourceHierarchyAttribute(d.Get("resource_hierarchy_attribute.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-resource_hierarchy_attribute").GetDiag()
		}
		createIamRegistrationOptions.SetResourceHierarchyAttribute(resourceHierarchyAttributeModel)
	}
	if _, ok := d.GetOk("supported_anonymous_accesses"); ok {
		var supportedAnonymousAccesses []partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccess
		for _, v := range d.Get("supported_anonymous_accesses").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAnonymousAccessesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAnonymousAccess(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-supported_anonymous_accesses").GetDiag()
			}
			supportedAnonymousAccesses = append(supportedAnonymousAccesses, *supportedAnonymousAccessesItem)
		}
		createIamRegistrationOptions.SetSupportedAnonymousAccesses(supportedAnonymousAccesses)
	}
	if _, ok := d.GetOk("supported_attributes"); ok {
		var supportedAttributes []partnercentersellv1.IamServiceRegistrationSupportedAttribute
		for _, v := range d.Get("supported_attributes").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAttributesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAttribute(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-supported_attributes").GetDiag()
			}
			supportedAttributes = append(supportedAttributes, *supportedAttributesItem)
		}
		createIamRegistrationOptions.SetSupportedAttributes(supportedAttributes)
	}
	if _, ok := d.GetOk("supported_authorization_subjects"); ok {
		var supportedAuthorizationSubjects []partnercentersellv1.IamServiceRegistrationSupportedAuthorizationSubject
		for _, v := range d.Get("supported_authorization_subjects").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAuthorizationSubjectsItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAuthorizationSubject(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-supported_authorization_subjects").GetDiag()
			}
			supportedAuthorizationSubjects = append(supportedAuthorizationSubjects, *supportedAuthorizationSubjectsItem)
		}
		createIamRegistrationOptions.SetSupportedAuthorizationSubjects(supportedAuthorizationSubjects)
	}
	if _, ok := d.GetOk("supported_roles"); ok {
		var supportedRoles []partnercentersellv1.IamServiceRegistrationSupportedRole
		for _, v := range d.Get("supported_roles").([]interface{}) {
			value := v.(map[string]interface{})
			supportedRolesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedRole(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-supported_roles").GetDiag()
			}
			supportedRoles = append(supportedRoles, *supportedRolesItem)
		}
		createIamRegistrationOptions.SetSupportedRoles(supportedRoles)
	}
	if _, ok := d.GetOk("supported_network"); ok {
		supportedNetworkModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedNetwork(d.Get("supported_network.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "create", "parse-supported_network").GetDiag()
		}
		createIamRegistrationOptions.SetSupportedNetwork(supportedNetworkModel)
	}
	if _, ok := d.GetOk("env"); ok {
		createIamRegistrationOptions.SetEnv(d.Get("env").(string))
	}

	iamServiceRegistration, _, err := partnerCenterSellClient.CreateIamRegistrationWithContext(context, createIamRegistrationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateIamRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_iam_registration", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createIamRegistrationOptions.ProductID, *iamServiceRegistration.Name))

	return resourceIbmOnboardingIamRegistrationRead(context, d, meta)
}

func resourceIbmOnboardingIamRegistrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getIamRegistrationOptions := &partnercentersellv1.GetIamRegistrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "sep-id-parts").GetDiag()
	}

	getIamRegistrationOptions.SetProductID(parts[0])
	getIamRegistrationOptions.SetProgrammaticName(parts[1])
	if _, ok := d.GetOk("env"); ok {
		getIamRegistrationOptions.SetEnv(d.Get("env").(string))
	}

	iamServiceRegistration, response, err := partnerCenterSellClient.GetIamRegistrationWithContext(context, getIamRegistrationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetIamRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_iam_registration", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if parts[0] != "" {
		if err = d.Set("product_id", parts[0]); err != nil {
			err = fmt.Errorf("Error setting product_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-product_id").GetDiag()
		}
	}

	if err = d.Set("name", iamServiceRegistration.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-name").GetDiag()
	}
	if !core.IsNil(iamServiceRegistration.Enabled) {
		if err = d.Set("enabled", iamServiceRegistration.Enabled); err != nil {
			err = fmt.Errorf("Error setting enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-enabled").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.ServiceType) {
		if err = d.Set("service_type", iamServiceRegistration.ServiceType); err != nil {
			err = fmt.Errorf("Error setting service_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-service_type").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.Actions) {
		actions := []map[string]interface{}{}
		for _, actionsItem := range iamServiceRegistration.Actions {
			actionsItemMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionToMap(&actionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "actions-to-map").GetDiag()
			}
			actions = append(actions, actionsItemMap)
		}
		if err = d.Set("actions", actions); err != nil {
			err = fmt.Errorf("Error setting actions: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-actions").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.AdditionalPolicyScopes) {
		if err = d.Set("additional_policy_scopes", iamServiceRegistration.AdditionalPolicyScopes); err != nil {
			err = fmt.Errorf("Error setting additional_policy_scopes: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-additional_policy_scopes").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.DisplayName) {
		displayNameMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(iamServiceRegistration.DisplayName)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "display_name-to-map").GetDiag()
		}
		if err = d.Set("display_name", []map[string]interface{}{displayNameMap}); err != nil {
			err = fmt.Errorf("Error setting display_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-display_name").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.ParentIds) {
		if err = d.Set("parent_ids", iamServiceRegistration.ParentIds); err != nil {
			err = fmt.Errorf("Error setting parent_ids: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-parent_ids").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.ResourceHierarchyAttribute) {
		resourceHierarchyAttributeMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationResourceHierarchyAttributeToMap(iamServiceRegistration.ResourceHierarchyAttribute)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "resource_hierarchy_attribute-to-map").GetDiag()
		}
		if err = d.Set("resource_hierarchy_attribute", []map[string]interface{}{resourceHierarchyAttributeMap}); err != nil {
			err = fmt.Errorf("Error setting resource_hierarchy_attribute: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-resource_hierarchy_attribute").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.SupportedAnonymousAccesses) {
		supportedAnonymousAccesses := []map[string]interface{}{}
		for _, supportedAnonymousAccessesItem := range iamServiceRegistration.SupportedAnonymousAccesses {
			supportedAnonymousAccessesItemMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessToMap(&supportedAnonymousAccessesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "supported_anonymous_accesses-to-map").GetDiag()
			}
			supportedAnonymousAccesses = append(supportedAnonymousAccesses, supportedAnonymousAccessesItemMap)
		}
		if err = d.Set("supported_anonymous_accesses", supportedAnonymousAccesses); err != nil {
			err = fmt.Errorf("Error setting supported_anonymous_accesses: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-supported_anonymous_accesses").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.SupportedAttributes) {
		supportedAttributes := []map[string]interface{}{}
		for _, supportedAttributesItem := range iamServiceRegistration.SupportedAttributes {
			supportedAttributesItemMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAttributeToMap(&supportedAttributesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "supported_attributes-to-map").GetDiag()
			}
			supportedAttributes = append(supportedAttributes, supportedAttributesItemMap)
		}
		if err = d.Set("supported_attributes", supportedAttributes); err != nil {
			err = fmt.Errorf("Error setting supported_attributes: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-supported_attributes").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.SupportedAuthorizationSubjects) {
		supportedAuthorizationSubjects := []map[string]interface{}{}
		for _, supportedAuthorizationSubjectsItem := range iamServiceRegistration.SupportedAuthorizationSubjects {
			supportedAuthorizationSubjectsItemMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAuthorizationSubjectToMap(&supportedAuthorizationSubjectsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "supported_authorization_subjects-to-map").GetDiag()
			}
			supportedAuthorizationSubjects = append(supportedAuthorizationSubjects, supportedAuthorizationSubjectsItemMap)
		}
		if err = d.Set("supported_authorization_subjects", supportedAuthorizationSubjects); err != nil {
			err = fmt.Errorf("Error setting supported_authorization_subjects: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-supported_authorization_subjects").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.SupportedRoles) {
		supportedRoles := []map[string]interface{}{}
		for _, supportedRolesItem := range iamServiceRegistration.SupportedRoles {
			supportedRolesItemMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedRoleToMap(&supportedRolesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "supported_roles-to-map").GetDiag()
			}
			supportedRoles = append(supportedRoles, supportedRolesItemMap)
		}
		if err = d.Set("supported_roles", supportedRoles); err != nil {
			err = fmt.Errorf("Error setting supported_roles: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-supported_roles").GetDiag()
		}
	}
	if !core.IsNil(iamServiceRegistration.SupportedNetwork) {
		supportedNetworkMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedNetworkToMap(iamServiceRegistration.SupportedNetwork)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "supported_network-to-map").GetDiag()
		}
		if err = d.Set("supported_network", []map[string]interface{}{supportedNetworkMap}); err != nil {
			err = fmt.Errorf("Error setting supported_network: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "read", "set-supported_network").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingIamRegistrationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateIamRegistrationOptions := &partnercentersellv1.UpdateIamRegistrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "sep-id-parts").GetDiag()
	}

	updateIamRegistrationOptions.SetProductID(parts[0])
	updateIamRegistrationOptions.SetProgrammaticName(parts[1])
	if _, ok := d.GetOk("env"); ok {
		updateIamRegistrationOptions.SetEnv(d.Get("env").(string))
	}

	hasChange := false

	patchVals := &partnercentersellv1.IamServiceRegistrationPatch{}
	if d.HasChange("product_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "product_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_onboarding_iam_registration", "update", "product_id-forces-new").GetDiag()
	}
	if d.HasChange("enabled") {
		newEnabled := d.Get("enabled").(bool)
		patchVals.Enabled = &newEnabled
		hasChange = true
	}
	if d.HasChange("service_type") {
		newServiceType := d.Get("service_type").(string)
		patchVals.ServiceType = &newServiceType
		hasChange = true
	}
	if d.HasChange("actions") {
		var actions []partnercentersellv1.IamServiceRegistrationAction
		for _, v := range d.Get("actions").([]interface{}) {
			value := v.(map[string]interface{})
			actionsItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationAction(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-actions").GetDiag()
			}
			actions = append(actions, *actionsItem)
		}
		patchVals.Actions = actions
		hasChange = true
	}
	if d.HasChange("additional_policy_scopes") {
		var additionalPolicyScopes []string
		for _, v := range d.Get("additional_policy_scopes").([]interface{}) {
			additionalPolicyScopesItem := v.(string)
			additionalPolicyScopes = append(additionalPolicyScopes, additionalPolicyScopesItem)
		}
		patchVals.AdditionalPolicyScopes = additionalPolicyScopes
		hasChange = true
	}
	if d.HasChange("display_name") {
		displayName, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(d.Get("display_name.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-display_name").GetDiag()
		}
		patchVals.DisplayName = displayName
		hasChange = true
	}
	if d.HasChange("parent_ids") {
		var parentIds []string
		for _, v := range d.Get("parent_ids").([]interface{}) {
			parentIdsItem := v.(string)
			parentIds = append(parentIds, parentIdsItem)
		}
		patchVals.ParentIds = parentIds
		hasChange = true
	}
	if d.HasChange("resource_hierarchy_attribute") {
		resourceHierarchyAttribute, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationResourceHierarchyAttribute(d.Get("resource_hierarchy_attribute.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-resource_hierarchy_attribute").GetDiag()
		}
		patchVals.ResourceHierarchyAttribute = resourceHierarchyAttribute
		hasChange = true
	}
	if d.HasChange("supported_anonymous_accesses") {
		var supportedAnonymousAccesses []partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccess
		for _, v := range d.Get("supported_anonymous_accesses").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAnonymousAccessesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAnonymousAccess(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-supported_anonymous_accesses").GetDiag()
			}
			supportedAnonymousAccesses = append(supportedAnonymousAccesses, *supportedAnonymousAccessesItem)
		}
		patchVals.SupportedAnonymousAccesses = supportedAnonymousAccesses
		hasChange = true
	}
	if d.HasChange("supported_attributes") {
		var supportedAttributes []partnercentersellv1.IamServiceRegistrationSupportedAttribute
		for _, v := range d.Get("supported_attributes").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAttributesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAttribute(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-supported_attributes").GetDiag()
			}
			supportedAttributes = append(supportedAttributes, *supportedAttributesItem)
		}
		patchVals.SupportedAttributes = supportedAttributes
		hasChange = true
	}
	if d.HasChange("supported_authorization_subjects") {
		var supportedAuthorizationSubjects []partnercentersellv1.IamServiceRegistrationSupportedAuthorizationSubject
		for _, v := range d.Get("supported_authorization_subjects").([]interface{}) {
			value := v.(map[string]interface{})
			supportedAuthorizationSubjectsItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAuthorizationSubject(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-supported_authorization_subjects").GetDiag()
			}
			supportedAuthorizationSubjects = append(supportedAuthorizationSubjects, *supportedAuthorizationSubjectsItem)
		}
		patchVals.SupportedAuthorizationSubjects = supportedAuthorizationSubjects
		hasChange = true
	}
	if d.HasChange("supported_roles") {
		var supportedRoles []partnercentersellv1.IamServiceRegistrationSupportedRole
		for _, v := range d.Get("supported_roles").([]interface{}) {
			value := v.(map[string]interface{})
			supportedRolesItem, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedRole(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-supported_roles").GetDiag()
			}
			supportedRoles = append(supportedRoles, *supportedRolesItem)
		}
		patchVals.SupportedRoles = supportedRoles
		hasChange = true
	}
	if d.HasChange("supported_network") {
		supportedNetwork, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedNetwork(d.Get("supported_network.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "update", "parse-supported_network").GetDiag()
		}
		patchVals.SupportedNetwork = supportedNetwork
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateIamRegistrationOptions.IamRegistrationPatch = ResourceIbmOnboardingIamRegistrationIamServiceRegistrationPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateIamRegistrationWithContext(context, updateIamRegistrationOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateIamRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_iam_registration", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingIamRegistrationRead(context, d, meta)
}

func resourceIbmOnboardingIamRegistrationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteIamRegistrationOptions := &partnercentersellv1.DeleteIamRegistrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_iam_registration", "delete", "sep-id-parts").GetDiag()
	}

	deleteIamRegistrationOptions.SetProductID(parts[0])
	deleteIamRegistrationOptions.SetProgrammaticName(parts[1])
	if _, ok := d.GetOk("env"); ok {
		deleteIamRegistrationOptions.SetEnv(d.Get("env").(string))
	}

	_, err = partnerCenterSellClient.DeleteIamRegistrationWithContext(context, deleteIamRegistrationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteIamRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_iam_registration", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationAction(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationAction, error) {
	model := &partnercentersellv1.IamServiceRegistrationAction{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["roles"] != nil {
		roles := []string{}
		for _, rolesItem := range modelMap["roles"].([]interface{}) {
			roles = append(roles, rolesItem.(string))
		}
		model.Roles = roles
	}
	if modelMap["description"] != nil && len(modelMap["description"].([]interface{})) > 0 && modelMap["description"].([]interface{})[0] != nil {
		DescriptionModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDescriptionObject(modelMap["description"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Description = DescriptionModel
	}
	if modelMap["display_name"] != nil && len(modelMap["display_name"].([]interface{})) > 0 && modelMap["display_name"].([]interface{})[0] != nil {
		DisplayNameModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(modelMap["display_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DisplayName = DisplayNameModel
	}
	if modelMap["options"] != nil && len(modelMap["options"].([]interface{})) > 0 && modelMap["options"].([]interface{})[0] != nil {
		OptionsModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationActionOptions(modelMap["options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Options = OptionsModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDescriptionObject(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationDescriptionObject, error) {
	model := &partnercentersellv1.IamServiceRegistrationDescriptionObject{}
	if modelMap["default"] != nil && modelMap["default"].(string) != "" {
		model.Default = core.StringPtr(modelMap["default"].(string))
	}
	if modelMap["en"] != nil && modelMap["en"].(string) != "" {
		model.En = core.StringPtr(modelMap["en"].(string))
	}
	if modelMap["de"] != nil && modelMap["de"].(string) != "" {
		model.De = core.StringPtr(modelMap["de"].(string))
	}
	if modelMap["es"] != nil && modelMap["es"].(string) != "" {
		model.Es = core.StringPtr(modelMap["es"].(string))
	}
	if modelMap["fr"] != nil && modelMap["fr"].(string) != "" {
		model.Fr = core.StringPtr(modelMap["fr"].(string))
	}
	if modelMap["it"] != nil && modelMap["it"].(string) != "" {
		model.It = core.StringPtr(modelMap["it"].(string))
	}
	if modelMap["ja"] != nil && modelMap["ja"].(string) != "" {
		model.Ja = core.StringPtr(modelMap["ja"].(string))
	}
	if modelMap["ko"] != nil && modelMap["ko"].(string) != "" {
		model.Ko = core.StringPtr(modelMap["ko"].(string))
	}
	if modelMap["pt_br"] != nil && modelMap["pt_br"].(string) != "" {
		model.PtBr = core.StringPtr(modelMap["pt_br"].(string))
	}
	if modelMap["zh_tw"] != nil && modelMap["zh_tw"].(string) != "" {
		model.ZhTw = core.StringPtr(modelMap["zh_tw"].(string))
	}
	if modelMap["zh_cn"] != nil && modelMap["zh_cn"].(string) != "" {
		model.ZhCn = core.StringPtr(modelMap["zh_cn"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationDisplayNameObject, error) {
	model := &partnercentersellv1.IamServiceRegistrationDisplayNameObject{}
	if modelMap["default"] != nil && modelMap["default"].(string) != "" {
		model.Default = core.StringPtr(modelMap["default"].(string))
	}
	if modelMap["en"] != nil && modelMap["en"].(string) != "" {
		model.En = core.StringPtr(modelMap["en"].(string))
	}
	if modelMap["de"] != nil && modelMap["de"].(string) != "" {
		model.De = core.StringPtr(modelMap["de"].(string))
	}
	if modelMap["es"] != nil && modelMap["es"].(string) != "" {
		model.Es = core.StringPtr(modelMap["es"].(string))
	}
	if modelMap["fr"] != nil && modelMap["fr"].(string) != "" {
		model.Fr = core.StringPtr(modelMap["fr"].(string))
	}
	if modelMap["it"] != nil && modelMap["it"].(string) != "" {
		model.It = core.StringPtr(modelMap["it"].(string))
	}
	if modelMap["ja"] != nil && modelMap["ja"].(string) != "" {
		model.Ja = core.StringPtr(modelMap["ja"].(string))
	}
	if modelMap["ko"] != nil && modelMap["ko"].(string) != "" {
		model.Ko = core.StringPtr(modelMap["ko"].(string))
	}
	if modelMap["pt_br"] != nil && modelMap["pt_br"].(string) != "" {
		model.PtBr = core.StringPtr(modelMap["pt_br"].(string))
	}
	if modelMap["zh_tw"] != nil && modelMap["zh_tw"].(string) != "" {
		model.ZhTw = core.StringPtr(modelMap["zh_tw"].(string))
	}
	if modelMap["zh_cn"] != nil && modelMap["zh_cn"].(string) != "" {
		model.ZhCn = core.StringPtr(modelMap["zh_cn"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationActionOptions(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationActionOptions, error) {
	model := &partnercentersellv1.IamServiceRegistrationActionOptions{}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationResourceHierarchyAttribute(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationResourceHierarchyAttribute, error) {
	model := &partnercentersellv1.IamServiceRegistrationResourceHierarchyAttribute{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAnonymousAccess(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccess, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccess{}
	if modelMap["attributes"] != nil && len(modelMap["attributes"].([]interface{})) > 0 && modelMap["attributes"].([]interface{})[0] != nil {
		AttributesModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAnonymousAccessAttributes(modelMap["attributes"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Attributes = AttributesModel
	}
	if modelMap["roles"] != nil {
		roles := []string{}
		for _, rolesItem := range modelMap["roles"].([]interface{}) {
			roles = append(roles, rolesItem.(string))
		}
		model.Roles = roles
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAnonymousAccessAttributes(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccessAttributes, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccessAttributes{}
	model.AccountID = core.StringPtr(modelMap["account_id"].(string))
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	if modelMap["additional_properties"] != nil {
		model.AdditionalProperties = make(map[string]string)
		for key, value := range modelMap["additional_properties"].(map[string]interface{}) {
			if str, ok := value.(string); ok {
				model.AdditionalProperties[key] = str
			}
		}
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAttribute(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedAttribute, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedAttribute{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["options"] != nil && len(modelMap["options"].([]interface{})) > 0 && modelMap["options"].([]interface{})[0] != nil {
		OptionsModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptions(modelMap["options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Options = OptionsModel
	}
	if modelMap["display_name"] != nil && len(modelMap["display_name"].([]interface{})) > 0 && modelMap["display_name"].([]interface{})[0] != nil {
		DisplayNameModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(modelMap["display_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DisplayName = DisplayNameModel
	}
	if modelMap["description"] != nil && len(modelMap["description"].([]interface{})) > 0 && modelMap["description"].([]interface{})[0] != nil {
		DescriptionModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDescriptionObject(modelMap["description"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Description = DescriptionModel
	}
	if modelMap["ui"] != nil && len(modelMap["ui"].([]interface{})) > 0 && modelMap["ui"].([]interface{})[0] != nil {
		UiModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUi(modelMap["ui"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ui = UiModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptions(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributesOptions, error) {
	model := &partnercentersellv1.SupportedAttributesOptions{}
	if modelMap["operators"] != nil {
		operators := []string{}
		for _, operatorsItem := range modelMap["operators"].([]interface{}) {
			operators = append(operators, operatorsItem.(string))
		}
		model.Operators = operators
	}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	if modelMap["supported_patterns"] != nil {
		supportedPatterns := []string{}
		for _, supportedPatternsItem := range modelMap["supported_patterns"].([]interface{}) {
			supportedPatterns = append(supportedPatterns, supportedPatternsItem.(string))
		}
		model.SupportedPatterns = supportedPatterns
	}
	if modelMap["policy_types"] != nil {
		policyTypes := []string{}
		for _, policyTypesItem := range modelMap["policy_types"].([]interface{}) {
			policyTypes = append(policyTypes, policyTypesItem.(string))
		}
		model.PolicyTypes = policyTypes
	}
	if modelMap["is_empty_value_supported"] != nil {
		model.IsEmptyValueSupported = core.BoolPtr(modelMap["is_empty_value_supported"].(bool))
	}
	if modelMap["is_string_exists_false_value_supported"] != nil {
		model.IsStringExistsFalseValueSupported = core.BoolPtr(modelMap["is_string_exists_false_value_supported"].(bool))
	}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["resource_hierarchy"] != nil && len(modelMap["resource_hierarchy"].([]interface{})) > 0 && modelMap["resource_hierarchy"].([]interface{})[0] != nil {
		ResourceHierarchyModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchy(modelMap["resource_hierarchy"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceHierarchy = ResourceHierarchyModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchy(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributesOptionsResourceHierarchy, error) {
	model := &partnercentersellv1.SupportedAttributesOptionsResourceHierarchy{}
	if modelMap["key"] != nil && len(modelMap["key"].([]interface{})) > 0 && modelMap["key"].([]interface{})[0] != nil {
		KeyModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchyKey(modelMap["key"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Key = KeyModel
	}
	if modelMap["value"] != nil && len(modelMap["value"].([]interface{})) > 0 && modelMap["value"].([]interface{})[0] != nil {
		ValueModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchyValue(modelMap["value"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Value = ValueModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchyKey(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributesOptionsResourceHierarchyKey, error) {
	model := &partnercentersellv1.SupportedAttributesOptionsResourceHierarchyKey{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributesOptionsResourceHierarchyValue(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributesOptionsResourceHierarchyValue, error) {
	model := &partnercentersellv1.SupportedAttributesOptionsResourceHierarchyValue{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUi(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributeUi, error) {
	model := &partnercentersellv1.SupportedAttributeUi{}
	if modelMap["input_type"] != nil && modelMap["input_type"].(string) != "" {
		model.InputType = core.StringPtr(modelMap["input_type"].(string))
	}
	if modelMap["input_details"] != nil && len(modelMap["input_details"].([]interface{})) > 0 && modelMap["input_details"].([]interface{})[0] != nil {
		InputDetailsModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputDetails(modelMap["input_details"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.InputDetails = InputDetailsModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputDetails(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributeUiInputDetails, error) {
	model := &partnercentersellv1.SupportedAttributeUiInputDetails{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["values"] != nil {
		values := []partnercentersellv1.SupportedAttributeUiInputValue{}
		for _, valuesItem := range modelMap["values"].([]interface{}) {
			valuesItemModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputValue(valuesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			values = append(values, *valuesItemModel)
		}
		model.Values = values
	}
	if modelMap["gst"] != nil && len(modelMap["gst"].([]interface{})) > 0 && modelMap["gst"].([]interface{})[0] != nil {
		GstModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputGst(modelMap["gst"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Gst = GstModel
	}
	if modelMap["url"] != nil && len(modelMap["url"].([]interface{})) > 0 && modelMap["url"].([]interface{})[0] != nil {
		URLModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputURL(modelMap["url"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.URL = URLModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputValue(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributeUiInputValue, error) {
	model := &partnercentersellv1.SupportedAttributeUiInputValue{}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["display_name"] != nil && len(modelMap["display_name"].([]interface{})) > 0 && modelMap["display_name"].([]interface{})[0] != nil {
		DisplayNameModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(modelMap["display_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DisplayName = DisplayNameModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputGst(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributeUiInputGst, error) {
	model := &partnercentersellv1.SupportedAttributeUiInputGst{}
	if modelMap["query"] != nil && modelMap["query"].(string) != "" {
		model.Query = core.StringPtr(modelMap["query"].(string))
	}
	if modelMap["value_property_name"] != nil && modelMap["value_property_name"].(string) != "" {
		model.ValuePropertyName = core.StringPtr(modelMap["value_property_name"].(string))
	}
	if modelMap["label_property_name"] != nil && modelMap["label_property_name"].(string) != "" {
		model.LabelPropertyName = core.StringPtr(modelMap["label_property_name"].(string))
	}
	if modelMap["input_option_label"] != nil && modelMap["input_option_label"].(string) != "" {
		model.InputOptionLabel = core.StringPtr(modelMap["input_option_label"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedAttributeUiInputURL(modelMap map[string]interface{}) (*partnercentersellv1.SupportedAttributeUiInputURL, error) {
	model := &partnercentersellv1.SupportedAttributeUiInputURL{}
	if modelMap["url_endpoint"] != nil && modelMap["url_endpoint"].(string) != "" {
		model.UrlEndpoint = core.StringPtr(modelMap["url_endpoint"].(string))
	}
	if modelMap["input_option_label"] != nil && modelMap["input_option_label"].(string) != "" {
		model.InputOptionLabel = core.StringPtr(modelMap["input_option_label"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedAuthorizationSubject(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedAuthorizationSubject, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedAuthorizationSubject{}
	if modelMap["attributes"] != nil && len(modelMap["attributes"].([]interface{})) > 0 && modelMap["attributes"].([]interface{})[0] != nil {
		AttributesModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportAuthorizationSubjectAttribute(modelMap["attributes"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Attributes = AttributesModel
	}
	if modelMap["roles"] != nil {
		roles := []string{}
		for _, rolesItem := range modelMap["roles"].([]interface{}) {
			roles = append(roles, rolesItem.(string))
		}
		model.Roles = roles
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportAuthorizationSubjectAttribute(modelMap map[string]interface{}) (*partnercentersellv1.SupportAuthorizationSubjectAttribute, error) {
	model := &partnercentersellv1.SupportAuthorizationSubjectAttribute{}
	if modelMap["service_name"] != nil && modelMap["service_name"].(string) != "" {
		model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	}
	if modelMap["resource_type"] != nil && modelMap["resource_type"].(string) != "" {
		model.ResourceType = core.StringPtr(modelMap["resource_type"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedRole(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedRole, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedRole{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["description"] != nil && len(modelMap["description"].([]interface{})) > 0 && modelMap["description"].([]interface{})[0] != nil {
		DescriptionModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDescriptionObject(modelMap["description"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Description = DescriptionModel
	}
	if modelMap["display_name"] != nil && len(modelMap["display_name"].([]interface{})) > 0 && modelMap["display_name"].([]interface{})[0] != nil {
		DisplayNameModel, err := ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationDisplayNameObject(modelMap["display_name"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DisplayName = DisplayNameModel
	}
	if modelMap["options"] != nil && len(modelMap["options"].([]interface{})) > 0 && modelMap["options"].([]interface{})[0] != nil {
		OptionsModel, err := ResourceIbmOnboardingIamRegistrationMapToSupportedRoleOptions(modelMap["options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Options = OptionsModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToSupportedRoleOptions(modelMap map[string]interface{}) (*partnercentersellv1.SupportedRoleOptions, error) {
	model := &partnercentersellv1.SupportedRoleOptions{}
	model.AccessPolicy = core.BoolPtr(modelMap["access_policy"].(bool))
	if modelMap["policy_type"] != nil {
		policyType := []string{}
		for _, policyTypeItem := range modelMap["policy_type"].([]interface{}) {
			policyType = append(policyType, policyTypeItem.(string))
		}
		model.PolicyType = policyType
	}
	if modelMap["account_type"] != nil && modelMap["account_type"].(string) != "" {
		model.AccountType = core.StringPtr(modelMap["account_type"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToIamServiceRegistrationSupportedNetwork(modelMap map[string]interface{}) (*partnercentersellv1.IamServiceRegistrationSupportedNetwork, error) {
	model := &partnercentersellv1.IamServiceRegistrationSupportedNetwork{}
	if modelMap["environment_attributes"] != nil {
		environmentAttributes := []partnercentersellv1.EnvironmentAttribute{}
		for _, environmentAttributesItem := range modelMap["environment_attributes"].([]interface{}) {
			environmentAttributesItemModel, err := ResourceIbmOnboardingIamRegistrationMapToEnvironmentAttribute(environmentAttributesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			environmentAttributes = append(environmentAttributes, *environmentAttributesItemModel)
		}
		model.EnvironmentAttributes = environmentAttributes
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToEnvironmentAttribute(modelMap map[string]interface{}) (*partnercentersellv1.EnvironmentAttribute, error) {
	model := &partnercentersellv1.EnvironmentAttribute{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["values"] != nil {
		values := []string{}
		for _, valuesItem := range modelMap["values"].([]interface{}) {
			values = append(values, valuesItem.(string))
		}
		model.Values = values
	}
	if modelMap["options"] != nil && len(modelMap["options"].([]interface{})) > 0 && modelMap["options"].([]interface{})[0] != nil {
		OptionsModel, err := ResourceIbmOnboardingIamRegistrationMapToEnvironmentAttributeOptions(modelMap["options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Options = OptionsModel
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationMapToEnvironmentAttributeOptions(modelMap map[string]interface{}) (*partnercentersellv1.EnvironmentAttributeOptions, error) {
	model := &partnercentersellv1.EnvironmentAttributeOptions{}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	return model, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionToMap(model *partnercentersellv1.IamServiceRegistrationAction) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Roles != nil {
		modelMap["roles"] = model.Roles
	}
	if model.Description != nil {
		descriptionMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectToMap(model.Description)
		if err != nil {
			return modelMap, err
		}
		modelMap["description"] = []map[string]interface{}{descriptionMap}
	}
	if model.DisplayName != nil {
		displayNameMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(model.DisplayName)
		if err != nil {
			return modelMap, err
		}
		modelMap["display_name"] = []map[string]interface{}{displayNameMap}
	}
	if model.Options != nil {
		optionsMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionOptionsToMap(model.Options)
		if err != nil {
			return modelMap, err
		}
		modelMap["options"] = []map[string]interface{}{optionsMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectToMap(model *partnercentersellv1.IamServiceRegistrationDescriptionObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Default != nil {
		modelMap["default"] = *model.Default
	}
	if model.En != nil {
		modelMap["en"] = *model.En
	}
	if model.De != nil {
		modelMap["de"] = *model.De
	}
	if model.Es != nil {
		modelMap["es"] = *model.Es
	}
	if model.Fr != nil {
		modelMap["fr"] = *model.Fr
	}
	if model.It != nil {
		modelMap["it"] = *model.It
	}
	if model.Ja != nil {
		modelMap["ja"] = *model.Ja
	}
	if model.Ko != nil {
		modelMap["ko"] = *model.Ko
	}
	if model.PtBr != nil {
		modelMap["pt_br"] = *model.PtBr
	}
	if model.ZhTw != nil {
		modelMap["zh_tw"] = *model.ZhTw
	}
	if model.ZhCn != nil {
		modelMap["zh_cn"] = *model.ZhCn
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(model *partnercentersellv1.IamServiceRegistrationDisplayNameObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Default != nil {
		modelMap["default"] = *model.Default
	}
	if model.En != nil {
		modelMap["en"] = *model.En
	}
	if model.De != nil {
		modelMap["de"] = *model.De
	}
	if model.Es != nil {
		modelMap["es"] = *model.Es
	}
	if model.Fr != nil {
		modelMap["fr"] = *model.Fr
	}
	if model.It != nil {
		modelMap["it"] = *model.It
	}
	if model.Ja != nil {
		modelMap["ja"] = *model.Ja
	}
	if model.Ko != nil {
		modelMap["ko"] = *model.Ko
	}
	if model.PtBr != nil {
		modelMap["pt_br"] = *model.PtBr
	}
	if model.ZhTw != nil {
		modelMap["zh_tw"] = *model.ZhTw
	}
	if model.ZhCn != nil {
		modelMap["zh_cn"] = *model.ZhCn
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionOptionsToMap(model *partnercentersellv1.IamServiceRegistrationActionOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationResourceHierarchyAttributeToMap(model *partnercentersellv1.IamServiceRegistrationResourceHierarchyAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessToMap(model *partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccess) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Attributes != nil {
		attributesMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessAttributesToMap(model.Attributes)
		if err != nil {
			return modelMap, err
		}
		modelMap["attributes"] = []map[string]interface{}{attributesMap}
	}
	if model.Roles != nil {
		modelMap["roles"] = model.Roles
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessAttributesToMap(model *partnercentersellv1.IamServiceRegistrationSupportedAnonymousAccessAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["account_id"] = *model.AccountID
	modelMap["service_name"] = *model.ServiceName
	additionalProperties := make(map[string]interface{})
	for k, v := range model.AdditionalProperties {
		additionalProperties[k] = flex.Stringify(v)
	}
	modelMap["additional_properties"] = additionalProperties
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAttributeToMap(model *partnercentersellv1.IamServiceRegistrationSupportedAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Options != nil {
		optionsMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsToMap(model.Options)
		if err != nil {
			return modelMap, err
		}
		modelMap["options"] = []map[string]interface{}{optionsMap}
	}
	if model.DisplayName != nil {
		displayNameMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(model.DisplayName)
		if err != nil {
			return modelMap, err
		}
		modelMap["display_name"] = []map[string]interface{}{displayNameMap}
	}
	if model.Description != nil {
		descriptionMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectToMap(model.Description)
		if err != nil {
			return modelMap, err
		}
		modelMap["description"] = []map[string]interface{}{descriptionMap}
	}
	if model.Ui != nil {
		uiMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributeUiToMap(model.Ui)
		if err != nil {
			return modelMap, err
		}
		modelMap["ui"] = []map[string]interface{}{uiMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsToMap(model *partnercentersellv1.SupportedAttributesOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Operators != nil {
		modelMap["operators"] = model.Operators
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.SupportedPatterns != nil {
		modelMap["supported_patterns"] = model.SupportedPatterns
	}
	if model.PolicyTypes != nil {
		modelMap["policy_types"] = model.PolicyTypes
	}
	if model.IsEmptyValueSupported != nil {
		modelMap["is_empty_value_supported"] = *model.IsEmptyValueSupported
	}
	if model.IsStringExistsFalseValueSupported != nil {
		modelMap["is_string_exists_false_value_supported"] = *model.IsStringExistsFalseValueSupported
	}
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.ResourceHierarchy != nil {
		resourceHierarchyMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyToMap(model.ResourceHierarchy)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_hierarchy"] = []map[string]interface{}{resourceHierarchyMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyToMap(model *partnercentersellv1.SupportedAttributesOptionsResourceHierarchy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		keyMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyKeyToMap(model.Key)
		if err != nil {
			return modelMap, err
		}
		modelMap["key"] = []map[string]interface{}{keyMap}
	}
	if model.Value != nil {
		valueMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyValueToMap(model.Value)
		if err != nil {
			return modelMap, err
		}
		modelMap["value"] = []map[string]interface{}{valueMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyKeyToMap(model *partnercentersellv1.SupportedAttributesOptionsResourceHierarchyKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyValueToMap(model *partnercentersellv1.SupportedAttributesOptionsResourceHierarchyValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiToMap(model *partnercentersellv1.SupportedAttributeUi) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InputType != nil {
		modelMap["input_type"] = *model.InputType
	}
	if model.InputDetails != nil {
		inputDetailsMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputDetailsToMap(model.InputDetails)
		if err != nil {
			return modelMap, err
		}
		modelMap["input_details"] = []map[string]interface{}{inputDetailsMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputDetailsToMap(model *partnercentersellv1.SupportedAttributeUiInputDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Values != nil {
		values := []map[string]interface{}{}
		for _, valuesItem := range model.Values {
			valuesItemMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputValueToMap(&valuesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			values = append(values, valuesItemMap)
		}
		modelMap["values"] = values
	}
	if model.Gst != nil {
		gstMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputGstToMap(model.Gst)
		if err != nil {
			return modelMap, err
		}
		modelMap["gst"] = []map[string]interface{}{gstMap}
	}
	if model.URL != nil {
		urlMap, err := ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputURLToMap(model.URL)
		if err != nil {
			return modelMap, err
		}
		modelMap["url"] = []map[string]interface{}{urlMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputValueToMap(model *partnercentersellv1.SupportedAttributeUiInputValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.DisplayName != nil {
		displayNameMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(model.DisplayName)
		if err != nil {
			return modelMap, err
		}
		modelMap["display_name"] = []map[string]interface{}{displayNameMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputGstToMap(model *partnercentersellv1.SupportedAttributeUiInputGst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Query != nil {
		modelMap["query"] = *model.Query
	}
	if model.ValuePropertyName != nil {
		modelMap["value_property_name"] = *model.ValuePropertyName
	}
	if model.LabelPropertyName != nil {
		modelMap["label_property_name"] = *model.LabelPropertyName
	}
	if model.InputOptionLabel != nil {
		modelMap["input_option_label"] = *model.InputOptionLabel
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputURLToMap(model *partnercentersellv1.SupportedAttributeUiInputURL) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UrlEndpoint != nil {
		modelMap["url_endpoint"] = *model.UrlEndpoint
	}
	if model.InputOptionLabel != nil {
		modelMap["input_option_label"] = *model.InputOptionLabel
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAuthorizationSubjectToMap(model *partnercentersellv1.IamServiceRegistrationSupportedAuthorizationSubject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Attributes != nil {
		attributesMap, err := ResourceIbmOnboardingIamRegistrationSupportAuthorizationSubjectAttributeToMap(model.Attributes)
		if err != nil {
			return modelMap, err
		}
		modelMap["attributes"] = []map[string]interface{}{attributesMap}
	}
	if model.Roles != nil {
		modelMap["roles"] = model.Roles
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportAuthorizationSubjectAttributeToMap(model *partnercentersellv1.SupportAuthorizationSubjectAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServiceName != nil {
		modelMap["service_name"] = *model.ServiceName
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedRoleToMap(model *partnercentersellv1.IamServiceRegistrationSupportedRole) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Description != nil {
		descriptionMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectToMap(model.Description)
		if err != nil {
			return modelMap, err
		}
		modelMap["description"] = []map[string]interface{}{descriptionMap}
	}
	if model.DisplayName != nil {
		displayNameMap, err := ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectToMap(model.DisplayName)
		if err != nil {
			return modelMap, err
		}
		modelMap["display_name"] = []map[string]interface{}{displayNameMap}
	}
	if model.Options != nil {
		optionsMap, err := ResourceIbmOnboardingIamRegistrationSupportedRoleOptionsToMap(model.Options)
		if err != nil {
			return modelMap, err
		}
		modelMap["options"] = []map[string]interface{}{optionsMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationSupportedRoleOptionsToMap(model *partnercentersellv1.SupportedRoleOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["access_policy"] = *model.AccessPolicy
	if model.PolicyType != nil {
		modelMap["policy_type"] = model.PolicyType
	}
	if model.AccountType != nil {
		modelMap["account_type"] = *model.AccountType
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedNetworkToMap(model *partnercentersellv1.IamServiceRegistrationSupportedNetwork) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnvironmentAttributes != nil {
		environmentAttributes := []map[string]interface{}{}
		for _, environmentAttributesItem := range model.EnvironmentAttributes {
			environmentAttributesItemMap, err := ResourceIbmOnboardingIamRegistrationEnvironmentAttributeToMap(&environmentAttributesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			environmentAttributes = append(environmentAttributes, environmentAttributesItemMap)
		}
		modelMap["environment_attributes"] = environmentAttributes
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationEnvironmentAttributeToMap(model *partnercentersellv1.EnvironmentAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	if model.Options != nil {
		optionsMap, err := ResourceIbmOnboardingIamRegistrationEnvironmentAttributeOptionsToMap(model.Options)
		if err != nil {
			return modelMap, err
		}
		modelMap["options"] = []map[string]interface{}{optionsMap}
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationEnvironmentAttributeOptionsToMap(model *partnercentersellv1.EnvironmentAttributeOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	return modelMap, nil
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationPatchAsPatch(patchVals *partnercentersellv1.IamServiceRegistrationPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "enabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["enabled"] = nil
	} else if !exists {
		delete(patch, "enabled")
	}
	path = "service_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service_type"] = nil
	} else if !exists {
		delete(patch, "service_type")
	}
	path = "actions"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["actions"] = nil
	} else if exists && patch["actions"] != nil {
		actionsList := patch["actions"].([]map[string]interface{})
		for i, actionsItem := range actionsList {
			ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionAsPatch(actionsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "actions")
	}
	path = "additional_policy_scopes"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["additional_policy_scopes"] = nil
	} else if !exists {
		delete(patch, "additional_policy_scopes")
	}
	path = "display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if !exists {
		delete(patch, "display_name")
	}
	path = "parent_ids"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["parent_ids"] = nil
	} else if !exists {
		delete(patch, "parent_ids")
	}
	path = "resource_hierarchy_attribute"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["resource_hierarchy_attribute"] = nil
	} else if exists && patch["resource_hierarchy_attribute"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationResourceHierarchyAttributeAsPatch(patch["resource_hierarchy_attribute"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "resource_hierarchy_attribute")
	}
	path = "supported_anonymous_accesses"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_anonymous_accesses"] = nil
	} else if exists && patch["supported_anonymous_accesses"] != nil {
		supported_anonymous_accessesList := patch["supported_anonymous_accesses"].([]map[string]interface{})
		for i, supported_anonymous_accessesItem := range supported_anonymous_accessesList {
			ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessAsPatch(supported_anonymous_accessesItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "supported_anonymous_accesses")
	}
	path = "supported_attributes"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_attributes"] = nil
	} else if exists && patch["supported_attributes"] != nil {
		supported_attributesList := patch["supported_attributes"].([]map[string]interface{})
		for i, supported_attributesItem := range supported_attributesList {
			ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAttributeAsPatch(supported_attributesItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "supported_attributes")
	}
	path = "supported_authorization_subjects"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_authorization_subjects"] = nil
	} else if exists && patch["supported_authorization_subjects"] != nil {
		supported_authorization_subjectsList := patch["supported_authorization_subjects"].([]map[string]interface{})
		for i, supported_authorization_subjectsItem := range supported_authorization_subjectsList {
			ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAuthorizationSubjectAsPatch(supported_authorization_subjectsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "supported_authorization_subjects")
	}
	path = "supported_roles"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_roles"] = nil
	} else if exists && patch["supported_roles"] != nil {
		supported_rolesList := patch["supported_roles"].([]map[string]interface{})
		for i, supported_rolesItem := range supported_rolesList {
			ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedRoleAsPatch(supported_rolesItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "supported_roles")
	}
	path = "supported_network"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_network"] = nil
	} else if exists && patch["supported_network"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedNetworkAsPatch(patch["supported_network"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "supported_network")
	}

	return patch
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedNetworkAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".environment_attributes"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["environment_attributes"] = nil
	} else if exists && patch["environment_attributes"] != nil {
		environment_attributesList := patch["environment_attributes"].([]map[string]interface{})
		for i, environment_attributesItem := range environment_attributesList {
			ResourceIbmOnboardingIamRegistrationEnvironmentAttributeAsPatch(environment_attributesItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "environment_attributes")
	}
}

func ResourceIbmOnboardingIamRegistrationEnvironmentAttributeAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".values"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["values"] = nil
	} else if !exists {
		delete(patch, "values")
	}
	path = rootPath + ".options"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options"] = nil
	} else if exists && patch["options"] != nil {
		ResourceIbmOnboardingIamRegistrationEnvironmentAttributeOptionsAsPatch(patch["options"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "options")
	}
}

func ResourceIbmOnboardingIamRegistrationEnvironmentAttributeOptionsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	} else if !exists {
		delete(patch, "hidden")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedRoleAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["id"] = nil
	} else if !exists {
		delete(patch, "id")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if exists && patch["description"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectAsPatch(patch["description"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if exists && patch["display_name"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectAsPatch(patch["display_name"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "display_name")
	}
	path = rootPath + ".options"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options"] = nil
	} else if exists && patch["options"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedRoleOptionsAsPatch(patch["options"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "options")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedRoleOptionsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".policy_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["policy_type"] = nil
	} else if !exists {
		delete(patch, "policy_type")
	}
	path = rootPath + ".account_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["account_type"] = nil
	} else if !exists {
		delete(patch, "account_type")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAuthorizationSubjectAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".attributes"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["attributes"] = nil
	} else if exists && patch["attributes"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportAuthorizationSubjectAttributeAsPatch(patch["attributes"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "attributes")
	}
	path = rootPath + ".roles"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["roles"] = nil
	} else if !exists {
		delete(patch, "roles")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportAuthorizationSubjectAttributeAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".service_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["service_name"] = nil
	} else if !exists {
		delete(patch, "service_name")
	}
	path = rootPath + ".resource_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["resource_type"] = nil
	} else if !exists {
		delete(patch, "resource_type")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAttributeAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".options"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options"] = nil
	} else if exists && patch["options"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsAsPatch(patch["options"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "options")
	}
	path = rootPath + ".display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if exists && patch["display_name"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectAsPatch(patch["display_name"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "display_name")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if exists && patch["description"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectAsPatch(patch["description"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".ui"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ui"] = nil
	} else if exists && patch["ui"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributeUiAsPatch(patch["ui"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "ui")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".input_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["input_type"] = nil
	} else if !exists {
		delete(patch, "input_type")
	}
	path = rootPath + ".input_details"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["input_details"] = nil
	} else if exists && patch["input_details"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputDetailsAsPatch(patch["input_details"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "input_details")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputDetailsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	} else if !exists {
		delete(patch, "type")
	}
	path = rootPath + ".values"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["values"] = nil
	} else if exists && patch["values"] != nil {
		valuesList := patch["values"].([]map[string]interface{})
		for i, valuesItem := range valuesList {
			ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputValueAsPatch(valuesItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "values")
	}
	path = rootPath + ".gst"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["gst"] = nil
	} else if exists && patch["gst"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputGstAsPatch(patch["gst"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "gst")
	}
	path = rootPath + ".url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["url"] = nil
	} else if exists && patch["url"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputURLAsPatch(patch["url"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "url")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputURLAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".url_endpoint"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["url_endpoint"] = nil
	} else if !exists {
		delete(patch, "url_endpoint")
	}
	path = rootPath + ".input_option_label"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["input_option_label"] = nil
	} else if !exists {
		delete(patch, "input_option_label")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputGstAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".query"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["query"] = nil
	} else if !exists {
		delete(patch, "query")
	}
	path = rootPath + ".value_property_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value_property_name"] = nil
	} else if !exists {
		delete(patch, "value_property_name")
	}
	path = rootPath + ".label_property_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["label_property_name"] = nil
	} else if !exists {
		delete(patch, "label_property_name")
	}
	path = rootPath + ".input_option_label"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["input_option_label"] = nil
	} else if !exists {
		delete(patch, "input_option_label")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributeUiInputValueAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if !exists {
		delete(patch, "value")
	}
	path = rootPath + ".display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if exists && patch["display_name"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectAsPatch(patch["display_name"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "display_name")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".operators"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["operators"] = nil
	} else if !exists {
		delete(patch, "operators")
	}
	path = rootPath + ".hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	} else if !exists {
		delete(patch, "hidden")
	}
	path = rootPath + ".supported_patterns"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["supported_patterns"] = nil
	} else if !exists {
		delete(patch, "supported_patterns")
	}
	path = rootPath + ".policy_types"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["policy_types"] = nil
	} else if !exists {
		delete(patch, "policy_types")
	}
	path = rootPath + ".is_empty_value_supported"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["is_empty_value_supported"] = nil
	} else if !exists {
		delete(patch, "is_empty_value_supported")
	}
	path = rootPath + ".is_string_exists_false_value_supported"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["is_string_exists_false_value_supported"] = nil
	} else if !exists {
		delete(patch, "is_string_exists_false_value_supported")
	}
	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".resource_hierarchy"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["resource_hierarchy"] = nil
	} else if exists && patch["resource_hierarchy"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyAsPatch(patch["resource_hierarchy"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "resource_hierarchy")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if exists && patch["key"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyKeyAsPatch(patch["key"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if exists && patch["value"] != nil {
		ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyValueAsPatch(patch["value"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "value")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyValueAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
}

func ResourceIbmOnboardingIamRegistrationSupportedAttributesOptionsResourceHierarchyKeyAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if !exists {
		delete(patch, "value")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationSupportedAnonymousAccessAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".attributes"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["attributes"] = nil
	} else if !exists {
		delete(patch, "attributes")
	}
	path = rootPath + ".roles"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["roles"] = nil
	} else if !exists {
		delete(patch, "roles")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationResourceHierarchyAttributeAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	} else if !exists {
		delete(patch, "key")
	}
	path = rootPath + ".value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	} else if !exists {
		delete(patch, "value")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["id"] = nil
	} else if !exists {
		delete(patch, "id")
	}
	path = rootPath + ".roles"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["roles"] = nil
	} else if !exists {
		delete(patch, "roles")
	}
	path = rootPath + ".description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	} else if exists && patch["description"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectAsPatch(patch["description"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "description")
	}
	path = rootPath + ".display_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["display_name"] = nil
	} else if exists && patch["display_name"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectAsPatch(patch["display_name"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "display_name")
	}
	path = rootPath + ".options"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["options"] = nil
	} else if exists && patch["options"] != nil {
		ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionOptionsAsPatch(patch["options"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "options")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationActionOptionsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".hidden"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["hidden"] = nil
	} else if !exists {
		delete(patch, "hidden")
	}
}

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDisplayNameObjectAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".default"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["default"] = nil
	} else if !exists {
		delete(patch, "default")
	}
	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
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

func ResourceIbmOnboardingIamRegistrationIamServiceRegistrationDescriptionObjectAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".default"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["default"] = nil
	} else if !exists {
		delete(patch, "default")
	}
	path = rootPath + ".en"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["en"] = nil
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
