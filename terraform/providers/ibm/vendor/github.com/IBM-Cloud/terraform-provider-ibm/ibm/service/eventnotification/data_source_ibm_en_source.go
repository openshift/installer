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
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_source", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetSourceOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("source_id").(string))

	result, _, err := enClient.GetSourceWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTemplateWithContext failed: %s", err.Error()), "(Data) ibm_en_source", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_source", "read")
		return tfErr.GetDiag()
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_en_source", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("enabled", result.Enabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enabled flag: %s", err), "(Data) ibm_en_source", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_source", "read")
		return tfErr.GetDiag()
	}

	return nil
}
