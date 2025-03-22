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

func DataSourceIbmBackupRecoveryDataSourceConnectors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryDataSourceConnectorsRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"connector_ids": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"connection_id"},
				Description:   "Specifies the unique IDs of the connectors which are to be fetched.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connector_names": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the names of the connectors which are to be fetched.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connection_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"connector_ids"},
				Description:   "Specifies the ID of the connection, connectors belonging to which are to be fetched.",
			},
			"connectors": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_side_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the IP of the connector's NIC facing the cluster.",
						},
						"connection_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the ID of the connection to which this connector belongs.",
						},
						"connector_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the unique ID of the connector.",
						},
						"connector_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the connector. The name of a connector need not be unique within a tenant or across tenants. The name of the connector can be updated as needed.",
						},
						"connectivity_status": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies status information for the data-source connector. For example if it's currently connected to the cluster, when it last connected to the cluster successfully, etc.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_connected": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether the connector is currently connected to the cluster.",
									},
									"last_connected_timestamp_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the last timestamp in UNIX time (seconds) when the connector had successfully connected to the cluster. This property can be present even if the connector is currently disconnected.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies error message when the connector is unable to connect to the cluster.",
									},
								},
							},
						},
						"software_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the connector's software version.",
						},
						"tenant_side_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the IP of the connector's NIC facing the sources of the tenant to which the connector belongs.",
						},
						"upgrade_status": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies upgrade status for the data-source connector. For example when the upgrade started, current status of the upgrade, errors for upgrade failure etc.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_status_fetched_timestamp_msecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the last timestamp in UNIX time (milliseconds) when the connector upgrade status was fetched.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies error message for upgrade failure.",
									},
									"start_timestamp_m_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the last timestamp in UNIX time (milliseconds) when the connector upgrade was triggered.",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the last fetched upgrade status of the connector.",
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

func dataSourceIbmBackupRecoveryDataSourceConnectorsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_data_source_connectors", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDataSourceConnectorsOptions := &backuprecoveryv1.GetDataSourceConnectorsOptions{}

	getDataSourceConnectorsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("connector_ids"); ok {
		var connectorIds []string
		for _, v := range d.Get("connector_ids").([]interface{}) {
			connectorIdsItem := v.(string)
			connectorIds = append(connectorIds, connectorIdsItem)
		}
		getDataSourceConnectorsOptions.SetConnectorIds(connectorIds)
	}
	if _, ok := d.GetOk("connector_names"); ok {
		var connectorNames []string
		for _, v := range d.Get("connector_names").([]interface{}) {
			connectorNamesItem := v.(string)
			connectorNames = append(connectorNames, connectorNamesItem)
		}
		getDataSourceConnectorsOptions.SetConnectorNames(connectorNames)
	}
	if _, ok := d.GetOk("connection_id"); ok {
		getDataSourceConnectorsOptions.SetConnectionID(d.Get("connection_id").(string))
	}

	dataSourceConnectorList, _, err := backupRecoveryClient.GetDataSourceConnectorsWithContext(context, getDataSourceConnectorsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataSourceConnectorsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_data_source_connectors", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryDataSourceConnectorsID(d))

	if !core.IsNil(dataSourceConnectorList) {
		if !core.IsNil(dataSourceConnectorList.Connectors) {
			connectors := []map[string]interface{}{}
			for _, connectorsItem := range dataSourceConnectorList.Connectors {
				connectorsItemMap, err := DataSourceIbmBackupRecoveryDataSourceConnectorsDataSourceConnectorToMap(&connectorsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_data_source_connectors", "read", "connectors-to-map").GetDiag()
				}
				connectors = append(connectors, connectorsItemMap)
			}
			if err = d.Set("connectors", connectors); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting connectors: %s", err), "(Data) ibm_backup_recovery_data_source_connectors", "read", "set-connectors").GetDiag()
			}
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryDataSourceConnectorsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryDataSourceConnectorsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryDataSourceConnectorsDataSourceConnectorToMap(model *backuprecoveryv1.DataSourceConnector) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSideIp != nil {
		modelMap["cluster_side_ip"] = *model.ClusterSideIp
	}
	if model.ConnectionID != nil {
		modelMap["connection_id"] = *model.ConnectionID
	}
	if model.ConnectorID != nil {
		modelMap["connection_id"] = *model.ConnectorID
	}
	if model.ConnectorName != nil {
		modelMap["connector_name"] = *model.ConnectorName
	}
	if model.ConnectivityStatus != nil {
		connectivityStatusMap, err := DataSourceIbmBackupRecoveryDataSourceConnectorsConnectorConnectivityStatusToMap(model.ConnectivityStatus)
		if err != nil {
			return modelMap, err
		}
		modelMap["connectivity_status"] = []map[string]interface{}{connectivityStatusMap}
	}
	if model.SoftwareVersion != nil {
		modelMap["software_version"] = *model.SoftwareVersion
	}
	if model.TenantSideIp != nil {
		modelMap["tenant_side_ip"] = *model.TenantSideIp
	}
	if model.UpgradeStatus != nil {
		upgradeStatusMap, err := DataSourceIbmBackupRecoveryDataSourceConnectorsConnectorUpgradeStatusToMap(model.UpgradeStatus)
		if err != nil {
			return modelMap, err
		}
		modelMap["upgrade_status"] = []map[string]interface{}{upgradeStatusMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryDataSourceConnectorsConnectorConnectivityStatusToMap(model *backuprecoveryv1.DataSourceConnectorConnectivityStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["is_connected"] = *model.IsConnected
	if model.LastConnectedTimestampSecs != nil {
		modelMap["last_connected_timestamp_secs"] = flex.IntValue(model.LastConnectedTimestampSecs)
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryDataSourceConnectorsConnectorUpgradeStatusToMap(model *backuprecoveryv1.DataSourceConnectorUpgradeStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LastStatusFetchedTimestampMsecs != nil {
		modelMap["last_status_fetched_timestamp_msecs"] = flex.IntValue(model.LastStatusFetchedTimestampMsecs)
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.StartTimestampMSecs != nil {
		modelMap["start_timestamp_m_secs"] = flex.IntValue(model.StartTimestampMSecs)
	}
	modelMap["status"] = *model.Status
	return modelMap, nil
}
