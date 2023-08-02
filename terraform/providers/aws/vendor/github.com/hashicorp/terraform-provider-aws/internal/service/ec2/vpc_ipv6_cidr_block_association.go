package ec2

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_vpc_ipv6_cidr_block_association")
func ResourceVPCIPv6CIDRBlockAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVPCIPv6CIDRBlockAssociationCreate,
		ReadWithoutTimeout:   resourceVPCIPv6CIDRBlockAssociationRead,
		DeleteWithoutTimeout: resourceVPCIPv6CIDRBlockAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// ipv6_cidr_block can be set by a value returned from IPAM or explicitly in config.
			if diff.Id() != "" && diff.HasChange("ipv6_cidr_block") {
				// If netmask is set then ipv6_cidr_block is derived from IPAM, ignore changes.
				if diff.Get("ipv6_netmask_length") != 0 {
					return diff.Clear("ipv6_cidr_block")
				}
				return diff.ForceNew("ipv6_cidr_block")
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					verify.ValidIPv6CIDRNetworkAddress,
					validation.IsCIDRNetwork(VPCCIDRMaxIPv6, VPCCIDRMaxIPv6)),
			},
			// ipam parameters are not required by the API but other usage mechanisms are not implemented yet. TODO ipv6 options:
			// --amazon-provided-ipv6-cidr-block
			// --ipv6-pool
			"ipv6_ipam_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_netmask_length": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IntInSlice([]int{VPCCIDRMaxIPv6}),
				ConflictsWith: []string{"ipv6_cidr_block"},
				// This RequiredWith setting should be applied once L57 is completed
				// RequiredWith:  []string{"ipv6_ipam_pool_id"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceVPCIPv6CIDRBlockAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	vpcID := d.Get("vpc_id").(string)
	input := &ec2.AssociateVpcCidrBlockInput{
		VpcId: aws.String(vpcID),
	}

	if v, ok := d.GetOk("ipv6_cidr_block"); ok {
		input.Ipv6CidrBlock = aws.String(v.(string))
	}

	if v, ok := d.GetOk("ipv6_ipam_pool_id"); ok {
		input.Ipv6IpamPoolId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("ipv6_netmask_length"); ok {
		input.Ipv6NetmaskLength = aws.Int64(int64(v.(int)))
	}

	log.Printf("[DEBUG] Creating EC2 VPC IPv6 CIDR Block Association: %s", input)
	output, err := conn.AssociateVpcCidrBlockWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 VPC (%s) IPv6 CIDR Block Association: %s", vpcID, err)
	}

	d.SetId(aws.StringValue(output.Ipv6CidrBlockAssociation.AssociationId))

	_, err = WaitVPCIPv6CIDRBlockAssociationCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPC (%s) IPv6 CIDR block (%s) to become associated: %s", vpcID, d.Id(), err)
	}

	return append(diags, resourceVPCIPv6CIDRBlockAssociationRead(ctx, d, meta)...)
}

func resourceVPCIPv6CIDRBlockAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	vpcIpv6CidrBlockAssociation, vpc, err := FindVPCIPv6CIDRBlockAssociationByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VPC IPv6 CIDR Block Association %s not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC IPv6 CIDR Block Association (%s): %s", d.Id(), err)
	}

	d.Set("ipv6_cidr_block", vpcIpv6CidrBlockAssociation.Ipv6CidrBlock)
	d.Set("vpc_id", vpc.VpcId)

	return diags
}

func resourceVPCIPv6CIDRBlockAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[DEBUG] Deleting VPC IPv6 CIDR Block Association: %s", d.Id())
	_, err := conn.DisassociateVpcCidrBlockWithContext(ctx, &ec2.DisassociateVpcCidrBlockInput{
		AssociationId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidVPCCIDRBlockAssociationIDNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 VPC IPv6 CIDR Block Association (%s): %s", d.Id(), err)
	}

	_, err = WaitVPCIPv6CIDRBlockAssociationDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPC IPv6 CIDR block (%s) to become disassociated: %s", d.Id(), err)
	}

	return diags
}
