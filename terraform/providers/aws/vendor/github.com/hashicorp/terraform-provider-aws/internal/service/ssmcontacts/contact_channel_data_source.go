package ssmcontacts

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_ssmcontacts_contact_channel")
func DataSourceContactChannel() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceContactChannelRead,

		Schema: map[string]*schema.Schema{
			"activation_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"simple_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	DSNameContactChannel = "Contact Channel Data Source"
)

func dataSourceContactChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	arn := d.Get("arn").(string)

	out, err := findContactChannelByID(ctx, conn, arn)
	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionReading, DSNameContactChannel, arn, err)
	}

	d.SetId(aws.ToString(out.ContactChannelArn))

	if err := setContactChannelResourceData(d, out); err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionSetting, ResNameContactChannel, d.Id(), err)
	}

	return nil
}
