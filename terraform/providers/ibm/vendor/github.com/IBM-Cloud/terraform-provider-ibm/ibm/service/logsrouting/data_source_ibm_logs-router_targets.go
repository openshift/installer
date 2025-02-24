// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.90.1-64fd3296-20240515-180710
 */

package logsrouting

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-router-go-sdk/ibmcloudlogsroutingv0"
)

func DataSourceIBMLogsRouterTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMLogsRouterTargetsRead,

		Schema: map[string]*schema.Schema{
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance ID of the tenant.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional: Name of the tenant target.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region where the tenant for these targets exist.",
			},
			"targets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of target of a tenant.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the target.",
						},
						"log_sink_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud resource name of the log-sink target instance.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this tenant target. The name is unique across all targets for this tenant.",
						},
						"etag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource version identifier.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of log-sink. Identical to the <code>service-name</code> segment of <code>log_sink_crn</code>.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time stamp the target was originally created.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time stamp the target was last updated.",
						},
						"parameters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of properties returned from a successful list operation for a log-sink of type IBM Log Analysis (logdna).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host name of the log-sink.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network port of the log-sink.",
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

func dataSourceIBMLogsRouterTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_router_targets", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listTenantTargetsOptions := &ibmcloudlogsroutingv0.ListTenantTargetsOptions{}

	listTenantTargetsOptions.SetTenantID(core.UUIDPtr(strfmt.UUID(d.Get("tenant_id").(string))))
	if _, ok := d.GetOk("name"); ok {
		listTenantTargetsOptions.SetName(d.Get("name").(string))
	}

	if _, ok := d.GetOk("region"); ok {
		listTenantTargetsOptions.SetRegion(d.Get("region").(string))
	}

	targetTypeCollection, _, err := ibmCloudLogsRoutingClient.ListTenantTargetsWithContext(context, listTenantTargetsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListTenantTargetsWithContext failed: %s", err.Error()), "(Data) ibm_logs_router_targets", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMLogsRouterTargetsID(d))

	targets := []map[string]interface{}{}
	if targetTypeCollection.Targets != nil {
		for _, modelItem := range targetTypeCollection.Targets {
			modelMap, err := DataSourceIBMLogsRouterTargetsTargetTypeToMap(modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_router_targets", "read", "targets-to-map").GetDiag()
			}
			targets = append(targets, modelMap)
		}
	}
	if err = d.Set("targets", targets); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting targets: %s", err), "(Data) ibm_logs_router_targets", "read", "set-targets").GetDiag()
	}

	return nil
}

// dataSourceIBMLogsRouterTargetsID returns a reasonable ID for the list.
func dataSourceIBMLogsRouterTargetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMLogsRouterTargetsTargetTypeToMap(model ibmcloudlogsroutingv0.TargetTypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogDna); ok {
		return DataSourceIBMLogsRouterTargetsTargetTypeLogDnaToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogDna))
	} else if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogs); ok {
		return DataSourceIBMLogsRouterTargetsTargetTypeLogsToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogs))
	} else if _, ok := model.(*ibmcloudlogsroutingv0.TargetType); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ibmcloudlogsroutingv0.TargetType)
		if model.ID != nil {
			modelMap["id"] = model.ID.String()
		}
		if model.LogSinkCRN != nil {
			modelMap["log_sink_crn"] = *model.LogSinkCRN
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.Etag != nil {
			modelMap["etag"] = *model.Etag
		}
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.CreatedAt != nil {
			modelMap["created_at"] = *model.CreatedAt
		}
		if model.UpdatedAt != nil {
			modelMap["updated_at"] = *model.UpdatedAt
		}
		if model.Parameters != nil {
			parametersMap, err := DataSourceIBMLogsRouterTargetsTargetParametersTypeLogDnaToMap(model.Parameters)
			if err != nil {
				return modelMap, err
			}
			modelMap["parameters"] = []map[string]interface{}{parametersMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ibmcloudlogsroutingv0.TargetTypeIntf subtype encountered")
	}
}

func DataSourceIBMLogsRouterTargetsTargetParametersTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}

func DataSourceIBMLogsRouterTargetsTargetTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := DataSourceIBMLogsRouterTargetsTargetParametersTypeLogDnaToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func DataSourceIBMLogsRouterTargetsTargetTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := DataSourceIBMLogsRouterTargetsTargetParametersTypeLogsToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func DataSourceIBMLogsRouterTargetsTargetParametersTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}
