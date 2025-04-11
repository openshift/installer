// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

func DataSourceIBMAtrackerTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMAtrackerTargetsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the target resource.",
			},
			"targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of target resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uuid of the target resource.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the target resource.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the target resource.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the target.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Included this optional field if you used it to create a target in a different region other than the one you are connected.",
						},
						"encrypt_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "use encryption_key instead",
							Description: "The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This credential is masked in the response.",
						},
						"encryption_key": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Computed:    true,
							Description: "The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This credential is masked in the response.",
						},
						"cos_endpoint": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Property values for a Cloud Object Storage Endpoint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host name of the Cloud Object Storage endpoint.",
									},
									"target_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the Cloud Object Storage instance.",
									},
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bucket name under the Cloud Object Storage instance.",
									},
									"api_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
										Description: "The IAM API key that has writer access to the Cloud Object Storage instance. This credential is masked in the response. This is required if service_to_service is not enabled.",
									},
									"service_to_service_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag to true if service to service is enabled and do not supply an apikey.",
									},
								},
							},
						},
						"logdna_endpoint": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Property values for a LogDNA Endpoint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the LogDNA instance.",
									},
									"ingestion_key": {
										Type:        schema.TypeString,
										Sensitive:   true,
										Computed:    true,
										Description: "The LogDNA ingestion key is used for routing logs to a specific LogDNA instance.",
									},
								},
							},
						},
						"eventstreams_endpoint": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Property values for the Event Streams Endpoint in responses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the Event Streams instance.",
									},
									"brokers": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of broker endpoints.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"topic": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The messsage hub topic defined in the Event Streams instance.",
									},
									"api_key": &schema.Schema{ // pragma: allowlist secret
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
										Description: "The user password (api key) for the message hub topic in the Event Streams instance. This is required if service_to_service is not enabled.",
									},
									"service_to_service_enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag to true if service to service is enabled and do not supply an apikey.",
									},
								},
							},
						},
						"cloudlogs_endpoint": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Property values for the IBM Cloud Logs endpoint in responses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the IBM Cloud Logs instance",
									},
								},
							},
						},
						"cos_write_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Deprecated:  "use write_status instead",
							Description: "The status of the write attempt with the provided cos_endpoint parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status such as failed or success.",
									},
									"last_failure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The timestamp of the failure.",
									},
									"reason_for_last_failure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detailed description of the cause of the failure.",
									},
								},
							},
						},
						"write_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The status of the write attempt to the target with the provided endpoint parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status such as failed or success.",
									},
									"last_failure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The timestamp of the failure.",
									},
									"reason_for_last_failure": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detailed description of the cause of the failure.",
									},
								},
							},
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "use created_at instead",
							Description: "The timestamp of the target creation time.",
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "use updated_at instead",
							Description: "The timestamp of the target last updated time.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the target creation time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the target last updated time.",
						},
						"api_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The API version of the target.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMAtrackerTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClientv2, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listTargetsOptions := &atrackerv2.ListTargetsOptions{}

	targetList, response, err := atrackerClientv2.ListTargetsWithContext(context, listTargetsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListTargetsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListTargetsWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchTargets []atrackerv2.Target
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range targetList.Targets {
			if *data.Name == name {
				matchTargets = append(matchTargets, data)
			}
		}
	} else {
		matchTargets = targetList.Targets
	}
	targetList.Targets = matchTargets

	if suppliedFilter {
		if len(targetList.Targets) == 0 {
			return diag.FromErr(fmt.Errorf("no Targets found with name %s", name))
		}
		d.SetId(name)
	} else {
		d.SetId(DataSourceIBMAtrackerTargetsID(d))
	}

	targets := []map[string]interface{}{}
	if targetList.Targets != nil {
		for _, modelItem := range targetList.Targets {
			modelMap, err := DataSourceIBMAtrackerTargetsTargetToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			targets = append(targets, modelMap)
		}
	}
	if err = d.Set("targets", targets); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
	}
	return nil

}

// DataSourceIBMAtrackerTargetsID returns a reasonable ID for the list.
func DataSourceIBMAtrackerTargetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMAtrackerTargetsTargetToMap(model *atrackerv2.Target) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	if model.Region != nil {
		modelMap["region"] = *model.Region
	}
	if model.CosEndpoint != nil {
		cosEndpointMap, err := DataSourceIBMAtrackerTargetsCosEndpointToMap(model.CosEndpoint)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_endpoint"] = []map[string]interface{}{cosEndpointMap}
	}
	if model.LogdnaEndpoint != nil {
		logdnaEndpointMap, err := DataSourceIBMAtrackerTargetsLogdnaEndpointToMap(model.LogdnaEndpoint)
		if err != nil {
			return modelMap, err
		}
		modelMap["logdna_endpoint"] = []map[string]interface{}{logdnaEndpointMap}
	}
	if model.EventstreamsEndpoint != nil {
		eventstreamsEndpointMap, err := DataSourceIBMAtrackerTargetsEventstreamsEndpointToMap(model.EventstreamsEndpoint)
		if err != nil {
			return modelMap, err
		}
		modelMap["eventstreams_endpoint"] = []map[string]interface{}{eventstreamsEndpointMap}
	}
	if model.CloudlogsEndpoint != nil {
		cloudlogsEndpointMap, err := DataSourceIBMAtrackerTargetsCloudlogsEndpointToMap(model.CloudlogsEndpoint)
		if err != nil {
			return modelMap, err
		}
		modelMap["cloudlogs_endpoint"] = []map[string]interface{}{cloudlogsEndpointMap}
	}
	if model.WriteStatus != nil {
		writeStatusMap, err := DataSourceIBMAtrackerTargetsWriteStatusToMap(model.WriteStatus)
		if err != nil {
			return modelMap, err
		}
		modelMap["write_status"] = []map[string]interface{}{writeStatusMap}
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.APIVersion != nil {
		modelMap["api_version"] = *model.APIVersion
	}
	// TODO: Deprecated, to remove
	modelMap["encryption_key"] = REDACTED_TEXT
	return modelMap, nil
}

func DataSourceIBMAtrackerTargetsCosEndpointToMap(model *atrackerv2.CosEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Endpoint != nil {
		modelMap["endpoint"] = *model.Endpoint
	}
	if model.TargetCRN != nil {
		modelMap["target_crn"] = *model.TargetCRN
	}
	if model.Bucket != nil {
		modelMap["bucket"] = *model.Bucket
	}
	if model.ServiceToServiceEnabled != nil {
		modelMap["service_to_service_enabled"] = *model.ServiceToServiceEnabled
	}
	return modelMap, nil
}

func DataSourceIBMAtrackerTargetsLogdnaEndpointToMap(model *atrackerv2.LogdnaEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetCRN != nil {
		modelMap["target_crn"] = *model.TargetCRN
	}
	return modelMap, nil
}

func DataSourceIBMAtrackerTargetsEventstreamsEndpointToMap(model *atrackerv2.EventstreamsEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetCRN != nil {
		modelMap["target_crn"] = *model.TargetCRN
	}
	if model.Brokers != nil {
		modelMap["brokers"] = model.Brokers
	}
	if model.Topic != nil {
		modelMap["topic"] = *model.Topic
	}
	if model.APIKey != nil {
		modelMap["api_key"] = *model.APIKey // pragma: allowlist secret
	}
	if model.ServiceToServiceEnabled != nil {
		modelMap["service_to_service_enabled"] = *model.ServiceToServiceEnabled
	}
	return modelMap, nil
}

func DataSourceIBMAtrackerTargetsCloudlogsEndpointToMap(model *atrackerv2.CloudLogsEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetCRN != nil {
		modelMap["target_crn"] = *model.TargetCRN
	}
	return modelMap, nil
}

func DataSourceIBMAtrackerTargetsWriteStatusToMap(model *atrackerv2.WriteStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.LastFailure != nil {
		modelMap["last_failure"] = model.LastFailure.String()
	}
	if model.ReasonForLastFailure != nil {
		modelMap["reason_for_last_failure"] = *model.ReasonForLastFailure
	}
	return modelMap, nil
}
