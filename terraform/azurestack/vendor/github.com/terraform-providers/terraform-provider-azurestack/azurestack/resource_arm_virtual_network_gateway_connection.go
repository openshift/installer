package azurestack

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/terraform-providers/terraform-provider-azurestack/azurestack/helpers/azure"
)

func resourceArmVirtualNetworkGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkGatewayConnectionCreateUpdate,
		Read:   resourceArmVirtualNetworkGatewayConnectionRead,
		Update: resourceArmVirtualNetworkGatewayConnectionCreateUpdate,
		Delete: resourceArmVirtualNetworkGatewayConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ExpressRoute),
					string(network.IPsec),
					string(network.Vnet2Vnet),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"virtual_network_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"authorization_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},

			"express_route_circuit_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"peer_virtual_network_gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"local_network_gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"enable_bgp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"routing_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkGatewayConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureStack Virtual Network Gateway Connection creation.")

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties, err := getArmVirtualNetworkGatewayConnectionProperties(d)
	if err != nil {
		return err
	}

	connection := network.VirtualNetworkGatewayConnection{
		Name:     &name,
		Location: &location,
		Tags:     *expandTags(tags),
		VirtualNetworkGatewayConnectionPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, connection)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AzureStack Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureStack Virtual Network Gateway Connection %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceArmVirtualNetworkGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureStack Virtual Network Gateway Connection %q: %+v", name, err)
	}

	conn := *resp.VirtualNetworkGatewayConnectionPropertiesFormat

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	if string(conn.ConnectionType) != "" {
		d.Set("type", string(conn.ConnectionType))
	}

	if conn.VirtualNetworkGateway1 != nil {
		d.Set("virtual_network_gateway_id", conn.VirtualNetworkGateway1.ID)
	}

	if conn.AuthorizationKey != nil {
		d.Set("authorization_key", conn.AuthorizationKey)
	}

	if conn.Peer != nil {
		d.Set("express_route_circuit_id", conn.Peer.ID)
	}

	if conn.VirtualNetworkGateway2 != nil {
		d.Set("peer_virtual_network_gateway_id", conn.VirtualNetworkGateway2.ID)
	}

	if conn.LocalNetworkGateway2 != nil {
		d.Set("local_network_gateway_id", conn.LocalNetworkGateway2.ID)
	}

	if conn.EnableBgp != nil {
		d.Set("enable_bgp", conn.EnableBgp)
	}

	if conn.RoutingWeight != nil {
		d.Set("routing_weight", conn.RoutingWeight)
	}

	// Get Shared Key
	sharedKeyResp, err := client.GetSharedKey(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making GetSharedKey request on AzureStack Virtual Network Gateway Connection %q: %+v", name, err)
	}

	if sharedKeyResp.Value != nil {
		d.Set("shared_key", sharedKeyResp.Value)
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmVirtualNetworkGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Deleting Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getArmVirtualNetworkGatewayConnectionProperties(d *schema.ResourceData) (*network.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := network.VirtualNetworkGatewayConnectionType(d.Get("type").(string))

	props := &network.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType: connectionType,
	}

	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndVirtualNetworkGatewayFromId(virtualNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("Error Getting VirtualNetworkGateway Name and Group:: %+v", err)
		}

		props.VirtualNetworkGateway1 = &network.VirtualNetworkGateway{
			ID:   &virtualNetworkGatewayId,
			Name: &name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		authorizationKey := v.(string)
		props.AuthorizationKey = &authorizationKey
	}

	if v, ok := d.GetOk("express_route_circuit_id"); ok {
		expressRouteCircuitId := v.(string)
		props.Peer = &network.SubResource{
			ID: &expressRouteCircuitId,
		}
	}

	if v, ok := d.GetOk("peer_virtual_network_gateway_id"); ok {
		peerVirtualNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndVirtualNetworkGatewayFromId(peerVirtualNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("Error Getting VirtualNetworkGateway Name and Group:: %+v", err)
		}

		props.VirtualNetworkGateway2 = &network.VirtualNetworkGateway{
			ID:   &peerVirtualNetworkGatewayId,
			Name: &name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("local_network_gateway_id"); ok {
		localNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("Error Getting LocalNetworkGateway Name and Group:: %+v", err)
		}

		props.LocalNetworkGateway2 = &network.LocalNetworkGateway{
			ID:   &localNetworkGatewayId,
			Name: &name,
			LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &network.AddressSpace{},
			},
		}
	}

	props.EnableBgp = utils.Bool(d.Get("enable_bgp").(bool))

	if v, ok := d.GetOk("routing_weight"); ok {
		routingWeight := int32(v.(int))
		props.RoutingWeight = &routingWeight
	}

	if v, ok := d.GetOk("shared_key"); ok {
		props.SharedKey = utils.String(v.(string))
	}

	if props.ConnectionType == network.ExpressRoute {
		if props.Peer == nil || props.Peer.ID == nil {
			return nil, fmt.Errorf("`express_route_circuit_id` must be specified when `type` is set to `ExpressRoute")
		}
	}

	if props.ConnectionType == network.IPsec {
		if props.LocalNetworkGateway2 == nil || props.LocalNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` and `shared_key` must be specified when `type` is set to `IPsec")
		}

		if props.SharedKey == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` and `shared_key` must be specified when `type` is set to `IPsec")
		}
	}

	if props.ConnectionType == network.Vnet2Vnet {
		if props.VirtualNetworkGateway2 == nil || props.VirtualNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`peer_virtual_network_gateway_id` and `shared_key` must be specified when `type` is set to `Vnet2Vnet")
		}
	}

	return props, nil
}

func resourceGroupAndVirtualNetworkGatewayConnectionFromId(virtualNetworkGatewayConnectionId string) (string, string, error) {
	id, err := parseAzureResourceID(virtualNetworkGatewayConnectionId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["connections"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}
