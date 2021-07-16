package ironic

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schema resource definition for an Ironic node.
func resourceNodeV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeV1Create,
		Read:   resourceNodeV1Read,
		Update: resourceNodeV1Update,
		Delete: resourceNodeV1Delete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"boot_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"clean": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"conductor_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"console_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deploy_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"driver": {
				Type:     schema.TypeString,
				Required: true,
			},
			"driver_info": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					/* FIXME: Password updates aren't considered. How can I know if the *local* data changed? */
					/* FIXME: Support drivers other than IPMI */
					if k == "driver_info.ipmi_password" && old == "******" {
						return true
					}

					return false
				},

				// driver_info could contain passwords
				Sensitive: true,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"root_device": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"inspect_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inspect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"manage": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"management_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"raid_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rescue_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"provision_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_power_state": {
				Type:     schema.TypeString,
				Optional: true,

				// If power_state is same as target_power_state, we have no changes to apply
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					return new == d.Get("power_state").(string)
				},
			},
			"power_state_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Create a node, including driving Ironic's state machine
func resourceNodeV1Create(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	// Create the node object in Ironic
	createOpts := schemaToCreateOpts(d)
	result, err := nodes.Create(client, createOpts).Extract()
	if err != nil {
		d.SetId("")
		return err
	}

	// Setting the ID is what tells terraform we were successful in creating the node
	log.Printf("[DEBUG] Node created with ID %s\n", d.Id())
	d.SetId(result.UUID)

	// Create ports as part of the node object - you may also use the native port resource
	portSet := d.Get("ports").(*schema.Set)
	if portSet != nil {
		portList := portSet.List()
		for _, portInterface := range portList {
			port := portInterface.(map[string]interface{})

			// Terraform map can't handle bool... seriously.
			var pxeEnabled bool
			if port["pxe_enabled"] != nil {
				if port["pxe_enabled"] == "true" {
					pxeEnabled = true
				} else {
					pxeEnabled = false
				}

			}
			// FIXME: All values other than address and pxe
			portCreateOpts := ports.CreateOpts{
				NodeUUID:   d.Id(),
				Address:    port["address"].(string),
				PXEEnabled: &pxeEnabled,
			}
			_, err := ports.Create(client, portCreateOpts).Extract()
			if err != nil {
				_ = resourcePortV1Read(d, meta)
				return err
			}
		}
	}

	// Make node manageable
	if d.Get("manage").(bool) || d.Get("clean").(bool) || d.Get("inspect").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "manage", nil, nil); err != nil {
			return fmt.Errorf("could not manage: %s", err)
		}
	}

	// Clean node
	if d.Get("clean").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "clean", nil, nil); err != nil {
			return fmt.Errorf("could not clean: %s", err)
		}
	}

	// Inspect node
	if d.Get("inspect").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "inspect", nil, nil); err != nil {
			return fmt.Errorf("could not inspect: %s", err)
		}
	}

	// Make node available
	if d.Get("available").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "provide", nil, nil); err != nil {
			return fmt.Errorf("could not make node available: %s", err)
		}
	}

	// Change power state, if required
	if targetPowerState := d.Get("target_power_state").(string); targetPowerState != "" {
		err := changePowerState(client, d, nodes.TargetPowerState(targetPowerState))
		if err != nil {
			return fmt.Errorf("could not change power state: %s", err)
		}
	}

	return resourceNodeV1Read(d, meta)
}

// Read the node's data from Ironic
func resourceNodeV1Read(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	node, err := nodes.Get(client, d.Id()).Extract()
	if err != nil {
		d.SetId("")
		return err
	}

	// TODO: Ironic's Create is different than the Node object itself, GET returns things like the
	//  RaidConfig, we need to add those and handle them in CREATE
	err = d.Set("boot_interface", node.BootInterface)
	if err != nil {
		return err
	}
	err = d.Set("conductor_group", node.ConductorGroup)
	if err != nil {
		return err
	}
	err = d.Set("console_interface", node.ConsoleInterface)
	if err != nil {
		return err
	}
	err = d.Set("deploy_interface", node.DeployInterface)
	if err != nil {
		return err
	}
	err = d.Set("driver", node.Driver)
	if err != nil {
		return err
	}
	err = d.Set("driver_info", node.DriverInfo)
	if err != nil {
		return err
	}
	err = d.Set("extra", node.Extra)
	if err != nil {
		return err
	}
	err = d.Set("inspect_interface", node.InspectInterface)
	if err != nil {
		return err
	}
	err = d.Set("instance_uuid", node.InstanceUUID)
	if err != nil {
		return err
	}
	err = d.Set("management_interface", node.ManagementInterface)
	if err != nil {
		return err
	}
	err = d.Set("name", node.Name)
	if err != nil {
		return err
	}
	err = d.Set("network_interface", node.NetworkInterface)
	if err != nil {
		return err
	}
	err = d.Set("owner", node.Owner)
	if err != nil {
		return err
	}
	err = d.Set("power_interface", node.PowerInterface)
	if err != nil {
		return err
	}
	err = d.Set("power_state", node.PowerState)
	if err != nil {
		return err
	}
	err = d.Set("root_device", node.Properties["root_device"])
	if err != nil {
		return err
	}
	delete(node.Properties, "root_device")
	err = d.Set("properties", node.Properties)
	if err != nil {
		return err
	}
	err = d.Set("raid_interface", node.RAIDInterface)
	if err != nil {
		return err
	}
	err = d.Set("rescue_interface", node.RescueInterface)
	if err != nil {
		return err
	}
	err = d.Set("resource_class", node.ResourceClass)
	if err != nil {
		return err
	}
	err = d.Set("storage_interface", node.StorageInterface)
	if err != nil {
		return err
	}
	err = d.Set("vendor_interface", node.VendorInterface)
	if err != nil {
		return err
	}
	return d.Set("provision_state", node.ProvisionState)
}

// Update a node's state based on the terraform config - TODO: handle everything
func resourceNodeV1Update(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	d.Partial(true)

	stringFields := []string{
		"boot_interface",
		"conductor_group",
		"console_interface",
		"deploy_interface",
		"driver",
		"inspect_interface",
		"management_interface",
		"name",
		"network_interface",
		"owner",
		"power_interface",
		"raid_interface",
		"rescue_interface",
		"resource_class",
		"storage_interface",
		"vendor_interface",
	}

	for _, field := range stringFields {
		if d.HasChange(field) {
			opts := nodes.UpdateOpts{
				nodes.UpdateOperation{
					Op:    nodes.ReplaceOp,
					Path:  fmt.Sprintf("/%s", field),
					Value: d.Get(field).(string),
				},
			}

			if _, err := UpdateNode(client, d.Id(), opts); err != nil {
				return err
			}
		}
	}

	// Make node manageable
	if (d.HasChange("manage") && d.Get("manage").(bool)) ||
		(d.HasChange("clean") && d.Get("clean").(bool)) ||
		(d.HasChange("inspect") && d.Get("inspect").(bool)) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "manage", nil, nil); err != nil {
			return fmt.Errorf("could not manage: %s", err)
		}
	}

	// Update power state if required
	if targetPowerState := d.Get("target_power_state").(string); d.HasChange("target_power_state") && targetPowerState != "" {
		if err := changePowerState(client, d, nodes.TargetPowerState(targetPowerState)); err != nil {
			return err
		}
	}

	// Clean node
	if d.HasChange("clean") && d.Get("clean").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "clean", nil, nil); err != nil {
			return fmt.Errorf("could not clean: %s", err)
		}
	}

	// Inspect node
	if d.HasChange("inspect") && d.Get("inspect").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "inspect", nil, nil); err != nil {
			return fmt.Errorf("could not inspect: %s", err)
		}
	}

	// Make node available
	if d.HasChange("available") && d.Get("available").(bool) {
		if err := ChangeProvisionStateToTarget(client, d.Id(), "provide", nil, nil); err != nil {
			return fmt.Errorf("could not make node available: %s", err)
		}
	}

	if d.HasChange("properties") || d.HasChange("root_device") {
		properties := propertiesMerge(d, "root_device")
		opts := nodes.UpdateOpts{
			nodes.UpdateOperation{
				Op:    nodes.AddOp,
				Path:  "/properties",
				Value: properties,
			},
		}
		if _, err := UpdateNode(client, d.Id(), opts); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceNodeV1Read(d, meta)
}

// Delete a node from Ironic
func resourceNodeV1Delete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	if err := ChangeProvisionStateToTarget(client, d.Id(), "deleted", nil, nil); err != nil {
		return err
	}

	return nodes.Delete(client, d.Id()).ExtractErr()
}

func propertiesMerge(d *schema.ResourceData, key string) map[string]interface{} {
	properties := d.Get("properties").(map[string]interface{})
	properties[key] = d.Get(key).(map[string]interface{})
	return properties
}

// Convert terraform schema to gophercloud CreateOpts
// TODO: Is there a better way to do this? Annotations?
func schemaToCreateOpts(d *schema.ResourceData) *nodes.CreateOpts {
	properties := propertiesMerge(d, "root_device")
	return &nodes.CreateOpts{
		BootInterface:       d.Get("boot_interface").(string),
		ConductorGroup:      d.Get("conductor_group").(string),
		ConsoleInterface:    d.Get("console_interface").(string),
		DeployInterface:     d.Get("deploy_interface").(string),
		Driver:              d.Get("driver").(string),
		DriverInfo:          d.Get("driver_info").(map[string]interface{}),
		Extra:               d.Get("extra").(map[string]interface{}),
		InspectInterface:    d.Get("inspect_interface").(string),
		ManagementInterface: d.Get("management_interface").(string),
		Name:                d.Get("name").(string),
		NetworkInterface:    d.Get("network_interface").(string),
		Owner:               d.Get("owner").(string),
		PowerInterface:      d.Get("power_interface").(string),
		Properties:          properties,
		RAIDInterface:       d.Get("raid_interface").(string),
		RescueInterface:     d.Get("rescue_interface").(string),
		ResourceClass:       d.Get("resource_class").(string),
		StorageInterface:    d.Get("storage_interface").(string),
		VendorInterface:     d.Get("vendor_interface").(string),
	}
}

// UpdateNode wraps gophercloud's update function, so we are able to retry on 409 when Ironic is busy.
func UpdateNode(client *gophercloud.ServiceClient, uuid string, opts nodes.UpdateOpts) (node *nodes.Node, err error) {
	interval := 5 * time.Second
	for retries := 0; retries < 5; retries++ {
		node, err = nodes.Update(client, uuid, opts).Extract()
		if _, ok := err.(gophercloud.ErrDefault409); ok {
			log.Printf("[DEBUG] Failed to update node: ironic is busy, will try again in %s", interval.String())
			time.Sleep(interval)
			interval *= 2
		} else {
			return
		}
	}

	return
}

// Call Ironic's API and change the power state of the node
func changePowerState(client *gophercloud.ServiceClient, d *schema.ResourceData, target nodes.TargetPowerState) error {
	opts := nodes.PowerStateOpts{
		Target: target,
	}

	timeout := d.Get("power_state_timeout").(int)
	if timeout != 0 {
		opts.Timeout = timeout
	} else {
		timeout = 300 // used below for how long to wait for Ironic to finish
	}

	interval := 5 * time.Second
	for retries := 0; retries < 5; retries++ {
		err := nodes.ChangePowerState(client, d.Id(), opts).ExtractErr()
		if _, ok := err.(gophercloud.ErrDefault409); ok {
			log.Printf("[DEBUG] Failed to change power state: ironic is busy, will try again in %s", interval.String())
			time.Sleep(interval)
			interval *= 2
		} else {
			break
		}
	}

	// Wait for target_power_state to be empty, i.e. Ironic thinks it's finished
	checkInterval := 5

	for {
		node, err := nodes.Get(client, d.Id()).Extract()
		if err != nil {
			return err
		}

		if node.TargetPowerState == "" {
			break
		}

		time.Sleep(time.Duration(checkInterval) * time.Second)
		timeout -= checkInterval
		if timeout <= 0 {
			return fmt.Errorf("timed out waiting for power state change")
		}
	}

	return nil
}
