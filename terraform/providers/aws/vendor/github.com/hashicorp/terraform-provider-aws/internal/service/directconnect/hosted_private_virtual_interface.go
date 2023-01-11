package directconnect

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceHostedPrivateVirtualInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostedPrivateVirtualInterfaceCreate,
		Read:   resourceHostedPrivateVirtualInterfaceRead,
		Delete: resourceHostedPrivateVirtualInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceHostedPrivateVirtualInterfaceImport,
		},

		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					directconnect.AddressFamilyIpv4,
					directconnect.AddressFamilyIpv6,
				}, false),
			},
			"amazon_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"amazon_side_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"bgp_auth_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"customer_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jumbo_frame_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mtu": {
				Type:         schema.TypeInt,
				Default:      1500,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1500, 9001}),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceHostedPrivateVirtualInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	req := &directconnect.AllocatePrivateVirtualInterfaceInput{
		ConnectionId: aws.String(d.Get("connection_id").(string)),
		NewPrivateVirtualInterfaceAllocation: &directconnect.NewPrivateVirtualInterfaceAllocation{
			AddressFamily:        aws.String(d.Get("address_family").(string)),
			Asn:                  aws.Int64(int64(d.Get("bgp_asn").(int))),
			Mtu:                  aws.Int64(int64(d.Get("mtu").(int))),
			VirtualInterfaceName: aws.String(d.Get("name").(string)),
			Vlan:                 aws.Int64(int64(d.Get("vlan").(int))),
		},
		OwnerAccount: aws.String(d.Get("owner_account_id").(string)),
	}
	if v, ok := d.GetOk("amazon_address"); ok {
		req.NewPrivateVirtualInterfaceAllocation.AmazonAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("bgp_auth_key"); ok {
		req.NewPrivateVirtualInterfaceAllocation.AuthKey = aws.String(v.(string))
	}
	if v, ok := d.GetOk("customer_address"); ok {
		req.NewPrivateVirtualInterfaceAllocation.CustomerAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("mtu"); ok {
		req.NewPrivateVirtualInterfaceAllocation.Mtu = aws.Int64(int64(v.(int)))
	}

	log.Printf("[DEBUG] Creating Direct Connect hosted private virtual interface: %s", req)
	resp, err := conn.AllocatePrivateVirtualInterface(req)
	if err != nil {
		return fmt.Errorf("error creating Direct Connect hosted private virtual interface: %s", err)
	}

	d.SetId(aws.StringValue(resp.VirtualInterfaceId))

	if err := hostedPrivateVirtualInterfaceWaitUntilAvailable(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	return resourceHostedPrivateVirtualInterfaceRead(d, meta)
}

func resourceHostedPrivateVirtualInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	vif, err := virtualInterfaceRead(d.Id(), conn)
	if err != nil {
		return err
	}
	if vif == nil {
		log.Printf("[WARN] Direct Connect hosted private virtual interface (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("address_family", vif.AddressFamily)
	d.Set("amazon_address", vif.AmazonAddress)
	d.Set("amazon_side_asn", strconv.FormatInt(aws.Int64Value(vif.AmazonSideAsn), 10))
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "directconnect",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("dxvif/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("aws_device", vif.AwsDeviceV2)
	d.Set("bgp_asn", vif.Asn)
	d.Set("bgp_auth_key", vif.AuthKey)
	d.Set("connection_id", vif.ConnectionId)
	d.Set("customer_address", vif.CustomerAddress)
	d.Set("jumbo_frame_capable", vif.JumboFrameCapable)
	d.Set("mtu", vif.Mtu)
	d.Set("name", vif.VirtualInterfaceName)
	d.Set("owner_account_id", vif.OwnerAccount)
	d.Set("vlan", vif.Vlan)

	return nil
}

func resourceHostedPrivateVirtualInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	return virtualInterfaceDelete(d, meta)
}

func resourceHostedPrivateVirtualInterfaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	vif, err := virtualInterfaceRead(d.Id(), conn)
	if err != nil {
		return nil, err
	}
	if vif == nil {
		return nil, fmt.Errorf("virtual interface (%s) not found", d.Id())
	}

	if vifType := aws.StringValue(vif.VirtualInterfaceType); vifType != "private" {
		return nil, fmt.Errorf("virtual interface (%s) has incorrect type: %s", d.Id(), vifType)
	}

	return []*schema.ResourceData{d}, nil
}

func hostedPrivateVirtualInterfaceWaitUntilAvailable(conn *directconnect.DirectConnect, vifId string, timeout time.Duration) error {
	return virtualInterfaceWaitUntilAvailable(
		conn,
		vifId,
		timeout,
		[]string{
			directconnect.VirtualInterfaceStatePending,
		},
		[]string{
			directconnect.VirtualInterfaceStateAvailable,
			directconnect.VirtualInterfaceStateConfirming,
			directconnect.VirtualInterfaceStateDown,
		})
}
