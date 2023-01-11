package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceLocalGatewayVirtualInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalGatewayVirtualInterfaceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_bgp_asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_gateway_virtual_interface_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"peer_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_bgp_asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"vlan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLocalGatewayVirtualInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &ec2.DescribeLocalGatewayVirtualInterfacesInput{}

	if v, ok := d.GetOk("id"); ok {
		input.LocalGatewayVirtualInterfaceIds = []*string{aws.String(v.(string))}
	}

	input.Filters = append(input.Filters, BuildTagFilterList(
		Tags(tftags.New(d.Get("tags").(map[string]interface{}))),
	)...)

	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)

	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	output, err := conn.DescribeLocalGatewayVirtualInterfaces(input)

	if err != nil {
		return fmt.Errorf("error describing EC2 Local Gateway Virtual Interfaces: %w", err)
	}

	if output == nil || len(output.LocalGatewayVirtualInterfaces) == 0 {
		return fmt.Errorf("no matching EC2 Local Gateway Virtual Interface found")
	}

	if len(output.LocalGatewayVirtualInterfaces) > 1 {
		return fmt.Errorf("multiple EC2 Local Gateway Virtual Interfaces matched; use additional constraints to reduce matches to a single EC2 Local Gateway Virtual Interface")
	}

	localGatewayVirtualInterface := output.LocalGatewayVirtualInterfaces[0]

	d.SetId(aws.StringValue(localGatewayVirtualInterface.LocalGatewayVirtualInterfaceId))
	d.Set("local_address", localGatewayVirtualInterface.LocalAddress)
	d.Set("local_bgp_asn", localGatewayVirtualInterface.LocalBgpAsn)
	d.Set("local_gateway_id", localGatewayVirtualInterface.LocalGatewayId)
	d.Set("peer_address", localGatewayVirtualInterface.PeerAddress)
	d.Set("peer_bgp_asn", localGatewayVirtualInterface.PeerBgpAsn)

	if err := d.Set("tags", KeyValueTags(localGatewayVirtualInterface.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	d.Set("vlan", localGatewayVirtualInterface.Vlan)

	return nil
}
