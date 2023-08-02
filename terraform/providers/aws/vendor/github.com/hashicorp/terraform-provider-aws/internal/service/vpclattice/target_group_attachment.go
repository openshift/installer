package vpclattice

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_vpclattice_target_group_attachment", name="Target Group Attachment")
func resourceTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTargetGroupAttachmentCreate,
		ReadWithoutTimeout:   resourceTargetGroupAttachmentRead,
		DeleteWithoutTimeout: resourceTargetGroupAttachmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 2048),
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsPortNumber,
						},
					},
				},
			},
			"target_group_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTargetGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	targetGroupID := d.Get("target_group_identifier").(string)
	target := expandTarget(d.Get("target").([]interface{})[0].(map[string]interface{}))
	targetID := aws.ToString(target.Id)
	targetPort := int(aws.ToInt32(target.Port))
	id := strings.Join([]string{targetGroupID, targetID, strconv.Itoa(targetPort)}, "/")
	input := &vpclattice.RegisterTargetsInput{
		TargetGroupIdentifier: aws.String(targetGroupID),
		Targets:               []types.Target{target},
	}

	_, err := conn.RegisterTargets(ctx, input)

	if err != nil {
		return diag.Errorf("creating VPC Lattice Target Group Attachment (%s): %s", id, err)
	}

	d.SetId(id)

	if _, err := waitTargetGroupAttachmentCreated(ctx, conn, targetGroupID, targetID, targetPort, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for VPC Lattice Target Group Attachment (%s) create: %s", id, err)
	}

	return resourceTargetGroupAttachmentRead(ctx, d, meta)
}

func resourceTargetGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	targetGroupID := d.Get("target_group_identifier").(string)
	target := expandTarget(d.Get("target").([]interface{})[0].(map[string]interface{}))
	targetID := aws.ToString(target.Id)
	targetPort := int(aws.ToInt32(target.Port))

	output, err := findTargetByThreePartKey(ctx, conn, targetGroupID, targetID, targetPort)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] VPC Lattice Target Group Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading VPC Lattice Target Group Attachment (%s): %s", d.Id(), err)
	}

	if err := d.Set("target", []interface{}{flattenTargetSummary(output)}); err != nil {
		return diag.Errorf("setting target: %s", err)
	}
	d.Set("target_group_identifier", targetGroupID)

	return nil
}

func resourceTargetGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	targetGroupID := d.Get("target_group_identifier").(string)
	target := expandTarget(d.Get("target").([]interface{})[0].(map[string]interface{}))
	targetID := aws.ToString(target.Id)
	targetPort := int(aws.ToInt32(target.Port))

	log.Printf("[INFO] Deleting VPC Lattice Target Group Attachment: %s", d.Id())
	_, err := conn.DeregisterTargets(ctx, &vpclattice.DeregisterTargetsInput{
		TargetGroupIdentifier: aws.String(targetGroupID),
		Targets:               []types.Target{target},
	})

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting VPC Lattice Target Group Attachment (%s): %s", d.Id(), err)
	}

	if _, err := waitTargetGroupAttachmentDeleted(ctx, conn, targetGroupID, targetID, targetPort, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for VPC Lattice Target Group Attachment (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func findTargetByThreePartKey(ctx context.Context, conn *vpclattice.Client, targetGroupID, targetID string, targetPort int) (*types.TargetSummary, error) {
	input := &vpclattice.ListTargetsInput{
		TargetGroupIdentifier: aws.String(targetGroupID),
		Targets: []types.Target{{
			Id: aws.String(targetID),
		}},
	}
	if targetPort > 0 {
		input.Targets[0].Port = aws.Int32(int32(targetPort))
	}

	paginator := vpclattice.NewListTargetsPaginator(conn, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)

		if errs.IsA[*types.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: input,
			}
		}

		if err != nil {
			return nil, err
		}

		if output != nil && len(output.Items) == 1 {
			return &(output.Items[0]), nil
		}
	}

	return nil, &retry.NotFoundError{}
}

func statusTarget(ctx context.Context, conn *vpclattice.Client, targetGroupID, targetID string, targetPort int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findTargetByThreePartKey(ctx, conn, targetGroupID, targetID, targetPort)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, string(output.Status), nil
	}
}

func waitTargetGroupAttachmentCreated(ctx context.Context, conn *vpclattice.Client, targetGroupID, targetID string, targetPort int, timeout time.Duration) (*types.TargetSummary, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   enum.Slice(types.TargetStatusInitial),
		Target:                    enum.Slice(types.TargetStatusHealthy, types.TargetStatusUnhealthy, types.TargetStatusUnused, types.TargetStatusUnavailable),
		Refresh:                   statusTarget(ctx, conn, targetGroupID, targetID, targetPort),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*types.TargetSummary); ok {
		tfresource.SetLastError(err, errors.New(aws.ToString(output.ReasonCode)))

		return output, err
	}

	return nil, err
}

func waitTargetGroupAttachmentDeleted(ctx context.Context, conn *vpclattice.Client, targetGroupID, targetID string, targetPort int, timeout time.Duration) (*types.TargetSummary, error) {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.TargetStatusDraining, types.TargetStatusInitial),
		Target:  []string{},
		Refresh: statusTarget(ctx, conn, targetGroupID, targetID, targetPort),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*types.TargetSummary); ok {
		tfresource.SetLastError(err, errors.New(aws.ToString(output.ReasonCode)))

		return output, err
	}

	return nil, err
}

func flattenTargetSummary(apiObject *types.TargetSummary) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Id; v != nil {
		tfMap["id"] = aws.ToString(v)
	}

	if v := apiObject.Port; v != nil {
		tfMap["port"] = aws.ToInt32(v)
	}

	return tfMap
}

func expandTarget(tfMap map[string]interface{}) types.Target {
	apiObject := types.Target{}

	if v, ok := tfMap["id"].(string); ok && v != "" {
		apiObject.Id = aws.String(v)
	}

	if v, ok := tfMap["port"].(int); ok && v != 0 {
		apiObject.Port = aws.Int32(int32(v))
	}

	return apiObject
}
