package ec2

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr_block_associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"association_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"dhcp_options_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_dns_hostnames": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_dns_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	// We specify "default" as boolean, but EC2 filters want
	// it to be serialized as a string. Note that setting it to
	// "false" here does not actually filter by it *not* being
	// the default, because Terraform can't distinguish between
	// "false" and "not set".
	isDefaultStr := ""
	if d.Get("default").(bool) {
		isDefaultStr = "true"
	}
	input := &ec2.DescribeVpcsInput{
		Filters: BuildAttributeFilterList(
			map[string]string{
				"cidr":            d.Get("cidr_block").(string),
				"dhcp-options-id": d.Get("dhcp_options_id").(string),
				"isDefault":       isDefaultStr,
				"state":           d.Get("state").(string),
			},
		),
	}

	if v, ok := d.GetOk("id"); ok {
		input.VpcIds = aws.StringSlice([]string{v.(string)})
	}

	if tags, tagsOk := d.GetOk("tags"); tagsOk {
		input.Filters = append(input.Filters, BuildTagFilterList(
			Tags(tftags.New(tags.(map[string]interface{}))),
		)...)
	}

	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	vpc, err := FindVPC(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 VPC", err)
	}

	d.SetId(aws.StringValue(vpc.VpcId))

	ownerID := aws.StringValue(vpc.OwnerId)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: ownerID,
		Resource:  fmt.Sprintf("vpc/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("default", vpc.IsDefault)
	d.Set("dhcp_options_id", vpc.DhcpOptionsId)
	d.Set("instance_tenancy", vpc.InstanceTenancy)
	d.Set("owner_id", ownerID)

	if v, err := FindVPCAttribute(conn, d.Id(), ec2.VpcAttributeNameEnableDnsHostnames); err != nil {
		return fmt.Errorf("error reading EC2 VPC (%s) Attribute (%s): %w", d.Id(), ec2.VpcAttributeNameEnableDnsHostnames, err)
	} else {
		d.Set("enable_dns_hostnames", v)
	}

	if v, err := FindVPCAttribute(conn, d.Id(), ec2.VpcAttributeNameEnableDnsSupport); err != nil {
		return fmt.Errorf("error reading EC2 VPC (%s) Attribute (%s): %w", d.Id(), ec2.VpcAttributeNameEnableDnsSupport, err)
	} else {
		d.Set("enable_dns_support", v)
	}

	if v, err := FindVPCMainRouteTable(conn, d.Id()); err != nil {
		log.Printf("[WARN] Error reading EC2 VPC (%s) main Route Table: %s", d.Id(), err)
		d.Set("main_route_table_id", nil)
	} else {
		d.Set("main_route_table_id", v.RouteTableId)
	}

	cidrAssociations := []interface{}{}
	for _, v := range vpc.CidrBlockAssociationSet {
		association := map[string]interface{}{
			"association_id": aws.StringValue(v.AssociationId),
			"cidr_block":     aws.StringValue(v.CidrBlock),
			"state":          aws.StringValue(v.CidrBlockState.State),
		}
		cidrAssociations = append(cidrAssociations, association)
	}
	if err := d.Set("cidr_block_associations", cidrAssociations); err != nil {
		return fmt.Errorf("error setting cidr_block_associations: %w", err)
	}

	if len(vpc.Ipv6CidrBlockAssociationSet) > 0 {
		d.Set("ipv6_association_id", vpc.Ipv6CidrBlockAssociationSet[0].AssociationId)
		d.Set("ipv6_cidr_block", vpc.Ipv6CidrBlockAssociationSet[0].Ipv6CidrBlock)
	} else {
		d.Set("ipv6_association_id", nil)
		d.Set("ipv6_cidr_block", nil)
	}

	if err := d.Set("tags", KeyValueTags(vpc.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
