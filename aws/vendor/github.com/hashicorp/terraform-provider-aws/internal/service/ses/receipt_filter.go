package ses

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_ses_receipt_filter")
func ResourceReceiptFilter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceReceiptFilterCreate,
		ReadWithoutTimeout:   resourceReceiptFilterRead,
		DeleteWithoutTimeout: resourceReceiptFilterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z._-]+$`), "must contain only alphanumeric, period, underscore, and hyphen characters"),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]`), "must begin with a alphanumeric character"),
					validation.StringMatch(regexp.MustCompile(`[0-9a-zA-Z]$`), "must end with a alphanumeric character"),
				),
			},

			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
					validation.IsIPv4Address,
				),
			},

			"policy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					ses.ReceiptFilterPolicyBlock,
					ses.ReceiptFilterPolicyAllow,
				}, false),
			},
		},
	}
}

func resourceReceiptFilterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn(ctx)

	name := d.Get("name").(string)

	createOpts := &ses.CreateReceiptFilterInput{
		Filter: &ses.ReceiptFilter{
			Name: aws.String(name),
			IpFilter: &ses.ReceiptIpFilter{
				Cidr:   aws.String(d.Get("cidr").(string)),
				Policy: aws.String(d.Get("policy").(string)),
			},
		},
	}

	_, err := conn.CreateReceiptFilterWithContext(ctx, createOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating SES receipt filter: %s", err)
	}

	d.SetId(name)

	return append(diags, resourceReceiptFilterRead(ctx, d, meta)...)
}

func resourceReceiptFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn(ctx)

	listOpts := &ses.ListReceiptFiltersInput{}

	response, err := conn.ListReceiptFiltersWithContext(ctx, listOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading SES Receipt Filter (%s): %s", d.Id(), err)
	}

	var filter *ses.ReceiptFilter

	for _, responseFilter := range response.Filters {
		if aws.StringValue(responseFilter.Name) == d.Id() {
			filter = responseFilter
			break
		}
	}

	if filter == nil {
		log.Printf("[WARN] SES Receipt Filter (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	d.Set("cidr", filter.IpFilter.Cidr)
	d.Set("policy", filter.IpFilter.Policy)
	d.Set("name", filter.Name)

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "ses",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("receipt-filter/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	return diags
}

func resourceReceiptFilterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SESConn(ctx)

	deleteOpts := &ses.DeleteReceiptFilterInput{
		FilterName: aws.String(d.Id()),
	}

	_, err := conn.DeleteReceiptFilterWithContext(ctx, deleteOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting SES receipt filter: %s", err)
	}

	return diags
}
