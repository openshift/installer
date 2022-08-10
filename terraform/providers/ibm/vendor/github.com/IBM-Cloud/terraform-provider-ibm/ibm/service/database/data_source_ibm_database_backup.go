// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"backup_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Backup ID.",
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
	}
}

func DataSourceIBMDatabaseBackupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupInfoOptions := &clouddatabasesv5.GetBackupInfoOptions{}

	getBackupInfoOptions.SetBackupID(d.Get("backup_id").(string))

	backup, response, err := cloudDatabasesClient.GetBackupInfoWithContext(context, getBackupInfoOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBackupInfoWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBackupInfoWithContext failed %s\n%s", err, response))
	}

	d.SetId(*backup.Backup.ID)

	if err = d.Set("backup_id", backup.Backup.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if err = d.Set("deployment_id", backup.Backup.DeploymentID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting deployment_id: %s", err))
	}

	if err = d.Set("type", backup.Backup.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("status", backup.Backup.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if err = d.Set("is_downloadable", backup.Backup.IsDownloadable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_downloadable: %s", err))
	}

	if err = d.Set("is_restorable", backup.Backup.IsRestorable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_restorable: %s", err))
	}

	if err = d.Set("download_link", backup.Backup.DownloadLink); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting download_link: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(backup.Backup.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	return nil
}
