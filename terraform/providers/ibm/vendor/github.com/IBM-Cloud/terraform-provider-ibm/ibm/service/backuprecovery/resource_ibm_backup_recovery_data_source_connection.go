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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryDataSourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryDataSourceConnectionCreate,
		ReadContext:   resourceIbmBackupRecoveryDataSourceConnectionRead,
		UpdateContext: resourceIbmBackupRecoveryDataSourceConnectionUpdate,
		DeleteContext: resourceIbmBackupRecoveryDataSourceConnectionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tenant accessing the cluster.",
			},
			"connection_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "connection Id",
			},
			"connection_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the connection. For a given tenant, different connections can't have the same name. However, two (or more) different tenants can each have a connection with the same name.",
			},
			"connector_ids": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the IDs of the connectors in this connection.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	}
}

func resourceIbmBackupRecoveryDataSourceConnectionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createDataSourceConnectionOptions := &backuprecoveryv1.CreateDataSourceConnectionOptions{}

	createDataSourceConnectionOptions.SetConnectionName(d.Get("connection_name").(string))
	if _, ok := d.GetOk("x_ibm_tenant_id"); ok {
		createDataSourceConnectionOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	}

	dataSourceConnection, _, err := backupRecoveryClient.CreateDataSourceConnectionWithContext(context, createDataSourceConnectionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDataSourceConnectionWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connection", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	connectionId := fmt.Sprintf("%s::%s", tenantId, *dataSourceConnection.ConnectionID)
	d.SetId(connectionId)
	d.Set("registration_token", dataSourceConnection.RegistrationToken)

	return resourceIbmBackupRecoveryDataSourceConnectionRead(context, d, meta)
}

func resourceIbmBackupRecoveryDataSourceConnectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tenantId := d.Get("x_ibm_tenant_id").(string)
	connectionId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		connectionId = ParseId(d.Id(), "id")
	}

	getDataSourceConnectionsOptions := &backuprecoveryv1.GetDataSourceConnectionsOptions{}
	getDataSourceConnectionsOptions.ConnectionIds = []string{connectionId}
	getDataSourceConnectionsOptions.SetXIBMTenantID(tenantId)

	dataSourceConnectionList, response, err := backupRecoveryClient.GetDataSourceConnectionsWithContext(context, getDataSourceConnectionsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataSourceConnectionsWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connection", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("connection_id", dataSourceConnectionList.Connections[0].ConnectionID); err != nil {
		err = fmt.Errorf("Error setting connection_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-connection_id").GetDiag()
	}

	if err = d.Set("connection_name", dataSourceConnectionList.Connections[0].ConnectionName); err != nil {
		err = fmt.Errorf("Error setting connection_name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-connection_name").GetDiag()
	}
	if !core.IsNil(dataSourceConnectionList.Connections[0].ConnectorIds) {
		if err = d.Set("connector_ids", dataSourceConnectionList.Connections[0].ConnectorIds); err != nil {
			err = fmt.Errorf("Error setting connector_ids: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-connector_ids").GetDiag()
		}
	} else {
		if err = d.Set("connector_ids", []string{}); err != nil {
			err = fmt.Errorf("Error setting connector_ids: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-connector_ids").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectionList.Connections[0].RegistrationToken) {
		if err = d.Set("registration_token", dataSourceConnectionList.Connections[0].RegistrationToken); err != nil {
			err = fmt.Errorf("Error setting registration_token: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-registration_token").GetDiag()
		}
	}

	if !core.IsNil(dataSourceConnectionList.Connections[0].TenantID) {
		if err = d.Set("x_ibm_tenant_id", dataSourceConnectionList.Connections[0].TenantID); err != nil {
			err = fmt.Errorf("Error setting x_ibm_tenant_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-x_ibm_tenant_id").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectionList.Connections[0].TenantID) {
		if err = d.Set("tenant_id", dataSourceConnectionList.Connections[0].TenantID); err != nil {
			err = fmt.Errorf("Error setting tenant_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-tenant_id").GetDiag()
		}
	}
	if !core.IsNil(dataSourceConnectionList.Connections[0].UpgradingConnectorID) {
		if err = d.Set("upgrading_connector_id", dataSourceConnectionList.Connections[0].UpgradingConnectorID); err != nil {
			err = fmt.Errorf("Error setting upgrading_connector_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "read", "set-upgrading_connector_id").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryDataSourceConnectionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	patchDataSourceConnectionOptions := &backuprecoveryv1.PatchDataSourceConnectionOptions{}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	connectionId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		connectionId = ParseId(d.Id(), "id")
	}
	patchDataSourceConnectionOptions.SetConnectionID(connectionId)
	patchDataSourceConnectionOptions.SetXIBMTenantID(tenantId)

	hasChange := false

	if d.HasChange("connection_name") {
		patchDataSourceConnectionOptions.SetConnectionName(d.Get("connection_name").(string))
		hasChange = true
	}

	if hasChange {
		_, _, err = backupRecoveryClient.PatchDataSourceConnectionWithContext(context, patchDataSourceConnectionOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("PatchDataSourceConnectionWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connection", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmBackupRecoveryDataSourceConnectionRead(context, d, meta)
}

func resourceIbmBackupRecoveryDataSourceConnectionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_data_source_connection", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteDataSourceConnectionOptions := &backuprecoveryv1.DeleteDataSourceConnectionOptions{}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	connectionId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		connectionId = ParseId(d.Id(), "id")
	}
	deleteDataSourceConnectionOptions.SetConnectionID(connectionId)
	deleteDataSourceConnectionOptions.SetXIBMTenantID(tenantId)

	_, err = backupRecoveryClient.DeleteDataSourceConnectionWithContext(context, deleteDataSourceConnectionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDataSourceConnectionWithContext failed: %s", err.Error()), "ibm_backup_recovery_data_source_connection", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")
	return nil
}
