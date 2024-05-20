package openstack

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func dataSourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSubnetV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"dhcp_enabled": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"dhcp_disabled"},
				Optional:      true,
			},

			"dhcp_disabled": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"dhcp_enabled"},
				Optional:      true,
				Deprecated:    "use dhcp_enabled instead",
			},

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: descriptions["tenant_id"],
			},

			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			// Computed values
			"allocation_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enable_dhcp": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dns_nameservers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"host_routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ipv6_address_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"slaac", "dhcpv6-stateful", "dhcpv6-stateless",
				}, false),
			},

			"ipv6_ra_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"slaac", "dhcpv6-stateful", "dhcpv6-stateless",
				}, false),
			},

			"subnetpool_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

func dataSourceNetworkingSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := subnets.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOkExists("dhcp_enabled"); ok {
		enableDHCP := v.(bool)
		listOpts.EnableDHCP = &enableDHCP
	}

	if v, ok := d.GetOkExists("dhcp_disabled"); ok {
		enableDHCP := !v.(bool)
		listOpts.EnableDHCP = &enableDHCP
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("ip_version"); ok {
		listOpts.IPVersion = v.(int)
	}

	if v, ok := d.GetOk("gateway_ip"); ok {
		listOpts.GatewayIP = v.(string)
	}

	if v, ok := d.GetOk("cidr"); ok {
		listOpts.CIDR = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("ipv6_address_mode"); ok {
		listOpts.IPv6AddressMode = v.(string)
	}

	if v, ok := d.GetOk("ipv6_ra_mode"); ok {
		listOpts.IPv6RAMode = v.(string)
	}

	if v, ok := d.GetOk("subnetpool_id"); ok {
		listOpts.SubnetPoolID = v.(string)
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := subnets.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_subnet_v2: %s", err)
	}

	allSubnets, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return diag.Errorf("Unable to extract openstack_networking_subnet_v2: %s", err)
	}

	if len(allSubnets) < 1 {
		return diag.Errorf("Your query returned no openstack_networking_subnet_v2. " +
			"Please change your search criteria and try again.")
	}

	if len(allSubnets) > 1 {
		return diag.Errorf("Your query returned more than one openstack_networking_subnet_v2." +
			" Please try a more specific search criteria")
	}

	subnet := allSubnets[0]

	log.Printf("[DEBUG] Retrieved openstack_networking_subnet_v2 %s: %+v", subnet.ID, subnet)
	d.SetId(subnet.ID)

	d.Set("name", subnet.Name)
	d.Set("description", subnet.Description)
	d.Set("tenant_id", subnet.TenantID)
	d.Set("network_id", subnet.NetworkID)
	d.Set("cidr", subnet.CIDR)
	d.Set("ip_version", subnet.IPVersion)
	d.Set("ipv6_address_mode", subnet.IPv6AddressMode)
	d.Set("ipv6_ra_mode", subnet.IPv6RAMode)
	d.Set("gateway_ip", subnet.GatewayIP)
	d.Set("enable_dhcp", subnet.EnableDHCP)
	d.Set("subnetpool_id", subnet.SubnetPoolID)
	d.Set("all_tags", subnet.Tags)
	d.Set("region", GetRegion(d, config))

	if err := d.Set("dns_nameservers", subnet.DNSNameservers); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_subnet_v2 dns_nameservers: %s", err)
	}

	if err := d.Set("host_routes", subnet.HostRoutes); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_subnet_v2 host_routes: %s", err)
	}

	allocationPools := flattenNetworkingSubnetV2AllocationPools(subnet.AllocationPools)
	if err := d.Set("allocation_pools", allocationPools); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_subnet_v2 allocation_pools: %s", err)
	}

	return nil
}
