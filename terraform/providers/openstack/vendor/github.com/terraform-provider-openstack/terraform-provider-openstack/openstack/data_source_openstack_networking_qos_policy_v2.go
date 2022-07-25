package openstack

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
)

func dataSourceNetworkingQoSPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingQoSPolicyV2Read,
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
				Computed: true,
				ForceNew: false,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"revision_number": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: false,
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

func dataSourceNetworkingQoSPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := policies.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("shared"); ok {
		shared := v.(bool)
		listOpts.Shared = &shared
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		isDefault := v.(bool)
		listOpts.IsDefault = &isDefault
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := policies.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_qos_policy_v2: %s", err)
	}

	allPolicies, err := policies.ExtractPolicies(pages)
	if err != nil {
		return diag.Errorf("Unable to extract openstack_networking_qos_policy_v2: %s", err)
	}

	if len(allPolicies) < 1 {
		return diag.Errorf("Your query returned no openstack_networking_qos_policy_v2. " +
			"Please change your search criteria and try again.")
	}

	if len(allPolicies) > 1 {
		return diag.Errorf("Your query returned more than one openstack_networking_qos_policy_v2." +
			" Please try a more specific search criteria")
	}

	policy := allPolicies[0]

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_policy_v2 %s: %+v", policy.ID, policy)
	d.SetId(policy.ID)

	d.Set("name", policy.Name)
	d.Set("project_id", policy.ProjectID)
	d.Set("shared", policy.Shared)
	d.Set("is_default", policy.IsDefault)
	d.Set("description", policy.Description)
	d.Set("revision_number", policy.RevisionNumber)
	d.Set("region", GetRegion(d, config))

	if err := d.Set("created_at", policy.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_qos_policy_v2 created_at: %s", err)
	}
	if err := d.Set("updated_at", policy.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_qos_policy_v2 updated_at: %s", err)
	}

	return nil
}
