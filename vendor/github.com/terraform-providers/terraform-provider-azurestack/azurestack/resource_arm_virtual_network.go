package azurestack

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var virtualNetworkResourceName = "azurestack_virtual_network"

func resourceArmVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkCreate,
		Read:   resourceArmVirtualNetworkRead,
		Update: resourceArmVirtualNetworkCreate,
		Delete: resourceArmVirtualNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subnet": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceAzureSubnetHash,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network creation.")

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	vnetProperties, vnetPropsErr := getVirtualNetworkProperties(ctx, d, meta)
	if vnetPropsErr != nil {
		return vnetPropsErr
	}

	vnet := network.VirtualNetwork{
		Name:                           &name,
		Location:                       &location,
		VirtualNetworkPropertiesFormat: vnetProperties,
		Tags:                           *expandTags(tags),
	}

	networkSecurityGroupNames := make([]string, 0)
	for _, subnet := range *vnet.VirtualNetworkPropertiesFormat.Subnets {
		if subnet.NetworkSecurityGroup != nil {
			nsgName, err := parseNetworkSecurityGroupName(*subnet.NetworkSecurityGroup.ID)
			if err != nil {
				return err
			}

			if !sliceContainsValue(networkSecurityGroupNames, nsgName) {
				networkSecurityGroupNames = append(networkSecurityGroupNames, nsgName)
			}
		}
	}

	azureStackLockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer azureStackUnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	future, err := client.CreateOrUpdate(ctx, resGroup, name, vnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Network %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkRead(d, meta)
}

func resourceArmVirtualNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualNetworks"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	vnet := *resp.VirtualNetworkPropertiesFormat

	// update appropriate values
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("address_space", vnet.AddressSpace.AddressPrefixes)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	subnets := &schema.Set{
		F: resourceAzureSubnetHash,
	}

	for _, subnet := range *vnet.Subnets {
		s := map[string]interface{}{}

		s["name"] = *subnet.Name
		s["address_prefix"] = *subnet.SubnetPropertiesFormat.AddressPrefix
		if subnet.SubnetPropertiesFormat.NetworkSecurityGroup != nil {
			s["security_group"] = *subnet.SubnetPropertiesFormat.NetworkSecurityGroup.ID
		}

		subnets.Add(s)
	}
	d.Set("subnet", subnets)

	if vnet.DhcpOptions != nil && vnet.DhcpOptions.DNSServers != nil {
		d.Set("dns_servers", *vnet.DhcpOptions.DNSServers)
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualNetworks"]

	nsgNames, err := expandAzureStackVirtualNetworkVirtualNetworkSecurityGroupNames(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Error parsing Network Security Group ID's: %+v", err)
	}

	azureStackLockMultipleByName(&nsgNames, virtualNetworkResourceName)
	defer azureStackUnlockMultipleByName(&nsgNames, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getVirtualNetworkProperties(ctx context.Context, d *schema.ResourceData, meta interface{}) (*network.VirtualNetworkPropertiesFormat, error) {
	// first; get address space prefixes:
	prefixes := []string{}
	for _, prefix := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, prefix.(string))
	}

	// then; the dns servers:
	dnses := []string{}
	for _, dns := range d.Get("dns_servers").([]interface{}) {
		dnses = append(dnses, dns.(string))
	}

	// then; the subnets:
	subnets := []network.Subnet{}
	if subs := d.Get("subnet").(*schema.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			//since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			resGroup := d.Get("resource_group_name").(string)
			vnetName := d.Get("name").(string)
			subnetObj, err := getExistingSubnet(ctx, resGroup, vnetName, name, meta)
			if err != nil {
				return nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props ")

			prefix := subnet["address_prefix"].(string)
			secGroup := subnet["security_group"].(string)

			//set the props from config and leave the rest intact
			subnetObj.Name = &name
			if subnetObj.SubnetPropertiesFormat == nil {
				subnetObj.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{}
			}

			subnetObj.SubnetPropertiesFormat.AddressPrefix = &prefix

			if secGroup != "" {
				subnetObj.SubnetPropertiesFormat.NetworkSecurityGroup = &network.SecurityGroup{
					ID: &secGroup,
				}
			} else {
				subnetObj.SubnetPropertiesFormat.NetworkSecurityGroup = nil
			}

			subnets = append(subnets, *subnetObj)
		}
	}

	properties := &network.VirtualNetworkPropertiesFormat{
		AddressSpace: &network.AddressSpace{
			AddressPrefixes: &prefixes,
		},
		DhcpOptions: &network.DhcpOptions{
			DNSServers: &dnses,
		},
		Subnets: &subnets,
	}
	// finally; return the struct:
	return properties, nil
}

func resourceAzureSubnetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["name"].(string))
	buf.WriteString(m["address_prefix"].(string))
	if v, ok := m["security_group"]; ok {
		buf.WriteString(v.(string))
	}
	return hashcode.String(buf.String())
}

func getExistingSubnet(ctx context.Context, resGroup string, vnetName string, subnetName string, meta interface{}) (*network.Subnet, error) {
	//attempt to retrieve existing subnet from the server
	existingSubnet := network.Subnet{}
	subnetClient := meta.(*ArmClient).subnetClient
	resp, err := subnetClient.Get(ctx, resGroup, vnetName, subnetName, "")

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return &existingSubnet, nil
		}
		//raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	existingSubnet.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{}
	existingSubnet.SubnetPropertiesFormat.AddressPrefix = resp.SubnetPropertiesFormat.AddressPrefix

	if resp.SubnetPropertiesFormat.NetworkSecurityGroup != nil {
		existingSubnet.SubnetPropertiesFormat.NetworkSecurityGroup = resp.SubnetPropertiesFormat.NetworkSecurityGroup
	}

	if resp.SubnetPropertiesFormat.RouteTable != nil {
		existingSubnet.SubnetPropertiesFormat.RouteTable = resp.SubnetPropertiesFormat.RouteTable
	}

	if resp.SubnetPropertiesFormat.IPConfigurations != nil {
		existingSubnet.SubnetPropertiesFormat.IPConfigurations = resp.SubnetPropertiesFormat.IPConfigurations
	}

	return &existingSubnet, nil
}

func expandAzureStackVirtualNetworkVirtualNetworkSecurityGroupNames(d *schema.ResourceData) ([]string, error) {
	nsgNames := make([]string, 0)

	if v, ok := d.GetOk("subnet"); ok {
		subnets := v.(*schema.Set).List()
		for _, subnet := range subnets {
			subnet, ok := subnet.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR] Subnet should be a Hash - was '%+v'", subnet)
			}

			networkSecurityGroupId := subnet["security_group"].(string)
			if networkSecurityGroupId != "" {
				nsgName, err := parseNetworkSecurityGroupName(networkSecurityGroupId)
				if err != nil {
					return nil, err
				}

				if !sliceContainsValue(nsgNames, nsgName) {
					nsgNames = append(nsgNames, nsgName)
				}
			}
		}
	}

	return nsgNames, nil
}
