package maintenance

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2021-05-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmMaintenanceConfigurationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"visibility": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"window": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"expiration_date_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"recur_every": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"properties": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := maintenanceconfigurations.NewMaintenanceConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.ResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("scope", props.MaintenanceScope)
			d.Set("visibility", props.Visibility)
			d.Set("properties", props.ExtensionProperties)

			window := flattenMaintenanceConfigurationWindow(props.MaintenanceWindow)
			if err := d.Set("window", window); err != nil {
				return fmt.Errorf("setting `window`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(model.Location))
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}
