package packet

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketVolumeCreate,
		Read:   resourcePacketVolumeRead,
		Update: resourcePacketVolumeUpdate,
		Delete: resourcePacketVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"facility": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"plan": {
				Type:     schema.TypeString,
				Required: true,
			},

			"billing_cycle": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"snapshot_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_frequency": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"snapshot_count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

func resourcePacketVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.VolumeCreateRequest{
		PlanID:     d.Get("plan").(string),
		FacilityID: d.Get("facility").(string),
		Size:       d.Get("size").(int),
		Locked:     d.Get("locked").(bool),
	}

	if attr, ok := d.GetOk("billing_cycle"); ok {
		createRequest.BillingCycle = attr.(string)
	} else {
		createRequest.BillingCycle = "hourly"
	}

	if attr, ok := d.GetOk("description"); ok {
		createRequest.Description = attr.(string)
	}

	snapshot_count := d.Get("snapshot_policies.#").(int)
	if snapshot_count > 0 {
		createRequest.SnapshotPolicies = make([]*packngo.SnapshotPolicy, 0, snapshot_count)
		for i := 0; i < snapshot_count; i++ {
			policy := new(packngo.SnapshotPolicy)
			policy.SnapshotFrequency = d.Get(fmt.Sprintf("snapshot_policies.%d.snapshot_frequency", i)).(string)
			policy.SnapshotCount = d.Get(fmt.Sprintf("snapshot_policies.%d.snapshot_count", i)).(int)
			createRequest.SnapshotPolicies = append(createRequest.SnapshotPolicies, policy)
		}
	}

	newVolume, _, err := client.Volumes.Create(createRequest, d.Get("project_id").(string))
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(newVolume.ID)

	err = waitForVolumeState(newVolume.ID, "active", []string{"queued", "provisioning"}, meta)
	if err != nil {
		d.SetId("")
		return err
	}

	return resourcePacketVolumeRead(d, meta)
}

func waitForVolumeState(volumeID string, target string, pending []string, meta interface{}) error {
	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  []string{target},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*packngo.Client)
			v, _, err := client.Volumes.Get(volumeID, &packngo.GetOptions{Includes: []string{"project", "snapshot_policies", "facility"}})
			if err == nil {
				return 42, v.State, nil
			}
			return 42, "error", err
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	return err
}

func resourcePacketVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	volume, _, err := client.Volumes.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "snapshot_policies", "facility"}})
	if err != nil {
		err = friendlyError(err)

		// If the volume somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
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

	attachments := make([]*packngo.VolumeAttachment, 0, len(volume.Attachments))
	for _, attachment := range volume.Attachments {
		attachments = append(attachments, attachment)
	}
	d.Set("attachments", attachments)

	return nil
}

func resourcePacketVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if d.HasChange("locked") {
		// the change is true => false, i.e. unlock
		if !d.Get("locked").(bool) {
			if _, err := client.Volumes.Unlock(d.Id()); err != nil {
				return friendlyError(err)
			}
		}
	}

	updateRequest := &packngo.VolumeUpdateRequest{}

	sendAttrUpdate := false

	if d.HasChange("description") {
		sendAttrUpdate = true
		vDesc := d.Get("description").(string)
		updateRequest.Description = &vDesc
	}
	if d.HasChange("plan") {
		sendAttrUpdate = true
		vPlan := d.Get("plan").(string)
		updateRequest.PlanID = &vPlan
	}
	if d.HasChange("size") {
		sendAttrUpdate = true
		vSize := d.Get("size").(int)
		updateRequest.Size = &vSize
	}
	if d.HasChange("billing_cycle") {
		sendAttrUpdate = true
		vCycle := d.Get("billing_cycle").(string)
		updateRequest.BillingCycle = &vCycle
	}

	if sendAttrUpdate {
		_, _, err := client.Volumes.Update(d.Id(), updateRequest)
		if err != nil {
			return friendlyError(err)
		}
	}
	if d.HasChange("locked") {
		// the change is false => true, i.e. lock
		if d.Get("locked").(bool) {
			if _, err := client.Volumes.Lock(d.Id()); err != nil {
				return friendlyError(err)
			}
		}
	}

	return resourcePacketVolumeRead(d, meta)
}

func resourcePacketVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if _, err := client.Volumes.Delete(d.Id()); err != nil {
		return friendlyError(err)
	}

	return nil
}
