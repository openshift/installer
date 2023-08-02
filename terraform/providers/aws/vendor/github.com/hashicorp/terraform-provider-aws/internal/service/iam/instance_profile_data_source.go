package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_iam_instance_profile")
func DataSourceInstanceProfile() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceInstanceProfileRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInstanceProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IAMConn(ctx)

	name := d.Get("name").(string)
	instanceProfile, err := FindInstanceProfileByName(ctx, conn, name)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IAM Instance Profile (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(instanceProfile.InstanceProfileId))
	d.Set("arn", instanceProfile.Arn)
	d.Set("create_date", fmt.Sprintf("%v", instanceProfile.CreateDate))
	d.Set("path", instanceProfile.Path)
	if len(instanceProfile.Roles) > 0 {
		role := instanceProfile.Roles[0]
		d.Set("role_arn", role.Arn)
		d.Set("role_id", role.RoleId)
		d.Set("role_name", role.RoleName)
	}

	return diags
}
