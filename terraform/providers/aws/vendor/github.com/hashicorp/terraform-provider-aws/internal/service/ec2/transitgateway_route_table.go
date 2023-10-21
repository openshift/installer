package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ec2_transit_gateway_route_table", name="Transit Gateway Route Table")
// @Tags(identifierAttribute="id")
func ResourceTransitGatewayRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTransitGatewayRouteTableCreate,
		ReadWithoutTimeout:   resourceTransitGatewayRouteTableRead,
		UpdateWithoutTimeout: resourceTransitGatewayRouteTableUpdate,
		DeleteWithoutTimeout: resourceTransitGatewayRouteTableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_association_route_table": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_propagation_route_table": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"transit_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceTransitGatewayRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.CreateTransitGatewayRouteTableInput{
		TransitGatewayId:  aws.String(d.Get("transit_gateway_id").(string)),
		TagSpecifications: getTagSpecificationsIn(ctx, ec2.ResourceTypeTransitGatewayRouteTable),
	}

	log.Printf("[DEBUG] Creating EC2 Transit Gateway Route Table: %s", input)
	output, err := conn.CreateTransitGatewayRouteTableWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 Transit Gateway Route Table: %s", err)
	}

	d.SetId(aws.StringValue(output.TransitGatewayRouteTable.TransitGatewayRouteTableId))

	if _, err := WaitTransitGatewayRouteTableCreated(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 Transit Gateway Route Table (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceTransitGatewayRouteTableRead(ctx, d, meta)...)
}

func resourceTransitGatewayRouteTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	transitGatewayRouteTable, err := FindTransitGatewayRouteTableByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Transit Gateway Route Table (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Transit Gateway Route Table (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("transit-gateway-route-table/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("default_association_route_table", transitGatewayRouteTable.DefaultAssociationRouteTable)
	d.Set("default_propagation_route_table", transitGatewayRouteTable.DefaultPropagationRouteTable)
	d.Set("transit_gateway_id", transitGatewayRouteTable.TransitGatewayId)

	SetTagsOut(ctx, transitGatewayRouteTable.Tags)

	return diags
}

func resourceTransitGatewayRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Tags only.

	return append(diags, resourceTransitGatewayRouteTableRead(ctx, d, meta)...)
}

func resourceTransitGatewayRouteTableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[DEBUG] Deleting EC2 Transit Gateway Route Table: %s", d.Id())
	_, err := conn.DeleteTransitGatewayRouteTableWithContext(ctx, &ec2.DeleteTransitGatewayRouteTableInput{
		TransitGatewayRouteTableId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidRouteTableIDNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 Transit Gateway Route Table (%s): %s", d.Id(), err)
	}

	if _, err := WaitTransitGatewayRouteTableDeleted(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 Transit Gateway Route Table (%s) delete: %s", d.Id(), err)
	}

	return diags
}
