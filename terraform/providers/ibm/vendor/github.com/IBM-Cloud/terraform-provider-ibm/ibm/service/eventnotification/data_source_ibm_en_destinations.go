// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnDestinations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnDestinationsRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the destinations by name or type.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of destinations.",
			},
			"destinations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of destinations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination description.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination type Email/SMS/Webhook.",
						},
						"subscription_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subscription count.",
						},
						"subscription_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Names of subscriptions.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated at.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMEnDestinationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.ListDestinationsOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var destinationList *en.DestinationList

	finalList := []en.DestinationListItem{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, response, err := enClient.ListDestinationsWithContext(context, options)

		destinationList = result

		if err != nil {
			return diag.FromErr(fmt.Errorf("ListDestinationsWithContext failed %s\n%s", err, response))
		}

		offset = offset + limit

		finalList = append(finalList, result.Destinations...)

		if offset > *result.TotalCount {
			break
		}
	}

	destinationList.Destinations = finalList

	d.SetId(fmt.Sprintf("destinations/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(destinationList.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
	}

	if destinationList.Destinations != nil {
		if err = d.Set("destinations", enFlattenDestinationsList(destinationList.Destinations)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting destinations %s", err))
		}
	}

	return nil
}

func enFlattenDestinationsList(result []en.DestinationListItem) (destinations []map[string]interface{}) {
	for _, destinationsItem := range result {
		destinations = append(destinations, enDestinationListToMap(destinationsItem))
	}

	return destinations
}

func enDestinationListToMap(destinationItem en.DestinationListItem) (destination map[string]interface{}) {
	destination = map[string]interface{}{}

	if destinationItem.ID != nil {
		destination["id"] = destinationItem.ID
	}
	if destinationItem.Name != nil {
		destination["name"] = destinationItem.Name
	}
	if destinationItem.Description != nil {
		destination["description"] = destinationItem.Description
	}
	if destinationItem.Type != nil {
		destination["type"] = destinationItem.Type
	}
	if destinationItem.SubscriptionCount != nil {
		destination["subscription_count"] = destinationItem.SubscriptionCount
	}
	if destinationItem.SubscriptionNames != nil {
		destination["subscription_names"] = destinationItem.SubscriptionNames
	}
	if destinationItem.UpdatedAt != nil {
		destination["updated_at"] = flex.DateTimeToString(destinationItem.UpdatedAt)
	}

	return destination
}
