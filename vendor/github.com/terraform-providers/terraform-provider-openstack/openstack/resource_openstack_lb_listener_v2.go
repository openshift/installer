package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	neutronlisteners "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerV2Create,
		Read:   resourceListenerV2Read,
		Update: resourceListenerV2Update,
		Delete: resourceListenerV2Delete,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS", "TERMINATED_HTTPS",
				}, false),
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"timeout_client_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_connect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_tcp_inspect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"insert_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceListenerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBV2LoadBalancer(lbClient, d.Get("loadbalancer_id").(string), "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	// Choose either the Octavia or Neutron create options.
	createOpts, err := chooseLBV2ListenerCreateOpts(d, config)
	if err != nil {
		return fmt.Errorf("Error building openstack_lb_listener_v2 create options: %s", err)
	}

	log.Printf("[DEBUG] openstack_lb_listener_v2 create options: %#v", createOpts)
	var listener *neutronlisteners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = neutronlisteners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating openstack_lb_listener_v2: %s", err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	d.SetId(listener.ID)

	return resourceListenerV2Read(d, meta)
}

func resourceListenerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Use Octavia listener body if Octavia/LBaaS is enabled.
	if config.UseOctavia {
		listener, err := octavialisteners.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "openstack_lb_listener_v2")
		}

		log.Printf("[DEBUG] Retrieved openstack_lb_listener_v2 %s: %#v", d.Id(), listener)

		d.Set("name", listener.Name)
		d.Set("protocol", listener.Protocol)
		d.Set("tenant_id", listener.ProjectID)
		d.Set("description", listener.Description)
		d.Set("protocol_port", listener.ProtocolPort)
		d.Set("admin_state_up", listener.AdminStateUp)
		d.Set("default_pool_id", listener.DefaultPoolID)
		d.Set("connection_limit", listener.ConnLimit)
		d.Set("timeout_client_data", listener.TimeoutClientData)
		d.Set("timeout_member_connect", listener.TimeoutMemberConnect)
		d.Set("timeout_member_data", listener.TimeoutMemberData)
		d.Set("timeout_tcp_inspect", listener.TimeoutTCPInspect)
		d.Set("sni_container_refs", listener.SniContainerRefs)
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
		d.Set("region", GetRegion(d, config))

		// Required by import.
		if len(listener.Loadbalancers) > 0 {
			d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
		}

		if err := d.Set("insert_headers", listener.InsertHeaders); err != nil {
			return fmt.Errorf("Unable to set openstack_lb_listener_v2 insert_headers: %s", err)
		}

		return nil
	}

	// Use Neutron/Networking in other case.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "openstack_lb_listener_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_listener_v2 %s: %#v", d.Id(), listener)

	// Required by import.
	if len(listener.Loadbalancers) > 0 {
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
	}

	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("tenant_id", listener.TenantID)
	d.Set("description", listener.Description)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("connection_limit", listener.ConnLimit)
	d.Set("sni_container_refs", listener.SniContainerRefs)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceListenerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	updateOpts, err := chooseLBV2ListenerUpdateOpts(d, config)
	if err != nil {
		return fmt.Errorf("Error building openstack_lb_listener_v2 update options: %s", err)
	}
	if updateOpts == nil {
		log.Printf("[DEBUG] openstack_lb_listener_v2 %s: nothing to update", d.Id())
		return resourceListenerV2Read(d, meta)
	}

	log.Printf("[DEBUG] openstack_lb_listener_v2 %s update options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = neutronlisteners.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error updating openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	return resourceListenerV2Read(d, meta)

}

func resourceListenerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Unable to retrieve openstack_lb_listener_v2")
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting openstack_lb_listener_v2 %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronlisteners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_lb_listener_v2")
	}

	// Wait for the listener to become DELETED.
	err = waitForLBV2Listener(lbClient, listener, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return err
	}

	return nil
}
