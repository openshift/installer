// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabaseBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMDatabaseBackupsRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the deployment this backup relates to.",
			},
			"backups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An array of backups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of this backup.",
						},
						"deployment_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the deployment this backup relates to.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of backup.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this backup.",
						},
						"is_downloadable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is this backup available to download?.",
						},
						"is_restorable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Can this backup be used to restore an instance?.",
						},
						"download_link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI which is currently available for file downloading.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time when this backup was created.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMDatabaseBackupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	listDeploymentBackupsOptions := &clouddatabasesv5.ListDeploymentBackupsOptions{}
	listDeploymentBackupsOptions.SetID(d.Get("deployment_id").(string))

	backups, response, err := cloudDatabasesClient.ListDeploymentBackupsWithContext(context, listDeploymentBackupsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDeploymentBackupsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListDeploymentBackupsWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchBackups []clouddatabasesv5.Backup
	var deploymentID string
	var suppliedFilter bool

	if v, ok := d.GetOk("deployment_id"); ok {
		deploymentID = v.(string)
		suppliedFilter = true
		for _, data := range backups.Backups {
			if *data.DeploymentID == deploymentID {
				matchBackups = append(matchBackups, data)
			}
		}
	} else {
		matchBackups = backups.Backups
	}
	backups.Backups = matchBackups

	if suppliedFilter {
		if len(backups.Backups) == 0 {
			return diag.FromErr(fmt.Errorf("no Backups found with deploymentID %s", deploymentID))
		}
		d.SetId(deploymentID)
	} else {
		d.SetId(DataSourceIBMDatabaseBackupsID(d))
	}

	backups2 := []map[string]interface{}{}
	if backups.Backups != nil {
		for _, modelItem := range backups.Backups {
			modelMap, err := DataSourceIBMDatabaseBackupsBackupToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			backups2 = append(backups2, modelMap)
		}
	}
	if err = d.Set("backups", backups2); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting backups %s", err))
	}

	return nil
}

// DataSourceIBMDatabaseBackupsID returns a reasonable ID for the list.
func DataSourceIBMDatabaseBackupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMDatabaseBackupsBackupToMap(model *clouddatabasesv5.Backup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["backup_id"] = *model.ID
	}
	if model.DeploymentID != nil {
		modelMap["deployment_id"] = *model.DeploymentID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.IsDownloadable != nil {
		modelMap["is_downloadable"] = *model.IsDownloadable
	}
	if model.IsRestorable != nil {
		modelMap["is_restorable"] = *model.IsRestorable
	}
	if model.DownloadLink != nil {
		modelMap["download_link"] = *model.DownloadLink
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	return modelMap, nil
}
