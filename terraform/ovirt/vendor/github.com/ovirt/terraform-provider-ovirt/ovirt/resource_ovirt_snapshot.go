package ovirt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtSnapshotCreate,
		Read:   resourceOvirtSnapshotRead,
		Delete: resourceOvirtSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"save_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOvirtSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	builder := ovirtsdk4.NewSnapshotBuilder()

	builder.Description(d.Get("description").(string)).
		PersistMemorystate(d.Get("save_memory").(bool))

	vmID := d.Get("vm_id").(string)
	snapshotsService := conn.SystemService().VmsService().VmService(vmID).SnapshotsService()

	resp, err := snapshotsService.
		Add().
		Snapshot(builder.MustBuild()).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error creating Snapshot for VM (%s): %s", vmID, err)
		return nil
	}

	snapshotID := resp.MustSnapshot().MustId()
	d.SetId(vmID + ":" + snapshotID)

	// Wait for snapshot is OK
	log.Printf("[DEBUG] Snapshot (%s) is created and wait for ready (status is OK)", d.Id())
	okStateConf := &resource.StateChangeConf{
		Pending:    []string{string(ovirtsdk4.SNAPSHOTSTATUS_LOCKED)},
		Target:     []string{string(ovirtsdk4.SNAPSHOTSTATUS_OK)},
		Refresh:    SnapshotStateRefreshFunc(conn, vmID, snapshotID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = okStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for Snapshot (%s) to become OK: %s", d.Id(), err)
		return err
	}

	return resourceOvirtSnapshotRead(d, meta)
}

func resourceOvirtSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	parts, err := parseResourceID(d.Id(), 2)
	if err != nil {
		return err
	}
	vmID, snapshotID := parts[0], parts[1]

	d.Set("vm_id", vmID)

	snapshotService := conn.SystemService().
		VmsService().
		VmService(vmID).
		SnapshotsService().
		SnapshotService(snapshotID)

	snapshotResp, err := snapshotService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Failed to get Snapshot (%s): %s", d.Id(), err)
		return err
	}
	snapshot := snapshotResp.MustSnapshot()

	d.Set("description", snapshot.MustDescription())
	d.Set("save_memory", snapshot.MustPersistMemorystate())
	d.Set("status", string(snapshot.MustSnapshotStatus()))
	d.Set("type", string(snapshot.MustSnapshotType()))
	d.Set("date", snapshot.MustDate().Format(time.RFC3339))

	return nil
}

func resourceOvirtSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	parts, err := parseResourceID(d.Id(), 2)
	if err != nil {
		return err
	}
	vmID, snapshotID := parts[0], parts[1]

	snapshotService := conn.SystemService().
		VmsService().
		VmService(vmID).
		SnapshotsService().
		SnapshotService(snapshotID)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		log.Printf("[DEBUG] Now to remove Snapshot (%s)", d.Id())
		_, err := snapshotService.Remove().Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Wait until NotFoundError raises
				log.Printf("[DEBUG] Snapshot (%s) has been removed", d.Id())
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error removing Snapshot (%s): %s", d.Id(), err))
		}
		return resource.RetryableError(fmt.Errorf("Snapshot (%s) is still being removed", d.Id()))
	})
}

// SnapshotStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt Snapshot.
func SnapshotStateRefreshFunc(conn *ovirtsdk4.Connection, vmID, snapshotID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			VmsService().
			VmService(vmID).
			SnapshotsService().
			SnapshotService(snapshotID).
			Get().
			Send()

		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				return nil, "", nil
			}
			return nil, "", err
		}

		return r.MustSnapshot, string(r.MustSnapshot().MustSnapshotStatus()), nil
	}
}
