// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func ResourceIBMEnFCMSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnFCMSubscriptionCreate,
		ReadContext:   resourceIBMEnFCMSubscriptionRead,
		UpdateContext: resourceIBMEnFCMSubscriptionUpdate,
		DeleteContext: resourceIBMEnFCMSubscriptionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subscription name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscription description.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination ID.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Topic ID.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription ID.",
			},
			"destination_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Destination.",
			},
			"destination_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Destintion name.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the topic.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func resourceIBMEnFCMSubscriptionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.CreateSubscriptionOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	options.SetName(d.Get("name").(string))
	options.SetTopicID(d.Get("topic_id").(string))
	options.SetDestinationID(d.Get("destination_id").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	result, _, err := enClient.CreateSubscriptionWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("CreateSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnFCMSubscriptionRead(context, d, meta)
}

func resourceIBMEnFCMSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.GetSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "read")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("GetSubscriptionWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("subscription_id", result.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}

	if result.From != nil {
		if err = d.Set("from", result.From); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting from: %s", err))
		}
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("destination_type", result.DestinationType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_type: %s", err))
	}

	if result.DestinationName != nil {
		if err = d.Set("destination_name", result.DestinationName); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_name: %s", err))
		}
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_id: %s", err))
	}

	if result.TopicName != nil {
		if err = d.Set("topic_name", result.TopicName); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_name: %s", err))
		}
	}

	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnFCMSubscriptionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.UpdateSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "update")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	if ok := d.HasChanges("name", "description", "attributes"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		_, _, err := enClient.UpdateSubscriptionWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
			// return diag.FromErr(fmt.Errorf("UpdateSubscriptionWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnFCMSubscriptionRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnFCMSubscriptionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.DeleteSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription", "delete")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSubscriptionWithContext: failed: %s", err.Error()), "ibm_en_subscription", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("DeleteSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
