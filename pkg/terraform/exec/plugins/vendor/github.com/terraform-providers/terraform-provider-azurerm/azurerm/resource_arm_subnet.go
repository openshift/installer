package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var subnetResourceName = "azurerm_subnet"

func resourceArmSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubnetCreateUpdate,
		Read:   resourceArmSubnetRead,
		Update: resourceArmSubnetCreateUpdate,
		Delete: resourceArmSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"virtual_network_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use the `address_prefixes` property instead.",
				ConflictsWith: []string{"address_prefixes"},
			},

			"address_prefixes": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"address_prefix"},
			},

			"network_security_group_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use the `azurerm_subnet_network_security_group_association` resource instead.",
			},

			"route_table_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use the `azurerm_subnet_route_table_association` resource instead.",
			},

			"ip_configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"service_endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"delegation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_delegation": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Microsoft.Batch/batchAccounts",
											"Microsoft.ContainerInstance/containerGroups",
											"Microsoft.HardwareSecurityModules/dedicatedHSMs",
											"Microsoft.Logic/integrationServiceEnvironments",
											"Microsoft.Netapp/volumes",
											"Microsoft.ServiceFabricMesh/networks",
											"Microsoft.Sql/managedInstances",
											"Microsoft.Sql/servers",
											"Microsoft.Web/serverFarms",
										}, false),
									},
									"actions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Microsoft.Network/virtualNetworks/subnets/action",
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmSubnetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Subnet creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, vnetName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Subnet %q (Virtual Network %q / Resource Group %q): %s", name, vnetName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_subnet", *existing.ID)
		}
	}

	var prefixSet bool
	properties := network.SubnetPropertiesFormat{}
	if value, ok := d.GetOk("address_prefixes"); ok {
		var addressPrefixes []string
		for _, item := range value.([]interface{}) {
			addressPrefixes = append(addressPrefixes, item.(string))
		}
		properties.AddressPrefixes = &addressPrefixes
		prefixSet = len(addressPrefixes) > 0
	}
	if value, ok := d.GetOk("address_prefix"); ok {
		addressPrefix := value.(string)
		properties.AddressPrefix = &addressPrefix
		prefixSet = len(addressPrefix) > 0
	}
	if !prefixSet {
		return fmt.Errorf("[ERROR] either address_prefix or address_prefixes is required")
	}

	azureRMLockByName(vnetName, virtualNetworkResourceName)
	defer azureRMUnlockByName(vnetName, virtualNetworkResourceName)

	if v, ok := d.GetOk("network_security_group_id"); ok {
		nsgId := v.(string)
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &nsgId,
		}

		networkSecurityGroupName, err := parseNetworkSecurityGroupName(nsgId)
		if err != nil {
			return err
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	} else {
		properties.NetworkSecurityGroup = nil
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		rtId := v.(string)
		properties.RouteTable = &network.RouteTable{
			ID: &rtId,
		}

		routeTableName, err := parseRouteTableName(rtId)
		if err != nil {
			return err
		}

		azureRMLockByName(routeTableName, routeTableResourceName)
		defer azureRMUnlockByName(routeTableName, routeTableResourceName)
	} else {
		properties.RouteTable = nil
	}

	serviceEndpoints := expandSubnetServiceEndpoints(d)
	properties.ServiceEndpoints = &serviceEndpoints

	delegations := expandSubnetDelegation(d)
	properties.Delegations = &delegations

	subnet := network.Subnet{
		Name:                   &name,
		SubnetPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, vnetName, name, subnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, vnetName, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Subnet %q (Virtual Network %q / Resource Group %q)", vnetName, name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSubnetRead(d, meta)
}

func resourceArmSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["subnets"]

	resp, err := client.Get(ctx, resGroup, vnetName, name, "")

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Subnet %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("virtual_network_name", vnetName)

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

		var securityGroupId *string
		if props.NetworkSecurityGroup != nil {
			securityGroupId = props.NetworkSecurityGroup.ID
		}
		d.Set("network_security_group_id", securityGroupId)

		var routeTableId string
		if props.RouteTable != nil && props.RouteTable.ID != nil {
			routeTableId = *props.RouteTable.ID
		}
		d.Set("route_table_id", routeTableId)

		ips := flattenSubnetIPConfigurations(props.IPConfigurations)
		if err := d.Set("ip_configurations", ips); err != nil {
			return err
		}

		serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
		if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
			return err
		}

		delegation := flattenSubnetDelegation(props.Delegations)
		if err := d.Set("delegation", delegation); err != nil {
			return fmt.Errorf("Error flattening `delegation`: %+v", err)
		}
	}

	return nil
}

func resourceArmSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["subnets"]
	vnetName := id.Path["virtualNetworks"]

	if v, ok := d.GetOk("network_security_group_id"); ok {
		networkSecurityGroupId := v.(string)
		networkSecurityGroupName, err2 := parseNetworkSecurityGroupName(networkSecurityGroupId)
		if err2 != nil {
			return err2
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		rtId := v.(string)
		routeTableName, err2 := parseRouteTableName(rtId)
		if err2 != nil {
			return err2
		}

		azureRMLockByName(routeTableName, routeTableResourceName)
		defer azureRMUnlockByName(routeTableName, routeTableResourceName)
	}

	azureRMLockByName(vnetName, virtualNetworkResourceName)
	defer azureRMUnlockByName(vnetName, virtualNetworkResourceName)

	azureRMLockByName(name, subnetResourceName)
	defer azureRMUnlockByName(name, subnetResourceName)

	future, err := client.Delete(ctx, resGroup, vnetName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion for Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	return nil
}

func expandSubnetServiceEndpoints(d *schema.ResourceData) []network.ServiceEndpointPropertiesFormat {
	serviceEndpoints := d.Get("service_endpoints").([]interface{})
	endpoints := make([]network.ServiceEndpointPropertiesFormat, 0)

	for _, svcEndpointRaw := range serviceEndpoints {
		if svc, ok := svcEndpointRaw.(string); ok {
			endpoint := network.ServiceEndpointPropertiesFormat{
				Service: &svc,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

func flattenSubnetServiceEndpoints(serviceEndpoints *[]network.ServiceEndpointPropertiesFormat) []string {
	endpoints := make([]string, 0)

	if serviceEndpoints == nil {
		return endpoints
	}

	for _, endpoint := range *serviceEndpoints {
		if endpoint.Service != nil {
			endpoints = append(endpoints, *endpoint.Service)
		}
	}

	return endpoints
}

func flattenSubnetIPConfigurations(ipConfigurations *[]network.IPConfiguration) []string {
	ips := make([]string, 0)

	if ipConfigurations != nil {
		for _, ip := range *ipConfigurations {
			ips = append(ips, *ip.ID)
		}
	}

	return ips
}

func expandSubnetDelegation(d *schema.ResourceData) []network.Delegation {
	delegations := d.Get("delegation").([]interface{})
	retDelegations := make([]network.Delegation, 0)

	for _, deleValue := range delegations {
		deleData := deleValue.(map[string]interface{})
		deleName := deleData["name"].(string)
		srvDelegations := deleData["service_delegation"].([]interface{})
		srvDelegation := srvDelegations[0].(map[string]interface{})
		srvName := srvDelegation["name"].(string)
		srvActions := srvDelegation["actions"].([]interface{})

		retSrvActions := make([]string, 0)
		for _, srvAction := range srvActions {
			srvActionData := srvAction.(string)
			retSrvActions = append(retSrvActions, srvActionData)
		}

		retDelegation := network.Delegation{
			Name: &deleName,
			ServiceDelegationPropertiesFormat: &network.ServiceDelegationPropertiesFormat{
				ServiceName: &srvName,
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return retDelegations
}

func flattenSubnetDelegation(delegations *[]network.Delegation) []interface{} {
	if delegations == nil {
		return []interface{}{}
	}

	retDeles := make([]interface{}, 0)

	for _, dele := range *delegations {
		retDele := make(map[string]interface{})
		if v := dele.Name; v != nil {
			retDele["name"] = *v
		}

		svcDeles := make([]interface{}, 0)
		svcDele := make(map[string]interface{})
		if props := dele.ServiceDelegationPropertiesFormat; props != nil {
			if v := props.ServiceName; v != nil {
				svcDele["name"] = *v
			}

			if v := props.Actions; v != nil {
				svcDele["actions"] = *v
			}
		}

		svcDeles = append(svcDeles, svcDele)

		retDele["service_delegation"] = svcDeles

		retDeles = append(retDeles, retDele)
	}

	return retDeles
}
