package cos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	rc "github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCOSBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCOSBackupPolicyCreate,
		ReadContext:   resourceIBMCOSBackupPolicyRead,
		DeleteContext: resourceIBMCOSBackupPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket Crn of the source bucket.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the backup policy to be applied on the source bucket.",
			},
			"target_backup_vault_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CRN for a COS BackupVault.",
			},
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"continuous"}),
				Description:  "The type of backup to support.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the backup policy applied on the source bucket",
			},
		},
	}
}

func resourceIBMCOSBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	policyName := d.Get("policy_name").(string)
	targetBackupVaultCRN := d.Get("target_backup_vault_crn").(string)
	backupType := d.Get("backup_type").(string)
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	createBucketBackupPolicyOptions := &rc.CreateBackupPolicyOptions{
		Bucket:               aws.String(bucketName),
		PolicyName:           aws.String(policyName),
		TargetBackupVaultCrn: aws.String(targetBackupVaultCRN),
		BackupType:           aws.String(backupType),
	}
	res, _, err := rcClient.CreateBackupPolicy(createBucketBackupPolicyOptions)
	if err != nil {
		return diag.Errorf("Failed to create the backup policy %s, %v", bucketName, err)
	}
	policyId := *res.PolicyID
	policySetID := fmt.Sprintf("%s:%s:%s:target:%s", bucketName, policyId, policyName, targetBackupVaultCRN)
	d.SetId(policySetID)
	return resourceIBMCOSBackupPolicyRead(ctx, d, meta)
}

func resourceIBMCOSBackupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketName, err := parseBackupPolicyID(d.Id(), "bucketName")
	if err != nil {
		return diag.Errorf("Failed to get bucket name from backup policy Id : %v", err)
	}
	policyId, err := parseBackupPolicyID(d.Id(), "policyId")
	if err != nil {
		return diag.Errorf("Failed to get policy name from backup policy Id :  %v", err)
	}
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	d.Set("policy_id", policyId)
	getBucketBackupPolicyOptions := &rc.GetBackupPolicyOptions{
		Bucket:   aws.String(bucketName),
		PolicyID: aws.String(policyId),
	}

	res, response, err := rcClient.GetBackupPolicy(getBucketBackupPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error while reading the backup policy vault : %v", err)
	}
	if res != nil {
		if res.PolicyName != nil {
			d.Set("policy_name", aws.String(*res.PolicyName))
		}
		if res.TargetBackupVaultCrn != nil {
			d.Set("target_backup_vault_crn", aws.String(*res.TargetBackupVaultCrn))
		}
		if res.BackupType != nil {
			d.Set("backup_type", aws.String(*res.BackupType))
		}
	}
	return nil
}

func resourceIBMCOSBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketName, err := parseBackupPolicyID(d.Id(), "bucketName")
	if err != nil {
		return diag.Errorf("Failed to get bucket name from backup policy Id :  %v", err)
	}
	policyId, err := parseBackupPolicyID(d.Id(), "policyId")
	if err != nil {
		return diag.Errorf("Failed to get policy name from backup policy Id :  %v", err)
	}
	rcClient, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return diag.Errorf("Failed to create rc client %v", err)
	}
	deleteBucketBackupPolicyOptions := &rc.DeleteBackupPolicyOptions{
		Bucket:   aws.String(bucketName),
		PolicyID: aws.String(policyId),
	}
	_, deleteErr := rcClient.DeleteBackupPolicy(deleteBucketBackupPolicyOptions)
	if deleteErr != nil {
		return diag.Errorf("Error deleting the backup policy %v", deleteErr)
	}
	return nil
}

func parseBackupPolicyID(id string, info string) (backupPolicy string, err error) {
	if info == "bucketName" {
		return strings.Split(id, ":")[0], nil
	}
	if info == "policyId" {
		return strings.Split(id, ":")[1], nil
	}
	if info == "policyName" {
		return strings.Split(id, ":")[2], nil
	}
	if info == "targebucketCrn" {
		return strings.Split(id, ":target:")[1], nil
	}
	return "", errors.New("Backup policy ID is null")
}
