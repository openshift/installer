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

func DataSourceIBMEnIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnIntegrationRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"integration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for Integration.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of key integration kms/hs-crypto.",
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public or private endpoint for kms/hpcs",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the kms/hpcs instance",
						},
						"root_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of root key id",
						},
					},
				},
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func dataSourceIBMEnIntegrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetIntegrationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("integration_id").(string))

	result, response, err := enClient.GetIntegrationWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetIntegration failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if result.Metadata != nil {
		d.Set("metadata", enKMSIntegrationFlattenMetadata(result.Metadata))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func enKMSIntegrationFlattenMetadata(metadataItem *en.IntegrationMetadata) []interface{} {

	metadataMap := make(map[string]interface{})
	if metadataItem.Endpoint != nil {
		metadataMap["endpoint"] = metadataItem.Endpoint
	}

	if metadataItem.CRN != nil {
		metadataMap["crn"] = metadataItem.CRN
	}

	if metadataItem.RootKeyID != nil {
		metadataMap["root_key_id"] = metadataItem.RootKeyID
	}
	return []interface{}{metadataMap}
}
