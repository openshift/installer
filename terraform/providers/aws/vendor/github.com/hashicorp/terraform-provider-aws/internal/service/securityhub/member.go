package securityhub

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	// Associated is the member status naming for regions that do not support Organizations
	memberStatusAssociated = "Associated"
	memberStatusInvited    = "Invited"
	memberStatusEnabled    = "Enabled"
	memberStatusResigned   = "Resigned"
)

func ResourceMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceMemberCreate,
		Read:   resourceMemberRead,
		Delete: resourceMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"invite": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"master_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMemberCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn
	log.Printf("[DEBUG] Creating Security Hub member %s", d.Get("account_id").(string))

	resp, err := conn.CreateMembers(&securityhub.CreateMembersInput{
		AccountDetails: []*securityhub.AccountDetails{
			{
				AccountId: aws.String(d.Get("account_id").(string)),
				Email:     aws.String(d.Get("email").(string)),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("Error creating Security Hub member %s: %s", d.Get("account_id").(string), err)
	}

	if len(resp.UnprocessedAccounts) > 0 {
		return fmt.Errorf("Error creating Security Hub member %s: UnprocessedAccounts is not empty", d.Get("account_id").(string))
	}

	d.SetId(d.Get("account_id").(string))

	if d.Get("invite").(bool) {
		log.Printf("[INFO] Inviting Security Hub member %s", d.Id())
		iresp, err := conn.InviteMembers(&securityhub.InviteMembersInput{
			AccountIds: []*string{aws.String(d.Get("account_id").(string))},
		})

		if err != nil {
			return fmt.Errorf("Error inviting Security Hub member %s: %s", d.Id(), err)
		}

		if len(iresp.UnprocessedAccounts) > 0 {
			return fmt.Errorf("Error inviting Security Hub member %s: UnprocessedAccounts is not empty", d.Id())
		}
	}

	return resourceMemberRead(d, meta)
}

func resourceMemberRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	log.Printf("[DEBUG] Reading Security Hub member %s", d.Id())
	resp, err := conn.GetMembers(&securityhub.GetMembersInput{
		AccountIds: []*string{aws.String(d.Id())},
	})

	if err != nil {
		if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
			log.Printf("[WARN] Security Hub member (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if len(resp.Members) == 0 {
		log.Printf("[WARN] Security Hub member (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	member := resp.Members[0]

	d.Set("account_id", member.AccountId)
	d.Set("email", member.Email)
	d.Set("master_id", member.MasterId)

	status := aws.StringValue(member.MemberStatus)
	d.Set("member_status", status)

	invited := status == memberStatusInvited || status == memberStatusEnabled || status == memberStatusAssociated || status == memberStatusResigned
	d.Set("invite", invited)

	return nil
}

func resourceMemberDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecurityHubConn

	_, err := conn.DisassociateMembers(&securityhub.DisassociateMembersInput{
		AccountIds: []*string{aws.String(d.Id())},
	})
	if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error disassociating Security Hub member %s: %w", d.Id(), err)
	}

	resp, err := conn.DeleteMembers(&securityhub.DeleteMembersInput{
		AccountIds: []*string{aws.String(d.Id())},
	})

	if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error deleting Security Hub member %s: %w", d.Id(), err)
	}

	if len(resp.UnprocessedAccounts) > 0 {
		return fmt.Errorf("Error deleting Security Hub member %s: UnprocessedAccounts is not empty", d.Get("account_id").(string))
	}

	return nil
}
