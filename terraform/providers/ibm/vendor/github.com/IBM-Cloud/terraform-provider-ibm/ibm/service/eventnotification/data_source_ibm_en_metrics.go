// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnMetricsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"destination_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination type. Allowed values are [smtp_custom].",
			},
			"gte": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "GTE (greater than equal), start timestamp in UTC.",
			},
			"lte": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "LTE (less than equal), end timestamp in UTC.",
			},
			"destination_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier for Destination.",
			},
			"source_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier for Source.",
			},
			"email_to": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Receiver email id.",
			},
			"notification_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notification Id.",
			},
			"subject": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email subject.",
			},
			"metrics": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "array of metrics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "key.",
						},
						"doc_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "doc count.",
						},
						"histogram": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Payload describing histogram.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"buckets": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of buckets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"doc_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total count.",
												},
												"key_as_string": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Timestamp.",
												},
											},
										},
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

func dataSourceIBMEnMetricsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getMetricsOptions := &eventnotificationsv1.GetMetricsOptions{}

	getMetricsOptions.SetInstanceID(d.Get("instance_id").(string))
	getMetricsOptions.SetDestinationType(d.Get("destination_type").(string))
	getMetricsOptions.SetGte(d.Get("gte").(string))
	getMetricsOptions.SetLte(d.Get("lte").(string))
	if _, ok := d.GetOk("destination_id"); ok {
		getMetricsOptions.SetDestinationID(d.Get("destination_id").(string))
	}
	if _, ok := d.GetOk("source_id"); ok {
		getMetricsOptions.SetSourceID(d.Get("source_id").(string))
	}
	if _, ok := d.GetOk("email_to"); ok {
		getMetricsOptions.SetEmailTo(d.Get("email_to").(string))
	}
	if _, ok := d.GetOk("notification_id"); ok {
		getMetricsOptions.SetNotificationID(d.Get("notification_id").(string))
	}
	if _, ok := d.GetOk("subject"); ok {
		getMetricsOptions.SetSubject(d.Get("subject").(string))
	}

	metricsres, _, err := eventNotificationsClient.GetMetricsWithContext(context, getMetricsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetMetricsWithContext failed: %s", err.Error()), "(Data) ibm_en_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMEnMetricsID(d))

	metrics := []map[string]interface{}{}
	if metricsres.Metrics != nil {
		for _, modelItem := range metricsres.Metrics {
			modelMap, err := dataSourceIBMEnMetricsMetricToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_metrics", "read")
				return tfErr.GetDiag()
			}
			metrics = append(metrics, modelMap)
		}
	}
	if err = d.Set("metrics", metrics); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting metrics: %s", err), "(Data) ibm_en_metrics", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIBMEnMetricsID returns a reasonable ID for the list.
func dataSourceIBMEnMetricsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMEnMetricsMetricToMap(model *eventnotificationsv1.Metric) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.DocCount != nil {
		modelMap["doc_count"] = flex.IntValue(model.DocCount)
	}
	if model.Histogram != nil {
		histogramMap, err := dataSourceIBMEnMetricsHistrogramToMap(model.Histogram)
		if err != nil {
			return modelMap, err
		}
		modelMap["histogram"] = []map[string]interface{}{histogramMap}
	}
	return modelMap, nil
}

func dataSourceIBMEnMetricsHistrogramToMap(model *eventnotificationsv1.Histrogram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Buckets != nil {
		buckets := []map[string]interface{}{}
		for _, bucketsItem := range model.Buckets {
			bucketsItemMap, err := dataSourceIBMEnMetricsBucketsToMap(&bucketsItem)
			if err != nil {
				return modelMap, err
			}
			buckets = append(buckets, bucketsItemMap)
		}
		modelMap["buckets"] = buckets
	}
	return modelMap, nil
}

func dataSourceIBMEnMetricsBucketsToMap(model *eventnotificationsv1.Buckets) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DocCount != nil {
		modelMap["doc_count"] = flex.IntValue(model.DocCount)
	}
	if model.KeyAsString != nil {
		modelMap["key_as_string"] = model.KeyAsString.String()
	}
	return modelMap, nil
}
