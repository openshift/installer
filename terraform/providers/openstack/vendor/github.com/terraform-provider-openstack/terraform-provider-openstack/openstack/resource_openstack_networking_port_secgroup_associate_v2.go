package openstack

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceNetworkingPortSecGroupAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingPortSecGroupAssociateV2Create,
		ReadContext:   resourceNetworkingPortSecGroupAssociateV2Read,
		UpdateContext: resourceNetworkingPortSecGroupAssociateV2Update,
		DeleteContext: resourceNetworkingPortSecGroupAssociateV2Delete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"enforce": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"all_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceNetworkingPortSecGroupAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	portID := d.Get("port_id").(string)

	port, err := ports.Get(networkingClient, portID).Extract()
	if err != nil {
		return diag.Errorf("Unable to get %s Port: %s", portID, err)
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", portID, port)

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOkExists("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		updateOpts.SecurityGroups = &securityGroups
	} else {
		// append security groups
		sg := sliceUnion(port.SecurityGroups, securityGroups)
		updateOpts.SecurityGroups = &sg
	}

	log.Printf("[DEBUG] Port Security Group Associate Options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(networkingClient, portID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error associating %s port with '%s' security groups: %s", portID, strings.Join(securityGroups, ","), err)
	}

	d.SetId(portID)

	return resourceNetworkingPortSecGroupAssociateV2Read(ctx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	port, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error fetching port security groups"))
	}

	var enforce bool
	if v, ok := d.GetOkExists("enforce"); ok {
		enforce = v.(bool)
	}

	d.Set("all_security_group_ids", port.SecurityGroups)

	if enforce {
		d.Set("security_group_ids", port.SecurityGroups)
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		desiredSet := d.Get("security_group_ids").(*schema.Set)
		actualSet := allSet.Intersection(desiredSet)
		if !actualSet.Equal(desiredSet) {
			d.Set("security_group_ids", expandToStringSlice(actualSet.List()))
		}
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingPortSecGroupAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOkExists("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
		updateOpts.SecurityGroups = &securityGroups
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldIDs, newIDs := d.GetChange("security_group_ids")
		oldSet, newSet := oldIDs.(*schema.Set), newIDs.(*schema.Set)

		allWithoutOld := allSet.Difference(oldSet)

		newSecurityGroups := expandToStringSlice(allWithoutOld.Union(newSet).List())

		updateOpts.SecurityGroups = &newSecurityGroups
	}

	if d.HasChange("security_group_ids") || d.HasChange("enforce") {
		log.Printf("[DEBUG] Port Security Group Update Options: %#v", updateOpts.SecurityGroups)
		_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenStack Neutron Port: %s", err)
		}
	}

	return resourceNetworkingPortSecGroupAssociateV2Read(ctx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOkExists("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		updateOpts.SecurityGroups = &[]string{}
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldSet := d.Get("security_group_ids").(*schema.Set)

		allWithoutOld := allSet.Difference(oldSet)

		newSecurityGroups := expandToStringSlice(allWithoutOld.List())

		updateOpts.SecurityGroups = &newSecurityGroups
	}

	log.Printf("[DEBUG] Port security groups disassociation options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error disassociating port security groups"))
	}

	return nil
}
