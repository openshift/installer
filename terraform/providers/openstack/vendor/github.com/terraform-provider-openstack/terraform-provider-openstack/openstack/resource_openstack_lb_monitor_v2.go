package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	octaviamonitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	neutronmonitors "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func resourceMonitorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorV2Create,
		ReadContext:   resourceMonitorV2Read,
		UpdateContext: resourceMonitorV2Update,
		DeleteContext: resourceMonitorV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMonitorV2Import,
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

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP-CONNECT", "HTTP", "HTTPS", "TLS-HELLO", "PING",
				}, false),
			},

			"delay": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"max_retries_down": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceMonitorV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Choose either the Octavia or Neutron create options.
	createOpts := chooseLBV2MonitorCreateOpts(d, config)

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent openstack_lb_pool_v2 %s: %s", poolID, err)
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] openstack_lb_monitor_v2 create options: %#v", createOpts)
	var monitor *neutronmonitors.Monitor
	err = resource.Retry(timeout, func() *resource.RetryError {
		monitor, err = neutronmonitors.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to create openstack_lb_monitor_v2: %s", err)
	}

	// Wait for monitor to become active before continuing
	err = waitForLBV2Monitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(monitor.ID)

	return resourceMonitorV2Read(ctx, d, meta)
}

func resourceMonitorV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Use Octavia monitor body if Octavia/LBaaS is enabled.
	if config.UseOctavia {
		monitor, err := octaviamonitors.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "monitor"))
		}

		log.Printf("[DEBUG] Retrieved openstack_lb_monitor_v2 %s: %#v", d.Id(), monitor)

		d.Set("tenant_id", monitor.ProjectID)
		d.Set("type", monitor.Type)
		d.Set("delay", monitor.Delay)
		d.Set("timeout", monitor.Timeout)
		d.Set("max_retries", monitor.MaxRetries)
		d.Set("max_retries_down", monitor.MaxRetriesDown)
		d.Set("url_path", monitor.URLPath)
		d.Set("http_method", monitor.HTTPMethod)
		d.Set("expected_codes", monitor.ExpectedCodes)
		d.Set("admin_state_up", monitor.AdminStateUp)
		d.Set("name", monitor.Name)
		d.Set("region", GetRegion(d, config))

		// OpenContrail workaround (https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/762)
		if len(monitor.Pools) > 0 && monitor.Pools[0].ID != "" {
			d.Set("pool_id", monitor.Pools[0].ID)
		}

		return nil
	}

	// Use Neutron/Networking in other case.
	monitor, err := neutronmonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "monitor"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_monitor_v2 %s: %#v", d.Id(), monitor)

	// OpenContrail workaround (https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/762)
	if len(monitor.Pools) > 0 && monitor.Pools[0].ID != "" {
		d.Set("pool_id", monitor.Pools[0].ID)
	}

	d.Set("tenant_id", monitor.TenantID)
	d.Set("type", monitor.Type)
	d.Set("delay", monitor.Delay)
	d.Set("timeout", monitor.Timeout)
	d.Set("max_retries", monitor.MaxRetries)
	d.Set("url_path", monitor.URLPath)
	d.Set("http_method", monitor.HTTPMethod)
	d.Set("expected_codes", monitor.ExpectedCodes)
	d.Set("admin_state_up", monitor.AdminStateUp)
	d.Set("name", monitor.Name)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceMonitorV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	updateOpts := chooseLBV2MonitorUpdateOpts(d, config)
	if updateOpts == nil {
		log.Printf("[DEBUG] openstack_lb_monitor_v2 %s: nothing to update", d.Id())
		return resourceMonitorV2Read(ctx, d, meta)
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent openstack_lb_pool_v2 %s: %s", poolID, err)
	}

	// Get a clean copy of the monitor.
	monitor, err := neutronmonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_lb_monitor_v2 %s: %s", d.Id(), err)
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for monitor to become active before continuing.
	err = waitForLBV2Monitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] openstack_lb_monitor_v2 %s update options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = neutronmonitors.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to update openstack_lb_monitor_v2 %s: %s", d.Id(), err)
	}

	// Wait for monitor to become active before continuing
	err = waitForLBV2Monitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMonitorV2Read(ctx, d, meta)
}

func resourceMonitorV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent openstack_lb_pool_v2 (%s)"+
			" for the openstack_lb_monitor_v2: %s", poolID, err)
	}

	// Get a clean copy of the monitor.
	monitor, err := neutronmonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve openstack_lb_monitor_v2"))
	}

	// Wait for parent pool to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting openstack_lb_monitor_v2 %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronmonitors.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_monitor_v2"))
	}

	// Wait for monitor to become DELETED
	err = waitForLBV2Monitor(ctx, lbClient, parentPool, monitor, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMonitorV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	monitorID := parts[0]

	if len(monitorID) == 0 {
		return nil, fmt.Errorf("Invalid format specified for openstack_lb_monitor_v2. Format must be <monitorID>[/<poolID>]")
	}

	d.SetId(monitorID)

	if len(parts) == 2 {
		d.Set("pool_id", parts[1])
	}

	return []*schema.ResourceData{d}, nil
}
