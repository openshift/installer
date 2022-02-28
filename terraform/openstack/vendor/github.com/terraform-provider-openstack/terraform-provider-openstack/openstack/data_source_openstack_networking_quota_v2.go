package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
)

func dataSourceNetworkingQuotaV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingQuotaV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type: schema.TypeString,

				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"floatingip": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"rbac_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"router": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"security_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"security_group_rule": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"subnet": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"subnetpool": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingQuotaV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	projectID := d.Get("project_id").(string)

	q, err := quotas.Get(networkingClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_quota_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_quota_v2 %s: %#v", d.Id(), q)

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)
	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("floatingip", q.FloatingIP)
	d.Set("network", q.Network)
	d.Set("port", q.Port)
	d.Set("rbac_policy", q.RBACPolicy)
	d.Set("router", q.Router)
	d.Set("security_group", q.SecurityGroup)
	d.Set("security_group_rule", q.SecurityGroupRule)
	d.Set("subnet", q.Subnet)
	d.Set("subnetpool", q.SubnetPool)

	return nil
}
