package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
)

func resourceNetworkingQoSPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingQoSPolicyV2Create,
		ReadContext:   resourceNetworkingQoSPolicyV2Read,
		UpdateContext: resourceNetworkingQoSPolicyV2Update,
		DeleteContext: resourceNetworkingQoSPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"revision_number": {
				Type:     schema.TypeInt,
				ForceNew: false,
				Computed: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkingQoSPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := QoSPolicyCreateOpts{
		policies.CreateOpts{
			Name:        d.Get("name").(string),
			ProjectID:   d.Get("project_id").(string),
			Shared:      d.Get("shared").(bool),
			Description: d.Get("description").(string),
			IsDefault:   d.Get("is_default").(bool),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] openstack_networking_qos_policy_v2 create options: %#v", createOpts)
	p, err := policies.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_qos_policy_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_qos_policy_v2 %s to become available.", p.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingQoSPolicyV2StateRefreshFunc(networkingClient, p.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_policy_v2 %s to become available: %s", p.ID, err)
	}

	d.SetId(p.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "qos/policies", p.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_qos_policy_v2 %s: %s", p.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_qos_policy_v2 %s", tags, p.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_qos_policy_v2 %s: %#v", p.ID, p)

	return resourceNetworkingQoSPolicyV2Read(ctx, d, meta)
}

func resourceNetworkingQoSPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	p, err := policies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_qos_policy_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_policy_v2 %s: %#v", d.Id(), p)

	d.Set("name", p.Name)
	d.Set("project_id", p.ProjectID)
	d.Set("shared", p.Shared)
	d.Set("is_default", p.IsDefault)
	d.Set("description", p.Description)
	d.Set("revision_number", p.RevisionNumber)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, p.Tags)

	if err := d.Set("created_at", p.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_qos_policy_v2 created_at: %s", err)
	}
	if err := d.Set("updated_at", p.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_qos_policy_v2 updated_at: %s", err)
	}

	return nil
}

func resourceNetworkingQoSPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var hasChange bool
	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("shared") {
		hasChange = true
		v := d.Get("shared").(bool)
		updateOpts.Shared = &v
	}

	if d.HasChange("description") {
		hasChange = true
		v := d.Get("description").(string)
		updateOpts.Description = &v
	}

	if d.HasChange("is_default") {
		hasChange = true
		v := d.Get("is_default").(bool)
		updateOpts.IsDefault = &v
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_qos_policy_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = policies.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_qos_policy_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "qos/policies", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_qos_policy_v2 %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_qos_policy_v2 %s", tags, d.Id())
	}

	return resourceNetworkingQoSPolicyV2Read(ctx, d, meta)
}

func resourceNetworkingQoSPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := policies.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_qos_policy_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingQoSPolicyV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_qos_policy_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
