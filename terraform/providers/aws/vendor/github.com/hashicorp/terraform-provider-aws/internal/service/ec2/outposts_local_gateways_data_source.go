package ec2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// @SDKDataSource("aws_ec2_local_gateways")
func DataSourceLocalGateways() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceLocalGatewaysRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceLocalGatewaysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.DescribeLocalGatewaysInput{}

	input.Filters = append(input.Filters, BuildTagFilterList(
		Tags(tftags.New(ctx, d.Get("tags").(map[string]interface{}))),
	)...)

	input.Filters = append(input.Filters, BuildFiltersDataSource(
		d.Get("filter").(*schema.Set),
	)...)

	if len(input.Filters) == 0 {
		input.Filters = nil
	}

	output, err := FindLocalGateways(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Local Gateways: %s", err)
	}

	var gatewayIDs []string

	for _, v := range output {
		gatewayIDs = append(gatewayIDs, aws.StringValue(v.LocalGatewayId))
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("ids", gatewayIDs)

	return diags
}
