package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
)

func resourceNetworkingQuotaV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingQuotaV2Create,
		ReadContext:   resourceNetworkingQuotaV2Read,
		UpdateContext: resourceNetworkingQuotaV2Update,
		Delete:        schema.RemoveFromState,
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

			"floatingip": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"rbac_policy": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"router": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_group_rule": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"subnet": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"subnetpool": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingQuotaV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts quotas.UpdateOpts

	projectID := d.Get("project_id").(string)
	floatingIP, ok := d.GetOk("floatingip")
	if ok {
		pfloatingIP := floatingIP.(int)
		updateOpts.FloatingIP = &pfloatingIP
	}
	network, ok := d.GetOk("network")
	if ok {
		pnetwork := network.(int)
		updateOpts.Network = &pnetwork
	}
	port, ok := d.GetOk("port")
	if ok {
		pport := port.(int)
		updateOpts.Port = &pport
	}
	rbacPolicy, ok := d.GetOk("rbac_policy")
	if ok {
		prbacPolicy := rbacPolicy.(int)
		updateOpts.RBACPolicy = &prbacPolicy
	}
	router, ok := d.GetOk("router")
	if ok {
		prouter := router.(int)
		updateOpts.Router = &prouter
	}
	securityGroup, ok := d.GetOk("security_group")
	if ok {
		psecurityGroup := securityGroup.(int)
		updateOpts.SecurityGroup = &psecurityGroup
	}
	securityGroupRule := d.Get("security_group_rule")
	if ok {
		psecurityGroupRule := securityGroupRule.(int)
		updateOpts.SecurityGroupRule = &psecurityGroupRule
	}
	subnet, ok := d.GetOk("subnet")
	if ok {
		psubnet := subnet.(int)
		updateOpts.Subnet = &psubnet
	}
	subnetPool, ok := d.GetOk("subnetpool")
	if ok {
		psubnetPool := subnetPool.(int)
		updateOpts.SubnetPool = &psubnetPool
	}

	q, err := quotas.Update(networkingClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_quota_v2: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_networking_quota_v2 %#v", q)

	return resourceNetworkingQuotaV2Read(ctx, d, meta)
}

func resourceNetworkingQuotaV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Depending on the provider version the resource was created, the resource id
	// can be either <project_id> or <project_id>/<region>. This parses the project_id
	// in both cases
	projectID := strings.Split(d.Id(), "/")[0]

	q, err := quotas.Get(networkingClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_quota_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_quota_v2 %s: %#v", d.Id(), q)

	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("floatingip", q.FloatingIP)
	d.Set("network", q.Network)
	d.Set("port", q.Port)
	d.Set("rbac_policy", q.RBACPolicy)
	d.Set("router", q.Router)
	d.Set("security_group", q.SecurityGroup)
	d.Set("security_group_rule", q.SecurityGroupRule)
	d.Set("subnet", q.Subnet)
	d.Set("subnetpool", q.SubnetPool)

	return nil
}

func resourceNetworkingQuotaV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts quotas.UpdateOpts
	)

	if d.HasChange("floatingip") {
		hasChange = true
		floatingIP := d.Get("floatingip").(int)
		updateOpts.FloatingIP = &floatingIP
	}

	if d.HasChange("network") {
		hasChange = true
		network := d.Get("network").(int)
		updateOpts.Network = &network
	}

	if d.HasChange("port") {
		hasChange = true
		port := d.Get("port").(int)
		updateOpts.Port = &port
	}

	if d.HasChange("rbac_policy") {
		hasChange = true
		rbacPolicy := d.Get("rbac_policy").(int)
		updateOpts.RBACPolicy = &rbacPolicy
	}

	if d.HasChange("router") {
		hasChange = true
		router := d.Get("router").(int)
		updateOpts.Router = &router
	}

	if d.HasChange("security_group") {
		hasChange = true
		securityGroup := d.Get("security_group").(int)
		updateOpts.SecurityGroup = &securityGroup
	}

	if d.HasChange("security_group_rule") {
		hasChange = true
		securityGroupRule := d.Get("security_group_rule").(int)
		updateOpts.SecurityGroupRule = &securityGroupRule
	}

	if d.HasChange("subnet") {
		hasChange = true
		subnet := d.Get("subnet").(int)
		updateOpts.Subnet = &subnet
	}

	if d.HasChange("subnetpool") {
		hasChange = true
		subnetPool := d.Get("subnetpool").(int)
		updateOpts.SubnetPool = &subnetPool
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_quota_v2 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)
		_, err := quotas.Update(networkingClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_quota_v2: %s", err)
		}
	}

	return resourceNetworkingQuotaV2Read(ctx, d, meta)
}
