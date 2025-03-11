// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmOffering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmOfferingCreate,
		ReadContext:   resourceIBMCmOfferingRead,
		UpdateContext: resourceIBMCmOfferingUpdate,
		DeleteContext: resourceIBMCmOfferingDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Offering identifier.  Provide this when an offering already exists and you wish to use it as a terraform resource.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url for this specific offering.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn for this specific offering.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display Name in the requested language.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The programmatic name of this offering.",
			},
			"offering_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL for an icon associated with this offering.",
			},
			"offering_docs_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL for an additional docs with this offering.",
			},
			"offering_support_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Deprecated:  "This argument is deprecated",
				Description: "[deprecated] - Use offering.support instead.  URL to be displayed in the Consumption UI for getting support on this offering.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of tags associated with this catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"keywords": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of keywords associated with offering, typically used to search for it.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"deprecate": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Deprecate this offering.",
			},
			"rating": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Repository info for offerings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"one_star_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "One start rating.",
						},
						"two_star_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Two start rating.",
						},
						"three_star_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Three start rating.",
						},
						"four_star_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Four start rating.",
						},
					},
				},
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this catalog was created.",
			},
			"updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this catalog was last updated.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Short description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"long_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Long description in the requested language.",
			},
			"long_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "list of features associated with this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Heading.",
						},
						"title_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Feature description.",
						},
						"description_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"kinds": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of kind.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique ID.",
						},
						"format_kind": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "content kind, e.g., helm, vm image.",
						},
						"install_kind": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "install kind, e.g., helm, operator, terraform.",
						},
						"target_kind": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "target cloud to install, e.g., iks, open_shift_iks.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Open ended metadata information.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of tags associated with this catalog.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"additional_features": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of features associated with this offering.",
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
						"created": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The date and time this catalog was created.",
						},
						"updated": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The date and time this catalog was last updated.",
						},
						"versions": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of versions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Unique ID.",
									},
									"rev": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Cloudant revision.",
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Version's CRN.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Version of content type.",
									},
									"flavor": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Version Flavor Information.  Only supported for Product kind Solution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Programmatic name for this flavor.",
												},
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label for this flavor.",
												},
												"label_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"index": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Order that this flavor should appear when listed for a single version.",
												},
											},
										},
									},
									"sha": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "hash of the content.",
									},
									"created": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The date and time this version was created.",
									},
									"updated": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The date and time this version was last updated.",
									},
									"offering_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Offering ID.",
									},
									"catalog_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Catalog ID.",
									},
									"kind_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kind ID.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of tags associated with this catalog.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"repo_url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Content's repo URL.",
									},
									"source_url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Content's source URL (e.g git repo).",
									},
									"tgz_url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "File used to on-board this version.",
									},
									"configuration": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "List of user solicited overrides.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Configuration key.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Value type (string, boolean, int).",
												},
												"default_value": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The default value as a JSON encoded string.  To use a secret when the type is password, specify a JSON encoded value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.",
												},
												"display_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Display name for configuration type.",
												},
												"value_constraint": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Constraint associated with value, e.g., for string type - regx:[a-z].",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Key description.",
												},
												"required": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Is key required to install.",
												},
												"options": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "List of options of type.",
													Elem:        &schema.Schema{Type: schema.TypeMap},
												},
												"hidden": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Hide values.",
												},
												"custom_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Render type.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "ID of the widget type.",
															},
															"grouping": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Determines where this configuration type is rendered (3 sections today - Target, Resource, and Deployment).",
															},
															"original_grouping": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Original grouping type for this configuration (3 types - Target, Resource, and Deployment).",
															},
															"grouping_index": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Determines the order that this configuration item shows in that particular grouping.",
															},
															"config_constraints": &schema.Schema{
																Type:        schema.TypeMap,
																Optional:    true,
																Computed:    true,
																Description: "Map of constraint parameters that will be passed to the custom widget.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"associations": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "List of parameters that are associated with this configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"parameters": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			Description: "Parameters for this association.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Name of this parameter.",
																					},
																					"options_refresh": &schema.Schema{
																						Type:        schema.TypeBool,
																						Optional:    true,
																						Computed:    true,
																						Description: "Refresh options.",
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
												"type_metadata": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The original type, as found in the source being onboarded.",
												},
											},
										},
									},
									"outputs": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of output values for this version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Output key.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Output description.",
												},
											},
										},
									},
									"iam_permissions": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of IAM permissions that are required to consume this version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Service name.",
												},
												"role_crns": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Role CRNs for this permission.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"resources": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Resources for this permission.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Resource name.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Resource description.",
															},
															"role_crns": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Role CRNs for this permission.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"metadata": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Generic data to be included with content being onboarded. Required for virtual server image for VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Version source URL.",
												},
												"working_directory": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Working directory of source files.",
												},
												"example_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Working directory of source files.",
												},
												"start_deploy_time": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time validation started.",
												},
												"end_deploy_time": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The time validation ended.",
												},
												"est_deploy_time": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "The estimated time validation takes.",
												},
												"usage": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Usage text for the version.",
												},
												"usage_template": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Usage text for the version.",
												},
												"modules": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Terraform modules.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the module.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of the module.",
															},
															"offering_reference": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																Description: "Terraform modules.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Name of the offering module.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "ID of the offering module.",
																		},
																		"kind": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Kind of the offeringmodule.",
																		},
																		"version": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Version of the offering module.",
																		},
																		"flavor": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Flavor of the module.",
																		},
																		"flavors": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Flavors of the module.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"catalog_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Catalog ID of the module reference.",
																		},
																		"metadata": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metadata of the module.",
																		},
																	},
																},
															},
														},
													},
												},
												"version_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Version name.",
												},
												"terraform_version": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Terraform version.",
												},
												"validated_terraform_version": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Version name.",
												},
												"vsi_vpc": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "VSI VPC version information",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"operating_system": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Operating system included in this image. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"dedicated_host_only": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
																		},
																		"vendor": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Vendor of the operating system. Required for virtual server image for VPC.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
																		},
																		"href": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "URL for this operating system. Required for virtual server image for VPC.",
																		},
																		"display_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
																		},
																		"family": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Software family for this operating system. Required for virtual server image for VPC.",
																		},
																		"version": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Major release version of this operating system. Required for virtual server image for VPC.",
																		},
																		"architecture": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Operating system architecture. Required for virtual server image for VPC.",
																		},
																	},
																},
															},
															"file": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Details for the stored image file. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"size": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
																		},
																	},
																},
															},
															"minimum_provisioned_size": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
															},
															"images": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Image operating system. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
																		},
																		"region": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Region the virtual server image is available in. Required for virtual server image for VPC.",
																		},
																	},
																},
															},
														},
													},
												},
												"operating_system": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Operating system included in this image. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"dedicated_host_only": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
															},
															"vendor": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Vendor of the operating system. Required for virtual server image for VPC.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "URL for this operating system. Required for virtual server image for VPC.",
															},
															"display_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
															},
															"family": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Software family for this operating system. Required for virtual server image for VPC.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Major release version of this operating system. Required for virtual server image for VPC.",
															},
															"architecture": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Operating system architecture. Required for virtual server image for VPC.",
															},
														},
													},
												},
												"file": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Details for the stored image file. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
															},
														},
													},
												},
												"minimum_provisioned_size": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
												},
												"images": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Image operating system. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
															},
															"region": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Region the virtual server image is available in. Required for virtual server image for VPC.",
															},
														},
													},
												},
											},
										},
									},
									"validation": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Validation response.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"validated": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Date and time of last successful validation.",
												},
												"requested": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Date and time of last validation was requested.",
												},
												"state": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Current validation state - <empty>, in_progress, valid, invalid, expired.",
												},
												"last_operation": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Last operation (e.g. submit_deployment, generate_installer, install_offering.",
												},
												"target": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "Validation target information (e.g. cluster_id, region, namespace, etc).  Values will vary by Content type.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Any message needing to be conveyed as part of the validation job.",
												},
											},
										},
									},
									"required_resources": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Resource requirments for installation.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type of requirement.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "mem, disk, cores, and nodes can be parsed as an int.  targetVersion will be a semver range value.",
												},
											},
										},
									},
									"single_instance": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Denotes if single instance can be deployed to a given cluster.",
									},
									"install": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Script information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instructions": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
												},
												"instructions_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"script": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional script that needs to be run post any pre-condition script.",
												},
												"script_permission": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional iam permissions that are required on the target cluster to run this script.",
												},
												"delete_script": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional script that if run will remove the installed version.",
												},
												"scope": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
												},
											},
										},
									},
									"pre_install": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Optional pre-install instructions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instructions": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
												},
												"instructions_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"script": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional script that needs to be run post any pre-condition script.",
												},
												"script_permission": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional iam permissions that are required on the target cluster to run this script.",
												},
												"delete_script": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional script that if run will remove the installed version.",
												},
												"scope": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
												},
											},
										},
									},
									"entitlement": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Entitlement license info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"provider_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Provider name.",
												},
												"provider_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Provider ID.",
												},
												"product_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Product ID.",
												},
												"part_numbers": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "list of license entitlement part numbers, eg. D1YGZLL,D1ZXILL.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"image_repo_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Image repository name.",
												},
											},
										},
									},
									"licenses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of licenses the product was built with.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "License ID.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "license name.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "type of license e.g., Apache xxx.",
												},
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "URL for the license text.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "License description.",
												},
											},
										},
									},
									"image_manifest_url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If set, denotes a url to a YAML file with list of container images used by this version.",
									},
									"deprecated": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "read only field, indicating if this version is deprecated.",
									},
									"package_version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Version of the package used to create this version.",
									},
									"state": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Offering state.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"current": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
												"current_entered": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Date and time of current request.",
												},
												"pending": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
												"pending_requested": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Date and time of pending request.",
												},
												"previous": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
											},
										},
									},
									"version_locator": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A dotted value of `catalogID`.`versionID`.",
									},
									"long_description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Long description for version.",
									},
									"long_description_i18n": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "A map of translated strings, by language code.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"whitelisted_accounts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Whitelisted accounts for version.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"image_pull_key_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the image pull key to use from Offering.ImagePullKeys.",
									},
									"deprecate_pending": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Optional:    true,
										Description: "Deprecation information for a Version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deprecate_date": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Date of deprecation.",
												},
												"deprecate_state": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Deprecation state.",
												},
												"description": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"solution_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Version Solution Information.  Only supported for Product kind Solution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"architecture_diagrams": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Architecture diagrams for this solution.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"diagram": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Offering Media information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "URL of the specified media item.",
																		},
																		"api_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "CM API specific URL of the specified media item.",
																		},
																		"url_proxy": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Offering URL proxy information.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "URL of the specified media item being proxied.",
																					},
																					"sha": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "SHA256 fingerprint of image.",
																					},
																				},
																			},
																		},
																		"caption": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Caption for this media item.",
																		},
																		"caption_i18n": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "A map of translated strings, by language code.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of this media item.",
																		},
																		"thumbnail_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Thumbnail URL for this media item.",
																		},
																	},
																},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of this diagram.",
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
												"features": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Features - titles only.",
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
												"cost_estimate": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cost estimate definition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cost estimate version.",
															},
															"currency": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Cost estimate currency.",
															},
															"projects": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cost estimate projects.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Project name.",
																		},
																		"metadata": &schema.Schema{
																			Type:        schema.TypeMap,
																			Computed:    true,
																			Description: "Project metadata.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"past_breakdown": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Cost breakdown definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"total_hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total hourly cost.",
																					},
																					"total_monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total monthly cost.",
																					},
																					"resources": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Resources.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Resource name.",
																								},
																								"metadata": &schema.Schema{
																									Type:        schema.TypeMap,
																									Computed:    true,
																									Description: "Resource metadata.",
																									Elem:        &schema.Schema{Type: schema.TypeString},
																								},
																								"hourly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Hourly cost.",
																								},
																								"monthly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Monthly cost.",
																								},
																								"cost_components": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Cost components.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component name.",
																											},
																											"unit": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component unit.",
																											},
																											"hourly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly quantity.",
																											},
																											"monthly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly quantity.",
																											},
																											"price": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component price.",
																											},
																											"hourly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly cost.",
																											},
																											"monthly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly cist.",
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
																		"breakdown": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Cost breakdown definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"total_hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total hourly cost.",
																					},
																					"total_monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total monthly cost.",
																					},
																					"resources": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Resources.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Resource name.",
																								},
																								"metadata": &schema.Schema{
																									Type:        schema.TypeMap,
																									Computed:    true,
																									Description: "Resource metadata.",
																									Elem:        &schema.Schema{Type: schema.TypeString},
																								},
																								"hourly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Hourly cost.",
																								},
																								"monthly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Monthly cost.",
																								},
																								"cost_components": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Cost components.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component name.",
																											},
																											"unit": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component unit.",
																											},
																											"hourly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly quantity.",
																											},
																											"monthly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly quantity.",
																											},
																											"price": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component price.",
																											},
																											"hourly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly cost.",
																											},
																											"monthly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly cist.",
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
																		"diff": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Cost breakdown definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"total_hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total hourly cost.",
																					},
																					"total_monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Total monthly cost.",
																					},
																					"resources": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Resources.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Resource name.",
																								},
																								"metadata": &schema.Schema{
																									Type:        schema.TypeMap,
																									Computed:    true,
																									Description: "Resource metadata.",
																									Elem:        &schema.Schema{Type: schema.TypeString},
																								},
																								"hourly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Hourly cost.",
																								},
																								"monthly_cost": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Monthly cost.",
																								},
																								"cost_components": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Cost components.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component name.",
																											},
																											"unit": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component unit.",
																											},
																											"hourly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly quantity.",
																											},
																											"monthly_quantity": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly quantity.",
																											},
																											"price": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component price.",
																											},
																											"hourly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component hourly cost.",
																											},
																											"monthly_cost": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Cost component monthly cist.",
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
																		"summary": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Cost summary definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"total_detected_resources": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Total detected resources.",
																					},
																					"total_supported_resources": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Total supported resources.",
																					},
																					"total_unsupported_resources": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Total unsupported resources.",
																					},
																					"total_usage_based_resources": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Total usage based resources.",
																					},
																					"total_no_price_resources": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Total no price resources.",
																					},
																					"unsupported_resource_counts": &schema.Schema{
																						Type:        schema.TypeMap,
																						Computed:    true,
																						Description: "Unsupported resource counts.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"no_price_resource_counts": &schema.Schema{
																						Type:        schema.TypeMap,
																						Computed:    true,
																						Description: "No price resource counts.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																				},
																			},
																		},
																	},
																},
															},
															"summary": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cost summary definition.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"total_detected_resources": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Total detected resources.",
																		},
																		"total_supported_resources": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Total supported resources.",
																		},
																		"total_unsupported_resources": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Total unsupported resources.",
																		},
																		"total_usage_based_resources": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Total usage based resources.",
																		},
																		"total_no_price_resources": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Total no price resources.",
																		},
																		"unsupported_resource_counts": &schema.Schema{
																			Type:        schema.TypeMap,
																			Computed:    true,
																			Description: "Unsupported resource counts.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"no_price_resource_counts": &schema.Schema{
																			Type:        schema.TypeMap,
																			Computed:    true,
																			Description: "No price resource counts.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Total hourly cost.",
															},
															"total_monthly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Total monthly cost.",
															},
															"past_total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Past total hourly cost.",
															},
															"past_total_monthly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Past total monthly cost.",
															},
															"diff_total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Difference in total hourly cost.",
															},
															"diff_total_monthly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Difference in total monthly cost.",
															},
															"time_generated": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "When this estimate was generated.",
															},
														},
													},
												},
												"dependencies": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Dependencies for this solution.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"catalog_id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional - If not specified, assumes the Public Catalog.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional - Offering ID - not required if name is set.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional - Programmatic Offering name.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Required - Semver value or range.",
															},
															"flavors": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Optional - List of dependent flavors in the specified range.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"is_consumable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the version able to be shared.",
									},
								},
							},
						},
						"plans": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of plans.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "unique id.",
									},
									"label": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display Name in the requested language.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The programmatic name of this offering.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Short description in the requested language.",
									},
									"long_description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Long description in the requested language.",
									},
									"metadata": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "open ended metadata information.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "list of tags associated with this catalog.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"additional_features": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "list of features associated with this offering.",
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
									"created": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "the date'time this catalog was created.",
									},
									"updated": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "the date'time this catalog was last updated.",
									},
									"deployments": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "list of deployments.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "unique id.",
												},
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Display Name in the requested language.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The programmatic name of this offering.",
												},
												"short_description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Short description in the requested language.",
												},
												"long_description": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Long description in the requested language.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "open ended metadata information.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"tags": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "list of tags associated with this catalog.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"created": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "the date'time this catalog was created.",
												},
												"updated": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "the date'time this catalog was last updated.",
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
			"pc_managed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Offering is managed by Partner Center.",
			},
			"publish_approved": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Offering has been approved to publish to permitted to IBM or Public Catalog.",
			},
			"share_with_all": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Denotes public availability of an Offering - if share_enabled is true.",
			},
			"share_with_ibm": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Denotes IBM employee availability of an Offering - if share_enabled is true.",
			},
			"share_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Denotes sharing including access list availability of an Offering is enabled.",
			},
			"share_with_access_list": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of account IDs to add to this offering's access list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"public_original_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The original offering CRN that this publish entry came from.",
			},
			"publish_public_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The crn of the public catalog entry of this offering.",
			},
			"portal_approval_record": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The portal's approval record ID.",
			},
			"portal_ui_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The portal UI URL.",
			},
			"catalog_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the catalog.",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Map of metadata values for this offering.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"disclaimer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A disclaimer for this offering.",
			},
			"hidden": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determine if this offering should be displayed in the Consumption UI.",
			},
			// "provider": &schema.Schema{
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Deprecated:  "This argument is deprecated",
			// 	Description: "Deprecated - Provider of this offering.",
			// },
			"provider_info": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Information on the provider for this offering, or omitted if no provider information is given.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The id of this provider.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of this provider.",
						},
					},
				},
			},
			"repo_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Repository info for offerings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Token for private repos.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public or enterprise GitHub.",
						},
					},
				},
			},
			"image_pull_keys": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Image pull keys for this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Key name.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Key value.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Key description.",
						},
					},
				},
			},
			"support": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Offering Support information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL to be displayed in the Consumption UI for getting support on this offering.",
						},
						"process": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Support process as provided by an ISV.",
						},
						"process_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"locations": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of country codes indicating where support is provided.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"support_details": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of support options (e.g. email, phone, slack, other).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the current support detail.",
									},
									"contact": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact for the current support detail.",
									},
									"response_wait_time": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Time descriptor.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Amount of time to wait in unit 'type'.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Valid values are hour or day.",
												},
											},
										},
									},
									"availability": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Times when support is available.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"times": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A list of support times.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The day of the week, represented as an integer.",
															},
															"start_time": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "HOURS:MINUTES:SECONDS using 24 hour time (e.g. 8:15:00).",
															},
															"end_time": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "HOURS:MINUTES:SECONDS using 24 hour time (e.g. 8:15:00).",
															},
														},
													},
												},
												"timezone": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Timezone (e.g. America/New_York).",
												},
												"always_available": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is this support always available.",
												},
											},
										},
									},
								},
							},
						},
						"support_escalation": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Support escalation policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"escalation_wait_time": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Time descriptor.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Amount of time to wait in unit 'type'.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Valid values are hour or day.",
												},
											},
										},
									},
									"response_wait_time": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Time descriptor.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Amount of time to wait in unit 'type'.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Valid values are hour or day.",
												},
											},
										},
									},
									"contact": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Escalation contact.",
									},
								},
							},
						},
						"support_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Support type for this product.",
						},
					},
				},
			},
			"media": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of media items related to this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "URL of the specified media item.",
						},
						"api_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "CM API specific URL of the specified media item.",
						},
						"url_proxy": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Offering URL proxy information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "URL of the specified media item being proxied.",
									},
									"sha": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "SHA256 fingerprint of image.",
									},
								},
							},
						},
						"caption": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Caption for this media item.",
						},
						"caption_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Type of this media item.",
						},
						"thumbnail_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Thumbnail URL for this media item.",
						},
					},
				},
			},
			"deprecate_pending": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Deprecation information for an Offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deprecate_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date of deprecation.",
						},
						"deprecate_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deprecation state.",
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"product_kind": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The product kind.  Valid values are module, solution, or empty string.",
			},
			"badges": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of badges for this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the current badge.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name for the current badge.",
						},
						"label_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the current badge.",
						},
						"description_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"icon": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Icon for the current badge.",
						},
						"authority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authority for the current badge.",
						},
						"tag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag for the current badge.",
						},
						"learn_more_links": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Learn more links for a badge.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_party": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "First party link.",
									},
									"third_party": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Third party link.",
									},
								},
							},
						},
						"constraints": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An optional set of constraints indicating which versions in an Offering have this particular badge.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the current constraint.",
									},
									"rule": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule for the current constraint.",
									},
								},
							},
						},
					},
				},
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"offering_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed Offering ID.",
			},
		},
	}
}

func resourceIBMCmOfferingCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if _, ok := d.GetOk("offering_id"); ok {
		getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}
		getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
		getOfferingOptions.SetOfferingID(d.Get("offering_id").(string))

		// Don't need to get full offering with all versions here because we call resourceIBMCmOfferingRead which does
		// This get is just to see if the offering exists and set state id
		offering, response, err := catalogManagementClient.GetOfferingWithContext(context, getOfferingOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		d.SetId(*offering.ID)

		err = handleShareOfferingAfterCreate(catalogManagementClient, *offering, d, context)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		return resourceIBMCmOfferingRead(context, d, meta)
	}

	createOfferingOptions := &catalogmanagementv1.CreateOfferingOptions{}

	createOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	if _, ok := d.GetOk("url"); ok {
		createOfferingOptions.SetURL(d.Get("url").(string))
	}
	if _, ok := d.GetOk("crn"); ok {
		createOfferingOptions.SetCRN(d.Get("crn").(string))
	}
	if _, ok := d.GetOk("label"); ok {
		createOfferingOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createOfferingOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("offering_icon_url"); ok {
		createOfferingOptions.SetOfferingIconURL(d.Get("offering_icon_url").(string))
	}
	if _, ok := d.GetOk("offering_docs_url"); ok {
		createOfferingOptions.SetOfferingDocsURL(d.Get("offering_docs_url").(string))
	}
	if _, ok := d.GetOk("offering_support_url"); ok {
		createOfferingOptions.SetOfferingSupportURL(d.Get("offering_support_url").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createOfferingOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("keywords"); ok {
		createOfferingOptions.SetKeywords(SIToSS(d.Get("keywords").([]interface{})))
	}
	if _, ok := d.GetOk("rating"); ok {
		ratingModel, err := resourceIBMCmOfferingMapToRating(d.Get("rating.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetRating(ratingModel)
	}
	if _, ok := d.GetOk("created"); ok {
		fmtDateTimeCreated, err := core.ParseDateTime(d.Get("created").(string))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetCreated(&fmtDateTimeCreated)
	}
	if _, ok := d.GetOk("updated"); ok {
		fmtDateTimeUpdated, err := core.ParseDateTime(d.Get("updated").(string))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetUpdated(&fmtDateTimeUpdated)
	}
	if _, ok := d.GetOk("short_description"); ok {
		createOfferingOptions.SetShortDescription(d.Get("short_description").(string))
	}
	if _, ok := d.GetOk("long_description"); ok {
		createOfferingOptions.SetLongDescription(d.Get("long_description").(string))
	}
	if _, ok := d.GetOk("features"); ok {
		var features []catalogmanagementv1.Feature
		for _, e := range d.Get("features").([]interface{}) {
			value := e.(map[string]interface{})
			featuresItem, err := resourceIBMCmOfferingMapToFeature(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			features = append(features, *featuresItem)
		}
		createOfferingOptions.SetFeatures(features)
	}
	if _, ok := d.GetOk("kinds"); ok {
		var kinds []catalogmanagementv1.Kind
		for _, e := range d.Get("kinds").([]interface{}) {
			value := e.(map[string]interface{})
			kindsItem, err := resourceIBMCmOfferingMapToKind(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			kinds = append(kinds, *kindsItem)
		}
		createOfferingOptions.SetKinds(kinds)
	}
	if _, ok := d.GetOk("pc_managed"); ok {
		createOfferingOptions.SetPcManaged(d.Get("pc_managed").(bool))
	}
	if _, ok := d.GetOk("publish_approved"); ok {
		createOfferingOptions.SetPublishApproved(d.Get("publish_approved").(bool))
	}
	if _, ok := d.GetOk("share_with_all"); ok {
		createOfferingOptions.SetShareWithAll(d.Get("share_with_all").(bool))
	}
	if _, ok := d.GetOk("share_with_ibm"); ok {
		createOfferingOptions.SetShareWithIBM(d.Get("share_with_ibm").(bool))
	}
	if _, ok := d.GetOk("share_enabled"); ok {
		createOfferingOptions.SetShareEnabled(d.Get("share_enabled").(bool))
	}
	if _, ok := d.GetOk("public_original_crn"); ok {
		createOfferingOptions.SetPublicOriginalCRN(d.Get("public_original_crn").(string))
	}
	if _, ok := d.GetOk("publish_public_crn"); ok {
		createOfferingOptions.SetPublishPublicCRN(d.Get("publish_public_crn").(string))
	}
	if _, ok := d.GetOk("portal_approval_record"); ok {
		createOfferingOptions.SetPortalApprovalRecord(d.Get("portal_approval_record").(string))
	}
	if _, ok := d.GetOk("portal_ui_url"); ok {
		createOfferingOptions.SetPortalUIURL(d.Get("portal_ui_url").(string))
	}
	if _, ok := d.GetOk("catalog_id"); ok {
		createOfferingOptions.SetCatalogID(d.Get("catalog_id").(string))
	}
	if _, ok := d.GetOk("catalog_name"); ok {
		createOfferingOptions.SetCatalogName(d.Get("catalog_name").(string))
	}
	if _, ok := d.GetOk("disclaimer"); ok {
		createOfferingOptions.SetDisclaimer(d.Get("disclaimer").(string))
	}
	if _, ok := d.GetOk("hidden"); ok {
		createOfferingOptions.SetHidden(d.Get("hidden").(bool))
	}
	if _, ok := d.GetOk("provider_info"); ok {
		providerInfoModel, err := resourceIBMCmOfferingMapToProviderInfo(d.Get("provider_info.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetProviderInfo(providerInfoModel)
	}
	if _, ok := d.GetOk("repo_info"); ok {
		repoInfoModel, err := resourceIBMCmOfferingMapToRepoInfo(d.Get("repo_info.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetRepoInfo(repoInfoModel)
	}
	if _, ok := d.GetOk("image_pull_keys"); ok {
		var imagePullKeys []catalogmanagementv1.ImagePullKey
		for _, e := range d.Get("image_pull_keys").([]interface{}) {
			value := e.(map[string]interface{})
			imagePullKeysItem, err := resourceIBMCmOfferingMapToImagePullKey(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			imagePullKeys = append(imagePullKeys, *imagePullKeysItem)
		}
		createOfferingOptions.SetImagePullKeys(imagePullKeys)
	}
	if _, ok := d.GetOk("support"); ok {
		supportModel, err := resourceIBMCmOfferingMapToSupport(d.Get("support.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createOfferingOptions.SetSupport(supportModel)
	}
	if _, ok := d.GetOk("media"); ok {
		var media []catalogmanagementv1.MediaItem
		for _, e := range d.Get("media").([]interface{}) {
			value := e.(map[string]interface{})
			mediaItem, err := resourceIBMCmOfferingMapToMediaItem(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			media = append(media, *mediaItem)
		}
		createOfferingOptions.SetMedia(media)
	}
	if _, ok := d.GetOk("product_kind"); ok {
		createOfferingOptions.SetProductKind(d.Get("product_kind").(string))
	}
	if _, ok := d.GetOk("badges"); ok {
		var badges []catalogmanagementv1.Badge
		for _, e := range d.Get("badges").([]interface{}) {
			value := e.(map[string]interface{})
			badgesItem, err := resourceIBMCmOfferingMapToBadge(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			badges = append(badges, *badgesItem)
		}
		createOfferingOptions.SetBadges(badges)
	}

	offering, response, err := catalogManagementClient.CreateOfferingWithContext(context, createOfferingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*offering.ID)

	err = handleShareOfferingAfterCreate(catalogManagementClient, *offering, d, context)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMCmOfferingRead(context, d, meta)
}

func resourceIBMCmOfferingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}

	getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getOfferingOptions.SetOfferingID(d.Id())

	offering, response, err := FetchOfferingWithAllVersions(context, catalogManagementClient, getOfferingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_id", getOfferingOptions.CatalogIdentifier); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("url", offering.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", offering.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("label", offering.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("label_i18n", offering.LabelI18n); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label_i18n: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", offering.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("offering_icon_url", offering.OfferingIconURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_icon_url: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("offering_docs_url", offering.OfferingDocsURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_docs_url: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("offering_support_url", offering.OfferingSupportURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_support_url: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	tags := []string{}
	if offering.Tags != nil {
		specifiedTags := SIToSS(d.Get("tags").([]interface{}))
		for _, specifiedTag := range specifiedTags {
			for _, offeringTag := range offering.Tags {
				if offeringTag == specifiedTag {
					tags = append(tags, specifiedTag)
					break
				}
			}
		}
	}
	if err = d.Set("tags", tags); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting tags: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	keywords := []string{}
	if offering.Keywords != nil {
		keywords = offering.Keywords
	}
	if err = d.Set("keywords", keywords); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting keywords: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("metadata", offering.Metadata); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering metadata: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if offering.Rating != nil {
		ratingMap, err := resourceIBMCmOfferingRatingToMap(offering.Rating)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("rating", []map[string]interface{}{ratingMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rating: %s", err), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("created", flex.DateTimeToString(offering.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated", flex.DateTimeToString(offering.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("short_description", offering.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("short_description_i18n", offering.ShortDescriptionI18n); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description_i18n: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("long_description", offering.LongDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("long_description_i18n", offering.LongDescriptionI18n); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description_i18n: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	features := []map[string]interface{}{}
	if offering.Features != nil {
		for _, featuresItem := range offering.Features {
			featuresItemMap, err := resourceIBMCmOfferingFeatureToMap(&featuresItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			features = append(features, featuresItemMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting features: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	kinds := []map[string]interface{}{}
	if offering.Kinds != nil {
		for _, kindsItem := range offering.Kinds {
			kindsItemMap, err := resourceIBMCmOfferingKindToMap(&kindsItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			kinds = append(kinds, kindsItemMap)
		}
	}
	if err = d.Set("kinds", kinds); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kinds: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("pc_managed", offering.PcManaged); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting pc_managed: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("publish_approved", offering.PublishApproved); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish_approved: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("share_with_all", offering.ShareWithAll); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_with_all: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("share_with_ibm", offering.ShareWithIBM); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_with_ibm: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("share_enabled", offering.ShareEnabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_enabled: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("public_original_crn", offering.PublicOriginalCRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting public_original_crn: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("publish_public_crn", offering.PublishPublicCRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish_public_crn: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("portal_approval_record", offering.PortalApprovalRecord); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting portal_approval_record: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("portal_ui_url", offering.PortalUIURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting portal_ui_url: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("catalog_id", offering.CatalogID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("catalog_name", offering.CatalogName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_name: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("disclaimer", offering.Disclaimer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disclaimer: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("hidden", offering.Hidden); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting hidden: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if offering.ProviderInfo != nil {
		providerInfoMap, err := resourceIBMCmOfferingProviderInfoToMap(offering.ProviderInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("provider_info", []map[string]interface{}{providerInfoMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting provider_info: %s", err), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if offering.RepoInfo != nil {
		repoInfoMap, err := resourceIBMCmOfferingRepoInfoToMap(offering.RepoInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("repo_info", []map[string]interface{}{repoInfoMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting repo_info: %s", err), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	imagePullKeys := []map[string]interface{}{}
	if offering.ImagePullKeys != nil {
		for _, imagePullKeysItem := range offering.ImagePullKeys {
			imagePullKeysItemMap, err := resourceIBMCmOfferingImagePullKeyToMap(&imagePullKeysItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			imagePullKeys = append(imagePullKeys, imagePullKeysItemMap)
		}
	}
	if err = d.Set("image_pull_keys", imagePullKeys); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_pull_keys: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if offering.Support != nil {
		supportMap, err := resourceIBMCmOfferingSupportToMap(offering.Support)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("support", []map[string]interface{}{supportMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting support: %s", err), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	media := []map[string]interface{}{}
	if offering.Media != nil {
		for _, mediaItem := range offering.Media {
			mediaItemMap, err := resourceIBMCmOfferingMediaItemToMap(&mediaItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			media = append(media, mediaItemMap)
		}
	}
	if err = d.Set("media", media); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting media: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	deprecatePendingMap := make(map[string]interface{})
	if offering.DeprecatePending != nil {
		deprecatePendingMap, err = resourceIBMCmOfferingDeprecatePendingToMap(offering.DeprecatePending)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("deprecate_pending", []map[string]interface{}{deprecatePendingMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting deprecate_pending: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("product_kind", offering.ProductKind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting product_kind: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	badges := []map[string]interface{}{}
	if offering.Badges != nil {
		for _, badgesItem := range offering.Badges {
			badgesItemMap, err := resourceIBMCmOfferingBadgeToMap(&badgesItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			badges = append(badges, badgesItemMap)
		}
	}
	if err = d.Set("badges", badges); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting badges: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("rev", offering.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("offering_identifier", offering.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_id: %s", err), "ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIBMCmOfferingUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateOfferingOptions := &catalogmanagementv1.UpdateOfferingOptions{}

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}

	getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getOfferingOptions.SetOfferingID(d.Id())
	offering, response, err := catalogManagementClient.GetOfferingWithContext(context, getOfferingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
	updateOfferingOptions.SetOfferingID(*offering.ID)
	ifMatch := fmt.Sprintf("\"%s\"", *offering.Rev)
	updateOfferingOptions.IfMatch = &ifMatch

	hasChange := false

	if d.HasChange("label") {
		var method string
		if offering.Label == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/label"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("label"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("label_i18n") {
		var method string
		if offering.LabelI18n == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/label_i18n"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("label_i18n"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("name") {
		var method string
		if offering.Name == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/name"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("name"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("offering_icon_url") {
		var method string
		if offering.OfferingIconURL == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/offering_icon_url"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("offering_icon_url"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("offering_docs_url") {
		var method string
		if offering.OfferingDocsURL == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/offering_docs_url"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("offering_docs_url"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("offering_support_url") {
		var method string
		if offering.OfferingSupportURL == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/offering_support_url"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("offering_support_url"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("tags") {
		var method string
		if offering.Tags == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/tags"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("tags"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("keywords") {
		var method string
		if offering.Keywords == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/keywords"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("keywords"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("short_description") {
		var method string
		if offering.ShortDescription == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/short_description"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("short_description"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("short_description_i18n") {
		var method string
		if offering.ShortDescriptionI18n == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/short_description_i18n"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("short_description_i18n"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("long_description") {
		var method string
		if offering.LongDescription == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/long_description"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("long_description"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("long_description_i18n") {
		var method string
		if offering.LongDescriptionI18n == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/long_description_i18n"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("long_description_i18n"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("features") {
		var method string
		if offering.Features == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/features"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("features"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("kinds") {
		var method string
		if offering.Kinds == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/kinds"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("kinds"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("public_original_crn") {
		var method string
		if offering.PublicOriginalCRN == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/public_original_crn"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("public_original_crn"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("publish_public_crn") {
		var method string
		if offering.PublishPublicCRN == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/publish_public_crn"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("publish_public_crn"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("portal_approval_record") {
		var method string
		if offering.PortalApprovalRecord == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/portal_approval_record"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("portal_approval_record"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("portal_ui_url") {
		var method string
		if offering.PortalUIURL == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/portal_ui_url"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("portal_ui_url"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("metadata") {
		var method string
		if offering.Metadata == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/metadata"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("metadata"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("disclaimer") {
		var method string
		if offering.Disclaimer == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/disclaimer"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("disclaimer"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("hidden") {
		var method string
		if offering.Hidden == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/hidden"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("hidden"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("provider") {
		var method string
		if offering.Provider == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/provider"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("provider"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("provider_info") {
		var method string
		if offering.ProviderInfo == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/provider_info"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("provider_info.0"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("image_pull_keys") {
		var method string
		if offering.ImagePullKeys == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/image_pull_keys"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("image_pull_keys"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("support") {
		var method string
		if offering.Support == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/support"
		supportModel, err := resourceIBMCmOfferingMapToSupport(d.Get("support.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: supportModel,
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("media") {
		var method string
		if offering.Media == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/media"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("media"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("product_kind") {
		var method string
		if offering.ProductKind == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := "/product_kind"
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("product_kind"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}

	if d.HasChange("share_with_access_list") {
		// Find removed accounts, ones that are present in old list but not in new list
		oldAccountList, newAccountList := d.GetChange("share_with_access_list")
		accountsToRemove := make([]string, 0)
		for _, oldAccount := range oldAccountList.([]interface{}) {
			presentInNewList := false
			for _, newAccount := range newAccountList.([]interface{}) {
				if newAccount.(string) == oldAccount.(string) {
					// found id in both lists, do not remove it
					presentInNewList = true
					break
				}
			}
			if !presentInNewList {
				accountsToRemove = append(accountsToRemove, oldAccount.(string))
			}
		}
		// Delete accounts from access list that are no longer specified, if there are any to delete
		if len(accountsToRemove) > 0 {
			deleteOfferingAccessListOptions := catalogmanagementv1.DeleteOfferingAccessListOptions{}
			deleteOfferingAccessListOptions.SetCatalogIdentifier(*offering.CatalogID)
			deleteOfferingAccessListOptions.SetOfferingID(*offering.ID)
			deleteOfferingAccessListOptions.SetAccesses(accountsToRemove)
			_, response, err = catalogManagementClient.DeleteOfferingAccessListWithContext(context, &deleteOfferingAccessListOptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteOfferingAccessListWithContext failed %s\n%s", err, response), "ibm_cm_offering", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}

		// Share with accounts if there are any
		if len(newAccountList.([]interface{})) > 0 {
			addOfferingAccessListOptions := catalogmanagementv1.AddOfferingAccessListOptions{}
			addOfferingAccessListOptions.SetCatalogIdentifier(*offering.CatalogID)
			addOfferingAccessListOptions.SetOfferingID(*offering.ID)
			addOfferingAccessListOptions.SetAccesses(SIToSS(newAccountList.([]interface{})))
			_, response, err = catalogManagementClient.AddOfferingAccessListWithContext(context, &addOfferingAccessListOptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("AddOfferingAccessListWithContext failed %s\n%s", err, response), "ibm_cm_offering", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
	}

	publishStatusChanged := false
	shareOfferingOptions := catalogmanagementv1.ShareOfferingOptions{}
	shareOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
	shareOfferingOptions.SetOfferingID(*offering.ID)
	shareOfferingOptions.SetEnabled(d.Get("share_enabled").(bool))

	if d.HasChange("share_enabled") {
		publishStatusChanged = true
	}

	if d.HasChange("share_with_ibm") {
		publishStatusChanged = true
		shareOfferingOptions.SetIBM(d.Get("share_with_ibm").(bool))
	}

	if d.HasChange("share_with_all") {
		publishStatusChanged = true
		shareOfferingOptions.SetPublic(d.Get("share_with_all").(bool))
	}

	if publishStatusChanged {
		_, response, err = catalogManagementClient.ShareOfferingWithContext(context, &shareOfferingOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ShareOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if d.HasChange("deprecate") && d.Get("deprecate") != nil {
		deprecateOfferingOptions := &catalogmanagementv1.DeprecateOfferingOptions{}
		deprecateOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
		deprecateOfferingOptions.SetOfferingID(*offering.ID)
		deprecateOfferingOptions.SetSetting(strconv.FormatBool(d.Get("deprecate").(bool)))

		response, err := catalogManagementClient.DeprecateOfferingWithContext(context, deprecateOfferingOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateOfferingWithContext failed trying to deprecate offering %s\n%s", err, response), "ibm_cm_offering", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if hasChange {
		_, response, err := catalogManagementClient.UpdateOfferingWithContext(context, updateOfferingOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCmOfferingRead(context, d, meta)
}

func resourceIBMCmOfferingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_offering", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteOfferingOptions := &catalogmanagementv1.DeleteOfferingOptions{}

	deleteOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	deleteOfferingOptions.SetOfferingID(d.Id())

	response, err := catalogManagementClient.DeleteOfferingWithContext(context, deleteOfferingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteOfferingWithContext failed %s\n%s", err, response), "ibm_cm_offering", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIBMCmOfferingMapToRating(modelMap map[string]interface{}) (*catalogmanagementv1.Rating, error) {
	model := &catalogmanagementv1.Rating{}
	if modelMap["one_star_count"] != nil {
		model.OneStarCount = core.Int64Ptr(int64(modelMap["one_star_count"].(int)))
	}
	if modelMap["two_star_count"] != nil {
		model.TwoStarCount = core.Int64Ptr(int64(modelMap["two_star_count"].(int)))
	}
	if modelMap["three_star_count"] != nil {
		model.ThreeStarCount = core.Int64Ptr(int64(modelMap["three_star_count"].(int)))
	}
	if modelMap["four_star_count"] != nil {
		model.FourStarCount = core.Int64Ptr(int64(modelMap["four_star_count"].(int)))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToFeature(modelMap map[string]interface{}) (*catalogmanagementv1.Feature, error) {
	model := &catalogmanagementv1.Feature{}
	if modelMap["title"] != nil && modelMap["title"].(string) != "" {
		model.Title = core.StringPtr(modelMap["title"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToKind(modelMap map[string]interface{}) (*catalogmanagementv1.Kind, error) {
	model := &catalogmanagementv1.Kind{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["format_kind"] != nil && modelMap["format_kind"].(string) != "" {
		model.FormatKind = core.StringPtr(modelMap["format_kind"].(string))
	}
	if modelMap["install_kind"] != nil && modelMap["install_kind"].(string) != "" {
		model.InstallKind = core.StringPtr(modelMap["install_kind"].(string))
	}
	if modelMap["target_kind"] != nil && modelMap["target_kind"].(string) != "" {
		model.TargetKind = core.StringPtr(modelMap["target_kind"].(string))
	}
	if modelMap["tags"] != nil {
		tags := []string{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		model.Tags = tags
	}
	if modelMap["additional_features"] != nil {
		additionalFeatures := []catalogmanagementv1.Feature{}
		for _, additionalFeaturesItem := range modelMap["additional_features"].([]interface{}) {
			additionalFeaturesItemModel, err := resourceIBMCmOfferingMapToFeature(additionalFeaturesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			additionalFeatures = append(additionalFeatures, *additionalFeaturesItemModel)
		}
		model.AdditionalFeatures = additionalFeatures
	}
	if modelMap["created"] != nil {

	}
	if modelMap["updated"] != nil {

	}
	if modelMap["versions"] != nil {
		versions := []catalogmanagementv1.Version{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versionsItemModel, err := resourceIBMCmOfferingMapToVersion(versionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			versions = append(versions, *versionsItemModel)
		}
		model.Versions = versions
	}
	return model, nil
}

func resourceIBMCmOfferingMapToVersion(modelMap map[string]interface{}) (*catalogmanagementv1.Version, error) {
	model := &catalogmanagementv1.Version{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["rev"] != nil && modelMap["rev"].(string) != "" {
		model.Rev = core.StringPtr(modelMap["rev"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	if modelMap["flavor"] != nil && len(modelMap["flavor"].([]interface{})) > 0 {
		FlavorModel, err := resourceIBMCmOfferingMapToFlavor(modelMap["flavor"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Flavor = FlavorModel
	}
	if modelMap["sha"] != nil && modelMap["sha"].(string) != "" {
		model.Sha = core.StringPtr(modelMap["sha"].(string))
	}
	if modelMap["created"] != nil {

	}
	if modelMap["updated"] != nil {

	}
	if modelMap["offering_id"] != nil && modelMap["offering_id"].(string) != "" {
		model.OfferingID = core.StringPtr(modelMap["offering_id"].(string))
	}
	if modelMap["catalog_id"] != nil && modelMap["catalog_id"].(string) != "" {
		model.CatalogID = core.StringPtr(modelMap["catalog_id"].(string))
	}
	if modelMap["kind_id"] != nil && modelMap["kind_id"].(string) != "" {
		model.KindID = core.StringPtr(modelMap["kind_id"].(string))
	}
	if modelMap["tags"] != nil {
		tags := []string{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		model.Tags = tags
	}
	if modelMap["repo_url"] != nil && modelMap["repo_url"].(string) != "" {
		model.RepoURL = core.StringPtr(modelMap["repo_url"].(string))
	}
	if modelMap["source_url"] != nil && modelMap["source_url"].(string) != "" {
		model.SourceURL = core.StringPtr(modelMap["source_url"].(string))
	}
	if modelMap["tgz_url"] != nil && modelMap["tgz_url"].(string) != "" {
		model.TgzURL = core.StringPtr(modelMap["tgz_url"].(string))
	}
	if modelMap["configuration"] != nil {
		configuration := []catalogmanagementv1.Configuration{}
		for _, configurationItem := range modelMap["configuration"].([]interface{}) {
			configurationItemModel, err := resourceIBMCmOfferingMapToConfiguration(configurationItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			configuration = append(configuration, *configurationItemModel)
		}
		model.Configuration = configuration
	}
	if modelMap["outputs"] != nil {
		outputs := []catalogmanagementv1.Output{}
		for _, outputsItem := range modelMap["outputs"].([]interface{}) {
			outputsItemModel, err := resourceIBMCmOfferingMapToOutput(outputsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			outputs = append(outputs, *outputsItemModel)
		}
		model.Outputs = outputs
	}
	if modelMap["iam_permissions"] != nil {
		iamPermissions := []catalogmanagementv1.IamPermission{}
		for _, iamPermissionsItem := range modelMap["iam_permissions"].([]interface{}) {
			iamPermissionsItemModel, err := resourceIBMCmOfferingMapToIamPermission(iamPermissionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			iamPermissions = append(iamPermissions, *iamPermissionsItemModel)
		}
		model.IamPermissions = iamPermissions
	}
	if modelMap["metadata"] != nil {
		metadataModel := modelMap["metadata"].([]interface{})[0].(map[string]interface{})
		model.Metadata = metadataModel
		// FlavorModel, err := resourceIBMCmOfferingMapToFlavor(modelMap["flavor"].([]interface{})[0].(map[string]interface{}))
		// if err != nil {
		// 	return model, err
		// }
		// model.Flavor = FlavorModel
	}
	if modelMap["validation"] != nil && len(modelMap["validation"].([]interface{})) > 0 {
		ValidationModel, err := resourceIBMCmOfferingMapToValidation(modelMap["validation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Validation = ValidationModel
	}
	if modelMap["required_resources"] != nil {
		requiredResources := []catalogmanagementv1.Resource{}
		for _, requiredResourcesItem := range modelMap["required_resources"].([]interface{}) {
			requiredResourcesItemModel, err := resourceIBMCmOfferingMapToResource(requiredResourcesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			requiredResources = append(requiredResources, *requiredResourcesItemModel)
		}
		model.RequiredResources = requiredResources
	}
	if modelMap["single_instance"] != nil {
		model.SingleInstance = core.BoolPtr(modelMap["single_instance"].(bool))
	}
	if modelMap["install"] != nil && len(modelMap["install"].([]interface{})) > 0 {
		InstallModel, err := resourceIBMCmOfferingMapToScript(modelMap["install"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Install = InstallModel
	}
	if modelMap["pre_install"] != nil {
		preInstall := []catalogmanagementv1.Script{}
		for _, preInstallItem := range modelMap["pre_install"].([]interface{}) {
			preInstallItemModel, err := resourceIBMCmOfferingMapToScript(preInstallItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			preInstall = append(preInstall, *preInstallItemModel)
		}
		model.PreInstall = preInstall
	}
	if modelMap["entitlement"] != nil && len(modelMap["entitlement"].([]interface{})) > 0 {
		EntitlementModel, err := resourceIBMCmOfferingMapToVersionEntitlement(modelMap["entitlement"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Entitlement = EntitlementModel
	}
	if modelMap["licenses"] != nil {
		licenses := []catalogmanagementv1.License{}
		for _, licensesItem := range modelMap["licenses"].([]interface{}) {
			licensesItemModel, err := resourceIBMCmOfferingMapToLicense(licensesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			licenses = append(licenses, *licensesItemModel)
		}
		model.Licenses = licenses
	}
	if modelMap["image_manifest_url"] != nil && modelMap["image_manifest_url"].(string) != "" {
		model.ImageManifestURL = core.StringPtr(modelMap["image_manifest_url"].(string))
	}
	if modelMap["deprecated"] != nil {
		model.Deprecated = core.BoolPtr(modelMap["deprecated"].(bool))
	}
	if modelMap["package_version"] != nil && modelMap["package_version"].(string) != "" {
		model.PackageVersion = core.StringPtr(modelMap["package_version"].(string))
	}
	if modelMap["state"] != nil && len(modelMap["state"].([]interface{})) > 0 {
		StateModel, err := resourceIBMCmOfferingMapToState(modelMap["state"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.State = StateModel
	}
	if modelMap["version_locator"] != nil && modelMap["version_locator"].(string) != "" {
		model.VersionLocator = core.StringPtr(modelMap["version_locator"].(string))
	}
	if modelMap["long_description"] != nil && modelMap["long_description"].(string) != "" {
		model.LongDescription = core.StringPtr(modelMap["long_description"].(string))
	}
	if modelMap["whitelisted_accounts"] != nil {
		whitelistedAccounts := []string{}
		for _, whitelistedAccountsItem := range modelMap["whitelisted_accounts"].([]interface{}) {
			whitelistedAccounts = append(whitelistedAccounts, whitelistedAccountsItem.(string))
		}
		model.WhitelistedAccounts = whitelistedAccounts
	}
	if modelMap["image_pull_key_name"] != nil && modelMap["image_pull_key_name"].(string) != "" {
		model.ImagePullKeyName = core.StringPtr(modelMap["image_pull_key_name"].(string))
	}
	if modelMap["deprecate_pending"] != nil && len(modelMap["deprecate_pending"].([]interface{})) > 0 {
		DeprecatePendingModel, err := resourceIBMCmOfferingMapToDeprecatePending(modelMap["deprecate_pending"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DeprecatePending = DeprecatePendingModel
	}
	if modelMap["solution_info"] != nil && len(modelMap["solution_info"].([]interface{})) > 0 {
		SolutionInfoModel, err := resourceIBMCmOfferingMapToSolutionInfo(modelMap["solution_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SolutionInfo = SolutionInfoModel
	}
	if modelMap["is_consumable"] != nil {
		model.IsConsumable = core.BoolPtr(modelMap["is_consumable"].(bool))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToFlavor(modelMap map[string]interface{}) (*catalogmanagementv1.Flavor, error) {
	model := &catalogmanagementv1.Flavor{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["index"] != nil {
		model.Index = core.Int64Ptr(int64(modelMap["index"].(int)))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToConfiguration(modelMap map[string]interface{}) (*catalogmanagementv1.Configuration, error) {
	model := &catalogmanagementv1.Configuration{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["default_value"] != nil {

	}
	if modelMap["display_name"] != nil && modelMap["display_name"].(string) != "" {
		model.DisplayName = core.StringPtr(modelMap["display_name"].(string))
	}
	if modelMap["value_constraint"] != nil && modelMap["value_constraint"].(string) != "" {
		model.ValueConstraint = core.StringPtr(modelMap["value_constraint"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["required"] != nil {
		model.Required = core.BoolPtr(modelMap["required"].(bool))
	}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	if modelMap["custom_config"] != nil && len(modelMap["custom_config"].([]interface{})) > 0 {
		CustomConfigModel, err := resourceIBMCmOfferingMapToRenderType(modelMap["custom_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.CustomConfig = CustomConfigModel
	}
	if modelMap["type_metadata"] != nil && modelMap["type_metadata"].(string) != "" {
		model.TypeMetadata = core.StringPtr(modelMap["type_metadata"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToRenderType(modelMap map[string]interface{}) (*catalogmanagementv1.RenderType, error) {
	model := &catalogmanagementv1.RenderType{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["grouping"] != nil && modelMap["grouping"].(string) != "" {
		model.Grouping = core.StringPtr(modelMap["grouping"].(string))
	}
	if modelMap["original_grouping"] != nil && modelMap["original_grouping"].(string) != "" {
		model.OriginalGrouping = core.StringPtr(modelMap["original_grouping"].(string))
	}
	if modelMap["grouping_index"] != nil {
		model.GroupingIndex = core.Int64Ptr(int64(modelMap["grouping_index"].(int)))
	}
	if modelMap["config_constraints"] != nil {
		model.ConfigConstraints = modelMap["config_constraints"].(map[string]interface{})
	}
	if modelMap["associations"] != nil && len(modelMap["associations"].([]interface{})) > 0 {
		AssociationsModel, err := resourceIBMCmOfferingMapToRenderTypeAssociations(modelMap["associations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Associations = AssociationsModel
	}
	return model, nil
}

func resourceIBMCmOfferingMapToRenderTypeAssociations(modelMap map[string]interface{}) (*catalogmanagementv1.RenderTypeAssociations, error) {
	model := &catalogmanagementv1.RenderTypeAssociations{}
	if modelMap["parameters"] != nil {
		parameters := []catalogmanagementv1.RenderTypeAssociationsParametersItem{}
		for _, parametersItem := range modelMap["parameters"].([]interface{}) {
			parametersItemModel, err := resourceIBMCmOfferingMapToRenderTypeAssociationsParametersItem(parametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			parameters = append(parameters, *parametersItemModel)
		}
		model.Parameters = parameters
	}
	return model, nil
}

func resourceIBMCmOfferingMapToRenderTypeAssociationsParametersItem(modelMap map[string]interface{}) (*catalogmanagementv1.RenderTypeAssociationsParametersItem, error) {
	model := &catalogmanagementv1.RenderTypeAssociationsParametersItem{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["options_refresh"] != nil {
		model.OptionsRefresh = core.BoolPtr(modelMap["options_refresh"].(bool))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToOutput(modelMap map[string]interface{}) (*catalogmanagementv1.Output, error) {
	model := &catalogmanagementv1.Output{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToIamPermission(modelMap map[string]interface{}) (*catalogmanagementv1.IamPermission, error) {
	model := &catalogmanagementv1.IamPermission{}
	if modelMap["service_name"] != nil && modelMap["service_name"].(string) != "" {
		model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	}
	if modelMap["role_crns"] != nil {
		roleCrns := []string{}
		for _, roleCrnsItem := range modelMap["role_crns"].([]interface{}) {
			roleCrns = append(roleCrns, roleCrnsItem.(string))
		}
		model.RoleCrns = roleCrns
	}
	if modelMap["resources"] != nil {
		resources := []catalogmanagementv1.IamResource{}
		for _, resourcesItem := range modelMap["resources"].([]interface{}) {
			resourcesItemModel, err := resourceIBMCmOfferingMapToIamResource(resourcesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			resources = append(resources, *resourcesItemModel)
		}
		model.Resources = resources
	}
	return model, nil
}

func resourceIBMCmOfferingMapToIamResource(modelMap map[string]interface{}) (*catalogmanagementv1.IamResource, error) {
	model := &catalogmanagementv1.IamResource{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["role_crns"] != nil {
		roleCrns := []string{}
		for _, roleCrnsItem := range modelMap["role_crns"].([]interface{}) {
			roleCrns = append(roleCrns, roleCrnsItem.(string))
		}
		model.RoleCrns = roleCrns
	}
	return model, nil
}

func resourceIBMCmOfferingMapToValidation(modelMap map[string]interface{}) (*catalogmanagementv1.Validation, error) {
	model := &catalogmanagementv1.Validation{}
	if modelMap["validated"] != nil {

	}
	if modelMap["requested"] != nil {

	}
	if modelMap["state"] != nil && modelMap["state"].(string) != "" {
		model.State = core.StringPtr(modelMap["state"].(string))
	}
	if modelMap["last_operation"] != nil && modelMap["last_operation"].(string) != "" {
		model.LastOperation = core.StringPtr(modelMap["last_operation"].(string))
	}
	if modelMap["message"] != nil && modelMap["message"].(string) != "" {
		model.Message = core.StringPtr(modelMap["message"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToResource(modelMap map[string]interface{}) (*catalogmanagementv1.Resource, error) {
	model := &catalogmanagementv1.Resource{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["value"] != nil {

	}
	return model, nil
}

func resourceIBMCmOfferingMapToScript(modelMap map[string]interface{}) (*catalogmanagementv1.Script, error) {
	model := &catalogmanagementv1.Script{}
	if modelMap["instructions"] != nil && modelMap["instructions"].(string) != "" {
		model.Instructions = core.StringPtr(modelMap["instructions"].(string))
	}
	if modelMap["script"] != nil && modelMap["script"].(string) != "" {
		model.Script = core.StringPtr(modelMap["script"].(string))
	}
	if modelMap["script_permission"] != nil && modelMap["script_permission"].(string) != "" {
		model.ScriptPermission = core.StringPtr(modelMap["script_permission"].(string))
	}
	if modelMap["delete_script"] != nil && modelMap["delete_script"].(string) != "" {
		model.DeleteScript = core.StringPtr(modelMap["delete_script"].(string))
	}
	if modelMap["scope"] != nil && modelMap["scope"].(string) != "" {
		model.Scope = core.StringPtr(modelMap["scope"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToVersionEntitlement(modelMap map[string]interface{}) (*catalogmanagementv1.VersionEntitlement, error) {
	model := &catalogmanagementv1.VersionEntitlement{}
	if modelMap["provider_name"] != nil && modelMap["provider_name"].(string) != "" {
		model.ProviderName = core.StringPtr(modelMap["provider_name"].(string))
	}
	if modelMap["provider_id"] != nil && modelMap["provider_id"].(string) != "" {
		model.ProviderID = core.StringPtr(modelMap["provider_id"].(string))
	}
	if modelMap["product_id"] != nil && modelMap["product_id"].(string) != "" {
		model.ProductID = core.StringPtr(modelMap["product_id"].(string))
	}
	if modelMap["part_numbers"] != nil {
		partNumbers := []string{}
		for _, partNumbersItem := range modelMap["part_numbers"].([]interface{}) {
			partNumbers = append(partNumbers, partNumbersItem.(string))
		}
		model.PartNumbers = partNumbers
	}
	if modelMap["image_repo_name"] != nil && modelMap["image_repo_name"].(string) != "" {
		model.ImageRepoName = core.StringPtr(modelMap["image_repo_name"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToLicense(modelMap map[string]interface{}) (*catalogmanagementv1.License, error) {
	model := &catalogmanagementv1.License{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToState(modelMap map[string]interface{}) (*catalogmanagementv1.State, error) {
	model := &catalogmanagementv1.State{}
	if modelMap["current"] != nil && modelMap["current"].(string) != "" {
		model.Current = core.StringPtr(modelMap["current"].(string))
	}
	if modelMap["current_entered"] != nil {

	}
	if modelMap["pending"] != nil && modelMap["pending"].(string) != "" {
		model.Pending = core.StringPtr(modelMap["pending"].(string))
	}
	if modelMap["pending_requested"] != nil {

	}
	if modelMap["previous"] != nil && modelMap["previous"].(string) != "" {
		model.Previous = core.StringPtr(modelMap["previous"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToDeprecatePending(modelMap map[string]interface{}) (*catalogmanagementv1.DeprecatePending, error) {
	model := &catalogmanagementv1.DeprecatePending{}
	if modelMap["deprecate_date"] != nil {

	}
	if modelMap["deprecate_state"] != nil && modelMap["deprecate_state"].(string) != "" {
		model.DeprecateState = core.StringPtr(modelMap["deprecate_state"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSolutionInfo(modelMap map[string]interface{}) (*catalogmanagementv1.SolutionInfo, error) {
	model := &catalogmanagementv1.SolutionInfo{}
	if modelMap["architecture_diagrams"] != nil {
		architectureDiagrams := []catalogmanagementv1.ArchitectureDiagram{}
		for _, architectureDiagramsItem := range modelMap["architecture_diagrams"].([]interface{}) {
			architectureDiagramsItemModel, err := resourceIBMCmOfferingMapToArchitectureDiagram(architectureDiagramsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			architectureDiagrams = append(architectureDiagrams, *architectureDiagramsItemModel)
		}
		model.ArchitectureDiagrams = architectureDiagrams
	}
	if modelMap["features"] != nil {
		features := []catalogmanagementv1.Feature{}
		for _, featuresItem := range modelMap["features"].([]interface{}) {
			featuresItemModel, err := resourceIBMCmOfferingMapToFeature(featuresItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			features = append(features, *featuresItemModel)
		}
		model.Features = features
	}
	if modelMap["cost_estimate"] != nil && len(modelMap["cost_estimate"].([]interface{})) > 0 {
		CostEstimateModel, err := resourceIBMCmOfferingMapToCostEstimate(modelMap["cost_estimate"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.CostEstimate = CostEstimateModel
	}
	if modelMap["dependencies"] != nil {
		dependencies := []catalogmanagementv1.OfferingReference{}
		for _, dependenciesItem := range modelMap["dependencies"].([]interface{}) {
			dependenciesItemModel, err := resourceIBMCmOfferingMapToDependency(dependenciesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			dependencies = append(dependencies, *dependenciesItemModel)
		}
		model.Dependencies = dependencies
	}
	return model, nil
}

func resourceIBMCmOfferingMapToArchitectureDiagram(modelMap map[string]interface{}) (*catalogmanagementv1.ArchitectureDiagram, error) {
	model := &catalogmanagementv1.ArchitectureDiagram{}
	if modelMap["diagram"] != nil && len(modelMap["diagram"].([]interface{})) > 0 {
		DiagramModel, err := resourceIBMCmOfferingMapToMediaItem(modelMap["diagram"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Diagram = DiagramModel
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToMediaItem(modelMap map[string]interface{}) (*catalogmanagementv1.MediaItem, error) {
	model := &catalogmanagementv1.MediaItem{}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["api_url"] != nil && modelMap["api_url"].(string) != "" {
		model.APIURL = core.StringPtr(modelMap["api_url"].(string))
	}
	if modelMap["url_proxy"] != nil && len(modelMap["url_proxy"].([]interface{})) > 0 {
		URLProxyModel, err := resourceIBMCmOfferingMapToURLProxy(modelMap["url_proxy"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.URLProxy = URLProxyModel
	}
	if modelMap["caption"] != nil && modelMap["caption"].(string) != "" {
		model.Caption = core.StringPtr(modelMap["caption"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["thumbnail_url"] != nil && modelMap["thumbnail_url"].(string) != "" {
		model.ThumbnailURL = core.StringPtr(modelMap["thumbnail_url"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToURLProxy(modelMap map[string]interface{}) (*catalogmanagementv1.URLProxy, error) {
	model := &catalogmanagementv1.URLProxy{}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["sha"] != nil && modelMap["sha"].(string) != "" {
		model.Sha = core.StringPtr(modelMap["sha"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToCostEstimate(modelMap map[string]interface{}) (*catalogmanagementv1.CostEstimate, error) {
	model := &catalogmanagementv1.CostEstimate{}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	if modelMap["currency"] != nil && modelMap["currency"].(string) != "" {
		model.Currency = core.StringPtr(modelMap["currency"].(string))
	}
	if modelMap["projects"] != nil {
		projects := []catalogmanagementv1.Project{}
		for _, projectsItem := range modelMap["projects"].([]interface{}) {
			projectsItemModel, err := resourceIBMCmOfferingMapToProject(projectsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			projects = append(projects, *projectsItemModel)
		}
		model.Projects = projects
	}
	if modelMap["summary"] != nil && len(modelMap["summary"].([]interface{})) > 0 {
		SummaryModel, err := resourceIBMCmOfferingMapToCostSummary(modelMap["summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Summary = SummaryModel
	}
	if modelMap["total_hourly_cost"] != nil && modelMap["total_hourly_cost"].(string) != "" {
		model.TotalHourlyCost = core.StringPtr(modelMap["total_hourly_cost"].(string))
	}
	if modelMap["total_monthly_cost"] != nil && modelMap["total_monthly_cost"].(string) != "" {
		model.TotalMonthlyCost = core.StringPtr(modelMap["total_monthly_cost"].(string))
	}
	if modelMap["past_total_hourly_cost"] != nil && modelMap["past_total_hourly_cost"].(string) != "" {
		model.PastTotalHourlyCost = core.StringPtr(modelMap["past_total_hourly_cost"].(string))
	}
	if modelMap["past_total_monthly_cost"] != nil && modelMap["past_total_monthly_cost"].(string) != "" {
		model.PastTotalMonthlyCost = core.StringPtr(modelMap["past_total_monthly_cost"].(string))
	}
	if modelMap["diff_total_hourly_cost"] != nil && modelMap["diff_total_hourly_cost"].(string) != "" {
		model.DiffTotalHourlyCost = core.StringPtr(modelMap["diff_total_hourly_cost"].(string))
	}
	if modelMap["diff_total_monthly_cost"] != nil && modelMap["diff_total_monthly_cost"].(string) != "" {
		model.DiffTotalMonthlyCost = core.StringPtr(modelMap["diff_total_monthly_cost"].(string))
	}
	if modelMap["time_generated"] != nil {

	}
	return model, nil
}

func resourceIBMCmOfferingMapToProject(modelMap map[string]interface{}) (*catalogmanagementv1.Project, error) {
	model := &catalogmanagementv1.Project{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["past_breakdown"] != nil && len(modelMap["past_breakdown"].([]interface{})) > 0 {
		PastBreakdownModel, err := resourceIBMCmOfferingMapToCostBreakdown(modelMap["past_breakdown"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PastBreakdown = PastBreakdownModel
	}
	if modelMap["breakdown"] != nil && len(modelMap["breakdown"].([]interface{})) > 0 {
		BreakdownModel, err := resourceIBMCmOfferingMapToCostBreakdown(modelMap["breakdown"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Breakdown = BreakdownModel
	}
	if modelMap["diff"] != nil && len(modelMap["diff"].([]interface{})) > 0 {
		DiffModel, err := resourceIBMCmOfferingMapToCostBreakdown(modelMap["diff"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Diff = DiffModel
	}
	if modelMap["summary"] != nil && len(modelMap["summary"].([]interface{})) > 0 {
		SummaryModel, err := resourceIBMCmOfferingMapToCostSummary(modelMap["summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Summary = SummaryModel
	}
	return model, nil
}

func resourceIBMCmOfferingMapToCostBreakdown(modelMap map[string]interface{}) (*catalogmanagementv1.CostBreakdown, error) {
	model := &catalogmanagementv1.CostBreakdown{}
	if modelMap["total_hourly_cost"] != nil && modelMap["total_hourly_cost"].(string) != "" {
		model.TotalHourlyCost = core.StringPtr(modelMap["total_hourly_cost"].(string))
	}
	if modelMap["total_monthly_cost"] != nil && modelMap["total_monthly_cost"].(string) != "" {
		model.TotalMonthlyCost = core.StringPtr(modelMap["total_monthly_cost"].(string))
	}
	if modelMap["resources"] != nil {
		resources := []catalogmanagementv1.CostResource{}
		for _, resourcesItem := range modelMap["resources"].([]interface{}) {
			resourcesItemModel, err := resourceIBMCmOfferingMapToCostResource(resourcesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			resources = append(resources, *resourcesItemModel)
		}
		model.Resources = resources
	}
	return model, nil
}

func resourceIBMCmOfferingMapToCostResource(modelMap map[string]interface{}) (*catalogmanagementv1.CostResource, error) {
	model := &catalogmanagementv1.CostResource{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["hourly_cost"] != nil && modelMap["hourly_cost"].(string) != "" {
		model.HourlyCost = core.StringPtr(modelMap["hourly_cost"].(string))
	}
	if modelMap["monthly_cost"] != nil && modelMap["monthly_cost"].(string) != "" {
		model.MonthlyCost = core.StringPtr(modelMap["monthly_cost"].(string))
	}
	if modelMap["cost_components"] != nil {
		costComponents := []catalogmanagementv1.CostComponent{}
		for _, costComponentsItem := range modelMap["cost_components"].([]interface{}) {
			costComponentsItemModel, err := resourceIBMCmOfferingMapToCostComponent(costComponentsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			costComponents = append(costComponents, *costComponentsItemModel)
		}
		model.CostComponents = costComponents
	}
	return model, nil
}

func resourceIBMCmOfferingMapToCostComponent(modelMap map[string]interface{}) (*catalogmanagementv1.CostComponent, error) {
	model := &catalogmanagementv1.CostComponent{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	if modelMap["hourly_quantity"] != nil && modelMap["hourly_quantity"].(string) != "" {
		model.HourlyQuantity = core.StringPtr(modelMap["hourly_quantity"].(string))
	}
	if modelMap["monthly_quantity"] != nil && modelMap["monthly_quantity"].(string) != "" {
		model.MonthlyQuantity = core.StringPtr(modelMap["monthly_quantity"].(string))
	}
	if modelMap["price"] != nil && modelMap["price"].(string) != "" {
		model.Price = core.StringPtr(modelMap["price"].(string))
	}
	if modelMap["hourly_cost"] != nil && modelMap["hourly_cost"].(string) != "" {
		model.HourlyCost = core.StringPtr(modelMap["hourly_cost"].(string))
	}
	if modelMap["monthly_cost"] != nil && modelMap["monthly_cost"].(string) != "" {
		model.MonthlyCost = core.StringPtr(modelMap["monthly_cost"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToCostSummary(modelMap map[string]interface{}) (*catalogmanagementv1.CostSummary, error) {
	model := &catalogmanagementv1.CostSummary{}
	if modelMap["total_detected_resources"] != nil {
		model.TotalDetectedResources = core.Int64Ptr(int64(modelMap["total_detected_resources"].(int)))
	}
	if modelMap["total_supported_resources"] != nil {
		model.TotalSupportedResources = core.Int64Ptr(int64(modelMap["total_supported_resources"].(int)))
	}
	if modelMap["total_unsupported_resources"] != nil {
		model.TotalUnsupportedResources = core.Int64Ptr(int64(modelMap["total_unsupported_resources"].(int)))
	}
	if modelMap["total_usage_based_resources"] != nil {
		model.TotalUsageBasedResources = core.Int64Ptr(int64(modelMap["total_usage_based_resources"].(int)))
	}
	if modelMap["total_no_price_resources"] != nil {
		model.TotalNoPriceResources = core.Int64Ptr(int64(modelMap["total_no_price_resources"].(int)))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToDependency(modelMap map[string]interface{}) (*catalogmanagementv1.OfferingReference, error) {
	model := &catalogmanagementv1.OfferingReference{}
	if modelMap["catalog_id"] != nil && modelMap["catalog_id"].(string) != "" {
		model.CatalogID = core.StringPtr(modelMap["catalog_id"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	if modelMap["flavors"] != nil {
		flavors := []string{}
		for _, flavorsItem := range modelMap["flavors"].([]interface{}) {
			flavors = append(flavors, flavorsItem.(string))
		}
		model.Flavors = flavors
	}
	return model, nil
}

func resourceIBMCmOfferingMapToProviderInfo(modelMap map[string]interface{}) (*catalogmanagementv1.ProviderInfo, error) {
	model := &catalogmanagementv1.ProviderInfo{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToRepoInfo(modelMap map[string]interface{}) (*catalogmanagementv1.RepoInfo, error) {
	model := &catalogmanagementv1.RepoInfo{}
	if modelMap["token"] != nil && modelMap["token"].(string) != "" {
		model.Token = core.StringPtr(modelMap["token"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToImagePullKey(modelMap map[string]interface{}) (*catalogmanagementv1.ImagePullKey, error) {
	model := &catalogmanagementv1.ImagePullKey{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSupport(modelMap map[string]interface{}) (*catalogmanagementv1.Support, error) {
	model := &catalogmanagementv1.Support{}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["process"] != nil && modelMap["process"].(string) != "" {
		model.Process = core.StringPtr(modelMap["process"].(string))
	}
	if modelMap["locations"] != nil {
		locations := []string{}
		for _, locationsItem := range modelMap["locations"].([]interface{}) {
			locations = append(locations, locationsItem.(string))
		}
		model.Locations = locations
	}
	if modelMap["support_details"] != nil {
		supportDetails := []catalogmanagementv1.SupportDetail{}
		for _, supportDetailsItem := range modelMap["support_details"].([]interface{}) {
			supportDetailsItemModel, err := resourceIBMCmOfferingMapToSupportDetail(supportDetailsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			supportDetails = append(supportDetails, *supportDetailsItemModel)
		}
		model.SupportDetails = supportDetails
	}
	if modelMap["support_escalation"] != nil && len(modelMap["support_escalation"].([]interface{})) > 0 {
		SupportEscalationModel, err := resourceIBMCmOfferingMapToSupportEscalation(modelMap["support_escalation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SupportEscalation = SupportEscalationModel
	}
	if modelMap["support_type"] != nil && modelMap["support_type"].(string) != "" {
		model.SupportType = core.StringPtr(modelMap["support_type"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSupportDetail(modelMap map[string]interface{}) (*catalogmanagementv1.SupportDetail, error) {
	model := &catalogmanagementv1.SupportDetail{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["contact"] != nil && modelMap["contact"].(string) != "" {
		model.Contact = core.StringPtr(modelMap["contact"].(string))
	}
	if modelMap["response_wait_time"] != nil && len(modelMap["response_wait_time"].([]interface{})) > 0 {
		ResponseWaitTimeModel, err := resourceIBMCmOfferingMapToSupportWaitTime(modelMap["response_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResponseWaitTime = ResponseWaitTimeModel
	}
	if modelMap["availability"] != nil && len(modelMap["availability"].([]interface{})) > 0 {
		AvailabilityModel, err := resourceIBMCmOfferingMapToSupportAvailability(modelMap["availability"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Availability = AvailabilityModel
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSupportWaitTime(modelMap map[string]interface{}) (*catalogmanagementv1.SupportWaitTime, error) {
	model := &catalogmanagementv1.SupportWaitTime{}
	if modelMap["value"] != nil {
		model.Value = core.Int64Ptr(int64(modelMap["value"].(int)))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSupportAvailability(modelMap map[string]interface{}) (*catalogmanagementv1.SupportAvailability, error) {
	model := &catalogmanagementv1.SupportAvailability{}
	if modelMap["times"] != nil {
		times := []catalogmanagementv1.SupportTime{}
		for _, timesItem := range modelMap["times"].([]interface{}) {
			timesItemModel, err := resourceIBMCmOfferingMapToSupportTime(timesItem.(map[string]interface{}))
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

func resourceIBMCmOfferingMapToSupportTime(modelMap map[string]interface{}) (*catalogmanagementv1.SupportTime, error) {
	model := &catalogmanagementv1.SupportTime{}
	if modelMap["day"] != nil {
		model.Day = core.Int64Ptr(int64(modelMap["day"].(int)))
	}
	if modelMap["start_time"] != nil && modelMap["start_time"].(string) != "" {
		model.StartTime = core.StringPtr(modelMap["start_time"].(string))
	}
	if modelMap["end_time"] != nil && modelMap["end_time"].(string) != "" {
		model.EndTime = core.StringPtr(modelMap["end_time"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToSupportEscalation(modelMap map[string]interface{}) (*catalogmanagementv1.SupportEscalation, error) {
	model := &catalogmanagementv1.SupportEscalation{}
	if modelMap["escalation_wait_time"] != nil && len(modelMap["escalation_wait_time"].([]interface{})) > 0 {
		EscalationWaitTimeModel, err := resourceIBMCmOfferingMapToSupportWaitTime(modelMap["escalation_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.EscalationWaitTime = EscalationWaitTimeModel
	}
	if modelMap["response_wait_time"] != nil && len(modelMap["response_wait_time"].([]interface{})) > 0 {
		ResponseWaitTimeModel, err := resourceIBMCmOfferingMapToSupportWaitTime(modelMap["response_wait_time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResponseWaitTime = ResponseWaitTimeModel
	}
	if modelMap["contact"] != nil && modelMap["contact"].(string) != "" {
		model.Contact = core.StringPtr(modelMap["contact"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToBadge(modelMap map[string]interface{}) (*catalogmanagementv1.Badge, error) {
	model := &catalogmanagementv1.Badge{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["icon"] != nil && modelMap["icon"].(string) != "" {
		model.Icon = core.StringPtr(modelMap["icon"].(string))
	}
	if modelMap["authority"] != nil && modelMap["authority"].(string) != "" {
		model.Authority = core.StringPtr(modelMap["authority"].(string))
	}
	if modelMap["tag"] != nil && modelMap["tag"].(string) != "" {
		model.Tag = core.StringPtr(modelMap["tag"].(string))
	}
	if modelMap["learn_more_links"] != nil && len(modelMap["learn_more_links"].([]interface{})) > 0 {
		LearnMoreLinksModel, err := resourceIBMCmOfferingMapToLearnMoreLinks(modelMap["learn_more_links"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LearnMoreLinks = LearnMoreLinksModel
	}
	if modelMap["constraints"] != nil {
		constraints := []catalogmanagementv1.Constraint{}
		for _, constraintsItem := range modelMap["constraints"].([]interface{}) {
			constraintsItemModel, err := resourceIBMCmOfferingMapToConstraint(constraintsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			constraints = append(constraints, *constraintsItemModel)
		}
		model.Constraints = constraints
	}
	return model, nil
}

func resourceIBMCmOfferingMapToLearnMoreLinks(modelMap map[string]interface{}) (*catalogmanagementv1.LearnMoreLinks, error) {
	model := &catalogmanagementv1.LearnMoreLinks{}
	if modelMap["first_party"] != nil && modelMap["first_party"].(string) != "" {
		model.FirstParty = core.StringPtr(modelMap["first_party"].(string))
	}
	if modelMap["third_party"] != nil && modelMap["third_party"].(string) != "" {
		model.ThirdParty = core.StringPtr(modelMap["third_party"].(string))
	}
	return model, nil
}

func resourceIBMCmOfferingMapToConstraint(modelMap map[string]interface{}) (*catalogmanagementv1.Constraint, error) {
	model := &catalogmanagementv1.Constraint{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	// if modelMap["rule"] != nil {

	// }
	return model, nil
}

func resourceIBMCmOfferingRatingToMap(model *catalogmanagementv1.Rating) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.OneStarCount != nil {
		modelMap["one_star_count"] = flex.IntValue(model.OneStarCount)
	}
	if model.TwoStarCount != nil {
		modelMap["two_star_count"] = flex.IntValue(model.TwoStarCount)
	}
	if model.ThreeStarCount != nil {
		modelMap["three_star_count"] = flex.IntValue(model.ThreeStarCount)
	}
	if model.FourStarCount != nil {
		modelMap["four_star_count"] = flex.IntValue(model.FourStarCount)
	}
	return modelMap, nil
}

func resourceIBMCmOfferingFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Title != nil {
		modelMap["title"] = model.Title
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingKindToMap(model *catalogmanagementv1.Kind) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.FormatKind != nil {
		modelMap["format_kind"] = model.FormatKind
	}
	if model.InstallKind != nil {
		modelMap["install_kind"] = model.InstallKind
	}
	if model.TargetKind != nil {
		modelMap["target_kind"] = model.TargetKind
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.AdditionalFeatures != nil {
		additionalFeatures := []map[string]interface{}{}
		for _, additionalFeaturesItem := range model.AdditionalFeatures {
			additionalFeaturesItemMap, err := resourceIBMCmOfferingFeatureToMap(&additionalFeaturesItem)
			if err != nil {
				return modelMap, err
			}
			additionalFeatures = append(additionalFeatures, additionalFeaturesItemMap)
		}
		modelMap["additional_features"] = additionalFeatures
	}
	if model.Created != nil {
		modelMap["created"] = model.Created.String()
	}
	if model.Updated != nil {
		modelMap["updated"] = model.Updated.String()
	}
	if model.Versions != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range model.Versions {
			versionsItemMap, err := resourceIBMCmOfferingVersionToMap(&versionsItem)
			if err != nil {
				return modelMap, err
			}
			versions = append(versions, versionsItemMap)
		}
		modelMap["versions"] = versions
	}
	return modelMap, nil
}

func resourceIBMCmOfferingVersionToMap(model *catalogmanagementv1.Version) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Rev != nil {
		modelMap["rev"] = model.Rev
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Flavor != nil {
		flavorMap, err := resourceIBMCmOfferingFlavorToMap(model.Flavor)
		if err != nil {
			return modelMap, err
		}
		modelMap["flavor"] = []map[string]interface{}{flavorMap}
	}
	if model.Sha != nil {
		modelMap["sha"] = model.Sha
	}
	if model.Created != nil {
		modelMap["created"] = model.Created.String()
	}
	if model.Updated != nil {
		modelMap["updated"] = model.Updated.String()
	}
	if model.OfferingID != nil {
		modelMap["offering_id"] = model.OfferingID
	}
	if model.CatalogID != nil {
		modelMap["catalog_id"] = model.CatalogID
	}
	if model.KindID != nil {
		modelMap["kind_id"] = model.KindID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.RepoURL != nil {
		modelMap["repo_url"] = model.RepoURL
	}
	if model.SourceURL != nil {
		modelMap["source_url"] = model.SourceURL
	}
	if model.TgzURL != nil {
		modelMap["tgz_url"] = model.TgzURL
	}
	if model.Configuration != nil {
		configuration := []map[string]interface{}{}
		for _, configurationItem := range model.Configuration {
			configurationItemMap, err := resourceIBMCmOfferingConfigurationToMap(&configurationItem)
			if err != nil {
				return modelMap, err
			}
			configuration = append(configuration, configurationItemMap)
		}
		modelMap["configuration"] = configuration
	}
	if model.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range model.Outputs {
			outputsItemMap, err := resourceIBMCmOfferingOutputToMap(&outputsItem)
			if err != nil {
				return modelMap, err
			}
			outputs = append(outputs, outputsItemMap)
		}
		modelMap["outputs"] = outputs
	}
	if model.IamPermissions != nil {
		iamPermissions := []map[string]interface{}{}
		for _, iamPermissionsItem := range model.IamPermissions {
			iamPermissionsItemMap, err := resourceIBMCmOfferingIamPermissionToMap(&iamPermissionsItem)
			if err != nil {
				return modelMap, err
			}
			iamPermissions = append(iamPermissions, iamPermissionsItemMap)
		}
		modelMap["iam_permissions"] = iamPermissions
	}
	metadata := []map[string]interface{}{}
	if model.Metadata != nil {
		var modelMapVSI map[string]interface{}
		var err error
		if model.Metadata["vsi_vpc"] != nil {
			modelMapVSI, err = dataSourceIBMCmVersionMetadataVSIToMap(model.Metadata["vsi_vpc"].(map[string]interface{}))
			if err != nil {
				return modelMap, err
			}
		}
		convertedMap := make(map[string]interface{}, len(model.Metadata))
		for k, v := range model.Metadata {
			if k == "vsi_vpc" {
				convertedMap[k] = []map[string]interface{}{modelMapVSI}
			} else {
				convertedMap[k] = v
			}
		}
		metadata = append(metadata, convertedMap)
	}
	modelMap["metadata"] = metadata

	if model.Validation != nil {
		validationMap, err := resourceIBMCmOfferingValidationToMap(model.Validation)
		if err != nil {
			return modelMap, err
		}
		modelMap["validation"] = []map[string]interface{}{validationMap}
	}
	if model.RequiredResources != nil {
		requiredResources := []map[string]interface{}{}
		for _, requiredResourcesItem := range model.RequiredResources {
			requiredResourcesItemMap, err := resourceIBMCmOfferingResourceToMap(&requiredResourcesItem)
			if err != nil {
				return modelMap, err
			}
			requiredResources = append(requiredResources, requiredResourcesItemMap)
		}
		modelMap["required_resources"] = requiredResources
	}
	if model.SingleInstance != nil {
		modelMap["single_instance"] = model.SingleInstance
	}
	if model.Install != nil {
		installMap, err := resourceIBMCmOfferingScriptToMap(model.Install)
		if err != nil {
			return modelMap, err
		}
		modelMap["install"] = []map[string]interface{}{installMap}
	}
	if model.PreInstall != nil {
		preInstall := []map[string]interface{}{}
		for _, preInstallItem := range model.PreInstall {
			preInstallItemMap, err := resourceIBMCmOfferingScriptToMap(&preInstallItem)
			if err != nil {
				return modelMap, err
			}
			preInstall = append(preInstall, preInstallItemMap)
		}
		modelMap["pre_install"] = preInstall
	}
	if model.Entitlement != nil {
		entitlementMap, err := resourceIBMCmOfferingVersionEntitlementToMap(model.Entitlement)
		if err != nil {
			return modelMap, err
		}
		modelMap["entitlement"] = []map[string]interface{}{entitlementMap}
	}
	if model.Licenses != nil {
		licenses := []map[string]interface{}{}
		for _, licensesItem := range model.Licenses {
			licensesItemMap, err := resourceIBMCmOfferingLicenseToMap(&licensesItem)
			if err != nil {
				return modelMap, err
			}
			licenses = append(licenses, licensesItemMap)
		}
		modelMap["licenses"] = licenses
	}
	if model.ImageManifestURL != nil {
		modelMap["image_manifest_url"] = model.ImageManifestURL
	}
	if model.Deprecated != nil {
		modelMap["deprecated"] = model.Deprecated
	}
	if model.PackageVersion != nil {
		modelMap["package_version"] = model.PackageVersion
	}
	if model.State != nil {
		stateMap, err := resourceIBMCmOfferingStateToMap(model.State)
		if err != nil {
			return modelMap, err
		}
		modelMap["state"] = []map[string]interface{}{stateMap}
	}
	if model.VersionLocator != nil {
		modelMap["version_locator"] = model.VersionLocator
	}
	if model.LongDescription != nil {
		modelMap["long_description"] = model.LongDescription
	}
	if model.WhitelistedAccounts != nil {
		modelMap["whitelisted_accounts"] = model.WhitelistedAccounts
	}
	if model.ImagePullKeyName != nil {
		modelMap["image_pull_key_name"] = model.ImagePullKeyName
	}
	if model.DeprecatePending != nil {
		deprecatePendingMap, err := resourceIBMCmOfferingDeprecatePendingToMap(model.DeprecatePending)
		if err != nil {
			return modelMap, err
		}
		modelMap["deprecate_pending"] = []map[string]interface{}{deprecatePendingMap}
	}
	if model.SolutionInfo != nil {
		solutionInfoMap, err := resourceIBMCmOfferingSolutionInfoToMap(model.SolutionInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["solution_info"] = []map[string]interface{}{solutionInfoMap}
	}
	if model.IsConsumable != nil {
		modelMap["is_consumable"] = model.IsConsumable
	}
	return modelMap, nil
}

func resourceIBMCmOfferingFlavorToMap(model *catalogmanagementv1.Flavor) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Label != nil {
		modelMap["label"] = model.Label
	}
	if model.Index != nil {
		modelMap["index"] = flex.IntValue(model.Index)
	}
	return modelMap, nil
}

func resourceIBMCmOfferingConfigurationToMap(model *catalogmanagementv1.Configuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.DefaultValue != nil {
		defaultValueJson, err := json.Marshal(model.DefaultValue)
		if err != nil {
			return nil, flex.FmtErrorf("[ERROR] Error marshalling the version configuration default_value: %s", err)
		}
		modelMap["default_value"] = string(defaultValueJson)
	}
	if model.DisplayName != nil {
		modelMap["display_name"] = model.DisplayName
	}
	if model.ValueConstraint != nil {
		modelMap["value_constraint"] = model.ValueConstraint
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Required != nil {
		modelMap["required"] = model.Required
	}
	if model.Options != nil {
		options := []map[string]interface{}{}
		for _, optionsItem := range model.Options {
			options = append(options, optionsItem.(map[string]interface{}))
		}
		modelMap["options"] = options
	}
	if model.Hidden != nil {
		modelMap["hidden"] = model.Hidden
	}
	if model.CustomConfig != nil {
		customConfigMap, err := resourceIBMCmOfferingRenderTypeToMap(model.CustomConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_config"] = []map[string]interface{}{customConfigMap}
	}
	if model.TypeMetadata != nil {
		modelMap["type_metadata"] = model.TypeMetadata
	}
	return modelMap, nil
}

func resourceIBMCmOfferingRenderTypeToMap(model *catalogmanagementv1.RenderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Grouping != nil {
		modelMap["grouping"] = model.Grouping
	}
	if model.OriginalGrouping != nil {
		modelMap["original_grouping"] = model.OriginalGrouping
	}
	if model.GroupingIndex != nil {
		modelMap["grouping_index"] = flex.IntValue(model.GroupingIndex)
	}
	if model.ConfigConstraints != nil {
		modelMap["config_constraints"] = flex.Flatten(model.ConfigConstraints)
	}
	if model.Associations != nil {
		associationsMap, err := resourceIBMCmOfferingRenderTypeAssociationsToMap(model.Associations)
		if err != nil {
			return modelMap, err
		}
		modelMap["associations"] = []map[string]interface{}{associationsMap}
	}
	return modelMap, nil
}

func resourceIBMCmOfferingRenderTypeAssociationsToMap(model *catalogmanagementv1.RenderTypeAssociations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := resourceIBMCmOfferingRenderTypeAssociationsParametersItemToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func resourceIBMCmOfferingRenderTypeAssociationsParametersItemToMap(model *catalogmanagementv1.RenderTypeAssociationsParametersItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.OptionsRefresh != nil {
		modelMap["options_refresh"] = model.OptionsRefresh
	}
	return modelMap, nil
}

func resourceIBMCmOfferingOutputToMap(model *catalogmanagementv1.Output) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingIamPermissionToMap(model *catalogmanagementv1.IamPermission) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServiceName != nil {
		modelMap["service_name"] = model.ServiceName
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := resourceIBMCmOfferingIamResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func resourceIBMCmOfferingIamResourceToMap(model *catalogmanagementv1.IamResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	return modelMap, nil
}

func resourceIBMCmOfferingValidationToMap(model *catalogmanagementv1.Validation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Validated != nil {
		modelMap["validated"] = model.Validated.String()
	}
	if model.Requested != nil {
		modelMap["requested"] = model.Requested.String()
	}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.LastOperation != nil {
		modelMap["last_operation"] = model.LastOperation
	}
	if model.Message != nil {
		modelMap["message"] = model.Message
	}
	return modelMap, nil
}

func resourceIBMCmOfferingResourceToMap(model *catalogmanagementv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIBMCmOfferingScriptToMap(model *catalogmanagementv1.Script) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instructions != nil {
		modelMap["instructions"] = model.Instructions
	}
	if model.Script != nil {
		modelMap["script"] = model.Script
	}
	if model.ScriptPermission != nil {
		modelMap["script_permission"] = model.ScriptPermission
	}
	if model.DeleteScript != nil {
		modelMap["delete_script"] = model.DeleteScript
	}
	if model.Scope != nil {
		modelMap["scope"] = model.Scope
	}
	return modelMap, nil
}

func resourceIBMCmOfferingVersionEntitlementToMap(model *catalogmanagementv1.VersionEntitlement) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProviderName != nil {
		modelMap["provider_name"] = model.ProviderName
	}
	if model.ProviderID != nil {
		modelMap["provider_id"] = model.ProviderID
	}
	if model.ProductID != nil {
		modelMap["product_id"] = model.ProductID
	}
	if model.PartNumbers != nil {
		modelMap["part_numbers"] = model.PartNumbers
	}
	if model.ImageRepoName != nil {
		modelMap["image_repo_name"] = model.ImageRepoName
	}
	return modelMap, nil
}

func resourceIBMCmOfferingLicenseToMap(model *catalogmanagementv1.License) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Current != nil {
		modelMap["current"] = model.Current
	}
	if model.CurrentEntered != nil {
		modelMap["current_entered"] = model.CurrentEntered.String()
	}
	if model.Pending != nil {
		modelMap["pending"] = model.Pending
	}
	if model.PendingRequested != nil {
		modelMap["pending_requested"] = model.PendingRequested.String()
	}
	if model.Previous != nil {
		modelMap["previous"] = model.Previous
	}
	return modelMap, nil
}

func resourceIBMCmOfferingDeprecatePendingToMap(model *catalogmanagementv1.DeprecatePending) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DeprecateDate != nil {
		modelMap["deprecate_date"] = model.DeprecateDate.String()
	}
	if model.DeprecateState != nil {
		modelMap["deprecate_state"] = model.DeprecateState
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSolutionInfoToMap(model *catalogmanagementv1.SolutionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchitectureDiagrams != nil {
		architectureDiagrams := []map[string]interface{}{}
		for _, architectureDiagramsItem := range model.ArchitectureDiagrams {
			architectureDiagramsItemMap, err := resourceIBMCmOfferingArchitectureDiagramToMap(&architectureDiagramsItem)
			if err != nil {
				return modelMap, err
			}
			architectureDiagrams = append(architectureDiagrams, architectureDiagramsItemMap)
		}
		modelMap["architecture_diagrams"] = architectureDiagrams
	}
	if model.Features != nil {
		features := []map[string]interface{}{}
		for _, featuresItem := range model.Features {
			featuresItemMap, err := resourceIBMCmOfferingFeatureToMap(&featuresItem)
			if err != nil {
				return modelMap, err
			}
			features = append(features, featuresItemMap)
		}
		modelMap["features"] = features
	}
	if model.CostEstimate != nil {
		costEstimateMap, err := resourceIBMCmOfferingCostEstimateToMap(model.CostEstimate)
		if err != nil {
			return modelMap, err
		}
		modelMap["cost_estimate"] = []map[string]interface{}{costEstimateMap}
	}
	if model.Dependencies != nil {
		dependencies := []map[string]interface{}{}
		for _, dependenciesItem := range model.Dependencies {
			dependenciesItemMap, err := resourceIBMCmOfferingDependencyToMap(&dependenciesItem)
			if err != nil {
				return modelMap, err
			}
			dependencies = append(dependencies, dependenciesItemMap)
		}
		modelMap["dependencies"] = dependencies
	}
	return modelMap, nil
}

func resourceIBMCmOfferingArchitectureDiagramToMap(model *catalogmanagementv1.ArchitectureDiagram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Diagram != nil {
		diagramMap, err := resourceIBMCmOfferingMediaItemToMap(model.Diagram)
		if err != nil {
			return modelMap, err
		}
		modelMap["diagram"] = []map[string]interface{}{diagramMap}
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingMediaItemToMap(model *catalogmanagementv1.MediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.APIURL != nil {
		modelMap["api_url"] = model.APIURL
	}
	if model.URLProxy != nil {
		urlProxyMap, err := resourceIBMCmOfferingURLProxyToMap(model.URLProxy)
		if err != nil {
			return modelMap, err
		}
		modelMap["url_proxy"] = []map[string]interface{}{urlProxyMap}
	}
	if model.Caption != nil {
		modelMap["caption"] = model.Caption
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.ThumbnailURL != nil {
		modelMap["thumbnail_url"] = model.ThumbnailURL
	}
	return modelMap, nil
}

func resourceIBMCmOfferingURLProxyToMap(model *catalogmanagementv1.URLProxy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Sha != nil {
		modelMap["sha"] = model.Sha
	}
	return modelMap, nil
}

func resourceIBMCmOfferingCostEstimateToMap(model *catalogmanagementv1.CostEstimate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Currency != nil {
		modelMap["currency"] = model.Currency
	}
	if model.Projects != nil {
		projects := []map[string]interface{}{}
		for _, projectsItem := range model.Projects {
			projectsItemMap, err := resourceIBMCmOfferingProjectToMap(&projectsItem)
			if err != nil {
				return modelMap, err
			}
			projects = append(projects, projectsItemMap)
		}
		modelMap["projects"] = projects
	}
	if model.Summary != nil {
		summaryMap, err := resourceIBMCmOfferingCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = model.TotalHourlyCost
	}
	if model.TotalMonthlyCost != nil {
		modelMap["total_monthly_cost"] = model.TotalMonthlyCost
	}
	if model.PastTotalHourlyCost != nil {
		modelMap["past_total_hourly_cost"] = model.PastTotalHourlyCost
	}
	if model.PastTotalMonthlyCost != nil {
		modelMap["past_total_monthly_cost"] = model.PastTotalMonthlyCost
	}
	if model.DiffTotalHourlyCost != nil {
		modelMap["diff_total_hourly_cost"] = model.DiffTotalHourlyCost
	}
	if model.DiffTotalMonthlyCost != nil {
		modelMap["diff_total_monthly_cost"] = model.DiffTotalMonthlyCost
	}
	if model.TimeGenerated != nil {
		modelMap["time_generated"] = model.TimeGenerated.String()
	}
	return modelMap, nil
}

func resourceIBMCmOfferingProjectToMap(model *catalogmanagementv1.Project) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.PastBreakdown != nil {
		pastBreakdownMap, err := resourceIBMCmOfferingCostBreakdownToMap(model.PastBreakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["past_breakdown"] = []map[string]interface{}{pastBreakdownMap}
	}
	if model.Breakdown != nil {
		breakdownMap, err := resourceIBMCmOfferingCostBreakdownToMap(model.Breakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["breakdown"] = []map[string]interface{}{breakdownMap}
	}
	if model.Diff != nil {
		diffMap, err := resourceIBMCmOfferingCostBreakdownToMap(model.Diff)
		if err != nil {
			return modelMap, err
		}
		modelMap["diff"] = []map[string]interface{}{diffMap}
	}
	if model.Summary != nil {
		summaryMap, err := resourceIBMCmOfferingCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	return modelMap, nil
}

func resourceIBMCmOfferingCostBreakdownToMap(model *catalogmanagementv1.CostBreakdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = model.TotalHourlyCost
	}
	if model.TotalMonthlyCost != nil {
		modelMap["total_monthly_cost"] = model.TotalMonthlyCost
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := resourceIBMCmOfferingCostResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func resourceIBMCmOfferingCostResourceToMap(model *catalogmanagementv1.CostResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = model.MonthlyCost
	}
	if model.CostComponents != nil {
		costComponents := []map[string]interface{}{}
		for _, costComponentsItem := range model.CostComponents {
			costComponentsItemMap, err := resourceIBMCmOfferingCostComponentToMap(&costComponentsItem)
			if err != nil {
				return modelMap, err
			}
			costComponents = append(costComponents, costComponentsItemMap)
		}
		modelMap["cost_components"] = costComponents
	}
	return modelMap, nil
}

func resourceIBMCmOfferingCostComponentToMap(model *catalogmanagementv1.CostComponent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Unit != nil {
		modelMap["unit"] = model.Unit
	}
	if model.HourlyQuantity != nil {
		modelMap["hourly_quantity"] = model.HourlyQuantity
	}
	if model.MonthlyQuantity != nil {
		modelMap["monthly_quantity"] = model.MonthlyQuantity
	}
	if model.Price != nil {
		modelMap["price"] = model.Price
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = model.MonthlyCost
	}
	return modelMap, nil
}

func resourceIBMCmOfferingCostSummaryToMap(model *catalogmanagementv1.CostSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalDetectedResources != nil {
		modelMap["total_detected_resources"] = flex.IntValue(model.TotalDetectedResources)
	}
	if model.TotalSupportedResources != nil {
		modelMap["total_supported_resources"] = flex.IntValue(model.TotalSupportedResources)
	}
	if model.TotalUnsupportedResources != nil {
		modelMap["total_unsupported_resources"] = flex.IntValue(model.TotalUnsupportedResources)
	}
	if model.TotalUsageBasedResources != nil {
		modelMap["total_usage_based_resources"] = flex.IntValue(model.TotalUsageBasedResources)
	}
	if model.TotalNoPriceResources != nil {
		modelMap["total_no_price_resources"] = flex.IntValue(model.TotalNoPriceResources)
	}
	return modelMap, nil
}

func resourceIBMCmOfferingDependencyToMap(model *catalogmanagementv1.OfferingReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CatalogID != nil {
		modelMap["catalog_id"] = model.CatalogID
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Flavors != nil {
		modelMap["flavors"] = model.Flavors
	}
	return modelMap, nil
}

func resourceIBMCmOfferingProviderInfoToMap(model *catalogmanagementv1.ProviderInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	return modelMap, nil
}

func resourceIBMCmOfferingRepoInfoToMap(model *catalogmanagementv1.RepoInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Token != nil {
		modelMap["token"] = model.Token
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIBMCmOfferingImagePullKeyToMap(model *catalogmanagementv1.ImagePullKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportToMap(model *catalogmanagementv1.Support) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Process != nil {
		modelMap["process"] = model.Process
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	if model.SupportDetails != nil {
		supportDetails := []map[string]interface{}{}
		for _, supportDetailsItem := range model.SupportDetails {
			supportDetailsItemMap, err := resourceIBMCmOfferingSupportDetailToMap(&supportDetailsItem)
			if err != nil {
				return modelMap, err
			}
			supportDetails = append(supportDetails, supportDetailsItemMap)
		}
		modelMap["support_details"] = supportDetails
	}
	if model.SupportEscalation != nil {
		supportEscalationMap, err := resourceIBMCmOfferingSupportEscalationToMap(model.SupportEscalation)
		if err != nil {
			return modelMap, err
		}
		modelMap["support_escalation"] = []map[string]interface{}{supportEscalationMap}
	}
	if model.SupportType != nil {
		modelMap["support_type"] = model.SupportType
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportDetailToMap(model *catalogmanagementv1.SupportDetail) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Contact != nil {
		modelMap["contact"] = model.Contact
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := resourceIBMCmOfferingSupportWaitTimeToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	if model.Availability != nil {
		availabilityMap, err := resourceIBMCmOfferingSupportAvailabilityToMap(model.Availability)
		if err != nil {
			return modelMap, err
		}
		modelMap["availability"] = []map[string]interface{}{availabilityMap}
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportWaitTimeToMap(model *catalogmanagementv1.SupportWaitTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = flex.IntValue(model.Value)
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportAvailabilityToMap(model *catalogmanagementv1.SupportAvailability) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Times != nil {
		times := []map[string]interface{}{}
		for _, timesItem := range model.Times {
			timesItemMap, err := resourceIBMCmOfferingSupportTimeToMap(&timesItem)
			if err != nil {
				return modelMap, err
			}
			times = append(times, timesItemMap)
		}
		modelMap["times"] = times
	}
	if model.Timezone != nil {
		modelMap["timezone"] = model.Timezone
	}
	if model.AlwaysAvailable != nil {
		modelMap["always_available"] = model.AlwaysAvailable
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportTimeToMap(model *catalogmanagementv1.SupportTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Day != nil {
		modelMap["day"] = flex.IntValue(model.Day)
	}
	if model.StartTime != nil {
		modelMap["start_time"] = model.StartTime
	}
	if model.EndTime != nil {
		modelMap["end_time"] = model.EndTime
	}
	return modelMap, nil
}

func resourceIBMCmOfferingSupportEscalationToMap(model *catalogmanagementv1.SupportEscalation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EscalationWaitTime != nil {
		escalationWaitTimeMap, err := resourceIBMCmOfferingSupportWaitTimeToMap(model.EscalationWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["escalation_wait_time"] = []map[string]interface{}{escalationWaitTimeMap}
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := resourceIBMCmOfferingSupportWaitTimeToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	if model.Contact != nil {
		modelMap["contact"] = model.Contact
	}
	return modelMap, nil
}

func resourceIBMCmOfferingBadgeToMap(model *catalogmanagementv1.Badge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Label != nil {
		modelMap["label"] = model.Label
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Icon != nil {
		modelMap["icon"] = model.Icon
	}
	if model.Authority != nil {
		modelMap["authority"] = model.Authority
	}
	if model.Tag != nil {
		modelMap["tag"] = model.Tag
	}
	if model.LearnMoreLinks != nil {
		learnMoreLinksMap, err := resourceIBMCmOfferingLearnMoreLinksToMap(model.LearnMoreLinks)
		if err != nil {
			return modelMap, err
		}
		modelMap["learn_more_links"] = []map[string]interface{}{learnMoreLinksMap}
	}
	if model.Constraints != nil {
		constraints := []map[string]interface{}{}
		for _, constraintsItem := range model.Constraints {
			constraintsItemMap, err := resourceIBMCmOfferingConstraintToMap(&constraintsItem)
			if err != nil {
				return modelMap, err
			}
			constraints = append(constraints, constraintsItemMap)
		}
		modelMap["constraints"] = constraints
	}
	return modelMap, nil
}

func resourceIBMCmOfferingLearnMoreLinksToMap(model *catalogmanagementv1.LearnMoreLinks) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FirstParty != nil {
		modelMap["first_party"] = model.FirstParty
	}
	if model.ThirdParty != nil {
		modelMap["third_party"] = model.ThirdParty
	}
	return modelMap, nil
}

func resourceIBMCmOfferingConstraintToMap(model *catalogmanagementv1.Constraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	// if model.Rule != nil {
	// 	modelMap["rule"] = model.Rule
	// }
	return modelMap, nil
}

func handleShareOfferingAfterCreate(catalogManagementClient *catalogmanagementv1.CatalogManagementV1, offering catalogmanagementv1.Offering, d *schema.ResourceData, context context.Context) error {
	if _, ok := d.GetOk("share_with_access_list"); ok {
		// Share with accounts if there are any
		if len(d.Get("share_with_access_list").([]interface{})) > 0 {
			addOfferingAccessListOptions := catalogmanagementv1.AddOfferingAccessListOptions{}
			addOfferingAccessListOptions.SetCatalogIdentifier(*offering.CatalogID)
			addOfferingAccessListOptions.SetOfferingID(*offering.ID)
			addOfferingAccessListOptions.SetAccesses(SIToSS(d.Get("share_with_access_list").([]interface{})))
			_, response, err := catalogManagementClient.AddOfferingAccessListWithContext(context, &addOfferingAccessListOptions)
			if err != nil {
				log.Printf("[DEBUG] AddOfferingAccessListWithContext failed %s\n%s", err, response)
				return flex.FmtErrorf("AddOfferingAccessListWithContext failed %s\n%s", err, response)
			}
		}
	}

	shareOffering := false
	shareOfferingOptions := catalogmanagementv1.ShareOfferingOptions{}
	shareOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
	shareOfferingOptions.SetOfferingID(*offering.ID)
	shareOfferingOptions.SetEnabled(d.Get("share_enabled").(bool))

	if _, ok := d.GetOk("share_with_ibm"); ok {
		shareOffering = true
	}

	if _, ok := d.GetOk("share_with_ibm"); ok {
		shareOffering = true
		shareOfferingOptions.SetIBM(d.Get("share_with_ibm").(bool))
	}

	if _, ok := d.GetOk("share_with_all"); ok {
		shareOffering = true
		shareOfferingOptions.SetPublic(d.Get("share_with_all").(bool))
	}

	if shareOffering {
		_, response, err := catalogManagementClient.ShareOfferingWithContext(context, &shareOfferingOptions)
		if err != nil {
			log.Printf("[DEBUG] ShareOfferingWithContext failed %s\n%s", err, response)
			return flex.FmtErrorf("ShareOfferingWithContext failed %s\n%s", err, response)
		}
	}

	return nil
}
