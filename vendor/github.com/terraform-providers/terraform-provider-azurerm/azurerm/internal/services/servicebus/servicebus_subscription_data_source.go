package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceServiceBusSubscription() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceBusSubscriptionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NamespaceName,
			},

			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TopicName(),
			},

			"auto_delete_on_idle": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_message_ttl": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"lock_duration": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dead_lettering_on_message_expiration": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dead_lettering_on_filter_evaluation_error": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_batched_operations": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"max_delivery_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"requires_session": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"forward_to": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"forward_dead_lettered_messages_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceBusSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.SubscriptionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSubscriptionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("topic_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.TopicName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := existing.SBSubscriptionProperties; props != nil {
		d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
		d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
		d.Set("lock_duration", props.LockDuration)
		d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
		d.Set("dead_lettering_on_filter_evaluation_error", props.DeadLetteringOnFilterEvaluationExceptions)
		d.Set("enable_batched_operations", props.EnableBatchedOperations)
		d.Set("requires_session", props.RequiresSession)
		d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
		d.Set("forward_to", props.ForwardTo)

		maxDeliveryCount := 0
		if props.MaxDeliveryCount != nil {
			maxDeliveryCount = int(*props.MaxDeliveryCount)
		}

		d.Set("max_delivery_count", maxDeliveryCount)
	}

	return nil
}
