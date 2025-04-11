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

func DataSourceIbmLogsDashboardFolders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsDashboardFoldersRead,

		Schema: map[string]*schema.Schema{
			"folders": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dashboard folder ID, uuid.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dashboard folder name, required.",
						},
						"parent_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dashboard folder parent ID, optional. If not set, the folder is a root folder, if set, the folder is a subfolder of the parent folder and needs to be a uuid.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsDashboardFoldersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard_folders", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	listDashboardFoldersOptions := &logsv0.ListDashboardFoldersOptions{}

	dashboardFolderCollection, _, err := logsClient.ListDashboardFoldersWithContext(context, listDashboardFoldersOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListDashboardFoldersWithContext failed: %s", err.Error()), "(Data) ibm_logs_dashboard_folders", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsDashboardFoldersID(d))

	folders := []map[string]interface{}{}
	if dashboardFolderCollection.Folders != nil {
		for _, modelItem := range dashboardFolderCollection.Folders {
			modelMap, err := DataSourceIbmLogsDashboardFoldersDashboardFolderToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_dashboard_folders", "read")
				return tfErr.GetDiag()
			}
			folders = append(folders, modelMap)
		}
	}
	if err = d.Set("folders", folders); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting folders: %s", err), "(Data) ibm_logs_dashboard_folders", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsDashboardFoldersID returns a reasonable ID for the list.
func dataSourceIbmLogsDashboardFoldersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsDashboardFoldersDashboardFolderToMap(model *logsv0.DashboardFolder) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID.String()
	}
	modelMap["name"] = *model.Name
	if model.ParentID != nil {
		modelMap["parent_id"] = model.ParentID.String()
	}
	return modelMap, nil
}
