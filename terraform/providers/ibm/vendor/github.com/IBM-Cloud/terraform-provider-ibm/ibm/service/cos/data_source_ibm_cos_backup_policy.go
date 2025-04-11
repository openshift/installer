// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	rc "github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCosBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCosBackupPolicyRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the source bucket",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the backup policy applied on the source bucket",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the backup policy to be applied on the source bucket.",
			},
			"target_backup_vault_crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for a COS BackupVault.",
			},
			"backup_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of backup to support.",
			},
		},
	}
}

func dataSourceIBMCosBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	bucketName := d.Get("bucket_name").(string)
	policySetID := d.Get("policy_id").(string)
	d.Set("policy_id", policySetID)
	var policy_name, target_backup_vault_crn string
	getBucketBackupPolicyOptions := &rc.GetBackupPolicyOptions{
		Bucket:   aws.String(bucketName),
		PolicyID: aws.String(policySetID),
	}

	res, _, err := rcClient.GetBackupPolicy(getBucketBackupPolicyOptions)
	if err != nil {
		return err
	}
	if res != nil {
		if res.PolicyName != nil {
			d.Set("policy_name", aws.String(*res.PolicyName))
			policy_name = *res.PolicyName
		}
		if res.TargetBackupVaultCrn != nil {
			d.Set("target_backup_vault_crn", aws.String(*res.TargetBackupVaultCrn))
			target_backup_vault_crn = *res.TargetBackupVaultCrn
		}
		if res.BackupType != nil {
			d.Set("backup_type", aws.String(*res.BackupType))
		}
	}
	policyID := fmt.Sprintf("%s:%s:%s:target:%s", bucketName, policySetID, policy_name, target_backup_vault_crn)
	d.SetId(policyID)
	return nil
}
