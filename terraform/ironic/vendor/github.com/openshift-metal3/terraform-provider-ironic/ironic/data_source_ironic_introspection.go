package ironic

import (
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/baremetalintrospection/v1/introspection"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	err = d.Set("finished", status.Finished)
	if err != nil {
		return err
	}
	err = d.Set("finished_at", status.FinishedAt.Format("2006-01-02T15:04:05"))
	if err != nil {
		return err
	}
	err = d.Set("started_at", status.StartedAt.Format("2006-01-02T15:04:05"))
	if err != nil {
		return err
	}
	err = d.Set("error", status.Error)
	if err != nil {
		return err
	}
	err = d.Set("state", status.State)
	if err != nil {
		return err
	}

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
		err = d.Set("interfaces", interfaces)
		if err != nil {
			return err
		}

		// CPU data
		err = d.Set("cpu_arch", data.CPUArch)
		if err != nil {
			return err
		}
		err = d.Set("cpu_count", data.CPUs)
		if err != nil {
			return err
		}

		// Memory info
		err = d.Set("memory_mb", data.MemoryMB)
		if err != nil {
			return err
		}
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
