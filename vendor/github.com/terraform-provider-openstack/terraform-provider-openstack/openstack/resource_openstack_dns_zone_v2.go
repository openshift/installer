package openstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneV2Create,
		Read:   resourceDNSZoneV2Read,
		Update: resourceDNSZoneV2Update,
		Delete: resourceDNSZoneV2Delete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// Allow import from different project with id:project_id
				parts := strings.Split(d.Id(), ":")
				if parts[0] == "" || len(parts) > 2 {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected zone <id> or <id>:<project_id>", d.Id())
				} else if len(parts) == 2 {
					d.Set("project_id", parts[1])
				}

				d.SetId(parts[0])
				return []*schema.ResourceData{d}, nil
			},
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
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PRIMARY", "SECONDARY",
				}, false),
			},

			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"masters": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDNSZoneV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	createOpts := ZoneCreateOpts{
		zones.CreateOpts{
			Name:        d.Get("name").(string),
			Type:        d.Get("type").(string),
			TTL:         d.Get("ttl").(int),
			Email:       d.Get("email").(string),
			Description: d.Get("description").(string),
			Attributes:  expandToMapStringString(d.Get("attributes").(map[string]interface{})),
			Masters:     expandToStringSlice(d.Get("masters").(*schema.Set).List()),
		},
		MapValueSpecs(d),
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return fmt.Errorf("Error setting dns client auth headers: %s", err)
	}

	log.Printf("[DEBUG] openstack_dns_zone_v2 create options: %#v", createOpts)
	n, err := zones.Create(dnsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_dns_zone_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_dns_zone_v2 %s to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsZoneV2RefreshFunc(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_zone_v2 %s to become active: %s", d.Id(), err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created OpenStack DNS Zone %s: %#v", n.ID, n)
	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return fmt.Errorf("Error setting dns client auth headers: %s", err)
	}

	n, err := zones.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_dns_zone_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_zone_v2 %s: %#v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("email", n.Email)
	d.Set("description", n.Description)
	d.Set("ttl", n.TTL)
	d.Set("type", n.Type)
	d.Set("attributes", n.Attributes)
	d.Set("masters", n.Masters)
	d.Set("region", GetRegion(d, config))
	d.Set("project_id", n.ProjectID)

	return nil
}

func resourceDNSZoneV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	var updateOpts zones.UpdateOpts
	if d.HasChange("email") {
		updateOpts.Email = d.Get("email").(string)
	}

	if d.HasChange("ttl") {
		updateOpts.TTL = d.Get("ttl").(int)
	}

	if d.HasChange("masters") {
		updateOpts.Masters = expandToStringSlice(d.Get("masters").(*schema.Set).List())
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return fmt.Errorf("Error setting dns client auth headers: %s", err)
	}

	log.Printf("[DEBUG] Updating openstack_dns_zone_v2 %s with options: %#v", d.Id(), updateOpts)

	_, err = zones.Update(dnsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating openstack_dns_zone_v2 %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsZoneV2RefreshFunc(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_zone_v2 %s to become active: %s", d.Id(), err)
	}

	return resourceDNSZoneV2Read(d, meta)
}

func resourceDNSZoneV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return fmt.Errorf("Error setting dns client auth headers: %s", err)
	}

	_, err = zones.Delete(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_dns_zone_v2")
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING"},
		Refresh:    dnsZoneV2RefreshFunc(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_zone_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
