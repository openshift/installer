package ec2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_vpn_connection_route")
func ResourceVPNConnectionRoute() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVPNConnectionRouteCreate,
		ReadWithoutTimeout:   resourceVPNConnectionRouteRead,
		DeleteWithoutTimeout: resourceVPNConnectionRouteDelete,

		Schema: map[string]*schema.Schema{
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVPNConnectionRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	cidrBlock := d.Get("destination_cidr_block").(string)
	vpnConnectionID := d.Get("vpn_connection_id").(string)
	id := VPNConnectionRouteCreateResourceID(cidrBlock, vpnConnectionID)
	input := &ec2.CreateVpnConnectionRouteInput{
		DestinationCidrBlock: aws.String(cidrBlock),
		VpnConnectionId:      aws.String(vpnConnectionID),
	}

	log.Printf("[DEBUG] Creating EC2 VPN Connection Route: %s", input)
	_, err := conn.CreateVpnConnectionRouteWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 VPN Connection Route (%s): %s", id, err)
	}

	d.SetId(id)

	if _, err := WaitVPNConnectionRouteCreated(ctx, conn, vpnConnectionID, cidrBlock); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPN Connection Route (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceVPNConnectionRouteRead(ctx, d, meta)...)
}

func resourceVPNConnectionRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	cidrBlock, vpnConnectionID, err := VPNConnectionRouteParseResourceID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPN Connection Route (%s): %s", d.Id(), err)
	}

	_, err = FindVPNConnectionRouteByVPNConnectionIDAndCIDR(ctx, conn, vpnConnectionID, cidrBlock)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VPN Connection Route (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPN Connection Route (%s): %s", d.Id(), err)
	}

	d.Set("destination_cidr_block", cidrBlock)
	d.Set("vpn_connection_id", vpnConnectionID)

	return diags
}

func resourceVPNConnectionRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	cidrBlock, vpnConnectionID, err := VPNConnectionRouteParseResourceID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 VPN Connection Route (%s): %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleting EC2 VPN Connection Route: %s", d.Id())
	_, err = conn.DeleteVpnConnectionRouteWithContext(ctx, &ec2.DeleteVpnConnectionRouteInput{
		DestinationCidrBlock: aws.String(cidrBlock),
		VpnConnectionId:      aws.String(vpnConnectionID),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidVPNConnectionIDNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 VPN Connection Route (%s): %s", d.Id(), err)
	}

	if _, err := WaitVPNConnectionRouteDeleted(ctx, conn, vpnConnectionID, cidrBlock); err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 VPN Connection Route (%s): waiting for completion: %s", d.Id(), err)
	}

	return diags
}

const vpnConnectionRouteResourceIDSeparator = ":"

func VPNConnectionRouteCreateResourceID(cidrBlock, vpcConnectionID string) string {
	parts := []string{cidrBlock, vpcConnectionID}
	id := strings.Join(parts, vpnConnectionRouteResourceIDSeparator)

	return id
}

func VPNConnectionRouteParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, vpnConnectionRouteResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected DestinationCIDRBlock%[2]sVPNConnectionID", id, vpnConnectionRouteResourceIDSeparator)
}
