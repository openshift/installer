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

func DataSourceIBMCmVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmVersionRead,

		Schema: map[string]*schema.Schema{
			"version_loc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A dotted value of `catalogID`.`versionID`.",
			},
			"version_id": &schema.Schema{
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"whitelisted_accounts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Whitelisted accounts for version.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
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
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
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
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
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
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
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
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"no_price_resource_counts": &schema.Schema{
																Type:        schema.TypeMap,
																Computed:    true,
																Description: "No price resource counts.",
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"no_price_resource_counts": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "No price resource counts.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
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
			"is_consumable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the version able to be shared.",
			},
		},
	}
}

func dataSourceIBMCmVersionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

	getVersionOptions.SetVersionLocID(d.Get("version_loc_id").(string))

	offering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVersionWithContext failed %s\n%s", err, response), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	version := offering.Kinds[0].Versions[0]

	d.SetId(*getVersionOptions.VersionLocID)

	if err = d.Set("version_id", version.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_id: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("rev", version.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", version.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("version", version.Version); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	flavor := []map[string]interface{}{}
	if version.Flavor != nil {
		modelMap, err := dataSourceIBMCmVersionFlavorToMap(version.Flavor)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		flavor = append(flavor, modelMap)
	}
	if err = d.Set("flavor", flavor); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting flavor: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("sha", version.Sha); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting sha: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created", flex.DateTimeToString(version.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated", flex.DateTimeToString(version.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offering_id", version.OfferingID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting offering_id: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_id", version.CatalogID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("kind_id", version.KindID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kind_id: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("repo_url", version.RepoURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting repo_url: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("source_url", version.SourceURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting source_url: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("tgz_url", version.TgzURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting tgz_url: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	configuration := []map[string]interface{}{}
	if version.Configuration != nil {
		for _, modelItem := range version.Configuration {
			modelMap, err := dataSourceIBMCmVersionConfigurationToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			configuration = append(configuration, modelMap)
		}
	}
	if err = d.Set("configuration", configuration); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting configuration: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	outputs := []map[string]interface{}{}
	if version.Outputs != nil {
		for _, modelItem := range version.Outputs {
			modelMap, err := dataSourceIBMCmVersionOutputToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			outputs = append(outputs, modelMap)
		}
	}
	if err = d.Set("outputs", outputs); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting outputs: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	iamPermissions := []map[string]interface{}{}
	if version.IamPermissions != nil {
		for _, modelItem := range version.IamPermissions {
			modelMap, err := dataSourceIBMCmVersionIamPermissionToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			iamPermissions = append(iamPermissions, modelMap)
		}
	}
	if err = d.Set("iam_permissions", iamPermissions); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting iam_permissions: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	metadata := []map[string]interface{}{}
	if version.Metadata != nil {
		var modelMapVSI map[string]interface{}
		if version.Metadata["vsi_vpc"] != nil {
			modelMapVSI, err = dataSourceIBMCmVersionMetadataVSIToMap(version.Metadata["vsi_vpc"].(map[string]interface{}))
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
		convertedMap := make(map[string]interface{}, len(version.Metadata))
		for k, v := range version.Metadata {
			if k == "vsi_vpc" {
				convertedMap[k] = []map[string]interface{}{modelMapVSI}
			} else {
				convertedMap[k] = v
			}
		}
		metadata = append(metadata, convertedMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metadata: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	validation := []map[string]interface{}{}
	if version.Validation != nil {
		modelMap, err := dataSourceIBMCmVersionValidationToMap(version.Validation)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		validation = append(validation, modelMap)
	}
	if err = d.Set("validation", validation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validation: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	requiredResources := []map[string]interface{}{}
	if version.RequiredResources != nil {
		for _, modelItem := range version.RequiredResources {
			modelMap, err := dataSourceIBMCmVersionResourceToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			requiredResources = append(requiredResources, modelMap)
		}
	}
	if err = d.Set("required_resources", requiredResources); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting required_resources: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("single_instance", version.SingleInstance); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting single_instance: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	install := []map[string]interface{}{}
	if version.Install != nil {
		modelMap, err := dataSourceIBMCmVersionScriptToMap(version.Install)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		install = append(install, modelMap)
	}
	if err = d.Set("install", install); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting install: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	preInstall := []map[string]interface{}{}
	if version.PreInstall != nil {
		for _, modelItem := range version.PreInstall {
			modelMap, err := dataSourceIBMCmVersionScriptToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			preInstall = append(preInstall, modelMap)
		}
	}
	if err = d.Set("pre_install", preInstall); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting pre_install: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	entitlement := []map[string]interface{}{}
	if version.Entitlement != nil {
		modelMap, err := dataSourceIBMCmVersionVersionEntitlementToMap(version.Entitlement)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		entitlement = append(entitlement, modelMap)
	}
	if err = d.Set("entitlement", entitlement); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting entitlement: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	licenses := []map[string]interface{}{}
	if version.Licenses != nil {
		for _, modelItem := range version.Licenses {
			modelMap, err := dataSourceIBMCmVersionLicenseToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			licenses = append(licenses, modelMap)
		}
	}
	if err = d.Set("licenses", licenses); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting licenses: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("image_manifest_url", version.ImageManifestURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_manifest_url: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("deprecated", version.Deprecated); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting deprecated: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("package_version", version.PackageVersion); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting package_version: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	state := []map[string]interface{}{}
	if version.State != nil {
		modelMap, err := dataSourceIBMCmVersionStateToMap(version.State)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		state = append(state, modelMap)
	}
	if err = d.Set("state", state); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("version_locator", version.VersionLocator); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_locator: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("long_description", version.LongDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if version.LongDescriptionI18n != nil {
		if err = d.Set("long_description_i18n", version.LongDescriptionI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting long_description_i18n: %s", err), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("image_pull_key_name", version.ImagePullKeyName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_pull_key_name: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deprecatePending := []map[string]interface{}{}
	if version.DeprecatePending != nil {
		modelMap, err := dataSourceIBMCmVersionDeprecatePendingToMap(version.DeprecatePending)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		deprecatePending = append(deprecatePending, modelMap)
	}
	if err = d.Set("deprecate_pending", deprecatePending); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting deprecate_pending: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	solutionInfo := []map[string]interface{}{}
	if version.SolutionInfo != nil {
		modelMap, err := dataSourceIBMCmVersionSolutionInfoToMap(version.SolutionInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_version", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		solutionInfo = append(solutionInfo, modelMap)
	}
	if err = d.Set("solution_info", solutionInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting solution_info: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("is_consumable", version.IsConsumable); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting is_consumable: %s", err), "(Data) ibm_cm_version", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIBMCmVersionFlavorToMap(model *catalogmanagementv1.Flavor) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionConfigurationToMap(model *catalogmanagementv1.Configuration) (map[string]interface{}, error) {
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
			return nil, fmt.Errorf("[ERROR] Error marshalling the version configuration default_value: %s", err)
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
		customConfigMap, err := dataSourceIBMCmVersionRenderTypeToMap(model.CustomConfig)
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

func dataSourceIBMCmVersionRenderTypeToMap(model *catalogmanagementv1.RenderType) (map[string]interface{}, error) {
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
		associationsMap, err := dataSourceIBMCmVersionRenderTypeAssociationsToMap(model.Associations)
		if err != nil {
			return modelMap, err
		}
		modelMap["associations"] = []map[string]interface{}{associationsMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionRenderTypeAssociationsToMap(model *catalogmanagementv1.RenderTypeAssociations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := dataSourceIBMCmVersionRenderTypeAssociationsParametersItemToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionRenderTypeAssociationsParametersItemToMap(model *catalogmanagementv1.RenderTypeAssociationsParametersItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.OptionsRefresh != nil {
		modelMap["options_refresh"] = *model.OptionsRefresh
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionOutputToMap(model *catalogmanagementv1.Output) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionIamPermissionToMap(model *catalogmanagementv1.IamPermission) (map[string]interface{}, error) {
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
			resourcesItemMap, err := dataSourceIBMCmVersionIamResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionIamResourceToMap(model *catalogmanagementv1.IamResource) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionValidationToMap(model *catalogmanagementv1.Validation) (map[string]interface{}, error) {
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
		// for k, v := range model.Target {
		// }
		modelMap["target"] = flex.Flatten(targetMap)
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionResourceToMap(model *catalogmanagementv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionScriptToMap(model *catalogmanagementv1.Script) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionVersionEntitlementToMap(model *catalogmanagementv1.VersionEntitlement) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionLicenseToMap(model *catalogmanagementv1.License) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionDeprecatePendingToMap(model *catalogmanagementv1.DeprecatePending) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionSolutionInfoToMap(model *catalogmanagementv1.SolutionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchitectureDiagrams != nil {
		architectureDiagrams := []map[string]interface{}{}
		for _, architectureDiagramsItem := range model.ArchitectureDiagrams {
			architectureDiagramsItemMap, err := dataSourceIBMCmVersionArchitectureDiagramToMap(&architectureDiagramsItem)
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
			featuresItemMap, err := dataSourceIBMCmVersionFeatureToMap(&featuresItem)
			if err != nil {
				return modelMap, err
			}
			features = append(features, featuresItemMap)
		}
		modelMap["features"] = features
	}
	if model.CostEstimate != nil {
		costEstimateMap, err := dataSourceIBMCmVersionCostEstimateToMap(model.CostEstimate)
		if err != nil {
			return modelMap, err
		}
		modelMap["cost_estimate"] = []map[string]interface{}{costEstimateMap}
	}
	if model.Dependencies != nil {
		dependencies := []map[string]interface{}{}
		for _, dependenciesItem := range model.Dependencies {
			dependenciesItemMap, err := dataSourceIBMCmVersionDependencyToMap(&dependenciesItem)
			if err != nil {
				return modelMap, err
			}
			dependencies = append(dependencies, dependenciesItemMap)
		}
		modelMap["dependencies"] = dependencies
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionArchitectureDiagramToMap(model *catalogmanagementv1.ArchitectureDiagram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Diagram != nil {
		diagramMap, err := dataSourceIBMCmVersionMediaItemToMap(model.Diagram)
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

func dataSourceIBMCmVersionMediaItemToMap(model *catalogmanagementv1.MediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.APIURL != nil {
		modelMap["api_url"] = *model.APIURL
	}
	if model.URLProxy != nil {
		urlProxyMap, err := dataSourceIBMCmVersionURLProxyToMap(model.URLProxy)
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

func dataSourceIBMCmVersionURLProxyToMap(model *catalogmanagementv1.URLProxy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Sha != nil {
		modelMap["sha"] = *model.Sha
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionCostEstimateToMap(model *catalogmanagementv1.CostEstimate) (map[string]interface{}, error) {
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
			projectsItemMap, err := dataSourceIBMCmVersionProjectToMap(&projectsItem)
			if err != nil {
				return modelMap, err
			}
			projects = append(projects, projectsItemMap)
		}
		modelMap["projects"] = projects
	}
	if model.Summary != nil {
		summaryMap, err := dataSourceIBMCmVersionCostSummaryToMap(model.Summary)
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

func dataSourceIBMCmVersionProjectToMap(model *catalogmanagementv1.Project) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Metadata != nil {
		metadataMap := make(map[string]interface{}, len(model.Metadata))
		// for k, v := range model.Metadata {
		// }
		modelMap["metadata"] = flex.Flatten(metadataMap)
	}
	if model.PastBreakdown != nil {
		pastBreakdownMap, err := dataSourceIBMCmVersionCostBreakdownToMap(model.PastBreakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["past_breakdown"] = []map[string]interface{}{pastBreakdownMap}
	}
	if model.Breakdown != nil {
		breakdownMap, err := dataSourceIBMCmVersionCostBreakdownToMap(model.Breakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["breakdown"] = []map[string]interface{}{breakdownMap}
	}
	if model.Diff != nil {
		diffMap, err := dataSourceIBMCmVersionCostBreakdownToMap(model.Diff)
		if err != nil {
			return modelMap, err
		}
		modelMap["diff"] = []map[string]interface{}{diffMap}
	}
	if model.Summary != nil {
		summaryMap, err := dataSourceIBMCmVersionCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionCostBreakdownToMap(model *catalogmanagementv1.CostBreakdown) (map[string]interface{}, error) {
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
			resourcesItemMap, err := dataSourceIBMCmVersionCostResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionCostResourceToMap(model *catalogmanagementv1.CostResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Metadata != nil {
		metadataMap := make(map[string]interface{}, len(model.Metadata))
		// for k, v := range model.Metadata {
		// }
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
			costComponentsItemMap, err := dataSourceIBMCmVersionCostComponentToMap(&costComponentsItem)
			if err != nil {
				return modelMap, err
			}
			costComponents = append(costComponents, costComponentsItemMap)
		}
		modelMap["cost_components"] = costComponents
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionCostComponentToMap(model *catalogmanagementv1.CostComponent) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionCostSummaryToMap(model *catalogmanagementv1.CostSummary) (map[string]interface{}, error) {
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

func dataSourceIBMCmVersionDependencyToMap(model *catalogmanagementv1.OfferingReference) (map[string]interface{}, error) {
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
