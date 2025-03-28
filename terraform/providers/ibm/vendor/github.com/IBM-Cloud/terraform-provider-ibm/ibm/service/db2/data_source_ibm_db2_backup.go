// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package db2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"
)

func DataSourceIbmDb2Backup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmDb2BackupRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Encoded CRN deployment id.",
			},
			"backups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CRN of the db2 instance.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the type of execution of backup.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the backup.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp of the backup created.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the backup or data set.",
						},
						"duration": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of the backup operation in seconds.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmDb2BackupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db2saasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_backup", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDb2SaasBackupOptions := &db2saasv1.GetDb2SaasBackupOptions{}

	getDb2SaasBackupOptions.SetXDbProfile(d.Get("deployment_id").(string))

	successGetBackups, _, err := db2saasClient.GetDb2SaasBackupWithContext(context, getDb2SaasBackupOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDb2SaasBackupWithContext failed: %s", err.Error()), "(Data) ibm_db2_backup", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmDb2BackupID(d))

	backups := []map[string]interface{}{}
	for _, backupsItem := range successGetBackups.Backups {
		backupsItemMap, err := DataSourceIbmDb2BackupBackupToMap(&backupsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_backup", "read", "backups-to-map").GetDiag()
		}
		backups = append(backups, backupsItemMap)
	}
	if err = d.Set("backups", backups); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting backups: %s", err), "(Data) ibm_db2_backup", "read", "set-backups").GetDiag()
	}

	return nil
}

// dataSourceIbmDb2SaasBackupID returns a reasonable ID for the list.
func dataSourceIbmDb2BackupID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmDb2BackupBackupToMap(model *db2saasv1.Backup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["type"] = *model.Type
	modelMap["status"] = *model.Status
	modelMap["created_at"] = *model.CreatedAt
	modelMap["size"] = flex.IntValue(model.Size)
	modelMap["duration"] = flex.IntValue(model.Duration)
	return modelMap, nil
}
