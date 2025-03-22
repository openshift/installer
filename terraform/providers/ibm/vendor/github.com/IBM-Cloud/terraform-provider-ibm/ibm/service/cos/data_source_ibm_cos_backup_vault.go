// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	rc "github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCosBackupVault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCosBackupVaultRead,

		Schema: map[string]*schema.Schema{
			"backup_vault_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"activity_tracking_management_events": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"metrics_monitoring_usage_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"kms_key_crn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMCosBackupVaultRead(d *schema.ResourceData, meta interface{}) error {
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	backupVaultName := d.Get("backup_vault_name").(string)
	instanceCRN := d.Get("service_instance_id").(string)
	region := d.Get("region").(string)
	getBucketBackupVaultOptions := &rc.GetBackupVaultOptions{
		BackupVaultName: aws.String(backupVaultName),
	}

	d.Set("service_instance_id", instanceCRN)
	d.Set("region", region)
	res, _, err := rcClient.GetBackupVault(getBucketBackupVaultOptions)
	if err != nil {
		return err
	}

	vaultID := fmt.Sprintf("%s:%s:%s:meta:%s", strings.Replace(instanceCRN, "::", "", -1), "backup-vault", backupVaultName, region)
	d.SetId(vaultID)
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
			d.Set("kms_key_crn", aws.String(*res.SseKpCustomerRootKeyCrn))
		}

	}

	return nil
}
