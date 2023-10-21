package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_ec2_availability_zone_group")
func ResourceAvailabilityZoneGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAvailabilityZoneGroupCreate,
		ReadWithoutTimeout:   resourceAvailabilityZoneGroupRead,
		UpdateWithoutTimeout: resourceAvailabilityZoneGroupUpdate,
		DeleteWithoutTimeout: schema.NoopContext,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opt_in_status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					ec2.AvailabilityZoneOptInStatusOptedIn,
					ec2.AvailabilityZoneOptInStatusNotOptedIn,
				}, false),
			},
		},
	}
}

func resourceAvailabilityZoneGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	groupName := d.Get("group_name").(string)
	availabilityZone, err := FindAvailabilityZoneGroupByName(ctx, conn, groupName)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 Availability Zone Group (%s): %s", groupName, err)
	}

	if v := d.Get("opt_in_status").(string); v != aws.StringValue(availabilityZone.OptInStatus) {
		if err := modifyAvailabilityZoneOptInStatus(ctx, conn, groupName, v); err != nil {
			return sdkdiag.AppendErrorf(diags, "creating EC2 Availability Zone Group (%s): %s", groupName, err)
		}
	}

	d.SetId(groupName)

	return append(diags, resourceAvailabilityZoneGroupRead(ctx, d, meta)...)
}

func resourceAvailabilityZoneGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	availabilityZone, err := FindAvailabilityZoneGroupByName(ctx, conn, d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Availability Zone Group (%s): %s", d.Id(), err)
	}

	if aws.StringValue(availabilityZone.OptInStatus) == ec2.AvailabilityZoneOptInStatusOptInNotRequired {
		return sdkdiag.AppendErrorf(diags, "unnecessary handling of EC2 Availability Zone Group (%s), status: %s", d.Id(), ec2.AvailabilityZoneOptInStatusOptInNotRequired)
	}

	d.Set("group_name", availabilityZone.GroupName)
	d.Set("opt_in_status", availabilityZone.OptInStatus)

	return diags
}

func resourceAvailabilityZoneGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	if err := modifyAvailabilityZoneOptInStatus(ctx, conn, d.Id(), d.Get("opt_in_status").(string)); err != nil {
		return sdkdiag.AppendErrorf(diags, "updating EC2 Availability Zone Group (%s): %s", d.Id(), err)
	}

	return append(diags, resourceAvailabilityZoneGroupRead(ctx, d, meta)...)
}

func modifyAvailabilityZoneOptInStatus(ctx context.Context, conn *ec2.EC2, groupName, optInStatus string) error {
	input := &ec2.ModifyAvailabilityZoneGroupInput{
		GroupName:   aws.String(groupName),
		OptInStatus: aws.String(optInStatus),
	}

	if _, err := conn.ModifyAvailabilityZoneGroupWithContext(ctx, input); err != nil {
		return err
	}

	waiter := WaitAvailabilityZoneGroupOptedIn
	if optInStatus == ec2.AvailabilityZoneOptInStatusNotOptedIn {
		waiter = WaitAvailabilityZoneGroupNotOptedIn
	}

	if _, err := waiter(ctx, conn, groupName); err != nil {
		return fmt.Errorf("waiting for completion: %w", err)
	}

	return nil
}
