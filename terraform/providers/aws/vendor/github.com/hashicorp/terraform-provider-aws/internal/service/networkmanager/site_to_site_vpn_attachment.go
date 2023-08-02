package networkmanager

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_networkmanager_site_to_site_vpn_attachment", name="Site To Site VPN Attachment")
// @Tags(identifierAttribute="arn")
func ResourceSiteToSiteVPNAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSiteToSiteVPNAttachmentCreate,
		ReadWithoutTimeout:   resourceSiteToSiteVPNAttachmentRead,
		UpdateWithoutTimeout: resourceSiteToSiteVPNAttachmentUpdate,
		DeleteWithoutTimeout: resourceSiteToSiteVPNAttachmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_policy_rule_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_network_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpn_connection_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^arn:[^:]{1,63}:ec2:[^:]{0,63}:[^:]{0,63}:vpn-connection\/vpn-[0-9a-f]{8,17}$`), "Must be valid VPN ARN"),
			},
		},
	}
}

func resourceSiteToSiteVPNAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)

	coreNetworkID := d.Get("core_network_id").(string)
	vpnConnectionARN := d.Get("vpn_connection_arn").(string)
	input := &networkmanager.CreateSiteToSiteVpnAttachmentInput{
		CoreNetworkId:    aws.String(coreNetworkID),
		Tags:             GetTagsIn(ctx),
		VpnConnectionArn: aws.String(vpnConnectionARN),
	}

	output, err := conn.CreateSiteToSiteVpnAttachmentWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating Network Manager Site To Site VPN (%s) Attachment (%s): %s", vpnConnectionARN, coreNetworkID, err)
	}

	d.SetId(aws.StringValue(output.SiteToSiteVpnAttachment.Attachment.AttachmentId))

	if _, err := waitSiteToSiteVPNAttachmentCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for Network Manager Site To Site VPN Attachment (%s) create: %s", d.Id(), err)
	}

	return resourceSiteToSiteVPNAttachmentRead(ctx, d, meta)
}

func resourceSiteToSiteVPNAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)

	vpnAttachment, err := FindSiteToSiteVPNAttachmentByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Network Manager Site To Site VPN Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Network Manager Site To Site VPN Attachment (%s): %s", d.Id(), err)
	}

	a := vpnAttachment.Attachment
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "networkmanager",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("attachment/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("attachment_policy_rule_number", a.AttachmentPolicyRuleNumber)
	d.Set("attachment_type", a.AttachmentType)
	d.Set("core_network_arn", a.CoreNetworkArn)
	d.Set("core_network_id", a.CoreNetworkId)
	d.Set("edge_location", a.EdgeLocation)
	d.Set("owner_account_id", a.OwnerAccountId)
	d.Set("resource_arn", a.ResourceArn)
	d.Set("segment_name", a.SegmentName)
	d.Set("state", a.State)
	d.Set("vpn_connection_arn", a.ResourceArn)

	SetTagsOut(ctx, a.Tags)

	return nil
}

func resourceSiteToSiteVPNAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Tags only.
	return resourceSiteToSiteVPNAttachmentRead(ctx, d, meta)
}

func resourceSiteToSiteVPNAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)

	// If ResourceAttachmentAccepter is used, then VPN Attachment state
	// is never updated from StatePendingAttachmentAcceptance and the delete fails
	output, err := FindSiteToSiteVPNAttachmentByID(ctx, conn, d.Id())

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Network Manager Site To Site VPN Attachment (%s): %s", d.Id(), err)
	}

	if state := aws.StringValue(output.Attachment.State); state == networkmanager.AttachmentStatePendingAttachmentAcceptance || state == networkmanager.AttachmentStatePendingTagAcceptance {
		return diag.Errorf("cannot delete Network Manager Site To Site VPN Attachment (%s) in state: %s", d.Id(), state)
	}

	log.Printf("[DEBUG] Deleting Network Manager Site To Site VPN Attachment: %s", d.Id())
	_, err = conn.DeleteAttachmentWithContext(ctx, &networkmanager.DeleteAttachmentInput{
		AttachmentId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Network Manager Site To Site VPN Attachment (%s): %s", d.Id(), err)
	}

	if _, err := waitSiteToSiteVPNAttachmentDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for Network Manager Site To Site VPN Attachment (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindSiteToSiteVPNAttachmentByID(ctx context.Context, conn *networkmanager.NetworkManager, id string) (*networkmanager.SiteToSiteVpnAttachment, error) {
	input := &networkmanager.GetSiteToSiteVpnAttachmentInput{
		AttachmentId: aws.String(id),
	}

	output, err := conn.GetSiteToSiteVpnAttachmentWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.SiteToSiteVpnAttachment == nil || output.SiteToSiteVpnAttachment.Attachment == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.SiteToSiteVpnAttachment, nil
}

func statusSiteToSiteVPNAttachmentState(ctx context.Context, conn *networkmanager.NetworkManager, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindSiteToSiteVPNAttachmentByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Attachment.State), nil
	}
}

func waitSiteToSiteVPNAttachmentCreated(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.SiteToSiteVpnAttachment, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{networkmanager.AttachmentStateCreating, networkmanager.AttachmentStatePendingNetworkUpdate},
		Target:  []string{networkmanager.AttachmentStateAvailable, networkmanager.AttachmentStatePendingAttachmentAcceptance},
		Timeout: timeout,
		Refresh: statusSiteToSiteVPNAttachmentState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.SiteToSiteVpnAttachment); ok {
		return output, err
	}

	return nil, err
}

func waitSiteToSiteVPNAttachmentDeleted(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.SiteToSiteVpnAttachment, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        []string{networkmanager.AttachmentStateDeleting},
		Target:         []string{},
		Timeout:        timeout,
		Refresh:        statusSiteToSiteVPNAttachmentState(ctx, conn, id),
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.SiteToSiteVpnAttachment); ok {
		return output, err
	}

	return nil, err
}

func waitSiteToSiteVPNAttachmentAvailable(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.SiteToSiteVpnAttachment, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{networkmanager.AttachmentStateCreating, networkmanager.AttachmentStatePendingAttachmentAcceptance, networkmanager.AttachmentStatePendingNetworkUpdate},
		Target:  []string{networkmanager.AttachmentStateAvailable},
		Timeout: timeout,
		Refresh: statusSiteToSiteVPNAttachmentState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.SiteToSiteVpnAttachment); ok {
		return output, err
	}

	return nil, err
}
