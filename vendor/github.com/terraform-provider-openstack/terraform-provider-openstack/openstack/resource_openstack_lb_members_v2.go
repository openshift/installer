package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	neutronpools "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func resourceMembersV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceMembersV2Create,
		Read:   resourceMembersV2Read,
		Update: resourceMembersV2Update,
		Delete: resourceMembersV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"address": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 256),
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"admin_state_up": {
							Type:     schema.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMembersV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := expandLBMembersV2(d.Get("member").(*schema.Set))
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
	}

	// Wait for parent pool to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create members")
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = octaviapools.BatchUpdateMembers(lbClient, poolID, createOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating members: %s", err)
	}

	// Wait for parent pool to become active before continuing
	err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	d.SetId(poolID)

	return resourceMembersV2Read(d, meta)
}

func resourceMembersV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	allPages, err := octaviapools.ListMembers(lbClient, d.Id(), octaviapools.ListMembersOpts{}).AllPages()
	if err != nil {
		return CheckDeleted(d, err, "Error getting openstack_lb_members_v2")
	}

	members, err := octaviapools.ExtractMembers(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_lb_members_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved members for the %s pool: %#v", d.Id(), members)

	d.Set("pool_id", d.Id())
	d.Set("member", flattenLBMembersV2(members))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceMembersV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if d.HasChange("member") {
		updateOpts := expandLBMembersV2(d.Get("member").(*schema.Set))

		// Get a clean copy of the parent pool.
		parentPool, err := neutronpools.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Unable to retrieve parent pool %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Updating %s pool members with options: %#v", d.Id(), updateOpts)
		err = resource.Retry(timeout, func() *resource.RetryError {
			err = octaviapools.BatchUpdateMembers(lbClient, d.Id(), updateOpts).ExtractErr()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Unable to update member %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing
		err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return err
		}
	}

	return resourceMembersV2Read(d, meta)
}

func resourceMembersV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the parent pool.
	parentPool, err := neutronpools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, fmt.Sprintf("Unable to retrieve parent pool (%s) for the member", d.Id()))
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return CheckDeleted(d, err, "Error waiting for the members' pool status")
	}

	log.Printf("[DEBUG] Attempting to delete %s pool members", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = octaviapools.BatchUpdateMembers(lbClient, d.Id(), []octaviapools.BatchUpdateMemberOpts{}).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return CheckDeleted(d, err, "Error deleting members")
	}

	// Wait for parent pool to become active before continuing.
	err = waitForLBV2Pool(lbClient, parentPool, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return CheckDeleted(d, err, "Error waiting for the members' pool status")
	}

	return nil
}
