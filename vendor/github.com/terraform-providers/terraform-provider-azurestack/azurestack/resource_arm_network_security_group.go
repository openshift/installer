package azurestack

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var networkSecurityGroupResourceName = "azurestack_network_security_group"

func resourceArmNetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkSecurityGroupCreate,
		Read:   resourceArmNetworkSecurityGroupRead,
		Update: resourceArmNetworkSecurityGroupCreate,
		Delete: resourceArmNetworkSecurityGroupDelete,
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

			"security_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 140),
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"*",
								"Tcp",
								"Udp",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"source_port_range": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// The Following attributes are not included in the profile for
						// Azure Stack
						// destination_port_ranges
						// source_address_prefixes
						// source_application_security_group_ids
						// destination_address_prefixes
						// destination_application_security_group_ids

						// "source_port_ranges": {
						// 	Type:     schema.TypeString,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },

						"destination_port_range": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// "destination_port_ranges": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },

						"source_address_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// "source_address_prefixes": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },

						"destination_address_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// "destination_address_prefixes": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },

						// "destination_application_security_group_ids": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },
						//
						// "source_application_security_group_ids": {
						// 	Type:     schema.TypeSet,
						// 	Optional: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Set:      schema.HashString,
						// },

						// Constants not defined in the profile 2017-03-09
						"access": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Allow",
								"Deny",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 4096),
						},

						// Constants not defined in the profile 2017-03-09
						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Inbound",
								"Outbound",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sgRules, sgErr := expandAzureStackSecurityRules(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Network Security Group Rules: %+v", sgErr)
	}

	azureStackLockByName(name, networkSecurityGroupResourceName)
	defer azureStackUnlockByName(name, networkSecurityGroupResourceName)

	sg := network.SecurityGroup{
		Name:     &name,
		Location: &location,
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			SecurityRules: &sgRules,
		},
		Tags: *expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, sg)
	if err != nil {
		return fmt.Errorf("Error creating/updating NSG %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the completion of NSG %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read NSG %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkSecurityGroupRead(d, meta)
}

func resourceArmNetworkSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkSecurityGroups"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Network Security Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}

	if props := resp.SecurityGroupPropertiesFormat; props != nil {
		flattenedRules := flattenNetworkSecurityRules(props.SecurityRules)
		if err := d.Set("security_rule", flattenedRules); err != nil {
			return fmt.Errorf("Error flattening `security_rule`: %+v", err)
		}
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmNetworkSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkSecurityGroups"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Network Security Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error deleting Network Security Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return err
}

func expandAzureStackSecurityRules(d *schema.ResourceData) ([]network.SecurityRule, error) {
	sgRules := d.Get("security_rule").(*schema.Set).List()
	rules := make([]network.SecurityRule, 0)

	for _, sgRaw := range sgRules {
		sgRule := sgRaw.(map[string]interface{})

		// if err := validateSecurityRule(sgRule); err != nil {
		// 	return nil, err
		// }

		name := sgRule["name"].(string)
		source_port_range := sgRule["source_port_range"].(string)
		destination_port_range := sgRule["destination_port_range"].(string)
		source_address_prefix := sgRule["source_address_prefix"].(string)
		destination_address_prefix := sgRule["destination_address_prefix"].(string)
		priority := int32(sgRule["priority"].(int))
		access := sgRule["access"].(string)
		direction := sgRule["direction"].(string)
		protocol := sgRule["protocol"].(string)

		properties := network.SecurityRulePropertiesFormat{
			SourcePortRange:          &source_port_range,
			DestinationPortRange:     &destination_port_range,
			SourceAddressPrefix:      &source_address_prefix,
			DestinationAddressPrefix: &destination_address_prefix,
			Priority:                 &priority,
			Access:                   network.SecurityRuleAccess(access),
			Direction:                network.SecurityRuleDirection(direction),
			Protocol:                 network.SecurityRuleProtocol(protocol),
		}

		if v := sgRule["description"].(string); v != "" {
			properties.Description = &v
		}

		// Source and destination port ranges are not supported by the profile

		// if r, ok := sgRule["source_port_ranges"].(*schema.Set); ok && r.Len() > 0 {
		// 	var sourcePortRanges []string
		// 	for _, v := range r.List() {
		// 		s := v.(string)
		// 		sourcePortRanges = append(sourcePortRanges, s)
		// 	}
		// 	properties.SourcePortRanges = &sourcePortRanges
		// }

		// if r, ok := sgRule["destination_port_ranges"].(*schema.Set); ok && r.Len() > 0 {
		// 	var destinationPortRanges []string
		// 	for _, v := range r.List() {
		// 		s := v.(string)
		// 		destinationPortRanges = append(destinationPortRanges, s)
		// 	}
		// properties.DestinationPortRanges = &destinationPortRanges
		// }

		// if r, ok := sgRule["source_address_prefixes"].(*schema.Set); ok && r.Len() > 0 {
		// 	var sourceAddressPrefixes []string
		// 	for _, v := range r.List() {
		// 		s := v.(string)
		// 		sourceAddressPrefixes = append(sourceAddressPrefixes, s)
		// 	}
		// 	// properties.SourceAddressPrefixes = &sourceAddressPrefixes
		// }

		// if r, ok := sgRule["destination_address_prefixes"].(*schema.Set); ok && r.Len() > 0 {
		// 	var destinationAddressPrefixes []string
		// 	for _, v := range r.List() {
		// 		s := v.(string)
		// 		destinationAddressPrefixes = append(destinationAddressPrefixes, s)
		// 	}
		// 	// properties.DestinationAddressPrefixes = &destinationAddressPrefixes
		// }

		// if r, ok := sgRule["source_application_security_group_ids"].(*schema.Set); ok && r.Len() > 0 {
		// 	var sourceApplicationSecurityGroups []network.ApplicationSecurityGroup
		// 	for _, v := range r.List() {
		// 		sg := network.ApplicationSecurityGroup{
		// 			ID: utils.String(v.(string)),
		// 		}
		// 		sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		// 	}
		// 	properties.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
		// }

		// if r, ok := sgRule["destination_application_security_group_ids"].(*schema.Set); ok && r.Len() > 0 {
		// 	var destinationApplicationSecurityGroups []network.ApplicationSecurityGroup
		// 	for _, v := range r.List() {
		// 		sg := network.ApplicationSecurityGroup{
		// 			ID: utils.String(v.(string)),
		// 		}
		// 		destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
		// 	}
		// 	properties.DestinationApplicationSecurityGroups = &destinationApplicationSecurityGroups
		// }

		rules = append(rules, network.SecurityRule{
			Name:                         &name,
			SecurityRulePropertiesFormat: &properties,
		})
	}

	return rules, nil
}

func flattenNetworkSecurityRules(rules *[]network.SecurityRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if rules != nil {
		for _, rule := range *rules {
			sgRule := make(map[string]interface{})
			sgRule["name"] = *rule.Name

			if props := rule.SecurityRulePropertiesFormat; props != nil {
				if props.Description != nil {
					sgRule["description"] = *props.Description
				}

				if props.DestinationAddressPrefix != nil {
					sgRule["destination_address_prefix"] = *props.DestinationAddressPrefix
				}
				// if props.DestinationAddressPrefixes != nil {
				// 	sgRule["destination_address_prefixes"] = sliceToSet(*props.DestinationAddressPrefixes)
				// }
				if props.DestinationPortRange != nil {
					sgRule["destination_port_range"] = *props.DestinationPortRange
				}
				// if props.DestinationPortRanges != nil {
				// 	sgRule["destination_port_ranges"] = sliceToSet(*props.DestinationPortRanges)
				// }

				// destinationApplicationSecurityGroups := make([]string, 0)
				// if props.DestinationApplicationSecurityGroups != nil {
				// 	for _, g := range *props.DestinationApplicationSecurityGroups {
				// 		destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, *g.ID)
				// 	}
				// }
				// sgRule["destination_application_security_group_ids"] = sliceToSet(destinationApplicationSecurityGroups)

				if props.SourceAddressPrefix != nil {
					sgRule["source_address_prefix"] = *props.SourceAddressPrefix
				}
				// if props.SourceAddressPrefixes != nil {
				// 	sgRule["source_address_prefixes"] = sliceToSet(*props.SourceAddressPrefixes)
				// }

				// sourceApplicationSecurityGroups := make([]string, 0)
				// if props.SourceApplicationSecurityGroups != nil {
				// 	for _, g := range *props.SourceApplicationSecurityGroups {
				// 		sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, *g.ID)
				// 	}
				// }
				// sgRule["source_application_security_group_ids"] = sliceToSet(sourceApplicationSecurityGroups)

				if props.SourcePortRange != nil {
					sgRule["source_port_range"] = *props.SourcePortRange
				}
				// if props.SourcePortRanges != nil {
				// 	sgRule["source_port_ranges"] = sliceToSet(*props.SourcePortRanges)
				// }

				sgRule["protocol"] = string(props.Protocol)
				sgRule["priority"] = int(*props.Priority)
				sgRule["access"] = string(props.Access)
				sgRule["direction"] = string(props.Direction)
			}

			result = append(result, sgRule)
		}
	}

	return result
}
