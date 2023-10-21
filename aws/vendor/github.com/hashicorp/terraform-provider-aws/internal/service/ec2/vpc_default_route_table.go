package ec2

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_default_route_table", name="Route Table")
// @Tags(identifierAttribute="id")
func ResourceDefaultRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceDefaultRouteTableCreate,
		ReadWithoutTimeout:   resourceDefaultRouteTableRead,
		UpdateWithoutTimeout: resourceRouteTableUpdate,
		DeleteWithoutTimeout: schema.NoopContext,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDefaultRouteTableImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},

		//
		// The top-level attributes must be a superset of the aws_route_table resource's attributes as common CRUD handlers are used.
		//
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"propagating_vgws": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route": {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Computed:   true,
				Optional:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						///
						// Destinations.
						///
						"cidr_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidIPv4CIDRNetworkAddress,
						},
						"destination_prefix_list_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_cidr_block": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidIPv6CIDRNetworkAddress,
						},
						//
						// Targets.
						// These target attributes are a subset of the aws_route_table resource's target attributes
						// as there are some targets that are not allowed in the default route table for a VPC.
						//
						"core_network_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"egress_only_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"transit_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_endpoint_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_peering_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceRouteTableHash,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceDefaultRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	routeTableID := d.Get("default_route_table_id").(string)

	routeTable, err := FindRouteTableByID(ctx, conn, routeTableID)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Default Route Table (%s): %s", routeTableID, err)
	}

	d.SetId(aws.StringValue(routeTable.RouteTableId))

	// Remove all existing VGW associations.
	for _, v := range routeTable.PropagatingVgws {
		if err := routeTableDisableVGWRoutePropagation(ctx, conn, d.Id(), aws.StringValue(v.GatewayId)); err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	// Delete all existing routes.
	for _, v := range routeTable.Routes {
		if gatewayID := aws.StringValue(v.GatewayId); gatewayID == gatewayIDLocal || gatewayID == gatewayIDVPCLattice {
			continue
		}

		if aws.StringValue(v.Origin) == ec2.RouteOriginEnableVgwRoutePropagation {
			continue
		}

		if v.DestinationPrefixListId != nil && strings.HasPrefix(aws.StringValue(v.GatewayId), "vpce-") {
			// Skipping because VPC endpoint routes are handled separately
			// See aws_vpc_endpoint
			continue
		}

		input := &ec2.DeleteRouteInput{
			RouteTableId: aws.String(d.Id()),
		}

		var destination string
		var routeFinder RouteFinder

		if v.DestinationCidrBlock != nil {
			input.DestinationCidrBlock = v.DestinationCidrBlock
			destination = aws.StringValue(v.DestinationCidrBlock)
			routeFinder = FindRouteByIPv4Destination
		} else if v.DestinationIpv6CidrBlock != nil {
			input.DestinationIpv6CidrBlock = v.DestinationIpv6CidrBlock
			destination = aws.StringValue(v.DestinationIpv6CidrBlock)
			routeFinder = FindRouteByIPv6Destination
		} else if v.DestinationPrefixListId != nil {
			input.DestinationPrefixListId = v.DestinationPrefixListId
			destination = aws.StringValue(v.DestinationPrefixListId)
			routeFinder = FindRouteByPrefixListIDDestination
		}

		_, err := conn.DeleteRouteWithContext(ctx, input)

		if tfawserr.ErrCodeEquals(err, errCodeInvalidRouteNotFound) {
			continue
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deleting Route in EC2 Default Route Table (%s) with destination (%s): %s", d.Id(), destination, err)
		}

		if _, err := WaitRouteDeleted(ctx, conn, routeFinder, routeTableID, destination, d.Timeout(schema.TimeoutCreate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Route in EC2 Default Route Table (%s) with destination (%s) delete: %s", d.Id(), destination, err)
		}
	}

	// Add new VGW associations.
	if v, ok := d.GetOk("propagating_vgws"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(string)

			if err := routeTableEnableVGWRoutePropagation(ctx, conn, d.Id(), v, d.Timeout(schema.TimeoutCreate)); err != nil {
				return sdkdiag.AppendFromErr(diags, err)
			}
		}
	}

	// Add new routes.
	if v, ok := d.GetOk("route"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(map[string]interface{})

			if err := routeTableAddRoute(ctx, conn, d.Id(), v, d.Timeout(schema.TimeoutCreate)); err != nil {
				return sdkdiag.AppendFromErr(diags, err)
			}
		}
	}

	if err := createTags(ctx, conn, d.Id(), GetTagsIn(ctx)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting EC2 Default Route Table (%s) tags: %s", d.Id(), err)
	}

	return append(diags, resourceDefaultRouteTableRead(ctx, d, meta)...)
}

func resourceDefaultRouteTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Set("default_route_table_id", d.Id())

	// Re-use regular AWS Route Table READ.
	// This is an extra API call but saves us from trying to manually keep parity.
	return resourceRouteTableRead(ctx, d, meta)
}

func resourceDefaultRouteTableImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	routeTable, err := FindMainRouteTableByVPCID(ctx, conn, d.Id())

	if err != nil {
		return nil, err
	}

	d.SetId(aws.StringValue(routeTable.RouteTableId))

	return []*schema.ResourceData{d}, nil
}
