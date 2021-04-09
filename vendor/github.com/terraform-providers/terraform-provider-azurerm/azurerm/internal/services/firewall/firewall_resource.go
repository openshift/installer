package firewall

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var azureFirewallResourceName = "azurerm_firewall"

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallCreateUpdate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallCreateUpdate,
		Delete: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			// TODO 3.0: change this to required
			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AZFWHub),
					string(network.AZFWVNet),
				}, false),
			},

			// TODO 3.0: change this to required
			"sku_tier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Premium),
					string(network.Standard),
				}, false),
			},

			"firewall_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.FirewallPolicyID,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.FirewallSubnetName,
						},
						"public_ip_address_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.PublicIpAddressID,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"management_ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.FirewallManagementSubnetName,
						},
						"public_ip_address_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.PublicIpAddressID,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"threat_intel_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(network.AzureFirewallThreatIntelModeAlert),
				ValidateFunc: validation.StringInSlice([]string{
					// TODO 3.0: remove the default value and the `""` below. So if it is not specified
					// in config, it will not be send in request, which is required in case of vhub.
					"",
					string(network.AzureFirewallThreatIntelModeOff),
					string(network.AzureFirewallThreatIntelModeAlert),
					string(network.AzureFirewallThreatIntelModeDeny),
				}, false),
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},

			"virtual_hub": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_hub_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.VirtualHubID,
						},
						"public_ip_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      1,
						},
						"public_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"zones": azure.SchemaZones(),

			"tags": tags.Schema(),
		},
	}
}

func resourceFirewallCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Firewall %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if d.IsNewResource() {
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall", *existing.ID)
		}
	}

	if err := validateFirewallIPConfigurationSettings(d.Get("ip_configuration").([]interface{})); err != nil {
		return fmt.Errorf("Error validating Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	i := d.Get("ip_configuration").([]interface{})
	ipConfigs, subnetToLock, vnetToLock, err := expandFirewallIPConfigurations(i)
	if err != nil {
		return fmt.Errorf("Error building list of Azure Firewall IP Configurations: %+v", err)
	}
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	parameters := network.AzureFirewall{
		Location: &location,
		Tags:     tags.Expand(t),
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations:     ipConfigs,
			ThreatIntelMode:      network.AzureFirewallThreatIntelMode(d.Get("threat_intel_mode").(string)),
			AdditionalProperties: expandFirewallDNSServers(d.Get("dns_servers").([]interface{})),
		},
		Zones: zones,
	}

	m := d.Get("management_ip_configuration").([]interface{})
	if len(m) == 1 {
		mgmtIPConfig, mgmtSubnetName, mgmtVirtualNetworkName, err := expandFirewallIPConfigurations(m)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Firewall Management IP Configurations: %+v", err)
		}

		if !utils.SliceContainsValue(*subnetToLock, (*mgmtSubnetName)[0]) {
			*subnetToLock = append(*subnetToLock, (*mgmtSubnetName)[0])
		}

		if !utils.SliceContainsValue(*vnetToLock, (*mgmtVirtualNetworkName)[0]) {
			*vnetToLock = append(*vnetToLock, (*mgmtVirtualNetworkName)[0])
		}
		if *mgmtIPConfig != nil {
			parameters.ManagementIPConfiguration = &(*mgmtIPConfig)[0]
		}
	}

	if threatIntelMode := d.Get("threat_intel_mode").(string); threatIntelMode != "" {
		parameters.AzureFirewallPropertiesFormat.ThreatIntelMode = network.AzureFirewallThreatIntelMode(threatIntelMode)
	}

	if policyId := d.Get("firewall_policy_id").(string); policyId != "" {
		parameters.AzureFirewallPropertiesFormat.FirewallPolicy = &network.SubResource{ID: &policyId}
	}

	vhub, hubIpAddresses, ok := expandFirewallVirtualHubSetting(existing, d.Get("virtual_hub").([]interface{}))
	if ok {
		parameters.AzureFirewallPropertiesFormat.VirtualHub = vhub
		parameters.AzureFirewallPropertiesFormat.HubIPAddresses = hubIpAddresses
	}

	// TODO 3.0: no need to test since sku_name is required
	if skuName := d.Get("sku_name").(string); skuName != "" {
		if parameters.Sku == nil {
			parameters.Sku = &network.AzureFirewallSku{}
		}
		parameters.Sku.Name = network.AzureFirewallSkuName(skuName)
	}

	// TODO 3.0: no need to test since sku_tier is required
	if skuTier := d.Get("sku_tier").(string); skuTier != "" {
		if parameters.Sku == nil {
			parameters.Sku = &network.AzureFirewallSku{}
		}
		parameters.Sku.Tier = network.AzureFirewallSkuTier(skuTier)
	}

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(vnetToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetToLock, SubnetResourceName)

	if !d.IsNewResource() {
		exists, err2 := client.Get(ctx, resourceGroup, name)
		if err2 != nil {
			if utils.ResponseWasNotFound(exists.Response) {
				return fmt.Errorf("Error retrieving existing Firewall %q (Resource Group %q): firewall not found in resource group", name, resourceGroup)
			}
			return fmt.Errorf("Error retrieving existing Firewall %q (Resource Group %q): %s", name, resourceGroup, err2)
		}
		if exists.AzureFirewallPropertiesFormat == nil {
			return fmt.Errorf("Error retrieving existing rules (Firewall %q / Resource Group %q): `props` was nil", name, resourceGroup)
		}
		props := *exists.AzureFirewallPropertiesFormat
		parameters.AzureFirewallPropertiesFormat.ApplicationRuleCollections = props.ApplicationRuleCollections
		parameters.AzureFirewallPropertiesFormat.NetworkRuleCollections = props.NetworkRuleCollections
		parameters.AzureFirewallPropertiesFormat.NatRuleCollections = props.NatRuleCollections
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceFirewallRead(d, meta)
}

func resourceFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.AzureFirewallName

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if err := d.Set("ip_configuration", flattenFirewallIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}
		managementIPConfigs := make([]interface{}, 0)
		if props.ManagementIPConfiguration != nil {
			managementIPConfigs = flattenFirewallIPConfigurations(&[]network.AzureFirewallIPConfiguration{
				*props.ManagementIPConfiguration,
			})
		}
		if err := d.Set("management_ip_configuration", managementIPConfigs); err != nil {
			return fmt.Errorf("Error setting `management_ip_configuration`: %+v", err)
		}

		d.Set("threat_intel_mode", string(props.ThreatIntelMode))

		if err := d.Set("dns_servers", flattenFirewallDNSServers(props.AdditionalProperties)); err != nil {
			return fmt.Errorf("Error setting `dns_servers`: %+v", err)
		}

		if policy := props.FirewallPolicy; policy != nil {
			d.Set("firewall_policy_id", policy.ID)
		}

		if sku := props.Sku; sku != nil {
			d.Set("sku_name", string(sku.Name))
			d.Set("sku_tier", string(sku.Tier))
		}

		if err := d.Set("virtual_hub", flattenFirewallVirtualHubSetting(props)); err != nil {
			return fmt.Errorf("Error setting `virtual_hub`: %+v", err)
		}
	}

	if err := d.Set("zones", azure.FlattenZones(read.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if config.Subnet == nil || config.Subnet.ID == nil {
					continue
				}

				parsedSubnetID, err2 := azure.ParseAzureResourceID(*config.Subnet.ID)
				if err2 != nil {
					return err2
				}
				subnetName := parsedSubnetID.Path["subnets"]

				if !utils.SliceContainsValue(subnetNamesToLock, subnetName) {
					subnetNamesToLock = append(subnetNamesToLock, subnetName)
				}

				virtualNetworkName := parsedSubnetID.Path["virtualNetworks"]
				if !utils.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
				}
			}
		}

		if mconfig := props.ManagementIPConfiguration; mconfig != nil {
			if mconfig.Subnet != nil && mconfig.Subnet.ID != nil {
				parsedSubnetID, err2 := azure.ParseAzureResourceID(*mconfig.Subnet.ID)
				if err2 != nil {
					return err2
				}
				subnetName := parsedSubnetID.Path["subnets"]

				if !utils.SliceContainsValue(subnetNamesToLock, subnetName) {
					subnetNamesToLock = append(subnetNamesToLock, subnetName)
				}

				virtualNetworkName := parsedSubnetID.Path["virtualNetworks"]
				if !utils.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
				}
			}
		}
	}

	locks.ByName(name, azureFirewallResourceName)
	defer locks.UnlockByName(name, azureFirewallResourceName)

	locks.MultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)

	locks.MultipleByName(&subnetNamesToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(&subnetNamesToLock, SubnetResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandFirewallIPConfigurations(configs []interface{}) (*[]network.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0)
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnetId := data["subnet_id"].(string)
		pubID := data["public_ip_address_id"].(string)

		ipConfig := network.AzureFirewallIPConfiguration{
			Name: utils.String(name),
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.SubResource{
					ID: utils.String(pubID),
				},
			},
		}

		if subnetId != "" {
			subnetID, err := azure.ParseAzureResourceID(subnetId)
			if err != nil {
				return nil, nil, nil, err
			}

			subnetName := subnetID.Path["subnets"]
			virtualNetworkName := subnetID.Path["virtualNetworks"]

			if !utils.SliceContainsValue(subnetNamesToLock, subnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetName)
			}

			if !utils.SliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
			}

			ipConfig.AzureFirewallIPConfigurationPropertiesFormat.Subnet = &network.SubResource{
				ID: utils.String(subnetId),
			}
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return &ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func flattenFirewallIPConfigurations(input *[]network.AzureFirewallIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		afIPConfig := make(map[string]interface{})
		props := v.AzureFirewallIPConfigurationPropertiesFormat
		if props == nil {
			continue
		}

		if name := v.Name; name != nil {
			afIPConfig["name"] = *name
		}

		if subnet := props.Subnet; subnet != nil {
			if id := subnet.ID; id != nil {
				afIPConfig["subnet_id"] = *id
			}
		}

		if ipAddress := props.PrivateIPAddress; ipAddress != nil {
			afIPConfig["private_ip_address"] = *ipAddress
		}

		if pip := props.PublicIPAddress; pip != nil {
			if id := pip.ID; id != nil {
				afIPConfig["public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func expandFirewallDNSServers(input []interface{}) map[string]*string {
	if len(input) == 0 {
		return nil
	}

	var servers []string
	for _, server := range input {
		servers = append(servers, server.(string))
	}

	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/11278
	return map[string]*string{
		"Network.DNS.EnableProxy": utils.String("true"),
		"Network.DNS.Servers":     utils.String(strings.Join(servers, ",")),
	}
}

func flattenFirewallDNSServers(input map[string]*string) []interface{} {
	if len(input) == 0 {
		return nil
	}

	enabled := false
	if enabledPtr := input["Network.DNS.EnableProxy"]; enabledPtr != nil {
		enabled = *enabledPtr == "true"
	}

	if !enabled {
		return nil
	}

	servers := []string{}
	if serversPtr := input["Network.DNS.Servers"]; serversPtr != nil {
		servers = strings.Split(*serversPtr, ",")
	}
	return utils.FlattenStringSlice(&servers)
}

func expandFirewallVirtualHubSetting(existing network.AzureFirewall, input []interface{}) (vhub *network.SubResource, ipAddresses *network.HubIPAddresses, ok bool) {
	if len(input) == 0 {
		return nil, nil, false
	}

	b := input[0].(map[string]interface{})

	// The API requires both "Count" and "Addresses" for the "PublicIPs" setting.
	// The "Count" means how many PIP to provision.
	// The "Addresses" means differently in different cases:
	// - Create: only "Count" is needed, "Addresses" is not necessary
	// - Update: both "Count" and "Addresses" are needed:
	//     Scale up: "Addresses" should remain same as before scaling up
	//     Scale down: "Addresses" should indicate the addresses to be retained (in this case we retain the first new "Count" ones)
	newCount := b["public_ip_count"].(int)
	var addresses *[]network.AzureFirewallPublicIPAddress
	if prop := existing.AzureFirewallPropertiesFormat; prop != nil {
		if ipaddress := prop.HubIPAddresses; ipaddress != nil {
			if pips := ipaddress.PublicIPs; pips != nil {
				if count := pips.Count; count != nil {
					oldCount := int(*count)
					addresses = pips.Addresses

					// In case of scale down, keep the first new "Count" addresses.
					if oldCount > newCount {
						keptAddresses := make([]network.AzureFirewallPublicIPAddress, newCount)
						for i := 0; i < newCount; i++ {
							keptAddresses[i] = (*addresses)[i]
						}
						addresses = &keptAddresses
					}
				}
			}
		}
	}

	vhub = &network.SubResource{ID: utils.String(b["virtual_hub_id"].(string))}
	ipAddresses = &network.HubIPAddresses{
		PublicIPs: &network.HubPublicIPAddresses{
			Count:     utils.Int32(int32(b["public_ip_count"].(int))),
			Addresses: addresses,
		},
	}

	return vhub, ipAddresses, true
}

func flattenFirewallVirtualHubSetting(props *network.AzureFirewallPropertiesFormat) []interface{} {
	if props.VirtualHub == nil {
		return nil
	}

	var vhubId string
	if props.VirtualHub.ID != nil {
		vhubId = *props.VirtualHub.ID
	}

	var (
		publicIpCount int
		publicIps     []string
		privateIp     string
	)
	if hubIP := props.HubIPAddresses; hubIP != nil {
		if hubIP.PrivateIPAddress != nil {
			privateIp = *hubIP.PrivateIPAddress
		}
		if pubIPs := hubIP.PublicIPs; pubIPs != nil {
			if pubIPs.Count != nil {
				publicIpCount = int(*pubIPs.Count)
			}
			if pubIPs.Addresses != nil {
				for _, addr := range *pubIPs.Addresses {
					if addr.Address != nil {
						publicIps = append(publicIps, *addr.Address)
					}
				}
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"virtual_hub_id":      vhubId,
			"public_ip_count":     publicIpCount,
			"public_ip_addresses": publicIps,
			"private_ip_address":  privateIp,
		},
	}
}

func validateFirewallIPConfigurationSettings(configs []interface{}) error {
	if len(configs) == 0 {
		return nil
	}

	subnetNumber := 0

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		if subnet, exist := data["subnet_id"].(string); exist && subnet != "" {
			subnetNumber++
		}
	}

	if subnetNumber != 1 {
		return fmt.Errorf(`The "ip_configuration" is invalid, %d "subnet_id" have been set, one "subnet_id" should be set among all "ip_configuration" blocks`, subnetNumber)
	}

	return nil
}
