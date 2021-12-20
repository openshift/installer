package azurestack

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetworkInterfaceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"network_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"application_gateway_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_inbound_nat_rules_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_dns_name_label": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"applied_dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			/**
			 * As of 2018-01-06: AN (aka. SR-IOV) on Azure is GA on Windows and Linux.
			 *
			 * Refer to: https://azure.microsoft.com/en-us/blog/maximize-your-vm-s-performance-with-accelerated-networking-now-generally-available-for-both-windows-and-linux/
			 *
			 * Refer to: https://docs.microsoft.com/en-us/azure/virtual-network/create-vm-accelerated-networking-cli
			 * For details, VM configuration and caveats.
			 */

			// enable_accelerated_networking is not supported in the profile 2017-03-09
			// "enable_accelerated_networking": {
			// 	Type:     schema.TypeBool,
			// 	Computed: true,
			// },

			"enable_ip_forwarding": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	if iface := resp.InterfacePropertiesFormat; iface != nil {
		d.Set("mac_address", iface.MacAddress)
		d.Set("enable_ip_forwarding", iface.EnableIPForwarding)

		if iface.IPConfigurations != nil && len(*iface.IPConfigurations) > 0 {
			configs := *iface.IPConfigurations

			if configs[0].InterfaceIPConfigurationPropertiesFormat != nil {
				privateIPAddress := configs[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
				d.Set("private_ip_address", privateIPAddress)
			}

			addresses := make([]interface{}, 0)
			for _, config := range configs {
				if config.InterfaceIPConfigurationPropertiesFormat != nil {
					addresses = append(addresses, *config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
				}
			}

			if err := d.Set("private_ip_addresses", addresses); err != nil {
				return fmt.Errorf("Error setting `private_ip_addresses`: %+v", err)
			}
		}

		if iface.IPConfigurations != nil {
			if err := d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(iface.IPConfigurations)); err != nil {
				return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
			}
		}

		if iface.VirtualMachine != nil {
			d.Set("virtual_machine_id", iface.VirtualMachine.ID)
		} else {
			d.Set("virtual_machine_id", "")
		}

		if dnsSettings := iface.DNSSettings; dnsSettings != nil {
			d.Set("applied_dns_servers", dnsSettings.AppliedDNSServers)
			d.Set("dns_servers", dnsSettings.DNSServers)
			d.Set("internal_fqdn", dnsSettings.InternalFqdn)
			d.Set("internal_dns_name_label", dnsSettings.InternalDNSNameLabel)
		}

		if iface.NetworkSecurityGroup != nil {
			d.Set("network_security_group_id", resp.NetworkSecurityGroup.ID)
		} else {
			d.Set("network_security_group_id", "")
		}
	}

	// enable_accelerated_networking is not supported in the profile used for
	// AzureStack
	// d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)

	flattenAndSetTags(d, &resp.Tags)

	return nil
}
