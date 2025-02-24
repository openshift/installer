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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsViewFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsViewFolderRead,

		Schema: map[string]*schema.Schema{
			"logs_view_folder_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Folder ID.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Folder name.",
			},
		},
	}
}

func dataSourceIbmLogsViewFolderRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getViewFolderOptions := &logsv0.GetViewFolderOptions{}

	getViewFolderOptions.SetID(core.UUIDPtr(strfmt.UUID(d.Get("logs_view_folder_id").(string))))

	viewFolder, _, err := logsClient.GetViewFolderWithContext(context, getViewFolderOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetViewFolderWithContext failed: %s", err.Error()), "(Data) ibm_logs_view_folder", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getViewFolderOptions.ID))

	if err = d.Set("name", viewFolder.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_view_folder", "read")
		return tfErr.GetDiag()
	}

	return nil
}
