package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client"
)

func (p *provider) waitForIPDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.waitForIPDataSourceRead,
		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:        schema.TypeString,
				Description: "ID of the oVirt VM.",
				Required:    true,
			},
		},
		Description: `This data source will wait for the VM to report an IP address.`,
	}
}

func (p *provider) waitForIPDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	var vmID = data.Get("vm_id").(string)
	_, err := client.WaitForNonLocalVMIPAddress(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags("waiting for IP", err)
	}
	data.SetId(vmID)
	return nil
}
