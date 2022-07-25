package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func resourceNetworkingQoSMinimumBandwidthRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingQoSMinimumBandwidthRuleV2Create,
		ReadContext:   resourceNetworkingQoSMinimumBandwidthRuleV2Read,
		UpdateContext: resourceNetworkingQoSMinimumBandwidthRuleV2Update,
		DeleteContext: resourceNetworkingQoSMinimumBandwidthRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

			"min_kbps": {
				Type:     schema.TypeInt,
				Required: true,
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

func resourceNetworkingQoSMinimumBandwidthRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := rules.CreateMinimumBandwidthRuleOpts{
		MinKBps:   d.Get("min_kbps").(int),
		Direction: d.Get("direction").(string),
	}
	qosPolicyID := d.Get("qos_policy_id").(string)

	log.Printf("[DEBUG] openstack_networking_qos_minimum_bandwidth_rule_v2 create options: %#v", createOpts)
	r, err := rules.CreateMinimumBandwidthRule(networkingClient, qosPolicyID, createOpts).ExtractMinimumBandwidthRule()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_qos_minimum_bandwidth_rule_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_qos_minimum_bandwidth_rule_v2 %s to become available.", r.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingQoSMinimumBandwidthRuleV2StateRefreshFunc(networkingClient, qosPolicyID, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_minimum_bandwidth_rule_v2 %s to become available: %s", r.ID, err)
	}

	id := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, r.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_networking_qos_minimum_bandwidth_rule_v2 %s: %#v", id, r)

	return resourceNetworkingQoSMinimumBandwidthRuleV2Read(ctx, d, meta)
}

func resourceNetworkingQoSMinimumBandwidthRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_minimum_bandwidth_rule_v2 ID %s: %s", d.Id(), err)
	}

	r, err := rules.GetMinimumBandwidthRule(networkingClient, qosPolicyID, qosRuleID).ExtractMinimumBandwidthRule()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_qos_minimum_bandwidth_rule_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_minimum_bandwidth_rule_v2 %s: %#v", d.Id(), r)

	d.Set("qos_policy_id", qosPolicyID)
	d.Set("min_kbps", r.MinKBps)
	d.Set("direction", r.Direction)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingQoSMinimumBandwidthRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_minimum_bandwidth_rule_v2 ID %s: %s", d.Id(), err)
	}

	var hasChange bool
	var updateOpts rules.UpdateMinimumBandwidthRuleOpts

	if d.HasChange("min_kbps") {
		hasChange = true
		minKBps := d.Get("min_kbps").(int)
		updateOpts.MinKBps = &minKBps
	}

	if d.HasChange("direction") {
		hasChange = true
		updateOpts.Direction = d.Get("direction").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_qos_minimum_bandwidth_rule_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = rules.UpdateMinimumBandwidthRule(networkingClient, qosPolicyID, qosRuleID, updateOpts).ExtractMinimumBandwidthRule()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_qos_minimum_bandwidth_rule_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkingQoSMinimumBandwidthRuleV2Read(ctx, d, meta)
}

func resourceNetworkingQoSMinimumBandwidthRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_minimum_bandwidth_rule_v2 ID %s: %s", d.Id(), err)
	}

	if err := rules.DeleteMinimumBandwidthRule(networkingClient, qosPolicyID, qosRuleID).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_qos_minimum_bandwidth_rule_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingQoSMinimumBandwidthRuleV2StateRefreshFunc(networkingClient, qosPolicyID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_minimum_bandwidth_rule_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
