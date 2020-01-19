package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/baremetalintrospection/v1/introspection"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

// Schema resource for an introspection data source, that has some selected details about the node exposed.
func dataSourceIronicIntrospection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIronicIntrospectionRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"finished": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Description: "A list of interfaces that were discovered",
			},
			"cpu_arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CPU architecture (e.g., x86_64)",
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of CPU's",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory in megabytes",
			},
		},
	}
}

func dataSourceIronicIntrospectionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetInspectorClient()
	if err != nil {
		return err
	}

	uuid := d.Get("uuid").(string)

	status, err := introspection.GetIntrospectionStatus(client, uuid).Extract()
	if err != nil {
		return fmt.Errorf("could not get introspection status: %s", err.Error())
	}

	d.Set("finished", status.Finished)
	d.Set("finished_at", status.FinishedAt)
	d.Set("started_at", status.StartedAt)
	d.Set("error", status.Error)
	d.Set("state", status.State)

	if status.Finished {
		data, err := introspection.GetIntrospectionData(client, uuid).Extract()
		if err != nil {
			return fmt.Errorf("could not get introspection data: %s", err.Error())
		}

		// Network interface data
		var interfaces []map[string]string
		for k, v := range data.AllInterfaces {
			interfaces = append(interfaces, map[string]string{
				"name": k,
				"mac":  v.MAC,
				"ip":   v.IP,
			})
		}
		d.Set("interfaces", interfaces)

		// CPU data
		d.Set("cpu_arch", data.CPUArch)
		d.Set("cpu_count", data.CPUs)

		// Memory info
		d.Set("memory_mb", data.MemoryMB)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
