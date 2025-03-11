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

func DataSourceIbmLogsE2m() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsE2mRead,

		Schema: map[string]*schema.Schema{
			"logs_e2m_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of e2m to be deleted.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the E2M.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the E2M.",
			},
			"create_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E2M create time.",
			},
			"update_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E2M update time.",
			},
			"permutations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Represents the limit of the permutations and if the limit was exceeded.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "E2M permutation limit.",
						},
						"has_exceeded_limit": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to indicate if limit was exceeded.",
						},
					},
				},
			},
			"metric_labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "E2M metric labels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_label": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric label target alias name.",
						},
						"source_field": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Metric label source field.",
						},
					},
				},
			},
			"metric_fields": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "E2M metric fields.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_base_metric_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target metric field alias name.",
						},
						"source_field": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source field.",
						},
						"aggregations": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Represents Aggregation type list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is enabled.",
									},
									"agg_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Aggregation type.",
									},
									"target_metric_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target metric field alias name.",
									},
									"samples": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "E2M sample type metadata.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sample_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Sample type min/max.",
												},
											},
										},
									},
									"histogram": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "E2M aggregate histogram type metadata.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"buckets": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Buckets of the E2M.",
													Elem: &schema.Schema{
														Type: schema.TypeFloat,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E2M type.",
			},
			"is_internal": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A flag that represents if the e2m is for internal usage.",
			},
			"logs_query": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "E2M logs query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lucene": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Lucene query.",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias.",
						},
						"applicationname_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Application name filters.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"subsystemname_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subsystem names filters.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"severity_filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Severity type filters.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsE2mRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_e2m", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getE2mOptions := &logsv0.GetE2mOptions{}

	getE2mOptions.SetID(d.Get("logs_e2m_id").(string))

	event2MetricIntf, _, err := logsClient.GetE2mWithContext(context, getE2mOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetE2mWithContext failed: %s", err.Error()), "(Data) ibm_logs_e2m", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	event2Metric := event2MetricIntf.(*logsv0.Event2Metric)

	d.SetId(fmt.Sprintf("%s", *getE2mOptions.ID))

	if err = d.Set("name", event2Metric.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("description", event2Metric.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("create_time", event2Metric.CreateTime); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting create_time: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("update_time", event2Metric.UpdateTime); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting update_time: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	permutations := []map[string]interface{}{}
	if event2Metric.Permutations != nil {
		modelMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2E2mPermutationsToMap(event2Metric.Permutations)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_e2m", "read")
			return tfErr.GetDiag()
		}
		permutations = append(permutations, modelMap)
	}
	if err = d.Set("permutations", permutations); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting permutations: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	metricLabels := []map[string]interface{}{}
	if event2Metric.MetricLabels != nil {
		for _, modelItem := range event2Metric.MetricLabels {
			modelMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2MetricLabelToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_e2m", "read")
				return tfErr.GetDiag()
			}
			metricLabels = append(metricLabels, modelMap)
		}
	}
	if err = d.Set("metric_labels", metricLabels); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metric_labels: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	metricFields := []map[string]interface{}{}
	if event2Metric.MetricFields != nil {
		for _, modelItem := range event2Metric.MetricFields {
			modelMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2MetricFieldToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_e2m", "read")
				return tfErr.GetDiag()
			}
			metricFields = append(metricFields, modelMap)
		}
	}
	if err = d.Set("metric_fields", metricFields); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metric_fields: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("type", event2Metric.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("is_internal", event2Metric.IsInternal); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting is_internal: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	logsQuery := []map[string]interface{}{}
	if event2Metric.LogsQuery != nil {
		modelMap, err := DataSourceIbmLogsE2mApisLogs2metricsV2LogsQueryToMap(event2Metric.LogsQuery)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_e2m", "read")
			return tfErr.GetDiag()
		}
		logsQuery = append(logsQuery, modelMap)
	}
	if err = d.Set("logs_query", logsQuery); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting logs_query: %s", err), "(Data) ibm_logs_e2m", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2E2mPermutationsToMap(model *logsv0.ApisEvents2metricsV2E2mPermutations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Limit != nil {
		modelMap["limit"] = flex.IntValue(model.Limit)
	}
	if model.HasExceededLimit != nil {
		modelMap["has_exceeded_limit"] = *model.HasExceededLimit
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2MetricLabelToMap(model *logsv0.ApisEvents2metricsV2MetricLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetLabel != nil {
		modelMap["target_label"] = *model.TargetLabel
	}
	if model.SourceField != nil {
		modelMap["source_field"] = *model.SourceField
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2MetricFieldToMap(model *logsv0.ApisEvents2metricsV2MetricField) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetBaseMetricName != nil {
		modelMap["target_base_metric_name"] = *model.TargetBaseMetricName
	}
	if model.SourceField != nil {
		modelMap["source_field"] = *model.SourceField
	}
	if model.Aggregations != nil {
		aggregations := []map[string]interface{}{}
		for _, aggregationsItem := range model.Aggregations {
			aggregationsItemMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2AggregationToMap(aggregationsItem)
			if err != nil {
				return modelMap, err
			}
			aggregations = append(aggregations, aggregationsItemMap)
		}
		modelMap["aggregations"] = aggregations
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2AggregationToMap(model logsv0.ApisEvents2metricsV2AggregationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples); ok {
		return DataSourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataSamplesToMap(model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples))
	} else if _, ok := model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram); ok {
		return DataSourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataHistogramToMap(model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram))
	} else if _, ok := model.(*logsv0.ApisEvents2metricsV2Aggregation); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisEvents2metricsV2Aggregation)
		if model.Enabled != nil {
			modelMap["enabled"] = *model.Enabled
		}
		if model.AggType != nil {
			modelMap["agg_type"] = *model.AggType
		}
		if model.TargetMetricName != nil {
			modelMap["target_metric_name"] = *model.TargetMetricName
		}
		if model.Samples != nil {
			samplesMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model.Samples)
			if err != nil {
				return modelMap, err
			}
			modelMap["samples"] = []map[string]interface{}{samplesMap}
		}
		if model.Histogram != nil {
			histogramMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model.Histogram)
			if err != nil {
				return modelMap, err
			}
			modelMap["histogram"] = []map[string]interface{}{histogramMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisEvents2metricsV2AggregationIntf subtype encountered")
	}
}

func DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model *logsv0.ApisEvents2metricsV2E2mAggSamples) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SampleType != nil {
		modelMap["sample_type"] = *model.SampleType
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model *logsv0.ApisEvents2metricsV2E2mAggHistogram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Buckets != nil {
		var buckets []interface{}
		for _, item := range model.Buckets {
			buckets = append(buckets, float64(item))
		}
		modelMap["buckets"] = buckets
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataSamplesToMap(model *logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.AggType != nil {
		modelMap["agg_type"] = *model.AggType
	}
	if model.TargetMetricName != nil {
		modelMap["target_metric_name"] = *model.TargetMetricName
	}
	if model.Samples != nil {
		samplesMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model.Samples)
		if err != nil {
			return modelMap, err
		}
		modelMap["samples"] = []map[string]interface{}{samplesMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataHistogramToMap(model *logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap["enabled"] = *model.Enabled
	}
	if model.AggType != nil {
		modelMap["agg_type"] = *model.AggType
	}
	if model.TargetMetricName != nil {
		modelMap["target_metric_name"] = *model.TargetMetricName
	}
	if model.Histogram != nil {
		histogramMap, err := DataSourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model.Histogram)
		if err != nil {
			return modelMap, err
		}
		modelMap["histogram"] = []map[string]interface{}{histogramMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsE2mApisLogs2metricsV2LogsQueryToMap(model *logsv0.ApisLogs2metricsV2LogsQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Lucene != nil {
		modelMap["lucene"] = *model.Lucene
	}
	if model.Alias != nil {
		modelMap["alias"] = *model.Alias
	}
	if model.ApplicationnameFilters != nil {
		modelMap["applicationname_filters"] = model.ApplicationnameFilters
	}
	if model.SubsystemnameFilters != nil {
		modelMap["subsystemname_filters"] = model.SubsystemnameFilters
	}
	if model.SeverityFilters != nil {
		modelMap["severity_filters"] = model.SeverityFilters
	}
	return modelMap, nil
}
