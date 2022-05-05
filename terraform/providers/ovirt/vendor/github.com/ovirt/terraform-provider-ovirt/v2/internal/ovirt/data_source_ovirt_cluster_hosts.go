package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) clusterHostsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.clusterHostsDataSourceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "oVirt cluster ID in the Data Center.",
				ValidateDiagFunc: validateUUID,
			},
			"hosts": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the host.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of the host.",
						},
					},
				},
			},
		},
		Description: `A set of all hosts of a Cluster.`,
	}
}

func (p *provider) clusterHostsDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID := data.Get("cluster_id").(string)
	allHosts, err := client.ListHosts()

	if err != nil {
		return errorToDiags("list all hosts", err)
	}

	hosts := make([]map[string]interface{}, 0)

	for _, host := range allHosts {
		if string(host.ClusterID()) == clusterID {
			hostMap := make(map[string]interface{}, 0)
			hostMap["id"] = host.ID()
			hostMap["status"] = host.Status()
			hosts = append(hosts, hostMap)
		}
	}

	if err := data.Set("hosts", hosts); err != nil {
		return errorToDiags("set hosts", err)
	}

	data.SetId(clusterID)

	return nil
}
