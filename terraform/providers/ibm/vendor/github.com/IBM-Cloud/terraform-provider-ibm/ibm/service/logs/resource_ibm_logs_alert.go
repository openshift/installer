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

func ResourceIbmLogsAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsAlertCreate,
		ReadContext:   resourceIbmLogsAlertRead,
		UpdateContext: resourceIbmLogsAlertUpdate,
		DeleteContext: resourceIbmLogsAlertDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_alert", "name"),
				Description:  "Alert name.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_alert", "description"),
				Description:  "Alert description.",
			},
			"is_active": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Alert is active.",
			},
			"severity": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_alert", "severity"),
				Description:  "Alert severity.",
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Alert expiration date.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"year": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Year.",
						},
						"month": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Month of the year.",
						},
						"day": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Day of the month.",
						},
					},
				},
			},
			"condition": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Alert condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"immediate": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for immediate standard alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"less_than": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for less than alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for more than alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "The evaluation window for the alert condition.",
									},
								},
							},
						},
						"more_than_usual": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for more than usual alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for new value alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for flow alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stages": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The Flow alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"groups": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of groups of alerts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"alerts": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "List of alerts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"op": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Operator for the alerts.",
																		},
																		"values": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of alerts.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"id": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The alert ID.",
																					},
																					"not": &schema.Schema{
																						Type:        schema.TypeBool,
																						Optional:    true,
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
																Optional:    true,
																Description: "Operator for the alerts.",
															},
														},
													},
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Timeframe for the flow.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ms": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Should suppression be enforced on the flow alert.",
									},
								},
							},
						},
						"unique_count": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for unique count alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Condition for less than usual alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameters": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The Less than alert condition parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"threshold": &schema.Schema{
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "The threshold for the alert condition.",
												},
												"timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The timeframe for the alert condition.",
												},
												"group_by": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The group by fields for the alert condition.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"metric_alert_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The lucene metric alert parameters if it is a lucene metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_field": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric field of the metric alert.",
															},
															"metric_source": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The metric source of the metric alert.",
															},
															"arithmetic_operator": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator modifier of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"metric_alert_promql_parameters": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The promql metric alert parameters if is is a promql metric alert.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"promql_text": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "The promql text of the metric alert by fields for the alert condition.",
															},
															"arithmetic_operator_modifier": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The arithmetic operator of the metric promql alert.",
															},
															"sample_threshold_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The threshold percentage.",
															},
															"non_null_percentage": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Non null percentage of the evaluation.",
															},
															"swap_null_values": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Should we swap null values with zero.",
															},
														},
													},
												},
												"ignore_infinity": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Should the evaluation ignore infinity value.",
												},
												"relative_timeframe": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The relative timeframe for time relative alerts.",
												},
												"cardinality_fields": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Cardinality fields for unique count alert.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"related_extended_data": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Deadman configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cleanup_deadman_duration": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Cleanup deadman duration.",
															},
															"should_trigger_deadman": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
				Optional:    true,
				Description: "Alert notification groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_by_fields": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Group by fields to group the values by.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"notifications": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Webhook target settings for the the notification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retriggering_period_seconds": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Retriggering period of the alert in seconds.",
									},
									"notify_on": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Notify on setting.",
									},
									"integration_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Integration ID.",
									},
									"recipients": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Recipients.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"emails": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Email addresses.",
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
			"filters": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Alert filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severities": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The severity of the logs to filter.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The metadata filters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"applications": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The applications to filter.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"subsystems": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The subsystems to filter.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The alias of the filter.",
						},
						"text": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The text to filter.",
						},
						"ratio_alerts": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The ratio alerts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The alias of the filter.",
									},
									"text": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The text to filter.",
									},
									"severities": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The severities to filter.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"applications": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The applications to filter.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"subsystems": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The subsystems to filter.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The group by fields.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"filter_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the filter.",
						},
					},
				},
			},
			"active_when": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "When should the alert be active.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeframes": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Activity timeframes of the alert.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_of_week": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "Days of the week for activity.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"range": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Time range in the day of the week.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Start time.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"hours": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Hours of the day.",
															},
															"minutes": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minutes of the hour.",
															},
															"seconds": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Seconds of the minute.",
															},
														},
													},
												},
												"end": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Start time.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"hours": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Hours of the day.",
															},
															"minutes": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minutes of the hour.",
															},
															"seconds": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
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
				Optional:    true,
				Description: "JSON keys to include in the alert notification, if left empty get the full log text in the alert notification.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"meta_labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Meta labels to add to the alert.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key of the label.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the label.",
						},
					},
				},
			},
			"meta_labels_strings": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Meta labels to add to the alert as string with ':' separator.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"incident_settings": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Incident settings, will create the incident based on this configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retriggering_period_seconds": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The retriggering period of the alert in seconds.",
						},
						"notify_on": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notify on setting.",
						},
						"use_as_notification_settings": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use these settings for all notificaion webhook.",
						},
					},
				},
			},
			"unique_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert unique identifier.",
			},
			"alert_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alert Id.",
			},
		},
	}
}

func ResourceIbmLogsAlertValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9_\-\s]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "severity",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "critical, error, info_or_unspecified, warning",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_alert", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsAlertCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_alert", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createAlertOptions := &logsv0.CreateAlertOptions{}

	createAlertOptions.SetName(d.Get("name").(string))
	createAlertOptions.SetIsActive(d.Get("is_active").(bool))
	createAlertOptions.SetSeverity(d.Get("severity").(string))
	conditionModel, err := ResourceIbmLogsAlertMapToAlertsV2AlertCondition(d.Get("condition.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createAlertOptions.SetCondition(conditionModel)
	var notificationGroups []logsv0.AlertsV2AlertNotificationGroups
	for _, v := range d.Get("notification_groups").([]interface{}) {
		if v != nil {
			value := v.(map[string]interface{})
			notificationGroupsItem, err := ResourceIbmLogsAlertMapToAlertsV2AlertNotificationGroups(value)
			if err != nil {
				return diag.FromErr(err)
			}
			notificationGroups = append(notificationGroups, *notificationGroupsItem)
		}
	}
	createAlertOptions.SetNotificationGroups(notificationGroups)
	filtersModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertFilters(d.Get("filters.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createAlertOptions.SetFilters(filtersModel)
	if _, ok := d.GetOk("description"); ok {
		createAlertOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("expiration"); ok {
		expirationModel, err := ResourceIbmLogsAlertMapToAlertsV1Date(d.Get("expiration.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createAlertOptions.SetExpiration(expirationModel)
	}
	if _, ok := d.GetOk("active_when"); ok {
		activeWhenModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertActiveWhen(d.Get("active_when.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createAlertOptions.SetActiveWhen(activeWhenModel)
	}
	if _, ok := d.GetOk("notification_payload_filters"); ok {
		var notificationPayloadFilters []string
		for _, v := range d.Get("notification_payload_filters").([]interface{}) {
			notificationPayloadFiltersItem := v.(string)
			notificationPayloadFilters = append(notificationPayloadFilters, notificationPayloadFiltersItem)
		}
		createAlertOptions.SetNotificationPayloadFilters(notificationPayloadFilters)
	}
	if _, ok := d.GetOk("meta_labels"); ok {
		var metaLabels []logsv0.AlertsV1MetaLabel
		for _, v := range d.Get("meta_labels").([]interface{}) {
			value := v.(map[string]interface{})
			metaLabelsItem, err := ResourceIbmLogsAlertMapToAlertsV1MetaLabel(value)
			if err != nil {
				return diag.FromErr(err)
			}
			metaLabels = append(metaLabels, *metaLabelsItem)
		}
		createAlertOptions.SetMetaLabels(metaLabels)
	}
	if _, ok := d.GetOk("meta_labels_strings"); ok {
		var metaLabelsStrings []string
		for _, v := range d.Get("meta_labels_strings").([]interface{}) {
			metaLabelsStringsItem := v.(string)
			metaLabelsStrings = append(metaLabelsStrings, metaLabelsStringsItem)
		}
		createAlertOptions.SetMetaLabelsStrings(metaLabelsStrings)
	}
	if _, ok := d.GetOk("incident_settings"); ok {
		incidentSettingsModel, err := ResourceIbmLogsAlertMapToAlertsV2AlertIncidentSettings(d.Get("incident_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createAlertOptions.SetIncidentSettings(incidentSettingsModel)
	}

	alert, _, err := logsClient.CreateAlertWithContext(context, createAlertOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateAlertWithContext failed: %s", err.Error()), "ibm_logs_alert", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	alertId := fmt.Sprintf("%s/%s/%s", region, instanceId, *alert.ID)
	d.SetId(alertId)

	return resourceIbmLogsAlertRead(context, d, meta)
}

func resourceIbmLogsAlertRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_alert", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, region, instanceId, alertId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getAlertOptions := &logsv0.GetAlertOptions{}

	getAlertOptions.SetID(core.UUIDPtr(strfmt.UUID(alertId)))

	alert, response, err := logsClient.GetAlertWithContext(context, getAlertOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAlertWithContext failed: %s", err.Error()), "ibm_logs_alert", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("alert_id", alertId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting alert_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("name", alert.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(alert.Description) {
		if err = d.Set("description", alert.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if err = d.Set("is_active", alert.IsActive); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_active: %s", err))
	}
	if err = d.Set("severity", alert.Severity); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting severity: %s", err))
	}
	if !core.IsNil(alert.Expiration) {
		expirationMap, err := ResourceIbmLogsAlertAlertsV1DateToMap(alert.Expiration)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("expiration", []map[string]interface{}{expirationMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting expiration: %s", err))
		}
	}
	conditionMap, err := ResourceIbmLogsAlertAlertsV2AlertConditionToMap(alert.Condition)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("condition", []map[string]interface{}{conditionMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting condition: %s", err))
	}
	notificationGroups := []map[string]interface{}{}
	for _, notificationGroupsItem := range alert.NotificationGroups {
		notificationGroupsItemMap, err := ResourceIbmLogsAlertAlertsV2AlertNotificationGroupsToMap(&notificationGroupsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		notificationGroups = append(notificationGroups, notificationGroupsItemMap)
	}
	if err = d.Set("notification_groups", notificationGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting notification_groups: %s", err))
	}
	if alert.Filters != nil {
		filtersMap, err := ResourceIbmLogsAlertAlertsV1AlertFiltersToMap(alert.Filters)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("filters", []map[string]interface{}{filtersMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting filters: %s", err))
		}
	}

	if !core.IsNil(alert.ActiveWhen) {
		activeWhenMap, err := ResourceIbmLogsAlertAlertsV1AlertActiveWhenToMap(alert.ActiveWhen)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("active_when", []map[string]interface{}{activeWhenMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting active_when: %s", err))
		}
	}
	if !core.IsNil(alert.NotificationPayloadFilters) {
		if err = d.Set("notification_payload_filters", alert.NotificationPayloadFilters); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting notification_payload_filters: %s", err))
		}
	}
	if !core.IsNil(alert.MetaLabels) {
		metaLabels := []map[string]interface{}{}
		for _, metaLabelsItem := range alert.MetaLabels {
			metaLabelsItemMap, err := ResourceIbmLogsAlertAlertsV1MetaLabelToMap(&metaLabelsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			metaLabels = append(metaLabels, metaLabelsItemMap)
		}
		if err = d.Set("meta_labels", metaLabels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting meta_labels: %s", err))
		}
	}
	if !core.IsNil(alert.MetaLabelsStrings) {
		if err = d.Set("meta_labels_strings", alert.MetaLabelsStrings); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting meta_labels_strings: %s", err))
		}
	}
	if !core.IsNil(alert.IncidentSettings) {
		incidentSettingsMap, err := ResourceIbmLogsAlertAlertsV2AlertIncidentSettingsToMap(alert.IncidentSettings)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("incident_settings", []map[string]interface{}{incidentSettingsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting incident_settings: %s", err))
		}
	}
	if !core.IsNil(alert.UniqueIdentifier) {
		if err = d.Set("unique_identifier", alert.UniqueIdentifier); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting unique_identifier: %s", err))
		}
	}

	return nil
}

func resourceIbmLogsAlertUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_alert", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, alertId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updateAlertOptions := &logsv0.UpdateAlertOptions{}

	updateAlertOptions.SetID(core.UUIDPtr(strfmt.UUID(alertId)))

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("is_active") ||
		d.HasChange("severity") ||
		d.HasChange("condition") ||
		d.HasChange("notification_groups") ||
		d.HasChange("filters") ||
		d.HasChange("description") ||
		d.HasChange("expiration") ||
		d.HasChange("active_when") ||
		d.HasChange("notification_payload_filters") ||
		d.HasChange("meta_labels") ||
		d.HasChange("meta_labels_strings") ||
		d.HasChange("incident_settings") {

		updateAlertOptions.SetName(d.Get("name").(string))
		updateAlertOptions.SetIsActive(d.Get("is_active").(bool))
		updateAlertOptions.SetSeverity(d.Get("severity").(string))

		conditionModel, err := ResourceIbmLogsAlertMapToAlertsV2AlertCondition(d.Get("condition.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateAlertOptions.SetCondition(conditionModel)

		var notificationGroups []logsv0.AlertsV2AlertNotificationGroups
		for _, v := range d.Get("notification_groups").([]interface{}) {
			if v != nil {
				value := v.(map[string]interface{})
				notificationGroupsItem, err := ResourceIbmLogsAlertMapToAlertsV2AlertNotificationGroups(value)
				if err != nil {
					return diag.FromErr(err)
				}
				notificationGroups = append(notificationGroups, *notificationGroupsItem)
			}
		}
		updateAlertOptions.SetNotificationGroups(notificationGroups)

		filtersModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertFilters(d.Get("filters.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateAlertOptions.SetFilters(filtersModel)

		if _, ok := d.GetOk("description"); ok {
			updateAlertOptions.SetDescription(d.Get("description").(string))
		}

		if _, ok := d.GetOk("expiration"); ok {
			expirationModel, err := ResourceIbmLogsAlertMapToAlertsV1Date(d.Get("expiration.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			updateAlertOptions.SetExpiration(expirationModel)
		}
		if _, ok := d.GetOk("active_when"); ok {
			activeWhenModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertActiveWhen(d.Get("active_when.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			updateAlertOptions.SetActiveWhen(activeWhenModel)
		}
		if _, ok := d.GetOk("notification_payload_filters"); ok {
			var notificationPayloadFilters []string
			for _, v := range d.Get("notification_payload_filters").([]interface{}) {
				notificationPayloadFiltersItem := v.(string)
				notificationPayloadFilters = append(notificationPayloadFilters, notificationPayloadFiltersItem)
			}
			updateAlertOptions.SetNotificationPayloadFilters(notificationPayloadFilters)
		}

		if _, ok := d.GetOk("meta_labels"); ok {
			var metaLabels []logsv0.AlertsV1MetaLabel
			for _, v := range d.Get("meta_labels").([]interface{}) {
				value := v.(map[string]interface{})
				metaLabelsItem, err := ResourceIbmLogsAlertMapToAlertsV1MetaLabel(value)
				if err != nil {
					return diag.FromErr(err)
				}
				metaLabels = append(metaLabels, *metaLabelsItem)
			}
			updateAlertOptions.SetMetaLabels(metaLabels)
		}
		if _, ok := d.GetOk("meta_labels_strings"); ok {
			var metaLabelsStrings []string
			for _, v := range d.Get("meta_labels_strings").([]interface{}) {
				metaLabelsStringsItem := v.(string)
				metaLabelsStrings = append(metaLabelsStrings, metaLabelsStringsItem)
			}
			updateAlertOptions.SetMetaLabelsStrings(metaLabelsStrings)
		}
		if _, ok := d.GetOk("incident_settings"); ok {
			incidentSettingsModel, err := ResourceIbmLogsAlertMapToAlertsV2AlertIncidentSettings(d.Get("incident_settings.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			updateAlertOptions.SetIncidentSettings(incidentSettingsModel)
		}
		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdateAlertWithContext(context, updateAlertOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateAlertWithContext failed: %s", err.Error()), "ibm_logs_alert", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsAlertRead(context, d, meta)
}

func resourceIbmLogsAlertDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_alert", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, alertId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAlertOptions := &logsv0.DeleteAlertOptions{}

	deleteAlertOptions.SetID(core.UUIDPtr(strfmt.UUID(alertId)))

	_, err = logsClient.DeleteAlertWithContext(context, deleteAlertOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteAlertWithContext failed: %s", err.Error()), "ibm_logs_alert", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertCondition(modelMap map[string]interface{}) (logsv0.AlertsV2AlertConditionIntf, error) {
	model := &logsv0.AlertsV2AlertCondition{}
	if modelMap["immediate"] != nil && len(modelMap["immediate"].([]interface{})) > 0 {
		ImmediateModel, err := ResourceIbmLogsAlertMapToAlertsV2ImmediateConditionEmpty(modelMap["immediate"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Immediate = ImmediateModel
	}
	if modelMap["less_than"] != nil && len(modelMap["less_than"].([]interface{})) > 0 {
		LessThanModel, err := ResourceIbmLogsAlertMapToAlertsV2LessThanCondition(modelMap["less_than"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LessThan = LessThanModel
	}
	if modelMap["more_than"] != nil && len(modelMap["more_than"].([]interface{})) > 0 {
		MoreThanModel, err := ResourceIbmLogsAlertMapToAlertsV2MoreThanCondition(modelMap["more_than"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MoreThan = MoreThanModel
	}
	if modelMap["more_than_usual"] != nil && len(modelMap["more_than_usual"].([]interface{})) > 0 {
		MoreThanUsualModel, err := ResourceIbmLogsAlertMapToAlertsV2MoreThanUsualCondition(modelMap["more_than_usual"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MoreThanUsual = MoreThanUsualModel
	}
	if modelMap["new_value"] != nil && len(modelMap["new_value"].([]interface{})) > 0 {
		NewValueModel, err := ResourceIbmLogsAlertMapToAlertsV2NewValueCondition(modelMap["new_value"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.NewValue = NewValueModel
	}
	if modelMap["flow"] != nil && len(modelMap["flow"].([]interface{})) > 0 {
		FlowModel, err := ResourceIbmLogsAlertMapToAlertsV2FlowCondition(modelMap["flow"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Flow = FlowModel
	}
	if modelMap["unique_count"] != nil && len(modelMap["unique_count"].([]interface{})) > 0 {
		UniqueCountModel, err := ResourceIbmLogsAlertMapToAlertsV2UniqueCountCondition(modelMap["unique_count"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.UniqueCount = UniqueCountModel
	}
	if modelMap["less_than_usual"] != nil && len(modelMap["less_than_usual"].([]interface{})) > 0 {
		LessThanUsualModel, err := ResourceIbmLogsAlertMapToAlertsV2LessThanUsualCondition(modelMap["less_than_usual"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LessThanUsual = LessThanUsualModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2ImmediateConditionEmpty(modelMap interface{}) (*logsv0.AlertsV2ImmediateConditionEmpty, error) {
	model := &logsv0.AlertsV2ImmediateConditionEmpty{}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2LessThanCondition(modelMap []interface{}) (*logsv0.AlertsV2LessThanCondition, error) {
	model := &logsv0.AlertsV2LessThanCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMap []interface{}) (*logsv0.AlertsV2ConditionParameters, error) {
	model := &logsv0.AlertsV2ConditionParameters{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.Threshold = core.Float64Ptr(modelMapElement["threshold"].(float64))
		model.Timeframe = core.StringPtr(modelMapElement["timeframe"].(string))
		if modelMapElement["group_by"] != nil {
			groupBy := []string{}
			for _, groupByItem := range modelMapElement["group_by"].([]interface{}) {
				groupBy = append(groupBy, groupByItem.(string))
			}
			model.GroupBy = groupBy
		}
		if modelMapElement["metric_alert_parameters"] != nil && len(modelMapElement["metric_alert_parameters"].([]interface{})) > 0 {
			MetricAlertParametersModel, err := ResourceIbmLogsAlertMapToAlertsV1MetricAlertConditionParameters(modelMapElement["metric_alert_parameters"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.MetricAlertParameters = MetricAlertParametersModel
		}
		if modelMapElement["metric_alert_promql_parameters"] != nil && len(modelMapElement["metric_alert_promql_parameters"].([]interface{})) > 0 {
			MetricAlertPromqlParametersModel, err := ResourceIbmLogsAlertMapToAlertsV1MetricAlertPromqlConditionParameters(modelMapElement["metric_alert_promql_parameters"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.MetricAlertPromqlParameters = MetricAlertPromqlParametersModel
		}
		if modelMapElement["ignore_infinity"] != nil {
			model.IgnoreInfinity = core.BoolPtr(modelMapElement["ignore_infinity"].(bool))
		}
		if modelMapElement["relative_timeframe"] != nil && modelMapElement["relative_timeframe"].(string) != "" {
			model.RelativeTimeframe = core.StringPtr(modelMapElement["relative_timeframe"].(string))
		}
		if modelMapElement["cardinality_fields"] != nil {
			cardinalityFields := []string{}
			for _, cardinalityFieldsItem := range modelMapElement["cardinality_fields"].([]interface{}) {
				cardinalityFields = append(cardinalityFields, cardinalityFieldsItem.(string))
			}
			model.CardinalityFields = cardinalityFields
		}
		if modelMapElement["related_extended_data"] != nil && len(modelMapElement["related_extended_data"].([]interface{})) > 0 {
			RelatedExtendedDataModel, err := ResourceIbmLogsAlertMapToAlertsV1RelatedExtendedData(modelMapElement["related_extended_data"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.RelatedExtendedData = RelatedExtendedDataModel
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1MetricAlertConditionParameters(modelMap []interface{}) (*logsv0.AlertsV1MetricAlertConditionParameters, error) {
	model := &logsv0.AlertsV1MetricAlertConditionParameters{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.MetricField = core.StringPtr(modelMapElement["metric_field"].(string))
		model.MetricSource = core.StringPtr(modelMapElement["metric_source"].(string))
		model.ArithmeticOperator = core.StringPtr(modelMapElement["arithmetic_operator"].(string))
		if modelMapElement["arithmetic_operator_modifier"] != nil {
			model.ArithmeticOperatorModifier = core.Int64Ptr(int64(modelMapElement["arithmetic_operator_modifier"].(int)))
		}
		if modelMapElement["sample_threshold_percentage"] != nil {
			model.SampleThresholdPercentage = core.Int64Ptr(int64(modelMapElement["sample_threshold_percentage"].(int)))
		}
		if modelMapElement["non_null_percentage"] != nil {
			model.NonNullPercentage = core.Int64Ptr(int64(modelMapElement["non_null_percentage"].(int)))
		}
		if modelMapElement["swap_null_values"] != nil {
			model.SwapNullValues = core.BoolPtr(modelMapElement["swap_null_values"].(bool))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1MetricAlertPromqlConditionParameters(modelMap []interface{}) (*logsv0.AlertsV1MetricAlertPromqlConditionParameters, error) {
	model := &logsv0.AlertsV1MetricAlertPromqlConditionParameters{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		model.PromqlText = core.StringPtr(modelMapElement["promql_text"].(string))
		if modelMapElement["arithmetic_operator_modifier"] != nil {
			model.ArithmeticOperatorModifier = core.Int64Ptr(int64(modelMapElement["arithmetic_operator_modifier"].(int)))
		}
		model.SampleThresholdPercentage = core.Int64Ptr(int64(modelMapElement["sample_threshold_percentage"].(int)))
		if modelMapElement["non_null_percentage"] != nil {
			model.NonNullPercentage = core.Int64Ptr(int64(modelMapElement["non_null_percentage"].(int)))
		}
		if modelMapElement["swap_null_values"] != nil {
			model.SwapNullValues = core.BoolPtr(modelMapElement["swap_null_values"].(bool))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1RelatedExtendedData(modelMap []interface{}) (*logsv0.AlertsV1RelatedExtendedData, error) {
	model := &logsv0.AlertsV1RelatedExtendedData{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["cleanup_deadman_duration"] != nil && modelMapElement["cleanup_deadman_duration"].(string) != "" {
			model.CleanupDeadmanDuration = core.StringPtr(modelMapElement["cleanup_deadman_duration"].(string))
		}
		if modelMapElement["should_trigger_deadman"] != nil {
			model.ShouldTriggerDeadman = core.BoolPtr(modelMapElement["should_trigger_deadman"].(bool))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2MoreThanCondition(modelMap []interface{}) (*logsv0.AlertsV2MoreThanCondition, error) {
	model := &logsv0.AlertsV2MoreThanCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
		if modelMapElement["evaluation_window"] != nil && modelMapElement["evaluation_window"].(string) != "" {
			model.EvaluationWindow = core.StringPtr(modelMapElement["evaluation_window"].(string))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2MoreThanUsualCondition(modelMap []interface{}) (*logsv0.AlertsV2MoreThanUsualCondition, error) {
	model := &logsv0.AlertsV2MoreThanUsualCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2NewValueCondition(modelMap []interface{}) (*logsv0.AlertsV2NewValueCondition, error) {
	model := &logsv0.AlertsV2NewValueCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2FlowCondition(modelMap []interface{}) (*logsv0.AlertsV2FlowCondition, error) {
	model := &logsv0.AlertsV2FlowCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["stages"] != nil {
			stages := []logsv0.AlertsV1FlowStage{}
			for _, stagesItem := range modelMapElement["stages"].([]interface{}) {
				stagesItemModel, err := ResourceIbmLogsAlertMapToAlertsV1FlowStage(stagesItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				stages = append(stages, *stagesItemModel)
			}
			model.Stages = stages
		}
		if modelMapElement["parameters"] != nil && len(modelMapElement["parameters"].([]interface{})) > 0 {
			ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
			if err != nil {
				return model, err
			}
			model.Parameters = ParametersModel
		}
		if modelMapElement["enforce_suppression"] != nil {
			model.EnforceSuppression = core.BoolPtr(modelMapElement["enforce_suppression"].(bool))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1FlowStage(modelMap map[string]interface{}) (*logsv0.AlertsV1FlowStage, error) {
	model := &logsv0.AlertsV1FlowStage{}
	if modelMap["groups"] != nil {
		groups := []logsv0.AlertsV1FlowGroup{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groupsItemModel, err := ResourceIbmLogsAlertMapToAlertsV1FlowGroup(groupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			groups = append(groups, *groupsItemModel)
		}
		model.Groups = groups
	}
	if modelMap["timeframe"] != nil && len(modelMap["timeframe"].([]interface{})) > 0 {
		TimeframeModel, err := ResourceIbmLogsAlertMapToAlertsV1FlowTimeframe(modelMap["timeframe"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Timeframe = TimeframeModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1FlowGroup(modelMap map[string]interface{}) (*logsv0.AlertsV1FlowGroup, error) {
	model := &logsv0.AlertsV1FlowGroup{}
	if modelMap["alerts"] != nil && len(modelMap["alerts"].([]interface{})) > 0 {
		AlertsModel, err := ResourceIbmLogsAlertMapToAlertsV1FlowAlerts(modelMap["alerts"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Alerts = AlertsModel
	}
	if modelMap["next_op"] != nil && modelMap["next_op"].(string) != "" {
		model.NextOp = core.StringPtr(modelMap["next_op"].(string))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1FlowAlerts(modelMap []interface{}) (*logsv0.AlertsV1FlowAlerts, error) {
	model := &logsv0.AlertsV1FlowAlerts{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["op"] != nil && modelMapElement["op"].(string) != "" {
			model.Op = core.StringPtr(modelMapElement["op"].(string))
		}
		if modelMapElement["values"] != nil {
			values := []logsv0.AlertsV1FlowAlert{}
			for _, valuesItem := range modelMapElement["values"].([]interface{}) {
				valuesItemModel, err := ResourceIbmLogsAlertMapToAlertsV1FlowAlert(valuesItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				values = append(values, *valuesItemModel)
			}
			model.Values = values
		}
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1FlowAlert(modelMap map[string]interface{}) (*logsv0.AlertsV1FlowAlert, error) {
	model := &logsv0.AlertsV1FlowAlert{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["not"] != nil {
		model.Not = core.BoolPtr(modelMap["not"].(bool))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1FlowTimeframe(modelMap []interface{}) (*logsv0.AlertsV1FlowTimeframe, error) {
	model := &logsv0.AlertsV1FlowTimeframe{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["ms"] != nil {
			model.Ms = core.Int64Ptr(int64(modelMapElement["ms"].(int)))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2UniqueCountCondition(modelMap []interface{}) (*logsv0.AlertsV2UniqueCountCondition, error) {
	model := &logsv0.AlertsV2UniqueCountCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2LessThanUsualCondition(modelMap []interface{}) (*logsv0.AlertsV2LessThanUsualCondition, error) {
	model := &logsv0.AlertsV2LessThanUsualCondition{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		ParametersModel, err := ResourceIbmLogsAlertMapToAlertsV2ConditionParameters(modelMapElement["parameters"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionImmediate(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionImmediate, error) {
	model := &logsv0.AlertsV2AlertConditionConditionImmediate{}
	if modelMap["immediate"] != nil && len(modelMap["immediate"].([]interface{})) > 0 {
		ImmediateModel, err := ResourceIbmLogsAlertMapToAlertsV2ImmediateConditionEmpty(modelMap["immediate"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Immediate = ImmediateModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionLessThan(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionLessThan, error) {
	model := &logsv0.AlertsV2AlertConditionConditionLessThan{}
	if modelMap["less_than"] != nil && len(modelMap["less_than"].([]interface{})) > 0 {
		LessThanModel, err := ResourceIbmLogsAlertMapToAlertsV2LessThanCondition(modelMap["less_than"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LessThan = LessThanModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionMoreThan(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionMoreThan, error) {
	model := &logsv0.AlertsV2AlertConditionConditionMoreThan{}
	if modelMap["more_than"] != nil && len(modelMap["more_than"].([]interface{})) > 0 {
		MoreThanModel, err := ResourceIbmLogsAlertMapToAlertsV2MoreThanCondition(modelMap["more_than"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MoreThan = MoreThanModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionMoreThanUsual(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionMoreThanUsual, error) {
	model := &logsv0.AlertsV2AlertConditionConditionMoreThanUsual{}
	if modelMap["more_than_usual"] != nil && len(modelMap["more_than_usual"].([]interface{})) > 0 {
		MoreThanUsualModel, err := ResourceIbmLogsAlertMapToAlertsV2MoreThanUsualCondition(modelMap["more_than_usual"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.MoreThanUsual = MoreThanUsualModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionNewValue(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionNewValue, error) {
	model := &logsv0.AlertsV2AlertConditionConditionNewValue{}
	if modelMap["new_value"] != nil && len(modelMap["new_value"].([]interface{})) > 0 {
		NewValueModel, err := ResourceIbmLogsAlertMapToAlertsV2NewValueCondition(modelMap["new_value"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.NewValue = NewValueModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionFlow(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionFlow, error) {
	model := &logsv0.AlertsV2AlertConditionConditionFlow{}
	if modelMap["flow"] != nil && len(modelMap["flow"].([]interface{})) > 0 {
		FlowModel, err := ResourceIbmLogsAlertMapToAlertsV2FlowCondition(modelMap["flow"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Flow = FlowModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionUniqueCount(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionUniqueCount, error) {
	model := &logsv0.AlertsV2AlertConditionConditionUniqueCount{}
	if modelMap["unique_count"] != nil && len(modelMap["unique_count"].([]interface{})) > 0 {
		UniqueCountModel, err := ResourceIbmLogsAlertMapToAlertsV2UniqueCountCondition(modelMap["unique_count"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.UniqueCount = UniqueCountModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertConditionConditionLessThanUsual(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertConditionConditionLessThanUsual, error) {
	model := &logsv0.AlertsV2AlertConditionConditionLessThanUsual{}
	if modelMap["less_than_usual"] != nil && len(modelMap["less_than_usual"].([]interface{})) > 0 {
		LessThanUsualModel, err := ResourceIbmLogsAlertMapToAlertsV2LessThanUsualCondition(modelMap["less_than_usual"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.LessThanUsual = LessThanUsualModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertNotificationGroups(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertNotificationGroups, error) {
	model := &logsv0.AlertsV2AlertNotificationGroups{}
	if modelMap["group_by_fields"] != nil {
		groupByFields := []string{}
		for _, groupByFieldsItem := range modelMap["group_by_fields"].([]interface{}) {
			groupByFields = append(groupByFields, groupByFieldsItem.(string))
		}
		model.GroupByFields = groupByFields
	}
	if modelMap["notifications"] != nil {
		notifications := []logsv0.AlertsV2AlertNotificationIntf{}
		for _, notificationsItem := range modelMap["notifications"].([]interface{}) {
			notificationsItemModel, err := ResourceIbmLogsAlertMapToAlertsV2AlertNotification(notificationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			notifications = append(notifications, notificationsItemModel)
		}
		model.Notifications = notifications
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertNotification(modelMap map[string]interface{}) (logsv0.AlertsV2AlertNotificationIntf, error) {
	model := &logsv0.AlertsV2AlertNotification{}

	if modelMap["retriggering_period_seconds"] != nil && modelMap["retriggering_period_seconds"] != 0 {
		model.RetriggeringPeriodSeconds = core.Int64Ptr(int64(modelMap["retriggering_period_seconds"].(int)))
	}
	if modelMap["notify_on"] != nil && modelMap["notify_on"].(string) != "" {
		model.NotifyOn = core.StringPtr(modelMap["notify_on"].(string))
	}
	if modelMap["integration_id"] != nil {
		model.IntegrationID = core.Int64Ptr(int64(modelMap["integration_id"].(int)))
	}
	if modelMap["recipients"] != nil && len(modelMap["recipients"].([]interface{})) > 0 {
		RecipientsModel, err := ResourceIbmLogsAlertMapToAlertsV2Recipients(modelMap["recipients"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Recipients = RecipientsModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2Recipients(modelMap []interface{}) (*logsv0.AlertsV2Recipients, error) {
	model := &logsv0.AlertsV2Recipients{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["emails"] != nil {
			emails := []string{}
			for _, emailsItem := range modelMapElement["emails"].([]interface{}) {
				emails = append(emails, emailsItem.(string))
			}
			model.Emails = emails
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertNotificationIntegrationTypeIntegrationID(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID, error) {
	model := &logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID{}
	if modelMap["retriggering_period_seconds"] != nil {
		model.RetriggeringPeriodSeconds = core.Int64Ptr(int64(modelMap["retriggering_period_seconds"].(int)))
	}
	if modelMap["notify_on"] != nil && modelMap["notify_on"].(string) != "" {
		model.NotifyOn = core.StringPtr(modelMap["notify_on"].(string))
	}
	if modelMap["integration_id"] != nil {
		model.IntegrationID = core.Int64Ptr(int64(modelMap["integration_id"].(int)))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertNotificationIntegrationTypeRecipients(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients, error) {
	model := &logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients{}
	if modelMap["retriggering_period_seconds"] != nil {
		model.RetriggeringPeriodSeconds = core.Int64Ptr(int64(modelMap["retriggering_period_seconds"].(int)))
	}
	if modelMap["notify_on"] != nil && modelMap["notify_on"].(string) != "" {
		model.NotifyOn = core.StringPtr(modelMap["notify_on"].(string))
	}
	if modelMap["recipients"] != nil && len(modelMap["recipients"].([]interface{})) > 0 {
		RecipientsModel, err := ResourceIbmLogsAlertMapToAlertsV2Recipients(modelMap["recipients"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Recipients = RecipientsModel
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1AlertFilters(modelMap map[string]interface{}) (*logsv0.AlertsV1AlertFilters, error) {
	model := &logsv0.AlertsV1AlertFilters{}
	if modelMap["severities"] != nil {
		severities := []string{}
		for _, severitiesItem := range modelMap["severities"].([]interface{}) {
			severities = append(severities, severitiesItem.(string))
		}
		model.Severities = severities
	}
	if modelMap["metadata"] != nil && len(modelMap["metadata"].([]interface{})) > 0 {
		MetadataModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertFiltersMetadataFilters(modelMap["metadata"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Metadata = MetadataModel
	}
	if modelMap["alias"] != nil && modelMap["alias"].(string) != "" {
		model.Alias = core.StringPtr(modelMap["alias"].(string))
	}
	if modelMap["text"] != nil && modelMap["text"].(string) != "" {
		model.Text = core.StringPtr(modelMap["text"].(string))
	}
	if modelMap["ratio_alerts"] != nil {
		ratioAlerts := []logsv0.AlertsV1AlertFiltersRatioAlert{}
		for _, ratioAlertsItem := range modelMap["ratio_alerts"].([]interface{}) {
			ratioAlertsItemModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertFiltersRatioAlert(ratioAlertsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			ratioAlerts = append(ratioAlerts, *ratioAlertsItemModel)
		}
		model.RatioAlerts = ratioAlerts
	}
	if modelMap["filter_type"] != nil && modelMap["filter_type"].(string) != "" {
		model.FilterType = core.StringPtr(modelMap["filter_type"].(string))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1AlertFiltersMetadataFilters(modelMap []interface{}) (*logsv0.AlertsV1AlertFiltersMetadataFilters, error) {
	model := &logsv0.AlertsV1AlertFiltersMetadataFilters{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["applications"] != nil {
			applications := []string{}
			for _, applicationsItem := range modelMapElement["applications"].([]interface{}) {
				applications = append(applications, applicationsItem.(string))
			}
			model.Applications = applications
		}
		if modelMapElement["subsystems"] != nil {
			subsystems := []string{}
			for _, subsystemsItem := range modelMapElement["subsystems"].([]interface{}) {
				subsystems = append(subsystems, subsystemsItem.(string))
			}
			model.Subsystems = subsystems
		}
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1AlertFiltersRatioAlert(modelMap map[string]interface{}) (*logsv0.AlertsV1AlertFiltersRatioAlert, error) {
	model := &logsv0.AlertsV1AlertFiltersRatioAlert{}
	if modelMap["alias"] != nil && modelMap["alias"].(string) != "" {
		model.Alias = core.StringPtr(modelMap["alias"].(string))
	}
	if modelMap["text"] != nil && modelMap["text"].(string) != "" {
		model.Text = core.StringPtr(modelMap["text"].(string))
	}
	if modelMap["severities"] != nil {
		severities := []string{}
		for _, severitiesItem := range modelMap["severities"].([]interface{}) {
			severities = append(severities, severitiesItem.(string))
		}
		model.Severities = severities
	}
	if modelMap["applications"] != nil {
		applications := []string{}
		for _, applicationsItem := range modelMap["applications"].([]interface{}) {
			applications = append(applications, applicationsItem.(string))
		}
		model.Applications = applications
	}
	if modelMap["subsystems"] != nil {
		subsystems := []string{}
		for _, subsystemsItem := range modelMap["subsystems"].([]interface{}) {
			subsystems = append(subsystems, subsystemsItem.(string))
		}
		model.Subsystems = subsystems
	}
	if modelMap["group_by"] != nil {
		groupBy := []string{}
		for _, groupByItem := range modelMap["group_by"].([]interface{}) {
			groupBy = append(groupBy, groupByItem.(string))
		}
		model.GroupBy = groupBy
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1Date(modelMap map[string]interface{}) (*logsv0.AlertsV1Date, error) {
	model := &logsv0.AlertsV1Date{}
	if modelMap["year"] != nil {
		model.Year = core.Int64Ptr(int64(modelMap["year"].(int)))
	}
	if modelMap["month"] != nil {
		model.Month = core.Int64Ptr(int64(modelMap["month"].(int)))
	}
	if modelMap["day"] != nil {
		model.Day = core.Int64Ptr(int64(modelMap["day"].(int)))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1AlertActiveWhen(modelMap map[string]interface{}) (*logsv0.AlertsV1AlertActiveWhen, error) {
	model := &logsv0.AlertsV1AlertActiveWhen{}
	timeframes := []logsv0.AlertsV1AlertActiveTimeframe{}
	for _, timeframesItem := range modelMap["timeframes"].([]interface{}) {
		timeframesItemModel, err := ResourceIbmLogsAlertMapToAlertsV1AlertActiveTimeframe(timeframesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		timeframes = append(timeframes, *timeframesItemModel)
	}
	model.Timeframes = timeframes
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1AlertActiveTimeframe(modelMap map[string]interface{}) (*logsv0.AlertsV1AlertActiveTimeframe, error) {
	model := &logsv0.AlertsV1AlertActiveTimeframe{}
	daysOfWeek := []string{}
	for _, daysOfWeekItem := range modelMap["days_of_week"].([]interface{}) {
		daysOfWeek = append(daysOfWeek, daysOfWeekItem.(string))
	}
	model.DaysOfWeek = daysOfWeek
	RangeModel, err := ResourceIbmLogsAlertMapToAlertsV1TimeRange(modelMap["range"].([]interface{}))
	if err != nil {
		return model, err
	}
	model.Range = RangeModel
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1TimeRange(modelMap []interface{}) (*logsv0.AlertsV1TimeRange, error) {
	model := &logsv0.AlertsV1TimeRange{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		StartModel, err := ResourceIbmLogsAlertMapToAlertsV1Time(modelMapElement["start"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.Start = StartModel
		EndModel, err := ResourceIbmLogsAlertMapToAlertsV1Time(modelMapElement["end"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.End = EndModel
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1Time(modelMap []interface{}) (*logsv0.AlertsV1Time, error) {
	model := &logsv0.AlertsV1Time{}
	if len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["hours"] != nil {
			model.Hours = core.Int64Ptr(int64(modelMapElement["hours"].(int)))
		}
		if modelMapElement["minutes"] != nil {
			model.Minutes = core.Int64Ptr(int64(modelMapElement["minutes"].(int)))
		}
		if modelMapElement["seconds"] != nil {
			model.Seconds = core.Int64Ptr(int64(modelMapElement["seconds"].(int)))
		}
	}

	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV1MetaLabel(modelMap map[string]interface{}) (*logsv0.AlertsV1MetaLabel, error) {
	model := &logsv0.AlertsV1MetaLabel{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func ResourceIbmLogsAlertMapToAlertsV2AlertIncidentSettings(modelMap map[string]interface{}) (*logsv0.AlertsV2AlertIncidentSettings, error) {
	model := &logsv0.AlertsV2AlertIncidentSettings{}
	if modelMap["retriggering_period_seconds"] != nil {
		model.RetriggeringPeriodSeconds = core.Int64Ptr(int64(modelMap["retriggering_period_seconds"].(int)))
	}
	if modelMap["notify_on"] != nil && modelMap["notify_on"].(string) != "" {
		model.NotifyOn = core.StringPtr(modelMap["notify_on"].(string))
	}
	if modelMap["use_as_notification_settings"] != nil {
		model.UseAsNotificationSettings = core.BoolPtr(modelMap["use_as_notification_settings"].(bool))
	}
	return model, nil
}

func ResourceIbmLogsAlertAlertsV1DateToMap(model *logsv0.AlertsV1Date) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV2AlertConditionToMap(model logsv0.AlertsV2AlertConditionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionImmediate); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionImmediateToMap(model.(*logsv0.AlertsV2AlertConditionConditionImmediate))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThan); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThan); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThan))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionMoreThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionNewValue); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionNewValueToMap(model.(*logsv0.AlertsV2AlertConditionConditionNewValue))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionFlow); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionFlowToMap(model.(*logsv0.AlertsV2AlertConditionConditionFlow))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionUniqueCountToMap(model.(*logsv0.AlertsV2AlertConditionConditionUniqueCount))
	} else if _, ok := model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual); ok {
		return ResourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanUsualToMap(model.(*logsv0.AlertsV2AlertConditionConditionLessThanUsual))
	} else if _, ok := model.(*logsv0.AlertsV2AlertCondition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.AlertsV2AlertCondition)
		if model.Immediate != nil {
			immediateMap, err := ResourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
			if err != nil {
				return modelMap, err
			}
			modelMap["immediate"] = []map[string]interface{}{immediateMap}
		}
		if model.LessThan != nil {
			lessThanMap, err := ResourceIbmLogsAlertAlertsV2LessThanConditionToMap(model.LessThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["less_than"] = []map[string]interface{}{lessThanMap}
		}
		if model.MoreThan != nil {
			moreThanMap, err := ResourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model.MoreThan)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than"] = []map[string]interface{}{moreThanMap}
		}
		if model.MoreThanUsual != nil {
			moreThanUsualMap, err := ResourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
			if err != nil {
				return modelMap, err
			}
			modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
		}
		if model.NewValue != nil {
			newValueMap, err := ResourceIbmLogsAlertAlertsV2NewValueConditionToMap(model.NewValue)
			if err != nil {
				return modelMap, err
			}
			modelMap["new_value"] = []map[string]interface{}{newValueMap}
		}
		if model.Flow != nil {
			flowMap, err := ResourceIbmLogsAlertAlertsV2FlowConditionToMap(model.Flow)
			if err != nil {
				return modelMap, err
			}
			modelMap["flow"] = []map[string]interface{}{flowMap}
		}
		if model.UniqueCount != nil {
			uniqueCountMap, err := ResourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model.UniqueCount)
			if err != nil {
				return modelMap, err
			}
			modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
		}
		if model.LessThanUsual != nil {
			lessThanUsualMap, err := ResourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
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

func ResourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model *logsv0.AlertsV2ImmediateConditionEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2LessThanConditionToMap(model *logsv0.AlertsV2LessThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model *logsv0.AlertsV2ConditionParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["threshold"] = *model.Threshold
	modelMap["timeframe"] = *model.Timeframe
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	if model.MetricAlertParameters != nil {
		metricAlertParametersMap, err := ResourceIbmLogsAlertAlertsV1MetricAlertConditionParametersToMap(model.MetricAlertParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["metric_alert_parameters"] = []map[string]interface{}{metricAlertParametersMap}
	}
	if model.MetricAlertPromqlParameters != nil {
		metricAlertPromqlParametersMap, err := ResourceIbmLogsAlertAlertsV1MetricAlertPromqlConditionParametersToMap(model.MetricAlertPromqlParameters)
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
		relatedExtendedDataMap, err := ResourceIbmLogsAlertAlertsV1RelatedExtendedDataToMap(model.RelatedExtendedData)
		if err != nil {
			return modelMap, err
		}
		modelMap["related_extended_data"] = []map[string]interface{}{relatedExtendedDataMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1MetricAlertConditionParametersToMap(model *logsv0.AlertsV1MetricAlertConditionParameters) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV1MetricAlertPromqlConditionParametersToMap(model *logsv0.AlertsV1MetricAlertPromqlConditionParameters) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV1RelatedExtendedDataToMap(model *logsv0.AlertsV1RelatedExtendedData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CleanupDeadmanDuration != nil {
		modelMap["cleanup_deadman_duration"] = *model.CleanupDeadmanDuration
	}
	if model.ShouldTriggerDeadman != nil {
		modelMap["should_trigger_deadman"] = *model.ShouldTriggerDeadman
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model *logsv0.AlertsV2MoreThanCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	if model.EvaluationWindow != nil {
		modelMap["evaluation_window"] = *model.EvaluationWindow
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model *logsv0.AlertsV2MoreThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2NewValueConditionToMap(model *logsv0.AlertsV2NewValueCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2FlowConditionToMap(model *logsv0.AlertsV2FlowCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Stages != nil {
		stages := []map[string]interface{}{}
		for _, stagesItem := range model.Stages {
			stagesItemMap, err := ResourceIbmLogsAlertAlertsV1FlowStageToMap(&stagesItem)
			if err != nil {
				return modelMap, err
			}
			stages = append(stages, stagesItemMap)
		}
		modelMap["stages"] = stages
	}
	if model.Parameters != nil {
		parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
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

func ResourceIbmLogsAlertAlertsV1FlowStageToMap(model *logsv0.AlertsV1FlowStage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Groups != nil {
		groups := []map[string]interface{}{}
		for _, groupsItem := range model.Groups {
			groupsItemMap, err := ResourceIbmLogsAlertAlertsV1FlowGroupToMap(&groupsItem)
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Timeframe != nil {
		timeframeMap, err := ResourceIbmLogsAlertAlertsV1FlowTimeframeToMap(model.Timeframe)
		if err != nil {
			return modelMap, err
		}
		modelMap["timeframe"] = []map[string]interface{}{timeframeMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1FlowGroupToMap(model *logsv0.AlertsV1FlowGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Alerts != nil {
		alertsMap, err := ResourceIbmLogsAlertAlertsV1FlowAlertsToMap(model.Alerts)
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

func ResourceIbmLogsAlertAlertsV1FlowAlertsToMap(model *logsv0.AlertsV1FlowAlerts) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Op != nil {
		modelMap["op"] = *model.Op
	}
	if model.Values != nil {
		values := []map[string]interface{}{}
		for _, valuesItem := range model.Values {
			valuesItemMap, err := ResourceIbmLogsAlertAlertsV1FlowAlertToMap(&valuesItem)
			if err != nil {
				return modelMap, err
			}
			values = append(values, valuesItemMap)
		}
		modelMap["values"] = values
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1FlowAlertToMap(model *logsv0.AlertsV1FlowAlert) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Not != nil {
		modelMap["not"] = *model.Not
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1FlowTimeframeToMap(model *logsv0.AlertsV1FlowTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Ms != nil {
		modelMap["ms"] = flex.IntValue(model.Ms)
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model *logsv0.AlertsV2UniqueCountCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model *logsv0.AlertsV2LessThanUsualCondition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	parametersMap, err := ResourceIbmLogsAlertAlertsV2ConditionParametersToMap(model.Parameters)
	if err != nil {
		return modelMap, err
	}
	modelMap["parameters"] = []map[string]interface{}{parametersMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionImmediateToMap(model *logsv0.AlertsV2AlertConditionConditionImmediate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Immediate != nil {
		immediateMap, err := ResourceIbmLogsAlertAlertsV2ImmediateConditionEmptyToMap(model.Immediate)
		if err != nil {
			return modelMap, err
		}
		modelMap["immediate"] = []map[string]interface{}{immediateMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanToMap(model *logsv0.AlertsV2AlertConditionConditionLessThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThan != nil {
		lessThanMap, err := ResourceIbmLogsAlertAlertsV2LessThanConditionToMap(model.LessThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than"] = []map[string]interface{}{lessThanMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThan != nil {
		moreThanMap, err := ResourceIbmLogsAlertAlertsV2MoreThanConditionToMap(model.MoreThan)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than"] = []map[string]interface{}{moreThanMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionMoreThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionMoreThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreThanUsual != nil {
		moreThanUsualMap, err := ResourceIbmLogsAlertAlertsV2MoreThanUsualConditionToMap(model.MoreThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["more_than_usual"] = []map[string]interface{}{moreThanUsualMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionNewValueToMap(model *logsv0.AlertsV2AlertConditionConditionNewValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewValue != nil {
		newValueMap, err := ResourceIbmLogsAlertAlertsV2NewValueConditionToMap(model.NewValue)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_value"] = []map[string]interface{}{newValueMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionFlowToMap(model *logsv0.AlertsV2AlertConditionConditionFlow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Flow != nil {
		flowMap, err := ResourceIbmLogsAlertAlertsV2FlowConditionToMap(model.Flow)
		if err != nil {
			return modelMap, err
		}
		modelMap["flow"] = []map[string]interface{}{flowMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionUniqueCountToMap(model *logsv0.AlertsV2AlertConditionConditionUniqueCount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UniqueCount != nil {
		uniqueCountMap, err := ResourceIbmLogsAlertAlertsV2UniqueCountConditionToMap(model.UniqueCount)
		if err != nil {
			return modelMap, err
		}
		modelMap["unique_count"] = []map[string]interface{}{uniqueCountMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertConditionConditionLessThanUsualToMap(model *logsv0.AlertsV2AlertConditionConditionLessThanUsual) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LessThanUsual != nil {
		lessThanUsualMap, err := ResourceIbmLogsAlertAlertsV2LessThanUsualConditionToMap(model.LessThanUsual)
		if err != nil {
			return modelMap, err
		}
		modelMap["less_than_usual"] = []map[string]interface{}{lessThanUsualMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertNotificationGroupsToMap(model *logsv0.AlertsV2AlertNotificationGroups) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GroupByFields != nil {
		modelMap["group_by_fields"] = model.GroupByFields
	}
	if model.Notifications != nil {
		notifications := []map[string]interface{}{}
		for _, notificationsItem := range model.Notifications {
			notificationsItemMap, err := ResourceIbmLogsAlertAlertsV2AlertNotificationToMap(notificationsItem)
			if err != nil {
				return modelMap, err
			}
			notifications = append(notifications, notificationsItemMap)
		}
		modelMap["notifications"] = notifications
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertNotificationToMap(model logsv0.AlertsV2AlertNotificationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID); ok {
		return ResourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID))
	} else if _, ok := model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients); ok {
		return ResourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model.(*logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients))
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
			recipientsMap, err := ResourceIbmLogsAlertAlertsV2RecipientsToMap(model.Recipients)
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

func ResourceIbmLogsAlertAlertsV2RecipientsToMap(model *logsv0.AlertsV2Recipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Emails != nil {
		modelMap["emails"] = model.Emails
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeIntegrationIDToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeIntegrationID) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV2AlertNotificationIntegrationTypeRecipientsToMap(model *logsv0.AlertsV2AlertNotificationIntegrationTypeRecipients) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RetriggeringPeriodSeconds != nil {
		modelMap["retriggering_period_seconds"] = flex.IntValue(model.RetriggeringPeriodSeconds)
	}
	if model.NotifyOn != nil {
		modelMap["notify_on"] = *model.NotifyOn
	}
	if model.Recipients != nil {
		recipientsMap, err := ResourceIbmLogsAlertAlertsV2RecipientsToMap(model.Recipients)
		if err != nil {
			return modelMap, err
		}
		modelMap["recipients"] = []map[string]interface{}{recipientsMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1AlertFiltersToMap(model *logsv0.AlertsV1AlertFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Severities != nil {
		modelMap["severities"] = model.Severities
	}
	if model.Metadata != nil {
		metadataMap, err := ResourceIbmLogsAlertAlertsV1AlertFiltersMetadataFiltersToMap(model.Metadata)
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
			ratioAlertsItemMap, err := ResourceIbmLogsAlertAlertsV1AlertFiltersRatioAlertToMap(&ratioAlertsItem)
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

func ResourceIbmLogsAlertAlertsV1AlertFiltersMetadataFiltersToMap(model *logsv0.AlertsV1AlertFiltersMetadataFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	if model.Subsystems != nil {
		modelMap["subsystems"] = model.Subsystems
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1AlertFiltersRatioAlertToMap(model *logsv0.AlertsV1AlertFiltersRatioAlert) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV1AlertActiveWhenToMap(model *logsv0.AlertsV1AlertActiveWhen) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	timeframes := []map[string]interface{}{}
	for _, timeframesItem := range model.Timeframes {
		timeframesItemMap, err := ResourceIbmLogsAlertAlertsV1AlertActiveTimeframeToMap(&timeframesItem)
		if err != nil {
			return modelMap, err
		}
		timeframes = append(timeframes, timeframesItemMap)
	}
	modelMap["timeframes"] = timeframes
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1AlertActiveTimeframeToMap(model *logsv0.AlertsV1AlertActiveTimeframe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["days_of_week"] = model.DaysOfWeek
	rangeVarMap, err := ResourceIbmLogsAlertAlertsV1TimeRangeToMap(model.Range)
	if err != nil {
		return modelMap, err
	}
	modelMap["range"] = []map[string]interface{}{rangeVarMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1TimeRangeToMap(model *logsv0.AlertsV1TimeRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	startMap, err := ResourceIbmLogsAlertAlertsV1TimeToMap(model.Start)
	if err != nil {
		return modelMap, err
	}
	modelMap["start"] = []map[string]interface{}{startMap}
	endMap, err := ResourceIbmLogsAlertAlertsV1TimeToMap(model.End)
	if err != nil {
		return modelMap, err
	}
	modelMap["end"] = []map[string]interface{}{endMap}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV1TimeToMap(model *logsv0.AlertsV1Time) (map[string]interface{}, error) {
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

func ResourceIbmLogsAlertAlertsV1MetaLabelToMap(model *logsv0.AlertsV1MetaLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmLogsAlertAlertsV2AlertIncidentSettingsToMap(model *logsv0.AlertsV2AlertIncidentSettings) (map[string]interface{}, error) {
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
