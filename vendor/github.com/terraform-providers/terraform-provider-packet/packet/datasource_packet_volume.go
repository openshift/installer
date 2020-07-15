package packet

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketVolumeRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"volume_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"volume_id"},
			},
			"volume_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id", "name"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"facility": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"billing_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"snapshot_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"device_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	nameRaw, nameOK := d.GetOk("name")
	projectIdRaw, projectIdOK := d.GetOk("project_id")
	volumeIdRaw, volumeIdOK := d.GetOk("volume_id")

	if !volumeIdOK && !nameOK {
		return fmt.Errorf("You must supply volume_id or name")
	}
	var volume *packngo.Volume
	if nameOK {
		if !projectIdOK {
			return fmt.Errorf("If you lookup via name, you must supply project_id")
		}
		name := nameRaw.(string)
		projectId := projectIdRaw.(string)

		vs, _, err := client.Volumes.List(projectId, &packngo.ListOptions{Includes: []string{"attachments.device"}})
		if err != nil {
			return err
		}

		volume, err = findVolumeByName(vs, name)
		if err != nil {
			return err
		}
	} else {
		volumeId := volumeIdRaw.(string)
		var err error
		volume, _, err = client.Volumes.Get(volumeId, &packngo.GetOptions{Includes: []string{"attachments.device"}})
		if err != nil {
			return err
		}
	}

	d.Set("name", volume.Name)
	d.Set("description", volume.Description)
	d.Set("size", volume.Size)
	d.Set("plan", volume.Plan.Slug)
	d.Set("facility", volume.Facility.Code)
	d.Set("state", volume.State)
	d.Set("billing_cycle", volume.BillingCycle)
	d.Set("locked", volume.Locked)
	d.Set("created", volume.Created)
	d.Set("updated", volume.Updated)
	d.Set("project_id", volume.Project.ID)

	snapshot_policies := make([]map[string]interface{}, 0, len(volume.SnapshotPolicies))
	for _, snapshot_policy := range volume.SnapshotPolicies {
		policy := map[string]interface{}{
			"snapshot_frequency": snapshot_policy.SnapshotFrequency,
			"snapshot_count":     snapshot_policy.SnapshotCount,
		}
		snapshot_policies = append(snapshot_policies, policy)
	}
	d.Set("snapshot_policies", snapshot_policies)

	deviceIds := []string{}

	for _, a := range volume.Attachments {
		deviceIds = append(deviceIds, path.Base(a.Device.Href))
	}

	d.Set("device_ids", deviceIds)
	d.SetId(volume.ID)

	return nil
}

func findVolumeByName(volumes []packngo.Volume, name string) (*packngo.Volume, error) {
	results := make([]packngo.Volume, 0)
	for _, v := range volumes {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no volume found with name %s", name)
	}
	return nil, fmt.Errorf("too many volumes found with hostname %s (found %d, expected 1)", name, len(results))
}
