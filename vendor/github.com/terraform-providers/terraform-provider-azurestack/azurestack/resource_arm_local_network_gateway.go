package azurestack

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLocalNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLocalNetworkGatewayCreate,
		Read:   resourceArmLocalNetworkGatewayRead,
		Update: resourceArmLocalNetworkGatewayCreate,
		Delete: resourceArmLocalNetworkGatewayDelete,
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

			"gateway_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"bgp_peering_address": {
							Type:     schema.TypeString,
							Required: true,
						},

						"peer_weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmLocalNetworkGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).localNetConnClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	ipAddress := d.Get("gateway_address").(string)

	addressSpaces := expandLocalNetworkGatewayAddressSpaces(d)

	bgpSettings, err := expandLocalNetworkGatewayBGPSettings(d)
	if err != nil {
		return err
	}

	tags := d.Get("tags").(map[string]interface{})

	gateway := network.LocalNetworkGateway{
		Name:     &name,
		Location: &location,
		LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
			LocalNetworkAddressSpace: &network.AddressSpace{
				AddressPrefixes: &addressSpaces,
			},
			GatewayIPAddress: &ipAddress,
			BgpSettings:      bgpSettings,
		},
		Tags: *expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, gateway)
	if err != nil {
		return fmt.Errorf("Error creating Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Local Network Gateway ID %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLocalNetworkGatewayRead(d, meta)
}

func resourceArmLocalNetworkGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).localNetConnClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndLocalNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading the state of Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	if props := resp.LocalNetworkGatewayPropertiesFormat; props != nil {
		d.Set("gateway_address", props.GatewayIPAddress)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			if prefixes := lnas.AddressPrefixes; prefixes != nil {
				d.Set("address_space", *prefixes)
			}
		}
		flattenedSettings := flattenLocalNetworkGatewayBGPSettings(props.BgpSettings)
		if err := d.Set("bgp_settings", flattenedSettings); err != nil {
			return err
		}
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmLocalNetworkGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).localNetConnClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndLocalNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error issuing delete request for local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for completion of local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId string) (string, string, error) {
	id, err := parseAzureResourceID(localNetworkGatewayId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["localNetworkGateways"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func expandLocalNetworkGatewayBGPSettings(d *schema.ResourceData) (*network.BgpSettings, error) {
	v, exists := d.GetOk("bgp_settings")
	if !exists {
		return nil, nil
	}

	settings := v.([]interface{})
	setting := settings[0].(map[string]interface{})

	bgpSettings := network.BgpSettings{
		Asn:               utils.Int64(int64(setting["asn"].(int))),
		BgpPeeringAddress: utils.String(setting["bgp_peering_address"].(string)),
		PeerWeight:        utils.Int32(int32(setting["peer_weight"].(int))),
	}

	return &bgpSettings, nil
}

func expandLocalNetworkGatewayAddressSpaces(d *schema.ResourceData) []string {
	prefixes := make([]string, 0)

	for _, pref := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, pref.(string))
	}

	return prefixes
}

func flattenLocalNetworkGatewayBGPSettings(input *network.BgpSettings) []interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	if v := input.Asn; v != nil {
		output["asn"] = int(*v)
	}

	if v := input.BgpPeeringAddress; v != nil {
		output["bgp_peering_address"] = *v
	}

	if v := input.PeerWeight; v != nil {
		output["peer_weight"] = int(*v)
	}

	return []interface{}{output}
}
