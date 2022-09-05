package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsWorkspaceCreateUpdate,
		Read:   resourceLogAnalyticsWorkspaceRead,
		Update: resourceLogAnalyticsWorkspaceCreateUpdate,
		Delete: resourceLogAnalyticsWorkspaceDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceLogAnalyticsWorkspaceCustomDiff),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workspaces.ParseWorkspaceID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
			1: migration.WorkspaceV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"cmk_for_query_forced": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"internet_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"internet_query_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// TODO 4.0: Clean up lacluster "workaround" to make it more readable and easier to understand. (@WodansSon already has the code written for the clean up)
			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(workspaces.WorkspaceSkuNameEnumFree),
					string(workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight),
					string(workspaces.WorkspaceSkuNameEnumPerNode),
					string(workspaces.WorkspaceSkuNameEnumPremium),
					string(workspaces.WorkspaceSkuNameEnumStandalone),
					string(workspaces.WorkspaceSkuNameEnumStandard),
					string(workspaces.WorkspaceSkuNameEnumCapacityReservation),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, false),
			},

			"reservation_capacity_in_gb_per_day": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.IntBetween(100, 5000), validation.IntDivisibleBy(100)),
			},

			"retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
			},

			"daily_quota_gb": {
				Type:             pluginsdk.TypeFloat,
				Optional:         true,
				Default:          -1.0,
				DiffSuppressFunc: dailyQuotaGbDiffSuppressFunc,
				ValidateFunc:     validation.FloatAtLeast(-1.0),
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceLogAnalyticsWorkspaceCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	// Since sku needs to be a force new if the sku changes we need to have this
	// custom diff here because when you link the workspace to a cluster the
	// cluster changes the sku to LACluster, so we need to ignore the change
	// if it is LACluster else invoke the ForceNew as before...
	//
	// NOTE: Since LACluster is not in our enum the value is returned as ""
	if d.HasChange("sku") {
		old, new := d.GetChange("sku")
		log.Printf("[INFO] Log Analytics Workspace SKU: OLD: %q, NEW: %q", old, new)
		// If the old value is not LACluster(e.g. "") return ForceNew because they are
		// really changing the sku...
		if !strings.EqualFold(old.(string), "") {
			d.ForceNew("sku")
		}
	}

	return nil
}

func resourceLogAnalyticsWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace creation.")

	var isLACluster bool
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := workspaces.NewWorkspaceID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Log Analytics Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	skuName := d.Get("sku").(string)
	sku := &workspaces.WorkspaceSku{
		Name: workspaces.WorkspaceSkuNameEnum(skuName),
	}

	// (@WodansSon) - If the workspace is connected to a cluster via the linked service resource
	// the workspace SKU cannot be modified since the linked service owns the sku value within
	// the workspace once it is linked
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err == nil {
			if resp.Model != nil && resp.Model.Properties != nil {
				if azSku := resp.Model.Properties.Sku; azSku != nil {
					if strings.EqualFold(string(azSku.Name), "lacluster") {
						isLACluster = true
						log.Printf("[INFO] Log Analytics Workspace %q (Resource Group %q): SKU is linked to Log Analytics cluster", name, resourceGroup)
					}
				}
			}
		}
	}

	internetIngestionEnabled := workspaces.PublicNetworkAccessTypeDisabled
	if d.Get("internet_ingestion_enabled").(bool) {
		internetIngestionEnabled = workspaces.PublicNetworkAccessTypeEnabled
	}
	internetQueryEnabled := workspaces.PublicNetworkAccessTypeDisabled
	if d.Get("internet_query_enabled").(bool) {
		internetQueryEnabled = workspaces.PublicNetworkAccessTypeEnabled
	}

	retentionInDays := int64(d.Get("retention_in_days").(int))

	t := d.Get("tags").(map[string]interface{})

	if isLACluster {
		sku.Name = "lacluster"
	} else if skuName == "" {
		// Default value if sku is not defined
		sku.Name = workspaces.WorkspaceSkuNameEnumPerGBTwoZeroOneEight
	}

	parameters := workspaces.Workspace{
		Name:     &name,
		Location: location,
		Tags:     tags.Expand(t),
		Properties: &workspaces.WorkspaceProperties{
			Sku:                             sku,
			PublicNetworkAccessForIngestion: &internetIngestionEnabled,
			PublicNetworkAccessForQuery:     &internetQueryEnabled,
			RetentionInDays:                 &retentionInDays,
		},
	}

	if v, ok := d.GetOkExists("cmk_for_query_forced"); ok {
		parameters.Properties.ForceCmkForQuery = utils.Bool(v.(bool))
	}

	dailyQuotaGb, ok := d.GetOk("daily_quota_gb")
	if ok && strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) && (dailyQuotaGb != -1 && dailyQuotaGb != 0.5) {
		return fmt.Errorf("`Free` tier SKU quota is not configurable and is hard set to 0.5GB")
	} else if !strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) {
		parameters.Properties.WorkspaceCapping = &workspaces.WorkspaceCapping{
			DailyQuotaGb: utils.Float(dailyQuotaGb.(float64)),
		}
	}

	propName := "reservation_capacity_in_gb_per_day"
	capacityReservationLevel, ok := d.GetOk(propName)
	if ok {
		if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			parameters.Properties.Sku.CapacityReservationLevel = utils.Int64((int64(capacityReservationLevel.(int))))
		} else {
			return fmt.Errorf("`%s` can only be used with the `CapacityReservation` SKU", propName)
		}
	} else {
		if strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumCapacityReservation)) {
			return fmt.Errorf("`%s` must be set when using the `CapacityReservation` SKU", propName)
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsWorkspaceRead(d, meta)
}

func resourceLogAnalyticsWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	sharedKeysClient := meta.(*clients.Client).LogAnalytics.SharedKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspaces '%s': %+v", id.WorkspaceName, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if props := model.Properties; props != nil {
			d.Set("internet_ingestion_enabled", *props.PublicNetworkAccessForIngestion == workspaces.PublicNetworkAccessTypeEnabled)
			d.Set("internet_query_enabled", *props.PublicNetworkAccessForQuery == workspaces.PublicNetworkAccessTypeEnabled)
			d.Set("workspace_id", props.CustomerId)

			skuName := ""
			if sku := props.Sku; sku != nil {
				for _, v := range workspaces.PossibleValuesForWorkspaceSkuNameEnum() {
					if strings.EqualFold(v, string(sku.Name)) {
						skuName = v
					}
				}

				if capacityReservationLevel := sku.CapacityReservationLevel; capacityReservationLevel != nil {
					d.Set("reservation_capacity_in_gb_per_day", capacityReservationLevel)
				}
			}
			d.Set("sku", skuName)

			d.Set("cmk_for_query_forced", props.ForceCmkForQuery)
			d.Set("retention_in_days", props.RetentionInDays)

			if props.Sku != nil && strings.EqualFold(string(props.Sku.Name), string(workspaces.WorkspaceSkuNameEnumFree)) {
				// Special case for "Free" tier
				d.Set("daily_quota_gb", utils.Float(0.5))
			} else if workspaceCapping := props.WorkspaceCapping; workspaceCapping != nil {
				d.Set("daily_quota_gb", props.WorkspaceCapping.DailyQuotaGb)
			} else {
				d.Set("daily_quota_gb", utils.Float(-1))
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	sharedKeys, err := sharedKeysClient.GetSharedKeys(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", id.WorkspaceName, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return nil
}

func resourceLogAnalyticsWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}
	PermanentlyDeleteOnDestroy := meta.(*clients.Client).Features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy
	parameters := workspaces.DeleteOperationOptions{
		Force: utils.Bool(PermanentlyDeleteOnDestroy),
	}
	if err := client.DeleteThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("issuing AzureRM delete request for Log Analytics Workspaces '%s': %+v", id.WorkspaceName, err)
	}

	return nil
}

func dailyQuotaGbDiffSuppressFunc(_, _, _ string, d *pluginsdk.ResourceData) bool {
	// (@jackofallops) - 'free' is a legacy special case that is always set to 0.5GB
	if skuName := d.Get("sku").(string); strings.EqualFold(skuName, string(workspaces.WorkspaceSkuNameEnumFree)) {
		return true
	}

	return false
}
