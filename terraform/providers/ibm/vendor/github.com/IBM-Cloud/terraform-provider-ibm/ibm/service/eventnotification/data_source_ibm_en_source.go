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

func DataSourceIBMEnSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSourceRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for Source.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source description.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The Source enable flag.",
			},
		},
	}
}

func dataSourceIBMEnSourceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetSourceOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("source_id").(string))

	result, response, err := enClient.GetSourceWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetSource failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}

	if err = d.Set("enabled", result.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enabled flag: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}
