package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func resourceBlockStorageVolumeAttachV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageVolumeAttachV3Create,
		ReadContext:   resourceBlockStorageVolumeAttachV3Read,
		DeleteContext: resourceBlockStorageVolumeAttachV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"attach_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ro", "rw",
				}, false),
			},

			"initiator": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"multipath": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"wwpn": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"wwnn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Volume attachment information
			"data": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
			},

			"driver_volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mount_point_base": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBlockStorageVolumeAttachV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	// initialize the connection
	volumeID := d.Get("volume_id").(string)
	connOpts := &volumeactions.InitializeConnectionOpts{}
	if v, ok := d.GetOk("host_name"); ok {
		connOpts.Host = v.(string)
	}

	if v, ok := d.GetOk("multipath"); ok {
		multipath := v.(bool)
		connOpts.Multipath = &multipath
	}

	if v, ok := d.GetOk("ip_address"); ok {
		connOpts.IP = v.(string)
	}

	if v, ok := d.GetOk("initiator"); ok {
		connOpts.Initiator = v.(string)
	}

	if v, ok := d.GetOk("os_type"); ok {
		connOpts.OSType = v.(string)
	}

	if v, ok := d.GetOk("platform"); ok {
		connOpts.Platform = v.(string)
	}

	if v, ok := d.GetOk("wwnns"); ok {
		connOpts.Wwnns = v.(string)
	}

	if v, ok := d.GetOk("wwpns"); ok {
		var wwpns []string
		for _, i := range v.([]string) {
			wwpns = append(wwpns, i)
		}

		connOpts.Wwpns = wwpns
	}

	connInfo, err := volumeactions.InitializeConnection(client, volumeID, connOpts).Extract()
	if err != nil {
		return diag.Errorf(
			"Unable to initialize connection for openstack_blockstorage_volume_attach_v3: %s", err)
	}

	// Only uncomment this when debugging since connInfo contains sensitive information.
	// log.Printf("[DEBUG] Volume Connection for %s: %#v", volumeID, connInfo)

	// Because this information is only returned upon creation,
	// it must be set in Create.
	if v, ok := connInfo["data"]; ok {
		data := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			if v, ok := value.(string); ok {
				data[key] = v
			}
		}

		d.Set("data", data)
	}

	if v, ok := connInfo["driver_volume_type"]; ok {
		d.Set("driver_volume_type", v)
	}

	if v, ok := connInfo["mount_point_base"]; ok {
		d.Set("mount_point_base", v)
	}

	// Once the connection has been made, tell Cinder to mark the volume as attached.
	attachMode, err := expandBlockStorageV3AttachMode(d.Get("attach_mode").(string))
	if err != nil {
		return nil
	}

	attachOpts := &volumeactions.AttachOpts{
		HostName:   d.Get("host_name").(string),
		MountPoint: d.Get("device").(string),
		Mode:       attachMode,
	}

	log.Printf("[DEBUG] openstack_blockstorage_volume_attach_v3 attach options: %#v", attachOpts)

	if err := volumeactions.Attach(client, volumeID, attachOpts).ExtractErr(); err != nil {
		return diag.Errorf(
			"Error attaching openstack_blockstorage_volume_attach_v3 for volume %s: %s", volumeID, err)
	}

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for openstack_blockstorage_volume_attach_v3 volume %s to become available", volumeID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "attaching"},
		Target:     []string{"in-use"},
		Refresh:    blockStorageVolumeV3StateRefreshFunc(client, volumeID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_blockstorage_volume_attach_v3 volume %s to become in-use: %s", volumeID, err)
	}

	// Once the volume has been marked as attached,
	// retrieve a fresh copy of it with all information now available.
	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		return diag.Errorf(
			"Unable to retrieve openstack_blockstorage_volume_attach_v3 volume %s: %s", volumeID, err)
	}

	// Search for the attachmentID
	var attachmentID string
	hostName := d.Get("host_name").(string)
	for _, attachment := range volume.Attachments {
		if hostName != "" && hostName == attachment.HostName {
			attachmentID = attachment.AttachmentID
		}
	}

	if attachmentID == "" {
		return diag.Errorf(
			"Unable to determine attachment ID for openstack_blockstorage_volume_attach_v3 volume %s", volumeID)
	}

	// The ID must be a combination of the volume and attachment ID
	// since a volume ID is required to retrieve an attachment ID.
	id := fmt.Sprintf("%s/%s", volumeID, attachmentID)
	d.SetId(id)

	return resourceBlockStorageVolumeAttachV3Read(ctx, d, meta)
}

func resourceBlockStorageVolumeAttachV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	volumeID, attachmentID, err := blockStorageVolumeAttachV3ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		return diag.Errorf(
			"Unable to retrieve openstack_blockstorage_volume_attach_v3 volume %s: %s", volumeID, err)
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_attach_v3 volume %s: %#v", volumeID, volume)

	var attachment volumes.Attachment
	for _, v := range volume.Attachments {
		if attachmentID == v.AttachmentID {
			attachment = v
		}
	}

	log.Printf(
		"[DEBUG] Retrieved openstack_blockstorage_volume_attach_v3 attachment %s: %#v", d.Id(), attachment)

	return nil
}

func resourceBlockStorageVolumeAttachV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	volumeID, attachmentID, err := blockStorageVolumeAttachV3ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing openstack_blockstorage_volume_attach_v3: %s", err)
	}

	// Terminate the connection
	termOpts := &volumeactions.TerminateConnectionOpts{}
	if v, ok := d.GetOk("host_name"); ok {
		termOpts.Host = v.(string)
	}

	if v, ok := d.GetOk("multipath"); ok {
		multipath := v.(bool)
		termOpts.Multipath = &multipath
	}

	if v, ok := d.GetOk("ip_address"); ok {
		termOpts.IP = v.(string)
	}

	if v, ok := d.GetOk("initiator"); ok {
		termOpts.Initiator = v.(string)
	}

	if v, ok := d.GetOk("os_type"); ok {
		termOpts.OSType = v.(string)
	}

	if v, ok := d.GetOk("platform"); ok {
		termOpts.Platform = v.(string)
	}

	if v, ok := d.GetOk("wwnns"); ok {
		termOpts.Wwnns = v.(string)
	}

	if v, ok := d.GetOk("wwpns"); ok {
		var wwpns []string
		for _, i := range v.([]string) {
			wwpns = append(wwpns, i)
		}

		termOpts.Wwpns = wwpns
	}

	err = volumeactions.TerminateConnection(client, volumeID, termOpts).ExtractErr()
	if err != nil {
		return diag.Errorf(
			"Error terminating openstack_blockstorage_volume_attach_v3 connection %s: %s", d.Id(), err)
	}

	// Detach the volume
	detachOpts := volumeactions.DetachOpts{
		AttachmentID: attachmentID,
	}

	log.Printf(
		"[DEBUG] openstack_blockstorage_volume_attach_v3 detachment options %s: %#v", d.Id(), detachOpts)

	if err := volumeactions.Detach(client, volumeID, detachOpts).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in-use", "attaching", "detaching"},
		Target:     []string{"available"},
		Refresh:    blockStorageVolumeV3StateRefreshFunc(client, volumeID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_blockstorage_volume_attach_v3 volume %s to become available: %s", volumeID, err)
	}

	return nil
}
