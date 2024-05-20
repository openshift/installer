package openstack

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas/monitors"
)

func resourceLBMonitorV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBMonitorV1Create,
		ReadContext:   resourceLBMonitorV1Read,
		UpdateContext: resourceLBMonitorV1Update,
		DeleteContext: resourceLBMonitorV1Delete,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"admin_state_up": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceLBMonitorV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := monitors.CreateOpts{
		TenantID:      d.Get("tenant_id").(string),
		Delay:         d.Get("delay").(int),
		Timeout:       d.Get("timeout").(int),
		MaxRetries:    d.Get("max_retries").(int),
		URLPath:       d.Get("url_path").(string),
		ExpectedCodes: d.Get("expected_codes").(string),
		HTTPMethod:    d.Get("http_method").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		monitorType := resourceLBMonitorV1DetermineType(v.(string))
		createOpts.Type = monitorType
	}

	asuRaw := d.Get("admin_state_up").(string)
	if asuRaw != "" {
		asu, err := strconv.ParseBool(asuRaw)
		if err != nil {
			return diag.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
		}
		createOpts.AdminStateUp = &asu
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	m, err := monitors.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating OpenStack LB Monitor: %s", err)
	}
	log.Printf("[INFO] LB Monitor ID: %s", m.ID)

	log.Printf("[DEBUG] Waiting for OpenStack LB Monitor (%s) to become available.", m.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForLBMonitorActive(networkingClient, m.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(m.ID)

	return resourceLBMonitorV1Read(ctx, d, meta)
}

func resourceLBMonitorV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	m, err := monitors.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "LB monitor"))
	}

	log.Printf("[DEBUG] Retrieved OpenStack LB Monitor %s: %+v", d.Id(), m)

	d.Set("type", m.Type)
	d.Set("delay", m.Delay)
	d.Set("timeout", m.Timeout)
	d.Set("max_retries", m.MaxRetries)
	d.Set("tenant_id", m.TenantID)
	d.Set("url_path", m.URLPath)
	d.Set("http_method", m.HTTPMethod)
	d.Set("expected_codes", m.ExpectedCodes)
	d.Set("admin_state_up", strconv.FormatBool(m.AdminStateUp))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceLBMonitorV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	updateOpts := monitors.UpdateOpts{
		Delay:         d.Get("delay").(int),
		Timeout:       d.Get("timeout").(int),
		MaxRetries:    d.Get("max_retries").(int),
		URLPath:       d.Get("url_path").(string),
		HTTPMethod:    d.Get("http_method").(string),
		ExpectedCodes: d.Get("expected_codes").(string),
	}

	if d.HasChange("admin_state_up") {
		asuRaw := d.Get("admin_state_up").(string)
		if asuRaw != "" {
			asu, err := strconv.ParseBool(asuRaw)
			if err != nil {
				return diag.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
			}
			updateOpts.AdminStateUp = &asu
		}
	}

	log.Printf("[DEBUG] Updating OpenStack LB Monitor %s with options: %+v", d.Id(), updateOpts)

	_, err = monitors.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating OpenStack LB Monitor: %s", err)
	}

	return resourceLBMonitorV1Read(ctx, d, meta)
}

func resourceLBMonitorV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "PENDING_DELETE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForLBMonitorDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting OpenStack LB Monitor: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceLBMonitorV1DetermineType(t string) monitors.MonitorType {
	var monitorType monitors.MonitorType
	switch t {
	case "PING":
		monitorType = monitors.TypePING
	case "TCP":
		monitorType = monitors.TypeTCP
	case "HTTP":
		monitorType = monitors.TypeHTTP
	case "HTTPS":
		monitorType = monitors.TypeHTTPS
	}

	return monitorType
}

func waitForLBMonitorActive(networkingClient *gophercloud.ServiceClient, monitorID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, err := monitors.Get(networkingClient, monitorID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The monitor resource has no Status attribute, so a successful Get is the best we can do
		log.Printf("[DEBUG] OpenStack LB Monitor: %+v", m)
		return m, "ACTIVE", nil
	}
}

func waitForLBMonitorDelete(networkingClient *gophercloud.ServiceClient, monitorID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenStack LB Monitor %s", monitorID)

		m, err := monitors.Get(networkingClient, monitorID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenStack LB Monitor %s", monitorID)
				return m, "DELETED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault409); ok {
				log.Printf("[DEBUG] OpenStack LB Monitor (%s) is waiting for Pool to delete.", monitorID)
				return m, "PENDING", nil
			}

			return m, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenStack LB Monitor: %+v", m)
		err = monitors.Delete(networkingClient, monitorID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenStack LB Monitor %s", monitorID)
				return m, "DELETED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault409); ok {
				log.Printf("[DEBUG] OpenStack LB Monitor (%s) is waiting for Pool to delete.", monitorID)
				return m, "PENDING", nil
			}

			return m, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenStack LB Monitor %s still active.", monitorID)
		return m, "ACTIVE", nil
	}
}
