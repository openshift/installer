package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

const (
	errEnableSNATWithoutExternalNet = "setting enable_snat for openstack_networking_router_v2 " +
		"requires external_network_id to be set"

	errExternalFixedIPWithoutExternalNet = "setting an external_fixed_ip for openstack_networking_router_v2 " +
		"requires external_network_id to be set"

	errExternalSubnetIDWithoutExternalNet = "setting external_subnet_ids for openstack_networking_router_v2 " +
		"requires external_network_id to be set"
)

func resourceNetworkingRouterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterV2Create,
		ReadContext:   resourceNetworkingRouterV2Read,
		UpdateContext: resourceNetworkingRouterV2Update,
		DeleteContext: resourceNetworkingRouterV2Delete,
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
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"external_gateway": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      false,
				Computed:      true,
				Deprecated:    "use external_network_id instead",
				ConflictsWith: []string{"external_network_id"},
			},

			"external_network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      false,
				Computed:      true,
				ConflictsWith: []string{"external_gateway"},
			},

			"enable_snat": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"external_fixed_ip": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"external_subnet_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"external_fixed_ip"},
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"availability_zone_hints": {
				Type:     schema.TypeList,
				Computed: true,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"set_router_gateway_after_create": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkingRouterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := RouterCreateOpts{
		routers.CreateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			TenantID:              d.Get("tenant_id").(string),
			AvailabilityZoneHints: resourceNetworkingAvailabilityZoneHintsV2(d),
		},
		MapValueSpecs(d),
	}

	if asuRaw, ok := d.GetOk("admin_state_up"); ok {
		asu := asuRaw.(bool)
		createOpts.AdminStateUp = &asu
	}

	if dRaw, ok := d.GetOkExists("distributed"); ok {
		d := dRaw.(bool)
		createOpts.Distributed = &d
	}

	// Get Vendor_options
	vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
	var vendorUpdateGateway bool
	if vendorOptionsRaw.Len() > 0 {
		vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
		vendorUpdateGateway = vendorOptions["set_router_gateway_after_create"].(bool)
	}

	// Gateway settings
	var externalNetworkID string
	var gatewayInfo routers.GatewayInfo
	if v := d.Get("external_gateway").(string); v != "" {
		externalNetworkID = v
		gatewayInfo.NetworkID = externalNetworkID
	}

	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
		gatewayInfo.NetworkID = externalNetworkID
	}

	if esRaw, ok := d.GetOkExists("enable_snat"); ok {
		if externalNetworkID == "" {
			return diag.Errorf(errEnableSNATWithoutExternalNet)
		}
		es := esRaw.(bool)
		gatewayInfo.EnableSNAT = &es
	}

	externalFixedIPs := expandNetworkingRouterExternalFixedIPsV2(d.Get("external_fixed_ip").([]interface{}))
	if len(externalFixedIPs) > 0 {
		if externalNetworkID == "" {
			return diag.Errorf(errExternalFixedIPWithoutExternalNet)
		}
		gatewayInfo.ExternalFixedIPs = externalFixedIPs
	}

	externalSubnetIDs := expandNetworkingRouterExternalSubnetIDsV2(d.Get("external_subnet_ids").([]interface{}))

	// vendorUpdateGateway is a flag for certain vendor-specific virtual routers
	// which do not allow gateway settings to be set during router creation.
	// If this flag was not enabled, then we can safely set the gateway
	// information during create.
	if !vendorUpdateGateway && externalNetworkID != "" {
		createOpts.GatewayInfo = &gatewayInfo
	}

	var r *routers.Router
	log.Printf("[DEBUG] openstack_networking_router_v2 create options: %#v", createOpts)

	if len(externalSubnetIDs) == 0 {
		// router creation without a retry
		r, err = routers.Create(networkingClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating openstack_networking_router_v2: %s", err)
		}
	} else {
		if externalNetworkID == "" {
			return diag.Errorf(errExternalSubnetIDWithoutExternalNet)
		}

		// create a router in a loop with the first available external subnet
		for i, externalSubnetID := range externalSubnetIDs {
			gatewayInfo.ExternalFixedIPs = []routers.ExternalFixedIP{externalSubnetID}

			log.Printf("[DEBUG] openstack_networking_router_v2 create options (try %d): %#v", i+1, createOpts)

			r, err = routers.Create(networkingClient, createOpts).Extract()
			if err != nil {
				if retryOn409(err) {
					continue
				}
				return diag.Errorf("Error creating openstack_networking_router_v2: %s", err)
			}
			break
		}
		// handle the last error
		if err != nil {
			return diag.Errorf("Error creating openstack_networking_router_v2: %d subnets exhausted: %s", len(externalSubnetIDs), err)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_router_v2 %s to become available.", r.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    resourceNetworkingRouterV2StateRefreshFunc(networkingClient, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_router_v2 %s to become available: %s", r.ID, err)
	}

	d.SetId(r.ID)

	// If the vendorUpdateGateway flag was specified and if an external network
	// was specified, then set the gateway information after router creation.
	if vendorUpdateGateway && externalNetworkID != "" {
		log.Printf("[DEBUG] Adding external_network %s to openstack_networking_router_v2 %s", externalNetworkID, r.ID)

		var updateOpts routers.UpdateOpts
		updateOpts.GatewayInfo = &gatewayInfo

		log.Printf("[DEBUG] Assigning external_gateway to openstack_networking_router_v2 %s with options: %#v", r.ID, updateOpts)
		_, err = routers.Update(networkingClient, r.ID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_router_v2: %s", err)
		}
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "routers", r.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_router_v2 %s: %s", r.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_router_v2 %s", tags, r.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_router_v2 %s: %#v", r.ID, r)
	return resourceNetworkingRouterV2Read(ctx, d, meta)
}

func resourceNetworkingRouterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	r, err := routers.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving openstack_networking_router_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_router_v2 %s: %#v", d.Id(), r)

	// Basic settings.
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("admin_state_up", r.AdminStateUp)
	d.Set("distributed", r.Distributed)
	d.Set("tenant_id", r.TenantID)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, r.Tags)

	if err := d.Set("availability_zone_hints", r.AvailabilityZoneHints); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_router_v2 %s availability_zone_hints: %s", d.Id(), err)
	}

	// Gateway settings.
	d.Set("external_gateway", r.GatewayInfo.NetworkID)
	d.Set("external_network_id", r.GatewayInfo.NetworkID)
	d.Set("enable_snat", r.GatewayInfo.EnableSNAT)

	externalFixedIPs := flattenNetworkingRouterExternalFixedIPsV2(r.GatewayInfo.ExternalFixedIPs)
	if err = d.Set("external_fixed_ip", externalFixedIPs); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_router_v2 %s external_fixed_ip: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkingRouterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Id()
	config.MutexKV.Lock(routerID)
	defer config.MutexKV.Unlock(routerID)

	var hasChange bool
	var updateOpts routers.UpdateOpts
	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Gateway settings.
	var updateGatewaySettings bool
	var externalNetworkID string
	gatewayInfo := routers.GatewayInfo{}

	if v := d.Get("external_gateway").(string); v != "" {
		externalNetworkID = v
	}

	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
	}

	if externalNetworkID != "" {
		gatewayInfo.NetworkID = externalNetworkID
	}

	if d.HasChange("external_gateway") {
		updateGatewaySettings = true
	}

	if d.HasChange("external_network_id") {
		updateGatewaySettings = true
	}

	if d.HasChange("enable_snat") {
		updateGatewaySettings = true
		if externalNetworkID == "" {
			return diag.Errorf(errEnableSNATWithoutExternalNet)
		}

		enableSNAT := d.Get("enable_snat").(bool)
		gatewayInfo.EnableSNAT = &enableSNAT
	}

	if d.HasChange("external_fixed_ip") {
		updateGatewaySettings = true

		externalFixedIPs := expandNetworkingRouterExternalFixedIPsV2(d.Get("external_fixed_ip").([]interface{}))
		gatewayInfo.ExternalFixedIPs = externalFixedIPs
		if len(externalFixedIPs) > 0 {
			if externalNetworkID == "" {
				return diag.Errorf(errExternalFixedIPWithoutExternalNet)
			}
		}
	}

	if updateGatewaySettings {
		hasChange = true
		updateOpts.GatewayInfo = &gatewayInfo
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_router_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = routers.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_router_v2: %s", err)
		}
	}

	// Next, perform any required updates to the tags.
	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "routers", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_router_v2 %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_router_v2 %s", tags, d.Id())
	}

	return resourceNetworkingRouterV2Read(ctx, d, meta)
}

func resourceNetworkingRouterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := routers.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_router_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRouterV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting openstack_networking_router_v2: %s", err)
	}

	d.SetId("")
	return nil
}
