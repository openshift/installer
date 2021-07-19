package azurestack

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-10-01/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancerNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerNatRuleCreateUpdate,
		Read:   resourceArmLoadBalancerNatRuleRead,
		Update: resourceArmLoadBalancerNatRuleCreateUpdate,
		Delete: resourceArmLoadBalancerNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        ignoreCaseStateFunc,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.TransportProtocolTCP),
					string(network.TransportProtocolUDP),
				}, true),
			},

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"frontend_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"backend_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_ip_configuration_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"frontend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_ip_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLoadBalancerNatRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	loadBalancerID := d.Get("loadbalancer_id").(string)
	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return errwrap.Wrapf("Error Getting LoadBalancer By ID {{err}}", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer %q not found. Removing from state", d.Get("name").(string))
		return nil
	}

	newNatRule, err := expandAzureRmLoadBalancerNatRule(d, loadBalancer)
	if err != nil {
		return errwrap.Wrapf("Error Expanding NAT Rule {{err}}", err)
	}

	natRules := append(*loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules, *newNatRule)

	existingNatRule, existingNatRuleIndex, exists := findLoadBalancerNatRuleByName(loadBalancer, d.Get("name").(string))
	if exists {
		if d.Get("name").(string) == *existingNatRule.Name {
			// this nat rule is being updated/reapplied remove old copy from the slice
			natRules = append(natRules[:existingNatRuleIndex], natRules[existingNatRuleIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &natRules
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating / Updating LoadBalancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving LoadBalancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read LoadBalancer %q (Resource Group %q) ID", loadBalancerName, resGroup)
	}

	var natRuleId string
	for _, InboundNatRule := range *(*read.LoadBalancerPropertiesFormat).InboundNatRules {
		if *InboundNatRule.Name == d.Get("name").(string) {
			natRuleId = *InboundNatRule.ID
		}
	}

	if natRuleId != "" {
		d.SetId(natRuleId)
	} else {
		return fmt.Errorf("Cannot find created LoadBalancer NAT Rule ID %q", natRuleId)
	}

	// TODO: is this still needed?
	log.Printf("[DEBUG] Waiting for LoadBalancer (%s) to become available", loadBalancerName)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: loadbalancerStateRefreshFunc(ctx, client, resGroup, loadBalancerName),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for LoadBalancer (%s) to become available: %s", loadBalancerName, err)
	}

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func resourceArmLoadBalancerNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["inboundNatRules"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := findLoadBalancerNatRuleByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer Nat Rule %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := config.InboundNatRulePropertiesFormat; props != nil {
		d.Set("protocol", props.Protocol)
		d.Set("frontend_port", props.FrontendPort)
		d.Set("backend_port", props.BackendPort)
		d.Set("enable_floating_ip", props.EnableFloatingIP)

		if ipconfiguration := props.FrontendIPConfiguration; ipconfiguration != nil {
			fipID, err := parseAzureResourceID(*ipconfiguration.ID)
			if err != nil {
				return err
			}

			d.Set("frontend_ip_configuration_name", fipID.Path["frontendIPConfigurations"])
			d.Set("frontend_ip_configuration_id", ipconfiguration.ID)
		}

		if ipconfiguration := props.BackendIPConfiguration; ipconfiguration != nil {
			d.Set("backend_ip_configuration_id", ipconfiguration.ID)
		}
	}

	return nil
}

func resourceArmLoadBalancerNatRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	loadBalancerID := d.Get("loadbalancer_id").(string)
	armMutexKV.Lock(loadBalancerID)
	defer armMutexKV.Unlock(loadBalancerID)

	loadBalancer, exists, err := retrieveLoadBalancerById(loadBalancerID, meta)
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer By ID: %+v", err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	_, index, exists := findLoadBalancerNatRuleByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldNatRules := *loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules
	newNatRules := append(oldNatRules[:index], oldNatRules[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.InboundNatRules = &newNatRules

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer %q (Resource Group %q) %+v", loadBalancerName, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the completion of LoadBalancer updates for %q (Resource Group %q) %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving LoadBalancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read LoadBalancer %q (resource group %q) ID", loadBalancerName, resGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerNatRule(d *schema.ResourceData, lb *network.LoadBalancer) (*network.InboundNatRule, error) {

	properties := network.InboundNatRulePropertiesFormat{
		Protocol:     network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort: utils.Int32(int32(d.Get("frontend_port").(int))),
		BackendPort:  utils.Int32(int32(d.Get("backend_port").(int))),
	}

	if v, ok := d.GetOk("enable_floating_ip"); ok {
		properties.EnableFloatingIP = utils.Bool(v.(bool))
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := findLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &network.SubResource{
			ID: rule.ID,
		}
	}

	natRule := network.InboundNatRule{
		Name:                           utils.String(d.Get("name").(string)),
		InboundNatRulePropertiesFormat: &properties,
	}

	return &natRule, nil
}
