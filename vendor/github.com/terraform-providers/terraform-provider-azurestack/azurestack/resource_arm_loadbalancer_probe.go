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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/terraform-providers/terraform-provider-azurestack/azurestack/helpers/azure"
)

func resourceArmLoadBalancerProbe() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerProbeCreateUpdate,
		Read:   resourceArmLoadBalancerProbeRead,
		Update: resourceArmLoadBalancerProbeCreateUpdate,
		Delete: resourceArmLoadBalancerProbeDelete,
		Importer: &schema.ResourceImporter{
			State: loadBalancerSubResourceStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"protocol": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				StateFunc:        ignoreCaseStateFunc,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ProbeProtocolHTTP),
					string(network.ProbeProtocolTCP),
				}, true),
			},

			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"request_path": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"interval_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(5),
			},

			"number_of_probes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"load_balancer_rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerProbeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[INFO] LoadBalancer %q not found. Removing from state", d.Get("name").(string))
		return nil
	}

	newProbe := expandAzureRmLoadBalancerProbe(d)
	probes := append(*loadBalancer.LoadBalancerPropertiesFormat.Probes, *newProbe)

	existingProbe, existingProbeIndex, exists := findLoadBalancerProbeByName(loadBalancer, d.Get("name").(string))
	if exists {
		if d.Get("name").(string) == *existingProbe.Name {
			// this probe is being updated/reapplied remove old copy from the slice
			probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.Probes = &probes
	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("Error Getting LoadBalancer Name and Group: %+v", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read LoadBalancer %q (resource group %q) ID", loadBalancerName, resGroup)
	}

	var createdProbeId string
	for _, Probe := range *(*read.LoadBalancerPropertiesFormat).Probes {
		if *Probe.Name == d.Get("name").(string) {
			createdProbeId = *Probe.ID
		}
	}

	if createdProbeId == "" {
		return fmt.Errorf("Cannot find created LoadBalancer Probe ID %q", createdProbeId)
	}

	d.SetId(createdProbeId)

	log.Printf("[DEBUG] Waiting for LoadBalancer (%s) to become available", loadBalancerName)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: loadbalancerStateRefreshFunc(ctx, client, resGroup, loadBalancerName),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for LoadBalancer (%q - Resource Group %q) to become available: %+v", loadBalancerName, resGroup, err)
	}

	return resourceArmLoadBalancerProbeRead(d, meta)
}

func resourceArmLoadBalancerProbeRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["probes"]

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Get("loadbalancer_id").(string), meta)
	if err != nil {
		return errwrap.Wrapf("Error Getting LoadBalancer By ID {{err}}", err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer %q not found. Removing from state", name)
		return nil
	}

	config, _, exists := findLoadBalancerProbeByName(loadBalancer, name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer Probe %q not found. Removing from state", name)
		return nil
	}

	d.Set("name", config.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := config.ProbePropertiesFormat; properties != nil {
		d.Set("protocol", properties.Protocol)
		d.Set("interval_in_seconds", properties.IntervalInSeconds)
		d.Set("number_of_probes", properties.NumberOfProbes)
		d.Set("port", properties.Port)
		d.Set("request_path", properties.RequestPath)

		var loadBalancerRules []string
		if rules := properties.LoadBalancingRules; rules != nil {
			for _, ruleConfig := range *rules {
				loadBalancerRules = append(loadBalancerRules, *ruleConfig.ID)
			}
		}
		d.Set("load_balancer_rules", loadBalancerRules)
	}

	return nil
}

func resourceArmLoadBalancerProbeDelete(d *schema.ResourceData, meta interface{}) error {
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

	_, index, exists := findLoadBalancerProbeByName(loadBalancer, d.Get("name").(string))
	if !exists {
		return nil
	}

	oldProbes := *loadBalancer.LoadBalancerPropertiesFormat.Probes
	newProbes := append(oldProbes[:index], oldProbes[index+1:]...)
	loadBalancer.LoadBalancerPropertiesFormat.Probes = &newProbes

	resGroup, loadBalancerName, err := resourceGroupAndLBNameFromId(d.Get("loadbalancer_id").(string))
	if err != nil {
		return errwrap.Wrapf("Error Getting LoadBalancer Name and Group: {{err}}", err)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, loadBalancerName, *loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of LoadBalancer %q (Resource Group %q): %+v", loadBalancerName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, loadBalancerName, "")
	if err != nil {
		return errwrap.Wrapf("Error Getting LoadBalancer {{err}}", err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read LoadBalancer %s (resource group %s) ID", loadBalancerName, resGroup)
	}

	return nil
}

func expandAzureRmLoadBalancerProbe(d *schema.ResourceData) *network.Probe {

	properties := network.ProbePropertiesFormat{
		NumberOfProbes:    utils.Int32(int32(d.Get("number_of_probes").(int))),
		IntervalInSeconds: utils.Int32(int32(d.Get("interval_in_seconds").(int))),
		Port:              utils.Int32(int32(d.Get("port").(int))),
	}

	if v, ok := d.GetOk("protocol"); ok {
		properties.Protocol = network.ProbeProtocol(v.(string))
	}

	if v, ok := d.GetOk("request_path"); ok {
		properties.RequestPath = utils.String(v.(string))
	}

	return &network.Probe{
		Name:                  utils.String(d.Get("name").(string)),
		ProbePropertiesFormat: &properties,
	}
}
