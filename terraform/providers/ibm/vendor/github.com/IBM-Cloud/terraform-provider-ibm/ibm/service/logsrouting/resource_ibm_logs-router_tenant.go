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
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-router-go-sdk/ibmcloudlogsroutingv0"
)

func ResourceIBMLogsRouterTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMLogsRouterTenantCreate,
		ReadContext:   resourceIBMLogsRouterTenantRead,
		UpdateContext: resourceIBMLogsRouterTenantUpdate,
		DeleteContext: resourceIBMLogsRouterTenantDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
				ID := d.Id()
				parts := strings.Split(ID, "/")
				if len(parts) < 2 {
					return nil, fmt.Errorf("Invalid import format: please specify 'tenant_id/region'")
				}
				tenantID := parts[0]
				region := parts[1]
				d.SetId(tenantID)
				if err := d.Set("region", region); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_router_tenant", "name"),
				Description:  "The name for this tenant. The name is regionally unique across all tenants in the account.",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_router_tenant", "region"),
				Description:  "The region where this tenant exists.",
			},
			"targets": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "List of targets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the target.",
						},
						"log_sink_crn": &schema.Schema{
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: "Cloud resource name of the log-sink target instance.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this tenant target. The name is unique across all targets for this tenant.",
						},
						"etag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource version identifier.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							MaxItems:    1,
							Required:    true,
							Description: "List of properties returned from a successful list operation for a log-sink of type IBM Log Analysis (logdna).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Host name of the log-sink.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Network port of the log-sink.",
									},
									"access_credential": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Secret to connect to the log-sink",
									},
								},
							},
						},
					},
				},
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
			"etag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource version identifier.",
			},
		},
	}
}

func ResourceIBMLogsRouterTenantValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[a-z,A-Z,0-9,-,.]`,
			MinValueLength:             1,
			MaxValueLength:             35,
		},
		validate.ValidateSchema{
			Identifier:                 "region",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 -._:]+$`,
			MinValueLength:             3,
			MaxValueLength:             1000,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_router_tenant", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMLogsRouterTenantCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createTenantOptions := &ibmcloudlogsroutingv0.CreateTenantOptions{}

	createTenantOptions.SetName(d.Get("name").(string))
	createTenantOptions.SetRegion(d.Get("region").(string))
	var targets []ibmcloudlogsroutingv0.TargetTypePrototypeIntf
	for _, v := range d.Get("targets").([]interface{}) {
		value := v.(map[string]interface{})
		targetsItem, err := ResourceIBMLogsRouterTenantMapToTargetTypePrototype(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "create", "parse-targets").GetDiag()
		}
		targets = append(targets, targetsItem)
	}
	createTenantOptions.SetTargets(targets)

	tenant, _, err := ibmCloudLogsRoutingClient.CreateTenantWithContext(context, createTenantOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTenantWithContext failed: %s", err.Error()), "ibm_logs_router_tenant", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(tenant.ID.String())

	return resourceIBMLogsRouterTenantRead(context, d, meta)
}

func resourceIBMLogsRouterTenantRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTenantDetailOptions := &ibmcloudlogsroutingv0.GetTenantDetailOptions{}

	tenantId := strfmt.UUID(d.Id())
	getTenantDetailOptions.SetTenantID(&tenantId)
	getTenantDetailOptions.SetRegion(d.Get("region").(string))

	tenant, response, err := ibmCloudLogsRoutingClient.GetTenantDetailWithContext(context, getTenantDetailOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTenantDetailWithContext failed: %s", err.Error()), "ibm_logs_router_tenant", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", tenant.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-name").GetDiag()
	}
	targets := []map[string]interface{}{}
	for _, targetsItem := range tenant.Targets {
		targetsItemMap, err := ResourceIBMLogsRouterTenantTargetTypeToMap(targetsItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "targets-to-map").GetDiag()
		}
		targets = append(targets, targetsItemMap)
	}

	saveCredsTarget0 := d.Get("targets.0.parameters.0.access_credential").(string)
	saveCredsTarget1 := d.Get("targets.1.parameters.0.access_credential").(string)
	if len(targets) == 2 {
		if d.Get("targets.1.type").(string) == "logdna" {
			targets[0], targets[1] = targets[1], targets[0]
		}
	}

	if err = d.Set("targets", targets); err != nil {
		err = fmt.Errorf("Error setting targets: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-targets").GetDiag()
	}
	if err = d.Set("created_at", tenant.CreatedAt); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("updated_at", tenant.UpdatedAt); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("crn", tenant.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-crn").GetDiag()
	}
	if err = d.Set("etag", tenant.Etag); err != nil {
		err = fmt.Errorf("Error setting etag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-etag").GetDiag()
	}

	rewriteAccessCredential := false
	newCRNTarget0, check := d.GetOk("targets.0.log_sink_crn")
	if check {
		if crn, ok := newCRNTarget0.(string); ok && crn != "" {
			if strings.Contains(crn, ":logdna:") {
				model := &ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype{}
				hostTarget0 := d.Get("targets.0.parameters.0.host").(string)
				portTarget0 := int64(d.Get("targets.0.parameters.0.port").(int))
				model.Host = &hostTarget0
				model.Port = &portTarget0
				model.AccessCredential = &saveCredsTarget0
				parameters0Map, err := ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMapAccessCredential(model)
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-access_credential").GetDiag()
				}
				targets[0]["parameters"] = []map[string]interface{}{parameters0Map}
				rewriteAccessCredential = true
			}
		}
	}

	newCRNTarget1, check := d.GetOk("targets.1.log_sink_crn")
	if check {
		if crn, ok := newCRNTarget1.(string); ok && crn != "" {
			if strings.Contains(crn, ":logdna:") {
				model := &ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype{}
				hostTarget1 := d.Get("targets.1.parameters.0.host").(string)
				portTarget1 := int64(d.Get("targets.1.parameters.0.port").(int))
				model.Host = &hostTarget1
				model.Port = &portTarget1
				model.AccessCredential = &saveCredsTarget1
				parameters1Map, err := ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMapAccessCredential(model)
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-access_credential").GetDiag()
				}
				targets[1]["parameters"] = []map[string]interface{}{parameters1Map}
				rewriteAccessCredential = true
			}
		}
	}

	if rewriteAccessCredential {
		if err = d.Set("targets", targets); err != nil {
			err = fmt.Errorf("Error setting targets: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "read", "set-targets").GetDiag()
		}
	}

	return nil
}

func resourceIBMLogsRouterTenantUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateTenantOptions := &ibmcloudlogsroutingv0.UpdateTenantOptions{}

	tenantId := strfmt.UUID(d.Id())
	updateTenantOptions.SetTenantID(&tenantId)
	updateTenantOptions.SetRegion(d.Get("region").(string))

	hasChange := false
	hasChangeTarget0 := false
	hasChangeTarget1 := false

	patchVals := &ibmcloudlogsroutingv0.TenantPatch{}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	updateTenantOptions.SetIfMatch(d.Get("etag").(string))

	updateTarget0Options := &ibmcloudlogsroutingv0.UpdateTargetOptions{}
	target0ID := strfmt.UUID(d.Get("targets.0.id").(string))
	updateTarget0Options.SetTenantID(&tenantId)
	updateTarget0Options.SetTargetID(&target0ID)
	updateTarget0Options.SetRegion(d.Get("region").(string))

	patchValsTarget0 := &ibmcloudlogsroutingv0.TargetTypePatch{}
	if d.HasChange("targets.0.tenant_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "tenant_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_logs_router_tenant", "update", "tenant_id-forces-new").GetDiag()
	}
	if d.HasChange("targets.0.log_sink_crn") {
		newLogSinkCRN := d.Get("targets.0.log_sink_crn").(string)
		patchValsTarget0.LogSinkCRN = &newLogSinkCRN
		hasChangeTarget0 = true
	}
	if d.HasChange("targets.0.name") {
		newName := d.Get("targets.0.name").(string)
		patchValsTarget0.Name = &newName
		hasChangeTarget0 = true
	}
	if d.HasChange("targets.0.parameters") {
		parameters, err := ResourceIBMLogsRouterTargetMapToTargetParametersTypeLogDNAPrototype(d.Get("targets.0.parameters.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "update", "parse-parameters").GetDiag()
		}
		patchValsTarget0.Parameters = parameters
		hasChangeTarget0 = true
	}
	updateTarget0Options.SetIfMatch(d.Get("targets.0.etag").(string))

	updateTarget1Options := &ibmcloudlogsroutingv0.UpdateTargetOptions{}
	target1ID := strfmt.UUID(d.Get("targets.1.id").(string))
	updateTarget1Options.SetTenantID(&tenantId)
	updateTarget1Options.SetTargetID(&target1ID)
	updateTarget1Options.SetRegion(d.Get("region").(string))

	bodyModelMap := map[string]interface{}{}
	createTarget1Options := &ibmcloudlogsroutingv0.CreateTargetOptions{}
	target1Create := false
	if d.Get("targets.1.id").(string) == "" && d.Get("targets.1.log_sink_crn").(string) != "" {
		target1Create = true
		if _, ok := d.GetOk("targets.1.log_sink_crn"); ok {
			bodyModelMap["log_sink_crn"] = d.Get("targets.1.log_sink_crn")
		}
		if _, ok := d.GetOk("targets.1.name"); ok {
			bodyModelMap["name"] = d.Get("targets.1.name")
		}
		if _, ok := d.GetOk("targets.1.parameters"); ok {
			bodyModelMap["parameters"] = d.Get("targets.1.parameters")
		}
		createTarget1Options.SetTenantID(&tenantId)
		createTarget1Options.SetRegion(d.Get("region").(string))
		convertedModel, err := ResourceIBMLogsRouterTargetMapToTargetTypePrototypeTargetTypeLogDNAPrototype(bodyModelMap)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_target", "create", "parse-request-body").GetDiag()
		}
		createTarget1Options.TargetTypePrototype = convertedModel
	}

	patchValsTarget1 := &ibmcloudlogsroutingv0.TargetTypePatch{}
	if d.HasChange("targets.1.tenant_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "tenant_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_logs_router_tenant", "update", "tenant_id-forces-new").GetDiag()
	}
	if d.HasChange("targets.1.log_sink_crn") {
		newLogSinkCRN := d.Get("targets.1.log_sink_crn").(string)
		patchValsTarget1.LogSinkCRN = &newLogSinkCRN
		hasChangeTarget1 = true
	}
	if d.HasChange("targets.1.name") {
		newName := d.Get("targets.1.name").(string)
		patchValsTarget1.Name = &newName
		hasChangeTarget1 = true
	}
	if d.HasChange("targets.1.parameters") {
		parameters, err := ResourceIBMLogsRouterTargetMapToTargetParametersTypeLogDNAPrototype(d.Get("targets.1.parameters.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "update", "parse-parameters").GetDiag()
		}
		patchValsTarget1.Parameters = parameters
		hasChangeTarget1 = true
	}
	updateTarget1Options.SetIfMatch(d.Get("targets.1.etag").(string))

	if hasChange {
		updateTenantOptions.TenantPatch, _ = patchVals.AsPatch()

		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments.
		if _, exists := d.GetOk("name"); d.HasChange("name") && !exists {
			updateTenantOptions.TenantPatch["name"] = nil
		}

		_, _, err = ibmCloudLogsRoutingClient.UpdateTenantWithContext(context, updateTenantOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTenantWithContext failed: %s", err.Error()), "ibm_logs_router_tenant", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	target0Patched := false

	// second target will be created either logDNA or logs, but first target must be changed first
	// if there is a change to the existing target1, we have to apply that first either delete or udpate, then update target 0
	if target1Create && hasChangeTarget0 {
		updateTarget0Options.TargetTypePatch, _ = patchValsTarget0.AsPatch()

		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments.
		if _, exists := d.GetOk("targets.0.log_sink_crn"); d.HasChange("targets.0.log_sink_crn") && !exists {
			updateTarget0Options.TargetTypePatch["log_sink_crn"] = nil
		}
		if _, exists := d.GetOk("targets.0.name"); d.HasChange("targets.0.name") && !exists {
			updateTarget0Options.TargetTypePatch["name"] = nil
		}
		if _, exists := d.GetOk("targets.0.parameters"); d.HasChange("targets.0.parameters") && !exists {
			updateTarget0Options.TargetTypePatch["parameters"] = nil
		}
		target0Patched = true

		_, newCRN := d.GetChange("targets.0.log_sink_crn")
		if crn, ok := newCRN.(string); ok {
			if strings.Contains(crn, ":logs:") {
				_, _, err = ibmCloudLogsRoutingClient.UpdateLogsTargetWithContext(context, updateTarget0Options)
			} else {
				_, _, err = ibmCloudLogsRoutingClient.UpdateTargetWithContext(context, updateTarget0Options)
			}
		}

		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTargetWithContext failed: %s", err.Error()), "ibm_logs_router_target", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if target1Create {
		_, _, err := ibmCloudLogsRoutingClient.CreateTargetWithContext(context, createTarget1Options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTargetWithContext failed: %s", err.Error()), "ibm_logs_router_target", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	} else if hasChangeTarget1 {
		target1Delete := false
		updateTarget1Options.TargetTypePatch, _ = patchValsTarget1.AsPatch()

		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments.
		if _, exists := d.GetOk("targets.1.log_sink_crn"); d.HasChange("targets.1.log_sink_crn") && !exists {
			updateTarget1Options.TargetTypePatch["log_sink_crn"] = nil
		}
		if _, exists := d.GetOk("targets.1.name"); d.HasChange("targets.1.name") && !exists {
			updateTarget1Options.TargetTypePatch["name"] = nil
		}
		if _, exists := d.GetOk("targets.1.parameters"); d.HasChange("targets.1.parameters") && !exists {
			updateTarget1Options.TargetTypePatch["parameters"] = nil
			target1Delete = true
		}
		if target1Delete {
			deleteTargetOptions := &ibmcloudlogsroutingv0.DeleteTargetOptions{}
			deleteTargetOptions.SetTenantID(&tenantId)
			deleteTargetOptions.SetTargetID(&target1ID)
			deleteTargetOptions.SetRegion(d.Get("region").(string))
			_, err = ibmCloudLogsRoutingClient.DeleteTargetWithContext(context, deleteTargetOptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTargetWithContext failed: %s", err.Error()), "ibm_logs_router_target", "delete")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}

		} else {
			_, newCRN := d.GetChange("targets.1.log_sink_crn")
			if crn, ok := newCRN.(string); ok {
				if strings.Contains(crn, ":logs:") {
					_, _, err = ibmCloudLogsRoutingClient.UpdateLogsTargetWithContext(context, updateTarget1Options)
				} else {
					_, _, err = ibmCloudLogsRoutingClient.UpdateTargetWithContext(context, updateTarget1Options)
				}
			}
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTargetWithContext failed: %s", err.Error()), "ibm_logs_router_target", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
	}

	if hasChangeTarget0 && !target0Patched {
		updateTarget0Options.TargetTypePatch, _ = patchValsTarget0.AsPatch()

		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments.
		if _, exists := d.GetOk("targets.0.log_sink_crn"); d.HasChange("targets.0.log_sink_crn") && !exists {
			updateTarget0Options.TargetTypePatch["log_sink_crn"] = nil
		}
		if _, exists := d.GetOk("targets.0.name"); d.HasChange("targets.0.name") && !exists {
			updateTarget0Options.TargetTypePatch["name"] = nil
		}
		if _, exists := d.GetOk("targets.0.parameters"); d.HasChange("targets.0.parameters") && !exists {
			updateTarget0Options.TargetTypePatch["parameters"] = nil
		}
		target0Patched = true
		// update logs RM access credential
		_, newCRN := d.GetChange("targets.0.log_sink_crn")
		if crn, ok := newCRN.(string); ok {
			if strings.Contains(crn, ":logs:") {
				_, _, err = ibmCloudLogsRoutingClient.UpdateLogsTargetWithContext(context, updateTarget0Options)
			} else {
				_, _, err = ibmCloudLogsRoutingClient.UpdateTargetWithContext(context, updateTarget0Options)
			}
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTargetWithContext failed: %s", err.Error()), "ibm_logs_router_target", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMLogsRouterTenantRead(context, d, meta)
}

func resourceIBMLogsRouterTenantDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ibmCloudLogsRoutingClient, err := meta.(conns.ClientSession).IBMCloudLogsRoutingV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_router_tenant", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTenantOptions := &ibmcloudlogsroutingv0.DeleteTenantOptions{}

	tenantId := strfmt.UUID(d.Id())
	deleteTenantOptions.SetTenantID(&tenantId)
	deleteTenantOptions.SetRegion(d.Get("region").(string))

	_, err = ibmCloudLogsRoutingClient.DeleteTenantWithContext(context, deleteTenantOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTenantWithContext failed: %s", err.Error()), "ibm_logs_router_tenant", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMLogsRouterTenantMapToTargetTypePrototype(modelMap map[string]interface{}) (ibmcloudlogsroutingv0.TargetTypePrototypeIntf, error) {
	model := &ibmcloudlogsroutingv0.TargetTypePrototype{}
	if modelMap["log_sink_crn"] != nil && modelMap["log_sink_crn"].(string) != "" {
		model.LogSinkCRN = core.StringPtr(modelMap["log_sink_crn"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["parameters"] != nil && len(modelMap["parameters"].([]interface{})) > 0 {
		ParametersModel, err := ResourceIBMLogsRouterTenantMapToTargetParametersTypeLogDnaPrototype(modelMap["parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}
	return model, nil
}

func ResourceIBMLogsRouterTenantMapToTargetParametersTypeLogDnaPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype{}
	model.Host = core.StringPtr(modelMap["host"].(string))
	model.Port = core.Int64Ptr(int64(modelMap["port"].(int)))
	model.AccessCredential = core.StringPtr(modelMap["access_credential"].(string))
	return model, nil
}

func ResourceIBMLogsRouterTenantMapToTargetTypePrototypeTargetTypeLogDnaPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogDnaPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogDnaPrototype{}
	model.LogSinkCRN = core.StringPtr(modelMap["log_sink_crn"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["parameters"] != nil && len(modelMap["parameters"].([]interface{})) > 0 {
		ParametersModel, err := ResourceIBMLogsRouterTenantMapToTargetParametersTypeLogDnaPrototype(modelMap["parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}
	return model, nil
}

func ResourceIBMLogsRouterTenantMapToTargetTypePrototypeTargetTypeLogsPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogsPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogsPrototype{}
	model.LogSinkCRN = core.StringPtr(modelMap["log_sink_crn"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["parameters"] != nil && len(modelMap["parameters"].([]interface{})) > 0 {
		ParametersModel, err := ResourceIBMLogsRouterTenantMapToTargetParametersTypeLogsPrototype(modelMap["parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}
	return model, nil
}

func ResourceIBMLogsRouterTenantMapToTargetParametersTypeLogsPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetParametersTypeLogsPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetParametersTypeLogsPrototype{}
	model.Host = core.StringPtr(modelMap["host"].(string))
	model.Port = core.Int64Ptr(int64(modelMap["port"].(int)))
	return model, nil
}

func ResourceIBMLogsRouterTenantTargetTypeToMap(model ibmcloudlogsroutingv0.TargetTypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogDna); ok {
		return ResourceIBMLogsRouterTenantTargetTypeLogDnaToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogDna))
	} else if _, ok := model.(*ibmcloudlogsroutingv0.TargetTypeLogs); ok {
		return ResourceIBMLogsRouterTenantTargetTypeLogsToMap(model.(*ibmcloudlogsroutingv0.TargetTypeLogs))
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
			parametersMap, err := ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMap(model.Parameters)
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

func ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}

func ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMapAccessCredential(model *ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	modelMap["access_credential"] = *model.AccessCredential // pragma: whitelist secret
	return modelMap, nil
}

func ResourceIBMLogsRouterTenantTargetTypeLogDnaToMap(model *ibmcloudlogsroutingv0.TargetTypeLogDna) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := ResourceIBMLogsRouterTenantTargetParametersTypeLogDnaToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func ResourceIBMLogsRouterTenantTargetTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID.String()
	modelMap["log_sink_crn"] = *model.LogSinkCRN
	modelMap["name"] = *model.Name
	modelMap["etag"] = *model.Etag
	modelMap["type"] = *model.Type
	modelMap["created_at"] = *model.CreatedAt
	modelMap["updated_at"] = *model.UpdatedAt
	if model.Parameters != nil {
		parametersMap, err := ResourceIBMLogsRouterTenantTargetParametersTypeLogsToMap(model.Parameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["parameters"] = []map[string]interface{}{parametersMap}
	}
	return modelMap, nil
}

func ResourceIBMLogsRouterTenantTargetParametersTypeLogsToMap(model *ibmcloudlogsroutingv0.TargetParametersTypeLogs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["host"] = *model.Host
	modelMap["port"] = flex.IntValue(model.Port)
	return modelMap, nil
}

func ResourceIBMLogsRouterTargetMapToTargetParametersTypeLogDNAPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetParametersTypeLogDnaPrototype{}
	model.Host = core.StringPtr(modelMap["host"].(string))
	model.Port = core.Int64Ptr(int64(modelMap["port"].(int)))
	model.AccessCredential = core.StringPtr(modelMap["access_credential"].(string))
	return model, nil
}

func ResourceIBMLogsRouterTargetMapToTargetTypePrototypeTargetTypeLogDNAPrototype(modelMap map[string]interface{}) (*ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogDnaPrototype, error) {
	model := &ibmcloudlogsroutingv0.TargetTypePrototypeTargetTypeLogDnaPrototype{}
	model.LogSinkCRN = core.StringPtr(modelMap["log_sink_crn"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["parameters"] != nil && len(modelMap["parameters"].([]interface{})) > 0 {
		ParametersModel, err := ResourceIBMLogsRouterTargetMapToTargetParametersTypeLogDNAPrototype(modelMap["parameters"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Parameters = ParametersModel
	}
	return model, nil
}
