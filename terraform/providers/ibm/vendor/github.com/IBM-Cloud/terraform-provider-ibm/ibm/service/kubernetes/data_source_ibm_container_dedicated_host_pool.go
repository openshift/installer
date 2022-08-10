// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerDedicatedHostPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerDedicatedHostPoolRead,
		Schema: map[string]*schema.Schema{
			"host_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the dedicated host pool",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the dedicated host pool",
			},
			"metro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The metro to create the dedicated host pool in",
			},
			"flavor_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor class of the dedicated host pool",
			},
			"host_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of the hosts under the dedicated host pool",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the dedicated host pool",
			},
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zones of the dedicated host pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"vcpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"host_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The worker pools of the dedicated host pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMContainerDedicatedHostPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hostPoolID := d.Get("host_pool_id").(string)
	err := getIBMContainerDedicatedHostPool(hostPoolID, d, meta)
	if err != nil {
		return diag.Errorf("[ERROR] Error retrieving host pool details %v", err)
	}

	d.SetId(hostPoolID)
	return nil
}
