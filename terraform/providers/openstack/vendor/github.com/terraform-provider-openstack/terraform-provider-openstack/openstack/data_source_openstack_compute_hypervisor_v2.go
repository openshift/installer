package openstack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
)

func dataSourceComputeHypervisorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeHypervisorV2Read,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},

			"host_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeHypervisorV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	allPages, err := hypervisors.List(computeClient, hypervisors.ListOpts{}).AllPages()
	if err != nil {
		return diag.Errorf("Error listing compute hypervisors: %s", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		return diag.Errorf("Error extracting compute hypervisors: %s", err)
	}

	name := d.Get("hostname").(string)

	var refinedHypervisors []hypervisors.Hypervisor
	for _, hypervisor := range allHypervisors {
		if hypervisor.HypervisorHostname == name {
			refinedHypervisors = append(refinedHypervisors, hypervisor)
		}
	}

	if len(refinedHypervisors) < 1 {
		return diag.Errorf("Could not find any hypervisor with this name: %s", name)
	}
	if len(refinedHypervisors) > 1 {
		return diag.Errorf("More than one hypervisor found with this name: %s", name)
	}

	h := refinedHypervisors[0]

	d.SetId(h.ID)
	d.Set("hostname", h.HypervisorHostname)
	d.Set("host_ip", h.HostIP)
	d.Set("state", h.State)
	d.Set("status", h.Status)
	d.Set("type", h.HypervisorType)

	d.Set("vcpus", h.VCPUs)
	d.Set("memory", h.MemoryMB)
	d.Set("disk", h.LocalGB)

	return nil
}
