//nolint:dupl,revive
package ovirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

func (p *provider) diskAttachmentsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.diskAttachmentsDataSourceRead,
		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "oVirt ID of the VM to list the attachments for.",
				ValidateDiagFunc: validateUUID,
			},
			"attachments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the attachement.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the disk in this attachment.",
						},
						"disk_interface": {
							Type:     schema.TypeString,
							Computed: true,
							Description: fmt.Sprintf(
								"Type of interface of the attached disk. One of: `%s`.",
								strings.Join(ovirtclient.DiskInterfaceValues().Strings(), "`, `"),
							),
						},
					},
				},
			},
		},
		Description: `A set of all attachments of a VM.`,
	}
}

func (p *provider) diskAttachmentsDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	vmID := data.Get("vm_id").(string)
	diskAttachments, err := client.ListDiskAttachments(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("list disk attachments of VM %s", vmID), err)
	}

	attachments := make([]map[string]interface{}, 0)

	for _, diskAttachment := range diskAttachments {
		attachment := make(map[string]interface{}, 0)

		attachment["id"] = diskAttachment.ID()
		attachment["disk_id"] = diskAttachment.DiskID()
		attachment["disk_interface"] = diskAttachment.DiskInterface()

		attachments = append(attachments, attachment)
	}

	if err := data.Set("attachments", attachments); err != nil {
		return errorToDiags("set attachments", err)
	}

	data.SetId(vmID)

	return nil
}
