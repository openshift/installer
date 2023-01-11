package ec2

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceTransitGatewayRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransitGatewayRouteTableAssociationCreate,
		Read:   resourceTransitGatewayRouteTableAssociationRead,
		Delete: resourceTransitGatewayRouteTableAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_gateway_attachment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"transit_gateway_route_table_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceTransitGatewayRouteTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayAttachmentID := d.Get("transit_gateway_attachment_id").(string)
	transitGatewayRouteTableID := d.Get("transit_gateway_route_table_id").(string)
	id := TransitGatewayRouteTableAssociationCreateResourceID(transitGatewayRouteTableID, transitGatewayAttachmentID)
	input := &ec2.AssociateTransitGatewayRouteTableInput{
		TransitGatewayAttachmentId: aws.String(transitGatewayAttachmentID),
		TransitGatewayRouteTableId: aws.String(transitGatewayRouteTableID),
	}

	_, err := conn.AssociateTransitGatewayRouteTable(input)

	if err != nil {
		return fmt.Errorf("creating EC2 Transit Gateway Route Table Association (%s): %w", id, err)
	}

	d.SetId(id)

	if _, err := WaitTransitGatewayRouteTableAssociationCreated(conn, transitGatewayRouteTableID, transitGatewayAttachmentID); err != nil {
		return fmt.Errorf("waiting for EC2 Transit Gateway Route Table Association (%s) create: %w", d.Id(), err)
	}

	return resourceTransitGatewayRouteTableAssociationRead(d, meta)
}

func resourceTransitGatewayRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayRouteTableID, transitGatewayAttachmentID, err := TransitGatewayRouteTableAssociationParseResourceID(d.Id())

	if err != nil {
		return err
	}

	transitGatewayRouteTableAssociation, err := FindTransitGatewayRouteTableAssociationByTwoPartKey(conn, transitGatewayRouteTableID, transitGatewayAttachmentID)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Transit Gateway Route Table Association %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 Transit Gateway Route Table Association (%s): %w", d.Id(), err)
	}

	d.Set("resource_id", transitGatewayRouteTableAssociation.ResourceId)
	d.Set("resource_type", transitGatewayRouteTableAssociation.ResourceType)
	d.Set("transit_gateway_attachment_id", transitGatewayRouteTableAssociation.TransitGatewayAttachmentId)
	d.Set("transit_gateway_route_table_id", transitGatewayRouteTableID)

	return nil
}

func resourceTransitGatewayRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayRouteTableID, transitGatewayAttachmentID, err := TransitGatewayRouteTableAssociationParseResourceID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting EC2 Transit Gateway Route Table Association: %s", d.Id())
	_, err = conn.DisassociateTransitGatewayRouteTable(&ec2.DisassociateTransitGatewayRouteTableInput{
		TransitGatewayAttachmentId: aws.String(transitGatewayAttachmentID),
		TransitGatewayRouteTableId: aws.String(transitGatewayRouteTableID),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidRouteTableIDNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting EC2 Transit Gateway Route Table Association (%s): %w", d.Id(), err)
	}

	if _, err := WaitTransitGatewayRouteTableAssociationDeleted(conn, transitGatewayRouteTableID, transitGatewayAttachmentID); err != nil {
		return fmt.Errorf("waiting for EC2 Transit Gateway Route Table Association (%s) delete: %w", d.Id(), err)
	}

	return nil
}

// transitGatewayRouteTableAssociationUpdate is used by Transit Gateway attachment resources to modify their route table associations.
// The route table ID may be empty (e.g. when the Transit Gateway itself has default route table association disabled).
func transitGatewayRouteTableAssociationUpdate(conn *ec2.EC2, transitGatewayRouteTableID, transitGatewayAttachmentID string, associate bool) error {
	if transitGatewayRouteTableID == "" {
		// Do nothing if no route table was specified.
		return nil
	}

	id := TransitGatewayRouteTableAssociationCreateResourceID(transitGatewayRouteTableID, transitGatewayAttachmentID)
	_, err := FindTransitGatewayRouteTableAssociationByTwoPartKey(conn, transitGatewayRouteTableID, transitGatewayAttachmentID)

	if tfresource.NotFound(err) {
		if associate {
			input := &ec2.AssociateTransitGatewayRouteTableInput{
				TransitGatewayAttachmentId: aws.String(transitGatewayAttachmentID),
				TransitGatewayRouteTableId: aws.String(transitGatewayRouteTableID),
			}

			_, err := conn.AssociateTransitGatewayRouteTable(input)

			if err != nil {
				return fmt.Errorf("creating EC2 Transit Gateway Route Table Association (%s): %w", id, err)
			}

			if _, err := WaitTransitGatewayRouteTableAssociationCreated(conn, transitGatewayRouteTableID, transitGatewayAttachmentID); err != nil {
				return fmt.Errorf("waiting for EC2 Transit Gateway Route Table Association (%s) create: %w", id, err)
			}
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 Transit Gateway Route Table Association (%s): %w", id, err)
	}

	if !associate {
		// Disassociation must be done only on already associated state.
		if _, err := WaitTransitGatewayRouteTableAssociationCreated(conn, transitGatewayRouteTableID, transitGatewayAttachmentID); err != nil {
			return fmt.Errorf("waiting for EC2 Transit Gateway Route Table Association (%s) create: %w", id, err)
		}

		input := &ec2.DisassociateTransitGatewayRouteTableInput{
			TransitGatewayAttachmentId: aws.String(transitGatewayAttachmentID),
			TransitGatewayRouteTableId: aws.String(transitGatewayRouteTableID),
		}

		if _, err := conn.DisassociateTransitGatewayRouteTable(input); err != nil {
			return fmt.Errorf("deleting EC2 Transit Gateway Route Table Association (%s): %w", id, err)
		}

		if _, err := WaitTransitGatewayRouteTableAssociationDeleted(conn, transitGatewayRouteTableID, transitGatewayAttachmentID); err != nil {
			return fmt.Errorf("waiting for EC2 Transit Gateway Route Table Association (%s) delete: %w", id, err)
		}
	}

	return nil
}

const transitGatewayRouteTableAssociationIDSeparator = "_"

func TransitGatewayRouteTableAssociationCreateResourceID(transitGatewayRouteTableID, transitGatewayAttachmentID string) string {
	parts := []string{transitGatewayRouteTableID, transitGatewayAttachmentID}
	id := strings.Join(parts, transitGatewayRouteTableAssociationIDSeparator)

	return id
}

func TransitGatewayRouteTableAssociationParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, transitGatewayRouteTableAssociationIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected TRANSIT-GATEWAY-ROUTE-TABLE-ID%[2]sTRANSIT-GATEWAY-ATTACHMENT-ID", id, transitGatewayRouteTableAssociationIDSeparator)
}
