package ovirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var diskAttachmentsSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Meta-identifier for the disk attachments. Will always be the same as the VM ID after apply.",
	},
	"attachment": {
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"disk_id": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "ID of the disk to attach. This disk must either be shared or not yet attached to a different VM.",
					ForceNew:         false,
					ValidateDiagFunc: validateUUID,
				},
				"disk_interface": {
					Type:     schema.TypeString,
					Required: true,
					Description: fmt.Sprintf(
						"Type of interface to use for attaching disk. One of: `%s`.",
						strings.Join(ovirtclient.DiskInterfaceValues().Strings(), "`, `"),
					),
					ForceNew:         false,
					ValidateDiagFunc: validateDiskInterface,
				},
				"bootable": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: "Defines whether the disk is bootable.",
				},
				"active": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: "Defines whether the disk is active in the virtual machine it is attached to.",
				},
			},
		},
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the VM the disks should be attached to.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"detach_unmanaged": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"remove_unmanaged"},
		Description:   `Detach unmanaged disks from the VM. This is useful for detaching disks that have been inherited from the template or added manually. The detached disks will not be removed and can be used. To remove the disks instead, use ` + "`remove_unmanaged`.",
	},
	"remove_unmanaged": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"detach_unmanaged"},
		Description: `Completely remove attached disks that are not listed in this resources. This is useful for removing disks that have been inherited from the template or added manually.

~> Use with care! This option will delete all disks attached to the current VM that are not managed, not just detach them!`,
	},
}

func (p *provider) diskAttachmentsResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.diskAttachmentsCreateOrUpdate,
		ReadContext:   p.diskAttachmentsRead,
		UpdateContext: p.diskAttachmentsCreateOrUpdate,
		DeleteContext: p.diskAttachmentsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.diskAttachmentsImport,
		},
		Schema: diskAttachmentsSchema,
		Description: `The ovirt_disk_attachments resource attaches multiple disks to a single VM in one operation. It also allows for removing all attachments that are not declared in an attachment block. This is useful for removing attachments that have been added from the template.

~> Do not use this resource on the same VM as ovirt_disk_attachment (singular). It will cause a ping-pong effect of resources being created and removed on each Terraform run.
`,
	}
}

func (p *provider) diskAttachmentsCreateOrUpdate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	desiredAttachments := data.Get("attachment").(*schema.Set)
	detachUnmanaged := data.Get("detach_unmanaged").(bool)
	removeUnmanaged := data.Get("remove_unmanaged").(bool)

	existingAttachments, err := client.ListDiskAttachments(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags("list existing disk attachments", err)
	}

	diags := diag.Diagnostics{}
	for _, desiredAttachmentInterface := range desiredAttachments.List() {
		desiredAttachment := desiredAttachmentInterface.(map[string]interface{})
		diags = append(
			diags, p.createOrUpdateDiskAttachment(
				client,
				existingAttachments,
				desiredAttachment,
				vmID,
			)...,
		)
	}

	if detachUnmanaged || removeUnmanaged {
		diags = p.cleanUnmanagedDiskAttachments(client, removeUnmanaged, existingAttachments, desiredAttachments, diags)
	}
	data.SetId(vmID)
	if err := data.Set("attachment", desiredAttachments); err != nil {
		diags = append(diags, errorToDiag("set attachment in Terraform", err))
	}
	return diags
}

func (p *provider) cleanUnmanagedDiskAttachments(
	client ovirtclient.Client,
	removeUnmanaged bool,
	existingAttachments []ovirtclient.DiskAttachment,
	desiredAttachments *schema.Set,
	diags diag.Diagnostics,
) diag.Diagnostics {
	for _, attachment := range existingAttachments {
		found := false
		for _, desiredAttachmentInterface := range desiredAttachments.List() {
			desiredAttachment := desiredAttachmentInterface.(map[string]interface{})
			if desiredAttachment["id"] == attachment.ID() {
				found = true
				break
			}
		}
		if !found {
			if err := attachment.Remove(); err != nil {
				diags = append(
					diags,
					errorToDiag(
						fmt.Sprintf("remove disk attachmend %s", attachment.ID()),
						err,
					),
				)
			} else if removeUnmanaged {
				if err := client.RemoveDisk(attachment.DiskID()); err != nil {
					diags = append(
						diags,
						errorToDiag(
							fmt.Sprintf("remove disk %s, please remove manually", attachment.DiskID()),
							err,
						),
					)
				}
			}
		}
	}
	return diags
}

// createOrUpdateDiskAttachment creates or updates a single disk attachment. If the ID is set it will attempt to find
// the attachment. If none is found, or the ID is not set, it will create the attachment. If the disk interface type
// is mismatched, the attachment will be recreated with the correct type.
func (p *provider) createOrUpdateDiskAttachment(
	client ovirtclient.Client,
	existingAttachments []ovirtclient.DiskAttachment,
	desiredAttachment map[string]interface{},
	vmID string,
) diag.Diagnostics {
	id := desiredAttachment["id"].(string)
	diskID := desiredAttachment["disk_id"].(string)
	diskInterfaceName := desiredAttachment["disk_interface"].(string)
	bootable := desiredAttachment["bootable"].(bool)
	active := desiredAttachment["active"].(bool)

	var foundExisting ovirtclient.DiskAttachment
	if id != "" {
		// The attachment is known in the Terraform state, let's try and find it.
		for _, existingAttachment := range existingAttachments {
			if existingAttachment.ID() == ovirtclient.DiskAttachmentID(id) {
				foundExisting = existingAttachment
				break
			}
		}
	}
	if foundExisting != nil {
		// If we found an existing attachment, check if all parameters match. Otherwise, remove the attachment
		// and let it be re-created below.
		if string(foundExisting.DiskID()) == diskID &&
			string(foundExisting.DiskInterface()) == diskInterfaceName &&
			foundExisting.Bootable() == bootable &&
			foundExisting.Active() == active {
			return nil
		}
		if err := foundExisting.Remove(); err != nil && !isNotFound(err) {
			return errorToDiags(
				fmt.Sprintf("remove existing disk interface %s", foundExisting.ID()),
				err,
			)
		}
		// Set the state to be empty. If the API call below fails, Terraform knows it will need to re-create the
		// attachment.
		desiredAttachment["id"] = ""
		desiredAttachment["disk_id"] = ""
		desiredAttachment["disk_interface"] = ""
		desiredAttachment["bootable"] = false
		desiredAttachment["active"] = false
	}

	var err error
	params := ovirtclient.CreateDiskAttachmentParams()
	params, err = params.WithBootable(bootable)
	if err != nil {
		return errorToDiags("set bootable flag for disk attachment.", err)
	}
	params, err = params.WithActive(active)
	if err != nil {
		return errorToDiags("set active flag for disk attachment.", err)
	}

	// Create or re-create disk attachment, then set it in the Terraform state.
	attachment, err := client.CreateDiskAttachment(
		ovirtclient.VMID(vmID),
		ovirtclient.DiskID(diskID),
		ovirtclient.DiskInterface(diskInterfaceName),
		params,
	)
	if err != nil {
		return errorToDiags(
			fmt.Sprintf("remove existing disk interface %s", foundExisting.ID()),
			err,
		)
	}
	desiredAttachment["id"] = attachment.ID()
	desiredAttachment["disk_id"] = attachment.DiskID()
	desiredAttachment["disk_interface"] = string(attachment.DiskInterface())
	desiredAttachment["bootable"] = attachment.Bootable()
	desiredAttachment["active"] = attachment.Active()

	return nil
}

func (p *provider) diskAttachmentsRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	vmID := data.Get("vm_id").(string)
	diskAttachments, err := client.ListDiskAttachments(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("listing disk attachments of VM %s", vmID), err)
	}

	// Go over the list of attachments and try to link them up. If not all attachments are found, Terraform will try to
	// create the missing attachments.
	attachments := data.Get("attachment").(*schema.Set)
	for _, attachmentInterface := range attachments.List() {
		attachment := attachmentInterface.(map[string]interface{})
		attachments.Remove(attachment)
		found := false
		for _, diskAttachment := range diskAttachments {
			if attachment["id"] == diskAttachment.ID() {
				found = true
				attachment["disk_id"] = diskAttachment.DiskID()
				attachment["disk_interface"] = string(diskAttachment.DiskInterface())
			}
		}
		if found {
			attachments.Add(attachment)
		}
	}

	// Go over the existing attachments. If any unmanaged attachments are found, detach_unmanaged and remove_unmanaged
	// are explicitly set to false. This will cause Terraform to run the update again and try to detach/remove the disk
	// again.
	for _, diskAttachment := range diskAttachments {
		found := false
		for _, attachmentInterface := range attachments.List() {
			attachment := attachmentInterface.(map[string]interface{})
			if attachment["id"] == diskAttachment.ID() {
				found = true
				break
			}
		}
		if !found {
			if err := data.Set("detach_unmanaged", false); err != nil {
				return errorToDiags("setting detach_unmanaged", err)
			}
			if err := data.Set("remove_unmanaged", false); err != nil {
				return errorToDiags("setting remove_unmanaged", err)
			}
		}
	}
	return nil
}

func (p *provider) diskAttachmentsDelete(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	diags := diag.Diagnostics{}
	vmID := data.Get("vm_id").(string)
	attachments := data.Get("attachment").(*schema.Set)
	for _, attachmentInterface := range attachments.List() {
		attachment := attachmentInterface.(map[string]interface{})
		if err := client.RemoveDiskAttachment(
			ovirtclient.VMID(vmID),
			ovirtclient.DiskAttachmentID(attachment["id"].(string)),
		); err != nil {
			if !isNotFound(err) {
				diags = append(diags, errorToDiag("remove disk attachment", err))
			} else {
				attachments.Remove(attachment)
			}
		} else {
			attachments.Remove(attachment)
		}
	}
	if err := data.Set("attachment", attachments); err != nil {
		diags = append(diags, errorToDiag("set attachment", err))
	}
	if !diags.HasError() {
		data.SetId("")
	}
	return diags
}

func (p *provider) diskAttachmentsImport(
	ctx context.Context,
	data *schema.ResourceData,
	i interface{},
) ([]*schema.ResourceData, error) {
	vmID := data.Id()
	if err := data.Set("vm_id", vmID); err != nil {
		return nil, err
	}

	diags := p.diskAttachmentsRead(ctx, data, i)
	if diags.HasError() {
		return []*schema.ResourceData{
			data,
		}, diagsToError(diags)
	}

	return []*schema.ResourceData{
		data,
	}, nil
}
