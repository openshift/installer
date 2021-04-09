package applicationinsights

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApplicationInsightsWebTests() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationInsightsWebTestsCreateUpdate,
		Read:   resourceApplicationInsightsWebTestsRead,
		Update: resourceApplicationInsightsWebTestsCreateUpdate,
		Delete: resourceApplicationInsightsWebTestsDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.WebTestID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"location": azure.SchemaLocation(),

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.Multistep),
					string(insights.Ping),
				}, false),
			},

			"frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
				ValidateFunc: validation.IntInSlice([]int{
					300,
					600,
					900,
				}),
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"retry_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"geo_locations": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateFunc:     validation.StringIsNotEmpty,
					StateFunc:        location.StateFunc,
					DiffSuppressFunc: location.DiffSuppressFunc,
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"configuration": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.XmlDiff,
			},

			"tags": tags.Schema(),

			"synthetic_monitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationInsightsWebTestsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights WebTest creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	appInsightsId, err := parse.ComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Application Insights WebTests %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_insights_web_test", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	description := d.Get("description").(string)
	frequency := int32(d.Get("frequency").(int))
	timeout := int32(d.Get("timeout").(int))
	isEnabled := d.Get("enabled").(bool)
	retryEnabled := d.Get("retry_enabled").(bool)
	geoLocationsRaw := d.Get("geo_locations").([]interface{})
	geoLocations := expandApplicationInsightsWebTestGeoLocations(geoLocationsRaw)
	testConf := d.Get("configuration").(string)

	t := d.Get("tags").(map[string]interface{})
	tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
	t[tagKey] = "Resource"

	webTest := insights.WebTest{
		Name:     &name,
		Location: &location,
		Kind:     insights.WebTestKind(kind),
		WebTestProperties: &insights.WebTestProperties{
			SyntheticMonitorID: &name,
			WebTestName:        &name,
			Description:        &description,
			Enabled:            &isEnabled,
			Frequency:          &frequency,
			Timeout:            &timeout,
			WebTestKind:        insights.WebTestKind(kind),
			RetryEnabled:       &retryEnabled,
			Locations:          &geoLocations,
			Configuration: &insights.WebTestPropertiesConfiguration{
				WebTest: &testConf,
			},
		},
		Tags: tags.Expand(t),
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, webTest)
	if err != nil {
		return fmt.Errorf("creating/updating Application Insights WebTest %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceApplicationInsightsWebTestsRead(d, meta)
}

func resourceApplicationInsightsWebTestsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebTestID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights WebTests %q", id)

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Application Insights WebTest %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Application Insights WebTests %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	appInsightsId := ""
	for i := range resp.Tags {
		if strings.HasPrefix(i, "hidden-link") {
			appInsightsId = strings.Split(i, ":")[1]
		}
	}
	d.Set("application_insights_id", appInsightsId)
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WebTestProperties; props != nil {
		// It is possible that the root level `kind` in response is empty in some cases (see PR #8372 for more info)
		if resp.Kind == "" {
			d.Set("kind", props.WebTestKind)
		}
		d.Set("synthetic_monitor_id", props.SyntheticMonitorID)
		d.Set("description", props.Description)
		d.Set("enabled", props.Enabled)
		d.Set("frequency", props.Frequency)
		d.Set("timeout", props.Timeout)
		d.Set("retry_enabled", props.RetryEnabled)

		if config := props.Configuration; config != nil {
			d.Set("configuration", config.WebTest)
		}

		if err := d.Set("geo_locations", flattenApplicationInsightsWebTestGeoLocations(props.Locations)); err != nil {
			return fmt.Errorf("Error setting `geo_locations`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceApplicationInsightsWebTestsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebTestID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting AzureRM Application Insights WebTest '%s' (resource group '%s')", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights WebTest '%s': %+v", id.Name, err)
	}

	return err
}

func expandApplicationInsightsWebTestGeoLocations(input []interface{}) []insights.WebTestGeolocation {
	locations := make([]insights.WebTestGeolocation, 0)

	for _, v := range input {
		lc := v.(string)
		loc := insights.WebTestGeolocation{
			Location: &lc,
		}
		locations = append(locations, loc)
	}

	return locations
}

func flattenApplicationInsightsWebTestGeoLocations(input *[]insights.WebTestGeolocation) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		if prop.Location != nil {
			results = append(results, azure.NormalizeLocation(*prop.Location))
		}
	}

	return results
}
