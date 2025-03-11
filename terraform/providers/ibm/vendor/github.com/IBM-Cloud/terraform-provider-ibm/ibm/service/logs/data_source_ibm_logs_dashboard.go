// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsDashboardRead,

		Schema: map[string]*schema.Schema{
			"dashboard_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dashboard.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the dashboard.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the dashboard.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Brief description or summary of the dashboard's purpose or content.",
			},
			"layout": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Layout configuration for the dashboard's visual elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sections": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The sections of the layout.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier of the section within the layout.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Unique identifier of the folder containing the dashboard.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The UUID value.",
												},
											},
										},
									},
									"rows": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The rows of the section.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier of the row within the layout.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Unique identifier of the folder containing the dashboard.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The UUID value.",
															},
														},
													},
												},
												"appearance": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The appearance of the row, such as height.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"height": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The height of the row.",
															},
														},
													},
												},
												"widgets": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The widgets of the row.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Widget identifier within the dashboard.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Unique identifier of the folder containing the dashboard.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The UUID value.",
																		},
																	},
																},
															},
															"title": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Widget title.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Widget description.",
															},
															"definition": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Widget definition, contains the widget type and its configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"line_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Line chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"legend": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Legend configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Whether to show the legend or not.",
																								},
																								"columns": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The columns to show in the legend.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"group_by_query": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Whether to group by the query or not.",
																								},
																							},
																						},
																					},
																					"tooltip": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Tooltip configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"show_labels": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Whether to show labels in the tooltip.",
																								},
																								"type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Tooltip type.",
																								},
																							},
																						},
																					},
																					"query_definitions": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Query definitions.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"id": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Unique identifier of the query within the widget.",
																								},
																								"query": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Data source specific query, defines from where and how to fetch the data.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"logs": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Logs specific query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"lucene_query": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"value": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"group_by": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Group by fields (deprecated).",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"aggregations": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Aggregations.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"count": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Count the number of entries.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{},
																																		},
																																	},
																																	"count_distinct": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Count the number of distinct values of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Sum values of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Calculate average value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Calculate minimum value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Calculate maximum value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Calculate percentile value of log field.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"percent": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Value in range (0, 100].",
																																				},
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Group by fields.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																												Computed:    true,
																												Description: "Metrics specific query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"promql_query": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "PromQL query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"value": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"filters": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
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
																												Computed:    true,
																												Description: "Dataprime language based query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"dataprime_query": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Dataprime query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"text": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "The query string.",
																																	},
																																},
																															},
																														},
																														"filters": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Filters to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"logs": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Extra filtering on top of the Lucene query.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"operator": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Operator to use for filtering the logs.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Selection criteria for the equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"all": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection of all values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{},
																																														},
																																													},
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Computed:    true,
																																																	Description: "List of values for the selection.",
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
																																									},
																																								},
																																							},
																																							"not_equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Selection criteria for the non-equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Computed:    true,
																																																	Description: "List of values for the selection.",
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
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																				"observation_field": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Field to count distinct values of.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"keypath": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Path within the dataset scope.",
																																								Elem: &schema.Schema{
																																									Type: schema.TypeString,
																																								},
																																							},
																																							"scope": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
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
																																		Computed:    true,
																																		Description: "Filtering to be applied to query results.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"label": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Label associated with the metric.",
																																				},
																																				"operator": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Operator to use for filtering the logs.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Selection criteria for the equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"all": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection of all values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{},
																																														},
																																													},
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Computed:    true,
																																																	Description: "List of values for the selection.",
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
																																									},
																																								},
																																							},
																																							"not_equals": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"selection": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Selection criteria for the non-equality comparison.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"list": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "Represents a selection from a list of values.",
																																														Elem: &schema.Resource{
																																															Schema: map[string]*schema.Schema{
																																																"values": &schema.Schema{
																																																	Type:        schema.TypeList,
																																																	Computed:    true,
																																																	Description: "List of values for the selection.",
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
																																									},
																																								},
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
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
																									Computed:    true,
																									Description: "Template for series name in legend and tooltip.",
																								},
																								"series_count_limit": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Maximum number of series to display.",
																								},
																								"unit": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Unit of the data.",
																								},
																								"scale_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Scale type.",
																								},
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Query name.",
																								},
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Whether data for this query should be visible on the chart.",
																								},
																								"color_scheme": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Color scheme for the series.",
																								},
																								"resolution": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Resolution of the data.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"interval": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Interval between data points.",
																											},
																											"buckets_presented": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Maximum number of data points to fetch.",
																											},
																										},
																									},
																								},
																								"data_mode_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Data mode type.",
																								},
																							},
																						},
																					},
																					"stacked_line": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Stacked lines.",
																					},
																				},
																			},
																		},
																		"data_table": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Data table widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filtering on top of the Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																												Computed:    true,
																												Description: "Grouping and aggregation.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"aggregations": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Aggregations.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"id": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Aggregation identifier, must be unique within grouping configuration.",
																																	},
																																	"name": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Aggregation name, used as column name.",
																																	},
																																	"is_visible": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Whether the aggregation is visible.",
																																	},
																																	"aggregation": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Aggregations.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"count": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Count the number of entries.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{},
																																					},
																																				},
																																				"count_distinct": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Count the number of distinct values of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																																					Computed:    true,
																																					Description: "Sum values of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																																					Computed:    true,
																																					Description: "Calculate average value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																																					Computed:    true,
																																					Description: "Calculate minimum value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																																					Computed:    true,
																																					Description: "Calculate maximum value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																																					Computed:    true,
																																					Description: "Calculate percentile value of log field.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"percent": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Value in range (0, 100].",
																																							},
																																							"observation_field": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Field to count distinct values of.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"keypath": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Path within the dataset scope.",
																																											Elem: &schema.Schema{
																																												Type: schema.TypeString,
																																											},
																																										},
																																										"scope": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
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
																															Computed:    true,
																															Description: "Fields to group by.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																									Computed:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filtering on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
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
																									Computed:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filtering on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
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
																						Computed:    true,
																						Description: "Number of results per page.",
																					},
																					"row_style": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Display style for table rows.",
																					},
																					"columns": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Columns to display, their order and width.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"field": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "References a field in result set. In case of aggregation, it references the aggregation identifier.",
																								},
																								"width": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Column width.",
																								},
																							},
																						},
																					},
																					"order_by": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Column used for ordering the results.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"field": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The field to order by.",
																								},
																								"order_direction": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The direction of the order: ascending or descending.",
																								},
																							},
																						},
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"gauge": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Gauge widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"metrics": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Aggregation. When AGGREGATION_UNSPECIFIED is selected, widget uses instant query. Otherwise, it uses range query.",
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filters applied on top of PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
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
																									Computed:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"logs_aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																												Computed:    true,
																												Description: "Extra filters applied on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																									Computed:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filters applied on top of Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
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
																						Computed:    true,
																						Description: "Minimum value of the gauge.",
																					},
																					"max": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Maximum value of the gauge.",
																					},
																					"show_inner_arc": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Show inner arc (styling).",
																					},
																					"show_outer_arc": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Show outer arc (styling).",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Query result value interpretation.",
																					},
																					"thresholds": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Thresholds for the gauge, values at which the gauge changes color.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"from": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Value at which the color should change.",
																								},
																								"color": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Color.",
																								},
																							},
																						},
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Data mode type.",
																					},
																					"threshold_by": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "What threshold color should be applied to: value or background.",
																					},
																				},
																			},
																		},
																		"pie_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Pie chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																												Computed:    true,
																												Description: "Extra filters on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
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
																									Computed:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filters on top of PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Field to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filters on top of Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
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
																						Computed:    true,
																						Description: "Maximum number of slices to display in the chart.",
																					},
																					"min_slice_percentage": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Minimum percentage of a slice to be displayed.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_stack": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Maximum number of slices per stack.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Template for stack labels.",
																								},
																							},
																						},
																					},
																					"label_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Label settings.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"label_source": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Source of the label.",
																								},
																								"is_visible": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Controls whether to show the label.",
																								},
																								"show_name": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Controls whether to show the name.",
																								},
																								"show_value": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Controls whether to show the value.",
																								},
																								"show_percentage": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Controls whether to show the percentage.",
																								},
																							},
																						},
																					},
																					"show_legend": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Controls whether to show the legend.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Template for group labels.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Unit of the data.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Color scheme name.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"bar_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Bar chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																												Computed:    true,
																												Description: "Extra filter on top of Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																												Computed:    true,
																												Description: "Fiel to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
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
																									Computed:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filter on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																												Computed:    true,
																												Description: "Labels to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Label to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Query based on Dataprime language.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filter on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
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
																						Computed:    true,
																						Description: "Maximum number of bars to present in the chart.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Template for bar labels.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_bar": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Maximum number of slices per bar.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Template for stack slice label.",
																								},
																							},
																						},
																					},
																					"scale_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Scale type.",
																					},
																					"colors_by": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Coloring mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"stack": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Each stack will have the same color across all groups.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"group_by": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Each group will have different color and stack color will be derived from group color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"aggregation": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
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
																						Computed:    true,
																						Description: "X axis mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"value": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Categorical axis.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"time": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Time based axis.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"interval": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Time interval.",
																											},
																											"buckets_presented": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
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
																						Computed:    true,
																						Description: "Unit of the data.",
																					},
																					"sort_by": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Sorting mode.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Supported vaues: classic, severity, cold, negative, green, red, blue.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"horizontal_bar_chart": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Horizontal bar chart widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"query": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Data source specific query, defines from where and how to fetch the data.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"logs": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Logs specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"lucene_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"aggregation": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Aggregations.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"count": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of entries.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{},
																															},
																														},
																														"count_distinct": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Count the number of distinct values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Sum values of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate average value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate minimum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate maximum value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Calculate percentile value of log field.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"percent": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Value in range (0, 100].",
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																												Computed:    true,
																												Description: "Extra filter on top of the Lucene query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"observation_field": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Field to count distinct values of.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"keypath": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Path within the dataset scope.",
																																		Elem: &schema.Schema{
																																			Type: schema.TypeString,
																																		},
																																	},
																																	"scope": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Scope of the dataset.",
																														},
																													},
																												},
																											},
																											"stacked_group_name_field": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Field to count distinct values of.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"keypath": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Path within the dataset scope.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"scope": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
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
																									Computed:    true,
																									Description: "Metrics specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"promql_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"value": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filter on top of the PromQL query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"label": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Label associated with the metric.",
																														},
																														"operator": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Operator to use for filtering the logs.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"all": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection of all values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{},
																																								},
																																							},
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																																			},
																																		},
																																	},
																																	"not_equals": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Non-equality comparison.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"selection": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Selection criteria for the non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"list": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Represents a selection from a list of values.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"values": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "List of values for the selection.",
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
																												Computed:    true,
																												Description: "Labels to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Label to stack by.",
																											},
																										},
																									},
																								},
																								"dataprime": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Dataprime specific query.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"dataprime_query": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"text": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The query string.",
																														},
																													},
																												},
																											},
																											"filters": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Extra filter on top of the Dataprime query.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"logs": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Extra filtering on top of the Lucene query.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"observation_field": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Field to count distinct values of.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"keypath": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Path within the dataset scope.",
																																					Elem: &schema.Schema{
																																						Type: schema.TypeString,
																																					},
																																				},
																																				"scope": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
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
																															Computed:    true,
																															Description: "Filtering to be applied to query results.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"label": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Label associated with the metric.",
																																	},
																																	"operator": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Operator to use for filtering the logs.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"all": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection of all values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{},
																																											},
																																										},
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
																																				"not_equals": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Non-equality comparison.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"selection": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Selection criteria for the non-equality comparison.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"list": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Represents a selection from a list of values.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"values": &schema.Schema{
																																														Type:        schema.TypeList,
																																														Computed:    true,
																																														Description: "List of values for the selection.",
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
																																						},
																																					},
																																				},
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
																												Computed:    true,
																												Description: "Fields to group by.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"stacked_group_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
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
																						Computed:    true,
																						Description: "Maximum number of bars to display in the chart.",
																					},
																					"group_name_template": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Template for bar labels.",
																					},
																					"stack_definition": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Stack definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_slices_per_bar": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Maximum number of slices per bar.",
																								},
																								"stack_name_template": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Template for stack slice label.",
																								},
																							},
																						},
																					},
																					"scale_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Scale type.",
																					},
																					"colors_by": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Coloring mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"stack": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Each stack will have the same color across all groups.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"group_by": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Each group will have different color and stack color will be derived from group color.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"aggregation": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
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
																						Computed:    true,
																						Description: "Unit of the data.",
																					},
																					"display_on_bar": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Whether to display values on the bars.",
																					},
																					"y_axis_view_by": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Y-axis view mode.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"category": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "View by category.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{},
																									},
																								},
																								"value": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
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
																						Computed:    true,
																						Description: "Sorting mode.",
																					},
																					"color_scheme": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Color scheme name.",
																					},
																					"data_mode_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Data mode type.",
																					},
																				},
																			},
																		},
																		"markdown": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Markdown widget.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"markdown_text": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Markdown text to render.",
																					},
																					"tooltip_text": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
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
																Computed:    true,
																Description: "Creation timestamp.",
															},
															"updated_at": &schema.Schema{
																Type:        schema.TypeString,
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
				Computed:    true,
				Description: "List of variables that can be used within the dashboard for dynamic content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable which can be used in templates.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Definition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"multi_select": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Multi-select value.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Variable value source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"logs_path": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Unique values for a given logs path.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"observation_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
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
																Computed:    true,
																Description: "Unique values for a given metric label.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metric name to source unique values from.",
																		},
																		"label": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Metric label to source unique values from.",
																		},
																	},
																},
															},
															"constant_list": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of constant values.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"values": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of constant values.",
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
												"selection": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "State of what is currently selected.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"all": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "All values are selected, usually translated to wildcard (*).",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"list": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specific values are selected.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"values": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Selected values.",
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
												"values_order_direction": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
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
							Computed:    true,
							Description: "Name used in variable UI.",
						},
					},
				},
			},
			"filters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of filters that can be applied to the dashboard's data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Filters to be applied to query results.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logs": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Extra filtering on top of the Lucene query.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Operator to use for filtering the logs.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"equals": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Selection criteria for the equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"all": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection of all values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{},
																						},
																					},
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "List of values for the selection.",
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
																	},
																},
															},
															"not_equals": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Non-equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Selection criteria for the non-equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "List of values for the selection.",
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
																	},
																},
															},
														},
													},
												},
												"observation_field": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Field to count distinct values of.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keypath": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Path within the dataset scope.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"scope": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
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
										Computed:    true,
										Description: "Filtering to be applied to query results.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Label associated with the metric.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Operator to use for filtering the logs.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"equals": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Selection criteria for the equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"all": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection of all values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{},
																						},
																					},
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "List of values for the selection.",
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
																	},
																},
															},
															"not_equals": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Non-equality comparison.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"selection": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Selection criteria for the non-equality comparison.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"list": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Represents a selection from a list of values.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"values": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "List of values for the selection.",
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
																	},
																},
															},
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
							Computed:    true,
							Description: "Indicates if the filter is currently enabled or not.",
						},
						"collapsed": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the filter's UI representation should be collapsed or expanded.",
						},
					},
				},
			},
			"annotations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of annotations that can be applied to the dashboard's visual elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier within the dashboard.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier within the dashboard.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the annotation.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the annotation is enabled.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Source of the annotation events.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Metrics source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"promql_query": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "PromQL query.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The PromQL query string.",
															},
														},
													},
												},
												"strategy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Strategy for turning metrics data into annotations.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_time_metric": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
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
													Computed:    true,
													Description: "Template for the annotation message.",
												},
												"labels": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Labels to display in the annotation.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"logs": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Logs source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lucene_query": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Lucene query.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The Lucene query string.",
															},
														},
													},
												},
												"strategy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Strategy for turning logs data into annotations.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"instant": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Event timestamp is extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
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
																Computed:    true,
																Description: "Event start and end timestamps are extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																		"end_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
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
																Computed:    true,
																Description: "Event start timestamp and duration are extracted from the log entry.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"start_timestamp_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Scope of the dataset.",
																					},
																				},
																			},
																		},
																		"duration_field": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Field to count distinct values of.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"keypath": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Path within the dataset scope.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"scope": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
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
													Computed:    true,
													Description: "Template for the annotation message.",
												},
												"label_fields": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Labels to display in the annotation.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keypath": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Path within the dataset scope.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"scope": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
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
				Computed:    true,
				Description: "Absolute time frame specifying a fixed start and end time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "From is the start of the time frame.",
						},
						"to": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "To is the end of the time frame.",
						},
					},
				},
			},
			"relative_time_frame": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Relative time frame specifying a duration from the current time.",
			},
			"folder_id": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unique identifier of the folder containing the dashboard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID value.",
						},
					},
				},
			},
			"folder_path": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Path of the folder containing the dashboard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The segments of the folder path.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"false": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Auto refresh interval is set to off.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"two_minutes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Auto refresh interval is set to two minutes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"five_minutes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Auto refresh interval is set to five minutes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
		},
	}
}

func dataSourceIbmLogsDashboardRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getDashboardOptions := &logsv0.GetDashboardOptions{}

	getDashboardOptions.SetDashboardID(d.Get("dashboard_id").(string))

	dashboardIntf, _, err := logsClient.GetDashboardWithContext(context, getDashboardOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDashboardWithContext failed: %s", err.Error()), "(Data) ibm_logs_dashboard", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	dashboard := dashboardIntf.(*logsv0.Dashboard)

	d.SetId(fmt.Sprintf("%s", *getDashboardOptions.DashboardID))

	if err = d.Set("href", dashboard.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", dashboard.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("description", dashboard.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	layout := []map[string]interface{}{}
	if dashboard.Layout != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstLayoutToMap(dashboard.Layout)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		layout = append(layout, modelMap)
	}
	if err = d.Set("layout", layout); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting layout: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	variables := []map[string]interface{}{}
	if dashboard.Variables != nil {
		for _, modelItem := range dashboard.Variables {
			modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstVariableToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
				return tfErr.GetDiag()
			}
			variables = append(variables, modelMap)
		}
	}
	if err = d.Set("variables", variables); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting variables: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	filters := []map[string]interface{}{}
	if dashboard.Filters != nil {
		for _, modelItem := range dashboard.Filters {
			modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
				return tfErr.GetDiag()
			}
			filters = append(filters, modelMap)
		}
	}
	if err = d.Set("filters", filters); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting filters: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	annotations := []map[string]interface{}{}
	if dashboard.Annotations != nil {
		for _, modelItem := range dashboard.Annotations {
			modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
				return tfErr.GetDiag()
			}
			annotations = append(annotations, modelMap)
		}
	}
	if err = d.Set("annotations", annotations); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting annotations: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	absoluteTimeFrame := []map[string]interface{}{}
	if dashboard.AbsoluteTimeFrame != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonTimeFrameToMap(dashboard.AbsoluteTimeFrame)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		absoluteTimeFrame = append(absoluteTimeFrame, modelMap)
	}
	if err = d.Set("absolute_time_frame", absoluteTimeFrame); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting absolute_time_frame: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("relative_time_frame", dashboard.RelativeTimeFrame); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting relative_time_frame: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	folderID := []map[string]interface{}{}
	if dashboard.FolderID != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1UUIDToMap(dashboard.FolderID)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		folderID = append(folderID, modelMap)
	}
	if err = d.Set("folder_id", folderID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting folder_id: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	folderPath := []map[string]interface{}{}
	if dashboard.FolderPath != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFolderPathToMap(dashboard.FolderPath)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		folderPath = append(folderPath, modelMap)
	}
	if err = d.Set("folder_path", folderPath); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting folder_path: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	falseVar := []map[string]interface{}{}
	if dashboard.False != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshOffEmptyToMap(dashboard.False)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		falseVar = append(falseVar, modelMap)
	}
	if err = d.Set("false", falseVar); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting false: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	twoMinutes := []map[string]interface{}{}
	if dashboard.TwoMinutes != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmptyToMap(dashboard.TwoMinutes)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		twoMinutes = append(twoMinutes, modelMap)
	}
	if err = d.Set("two_minutes", twoMinutes); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting two_minutes: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	fiveMinutes := []map[string]interface{}{}
	if dashboard.FiveMinutes != nil {
		modelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmptyToMap(dashboard.FiveMinutes)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard", "read")
			return tfErr.GetDiag()
		}
		fiveMinutes = append(fiveMinutes, modelMap)
	}
	if err = d.Set("five_minutes", fiveMinutes); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting five_minutes: %s", err), "(Data) ibm_logs_dashboard", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstLayoutToMap(model *logsv0.ApisDashboardsV1AstLayout) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Sections != nil {
		sections := []map[string]interface{}{}
		for _, sectionsItem := range model.Sections {
			sectionsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstSectionToMap(&sectionsItem)
			if err != nil {
				return modelMap, err
			}
			sections = append(sections, sectionsItemMap)
		}
		modelMap["sections"] = sections
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstSectionToMap(model *logsv0.ApisDashboardsV1AstSection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := DataSourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	if model.Rows != nil {
		rows := []map[string]interface{}{}
		for _, rowsItem := range model.Rows {
			rowsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstRowToMap(&rowsItem)
			if err != nil {
				return modelMap, err
			}
			rows = append(rows, rowsItemMap)
		}
		modelMap["rows"] = rows
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model *logsv0.ApisDashboardsV1UUID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["value"] = model.Value.String()
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstRowToMap(model *logsv0.ApisDashboardsV1AstRow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := DataSourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	appearanceMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstRowAppearanceToMap(model.Appearance)
	if err != nil {
		return modelMap, err
	}
	modelMap["appearance"] = []map[string]interface{}{appearanceMap}
	if model.Widgets != nil {
		widgets := []map[string]interface{}{}
		for _, widgetsItem := range model.Widgets {
			widgetsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetToMap(&widgetsItem)
			if err != nil {
				return modelMap, err
			}
			widgets = append(widgets, widgetsItemMap)
		}
		modelMap["widgets"] = widgets
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstRowAppearanceToMap(model *logsv0.ApisDashboardsV1AstRowAppearance) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["height"] = flex.IntValue(model.Height)
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetToMap(model *logsv0.ApisDashboardsV1AstWidget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	idMap, err := DataSourceIbmLogsDashboardApisDashboardsV1UUIDToMap(model.ID)
	if err != nil {
		return modelMap, err
	}
	modelMap["id"] = []map[string]interface{}{idMap}
	modelMap["title"] = *model.Title
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	definitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionToMap(model.Definition)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionToMap(model logsv0.ApisDashboardsV1AstWidgetDefinitionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueLineChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueDataTableToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueGaugeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValuePieChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueBarChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChartToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueMarkdownToMap(model.(*logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetDefinition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetDefinition)
		if model.LineChart != nil {
			lineChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model.LineChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["line_chart"] = []map[string]interface{}{lineChartMap}
		}
		if model.DataTable != nil {
			dataTableMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model.DataTable)
			if err != nil {
				return modelMap, err
			}
			modelMap["data_table"] = []map[string]interface{}{dataTableMap}
		}
		if model.Gauge != nil {
			gaugeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model.Gauge)
			if err != nil {
				return modelMap, err
			}
			modelMap["gauge"] = []map[string]interface{}{gaugeMap}
		}
		if model.PieChart != nil {
			pieChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model.PieChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["pie_chart"] = []map[string]interface{}{pieChartMap}
		}
		if model.BarChart != nil {
			barChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model.BarChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["bar_chart"] = []map[string]interface{}{barChartMap}
		}
		if model.HorizontalBarChart != nil {
			horizontalBarChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model.HorizontalBarChart)
			if err != nil {
				return modelMap, err
			}
			modelMap["horizontal_bar_chart"] = []map[string]interface{}{horizontalBarChartMap}
		}
		if model.Markdown != nil {
			markdownMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model.Markdown)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	legendMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLegendToMap(model.Legend)
	if err != nil {
		return modelMap, err
	}
	modelMap["legend"] = []map[string]interface{}{legendMap}
	tooltipMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartTooltipToMap(model.Tooltip)
	if err != nil {
		return modelMap, err
	}
	modelMap["tooltip"] = []map[string]interface{}{tooltipMap}
	if model.QueryDefinitions != nil {
		queryDefinitions := []map[string]interface{}{}
		for _, queryDefinitionsItem := range model.QueryDefinitions {
			queryDefinitionsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryDefinitionToMap(&queryDefinitionsItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLegendToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonLegend) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["is_visible"] = *model.IsVisible
	if model.Columns != nil {
		modelMap["columns"] = model.Columns
	}
	modelMap["group_by_query"] = *model.GroupByQuery
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartTooltipToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartTooltip) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ShowLabels != nil {
		modelMap["show_labels"] = *model.ShowLabels
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryToMap(model.Query)
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
	resolutionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartResolutionToMap(model.Resolution)
	if err != nil {
		return modelMap, err
	}
	modelMap["resolution"] = []map[string]interface{}{resolutionMap}
	if model.DataModeType != nil {
		modelMap["data_mode_type"] = *model.DataModeType
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsLineChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsLineChartQuery)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
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
			aggregationsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(aggregationsItem)
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
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
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
			groupBysItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupBysItem)
			if err != nil {
				return modelMap, err
			}
			groupBys = append(groupBys, groupBysItemMap)
		}
		modelMap["group_bys"] = groupBys
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonLuceneQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model logsv0.ApisDashboardsV1CommonLogsAggregationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCount); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCount))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountDistinctToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueSum); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueSumToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueSum))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueAverageToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMin); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMinToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMin))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMax); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMaxToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValueMax))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValuePercentileToMap(model.(*logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1CommonLogsAggregation); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1CommonLogsAggregation)
		if model.Count != nil {
			countMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model.Count)
			if err != nil {
				return modelMap, err
			}
			modelMap["count"] = []map[string]interface{}{countMap}
		}
		if model.CountDistinct != nil {
			countDistinctMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model.CountDistinct)
			if err != nil {
				return modelMap, err
			}
			modelMap["count_distinct"] = []map[string]interface{}{countDistinctMap}
		}
		if model.Sum != nil {
			sumMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model.Sum)
			if err != nil {
				return modelMap, err
			}
			modelMap["sum"] = []map[string]interface{}{sumMap}
		}
		if model.Average != nil {
			averageMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model.Average)
			if err != nil {
				return modelMap, err
			}
			modelMap["average"] = []map[string]interface{}{averageMap}
		}
		if model.Min != nil {
			minMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model.Min)
			if err != nil {
				return modelMap, err
			}
			modelMap["min"] = []map[string]interface{}{minMap}
		}
		if model.Max != nil {
			maxMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model.Max)
			if err != nil {
				return modelMap, err
			}
			modelMap["max"] = []map[string]interface{}{maxMap}
		}
		if model.Percentile != nil {
			percentileMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model.Percentile)
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

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationCountEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationCountDistinct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model *logsv0.ApisDashboardsV1CommonObservationField) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Keypath != nil {
		modelMap["keypath"] = model.Keypath
	}
	if model.Scope != nil {
		modelMap["scope"] = *model.Scope
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationSum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationAverage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationMin) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationMax) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationPercentile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["percent"] = *model.Percent
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueCount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Count != nil {
		countMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountEmptyToMap(model.Count)
		if err != nil {
			return modelMap, err
		}
		modelMap["count"] = []map[string]interface{}{countMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueCountDistinctToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueCountDistinct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CountDistinct != nil {
		countDistinctMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationCountDistinctToMap(model.CountDistinct)
		if err != nil {
			return modelMap, err
		}
		modelMap["count_distinct"] = []map[string]interface{}{countDistinctMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueSumToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueSum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Sum != nil {
		sumMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationSumToMap(model.Sum)
		if err != nil {
			return modelMap, err
		}
		modelMap["sum"] = []map[string]interface{}{sumMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueAverageToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueAverage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Average != nil {
		averageMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationAverageToMap(model.Average)
		if err != nil {
			return modelMap, err
		}
		modelMap["average"] = []map[string]interface{}{averageMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMinToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueMin) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Min != nil {
		minMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMinToMap(model.Min)
		if err != nil {
			return modelMap, err
		}
		modelMap["min"] = []map[string]interface{}{minMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValueMaxToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValueMax) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Max != nil {
		maxMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationMaxToMap(model.Max)
		if err != nil {
			return modelMap, err
		}
		modelMap["max"] = []map[string]interface{}{maxMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationValuePercentileToMap(model *logsv0.ApisDashboardsV1CommonLogsAggregationValuePercentile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Percentile != nil {
		percentileMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationPercentileToMap(model.Percentile)
		if err != nil {
			return modelMap, err
		}
		modelMap["percentile"] = []map[string]interface{}{percentileMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model *logsv0.ApisDashboardsV1AstFilterLogsFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Operator != nil {
		operatorMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model.Operator)
		if err != nil {
			return modelMap, err
		}
		modelMap["operator"] = []map[string]interface{}{operatorMap}
	}
	if model.ObservationField != nil {
		observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
		if err != nil {
			return modelMap, err
		}
		modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model logsv0.ApisDashboardsV1AstFilterOperatorIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueEquals); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueEqualsToMap(model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueEquals))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueNotEqualsToMap(model.(*logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterOperator); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterOperator)
		if model.Equals != nil {
			equalsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model.Equals)
			if err != nil {
				return modelMap, err
			}
			modelMap["equals"] = []map[string]interface{}{equalsMap}
		}
		if model.NotEquals != nil {
			notEqualsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model.NotEquals)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Selection != nil {
		selectionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionToMap(model.Selection)
		if err != nil {
			return modelMap, err
		}
		modelMap["selection"] = []map[string]interface{}{selectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionToMap(model logsv0.ApisDashboardsV1AstFilterEqualsSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueAllToMap(model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueListToMap(model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterEqualsSelection)
		if model.All != nil {
			allMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model.All)
			if err != nil {
				return modelMap, err
			}
			modelMap["all"] = []map[string]interface{}{allMap}
		}
		if model.List != nil {
			listMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model.List)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueAllToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueAll) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.All != nil {
		allMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionAllSelectionEmptyToMap(model.All)
		if err != nil {
			return modelMap, err
		}
		modelMap["all"] = []map[string]interface{}{allMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionValueListToMap(model *logsv0.ApisDashboardsV1AstFilterEqualsSelectionValueList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterNotEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Selection != nil {
		selectionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionToMap(model.Selection)
		if err != nil {
			return modelMap, err
		}
		modelMap["selection"] = []map[string]interface{}{selectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterNotEqualsSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstFilterNotEqualsSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterOperatorValueEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Equals != nil {
		equalsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterEqualsToMap(model.Equals)
		if err != nil {
			return modelMap, err
		}
		modelMap["equals"] = []map[string]interface{}{equalsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorValueNotEqualsToMap(model *logsv0.ApisDashboardsV1AstFilterOperatorValueNotEquals) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotEquals != nil {
		notEqualsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterNotEqualsToMap(model.NotEquals)
		if err != nil {
			return modelMap, err
		}
		modelMap["not_equals"] = []map[string]interface{}{notEqualsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonPromQlQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model *logsv0.ApisDashboardsV1AstFilterMetricsFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	if model.Operator != nil {
		operatorMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterOperatorToMap(model.Operator)
		if err != nil {
			return modelMap, err
		}
		modelMap["operator"] = []map[string]interface{}{operatorMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model *logsv0.ApisDashboardsV1CommonDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Text != nil {
		modelMap["text"] = *model.Text
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(model logsv0.ApisDashboardsV1AstFilterSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSourceValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstFilterSourceValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSourceValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstFilterSourceValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstFilterSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstFilterSource)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model.Metrics)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueLogsToMap(model *logsv0.ApisDashboardsV1AstFilterSourceValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceValueMetricsToMap(model *logsv0.ApisDashboardsV1AstFilterSourceValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartResolutionToMap(model *logsv0.ApisDashboardsV1AstWidgetsLineChartResolution) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.BucketsPresented != nil {
		modelMap["buckets_presented"] = flex.IntValue(model.BucketsPresented)
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["results_per_page"] = flex.IntValue(model.ResultsPerPage)
	modelMap["row_style"] = *model.RowStyle
	if model.Columns != nil {
		columns := []map[string]interface{}{}
		for _, columnsItem := range model.Columns {
			columnsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableColumnToMap(&columnsItem)
			if err != nil {
				return modelMap, err
			}
			columns = append(columns, columnsItemMap)
		}
		modelMap["columns"] = columns
	}
	if model.OrderBy != nil {
		orderByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonOrderingFieldToMap(model.OrderBy)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsDataTableQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsDataTableQuery)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	if model.Grouping != nil {
		groupingMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryGroupingToMap(model.Grouping)
		if err != nil {
			return modelMap, err
		}
		modelMap["grouping"] = []map[string]interface{}{groupingMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryGroupingToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryGrouping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Aggregations != nil {
		aggregations := []map[string]interface{}{}
		for _, aggregationsItem := range model.Aggregations {
			aggregationsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryAggregationToMap(&aggregationsItem)
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
			groupBysItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupBysItem)
			if err != nil {
				return modelMap, err
			}
			groupBys = append(groupBys, groupBysItemMap)
		}
		modelMap["group_bys"] = groupBys
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryAggregationToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableLogsQueryAggregation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	modelMap["is_visible"] = *model.IsVisible
	aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableColumnToMap(model *logsv0.ApisDashboardsV1AstWidgetsDataTableColumn) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["field"] = *model.Field
	if model.Width != nil {
		modelMap["width"] = flex.IntValue(model.Width)
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonOrderingFieldToMap(model *logsv0.ApisDashboardsV1CommonOrderingField) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Field != nil {
		modelMap["field"] = *model.Field
	}
	if model.OrderDirection != nil {
		modelMap["order_direction"] = *model.OrderDirection
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model *logsv0.ApisDashboardsV1AstWidgetsGauge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryToMap(model.Query)
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
			thresholdsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeThresholdToMap(&thresholdsItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsGaugeQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsGaugeQuery)
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	modelMap["aggregation"] = *model.Aggregation
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.LogsAggregation != nil {
		logsAggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.LogsAggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs_aggregation"] = []map[string]interface{}{logsAggregationMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeThresholdToMap(model *logsv0.ApisDashboardsV1AstWidgetsGaugeThreshold) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["from"] = *model.From
	modelMap["color"] = *model.Color
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["max_slices_per_chart"] = flex.IntValue(model.MaxSlicesPerChart)
	modelMap["min_slice_percentage"] = flex.IntValue(model.MinSlicePercentage)
	stackDefinitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartStackDefinitionToMap(model.StackDefinition)
	if err != nil {
		return modelMap, err
	}
	modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	labelDefinitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLabelDefinitionToMap(model.LabelDefinition)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsPieChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsPieChartQuery)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
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
			groupNamesFieldsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerStack != nil {
		modelMap["max_slices_per_stack"] = flex.IntValue(model.MaxSlicesPerStack)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartLabelDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsPieChartLabelDefinition) (map[string]interface{}, error) {
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryToMap(model.Query)
	if err != nil {
		return modelMap, err
	}
	modelMap["query"] = []map[string]interface{}{queryMap}
	modelMap["max_bars_per_chart"] = flex.IntValue(model.MaxBarsPerChart)
	modelMap["group_name_template"] = *model.GroupNameTemplate
	stackDefinitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartStackDefinitionToMap(model.StackDefinition)
	if err != nil {
		return modelMap, err
	}
	modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	modelMap["scale_type"] = *model.ScaleType
	colorsByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model.ColorsBy)
	if err != nil {
		return modelMap, err
	}
	modelMap["colors_by"] = []map[string]interface{}{colorsByMap}
	xAxisMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisToMap(model.XAxis)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsBarChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartQuery)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
	if err != nil {
		return modelMap, err
	}
	modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
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
			groupNamesFieldsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerBar != nil {
		modelMap["max_slices_per_bar"] = flex.IntValue(model.MaxSlicesPerBar)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model logsv0.ApisDashboardsV1AstWidgetsCommonColorsByIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueStackToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueGroupByToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueAggregationToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsBy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsCommonColorsBy)
		if model.Stack != nil {
			stackMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model.Stack)
			if err != nil {
				return modelMap, err
			}
			modelMap["stack"] = []map[string]interface{}{stackMap}
		}
		if model.GroupBy != nil {
			groupByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model.GroupBy)
			if err != nil {
				return modelMap, err
			}
			modelMap["group_by"] = []map[string]interface{}{groupByMap}
		}
		if model.Aggregation != nil {
			aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model.Aggregation)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueStackToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueStack) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Stack != nil {
		stackMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByStackEmptyToMap(model.Stack)
		if err != nil {
			return modelMap, err
		}
		modelMap["stack"] = []map[string]interface{}{stackMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueGroupByToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueGroupBy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GroupBy != nil {
		groupByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByGroupByEmptyToMap(model.GroupBy)
		if err != nil {
			return modelMap, err
		}
		modelMap["group_by"] = []map[string]interface{}{groupByMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByValueAggregationToMap(model *logsv0.ApisDashboardsV1AstWidgetsCommonColorsByValueAggregation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Aggregation != nil {
		aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByColorsByAggregationEmptyToMap(model.Aggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisToMap(model logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeValueToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeTimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxis); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsBarChartXAxis)
		if model.Value != nil {
			valueMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model.Value)
			if err != nil {
				return modelMap, err
			}
			modelMap["value"] = []map[string]interface{}{valueMap}
		}
		if model.Time != nil {
			timeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model.Time)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.BucketsPresented != nil {
		modelMap["buckets_presented"] = flex.IntValue(model.BucketsPresented)
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeValueToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		valueMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByValueEmptyToMap(model.Value)
		if err != nil {
			return modelMap, err
		}
		modelMap["value"] = []map[string]interface{}{valueMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisTypeTimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsBarChartXAxisTypeTime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Time != nil {
		timeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartXAxisXAxisByTimeToMap(model.Time)
		if err != nil {
			return modelMap, err
		}
		modelMap["time"] = []map[string]interface{}{timeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Query != nil {
		queryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryToMap(model.Query)
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
		stackDefinitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinitionToMap(model.StackDefinition)
		if err != nil {
			return modelMap, err
		}
		modelMap["stack_definition"] = []map[string]interface{}{stackDefinitionMap}
	}
	if model.ScaleType != nil {
		modelMap["scale_type"] = *model.ScaleType
	}
	if model.ColorsBy != nil {
		colorsByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonColorsByToMap(model.ColorsBy)
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
		yAxisViewByMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByToMap(model.YAxisViewBy)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryToMap(model logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprimeToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQuery); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQuery)
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model.Logs)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs"] = []map[string]interface{}{logsMap}
		}
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Dataprime != nil {
			dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model.Dataprime)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartLogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LuceneQuery != nil {
		luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonLuceneQueryToMap(model.LuceneQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	}
	if model.Aggregation != nil {
		aggregationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLogsAggregationToMap(model.Aggregation)
		if err != nil {
			return modelMap, err
		}
		modelMap["aggregation"] = []map[string]interface{}{aggregationMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterLogsFilterToMap(&filtersItem)
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
			groupNamesFieldsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&groupNamesFieldsItem)
			if err != nil {
				return modelMap, err
			}
			groupNamesFields = append(groupNamesFields, groupNamesFieldsItemMap)
		}
		modelMap["group_names_fields"] = groupNamesFields
	}
	if model.StackedGroupNameField != nil {
		stackedGroupNameFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StackedGroupNameField)
		if err != nil {
			return modelMap, err
		}
		modelMap["stacked_group_name_field"] = []map[string]interface{}{stackedGroupNameFieldMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsCommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterMetricsFilterToMap(&filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DataprimeQuery != nil {
		dataprimeQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonDataprimeQueryToMap(model.DataprimeQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime_query"] = []map[string]interface{}{dataprimeQueryMap}
	}
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(filtersItem)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogsToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartLogsQueryToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetricsToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartMetricsQueryToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprimeToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartQueryValueDataprime) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dataprime != nil {
		dataprimeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartDataprimeQueryToMap(model.Dataprime)
		if err != nil {
			return modelMap, err
		}
		modelMap["dataprime"] = []map[string]interface{}{dataprimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinitionToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartStackDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSlicesPerBar != nil {
		modelMap["max_slices_per_bar"] = flex.IntValue(model.MaxSlicesPerBar)
	}
	if model.StackNameTemplate != nil {
		modelMap["stack_name_template"] = *model.StackNameTemplate
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByToMap(model logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategoryToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValueToMap(model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewBy)
		if model.Category != nil {
			categoryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model.Category)
			if err != nil {
				return modelMap, err
			}
			modelMap["category"] = []map[string]interface{}{categoryMap}
		}
		if model.Value != nil {
			valueMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model.Value)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategoryToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewCategory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Category != nil {
		categoryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByCategoryEmptyToMap(model.Category)
		if err != nil {
			return modelMap, err
		}
		modelMap["category"] = []map[string]interface{}{categoryMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValueToMap(model *logsv0.ApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		valueMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartYAxisViewByYAxisViewByValueEmptyToMap(model.Value)
		if err != nil {
			return modelMap, err
		}
		modelMap["value"] = []map[string]interface{}{valueMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model *logsv0.ApisDashboardsV1AstWidgetsMarkdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["markdown_text"] = *model.MarkdownText
	if model.TooltipText != nil {
		modelMap["tooltip_text"] = *model.TooltipText
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueLineChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueLineChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LineChart != nil {
		lineChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsLineChartToMap(model.LineChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["line_chart"] = []map[string]interface{}{lineChartMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueDataTableToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueDataTable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DataTable != nil {
		dataTableMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsDataTableToMap(model.DataTable)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_table"] = []map[string]interface{}{dataTableMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueGaugeToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueGauge) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Gauge != nil {
		gaugeMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsGaugeToMap(model.Gauge)
		if err != nil {
			return modelMap, err
		}
		modelMap["gauge"] = []map[string]interface{}{gaugeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValuePieChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValuePieChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PieChart != nil {
		pieChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsPieChartToMap(model.PieChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["pie_chart"] = []map[string]interface{}{pieChartMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BarChart != nil {
		barChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsBarChartToMap(model.BarChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["bar_chart"] = []map[string]interface{}{barChartMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChartToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueHorizontalBarChart) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.HorizontalBarChart != nil {
		horizontalBarChartMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsHorizontalBarChartToMap(model.HorizontalBarChart)
		if err != nil {
			return modelMap, err
		}
		modelMap["horizontal_bar_chart"] = []map[string]interface{}{horizontalBarChartMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetDefinitionValueMarkdownToMap(model *logsv0.ApisDashboardsV1AstWidgetDefinitionValueMarkdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Markdown != nil {
		markdownMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstWidgetsMarkdownToMap(model.Markdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["markdown"] = []map[string]interface{}{markdownMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstVariableToMap(model *logsv0.ApisDashboardsV1AstVariable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	definitionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["display_name"] = *model.DisplayName
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionToMap(model logsv0.ApisDashboardsV1AstVariableDefinitionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionValueMultiSelectToMap(model.(*logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstVariableDefinition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstVariableDefinition)
		if model.MultiSelect != nil {
			multiSelectMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model.MultiSelect)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model *logsv0.ApisDashboardsV1AstMultiSelect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	selectionMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionToMap(model.Selection)
	if err != nil {
		return modelMap, err
	}
	modelMap["selection"] = []map[string]interface{}{selectionMap}
	modelMap["values_order_direction"] = *model.ValuesOrderDirection
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceToMap(model logsv0.ApisDashboardsV1AstMultiSelectSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueLogsPathToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueMetricLabelToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueConstantListToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstMultiSelectSource)
		if model.LogsPath != nil {
			logsPathMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model.LogsPath)
			if err != nil {
				return modelMap, err
			}
			modelMap["logs_path"] = []map[string]interface{}{logsPathMap}
		}
		if model.MetricLabel != nil {
			metricLabelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model.MetricLabel)
			if err != nil {
				return modelMap, err
			}
			modelMap["metric_label"] = []map[string]interface{}{metricLabelMap}
		}
		if model.ConstantList != nil {
			constantListMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model.ConstantList)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectLogsPathSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	observationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.ObservationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["observation_field"] = []map[string]interface{}{observationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectMetricLabelSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["metric_name"] = *model.MetricName
	modelMap["label"] = *model.Label
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model *logsv0.ApisDashboardsV1AstMultiSelectConstantListSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueLogsPathToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueLogsPath) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LogsPath != nil {
		logsPathMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectLogsPathSourceToMap(model.LogsPath)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs_path"] = []map[string]interface{}{logsPathMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueMetricLabelToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueMetricLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MetricLabel != nil {
		metricLabelMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectMetricLabelSourceToMap(model.MetricLabel)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_label"] = []map[string]interface{}{metricLabelMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSourceValueConstantListToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSourceValueConstantList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConstantList != nil {
		constantListMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectConstantListSourceToMap(model.ConstantList)
		if err != nil {
			return modelMap, err
		}
		modelMap["constant_list"] = []map[string]interface{}{constantListMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionToMap(model logsv0.ApisDashboardsV1AstMultiSelectSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueAllToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueListToMap(model.(*logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstMultiSelectSelection)
		if model.All != nil {
			allMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model.All)
			if err != nil {
				return modelMap, err
			}
			modelMap["all"] = []map[string]interface{}{allMap}
		}
		if model.List != nil {
			listMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model.List)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionAllSelectionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionListSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Values != nil {
		modelMap["values"] = model.Values
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueAllToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionValueAll) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.All != nil {
		allMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionAllSelectionEmptyToMap(model.All)
		if err != nil {
			return modelMap, err
		}
		modelMap["all"] = []map[string]interface{}{allMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionValueListToMap(model *logsv0.ApisDashboardsV1AstMultiSelectSelectionValueList) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.List != nil {
		listMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectSelectionListSelectionToMap(model.List)
		if err != nil {
			return modelMap, err
		}
		modelMap["list"] = []map[string]interface{}{listMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstVariableDefinitionValueMultiSelectToMap(model *logsv0.ApisDashboardsV1AstVariableDefinitionValueMultiSelect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MultiSelect != nil {
		multiSelectMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstMultiSelectToMap(model.MultiSelect)
		if err != nil {
			return modelMap, err
		}
		modelMap["multi_select"] = []map[string]interface{}{multiSelectMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFilterToMap(model *logsv0.ApisDashboardsV1AstFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Source != nil {
		sourceMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstFilterSourceToMap(model.Source)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationToMap(model *logsv0.ApisDashboardsV1AstAnnotation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = model.Href.String()
	}
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	modelMap["enabled"] = *model.Enabled
	sourceMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceToMap(model logsv0.ApisDashboardsV1AstAnnotationSourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueMetricsToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueLogsToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationSource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstAnnotationSource)
		if model.Metrics != nil {
			metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model.Metrics)
			if err != nil {
				return modelMap, err
			}
			modelMap["metrics"] = []map[string]interface{}{metricsMap}
		}
		if model.Logs != nil {
			logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model.Logs)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PromqlQuery != nil {
		promqlQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonPromQlQueryToMap(model.PromqlQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["promql_query"] = []map[string]interface{}{promqlQueryMap}
	}
	if model.Strategy != nil {
		strategyMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStrategyToMap(model.Strategy)
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

func DataSourceIbmLogsDashboardApisDashboardsV1CommonPromQlQueryToMap(model *logsv0.ApisDashboardsV1CommonPromQlQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStrategyToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStrategy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StartTimeMetric != nil {
		startTimeMetricMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmptyToMap(model.StartTimeMetric)
		if err != nil {
			return modelMap, err
		}
		modelMap["start_time_metric"] = []map[string]interface{}{startTimeMetricMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmptyToMap(model *logsv0.ApisDashboardsV1AstAnnotationMetricsSourceStartTimeMetricEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	luceneQueryMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonLuceneQueryToMap(model.LuceneQuery)
	if err != nil {
		return modelMap, err
	}
	modelMap["lucene_query"] = []map[string]interface{}{luceneQueryMap}
	strategyMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyToMap(model.Strategy)
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
			labelFieldsItemMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(&labelFieldsItem)
			if err != nil {
				return modelMap, err
			}
			labelFields = append(labelFields, labelFieldsItemMap)
		}
		modelMap["label_fields"] = labelFields
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonLuceneQueryToMap(model *logsv0.ApisDashboardsV1CommonLuceneQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyToMap(model logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstantToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueRangeToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration); ok {
		return DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueDurationToMap(model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration))
	} else if _, ok := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategy)
		if model.Instant != nil {
			instantMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model.Instant)
			if err != nil {
				return modelMap, err
			}
			modelMap["instant"] = []map[string]interface{}{instantMap}
		}
		if model.Range != nil {
			rangeVarMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model.Range)
			if err != nil {
				return modelMap, err
			}
			modelMap["range"] = []map[string]interface{}{rangeVarMap}
		}
		if model.Duration != nil {
			durationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model.Duration)
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

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyInstant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	timestampFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.TimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["timestamp_field"] = []map[string]interface{}{timestampFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startTimestampFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StartTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_timestamp_field"] = []map[string]interface{}{startTimestampFieldMap}
	endTimestampFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.EndTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["end_timestamp_field"] = []map[string]interface{}{endTimestampFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyDuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startTimestampFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.StartTimestampField)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_timestamp_field"] = []map[string]interface{}{startTimestampFieldMap}
	durationFieldMap, err := DataSourceIbmLogsDashboardApisDashboardsV1CommonObservationFieldToMap(model.DurationField)
	if err != nil {
		return modelMap, err
	}
	modelMap["duration_field"] = []map[string]interface{}{durationFieldMap}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstantToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueInstant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instant != nil {
		instantMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyInstantToMap(model.Instant)
		if err != nil {
			return modelMap, err
		}
		modelMap["instant"] = []map[string]interface{}{instantMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueRangeToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Range != nil {
		rangeVarMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyRangeToMap(model.Range)
		if err != nil {
			return modelMap, err
		}
		modelMap["range"] = []map[string]interface{}{rangeVarMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyValueDurationToMap(model *logsv0.ApisDashboardsV1AstAnnotationLogsSourceStrategyValueDuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Duration != nil {
		durationMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceStrategyDurationToMap(model.Duration)
		if err != nil {
			return modelMap, err
		}
		modelMap["duration"] = []map[string]interface{}{durationMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueMetricsToMap(model *logsv0.ApisDashboardsV1AstAnnotationSourceValueMetrics) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Metrics != nil {
		metricsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationMetricsSourceToMap(model.Metrics)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics"] = []map[string]interface{}{metricsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationSourceValueLogsToMap(model *logsv0.ApisDashboardsV1AstAnnotationSourceValueLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Logs != nil {
		logsMap, err := DataSourceIbmLogsDashboardApisDashboardsV1AstAnnotationLogsSourceToMap(model.Logs)
		if err != nil {
			return modelMap, err
		}
		modelMap["logs"] = []map[string]interface{}{logsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1CommonTimeFrameToMap(model *logsv0.ApisDashboardsV1CommonTimeFrame) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.From != nil {
		modelMap["from"] = model.From.String()
	}
	if model.To != nil {
		modelMap["to"] = model.To.String()
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstFolderPathToMap(model *logsv0.ApisDashboardsV1AstFolderPath) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Segments != nil {
		modelMap["segments"] = model.Segments
	}
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshOffEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshOffEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshTwoMinutesEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsDashboardApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmptyToMap(model *logsv0.ApisDashboardsV1AstDashboardAutoRefreshFiveMinutesEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}
