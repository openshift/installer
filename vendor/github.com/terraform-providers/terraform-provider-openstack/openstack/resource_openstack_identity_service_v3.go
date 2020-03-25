package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIdentityServiceV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityServiceV3Create,
		Read:   resourceIdentityServiceV3Read,
		Update: resourceIdentityServiceV3Update,
		Delete: resourceIdentityServiceV3Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceIdentityServiceV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	createOpts := services.CreateOpts{
		Extra: map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
		Type:    d.Get("type").(string),
		Enabled: &enabled,
	}

	log.Printf("[DEBUG] openstack_identity_service_v3 create options: %#v", createOpts)
	service, err := services.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_identity_service_v3: %s", err)
	}

	d.SetId(service.ID)

	return resourceIdentityServiceV3Read(d, meta)
}

func resourceIdentityServiceV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	service, err := services.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_identity_service_v3")
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_service_v3: %#v", service)

	name := ""
	description := ""
	if v, ok := service.Extra["name"].(string); ok {
		name = v
	}
	if v, ok := service.Extra["description"].(string); ok {
		description = v
	}

	d.Set("name", name)
	d.Set("description", description)
	d.Set("type", service.Type)
	d.Set("enabled", service.Enabled)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityServiceV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var updateOpts services.UpdateOpts

	// these options must always be set
	enabled := d.Get("enabled").(bool)
	updateOpts.Enabled = &enabled
	updateOpts.Type = d.Get("type").(string)

	if d.HasChange("name") || d.HasChange("description") {
		updateOpts.Extra = map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		}
	}

	_, err = services.Update(identityClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating openstack_identity_service_v3: %s", err)
	}

	return resourceIdentityServiceV3Read(d, meta)
}

func resourceIdentityServiceV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = services.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting openstack_identity_service_v3: %s", err)
	}

	return nil
}
