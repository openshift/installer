package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
)

func dataSourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"all_projects": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"serial": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"transferred_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"masters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDNSZoneV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return diag.FromErr(err)
	}

	listOpts := zones.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("email"); ok {
		listOpts.Email = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		listOpts.TTL = v.(int)
	}

	if v, ok := d.GetOk("type"); ok {
		listOpts.Type = v.(string)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		log.Printf("[DEBUG] unable to ser auth header: %s", err)
	}

	pages, err := zones.List(dnsClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve zones: %s", err)
	}

	allZones, err := zones.ExtractZones(pages)
	if err != nil {
		return diag.Errorf("Unable to extract zones: %s", err)
	}

	if len(allZones) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allZones) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	zone := allZones[0]

	log.Printf("[DEBUG] Retrieved DNS Zone %s: %+v", zone.ID, zone)
	d.SetId(zone.ID)

	// strings
	d.Set("name", zone.Name)
	d.Set("pool_id", zone.PoolID)
	d.Set("project_id", zone.ProjectID)
	d.Set("email", zone.Email)
	d.Set("description", zone.Description)
	d.Set("status", zone.Status)
	d.Set("type", zone.Type)
	d.Set("region", GetRegion(d, config))

	// ints
	d.Set("ttl", zone.TTL)
	d.Set("version", zone.Version)
	d.Set("serial", zone.Serial)

	// time.Times
	d.Set("created_at", zone.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", zone.UpdatedAt.Format(time.RFC3339))
	d.Set("transferred_at", zone.TransferredAt.Format(time.RFC3339))

	// maps
	err = d.Set("attributes", zone.Attributes)
	if err != nil {
		log.Printf("[DEBUG] Unable to set attributes: %s", err)
		return diag.FromErr(err)
	}

	// slices
	err = d.Set("masters", zone.Masters)
	if err != nil {
		log.Printf("[DEBUG] Unable to set masters: %s", err)
		return diag.FromErr(err)
	}

	return nil
}
