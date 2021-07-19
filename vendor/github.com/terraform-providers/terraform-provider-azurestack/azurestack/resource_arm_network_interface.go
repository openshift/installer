package azurestack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkInterfaceCreateUpdate,
		Read:   resourceArmNetworkInterfaceRead,
		Update: resourceArmNetworkInterfaceCreateUpdate,
		Delete: resourceArmNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"network_security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"subnet_id": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Dynamic),
								string(network.Static),
							}, true),
							StateFunc:        ignoreCaseStateFunc,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"application_gateway_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_inbound_nat_rules_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"application_security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_dns_name_label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"applied_dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_fqdn": {
				Type:     schema.TypeString,
				Optional: true,
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

			// enable_accelerated_networking is not in the profile for Azure Stack
			// "enable_accelerated_networking": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// 	Default:  false,
			// },

			"enable_ip_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkInterfaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureStack Network Interface creation.")

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	enableIpForwarding := d.Get("enable_ip_forwarding").(bool)
	// enableAcceleratedNetworking := d.Get("enable_accelerated_networking").(bool)
	tags := d.Get("tags").(map[string]interface{})

	properties := network.InterfacePropertiesFormat{
		EnableIPForwarding: &enableIpForwarding,
		Primary:            utils.Bool(true),
		// EnableAcceleratedNetworking: &enableAcceleratedNetworking,
	}

	if v, ok := d.GetOk("network_security_group_id"); ok {
		nsgId := v.(string)
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &nsgId,
		}

		networkSecurityGroupName, err := parseNetworkSecurityGroupName(nsgId)
		if err != nil {
			return err
		}

		azureStackLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureStackUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	fqdn, hasFQDN := d.GetOk("internal_fqdn")
	if hasDns || hasNameLabel || hasFQDN {
		ifaceDnsSettings := network.InterfaceDNSSettings{}

		if hasDns {
			var dnsServers []string
			dns := dns.(*schema.Set).List()
			for _, v := range dns {
				str := v.(string)
				dnsServers = append(dnsServers, str)
			}
			ifaceDnsSettings.DNSServers = &dnsServers
		}

		if hasNameLabel {
			name_label := nameLabel.(string)
			ifaceDnsSettings.InternalDNSNameLabel = &name_label
		}

		if hasFQDN {
			ifaceDnsSettings.InternalFqdn = utils.String(fqdn.(string))
		}

		properties.DNSSettings = &ifaceDnsSettings
	}

	ipConfigs, subnetnToLock, vnnToLock, sgErr := expandAzureStackNetworkInterfaceIpConfigurations(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Network Interface IP Configurations: %+v", sgErr)
	}

	azureStackLockMultipleByName(subnetnToLock, subnetResourceName)
	defer azureStackUnlockMultipleByName(subnetnToLock, subnetResourceName)

	azureStackLockMultipleByName(vnnToLock, virtualNetworkResourceName)
	defer azureStackUnlockMultipleByName(vnnToLock, virtualNetworkResourceName)

	if len(ipConfigs) > 0 {
		properties.IPConfigurations = &ipConfigs
	}

	iface := network.Interface{
		Name:                      &name,
		Location:                  &location,
		InterfacePropertiesFormat: &properties,
		Tags:                      *expandTags(tags),
	}

	data, _ := json.Marshal(&iface)
	fmt.Println(string(data))
	future, err := client.CreateOrUpdate(ctx, resGroup, name, iface)
	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read NIC %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkInterfaceRead(d, meta)
}

func resourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

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

func resourceArmNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	if v, ok := d.GetOk("network_security_group_id"); ok {
		networkSecurityGroupId := v.(string)
		networkSecurityGroupName, err := parseNetworkSecurityGroupName(networkSecurityGroupId)
		if err != nil {
			return err
		}

		azureStackLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureStackUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	configs := d.Get("ip_configuration").([]interface{})
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		subnetId, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return err
		}
		subnetName := subnetId.Path["subnets"]
		if !sliceContainsValue(subnetNamesToLock, subnetName) {
			subnetNamesToLock = append(subnetNamesToLock, subnetName)
		}

		virtualNetworkName := subnetId.Path["virtualNetworks"]
		if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
			virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
		}
	}

	azureStackLockMultipleByName(&subnetNamesToLock, subnetResourceName)
	defer azureStackUnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	azureStackLockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer azureStackUnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return err
}

func flattenNetworkInterfaceIPConfigurations(ipConfigs *[]network.InterfaceIPConfiguration) []interface{} {
	result := make([]interface{}, 0, len(*ipConfigs))
	for _, ipConfig := range *ipConfigs {
		niIPConfig := make(map[string]interface{})

		props := ipConfig.InterfaceIPConfigurationPropertiesFormat

		niIPConfig["name"] = *ipConfig.Name
		niIPConfig["subnet_id"] = *props.Subnet.ID
		niIPConfig["private_ip_address_allocation"] = strings.ToLower(string(props.PrivateIPAllocationMethod))

		if props.PrivateIPAllocationMethod == network.Static {
			niIPConfig["private_ip_address"] = *props.PrivateIPAddress
		}

		if props.PublicIPAddress != nil {
			niIPConfig["public_ip_address_id"] = *props.PublicIPAddress.ID
		}

		if props.Primary != nil {
			niIPConfig["primary"] = *props.Primary
		}

		// var poolsAG []interface{}
		// if props.ApplicationGatewayBackendAddressPools != nil {
		// 	for _, pool := range *props.ApplicationGatewayBackendAddressPools {
		// 		poolsAG = append(poolsAG, *pool.ID)
		// 	}
		// }
		// niIPConfig["application_gateway_backend_address_pools_ids"] = schema.NewSet(schema.HashString, poolsAG)

		var pools []interface{}
		if props.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *props.LoadBalancerBackendAddressPools {
				pools = append(pools, *pool.ID)
			}
		}
		niIPConfig["load_balancer_backend_address_pools_ids"] = schema.NewSet(schema.HashString, pools)

		var rules []interface{}
		if props.LoadBalancerInboundNatRules != nil {
			for _, rule := range *props.LoadBalancerInboundNatRules {
				rules = append(rules, *rule.ID)
			}
		}
		niIPConfig["load_balancer_inbound_nat_rules_ids"] = schema.NewSet(schema.HashString, rules)

		// securityGroups := make([]interface{}, 0)
		// if sgs := props.ApplicationSecurityGroups; sgs != nil {
		// 	for _, sg := range *sgs {
		// 		securityGroups = append(securityGroups, *sg.ID)
		// 	}
		// }
		// niIPConfig["application_security_group_ids"] = schema.NewSet(schema.HashString, securityGroups)

		result = append(result, niIPConfig)
	}
	return result
}

func expandAzureStackNetworkInterfaceIpConfigurations(d *schema.ResourceData) ([]network.InterfaceIPConfiguration, *[]string, *[]string, error) {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.InterfaceIPConfiguration, 0, len(configs))
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		private_ip_allocation_method := data["private_ip_address_allocation"].(string)

		allocationMethod := network.IPAllocationMethod(private_ip_allocation_method)
		properties := network.InterfaceIPConfigurationPropertiesFormat{
			Subnet: &network.Subnet{
				ID: &subnet_id,
			},
			PrivateIPAllocationMethod: allocationMethod,
		}

		subnetId, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return []network.InterfaceIPConfiguration{}, nil, nil, err
		}

		subnetName := subnetId.Path["subnets"]
		virtualNetworkName := subnetId.Path["virtualNetworks"]

		if !sliceContainsValue(subnetNamesToLock, subnetName) {
			subnetNamesToLock = append(subnetNamesToLock, subnetName)
		}

		if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
			virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &network.PublicIPAddress{
				ID: &v,
			}
		}

		// Come back to this later
		if v, ok := data["primary"]; ok {
			b := v.(bool)
			properties.Primary = &b
		}

		//also this
		// if v, ok := data["application_gateway_backend_address_pools_ids"]; ok {
		// 	var ids []network.ApplicationGatewayBackendAddressPool
		// 	pools := v.(*schema.Set).List()
		// 	for _, p := range pools {
		// 		pool_id := p.(string)
		// 		id := network.ApplicationGatewayBackendAddressPool{
		// 			ID: &pool_id,
		// 		}
		//
		// 		ids = append(ids, id)
		// 	}
		//
		// 	properties.ApplicationGatewayBackendAddressPools = &ids
		// }

		if v, ok := data["load_balancer_backend_address_pools_ids"]; ok {
			var ids []network.BackendAddressPool
			pools := v.(*schema.Set).List()
			for _, p := range pools {
				pool_id := p.(string)
				id := network.BackendAddressPool{
					ID: &pool_id,
				}

				ids = append(ids, id)
			}

			properties.LoadBalancerBackendAddressPools = &ids
		}

		if v, ok := data["load_balancer_inbound_nat_rules_ids"]; ok {
			var natRules []network.InboundNatRule
			rules := v.(*schema.Set).List()
			for _, r := range rules {
				rule_id := r.(string)
				rule := network.InboundNatRule{
					ID: &rule_id,
				}

				natRules = append(natRules, rule)
			}

			properties.LoadBalancerInboundNatRules = &natRules
		}

		// if v, ok := data["application_security_group_ids"]; ok {
		// 	var securityGroups []network.ApplicationSecurityGroup
		// 	rules := v.(*schema.Set).List()
		// 	for _, r := range rules {
		// 		groupId := r.(string)
		// 		group := network.ApplicationSecurityGroup{
		// 			ID: &groupId,
		// 		}
		//
		// 		securityGroups = append(securityGroups, group)
		// 	}
		//
		// 	properties.ApplicationSecurityGroups = &securityGroups
		// }

		name := data["name"].(string)
		ipConfig := network.InterfaceIPConfiguration{
			Name:                                     &name,
			InterfaceIPConfigurationPropertiesFormat: &properties,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	// Come back to this later
	// if we've got multiple IP Configurations - one must be designated Primary
	if len(ipConfigs) > 1 {
		hasPrimary := false
		for _, config := range ipConfigs {
			if config.Primary != nil && *config.Primary {
				hasPrimary = true
				break
			}
		}

		if !hasPrimary {
			return nil, nil, nil, fmt.Errorf("If multiple `ip_configurations` are specified - one must be designated as `primary`.")
		}
	}

	return ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func sliceContainsValue(input []string, value string) bool {
	for _, v := range input {
		if v == value {
			return true
		}
	}

	return false
}
