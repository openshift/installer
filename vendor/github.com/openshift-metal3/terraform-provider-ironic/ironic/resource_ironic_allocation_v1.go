package ironic

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/allocations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema resource definition for an Ironic allocation.
func resourceAllocationV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocationV1Create,
		Read:   resourceAllocationV1Read,
		Delete: resourceAllocationV1Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_class": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"candidate_nodes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"traits": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"node_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_error": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Create an allocation, including driving Ironic's state machine
func resourceAllocationV1Create(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	result, err := allocations.Create(client, allocationSchemaToCreateOpts(d)).Extract()
	if err != nil {
		return err
	}

	d.SetId(result.UUID)

	// Wait for state to change from allocating
	var state string
	timeout := 1 * time.Minute
	checkInterval := 2 * time.Second

	getState := func() string {
		resourceAllocationV1Read(d, meta)
		return d.Get("state").(string)
	}

	for state = getState(); state == "allocating"; state = getState() {
		log.Printf("[DEBUG] Requested allocation %s; current state is '%s'\n", d.Id(), state)

		time.Sleep(checkInterval)
		checkInterval += 2
		timeout -= checkInterval
		if timeout < 0 {
			return fmt.Errorf("timed out waiting for allocation")
		}
	}

	if state == "error" {
		err := d.Get("last_error").(string)
		resourceAllocationV1Delete(d, meta)
		d.SetId("")
		return fmt.Errorf("error creating resource: %s", err)
	}

	return nil
}

// Read the allocation's data from Ironic
func resourceAllocationV1Read(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	result, err := allocations.Get(client, d.Id()).Extract()
	if err != nil {
		return err
	}

	d.Set("name", result.Name)
	d.Set("resource_class", result.ResourceClass)
	d.Set("candidate_nodes", result.CandidateNodes)
	d.Set("traits", result.Traits)
	d.Set("extra", result.Extra)
	d.Set("node_uuid", result.NodeUUID)
	d.Set("state", result.State)
	d.Set("last_error", result.LastError)

	return nil
}

// Delete an allocation from Ironic if it exists
func resourceAllocationV1Delete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	_, err = allocations.Get(client, d.Id()).Extract()
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		return nil
	}

	return allocations.Delete(client, d.Id()).ExtractErr()
}

func allocationSchemaToCreateOpts(d *schema.ResourceData) *allocations.CreateOpts {
	candidateNodesRaw := d.Get("candidate_nodes").([]interface{})
	traitsRaw := d.Get("traits").([]interface{})
	extraRaw := d.Get("extra").(map[string]interface{})

	candidateNodes := make([]string, len(candidateNodesRaw))
	for i := range candidateNodesRaw {
		candidateNodes[i] = candidateNodesRaw[i].(string)
	}

	traits := make([]string, len(traitsRaw))
	for i := range traitsRaw {
		traits[i] = traitsRaw[i].(string)
	}

	extra := make(map[string]string)
	for k, v := range extraRaw {
		extra[k] = v.(string)
	}

	return &allocations.CreateOpts{
		Name:           d.Get("name").(string),
		ResourceClass:  d.Get("resource_class").(string),
		CandidateNodes: candidateNodes,
		Traits:         traits,
		Extra:          extra,
	}
}
