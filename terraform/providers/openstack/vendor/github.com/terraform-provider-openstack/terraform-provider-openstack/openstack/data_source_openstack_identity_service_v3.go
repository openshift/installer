package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
)

func dataSourceIdentityServiceV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityServiceV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityServiceV3Read performs the service lookup.
func dataSourceIdentityServiceV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	name := d.Get("name").(string)
	listOpts := services.ListOpts{
		Name:        name,
		ServiceType: d.Get("type").(string),
	}

	log.Printf("[DEBUG] openstack_identity_service_v3 list options: %#v", listOpts)

	var service services.Service
	allPages, err := services.List(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_service_v3: %s", err)
	}

	allServices, err := services.ExtractServices(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_identity_service_v3: %s", err)
	}

	// filter by enabled, when the enabled is specified
	if v, ok := d.GetOkExists("enabled"); ok {
		var filteredServices []services.Service
		for _, svc := range allServices {
			if svc.Enabled == v.(bool) {
				filteredServices = append(filteredServices, svc)
			}
		}
		allServices = filteredServices
	}

	if len(allServices) < 1 {
		return diag.Errorf("Your openstack_identity_service_v3 query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allServices) > 1 {
		return diag.Errorf("Your openstack_identity_service_v3 query returned more than one result")
	}
	service = allServices[0]

	description := ""
	if v, ok := service.Extra["name"].(string); ok {
		name = v
	}
	if v, ok := service.Extra["description"].(string); ok {
		description = v
	}

	log.Printf("[DEBUG] openstack_identity_service_v3 details: %#v", service)

	d.SetId(service.ID)

	d.Set("name", name)
	d.Set("type", service.Type)
	d.Set("description", description)
	d.Set("enabled", service.Enabled)

	d.Set("region", GetRegion(d, config))

	return nil
}
