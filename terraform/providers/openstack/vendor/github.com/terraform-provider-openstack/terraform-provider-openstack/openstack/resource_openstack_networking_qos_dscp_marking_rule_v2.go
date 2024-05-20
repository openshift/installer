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

func resourceNetworkingQoSDSCPMarkingRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingQoSDSCPMarkingRuleV2Create,
		ReadContext:   resourceNetworkingQoSDSCPMarkingRuleV2Read,
		UpdateContext: resourceNetworkingQoSDSCPMarkingRuleV2Update,
		DeleteContext: resourceNetworkingQoSDSCPMarkingRuleV2Delete,
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

			"dscp_mark": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceNetworkingQoSDSCPMarkingRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := rules.CreateDSCPMarkingRuleOpts{
		DSCPMark: d.Get("dscp_mark").(int),
	}
	qosPolicyID := d.Get("qos_policy_id").(string)

	log.Printf("[DEBUG] openstack_networking_qos_dscp_marking_rule_v2 create options: %#v", createOpts)
	r, err := rules.CreateDSCPMarkingRule(networkingClient, qosPolicyID, createOpts).ExtractDSCPMarkingRule()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_qos_dscp_marking_rule_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_qos_dscp_marking_rule_v2 %s to become available.", r.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingQoSDSCPMarkingRuleV2StateRefreshFunc(networkingClient, qosPolicyID, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_dscp_marking_rule_v2 %s to become available: %s", r.ID, err)
	}

	id := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, r.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_networking_qos_dscp_marking_rule_v2 %s: %#v", id, r)

	return resourceNetworkingQoSDSCPMarkingRuleV2Read(ctx, d, meta)
}

func resourceNetworkingQoSDSCPMarkingRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_dscp_marking_rule_v2 ID %s: %s", d.Id(), err)
	}

	r, err := rules.GetDSCPMarkingRule(networkingClient, qosPolicyID, qosRuleID).ExtractDSCPMarkingRule()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_qos_dscp_marking_rule_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_dscp_marking_rule_v2 %s: %#v", d.Id(), r)

	d.Set("qos_policy_id", qosPolicyID)
	d.Set("dscp_mark", r.DSCPMark)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingQoSDSCPMarkingRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_dscp_marking_rule_v2 ID %s: %s", d.Id(), err)
	}

	if d.HasChange("dscp_mark") {
		dscpMark := d.Get("dscp_mark").(int)
		updateOpts := rules.UpdateDSCPMarkingRuleOpts{
			DSCPMark: &dscpMark,
		}
		log.Printf("[DEBUG] openstack_networking_qos_dscp_marking_rule_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = rules.UpdateDSCPMarkingRule(networkingClient, qosPolicyID, qosRuleID, updateOpts).ExtractDSCPMarkingRule()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_qos_dscp_marking_rule_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkingQoSDSCPMarkingRuleV2Read(ctx, d, meta)
}

func resourceNetworkingQoSDSCPMarkingRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_qos_dscp_marking_rule_v2 ID %s: %s", d.Id(), err)
	}

	if err := rules.DeleteDSCPMarkingRule(networkingClient, qosPolicyID, qosRuleID).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_qos_dscp_marking_rule_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingQoSDSCPMarkingRuleV2StateRefreshFunc(networkingClient, qosPolicyID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_dscp_marking_rule_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
