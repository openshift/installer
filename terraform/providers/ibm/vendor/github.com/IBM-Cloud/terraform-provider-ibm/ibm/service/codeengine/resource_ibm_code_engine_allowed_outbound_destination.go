// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIbmCodeEngineAllowedOutboundDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineAllowedOutboundDestinationCreate,
		ReadContext:   resourceIbmCodeEngineAllowedOutboundDestinationRead,
		UpdateContext: resourceIbmCodeEngineAllowedOutboundDestinationUpdate,
		DeleteContext: resourceIbmCodeEngineAllowedOutboundDestinationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_allowed_outbound_destination", "project_id"),
				Description:  "The ID of the project.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_allowed_outbound_destination", "type"),
				Description:  "Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_allowed_outbound_destination", "cidr_block"),
				Description:  "The IPv4 address range.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_allowed_outbound_destination", "name"),
				Description:  "The name of the CIDR block.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the allowed outbound destination, which is used to achieve optimistic locking.",
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineAllowedOutboundDestinationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "project_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "cidr_block",
			Regexp:                     `^(cidr_block)$`,
		},
		validate.ValidateSchema{
			Identifier:                 "cidr_block",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$`,
			MinValueLength:             0,
			MaxValueLength:             18,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z]([-a-z0-9]*[a-z0-9])?$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_allowed_outbound_destination", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineAllowedOutboundDestinationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bodyModelMap := map[string]interface{}{}
	createAllowedOutboundDestinationOptions := &codeenginev2.CreateAllowedOutboundDestinationOptions{}

	bodyModelMap["type"] = d.Get("type")
	if _, ok := d.GetOk("cidr_block"); ok {
		bodyModelMap["cidr_block"] = d.Get("cidr_block")
	}
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	createAllowedOutboundDestinationOptions.SetProjectID(d.Get("project_id").(string))
	convertedModel, err := ResourceIbmCodeEngineAllowedOutboundDestinationMapToAllowedOutboundDestinationPrototype(bodyModelMap)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "create", "parse-request-body").GetDiag()
	}
	createAllowedOutboundDestinationOptions.AllowedOutboundDestination = convertedModel

	allowedOutboundDestinationIntf, _, err := codeEngineClient.CreateAllowedOutboundDestinationWithContext(context, createAllowedOutboundDestinationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateAllowedOutboundDestinationWithContext failed: %s", err.Error()), "ibm_code_engine_allowed_outbound_destination", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allowedOutboundDestination := allowedOutboundDestinationIntf.(*codeenginev2.AllowedOutboundDestination)
	d.SetId(fmt.Sprintf("%s/%s", *createAllowedOutboundDestinationOptions.ProjectID, *allowedOutboundDestination.Name))

	return resourceIbmCodeEngineAllowedOutboundDestinationRead(context, d, meta)
}

func resourceIbmCodeEngineAllowedOutboundDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAllowedOutboundDestinationOptions := &codeenginev2.GetAllowedOutboundDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "sep-id-parts").GetDiag()
	}

	getAllowedOutboundDestinationOptions.SetProjectID(parts[0])
	getAllowedOutboundDestinationOptions.SetName(parts[1])

	allowedOutboundDestinationIntf, response, err := codeEngineClient.GetAllowedOutboundDestinationWithContext(context, getAllowedOutboundDestinationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAllowedOutboundDestinationWithContext failed: %s", err.Error()), "ibm_code_engine_allowed_outbound_destination", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allowedOutboundDestination := allowedOutboundDestinationIntf.(*codeenginev2.AllowedOutboundDestination)
	if err = d.Set("type", allowedOutboundDestination.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "set-type").GetDiag()
	}
	if !core.IsNil(allowedOutboundDestination.CidrBlock) {
		if err = d.Set("cidr_block", allowedOutboundDestination.CidrBlock); err != nil {
			err = fmt.Errorf("Error setting cidr_block: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "set-cidr_block").GetDiag()
		}
	}
	if !core.IsNil(allowedOutboundDestination.Name) {
		if err = d.Set("name", allowedOutboundDestination.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("entity_tag", allowedOutboundDestination.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "read", "set-entity_tag").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_allowed_outbound_destination", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineAllowedOutboundDestinationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateAllowedOutboundDestinationOptions := &codeenginev2.UpdateAllowedOutboundDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "update", "sep-id-parts").GetDiag()
	}

	updateAllowedOutboundDestinationOptions.SetProjectID(parts[0])
	updateAllowedOutboundDestinationOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.AllowedOutboundDestinationPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_allowed_outbound_destination", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("type") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "type")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_allowed_outbound_destination", "update", "type-forces-new").GetDiag()
	}
	typeVal := d.Get("type").(string)
	patchVals.Type = &typeVal

	if d.HasChange("cidr_block") {
		newCidrBlock := d.Get("cidr_block").(string)
		patchVals.CidrBlock = &newCidrBlock
		hasChange = true
	}
	updateAllowedOutboundDestinationOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateAllowedOutboundDestinationOptions.AllowedOutboundDestination = ResourceIbmCodeEngineAllowedOutboundDestinationAllowedOutboundDestinationPatchAsPatch(patchVals, d)

		_, _, err = codeEngineClient.UpdateAllowedOutboundDestinationWithContext(context, updateAllowedOutboundDestinationOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateAllowedOutboundDestinationWithContext failed: %s", err.Error()), "ibm_code_engine_allowed_outbound_destination", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineAllowedOutboundDestinationRead(context, d, meta)
}

func resourceIbmCodeEngineAllowedOutboundDestinationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteAllowedOutboundDestinationOptions := &codeenginev2.DeleteAllowedOutboundDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_allowed_outbound_destination", "delete", "sep-id-parts").GetDiag()
	}

	deleteAllowedOutboundDestinationOptions.SetProjectID(parts[0])
	deleteAllowedOutboundDestinationOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteAllowedOutboundDestinationWithContext(context, deleteAllowedOutboundDestinationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteAllowedOutboundDestinationWithContext failed: %s", err.Error()), "ibm_code_engine_allowed_outbound_destination", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmCodeEngineAllowedOutboundDestinationMapToAllowedOutboundDestinationPrototype(modelMap map[string]interface{}) (codeenginev2.AllowedOutboundDestinationPrototypeIntf, error) {
	model := &codeenginev2.AllowedOutboundDestinationPrototype{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	if modelMap["cidr_block"] != nil && modelMap["cidr_block"].(string) != "" {
		model.CidrBlock = core.StringPtr(modelMap["cidr_block"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmCodeEngineAllowedOutboundDestinationMapToAllowedOutboundDestinationPrototypeCidrBlockDataPrototype(modelMap map[string]interface{}) (*codeenginev2.AllowedOutboundDestinationPrototypeCidrBlockDataPrototype, error) {
	model := &codeenginev2.AllowedOutboundDestinationPrototypeCidrBlockDataPrototype{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.CidrBlock = core.StringPtr(modelMap["cidr_block"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIbmCodeEngineAllowedOutboundDestinationAllowedOutboundDestinationPatchAsPatch(patchVals *codeenginev2.AllowedOutboundDestinationPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "cidr_block"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["cidr_block"] = nil
	}

	return patch
}
