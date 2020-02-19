package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/drivers"
	noauthintrospection "github.com/gophercloud/gophercloud/openstack/baremetalintrospection/noauth"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"net/http"
	"sync"
	"time"
)

// Clients stores the client connection information for Ironic and Inspector
type Clients struct {
	ironic    *gophercloud.ServiceClient
	inspector *gophercloud.ServiceClient

	// Boolean that determines if Ironic API was previously determined to be available, we don't need to try every time.
	ironicUp bool

	// Boolean that determines we've already waited and the API never came up, we don't need to wait again.
	ironicFailed bool

	// Mutex so that only one resource being created by terraform checks at a time. There's no reason to have multiple
	// resources calling out to the API.
	ironicMux sync.Mutex

	// Boolean that determines if Inspector API was previously determined to be available, we don't need to try every time.
	inspectorUp bool

	// Boolean that determines that we've already waited, and inspector API did not come up.
	inspectorFailed bool

	// Mutex so that only one resource being created by terraform checks at a time. There's no reason to have multiple
	// resources calling out to the API.
	inspectorMux sync.Mutex

	timeout int
}

// GetIronicClient returns the API client for Ironic, optionally retrying to reach the API if timeout is set.
func (c *Clients) GetIronicClient() (*gophercloud.ServiceClient, error) {
	c.ironicMux.Lock()
	defer c.ironicMux.Unlock()

	// Ironic is UP, or user didn't ask us to check
	if c.ironicUp || c.timeout == 0 {
		return c.ironic, nil
	}

	if c.ironicFailed {
		return nil, fmt.Errorf("could not contact API: timeout reached")
	}

	// Wait the specified interval for Ironic to come up
	err := waitForAPI(c.timeout, c.ironic)
	if err != nil {
		c.ironicFailed = true
		return nil, err

	}

	c.ironicUp = true
	return c.ironic, nil
}

// GetInspectorClient returns the API client for Ironic, optionally retrying to reach the API if timeout is set.
func (c *Clients) GetInspectorClient() (*gophercloud.ServiceClient, error) {
	c.inspectorMux.Lock()
	defer c.inspectorMux.Unlock()

	if c.inspector == nil {
		return nil, fmt.Errorf("no inspector endpoint was specified")
	} else if c.inspectorUp || c.timeout == 0 {
		return c.inspector, nil
	} else if c.inspectorFailed {
		return nil, fmt.Errorf("could not contact API: timeout reached")
	}

	err := waitForAPI(c.timeout, c.inspector)
	if err != nil {
		c.inspectorFailed = true
		return nil, err
	}

	c.inspectorUp = true
	return c.inspector, nil
}

// Provider Ironic
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_ENDPOINT", ""),
				Description: descriptions["url"],
			},
			"inspector": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_INSPECTOR_ENDPOINT", ""),
				Description: descriptions["inspector"],
			},
			"microversion": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_MICROVERSION", "1.52"),
				Description: descriptions["microversion"],
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["timeout"],
				Default:     0,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ironic_node_v1":       resourceNodeV1(),
			"ironic_port_v1":       resourcePortV1(),
			"ironic_allocation_v1": resourceAllocationV1(),
			"ironic_deployment":    resourceDeployment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ironic_introspection": dataSourceIronicIntrospection(),
		},
		ConfigureFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"url":          "The authentication endpoint for Ironic",
		"inspector":    "The endpoint for Ironic inspector",
		"microversion": "The microversion to use for Ironic",
		"timeout":      "Wait at least the specified number of seconds for the API to become available",
	}
}

// Creates a noauth Ironic client
func configureProvider(schema *schema.ResourceData) (interface{}, error) {
	var clients Clients

	url := schema.Get("url").(string)
	if url == "" {
		return nil, fmt.Errorf("url is required for ironic provider")
	}
	log.Printf("[DEBUG] Ironic endpoint is %s", url)

	ironic, err := noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: url,
	})
	if err != nil {
		return nil, err
	}
	ironic.Microversion = schema.Get("microversion").(string)
	clients.ironic = ironic

	inspectorURL := schema.Get("inspector").(string)
	if inspectorURL != "" {
		log.Printf("[DEBUG] Inspector endpoint is %s", inspectorURL)
		inspector, err := noauthintrospection.NewBareMetalIntrospectionNoAuth(noauthintrospection.EndpointOpts{
			IronicInspectorEndpoint: inspectorURL,
		})
		if err != nil {
			return nil, fmt.Errorf("could not configure inspector endpoint: %s", err.Error())
		}
		clients.inspector = inspector
	}

	clients.timeout = schema.Get("timeout").(int)

	return &clients, err
}

// Retries an API until a timeout is reached. The timeout is approximate, we calculate approximately
// the number of attempts that would equal that timeout. It's not exact, but it will be *at least* the
// value specified.
func waitForAPI(timeout int, client *gophercloud.ServiceClient) error {
	var maxTries int
	if timeout < 5 {
		maxTries = 1
	} else {
		maxTries = timeout / 5
	}
	log.Printf("[DEBUG] Waiting for Ironic API to become available - max attempts is %d", maxTries)

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	for tries := 1; tries <= maxTries; tries++ {
		// FIXME(stbenjam): After we sleep at the end of the loop, *something* is changing log output to be discarded, and
		// when we get woken up again, our logs aren't printed anymore until the end of the loop. Very odd error, possibly
		// something in terraform is doing this.
		log.Printf("[DEBUG] Trying to connect to API, attempt %d of %d\n", tries, maxTries)

		r, err := httpClient.Get(client.Endpoint)
		if err == nil {
			statusCode := r.StatusCode
			r.Body.Close()
			if statusCode == http.StatusOK {
				log.Printf("[DEBUG] API successfully connected, waiting for conductor...")
				driverCount := 0
				drivers.ListDrivers(client, drivers.ListDriversOpts{
					Detail: false,
				}).EachPage(func(page pagination.Page) (bool, error) {
					actual, err := drivers.ExtractDrivers(page)
					if err != nil {
						return false, err
					}
					driverCount += len(actual)
					return true, nil
				})
				// If we have any drivers, conductor is up.
				if driverCount > 0 {
					return nil
				}
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("could not contact API: timeout reached")
}
