// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.0-fa797aec-20240814-142622
 */

package backuprecovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func DataSourceIbmBackupRecoveryDownloadIndexedFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryDownloadIndexedFilesRead,

		Schema: map[string]*schema.Schema{
			"snapshots_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the snapshot id to download from.",
			},
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"file_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the path to the file to download. If no path is specified and snapshot environment is kVMWare, VMX file for VMware will be downloaded. For other snapshot environments, this field must be specified.",
			},
			"nvram_file": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies if NVRAM file for VMware should be downloaded.",
			},
			"retry_attempt": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of attempts the protection run took to create this file.",
			},
			"start_offset": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the start offset of file chunk to be downloaded.",
			},
			"length": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the length of bytes to download. This can not be greater than 8MB (8388608 byets).",
			},
		},
	}
}

func dataSourceIbmBackupRecoveryDownloadIndexedFilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_download_indexed_files", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	downloadIndexedFileOptions := &backuprecoveryv1.DownloadIndexedFileOptions{}

	downloadIndexedFileOptions.SetSnapshotsID(d.Get("snapshots_id").(string))
	downloadIndexedFileOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("file_path"); ok {
		downloadIndexedFileOptions.SetFilePath(d.Get("file_path").(string))
	}
	if _, ok := d.GetOk("nvram_file"); ok {
		downloadIndexedFileOptions.SetNvramFile(d.Get("nvram_file").(bool))
	}
	if _, ok := d.GetOk("retry_attempt"); ok {
		downloadIndexedFileOptions.SetRetryAttempt(int64(d.Get("retry_attempt").(int)))
	}
	if _, ok := d.GetOk("start_offset"); ok {
		downloadIndexedFileOptions.SetStartOffset(int64(d.Get("start_offset").(int)))
	}
	if _, ok := d.GetOk("length"); ok {
		downloadIndexedFileOptions.SetLength(int64(d.Get("length").(int)))
	}

	_, err = backupRecoveryClient.DownloadIndexedFileWithContext(context, downloadIndexedFileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DownloadIndexedFileWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_download_indexed_files", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryDownloadIndexedFilesID(d))

	return nil
}

// dataSourceIbmBackupRecoveryDownloadIndexedFilesID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryDownloadIndexedFilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
