package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/quotas"
)

func resourceLoadBalancerQuotaV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerQuotaV2Create,
		ReadContext:   resourceLoadBalancerQuotaV2Read,
		UpdateContext: resourceLoadBalancerQuotaV2Update,
		DeleteContext: resourceLoadBalancerQuotaV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"loadbalancer": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"listener": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"member": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"pool": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"health_monitor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"l7_policy": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"l7_rule": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceLoadBalancerQuotaV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	if lbClient.Type != octaviaLBClientType {
		return diag.Errorf("Error creating openstack_lb_quota_v2: Only available when using octavia")
	}

	region := GetRegion(d, config)
	projectID := d.Get("project_id").(string)
	loadbalancer := d.Get("loadbalancer").(int)
	listener := d.Get("listener").(int)
	member := d.Get("member").(int)
	pool := d.Get("pool").(int)
	healthmonitor := d.Get("health_monitor").(int)

	updateOpts := quotas.UpdateOpts{
		Loadbalancer:  &loadbalancer,
		Listener:      &listener,
		Member:        &member,
		Pool:          &pool,
		Healthmonitor: &healthmonitor,
	}

	// l7_policy requires octavia minor version 2.19. Only set when specified
	if v, ok := d.GetOkExists("l7_policy"); ok {
		l7Policy := v.(int)
		updateOpts.L7Policy = &l7Policy
	}

	// l7_rule requires octavia minor version 2.19. Only set when specified
	if v, ok := d.GetOkExists("l7_rule"); ok {
		l7Rule := v.(int)
		updateOpts.L7Rule = &l7Rule
	}

	q, err := quotas.Update(lbClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_lb_quota_v2: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_lb_quota_v2 %#v", q)

	return resourceLoadBalancerQuotaV2Read(ctx, d, meta)
}

func resourceLoadBalancerQuotaV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	if lbClient.Type != octaviaLBClientType {
		return diag.Errorf("Error creating openstack_lb_quota_v2: Only available when using octavia")
	}

	// Pase projectID from resource id that is <project_id>/<region>
	projectID := strings.Split(d.Id(), "/")[0]

	q, err := quotas.Get(lbClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_lb_quota_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_quota_v2 %s: %#v", d.Id(), q)

	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("loadbalancer", q.Loadbalancer)
	d.Set("listener", q.Listener)
	d.Set("member", q.Member)
	d.Set("pool", q.Pool)
	d.Set("health_monitor", q.Healthmonitor)
	d.Set("l7_policy", q.L7Policy)
	d.Set("l7_rule", q.L7Rule)

	return nil
}

func resourceLoadBalancerQuotaV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	if lbClient.Type != octaviaLBClientType {
		return diag.Errorf("Error creating openstack_lb_quota_v2: Only available when using octavia")
	}

	var (
		hasChange  bool
		updateOpts quotas.UpdateOpts
	)

	if d.HasChange("loadbalancer") {
		hasChange = true
		loadbalancer := d.Get("loadbalancer").(int)
		updateOpts.Loadbalancer = &loadbalancer
	}

	if d.HasChange("listener") {
		hasChange = true
		listener := d.Get("listener").(int)
		updateOpts.Listener = &listener
	}

	if d.HasChange("member") {
		hasChange = true
		member := d.Get("member").(int)
		updateOpts.Member = &member
	}

	if d.HasChange("pool") {
		hasChange = true
		pool := d.Get("pool").(int)
		updateOpts.Pool = &pool
	}

	if d.HasChange("health_monitor") {
		hasChange = true
		healthmonitor := d.Get("health_monitor").(int)
		updateOpts.Healthmonitor = &healthmonitor
	}

	if d.HasChange("l7_policy") {
		hasChange = true
		l7Policy := d.Get("l7_policy").(int)
		updateOpts.L7Policy = &l7Policy
	}

	if d.HasChange("l7_rule") {
		hasChange = true
		l7Rule := d.Get("l7_rule").(int)
		updateOpts.L7Rule = &l7Rule
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_lb_quota_v2 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)
		_, err := quotas.Update(lbClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_lb_quota_v2: %s", err)
		}
	}

	return resourceLoadBalancerQuotaV2Read(ctx, d, meta)
}

func resourceLoadBalancerQuotaV2Delete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
