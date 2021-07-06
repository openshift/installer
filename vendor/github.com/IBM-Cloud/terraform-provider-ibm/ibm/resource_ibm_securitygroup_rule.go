// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSecurityGroupRuleCreate,
		Read:     resourceIBMSecurityGroupRuleRead,
		Delete:   resourceIBMSecurityGroupRuleDelete,
		Update:   resourceIBMSecurityGroupRuleUpdate,
		Exists:   resourceIBMSecurityGroupRuleExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Direction of rule: ingress or egress",
				ValidateFunc: validateSecurityRuleDirection,
			},
			"ether_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IP version IPv4 or IPv6",
				Default:      "IPv4",
				ValidateFunc: validateSecurityRuleEtherType,
			},
			"port_range_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port number minimum range",
			},
			"port_range_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port number max range",
			},
			"remote_group_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"remote_ip"},
				Description:   "remote group ID",
			},
			"remote_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"remote_group_id"},
				ValidateFunc:  validateRemoteIP,
				Description:   "Remote IP Address",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "icmp, tcp or udp",
				ValidateFunc: validateSecurityRuleProtocol,
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Security group ID",
			},
		},
	}
}

func findMatchingRule(sgID int, rule *datatypes.Network_SecurityGroup_Rule,
	service services.Network_SecurityGroup) (*datatypes.Network_SecurityGroup_Rule, error) {

	var filters []filter.Filter
	if rule.PortRangeMax != nil {
		filters = append(filters, filter.Path("rules.portRangeMax").Eq(rule.PortRangeMax))
	}
	if rule.PortRangeMin != nil {
		filters = append(filters, filter.Path("rules.portRangeMin").Eq(rule.PortRangeMin))
	}

	if rule.RemoteGroupId != nil {
		filters = append(filters, filter.Path("rules.remoteGroupId").Eq(rule.RemoteGroupId))
	}

	if rule.RemoteIp != nil {
		filters = append(filters, filter.Path("rules.remoteIp").Eq(rule.RemoteIp))
	}

	filters = append(filters, filter.Path("rules.direction").Eq(rule.Direction))

	if rule.Ethertype != nil {
		filters = append(filters, filter.Path("rules.ethertype").Eq(rule.Ethertype))
	}
	if rule.Protocol != nil {
		filters = append(filters, filter.Path("rules.protocol").Eq(rule.Protocol))
	}

	rules, err := service.Filter(filter.Build(filters...)).Id(sgID).GetRules()
	if err != nil {
		return nil, fmt.Errorf("Error fetching information for Security Group Rule: %s", err)
	}
	log.Printf("[INFO] rules %v", rules)

	if len(rules) == 0 {
		return nil, nil
	}
	return &rules[0], nil
}

func resourceIBMSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	sgID := d.Get("security_group_id").(int)

	sgrule := datatypes.Network_SecurityGroup_Rule{}

	sgrule.Direction = sl.String(d.Get("direction").(string))

	if d.Get("ether_type").(string) != "" {
		sgrule.Ethertype = sl.String(d.Get("ether_type").(string))
	}

	if _, ok := d.GetOk("port_range_min"); ok {
		sgrule.PortRangeMin = sl.Int(d.Get("port_range_min").(int))

	}

	if _, ok := d.GetOk("port_range_max"); ok {
		sgrule.PortRangeMax = sl.Int(d.Get("port_range_max").(int))
	}

	if d.Get("protocol").(string) != "" {
		sgrule.Protocol = sl.String(d.Get("protocol").(string))
	}

	if v, ok := d.GetOk("remote_group_id"); ok {
		sgrule.RemoteGroupId = sl.Int(v.(int))
	}

	if v, ok := d.GetOk("remote_ip"); ok {
		sgrule.RemoteIp = sl.String(v.(string))
	}

	// if only one of min/max is provided, set the other one to the provided
	if sgrule.PortRangeMin != nil && sgrule.PortRangeMax == nil {
		sgrule.PortRangeMax = sgrule.PortRangeMin
	}
	if sgrule.PortRangeMax != nil && sgrule.PortRangeMin == nil {
		sgrule.PortRangeMin = sgrule.PortRangeMax
	}

	matchingrule, err := findMatchingRule(sgID, &sgrule, service)
	if err != nil {
		return err
	}

	if matchingrule != nil {
		log.Printf("[INFO] rule exists")
		d.SetId(fmt.Sprintf("%d", *matchingrule.Id))
		return nil
	}

	opts := []datatypes.Network_SecurityGroup_Rule{
		sgrule,
	}
	log.Println("[INFO] creating security group rule")
	_, err = service.Id(sgID).AddRules(opts)
	if err != nil {
		return fmt.Errorf("Error creating Security Group Rule: %s", err)
	}

	matchingrule, err = findMatchingRule(sgID, &sgrule, service)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(*matchingrule.Id))

	return resourceIBMSecurityGroupRuleRead(d, meta)
}

func resourceIBMSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	sgID := d.Get("security_group_id").(int)
	matchingrules, err := service.Filter(filter.Build(
		filter.Path("rules.id").Eq(d.Id()))).Id(sgID).GetRules()
	if err != nil {
		// If the group is somehow already destroyed, mark as
		// succesfully gone
		if err, ok := err.(sl.Error); ok && err.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Security Group Rule: %s", err)
	}

	if len(matchingrules) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("direction", matchingrules[0].Direction)

	if matchingrules[0].Ethertype != nil {
		d.Set("ether_type", matchingrules[0].Ethertype)
	}
	if matchingrules[0].PortRangeMin != nil {
		d.Set("port_range_min", matchingrules[0].PortRangeMin)
	}
	if matchingrules[0].PortRangeMax != nil {
		d.Set("port_range_max", matchingrules[0].PortRangeMax)
	}
	if matchingrules[0].Protocol != nil {
		d.Set("protocol", matchingrules[0].Protocol)
	}

	if matchingrules[0].RemoteGroupId != nil {
		d.Set("remote_group_id", matchingrules[0].RemoteGroupId)
	}
	if matchingrules[0].RemoteIp != nil {
		d.Set("remote_ip", matchingrules[0].RemoteIp)
	}
	return nil
}

func resourceIBMSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)
	securityGroupID := d.Get("security_group_id").(int)
	matchingrules, err := service.Filter(filter.Build(
		filter.Path("rules.id").Eq(d.Id()))).Id(securityGroupID).GetRules()
	if err != nil {
		return fmt.Errorf("Error retrieving Security Group Rule: %s", err)
	}
	if d.HasChange("direction") {
		matchingrules[0].Direction = sl.String(d.Get("direction").(string))
	}
	if d.HasChange("ether_type") {
		matchingrules[0].Ethertype = sl.String(d.Get("ether_type").(string))
	}
	if d.HasChange("port_range_min") {
		matchingrules[0].PortRangeMin = sl.Int(d.Get("port_range_min").(int))
	}
	if d.HasChange("port_range_max") {
		matchingrules[0].PortRangeMax = sl.Int(d.Get("port_range_max").(int))
	}
	if d.HasChange("protocol") {
		matchingrules[0].Protocol = sl.String(d.Get("protocol").(string))
	}
	if d.HasChange("remote_group_ip") {
		matchingrules[0].RemoteGroupId = sl.Int(d.Get("remote_group_ip").(int))
	}
	if d.HasChange("remote_ip") {
		matchingrules[0].RemoteIp = sl.String(d.Get("remote_ip").(string))
	}
	_, err = service.Id(securityGroupID).EditRules([]datatypes.Network_SecurityGroup_Rule{matchingrules[0]})
	if err != nil {
		return fmt.Errorf("Couldn't update Security Group Rule: %s", err)
	}
	return resourceIBMSecurityGroupRuleRead(d, meta)
}

func resourceIBMSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)
	sgID := d.Get("security_group_id").(int)
	id, _ := strconv.Atoi(d.Id())
	_, err := service.Id(sgID).RemoveRules([]int{id})
	if err != nil {
		if err, ok := err.(sl.Error); ok && err.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting Security Group Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMSecurityGroupRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	sgID := d.Get("security_group_id").(int)
	matchingrules, err := service.Filter(filter.Build(
		filter.Path("rules.id").Eq(d.Id()))).Id(sgID).GetRules()
	if err != nil {
		// If the group is somehow already destroyed, mark as
		// succesfully gone
		if err, ok := err.(sl.Error); ok && err.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving Security Group Rule: %s", err)
	}

	if len(matchingrules) == 0 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
