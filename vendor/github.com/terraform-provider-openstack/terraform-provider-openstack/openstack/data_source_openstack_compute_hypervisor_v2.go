package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceComputeHypervisorV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeHypervisorV2Read,
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

func dataSourceComputeHypervisorV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	allPages, err := hypervisors.List(computeClient).AllPages()
	if err != nil {
		return fmt.Errorf("Error listing compute hypervisors: %s", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting compute hypervisors: %s", err)
	}

	name := d.Get("hostname").(string)

	var refinedHypervisors []hypervisors.Hypervisor
	for _, hypervisor := range allHypervisors {
		if hypervisor.HypervisorHostname == name {
			refinedHypervisors = append(refinedHypervisors, hypervisor)
		}
	}

	if len(refinedHypervisors) < 1 {
		return fmt.Errorf("Could not find any hypervisor with this name: %s", name)
	}
	if len(refinedHypervisors) > 1 {
		return fmt.Errorf("More than one hypervisor found with this name: %s", name)
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
