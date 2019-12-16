package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubnetRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"virtual_network_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"address_prefixes": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"network_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_configurations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"service_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	virtualNetworkName := d.Get("virtual_network_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, virtualNetworkName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Subnet %q (Virtual Network %q / Resource Group %q) was not found", name, virtualNetworkName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Subnet %q: %+v", name, err)
	}
	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_network_name", virtualNetworkName)

	if props := resp.SubnetPropertiesFormat; props != nil {
		if props.AddressPrefix != nil {
			d.Set("address_prefix", props.AddressPrefix)
		}
		if props.AddressPrefixes == nil {
			if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
				d.Set("address_prefixes", []string{*props.AddressPrefix})
			} else {
				d.Set("address_prefixes", []string{})
			}
		} else {
			d.Set("address_prefixes", props.AddressPrefixes)
		}

		if props.NetworkSecurityGroup != nil {
			d.Set("network_security_group_id", props.NetworkSecurityGroup.ID)
		} else {
			d.Set("network_security_group_id", "")
		}

		if props.RouteTable != nil {
			d.Set("route_table_id", props.RouteTable.ID)
		} else {
			d.Set("route_table_id", "")
		}

		if err := d.Set("ip_configurations", flattenSubnetIPConfigurations(props.IPConfigurations)); err != nil {
			return err
		}

		if err := d.Set("service_endpoints", flattenSubnetServiceEndpoints(props.ServiceEndpoints)); err != nil {
			return err
		}
	}

	return nil
}
