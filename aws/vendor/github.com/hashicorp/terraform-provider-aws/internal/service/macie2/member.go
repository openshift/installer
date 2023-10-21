package macie2

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_macie2_member", name="Member")
// @Tags
func ResourceMember() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMemberCreate,
		ReadWithoutTimeout:   resourceMemberRead,
		UpdateWithoutTimeout: resourceMemberUpdate,
		DeleteWithoutTimeout: resourceMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrTags:    tftags.TagsSchemaForceNew(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"relationship_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"administrator_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"invited_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(macie2.MacieStatus_Values(), false),
			},
			"invite": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"invitation_disable_email_notification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"invitation_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Second),
			Update: schema.DefaultTimeout(60 * time.Second),
		},
	}
}

func resourceMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn(ctx)

	accountId := d.Get("account_id").(string)
	input := &macie2.CreateMemberInput{
		Account: &macie2.AccountDetail{
			AccountId: aws.String(accountId),
			Email:     aws.String(d.Get("email").(string)),
		},
		Tags: GetTagsIn(ctx),
	}

	var err error
	err = retry.RetryContext(ctx, 4*time.Minute, func() *retry.RetryError {
		_, err := conn.CreateMemberWithContext(ctx, input)

		if tfawserr.ErrCodeEquals(err, macie2.ErrorCodeClientError) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateMemberWithContext(ctx, input)
	}

	if err != nil {
		return diag.Errorf("creating Macie Member: %s", err)
	}

	d.SetId(accountId)

	if !d.Get("invite").(bool) {
		return resourceMemberRead(ctx, d, meta)
	}

	// Invitation workflow

	inputInvite := &macie2.CreateInvitationsInput{
		AccountIds: []*string{aws.String(d.Id())},
	}

	if v, ok := d.GetOk("invitation_disable_email_notification"); ok {
		inputInvite.DisableEmailNotification = aws.Bool(v.(bool))
	}
	if v, ok := d.GetOk("invitation_message"); ok {
		inputInvite.Message = aws.String(v.(string))
	}

	log.Printf("[INFO] Inviting Macie2 Member: %s", inputInvite)

	var output *macie2.CreateInvitationsOutput
	err = retry.RetryContext(ctx, 4*time.Minute, func() *retry.RetryError {
		output, err = conn.CreateInvitationsWithContext(ctx, inputInvite)

		if tfawserr.ErrCodeEquals(err, macie2.ErrorCodeClientError) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateInvitationsWithContext(ctx, inputInvite)
	}

	if err != nil {
		return diag.Errorf("inviting Macie Member: %s", err)
	}

	if len(output.UnprocessedAccounts) != 0 {
		return diag.Errorf("inviting Macie Member: %s: %s", aws.StringValue(output.UnprocessedAccounts[0].ErrorCode), aws.StringValue(output.UnprocessedAccounts[0].ErrorMessage))
	}

	if _, err = waitMemberInvited(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for Macie Member (%s) invitation: %s", d.Id(), err)
	}

	return resourceMemberRead(ctx, d, meta)
}

func resourceMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn(ctx)

	input := &macie2.GetMemberInput{
		Id: aws.String(d.Id()),
	}

	resp, err := conn.GetMemberWithContext(ctx, input)

	if !d.IsNewResource() && (tfawserr.ErrCodeEquals(err, macie2.ErrCodeResourceNotFoundException) ||
		tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") ||
		tfawserr.ErrMessageContains(err, macie2.ErrCodeConflictException, "member accounts are associated with your account") ||
		tfawserr.ErrMessageContains(err, macie2.ErrCodeValidationException, "account is not associated with your account")) {
		log.Printf("[WARN] Macie Member (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Macie Member (%s): %s", d.Id(), err)
	}

	d.Set("account_id", resp.AccountId)
	d.Set("email", resp.Email)
	d.Set("relationship_status", resp.RelationshipStatus)
	d.Set("administrator_account_id", resp.AdministratorAccountId)
	d.Set("master_account_id", resp.MasterAccountId)
	d.Set("invited_at", aws.TimeValue(resp.InvitedAt).Format(time.RFC3339))
	d.Set("updated_at", aws.TimeValue(resp.UpdatedAt).Format(time.RFC3339))
	d.Set("arn", resp.Arn)

	SetTagsOut(ctx, resp.Tags)

	status := aws.StringValue(resp.RelationshipStatus)
	log.Printf("[DEBUG] print resp.RelationshipStatus: %v", aws.StringValue(resp.RelationshipStatus))
	if status == macie2.RelationshipStatusEnabled ||
		status == macie2.RelationshipStatusInvited || status == macie2.RelationshipStatusEmailVerificationInProgress ||
		status == macie2.RelationshipStatusPaused {
		d.Set("invite", true)
	}

	if status == macie2.RelationshipStatusRemoved {
		d.Set("invite", false)
	}

	// To fake a result for status in order to avoid an error related to difference for ImportVerifyState
	// It sets to MacieStatusPaused because it can only be changed to PAUSED, normally when it's accepted its status is ENABLED
	status = macie2.MacieStatusEnabled
	if aws.StringValue(resp.RelationshipStatus) == macie2.RelationshipStatusPaused {
		status = macie2.MacieStatusPaused
	}
	d.Set("status", status)

	return nil
}

func resourceMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn(ctx)

	// Invitation workflow

	if d.HasChange("invite") {
		if d.Get("invite").(bool) {
			inputInvite := &macie2.CreateInvitationsInput{
				AccountIds: []*string{aws.String(d.Id())},
			}

			if v, ok := d.GetOk("invitation_disable_email_notification"); ok {
				inputInvite.DisableEmailNotification = aws.Bool(v.(bool))
			}
			if v, ok := d.GetOk("invitation_message"); ok {
				inputInvite.Message = aws.String(v.(string))
			}

			log.Printf("[INFO] Inviting Macie2 Member: %s", inputInvite)
			var output *macie2.CreateInvitationsOutput
			var err error
			err = retry.RetryContext(ctx, 4*time.Minute, func() *retry.RetryError {
				output, err = conn.CreateInvitationsWithContext(ctx, inputInvite)

				if tfawserr.ErrCodeEquals(err, macie2.ErrorCodeClientError) {
					return retry.RetryableError(err)
				}

				if err != nil {
					return retry.NonRetryableError(err)
				}

				return nil
			})

			if tfresource.TimedOut(err) {
				output, err = conn.CreateInvitationsWithContext(ctx, inputInvite)
			}

			if err != nil {
				return diag.Errorf("inviting Macie Member: %s", err)
			}

			if len(output.UnprocessedAccounts) != 0 {
				return diag.Errorf("inviting Macie Member: %s: %s", aws.StringValue(output.UnprocessedAccounts[0].ErrorCode), aws.StringValue(output.UnprocessedAccounts[0].ErrorMessage))
			}

			if _, err = waitMemberInvited(ctx, conn, d.Id()); err != nil {
				return diag.Errorf("waiting for Macie Member (%s) invitation: %s", d.Id(), err)
			}
		} else {
			input := &macie2.DisassociateMemberInput{
				Id: aws.String(d.Id()),
			}

			_, err := conn.DisassociateMemberWithContext(ctx, input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, macie2.ErrCodeResourceNotFoundException) ||
					tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") {
					return nil
				}
				return diag.Errorf("disassociating Macie Member invite (%s): %s", d.Id(), err)
			}
		}
	}

	// End Invitation workflow

	if d.HasChange("status") {
		input := &macie2.UpdateMemberSessionInput{
			Id:     aws.String(d.Id()),
			Status: aws.String(d.Get("status").(string)),
		}

		_, err := conn.UpdateMemberSessionWithContext(ctx, input)
		if err != nil {
			return diag.Errorf("updating Macie Member (%s): %s", d.Id(), err)
		}
	}

	return resourceMemberRead(ctx, d, meta)
}

func resourceMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn(ctx)

	input := &macie2.DeleteMemberInput{
		Id: aws.String(d.Id()),
	}

	_, err := conn.DeleteMemberWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, macie2.ErrCodeResourceNotFoundException) ||
			tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") ||
			tfawserr.ErrMessageContains(err, macie2.ErrCodeConflictException, "member accounts are associated with your account") ||
			tfawserr.ErrMessageContains(err, macie2.ErrCodeValidationException, "account is not associated with your account") {
			return nil
		}
		return diag.Errorf("deleting Macie Member (%s): %s", d.Id(), err)
	}
	return nil
}
