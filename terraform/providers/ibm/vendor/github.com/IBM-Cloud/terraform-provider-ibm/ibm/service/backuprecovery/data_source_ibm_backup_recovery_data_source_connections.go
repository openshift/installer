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

func DataSourceIbmBackupRecoveryDataSourceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryDataSourceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"connection_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the unique IDs of the connections which are to be fetched.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connection_names": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the names of the connections which are to be fetched.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connections": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the unique ID of the connection.",
						},
						"connection_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the connection. For a given tenant, different connections can't have the same name. However, two (or more) different tenants can each have a connection with the same name.",
						},
						"connector_ids": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the IDs of the connectors in this connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"registration_token": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies a token that can be used to register a connector against this connection.",
						},
						"tenant_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the tenant ID of the connection.",
						},
						"upgrading_connector_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the connector ID that is currently in upgrade.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoveryDataSourceConnectionsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_data_source_connections", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDataSourceConnectionsOptions := &backuprecoveryv1.GetDataSourceConnectionsOptions{}

	getDataSourceConnectionsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("connection_ids"); ok {
		var connectionIds []string
		for _, v := range d.Get("connection_ids").([]interface{}) {
			connectionIdsItem := v.(string)
			connectionIds = append(connectionIds, connectionIdsItem)
		}
		getDataSourceConnectionsOptions.SetConnectionIds(connectionIds)
	}
	if _, ok := d.GetOk("connection_names"); ok {
		var connectionNames []string
		for _, v := range d.Get("connection_names").([]interface{}) {
			connectionNamesItem := v.(string)
			connectionNames = append(connectionNames, connectionNamesItem)
		}
		getDataSourceConnectionsOptions.SetConnectionNames(connectionNames)
	}

	dataSourceConnectionList, _, err := backupRecoveryClient.GetDataSourceConnectionsWithContext(context, getDataSourceConnectionsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataSourceConnectionsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_data_source_connections", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryDataSourceConnectionsID(d))

	if !core.IsNil(dataSourceConnectionList.Connections) {
		connections := []map[string]interface{}{}
		for _, connectionsItem := range dataSourceConnectionList.Connections {
			connectionsItemMap, err := DataSourceIbmBackupRecoveryDataSourceConnectionsDataSourceConnectionToMap(&connectionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_data_source_connections", "read", "connections-to-map").GetDiag()
			}
			connections = append(connections, connectionsItemMap)
		}
		if err = d.Set("connections", connections); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting connections: %s", err), "(Data) ibm_backup_recovery_data_source_connections", "read", "set-connections").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryDataSourceConnectionsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryDataSourceConnectionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryDataSourceConnectionsDataSourceConnectionToMap(model *backuprecoveryv1.DataSourceConnection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConnectionID != nil {
		modelMap["connection_id"] = *model.ConnectionID
	}
	modelMap["connection_name"] = *model.ConnectionName
	if model.ConnectorIds != nil {
		modelMap["connector_ids"] = model.ConnectorIds
	}
	if model.RegistrationToken != nil {
		modelMap["registration_token"] = *model.RegistrationToken
	}
	if model.TenantID != nil {
		modelMap["tenant_id"] = *model.TenantID
	}
	if model.UpgradingConnectorID != nil {
		modelMap["upgrading_connector_id"] = *model.UpgradingConnectorID
	}
	return modelMap, nil
}
