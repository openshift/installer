package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/request"
)

func resourceDNSTransferRequestV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSTransferRequestV2Create,
		Read:   resourceDNSTransferRequestV2Read,
		Update: resourceDNSTransferRequestV2Update,
		Delete: resourceDNSTransferRequestV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceDNSTransferRequestV2Import,
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

			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"disable_status_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDNSTransferRequestV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	createOpts := TransferRequestCreateOpts{
		request.CreateOpts{
			TargetProjectID: d.Get("target_project_id").(string),
			Description:     d.Get("description").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] openstack_dns_transfer_request_v2 create options: %#v", createOpts)

	zoneID := d.Get("zone_id").(string)
	n, err := request.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_transfer_request_zone_v2: %s", err)
	}

	if d.Get("disable_status_check").(bool) {
		d.SetId(n.ID)

		log.Printf("[DEBUG] Created OpenStack Zone Transfer request %s: %#v", n.ID, n)
		return resourceDNSTransferRequestV2Read(d, meta)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become active: %s", d.Id(), err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created OpenStack Zone Transfer request %s: %#v", n.ID, n)
	return resourceDNSTransferRequestV2Read(d, meta)
}

func resourceDNSTransferRequestV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	n, err := request.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_dns_transfer_request_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_transfer_request_v2 %s: %#v", d.Id(), n)

	d.Set("region", GetRegion(d, config))
	d.Set("zone_id", n.ZoneID)
	d.Set("target_project_id", n.TargetProjectID)
	d.Set("description", n.Description)
	d.Set("key", n.Key)

	return nil
}

func resourceDNSTransferRequestV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	var updateOpts request.UpdateOpts
	changed := false
	if d.HasChange("target_project_id") {
		updateOpts.TargetProjectID = d.Get("target_project_id").(string)
		changed = true
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
		changed = true
	}

	if !changed {
		return resourceDNSTransferRequestV2Read(d, meta)
	}

	log.Printf("[DEBUG] Updating openstack_dns_transfer_request_v2 %s with options: %#v", d.Id(), updateOpts)

	_, err = request.Update(dnsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating openstack_dns_transfer_request_v2 %s: %s", d.Id(), err)
	}

	if d.Get("disable_status_check").(bool) {
		return resourceDNSTransferRequestV2Read(d, meta)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become active: %s", d.Id(), err)
	}

	return resourceDNSTransferRequestV2Read(d, meta)
}

func resourceDNSTransferRequestV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	err = request.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_dns_transfer_request_v2")
	}

	if d.Get("disable_status_check").(bool) {
		return nil
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceDNSTransferRequestV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	n, err := request.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving openstack_dns_transfer_request_v2 %s: %s", d.Id(), err)
	}

	d.Set("zone_id", n.ZoneID)

	return []*schema.ResourceData{d}, nil
}
