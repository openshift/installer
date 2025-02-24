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

func DataSourceIbmLogsEnrichments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsEnrichmentsRead,

		Schema: map[string]*schema.Schema{
			"enrichments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The enrichments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The enrichment ID.",
						},
						"field_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enrichment field name.",
						},
						"enrichment_type": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The enrichment type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"geo_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The geo ip enrichment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"suspicious_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The suspicious ip enrichment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"custom_enrichment": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The custom enrichment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID of the custom enrichment.",
												},
											},
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

func dataSourceIbmLogsEnrichmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_enrichments", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getEnrichmentsOptions := &logsv0.GetEnrichmentsOptions{}

	entrichmentCollection, _, err := logsClient.GetEnrichmentsWithContext(context, getEnrichmentsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEnrichmentsWithContext failed: %s", err.Error()), "(Data) ibm_logs_enrichments", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsEnrichmentsID(d))

	enrichments := []map[string]interface{}{}
	if entrichmentCollection.Enrichments != nil {
		for _, modelItem := range entrichmentCollection.Enrichments {
			modelMap, err := DataSourceIbmLogsEnrichmentsEnrichmentToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_enrichments", "read")
				return tfErr.GetDiag()
			}
			enrichments = append(enrichments, modelMap)
		}
	}
	if err = d.Set("enrichments", enrichments); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enrichments: %s", err), "(Data) ibm_logs_enrichments", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsEnrichmentsID returns a reasonable ID for the list.
func dataSourceIbmLogsEnrichmentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsEnrichmentsEnrichmentToMap(model *logsv0.Enrichment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	modelMap["field_name"] = *model.FieldName
	enrichmentTypeMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeToMap(model.EnrichmentType)
	if err != nil {
		return modelMap, err
	}
	modelMap["enrichment_type"] = []map[string]interface{}{enrichmentTypeMap}
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeToMap(model logsv0.EnrichmentV1EnrichmentTypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp); ok {
		return DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeGeoIpToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp); ok {
		return DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeSuspiciousIpToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment); ok {
		return DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeCustomEnrichmentToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentType); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.EnrichmentV1EnrichmentType)
		if model.GeoIp != nil {
			geoIpMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1GeoIpTypeEmptyToMap(model.GeoIp)
			if err != nil {
				return modelMap, err
			}
			modelMap["geo_ip"] = []map[string]interface{}{geoIpMap}
		}
		if model.SuspiciousIp != nil {
			suspiciousIpMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1SuspiciousIpTypeEmptyToMap(model.SuspiciousIp)
			if err != nil {
				return modelMap, err
			}
			modelMap["suspicious_ip"] = []map[string]interface{}{suspiciousIpMap}
		}
		if model.CustomEnrichment != nil {
			customEnrichmentMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1CustomEnrichmentTypeToMap(model.CustomEnrichment)
			if err != nil {
				return modelMap, err
			}
			modelMap["custom_enrichment"] = []map[string]interface{}{customEnrichmentMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.EnrichmentV1EnrichmentTypeIntf subtype encountered")
	}
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1GeoIpTypeEmptyToMap(model *logsv0.EnrichmentV1GeoIpTypeEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1SuspiciousIpTypeEmptyToMap(model *logsv0.EnrichmentV1SuspiciousIpTypeEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1CustomEnrichmentTypeToMap(model *logsv0.EnrichmentV1CustomEnrichmentType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeGeoIpToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GeoIp != nil {
		geoIpMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1GeoIpTypeEmptyToMap(model.GeoIp)
		if err != nil {
			return modelMap, err
		}
		modelMap["geo_ip"] = []map[string]interface{}{geoIpMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeSuspiciousIpToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SuspiciousIp != nil {
		suspiciousIpMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1SuspiciousIpTypeEmptyToMap(model.SuspiciousIp)
		if err != nil {
			return modelMap, err
		}
		modelMap["suspicious_ip"] = []map[string]interface{}{suspiciousIpMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsEnrichmentsEnrichmentV1EnrichmentTypeTypeCustomEnrichmentToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomEnrichment != nil {
		customEnrichmentMap, err := DataSourceIbmLogsEnrichmentsEnrichmentV1CustomEnrichmentTypeToMap(model.CustomEnrichment)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_enrichment"] = []map[string]interface{}{customEnrichmentMap}
	}
	return modelMap, nil
}
