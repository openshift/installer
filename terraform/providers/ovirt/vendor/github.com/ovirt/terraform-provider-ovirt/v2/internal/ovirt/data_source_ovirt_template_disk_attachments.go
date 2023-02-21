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

func (p *provider) templateDiskAttachmentsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.templateDisksDataSourceRead,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "oVirt ID of the template to list the disk attachments for.",
				ValidateDiagFunc: validateUUID,
			},
			"disk_attachments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the disk attachment.",
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
		Description: `A set of all disk attachments of a template.`,
	}
}

func (p *provider) templateDisksDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	templateID := data.Get("template_id").(string)
	templateDiskAttachments, err := client.ListTemplateDiskAttachments(ovirtclient.TemplateID(templateID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("list disk attachments for template with ID '%s'", templateID), err)
	}

	attachments := make([]map[string]interface{}, 0)

	for _, templateDiskAttachment := range templateDiskAttachments {
		attachment := make(map[string]interface{}, 0)

		attachment["id"] = templateDiskAttachment.ID()
		attachment["disk_id"] = templateDiskAttachment.DiskID()
		attachment["disk_interface"] = templateDiskAttachment.DiskInterface()

		attachments = append(attachments, attachment)
	}

	if err := data.Set("disk_attachments", attachments); err != nil {
		return errorToDiags("set attachments", err)
	}

	data.SetId(templateID)

	return nil
}
