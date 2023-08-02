package route53resolver

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_route53_resolver_query_log_config", name="Query Log Config")
// @Tags(identifierAttribute="arn")
func ResourceQueryLogConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceQueryLogConfigCreate,
		ReadWithoutTimeout:   resourceQueryLogConfigRead,
		UpdateWithoutTimeout: resourceQueryLogConfigUpdate,
		DeleteWithoutTimeout: resourceQueryLogConfigDelete,

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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validResolverName,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceQueryLogConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Route53ResolverConn(ctx)

	name := d.Get("name").(string)
	input := &route53resolver.CreateResolverQueryLogConfigInput{
		CreatorRequestId: aws.String(id.PrefixedUniqueId("tf-r53-resolver-query-log-config-")),
		DestinationArn:   aws.String(d.Get("destination_arn").(string)),
		Name:             aws.String(name),
		Tags:             GetTagsIn(ctx),
	}

	output, err := conn.CreateResolverQueryLogConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating Route53 Resolver Query Log Config (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.ResolverQueryLogConfig.Id))

	if _, err := waitQueryLogConfigCreated(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for Route53 Resolver Query Log Config (%s) create: %s", d.Id(), err)
	}

	return resourceQueryLogConfigRead(ctx, d, meta)
}

func resourceQueryLogConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Route53ResolverConn(ctx)

	queryLogConfig, err := FindResolverQueryLogConfigByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Route53 Resolver Query Log Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Route53 Resolver Query Log Config (%s): %s", d.Id(), err)
	}

	arn := aws.StringValue(queryLogConfig.Arn)
	d.Set("arn", arn)
	d.Set("destination_arn", queryLogConfig.DestinationArn)
	d.Set("name", queryLogConfig.Name)
	d.Set("owner_id", queryLogConfig.OwnerId)
	d.Set("share_status", queryLogConfig.ShareStatus)

	return nil
}

func resourceQueryLogConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Tags only.
	return resourceQueryLogConfigRead(ctx, d, meta)
}

func resourceQueryLogConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).Route53ResolverConn(ctx)

	log.Printf("[DEBUG] Deleting Route53 Resolver Query Log Config: %s", d.Id())
	_, err := conn.DeleteResolverQueryLogConfigWithContext(ctx, &route53resolver.DeleteResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Route53 Resolver Query Log Config (%s): %s", d.Id(), err)
	}

	if _, err := waitQueryLogConfigDeleted(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for Route53 Resolver Query Log Config (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindResolverQueryLogConfigByID(ctx context.Context, conn *route53resolver.Route53Resolver, id string) (*route53resolver.ResolverQueryLogConfig, error) {
	input := &route53resolver.GetResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: aws.String(id),
	}

	output, err := conn.GetResolverQueryLogConfigWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, route53resolver.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.ResolverQueryLogConfig == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.ResolverQueryLogConfig, nil
}

func statusQueryLogConfig(ctx context.Context, conn *route53resolver.Route53Resolver, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindResolverQueryLogConfigByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

const (
	queryLogConfigCreatedTimeout = 5 * time.Minute
	queryLogConfigDeletedTimeout = 5 * time.Minute
)

func waitQueryLogConfigCreated(ctx context.Context, conn *route53resolver.Route53Resolver, id string) (*route53resolver.ResolverQueryLogConfig, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{route53resolver.ResolverQueryLogConfigStatusCreating},
		Target:  []string{route53resolver.ResolverQueryLogConfigStatusCreated},
		Refresh: statusQueryLogConfig(ctx, conn, id),
		Timeout: queryLogConfigCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*route53resolver.ResolverQueryLogConfig); ok {
		return output, err
	}

	return nil, err
}

func waitQueryLogConfigDeleted(ctx context.Context, conn *route53resolver.Route53Resolver, id string) (*route53resolver.ResolverQueryLogConfig, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{route53resolver.ResolverQueryLogConfigStatusDeleting},
		Target:  []string{},
		Refresh: statusQueryLogConfig(ctx, conn, id),
		Timeout: queryLogConfigDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*route53resolver.ResolverQueryLogConfig); ok {
		return output, err
	}

	return nil, err
}
