package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

func (p *provider) affinityGroupDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.affinityGroupDataSourceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "oVirt cluster ID in the Data Center.",
				ValidateDiagFunc: validateUUID,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of the affinity group to look for.",
				ValidateDiagFunc: validateNonEmpty,
			},
		},
		Description: `Search oVirt affinity groups by name.`,
	}
}

func (p *provider) affinityGroupDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	clusterID := data.Get("cluster_id").(string)
	affinityGroupName := data.Get("name").(string)

	foundAffinityGroup, err := client.GetAffinityGroupByName(ovirtclient.ClusterID(clusterID), affinityGroupName)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to get affinity group by name",
				Detail:   err.Error(),
			},
		}
	}

	data.SetId(string(foundAffinityGroup.ID()))
	return nil
}
