package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceSpotPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotPriceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"filter": DataSourceFiltersSchema(),
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spot_price_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSpotPriceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	now := time.Now()
	input := &ec2.DescribeSpotPriceHistoryInput{
		StartTime: &now,
	}

	if v, ok := d.GetOk("instance_type"); ok {
		instanceType := v.(string)
		input.InstanceTypes = []*string{
			aws.String(instanceType),
		}
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		availabilityZone := v.(string)
		input.AvailabilityZone = aws.String(availabilityZone)
	}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = BuildFiltersDataSource(v.(*schema.Set))
	}

	var foundSpotPrice []*ec2.SpotPrice

	err := conn.DescribeSpotPriceHistoryPages(input, func(output *ec2.DescribeSpotPriceHistoryOutput, lastPage bool) bool {
		foundSpotPrice = append(foundSpotPrice, output.SpotPriceHistory...)
		return true
	})
	if err != nil {
		return fmt.Errorf("error reading EC2 Spot Price History: %w", err)
	}

	if len(foundSpotPrice) == 0 {
		return fmt.Errorf("no EC2 Spot Price History found matching criteria; try different search")
	}

	if len(foundSpotPrice) > 1 {
		return fmt.Errorf("multiple EC2 Spot Price History results found matching criteria; try different search")
	}

	resultSpotPrice := foundSpotPrice[0]

	d.Set("spot_price", resultSpotPrice.SpotPrice)
	d.Set("spot_price_timestamp", (*resultSpotPrice.Timestamp).Format(time.RFC3339))
	d.SetId(meta.(*conns.AWSClient).Region)

	return nil
}
