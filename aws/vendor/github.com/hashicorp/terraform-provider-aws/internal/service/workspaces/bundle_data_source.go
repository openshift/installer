package workspaces

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_workspaces_bundle")
func DataSourceBundle() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceWorkspaceBundleRead,

		Schema: map[string]*schema.Schema{
			"bundle_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"owner", "name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bundle_id"},
			},
			"owner": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bundle_id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compute_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"user_storage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"root_storage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWorkspaceBundleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WorkSpacesConn(ctx)

	var bundle *workspaces.WorkspaceBundle

	if bundleID, ok := d.GetOk("bundle_id"); ok {
		resp, err := conn.DescribeWorkspaceBundlesWithContext(ctx, &workspaces.DescribeWorkspaceBundlesInput{
			BundleIds: []*string{aws.String(bundleID.(string))},
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading WorkSpaces Workspace Bundle (%s): %s", bundleID, err)
		}

		if len(resp.Bundles) != 1 {
			return sdkdiag.AppendErrorf(diags, "expected 1 result for WorkSpaces Workspace Bundle %q, found %d", bundleID, len(resp.Bundles))
		}

		bundle = resp.Bundles[0]

		if bundle == nil {
			return sdkdiag.AppendErrorf(diags, "no WorkSpaces Workspace Bundle with ID %q found", bundleID)
		}
	}

	if name, ok := d.GetOk("name"); ok {
		id := name
		input := &workspaces.DescribeWorkspaceBundlesInput{}

		if owner, ok := d.GetOk("owner"); ok {
			id = fmt.Sprintf("%s:%s", owner, id)
			input.Owner = aws.String(owner.(string))
		}

		name := name.(string)
		err := conn.DescribeWorkspaceBundlesPagesWithContext(ctx, input, func(out *workspaces.DescribeWorkspaceBundlesOutput, lastPage bool) bool {
			for _, b := range out.Bundles {
				if aws.StringValue(b.Name) == name {
					bundle = b
					return true
				}
			}

			return !lastPage
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading WorkSpaces Workspace Bundle (%s): %s", id, err)
		}

		if bundle == nil {
			return sdkdiag.AppendErrorf(diags, "no WorkSpaces Workspace Bundle with name %q found", name)
		}
	}

	d.SetId(aws.StringValue(bundle.BundleId))
	d.Set("bundle_id", bundle.BundleId)
	d.Set("description", bundle.Description)
	d.Set("name", bundle.Name)
	d.Set("owner", bundle.Owner)

	computeType := make([]map[string]interface{}, 1)
	if bundle.ComputeType != nil {
		computeType[0] = map[string]interface{}{
			"name": aws.StringValue(bundle.ComputeType.Name),
		}
	}
	if err := d.Set("compute_type", computeType); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting compute_type: %s", err)
	}

	rootStorage := make([]map[string]interface{}, 1)
	if bundle.RootStorage != nil {
		rootStorage[0] = map[string]interface{}{
			"capacity": aws.StringValue(bundle.RootStorage.Capacity),
		}
	}
	if err := d.Set("root_storage", rootStorage); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting root_storage: %s", err)
	}

	userStorage := make([]map[string]interface{}, 1)
	if bundle.UserStorage != nil {
		userStorage[0] = map[string]interface{}{
			"capacity": aws.StringValue(bundle.UserStorage.Capacity),
		}
	}
	if err := d.Set("user_storage", userStorage); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting user_storage: %s", err)
	}

	return diags
}
