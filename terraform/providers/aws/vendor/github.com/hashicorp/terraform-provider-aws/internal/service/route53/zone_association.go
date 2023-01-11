package route53

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceZoneAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneAssociationCreate,
		Read:   resourceZoneAssociationRead,
		Delete: resourceZoneAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"owning_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceZoneAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn

	vpcRegion := meta.(*conns.AWSClient).Region
	vpcID := d.Get("vpc_id").(string)
	zoneID := d.Get("zone_id").(string)

	if v, ok := d.GetOk("vpc_region"); ok {
		vpcRegion = v.(string)
	}

	input := &route53.AssociateVPCWithHostedZoneInput{
		HostedZoneId: aws.String(zoneID),
		VPC: &route53.VPC{
			VPCId:     aws.String(vpcID),
			VPCRegion: aws.String(vpcRegion),
		},
		Comment: aws.String("Managed by Terraform"),
	}

	output, err := conn.AssociateVPCWithHostedZone(input)

	if err != nil {
		return fmt.Errorf("associating Route 53 Hosted Zone (%s) to EC2 VPC (%s): %w", zoneID, vpcID, err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", zoneID, vpcID, vpcRegion))

	if output != nil && output.ChangeInfo != nil && output.ChangeInfo.Id != nil {
		wait := resource.StateChangeConf{
			Delay:      30 * time.Second,
			Pending:    []string{route53.ChangeStatusPending},
			Target:     []string{route53.ChangeStatusInsync},
			Timeout:    10 * time.Minute,
			MinTimeout: 2 * time.Second,
			Refresh:    resourceZoneAssociationRefreshFunc(conn, CleanChangeID(aws.StringValue(output.ChangeInfo.Id)), d.Id()),
		}

		if _, err := wait.WaitForState(); err != nil {
			return fmt.Errorf("waiting for Route 53 Zone Association (%s) synchronization: %w", d.Id(), err)
		}
	}

	return resourceZoneAssociationRead(d, meta)
}

func resourceZoneAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn

	zoneID, vpcID, vpcRegion, err := ZoneAssociationParseID(d.Id())

	if err != nil {
		return err
	}

	// Continue supporting older resources without VPC Region in ID
	if vpcRegion == "" {
		vpcRegion = d.Get("vpc_region").(string)
	}

	if vpcRegion == "" {
		vpcRegion = meta.(*conns.AWSClient).Region
	}

	hostedZoneSummary, err := GetZoneAssociation(conn, zoneID, vpcID, vpcRegion)

	if tfawserr.ErrMessageContains(err, "AccessDenied", "is not owned by you") && !d.IsNewResource() {
		log.Printf("[WARN] Route 53 Zone Association (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("getting Route 53 Zone Association (%s): %w", d.Id(), err)
	}

	if hostedZoneSummary == nil {
		if d.IsNewResource() {
			return fmt.Errorf("getting Route 53 Zone Association (%s): missing after creation", d.Id())
		}

		log.Printf("[WARN] Route 53 Hosted Zone (%s) Association (%s) not found, removing from state", zoneID, vpcID)
		d.SetId("")
		return nil
	}

	d.Set("vpc_id", vpcID)
	d.Set("vpc_region", vpcRegion)
	d.Set("zone_id", zoneID)
	d.Set("owning_account", hostedZoneSummary.Owner.OwningAccount)

	return nil
}

func resourceZoneAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn

	zoneID, vpcID, vpcRegion, err := ZoneAssociationParseID(d.Id())

	if err != nil {
		return err
	}

	// Continue supporting older resources without VPC Region in ID
	if vpcRegion == "" {
		vpcRegion = d.Get("vpc_region").(string)
	}

	if vpcRegion == "" {
		vpcRegion = meta.(*conns.AWSClient).Region
	}

	input := &route53.DisassociateVPCFromHostedZoneInput{
		HostedZoneId: aws.String(zoneID),
		VPC: &route53.VPC{
			VPCId:     aws.String(vpcID),
			VPCRegion: aws.String(vpcRegion),
		},
		Comment: aws.String("Managed by Terraform"),
	}

	_, err = conn.DisassociateVPCFromHostedZone(input)

	if tfawserr.ErrCodeEquals(err, route53.ErrCodeNoSuchHostedZone) {
		return nil
	}

	if tfawserr.ErrCodeEquals(err, route53.ErrCodeVPCAssociationNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("disassociating Route 53 Hosted Zone (%s) from EC2 VPC (%s): %w", zoneID, vpcID, err)
	}

	return nil
}

func ZoneAssociationParseID(id string) (string, string, string, error) {
	parts := strings.Split(id, ":")

	if len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != "" {
		return parts[0], parts[1], parts[2], nil
	}

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("Unexpected format of ID (%q), expected ZONEID:VPCID or ZONEID:VPCID:VPCREGION", id)
	}

	return parts[0], parts[1], "", nil
}

func resourceZoneAssociationRefreshFunc(conn *route53.Route53, changeId, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		changeRequest := &route53.GetChangeInput{
			Id: aws.String(changeId),
		}
		result, state, err := resourceGoWait(conn, changeRequest)
		if tfawserr.ErrCodeEquals(err, "AccessDenied") {
			log.Printf("[WARN] AccessDenied when trying to get Route 53 change progress for %s - ignoring due to likely cross account issue", id)
			return true, route53.ChangeStatusInsync, nil
		}
		return result, state, err
	}
}

func GetZoneAssociation(conn *route53.Route53, zoneID, vpcID, vpcRegion string) (*route53.HostedZoneSummary, error) {
	input := &route53.ListHostedZonesByVPCInput{
		VPCId:     aws.String(vpcID),
		VPCRegion: aws.String(vpcRegion),
	}

	for {
		output, err := conn.ListHostedZonesByVPC(input)

		if err != nil {
			return nil, err
		}

		var associatedHostedZoneSummary *route53.HostedZoneSummary
		for _, hostedZoneSummary := range output.HostedZoneSummaries {
			if zoneID == aws.StringValue(hostedZoneSummary.HostedZoneId) {
				associatedHostedZoneSummary = hostedZoneSummary
				return associatedHostedZoneSummary, nil
			}
		}
		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}

	return nil, nil
}
