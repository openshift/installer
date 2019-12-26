package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFirewallApplicationRuleCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallApplicationRuleCollectionCreateUpdate,
		Read:   resourceArmFirewallApplicationRuleCollectionRead,
		Update: resourceArmFirewallApplicationRuleCollectionCreateUpdate,
		Delete: resourceArmFirewallApplicationRuleCollectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureFirewallName,
			},

			"azure_firewall_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureFirewallName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallRCActionTypeAllow),
					string(network.AzureFirewallRCActionTypeDeny),
				}, false),
			},

			"rule": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"source_addresses": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"fqdn_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"target_fqdns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"protocol": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTP),
											string(network.AzureFirewallApplicationRuleProtocolTypeHTTPS),
										}, false),
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.PortNumber,
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

func resourceArmFirewallApplicationRuleCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	firewallName := d.Get("azure_firewall_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	applicationRules, err := expandArmFirewallApplicationRules(d.Get("rule").([]interface{}))
	if err != nil {
		return fmt.Errorf("Error expanding Firewall Application Rules: %+v", err)
	}

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	if firewall.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("Error retrieving Application Rule Collections (Firewall %q / Resource Group %q): `properties` was nil", firewallName, resourceGroup)
	}
	props := *firewall.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error retrieving Application Rule Collections (Firewall %q / Resource Group %q): `properties.ApplicationRuleCollections` was nil", firewallName, resourceGroup)
	}
	ruleCollections := *props.ApplicationRuleCollections

	priority := d.Get("priority").(int)
	newRuleCollection := network.AzureFirewallApplicationRuleCollection{
		Name: utils.String(name),
		AzureFirewallApplicationRuleCollectionPropertiesFormat: &network.AzureFirewallApplicationRuleCollectionPropertiesFormat{
			Action: &network.AzureFirewallRCAction{
				Type: network.AzureFirewallRCActionType(d.Get("action").(string)),
			},
			Priority: utils.Int32(int32(priority)),
			Rules:    &applicationRules,
		},
	}

	index := -1
	var id string
	for i, v := range ruleCollections {
		if v.Name == nil || v.ID == nil {
			continue
		}

		if *v.Name == name {
			index = i
			id = *v.ID
			break
		}
	}

	if !d.IsNewResource() {
		if index == -1 {
			return fmt.Errorf("Error locating Application Rule Collection %q (Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
		}

		ruleCollections[index] = newRuleCollection
	} else {
		if requireResourcesToBeImported && d.IsNewResource() {
			if index != -1 {
				return tf.ImportAsExistsError("azurerm_firewall_application_rule_collection", id)
			}
		}

		ruleCollections = append(ruleCollections, newRuleCollection)
	}

	firewall.AzureFirewallPropertiesFormat.ApplicationRuleCollections = &ruleCollections

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("Error creating/updating Application Rule Collection %q in Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Application Rule Collection %q of Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		return fmt.Errorf("Error retrieving Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	var collectionID string
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if collections := props.ApplicationRuleCollections; collections != nil {
			for _, collection := range *collections {
				if collection.Name == nil {
					continue
				}

				if *collection.Name == name {
					collectionID = *collection.ID
					break
				}
			}
		}
	}

	if collectionID == "" {
		return fmt.Errorf("Cannot find ID for Application Rule Collection %q (Azure Firewall %q / Resource Group %q)", name, firewallName, resourceGroup)
	}
	d.SetId(collectionID)

	return resourceArmFirewallApplicationRuleCollectionRead(d, meta)
}

func resourceArmFirewallApplicationRuleCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["applicationRuleCollections"]

	read, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure Firewall %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Azure Firewall %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.AzureFirewallPropertiesFormat == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	props := *read.AzureFirewallPropertiesFormat

	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	var rule *network.AzureFirewallApplicationRuleCollection
	for _, r := range *props.ApplicationRuleCollections {
		if r.Name == nil {
			continue
		}

		if *r.Name == name {
			rule = &r
			break
		}
	}

	if rule == nil {
		log.Printf("[DEBUG] Application Rule Collection %q was not found on Firewall %q (Resource Group %q) - removing from state!", name, firewallName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("azure_firewall_name", firewallName)
	d.Set("resource_group_name", resourceGroup)

	if props := rule.AzureFirewallApplicationRuleCollectionPropertiesFormat; props != nil {
		if action := props.Action; action != nil {
			d.Set("action", string(action.Type))
		}

		if priority := props.Priority; priority != nil {
			d.Set("priority", int(*priority))
		}

		flattenedRules := flattenFirewallApplicationRuleCollectionRules(props.Rules)
		if err := d.Set("rule", flattenedRules); err != nil {
			return fmt.Errorf("Error setting `rule`: %+v", err)
		}
	}

	return nil
}

func resourceArmFirewallApplicationRuleCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).azureFirewallsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	firewallName := id.Path["azureFirewalls"]
	name := id.Path["applicationRuleCollections"]

	azureRMLockByName(firewallName, azureFirewallResourceName)
	defer azureRMUnlockByName(firewallName, azureFirewallResourceName)

	firewall, err := client.Get(ctx, resourceGroup, firewallName)
	if err != nil {
		if utils.ResponseWasNotFound(firewall.Response) {
			// assume deleted
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall %q (Resource Group %q): %+v", firewallName, resourceGroup, err)
	}

	props := firewall.AzureFirewallPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props` was nil", name, firewallName, resourceGroup)
	}
	if props.ApplicationRuleCollections == nil {
		return fmt.Errorf("Error retrieving Application Rule Collection %q (Firewall %q / Resource Group %q): `props.ApplicationRuleCollections` was nil", name, firewallName, resourceGroup)
	}

	applicationRules := make([]network.AzureFirewallApplicationRuleCollection, 0)
	for _, rule := range *props.ApplicationRuleCollections {
		if rule.Name == nil {
			continue
		}

		if *rule.Name != name {
			applicationRules = append(applicationRules, rule)
		}
	}
	props.ApplicationRuleCollections = &applicationRules

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallName, firewall)
	if err != nil {
		return fmt.Errorf("Error deleting Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Application Rule Collection %q from Firewall %q (Resource Group %q): %+v", name, firewallName, resourceGroup, err)
	}

	return nil
}

func expandArmFirewallApplicationRules(inputs []interface{}) ([]network.AzureFirewallApplicationRule, error) {
	outputs := make([]network.AzureFirewallApplicationRule, 0)

	for _, input := range inputs {
		rule := input.(map[string]interface{})

		ruleName := rule["name"].(string)
		ruleDescription := rule["description"].(string)
		ruleSourceAddresses := rule["source_addresses"].(*schema.Set).List()
		ruleFqdnTags := rule["fqdn_tags"].(*schema.Set).List()
		ruleTargetFqdns := rule["target_fqdns"].(*schema.Set).List()

		output := network.AzureFirewallApplicationRule{
			Name:            utils.String(ruleName),
			Description:     utils.String(ruleDescription),
			SourceAddresses: utils.ExpandStringArray(ruleSourceAddresses),
			FqdnTags:        utils.ExpandStringArray(ruleFqdnTags),
			TargetFqdns:     utils.ExpandStringArray(ruleTargetFqdns),
		}

		ruleProtocols := make([]network.AzureFirewallApplicationRuleProtocol, 0)
		protocols := rule["protocol"].([]interface{})
		for _, v := range protocols {
			protocol := v.(map[string]interface{})
			port := protocol["port"].(int)
			ruleProtocol := network.AzureFirewallApplicationRuleProtocol{
				Port:         utils.Int32(int32(port)),
				ProtocolType: network.AzureFirewallApplicationRuleProtocolType(protocol["type"].(string)),
			}
			ruleProtocols = append(ruleProtocols, ruleProtocol)
		}
		output.Protocols = &ruleProtocols
		if len(*output.FqdnTags) > 0 {
			if len(*output.TargetFqdns) > 0 || len(*output.Protocols) > 0 {
				return outputs, fmt.Errorf("`fqdn_tags` cannot be used with `target_fqdns` or `protocol`")
			}
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func flattenFirewallApplicationRuleCollectionRules(rules *[]network.AzureFirewallApplicationRule) []map[string]interface{} {
	outputs := make([]map[string]interface{}, 0)
	if rules == nil {
		return outputs
	}

	for _, rule := range *rules {
		output := make(map[string]interface{})
		if ruleName := rule.Name; ruleName != nil {
			output["name"] = *ruleName
		}
		if ruleDescription := rule.Description; ruleDescription != nil {
			output["description"] = *ruleDescription
		}
		if ruleSourceAddresses := rule.SourceAddresses; ruleSourceAddresses != nil {
			output["source_addresses"] = set.FromStringSlice(*ruleSourceAddresses)
		}
		if ruleFqdnTags := rule.FqdnTags; ruleFqdnTags != nil {
			output["fqdn_tags"] = set.FromStringSlice(*ruleFqdnTags)
		}
		if ruleTargetFqdns := rule.TargetFqdns; ruleTargetFqdns != nil {
			output["target_fqdns"] = set.FromStringSlice(*ruleTargetFqdns)
		}
		protocols := make([]map[string]interface{}, 0)
		if ruleProtocols := rule.Protocols; ruleProtocols != nil {
			for _, p := range *ruleProtocols {
				protocol := make(map[string]interface{})
				if port := p.Port; port != nil {
					protocol["port"] = int(*port)
				}
				protocol["type"] = string(p.ProtocolType)
				protocols = append(protocols, protocol)
			}
		}
		output["protocol"] = protocols
		outputs = append(outputs, output)
	}
	return outputs
}
