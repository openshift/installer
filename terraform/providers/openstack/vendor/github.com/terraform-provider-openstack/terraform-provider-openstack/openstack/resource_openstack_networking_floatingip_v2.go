package openstack

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func resourceNetworkingFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkFloatingIPV2Create,
		ReadContext:   resourceNetworkFloatingIPV2Read,
		UpdateContext: resourceNetworkFloatingIPV2Update,
		DeleteContext: resourceNetworkFloatingIPV2Delete,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_POOL_NAME", nil),
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subnet_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"subnet_id"},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
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

			"dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"dns_domain": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^$|\.$`), "fully-qualified (unambiguous) DNS domain names must have a dot at the end"),
			},
		},
	}
}

func resourceNetworkFloatingIPV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	poolName := d.Get("pool").(string)
	poolID, err := networkingNetworkV2ID(d, meta, poolName)
	if err != nil {
		return diag.Errorf("Error retrieving ID for openstack_networking_floatingip_v2 pool name %s: %s", poolName, err)
	}
	if len(poolID) == 0 {
		return diag.Errorf("No network found with name: %s", poolName)
	}

	subnetID := d.Get("subnet_id").(string)
	var subnetIDs []string
	if v, ok := d.Get("subnet_ids").([]interface{}); ok {
		subnetIDs = make([]string, len(v))
		for i, v := range v {
			subnetIDs[i] = v.(string)
		}
	}

	if subnetID == "" && len(subnetIDs) > 0 {
		subnetID = subnetIDs[0]
	}

	createOpts := &floatingips.CreateOpts{
		FloatingNetworkID: poolID,
		Description:       d.Get("description").(string),
		FloatingIP:        d.Get("address").(string),
		PortID:            d.Get("port_id").(string),
		TenantID:          d.Get("tenant_id").(string),
		FixedIP:           d.Get("fixed_ip").(string),
		SubnetID:          subnetID,
	}

	var finalCreateOpts floatingips.CreateOptsBuilder
	finalCreateOpts = FloatingIPCreateOpts{
		createOpts,
		MapValueSpecs(d),
	}

	dnsName := d.Get("dns_name").(string)
	dnsDomain := d.Get("dns_domain").(string)
	if dnsName != "" || dnsDomain != "" {
		finalCreateOpts = dns.FloatingIPCreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			DNSName:           dnsName,
			DNSDomain:         dnsDomain,
		}
	}

	var fip floatingIPExtended

	log.Printf("[DEBUG] openstack_networking_floatingip_v2 create options: %#v", finalCreateOpts)

	if len(subnetIDs) == 0 {
		// floating IP allocation without a retry
		err = floatingips.Create(networkingClient, finalCreateOpts).ExtractInto(&fip)
		if err != nil {
			return diag.Errorf("Error creating openstack_networking_floatingip_v2: %s", err)
		}
	} else {
		// create a floatingip in a loop with the first available external subnet
		for i, subnetID := range subnetIDs {
			createOpts.SubnetID = subnetID

			log.Printf("[DEBUG] openstack_networking_floatingip_v2 create options (try %d): %#v", i+1, finalCreateOpts)

			err = floatingips.Create(networkingClient, finalCreateOpts).ExtractInto(&fip)
			if err != nil {
				if retryOn409(err) {
					continue
				}
				return diag.Errorf("Error creating openstack_networking_floatingip_v2: %s", err)
			}
			break
		}
		// handle the last error
		if err != nil {
			return diag.Errorf("Error creating openstack_networking_floatingip_v2: %d subnets exhausted: %s", len(subnetIDs), err)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_floatingip_v2 %s to become available.", fip.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    networkingFloatingIPV2StateRefreshFunc(networkingClient, fip.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_floatingip_v2 %s to become available: %s", fip.ID, err)
	}

	d.SetId(fip.ID)

	if createOpts.SubnetID != "" {
		// resourceNetworkFloatingIPV2Read doesn't handle this, since FIP GET request doesn't provide this info.
		d.Set("subnet_id", createOpts.SubnetID)
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "floatingips", fip.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_floatingip_v2 %s: %s", fip.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_floatingip_v2 %s", tags, fip.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_floatingip_v2 %s: %#v", fip.ID, fip)
	return resourceNetworkFloatingIPV2Read(ctx, d, meta)
}

func resourceNetworkFloatingIPV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	var fip floatingIPExtended

	err = floatingips.Get(networkingClient, d.Id()).ExtractInto(&fip)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_floatingip_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_floatingip_v2 %s: %#v", d.Id(), fip)

	d.Set("description", fip.Description)
	d.Set("address", fip.FloatingIP.FloatingIP)
	d.Set("port_id", fip.PortID)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("tenant_id", fip.TenantID)
	d.Set("dns_name", fip.DNSName)
	d.Set("dns_domain", fip.DNSDomain)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, fip.Tags)

	poolName, err := networkingNetworkV2Name(d, meta, fip.FloatingNetworkID)
	if err != nil {
		return diag.Errorf("Error retrieving pool name for openstack_networking_floatingip_v2 %s: %s", d.Id(), err)
	}
	d.Set("pool", poolName)

	return nil
}

func resourceNetworkFloatingIPV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	var hasChange bool
	var updateOpts floatingips.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	// fixed_ip_address cannot be specified without a port_id
	if d.HasChange("port_id") || d.HasChange("fixed_ip") {
		hasChange = true
		portID := d.Get("port_id").(string)
		updateOpts.PortID = &portID
	}

	if d.HasChange("fixed_ip") {
		hasChange = true
		fixedIP := d.Get("fixed_ip").(string)
		updateOpts.FixedIP = fixedIP
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_floatingip_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_floatingip_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "floatingips", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_floatingip_v2 %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_floatingip_v2 %s", tags, d.Id())
	}

	return resourceNetworkFloatingIPV2Read(ctx, d, meta)
}

func resourceNetworkFloatingIPV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	if err := floatingips.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_floatingip_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DOWN"},
		Target:     []string{"DELETED"},
		Refresh:    networkingFloatingIPV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_floatingip_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
