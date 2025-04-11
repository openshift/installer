// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
	"github.com/go-openapi/strfmt"
)

func ResourceIbmLogsDashboardFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsDashboardFolderCreate,
		ReadContext:   resourceIbmLogsDashboardFolderRead,
		UpdateContext: resourceIbmLogsDashboardFolderUpdate,
		DeleteContext: resourceIbmLogsDashboardFolderDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_dashboard_folder", "name"),
				Description:  "The dashboard folder name, required.",
			},
			"parent_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The dashboard folder parent ID, optional. If not set, the folder is a root folder, if set, the folder is a subfolder of the parent folder and needs to be a uuid.",
			},
			"dashboard_folder_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dashboard Folder Id.",
			},
		},
	}
}

func ResourceIbmLogsDashboardFolderValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_dashboard_folder", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsDashboardFolderCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard_folder", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	createDashboardFolderOptions := &logsv0.CreateDashboardFolderOptions{}

	createDashboardFolderOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("parent_id"); ok {
		createDashboardFolderOptions.SetParentID(core.UUIDPtr(strfmt.UUID(d.Get("parent_id").(string))))
	}

	dashboardFolder, _, err := logsClient.CreateDashboardFolderWithContext(context, createDashboardFolderOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDashboardFolderWithContext failed: %s", err.Error()), "ibm_logs_dashboard_folder", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	dashboardFolderId := fmt.Sprintf("%s/%s/%s", region, instanceId, *dashboardFolder.ID)
	d.SetId(dashboardFolderId)

	return resourceIbmLogsDashboardFolderRead(context, d, meta)
}

func resourceIbmLogsDashboardFolderRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, dashboardFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	listDashboardFoldersOptions := &logsv0.ListDashboardFoldersOptions{}

	dashboardFolderCollection, response, err := logsClient.ListDashboardFoldersWithContext(context, listDashboardFoldersOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListDashboardFoldersWithContext failed: %s", err.Error()), "ibm_logs_dashboard_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("dashboard_folder_id", dashboardFolderId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dashboard_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if dashboardFolderCollection != nil && len(dashboardFolderCollection.Folders) > 0 {

		for _, folder := range dashboardFolderCollection.Folders {

			if fmt.Sprintf("%s", folder.ID) == dashboardFolderId {
				if err = d.Set("name", folder.Name); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
				}
				if !core.IsNil(folder.ParentID) {
					if err = d.Set("parent_id", folder.ParentID); err != nil {
						return diag.FromErr(fmt.Errorf("Error setting parent_id: %s", err))
					}
				}
			}
		}
	}

	return nil
}

func resourceIbmLogsDashboardFolderUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard_folder", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, dashboardFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	replaceDashboardFolderOptions := &logsv0.ReplaceDashboardFolderOptions{}

	replaceDashboardFolderOptions.SetFolderID(core.UUIDPtr(strfmt.UUID(dashboardFolderId)))

	hasChange := false

	if d.HasChange("name") || d.HasChange("parent_id") {

		replaceDashboardFolderOptions.SetName(d.Get("name").(string))
		if _, ok := d.GetOk("parent_id"); ok {
			replaceDashboardFolderOptions.SetParentID(core.UUIDPtr(strfmt.UUID(d.Get("parent_id").(string))))
		}
		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.ReplaceDashboardFolderWithContext(context, replaceDashboardFolderOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceDashboardFolderWithContext failed: %s", err.Error()), "ibm_logs_dashboard_folder", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsDashboardFolderRead(context, d, meta)
}

func resourceIbmLogsDashboardFolderDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_dashboard_folder", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, dashboardFolderId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDashboardFolderOptions := &logsv0.DeleteDashboardFolderOptions{}

	deleteDashboardFolderOptions.SetFolderID(core.UUIDPtr(strfmt.UUID(dashboardFolderId)))

	_, err = logsClient.DeleteDashboardFolderWithContext(context, deleteDashboardFolderOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDashboardFolderWithContext failed: %s", err.Error()), "ibm_logs_dashboard_folder", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
