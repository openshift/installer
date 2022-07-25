package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func resourceComputeVolumeAttachV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeVolumeAttachV2Create,
		ReadContext:   resourceComputeVolumeAttachV2Read,
		DeleteContext: resourceComputeVolumeAttachV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_volume_confirmation": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceComputeVolumeAttachV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var (
		blockStorageClient       *gophercloud.ServiceClient
		ignoreVolumeConfirmation bool
	)

	// Get vendor_options and decide if BlockStorage V3 client should be initialized.
	vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
	if vendorOptionsRaw.Len() > 0 {
		vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
		ignoreVolumeConfirmation = vendorOptions["ignore_volume_confirmation"].(bool)
	}
	if !ignoreVolumeConfirmation {
		blockStorageClient, err = config.BlockStorageV3Client(GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating OpenStack block storage client: %s", err)
		}
	}

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)

	var device string
	if v, ok := d.GetOk("device"); ok {
		device = v.(string)
	}

	attachOpts := volumeattach.CreateOpts{
		Device:   device,
		VolumeID: volumeID,
	}

	log.Printf("[DEBUG] openstack_compute_volume_attach_v2 attach options %s: %#v", instanceID, attachOpts)

	multiattach := d.Get("multiattach").(bool)
	if multiattach {
		computeClient.Microversion = "2.60"
	}

	var attachment *volumeattach.VolumeAttachment
	timeout := d.Timeout(schema.TimeoutCreate)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		attachment, err = volumeattach.Create(computeClient, instanceID, attachOpts).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault400); ok && multiattach {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating openstack_compute_volume_attach_v2 %s: %s", instanceID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeVolumeAttachV2AttachFunc(computeClient, blockStorageClient, instanceID, attachment.ID, volumeID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error attaching openstack_compute_volume_attach_v2 %s: %s", instanceID, err)
	}

	// Use the instance ID and attachment ID as the resource ID.
	// This is because an attachment cannot be retrieved just by its ID alone.
	id := fmt.Sprintf("%s/%s", instanceID, attachment.ID)

	d.SetId(id)

	return resourceComputeVolumeAttachV2Read(ctx, d, meta)
}

func resourceComputeVolumeAttachV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceID, attachmentID, err := computeVolumeAttachV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_volume_attach_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_volume_attach_v2 %s: %#v", d.Id(), attachment)

	d.Set("instance_id", attachment.ServerID)
	d.Set("volume_id", attachment.VolumeID)
	d.Set("device", attachment.Device)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeVolumeAttachV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceID, attachmentID, err := computeVolumeAttachV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeVolumeAttachV2DetachFunc(computeClient, instanceID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error detaching openstack_compute_volume_attach_v2"))
	}

	return nil
}
