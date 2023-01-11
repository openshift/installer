package macie2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceInvitationAccepter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceInvitationAccepterCreate,
		ReadWithoutTimeout:   resourceInvitationAccepterRead,
		DeleteWithoutTimeout: resourceInvitationAccepterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"administrator_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"invitation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceInvitationAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	adminAccountID := d.Get("administrator_account_id").(string)
	var invitationID string

	listInvitationsInput := &macie2.ListInvitationsInput{}

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := conn.ListInvitationsPages(listInvitationsInput, func(page *macie2.ListInvitationsOutput, lastPage bool) bool {
			for _, invitation := range page.Invitations {
				if aws.StringValue(invitation.AccountId) == adminAccountID {
					invitationID = aws.StringValue(invitation.InvitationId)
					return false
				}
			}
			return !lastPage
		})

		if err != nil {
			return resource.NonRetryableError(err)
		}

		if invitationID == "" {
			return resource.RetryableError(fmt.Errorf("unable to find pending Macie Invitation for administrator account ID (%s)", adminAccountID))
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		err = conn.ListInvitationsPages(listInvitationsInput, func(page *macie2.ListInvitationsOutput, lastPage bool) bool {
			for _, invitation := range page.Invitations {
				if aws.StringValue(invitation.AccountId) == adminAccountID {
					invitationID = aws.StringValue(invitation.InvitationId)
					return false
				}
			}
			return !lastPage
		})
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing Macie InvitationAccepter (%s): %w", d.Id(), err))
	}

	acceptInvitationInput := &macie2.AcceptInvitationInput{
		InvitationId:           aws.String(invitationID),
		AdministratorAccountId: aws.String(adminAccountID),
	}

	_, err = conn.AcceptInvitationWithContext(ctx, acceptInvitationInput)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error accepting Macie InvitationAccepter (%s): %w", d.Id(), err))
	}

	d.SetId(adminAccountID)

	return resourceInvitationAccepterRead(ctx, d, meta)
}

func resourceInvitationAccepterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	var err error

	input := &macie2.GetAdministratorAccountInput{}

	output, err := conn.GetAdministratorAccountWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, macie2.ErrCodeResourceNotFoundException) ||
		tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") {
		log.Printf("[WARN] Macie InvitationAccepter (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Macie InvitationAccepter (%s): %w", d.Id(), err))
	}

	if output == nil || output.Administrator == nil {
		return diag.FromErr(fmt.Errorf("error reading Macie InvitationAccepter (%s): %w", d.Id(), err))
	}

	d.Set("administrator_account_id", output.Administrator.AccountId)
	d.Set("invitation_id", output.Administrator.InvitationId)
	return nil
}

func resourceInvitationAccepterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Macie2Conn

	input := &macie2.DisassociateFromAdministratorAccountInput{}

	_, err := conn.DisassociateFromAdministratorAccountWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, macie2.ErrCodeResourceNotFoundException) ||
			tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error disassociating Macie InvitationAccepter (%s): %w", d.Id(), err))
	}
	return nil
}
