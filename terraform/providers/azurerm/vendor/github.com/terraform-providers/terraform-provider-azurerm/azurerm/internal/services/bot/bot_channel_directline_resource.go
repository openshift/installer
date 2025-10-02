package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBotChannelDirectline() *schema.Resource {
	return &schema.Resource{
		Create: resourceBotChannelDirectlineCreate,
		Read:   resourceBotChannelDirectlineRead,
		Delete: resourceBotChannelDirectlineDelete,
		Update: resourceBotChannelDirectlineUpdate,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BotChannelID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"site": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"v1_allowed": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"v3_allowed": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"enhanced_authentication_enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},

						"trusted_origins": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"key2": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceBotChannelDirectlineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameDirectLineChannel1))
	existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ChannelName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		// a "Default Site" site gets created and returned.. so let's check it's not just that
		if props := existing.Properties; props != nil {
			directLineChannel, ok := props.AsDirectLineChannel()
			if ok && directLineChannel.Properties != nil {
				sites := filterSites(directLineChannel.Properties.Sites)
				if len(sites) != 0 {
					return tf.ImportAsExistsError("azurerm_bot_channel_directline", resourceId.ID())
				}
			}
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineChannel{
			Properties: &botservice.DirectLineChannelProperties{
				Sites: expandDirectlineSites(d.Get("site").(*schema.Set).List()),
			},
			ChannelName: botservice.ChannelNameDirectLineChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("creating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}
	d.SetId(resourceId.ID())

	// Unable to create a new site with enhanced_authentication_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	return resourceBotChannelDirectlineRead(d, meta)
}

func resourceBotChannelDirectlineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineChannel1))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Directline Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroup, id.BotServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Channel Directline for Bot %q (Resource Group %q): %+v", id.ResourceGroup, id.BotServiceName, err)
	}

	channelsResp, err := client.ListWithKeys(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineChannel)
	if err != nil {
		return fmt.Errorf("listing Keys for Directline Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroup, id.BotServiceName, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := channelsResp.Properties; props != nil {
		if channel, ok := props.AsDirectLineChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("site", flattenDirectlineSites(filterSites(channelProps.Sites)))
			}
		}
	}

	return nil
}

func resourceBotChannelDirectlineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineChannel{
			Properties: &botservice.DirectLineChannelProperties{
				Sites: expandDirectlineSites(d.Get("site").(*schema.Set).List()),
			},
			ChannelName: botservice.ChannelNameDirectLineChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceBotChannelDirectlineRead(d, meta)
}

func resourceBotChannelDirectlineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameDirectLineChannel1))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Directline Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandDirectlineSites(input []interface{}) *[]botservice.DirectLineSite {
	sites := make([]botservice.DirectLineSite, len(input))

	for _, element := range input {
		if element == nil {
			continue
		}

		site := element.(map[string]interface{})
		expanded := botservice.DirectLineSite{}

		if v, ok := site["name"].(string); ok {
			expanded.SiteName = &v
		}
		if v, ok := site["enabled"].(bool); ok {
			expanded.IsEnabled = &v
		}
		if v, ok := site["v1_allowed"].(bool); ok {
			expanded.IsV1Enabled = &v
		}
		if v, ok := site["v3_allowed"].(bool); ok {
			expanded.IsV3Enabled = &v
		}
		if v, ok := site["enhanced_authentication_enabled"].(bool); ok {
			expanded.IsSecureSiteEnabled = &v
		}
		if v, ok := site["trusted_origins"].(*schema.Set); ok {
			origins := v.List()
			items := make([]string, len(origins))
			for i, raw := range origins {
				items[i] = raw.(string)
			}
			expanded.TrustedOrigins = &items
		}

		sites = append(sites, expanded)
	}

	return &sites
}

func flattenDirectlineSites(input []botservice.DirectLineSite) []interface{} {
	sites := make([]interface{}, len(input))

	for i, element := range input {
		site := make(map[string]interface{})

		if v := element.SiteName; v != nil {
			site["name"] = *v
		}

		if element.Key != nil {
			site["key"] = *element.Key
		}

		if element.Key2 != nil {
			site["key2"] = *element.Key2
		}

		if element.IsEnabled != nil {
			site["enabled"] = *element.IsEnabled
		}

		if element.IsV1Enabled != nil {
			site["v1_allowed"] = *element.IsV1Enabled
		}

		if element.IsV3Enabled != nil {
			site["v3_allowed"] = *element.IsV3Enabled
		}

		if element.IsSecureSiteEnabled != nil {
			site["enhanced_authentication_enabled"] = *element.IsSecureSiteEnabled
		}

		if element.TrustedOrigins != nil {
			site["trusted_origins"] = *element.TrustedOrigins
		}

		sites[i] = site
	}

	return sites
}

// When creating a new directline channel, a Default Site is created
// There is a race condition where this site is not removed before the create request is completed
func filterSites(sites *[]botservice.DirectLineSite) []botservice.DirectLineSite {
	filtered := make([]botservice.DirectLineSite, 0)
	for _, site := range *sites {
		if *site.SiteName == "Default Site" {
			continue
		}
		filtered = append(filtered, site)
	}
	return filtered
}
