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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsE2m() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsE2mCreate,
		ReadContext:   resourceIbmLogsE2mRead,
		UpdateContext: resourceIbmLogsE2mUpdate,
		DeleteContext: resourceIbmLogsE2mDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_e2m", "name"),
				Description:  "Name of the E2M.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_e2m", "description"),
				Description:  "Description of the E2M.",
			},
			"metric_labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "E2M metric labels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_label": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metric label target alias name.",
						},
						"source_field": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metric label source field.",
						},
					},
				},
			},
			"metric_fields": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "E2M metric fields.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_base_metric_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target metric field alias name.",
						},
						"source_field": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source field.",
						},
						"aggregations": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Represents Aggregation type list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is enabled.",
									},
									"agg_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Aggregation type.",
									},
									"target_metric_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target metric field alias name.",
									},
									"samples": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "E2M sample type metadata.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sample_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sample type min/max.",
												},
											},
										},
									},
									"histogram": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "E2M aggregate histogram type metadata.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"buckets": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Buckets of the E2M.",
													Elem:        &schema.Schema{Type: schema.TypeFloat},
												},
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_e2m", "type"),
				Description:  "E2M type.",
			},
			"logs_query": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "E2M logs query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lucene": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Lucene query.",
						},
						"alias": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alias.",
						},
						"applicationname_filters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Application name filters.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"subsystemname_filters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Subsystem names filters.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"severity_filters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Severity type filters.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
			"is_internal": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A flag that represents if the e2m is for internal usage.",
			},
			"e2m_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "event to metrics Id.",
			},
		},
	}
}

func ResourceIbmLogsE2mValidator() *validate.ResourceValidator {
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
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "logs2metrics, unspecified",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_e2m", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsE2mCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	bodyModelMap := map[string]interface{}{}
	createE2mOptions := &logsv0.CreateE2mOptions{}

	bodyModelMap["name"] = d.Get("name")
	if _, ok := d.GetOk("description"); ok {
		bodyModelMap["description"] = d.Get("description")
	}
	if _, ok := d.GetOk("permutations_limit"); ok {
		bodyModelMap["permutations_limit"] = d.Get("permutations_limit")
	}
	if _, ok := d.GetOk("metric_labels"); ok {
		bodyModelMap["metric_labels"] = d.Get("metric_labels")
	}
	if _, ok := d.GetOk("metric_fields"); ok {
		bodyModelMap["metric_fields"] = d.Get("metric_fields")
	}
	if _, ok := d.GetOk("type"); ok {
		bodyModelMap["type"] = d.Get("type")
	}
	if _, ok := d.GetOk("logs_query"); ok {
		bodyModelMap["logs_query"] = d.Get("logs_query")
	}
	convertedModel, err := ResourceIbmLogsE2mMapToEvent2MetricPrototype(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "create")
		return tfErr.GetDiag()
	}
	createE2mOptions.Event2MetricPrototype = convertedModel

	event2MetricIntf, _, err := logsClient.CreateE2mWithContext(context, createE2mOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateE2mWithContext failed: %s", err.Error()), "ibm_logs_e2m", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	event2Metric := event2MetricIntf.(*logsv0.Event2Metric)

	event2MetricId := fmt.Sprintf("%s/%s/%s", region, instanceId, *event2Metric.ID)
	d.SetId(event2MetricId)

	return resourceIbmLogsE2mRead(context, d, meta)
}

func resourceIbmLogsE2mRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, e2mId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getE2mOptions := &logsv0.GetE2mOptions{}

	getE2mOptions.SetID(e2mId)

	event2MetricIntf, response, err := logsClient.GetE2mWithContext(context, getE2mOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetE2mWithContext failed: %s", err.Error()), "ibm_logs_e2m", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	event2Metric := event2MetricIntf.(*logsv0.Event2Metric)

	if err = d.Set("e2m_id", e2mId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting e2m_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("name", event2Metric.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(event2Metric.Description) {
		if err = d.Set("description", event2Metric.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(event2Metric.MetricLabels) {
		metricLabels := []map[string]interface{}{}
		for _, metricLabelsItem := range event2Metric.MetricLabels {
			metricLabelsItemMap, err := ResourceIbmLogsE2mApisEvents2metricsV2MetricLabelToMap(&metricLabelsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			metricLabels = append(metricLabels, metricLabelsItemMap)
		}
		if err = d.Set("metric_labels", metricLabels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting metric_labels: %s", err))
		}
	}
	if !core.IsNil(event2Metric.MetricFields) {
		metricFields := []map[string]interface{}{}
		for _, metricFieldsItem := range event2Metric.MetricFields {
			metricFieldsItemMap, err := ResourceIbmLogsE2mApisEvents2metricsV2MetricFieldToMap(&metricFieldsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			metricFields = append(metricFields, metricFieldsItemMap)
		}
		if err = d.Set("metric_fields", metricFields); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting metric_fields: %s", err))
		}
	}
	if !core.IsNil(event2Metric.Type) {
		if err = d.Set("type", event2Metric.Type); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
		}
	}
	if !core.IsNil(event2Metric.LogsQuery) {
		logsQueryMap, err := ResourceIbmLogsE2mApisLogs2metricsV2LogsQueryToMap(event2Metric.LogsQuery)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("logs_query", []map[string]interface{}{logsQueryMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting logs_query: %s", err))
		}
	}
	if !core.IsNil(event2Metric.CreateTime) {
		if err = d.Set("create_time", event2Metric.CreateTime); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting create_time: %s", err))
		}
	}
	if !core.IsNil(event2Metric.UpdateTime) {
		if err = d.Set("update_time", event2Metric.UpdateTime); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting update_time: %s", err))
		}
	}
	if !core.IsNil(event2Metric.Permutations) {
		permutationsMap, err := ResourceIbmLogsE2mApisEvents2metricsV2E2mPermutationsToMap(event2Metric.Permutations)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("permutations", []map[string]interface{}{permutationsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting permutations: %s", err))
		}
	}
	if !core.IsNil(event2Metric.IsInternal) {
		if err = d.Set("is_internal", event2Metric.IsInternal); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting is_internal: %s", err))
		}
	}

	return nil
}

func resourceIbmLogsE2mUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, e2mId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	replaceE2mOptions := &logsv0.ReplaceE2mOptions{}

	replaceE2mOptions.SetID(e2mId)

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("permutations_limit") ||
		d.HasChange("metric_labels") ||
		d.HasChange("metric_fields") ||
		d.HasChange("type") ||
		d.HasChange("logs_query") {

		bodyModelMap := map[string]interface{}{}
		bodyModelMap["name"] = d.Get("name")
		if _, ok := d.GetOk("description"); ok {
			bodyModelMap["description"] = d.Get("description")
		}
		if _, ok := d.GetOk("permutations_limit"); ok {
			bodyModelMap["permutations_limit"] = d.Get("permutations_limit")
		}
		if _, ok := d.GetOk("metric_labels"); ok {
			bodyModelMap["metric_labels"] = d.Get("metric_labels")
		}
		if _, ok := d.GetOk("metric_fields"); ok {
			bodyModelMap["metric_fields"] = d.Get("metric_fields")
		}
		if _, ok := d.GetOk("type"); ok {
			bodyModelMap["type"] = d.Get("type")
		}
		if _, ok := d.GetOk("logs_query"); ok {
			bodyModelMap["logs_query"] = d.Get("logs_query")
		}
		convertedModel, err := ResourceIbmLogsE2mMapToEvent2MetricPrototype(bodyModelMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "update")
			return tfErr.GetDiag()
		}
		replaceE2mOptions.Event2MetricPrototype = convertedModel

		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.ReplaceE2mWithContext(context, replaceE2mOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceE2mWithContext failed: %s", err.Error()), "ibm_logs_e2m", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsE2mRead(context, d, meta)
}

func resourceIbmLogsE2mDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_e2m", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, e2mId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteE2mOptions := &logsv0.DeleteE2mOptions{}

	deleteE2mOptions.SetID(e2mId)

	_, err = logsClient.DeleteE2mWithContext(context, deleteE2mOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteE2mWithContext failed: %s", err.Error()), "ibm_logs_e2m", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricLabel(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2MetricLabel, error) {
	model := &logsv0.ApisEvents2metricsV2MetricLabel{}
	if modelMap["target_label"] != nil && modelMap["target_label"].(string) != "" {
		model.TargetLabel = core.StringPtr(modelMap["target_label"].(string))
	}
	if modelMap["source_field"] != nil && modelMap["source_field"].(string) != "" {
		model.SourceField = core.StringPtr(modelMap["source_field"].(string))
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricField(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2MetricField, error) {
	model := &logsv0.ApisEvents2metricsV2MetricField{}
	if modelMap["target_base_metric_name"] != nil && modelMap["target_base_metric_name"].(string) != "" {
		model.TargetBaseMetricName = core.StringPtr(modelMap["target_base_metric_name"].(string))
	}
	if modelMap["source_field"] != nil && modelMap["source_field"].(string) != "" {
		model.SourceField = core.StringPtr(modelMap["source_field"].(string))
	}
	if modelMap["aggregations"] != nil {
		aggregations := []logsv0.ApisEvents2metricsV2AggregationIntf{}
		for _, aggregationsItem := range modelMap["aggregations"].([]interface{}) {
			aggregationsItemModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2Aggregation(aggregationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			aggregations = append(aggregations, aggregationsItemModel)
		}
		model.Aggregations = aggregations
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2Aggregation(modelMap map[string]interface{}) (logsv0.ApisEvents2metricsV2AggregationIntf, error) {
	model := &logsv0.ApisEvents2metricsV2Aggregation{}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["agg_type"] != nil && modelMap["agg_type"].(string) != "" {
		model.AggType = core.StringPtr(modelMap["agg_type"].(string))
	}
	if modelMap["target_metric_name"] != nil && modelMap["target_metric_name"].(string) != "" {
		model.TargetMetricName = core.StringPtr(modelMap["target_metric_name"].(string))
	}
	if modelMap["samples"] != nil && len(modelMap["samples"].([]interface{})) > 0 {
		SamplesModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggSamples(modelMap["samples"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Samples = SamplesModel
	}
	if modelMap["histogram"] != nil && len(modelMap["histogram"].([]interface{})) > 0 {
		HistogramModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggHistogram(modelMap["histogram"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Histogram = HistogramModel
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggSamples(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2E2mAggSamples, error) {
	model := &logsv0.ApisEvents2metricsV2E2mAggSamples{}
	if modelMap["sample_type"] != nil && modelMap["sample_type"].(string) != "" {
		model.SampleType = core.StringPtr(modelMap["sample_type"].(string))
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggHistogram(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2E2mAggHistogram, error) {
	model := &logsv0.ApisEvents2metricsV2E2mAggHistogram{}
	if modelMap["buckets"] != nil {
		buckets := []float32{}
		for _, bucketsItem := range modelMap["buckets"].([]interface{}) {
			buckets = append(buckets, float32(bucketsItem.(float64)))
		}
		model.Buckets = buckets
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2AggregationAggMetadataSamples(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples, error) {
	model := &logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples{}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["agg_type"] != nil && modelMap["agg_type"].(string) != "" {
		model.AggType = core.StringPtr(modelMap["agg_type"].(string))
	}
	if modelMap["target_metric_name"] != nil && modelMap["target_metric_name"].(string) != "" {
		model.TargetMetricName = core.StringPtr(modelMap["target_metric_name"].(string))
	}
	if modelMap["samples"] != nil && len(modelMap["samples"].([]interface{})) > 0 {
		SamplesModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggSamples(modelMap["samples"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Samples = SamplesModel
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisEvents2metricsV2AggregationAggMetadataHistogram(modelMap map[string]interface{}) (*logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram, error) {
	model := &logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram{}
	if modelMap["enabled"] != nil {
		model.Enabled = core.BoolPtr(modelMap["enabled"].(bool))
	}
	if modelMap["agg_type"] != nil && modelMap["agg_type"].(string) != "" {
		model.AggType = core.StringPtr(modelMap["agg_type"].(string))
	}
	if modelMap["target_metric_name"] != nil && modelMap["target_metric_name"].(string) != "" {
		model.TargetMetricName = core.StringPtr(modelMap["target_metric_name"].(string))
	}
	if modelMap["histogram"] != nil && len(modelMap["histogram"].([]interface{})) > 0 {
		HistogramModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2E2mAggHistogram(modelMap["histogram"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Histogram = HistogramModel
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToApisLogs2metricsV2LogsQuery(modelMap map[string]interface{}) (*logsv0.ApisLogs2metricsV2LogsQuery, error) {
	model := &logsv0.ApisLogs2metricsV2LogsQuery{}
	if modelMap["lucene"] != nil && modelMap["lucene"].(string) != "" {
		model.Lucene = core.StringPtr(modelMap["lucene"].(string))
	}
	if modelMap["alias"] != nil && modelMap["alias"].(string) != "" {
		model.Alias = core.StringPtr(modelMap["alias"].(string))
	}
	if modelMap["applicationname_filters"] != nil {
		applicationnameFilters := []string{}
		for _, applicationnameFiltersItem := range modelMap["applicationname_filters"].([]interface{}) {
			applicationnameFilters = append(applicationnameFilters, applicationnameFiltersItem.(string))
		}
		model.ApplicationnameFilters = applicationnameFilters
	}
	if modelMap["subsystemname_filters"] != nil {
		subsystemnameFilters := []string{}
		for _, subsystemnameFiltersItem := range modelMap["subsystemname_filters"].([]interface{}) {
			subsystemnameFilters = append(subsystemnameFilters, subsystemnameFiltersItem.(string))
		}
		model.SubsystemnameFilters = subsystemnameFilters
	}
	if modelMap["severity_filters"] != nil {
		severityFilters := []string{}
		for _, severityFiltersItem := range modelMap["severity_filters"].([]interface{}) {
			severityFilters = append(severityFilters, severityFiltersItem.(string))
		}
		model.SeverityFilters = severityFilters
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToEvent2MetricPrototype(modelMap map[string]interface{}) (logsv0.Event2MetricPrototypeIntf, error) {
	model := &logsv0.Event2MetricPrototype{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["permutations_limit"] != nil {
		model.PermutationsLimit = core.Int64Ptr(int64(modelMap["permutations_limit"].(int)))
	}
	if modelMap["metric_labels"] != nil {
		metricLabels := []logsv0.ApisEvents2metricsV2MetricLabel{}
		for _, metricLabelsItem := range modelMap["metric_labels"].([]interface{}) {
			metricLabelsItemModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricLabel(metricLabelsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			metricLabels = append(metricLabels, *metricLabelsItemModel)
		}
		model.MetricLabels = metricLabels
	}
	if modelMap["metric_fields"] != nil {
		metricFields := []logsv0.ApisEvents2metricsV2MetricField{}
		for _, metricFieldsItem := range modelMap["metric_fields"].([]interface{}) {
			metricFieldsItemModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricField(metricFieldsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			metricFields = append(metricFields, *metricFieldsItemModel)
		}
		model.MetricFields = metricFields
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["logs_query"] != nil && len(modelMap["logs_query"].([]interface{})) > 0 {
		LogsQueryModel, err := ResourceIbmLogsE2mMapToApisLogs2metricsV2LogsQuery(modelMap["logs_query"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogsQuery = LogsQueryModel
	}
	return model, nil
}

func ResourceIbmLogsE2mMapToEvent2MetricPrototypeApisEvents2metricsV2E2mCreateParamsQueryLogsQuery(modelMap map[string]interface{}) (*logsv0.Event2MetricPrototypeApisEvents2metricsV2E2mCreateParamsQueryLogsQuery, error) {
	model := &logsv0.Event2MetricPrototypeApisEvents2metricsV2E2mCreateParamsQueryLogsQuery{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["permutations_limit"] != nil {
		model.PermutationsLimit = core.Int64Ptr(int64(modelMap["permutations_limit"].(int)))
	}
	if modelMap["metric_labels"] != nil {
		metricLabels := []logsv0.ApisEvents2metricsV2MetricLabel{}
		for _, metricLabelsItem := range modelMap["metric_labels"].([]interface{}) {
			metricLabelsItemModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricLabel(metricLabelsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			metricLabels = append(metricLabels, *metricLabelsItemModel)
		}
		model.MetricLabels = metricLabels
	}
	if modelMap["metric_fields"] != nil {
		metricFields := []logsv0.ApisEvents2metricsV2MetricField{}
		for _, metricFieldsItem := range modelMap["metric_fields"].([]interface{}) {
			metricFieldsItemModel, err := ResourceIbmLogsE2mMapToApisEvents2metricsV2MetricField(metricFieldsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			metricFields = append(metricFields, *metricFieldsItemModel)
		}
		model.MetricFields = metricFields
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["logs_query"] != nil && len(modelMap["logs_query"].([]interface{})) > 0 {
		LogsQueryModel, err := ResourceIbmLogsE2mMapToApisLogs2metricsV2LogsQuery(modelMap["logs_query"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogsQuery = LogsQueryModel
	}
	return model, nil
}

func ResourceIbmLogsE2mApisEvents2metricsV2MetricLabelToMap(model *logsv0.ApisEvents2metricsV2MetricLabel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetLabel != nil {
		modelMap["target_label"] = *model.TargetLabel
	}
	if model.SourceField != nil {
		modelMap["source_field"] = *model.SourceField
	}
	return modelMap, nil
}

func ResourceIbmLogsE2mApisEvents2metricsV2MetricFieldToMap(model *logsv0.ApisEvents2metricsV2MetricField) (map[string]interface{}, error) {
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
			aggregationsItemMap, err := ResourceIbmLogsE2mApisEvents2metricsV2AggregationToMap(aggregationsItem)
			if err != nil {
				return modelMap, err
			}
			aggregations = append(aggregations, aggregationsItemMap)
		}
		modelMap["aggregations"] = aggregations
	}
	return modelMap, nil
}

func ResourceIbmLogsE2mApisEvents2metricsV2AggregationToMap(model logsv0.ApisEvents2metricsV2AggregationIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples); ok {
		return ResourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataSamplesToMap(model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples))
	} else if _, ok := model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram); ok {
		return ResourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataHistogramToMap(model.(*logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram))
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
			samplesMap, err := ResourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model.Samples)
			if err != nil {
				return modelMap, err
			}
			modelMap["samples"] = []map[string]interface{}{samplesMap}
		}
		if model.Histogram != nil {
			histogramMap, err := ResourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model.Histogram)
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

func ResourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model *logsv0.ApisEvents2metricsV2E2mAggSamples) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SampleType != nil {
		modelMap["sample_type"] = *model.SampleType
	}
	return modelMap, nil
}

func ResourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model *logsv0.ApisEvents2metricsV2E2mAggHistogram) (map[string]interface{}, error) {
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

func ResourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataSamplesToMap(model *logsv0.ApisEvents2metricsV2AggregationAggMetadataSamples) (map[string]interface{}, error) {
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
		samplesMap, err := ResourceIbmLogsE2mApisEvents2metricsV2E2mAggSamplesToMap(model.Samples)
		if err != nil {
			return modelMap, err
		}
		modelMap["samples"] = []map[string]interface{}{samplesMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsE2mApisEvents2metricsV2AggregationAggMetadataHistogramToMap(model *logsv0.ApisEvents2metricsV2AggregationAggMetadataHistogram) (map[string]interface{}, error) {
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
		histogramMap, err := ResourceIbmLogsE2mApisEvents2metricsV2E2mAggHistogramToMap(model.Histogram)
		if err != nil {
			return modelMap, err
		}
		modelMap["histogram"] = []map[string]interface{}{histogramMap}
	}
	return modelMap, nil
}
func ResourceIbmLogsE2mApisLogs2metricsV2LogsQueryToMap(model *logsv0.ApisLogs2metricsV2LogsQuery) (map[string]interface{}, error) {
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

func ResourceIbmLogsE2mApisEvents2metricsV2E2mPermutationsToMap(model *logsv0.ApisEvents2metricsV2E2mPermutations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Limit != nil {
		modelMap["limit"] = flex.IntValue(model.Limit)
	}
	if model.HasExceededLimit != nil {
		modelMap["has_exceeded_limit"] = *model.HasExceededLimit
	}
	return modelMap, nil
}
