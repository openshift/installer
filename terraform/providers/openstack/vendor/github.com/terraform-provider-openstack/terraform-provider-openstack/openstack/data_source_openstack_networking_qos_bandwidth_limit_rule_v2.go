package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func dataSourceNetworkingQoSBandwidthLimitRuleV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingQoSBandwidthLimitRuleV2Read,
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
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"max_burst_kbps": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"direction": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},
		},
	}
}

func dataSourceNetworkingQoSBandwidthLimitRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := rules.BandwidthLimitRulesListOpts{}

	if v, ok := d.GetOk("max_kbps"); ok {
		listOpts.MaxKBps = v.(int)
	}

	if v, ok := d.GetOk("max_burst_kbps"); ok {
		listOpts.MaxBurstKBps = v.(int)
	}

	qosPolicyID := d.Get("qos_policy_id").(string)

	pages, err := rules.ListBandwidthLimitRules(networkingClient, qosPolicyID, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_qos_bandwidth_limit_rule_v2: %s", err)
	}

	allRules, err := rules.ExtractBandwidthLimitRules(pages)
	if err != nil {
		return diag.Errorf("Unable to extract openstack_networking_qos_bandwidth_limit_rule_v2: %s", err)
	}

	if len(allRules) < 1 {
		return diag.Errorf("Your query returned no openstack_networking_qos_bandwidth_limit_rule_v2. " +
			"Please change your search criteria and try again.")
	}

	if len(allRules) > 1 {
		return diag.Errorf("Your query returned more than one openstack_networking_qos_bandwidth_limit_rule_v2." +
			" Please try a more specific search criteria")
	}

	rule := allRules[0]
	id := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, rule.ID)

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_bandwidth_limit_rule_v2 %s: %+v", id, rule)
	d.SetId(id)

	d.Set("qos_policy_id", qosPolicyID)
	d.Set("max_kbps", rule.MaxKBps)
	d.Set("max_burst_kbps", rule.MaxBurstKBps)
	d.Set("direction", rule.Direction)
	d.Set("region", GetRegion(d, config))

	return nil
}
