package ec2

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceCoIPPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCoIPPoolRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"local_gateway_route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pool_cidrs": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tftags.TagsSchemaComputed(),

			"filter": CustomFiltersSchema(),
		},
	}
}

func dataSourceCoIPPoolRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	req := &ec2.DescribeCoipPoolsInput{}

	if v, ok := d.GetOk("pool_id"); ok {
		req.PoolIds = []*string{aws.String(v.(string))}
	}

	filters := map[string]string{}

	if v, ok := d.GetOk("local_gateway_route_table_id"); ok {
		filters["coip-pool.local-gateway-route-table-id"] = v.(string)
	}

	req.Filters = BuildAttributeFilterList(filters)

	if tags, tagsOk := d.GetOk("tags"); tagsOk {
		req.Filters = append(req.Filters, BuildTagFilterList(
			Tags(tftags.New(tags.(map[string]interface{}))),
		)...)
	}

	req.Filters = append(req.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	log.Printf("[DEBUG] Reading AWS COIP Pool: %s", req)
	resp, err := conn.DescribeCoipPools(req)
	if err != nil {
		return fmt.Errorf("describing EC2 COIP Pools: %w", err)
	}
	if resp == nil || len(resp.CoipPools) == 0 {
		return fmt.Errorf("no matching COIP Pool found")
	}
	if len(resp.CoipPools) > 1 {
		return fmt.Errorf("multiple Coip Pools matched; use additional constraints to reduce matches to a single COIP Pool")
	}

	coip := resp.CoipPools[0]

	d.SetId(aws.StringValue(coip.PoolId))

	d.Set("local_gateway_route_table_id", coip.LocalGatewayRouteTableId)
	d.Set("arn", coip.PoolArn)

	if err := d.Set("pool_cidrs", aws.StringValueSlice(coip.PoolCidrs)); err != nil {
		return fmt.Errorf("setting pool_cidrs: %w", err)
	}

	d.Set("pool_id", coip.PoolId)

	if err := d.Set("tags", KeyValueTags(coip.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	return nil
}
