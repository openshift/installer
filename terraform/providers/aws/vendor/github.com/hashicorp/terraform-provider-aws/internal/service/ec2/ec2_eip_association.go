package ec2

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceEIPAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPAssociationCreate,
		Read:   resourceEIPAssociationRead,
		Delete: resourceEIPAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"allow_reassociation": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEIPAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	input := &ec2.AssociateAddressInput{}

	if v, ok := d.GetOk("allocation_id"); ok {
		input.AllocationId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("allow_reassociation"); ok {
		input.AllowReassociation = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		input.InstanceId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("network_interface_id"); ok {
		input.NetworkInterfaceId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		input.PrivateIpAddress = aws.String(v.(string))
	}

	if v, ok := d.GetOk("public_ip"); ok {
		input.PublicIp = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating EC2 EIP Association: %s", input)
	output, err := conn.AssociateAddress(input)

	if err != nil {
		return fmt.Errorf("creating EC2 EIP Association: %w", err)
	}

	if output.AssociationId != nil {
		d.SetId(aws.StringValue(output.AssociationId))

		_, err = tfresource.RetryWhen(propagationTimeout,
			func() (interface{}, error) {
				return FindEIPByAssociationID(conn, d.Id())
			},
			func(err error) (bool, error) {
				if tfresource.NotFound(err) {
					return true, err
				}

				// "InvalidInstanceID: The pending instance 'i-0504e5b44ea06d599' is not in a valid state for this operation."
				if tfawserr.ErrMessageContains(err, errCodeInvalidInstanceID, "pending instance") {
					return true, err
				}

				return false, err
			},
		)

		if err != nil {
			return fmt.Errorf("waiting for EC2 EIP Association (%s) create: %w", d.Id(), err)
		}
	} else {
		// EC2-Classic.
		publicIP := aws.StringValue(input.PublicIp)
		d.SetId(publicIP)

		instanceID := aws.StringValue(input.InstanceId)
		if err := waitForAddressAssociationClassic(conn, publicIP, instanceID); err != nil {
			return fmt.Errorf("waiting for EC2 EIP (%s) to associate with EC2-Classic Instance (%s): %w", publicIP, instanceID, err)
		}
	}

	return resourceEIPAssociationRead(d, meta)
}

func resourceEIPAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	var err error
	var address *ec2.Address

	if eipAssociationID(d.Id()).IsVPC() {
		address, err = FindEIPByAssociationID(conn, d.Id())
	} else {
		address, err = FindEIPByPublicIP(conn, d.Id())
	}

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 EIP Association (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 EIP Association (%s): %w", d.Id(), err)
	}

	d.Set("allocation_id", address.AllocationId)
	d.Set("instance_id", address.InstanceId)
	d.Set("network_interface_id", address.NetworkInterfaceId)
	d.Set("private_ip_address", address.PrivateIpAddress)
	d.Set("public_ip", address.PublicIp)

	return nil
}

func resourceEIPAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	input := &ec2.DisassociateAddressInput{}

	if eipAssociationID(d.Id()).IsVPC() {
		input.AssociationId = aws.String(d.Id())
	} else {
		input.PublicIp = aws.String(d.Id())
	}

	log.Printf("[DEBUG] Deleting EC2 EIP Association: %s", d.Id())
	_, err := conn.DisassociateAddress(input)

	if tfawserr.ErrCodeEquals(err, errCodeInvalidAssociationIDNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting EC2 EIP Association (%s): %w", d.Id(), err)
	}

	return nil
}

type eipAssociationID string

// IsVPC returns whether or not the associated EIP is in the VPC domain.
func (id eipAssociationID) IsVPC() bool {
	return strings.HasPrefix(string(id), "eipassoc-")
}
