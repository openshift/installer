// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsEnrichment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsEnrichmentCreate,
		ReadContext:   resourceIbmLogsEnrichmentRead,
		UpdateContext: resourceIbmLogsEnrichmentUpdate,
		DeleteContext: resourceIbmLogsEnrichmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"field_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_enrichment", "field_name"),
				Description:  "The enrichment field name.",
			},
			"enrichment_type": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "The enrichment type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_ip": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The geo ip enrichment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"suspicious_ip": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The suspicious ip enrichment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"custom_enrichment": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The custom enrichment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The ID of the custom enrichment.",
									},
								},
							},
						},
					},
				},
			},
			"enrichment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enrichment ID",
			},
		},
	}
}

func ResourceIbmLogsEnrichmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "field_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_enrichment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsEnrichmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no update feature available for enrichments
	return resourceIbmLogsEnrichmentRead(context, d, meta)
}

func resourceIbmLogsEnrichmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_enrichment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createEnrichmentOptions := &logsv0.CreateEnrichmentOptions{}

	createEnrichmentOptions.SetFieldName(d.Get("field_name").(string))
	enrichmentTypeModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1EnrichmentType(d.Get("enrichment_type.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createEnrichmentOptions.SetEnrichmentType(enrichmentTypeModel)

	enrichment, _, err := logsClient.CreateEnrichmentWithContext(context, createEnrichmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateEnrichmentWithContext failed: %s", err.Error()), "ibm_logs_enrichment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	enrichmentID := fmt.Sprintf("%s/%s/%d", region, instanceId, *enrichment.ID)
	d.SetId(enrichmentID)

	return resourceIbmLogsEnrichmentRead(context, d, meta)
}

func resourceIbmLogsEnrichmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_enrichment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, region, instanceId, enrichmentID, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getEnrichmentsOptions := &logsv0.GetEnrichmentsOptions{}

	entrichmentCollection, response, err := logsClient.GetEnrichmentsWithContext(context, getEnrichmentsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEnrichmentsWithContext failed: %s", err.Error()), "ibm_logs_enrichment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("enrichment_id", enrichmentID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dashboard_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	enrichmentIDInt, _ := strconv.ParseInt(enrichmentID, 10, 64)
	if entrichmentCollection != nil && len(entrichmentCollection.Enrichments) > 0 {

		for _, enrichment := range entrichmentCollection.Enrichments {

			if *enrichment.ID == enrichmentIDInt {

				if err = d.Set("field_name", enrichment.FieldName); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting field_name: %s", err))
				}
				enrichmentTypeMap, err := ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeToMap(enrichment.EnrichmentType)
				if err != nil {
					return diag.FromErr(err)
				}
				if err = d.Set("enrichment_type", []map[string]interface{}{enrichmentTypeMap}); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting enrichment_type: %s", err))
				}

			}
		}
	}

	return nil
}

func resourceIbmLogsEnrichmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_enrichment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, enrichmentID, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	removeEnrichmentsOptions := &logsv0.RemoveEnrichmentsOptions{}

	enrichmentIDInt, _ := strconv.ParseInt(enrichmentID, 10, 64)
	removeEnrichmentsOptions.SetID(enrichmentIDInt)

	_, err = logsClient.RemoveEnrichmentsWithContext(context, removeEnrichmentsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("RemoveEnrichmentsWithContext failed: %s", err.Error()), "ibm_logs_enrichment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1EnrichmentType(modelMap map[string]interface{}) (logsv0.EnrichmentV1EnrichmentTypeIntf, error) {
	model := &logsv0.EnrichmentV1EnrichmentType{}
	if modelMap["geo_ip"] != nil && len(modelMap["geo_ip"].([]interface{})) > 0 {
		GeoIpModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1GeoIpTypeEmpty(modelMap["geo_ip"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GeoIp = GeoIpModel
	}
	if modelMap["suspicious_ip"] != nil && len(modelMap["suspicious_ip"].([]interface{})) > 0 {
		SuspiciousIpModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1SuspiciousIpTypeEmpty(modelMap["suspicious_ip"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.SuspiciousIp = SuspiciousIpModel
	}
	if modelMap["custom_enrichment"] != nil && len(modelMap["custom_enrichment"].([]interface{})) > 0 {
		CustomEnrichmentModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1CustomEnrichmentType(modelMap["custom_enrichment"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.CustomEnrichment = CustomEnrichmentModel
	}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1GeoIpTypeEmpty(modelMap []interface{}) (*logsv0.EnrichmentV1GeoIpTypeEmpty, error) {
	model := &logsv0.EnrichmentV1GeoIpTypeEmpty{}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1SuspiciousIpTypeEmpty(modelMap []interface{}) (*logsv0.EnrichmentV1SuspiciousIpTypeEmpty, error) {
	model := &logsv0.EnrichmentV1SuspiciousIpTypeEmpty{}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1CustomEnrichmentType(modelMap []interface{}) (*logsv0.EnrichmentV1CustomEnrichmentType, error) {
	model := &logsv0.EnrichmentV1CustomEnrichmentType{}

	if modelMap != nil && len(modelMap) > 0 && modelMap[0] != nil {
		modelMapElement := modelMap[0].(map[string]interface{})
		if modelMapElement["id"] != nil {
			model.ID = core.Int64Ptr(int64(modelMapElement["id"].(int)))
		}
	}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1EnrichmentTypeTypeGeoIp(modelMap map[string]interface{}) (*logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp, error) {
	model := &logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp{}
	if modelMap["geo_ip"] != nil && len(modelMap["geo_ip"].([]interface{})) > 0 {
		GeoIpModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1GeoIpTypeEmpty(modelMap["geo_ip"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.GeoIp = GeoIpModel
	}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1EnrichmentTypeTypeSuspiciousIp(modelMap map[string]interface{}) (*logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp, error) {
	model := &logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp{}
	if modelMap["suspicious_ip"] != nil && len(modelMap["suspicious_ip"].([]interface{})) > 0 {
		SuspiciousIpModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1SuspiciousIpTypeEmpty(modelMap["suspicious_ip"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.SuspiciousIp = SuspiciousIpModel
	}
	return model, nil
}

func ResourceIbmLogsEnrichmentMapToEnrichmentV1EnrichmentTypeTypeCustomEnrichment(modelMap map[string]interface{}) (*logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment, error) {
	model := &logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment{}
	if modelMap["custom_enrichment"] != nil && len(modelMap["custom_enrichment"].([]interface{})) > 0 {
		CustomEnrichmentModel, err := ResourceIbmLogsEnrichmentMapToEnrichmentV1CustomEnrichmentType(modelMap["custom_enrichment"].([]interface{}))
		if err != nil {
			return model, err
		}
		model.CustomEnrichment = CustomEnrichmentModel
	}
	return model, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeToMap(model logsv0.EnrichmentV1EnrichmentTypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp); ok {
		return ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeGeoIpToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp); ok {
		return ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeSuspiciousIpToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment); ok {
		return ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeCustomEnrichmentToMap(model.(*logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment))
	} else if _, ok := model.(*logsv0.EnrichmentV1EnrichmentType); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.EnrichmentV1EnrichmentType)
		if model.GeoIp != nil {
			geoIpMap, err := ResourceIbmLogsEnrichmentEnrichmentV1GeoIpTypeEmptyToMap(model.GeoIp)
			if err != nil {
				return modelMap, err
			}
			modelMap["geo_ip"] = []map[string]interface{}{geoIpMap}
		}
		if model.SuspiciousIp != nil {
			suspiciousIpMap, err := ResourceIbmLogsEnrichmentEnrichmentV1SuspiciousIpTypeEmptyToMap(model.SuspiciousIp)
			if err != nil {
				return modelMap, err
			}
			modelMap["suspicious_ip"] = []map[string]interface{}{suspiciousIpMap}
		}
		if model.CustomEnrichment != nil {
			customEnrichmentMap, err := ResourceIbmLogsEnrichmentEnrichmentV1CustomEnrichmentTypeToMap(model.CustomEnrichment)
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

func ResourceIbmLogsEnrichmentEnrichmentV1GeoIpTypeEmptyToMap(model *logsv0.EnrichmentV1GeoIpTypeEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1SuspiciousIpTypeEmptyToMap(model *logsv0.EnrichmentV1SuspiciousIpTypeEmpty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1CustomEnrichmentTypeToMap(model *logsv0.EnrichmentV1CustomEnrichmentType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	return modelMap, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeGeoIpToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeGeoIp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.GeoIp != nil {
		geoIpMap, err := ResourceIbmLogsEnrichmentEnrichmentV1GeoIpTypeEmptyToMap(model.GeoIp)
		if err != nil {
			return modelMap, err
		}
		modelMap["geo_ip"] = []map[string]interface{}{geoIpMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeSuspiciousIpToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeSuspiciousIp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SuspiciousIp != nil {
		suspiciousIpMap, err := ResourceIbmLogsEnrichmentEnrichmentV1SuspiciousIpTypeEmptyToMap(model.SuspiciousIp)
		if err != nil {
			return modelMap, err
		}
		modelMap["suspicious_ip"] = []map[string]interface{}{suspiciousIpMap}
	}
	return modelMap, nil
}

func ResourceIbmLogsEnrichmentEnrichmentV1EnrichmentTypeTypeCustomEnrichmentToMap(model *logsv0.EnrichmentV1EnrichmentTypeTypeCustomEnrichment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomEnrichment != nil {
		customEnrichmentMap, err := ResourceIbmLogsEnrichmentEnrichmentV1CustomEnrichmentTypeToMap(model.CustomEnrichment)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_enrichment"] = []map[string]interface{}{customEnrichmentMap}
	}
	return modelMap, nil
}
