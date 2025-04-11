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

func ResourceIbmBackupRecoveryDataSourceConnectorPatch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryDataSourceConnectorPatchCreate,
		ReadContext:   resourceIbmBackupRecoveryDataSourceConnectorPatchRead,
		DeleteContext: resourceIbmBackupRecoveryDataSourceConnectorPatchDelete,
		UpdateContext: resourceIbmBackupRecoveryDataSourceConnectorPatchUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryDataSourceConnectorPatch,
		Schema: map[string]*schema.Schema{
			"connector_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the unique ID of the connector which is to be deleted.",
			},
			"x_ibm_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"connector_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the name of the connector. The name of a connector need not be unique within a tenant or across tenants. The name of the connector can be updated as needed.",
			},
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
							Optional:    true,
							Computed:    true,
							Description: "Specifies the last timestamp in UNIX time (seconds) when the connector had successfully connected to the cluster. This property can be present even if the connector is currently disconnected.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
	}
}

func checkDiffResourceIbmBackupRecoveryDataSourceConnectorPatch(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryDataSourceConnectorPatch().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_data_source_connector_patch cannot be updated.")
		}
	}
	return nil
}

func resourceIbmBackupRecoveryDataSourceConnectorPatchCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	patchDataSourceConnectorOptions := &backuprecoveryv1.PatchDataSourceConnectorOptions{}

	patchDataSourceConnectorOptions.SetConnectorID(d.Get("connector_id").(string))
	patchDataSourceConnectorOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("connector_name"); ok {
		patchDataSourceConnectorOptions.SetConnectorName(d.Get("connector_name").(string))
	}

	dataSourceConnector, _, err := backupRecoveryClient.PatchDataSourceConnectorWithContext(context, patchDataSourceConnectorOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("PatchDataSourceConnectorWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connector_patch", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*dataSourceConnector.ConnectorID)

	return resourceIbmBackupRecoveryDataSourceConnectorPatchRead(context, d, meta)
}

func resourceIbmBackupRecoveryDataSourceConnectorPatchRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDataSourceConnectorsOptions := &backuprecoveryv1.GetDataSourceConnectorsOptions{}

	getDataSourceConnectorsOptions.SetConnectorIds([]string{d.Id()})
	getDataSourceConnectorsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	dataSourceConnectorList, response, err := backupRecoveryClient.GetDataSourceConnectorsWithContext(context, getDataSourceConnectorsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataSourceConnectorsWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connector_patch", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(dataSourceConnectorList.Connectors[0].ConnectorName) {
		if err = d.Set("connector_name", dataSourceConnectorList.Connectors[0].ConnectorName); err != nil {
			err = fmt.Errorf("Error setting connector_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-connector_name").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].ClusterSideIp) {
		if err = d.Set("cluster_side_ip", dataSourceConnectorList.Connectors[0].ClusterSideIp); err != nil {
			err = fmt.Errorf("Error setting cluster_side_ip: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-cluster_side_ip").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].ConnectionID) {
		if err = d.Set("connection_id", dataSourceConnectorList.Connectors[0].ConnectionID); err != nil {
			err = fmt.Errorf("Error setting connection_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-connection_id").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].ConnectivityStatus) {
		connectivityStatusMap, err := ResourceIbmBackupRecoveryDataSourceConnectorPatchDataSourceConnectivityStatusToMap(dataSourceConnectorList.Connectors[0].ConnectivityStatus)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "connectivity_status-to-map").GetDiag()
		}
		if err = d.Set("connectivity_status", []map[string]interface{}{connectivityStatusMap}); err != nil {
			err = fmt.Errorf("Error setting connectivity_status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-connectivity_status").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].SoftwareVersion) {
		if err = d.Set("software_version", dataSourceConnectorList.Connectors[0].SoftwareVersion); err != nil {
			err = fmt.Errorf("Error setting software_version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-software_version").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].TenantSideIp) {
		if err = d.Set("tenant_side_ip", dataSourceConnectorList.Connectors[0].TenantSideIp); err != nil {
			err = fmt.Errorf("Error setting tenant_side_ip: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-tenant_side_ip").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectorList.Connectors[0].UpgradeStatus) {
		if err = d.Set("upgrade_status", dataSourceConnectorList.Connectors[0].UpgradeStatus); err != nil {
			err = fmt.Errorf("Error setting upgrade_status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connector_patch", "read", "set-upgrade_status").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryDataSourceConnectorID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func resourceIbmBackupRecoveryDataSourceConnectorPatchDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.

	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Delete Not Supported",
		Detail:   "The resource definition will be only be removed from the terraform statefile. This resource cannot be deleted from the backend. ",
	}
	diags = append(diags, warning)
	d.SetId("")
	return diags
}

func resourceIbmBackupRecoveryDataSourceConnectorPatchUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "update" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource update will only affect terraform state and not the actual backend resource",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource.",
	}
	// d.SetId("")
	diags = append(diags, warning)
	return diags
}

func ResourceIbmBackupRecoveryDataSourceConnectorPatchDataSourceConnectivityStatusToMap(model *backuprecoveryv1.DataSourceConnectorConnectivityStatus) (map[string]interface{}, error) {
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
