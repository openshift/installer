// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIbmIsShareReplicaOperations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareReplicaOperationsCreate,
		ReadContext:   resourceIbmIsShareReplicaOperationsRead,
		UpdateContext: resourceIbmIsShareReplicaOperationsUpdate,
		DeleteContext: resourceIbmIsShareReplicaOperationsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"share_replica": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file share identifier.",
			},
			"fallback_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "fail",
				ExactlyOneOf: []string{"split_share", "fallback_policy"},
				ValidateFunc: validate.InvokeValidator("ibm_is_share_replica_operations", "fallback_policy"),
				Description:  "The action to take if the failover request is accepted but cannot be performed or times out",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"fallback_policy"},
				//ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share_replica_operations", "timeout"),
				Description:  "The failover timeout in seconds",
			},
			"split_share": {
				Type:         schema.TypeBool,
				Default:      false,
				ForceNew:     true,
				Optional:     true,
				ExactlyOneOf: []string{"split_share", "fallback_policy"},
				Description:  "If set to true the replication relationship between source share and replica will be removed.",
			},
		},
	}
}

func ResourceIbmIsShareReplicaOperationsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	fbPolicyAllowedValues := "fail, split"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "timeout",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "300",
			MaxValue:                   "3600",
		},
		validate.ValidateSchema{
			Identifier:                 "fallback_policy",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              fbPolicyAllowedValues,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share_replica_operations", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareReplicaOperationsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	share_id := d.Get("share_replica").(string)

	splitShare := d.Get("split_share").(bool)

	if !splitShare {
		fallback_policy := d.Get("fallback_policy").(string)
		timeout := d.Get("timeout").(int)
		failOverShareOptions := &vpcv1.FailoverShareOptions{
			ShareID: &share_id,
		}
		failOverShareOptions.FallbackPolicy = &fallback_policy
		if timeout != 0 {
			failOverShareOptions.Timeout = core.Int64Ptr(int64(timeout))
		}
		response, err := vpcClient.FailoverShareWithContext(context, failOverShareOptions)
		if err != nil {
			log.Printf("[DEBUG] FailoverShareWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] FailoverShareWithContext failed %s\n%s", err, response))
		}
	} else {
		deleteShareSourceOptions := &vpcv1.DeleteShareSourceOptions{
			ShareID: &share_id,
		}
		response, err := vpcClient.DeleteShareSourceWithContext(context, deleteShareSourceOptions)
		if err != nil {
			log.Printf("[DEBUG] DeleteShareSourceWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] DeleteShareSourceWithContext failed %s\n%s", err, response))
		}
	}
	_, err = isWaitForShareReplicationJobDone(context, vpcClient, share_id, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(share_id)
	return nil
}

func isWaitForShareReplicationJobDone(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share (%s) to be available.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"active", "none"},
		Refresh:    isShareReplicationJobRefreshFunc(context, vpcClient, shareid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareReplicationJobRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareOptions := &vpcv1.GetShareOptions{}

		shareOptions.SetID(shareid)

		share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting share: %s\n%s", err, response)
		}
		if *share.ReplicationStatus == "active" || *share.ReplicationStatus == "none" {

			return share, *share.ReplicationStatus, nil

		}
		return share, "pending", nil
	}
}

func resourceIbmIsShareReplicaOperationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmIsShareReplicaOperationsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIbmIsShareReplicaOperationsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func isWaitForShareSplit(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share (%s) to be available.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"split_pending"},
		Target:     []string{"none"},
		Refresh:    isShareSplitRefreshFunc(context, vpcClient, shareid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareSplitRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareOptions := &vpcv1.GetShareOptions{}

		shareOptions.SetID(shareid)

		share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting share: %s\n%s", err, response)
		}
		d.Set("replication_status", *share.LifecycleState)
		if *share.LifecycleState == "none" {

			return share, *share.LifecycleState, nil

		}
		return share, "split_pending", nil
	}
}
