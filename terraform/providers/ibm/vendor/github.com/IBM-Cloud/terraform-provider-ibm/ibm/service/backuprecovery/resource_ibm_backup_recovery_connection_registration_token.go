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

func ResourceIbmBackupRecoveryConnectionRegistrationToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryConnectionRegistrationTokenCreate,
		ReadContext:   resourceIbmBackupRecoveryConnectionRegistrationTokenRead,
		DeleteContext: resourceIbmBackupRecoveryConnectionRegistrationTokenDelete,
		UpdateContext: resourceIbmBackupRecoveryConnectionRegistrationTokenUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryConnectionRegistrationToken,
		Schema: map[string]*schema.Schema{
			"connection_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the ID of the connection, connectors belonging to which are to be fetched.",
			},
			"x_ibm_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"registration_token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func checkDiffResourceIbmBackupRecoveryConnectionRegistrationToken(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryConnectionRegistrationToken().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_connection_registration_token cannot be updated.")
		}
	}
	return nil
}

func resourceIbmBackupRecoveryConnectionRegistrationTokenCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_connection_registration_token", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	generateDataSourceConnectionRegistrationTokenOptions := &backuprecoveryv1.GenerateDataSourceConnectionRegistrationTokenOptions{}

	generateDataSourceConnectionRegistrationTokenOptions.SetConnectionID(d.Get("connection_id").(string))
	generateDataSourceConnectionRegistrationTokenOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	connectionRegistrationTokenString, _, err := backupRecoveryClient.GenerateDataSourceConnectionRegistrationTokenWithContext(context, generateDataSourceConnectionRegistrationTokenOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GenerateDataSourceConnectionRegistrationTokenWithContext failed: %s", err.Error()), "ibm_backup_recovery_connection_registration_token", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(resourceIbmBackupRecoveryConnectionRegistrationTokenID(d))

	if !core.IsNil(connectionRegistrationTokenString) {
		if err = d.Set("registration_token", connectionRegistrationTokenString); err != nil {
			err = fmt.Errorf("Error setting registration_token: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_connection-registration-token", "read", "set-registration_token").GetDiag()
		}
	}

	return resourceIbmBackupRecoveryConnectionRegistrationTokenRead(context, d, meta)
}

func resourceIbmBackupRecoveryConnectionRegistrationTokenRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
func resourceIbmBackupRecoveryConnectionRegistrationTokenID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func resourceIbmBackupRecoveryConnectionRegistrationTokenDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryConnectionRegistrationTokenUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource update will only affect terraform state and not the actual backend resource",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource.",
	}
	diags = append(diags, warning)
	// d.SetId("")
	return diags
}
