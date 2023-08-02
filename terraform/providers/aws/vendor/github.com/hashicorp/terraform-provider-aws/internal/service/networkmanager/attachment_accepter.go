package networkmanager

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// AttachmentAccepter does not require AttachmentType. However, querying attachments for status updates requires knowing tyupe
// To facilitate querying and waiters on specific attachment types, attachment_type set to required

// @SDKResource("aws_networkmanager_attachment_accepter")
func ResourceAttachmentAccepter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAttachmentAccepterCreate,
		ReadWithoutTimeout:   resourceAttachmentAccepterRead,
		DeleteWithoutTimeout: schema.NoopContext,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_policy_rule_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// querying attachments requires knowing the type ahead of time
			// therefore type is required in provider, though not on the API
			"attachment_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					networkmanager.AttachmentTypeVpc,
					networkmanager.AttachmentTypeSiteToSiteVpn,
					networkmanager.AttachmentTypeConnect,
				}, false),
			},
			"core_network_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edge_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAttachmentAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)

	var state string
	attachmentID := d.Get("attachment_id").(string)
	attachmentType := d.Get("attachment_type").(string)

	switch attachmentType {
	case networkmanager.AttachmentTypeVpc:
		vpcAttachment, err := FindVPCAttachmentByID(ctx, conn, attachmentID)

		if err != nil {
			return diag.Errorf("reading Network Manager VPC Attachment (%s): %s", attachmentID, err)
		}

		state = aws.StringValue(vpcAttachment.Attachment.State)

		d.SetId(attachmentID)

	case networkmanager.AttachmentTypeSiteToSiteVpn:
		vpnAttachment, err := FindSiteToSiteVPNAttachmentByID(ctx, conn, attachmentID)

		if err != nil {
			return diag.Errorf("reading Network Manager Site To Site VPN Attachment (%s): %s", attachmentID, err)
		}

		state = aws.StringValue(vpnAttachment.Attachment.State)

		d.SetId(attachmentID)

	case networkmanager.AttachmentTypeConnect:
		connectAttachment, err := FindConnectAttachmentByID(ctx, conn, attachmentID)

		if err != nil {
			return diag.Errorf("reading Network Manager Connect Attachment (%s): %s", attachmentID, err)
		}

		state = aws.StringValue(connectAttachment.Attachment.State)

		d.SetId(attachmentID)

	default:
		return diag.Errorf("unsupported Network Manager Attachment type: %s", attachmentType)
	}

	if state == networkmanager.AttachmentStatePendingAttachmentAcceptance || state == networkmanager.AttachmentStatePendingTagAcceptance {
		input := &networkmanager.AcceptAttachmentInput{
			AttachmentId: aws.String(attachmentID),
		}

		_, err := conn.AcceptAttachmentWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("accepting Network Manager Attachment (%s): %s", attachmentID, err)
		}

		switch attachmentType {
		case networkmanager.AttachmentTypeVpc:
			if _, err := waitVPCAttachmentCreated(ctx, conn, attachmentID, d.Timeout(schema.TimeoutCreate)); err != nil {
				return diag.Errorf("waiting for Network Manager VPC Attachment (%s) create: %s", attachmentID, err)
			}

		case networkmanager.AttachmentTypeSiteToSiteVpn:
			if _, err := waitSiteToSiteVPNAttachmentAvailable(ctx, conn, attachmentID, d.Timeout(schema.TimeoutCreate)); err != nil {
				return diag.Errorf("waiting for Network Manager VPN Attachment (%s) create: %s", attachmentID, err)
			}

		case networkmanager.AttachmentTypeConnect:
			if _, err := waitConnectAttachmentAvailable(ctx, conn, attachmentID, d.Timeout(schema.TimeoutCreate)); err != nil {
				return diag.Errorf("waiting for Network Manager Connect Attachment (%s) create: %s", attachmentID, err)
			}
		}
	}

	return resourceAttachmentAccepterRead(ctx, d, meta)
}

func resourceAttachmentAccepterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)

	var a *networkmanager.Attachment

	switch aType := d.Get("attachment_type"); aType {
	case networkmanager.AttachmentTypeVpc:
		vpcAttachment, err := FindVPCAttachmentByID(ctx, conn, d.Id())

		if !d.IsNewResource() && tfresource.NotFound(err) {
			log.Printf("[WARN] Network Manager VPC Attachment %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		if err != nil {
			return diag.Errorf("reading Network Manager VPC Attachment (%s): %s", d.Id(), err)
		}

		a = vpcAttachment.Attachment

	case networkmanager.AttachmentTypeSiteToSiteVpn:
		vpnAttachment, err := FindSiteToSiteVPNAttachmentByID(ctx, conn, d.Id())

		if !d.IsNewResource() && tfresource.NotFound(err) {
			log.Printf("[WARN] Network Manager Site To Site VPN Attachment %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		if err != nil {
			return diag.Errorf("reading Network Manager Site To Site VPN Attachment (%s): %s", d.Id(), err)
		}

		a = vpnAttachment.Attachment

	case networkmanager.AttachmentTypeConnect:
		connectAttachment, err := FindConnectAttachmentByID(ctx, conn, d.Id())

		if !d.IsNewResource() && tfresource.NotFound(err) {
			log.Printf("[WARN] Network Manager Connect Attachment %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		if err != nil {
			return diag.Errorf("reading Network Manager Connect Attachment (%s): %s", d.Id(), err)
		}

		a = connectAttachment.Attachment
	}

	d.Set("attachment_policy_rule_number", a.AttachmentPolicyRuleNumber)
	d.Set("core_network_arn", a.CoreNetworkArn)
	d.Set("core_network_id", a.CoreNetworkId)
	d.Set("edge_location", a.EdgeLocation)
	d.Set("owner_account_id", a.OwnerAccountId)
	d.Set("resource_arn", a.ResourceArn)
	d.Set("segment_name", a.SegmentName)
	d.Set("state", a.State)

	return nil
}
