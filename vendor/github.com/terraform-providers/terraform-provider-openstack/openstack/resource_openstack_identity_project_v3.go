package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProjectV3Create,
		Read:   resourceIdentityProjectV3Read,
		Update: resourceIdentityProjectV3Update,
		Delete: resourceIdentityProjectV3Delete,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_domain": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIdentityProjectV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	isDomain := d.Get("is_domain").(bool)
	createOpts := projects.CreateOpts{
		Description: d.Get("description").(string),
		DomainID:    d.Get("domain_id").(string),
		Enabled:     &enabled,
		IsDomain:    &isDomain,
		Name:        d.Get("name").(string),
		ParentID:    d.Get("parent_id").(string),
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = expandToStringSlice(tags)
	}

	log.Printf("[DEBUG] openstack_identity_project_v3 create options: %#v", createOpts)
	project, err := projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_identity_project_v3: %s", err)
	}

	d.SetId(project.ID)

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	project, err := projects.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_identity_project_v3")
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_project_v3 %s: %#v", d.Id(), project)

	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("is_domain", project.IsDomain)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)
	d.Set("region", GetRegion(d, config))
	d.Set("tags", project.Tags)

	return nil
}

func resourceIdentityProjectV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts projects.UpdateOpts

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("is_domain") {
		hasChange = true
		isDomain := d.Get("is_domain").(bool)
		updateOpts.IsDomain = &isDomain
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("parent_id") {
		hasChange = true
		updateOpts.ParentID = d.Get("parent_id").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("tags") {
		hasChange = true
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := expandToStringSlice(tags)
			updateOpts.Tags = &tagsToUpdate
		} else {
			updateOpts.Tags = &[]string{}
		}
	}

	if hasChange {
		_, err := projects.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_identity_project_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = projects.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_identity_project_v3")
	}

	return nil
}
