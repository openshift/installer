package openstack

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/pagination"
)

func resourceIdentityEndpointV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityEndpointV3Create,
		ReadContext:   resourceIdentityEndpointV3Read,
		UpdateContext: resourceIdentityEndpointV3Update,
		DeleteContext: resourceIdentityEndpointV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Optional: true,
			},

			"endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
			},

			"url": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					_, err := url.ParseRequestURI(value)
					if err != nil {
						errors = append(errors, fmt.Errorf("URL is not valid: %s", err))
					}
					return
				},
			},

			"interface": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"admin", "internal", "public",
				}, false),
			},

			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityEndpointV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := endpoints.CreateOpts{
		Name:         d.Get("name").(string),
		Availability: identityEndpointAvailability(d.Get("interface").(string)),
		Region:       d.Get("endpoint_region").(string),
		URL:          d.Get("url").(string),
		ServiceID:    d.Get("service_id").(string),
	}

	log.Printf("[DEBUG] openstack_identity_endpoint_v3 create options: %#v", createOpts)
	endpoint, err := endpoints.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_endpoint_v3: %s", err)
	}

	d.SetId(endpoint.ID)

	return resourceIdentityEndpointV3Read(ctx, d, meta)
}

func resourceIdentityEndpointV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var endpoint endpoints.Endpoint
	err = endpoints.List(identityClient, nil).EachPage(func(page pagination.Page) (bool, error) {
		endpointList, err := endpoints.ExtractEndpoints(page)
		if err != nil {
			return false, err
		}
		for _, v := range endpointList {
			if v.ID == d.Id() {
				endpoint = v
				break
			}
		}
		return true, nil
	})

	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_endpoint_v3"))
	}

	if endpoint == (endpoints.Endpoint{}) {
		// Endpoint was not found
		d.SetId("")
		return nil
	}

	// Query services
	serviceType := d.Get("service_type").(string)
	serviceName := d.Get("service_name").(string)
	allServicePages, err := services.List(identityClient, services.ListOpts{ServiceType: serviceType, Name: serviceName}).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_endpoint_v3 services: %s", err)
	}

	allServices, err := services.ExtractServices(allServicePages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_identity_endpoint_v3 services: %s", err)
	}

	for _, service := range allServices {
		if endpoint.ServiceID == service.ID {
			if v, ok := service.Extra["name"].(string); ok {
				serviceName = v
			}
			serviceType = service.Type
		}
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_endpoint_v3: %#v", endpoint)

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

func resourceIdentityEndpointV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts endpoints.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Region = d.Get("name").(string)
	}

	if d.HasChange("endpoint_region") {
		hasChange = true
		updateOpts.Region = d.Get("endpoint_region").(string)
	}

	if d.HasChange("url") {
		hasChange = true
		updateOpts.URL = d.Get("url").(string)
	}

	if d.HasChange("service_id") {
		hasChange = true
		updateOpts.ServiceID = d.Get("service_id").(string)
	}

	if d.HasChange("interface") {
		hasChange = true

		updateOpts.Availability = identityEndpointAvailability(d.Get("interface").(string))
	}

	if hasChange {
		_, err := endpoints.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_endpoint_v3: %s", err)
		}
	}

	return resourceIdentityEndpointV3Read(ctx, d, meta)
}

func resourceIdentityEndpointV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = endpoints.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting openstack_identity_endpoint_v3: %s", err)
	}

	return nil
}
