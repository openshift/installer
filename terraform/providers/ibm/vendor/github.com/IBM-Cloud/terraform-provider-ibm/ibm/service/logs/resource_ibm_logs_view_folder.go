// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsViewFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsViewFolderCreate,
		ReadContext:   resourceIbmLogsViewFolderRead,
		UpdateContext: resourceIbmLogsViewFolderUpdate,
		DeleteContext: resourceIbmLogsViewFolderDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_view_folder", "name"),
				Description:  "Folder name.",
			},
			"view_folder_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "View Folder Id.",
			},
		},
	}
}

func ResourceIbmLogsViewFolderValidator() *validate.ResourceValidator {
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_view_folder", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsViewFolderCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_view_folder", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createViewFolderOptions := &logsv0.CreateViewFolderOptions{}

	createViewFolderOptions.SetName(d.Get("name").(string))

	viewFolder, _, err := logsClient.CreateViewFolderWithContext(context, createViewFolderOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateViewFolderWithContext failed: %s", err.Error()), "ibm_logs_view_folder", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	viewFolderId := fmt.Sprintf("%s/%s/%s", region, instanceId, *viewFolder.ID)
	d.SetId(viewFolderId)

	return resourceIbmLogsViewFolderRead(context, d, meta)
}

func resourceIbmLogsViewFolderRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_view_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, viewFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getViewFolderOptions := &logsv0.GetViewFolderOptions{}

	getViewFolderOptions.SetID(core.UUIDPtr(strfmt.UUID(viewFolderId)))

	viewFolder, response, err := logsClient.GetViewFolderWithContext(context, getViewFolderOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetViewFolderWithContext failed: %s", err.Error()), "ibm_logs_view_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("view_folder_id", viewFolderId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting view_folder_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("name", viewFolder.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	return nil
}

func resourceIbmLogsViewFolderUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_view_folder", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, viewFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	replaceViewFolderOptions := &logsv0.ReplaceViewFolderOptions{}

	replaceViewFolderOptions.SetID(core.UUIDPtr(strfmt.UUID(viewFolderId)))

	hasChange := false

	if d.HasChange("name") {
		replaceViewFolderOptions.SetName(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.ReplaceViewFolderWithContext(context, replaceViewFolderOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceViewFolderWithContext failed: %s", err.Error()), "ibm_logs_view_folder", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsViewFolderRead(context, d, meta)
}

func resourceIbmLogsViewFolderDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_view_folder", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, viewFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteViewFolderOptions := &logsv0.DeleteViewFolderOptions{}

	deleteViewFolderOptions.SetID(core.UUIDPtr(strfmt.UUID(viewFolderId)))

	_, err = logsClient.DeleteViewFolderWithContext(context, deleteViewFolderOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteViewFolderWithContext failed: %s", err.Error()), "ibm_logs_view_folder", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
