package memorydb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceSubnetGroup() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSubnetGroupRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).MemoryDBConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get("name").(string)

	group, err := FindSubnetGroupByName(ctx, conn, name)

	if err != nil {
		return diag.FromErr(tfresource.SingularDataSourceFindError("MemoryDB Subnet Group", err))
	}

	d.SetId(aws.StringValue(group.Name))

	var subnetIds []*string
	for _, subnet := range group.Subnets {
		subnetIds = append(subnetIds, subnet.Identifier)
	}

	d.Set("arn", group.ARN)
	d.Set("description", group.Description)
	d.Set("subnet_ids", flex.FlattenStringSet(subnetIds))
	d.Set("name", group.Name)
	d.Set("vpc_id", group.VpcId)

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return diag.Errorf("error listing tags for MemoryDB Subnet Group (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return diag.Errorf("error setting tags: %s", err)
	}

	return nil
}
