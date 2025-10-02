package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var VirtualNetworkResourceName = "azurerm_virtual_network"

func resourceVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualNetworkCreateUpdate,
		Read:   resourceVirtualNetworkRead,
		Update: resourceVirtualNetworkCreateUpdate,
		Delete: resourceVirtualNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"bgp_community": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.VirtualNetworkBgpCommunity,
			},

			"ddos_protection_plan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"vm_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"address_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							// TODO Remove this in the next major version release
							Deprecated: "Use the `address_prefixes` property instead.",
						},

						"address_prefixes": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"security_group": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceAzureSubnetHash,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualNetworkCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_network", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	vnetProperties, err := expandVirtualNetworkProperties(ctx, d, meta)
	if err != nil {
		return err
	}

	vnet := network.VirtualNetwork{
		Name:                           utils.String(id.Name),
		Location:                       utils.String(location),
		VirtualNetworkPropertiesFormat: vnetProperties,
		Tags:                           tags.Expand(t),
	}

	networkSecurityGroupNames := make([]string, 0)
	for _, subnet := range *vnet.VirtualNetworkPropertiesFormat.Subnets {
		if subnet.NetworkSecurityGroup != nil {
			parsedNsgID, err := parse.NetworkSecurityGroupID(*subnet.NetworkSecurityGroup.ID)
			if err != nil {
				return err
			}

			networkSecurityGroupName := parsedNsgID.Name
			if !utils.SliceContainsValue(networkSecurityGroupNames, networkSecurityGroupName) {
				networkSecurityGroupNames = append(networkSecurityGroupNames, networkSecurityGroupName)
			}
		}
	}

	locks.MultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer locks.UnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, vnet)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkRead(d, meta)
}

func resourceVirtualNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		d.Set("guid", props.ResourceGUID)

		if space := props.AddressSpace; space != nil {
			d.Set("address_space", utils.FlattenStringSlice(space.AddressPrefixes))
		}

		if err := d.Set("ddos_protection_plan", flattenVirtualNetworkDDoSProtectionPlan(props)); err != nil {
			return fmt.Errorf("Error setting `ddos_protection_plan`: %+v", err)
		}

		if err := d.Set("subnet", flattenVirtualNetworkSubnets(props.Subnets)); err != nil {
			return fmt.Errorf("Error setting `subnets`: %+v", err)
		}

		if err := d.Set("dns_servers", flattenVirtualNetworkDNSServers(props.DhcpOptions)); err != nil {
			return fmt.Errorf("Error setting `dns_servers`: %+v", err)
		}

		bgpCommunity := ""
		if p := props.BgpCommunities; p != nil {
			if v := p.VirtualNetworkCommunity; v != nil {
				bgpCommunity = *v
			}
		}
		if err := d.Set("bgp_community", bgpCommunity); err != nil {
			return fmt.Errorf("Error setting `bgp_community`: %+v", err)
		}

		d.Set("vm_protection_enabled", props.EnableVMProtection)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	nsgNames, err := expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d)
	if err != nil {
		return fmt.Errorf("Error parsing Network Security Group ID's: %+v", err)
	}

	locks.MultipleByName(&nsgNames, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&nsgNames, VirtualNetworkResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualNetworkProperties(ctx context.Context, d *schema.ResourceData, meta interface{}) (*network.VirtualNetworkPropertiesFormat, error) {
	subnets := make([]network.Subnet, 0)
	if subs := d.Get("subnet").(*schema.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			// since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			resGroup := d.Get("resource_group_name").(string)
			vnetName := d.Get("name").(string)
			subnetObj, err := getExistingSubnet(ctx, resGroup, vnetName, name, meta)
			if err != nil {
				return nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props ")

			var prefixSet int = 0
			properties := network.SubnetPropertiesFormat{}
			if value, ok := subnet["address_prefixes"]; ok {
				var addressPrefixes []string
				for _, item := range value.([]interface{}) {
					addressPrefixes = append(addressPrefixes, item.(string))
				}
				properties.AddressPrefixes = &addressPrefixes
				if len(addressPrefixes) > 0 {
					prefixSet++
				}
			}
			if value, ok := subnet["address_prefix"]; ok {
				addressPrefix := value.(string)
				properties.AddressPrefix = &addressPrefix
				if len(addressPrefix) > 0 {
					prefixSet++
				}
			}
			if properties.AddressPrefixes != nil && len(*properties.AddressPrefixes) == 1 {
				properties.AddressPrefix = &(*properties.AddressPrefixes)[0]
				properties.AddressPrefixes = nil
			}
			if prefixSet == 0 {
				return nil, fmt.Errorf("one of `address_prefix` or `address_prefixes` must be set")
			}
			if prefixSet == 2 {
				return nil, fmt.Errorf("only one of `address_prefix` or `address_prefixes` can be set")
			}

			secGroup := subnet["security_group"].(string)

			// set the props from config and leave the rest intact
			subnetObj.Name = &name
			if subnetObj.SubnetPropertiesFormat == nil {
				subnetObj.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{}
			}

			subnetObj.SubnetPropertiesFormat.AddressPrefixes = properties.AddressPrefixes
			subnetObj.SubnetPropertiesFormat.AddressPrefix = properties.AddressPrefix

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
			AddressPrefixes: utils.ExpandStringSlice(d.Get("address_space").([]interface{})),
		},
		DhcpOptions: &network.DhcpOptions{
			DNSServers: utils.ExpandStringSlice(d.Get("dns_servers").([]interface{})),
		},
		EnableVMProtection: utils.Bool(d.Get("vm_protection_enabled").(bool)),
		Subnets:            &subnets,
	}

	if v, ok := d.GetOk("ddos_protection_plan"); ok {
		rawList := v.([]interface{})

		var ddosPPlan map[string]interface{}
		if len(rawList) > 0 {
			ddosPPlan = rawList[0].(map[string]interface{})
		}

		if v, ok := ddosPPlan["id"]; ok {
			id := v.(string)
			properties.DdosProtectionPlan = &network.SubResource{
				ID: &id,
			}
		}

		if v, ok := ddosPPlan["enable"]; ok {
			enable := v.(bool)
			properties.EnableDdosProtection = &enable
		}
	}

	if v, ok := d.GetOk("bgp_community"); ok {
		properties.BgpCommunities = &network.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: utils.String(v.(string))}
	}

	return properties, nil
}

func flattenVirtualNetworkDDoSProtectionPlan(input *network.VirtualNetworkPropertiesFormat) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.DdosProtectionPlan == nil || input.DdosProtectionPlan.ID == nil || input.EnableDdosProtection == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id":     *input.DdosProtectionPlan.ID,
			"enable": *input.EnableDdosProtection,
		},
	}
}

func flattenVirtualNetworkSubnets(input *[]network.Subnet) *schema.Set {
	results := &schema.Set{
		F: resourceAzureSubnetHash,
	}

	if subnets := input; subnets != nil {
		for _, subnet := range *input {
			output := map[string]interface{}{}

			if id := subnet.ID; id != nil {
				output["id"] = *id
			}

			if name := subnet.Name; name != nil {
				output["name"] = *name
			}

			if props := subnet.SubnetPropertiesFormat; props != nil {
				if prefixes := props.AddressPrefixes; prefixes != nil {
					var addressPrefixes []interface{}
					for _, prefix := range *prefixes {
						addressPrefixes = append(addressPrefixes, prefix)
					}
					output["address_prefixes"] = addressPrefixes
				}

				if prefix := props.AddressPrefix; prefix != nil {
					output["address_prefix"] = *prefix
				}

				if nsg := props.NetworkSecurityGroup; nsg != nil {
					if nsg.ID != nil {
						output["security_group"] = *nsg.ID
					}
				}
			}

			results.Add(output)
		}
	}

	return results
}

func flattenVirtualNetworkDNSServers(input *network.DhcpOptions) []string {
	results := make([]string, 0)

	if input != nil {
		if servers := input.DNSServers; servers != nil {
			results = *servers
		}
	}

	return results
}

func resourceAzureSubnetHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
		if v, ok := m["address_prefixes"]; ok {
			for _, a := range v.([]interface{}) {
				buf.WriteString(a.(string))
			}
		}
		if v, ok := m["address_prefix"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["security_group"]; ok {
			buf.WriteString(v.(string))
		}
	}

	return schema.HashString(buf.String())
}

func getExistingSubnet(ctx context.Context, resGroup string, vnetName string, subnetName string, meta interface{}) (*network.Subnet, error) {
	subnetClient := meta.(*clients.Client).Network.SubnetsClient
	resp, err := subnetClient.Get(ctx, resGroup, vnetName, subnetName, "")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return &network.Subnet{}, nil
		}
		// raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	// Return it directly rather than copy the fields to prevent potential uncovered properties (for example, `ServiceEndpoints` mentioned in #1619)
	return &resp, nil
}

func expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d *schema.ResourceData) ([]string, error) {
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
				parsedNsgID, err := parse.NetworkSecurityGroupID(networkSecurityGroupId)
				if err != nil {
					return nil, err
				}

				networkSecurityGroupName := parsedNsgID.Name
				if !utils.SliceContainsValue(nsgNames, networkSecurityGroupName) {
					nsgNames = append(nsgNames, networkSecurityGroupName)
				}
			}
		}
	}

	return nsgNames, nil
}
