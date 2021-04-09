package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceLogAnalyticsWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogAnalyticsWorkspaceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"retention_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"daily_quota_gb": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "this property has been removed from the API and will be removed in version 3.0 of the provider",
			},

			"primary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceLogAnalyticsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	sharedKeysClient := meta.(*clients.Client).LogAnalytics.SharedKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Log Analytics workspaces %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics workspaces '%s': %+v", name, err)
	}

	id, err := parse.LogAnalyticsWorkspaceID(*resp.ID)
	if err != nil {
		return fmt.Errorf("Error parsing Log Analytics Workspace ID %q", *resp.ID)
	}
	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", "")
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)

	if workspaceCapping := resp.WorkspaceCapping; workspaceCapping != nil {
		d.Set("daily_quota_gb", resp.WorkspaceCapping.DailyQuotaGb)
	} else {
		d.Set("daily_quota_gb", utils.Float(-1))
	}

	sharedKeys, err := sharedKeysClient.GetSharedKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
