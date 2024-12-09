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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-router-go-sdk/ibmcloudlogsroutingv0"
)

func DataSourceIBMLogsRouterTenants() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMLogsRouterTenantsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Optional: The name of a tenant.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region where the tenants exist.",
			},
			"tenants": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tenants in the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the tenant.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time stamp the tenant was originally created.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time stamp the tenant was last updated.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud resource name of the tenant.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this tenant. The name is regionally unique across all tenants in the account.",
						},
						"etag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource version identifier.",
						},
						"targets": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of targets.",
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
				},
			},
		},
	}
}

func dataSourceIBMLogsRouterTenantsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_router_tenants", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listTenantsOptions := &ibmcloudlogsroutingv0.ListTenantsOptions{}

	if _, ok := d.GetOk("name"); ok {
		listTenantsOptions.SetName(d.Get("name").(string))
	}

	if _, ok := d.GetOk("region"); ok {
		listTenantsOptions.SetRegion(d.Get("region").(string))
	}

	tenantCollection, _, err := ibmCloudLogsRoutingClient.ListTenantsWithContext(context, listTenantsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListTenantsWithContext failed: %s", err.Error()), "(Data) ibm_logs_router_tenants", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMLogsRouterTenantsID(d))

	tenants := []map[string]interface{}{}
	if tenantCollection.Tenants != nil {
		for _, modelItem := range tenantCollection.Tenants {
			modelMap, err := DataSourceIBMLogsRouterTenantsTenantToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_router_tenants", "read", "tenants-to-map").GetDiag()
			}
			tenants = append(tenants, modelMap)
		}
	}
	if err = d.Set("tenants", tenants); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tenants: %s", err), "(Data) ibm_logs_router_tenants", "read", "set-tenants").GetDiag()
	}

	return nil
}

// dataSourceIBMLogsRouterTenantsID returns a reasonable ID for the list.
func dataSourceIBMLogsRouterTenantsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMLogsRouterTenantsTenantToMap(model *ibmcloudlogsroutingv0.Tenant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	modelMap["crn"] = *model.CRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	targets := []map[string]interface{}{}
	for _, targetsItem := range model.Targets {
		targetsItemMap, err := DataSourceIBMLogsRouterTenantsTargetTypeToMap(targetsItem)
		if err != nil {
			return modelMap, err
		}
		targets = append(targets, targetsItemMap)
	}
	modelMap["targets"] = targets
	return modelMap, nil
}

func DataSourceIBMLogsRouterTenantsTargetTypeToMap(model ibmcloudlogsroutingv0.TargetTypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogDna); ok {
		return DataSourceIBMLogsRouterTenantsTargetTypeLogDnaToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogDna))
	} else if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogs); ok {
		return DataSourceIBMLogsRouterTenantsTargetTypeLogsToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogs))
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
			parametersMap, err := DataSourceIBMLogsRouterTenantsTargetParametersTypeLogDnaToMap(model.Parameters)
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

func DataSourceIBMLogsRouterTenantsTargetParametersTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}

func DataSourceIBMLogsRouterTenantsTargetTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := DataSourceIBMLogsRouterTenantsTargetParametersTypeLogDnaToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func DataSourceIBMLogsRouterTenantsTargetTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := DataSourceIBMLogsRouterTenantsTargetParametersTypeLogsToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func DataSourceIBMLogsRouterTenantsTargetParametersTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}
