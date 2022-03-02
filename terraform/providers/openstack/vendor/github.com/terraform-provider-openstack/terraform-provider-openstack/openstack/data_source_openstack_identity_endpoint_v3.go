package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
)

func dataSourceIdentityEndpointV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityEndpointV3Read,

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

			"endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"interface": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public", "internal", "admin",
				}, false),
			},

			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityEndpointV3Read performs the endpoint lookup.
func dataSourceIdentityEndpointV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	listOpts := endpoints.ListOpts{
		Availability: identityEndpointAvailability(d.Get("interface").(string)),
		ServiceID:    d.Get("service_id").(string),
		RegionID:     d.Get("endpoint_region").(string),
	}

	log.Printf("[DEBUG] openstack_identity_endpoint_v3 list options: %#v", listOpts)

	var endpoint endpoints.Endpoint
	allPages, err := endpoints.List(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_endpoint_v3: %s", err)
	}

	allEndpoints, err := endpoints.ExtractEndpoints(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_identity_endpoint_v3: %s", err)
	}

	// filter by name, when the name is specified
	if v, ok := d.GetOkExists("name"); ok {
		var filteredEndpoints []endpoints.Endpoint
		for _, endpoint := range allEndpoints {
			if endpoint.Name == v.(string) {
				filteredEndpoints = append(filteredEndpoints, endpoint)
			}
		}
		allEndpoints = filteredEndpoints
	}

	if len(allEndpoints) < 1 {
		return diag.Errorf("Your openstack_identity_endpoint_v3 query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// Query services
	serviceType := d.Get("service_type").(string)
	serviceName := d.Get("service_name").(string)
	var filteredEndpoints []endpoints.Endpoint
	allServicePages, err := services.List(identityClient, services.ListOpts{ServiceType: serviceType, Name: serviceName}).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_endpoint_v3 services: %s", err)
	}

	allServices, err := services.ExtractServices(allServicePages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_identity_endpoint_v3 services: %s", err)
	}

	for _, endpoint := range allEndpoints {
		for _, service := range allServices {
			if endpoint.ServiceID == service.ID {
				// it is safe to assign these vars here, since if there are more than
				// one endpoint, an error will be returned
				if v, ok := service.Extra["name"].(string); ok {
					serviceName = v
				}
				serviceType = service.Type
				filteredEndpoints = append(filteredEndpoints, endpoint)
			}
		}
	}

	allEndpoints = filteredEndpoints

	if len(allEndpoints) < 1 {
		return diag.Errorf("Your openstack_identity_endpoint_v3 query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEndpoints) > 1 {
		return diag.Errorf("Your openstack_identity_endpoint_v3 query returned more than one result")
	}
	endpoint = allEndpoints[0]

	log.Printf("[DEBUG] openstack_identity_endpoint_v3 details: %#v", endpoint)

	d.SetId(endpoint.ID)

	d.Set("name", endpoint.Name)
	d.Set("interface", string(endpoint.Availability))
	d.Set("endpoint_region", endpoint.Region)
	d.Set("service_id", endpoint.ServiceID)
	d.Set("service_name", serviceName)
	d.Set("service_type", serviceType)
	d.Set("url", endpoint.URL)

	d.Set("region", GetRegion(d, config))

	return nil
}
