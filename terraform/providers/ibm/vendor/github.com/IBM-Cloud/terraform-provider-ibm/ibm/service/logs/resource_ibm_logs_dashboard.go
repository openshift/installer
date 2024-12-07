// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsDashboardCreate,
		ReadContext:   resourceIbmLogsDashboardRead,
		UpdateContext: resourceIbmLogsDashboardUpdate,
		DeleteContext: resourceIbmLogsDashboardDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"href": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_dashboard", "href"),
				Description:  "Unique identifier for the dashboard.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_dashboard", "name"),
				Description:  "Display name of the dashboard.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_dashboard", "description"),
				Description:  "Brief description or summary of the dashboard's purpose or content.",
			},
			"layout": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Layout configuration for the dashboard's visual elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sections": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The sections of the layout.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique identifier of the section within the layout.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Unique identifier of the folder containing the dashboard.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The UUID value.",
												},
											},
										},
									},
									"rows": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The rows of the section.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The unique identifier of the row within the layout.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Unique identifier of the folder containing the dashboard.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The UUID value.",
															},
														},
													},
												},
												"appearance": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "The appearance of the row, such as height.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"height": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The height of the row.",
															},
														},
													},
												},
												"widgets": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The widgets of the row.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Widget identifier within the dashboard.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Unique identifier of the folder containing the dashboard.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The UUID value.",
																		},
																	},
																},
															},
															"title": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Widget title.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Widget description.",
															},
															"definition": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Widget definition, contains the widget type and its configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"line_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Line chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"legend": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Legend configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Required:    true,
																									Description: "Whether to show the legend or not.",
																								},
																								"columns": &schema.Schema{
																									Type:        schema.TypeList,
																									Optional:    true,
																									Computed:    true,
																									Description: "The columns to show in the legend.",
																									Elem:        &schema.Schema{Type: schema.TypeString},
																								},
																								"group_by_query": &schema.Schema{
																									Type:        schema.TypeBool,
																									Required:    true,
																									Description: "Whether to group by the query or not.",
																								},
																							},
																						},
																					},
																					"tooltip": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Tooltip configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"show_labels": &schema.Schema{
																									Type:        schema.TypeBool,
																									Optional:    true,
																									Description: "Whether to show labels in the tooltip.",
																								},
																								"type": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Tooltip type.",
																								},
																							},
																						},
																					},
																					"query_definitions": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Query definitions.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"id": &schema.Schema{
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Unique identifier of the query within the widget.",
																								},
																								"query": &schema.Schema{
																									Type:        schema.TypeList,
																									MinItems:    1,
																									MaxItems:    1,
																									Required:    true,
																									Description: "Data source specific query, defines from where and how to fetch the data.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"logs": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Logs specific query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"lucene_query": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"value": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"group_by": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Group by fields (deprecated).",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"aggregations": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Aggregations.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"count": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Count the number of entries.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{},
																																		},
																																	},
																																	"count_distinct": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Count the number of distinct values of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"sum": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Sum values of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"average": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Calculate average value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"min": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Calculate minimum value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"max": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Calculate maximum value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"percentile": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Calculate percentile value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"percent": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Required:    true,
																																					Description: "Value in range (0, 100].",
																																				},
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MinItems:    1,
																																					MaxItems:    1,
																																					Required:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
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
																														"filters": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"group_bys": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Group by fields.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"metrics": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Metrics specific query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"promql_query": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "PromQL query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"value": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"filters": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"dataprime": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Dataprime language based query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"dataprime_query": &schema.Schema{
																															Type:        schema.TypeList,
																															MinItems:    1,
																															MaxItems:    1,
																															Required:    true,
																															Description: "Dataprime query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"text": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"filters": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Filters to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"logs": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Extra filtering on top of the Lucene query.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"operator": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Operator to use for filtering the logs.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Selection criteria for the equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"all": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection of all values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{},
																																														},
																																													},
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Optional:    true,
																																																	Description: "List of values for the selection.",
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
																																							"not_equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Selection criteria for the non-equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Optional:    true,
																																																	Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Optional:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem:        &schema.Schema{Type: schema.TypeString},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Optional:    true,
																																								Description: "Scope of the dataset.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"metrics": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Filtering to be applied to query results.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"label": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Label associated with the metric.",
																																				},
																																				"operator": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Operator to use for filtering the logs.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Selection criteria for the equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"all": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection of all values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{},
																																														},
																																													},
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Optional:    true,
																																																	Description: "List of values for the selection.",
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
																																							"not_equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Selection criteria for the non-equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														MaxItems:    1,
																																														Optional:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Optional:    true,
																																																	Description: "List of values for the selection.",
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
																									},
																								},
																								"series_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Template for series name in legend and tooltip.",
																								},
																								"series_count_limit": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Maximum number of series to display.",
																								},
																								"unit": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Computed:    true,
																									Description: "Unit of the data.",
																								},
																								"scale_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Scale type.",
																								},
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Query name.",
																								},
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Required:    true,
																									Description: "Whether data for this query should be visible on the chart.",
																								},
																								"color_scheme": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Color scheme for the series.",
																								},
																								"resolution": &schema.Schema{
																									Type:        schema.TypeList,
																									MinItems:    1,
																									MaxItems:    1,
																									Required:    true,
																									Description: "Resolution of the data.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"interval": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Interval between data points.",
																											},
																											"buckets_presented": &schema.Schema{
																												Type:        schema.TypeInt,
																												Optional:    true,
																												Description: "Maximum number of data points to fetch.",
																											},
																										},
																									},
																								},
																								"data_mode_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Computed:    true,
																									Description: "Data mode type.",
																								},
																							},
																						},
																					},
																					"stacked_line": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Stacked lines.",
																					},
																				},
																			},
																		},
																		"data_table": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Data table widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filtering on top of the Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"grouping": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Grouping and aggregation.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"aggregations": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Aggregations.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"id": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Required:    true,
																																		Description: "Aggregation identifier, must be unique within grouping configuration.",
																																	},
																																	"name": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Required:    true,
																																		Description: "Aggregation name, used as column name.",
																																	},
																																	"is_visible": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Required:    true,
																																		Description: "Whether the aggregation is visible.",
																																	},
																																	"aggregation": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Aggregations.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"count": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Count the number of entries.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{},
																																					},
																																				},
																																				"count_distinct": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Count the number of distinct values of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
																																										},
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"sum": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Sum values of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
																																										},
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"average": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Calculate average value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
																																										},
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"min": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Calculate minimum value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
																																										},
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"max": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Calculate maximum value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
																																										},
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"percentile": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Calculate percentile value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"percent": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Required:    true,
																																								Description: "Value in range (0, 100].",
																																							},
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MinItems:    1,
																																								MaxItems:    1,
																																								Required:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem:        &schema.Schema{Type: schema.TypeString},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Optional:    true,
																																											Description: "Scope of the dataset.",
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
																														"group_bys": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Fields to group by.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
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
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filtering on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filtering on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"metrics": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																						},
																					},
																					"results_per_page": &schema.Schema{
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "Number of results per page.",
																					},
																					"row_style": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Display style for table rows.",
																					},
																					"columns": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Columns to display, their order and width.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"field": &schema.Schema{
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "References a field in result set. In case of aggregation, it references the aggregation identifier.",
																								},
																								"width": &schema.Schema{
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: "Column width.",
																								},
																							},
																						},
																					},
																					"order_by": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Column used for ordering the results.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"field": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The field to order by.",
																								},
																								"order_direction": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The direction of the order: ascending or descending.",
																								},
																							},
																						},
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"gauge": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Gauge widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "Aggregation. When AGGREGATION_UNSPECIFIED is selected, widget uses instant query. Otherwise, it uses range query.",
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters applied on top of PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"logs_aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"sum": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"average": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"min": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"max": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"percentile": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Required:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
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
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters applied on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
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
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters applied on top of Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"metrics": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																						},
																					},
																					"min": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Required:    true,
																						Description: "Minimum value of the gauge.",
																					},
																					"max": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Required:    true,
																						Description: "Maximum value of the gauge.",
																					},
																					"show_inner_arc": &schema.Schema{
																						Type:        schema.TypeBool,
																						Required:    true,
																						Description: "Show inner arc (styling).",
																					},
																					"show_outer_arc": &schema.Schema{
																						Type:        schema.TypeBool,
																						Required:    true,
																						Description: "Show outer arc (styling).",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Query result value interpretation.",
																					},
																					"thresholds": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Thresholds for the gauge, values at which the gauge changes color.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"from": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Required:    true,
																									Description: "Value at which the color should change.",
																								},
																								"color": &schema.Schema{
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Color.",
																								},
																							},
																						},
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Data mode type.",
																					},
																					"threshold_by": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "What threshold color should be applied to: value or background.",
																					},
																				},
																			},
																		},
																		"pie_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Pie chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"sum": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"average": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"min": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"max": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"percentile": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Required:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
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
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names_fields": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters on top of PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Field to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filters on top of Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"metrics": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Field to stack by.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"max_slices_per_chart": &schema.Schema{
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "Maximum number of slices to display in the chart.",
																					},
																					"min_slice_percentage": &schema.Schema{
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "Minimum percentage of a slice to be displayed.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_stack": &schema.Schema{
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: "Maximum number of slices per stack.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Template for stack labels.",
																								},
																							},
																						},
																					},
																					"label_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Label settings.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"label_source": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Source of the label.",
																								},
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Optional:    true,
																									Description: "Controls whether to show the label.",
																								},
																								"show_name": &schema.Schema{
																									Type:        schema.TypeBool,
																									Optional:    true,
																									Description: "Controls whether to show the name.",
																								},
																								"show_value": &schema.Schema{
																									Type:        schema.TypeBool,
																									Optional:    true,
																									Description: "Controls whether to show the value.",
																								},
																								"show_percentage": &schema.Schema{
																									Type:        schema.TypeBool,
																									Optional:    true,
																									Description: "Controls whether to show the percentage.",
																								},
																							},
																						},
																					},
																					"show_legend": &schema.Schema{
																						Type:        schema.TypeBool,
																						Required:    true,
																						Description: "Controls whether to show the legend.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Template for group labels.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Unit of the data.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Color scheme name.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"bar_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Bar chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"sum": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"average": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"min": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"max": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"percentile": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Required:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
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
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names_fields": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fiel to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Labels to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Label to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MinItems:    1,
																												MaxItems:    1,
																												Required:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"metrics": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Field to stack by.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"max_bars_per_chart": &schema.Schema{
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "Maximum number of bars to present in the chart.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Template for bar labels.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_bar": &schema.Schema{
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: "Maximum number of slices per bar.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Template for stack slice label.",
																								},
																							},
																						},
																					},
																					"scale_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Scale type.",
																					},
																					"colors_by": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Coloring mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"stack": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each stack will have the same color across all groups.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"group_by": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each group will have different color and stack color will be derived from group color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"aggregation": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each aggregation will have different color and stack color will be derived from aggregation color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																							},
																						},
																					},
																					"x_axis": &schema.Schema{
																						Type:        schema.TypeList,
																						MinItems:    1,
																						MaxItems:    1,
																						Required:    true,
																						Description: "X axis mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"value": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Categorical axis.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"time": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Time based axis.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"interval": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Time interval.",
																											},
																											"buckets_presented": &schema.Schema{
																												Type:        schema.TypeInt,
																												Optional:    true,
																												Description: "Buckets presented.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Unit of the data.",
																					},
																					"sort_by": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Sorting mode.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Supported vaues: classic, severity, cold, negative, green, red, blue.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"horizontal_bar_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Horizontal bar chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"sum": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"average": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"min": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"max": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"percentile": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Required:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MinItems:    1,
																																		MaxItems:    1,
																																		Required:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
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
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of the Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Optional:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem:        &schema.Schema{Type: schema.TypeString},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Scope of the dataset.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names_fields": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Optional:    true,
																															Description: "Path within the dataset scope.",
																															Elem:        &schema.Schema{Type: schema.TypeString},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Optional:    true,
																																											Description: "List of values for the selection.",
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
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Labels to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Label to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Dataprime specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												MaxItems:    1,
																												Optional:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Optional:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Extra filter on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Optional:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem:        &schema.Schema{Type: schema.TypeString},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Optional:    true,
																																					Description: "Scope of the dataset.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"metrics": &schema.Schema{
																															Type:        schema.TypeList,
																															MaxItems:    1,
																															Optional:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Optional:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		MaxItems:    1,
																																		Optional:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					MaxItems:    1,
																																					Optional:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								MaxItems:    1,
																																								Optional:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											MaxItems:    1,
																																											Optional:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Optional:    true,
																																														Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"group_names": &schema.Schema{
																												Type:        schema.TypeList,
																												Optional:    true,
																												Description: "Fields to group by.",
																												Elem:        &schema.Schema{Type: schema.TypeString},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "Field to stack by.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"max_bars_per_chart": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
																						Description: "Maximum number of bars to display in the chart.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Template for bar labels.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_bar": &schema.Schema{
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: "Maximum number of slices per bar.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Template for stack slice label.",
																								},
																							},
																						},
																					},
																					"scale_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scale type.",
																					},
																					"colors_by": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Coloring mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"stack": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each stack will have the same color across all groups.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"group_by": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each group will have different color and stack color will be derived from group color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"aggregation": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "Each aggregation will have different color and stack color will be derived from aggregation color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																							},
																						},
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Unit of the data.",
																					},
																					"display_on_bar": &schema.Schema{
																						Type:        schema.TypeBool,
																						Optional:    true,
																						Description: "Whether to display values on the bars.",
																					},
																					"y_axis_view_by": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Y-axis view mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"category": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "View by category.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"value": &schema.Schema{
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "View by value.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																							},
																						},
																					},
																					"sort_by": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Sorting mode.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Color scheme name.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"markdown": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Markdown widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"markdown_text": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Markdown text to render.",
																					},
																					"tooltip_text": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Tooltip text on hover.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"created_at": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Creation timestamp.",
															},
															"updated_at": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Last update timestamp.",
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
			"variables": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of variables that can be used within the dashboard for dynamic content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the variable which can be used in templates.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Definition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"multi_select": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Multi-select value.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Variable value source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"logs_path": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Unique values for a given logs path.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"observation_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"metric_label": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Unique values for a given metric label.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Metric name to source unique values from.",
																		},
																		"label": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Metric label to source unique values from.",
																		},
																	},
																},
															},
															"constant_list": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "List of constant values.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"values": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "List of constant values.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
														},
													},
												},
												"selection": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "State of what is currently selected.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"all": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "All values are selected, usually translated to wildcard (*).",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"list": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specific values are selected.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"values": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Selected values.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
														},
													},
												},
												"values_order_direction": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The direction of the order: ascending or descending.",
												},
											},
										},
									},
								},
							},
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name used in variable UI.",
						},
					},
				},
			},
			"filters": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of filters that can be applied to the dashboard's data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Filters to be applied to query results.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logs": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Extra filtering on top of the Lucene query.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Operator to use for filtering the logs.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"equals": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Selection criteria for the equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"all": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection of all values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{},
																						},
																					},
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Optional:    true,
																									Description: "List of values for the selection.",
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
															"not_equals": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Non-equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Selection criteria for the non-equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Optional:    true,
																									Description: "List of values for the selection.",
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
														},
													},
												},
												"observation_field": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Field to count distinct values of.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keypath": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Path within the dataset scope.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"scope": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Scope of the dataset.",
															},
														},
													},
												},
											},
										},
									},
									"metrics": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Filtering to be applied to query results.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Label associated with the metric.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Operator to use for filtering the logs.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"equals": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Selection criteria for the equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"all": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection of all values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{},
																						},
																					},
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Optional:    true,
																									Description: "List of values for the selection.",
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
															"not_equals": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Non-equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Selection criteria for the non-equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Optional:    true,
																									Description: "List of values for the selection.",
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
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the filter is currently enabled or not.",
						},
						"collapsed": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the filter's UI representation should be collapsed or expanded.",
						},
					},
				},
			},
			"annotations": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of annotations that can be applied to the dashboard's visual elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique identifier within the dashboard.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier within the dashboard.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the annotation.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the annotation is enabled.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Source of the annotation events.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Metrics source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"promql_query": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "PromQL query.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The PromQL query string.",
															},
														},
													},
												},
												"strategy": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Strategy for turning metrics data into annotations.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_time_metric": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Take first data point and use its value as annotation timestamp (instead of point own timestamp).",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
														},
													},
												},
												"message_template": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template for the annotation message.",
												},
												"labels": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Labels to display in the annotation.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"logs": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Logs source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lucene_query": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Lucene query.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The Lucene query string.",
															},
														},
													},
												},
												"strategy": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Strategy for turning logs data into annotations.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"instant": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Event timestamp is extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"range": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Event start and end timestamps are extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																		"end_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"duration": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Event start timestamp and duration are extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																		"duration_field": &schema.Schema{
																			Type:        schema.TypeList,
																			MinItems:    1,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Path within the dataset scope.",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Scope of the dataset.",
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
												"message_template": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template for the annotation message.",
												},
												"label_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Labels to display in the annotation.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keypath": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Path within the dataset scope.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"scope": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Scope of the dataset.",
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
			"absolute_time_frame": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Absolute time frame specifying a fixed start and end time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "From is the start of the time frame.",
						},
						"to": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "To is the end of the time frame.",
						},
					},
				},
			},
			"relative_time_frame": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_dashboard", "relative_time_frame"),
				Description:  "Relative time frame specifying a duration from the current time.",
			},
			"folder_id": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Unique identifier of the folder containing the dashboard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The UUID value.",
						},
					},
				},
			},
			"folder_path": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Path of the folder containing the dashboard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segments": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The segments of the folder path.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"false": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Auto refresh interval is set to off.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"two_minutes": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Auto refresh interval is set to two minutes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"five_minutes": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Auto refresh interval is set to five minutes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"dashboard_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dashboard Id.",
			},
		},
	}
}

func ResourceIbmLogsDashboardValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "href",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9]{21}$`,
			MinValueLength:             21,
			MaxValueLength:             21,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             200,
		},
		validate.ValidateSchema{
			Identifier:                 "relative_time_frame",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[0-9]+[smhdw]?$`,
			MinValueLength:             2,
			MaxValueLength:             10,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_dashboard", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsDashboardCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	bodyModelMap := map[string]interface{}{}
	createDashboardOptions := &logsv0.CreateDashboardOptions{}

	if _, ok := d.GetOk("href"); ok {
		bodyModelMap["href"] = d.Get("href")
	}

	bodyModelMap["name"] = d.Get("name")
	if _, ok := d.GetOk("description"); ok {
		bodyModelMap["description"] = d.Get("description")
	}
	bodyModelMap["layout"] = d.Get("layout")
	if _, ok := d.GetOk("variables"); ok {
		bodyModelMap["variables"] = d.Get("variables")
	}
	if _, ok := d.GetOk("filters"); ok {
		bodyModelMap["filters"] = d.Get("filters")
	}
	if _, ok := d.GetOk("annotations"); ok {
		bodyModelMap["annotations"] = d.Get("annotations")
	}
	if _, ok := d.GetOk("absolute_time_frame"); ok {
		bodyModelMap["absolute_time_frame"] = d.Get("absolute_time_frame")
	}
	if _, ok := d.GetOk("relative_time_frame"); ok {
		bodyModelMap["relative_time_frame"] = d.Get("relative_time_frame")
	}
	if _, ok := d.GetOk("folder_id"); ok {
		bodyModelMap["folder_id"] = d.Get("folder_id")
	}
	if _, ok := d.GetOk("folder_path"); ok {
		bodyModelMap["folder_path"] = d.Get("folder_path")
	}
	if _, ok := d.GetOk("false"); ok {
		bodyModelMap["false"] = d.Get("false")
	}
	if _, ok := d.GetOk("two_minutes"); ok {
		bodyModelMap["two_minutes"] = d.Get("two_minutes")
	}
	if _, ok := d.GetOk("five_minutes"); ok {
		bodyModelMap["five_minutes"] = d.Get("five_minutes")
	}
	convertedModel, err := ResourceIbmLogsDashboardMapToDashboard(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "create")
		return tfErr.GetDiag()
	}
	createDashboardOptions.Dashboard = convertedModel

	dashboardIntf, _, err := logsClient.CreateDashboardWithContext(context, createDashboardOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDashboardWithContext failed: %s", err.Error()), "ibm_logs_dashboard", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	dashboard := dashboardIntf.(*logsv0.Dashboard)

	dashboardId := fmt.Sprintf("%s/%s/%s", region, instanceId, *dashboard.ID)
	d.SetId(dashboardId)

	return resourceIbmLogsDashboardRead(context, d, meta)
}

func resourceIbmLogsDashboardRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, dashboardId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getDashboardOptions := &logsv0.GetDashboardOptions{}

	getDashboardOptions.SetDashboardID(dashboardId)

	dashboardIntf, response, err := logsClient.GetDashboardWithContext(context, getDashboardOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDashboardWithContext failed: %s", err.Error()), "ibm_logs_dashboard", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	dashboard := dashboardIntf.(*logsv0.Dashboard)

	if err = d.Set("dashboard_id", dashboardId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dashboard_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if !core.IsNil(dashboard.Href) {
		if err = d.Set("href", dashboard.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if err = d.Set("name", dashboard.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(dashboard.Description) {
		if err = d.Set("description", dashboard.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	layoutMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstLayoutToMap(dashboard.Layout)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("layout", []map[string]interface{}{layoutMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting layout: %s", err))
	}
	if !core.IsNil(dashboard.Variables) {
		variables := []map[string]interface{}{}
		for _, variablesItem := range dashboard.Variables {
			variablesItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstVariableToMap(&variablesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			variables = append(variables, variablesItemMap)
		}
		if err = d.Set("variables", variables); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting variables: %s", err))
		}
	}
	if !core.IsNil(dashboard.Filters) {
		filters := []map[string]interface{}{}
		for _, filtersItem := range dashboard.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterToMap(&filtersItem)
			if err != nil {
				return diag.FromErr(err)
			}
			filters = append(filters, filtersItemMap)
		}
		if err = d.Set("filters", filters); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting filters: %s", err))
		}
	}
	if !core.IsNil(dashboard.Annotations) {
		annotations := []map[string]interface{}{}
		for _, annotationsItem := range dashboard.Annotations {
			annotationsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationToMap(&annotationsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			annotations = append(annotations, annotationsItemMap)
		}
		if err = d.Set("annotations", annotations); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting annotations: %s", err))
		}
	}
	if !core.IsNil(dashboard.AbsoluteTimeFrame) {
		absoluteTimeFrameMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonTimeFrameToMap(dashboard.AbsoluteTimeFrame)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("absolute_time_frame", []map[string]interface{}{absoluteTimeFrameMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting absolute_time_frame: %s", err))
		}
	}
	if !core.IsNil(dashboard.RelativeTimeFrame) {
		if err = d.Set("relative_time_frame", dashboard.RelativeTimeFrame); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting relative_time_frame: %s", err))
		}
	}
	if !core.IsNil(dashboard.FolderID) {
		folderIDMap, err := ResourceIbmLogsDashboardApisDashboardsV1UUIDToMap(dashboard.FolderID)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("folder_id", []map[string]interface{}{folderIDMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting folder_id: %s", err))
		}
	}
	if !core.IsNil(dashboard.FolderPath) {
		folderPathMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFolderPathToMap(dashboard.FolderPath)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("folder_path", []map[string]interface{}{folderPathMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting folder_path: %s", err))
		}
	}
	if !core.IsNil(dashboard.False) {
		falseVarMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshOffEmptyToMap(dashboard.False)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("false", []map[string]interface{}{falseVarMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting false: %s", err))
		}
	}
	if !core.IsNil(dashboard.TwoMinutes) {
		twoMinutesMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmptyToMap(dashboard.TwoMinutes)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("two_minutes", []map[string]interface{}{twoMinutesMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting two_minutes: %s", err))
		}
	}
	if !core.IsNil(dashboard.FiveMinutes) {
		fiveMinutesMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmptyToMap(dashboard.FiveMinutes)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("five_minutes", []map[string]interface{}{fiveMinutesMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting five_minutes: %s", err))
		}
	}

	return nil
}

func resourceIbmLogsDashboardUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, dashboardId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	replaceDashboardOptions := &logsv0.ReplaceDashboardOptions{}

	replaceDashboardOptions.SetDashboardID(dashboardId)

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("layout") ||
		d.HasChange("variables") ||
		d.HasChange("filters") ||
		d.HasChange("annotations") ||
		d.HasChange("absolute_time_frame") ||
		d.HasChange("relative_time_frame") ||
		d.HasChange("folder_id") ||
		d.HasChange("folder_path") ||
		d.HasChange("false") ||
		d.HasChange("two_minutes") ||
		d.HasChange("five_minutes") {

		bodyModelMap := map[string]interface{}{}

		if _, ok := d.GetOk("href"); ok {
			bodyModelMap["href"] = d.Get("href")
		}

		bodyModelMap["name"] = d.Get("name")
		if _, ok := d.GetOk("description"); ok {
			bodyModelMap["description"] = d.Get("description")
		}
		bodyModelMap["layout"] = d.Get("layout")
		if _, ok := d.GetOk("variables"); ok {
			bodyModelMap["variables"] = d.Get("variables")
		}
		if _, ok := d.GetOk("filters"); ok {
			bodyModelMap["filters"] = d.Get("filters")
		}
		if _, ok := d.GetOk("annotations"); ok {
			bodyModelMap["annotations"] = d.Get("annotations")
		}
		if _, ok := d.GetOk("absolute_time_frame"); ok {
			bodyModelMap["absolute_time_frame"] = d.Get("absolute_time_frame")
		}
		if _, ok := d.GetOk("relative_time_frame"); ok {
			bodyModelMap["relative_time_frame"] = d.Get("relative_time_frame")
		}
		if _, ok := d.GetOk("folder_id"); ok {
			bodyModelMap["folder_id"] = d.Get("folder_id")
		}
		if _, ok := d.GetOk("folder_path"); ok {
			bodyModelMap["folder_path"] = d.Get("folder_path")
		}
		if _, ok := d.GetOk("false"); ok {
			bodyModelMap["false"] = d.Get("false")
		}
		if _, ok := d.GetOk("two_minutes"); ok {
			bodyModelMap["two_minutes"] = d.Get("two_minutes")
		}
		if _, ok := d.GetOk("five_minutes"); ok {
			bodyModelMap["five_minutes"] = d.Get("five_minutes")
		}
		convertedModel, err := ResourceIbmLogsDashboardMapToDashboard(bodyModelMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "create")
			return tfErr.GetDiag()
		}
		replaceDashboardOptions.Dashboard = convertedModel

		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.ReplaceDashboardWithContext(context, replaceDashboardOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceDashboardWithContext failed: %s", err.Error()), "ibm_logs_dashboard", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsDashboardRead(context, d, meta)
}

func resourceIbmLogsDashboardDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, dashboardId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDashboardOptions := &logsv0.DeleteDashboardOptions{}

	deleteDashboardOptions.SetDashboardID(dashboardId)

	_, err = logsClient.DeleteDashboardWithContext(context, deleteDashboardOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDashboardWithContext failed: %s", err.Error()), "ibm_logs_dashboard", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstLayout, error) {
	model := &logsv0.ApisDashboardsV1AstLayout{}
	if modelMap["sections"] != nil {
		sections := []logsv0.ApisDashboardsV1AstSection{}
		for _, sectionsItem := range modelMap["sections"].([]interface{}) {
			sectionsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstSection(sectionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			sections = append(sections, *sectionsItemModel)
		}
		model.Sections = sections
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstSection(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstSection, error) {
	model := &logsv0.ApisDashboardsV1AstSection{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	IDModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap["id"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ID = IDModel
	if modelMap["rows"] != nil {
		rows := []logsv0.ApisDashboardsV1AstRow{}
		for _, rowsItem := range modelMap["rows"].([]interface{}) {
			rowsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstRow(rowsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rows = append(rows, *rowsItemModel)
		}
		model.Rows = rows
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1UUID, error) {
	model := &logsv0.ApisDashboardsV1UUID{}
	model.Value = core.UUIDPtr(strfmt.UUID(modelMap["value"].(string)))
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstRow(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstRow, error) {
	model := &logsv0.ApisDashboardsV1AstRow{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	IDModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap["id"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ID = IDModel
	AppearanceModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstRowAppearance(modelMap["appearance"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Appearance = AppearanceModel
	if modelMap["widgets"] != nil {
		widgets := []logsv0.ApisDashboardsV1AstWidget{}
		for _, widgetsItem := range modelMap["widgets"].([]interface{}) {
			widgetsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidget(widgetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			widgets = append(widgets, *widgetsItemModel)
		}
		model.Widgets = widgets
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstRowAppearance(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstRowAppearance, error) {
	model := &logsv0.ApisDashboardsV1AstRowAppearance{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.Height = core.Int64Ptr(int64(modelMapElement["height"].(int)))
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidget(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidget, error) {
	model := &logsv0.ApisDashboardsV1AstWidget{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	IDModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap["id"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ID = IDModel
	model.Title = core.StringPtr(modelMap["title"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	DefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinition(modelMap["definition"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	if modelMap["created_at"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["created_at"].(string))
		if err != nil {
			return model, err
		}
		model.CreatedAt = &dateTime
	}
	if modelMap["updated_at"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["updated_at"].(string))
		if err != nil {
			return model, err
		}
		model.UpdatedAt = &dateTime
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinition(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetDefinitionIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["line_chart"] != nil && len(modelMapElement["line_chart"].([]interface{})) > 0 {
			LineChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChart(modelMapElement["line_chart"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LineChart = LineChartModel
		}
		if modelMapElement["data_table"] != nil && len(modelMapElement["data_table"].([]interface{})) > 0 {
			DataTableModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTable(modelMapElement["data_table"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.DataTable = DataTableModel
		}
		if modelMapElement["gauge"] != nil && len(modelMapElement["gauge"].([]interface{})) > 0 {
			GaugeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGauge(modelMapElement["gauge"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Gauge = GaugeModel
		}
		if modelMapElement["pie_chart"] != nil && len(modelMapElement["pie_chart"].([]interface{})) > 0 {
			PieChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChart(modelMapElement["pie_chart"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.PieChart = PieChartModel
		}
		if modelMapElement["bar_chart"] != nil && len(modelMapElement["bar_chart"].([]interface{})) > 0 {
			BarChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChart(modelMapElement["bar_chart"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.BarChart = BarChartModel
		}
		if modelMapElement["horizontal_bar_chart"] != nil && len(modelMapElement["horizontal_bar_chart"].([]interface{})) > 0 {
			HorizontalBarChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChart(modelMapElement["horizontal_bar_chart"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.HorizontalBarChart = HorizontalBarChartModel
		}
		if modelMapElement["markdown"] != nil && len(modelMapElement["markdown"].([]interface{})) > 0 {
			MarkdownModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsMarkdown(modelMapElement["markdown"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Markdown = MarkdownModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChart(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChart{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		LegendModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLegend(modelMapElement["legend"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Legend = LegendModel
		TooltipModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartTooltip(modelMapElement["tooltip"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Tooltip = TooltipModel
		if modelMapElement["query_definitions"] != nil {
			queryDefinitions := []logsv0.ApisDashboardsV1AstWidgetsLineChartQueryDefinition{}
			for _, queryDefinitionsItem := range modelMapElement["query_definitions"].([]interface{}) {
				queryDefinitionsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQueryDefinition(queryDefinitionsItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				queryDefinitions = append(queryDefinitions, *queryDefinitionsItemModel)
			}
			model.QueryDefinitions = queryDefinitions
		}
		if modelMapElement["stacked_line"] != nil && modelMapElement["stacked_line"].(string) != "" {
			model.StackedLine = core.StringPtr(modelMapElement["stacked_line"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLegend(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonLegend, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonLegend{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.IsVisible = core.BoolPtr(modelMapElement["is_visible"].(bool))
		if modelMapElement["columns"] != nil {
			columns := []string{}
			for _, columnsItem := range modelMapElement["columns"].([]interface{}) {
				columns = append(columns, columnsItem.(string))
			}
			model.Columns = columns
		}
		model.GroupByQuery = core.BoolPtr(modelMapElement["group_by_query"].(bool))
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartTooltip(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartTooltip, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartTooltip{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["show_labels"] != nil {
			model.ShowLabels = core.BoolPtr(modelMapElement["show_labels"].(bool))
		}
		if modelMapElement["type"] != nil && modelMapElement["type"].(string) != "" {
			model.Type = core.StringPtr(modelMapElement["type"].(string))
		}

	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQueryDefinition(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryDefinition, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartQueryDefinition{}
	model.ID = core.UUIDPtr(strfmt.UUID(modelMap["id"].(string)))
	QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQuery(modelMap["query"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Query = QueryModel
	if modelMap["series_name_template"] != nil && modelMap["series_name_template"].(string) != "" {
		model.SeriesNameTemplate = core.StringPtr(modelMap["series_name_template"].(string))
	}
	if modelMap["series_count_limit"] != nil && modelMap["series_count_limit"].(string) != "" {
		model.SeriesCountLimit = core.StringPtr(modelMap["series_count_limit"].(string))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	if modelMap["scale_type"] != nil && modelMap["scale_type"].(string) != "" {
		model.ScaleType = core.StringPtr(modelMap["scale_type"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.IsVisible = core.BoolPtr(modelMap["is_visible"].(bool))
	if modelMap["color_scheme"] != nil && modelMap["color_scheme"].(string) != "" {
		model.ColorScheme = core.StringPtr(modelMap["color_scheme"].(string))
	}
	ResolutionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartResolution(modelMap["resolution"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Resolution = ResolutionModel
	if modelMap["data_mode_type"] != nil && modelMap["data_mode_type"].(string) != "" {
		model.DataModeType = core.StringPtr(modelMap["data_mode_type"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsLineChartQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartLogsQuery(modelMapElement["logs"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartMetricsQuery(modelMapElement["metrics"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartDataprimeQuery(modelMapElement["dataprime"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartLogsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartLogsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["lucene_query"] != nil && len(modelMapElement["lucene_query"].([]interface{})) > 0 {
			LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMapElement["lucene_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LuceneQuery = LuceneQueryModel
		}
		if modelMapElement["group_by"] != nil {
			groupBy := []string{}
			for _, groupByItem := range modelMapElement["group_by"].([]interface{}) {
				groupBy = append(groupBy, groupByItem.(string))
			}
			model.GroupBy = groupBy
		}
		if modelMapElement["aggregations"] != nil {
			aggregations := []logsv0.ApisDashboardsV1CommonLogsAggregationIntf{}
			for _, aggregationsItem := range modelMapElement["aggregations"].([]interface{}) {
				aggregationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(aggregationsItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				aggregations = append(aggregations, aggregationsItemModel)
			}
			model.Aggregations = aggregations
		}
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
		if modelMapElement["group_bys"] != nil {
			// groupBys := []logsv0.ApisDashboardsV1CommonObservationField{}
			// for _, groupBysItem := range modelMap["group_bys"].([]interface{}) {
			// 	groupBysItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(groupBysItem.(map[string]interface{}))
			// 	if err != nil {
			// 		return model, err
			// 	}
			// 	groupBys = append(groupBys, *groupBysItemModel)
			// }
			groupBys, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapElement["group_bys"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.GroupBys = groupBys
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonLuceneQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonLuceneQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["value"] != nil && modelMapElement["value"].(string) != "" {
			model.Value = core.StringPtr(modelMapElement["value"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMap map[string]interface{}) (logsv0.ApisDashboardsV1CommonLogsAggregationIntf, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregation{}
	if modelMap["count"] != nil && len(modelMap["count"].([]interface{})) > 0 {
		CountModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountEmpty(modelMap["count"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Count = CountModel
	}
	if modelMap["count_distinct"] != nil && len(modelMap["count_distinct"].([]interface{})) > 0 {
		CountDistinctModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountDistinct(modelMap["count_distinct"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.CountDistinct = CountDistinctModel
	}
	if modelMap["sum"] != nil && len(modelMap["sum"].([]interface{})) > 0 {
		SumModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationSum(modelMap["sum"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Sum = SumModel
	}
	if modelMap["average"] != nil && len(modelMap["average"].([]interface{})) > 0 {
		AverageModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationAverage(modelMap["average"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Average = AverageModel
	}
	if modelMap["min"] != nil && len(modelMap["min"].([]interface{})) > 0 {
		MinModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMin(modelMap["min"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Min = MinModel
	}
	if modelMap["max"] != nil && len(modelMap["max"].([]interface{})) > 0 {
		MaxModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMax(modelMap["max"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Max = MaxModel
	}
	if modelMap["percentile"] != nil && len(modelMap["percentile"].([]interface{})) > 0 {
		PercentileModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationPercentile(modelMap["percentile"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Percentile = PercentileModel
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationCountEmpty, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationCountEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountDistinct(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationCountDistinct, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationCountDistinct{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap[0].(map[string]interface{})["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap []interface{}) ([]logsv0.ApisDashboardsV1CommonObservationField, error) {
	models := []logsv0.ApisDashboardsV1CommonObservationField{}

	if modelMap != nil && len(modelMap) > 0 {
		for _, element := range modelMap {
			if element != nil {
				model := &logsv0.ApisDashboardsV1CommonObservationField{}
				modelMapElement := element.(map[string]interface{})
				if modelMapElement["keypath"] != nil {
					keypath := []string{}
					for _, keypathItem := range modelMapElement["keypath"].([]interface{}) {
						keypath = append(keypath, keypathItem.(string))
					}
					model.Keypath = keypath
				}
				if modelMapElement["scope"] != nil && modelMapElement["scope"].(string) != "" {
					model.Scope = core.StringPtr(modelMapElement["scope"].(string))
				}
				models = append(models, *model)
			}

		}

	}

	return models, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationSum(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationSum, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationSum{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationAverage(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationAverage, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationAverage{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMin(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationMin, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationMin{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMax(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationMax, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationMax{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationPercentile(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationPercentile, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationPercentile{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		model.Percent = core.Float64Ptr(modelMapObsevationField["percent"].(float64))
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueCount(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueCount, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueCount{}
	if modelMap["count"] != nil && len(modelMap["count"].([]interface{})) > 0 {
		CountModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountEmpty(modelMap["count"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Count = CountModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueCountDistinct(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct{}
	if modelMap["count_distinct"] != nil && len(modelMap["count_distinct"].([]interface{})) > 0 {
		CountDistinctModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationCountDistinct(modelMap["count_distinct"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.CountDistinct = CountDistinctModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueSum(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueSum, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueSum{}
	if modelMap["sum"] != nil && len(modelMap["sum"].([]interface{})) > 0 {
		SumModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationSum(modelMap["sum"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Sum = SumModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueAverage(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage{}
	if modelMap["average"] != nil && len(modelMap["average"].([]interface{})) > 0 {
		AverageModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationAverage(modelMap["average"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Average = AverageModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueMin(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueMin, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueMin{}
	if modelMap["min"] != nil && len(modelMap["min"].([]interface{})) > 0 {
		MinModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMin(modelMap["min"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Min = MinModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValueMax(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValueMax, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValueMax{}
	if modelMap["max"] != nil && len(modelMap["max"].([]interface{})) > 0 {
		MaxModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationMax(modelMap["max"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Max = MaxModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationValuePercentile(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile, error) {
	model := &logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile{}
	if modelMap["percentile"] != nil && len(modelMap["percentile"].([]interface{})) > 0 {
		PercentileModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregationPercentile(modelMap["percentile"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Percentile = PercentileModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterLogsFilter, error) {
	model := &logsv0.ApisDashboardsV1AstFilterLogsFilter{}
	if modelMap["operator"] != nil && len(modelMap["operator"].([]interface{})) > 0 {
		OperatorModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterOperator(modelMap["operator"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Operator = OperatorModel
	}
	if modelMap["observation_field"] != nil && len(modelMap["observation_field"].([]interface{})) > 0 && modelMap["observation_field"].([]interface{})[0] != nil {
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterOperator(modelMap []interface{}) (logsv0.ApisDashboardsV1AstFilterOperatorIntf, error) {
	model := &logsv0.ApisDashboardsV1AstFilterOperator{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["equals"] != nil && len(modelMapElement["equals"].([]interface{})) > 0 {
			EqualsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEquals(modelMapElement["equals"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Equals = EqualsModel
		}
		if modelMapElement["not_equals"] != nil && len(modelMapElement["not_equals"].([]interface{})) > 0 {
			NotEqualsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEquals(modelMapElement["not_equals"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.NotEquals = NotEqualsModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEquals(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterEquals, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEquals{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["selection"] != nil && len(modelMapElement["selection"].([]interface{})) > 0 {
			SelectionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelection(modelMapElement["selection"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Selection = SelectionModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelection(modelMap []interface{}) (logsv0.ApisDashboardsV1AstFilterEqualsSelectionIntf, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEqualsSelection{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["all"] != nil && len(modelMapElement["all"].([]interface{})) > 0 {
			AllModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty(modelMapElement["all"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.All = AllModel
		}
		if modelMapElement["list"] != nil && len(modelMapElement["list"].([]interface{})) > 0 {
			ListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionListSelection(modelMapElement["list"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.List = ListModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionListSelection(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterEqualsSelectionListSelection, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEqualsSelectionListSelection{}
	if modelMap[0] != nil {
		if modelMap[0].(map[string]interface{})["values"] != nil {
			values := []string{}
			for _, valuesItem := range modelMap[0].(map[string]interface{})["values"].([]interface{}) {
				values = append(values, valuesItem.(string))
			}
			model.Values = values
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionValueAll(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll{}
	if modelMap["all"] != nil && len(modelMap["all"].([]interface{})) > 0 {
		AllModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty(modelMap["all"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.All = AllModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionValueList(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList, error) {
	model := &logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList{}
	if modelMap["list"] != nil && len(modelMap["list"].([]interface{})) > 0 {
		ListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEqualsSelectionListSelection(modelMap["list"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.List = ListModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEquals(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterNotEquals, error) {
	model := &logsv0.ApisDashboardsV1AstFilterNotEquals{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["selection"] != nil && len(modelMapElement["selection"].([]interface{})) > 0 {
			SelectionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEqualsSelection(modelMapElement["selection"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Selection = SelectionModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEqualsSelection(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterNotEqualsSelection, error) {
	model := &logsv0.ApisDashboardsV1AstFilterNotEqualsSelection{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["list"] != nil && len(modelMapElement["list"].([]interface{})) > 0 {
			ListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEqualsSelectionListSelection(modelMapElement["list"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.List = ListModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEqualsSelectionListSelection(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFilterNotEqualsSelectionListSelection, error) {
	model := &logsv0.ApisDashboardsV1AstFilterNotEqualsSelectionListSelection{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["values"] != nil {
			values := []string{}
			for _, valuesItem := range modelMapElement["values"].([]interface{}) {
				values = append(values, valuesItem.(string))
			}
			model.Values = values
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterOperatorValueEquals(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterOperatorValueEquals, error) {
	model := &logsv0.ApisDashboardsV1AstFilterOperatorValueEquals{}
	if modelMap["equals"] != nil && len(modelMap["equals"].([]interface{})) > 0 {
		EqualsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterEquals(modelMap["equals"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Equals = EqualsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterOperatorValueNotEquals(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals, error) {
	model := &logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals{}
	if modelMap["not_equals"] != nil && len(modelMap["not_equals"].([]interface{})) > 0 {
		NotEqualsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterNotEquals(modelMap["not_equals"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.NotEquals = NotEqualsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartMetricsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartMetricsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["promql_query"] != nil && len(modelMapElement["promql_query"].([]interface{})) > 0 {
			PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMapElement["promql_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.PromqlQuery = PromqlQueryModel
		}
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonPromQlQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonPromQlQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["value"] != nil && modelMapElement["value"].(string) != "" {
			model.Value = core.StringPtr(modelMapElement["value"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterMetricsFilter, error) {
	model := &logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["operator"] != nil && len(modelMap["operator"].([]interface{})) > 0 {
		OperatorModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterOperator(modelMap["operator"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Operator = OperatorModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartDataprimeQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartDataprimeQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMapElement["dataprime_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DataprimeQuery = DataprimeQueryModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, filtersItemModel)
			}
			model.Filters = filters
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1CommonDataprimeQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["text"] != nil && modelMapElement["text"].(string) != "" {
			model.Text = core.StringPtr(modelMapElement["text"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(modelMap map[string]interface{}) (logsv0.ApisDashboardsV1AstFilterSourceIntf, error) {
	model := &logsv0.ApisDashboardsV1AstFilterSource{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(modelMap["logs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(modelMap["metrics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSourceValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterSourceValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstFilterSourceValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(modelMap["logs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSourceValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilterSourceValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstFilterSourceValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(modelMap["metrics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartLogsQuery(modelMap["logs"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartMetricsQuery(modelMap["metrics"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartDataprimeQuery(modelMap["dataprime"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChartResolution(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsLineChartResolution, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsLineChartResolution{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["interval"] != nil && modelMapElement["interval"].(string) != "" {
			model.Interval = core.StringPtr(modelMapElement["interval"].(string))
		}
		if modelMapElement["buckets_presented"] != nil {
			model.BucketsPresented = core.Int64Ptr(int64(modelMapElement["buckets_presented"].(int)))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTable(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTable, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTable{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableQuery(modelMapElement["query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Query = QueryModel
		model.ResultsPerPage = core.Int64Ptr(int64(modelMapElement["results_per_page"].(int)))
		model.RowStyle = core.StringPtr(modelMapElement["row_style"].(string))
		if modelMapElement["columns"] != nil {
			columns := []logsv0.ApisDashboardsV1AstWidgetsDataTableColumn{}
			for _, columnsItem := range modelMapElement["columns"].([]interface{}) {
				columnsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableColumn(columnsItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				columns = append(columns, *columnsItemModel)
			}
			model.Columns = columns
		}
		if modelMapElement["order_by"] != nil && len(modelMapElement["order_by"].([]interface{})) > 0 {
			OrderByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonOrderingField(modelMapElement["order_by"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.OrderBy = OrderByModel
		}
		if modelMapElement["data_mode_type"] != nil && modelMapElement["data_mode_type"].(string) != "" {
			model.DataModeType = core.StringPtr(modelMapElement["data_mode_type"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsDataTableQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQuery(modelMapElement["logs"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableMetricsQuery(modelMapElement["metrics"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableDataprimeQuery(modelMapElement["dataprime"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["lucene_query"] != nil && len(modelMapElement["lucene_query"].([]interface{})) > 0 {
			LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMapElement["lucene_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LuceneQuery = LuceneQueryModel
		}
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
		if modelMapElement["grouping"] != nil && len(modelMapElement["grouping"].([]interface{})) > 0 {
			GroupingModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping(modelMapElement["grouping"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Grouping = GroupingModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping{}
	if modelMap["aggregations"] != nil {
		aggregations := []logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation{}
		for _, aggregationsItem := range modelMap["aggregations"].([]interface{}) {
			aggregationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation(aggregationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			aggregations = append(aggregations, *aggregationsItemModel)
		}
		model.Aggregations = aggregations
	}
	if modelMap["group_bys"] != nil {
		// groupBys := []logsv0.ApisDashboardsV1CommonObservationField{}
		// for _, groupBysItem := range modelMap["group_bys"].([]interface{}) {
		// 	groupBysItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(groupBysItem.(map[string]interface{}))
		// 	if err != nil {
		// 		return model, err
		// 	}
		// 	groupBys = append(groupBys, *groupBysItemModel)
		// }
		groupBys, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["group_bys"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GroupBys = groupBys
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation{}
	model.ID = core.UUIDPtr(strfmt.UUID(modelMap["id"].(string)))
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.IsVisible = core.BoolPtr(modelMap["is_visible"].(bool))
	AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMap["aggregation"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Aggregation = AggregationModel
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableMetricsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableMetricsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMapElement["promql_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PromqlQuery = PromqlQueryModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableDataprimeQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableDataprimeQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMapElement["dataprime_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DataprimeQuery = DataprimeQueryModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, filtersItemModel)
			}
			model.Filters = filters
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableLogsQuery(modelMap["logs"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableMetricsQuery(modelMap["metrics"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableDataprimeQuery(modelMap["dataprime"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTableColumn(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsDataTableColumn, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsDataTableColumn{}
	model.Field = core.StringPtr(modelMap["field"].(string))
	if modelMap["width"] != nil {
		model.Width = core.Int64Ptr(int64(modelMap["width"].(int)))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonOrderingField(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonOrderingField, error) {
	model := &logsv0.ApisDashboardsV1CommonOrderingField{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["field"] != nil && modelMapElement["field"].(string) != "" {
			model.Field = core.StringPtr(modelMapElement["field"].(string))
		}
		if modelMapElement["order_direction"] != nil && modelMapElement["order_direction"].(string) != "" {
			model.OrderDirection = core.StringPtr(modelMapElement["order_direction"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGauge(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGauge, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGauge{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeQuery(modelMapElement["query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Query = QueryModel
		model.Min = core.Float64Ptr(modelMapElement["min"].(float64))
		model.Max = core.Float64Ptr(modelMapElement["max"].(float64))
		model.ShowInnerArc = core.BoolPtr(modelMapElement["show_inner_arc"].(bool))
		model.ShowOuterArc = core.BoolPtr(modelMapElement["show_outer_arc"].(bool))
		model.Unit = core.StringPtr(modelMapElement["unit"].(string))
		if modelMapElement["thresholds"] != nil {
			thresholds := []logsv0.ApisDashboardsV1AstWidgetsGaugeThreshold{}
			for _, thresholdsItem := range modelMapElement["thresholds"].([]interface{}) {
				thresholdsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeThreshold(thresholdsItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				thresholds = append(thresholds, *thresholdsItemModel)
			}
			model.Thresholds = thresholds
		}
		if modelMapElement["data_mode_type"] != nil && modelMapElement["data_mode_type"].(string) != "" {
			model.DataModeType = core.StringPtr(modelMapElement["data_mode_type"].(string))
		}
		model.ThresholdBy = core.StringPtr(modelMapElement["threshold_by"].(string))
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsGaugeQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeMetricsQuery(modelMapElement["metrics"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeLogsQuery(modelMapElement["logs"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeDataprimeQuery(modelMapElement["dataprime"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeMetricsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeMetricsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMapElement["promql_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PromqlQuery = PromqlQueryModel
		model.Aggregation = core.StringPtr(modelMapElement["aggregation"].(string))
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeLogsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeLogsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["lucene_query"] != nil && len(modelMapElement["lucene_query"].([]interface{})) > 0 {
			LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMapElement["lucene_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LuceneQuery = LuceneQueryModel
		}
		if modelMapElement["logs_aggregation"] != nil && len(modelMapElement["logs_aggregation"].([]interface{})) > 0 {
			LogsAggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMapElement["logs_aggregation"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.LogsAggregation = LogsAggregationModel
		}
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeDataprimeQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeDataprimeQuery{}
	DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMap["dataprime_query"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.DataprimeQuery = DataprimeQueryModel
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, filtersItemModel)
		}
		model.Filters = filters
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeMetricsQuery(modelMap["metrics"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeLogsQuery(modelMap["logs"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeDataprimeQuery(modelMap["dataprime"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGaugeThreshold(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsGaugeThreshold, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsGaugeThreshold{}
	model.From = core.Float64Ptr(modelMap["from"].(float64))
	model.Color = core.StringPtr(modelMap["color"].(string))
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChart(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChart{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartQuery(modelMapElement["query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Query = QueryModel
		model.MaxSlicesPerChart = core.Int64Ptr(int64(modelMapElement["max_slices_per_chart"].(int)))
		model.MinSlicePercentage = core.Int64Ptr(int64(modelMapElement["min_slice_percentage"].(int)))
		StackDefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartStackDefinition(modelMapElement["stack_definition"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StackDefinition = StackDefinitionModel
		LabelDefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartLabelDefinition(modelMapElement["label_definition"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LabelDefinition = LabelDefinitionModel
		model.ShowLegend = core.BoolPtr(modelMapElement["show_legend"].(bool))
		if modelMapElement["group_name_template"] != nil && modelMapElement["group_name_template"].(string) != "" {
			model.GroupNameTemplate = core.StringPtr(modelMapElement["group_name_template"].(string))
		}
		if modelMapElement["unit"] != nil && modelMapElement["unit"].(string) != "" {
			model.Unit = core.StringPtr(modelMapElement["unit"].(string))
		}
		model.ColorScheme = core.StringPtr(modelMapElement["color_scheme"].(string))
		if modelMapElement["data_mode_type"] != nil && modelMapElement["data_mode_type"].(string) != "" {
			model.DataModeType = core.StringPtr(modelMapElement["data_mode_type"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsPieChartQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartLogsQuery(modelMapElement["logs"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartMetricsQuery(modelMapElement["metrics"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartDataprimeQuery(modelMapElement["dataprime"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartLogsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartLogsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["lucene_query"] != nil && len(modelMapElement["lucene_query"].([]interface{})) > 0 {
			LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMapElement["lucene_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LuceneQuery = LuceneQueryModel
		}
		AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMapElement["aggregation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Aggregation = AggregationModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
		if modelMapElement["group_names_fields"] != nil {
			// groupNamesFields := []logsv0.ApisDashboardsV1CommonObservationField{}
			// for _, groupNamesFieldsItem := range modelMap["group_names_fields"].([]interface{}) {
			// 	groupNamesFieldsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(groupNamesFieldsItem.(map[string]interface{}))
			// 	if err != nil {
			// 		return model, err
			// 	}
			// 	groupNamesFields = append(groupNamesFields, *groupNamesFieldsItemModel)
			// }
			groupNamesFields, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapElement["group_names_fields"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.GroupNamesFields = groupNamesFields
		}
		if modelMapElement["stacked_group_name_field"] != nil && len(modelMapElement["stacked_group_name_field"].([]interface{})) > 0 && modelMapElement["stacked_group_name_field"].([]interface{})[0] != nil {
			StackedGroupNameFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapElement["stacked_group_name_field"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.StackedGroupNameField = &StackedGroupNameFieldModel[0]
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartMetricsQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartMetricsQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMapElement["promql_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PromqlQuery = PromqlQueryModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, *filtersItemModel)
			}
			model.Filters = filters
		}
		if modelMapElement["group_names"] != nil {
			groupNames := []string{}
			for _, groupNamesItem := range modelMapElement["group_names"].([]interface{}) {
				groupNames = append(groupNames, groupNamesItem.(string))
			}
			model.GroupNames = groupNames
		}
		if modelMapElement["stacked_group_name"] != nil && modelMapElement["stacked_group_name"].(string) != "" {
			model.StackedGroupName = core.StringPtr(modelMapElement["stacked_group_name"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartDataprimeQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartDataprimeQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMapElement["dataprime_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DataprimeQuery = DataprimeQueryModel
		if modelMapElement["filters"] != nil {
			filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
			for _, filtersItem := range modelMapElement["filters"].([]interface{}) {
				filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				filters = append(filters, filtersItemModel)
			}
			model.Filters = filters
		}
		if modelMapElement["group_names"] != nil {
			groupNames := []string{}
			for _, groupNamesItem := range modelMapElement["group_names"].([]interface{}) {
				groupNames = append(groupNames, groupNamesItem.(string))
			}
			model.GroupNames = groupNames
		}
		if modelMapElement["stacked_group_name"] != nil && modelMapElement["stacked_group_name"].(string) != "" {
			model.StackedGroupName = core.StringPtr(modelMapElement["stacked_group_name"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartLogsQuery(modelMap["logs"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartMetricsQuery(modelMap["metrics"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartDataprimeQuery(modelMap["dataprime"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartStackDefinition(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartStackDefinition, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartStackDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["max_slices_per_stack"] != nil {
			model.MaxSlicesPerStack = core.Int64Ptr(int64(modelMapElement["max_slices_per_stack"].(int)))
		}
		if modelMapElement["stack_name_template"] != nil && modelMapElement["stack_name_template"].(string) != "" {
			model.StackNameTemplate = core.StringPtr(modelMapElement["stack_name_template"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChartLabelDefinition(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsPieChartLabelDefinition, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsPieChartLabelDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["label_source"] != nil && modelMapElement["label_source"].(string) != "" {
			model.LabelSource = core.StringPtr(modelMapElement["label_source"].(string))
		}
		if modelMapElement["is_visible"] != nil {
			model.IsVisible = core.BoolPtr(modelMapElement["is_visible"].(bool))
		}
		if modelMapElement["show_name"] != nil {
			model.ShowName = core.BoolPtr(modelMapElement["show_name"].(bool))
		}
		if modelMapElement["show_value"] != nil {
			model.ShowValue = core.BoolPtr(modelMapElement["show_value"].(bool))
		}
		if modelMapElement["show_percentage"] != nil {
			model.ShowPercentage = core.BoolPtr(modelMapElement["show_percentage"].(bool))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChart(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChart{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartQuery(modelMapElement["query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Query = QueryModel
		model.MaxBarsPerChart = core.Int64Ptr(int64(modelMapElement["max_bars_per_chart"].(int)))
		model.GroupNameTemplate = core.StringPtr(modelMapElement["group_name_template"].(string))
		StackDefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartStackDefinition(modelMapElement["stack_definition"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StackDefinition = StackDefinitionModel
		model.ScaleType = core.StringPtr(modelMapElement["scale_type"].(string))
		ColorsByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsBy(modelMapElement["colors_by"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ColorsBy = ColorsByModel
		XAxisModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxis(modelMapElement["x_axis"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.XAxis = XAxisModel
		model.Unit = core.StringPtr(modelMapElement["unit"].(string))
		model.SortBy = core.StringPtr(modelMapElement["sort_by"].(string))
		model.ColorScheme = core.StringPtr(modelMapElement["color_scheme"].(string))
		if modelMapElement["data_mode_type"] != nil && modelMapElement["data_mode_type"].(string) != "" {
			model.DataModeType = core.StringPtr(modelMapElement["data_mode_type"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsBarChartQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartLogsQuery(modelMapElement["logs"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartMetricsQuery(modelMapElement["metrics"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartDataprimeQuery(modelMapElement["dataprime"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartLogsQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartLogsQuery{}
	if modelMap["lucene_query"] != nil && len(modelMap["lucene_query"].([]interface{})) > 0 {
		LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMap["lucene_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LuceneQuery = LuceneQueryModel
	}
	AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMap["aggregation"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Aggregation = AggregationModel
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names_fields"] != nil {
		// groupNamesFields := []logsv0.ApisDashboardsV1CommonObservationField{}
		// for _, groupNamesFieldsItem := range modelMap["group_names_fields"].([]interface{}) {
		// 	groupNamesFieldsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(groupNamesFieldsItem.(map[string]interface{}))
		// 	if err != nil {
		// 		return model, err
		// 	}
		// 	groupNamesFields = append(groupNamesFields, *groupNamesFieldsItemModel)
		// }
		groupNamesFields, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["group_names_fields"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GroupNamesFields = groupNamesFields
	}
	if modelMap["stacked_group_name_field"] != nil && len(modelMap["stacked_group_name_field"].([]interface{})) > 0 && modelMap["stacked_group_name_field"].([]interface{})[0] != nil {
		StackedGroupNameFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["stacked_group_name_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StackedGroupNameField = &StackedGroupNameFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartMetricsQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartMetricsQuery{}
	if modelMap["promql_query"] != nil && len(modelMap["promql_query"].([]interface{})) > 0 {
		PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMap["promql_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PromqlQuery = PromqlQueryModel
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names"] != nil {
		groupNames := []string{}
		for _, groupNamesItem := range modelMap["group_names"].([]interface{}) {
			groupNames = append(groupNames, groupNamesItem.(string))
		}
		model.GroupNames = groupNames
	}
	if modelMap["stacked_group_name"] != nil && modelMap["stacked_group_name"].(string) != "" {
		model.StackedGroupName = core.StringPtr(modelMap["stacked_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartDataprimeQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartDataprimeQuery{}
	DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMap["dataprime_query"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.DataprimeQuery = DataprimeQueryModel
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names"] != nil {
		groupNames := []string{}
		for _, groupNamesItem := range modelMap["group_names"].([]interface{}) {
			groupNames = append(groupNames, groupNamesItem.(string))
		}
		model.GroupNames = groupNames
	}
	if modelMap["stacked_group_name"] != nil && modelMap["stacked_group_name"].(string) != "" {
		model.StackedGroupName = core.StringPtr(modelMap["stacked_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartLogsQuery(modelMap["logs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartMetricsQuery(modelMap["metrics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartDataprimeQuery(modelMap["dataprime"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartStackDefinition(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartStackDefinition, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartStackDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["max_slices_per_bar"] != nil {
			model.MaxSlicesPerBar = core.Int64Ptr(int64(modelMapElement["max_slices_per_bar"].(int)))
		}
		if modelMapElement["stack_name_template"] != nil && modelMapElement["stack_name_template"].(string) != "" {
			model.StackNameTemplate = core.StringPtr(modelMapElement["stack_name_template"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsBy(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsCommonColorsByIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsBy{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["stack"] != nil && len(modelMapElement["stack"].([]interface{})) > 0 {
			StackModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty(modelMapElement["stack"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Stack = StackModel
		}
		if modelMapElement["group_by"] != nil && len(modelMapElement["group_by"].([]interface{})) > 0 {
			GroupByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty(modelMapElement["group_by"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.GroupBy = GroupByModel
		}
		if modelMapElement["aggregation"] != nil && len(modelMapElement["aggregation"].([]interface{})) > 0 {
			AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty(modelMapElement["aggregation"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Aggregation = AggregationModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByValueStack(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack{}
	if modelMap["stack"] != nil && len(modelMap["stack"].([]interface{})) > 0 {
		StackModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty(modelMap["stack"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Stack = StackModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy{}
	if modelMap["group_by"] != nil && len(modelMap["group_by"].([]interface{})) > 0 {
		GroupByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty(modelMap["group_by"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GroupBy = GroupByModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByValueAggregation(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation{}
	if modelMap["aggregation"] != nil && len(modelMap["aggregation"].([]interface{})) > 0 {
		AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty(modelMap["aggregation"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Aggregation = AggregationModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxis(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartXAxis{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["value"] != nil && len(modelMapElement["value"].([]interface{})) > 0 {
			ValueModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty(modelMapElement["value"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Value = ValueModel
		}
		if modelMapElement["time"] != nil && len(modelMapElement["time"].([]interface{})) > 0 {
			TimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime(modelMapElement["time"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Time = TimeModel
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime{}
	if modelMap["interval"] != nil && modelMap["interval"].(string) != "" {
		model.Interval = core.StringPtr(modelMap["interval"].(string))
	}
	if modelMap["buckets_presented"] != nil {
		model.BucketsPresented = core.Int64Ptr(int64(modelMap["buckets_presented"].(int)))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisTypeValue(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue{}
	if modelMap["value"] != nil && len(modelMap["value"].([]interface{})) > 0 {
		ValueModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty(modelMap["value"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Value = ValueModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisTypeTime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime{}
	if modelMap["time"] != nil && len(modelMap["time"].([]interface{})) > 0 {
		TimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime(modelMap["time"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Time = TimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChart(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChart{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["query"] != nil && len(modelMapElement["query"].([]interface{})) > 0 {
			QueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartQuery(modelMapElement["query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Query = QueryModel
		}
		if modelMapElement["max_bars_per_chart"] != nil {
			model.MaxBarsPerChart = core.Int64Ptr(int64(modelMapElement["max_bars_per_chart"].(int)))
		}
		if modelMapElement["group_name_template"] != nil && modelMapElement["group_name_template"].(string) != "" {
			model.GroupNameTemplate = core.StringPtr(modelMapElement["group_name_template"].(string))
		}
		if modelMapElement["stack_definition"] != nil && len(modelMapElement["stack_definition"].([]interface{})) > 0 {
			StackDefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition(modelMapElement["stack_definition"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.StackDefinition = StackDefinitionModel
		}
		if modelMapElement["scale_type"] != nil && modelMapElement["scale_type"].(string) != "" {
			model.ScaleType = core.StringPtr(modelMapElement["scale_type"].(string))
		}
		if modelMapElement["colors_by"] != nil && len(modelMapElement["colors_by"].([]interface{})) > 0 {
			ColorsByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonColorsBy(modelMapElement["colors_by"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.ColorsBy = ColorsByModel
		}
		if modelMapElement["unit"] != nil && modelMapElement["unit"].(string) != "" {
			model.Unit = core.StringPtr(modelMapElement["unit"].(string))
		}
		if modelMapElement["display_on_bar"] != nil {
			model.DisplayOnBar = core.BoolPtr(modelMapElement["display_on_bar"].(bool))
		}
		if modelMapElement["y_axis_view_by"] != nil && len(modelMapElement["y_axis_view_by"].([]interface{})) > 0 {
			YAxisViewByModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy(modelMapElement["y_axis_view_by"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.YAxisViewBy = YAxisViewByModel
		}
		if modelMapElement["sort_by"] != nil && modelMapElement["sort_by"].(string) != "" {
			model.SortBy = core.StringPtr(modelMapElement["sort_by"].(string))
		}
		model.ColorScheme = core.StringPtr(modelMapElement["color_scheme"].(string))
		if modelMapElement["data_mode_type"] != nil && modelMapElement["data_mode_type"].(string) != "" {
			model.DataModeType = core.StringPtr(modelMapElement["data_mode_type"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartQuery(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery(modelMapElement["logs"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery(modelMapElement["metrics"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["dataprime"] != nil && len(modelMapElement["dataprime"].([]interface{})) > 0 {
			DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery(modelMapElement["dataprime"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Dataprime = DataprimeModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery{}
	if modelMap["lucene_query"] != nil && len(modelMap["lucene_query"].([]interface{})) > 0 {
		LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonLuceneQuery(modelMap["lucene_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LuceneQuery = LuceneQueryModel
	}
	if modelMap["aggregation"] != nil && len(modelMap["aggregation"].([]interface{})) > 0 {
		AggregationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLogsAggregation(modelMap["aggregation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Aggregation = AggregationModel
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterLogsFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterLogsFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names_fields"] != nil {
		// groupNamesFields := []logsv0.ApisDashboardsV1CommonObservationField{}
		// for _, groupNamesFieldsItem := range modelMap["group_names_fields"].([]interface{}) {
		// 	groupNamesFieldsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(groupNamesFieldsItem.(map[string]interface{}))
		// 	if err != nil {
		// 		return model, err
		// 	}
		// 	groupNamesFields = append(groupNamesFields, *groupNamesFieldsItemModel)
		// }
		groupNamesFields, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["group_names_fields"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GroupNamesFields = groupNamesFields
	}
	if modelMap["stacked_group_name_field"] != nil && len(modelMap["stacked_group_name_field"].([]interface{})) > 0 && modelMap["stacked_group_name_field"].([]interface{})[0] != nil {
		StackedGroupNameFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMap["stacked_group_name_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StackedGroupNameField = &StackedGroupNameFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery{}
	if modelMap["promql_query"] != nil && len(modelMap["promql_query"].([]interface{})) > 0 {
		PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsCommonPromQlQuery(modelMap["promql_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PromqlQuery = PromqlQueryModel
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterMetricsFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterMetricsFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names"] != nil {
		groupNames := []string{}
		for _, groupNamesItem := range modelMap["group_names"].([]interface{}) {
			groupNames = append(groupNames, groupNamesItem.(string))
		}
		model.GroupNames = groupNames
	}
	if modelMap["stacked_group_name"] != nil && modelMap["stacked_group_name"].(string) != "" {
		model.StackedGroupName = core.StringPtr(modelMap["stacked_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery{}
	if modelMap["dataprime_query"] != nil && len(modelMap["dataprime_query"].([]interface{})) > 0 {
		DataprimeQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonDataprimeQuery(modelMap["dataprime_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DataprimeQuery = DataprimeQueryModel
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilterSourceIntf{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["group_names"] != nil {
		groupNames := []string{}
		for _, groupNamesItem := range modelMap["group_names"].([]interface{}) {
			groupNames = append(groupNames, groupNamesItem.(string))
		}
		model.GroupNames = groupNames
	}
	if modelMap["stacked_group_name"] != nil && modelMap["stacked_group_name"].(string) != "" {
		model.StackedGroupName = core.StringPtr(modelMap["stacked_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery(modelMap["logs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery(modelMap["metrics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime{}
	if modelMap["dataprime"] != nil && len(modelMap["dataprime"].([]interface{})) > 0 {
		DataprimeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery(modelMap["dataprime"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Dataprime = DataprimeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["max_slices_per_bar"] != nil {
			model.MaxSlicesPerBar = core.Int64Ptr(int64(modelMapElement["max_slices_per_bar"].(int)))
		}
		if modelMapElement["stack_name_template"] != nil && modelMapElement["stack_name_template"].(string) != "" {
			model.StackNameTemplate = core.StringPtr(modelMapElement["stack_name_template"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy(modelMap []interface{}) (logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByIntf, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["category"] != nil && len(modelMapElement["category"].([]interface{})) > 0 {
			CategoryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty(modelMapElement["category"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Category = CategoryModel
		}
		if modelMapElement["value"] != nil && len(modelMapElement["value"].([]interface{})) > 0 {
			ValueModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty(modelMapElement["value"].([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.Value = ValueModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory{}
	if modelMap["category"] != nil && len(modelMap["category"].([]interface{})) > 0 {
		CategoryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty(modelMap["category"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Category = CategoryModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue{}
	if modelMap["value"] != nil && len(modelMap["value"].([]interface{})) > 0 {
		ValueModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty(modelMap["value"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Value = ValueModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsMarkdown(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstWidgetsMarkdown, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetsMarkdown{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.MarkdownText = core.StringPtr(modelMapElement["markdown_text"].(string))
		if modelMapElement["tooltip_text"] != nil && modelMapElement["tooltip_text"].(string) != "" {
			model.TooltipText = core.StringPtr(modelMapElement["tooltip_text"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueLineChart(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart{}
	if modelMap["line_chart"] != nil && len(modelMap["line_chart"].([]interface{})) > 0 {
		LineChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsLineChart(modelMap["line_chart"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LineChart = LineChartModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueDataTable(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable{}
	if modelMap["data_table"] != nil && len(modelMap["data_table"].([]interface{})) > 0 {
		DataTableModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsDataTable(modelMap["data_table"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DataTable = DataTableModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueGauge(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge{}
	if modelMap["gauge"] != nil && len(modelMap["gauge"].([]interface{})) > 0 {
		GaugeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsGauge(modelMap["gauge"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Gauge = GaugeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValuePieChart(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart{}
	if modelMap["pie_chart"] != nil && len(modelMap["pie_chart"].([]interface{})) > 0 {
		PieChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsPieChart(modelMap["pie_chart"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.PieChart = PieChartModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueBarChart(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart{}
	if modelMap["bar_chart"] != nil && len(modelMap["bar_chart"].([]interface{})) > 0 {
		BarChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsBarChart(modelMap["bar_chart"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.BarChart = BarChartModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart{}
	if modelMap["horizontal_bar_chart"] != nil && len(modelMap["horizontal_bar_chart"].([]interface{})) > 0 {
		HorizontalBarChartModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsHorizontalBarChart(modelMap["horizontal_bar_chart"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.HorizontalBarChart = HorizontalBarChartModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetDefinitionValueMarkdown(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown, error) {
	model := &logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown{}
	if modelMap["markdown"] != nil && len(modelMap["markdown"].([]interface{})) > 0 {
		MarkdownModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstWidgetsMarkdown(modelMap["markdown"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Markdown = MarkdownModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstVariable, error) {
	model := &logsv0.ApisDashboardsV1AstVariable{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	DefinitionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariableDefinition(modelMap["definition"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	model.DisplayName = core.StringPtr(modelMap["display_name"].(string))
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariableDefinition(modelMap []interface{}) (logsv0.ApisDashboardsV1AstVariableDefinitionIntf, error) {
	model := &logsv0.ApisDashboardsV1AstVariableDefinition{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["multi_select"] != nil && len(modelMapElement["multi_select"].([]interface{})) > 0 {
			MultiSelectModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelect(modelMapElement["multi_select"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.MultiSelect = MultiSelectModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelect(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelect, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelect{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		SourceModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSource(modelMapElement["source"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Source = SourceModel
		SelectionModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelection(modelMapElement["selection"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Selection = SelectionModel
		model.ValuesOrderDirection = core.StringPtr(modelMapElement["values_order_direction"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSource(modelMap []interface{}) (logsv0.ApisDashboardsV1AstMultiSelectSourceIntf, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["logs_path"] != nil && len(modelMapElement["logs_path"].([]interface{})) > 0 {
			LogsPathModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectLogsPathSource(modelMapElement["logs_path"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LogsPath = LogsPathModel
		}
		if modelMapElement["metric_label"] != nil && len(modelMapElement["metric_label"].([]interface{})) > 0 {
			MetricLabelModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectMetricLabelSource(modelMapElement["metric_label"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.MetricLabel = MetricLabelModel
		}
		if modelMapElement["constant_list"] != nil && len(modelMapElement["constant_list"].([]interface{})) > 0 {
			ConstantListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectConstantListSource(modelMapElement["constant_list"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.ConstantList = ConstantListModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectLogsPathSource(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectLogsPathSource, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectLogsPathSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		ObservationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["observation_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ObservationField = &ObservationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectMetricLabelSource(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectMetricLabelSource, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectMetricLabelSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.MetricName = core.StringPtr(modelMapElement["metric_name"].(string))
		model.Label = core.StringPtr(modelMapElement["label"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectConstantListSource(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectConstantListSource, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectConstantListSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		values := []string{}
		for _, valuesItem := range modelMapElement["values"].([]interface{}) {
			values = append(values, valuesItem.(string))
		}
		model.Values = values
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSourceValueLogsPath(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath{}
	if modelMap["logs_path"] != nil && len(modelMap["logs_path"].([]interface{})) > 0 {
		LogsPathModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectLogsPathSource(modelMap["logs_path"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LogsPath = LogsPathModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSourceValueMetricLabel(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel{}
	if modelMap["metric_label"] != nil && len(modelMap["metric_label"].([]interface{})) > 0 {
		MetricLabelModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectMetricLabelSource(modelMap["metric_label"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MetricLabel = MetricLabelModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSourceValueConstantList(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList{}
	if modelMap["constant_list"] != nil && len(modelMap["constant_list"].([]interface{})) > 0 {
		ConstantListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectConstantListSource(modelMap["constant_list"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.ConstantList = ConstantListModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelection(modelMap []interface{}) (logsv0.ApisDashboardsV1AstMultiSelectSelectionIntf, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSelection{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["all"] != nil && len(modelMapElement["all"].([]interface{})) > 0 {
			AllModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty(modelMapElement["all"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.All = AllModel
		}
		if modelMapElement["list"] != nil && len(modelMapElement["list"].([]interface{})) > 0 {
			ListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionListSelection(modelMapElement["list"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.List = ListModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionListSelection(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSelectionListSelection, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSelectionListSelection{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["values"] != nil {
			values := []string{}
			for _, valuesItem := range modelMapElement["values"].([]interface{}) {
				values = append(values, valuesItem.(string))
			}
			model.Values = values
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionValueAll(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll{}
	if modelMap["all"] != nil && len(modelMap["all"].([]interface{})) > 0 {
		AllModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty(modelMap["all"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.All = AllModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionValueList(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList, error) {
	model := &logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList{}
	if modelMap["list"] != nil && len(modelMap["list"].([]interface{})) > 0 {
		ListModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelectSelectionListSelection(modelMap["list"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.List = ListModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariableDefinitionValueMultiSelect(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect, error) {
	model := &logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect{}
	if modelMap["multi_select"] != nil && len(modelMap["multi_select"].([]interface{})) > 0 {
		MultiSelectModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstMultiSelect(modelMap["multi_select"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MultiSelect = MultiSelectModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstFilter, error) {
	model := &logsv0.ApisDashboardsV1AstFilter{}
	if modelMap["source"] != nil && len(modelMap["source"].([]interface{})) > 0 {
		SourceModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilterSource(modelMap["source"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Source = SourceModel
	}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["collapsed"] != nil {
		model.Collapsed = core.BoolPtr(modelMap["collapsed"].(bool))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotation, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotation{}
	if modelMap["href"] != nil {
		model.Href = core.UUIDPtr(strfmt.UUID(modelMap["href"].(string)))
	}
	model.ID = core.UUIDPtr(strfmt.UUID(modelMap["id"].(string)))
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	SourceModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationSource(modelMap["source"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Source = SourceModel
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationSource(modelMap []interface{}) (logsv0.ApisDashboardsV1AstAnnotationSourceIntf, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["metrics"] != nil && len(modelMapElement["metrics"].([]interface{})) > 0 {
			MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSource(modelMapElement["metrics"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Metrics = MetricsModel
		}
		if modelMapElement["logs"] != nil && len(modelMapElement["logs"].([]interface{})) > 0 {
			LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSource(modelMapElement["logs"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Logs = LogsModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSource(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationMetricsSource, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationMetricsSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["promql_query"] != nil && len(modelMapElement["promql_query"].([]interface{})) > 0 {
			PromqlQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonPromQlQuery(modelMapElement["promql_query"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.PromqlQuery = PromqlQueryModel
		}
		if modelMapElement["strategy"] != nil && len(modelMapElement["strategy"].([]interface{})) > 0 {
			StrategyModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSourceStrategy(modelMapElement["strategy"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Strategy = StrategyModel
		}
		if modelMapElement["message_template"] != nil && modelMapElement["message_template"].(string) != "" {
			model.MessageTemplate = core.StringPtr(modelMapElement["message_template"].(string))
		}
		if modelMapElement["labels"] != nil {
			labels := []string{}
			for _, labelsItem := range modelMapElement["labels"].([]interface{}) {
				labels = append(labels, labelsItem.(string))
			}
			model.Labels = labels
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonPromQlQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonPromQlQuery, error) {
	model := &logsv0.ApisDashboardsV1CommonPromQlQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["value"] != nil && modelMapElement["value"].(string) != "" {
			model.Value = core.StringPtr(modelMapElement["value"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSourceStrategy(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStrategy, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStrategy{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["start_time_metric"] != nil && len(modelMapElement["start_time_metric"].([]interface{})) > 0 {
			StartTimeMetricModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty(modelMapElement["start_time_metric"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.StartTimeMetric = StartTimeMetricModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSource(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSource, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSource{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		LuceneQueryModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLuceneQuery(modelMapElement["lucene_query"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LuceneQuery = LuceneQueryModel
		StrategyModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategy(modelMapElement["strategy"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Strategy = StrategyModel
		if modelMapElement["message_template"] != nil && modelMapElement["message_template"].(string) != "" {
			model.MessageTemplate = core.StringPtr(modelMapElement["message_template"].(string))
		}
		if modelMapElement["label_fields"] != nil {
			// labelFields := []logsv0.ApisDashboardsV1CommonObservationField{}
			// for _, labelFieldsItem := range modelMapElement["label_fields"].([]interface{}) {
			// 	labelFieldsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(labelFieldsItem.(map[string]interface{}))
			// 	if err != nil {
			// 		return model, err
			// 	}
			// 	labelFields = append(labelFields, *labelFieldsItemModel)
			// }
			labelFields, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapElement["label_fields"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.LabelFields = labelFields
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonLuceneQuery(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonLuceneQuery, error) {
	model := &logsv0.ApisDashboardsV1CommonLuceneQuery{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["value"] != nil && modelMapElement["value"].(string) != "" {
			model.Value = core.StringPtr(modelMapElement["value"].(string))
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategy(modelMap []interface{}) (logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyIntf, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategy{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["instant"] != nil && len(modelMapElement["instant"].([]interface{})) > 0 {
			InstantModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyInstant(modelMapElement["instant"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Instant = InstantModel
		}
		if modelMapElement["range"] != nil && len(modelMapElement["range"].([]interface{})) > 0 {
			RangeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyRange(modelMapElement["range"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Range = RangeModel
		}
		if modelMapElement["duration"] != nil && len(modelMapElement["duration"].([]interface{})) > 0 {
			DurationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyDuration(modelMapElement["duration"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Duration = DurationModel
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyInstant(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyInstant, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyInstant{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		TimestampFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["timestamp_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.TimestampField = &TimestampFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyRange(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyRange, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyRange{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		StartTimestampFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["start_timestamp_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StartTimestampField = &StartTimestampFieldModel[0]

		EndTimestampFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["end_timestamp_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.EndTimestampField = &EndTimestampFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyDuration(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyDuration, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyDuration{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapObsevationField := modelMap[0].(map[string]interface{})
		StartTimestampFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["start_timestamp_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.StartTimestampField = &StartTimestampFieldModel[0]
		DurationFieldModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonObservationField(modelMapObsevationField["duration_field"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.DurationField = &DurationFieldModel[0]
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant{}
	if modelMap["instant"] != nil && len(modelMap["instant"].([]interface{})) > 0 {
		InstantModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyInstant(modelMap["instant"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Instant = InstantModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange{}
	if modelMap["range"] != nil && len(modelMap["range"].([]interface{})) > 0 {
		RangeModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyRange(modelMap["range"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Range = RangeModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration{}
	if modelMap["duration"] != nil && len(modelMap["duration"].([]interface{})) > 0 {
		DurationModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSourceStrategyDuration(modelMap["duration"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Duration = DurationModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationSourceValueMetrics(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics{}
	if modelMap["metrics"] != nil && len(modelMap["metrics"].([]interface{})) > 0 {
		MetricsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationMetricsSource(modelMap["metrics"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metrics = MetricsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationSourceValueLogs(modelMap map[string]interface{}) (*logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs, error) {
	model := &logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs{}
	if modelMap["logs"] != nil && len(modelMap["logs"].([]interface{})) > 0 {
		LogsModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotationLogsSource(modelMap["logs"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Logs = LogsModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1CommonTimeFrame(modelMap []interface{}) (*logsv0.ApisDashboardsV1CommonTimeFrame, error) {
	model := &logsv0.ApisDashboardsV1CommonTimeFrame{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["from"] != nil {
			dateTime, err := core.ParseDateTime(modelMapElement["from"].(string))
			if err != nil {
				return model, err
			}
			model.From = &dateTime
		}
		if modelMapElement["to"] != nil {
			dateTime, err := core.ParseDateTime(modelMapElement["to"].(string))
			if err != nil {
				return model, err
			}
			model.To = &dateTime
		}
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstFolderPath(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstFolderPath, error) {
	model := &logsv0.ApisDashboardsV1AstFolderPath{}
	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["segments"] != nil {
			segments := []string{}
			for _, segmentsItem := range modelMapElement["segments"].([]interface{}) {
				segments = append(segments, segmentsItem.(string))
			}
			model.Segments = segments
		}
	}

	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshOffEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstDashboardAutoRefreshOffEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstDashboardAutoRefreshOffEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty(modelMap []interface{}) (*logsv0.ApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty, error) {
	model := &logsv0.ApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty{}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboard(modelMap map[string]interface{}) (logsv0.DashboardIntf, error) {
	model := &logsv0.Dashboard{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["absolute_time_frame"] != nil && len(modelMap["absolute_time_frame"].([]interface{})) > 0 {
		AbsoluteTimeFrameModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonTimeFrame(modelMap["absolute_time_frame"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.AbsoluteTimeFrame = AbsoluteTimeFrameModel
	}
	if modelMap["relative_time_frame"] != nil && modelMap["relative_time_frame"].(string) != "" {
		model.RelativeTimeFrame = core.StringPtr(modelMap["relative_time_frame"].(string))
	}
	if modelMap["folder_id"] != nil && len(modelMap["folder_id"].([]interface{})) > 0 {
		FolderIDModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap["folder_id"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.FolderID = FolderIDModel
	}
	if modelMap["folder_path"] != nil && len(modelMap["folder_path"].([]interface{})) > 0 {
		FolderPathModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFolderPath(modelMap["folder_path"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.FolderPath = FolderPathModel
	}
	if modelMap["false"] != nil && len(modelMap["false"].([]interface{})) > 0 {
		FalseModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshOffEmpty(modelMap["false"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.False = FalseModel
	}
	if modelMap["two_minutes"] != nil && len(modelMap["two_minutes"].([]interface{})) > 0 {
		TwoMinutesModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty(modelMap["two_minutes"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.TwoMinutes = TwoMinutesModel
	}
	if modelMap["five_minutes"] != nil && len(modelMap["five_minutes"].([]interface{})) > 0 {
		FiveMinutesModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty(modelMap["five_minutes"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.FiveMinutes = FiveMinutesModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardTimeFrameAbsoluteTimeFrame(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardTimeFrameAbsoluteTimeFrame, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardTimeFrameAbsoluteTimeFrame{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["absolute_time_frame"] != nil && len(modelMap["absolute_time_frame"].([]interface{})) > 0 {
		AbsoluteTimeFrameModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1CommonTimeFrame(modelMap["absolute_time_frame"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.AbsoluteTimeFrame = AbsoluteTimeFrameModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardTimeFrameRelativeTimeFrame(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardTimeFrameRelativeTimeFrame, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardTimeFrameRelativeTimeFrame{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["relative_time_frame"] != nil && modelMap["relative_time_frame"].(string) != "" {
		model.RelativeTimeFrame = core.StringPtr(modelMap["relative_time_frame"].(string))
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardFolderFolderID(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardFolderFolderID, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardFolderFolderID{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["folder_id"] != nil && len(modelMap["folder_id"].([]interface{})) > 0 {
		FolderIDModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1UUID(modelMap["folder_id"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.FolderID = FolderIDModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardFolderFolderPath(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardFolderFolderPath, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardFolderFolderPath{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["folder_path"] != nil && len(modelMap["folder_path"].([]interface{})) > 0 {
		FolderPathModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFolderPath(modelMap["folder_path"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.FolderPath = FolderPathModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardAutoRefreshOff(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshOff, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshOff{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["false"] != nil && len(modelMap["false"].([]interface{})) > 0 {
		FalseModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshOffEmpty(modelMap["false"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.False = FalseModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutes(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutes, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutes{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["two_minutes"] != nil && len(modelMap["two_minutes"].([]interface{})) > 0 {
		TwoMinutesModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty(modelMap["two_minutes"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.TwoMinutes = TwoMinutesModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardMapToDashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutes(modelMap map[string]interface{}) (*logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutes, error) {
	model := &logsv0.DashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutes{}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	LayoutModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstLayout(modelMap["layout"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Layout = LayoutModel
	if modelMap["variables"] != nil {
		variables := []logsv0.ApisDashboardsV1AstVariable{}
		for _, variablesItem := range modelMap["variables"].([]interface{}) {
			variablesItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstVariable(variablesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			variables = append(variables, *variablesItemModel)
		}
		model.Variables = variables
	}
	if modelMap["filters"] != nil {
		filters := []logsv0.ApisDashboardsV1AstFilter{}
		for _, filtersItem := range modelMap["filters"].([]interface{}) {
			filtersItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstFilter(filtersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filters = append(filters, *filtersItemModel)
		}
		model.Filters = filters
	}
	if modelMap["annotations"] != nil {
		annotations := []logsv0.ApisDashboardsV1AstAnnotation{}
		for _, annotationsItem := range modelMap["annotations"].([]interface{}) {
			annotationsItemModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstAnnotation(annotationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			annotations = append(annotations, *annotationsItemModel)
		}
		model.Annotations = annotations
	}
	if modelMap["five_minutes"] != nil && len(modelMap["five_minutes"].([]interface{})) > 0 {
		FiveMinutesModel, err := ResourceIbmLogsDashboardMapToApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty(modelMap["five_minutes"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.FiveMinutes = FiveMinutesModel
	}
	return model, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstLayoutToMap(model *logsv0.ApisDashboardsV1AstLayout) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Sections != nil {
		sections := []map[string]interface{}{}
		for _, sectionsItem := range model.Sections {
			sectionsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstSectionToMap(&sectionsItem)
			if err != nil {
				return modelMap, err
			}
			sections = append(sections, sectionsItemMap)
		}
		modelMap["sections"] = sections
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstSectionToMap(model *logsv0.ApisDashboardsV1AstSection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := ResourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	if model.Rows != nil {
		rows := []map[string]interface{}{}
		for _, rowsItem := range model.Rows {
			rowsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstRowToMap(&rowsItem)
			if err != nil {
				return modelMap, err
			}
			rows = append(rows, rowsItemMap)
		}
		modelMap["rows"] = rows
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model *logsv0.ApisDashboardsV1UUID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = model.Value.String()
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstRowToMap(model *logsv0.ApisDashboardsV1AstRow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := ResourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	appearanceMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstRowAppearanceToMap(model.Appearance)
	if err != nil {
		return modelMap, err
	}
	modelMap["appearance"] = []map[string]interface{}{appearanceMap}
	if model.Widgets != nil {
		widgets := []map[string]interface{}{}
		for _, widgetsItem := range model.Widgets {
			widgetsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetToMap(&widgetsItem)
			if err != nil {
				return modelMap, err
			}
			widgets = append(widgets, widgetsItemMap)
		}
		modelMap["widgets"] = widgets
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstRowAppearanceToMap(model *logsv0.ApisDashboardsV1AstRowAppearance) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["height"] = flex.IntValue(model.Height)
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetToMap(model *logsv0.ApisDashboardsV1AstWidget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := ResourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	modelMap["title"] = *model.Title
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	definitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionToMap(model logsv0.ApisDashboardsV1AstWidgetDefinitionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueLineChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueDataTableToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueGaugeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValuePieChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueBarChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueMarkdownToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetDefinition)
		if model.LineChart != nil {
			lineChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model.LineChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["line_chart"] = []map[string]interface{}{lineChartMap}
		}
		if model.DataTable != nil {
			dataTableMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model.DataTable)
			if err != nil {
				return modelMap, err
			}
			modelMap["data_table"] = []map[string]interface{}{dataTableMap}
		}
		if model.Gauge != nil {
			gaugeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model.Gauge)
			if err != nil {
				return modelMap, err
			}
			modelMap["gauge"] = []map[string]interface{}{gaugeMap}
		}
		if model.PieChart != nil {
			pieChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model.PieChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["pie_chart"] = []map[string]interface{}{pieChartMap}
		}
		if model.BarChart != nil {
			barChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model.BarChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["bar_chart"] = []map[string]interface{}{barChartMap}
		}
		if model.HorizontalBarChart != nil {
			horizontalBarChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model.HorizontalBarChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["horizontal_bar_chart"] = []map[string]interface{}{horizontalBarChartMap}
		}
		if model.Markdown != nil {
			markdownMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model.Markdown)
			if err != nil {
				return modelMap, err
			}
			modelMap["markdown"] = []map[string]interface{}{markdownMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetDefinitionIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	legendMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLegendToMap(model.Legend)
	if err != nil {
		return modelMap, err
	}
	modelMap["legend"] = []map[string]interface{}{legendMap}
	tooltipMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartTooltipToMap(model.Tooltip)
	if err != nil {
		return modelMap, err
	}
	modelMap["tooltip"] = []map[string]interface{}{tooltipMap}
	if model.QueryDefinitions != nil {
		queryDefinitions := []map[string]interface{}{}
		for _, queryDefinitionsItem := range model.QueryDefinitions {
			queryDefinitionsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryDefinitionToMap(&queryDefinitionsItem)
			if err != nil {
				return modelMap, err
			}
			queryDefinitions = append(queryDefinitions, queryDefinitionsItemMap)
		}
		modelMap["query_definitions"] = queryDefinitions
	}
	if model.StackedLine != nil {
		modelMap["stacked_line"] = *model.StackedLine
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLegendToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonLegend) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["is_visible"] = *model.IsVisible
	if model.Columns != nil {
		modelMap["columns"] = model.Columns
	}
	modelMap["group_by_query"] = *model.GroupByQuery
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartTooltipToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartTooltip) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ShowLabels != nil {
		modelMap["show_labels"] = *model.ShowLabels
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	if model.SeriesNameTemplate != nil {
		modelMap["series_name_template"] = *model.SeriesNameTemplate
	}
	if model.SeriesCountLimit != nil {
		modelMap["series_count_limit"] = *model.SeriesCountLimit
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	if model.ScaleType != nil {
		modelMap["scale_type"] = *model.ScaleType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	modelMap["is_visible"] = *model.IsVisible
	if model.ColorScheme != nil {
		modelMap["color_scheme"] = *model.ColorScheme
	}
	resolutionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartResolutionToMap(model.Resolution)
	if err != nil {
		return modelMap, err
	}
	modelMap["resolution"] = []map[string]interface{}{resolutionMap}
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsLineChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQuery)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsLineChartQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	if model.Aggregations != nil {
		aggregations := []map[string]interface{}{}
		for _, aggregationsItem := range model.Aggregations {
			aggregationsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(aggregationsItem)
			if err != nil {
				return modelMap, err
			}
			aggregations = append(aggregations, aggregationsItemMap)
		}
		modelMap["aggregations"] = aggregations
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupBys != nil {
		groupBys := []map[string]interface{}{}
		for _, groupBysItem := range model.GroupBys {
			groupBysItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupBysItem)
			if err != nil {
				return modelMap, err
			}
			groupBys = append(groupBys, groupBysItemMap)
		}
		modelMap["group_bys"] = groupBys
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonLuceneQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model logsv0.ApisDashboardsV1CommonLogsAggregationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCount); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCount))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountDistinctToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueSum); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueSumToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueSum))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueAverageToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMin); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMinToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMin))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMax); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMaxToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMax))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValuePercentileToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregation); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1CommonLogsAggregation)
		if model.Count != nil {
			countMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model.Count)
			if err != nil {
				return modelMap, err
			}
			modelMap["count"] = []map[string]interface{}{countMap}
		}
		if model.CountDistinct != nil {
			countDistinctMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model.CountDistinct)
			if err != nil {
				return modelMap, err
			}
			modelMap["count_distinct"] = []map[string]interface{}{countDistinctMap}
		}
		if model.Sum != nil {
			sumMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model.Sum)
			if err != nil {
				return modelMap, err
			}
			modelMap["sum"] = []map[string]interface{}{sumMap}
		}
		if model.Average != nil {
			averageMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model.Average)
			if err != nil {
				return modelMap, err
			}
			modelMap["average"] = []map[string]interface{}{averageMap}
		}
		if model.Min != nil {
			minMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model.Min)
			if err != nil {
				return modelMap, err
			}
			modelMap["min"] = []map[string]interface{}{minMap}
		}
		if model.Max != nil {
			maxMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model.Max)
			if err != nil {
				return modelMap, err
			}
			modelMap["max"] = []map[string]interface{}{maxMap}
		}
		if model.Percentile != nil {
			percentileMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model.Percentile)
			if err != nil {
				return modelMap, err
			}
			modelMap["percentile"] = []map[string]interface{}{percentileMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1CommonLogsAggregationIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationCountEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationCountDistinct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model *logsv0.ApisDashboardsV1CommonObservationField) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Keypath != nil {
		modelMap["keypath"] = model.Keypath
	}
	if model.Scope != nil {
		modelMap["scope"] = *model.Scope
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationSum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationAverage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationMin) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationMax) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationPercentile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["percent"] = *model.Percent
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueCount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Count != nil {
		countMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model.Count)
		if err != nil {
			return modelMap, err
		}
		modelMap["count"] = []map[string]interface{}{countMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountDistinctToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CountDistinct != nil {
		countDistinctMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model.CountDistinct)
		if err != nil {
			return modelMap, err
		}
		modelMap["count_distinct"] = []map[string]interface{}{countDistinctMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueSumToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueSum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Sum != nil {
		sumMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model.Sum)
		if err != nil {
			return modelMap, err
		}
		modelMap["sum"] = []map[string]interface{}{sumMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueAverageToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Average != nil {
		averageMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model.Average)
		if err != nil {
			return modelMap, err
		}
		modelMap["average"] = []map[string]interface{}{averageMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMinToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueMin) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Min != nil {
		minMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model.Min)
		if err != nil {
			return modelMap, err
		}
		modelMap["min"] = []map[string]interface{}{minMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMaxToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueMax) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Max != nil {
		maxMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model.Max)
		if err != nil {
			return modelMap, err
		}
		modelMap["max"] = []map[string]interface{}{maxMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValuePercentileToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Percentile != nil {
		percentileMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model.Percentile)
		if err != nil {
			return modelMap, err
		}
		modelMap["percentile"] = []map[string]interface{}{percentileMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model *logsv0.ApisDashboardsV1AstFilterLogsFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Operator != nil {
		operatorMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model.Operator)
		if err != nil {
			return modelMap, err
		}
		modelMap["operator"] = []map[string]interface{}{operatorMap}
	}
	if model.ObservationField != nil {
		observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
		if err != nil {
			return modelMap, err
		}
		modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model logsv0.ApisDashboardsV1AstFilterOperatorIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueEquals); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueEqualsToMap(model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueEquals))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueNotEqualsToMap(model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperator); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterOperator)
		if model.Equals != nil {
			equalsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model.Equals)
			if err != nil {
				return modelMap, err
			}
			modelMap["equals"] = []map[string]interface{}{equalsMap}
		}
		if model.NotEquals != nil {
			notEqualsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model.NotEquals)
			if err != nil {
				return modelMap, err
			}
			modelMap["not_equals"] = []map[string]interface{}{notEqualsMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstFilterOperatorIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Selection != nil {
		selectionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionToMap(model.Selection)
		if err != nil {
			return modelMap, err
		}
		modelMap["selection"] = []map[string]interface{}{selectionMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionToMap(model logsv0.ApisDashboardsV1AstFilterEqualsSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueAllToMap(model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueListToMap(model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelection)
		if model.All != nil {
			allMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model.All)
			if err != nil {
				return modelMap, err
			}
			modelMap["all"] = []map[string]interface{}{allMap}
		}
		if model.List != nil {
			listMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model.List)
			if err != nil {
				return modelMap, err
			}
			modelMap["list"] = []map[string]interface{}{listMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstFilterEqualsSelectionIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueAllToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.All != nil {
		allMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model.All)
		if err != nil {
			return modelMap, err
		}
		modelMap["all"] = []map[string]interface{}{allMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueListToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterNotEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Selection != nil {
		selectionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionToMap(model.Selection)
		if err != nil {
			return modelMap, err
		}
		modelMap["selection"] = []map[string]interface{}{selectionMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterNotEqualsSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterNotEqualsSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterOperatorValueEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Equals != nil {
		equalsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model.Equals)
		if err != nil {
			return modelMap, err
		}
		modelMap["equals"] = []map[string]interface{}{equalsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueNotEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotEquals != nil {
		notEqualsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model.NotEquals)
		if err != nil {
			return modelMap, err
		}
		modelMap["not_equals"] = []map[string]interface{}{notEqualsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonPromQlQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model *logsv0.ApisDashboardsV1AstFilterMetricsFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	if model.Operator != nil {
		operatorMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model.Operator)
		if err != nil {
			return modelMap, err
		}
		modelMap["operator"] = []map[string]interface{}{operatorMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model *logsv0.ApisDashboardsV1CommonDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Text != nil {
		modelMap["text"] = *model.Text
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(model logsv0.ApisDashboardsV1AstFilterSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSourceValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstFilterSourceValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSourceValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstFilterSourceValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterSource)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstFilterSourceIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueLogsToMap(model *logsv0.ApisDashboardsV1AstFilterSourceValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueMetricsToMap(model *logsv0.ApisDashboardsV1AstFilterSourceValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartResolutionToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartResolution) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.BucketsPresented != nil {
		modelMap["buckets_presented"] = flex.IntValue(model.BucketsPresented)
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["results_per_page"] = flex.IntValue(model.ResultsPerPage)
	modelMap["row_style"] = *model.RowStyle
	if model.Columns != nil {
		columns := []map[string]interface{}{}
		for _, columnsItem := range model.Columns {
			columnsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableColumnToMap(&columnsItem)
			if err != nil {
				return modelMap, err
			}
			columns = append(columns, columnsItemMap)
		}
		modelMap["columns"] = columns
	}
	if model.OrderBy != nil {
		orderByMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonOrderingFieldToMap(model.OrderBy)
		if err != nil {
			return modelMap, err
		}
		modelMap["order_by"] = []map[string]interface{}{orderByMap}
	}
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsDataTableQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQuery)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsDataTableQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.Grouping != nil {
		groupingMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryGroupingToMap(model.Grouping)
		if err != nil {
			return modelMap, err
		}
		modelMap["grouping"] = []map[string]interface{}{groupingMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryGroupingToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Aggregations != nil {
		aggregations := []map[string]interface{}{}
		for _, aggregationsItem := range model.Aggregations {
			aggregationsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryAggregationToMap(&aggregationsItem)
			if err != nil {
				return modelMap, err
			}
			aggregations = append(aggregations, aggregationsItemMap)
		}
		modelMap["aggregations"] = aggregations
	}
	if model.GroupBys != nil {
		groupBys := []map[string]interface{}{}
		for _, groupBysItem := range model.GroupBys {
			groupBysItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupBysItem)
			if err != nil {
				return modelMap, err
			}
			groupBys = append(groupBys, groupBysItemMap)
		}
		modelMap["group_bys"] = groupBys
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryAggregationToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	modelMap["is_visible"] = *model.IsVisible
	aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableColumnToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableColumn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["field"] = *model.Field
	if model.Width != nil {
		modelMap["width"] = flex.IntValue(model.Width)
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonOrderingFieldToMap(model *logsv0.ApisDashboardsV1CommonOrderingField) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Field != nil {
		modelMap["field"] = *model.Field
	}
	if model.OrderDirection != nil {
		modelMap["order_direction"] = *model.OrderDirection
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model *logsv0.ApisDashboardsV1AstWidgetsGauge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["min"] = *model.Min
	modelMap["max"] = *model.Max
	modelMap["show_inner_arc"] = *model.ShowInnerArc
	modelMap["show_outer_arc"] = *model.ShowOuterArc
	modelMap["unit"] = *model.Unit
	if model.Thresholds != nil {
		thresholds := []map[string]interface{}{}
		for _, thresholdsItem := range model.Thresholds {
			thresholdsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeThresholdToMap(&thresholdsItem)
			if err != nil {
				return modelMap, err
			}
			thresholds = append(thresholds, thresholdsItemMap)
		}
		modelMap["thresholds"] = thresholds
	}
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	modelMap["threshold_by"] = *model.ThresholdBy
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsGaugeQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQuery)
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsGaugeQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	modelMap["aggregation"] = *model.Aggregation
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.LogsAggregation != nil {
		logsAggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.LogsAggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs_aggregation"] = []map[string]interface{}{logsAggregationMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeThresholdToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeThreshold) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["from"] = *model.From
	modelMap["color"] = *model.Color
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["max_slices_per_chart"] = flex.IntValue(model.MaxSlicesPerChart)
	modelMap["min_slice_percentage"] = flex.IntValue(model.MinSlicePercentage)
	stackDefinitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartStackDefinitionToMap(model.StackDefinition)
	if err != nil {
		return modelMap, err
	}
	modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	labelDefinitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLabelDefinitionToMap(model.LabelDefinition)
	if err != nil {
		return modelMap, err
	}
	modelMap["label_definition"] = []map[string]interface{}{labelDefinitionMap}
	modelMap["show_legend"] = *model.ShowLegend
	if model.GroupNameTemplate != nil {
		modelMap["group_name_template"] = *model.GroupNameTemplate
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	modelMap["color_scheme"] = *model.ColorScheme
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsPieChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQuery)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsPieChartQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNamesFields != nil {
		groupNamesFields := []map[string]interface{}{}
		for _, groupNamesFieldsItem := range model.GroupNamesFields {
			groupNamesFieldsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerStack != nil {
		modelMap["max_slices_per_stack"] = flex.IntValue(model.MaxSlicesPerStack)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLabelDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartLabelDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LabelSource != nil {
		modelMap["label_source"] = *model.LabelSource
	}
	if model.IsVisible != nil {
		modelMap["is_visible"] = *model.IsVisible
	}
	if model.ShowName != nil {
		modelMap["show_name"] = *model.ShowName
	}
	if model.ShowValue != nil {
		modelMap["show_value"] = *model.ShowValue
	}
	if model.ShowPercentage != nil {
		modelMap["show_percentage"] = *model.ShowPercentage
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["max_bars_per_chart"] = flex.IntValue(model.MaxBarsPerChart)
	modelMap["group_name_template"] = *model.GroupNameTemplate
	stackDefinitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartStackDefinitionToMap(model.StackDefinition)
	if err != nil {
		return modelMap, err
	}
	modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	modelMap["scale_type"] = *model.ScaleType
	colorsByMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model.ColorsBy)
	if err != nil {
		return modelMap, err
	}
	modelMap["colors_by"] = []map[string]interface{}{colorsByMap}
	xAxisMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisToMap(model.XAxis)
	if err != nil {
		return modelMap, err
	}
	modelMap["x_axis"] = []map[string]interface{}{xAxisMap}
	modelMap["unit"] = *model.Unit
	modelMap["sort_by"] = *model.SortBy
	modelMap["color_scheme"] = *model.ColorScheme
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsBarChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQuery)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsBarChartQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNamesFields != nil {
		groupNamesFields := []map[string]interface{}{}
		for _, groupNamesFieldsItem := range model.GroupNamesFields {
			groupNamesFieldsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerBar != nil {
		modelMap["max_slices_per_bar"] = flex.IntValue(model.MaxSlicesPerBar)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model logsv0.ApisDashboardsV1AstWidgetsCommonColorsByIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueStackToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueGroupByToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueAggregationToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsBy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsBy)
		if model.Stack != nil {
			stackMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model.Stack)
			if err != nil {
				return modelMap, err
			}
			modelMap["stack"] = []map[string]interface{}{stackMap}
		}
		if model.GroupBy != nil {
			groupByMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model.GroupBy)
			if err != nil {
				return modelMap, err
			}
			modelMap["group_by"] = []map[string]interface{}{groupByMap}
		}
		if model.Aggregation != nil {
			aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model.Aggregation)
			if err != nil {
				return modelMap, err
			}
			modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsCommonColorsByIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueStackToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Stack != nil {
		stackMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model.Stack)
		if err != nil {
			return modelMap, err
		}
		modelMap["stack"] = []map[string]interface{}{stackMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueGroupByToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GroupBy != nil {
		groupByMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model.GroupBy)
		if err != nil {
			return modelMap, err
		}
		modelMap["group_by"] = []map[string]interface{}{groupByMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueAggregationToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Aggregation != nil {
		aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model.Aggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisToMap(model logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeValueToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeTimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxis); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxis)
		if model.Value != nil {
			valueMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model.Value)
			if err != nil {
				return modelMap, err
			}
			modelMap["value"] = []map[string]interface{}{valueMap}
		}
		if model.Time != nil {
			timeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model.Time)
			if err != nil {
				return modelMap, err
			}
			modelMap["time"] = []map[string]interface{}{timeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.BucketsPresented != nil {
		modelMap["buckets_presented"] = flex.IntValue(model.BucketsPresented)
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeValueToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		valueMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model.Value)
		if err != nil {
			return modelMap, err
		}
		modelMap["value"] = []map[string]interface{}{valueMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeTimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Time != nil {
		timeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model.Time)
		if err != nil {
			return modelMap, err
		}
		modelMap["time"] = []map[string]interface{}{timeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Query != nil {
		queryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryToMap(model.Query)
		if err != nil {
			return modelMap, err
		}
		modelMap["query"] = []map[string]interface{}{queryMap}
	}
	if model.MaxBarsPerChart != nil {
		modelMap["max_bars_per_chart"] = flex.IntValue(model.MaxBarsPerChart)
	}
	if model.GroupNameTemplate != nil {
		modelMap["group_name_template"] = *model.GroupNameTemplate
	}
	if model.StackDefinition != nil {
		stackDefinitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinitionToMap(model.StackDefinition)
		if err != nil {
			return modelMap, err
		}
		modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	}
	if model.ScaleType != nil {
		modelMap["scale_type"] = *model.ScaleType
	}
	if model.ColorsBy != nil {
		colorsByMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model.ColorsBy)
		if err != nil {
			return modelMap, err
		}
		modelMap["colors_by"] = []map[string]interface{}{colorsByMap}
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	if model.DisplayOnBar != nil {
		modelMap["display_on_bar"] = *model.DisplayOnBar
	}
	if model.YAxisViewBy != nil {
		yAxisViewByMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByToMap(model.YAxisViewBy)
		if err != nil {
			return modelMap, err
		}
		modelMap["y_axis_view_by"] = []map[string]interface{}{yAxisViewByMap}
	}
	if model.SortBy != nil {
		modelMap["sort_by"] = *model.SortBy
	}
	modelMap["color_scheme"] = *model.ColorScheme
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQuery)
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model.Dataprime)
			if err != nil {
				return modelMap, err
			}
			modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.Aggregation != nil {
		aggregationMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNamesFields != nil {
		groupNamesFields := []map[string]interface{}{}
		for _, groupNamesFieldsItem := range model.GroupNamesFields {
			groupNamesFieldsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DataprimeQuery != nil {
		dataprimeQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.GroupNames != nil {
		modelMap["group_names"] = model.GroupNames
	}
	if model.StackedGroupName != nil {
		modelMap["stacked_group_name"] = *model.StackedGroupName
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerBar != nil {
		modelMap["max_slices_per_bar"] = flex.IntValue(model.MaxSlicesPerBar)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByToMap(model logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategoryToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValueToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy)
		if model.Category != nil {
			categoryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model.Category)
			if err != nil {
				return modelMap, err
			}
			modelMap["category"] = []map[string]interface{}{categoryMap}
		}
		if model.Value != nil {
			valueMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model.Value)
			if err != nil {
				return modelMap, err
			}
			modelMap["value"] = []map[string]interface{}{valueMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategoryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Category != nil {
		categoryMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model.Category)
		if err != nil {
			return modelMap, err
		}
		modelMap["category"] = []map[string]interface{}{categoryMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValueToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		valueMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model.Value)
		if err != nil {
			return modelMap, err
		}
		modelMap["value"] = []map[string]interface{}{valueMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model *logsv0.ApisDashboardsV1AstWidgetsMarkdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["markdown_text"] = *model.MarkdownText
	if model.TooltipText != nil {
		modelMap["tooltip_text"] = *model.TooltipText
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueLineChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LineChart != nil {
		lineChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model.LineChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["line_chart"] = []map[string]interface{}{lineChartMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueDataTableToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DataTable != nil {
		dataTableMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model.DataTable)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_table"] = []map[string]interface{}{dataTableMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueGaugeToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Gauge != nil {
		gaugeMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model.Gauge)
		if err != nil {
			return modelMap, err
		}
		modelMap["gauge"] = []map[string]interface{}{gaugeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValuePieChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PieChart != nil {
		pieChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model.PieChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["pie_chart"] = []map[string]interface{}{pieChartMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BarChart != nil {
		barChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model.BarChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["bar_chart"] = []map[string]interface{}{barChartMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.HorizontalBarChart != nil {
		horizontalBarChartMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model.HorizontalBarChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["horizontal_bar_chart"] = []map[string]interface{}{horizontalBarChartMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueMarkdownToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Markdown != nil {
		markdownMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model.Markdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["markdown"] = []map[string]interface{}{markdownMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstVariableToMap(model *logsv0.ApisDashboardsV1AstVariable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	definitionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["display_name"] = *model.DisplayName
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionToMap(model logsv0.ApisDashboardsV1AstVariableDefinitionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionValueMultiSelectToMap(model.(*logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstVariableDefinition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstVariableDefinition)
		if model.MultiSelect != nil {
			multiSelectMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model.MultiSelect)
			if err != nil {
				return modelMap, err
			}
			modelMap["multi_select"] = []map[string]interface{}{multiSelectMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstVariableDefinitionIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model *logsv0.ApisDashboardsV1AstMultiSelect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	selectionMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionToMap(model.Selection)
	if err != nil {
		return modelMap, err
	}
	modelMap["selection"] = []map[string]interface{}{selectionMap}
	modelMap["values_order_direction"] = *model.ValuesOrderDirection
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceToMap(model logsv0.ApisDashboardsV1AstMultiSelectSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueLogsPathToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueMetricLabelToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueConstantListToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstMultiSelectSource)
		if model.LogsPath != nil {
			logsPathMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model.LogsPath)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs_path"] = []map[string]interface{}{logsPathMap}
		}
		if model.MetricLabel != nil {
			metricLabelMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model.MetricLabel)
			if err != nil {
				return modelMap, err
			}
			modelMap["metric_label"] = []map[string]interface{}{metricLabelMap}
		}
		if model.ConstantList != nil {
			constantListMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model.ConstantList)
			if err != nil {
				return modelMap, err
			}
			modelMap["constant_list"] = []map[string]interface{}{constantListMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstMultiSelectSourceIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectLogsPathSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectMetricLabelSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["metric_name"] = *model.MetricName
	modelMap["label"] = *model.Label
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectConstantListSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["values"] = model.Values
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueLogsPathToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LogsPath != nil {
		logsPathMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model.LogsPath)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs_path"] = []map[string]interface{}{logsPathMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueMetricLabelToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MetricLabel != nil {
		metricLabelMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model.MetricLabel)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_label"] = []map[string]interface{}{metricLabelMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueConstantListToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConstantList != nil {
		constantListMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model.ConstantList)
		if err != nil {
			return modelMap, err
		}
		modelMap["constant_list"] = []map[string]interface{}{constantListMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionToMap(model logsv0.ApisDashboardsV1AstMultiSelectSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueAllToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueListToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelection)
		if model.All != nil {
			allMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model.All)
			if err != nil {
				return modelMap, err
			}
			modelMap["all"] = []map[string]interface{}{allMap}
		}
		if model.List != nil {
			listMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model.List)
			if err != nil {
				return modelMap, err
			}
			modelMap["list"] = []map[string]interface{}{listMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstMultiSelectSelectionIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueAllToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.All != nil {
		allMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model.All)
		if err != nil {
			return modelMap, err
		}
		modelMap["all"] = []map[string]interface{}{allMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueListToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionValueMultiSelectToMap(model *logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MultiSelect != nil {
		multiSelectMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model.MultiSelect)
		if err != nil {
			return modelMap, err
		}
		modelMap["multi_select"] = []map[string]interface{}{multiSelectMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFilterToMap(model *logsv0.ApisDashboardsV1AstFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Source != nil {
		sourceMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(model.Source)
		if err != nil {
			return modelMap, err
		}
		modelMap["source"] = []map[string]interface{}{sourceMap}
	}
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.Collapsed != nil {
		modelMap["collapsed"] = *model.Collapsed
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationToMap(model *logsv0.ApisDashboardsV1AstAnnotation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = model.Href.String()
	}
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	modelMap["enabled"] = *model.Enabled
	sourceMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceToMap(model logsv0.ApisDashboardsV1AstAnnotationSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstAnnotationSource)
		if model.Metrics != nil {
			metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Logs != nil {
			logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstAnnotationSourceIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Strategy != nil {
		strategyMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStrategyToMap(model.Strategy)
		if err != nil {
			return modelMap, err
		}
		modelMap["strategy"] = []map[string]interface{}{strategyMap}
	}
	if model.MessageTemplate != nil {
		modelMap["message_template"] = *model.MessageTemplate
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonPromQlQueryToMap(model *logsv0.ApisDashboardsV1CommonPromQlQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStrategyToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStrategy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StartTimeMetric != nil {
		startTimeMetricMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmptyToMap(model.StartTimeMetric)
		if err != nil {
			return modelMap, err
		}
		modelMap["start_time_metric"] = []map[string]interface{}{startTimeMetricMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmptyToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	luceneQueryMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonLuceneQueryToMap(model.LuceneQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	strategyMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyToMap(model.Strategy)
	if err != nil {
		return modelMap, err
	}
	modelMap["strategy"] = []map[string]interface{}{strategyMap}
	if model.MessageTemplate != nil {
		modelMap["message_template"] = *model.MessageTemplate
	}
	if model.LabelFields != nil {
		labelFields := []map[string]interface{}{}
		for _, labelFieldsItem := range model.LabelFields {
			labelFieldsItemMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&labelFieldsItem)
			if err != nil {
				return modelMap, err
			}
			labelFields = append(labelFields, labelFieldsItemMap)
		}
		modelMap["label_fields"] = labelFields
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonLuceneQueryToMap(model *logsv0.ApisDashboardsV1CommonLuceneQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyToMap(model logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstantToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueRangeToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration); ok {
		return ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueDurationToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategy)
		if model.Instant != nil {
			instantMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model.Instant)
			if err != nil {
				return modelMap, err
			}
			modelMap["instant"] = []map[string]interface{}{instantMap}
		}
		if model.Range != nil {
			rangeVarMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model.Range)
			if err != nil {
				return modelMap, err
			}
			modelMap["range"] = []map[string]interface{}{rangeVarMap}
		}
		if model.Duration != nil {
			durationMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model.Duration)
			if err != nil {
				return modelMap, err
			}
			modelMap["duration"] = []map[string]interface{}{durationMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyIntf subtype encountered")
	}
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyInstant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	timestampFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.TimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["timestamp_field"] = []map[string]interface{}{timestampFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startTimestampFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StartTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_timestamp_field"] = []map[string]interface{}{startTimestampFieldMap}
	endTimestampFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.EndTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["end_timestamp_field"] = []map[string]interface{}{endTimestampFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyDuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startTimestampFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StartTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_timestamp_field"] = []map[string]interface{}{startTimestampFieldMap}
	durationFieldMap, err := ResourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.DurationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["duration_field"] = []map[string]interface{}{durationFieldMap}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstantToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instant != nil {
		instantMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model.Instant)
		if err != nil {
			return modelMap, err
		}
		modelMap["instant"] = []map[string]interface{}{instantMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueRangeToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Range != nil {
		rangeVarMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model.Range)
		if err != nil {
			return modelMap, err
		}
		modelMap["range"] = []map[string]interface{}{rangeVarMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueDurationToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Duration != nil {
		durationMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model.Duration)
		if err != nil {
			return modelMap, err
		}
		modelMap["duration"] = []map[string]interface{}{durationMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueMetricsToMap(model *logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueLogsToMap(model *logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := ResourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1CommonTimeFrameToMap(model *logsv0.ApisDashboardsV1CommonTimeFrame) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.From != nil {
		modelMap["from"] = model.From.String()
	}
	if model.To != nil {
		modelMap["to"] = model.To.String()
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstFolderPathToMap(model *logsv0.ApisDashboardsV1AstFolderPath) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Segments != nil {
		modelMap["segments"] = model.Segments
	}
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshOffEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshOffEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}
