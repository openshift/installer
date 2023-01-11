package eks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceAddonVersion() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAddonVersionRead,

		Schema: map[string]*schema.Schema{
			"addon_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAddonVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EKSConn

	addonName := d.Get("addon_name").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)
	mostRecent := d.Get("most_recent").(bool)
	id := addonName

	versionInfo, err := FindAddonVersionByAddonNameAndKubernetesVersion(ctx, conn, id, kubernetesVersion, mostRecent)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading EKS Add-On version info (%s, %s): %w", id, kubernetesVersion, err))
	}

	d.SetId(id)

	d.Set("addon_name", addonName)
	d.Set("kubernetes_version", kubernetesVersion)
	d.Set("most_recent", mostRecent)
	d.Set("version", versionInfo.AddonVersion)

	return nil
}
