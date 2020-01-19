// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtDiskCreate,
		Read:   resourceOvirtDiskRead,
		Update: resourceOvirtDiskUpdate,
		Delete: resourceOvirtDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ovirtsdk4.DISKFORMAT_COW),
					string(ovirtsdk4.DISKFORMAT_RAW),
				}, false),
			},
			"quota_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "in GiB",
			},
			"shareable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"sparse": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			// "qcow_version" is the only field supporting Disk-Update
			// See: http://ovirt.github.io/ovirt-engine-api-model/4.3/#services/disk/methods/update
		},
	}
}

func resourceOvirtDiskCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	diskBuilder := ovirtsdk4.NewDiskBuilder().
		Name(d.Get("name").(string)).
		Format(ovirtsdk4.DiskFormat(d.Get("format").(string))).
		ProvisionedSize(int64(d.Get("size").(int)) * int64(math.Pow(2, 30))).
		StorageDomainsOfAny(
			ovirtsdk4.NewStorageDomainBuilder().
				Id(d.Get("storage_domain_id").(string)).
				MustBuild())
	if alias, ok := d.GetOk("alias"); ok {
		diskBuilder.Alias(alias.(string))
	}
	if _, ok := d.GetOk("quota_id"); ok {
		diskBuilder.Quota(ovirtsdk4.NewQuotaBuilder().
			Id(d.Get("quota_id").(string)).
			MustBuild())
	}
	if shareable, ok := d.GetOkExists("shareable"); ok {
		diskBuilder.Shareable(shareable.(bool))
	}
	if sparse, ok := d.GetOkExists("sparse"); ok {
		diskBuilder.Sparse(sparse.(bool))
	}
	disk, err := diskBuilder.Build()
	if err != nil {
		return err
	}

	addResp, err := conn.SystemService().DisksService().Add().Disk(disk).Send()
	if err != nil {
		log.Printf("[DEBUG] Error creating the Disk (%s)", d.Get("name").(string))
		return err
	}
	diskID := addResp.MustDisk().MustId()
	d.SetId(diskID)

	// Wait for disk is ready
	log.Printf("[DEBUG] Disk (%s) is created and wait for ready (status is OK)", d.Id())
	okStateConf := &resource.StateChangeConf{
		Pending:    []string{string(ovirtsdk4.DISKSTATUS_LOCKED)},
		Target:     []string{string(ovirtsdk4.DISKSTATUS_OK)},
		Refresh:    DiskStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = okStateConf.WaitForState()
	if err != nil {
		log.Printf("[DEBUG] Failed to wait for Disk (%s) to become OK: %s", d.Id(), err)
		return err
	}

	return resourceOvirtDiskRead(d, meta)
}

func resourceOvirtDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	vml, err := getAttachedVMsOfDisk(d.Id(), meta)
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
	}
	// Disk has not yet attached to any VM
	if len(vml) == 0 {
		return fmt.Errorf("Only the disks attached to VMs can be updated")
	}

	paramDisk := ovirtsdk4.NewDiskBuilder()

	attributeUpdate := false
	if d.HasChange("name") && d.Get("name").(string) != "" {
		paramDisk.Name(d.Get("name").(string))
		attributeUpdate = true
	}
	if d.HasChange("alias") && d.Get("alias").(string) != "" {
		paramDisk.Alias(d.Get("alias").(string))
		attributeUpdate = true
	}
	if d.HasChange("size") {
		oldSizeValue, newSizeValue := d.GetChange("size")
		oldSize := oldSizeValue.(int)
		newSize := newSizeValue.(int)
		if oldSize > newSize {
			return fmt.Errorf("Only size extension is supported")
		}
		paramDisk.ProvisionedSize(int64(newSize) * int64(math.Pow(2, 30)))
		attributeUpdate = true
	}

	if attributeUpdate {
		// Only retrieve the first VM
		vmID := vml[0].MustId()
		attachmentService := conn.SystemService().
			VmsService().
			VmService(vmID).
			DiskAttachmentsService().
			AttachmentService(d.Id())

		_, err := attachmentService.Update().DiskAttachment(
			ovirtsdk4.NewDiskAttachmentBuilder().
				Disk(
					paramDisk.MustBuild()).
				MustBuild()).
			Send()
		if err != nil {
			return err
		}

		// Wait for disk is ready
		log.Printf("[DEBUG] Disk (%s) is updated and wait for ready (status is OK)", d.Id())
		okStateConf := &resource.StateChangeConf{
			Target:     []string{string(ovirtsdk4.DISKSTATUS_OK)},
			Refresh:    DiskStateRefreshFunc(conn, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = okStateConf.WaitForState()
		if err != nil {
			log.Printf("[DEBUG] Failed to wait for Disk (%s) to become OK: %s", d.Id(), err)
			return err
		}
		log.Printf("[DEBUG] Disk (%s) has been successfully updated", d.Id())
	}

	return resourceOvirtDiskRead(d, meta)
}

func resourceOvirtDiskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	getDiskResp, err := conn.SystemService().DisksService().
		DiskService(d.Id()).Get().Send()
	if err != nil {
		return err
	}

	disk, ok := getDiskResp.Disk()
	if !ok {
		d.SetId("")
		return nil
	}
	d.Set("name", disk.MustName())
	d.Set("size", disk.MustProvisionedSize()/int64(math.Pow(2, 30)))
	d.Set("format", disk.MustFormat())

	if sds, ok := disk.StorageDomains(); ok {
		if len(sds.Slice()) > 0 {
			d.Set("storage_domain_id", sds.Slice()[0].MustId())
		}
	}
	if alias, ok := disk.Alias(); ok {
		d.Set("alias", alias)
	}
	if shareable, ok := disk.Shareable(); ok {
		d.Set("shareable", shareable)
	}
	if sparse, ok := disk.Sparse(); ok {
		d.Set("sparse", sparse)
	}

	return nil
}

func resourceOvirtDiskDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	diskService := conn.SystemService().
		DisksService().
		DiskService(d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Now to remove Disk (%s)", d.Id())
		_, e := diskService.Remove().Send()
		if e != nil {
			if _, ok := e.(*ovirtsdk4.NotFoundError); ok {
				log.Printf("[DEBUG] Disk (%s) has been removed", d.Id())
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error removing Disk (%s): %s", d.Id(), e))
		}
		return resource.RetryableError(fmt.Errorf("Disk (%s) is still being removed", d.Id()))
	})
}

func getAttachedVMsOfDisk(diskID string, meta interface{}) ([]*ovirtsdk4.Vm, error) {
	conn := meta.(*ovirtsdk4.Connection)

	diskService := conn.SystemService().DisksService().DiskService(diskID)
	getDiskResp, err := diskService.Get().
		Header("All-Content", "true").
		Send()
	if err != nil {
		return nil, err
	}

	if disk, ok := getDiskResp.Disk(); ok {
		if vmSlice, ok := disk.Vms(); ok {
			return vmSlice.Slice(), nil
		}
	}
	return nil, nil
}

// DiskStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt Disk.
func DiskStateRefreshFunc(conn *ovirtsdk4.Connection, diskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			DisksService().
			DiskService(diskID).
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

		return r.MustDisk(), string(r.MustDisk().MustStatus()), nil
	}
}
