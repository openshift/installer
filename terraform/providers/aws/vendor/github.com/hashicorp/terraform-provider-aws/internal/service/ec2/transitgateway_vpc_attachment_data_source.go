package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceTransitGatewayVPCAttachment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransitGatewayVPCAttachmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"appliance_mode_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_support": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTransitGatewayVPCAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &ec2.DescribeTransitGatewayVpcAttachmentsInput{}

	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)

	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	if v, ok := d.GetOk("id"); ok {
		input.TransitGatewayAttachmentIds = aws.StringSlice([]string{v.(string)})
	}

	transitGatewayVPCAttachment, err := FindTransitGatewayVPCAttachment(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 Transit Gateway VPC Attachment", err)
	}

	d.SetId(aws.StringValue(transitGatewayVPCAttachment.TransitGatewayAttachmentId))
	d.Set("appliance_mode_support", transitGatewayVPCAttachment.Options.ApplianceModeSupport)
	d.Set("dns_support", transitGatewayVPCAttachment.Options.DnsSupport)
	d.Set("ipv6_support", transitGatewayVPCAttachment.Options.Ipv6Support)
	d.Set("subnet_ids", aws.StringValueSlice(transitGatewayVPCAttachment.SubnetIds))
	d.Set("transit_gateway_id", transitGatewayVPCAttachment.TransitGatewayId)
	d.Set("vpc_id", transitGatewayVPCAttachment.VpcId)
	d.Set("vpc_owner_id", transitGatewayVPCAttachment.VpcOwnerId)

	if err := d.Set("tags", KeyValueTags(transitGatewayVPCAttachment.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	return nil
}
