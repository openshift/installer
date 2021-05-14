package openstack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceComputeAggregateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeAggregateV2Create,
		Read:   resourceComputeAggregateV2Read,
		Update: resourceComputeAggregateV2Update,
		Delete: resourceComputeAggregateV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metadata": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				DefaultFunc: func() (interface{}, error) { return map[string]interface{}{}, nil },
			},

			"hosts": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				DefaultFunc: func() (interface{}, error) { return []string{}, nil },
			},
		},
	}
}

func resourceComputeAggregateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	aggregate, err := aggregates.Create(computeClient, aggregates.CreateOpts{
		Name:             d.Get("name").(string),
		AvailabilityZone: d.Get("zone").(string),
	}).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack aggregate: %s", err)
	}
	idStr := strconv.Itoa(aggregate.ID)
	d.SetId(idStr)

	h, ok := d.GetOk("hosts")
	if ok {
		hosts := h.(*schema.Set)
		for _, host := range hosts.List() {
			_, err = aggregates.AddHost(computeClient, aggregate.ID, aggregates.AddHostOpts{Host: host.(string)}).Extract()
			if err != nil {
				return fmt.Errorf("Error adding host %s to Openstack aggregate: %s", host.(string), err)
			}
		}
	}

	_, err = aggregates.SetMetadata(computeClient, aggregate.ID, aggregates.SetMetadataOpts{Metadata: d.Get("metadata").(map[string]interface{})}).Extract()
	if err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	return nil
}

func resourceComputeAggregateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	aggregate, err := aggregates.Get(computeClient, id).Extract()
	if err != nil {
		return fmt.Errorf("Error getting host aggregate: %s", err)
	}

	// Metadata is redundant with Availability Zone
	metadata := aggregate.Metadata
	_, ok := metadata["availability_zone"]
	if ok {
		delete(metadata, "availability_zone")
	}

	d.Set("name", aggregate.Name)
	d.Set("zone", aggregate.AvailabilityZone)
	d.Set("hosts", aggregate.Hosts)
	d.Set("metadata", metadata)

	return nil
}

func resourceComputeAggregateV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	var updateOpts aggregates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("zone") {
		updateOpts.AvailabilityZone = d.Get("zone").(string)
	}

	if updateOpts != (aggregates.UpdateOpts{}) {
		_, err = aggregates.Update(computeClient, id, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenStack aggregate: %s", err)
		}
	}

	if d.HasChange("hosts") {
		o, n := d.GetChange("hosts")
		oldHosts, newHosts := o.(*schema.Set), n.(*schema.Set)
		hostsToDelete := oldHosts.Difference(newHosts)
		hostsToAdd := newHosts.Difference(oldHosts)
		for _, h := range hostsToDelete.List() {
			host := h.(string)
			log.Printf("[DEBUG] Removing host '%s' from aggregate '%s'", host, d.Get("name"))
			_, err = aggregates.RemoveHost(computeClient, id, aggregates.RemoveHostOpts{Host: host}).Extract()
			if err != nil {
				return fmt.Errorf("Error deleting host %s from Openstack aggregate: %s", host, err)
			}
		}
		for _, h := range hostsToAdd.List() {
			host := h.(string)
			log.Printf("[DEBUG] Adding host '%s' to aggregate '%s'", host, d.Get("name"))
			_, err = aggregates.AddHost(computeClient, id, aggregates.AddHostOpts{Host: host}).Extract()
			if err != nil {
				return fmt.Errorf("Error adding host %s to Openstack aggregate: %s", host, err)
			}
		}
	}

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")
		metadata := mapDiffWithNilValues(oldMetadata.(map[string]interface{}), newMetadata.(map[string]interface{}))
		_, err = aggregates.SetMetadata(computeClient, id, aggregates.SetMetadataOpts{Metadata: metadata}).Extract()
		if err != nil {
			return fmt.Errorf("Error setting metadata: %s", err)
		}
	}

	return nil
}

func resourceComputeAggregateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	// Openstack do not delete the host aggregate if it's not empty
	hostsToDelete := d.Get("hosts").(*schema.Set)
	for _, h := range hostsToDelete.List() {
		host := h.(string)
		log.Printf("[DEBUG] Removing host '%s' from aggregate '%s'", host, d.Get("name"))
		_, err = aggregates.RemoveHost(computeClient, id, aggregates.RemoveHostOpts{Host: host}).Extract()
		if err != nil {
			return fmt.Errorf("Error deleting host %s from Openstack aggregate: %s", host, err)
		}
	}

	err = aggregates.Delete(computeClient, id).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting Openstack aggregate: %s", err)
	}

	return nil
}
