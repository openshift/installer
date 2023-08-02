package vpclattice

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_vpclattice_access_log_subscription", name="Access Log Subscription")
// @Tags(identifierAttribute="arn")
func ResourceAccessLogSubscription() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAccessLogSubscriptionCreate,
		ReadWithoutTimeout:   resourceAccessLogSubscriptionRead,
		UpdateWithoutTimeout: resourceAccessLogSubscriptionUpdate,
		DeleteWithoutTimeout: resourceAccessLogSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"resource_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameAccessLogSubscription = "Access Log Subscription"
)

func resourceAccessLogSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	in := &vpclattice.CreateAccessLogSubscriptionInput{
		ClientToken:        aws.String(id.UniqueId()),
		DestinationArn:     aws.String(d.Get("destination_arn").(string)),
		ResourceIdentifier: aws.String(d.Get("resource_identifier").(string)),
		Tags:               GetTagsIn(ctx),
	}

	out, err := conn.CreateAccessLogSubscription(ctx, in)

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameAccessLogSubscription, d.Get("destination_arn").(string), err)
	}

	d.SetId(aws.ToString(out.Id))

	return resourceAccessLogSubscriptionRead(ctx, d, meta)
}

func resourceAccessLogSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	out, err := findAccessLogSubscriptionByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] VPCLattice AccessLogSubscription (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameAccessLogSubscription, d.Id(), err)
	}

	d.Set("arn", out.Arn)
	d.Set("destination_arn", out.DestinationArn)
	d.Set("resource_arn", out.ResourceArn)
	d.Set("resource_identifier", out.ResourceId)

	return nil
}

func resourceAccessLogSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Tags only.
	return resourceAccessLogSubscriptionRead(ctx, d, meta)
}

func resourceAccessLogSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient(ctx)

	log.Printf("[INFO] Deleting VPCLattice AccessLogSubscription %s", d.Id())
	_, err := conn.DeleteAccessLogSubscription(ctx, &vpclattice.DeleteAccessLogSubscriptionInput{
		AccessLogSubscriptionIdentifier: aws.String(d.Id()),
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.VPCLattice, create.ErrActionDeleting, ResNameAccessLogSubscription, d.Id(), err)
	}

	return nil
}

func findAccessLogSubscriptionByID(ctx context.Context, conn *vpclattice.Client, id string) (*vpclattice.GetAccessLogSubscriptionOutput, error) {
	in := &vpclattice.GetAccessLogSubscriptionInput{
		AccessLogSubscriptionIdentifier: aws.String(id),
	}
	out, err := conn.GetAccessLogSubscription(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Id == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
