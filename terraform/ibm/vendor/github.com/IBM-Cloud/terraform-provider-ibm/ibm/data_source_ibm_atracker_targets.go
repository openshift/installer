// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"
)

func dataSourceIBMAtrackerTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAtrackerTargetsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the target resource.",
			},
			"targets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of target resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uuid of the target resource.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the target resource.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the target resource.",
						},
						"target_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the target.",
						},
						"encrypt_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This credential is masked in the response.",
						},
						"cos_endpoint": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Property values for a Cloud Object Storage Endpoint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host name of the Cloud Object Storage endpoint.",
									},
									"target_crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the Cloud Object Storage instance.",
									},
									"bucket": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bucket name under the Cloud Object Storage instance.",
									},
									"api_key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
										Description: "The IAM API key that has writer access to the Cloud Object Storage instance. This credential is masked in the response.",
									},
								},
							},
						},
						"cos_write_status": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The status of the write attempt with the provided cos_endpoint parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status such as failed or success.",
									},
									"last_failure": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The timestamp of the failure.",
									},
									"reason_for_last_failure": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detailed description of the cause of the failure.",
									},
								},
							},
						},
						"created": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the target creation time.",
						},
						"updated": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the target last updated time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAtrackerTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	listTargetsOptions := &atrackerv1.ListTargetsOptions{}

	targetList, response, err := atrackerClient.ListTargetsWithContext(context, listTargetsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListTargetsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListTargetsWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchTargets []atrackerv1.Target
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
		d.SetId(dataSourceIBMAtrackerTargetsID(d))
	}

	if targetList.Targets != nil {
		err = d.Set("targets", dataSourceTargetListFlattenTargets(targetList.Targets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
		}
	}

	return nil
}

// dataSourceIBMAtrackerTargetsID returns a reasonable ID for the list.
func dataSourceIBMAtrackerTargetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceTargetListFlattenTargets(result []atrackerv1.Target) (targets []map[string]interface{}) {
	for _, targetsItem := range result {
		targets = append(targets, dataSourceTargetListTargetsToMap(targetsItem))
	}

	return targets
}

func dataSourceTargetListTargetsToMap(targetsItem atrackerv1.Target) (targetsMap map[string]interface{}) {
	targetsMap = map[string]interface{}{}

	if targetsItem.ID != nil {
		targetsMap["id"] = targetsItem.ID
	}
	if targetsItem.Name != nil {
		targetsMap["name"] = targetsItem.Name
	}
	if targetsItem.CRN != nil {
		targetsMap["crn"] = targetsItem.CRN
	}
	if targetsItem.TargetType != nil {
		targetsMap["target_type"] = targetsItem.TargetType
	}
	if targetsItem.EncryptKey != nil {
		targetsMap["encrypt_key"] = targetsItem.EncryptKey
	}
	if targetsItem.CosEndpoint != nil {
		cosEndpointList := []map[string]interface{}{}
		cosEndpointMap := dataSourceTargetListTargetsCosEndpointToMap(*targetsItem.CosEndpoint)
		cosEndpointList = append(cosEndpointList, cosEndpointMap)
		targetsMap["cos_endpoint"] = cosEndpointList
	}
	if targetsItem.CosWriteStatus != nil {
		cosWriteStatusList := []map[string]interface{}{}
		cosWriteStatusMap := dataSourceTargetListTargetsCosWriteStatusToMap(*targetsItem.CosWriteStatus)
		cosWriteStatusList = append(cosWriteStatusList, cosWriteStatusMap)
		targetsMap["cos_write_status"] = cosWriteStatusList
	}
	if targetsItem.Created != nil {
		targetsMap["created"] = targetsItem.Created.String()
	}
	if targetsItem.Updated != nil {
		targetsMap["updated"] = targetsItem.Updated.String()
	}

	return targetsMap
}

func dataSourceTargetListTargetsCosEndpointToMap(cosEndpointItem atrackerv1.CosEndpoint) (cosEndpointMap map[string]interface{}) {
	cosEndpointMap = map[string]interface{}{}

	if cosEndpointItem.Endpoint != nil {
		cosEndpointMap["endpoint"] = cosEndpointItem.Endpoint
	}
	if cosEndpointItem.TargetCRN != nil {
		cosEndpointMap["target_crn"] = cosEndpointItem.TargetCRN
	}
	if cosEndpointItem.Bucket != nil {
		cosEndpointMap["bucket"] = cosEndpointItem.Bucket
	}
	if cosEndpointItem.APIKey != nil {
		cosEndpointMap["api_key"] = cosEndpointItem.APIKey
	}

	return cosEndpointMap
}

func dataSourceTargetListTargetsCosWriteStatusToMap(cosWriteStatusItem atrackerv1.CosWriteStatus) (cosWriteStatusMap map[string]interface{}) {
	cosWriteStatusMap = map[string]interface{}{}

	if cosWriteStatusItem.Status != nil {
		cosWriteStatusMap["status"] = cosWriteStatusItem.Status
	}
	if cosWriteStatusItem.LastFailure != nil {
		cosWriteStatusMap["last_failure"] = cosWriteStatusItem.LastFailure.String()
	}
	if cosWriteStatusItem.ReasonForLastFailure != nil {
		cosWriteStatusMap["reason_for_last_failure"] = cosWriteStatusItem.ReasonForLastFailure
	}

	return cosWriteStatusMap
}
