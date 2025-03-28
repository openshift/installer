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

func DataSourceIbmBackupRecoveryDownloadFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryDownloadFilesRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies tenantId",
			},
			"recovery_download_files_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the id of a Recovery.",
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
			"file_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the downloaded type, i.e: error, success_files_list.",
			},
			"source_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the source on which restore is done.",
			},
			"start_time": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the start time of restore task.",
			},
			"include_tenants": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies if objects of all the organizations under the hierarchy of the logged in user's organization should be returned.",
			},
		},
	}
}

func dataSourceIbmBackupRecoveryDownloadFilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery_download_files", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	downloadFilesFromRecoveryOptions := &backuprecoveryv1.DownloadFilesFromRecoveryOptions{}

	downloadFilesFromRecoveryOptions.SetID(d.Get("recovery_download_files_id").(string))
	downloadFilesFromRecoveryOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("start_offset"); ok {
		downloadFilesFromRecoveryOptions.SetStartOffset(int64(d.Get("start_offset").(int)))
	}
	if _, ok := d.GetOk("length"); ok {
		downloadFilesFromRecoveryOptions.SetLength(int64(d.Get("length").(int)))
	}
	if _, ok := d.GetOk("file_type"); ok {
		downloadFilesFromRecoveryOptions.SetFileType(d.Get("file_type").(string))
	}
	if _, ok := d.GetOk("source_name"); ok {
		downloadFilesFromRecoveryOptions.SetSourceName(d.Get("source_name").(string))
	}
	if _, ok := d.GetOk("start_time"); ok {
		downloadFilesFromRecoveryOptions.SetStartTime(d.Get("start_time").(string))
	}
	if _, ok := d.GetOk("include_tenants"); ok {
		downloadFilesFromRecoveryOptions.SetIncludeTenants(d.Get("include_tenants").(bool))
	}

	_, err = backupRecoveryClient.DownloadFilesFromRecoveryWithContext(context, downloadFilesFromRecoveryOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DownloadFilesFromRecoveryWithContext failed: %s", err.Error()), "(Data) ibm_recovery_download_files", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryDownloadFilesID(d))

	return nil
}

// dataSourceIbmBackupRecoveryDownloadFilesID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryDownloadFilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
