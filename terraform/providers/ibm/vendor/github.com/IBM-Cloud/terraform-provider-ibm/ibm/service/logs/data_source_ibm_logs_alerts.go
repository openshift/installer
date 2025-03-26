// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsAlerts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsAlertsRead,

		Schema: map[string]*schema.Schema{
			"alerts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alerts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert name.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert description.",
						},
						"is_active": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Alert is active.",
						},
						"severity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert severity.",
						},
						"expiration": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alert expiration date.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"year": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Year.",
									},
									"month": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Month of the year.",
									},
									"day": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Day of the month.",
									},
								},
							},
						},
						"condition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alert condition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"immediate": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for immediate standard alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"less_than": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for less than alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"more_than": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for more than alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
												"evaluation_window": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The evaluation window for the alert condition.",
												},
											},
										},
									},
									"more_than_usual": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for more than usual alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"new_value": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for new value alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"flow": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for flow alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stages": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Flow alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"groups": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of groups of alerts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"alerts": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of alerts.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"op": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Operator for the alerts.",
																					},
																					"values": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "List of alerts.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"id": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The alert ID.",
																								},
																								"not": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "The alert not.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"next_op": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Operator for the alerts.",
																		},
																	},
																},
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Timeframe for the flow.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ms": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Timeframe in milliseconds.",
																		},
																	},
																},
															},
														},
													},
												},
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
												"enforce_suppression": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Should suppression be enforced on the flow alert.",
												},
											},
										},
									},
									"unique_count": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for unique count alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"less_than_usual": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition for less than usual alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Less than alert condition parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"threshold": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "The threshold for the alert condition.",
															},
															"timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The timeframe for the alert condition.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The group by fields for the alert condition.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"metric_alert_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The lucene metric alert parameters if it is a lucene metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"metric_field": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric field of the metric alert.",
																		},
																		"metric_source": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The metric source of the metric alert.",
																		},
																		"arithmetic_operator": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator modifier of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"metric_alert_promql_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The promql metric alert parameters if is is a promql metric alert.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"promql_text": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The promql text of the metric alert by fields for the alert condition.",
																		},
																		"arithmetic_operator_modifier": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The arithmetic operator of the metric promql alert.",
																		},
																		"sample_threshold_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "The threshold percentage.",
																		},
																		"non_null_percentage": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Non null percentage of the evaluation.",
																		},
																		"swap_null_values": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we swap null values with zero.",
																		},
																	},
																},
															},
															"ignore_infinity": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Should the evaluation ignore infinity value.",
															},
															"relative_timeframe": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative timeframe for time relative alerts.",
															},
															"cardinality_fields": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Cardinality fields for unique count alert.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"related_extended_data": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Deadman configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cleanup_deadman_duration": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Cleanup deadman duration.",
																		},
																		"should_trigger_deadman": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Should we trigger deadman.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"notification_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alert notification groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_by_fields": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Group by fields to group the values by.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"notifications": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Webhook target settings for the the notification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retriggering_period_seconds": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Retriggering period of the alert in seconds.",
												},
												"notify_on": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Notify on setting.",
												},
												"integration_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Integration ID.",
												},
												"recipients": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Recipients.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"emails": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Email addresses.",
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
						"filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alert filters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severities": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The severity of the logs to filter.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"metadata": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The metadata filters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"applications": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The applications to filter.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"subsystems": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The subsystems to filter.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"alias": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The alias of the filter.",
									},
									"text": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text to filter.",
									},
									"ratio_alerts": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The ratio alerts.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alias": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The alias of the filter.",
												},
												"text": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The text to filter.",
												},
												"severities": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The severities to filter.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"applications": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The applications to filter.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"subsystems": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The subsystems to filter.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The group by fields.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"filter_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the filter.",
									},
								},
							},
						},
						"active_when": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "When should the alert be active.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeframes": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Activity timeframes of the alert.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days_of_week": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Days of the week for activity.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"range": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Time range in the day of the week.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Start time.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"hours": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Hours of the day.",
																		},
																		"minutes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minutes of the hour.",
																		},
																		"seconds": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Seconds of the minute.",
																		},
																	},
																},
															},
															"end": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Start time.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"hours": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Hours of the day.",
																		},
																		"minutes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minutes of the hour.",
																		},
																		"seconds": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Seconds of the minute.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"notification_payload_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "JSON keys to include in the alert notification, if left empty get the full log text in the alert notification.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"meta_labels": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Meta labels to add to the alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of the label.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the label.",
									},
								},
							},
						},
						"meta_labels_strings": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Meta labels to add to the alert as string with ':' separator.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"unique_identifier": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert unique identifier.",
						},
						"incident_settings": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Incident settings, will create the incident based on this configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retriggering_period_seconds": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The retriggering period of the alert in seconds.",
									},
									"notify_on": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notify on setting.",
									},
									"use_as_notification_settings": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Use these settings for all notificaion webhook.",
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

func dataSourceIbmLogsAlertsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alerts", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getAlertsOptions := &logsv0.GetAlertsOptions{}

	alertCollection, _, err := logsClient.GetAlertsWithContext(context, getAlertsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAlertsWithContext failed: %s", err.Error()), "(Data) ibm_logs_alerts", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsAlertsID(d))

	alerts := []map[string]interface{}{}
	if alertCollection.Alerts != nil {
		for _, modelItem := range alertCollection.Alerts {
			modelMap, err := DataSourceIbmLogsAlertsAlertToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alerts", "read")
				return tfErr.GetDiag()
			}
			alerts = append(alerts, modelMap)
		}
	}
	if err = d.Set("alerts", alerts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting alerts: %s", err), "(Data) ibm_logs_alerts", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsAlertsID returns a reasonable ID for the list.
func dataSourceIbmLogsAlertsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsAlertsAlertToMap(model *logsv0.Alert) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["name"] = *model.Name
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	modelMap["is_active"] = *model.IsActive
	modelMap["severity"] = *model.Severity
	if model.Expiration != nil {
		expirationMap, err := DataSourceIbmLogsAlertsAlertsV1DateToMap(model.Expiration)
		if err != nil {
			return modelMap, err
		}
		modelMap["expiration"] = []map[string]interface{}{expirationMap}
	}
	conditionMap, err := DataSourceIbmLogsAlertsAlertsV2AlertConditionToMap(model.Condition)
	if err != nil {
		return modelMap, err
	}
	modelMap["condition"] = []map[string]interface{}{conditionMap}
	notificationGroups := []map[string]interface{}{}
	for _, notificationGroupsItem := range model.NotificationGroups {
		notificationGroupsItemMap, err := DataSourceIbmLogsAlertsAlertsV2AlertNotificationGroupsToMap(&notificationGroupsItem)
		if err != nil {
			return modelMap, err
		}
		notificationGroups = append(notificationGroups, notificationGroupsItemMap)
	}
	modelMap["notification_groups"] = notificationGroups
	filtersMap, err := DataSourceIbmLogsAlertsAlertsV1AlertFiltersToMap(model.Filters)
	if err != nil {
		return modelMap, err
	}
	modelMap["filters"] = []map[string]interface{}{filtersMap}
	if model.ActiveWhen != nil {
		activeWhenMap, err := DataSourceIbmLogsAlertsAlertsV1AlertActiveWhenToMap(model.ActiveWhen)
		if err != nil {
			return modelMap, err
		}
		modelMap["active_when"] = []map[string]interface{}{activeWhenMap}
	}
	if model.NotificationPayloadFilters != nil {
		modelMap["notification_payload_filters"] = model.NotificationPayloadFilters
	}
	if model.MetaLabels != nil {
		metaLabels := []map[string]interface{}{}
		for _, metaLabelsItem := range model.MetaLabels {
			metaLabelsItemMap, err := DataSourceIbmLogsAlertsAlertsV1MetaLabelToMap(&metaLabelsItem)
			if err != nil {
				return modelMap, err
			}
			metaLabels = append(metaLabels, metaLabelsItemMap)
		}
		modelMap["meta_labels"] = metaLabels
	}
	if model.MetaLabelsStrings != nil {
		modelMap["meta_labels_strings"] = model.MetaLabelsStrings
	}
	if model.UniqueIdentifier != nil {
		modelMap["unique_identifier"] = model.UniqueIdentifier.String()
	}
	if model.IncidentSettings != nil {
		incidentSettingsMap, err := DataSourceIbmLogsAlertsAlertsV2AlertIncidentSettingsToMap(model.IncidentSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["incident_settings"] = []map[string]interface{}{incidentSettingsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1DateToMap(model *logsv0.AlertsV1Date) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Year != nil {
		modelMap["year"] = flex.IntValue(model.Year)
	}
	if model.Month != nil {
		modelMap["month"] = flex.IntValue(model.Month)
	}
	if model.Day != nil {
		modelMap["day"] = flex.IntValue(model.Day)
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionToMap(model logsv0.AlertsV2AlertConditionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionImmediate); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionImmediateToMap(model.(*logsv0.AlertsV2AlertConditionConditionImmediate))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThan); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionLessThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThan); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionMoreThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionMoreThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionNewValue); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionNewValueToMap(model.(*logsv0.AlertsV2AlertConditionConditionNewValue))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionFlow); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionFlowToMap(model.(*logsv0.AlertsV2AlertConditionConditionFlow))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionUniqueCountToMap(model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionLessThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertCondition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.AlertsV2AlertCondition)
		if model.Immediate != nil {
			immediateMap, err := DataSourceIbmLogsAlertsAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
			if err != nil {
				return modelMap, err
			}
			modelMap["immediate"] = []map[string]interface{}{immediateMap}
		}
		if model.LessThan != nil {
			lessThanMap, err := DataSourceIbmLogsAlertsAlertsV2LessThanConditionToMap(model.LessThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["less_than"] = []map[string]interface{}{lessThanMap}
		}
		if model.MoreThan != nil {
			moreThanMap, err := DataSourceIbmLogsAlertsAlertsV2MoreThanConditionToMap(model.MoreThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than"] = []map[string]interface{}{moreThanMap}
		}
		if model.MoreThanUsual != nil {
			moreThanUsualMap, err := DataSourceIbmLogsAlertsAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
		}
		if model.NewValue != nil {
			newValueMap, err := DataSourceIbmLogsAlertsAlertsV2NewValueConditionToMap(model.NewValue)
			if err != nil {
				return modelMap, err
			}
			modelMap["new_value"] = []map[string]interface{}{newValueMap}
		}
		if model.Flow != nil {
			flowMap, err := DataSourceIbmLogsAlertsAlertsV2FlowConditionToMap(model.Flow)
			if err != nil {
				return modelMap, err
			}
			modelMap["flow"] = []map[string]interface{}{flowMap}
		}
		if model.UniqueCount != nil {
			uniqueCountMap, err := DataSourceIbmLogsAlertsAlertsV2UniqueCountConditionToMap(model.UniqueCount)
			if err != nil {
				return modelMap, err
			}
			modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
		}
		if model.LessThanUsual != nil {
			lessThanUsualMap, err := DataSourceIbmLogsAlertsAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
			if err != nil {
				return modelMap, err
			}
			modelMap["less_than_usual"] = []map[string]interface{}{lessThanUsualMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.AlertsV2AlertConditionIntf subtype encountered")
	}
}

func DataSourceIbmLogsAlertsAlertsV2ImmediateConditionEmptyToMap(model *logsv0.AlertsV2ImmediateConditionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2LessThanConditionToMap(model *logsv0.AlertsV2LessThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model *logsv0.AlertsV2ConditionParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["threshold"] = *model.Threshold
	modelMap["timeframe"] = *model.Timeframe
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	if model.MetricAlertParameters != nil {
		metricAlertParametersMap, err := DataSourceIbmLogsAlertsAlertsV1MetricAlertConditionParametersToMap(model.MetricAlertParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_alert_parameters"] = []map[string]interface{}{metricAlertParametersMap}
	}
	if model.MetricAlertPromqlParameters != nil {
		metricAlertPromqlParametersMap, err := DataSourceIbmLogsAlertsAlertsV1MetricAlertPromqlConditionParametersToMap(model.MetricAlertPromqlParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_alert_promql_parameters"] = []map[string]interface{}{metricAlertPromqlParametersMap}
	}
	if model.IgnoreInfinity != nil {
		modelMap["ignore_infinity"] = *model.IgnoreInfinity
	}
	if model.RelativeTimeframe != nil {
		modelMap["relative_timeframe"] = *model.RelativeTimeframe
	}
	if model.CardinalityFields != nil {
		modelMap["cardinality_fields"] = model.CardinalityFields
	}
	if model.RelatedExtendedData != nil {
		relatedExtendedDataMap, err := DataSourceIbmLogsAlertsAlertsV1RelatedExtendedDataToMap(model.RelatedExtendedData)
		if err != nil {
			return modelMap, err
		}
		modelMap["related_extended_data"] = []map[string]interface{}{relatedExtendedDataMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1MetricAlertConditionParametersToMap(model *logsv0.AlertsV1MetricAlertConditionParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["metric_field"] = *model.MetricField
	modelMap["metric_source"] = *model.MetricSource
	modelMap["arithmetic_operator"] = *model.ArithmeticOperator
	if model.ArithmeticOperatorModifier != nil {
		modelMap["arithmetic_operator_modifier"] = flex.IntValue(model.ArithmeticOperatorModifier)
	}
	if model.SampleThresholdPercentage != nil {
		modelMap["sample_threshold_percentage"] = flex.IntValue(model.SampleThresholdPercentage)
	}
	if model.NonNullPercentage != nil {
		modelMap["non_null_percentage"] = flex.IntValue(model.NonNullPercentage)
	}
	if model.SwapNullValues != nil {
		modelMap["swap_null_values"] = *model.SwapNullValues
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1MetricAlertPromqlConditionParametersToMap(model *logsv0.AlertsV1MetricAlertPromqlConditionParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["promql_text"] = *model.PromqlText
	if model.ArithmeticOperatorModifier != nil {
		modelMap["arithmetic_operator_modifier"] = flex.IntValue(model.ArithmeticOperatorModifier)
	}
	modelMap["sample_threshold_percentage"] = flex.IntValue(model.SampleThresholdPercentage)
	if model.NonNullPercentage != nil {
		modelMap["non_null_percentage"] = flex.IntValue(model.NonNullPercentage)
	}
	if model.SwapNullValues != nil {
		modelMap["swap_null_values"] = *model.SwapNullValues
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1RelatedExtendedDataToMap(model *logsv0.AlertsV1RelatedExtendedData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CleanupDeadmanDuration != nil {
		modelMap["cleanup_deadman_duration"] = *model.CleanupDeadmanDuration
	}
	if model.ShouldTriggerDeadman != nil {
		modelMap["should_trigger_deadman"] = *model.ShouldTriggerDeadman
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2MoreThanConditionToMap(model *logsv0.AlertsV2MoreThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	if model.EvaluationWindow != nil {
		modelMap["evaluation_window"] = *model.EvaluationWindow
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2MoreThanUsualConditionToMap(model *logsv0.AlertsV2MoreThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2NewValueConditionToMap(model *logsv0.AlertsV2NewValueCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2FlowConditionToMap(model *logsv0.AlertsV2FlowCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Stages != nil {
		stages := []map[string]interface{}{}
		for _, stagesItem := range model.Stages {
			stagesItemMap, err := DataSourceIbmLogsAlertsAlertsV1FlowStageToMap(&stagesItem)
			if err != nil {
				return modelMap, err
			}
			stages = append(stages, stagesItemMap)
		}
		modelMap["stages"] = stages
	}
	if model.Parameters != nil {
		parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	if model.EnforceSuppression != nil {
		modelMap["enforce_suppression"] = *model.EnforceSuppression
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1FlowStageToMap(model *logsv0.AlertsV1FlowStage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Groups != nil {
		groups := []map[string]interface{}{}
		for _, groupsItem := range model.Groups {
			groupsItemMap, err := DataSourceIbmLogsAlertsAlertsV1FlowGroupToMap(&groupsItem)
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Timeframe != nil {
		timeframeMap, err := DataSourceIbmLogsAlertsAlertsV1FlowTimeframeToMap(model.Timeframe)
		if err != nil {
			return modelMap, err
		}
		modelMap["timeframe"] = []map[string]interface{}{timeframeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1FlowGroupToMap(model *logsv0.AlertsV1FlowGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Alerts != nil {
		alertsMap, err := DataSourceIbmLogsAlertsAlertsV1FlowAlertsToMap(model.Alerts)
		if err != nil {
			return modelMap, err
		}
		modelMap["alerts"] = []map[string]interface{}{alertsMap}
	}
	if model.NextOp != nil {
		modelMap["next_op"] = *model.NextOp
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1FlowAlertsToMap(model *logsv0.AlertsV1FlowAlerts) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Op != nil {
		modelMap["op"] = *model.Op
	}
	if model.Values != nil {
		values := []map[string]interface{}{}
		for _, valuesItem := range model.Values {
			valuesItemMap, err := DataSourceIbmLogsAlertsAlertsV1FlowAlertToMap(&valuesItem)
			if err != nil {
				return modelMap, err
			}
			values = append(values, valuesItemMap)
		}
		modelMap["values"] = values
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1FlowAlertToMap(model *logsv0.AlertsV1FlowAlert) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Not != nil {
		modelMap["not"] = *model.Not
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1FlowTimeframeToMap(model *logsv0.AlertsV1FlowTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Ms != nil {
		modelMap["ms"] = flex.IntValue(model.Ms)
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2UniqueCountConditionToMap(model *logsv0.AlertsV2UniqueCountCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2LessThanUsualConditionToMap(model *logsv0.AlertsV2LessThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertsAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionImmediateToMap(model *logsv0.AlertsV2AlertConditionConditionImmediate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Immediate != nil {
		immediateMap, err := DataSourceIbmLogsAlertsAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
		if err != nil {
			return modelMap, err
		}
		modelMap["immediate"] = []map[string]interface{}{immediateMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionLessThanToMap(model *logsv0.AlertsV2AlertConditionConditionLessThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThan != nil {
		lessThanMap, err := DataSourceIbmLogsAlertsAlertsV2LessThanConditionToMap(model.LessThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than"] = []map[string]interface{}{lessThanMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionMoreThanToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThan != nil {
		moreThanMap, err := DataSourceIbmLogsAlertsAlertsV2MoreThanConditionToMap(model.MoreThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than"] = []map[string]interface{}{moreThanMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionMoreThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThanUsual != nil {
		moreThanUsualMap, err := DataSourceIbmLogsAlertsAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionNewValueToMap(model *logsv0.AlertsV2AlertConditionConditionNewValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewValue != nil {
		newValueMap, err := DataSourceIbmLogsAlertsAlertsV2NewValueConditionToMap(model.NewValue)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_value"] = []map[string]interface{}{newValueMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionFlowToMap(model *logsv0.AlertsV2AlertConditionConditionFlow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Flow != nil {
		flowMap, err := DataSourceIbmLogsAlertsAlertsV2FlowConditionToMap(model.Flow)
		if err != nil {
			return modelMap, err
		}
		modelMap["flow"] = []map[string]interface{}{flowMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionUniqueCountToMap(model *logsv0.AlertsV2AlertConditionConditionUniqueCount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UniqueCount != nil {
		uniqueCountMap, err := DataSourceIbmLogsAlertsAlertsV2UniqueCountConditionToMap(model.UniqueCount)
		if err != nil {
			return modelMap, err
		}
		modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertConditionConditionLessThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionLessThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThanUsual != nil {
		lessThanUsualMap, err := DataSourceIbmLogsAlertsAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than_usual"] = []map[string]interface{}{lessThanUsualMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertNotificationGroupsToMap(model *logsv0.AlertsV2AlertNotificationGroups) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GroupByFields != nil {
		modelMap["group_by_fields"] = model.GroupByFields
	}
	if model.Notifications != nil {
		notifications := []map[string]interface{}{}
		for _, notificationsItem := range model.Notifications {
			notificationsItemMap, err := DataSourceIbmLogsAlertsAlertsV2AlertNotificationToMap(notificationsItem)
			if err != nil {
				return modelMap, err
			}
			notifications = append(notifications, notificationsItemMap)
		}
		modelMap["notifications"] = notifications
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertNotificationToMap(model logsv0.AlertsV2AlertNotificationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID))
	} else if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients); ok {
		return DataSourceIbmLogsAlertsAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients))
	} else if _, ok := model.(*logsv0.AlertsV2AlertNotification); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.AlertsV2AlertNotification)
		if model.RetriggeringPeriodSeconds != nil {
			modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
		}
		if model.NotifyOn != nil {
			modelMap["notify_on"] = *model.NotifyOn
		}
		if model.IntegrationID != nil {
			modelMap["integration_id"] = flex.IntValue(model.IntegrationID)
		}
		if model.Recipients != nil {
			recipientsMap, err := DataSourceIbmLogsAlertsAlertsV2RecipientsToMap(model.Recipients)
			if err != nil {
				return modelMap, err
			}
			modelMap["recipients"] = []map[string]interface{}{recipientsMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.AlertsV2AlertNotificationIntf subtype encountered")
	}
}

func DataSourceIbmLogsAlertsAlertsV2RecipientsToMap(model *logsv0.AlertsV2Recipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Emails != nil {
		modelMap["emails"] = model.Emails
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RetriggeringPeriodSeconds != nil {
		modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
	}
	if model.NotifyOn != nil {
		modelMap["notify_on"] = *model.NotifyOn
	}
	if model.IntegrationID != nil {
		modelMap["integration_id"] = flex.IntValue(model.IntegrationID)
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RetriggeringPeriodSeconds != nil {
		modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
	}
	if model.NotifyOn != nil {
		modelMap["notify_on"] = *model.NotifyOn
	}
	if model.Recipients != nil {
		recipientsMap, err := DataSourceIbmLogsAlertsAlertsV2RecipientsToMap(model.Recipients)
		if err != nil {
			return modelMap, err
		}
		modelMap["recipients"] = []map[string]interface{}{recipientsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1AlertFiltersToMap(model *logsv0.AlertsV1AlertFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	if model.Metadata != nil {
		metadataMap, err := DataSourceIbmLogsAlertsAlertsV1AlertFiltersMetadataFiltersToMap(model.Metadata)
		if err != nil {
			return modelMap, err
		}
		modelMap["metadata"] = []map[string]interface{}{metadataMap}
	}
	if model.Alias != nil {
		modelMap["alias"] = *model.Alias
	}
	if model.Text != nil {
		modelMap["text"] = *model.Text
	}
	if model.RatioAlerts != nil {
		ratioAlerts := []map[string]interface{}{}
		for _, ratioAlertsItem := range model.RatioAlerts {
			ratioAlertsItemMap, err := DataSourceIbmLogsAlertsAlertsV1AlertFiltersRatioAlertToMap(&ratioAlertsItem)
			if err != nil {
				return modelMap, err
			}
			ratioAlerts = append(ratioAlerts, ratioAlertsItemMap)
		}
		modelMap["ratio_alerts"] = ratioAlerts
	}
	if model.FilterType != nil {
		modelMap["filter_type"] = *model.FilterType
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1AlertFiltersMetadataFiltersToMap(model *logsv0.AlertsV1AlertFiltersMetadataFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	if model.Subsystems != nil {
		modelMap["subsystems"] = model.Subsystems
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1AlertFiltersRatioAlertToMap(model *logsv0.AlertsV1AlertFiltersRatioAlert) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["alias"] = *model.Alias
	if model.Text != nil {
		modelMap["text"] = *model.Text
	}
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	if model.Subsystems != nil {
		modelMap["subsystems"] = model.Subsystems
	}
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1AlertActiveWhenToMap(model *logsv0.AlertsV1AlertActiveWhen) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	timeframes := []map[string]interface{}{}
	for _, timeframesItem := range model.Timeframes {
		timeframesItemMap, err := DataSourceIbmLogsAlertsAlertsV1AlertActiveTimeframeToMap(&timeframesItem)
		if err != nil {
			return modelMap, err
		}
		timeframes = append(timeframes, timeframesItemMap)
	}
	modelMap["timeframes"] = timeframes
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1AlertActiveTimeframeToMap(model *logsv0.AlertsV1AlertActiveTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["days_of_week"] = model.DaysOfWeek
	rangeVarMap, err := DataSourceIbmLogsAlertsAlertsV1TimeRangeToMap(model.Range)
	if err != nil {
		return modelMap, err
	}
	modelMap["range"] = []map[string]interface{}{rangeVarMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1TimeRangeToMap(model *logsv0.AlertsV1TimeRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startMap, err := DataSourceIbmLogsAlertsAlertsV1TimeToMap(model.Start)
	if err != nil {
		return modelMap, err
	}
	modelMap["start"] = []map[string]interface{}{startMap}
	endMap, err := DataSourceIbmLogsAlertsAlertsV1TimeToMap(model.End)
	if err != nil {
		return modelMap, err
	}
	modelMap["end"] = []map[string]interface{}{endMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1TimeToMap(model *logsv0.AlertsV1Time) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hours != nil {
		modelMap["hours"] = flex.IntValue(model.Hours)
	}
	if model.Minutes != nil {
		modelMap["minutes"] = flex.IntValue(model.Minutes)
	}
	if model.Seconds != nil {
		modelMap["seconds"] = flex.IntValue(model.Seconds)
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV1MetaLabelToMap(model *logsv0.AlertsV1MetaLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertsAlertsV2AlertIncidentSettingsToMap(model *logsv0.AlertsV2AlertIncidentSettings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RetriggeringPeriodSeconds != nil {
		modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
	}
	if model.NotifyOn != nil {
		modelMap["notify_on"] = *model.NotifyOn
	}
	if model.UseAsNotificationSettings != nil {
		modelMap["use_as_notification_settings"] = *model.UseAsNotificationSettings
	}
	return modelMap, nil
}
