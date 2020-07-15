package packet

import (
	"log"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketVolumeAttachmentCreate,
		Read:   resourcePacketVolumeAttachmentRead,
		Delete: resourcePacketVolumeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePacketVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	vID := d.Get("volume_id").(string)
	log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)\n", vID, dID)
	va, _, err := client.VolumeAttachments.Create(vID, dID)
	if err != nil {
		switch err.(type) {
		case *packngo.ErrorResponse:
			e := err.(*packngo.ErrorResponse)
			if len(e.Errors) == 1 {
				if e.Errors[0] == "Instance is already attached to this volume" {
					log.Printf("[DEBUG] Volume (%s) is already attached to Instance (%s)", vID, dID)
					break
				}
			}
		}
		return err
	}

	d.SetId(va.ID)
	return resourcePacketVolumeAttachmentRead(d, meta)
}

func resourcePacketVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	va, _, err := client.VolumeAttachments.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("device_id", filepath.Base(va.Device.Href))
	d.Set("volume_id", filepath.Base(va.Volume.Href))
	return nil
}

func resourcePacketVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	_, err := client.VolumeAttachments.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}
