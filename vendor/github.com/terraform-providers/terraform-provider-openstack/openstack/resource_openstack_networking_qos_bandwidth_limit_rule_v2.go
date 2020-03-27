package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func resourceNetworkingQoSBandwidthLimitRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingQoSBandwidthLimitRuleV2Create,
		Read:   resourceNetworkingQoSBandwidthLimitRuleV2Read,
		Update: resourceNetworkingQoSBandwidthLimitRuleV2Update,
		Delete: resourceNetworkingQoSBandwidthLimitRuleV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"qos_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"max_kbps": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"max_burst_kbps": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			"direction": {
				Type:     schema.TypeString,
				Default:  "egress",
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceNetworkingQoSBandwidthLimitRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := rules.CreateBandwidthLimitRuleOpts{
		MaxKBps:      d.Get("max_kbps").(int),
		MaxBurstKBps: d.Get("max_burst_kbps").(int),
		Direction:    d.Get("direction").(string),
	}
	qosPolicyID := d.Get("qos_policy_id").(string)

	log.Printf("[DEBUG] openstack_networking_qos_bandwidth_limit_rule_v2 create options: %#v", createOpts)
	r, err := rules.CreateBandwidthLimitRule(networkingClient, qosPolicyID, createOpts).ExtractBandwidthLimitRule()
	if err != nil {
		return fmt.Errorf("Error creating openstack_networking_qos_bandwidth_limit_rule_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_qos_bandwidth_limit_rule_v2 %s to become available.", r.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingQoSBandwidthLimitRuleV2StateRefreshFunc(networkingClient, qosPolicyID, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for openstack_networking_qos_bandwidth_limit_rule_v2 %s to become available: %s", r.ID, err)
	}

	id := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, r.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_networking_qos_bandwidth_limit_rule_v2 %s: %#v", id, r)

	return resourceNetworkingQoSBandwidthLimitRuleV2Read(d, meta)
}

func resourceNetworkingQoSBandwidthLimitRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading openstack_networking_qos_bandwidth_limit_rule_v2 ID %s: %s", d.Id(), err)
	}

	r, err := rules.GetBandwidthLimitRule(networkingClient, qosPolicyID, qosRuleID).ExtractBandwidthLimitRule()
	if err != nil {
		return CheckDeleted(d, err, "Error getting openstack_networking_qos_bandwidth_limit_rule_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_bandwidth_limit_rule_v2 %s: %#v", d.Id(), r)

	d.Set("qos_policy_id", qosPolicyID)
	d.Set("max_kbps", r.MaxKBps)
	d.Set("max_burst_kbps", r.MaxBurstKBps)
	d.Set("direction", r.Direction)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingQoSBandwidthLimitRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading openstack_networking_qos_bandwidth_limit_rule_v2 ID %s: %s", d.Id(), err)
	}

	var hasChange bool
	var updateOpts rules.UpdateBandwidthLimitRuleOpts

	if d.HasChange("max_kbps") {
		hasChange = true
		maxKBps := d.Get("max_kbps").(int)
		updateOpts.MaxKBps = &maxKBps
	}

	if d.HasChange("max_burst_kbps") {
		hasChange = true
		maxBurstKBps := d.Get("max_burst_kbps").(int)
		updateOpts.MaxBurstKBps = &maxBurstKBps
	}

	if d.HasChange("direction") {
		hasChange = true
		updateOpts.Direction = d.Get("direction").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_qos_bandwidth_limit_rule_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = rules.UpdateBandwidthLimitRule(networkingClient, qosPolicyID, qosRuleID, updateOpts).ExtractBandwidthLimitRule()
		if err != nil {
			return fmt.Errorf("Error updating openstack_networking_qos_bandwidth_limit_rule_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkingQoSBandwidthLimitRuleV2Read(d, meta)
}

func resourceNetworkingQoSBandwidthLimitRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading openstack_networking_qos_bandwidth_limit_rule_v2 ID %s: %s", d.Id(), err)
	}

	if err := rules.DeleteBandwidthLimitRule(networkingClient, qosPolicyID, qosRuleID).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "Error getting openstack_networking_qos_bandwidth_limit_rule_v2")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingQoSBandwidthLimitRuleV2StateRefreshFunc(networkingClient, qosPolicyID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for openstack_networking_qos_bandwidth_limit_rule_v2 %s to delete: %s", d.Id(), err)
	}

	return nil
}
