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

func DataSourceIbmLogsViewFolders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsViewFoldersRead,

		Schema: map[string]*schema.Schema{
			"view_folders": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of view folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Folder ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Folder name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsViewFoldersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view_folders", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	listViewFoldersOptions := &logsv0.ListViewFoldersOptions{}

	viewFolderCollection, _, err := logsClient.ListViewFoldersWithContext(context, listViewFoldersOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListViewFoldersWithContext failed: %s", err.Error()), "(Data) ibm_logs_view_folders", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsViewFoldersID(d))

	viewFolders := []map[string]interface{}{}
	if viewFolderCollection.ViewFolders != nil {
		for _, modelItem := range viewFolderCollection.ViewFolders {
			modelMap, err := DataSourceIbmLogsViewFoldersViewFolderToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view_folders", "read")
				return tfErr.GetDiag()
			}
			viewFolders = append(viewFolders, modelMap)
		}
	}
	if err = d.Set("view_folders", viewFolders); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting view_folders: %s", err), "(Data) ibm_logs_view_folders", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsViewFoldersID returns a reasonable ID for the list.
func dataSourceIbmLogsViewFoldersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsViewFoldersViewFolderToMap(model *logsv0.ViewFolder) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID.String()
	}
	modelMap["name"] = *model.Name
	return modelMap, nil
}
