package connect

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// @SDKDataSource("aws_connect_user_hierarchy_structure")
func DataSourceUserHierarchyStructure() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceUserHierarchyStructureRead,
		Schema: map[string]*schema.Schema{
			"hierarchy_structure": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level_one": func() *schema.Schema {
							schema := userHierarchyLevelDataSourceSchema()
							return schema
						}(),
						"level_two": func() *schema.Schema {
							schema := userHierarchyLevelDataSourceSchema()
							return schema
						}(),
						"level_three": func() *schema.Schema {
							schema := userHierarchyLevelDataSourceSchema()
							return schema
						}(),
						"level_four": func() *schema.Schema {
							schema := userHierarchyLevelDataSourceSchema()
							return schema
						}(),
						"level_five": func() *schema.Schema {
							schema := userHierarchyLevelDataSourceSchema()
							return schema
						}(),
					},
				},
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
	}
}

// Each level shares the same schema
func userHierarchyLevelDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"arn": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceUserHierarchyStructureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn(ctx)

	instanceID := d.Get("instance_id").(string)

	resp, err := conn.DescribeUserHierarchyStructureWithContext(ctx, &connect.DescribeUserHierarchyStructureInput{
		InstanceId: aws.String(instanceID),
	})

	if err != nil {
		return diag.Errorf("getting Connect User Hierarchy Structure for Connect Instance (%s): %s", instanceID, err)
	}

	if resp == nil || resp.HierarchyStructure == nil {
		return diag.Errorf("getting Connect User Hierarchy Structure for Connect Instance (%s): empty response", instanceID)
	}

	if err := d.Set("hierarchy_structure", flattenUserHierarchyStructure(resp.HierarchyStructure)); err != nil {
		return diag.Errorf("setting Connect User Hierarchy Structure for Connect Instance: (%s)", instanceID)
	}

	d.SetId(instanceID)

	return nil
}
