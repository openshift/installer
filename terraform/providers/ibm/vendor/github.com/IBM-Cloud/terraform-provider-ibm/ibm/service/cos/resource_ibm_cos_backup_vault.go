package cos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	rc "github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCOSBackupVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCOSBackupVaultCreate,
		ReadContext:   resourceIBMCOSBackupVaultRead,
		UpdateContext: resourceIBMCOSBackupVaultUpdate,
		DeleteContext: resourceIBMCOSBackupVaultDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_vault_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Backup Vault.",
			},
			"service_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id for the backup vault.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location where backup vault  to be created.",
			},
			"backup_vault_crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"activity_tracking_management_events": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Activity Tracking configuration.Whether to send notifications for management events on the BackupVault.",
			},
			"metrics_monitoring_usage_metrics": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Metrics Monitoring configuration.Whether usage metrics are collected for this BackupVault.",
			},
			"kms_key_crn": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The CRN for a KeyProtect root key.",
			},
		},
	}
}

func resourceIBMCOSBackupVaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	backupVaultName := d.Get("backup_vault_name").(string)
	instanceCRN := d.Get("service_instance_id").(string)
	region := d.Get("region").(string)
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	createBackupVault := &rc.CreateBackupVaultOptions{}
	createBackupVault.BackupVaultName = aws.String(backupVaultName)
	createBackupVault.ServiceInstanceID = aws.String(instanceCRN)
	createBackupVault.Region = aws.String(region)

	if managementEvents, ok := d.GetOkExists("activity_tracking_management_events"); ok {
		createBackupVault.ActivityTracking = &rc.BackupVaultActivityTracking{}
		managementEvents := managementEvents.(bool)
		createBackupVault.ActivityTracking.ManagementEvents = &managementEvents
	}
	if usageMetrics, ok := d.GetOkExists("metrics_monitoring_usage_metrics"); ok {
		um := usageMetrics.(bool)
		createBackupVault.MetricsMonitoring = &rc.BackupVaultMetricsMonitoring{}
		createBackupVault.MetricsMonitoring.UsageMetricsEnabled = &um
	}
	if key, ok := d.GetOk("kms_key_crn"); ok {
		createBackupVault.SseKpCustomerRootKeyCrn = aws.String(key.(string))
	}
	_, _, err = rcClient.CreateBackupVault(createBackupVault)
	if err != nil {
		return diag.Errorf("Failed to create the backup vault %s, %v", backupVaultName, err)
	}
	vaultID := fmt.Sprintf("%s:%s:%s:meta:%s", strings.Replace(instanceCRN, "::", "", -1), "backup-vault", backupVaultName, region)
	d.SetId(vaultID)
	return resourceIBMCOSBackupVaultUpdate(ctx, d, meta)
}

func resourceIBMCOSBackupVaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupVaultName := d.Get("backup_vault_name").(string)
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	backupVaultPatch := &rc.BackupVaultPatch{}
	hasChange := false
	if d.HasChange("activity_tracking_management_events") {
		hasChange = true
		_, newValue := d.GetChange("activity_tracking_management_events")
		if newValue != nil {
			backupVaultPatch.ActivityTracking = &rc.BackupVaultActivityTracking{}
			backupVaultPatch.ActivityTracking.ManagementEvents = aws.Bool(d.Get("activity_tracking_management_events").(bool))
		} else {
			backupVaultPatch.ActivityTracking = &rc.BackupVaultActivityTracking{}
		}
	}

	if d.HasChange("metrics_monitoring_usage_metrics") {
		hasChange = true
		_, newValue := d.GetChange("metrics_monitoring_usage_metrics")
		if newValue != nil {
			backupVaultPatch.MetricsMonitoring = &rc.BackupVaultMetricsMonitoring{}
			backupVaultPatch.MetricsMonitoring.UsageMetricsEnabled = aws.Bool(d.Get("metrics_monitoring_usage_metrics").(bool))
		} else {
			backupVaultPatch.MetricsMonitoring = &rc.BackupVaultMetricsMonitoring{}
		}
	}

	if hasChange == true {
		bucketPatchModelAsPatch, asPatchErr := backupVaultPatch.AsPatch()
		if asPatchErr != nil {
			return diag.Errorf("Unable to create the update patch for backup vault %v", asPatchErr)
		}
		updateBucketBackupVaultOptions := &rc.UpdateBackupVaultOptions{
			BackupVaultName:  aws.String(backupVaultName),
			BackupVaultPatch: bucketPatchModelAsPatch,
		}
		_, _, err = rcClient.UpdateBackupVault(updateBucketBackupVaultOptions)
		if err != nil {
			return diag.Errorf("Unable to update the backup vault configurations %v", err)
		}
	}
	return resourceIBMCOSBackupVaultRead(ctx, d, meta)
}

func resourceIBMCOSBackupVaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupVaultName := parseBackupVaultID(d.Id(), "backupVaultName")
	instanceCRN := parseBackupVaultID(d.Id(), "instanceCRN")
	region := parseBackupVaultID(d.Id(), "region")
	crn := parseBackupVaultID(d.Id(), "backupVaultCrn")
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	d.Set("service_instance_id", instanceCRN)
	d.Set("region", region)
	d.Set("backup_vault_crn", crn)
	getBucketBackupVaultOptions := &rc.GetBackupVaultOptions{
		BackupVaultName: aws.String(backupVaultName),
	}
	res, _, err := rcClient.GetBackupVault(getBucketBackupVaultOptions)
	if err != nil {
		return diag.Errorf("Error while reading the backup vault : %v", err)
	}
	if res != nil {

		if res.ActivityTracking != nil {
			if res.ActivityTracking.ManagementEvents != nil {
				d.Set("activity_tracking_management_events", *res.ActivityTracking.ManagementEvents)
			}
		} else {
			d.Set("activity_tracking_management_events", nil)
		}
		if res.MetricsMonitoring != nil {
			d.Set("metrics_monitoring_usage_metrics", *res.MetricsMonitoring.UsageMetricsEnabled)
		} else {
			d.Set("metrics_monitoring_usage_metrics", nil)
		}
		if res.SseKpCustomerRootKeyCrn != nil {
			d.Set("kms_key_crn", *res.SseKpCustomerRootKeyCrn)
		}
	}

	return nil
}

func resourceIBMCOSBackupVaultDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	backupVaultName := parseBackupVaultID(d.Id(), "backupVaultName")
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	deleteBucketBackupVaultOptions := &rc.DeleteBackupVaultOptions{
		BackupVaultName: aws.String(backupVaultName),
	}
	_, deleteErr := rcClient.DeleteBackupVault(deleteBucketBackupVaultOptions)
	if deleteErr != nil {
		return diag.Errorf("Error deleting the backup vault %v", deleteErr)
	}
	return nil
}

func parseBackupVaultID(id string, info string) string {
	backupVaultCrn := strings.Split(id, ":meta:")[0]

	if info == "backupVaultName" {
		return strings.Split(backupVaultCrn, ":backup-vault:")[1]
	}
	if info == "instanceCRN" {
		return fmt.Sprintf("%s::", strings.Split(backupVaultCrn, ":backup-vault:")[0])
	}
	if info == "backupVaultCrn" {
		return backupVaultCrn
	}
	if info == "region" {
		return strings.Split(id, ":meta:")[1]
	}
	return ""
}
