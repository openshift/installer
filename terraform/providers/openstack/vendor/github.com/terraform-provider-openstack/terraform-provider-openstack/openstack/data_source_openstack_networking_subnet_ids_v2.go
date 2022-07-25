package openstack

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/terraform/hashcode"
)

func dataSourceNetworkingSubnetIDsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSubnetIDsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"dhcp_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: descriptions["tenant_id"],
			},

			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ipv6_address_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, true),
			},

			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkingSubnetIDsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if v, ok := d.GetOk("sort_key"); ok {
		listOpts.SortKey = v.(string)
	}

	if v, ok := d.GetOk("sort_direction"); ok {
		listOpts.SortDir = v.(string)
	}

	pages, err := subnets.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_subnet_ids_v2: %s", err)
	}

	allSubnets, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return diag.Errorf("Unable to extract openstack_networking_subnet_ids_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved %d subnets in openstack_networking_subnet_ids_v2: %+v", len(allSubnets), allSubnets)

	var subnetIDs []string
	if nameRegex, ok := d.GetOk("name_regex"); !ok {
		subnetIDs = make([]string, len(allSubnets))
		for i, subnet := range allSubnets {
			subnetIDs[i] = subnet.ID
		}
	} else {
		r := regexp.MustCompile(nameRegex.(string))
		for _, subnet := range allSubnets {
			// Check for a very rare case where the response would include no
			// subnet name. No name means nothing to attempt a match against,
			// therefore we are skipping such subnet.
			if subnet.Name == "" {
				log.Printf("[WARN] Unable to find subnet name to match against "+
					"for %q subnet ID, nothing to do.",
					subnet.ID)
				continue
			}
			if r.MatchString(subnet.Name) {
				subnetIDs = append(subnetIDs, subnet.ID)
			}
		}

		log.Printf("[DEBUG] Subnet list filtered by regex: %s", d.Get("name_regex"))
		log.Printf("[DEBUG] Retrieved %d subnet IDs after filtering in openstack_networking_subnet_ids_v2: %+v", len(subnetIDs), subnetIDs)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(strings.Join(subnetIDs, ","))))
	d.Set("ids", subnetIDs)
	d.Set("region", GetRegion(d, config))

	return nil
}
