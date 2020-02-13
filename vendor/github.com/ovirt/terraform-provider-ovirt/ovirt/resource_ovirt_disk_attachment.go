// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtDiskAttachmentCreate,
		Read:   resourceOvirtDiskAttachmentRead,
		Update: resourceOvirtDiskAttachmentUpdate,
		Delete: resourceOvirtDiskAttachmentDelete,
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
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"bootable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// TODO: Add support for logical_name

			"pass_discard": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"use_scsi_reservation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func resourceOvirtDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	diskID := d.Get("disk_id").(string)
	diskService := conn.SystemService().DisksService().DiskService(diskID)

	// Retrieve the Disk entity
	getDiskResp, err := diskService.Get().Send()
	if err != nil {
		log.Printf("[DEBUG] Failed to get Disk (%s): %s", diskID, err)
		return err
	}

	// Build disk-attachment request
	attachmentBuilder := ovirtsdk4.NewDiskAttachmentBuilder().
		Disk(getDiskResp.MustDisk()).
		Interface(ovirtsdk4.DiskInterface(d.Get("interface").(string))).
		Bootable(d.Get("bootable").(bool)).
		Active(d.Get("active").(bool)).
		ReadOnly(d.Get("read_only").(bool)).
		UsesScsiReservation(d.Get("use_scsi_reservation").(bool))

	if passDiscard, ok := d.GetOkExists("pass_discard"); ok {
		attachmentBuilder.PassDiscard(passDiscard.(bool))
	}
	attachment, err := attachmentBuilder.Build()
	if err != nil {
		return err
	}

	vmID := d.Get("vm_id").(string)
	_, err = conn.SystemService().
		VmsService().
		VmService(vmID).
		DiskAttachmentsService().
		Add().
		Attachment(attachment).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Failed to attach Disk (%s) to VM (%s): %s", diskID, vmID, err)
		return err
	}

	d.SetId(vmID + ":" + diskID)

	return resourceOvirtDiskAttachmentRead(d, meta)
}

func resourceOvirtDiskAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	parts, err := parseResourceID(d.Id(), 2)
	if err != nil {
		return err
	}
	vmID, diskID := parts[0], parts[1]

	paramDiskAttachment := ovirtsdk4.NewDiskAttachmentBuilder()
	attributeUpdate := false

	if d.HasChange("active") {
		paramDiskAttachment.Active(d.Get("active").(bool))
		attributeUpdate = true
	}

	if d.HasChange("bootable") {
		paramDiskAttachment.Bootable(d.Get("bootable").(bool))
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := conn.SystemService().
			VmsService().
			VmService(vmID).
			DiskAttachmentsService().
			AttachmentService(diskID).
			Update().
			DiskAttachment(paramDiskAttachment.MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Failed to update DiskAttachment (%s): %s", d.Id(), err)
			return err
		}
	}
	return resourceOvirtDiskAttachmentRead(d, meta)
}

func resourceOvirtDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	// Disk ID is equals to its relevant Disk Attachment ID
	// Sess: https://github.com/oVirt/ovirt-engine/blob/68753f46f09419ddcdbb632453501273697d1a20/backend/manager/modules/restapi/types/src/main/java/org/ovirt/engine/api/restapi/types/DiskAttachmentMapper.java
	parts, err := parseResourceID(d.Id(), 2)
	if err != nil {
		return err
	}
	vmID, diskID := parts[0], parts[1]
	d.Set("vm_id", vmID)
	d.Set("disk_id", diskID)

	attachmentService := conn.SystemService().
		VmsService().
		VmService(vmID).
		DiskAttachmentsService().
		AttachmentService(diskID)
	attachmentResp, err := attachmentService.Get().Send()
	if err != nil {
		return err
	}
	attachment, ok := attachmentResp.Attachment()
	if !ok {
		d.SetId("")
		return nil
	}

	d.Set("active", attachment.MustActive())
	d.Set("bootable", attachment.MustBootable())
	d.Set("interface", attachment.MustInterface())
	d.Set("read_only", attachment.MustReadOnly())
	d.Set("use_scsi_reservation", attachment.MustUsesScsiReservation())

	if passDiscard, ok := attachment.PassDiscard(); ok {
		d.Set("pass_discard", passDiscard)
	}

	return nil
}

func resourceOvirtDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	parts, err := parseResourceID(d.Id(), 2)
	if err != nil {
		return err
	}
	vmID, diskID := parts[0], parts[1]

	diskAttachmentService := conn.SystemService().
		VmsService().
		VmService(vmID).
		DiskAttachmentsService().
		AttachmentService(diskID)

	// Ensure the disk attachment is not active
	getResp, err := diskAttachmentService.Get().Send()
	if err != nil {
		log.Printf("[DEBUG] Failed to get DiskAttachment (%s): %s", d.Id(), err)
		return err
	}
	// If it's active, deactivate first
	if getResp.MustAttachment().MustActive() {
		log.Printf("[DEBUG] Now to deactivate DiskAttachment (%s)", d.Id())
		_, err = diskAttachmentService.Update().
			DiskAttachment(
				ovirtsdk4.NewDiskAttachmentBuilder().
					Active(false).
					MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Failed to deactivate DiskAttachment (%s): %s", d.Id(), err)
			return err
		}

		log.Printf("[DEBUG] DiskAttachment (%s) is updated and wait for inactive", d.Id())
		inactiveStateConf := &resource.StateChangeConf{
			Target:     []string{"inactive"},
			Refresh:    DiskAttachmentStateRefreshFunc(conn, vmID, diskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = inactiveStateConf.WaitForState()
		if err != nil {
			log.Printf("[DEBUG] Failed to wait for DiskAttachment (%s) to become inactive: %s", d.Id(), err)
			return err
		}
		log.Printf("[DEBUG] Disk (%s) has became inactive", d.Id())
	}

	_, err = diskAttachmentService.Remove().Send()
	if err != nil {
		log.Printf("[DEBUG] Failed to remove DiskAttachment (%s): %s", d.Id(), err)
		return err
	}
	return nil
}

// DiskAttachmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt DiskAttachment.
func DiskAttachmentStateRefreshFunc(conn *ovirtsdk4.Connection, vmID, diskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			VmsService().
			VmService(vmID).
			DiskAttachmentsService().
			AttachmentService(diskID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// Sometimes oVirt has consistency issues and doesn't see
				// newly created Disk instance. Return an empty state.
				return nil, "", nil
			}
			return nil, "", err
		}

		attachmentState := "inactive"
		if r.MustAttachment().MustActive() {
			attachmentState = "active"
		}

		return r.MustAttachment(), attachmentState, nil
	}
}
