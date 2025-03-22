// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.0-fa797aec-20240814-142622
 */

package backuprecovery

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func DataSourceIbmBackupRecoveryDownloadAgent() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIbmBackupRecoveryDownloadAgentRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"file_path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the absolute path for download",
			},
			"platform": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the platform for which agent needs to be downloaded.",
			},
			"linux_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Linux agent parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the type of installer.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIbmBackupRecoveryDownloadAgentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_download_agent", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	downloadAgentOptions := &backuprecoveryv1.DownloadAgentOptions{}

	downloadAgentOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	downloadAgentOptions.SetPlatform(d.Get("platform").(string))
	if _, ok := d.GetOk("linux_params"); ok {
		linuxParamsModel, err := dataSourceIbmBackupRecoveryDownloadAgentMapToLinuxAgentParams(d.Get("linux_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_download_agent", "read", "parse-linux_params").GetDiag()
		}
		downloadAgentOptions.SetLinuxParams(linuxParamsModel)
	}

	typeString, _, err := backupRecoveryClient.DownloadAgentWithContext(context, downloadAgentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DownloadAgentWithContext failed: %s", err.Error()), "ibm_backup_recovery_download_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryAgentDownloadID(d))

	err = saveToFile(typeString, d.Get("file_path").(string))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_download_agent", "read", "parse-linux_params").GetDiag()
	}

	return nil
}

func dataSourceIbmBackupRecoveryDownloadAgentMapToLinuxAgentParams(modelMap map[string]interface{}) (*backuprecoveryv1.LinuxAgentParams, error) {
	model := &backuprecoveryv1.LinuxAgentParams{}
	model.PackageType = core.StringPtr(modelMap["package_type"].(string))
	return model, nil
}

func dataSourceIbmBackupRecoveryDownloadAgentLinuxAgentParamsToMap(model *backuprecoveryv1.LinuxAgentParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["package_type"] = *model.PackageType
	return modelMap, nil
}

func saveToFile(response io.ReadCloser, filePath string) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response)
	if err != nil {
		return err
	}

	err = response.Close()
	if err != nil {
		return err
	}

	return nil
}

// dataSourceIbmBackupRecoveryDownloadAgentID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryAgentDownloadID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
