package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
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
			"interfaces": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the interface.",
						},
						"ipv4_addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "IP v4 addresses of the interface.",
						},
						"ipv6_addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "IP v6 addresses of the interface.",
						},
					},
				},
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
	vmID := data.Get("vm_id").(string)

	result, err := client.WaitForNonLocalVMIPAddress(ovirtclient.VMID(vmID))
	if err != nil {
		return errorToDiags("waiting for IP", err)
	}

	ifaces := make([]map[string]interface{}, 0)
	for ifname, ips := range result {
		iface := make(map[string]interface{}, 0)
		iface["name"] = ifname

		ipv4Addresses := make([]string, 0)
		ipv6Addresses := make([]string, 0)
		for _, ip := range ips {
			ipv4 := ip.To4()
			if ipv4 != nil {
				ipv4Addresses = append(ipv4Addresses, ip.String())
			} else {
				ipv6Addresses = append(ipv6Addresses, ip.String())
			}
		}

		iface["ipv4_addresses"] = ipv4Addresses
		iface["ipv6_addresses"] = ipv6Addresses

		ifaces = append(ifaces, iface)
	}

	if err := data.Set("interfaces", ifaces); err != nil {
		return errorToDiags("set interfaces", err)
	}
	data.SetId(vmID)

	return nil
}
