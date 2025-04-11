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
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmOffering() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmOfferingRead,

		Schema: map[string]*schema.Schema{
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Offering identifier.",
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
				Computed:    true,
				Description: "Display Name in the requested language.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The programmatic name of this offering.",
			},
			"offering_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an icon associated with this offering.",
			},
			"offering_docs_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an additional docs with this offering.",
			},
			"offering_support_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "[deprecated] - Use offering.support instead.  URL to be displayed in the Consumption UI for getting support on this offering.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tags associated with this catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"keywords": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of keywords associated with offering, typically used to search for it.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Computed:    true,
				Description: "Short description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"long_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Long description in the requested language.",
			},
			"long_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"features": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of features associated with this offering.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
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
							Computed:    true,
							Description: "Unique ID.",
						},
						"format_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "content kind, e.g., helm, vm image.",
						},
						"install_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "install kind, e.g., helm, operator, terraform.",
						},
						"target_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "target cloud to install, e.g., iks, open_shift_iks.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Open ended metadata information.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tags associated with this catalog.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"additional_features": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of features associated with this offering.",
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
										Elem:        &schema.Schema{Type: schema.TypeString},
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
										Elem:        &schema.Schema{Type: schema.TypeString},
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
						"versions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of versions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version's CRN.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version of content type.",
									},
									"flavor": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Version Flavor Information.  Only supported for Product kind Solution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Programmatic name for this flavor.",
												},
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Label for this flavor.",
												},
												"label_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"index": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Order that this flavor should appear when listed for a single version.",
												},
											},
										},
									},
									"sha": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "hash of the content.",
									},
									"created": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time this version was created.",
									},
									"updated": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time this version was last updated.",
									},
									"offering_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Offering ID.",
									},
									"catalog_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Catalog ID.",
									},
									"kind_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kind ID.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of tags associated with this catalog.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"repo_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content's repo URL.",
									},
									"source_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content's source URL (e.g git repo).",
									},
									"tgz_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File used to on-board this version.",
									},
									"configuration": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of user solicited overrides.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Configuration key.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value type (string, boolean, int).",
												},
												"default_value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The default value as a JSON encoded string.  To use a secret when the type is password, specify a JSON encoded value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.",
												},
												"display_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Display name for configuration type.",
												},
												"value_constraint": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Constraint associated with value, e.g., for string type - regx:[a-z].",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key description.",
												},
												"required": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Is key required to install.",
												},
												"options": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of options of type.",
													Elem:        &schema.Schema{Type: schema.TypeMap},
												},
												"hidden": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Hide values.",
												},
												"custom_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Render type.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "ID of the widget type.",
															},
															"grouping": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Determines where this configuration type is rendered (3 sections today - Target, Resource, and Deployment).",
															},
															"original_grouping": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Original grouping type for this configuration (3 types - Target, Resource, and Deployment).",
															},
															"grouping_index": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Determines the order that this configuration item shows in that particular grouping.",
															},
															"config_constraints": &schema.Schema{
																Type:        schema.TypeMap,
																Computed:    true,
																Description: "Map of constraint parameters that will be passed to the custom widget.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"associations": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of parameters that are associated with this configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"parameters": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Parameters for this association.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Name of this parameter.",
																					},
																					"options_refresh": &schema.Schema{
																						Type:        schema.TypeBool,
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
													Computed:    true,
													Description: "The original type, as found in the source being onboarded.",
												},
											},
										},
									},
									"outputs": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of output values for this version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Output key.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Output description.",
												},
											},
										},
									},
									"iam_permissions": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of IAM permissions that are required to consume this version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Service name.",
												},
												"role_crns": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Role CRNs for this permission.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"resources": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Resources for this permission.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource name.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource description.",
															},
															"role_crns": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
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
										Computed:    true,
										Description: "Generic data to be included with content being onboarded. Required for virtual server image for VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_url": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Version source URL.",
												},
												"working_directory": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Working directory of source files.",
												},
												"example_name": &schema.Schema{
													Type:        schema.TypeString,
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
													Computed:    true,
													Description: "Usage text for the version.",
												},
												"usage_template": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Usage text for the version.",
												},
												"modules": &schema.Schema{
													Type:        schema.TypeList,
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
													Computed:    true,
													Description: "Version name.",
												},
												"terraform_version": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Terraform version.",
												},
												"validated_terraform_version": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
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
																Computed:    true,
																Description: "Operating system included in this image. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"dedicated_host_only": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
																		},
																		"vendor": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Vendor of the operating system. Required for virtual server image for VPC.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
																		},
																		"href": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "URL for this operating system. Required for virtual server image for VPC.",
																		},
																		"display_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
																		},
																		"family": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Software family for this operating system. Required for virtual server image for VPC.",
																		},
																		"version": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Major release version of this operating system. Required for virtual server image for VPC.",
																		},
																		"architecture": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Operating system architecture. Required for virtual server image for VPC.",
																		},
																	},
																},
															},
															"file": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Details for the stored image file. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"size": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
																		},
																	},
																},
															},
															"minimum_provisioned_size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
															},
															"images": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Image operating system. Required for virtual server image for VPC.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
																		},
																		"region": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
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
													Computed:    true,
													Description: "Operating system included in this image. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"dedicated_host_only": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
															},
															"vendor": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Vendor of the operating system. Required for virtual server image for VPC.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "URL for this operating system. Required for virtual server image for VPC.",
															},
															"display_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
															},
															"family": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Software family for this operating system. Required for virtual server image for VPC.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Major release version of this operating system. Required for virtual server image for VPC.",
															},
															"architecture": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operating system architecture. Required for virtual server image for VPC.",
															},
														},
													},
												},
												"file": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Details for the stored image file. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
															},
														},
													},
												},
												"minimum_provisioned_size": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
												},
												"images": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Image operating system. Required for virtual server image for VPC.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
															},
															"region": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
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
										Computed:    true,
										Description: "Validation response.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"validated": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Date and time of last successful validation.",
												},
												"requested": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Date and time of last validation was requested.",
												},
												"state": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Current validation state - <empty>, in_progress, valid, invalid, expired.",
												},
												"last_operation": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Last operation (e.g. submit_deployment, generate_installer, install_offering.",
												},
												"target": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "Validation target information (e.g. cluster_id, region, namespace, etc).  Values will vary by Content type.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Any message needing to be conveyed as part of the validation job.",
												},
											},
										},
									},
									"required_resources": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Resource requirments for installation.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of requirement.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "mem, disk, cores, and nodes can be parsed as an int.  targetVersion will be a semver range value.",
												},
											},
										},
									},
									"single_instance": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Denotes if single instance can be deployed to a given cluster.",
									},
									"install": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Script information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instructions": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
												},
												"instructions_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"script": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional script that needs to be run post any pre-condition script.",
												},
												"script_permission": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional iam permissions that are required on the target cluster to run this script.",
												},
												"delete_script": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional script that if run will remove the installed version.",
												},
												"scope": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
												},
											},
										},
									},
									"pre_install": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Optional pre-install instructions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instructions": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
												},
												"instructions_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"script": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional script that needs to be run post any pre-condition script.",
												},
												"script_permission": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional iam permissions that are required on the target cluster to run this script.",
												},
												"delete_script": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional script that if run will remove the installed version.",
												},
												"scope": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
												},
											},
										},
									},
									"entitlement": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Entitlement license info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"provider_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Provider name.",
												},
												"provider_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Provider ID.",
												},
												"product_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Product ID.",
												},
												"part_numbers": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "list of license entitlement part numbers, eg. D1YGZLL,D1ZXILL.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"image_repo_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Image repository name.",
												},
											},
										},
									},
									"licenses": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of licenses the product was built with.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "License ID.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "license name.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "type of license e.g., Apache xxx.",
												},
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "URL for the license text.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "License description.",
												},
											},
										},
									},
									"image_manifest_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "If set, denotes a url to a YAML file with list of container images used by this version.",
									},
									"deprecated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "read only field, indicating if this version is deprecated.",
									},
									"package_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version of the package used to create this version.",
									},
									"state": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Offering state.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"current": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
												"current_entered": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Date and time of current request.",
												},
												"pending": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
												"pending_requested": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Date and time of pending request.",
												},
												"previous": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "one of: new, validated, account-published, ibm-published, public-published.",
												},
											},
										},
									},
									"version_locator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A dotted value of `catalogID`.`versionID`.",
									},
									"long_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Long description for version.",
									},
									"long_description_i18n": &schema.Schema{
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "A map of translated strings, by language code.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"whitelisted_accounts": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Whitelisted accounts for version.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"image_pull_key_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the image pull key to use from Offering.ImagePullKeys.",
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
									"solution_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Version Solution Information.  Only supported for Product kind Solution.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"architecture_diagrams": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Architecture diagrams for this solution.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"diagram": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Offering Media information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "URL of the specified media item.",
																		},
																		"api_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "CM API specific URL of the specified media item.",
																		},
																		"url_proxy": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Offering URL proxy information.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "URL of the specified media item being proxied.",
																					},
																					"sha": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "SHA256 fingerprint of image.",
																					},
																				},
																			},
																		},
																		"caption": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Caption for this media item.",
																		},
																		"caption_i18n": &schema.Schema{
																			Type:        schema.TypeMap,
																			Computed:    true,
																			Description: "A map of translated strings, by language code.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of this media item.",
																		},
																		"thumbnail_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Thumbnail URL for this media item.",
																		},
																	},
																},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of this diagram.",
															},
															"description_i18n": &schema.Schema{
																Type:        schema.TypeMap,
																Computed:    true,
																Description: "A map of translated strings, by language code.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"features": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Features - titles only.",
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
																Elem:        &schema.Schema{Type: schema.TypeString},
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
																					"total_monthly_c_ost": &schema.Schema{
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
																					"total_monthly_c_ost": &schema.Schema{
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
																					"total_monthly_c_ost": &schema.Schema{
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
													Computed:    true,
													Description: "Dependencies for this solution.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"catalog_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Optional - If not specified, assumes the Public Catalog.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Optional - Offering ID - not required if name is set.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Optional - Programmatic Offering name.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Required - Semver value or range.",
															},
															"flavors": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
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
							Computed:    true,
							Description: "list of plans.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "unique id.",
									},
									"label": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display Name in the requested language.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The programmatic name of this offering.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Short description in the requested language.",
									},
									"long_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Long description in the requested language.",
									},
									"metadata": &schema.Schema{
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "open ended metadata information.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "list of tags associated with this catalog.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"additional_features": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "list of features associated with this offering.",
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
													Elem:        &schema.Schema{Type: schema.TypeString},
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
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"created": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the date'time this catalog was created.",
									},
									"updated": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the date'time this catalog was last updated.",
									},
									"deployments": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "list of deployments.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "unique id.",
												},
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Display Name in the requested language.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The programmatic name of this offering.",
												},
												"short_description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Short description in the requested language.",
												},
												"long_description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Long description in the requested language.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "open ended metadata information.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"tags": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "list of tags associated with this catalog.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"created": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "the date'time this catalog was created.",
												},
												"updated": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
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
				Computed:    true,
				Description: "Denotes public availability of an Offering - if share_enabled is true.",
			},
			"share_with_ibm": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes IBM employee availability of an Offering - if share_enabled is true.",
			},
			"share_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes sharing including access list availability of an Offering is enabled.",
			},
			"share_with_access_list": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of account IDs to add to this offering's access list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"public_original_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The original offering CRN that this publish entry came from.",
			},
			"publish_public_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the public catalog entry of this offering.",
			},
			"portal_approval_record": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The portal's approval record ID.",
			},
			"portal_ui_url": &schema.Schema{
				Type:        schema.TypeString,
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
				Computed:    true,
				Description: "Map of metadata values for this offering.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"disclaimer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A disclaimer for this offering.",
			},
			"hidden": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determine if this offering should be displayed in the Consumption UI.",
			},
			"provider_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on the provider for this offering, or omitted if no provider information is given.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of this provider.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
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
							Computed:    true,
							Description: "Token for private repos.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public or enterprise GitHub.",
						},
					},
				},
			},
			"image_pull_keys": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Image pull keys for this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key name.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key value.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
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
				Computed:    true,
				Description: "A list of media items related to this offering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the specified media item.",
						},
						"api_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CM API specific URL of the specified media item.",
						},
						"url_proxy": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Offering URL proxy information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of the specified media item being proxied.",
									},
									"sha": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SHA256 fingerprint of image.",
									},
								},
							},
						},
						"caption": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Caption for this media item.",
						},
						"caption_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of this media item.",
						},
						"thumbnail_url": &schema.Schema{
							Type:        schema.TypeString,
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

func dataSourceIBMCmOfferingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}

	getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getOfferingOptions.SetOfferingID(d.Get("offering_id").(string))

	offering, response, err := FetchOfferingWithAllVersions(context, catalogManagementClient, getOfferingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOfferingWithContext failed %s\n%s", err, response), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getOfferingOptions.OfferingID)

	if err = d.Set("rev", offering.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offering_identifier", offering.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_identifier: %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("url", offering.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", offering.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("label", offering.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if offering.LabelI18n != nil {
		if err = d.Set("label_i18n", offering.LabelI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label_i18n: : %s", err), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("name", offering.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offering_icon_url", offering.OfferingIconURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_icon_url: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offering_docs_url", offering.OfferingDocsURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_docs_url: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offering_support_url", offering.OfferingSupportURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_support_url: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tags := []string{}
	if offering.Tags != nil {
		tags = offering.Tags
	}
	if err = d.Set("tags", tags); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting tags: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	keywords := []string{}
	if offering.Keywords != nil {
		keywords = offering.Keywords
	}
	if err = d.Set("keywords", keywords); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting keywords: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	rating := []map[string]interface{}{}
	if offering.Rating != nil {
		modelMap, err := dataSourceIBMCmOfferingRatingToMap(offering.Rating)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		rating = append(rating, modelMap)
	}
	if err = d.Set("rating", rating); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rating : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created", flex.DateTimeToString(offering.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated", flex.DateTimeToString(offering.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("short_description", offering.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if offering.ShortDescriptionI18n != nil {
		if err = d.Set("short_description_i18n", offering.ShortDescriptionI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description_i18n: : %s", err), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("long_description", offering.LongDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if offering.LongDescriptionI18n != nil {
		if err = d.Set("long_description_i18n", offering.LongDescriptionI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description_i18n: : %s", err), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	features := []map[string]interface{}{}
	if offering.Features != nil {
		for _, modelItem := range offering.Features {
			modelMap, err := dataSourceIBMCmOfferingFeatureToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			features = append(features, modelMap)
		}
	}
	if err = d.Set("features", features); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting features : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	kinds := []map[string]interface{}{}
	if offering.Kinds != nil {
		for _, modelItem := range offering.Kinds {
			modelMap, err := dataSourceIBMCmOfferingKindToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			kinds = append(kinds, modelMap)
		}
	}
	if err = d.Set("kinds", kinds); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kinds : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("pc_managed", offering.PcManaged); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting pc_managed: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("publish_approved", offering.PublishApproved); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish_approved: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("share_with_all", offering.ShareWithAll); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_with_all: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("share_with_ibm", offering.ShareWithIBM); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_with_ibm: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("share_enabled", offering.ShareEnabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting share_enabled: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("public_original_crn", offering.PublicOriginalCRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting public_original_crn: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("publish_public_crn", offering.PublishPublicCRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish_public_crn: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("portal_approval_record", offering.PortalApprovalRecord); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting portal_approval_record: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("portal_ui_url", offering.PortalUIURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting portal_ui_url: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_id", offering.CatalogID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_name", offering.CatalogName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_name: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if offering.Metadata != nil {
		convertedMap := make(map[string]interface{}, len(offering.Metadata))
		for k, v := range offering.Metadata {
			convertedMap[k] = v
		}

		if err = d.Set("metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metadata: : %s", err), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("disclaimer", offering.Disclaimer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disclaimer: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("hidden", offering.Hidden); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting hidden: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	providerInfo := []map[string]interface{}{}
	if offering.ProviderInfo != nil {
		modelMap, err := dataSourceIBMCmOfferingProviderInfoToMap(offering.ProviderInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		providerInfo = append(providerInfo, modelMap)
	}
	if err = d.Set("provider_info", providerInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting provider_info : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	repoInfo := []map[string]interface{}{}
	if offering.RepoInfo != nil {
		modelMap, err := dataSourceIBMCmOfferingRepoInfoToMap(offering.RepoInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		repoInfo = append(repoInfo, modelMap)
	}
	if err = d.Set("repo_info", repoInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting repo_info : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	imagePullKeys := []map[string]interface{}{}
	if offering.ImagePullKeys != nil {
		for _, modelItem := range offering.ImagePullKeys {
			modelMap, err := dataSourceIBMCmOfferingImagePullKeyToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			imagePullKeys = append(imagePullKeys, modelMap)
		}
	}
	if err = d.Set("image_pull_keys", imagePullKeys); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_pull_keys : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	support := []map[string]interface{}{}
	if offering.Support != nil {
		modelMap, err := dataSourceIBMCmOfferingSupportToMap(offering.Support)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		support = append(support, modelMap)
	}
	if err = d.Set("support", support); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting support : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	media := []map[string]interface{}{}
	if offering.Media != nil {
		for _, modelItem := range offering.Media {
			modelMap, err := dataSourceIBMCmOfferingMediaItemToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			media = append(media, modelMap)
		}
	}
	if err = d.Set("media", media); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting media : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deprecatePending := []map[string]interface{}{}
	if offering.DeprecatePending != nil {
		modelMap, err := dataSourceIBMCmOfferingDeprecatePendingToMap(offering.DeprecatePending)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		deprecatePending = append(deprecatePending, modelMap)
	}
	if err = d.Set("deprecate_pending", deprecatePending); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting deprecate_pending : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("product_kind", offering.ProductKind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting product_kind: : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	badges := []map[string]interface{}{}
	if offering.Badges != nil {
		for _, modelItem := range offering.Badges {
			modelMap, err := dataSourceIBMCmOfferingBadgeToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_offering", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			badges = append(badges, modelMap)
		}
	}
	if err = d.Set("badges", badges); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting badges : %s", err), "(Data) ibm_cm_offering", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIBMCmOfferingRatingToMap(model *catalogmanagementv1.Rating) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.OneStarCount != nil {
		modelMap["one_star_count"] = *model.OneStarCount
	}
	if model.TwoStarCount != nil {
		modelMap["two_star_count"] = *model.TwoStarCount
	}
	if model.ThreeStarCount != nil {
		modelMap["three_star_count"] = *model.ThreeStarCount
	}
	if model.FourStarCount != nil {
		modelMap["four_star_count"] = *model.FourStarCount
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
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

func dataSourceIBMCmOfferingKindToMap(model *catalogmanagementv1.Kind) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.FormatKind != nil {
		modelMap["format_kind"] = *model.FormatKind
	}
	if model.InstallKind != nil {
		modelMap["install_kind"] = *model.InstallKind
	}
	if model.TargetKind != nil {
		modelMap["target_kind"] = *model.TargetKind
	}
	if model.Metadata != nil {
		metadataMap := make(map[string]interface{}, len(model.Metadata))
		modelMap["metadata"] = flex.Flatten(metadataMap)
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.AdditionalFeatures != nil {
		additionalFeatures := []map[string]interface{}{}
		for _, additionalFeaturesItem := range model.AdditionalFeatures {
			additionalFeaturesItemMap, err := dataSourceIBMCmOfferingFeatureToMap(&additionalFeaturesItem)
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
			versionsItemMap, err := dataSourceIBMCmOfferingVersionToMap(&versionsItem)
			if err != nil {
				return modelMap, err
			}
			versions = append(versions, versionsItemMap)
		}
		modelMap["versions"] = versions
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingVersionToMap(model *catalogmanagementv1.Version) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Rev != nil {
		modelMap["rev"] = *model.Rev
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Version != nil {
		modelMap["version"] = *model.Version
	}
	if model.Flavor != nil {
		flavorMap, err := dataSourceIBMCmOfferingFlavorToMap(model.Flavor)
		if err != nil {
			return modelMap, err
		}
		modelMap["flavor"] = []map[string]interface{}{flavorMap}
	}
	if model.Sha != nil {
		modelMap["sha"] = *model.Sha
	}
	if model.Created != nil {
		modelMap["created"] = model.Created.String()
	}
	if model.Updated != nil {
		modelMap["updated"] = model.Updated.String()
	}
	if model.OfferingID != nil {
		modelMap["offering_id"] = *model.OfferingID
	}
	if model.CatalogID != nil {
		modelMap["catalog_id"] = *model.CatalogID
	}
	if model.KindID != nil {
		modelMap["kind_id"] = *model.KindID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.RepoURL != nil {
		modelMap["repo_url"] = *model.RepoURL
	}
	if model.SourceURL != nil {
		modelMap["source_url"] = *model.SourceURL
	}
	if model.TgzURL != nil {
		modelMap["tgz_url"] = *model.TgzURL
	}
	if model.Configuration != nil {
		configuration := []map[string]interface{}{}
		for _, configurationItem := range model.Configuration {
			configurationItemMap, err := dataSourceIBMCmOfferingConfigurationToMap(&configurationItem)
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
			outputsItemMap, err := dataSourceIBMCmOfferingOutputToMap(&outputsItem)
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
			iamPermissionsItemMap, err := dataSourceIBMCmOfferingIamPermissionToMap(&iamPermissionsItem)
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
				return nil, err
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
		validationMap, err := dataSourceIBMCmOfferingValidationToMap(model.Validation)
		if err != nil {
			return modelMap, err
		}
		modelMap["validation"] = []map[string]interface{}{validationMap}
	}
	if model.RequiredResources != nil {
		requiredResources := []map[string]interface{}{}
		for _, requiredResourcesItem := range model.RequiredResources {
			requiredResourcesItemMap, err := dataSourceIBMCmOfferingResourceToMap(&requiredResourcesItem)
			if err != nil {
				return modelMap, err
			}
			requiredResources = append(requiredResources, requiredResourcesItemMap)
		}
		modelMap["required_resources"] = requiredResources
	}
	if model.SingleInstance != nil {
		modelMap["single_instance"] = *model.SingleInstance
	}
	if model.Install != nil {
		installMap, err := dataSourceIBMCmOfferingScriptToMap(model.Install)
		if err != nil {
			return modelMap, err
		}
		modelMap["install"] = []map[string]interface{}{installMap}
	}
	if model.PreInstall != nil {
		preInstall := []map[string]interface{}{}
		for _, preInstallItem := range model.PreInstall {
			preInstallItemMap, err := dataSourceIBMCmOfferingScriptToMap(&preInstallItem)
			if err != nil {
				return modelMap, err
			}
			preInstall = append(preInstall, preInstallItemMap)
		}
		modelMap["pre_install"] = preInstall
	}
	if model.Entitlement != nil {
		entitlementMap, err := dataSourceIBMCmOfferingVersionEntitlementToMap(model.Entitlement)
		if err != nil {
			return modelMap, err
		}
		modelMap["entitlement"] = []map[string]interface{}{entitlementMap}
	}
	if model.Licenses != nil {
		licenses := []map[string]interface{}{}
		for _, licensesItem := range model.Licenses {
			licensesItemMap, err := dataSourceIBMCmOfferingLicenseToMap(&licensesItem)
			if err != nil {
				return modelMap, err
			}
			licenses = append(licenses, licensesItemMap)
		}
		modelMap["licenses"] = licenses
	}
	if model.ImageManifestURL != nil {
		modelMap["image_manifest_url"] = *model.ImageManifestURL
	}
	if model.Deprecated != nil {
		modelMap["deprecated"] = *model.Deprecated
	}
	if model.PackageVersion != nil {
		modelMap["package_version"] = *model.PackageVersion
	}
	if model.State != nil {
		stateMap, err := dataSourceIBMCmOfferingStateToMap(model.State)
		if err != nil {
			return modelMap, err
		}
		modelMap["state"] = []map[string]interface{}{stateMap}
	}
	if model.VersionLocator != nil {
		modelMap["version_locator"] = *model.VersionLocator
	}
	if model.LongDescription != nil {
		modelMap["long_description"] = *model.LongDescription
	}
	if model.LongDescriptionI18n != nil {
		longDescriptionI18nMap := make(map[string]interface{}, len(model.LongDescriptionI18n))
		for k, v := range model.LongDescriptionI18n {
			longDescriptionI18nMap[k] = v
		}
		modelMap["long_description_i18n"] = flex.Flatten(longDescriptionI18nMap)
	}
	if model.WhitelistedAccounts != nil {
		modelMap["whitelisted_accounts"] = model.WhitelistedAccounts
	}
	if model.ImagePullKeyName != nil {
		modelMap["image_pull_key_name"] = *model.ImagePullKeyName
	}
	if model.DeprecatePending != nil {
		deprecatePendingMap, err := dataSourceIBMCmOfferingDeprecatePendingToMap(model.DeprecatePending)
		if err != nil {
			return modelMap, err
		}
		modelMap["deprecate_pending"] = []map[string]interface{}{deprecatePendingMap}
	}
	if model.SolutionInfo != nil {
		solutionInfoMap, err := dataSourceIBMCmOfferingSolutionInfoToMap(model.SolutionInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["solution_info"] = []map[string]interface{}{solutionInfoMap}
	}
	if model.IsConsumable != nil {
		modelMap["is_consumable"] = *model.IsConsumable
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingFlavorToMap(model *catalogmanagementv1.Flavor) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	if model.LabelI18n != nil {
		labelI18nMap := make(map[string]interface{}, len(model.LabelI18n))
		for k, v := range model.LabelI18n {
			labelI18nMap[k] = v
		}
		modelMap["label_i18n"] = flex.Flatten(labelI18nMap)
	}
	if model.Index != nil {
		modelMap["index"] = *model.Index
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingConfigurationToMap(model *catalogmanagementv1.Configuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.DefaultValue != nil {
		defaultValueJson, err := json.Marshal(model.DefaultValue)
		if err != nil {
			return nil, flex.FmtErrorf("[ERROR] Error marshalling the version configuration default_value: %s", err)
		}
		defaultValueString, _ := strconv.Unquote(string(defaultValueJson))
		modelMap["default_value"] = defaultValueString
	}
	if model.DisplayName != nil {
		modelMap["display_name"] = *model.DisplayName
	}
	if model.ValueConstraint != nil {
		modelMap["value_constraint"] = *model.ValueConstraint
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Required != nil {
		modelMap["required"] = *model.Required
	}
	if model.Options != nil {
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.CustomConfig != nil {
		customConfigMap, err := dataSourceIBMCmOfferingRenderTypeToMap(model.CustomConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_config"] = []map[string]interface{}{customConfigMap}
	}
	if model.TypeMetadata != nil {
		modelMap["type_metadata"] = *model.TypeMetadata
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingRenderTypeToMap(model *catalogmanagementv1.RenderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Grouping != nil {
		modelMap["grouping"] = *model.Grouping
	}
	if model.OriginalGrouping != nil {
		modelMap["original_grouping"] = *model.OriginalGrouping
	}
	if model.GroupingIndex != nil {
		modelMap["grouping_index"] = *model.GroupingIndex
	}
	if model.ConfigConstraints != nil {
		modelMap["config_constraints"] = flex.Flatten(model.ConfigConstraints)
	}
	if model.Associations != nil {
		associationsMap, err := dataSourceIBMCmOfferingRenderTypeAssociationsToMap(model.Associations)
		if err != nil {
			return modelMap, err
		}
		modelMap["associations"] = []map[string]interface{}{associationsMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingRenderTypeAssociationsToMap(model *catalogmanagementv1.RenderTypeAssociations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := dataSourceIBMCmOfferingRenderTypeAssociationsParametersItemToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingRenderTypeAssociationsParametersItemToMap(model *catalogmanagementv1.RenderTypeAssociationsParametersItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.OptionsRefresh != nil {
		modelMap["options_refresh"] = *model.OptionsRefresh
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingOutputToMap(model *catalogmanagementv1.Output) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingIamPermissionToMap(model *catalogmanagementv1.IamPermission) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServiceName != nil {
		modelMap["service_name"] = *model.ServiceName
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := dataSourceIBMCmOfferingIamResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingIamResourceToMap(model *catalogmanagementv1.IamResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingValidationToMap(model *catalogmanagementv1.Validation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Validated != nil {
		modelMap["validated"] = model.Validated.String()
	}
	if model.Requested != nil {
		modelMap["requested"] = model.Requested.String()
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.LastOperation != nil {
		modelMap["last_operation"] = *model.LastOperation
	}
	if model.Target != nil {
		targetMap := make(map[string]interface{}, len(model.Target))
		modelMap["target"] = flex.Flatten(targetMap)
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingResourceToMap(model *catalogmanagementv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingScriptToMap(model *catalogmanagementv1.Script) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instructions != nil {
		modelMap["instructions"] = *model.Instructions
	}
	if model.InstructionsI18n != nil {
		instructionsI18nMap := make(map[string]interface{}, len(model.InstructionsI18n))
		for k, v := range model.InstructionsI18n {
			instructionsI18nMap[k] = v
		}
		modelMap["instructions_i18n"] = flex.Flatten(instructionsI18nMap)
	}
	if model.Script != nil {
		modelMap["script"] = *model.Script
	}
	if model.ScriptPermission != nil {
		modelMap["script_permission"] = *model.ScriptPermission
	}
	if model.DeleteScript != nil {
		modelMap["delete_script"] = *model.DeleteScript
	}
	if model.Scope != nil {
		modelMap["scope"] = *model.Scope
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingVersionEntitlementToMap(model *catalogmanagementv1.VersionEntitlement) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProviderName != nil {
		modelMap["provider_name"] = *model.ProviderName
	}
	if model.ProviderID != nil {
		modelMap["provider_id"] = *model.ProviderID
	}
	if model.ProductID != nil {
		modelMap["product_id"] = *model.ProductID
	}
	if model.PartNumbers != nil {
		modelMap["part_numbers"] = model.PartNumbers
	}
	if model.ImageRepoName != nil {
		modelMap["image_repo_name"] = *model.ImageRepoName
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingLicenseToMap(model *catalogmanagementv1.License) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Current != nil {
		modelMap["current"] = *model.Current
	}
	if model.CurrentEntered != nil {
		modelMap["current_entered"] = model.CurrentEntered.String()
	}
	if model.Pending != nil {
		modelMap["pending"] = *model.Pending
	}
	if model.PendingRequested != nil {
		modelMap["pending_requested"] = model.PendingRequested.String()
	}
	if model.Previous != nil {
		modelMap["previous"] = *model.Previous
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingDeprecatePendingToMap(model *catalogmanagementv1.DeprecatePending) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DeprecateDate != nil {
		modelMap["deprecate_date"] = model.DeprecateDate.String()
	}
	if model.DeprecateState != nil {
		modelMap["deprecate_state"] = *model.DeprecateState
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingSolutionInfoToMap(model *catalogmanagementv1.SolutionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchitectureDiagrams != nil {
		architectureDiagrams := []map[string]interface{}{}
		for _, architectureDiagramsItem := range model.ArchitectureDiagrams {
			architectureDiagramsItemMap, err := dataSourceIBMCmOfferingArchitectureDiagramToMap(&architectureDiagramsItem)
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
			featuresItemMap, err := dataSourceIBMCmOfferingFeatureToMap(&featuresItem)
			if err != nil {
				return modelMap, err
			}
			features = append(features, featuresItemMap)
		}
		modelMap["features"] = features
	}
	if model.CostEstimate != nil {
		costEstimateMap, err := dataSourceIBMCmOfferingCostEstimateToMap(model.CostEstimate)
		if err != nil {
			return modelMap, err
		}
		modelMap["cost_estimate"] = []map[string]interface{}{costEstimateMap}
	}
	if model.Dependencies != nil {
		dependencies := []map[string]interface{}{}
		for _, dependenciesItem := range model.Dependencies {
			dependenciesItemMap, err := dataSourceIBMCmOfferingDependencyToMap(&dependenciesItem)
			if err != nil {
				return modelMap, err
			}
			dependencies = append(dependencies, dependenciesItemMap)
		}
		modelMap["dependencies"] = dependencies
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingArchitectureDiagramToMap(model *catalogmanagementv1.ArchitectureDiagram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Diagram != nil {
		diagramMap, err := dataSourceIBMCmOfferingMediaItemToMap(model.Diagram)
		if err != nil {
			return modelMap, err
		}
		modelMap["diagram"] = []map[string]interface{}{diagramMap}
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

func dataSourceIBMCmOfferingMediaItemToMap(model *catalogmanagementv1.MediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.APIURL != nil {
		modelMap["api_url"] = *model.APIURL
	}
	if model.URLProxy != nil {
		urlProxyMap, err := dataSourceIBMCmOfferingURLProxyToMap(model.URLProxy)
		if err != nil {
			return modelMap, err
		}
		modelMap["url_proxy"] = []map[string]interface{}{urlProxyMap}
	}
	if model.Caption != nil {
		modelMap["caption"] = *model.Caption
	}
	if model.CaptionI18n != nil {
		captionI18nMap := make(map[string]interface{}, len(model.CaptionI18n))
		for k, v := range model.CaptionI18n {
			captionI18nMap[k] = v
		}
		modelMap["caption_i18n"] = flex.Flatten(captionI18nMap)
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.ThumbnailURL != nil {
		modelMap["thumbnail_url"] = *model.ThumbnailURL
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingURLProxyToMap(model *catalogmanagementv1.URLProxy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Sha != nil {
		modelMap["sha"] = *model.Sha
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingCostEstimateToMap(model *catalogmanagementv1.CostEstimate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Version != nil {
		modelMap["version"] = *model.Version
	}
	if model.Currency != nil {
		modelMap["currency"] = *model.Currency
	}
	if model.Projects != nil {
		projects := []map[string]interface{}{}
		for _, projectsItem := range model.Projects {
			projectsItemMap, err := dataSourceIBMCmOfferingProjectToMap(&projectsItem)
			if err != nil {
				return modelMap, err
			}
			projects = append(projects, projectsItemMap)
		}
		modelMap["projects"] = projects
	}
	if model.Summary != nil {
		summaryMap, err := dataSourceIBMCmOfferingCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = *model.TotalHourlyCost
	}
	if model.TotalMonthlyCost != nil {
		modelMap["total_monthly_cost"] = *model.TotalMonthlyCost
	}
	if model.PastTotalHourlyCost != nil {
		modelMap["past_total_hourly_cost"] = *model.PastTotalHourlyCost
	}
	if model.PastTotalMonthlyCost != nil {
		modelMap["past_total_monthly_cost"] = *model.PastTotalMonthlyCost
	}
	if model.DiffTotalHourlyCost != nil {
		modelMap["diff_total_hourly_cost"] = *model.DiffTotalHourlyCost
	}
	if model.DiffTotalMonthlyCost != nil {
		modelMap["diff_total_monthly_cost"] = *model.DiffTotalMonthlyCost
	}
	if model.TimeGenerated != nil {
		modelMap["time_generated"] = model.TimeGenerated.String()
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingProjectToMap(model *catalogmanagementv1.Project) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Metadata != nil {
		metadataMap := make(map[string]interface{}, len(model.Metadata))
		modelMap["metadata"] = flex.Flatten(metadataMap)
	}
	if model.PastBreakdown != nil {
		pastBreakdownMap, err := dataSourceIBMCmOfferingCostBreakdownToMap(model.PastBreakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["past_breakdown"] = []map[string]interface{}{pastBreakdownMap}
	}
	if model.Breakdown != nil {
		breakdownMap, err := dataSourceIBMCmOfferingCostBreakdownToMap(model.Breakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["breakdown"] = []map[string]interface{}{breakdownMap}
	}
	if model.Diff != nil {
		diffMap, err := dataSourceIBMCmOfferingCostBreakdownToMap(model.Diff)
		if err != nil {
			return modelMap, err
		}
		modelMap["diff"] = []map[string]interface{}{diffMap}
	}
	if model.Summary != nil {
		summaryMap, err := dataSourceIBMCmOfferingCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingCostBreakdownToMap(model *catalogmanagementv1.CostBreakdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = *model.TotalHourlyCost
	}
	if model.TotalMonthlyCost != nil {
		modelMap["total_monthly_c_ost"] = *model.TotalMonthlyCost
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := dataSourceIBMCmOfferingCostResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingCostResourceToMap(model *catalogmanagementv1.CostResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Metadata != nil {
		metadataMap := make(map[string]interface{}, len(model.Metadata))
		modelMap["metadata"] = flex.Flatten(metadataMap)
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = *model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = *model.MonthlyCost
	}
	if model.CostComponents != nil {
		costComponents := []map[string]interface{}{}
		for _, costComponentsItem := range model.CostComponents {
			costComponentsItemMap, err := dataSourceIBMCmOfferingCostComponentToMap(&costComponentsItem)
			if err != nil {
				return modelMap, err
			}
			costComponents = append(costComponents, costComponentsItemMap)
		}
		modelMap["cost_components"] = costComponents
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingCostComponentToMap(model *catalogmanagementv1.CostComponent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	if model.HourlyQuantity != nil {
		modelMap["hourly_quantity"] = *model.HourlyQuantity
	}
	if model.MonthlyQuantity != nil {
		modelMap["monthly_quantity"] = *model.MonthlyQuantity
	}
	if model.Price != nil {
		modelMap["price"] = *model.Price
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = *model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = *model.MonthlyCost
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingCostSummaryToMap(model *catalogmanagementv1.CostSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalDetectedResources != nil {
		modelMap["total_detected_resources"] = *model.TotalDetectedResources
	}
	if model.TotalSupportedResources != nil {
		modelMap["total_supported_resources"] = *model.TotalSupportedResources
	}
	if model.TotalUnsupportedResources != nil {
		modelMap["total_unsupported_resources"] = *model.TotalUnsupportedResources
	}
	if model.TotalUsageBasedResources != nil {
		modelMap["total_usage_based_resources"] = *model.TotalUsageBasedResources
	}
	if model.TotalNoPriceResources != nil {
		modelMap["total_no_price_resources"] = *model.TotalNoPriceResources
	}
	if model.UnsupportedResourceCounts != nil {
		unsupportedResourceCountsMap := make(map[string]interface{}, len(model.UnsupportedResourceCounts))
		for k, v := range model.UnsupportedResourceCounts {
			unsupportedResourceCountsMap[k] = v
		}
		modelMap["unsupported_resource_counts"] = flex.Flatten(unsupportedResourceCountsMap)
	}
	if model.NoPriceResourceCounts != nil {
		noPriceResourceCountsMap := make(map[string]interface{}, len(model.NoPriceResourceCounts))
		for k, v := range model.NoPriceResourceCounts {
			noPriceResourceCountsMap[k] = v
		}
		modelMap["no_price_resource_counts"] = flex.Flatten(noPriceResourceCountsMap)
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingDependencyToMap(model *catalogmanagementv1.OfferingReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CatalogID != nil {
		modelMap["catalog_id"] = *model.CatalogID
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Version != nil {
		modelMap["version"] = *model.Version
	}
	if model.Flavors != nil {
		modelMap["flavors"] = model.Flavors
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingProviderInfoToMap(model *catalogmanagementv1.ProviderInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingRepoInfoToMap(model *catalogmanagementv1.RepoInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Token != nil {
		modelMap["token"] = *model.Token
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingImagePullKeyToMap(model *catalogmanagementv1.ImagePullKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingSupportToMap(model *catalogmanagementv1.Support) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Process != nil {
		modelMap["process"] = *model.Process
	}
	if model.ProcessI18n != nil {
		processI18nMap := make(map[string]interface{}, len(model.ProcessI18n))
		for k, v := range model.ProcessI18n {
			processI18nMap[k] = v
		}
		modelMap["process_i18n"] = flex.Flatten(processI18nMap)
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	if model.SupportDetails != nil {
		supportDetails := []map[string]interface{}{}
		for _, supportDetailsItem := range model.SupportDetails {
			supportDetailsItemMap, err := dataSourceIBMCmOfferingSupportDetailToMap(&supportDetailsItem)
			if err != nil {
				return modelMap, err
			}
			supportDetails = append(supportDetails, supportDetailsItemMap)
		}
		modelMap["support_details"] = supportDetails
	}
	if model.SupportEscalation != nil {
		supportEscalationMap, err := dataSourceIBMCmOfferingSupportEscalationToMap(model.SupportEscalation)
		if err != nil {
			return modelMap, err
		}
		modelMap["support_escalation"] = []map[string]interface{}{supportEscalationMap}
	}
	if model.SupportType != nil {
		modelMap["support_type"] = *model.SupportType
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingSupportDetailToMap(model *catalogmanagementv1.SupportDetail) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Contact != nil {
		modelMap["contact"] = *model.Contact
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := dataSourceIBMCmOfferingSupportWaitTimeToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	if model.Availability != nil {
		availabilityMap, err := dataSourceIBMCmOfferingSupportAvailabilityToMap(model.Availability)
		if err != nil {
			return modelMap, err
		}
		modelMap["availability"] = []map[string]interface{}{availabilityMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingSupportWaitTimeToMap(model *catalogmanagementv1.SupportWaitTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingSupportAvailabilityToMap(model *catalogmanagementv1.SupportAvailability) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Times != nil {
		times := []map[string]interface{}{}
		for _, timesItem := range model.Times {
			timesItemMap, err := dataSourceIBMCmOfferingSupportTimeToMap(&timesItem)
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

func dataSourceIBMCmOfferingSupportTimeToMap(model *catalogmanagementv1.SupportTime) (map[string]interface{}, error) {
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

func dataSourceIBMCmOfferingSupportEscalationToMap(model *catalogmanagementv1.SupportEscalation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EscalationWaitTime != nil {
		escalationWaitTimeMap, err := dataSourceIBMCmOfferingSupportWaitTimeToMap(model.EscalationWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["escalation_wait_time"] = []map[string]interface{}{escalationWaitTimeMap}
	}
	if model.ResponseWaitTime != nil {
		responseWaitTimeMap, err := dataSourceIBMCmOfferingSupportWaitTimeToMap(model.ResponseWaitTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["response_wait_time"] = []map[string]interface{}{responseWaitTimeMap}
	}
	if model.Contact != nil {
		modelMap["contact"] = *model.Contact
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingBadgeToMap(model *catalogmanagementv1.Badge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	if model.LabelI18n != nil {
		labelI18nMap := make(map[string]interface{}, len(model.LabelI18n))
		for k, v := range model.LabelI18n {
			labelI18nMap[k] = v
		}
		modelMap["label_i18n"] = flex.Flatten(labelI18nMap)
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
	if model.Icon != nil {
		modelMap["icon"] = *model.Icon
	}
	if model.Authority != nil {
		modelMap["authority"] = *model.Authority
	}
	if model.Tag != nil {
		modelMap["tag"] = *model.Tag
	}
	if model.LearnMoreLinks != nil {
		learnMoreLinksMap, err := dataSourceIBMCmOfferingLearnMoreLinksToMap(model.LearnMoreLinks)
		if err != nil {
			return modelMap, err
		}
		modelMap["learn_more_links"] = []map[string]interface{}{learnMoreLinksMap}
	}
	if model.Constraints != nil {
		constraints := []map[string]interface{}{}
		for _, constraintsItem := range model.Constraints {
			constraintsItemMap, err := dataSourceIBMCmOfferingConstraintToMap(&constraintsItem)
			if err != nil {
				return modelMap, err
			}
			constraints = append(constraints, constraintsItemMap)
		}
		modelMap["constraints"] = constraints
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingLearnMoreLinksToMap(model *catalogmanagementv1.LearnMoreLinks) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FirstParty != nil {
		modelMap["first_party"] = *model.FirstParty
	}
	if model.ThirdParty != nil {
		modelMap["third_party"] = *model.ThirdParty
	}
	return modelMap, nil
}

func dataSourceIBMCmOfferingConstraintToMap(model *catalogmanagementv1.Constraint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Rule != nil {
	}
	return modelMap, nil
}
