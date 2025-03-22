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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func DataSourceIbmBackupRecoveryConnectorsMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryConnectorsMetadataRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"connector_image_metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies information about the connector images for various platforms.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_image_file_list": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies info about connector images for the supported platforms.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the platform on which the image can be deployed.",
									},
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the URL to access the file.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoveryConnectorsMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_connectors_metadata", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getConnectorMetadataOptions := &backuprecoveryv1.GetConnectorMetadataOptions{}

	getConnectorMetadataOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	connectorMetadata, _, err := backupRecoveryClient.GetConnectorMetadataWithContext(context, getConnectorMetadataOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConnectorMetadataWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_connectors_metadata", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryConnectorsMetadataID(d))

	if !core.IsNil(connectorMetadata.ConnectorImageMetadata) {
		connectorImageMetadata := []map[string]interface{}{}
		connectorImageMetadataMap, err := DataSourceIbmBackupRecoveryConnectorsMetadataConnectorImageMetadataToMap(connectorMetadata.ConnectorImageMetadata)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_connectors_metadata", "read", "connector_image_metadata-to-map").GetDiag()
		}
		connectorImageMetadata = append(connectorImageMetadata, connectorImageMetadataMap)
		if err = d.Set("connector_image_metadata", connectorImageMetadata); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting connector_image_metadata: %s", err), "(Data) ibm_backup_recovery_connectors_metadata", "read", "set-connector_image_metadata").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryConnectorsMetadataID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryConnectorsMetadataID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryConnectorsMetadataConnectorImageMetadataToMap(model *backuprecoveryv1.ConnectorImageMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	connectorImageFileList := []map[string]interface{}{}
	for _, connectorImageFileListItem := range model.ConnectorImageFileList {
		connectorImageFileListItemMap, err := DataSourceIbmBackupRecoveryConnectorsMetadataConnectorImageFileToMap(&connectorImageFileListItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		connectorImageFileList = append(connectorImageFileList, connectorImageFileListItemMap)
	}
	modelMap["connector_image_file_list"] = connectorImageFileList
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryConnectorsMetadataConnectorImageFileToMap(model *backuprecoveryv1.ConnectorImageFile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["image_type"] = *model.ImageType
	modelMap["url"] = *model.URL
	return modelMap, nil
}
