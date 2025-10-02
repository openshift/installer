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

func resourceBotChannelSlack() *schema.Resource {
	return &schema.Resource{
		Create: resourceBotChannelSlackCreate,
		Read:   resourceBotChannelSlackRead,
		Delete: resourceBotChannelSlackDelete,
		Update: resourceBotChannelSlackUpdate,

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

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"verification_token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"landing_page_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceBotChannelSlackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameSlackChannel))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ChannelName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_slack", resourceId.ID())
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameSlackChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("creating Slack Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceBotChannelSlackRead(d, meta)
}

func resourceBotChannelSlackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Slack Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retreving Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsSlackChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("client_id", channelProps.ClientID)
			}
		}
	}

	return nil
}

func resourceBotChannelSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameSlackChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("updating Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceBotChannelSlackRead(d, meta)
}

func resourceBotChannelSlackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Slack Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
