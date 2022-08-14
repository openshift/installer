// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerDedicatedHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerDedicatedHostRead,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the dedicated host",
			},
			"host_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the dedicated host pool the dedicated host is associated with",
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor of the dedicated host",
			},
			"placement_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Describes if the placement on the dedicated host is enabled",
			},
			"life_cycle": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actual_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desired_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_details": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_details_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resources of the dedicated host",
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
						"consumed": {
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
					},
				},
			},
			"workers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The workers of the dedicated host",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_id": {
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
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zone of the dedicated host",
			},
		},
	}
}
func dataSourceIBMContainerDedicatedHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	hostID := d.Get("host_id").(string)
	hostPoolID := d.Get("host_pool_id").(string)
	id := fmt.Sprintf("%s/%s", hostPoolID, hostID)

	if err := getIBMContainerDedicatedHost(id, d, meta); err != nil {
		return diag.Errorf("[ERROR] getIBMContainerDedicatedHost failed: %v", err)
	}

	d.SetId(id)
	return nil
}
