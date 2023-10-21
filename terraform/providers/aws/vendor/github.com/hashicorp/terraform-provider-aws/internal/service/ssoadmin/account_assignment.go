package ssoadmin

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_ssoadmin_account_assignment")
func ResourceAccountAssignment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAccountAssignmentCreate,
		ReadWithoutTimeout:   resourceAccountAssignmentRead,
		DeleteWithoutTimeout: resourceAccountAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},

			"permission_set_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},

			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 47),
					validation.StringMatch(regexp.MustCompile(`^([0-9a-f]{10}-|)[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$`), "must match ([0-9a-f]{10}-|)[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}"),
				),
			},

			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ssoadmin.PrincipalType_Values(), false),
			},

			"target_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},

			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(ssoadmin.TargetType_Values(), false),
			},
		},
	}
}

func resourceAccountAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	instanceArn := d.Get("instance_arn").(string)
	permissionSetArn := d.Get("permission_set_arn").(string)
	principalID := d.Get("principal_id").(string)
	principalType := d.Get("principal_type").(string)
	targetID := d.Get("target_id").(string)
	targetType := d.Get("target_type").(string)

	// We need to check if the assignment exists before creating it
	// since the AWS SSO API doesn't prevent us from creating duplicates
	accountAssignment, err := FindAccountAssignment(ctx, conn, principalID, principalType, targetID, permissionSetArn, instanceArn)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing SSO Account Assignments for AccountId (%s) PermissionSet (%s): %s", targetID, permissionSetArn, err)
	}

	if accountAssignment != nil {
		return sdkdiag.AppendErrorf(diags, "creating SSO Account Assignment for %s (%s): already exists", principalType, principalID)
	}

	input := &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
		PrincipalId:      aws.String(principalID),
		PrincipalType:    aws.String(principalType),
		TargetId:         aws.String(targetID),
		TargetType:       aws.String(targetType),
	}

	output, err := conn.CreateAccountAssignmentWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating SSO Account Assignment for %s (%s): %s", principalType, principalID, err)
	}

	if output == nil || output.AccountAssignmentCreationStatus == nil {
		return sdkdiag.AppendErrorf(diags, "creating SSO Account Assignment for %s (%s): empty output", principalType, principalID)
	}

	status := output.AccountAssignmentCreationStatus

	_, err = waitAccountAssignmentCreated(ctx, conn, instanceArn, aws.StringValue(status.RequestId))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for SSO Account Assignment for %s (%s) to be created: %s", principalType, principalID, err)
	}

	d.SetId(fmt.Sprintf("%s,%s,%s,%s,%s,%s", principalID, principalType, targetID, targetType, permissionSetArn, instanceArn))

	return append(diags, resourceAccountAssignmentRead(ctx, d, meta)...)
}

func resourceAccountAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	idParts, err := ParseAccountAssignmentID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing SSO Account Assignment ID: %s", err)
	}

	principalID := idParts[0]
	principalType := idParts[1]
	targetID := idParts[2]
	targetType := idParts[3]
	permissionSetArn := idParts[4]
	instanceArn := idParts[5]

	accountAssignment, err := FindAccountAssignment(ctx, conn, principalID, principalType, targetID, permissionSetArn, instanceArn)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, ssoadmin.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] SSO Account Assignment for Principal (%s) not found, removing from state", principalID)
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading SSO Account Assignment for Principal (%s): %s", principalID, err)
	}

	if accountAssignment == nil {
		if d.IsNewResource() {
			return sdkdiag.AppendErrorf(diags, "reading SSO Account Assignment for Principal (%s): not found", principalID)
		}

		log.Printf("[WARN] SSO Account Assignment for Principal (%s) not found, removing from state", principalID)
		d.SetId("")
		return diags
	}

	d.Set("instance_arn", instanceArn)
	d.Set("permission_set_arn", accountAssignment.PermissionSetArn)
	d.Set("principal_id", accountAssignment.PrincipalId)
	d.Set("principal_type", accountAssignment.PrincipalType)
	d.Set("target_id", accountAssignment.AccountId)
	d.Set("target_type", targetType)

	return diags
}

func resourceAccountAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	idParts, err := ParseAccountAssignmentID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing SSO Account Assignment ID: %s", err)
	}

	principalID := idParts[0]
	principalType := idParts[1]
	targetID := idParts[2]
	targetType := idParts[3]
	permissionSetArn := idParts[4]
	instanceArn := idParts[5]

	input := &ssoadmin.DeleteAccountAssignmentInput{
		PrincipalId:      aws.String(principalID),
		InstanceArn:      aws.String(instanceArn),
		PermissionSetArn: aws.String(permissionSetArn),
		TargetType:       aws.String(targetType),
		TargetId:         aws.String(targetID),
		PrincipalType:    aws.String(principalType),
	}

	output, err := conn.DeleteAccountAssignmentWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, ssoadmin.ErrCodeResourceNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting SSO Account Assignment for Principal (%s): %s", principalID, err)
	}

	if output == nil || output.AccountAssignmentDeletionStatus == nil {
		return sdkdiag.AppendErrorf(diags, "deleting SSO Account Assignment for Principal (%s): empty output", principalID)
	}

	status := output.AccountAssignmentDeletionStatus

	_, err = waitAccountAssignmentDeleted(ctx, conn, instanceArn, aws.StringValue(status.RequestId))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for SSO Account Assignment for Principal (%s) to be deleted: %s", principalID, err)
	}

	return diags
}

func ParseAccountAssignmentID(id string) ([]string, error) {
	idParts := strings.Split(id, ",")
	if len(idParts) != 6 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" ||
		idParts[3] == "" || idParts[4] == "" || idParts[5] == "" {
		return nil, fmt.Errorf("unexpected format for ID (%q), expected PRINCIPAL_ID,PRINCIPAL_TYPE,TARGET_ID,TARGET_TYPE,PERMISSION_SET_ARN,INSTANCE_ARN", id)
	}
	return idParts, nil
}
