// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
)

func dataSourceVSphereComputeClusterHostGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereComputeClusterHostGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the host group in the cluster.",
			},
			"compute_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The managed object ID of the cluster.",
			},
			"host_system_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The managed object IDs of the hosts.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVSphereComputeClusterHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	cluster, name, err := resourceVSphereComputeClusterHostGroupObjects(d, meta)
	if err != nil {
		return fmt.Errorf("cannot locate resource: %s", err)
	}

	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return fmt.Errorf("cannot read cluster properties: %s", err)
	}

	hostSystemIDs := make([]string, len(props.Host))
	for i, host := range props.Host {
		hostSystemIDs[i] = host.Reference().Value
	}

	d.SetId(name)
	if err := d.Set("host_system_ids", hostSystemIDs); err != nil {
		return fmt.Errorf("cannot set host_system_ids: %s", err)
	}

	return nil
}
