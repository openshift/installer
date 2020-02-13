package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbacpolicies"
)

func resourceNetworkingRBACPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRBACPolicyV2Create,
		Read:   resourceNetworkingRBACPolicyV2Read,
		Update: resourceNetworkingRBACPolicyV2Update,
		Delete: resourceNetworkingRBACPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access_as_external", "access_as_shared",
				}, false),
			},

			"object_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"network", "qos_policy",
				}, false),
			},

			"target_tenant": {
				Type:     schema.TypeString,
				Required: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingRBACPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.PolicyAction(d.Get("action").(string)),
		ObjectType:   d.Get("object_type").(string),
		TargetTenant: d.Get("target_tenant").(string),
		ObjectID:     d.Get("object_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	rbac, err := rbacpolicies.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_networking_rbac_policy_v2: %s", err)
	}

	d.SetId(rbac.ID)

	return resourceNetworkingRBACPolicyV2Read(d, meta)
}

func resourceNetworkingRBACPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	rbac, err := rbacpolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_networking_rbac_policy_v2")
	}

	log.Printf("[DEBUG] Retrieved RBAC policy %s: %+v", d.Id(), rbac)

	d.Set("action", string(rbac.Action))
	d.Set("object_type", rbac.ObjectType)
	d.Set("target_tenant", rbac.TargetTenant)
	d.Set("object_id", rbac.ObjectID)
	d.Set("project_id", rbac.ProjectID)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRBACPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts rbacpolicies.UpdateOpts

	if d.HasChange("target_tenant") {
		updateOpts.TargetTenant = d.Get("target_tenant").(string)

		_, err := rbacpolicies.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_networking_rbac_policy_v2: %s", err)
		}
	}

	return resourceNetworkingRBACPolicyV2Read(d, meta)
}

func resourceNetworkingRBACPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = rbacpolicies.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_networking_rbac_policy_v2")
	}

	return nil
}
