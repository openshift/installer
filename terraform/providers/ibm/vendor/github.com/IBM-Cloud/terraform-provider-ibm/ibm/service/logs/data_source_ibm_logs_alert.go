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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsAlertRead,

		Schema: map[string]*schema.Schema{
			"logs_alert_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func dataSourceIbmLogsAlertRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getAlertOptions := &logsv0.GetAlertOptions{}

	getAlertOptions.SetID(core.UUIDPtr(strfmt.UUID(d.Get("logs_alert_id").(string))))

	alert, _, err := logsClient.GetAlertWithContext(context, getAlertOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAlertWithContext failed: %s", err.Error()), "(Data) ibm_logs_alert", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getAlertOptions.ID))

	if err = d.Set("name", alert.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("description", alert.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("is_active", alert.IsActive); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting is_active: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("severity", alert.Severity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting severity: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	expiration := []map[string]interface{}{}
	if alert.Expiration != nil {
		modelMap, err := DataSourceIbmLogsAlertAlertsV1DateToMap(alert.Expiration)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
			return tfErr.GetDiag()
		}
		expiration = append(expiration, modelMap)
	}
	if err = d.Set("expiration", expiration); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	condition := []map[string]interface{}{}
	if alert.Condition != nil {
		modelMap, err := DataSourceIbmLogsAlertAlertsV2AlertConditionToMap(alert.Condition)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
			return tfErr.GetDiag()
		}
		condition = append(condition, modelMap)
	}
	if err = d.Set("condition", condition); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting condition: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	notificationGroups := []map[string]interface{}{}
	if alert.NotificationGroups != nil {
		for _, modelItem := range alert.NotificationGroups {
			modelMap, err := DataSourceIbmLogsAlertAlertsV2AlertNotificationGroupsToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
				return tfErr.GetDiag()
			}
			notificationGroups = append(notificationGroups, modelMap)
		}
	}
	if err = d.Set("notification_groups", notificationGroups); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting notification_groups: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	filters := []map[string]interface{}{}
	if alert.Filters != nil {
		modelMap, err := DataSourceIbmLogsAlertAlertsV1AlertFiltersToMap(alert.Filters)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
			return tfErr.GetDiag()
		}
		filters = append(filters, modelMap)
	}
	if err = d.Set("filters", filters); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting filters: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	activeWhen := []map[string]interface{}{}
	if alert.ActiveWhen != nil {
		modelMap, err := DataSourceIbmLogsAlertAlertsV1AlertActiveWhenToMap(alert.ActiveWhen)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
			return tfErr.GetDiag()
		}
		activeWhen = append(activeWhen, modelMap)
	}
	if err = d.Set("active_when", activeWhen); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting active_when: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	metaLabels := []map[string]interface{}{}
	if alert.MetaLabels != nil {
		for _, modelItem := range alert.MetaLabels {
			modelMap, err := DataSourceIbmLogsAlertAlertsV1MetaLabelToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
				return tfErr.GetDiag()
			}
			metaLabels = append(metaLabels, modelMap)
		}
	}
	if err = d.Set("meta_labels", metaLabels); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting meta_labels: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("unique_identifier", alert.UniqueIdentifier); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting unique_identifier: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	incidentSettings := []map[string]interface{}{}
	if alert.IncidentSettings != nil {
		modelMap, err := DataSourceIbmLogsAlertAlertsV2AlertIncidentSettingsToMap(alert.IncidentSettings)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_alert", "read")
			return tfErr.GetDiag()
		}
		incidentSettings = append(incidentSettings, modelMap)
	}
	if err = d.Set("incident_settings", incidentSettings); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting incident_settings: %s", err), "(Data) ibm_logs_alert", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIbmLogsAlertAlertsV1DateToMap(model *logsv0.AlertsV1Date) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV2AlertConditionToMap(model logsv0.AlertsV2AlertConditionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionImmediate); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionImmediateToMap(model.(*logsv0.AlertsV2AlertConditionConditionImmediate))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThan); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThan); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionNewValue); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionNewValueToMap(model.(*logsv0.AlertsV2AlertConditionConditionNewValue))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionFlow); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionFlowToMap(model.(*logsv0.AlertsV2AlertConditionConditionFlow))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionUniqueCountToMap(model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertCondition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.AlertsV2AlertCondition)
		if model.Immediate != nil {
			immediateMap, err := DataSourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
			if err != nil {
				return modelMap, err
			}
			modelMap["immediate"] = []map[string]interface{}{immediateMap}
		}
		if model.LessThan != nil {
			lessThanMap, err := DataSourceIbmLogsAlertAlertsV2LessThanConditionToMap(model.LessThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["less_than"] = []map[string]interface{}{lessThanMap}
		}
		if model.MoreThan != nil {
			moreThanMap, err := DataSourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model.MoreThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than"] = []map[string]interface{}{moreThanMap}
		}
		if model.MoreThanUsual != nil {
			moreThanUsualMap, err := DataSourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
		}
		if model.NewValue != nil {
			newValueMap, err := DataSourceIbmLogsAlertAlertsV2NewValueConditionToMap(model.NewValue)
			if err != nil {
				return modelMap, err
			}
			modelMap["new_value"] = []map[string]interface{}{newValueMap}
		}
		if model.Flow != nil {
			flowMap, err := DataSourceIbmLogsAlertAlertsV2FlowConditionToMap(model.Flow)
			if err != nil {
				return modelMap, err
			}
			modelMap["flow"] = []map[string]interface{}{flowMap}
		}
		if model.UniqueCount != nil {
			uniqueCountMap, err := DataSourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model.UniqueCount)
			if err != nil {
				return modelMap, err
			}
			modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
		}
		if model.LessThanUsual != nil {
			lessThanUsualMap, err := DataSourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
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

func DataSourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model *logsv0.AlertsV2ImmediateConditionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2LessThanConditionToMap(model *logsv0.AlertsV2LessThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model *logsv0.AlertsV2ConditionParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["threshold"] = *model.Threshold
	modelMap["timeframe"] = *model.Timeframe
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	if model.MetricAlertParameters != nil {
		metricAlertParametersMap, err := DataSourceIbmLogsAlertAlertsV1MetricAlertConditionParametersToMap(model.MetricAlertParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_alert_parameters"] = []map[string]interface{}{metricAlertParametersMap}
	}
	if model.MetricAlertPromqlParameters != nil {
		metricAlertPromqlParametersMap, err := DataSourceIbmLogsAlertAlertsV1MetricAlertPromqlConditionParametersToMap(model.MetricAlertPromqlParameters)
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
		relatedExtendedDataMap, err := DataSourceIbmLogsAlertAlertsV1RelatedExtendedDataToMap(model.RelatedExtendedData)
		if err != nil {
			return modelMap, err
		}
		modelMap["related_extended_data"] = []map[string]interface{}{relatedExtendedDataMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1MetricAlertConditionParametersToMap(model *logsv0.AlertsV1MetricAlertConditionParameters) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV1MetricAlertPromqlConditionParametersToMap(model *logsv0.AlertsV1MetricAlertPromqlConditionParameters) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV1RelatedExtendedDataToMap(model *logsv0.AlertsV1RelatedExtendedData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CleanupDeadmanDuration != nil {
		modelMap["cleanup_deadman_duration"] = *model.CleanupDeadmanDuration
	}
	if model.ShouldTriggerDeadman != nil {
		modelMap["should_trigger_deadman"] = *model.ShouldTriggerDeadman
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model *logsv0.AlertsV2MoreThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	if model.EvaluationWindow != nil {
		modelMap["evaluation_window"] = *model.EvaluationWindow
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model *logsv0.AlertsV2MoreThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2NewValueConditionToMap(model *logsv0.AlertsV2NewValueCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2FlowConditionToMap(model *logsv0.AlertsV2FlowCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Stages != nil {
		stages := []map[string]interface{}{}
		for _, stagesItem := range model.Stages {
			stagesItemMap, err := DataSourceIbmLogsAlertAlertsV1FlowStageToMap(&stagesItem)
			if err != nil {
				return modelMap, err
			}
			stages = append(stages, stagesItemMap)
		}
		modelMap["stages"] = stages
	}
	if model.Parameters != nil {
		parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
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

func DataSourceIbmLogsAlertAlertsV1FlowStageToMap(model *logsv0.AlertsV1FlowStage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Groups != nil {
		groups := []map[string]interface{}{}
		for _, groupsItem := range model.Groups {
			groupsItemMap, err := DataSourceIbmLogsAlertAlertsV1FlowGroupToMap(&groupsItem)
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Timeframe != nil {
		timeframeMap, err := DataSourceIbmLogsAlertAlertsV1FlowTimeframeToMap(model.Timeframe)
		if err != nil {
			return modelMap, err
		}
		modelMap["timeframe"] = []map[string]interface{}{timeframeMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1FlowGroupToMap(model *logsv0.AlertsV1FlowGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Alerts != nil {
		alertsMap, err := DataSourceIbmLogsAlertAlertsV1FlowAlertsToMap(model.Alerts)
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

func DataSourceIbmLogsAlertAlertsV1FlowAlertsToMap(model *logsv0.AlertsV1FlowAlerts) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Op != nil {
		modelMap["op"] = *model.Op
	}
	if model.Values != nil {
		values := []map[string]interface{}{}
		for _, valuesItem := range model.Values {
			valuesItemMap, err := DataSourceIbmLogsAlertAlertsV1FlowAlertToMap(&valuesItem)
			if err != nil {
				return modelMap, err
			}
			values = append(values, valuesItemMap)
		}
		modelMap["values"] = values
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1FlowAlertToMap(model *logsv0.AlertsV1FlowAlert) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Not != nil {
		modelMap["not"] = *model.Not
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1FlowTimeframeToMap(model *logsv0.AlertsV1FlowTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Ms != nil {
		modelMap["ms"] = flex.IntValue(model.Ms)
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model *logsv0.AlertsV2UniqueCountCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model *logsv0.AlertsV2LessThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := DataSourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionImmediateToMap(model *logsv0.AlertsV2AlertConditionConditionImmediate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Immediate != nil {
		immediateMap, err := DataSourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
		if err != nil {
			return modelMap, err
		}
		modelMap["immediate"] = []map[string]interface{}{immediateMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanToMap(model *logsv0.AlertsV2AlertConditionConditionLessThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThan != nil {
		lessThanMap, err := DataSourceIbmLogsAlertAlertsV2LessThanConditionToMap(model.LessThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than"] = []map[string]interface{}{lessThanMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThan != nil {
		moreThanMap, err := DataSourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model.MoreThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than"] = []map[string]interface{}{moreThanMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThanUsual != nil {
		moreThanUsualMap, err := DataSourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionNewValueToMap(model *logsv0.AlertsV2AlertConditionConditionNewValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewValue != nil {
		newValueMap, err := DataSourceIbmLogsAlertAlertsV2NewValueConditionToMap(model.NewValue)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_value"] = []map[string]interface{}{newValueMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionFlowToMap(model *logsv0.AlertsV2AlertConditionConditionFlow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Flow != nil {
		flowMap, err := DataSourceIbmLogsAlertAlertsV2FlowConditionToMap(model.Flow)
		if err != nil {
			return modelMap, err
		}
		modelMap["flow"] = []map[string]interface{}{flowMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionUniqueCountToMap(model *logsv0.AlertsV2AlertConditionConditionUniqueCount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UniqueCount != nil {
		uniqueCountMap, err := DataSourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model.UniqueCount)
		if err != nil {
			return modelMap, err
		}
		modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionLessThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThanUsual != nil {
		lessThanUsualMap, err := DataSourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than_usual"] = []map[string]interface{}{lessThanUsualMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertNotificationGroupsToMap(model *logsv0.AlertsV2AlertNotificationGroups) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GroupByFields != nil {
		modelMap["group_by_fields"] = model.GroupByFields
	}
	if model.Notifications != nil {
		notifications := []map[string]interface{}{}
		for _, notificationsItem := range model.Notifications {
			notificationsItemMap, err := DataSourceIbmLogsAlertAlertsV2AlertNotificationToMap(notificationsItem)
			if err != nil {
				return modelMap, err
			}
			notifications = append(notifications, notificationsItemMap)
		}
		modelMap["notifications"] = notifications
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertNotificationToMap(model logsv0.AlertsV2AlertNotificationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID))
	} else if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients); ok {
		return DataSourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients))
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
			recipientsMap, err := DataSourceIbmLogsAlertAlertsV2RecipientsToMap(model.Recipients)
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

func DataSourceIbmLogsAlertAlertsV2RecipientsToMap(model *logsv0.AlertsV2Recipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Emails != nil {
		modelMap["emails"] = model.Emails
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RetriggeringPeriodSeconds != nil {
		modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
	}
	if model.NotifyOn != nil {
		modelMap["notify_on"] = *model.NotifyOn
	}
	if model.Recipients != nil {
		recipientsMap, err := DataSourceIbmLogsAlertAlertsV2RecipientsToMap(model.Recipients)
		if err != nil {
			return modelMap, err
		}
		modelMap["recipients"] = []map[string]interface{}{recipientsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1AlertFiltersToMap(model *logsv0.AlertsV1AlertFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	if model.Metadata != nil {
		metadataMap, err := DataSourceIbmLogsAlertAlertsV1AlertFiltersMetadataFiltersToMap(model.Metadata)
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
			ratioAlertsItemMap, err := DataSourceIbmLogsAlertAlertsV1AlertFiltersRatioAlertToMap(&ratioAlertsItem)
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

func DataSourceIbmLogsAlertAlertsV1AlertFiltersMetadataFiltersToMap(model *logsv0.AlertsV1AlertFiltersMetadataFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	if model.Subsystems != nil {
		modelMap["subsystems"] = model.Subsystems
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1AlertFiltersRatioAlertToMap(model *logsv0.AlertsV1AlertFiltersRatioAlert) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV1AlertActiveWhenToMap(model *logsv0.AlertsV1AlertActiveWhen) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	timeframes := []map[string]interface{}{}
	for _, timeframesItem := range model.Timeframes {
		timeframesItemMap, err := DataSourceIbmLogsAlertAlertsV1AlertActiveTimeframeToMap(&timeframesItem)
		if err != nil {
			return modelMap, err
		}
		timeframes = append(timeframes, timeframesItemMap)
	}
	modelMap["timeframes"] = timeframes
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1AlertActiveTimeframeToMap(model *logsv0.AlertsV1AlertActiveTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["days_of_week"] = model.DaysOfWeek
	rangeVarMap, err := DataSourceIbmLogsAlertAlertsV1TimeRangeToMap(model.Range)
	if err != nil {
		return modelMap, err
	}
	modelMap["range"] = []map[string]interface{}{rangeVarMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1TimeRangeToMap(model *logsv0.AlertsV1TimeRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startMap, err := DataSourceIbmLogsAlertAlertsV1TimeToMap(model.Start)
	if err != nil {
		return modelMap, err
	}
	modelMap["start"] = []map[string]interface{}{startMap}
	endMap, err := DataSourceIbmLogsAlertAlertsV1TimeToMap(model.End)
	if err != nil {
		return modelMap, err
	}
	modelMap["end"] = []map[string]interface{}{endMap}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV1TimeToMap(model *logsv0.AlertsV1Time) (map[string]interface{}, error) {
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

func DataSourceIbmLogsAlertAlertsV1MetaLabelToMap(model *logsv0.AlertsV1MetaLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmLogsAlertAlertsV2AlertIncidentSettingsToMap(model *logsv0.AlertsV2AlertIncidentSettings) (map[string]interface{}, error) {
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
