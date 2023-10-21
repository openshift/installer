package appmesh

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_appmesh_gateway_route")
func DataSourceGatewayRoute() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceGatewayRouteRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mesh_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mesh_owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec":         dataSourcePropertyFromResourceProperty(resourceGatewayRouteSpecSchema()),
			names.AttrTags: tftags.TagsSchemaComputed(),
			"virtual_gateway_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceGatewayRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AppMeshConn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	gatewayRouteName := d.Get("name").(string)
	gatewayRoute, err := FindGatewayRouteByFourPartKey(ctx, conn, d.Get("mesh_name").(string), d.Get("mesh_owner").(string), d.Get("virtual_gateway_name").(string), gatewayRouteName)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading App Mesh Gateway Route (%s): %s", gatewayRouteName, err)
	}

	d.SetId(aws.StringValue(gatewayRoute.GatewayRouteName))
	arn := aws.StringValue(gatewayRoute.Metadata.Arn)
	d.Set("arn", arn)
	d.Set("created_date", gatewayRoute.Metadata.CreatedAt.Format(time.RFC3339))
	d.Set("last_updated_date", gatewayRoute.Metadata.LastUpdatedAt.Format(time.RFC3339))
	d.Set("mesh_name", gatewayRoute.MeshName)
	meshOwner := aws.StringValue(gatewayRoute.Metadata.MeshOwner)
	d.Set("mesh_owner", meshOwner)
	d.Set("name", gatewayRoute.GatewayRouteName)
	d.Set("resource_owner", gatewayRoute.Metadata.ResourceOwner)
	if err := d.Set("spec", flattenGatewayRouteSpec(gatewayRoute.Spec)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting spec: %s", err)
	}
	d.Set("virtual_gateway_name", gatewayRoute.VirtualGatewayName)

	// https://docs.aws.amazon.com/app-mesh/latest/userguide/sharing.html#sharing-permissions
	// Owners and consumers can list tags and can tag/untag resources in a mesh that the account created.
	// They can't list tags and tag/untag resources in a mesh that aren't created by the account.
	var tags tftags.KeyValueTags

	if meshOwner == meta.(*conns.AWSClient).AccountID {
		tags, err = ListTags(ctx, conn, arn)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing tags for App Mesh Gateway Route (%s): %s", arn, err)
		}
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}
