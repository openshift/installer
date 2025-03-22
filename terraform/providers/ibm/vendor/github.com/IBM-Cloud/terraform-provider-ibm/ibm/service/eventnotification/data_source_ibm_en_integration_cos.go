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

func DataSourceIBMEnCOSIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnCOSIntegrationRead,

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
				Description: "The type of integration is collect_failed_events.",
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public or private endpoint for COS bucket",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the Cloud Object Storage instance",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud Object Storage bucket name",
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

func dataSourceIBMEnCOSIntegrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_integration_cos", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetIntegrationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("integration_id").(string))

	result, _, err := enClient.GetIntegrationWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetIntegrationWithContext failed: %s", err.Error()), "(Data) ibm_en_integration_cos", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("type", result.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_en_integration_cos", "read")
		return tfErr.GetDiag()
	}

	if result.Metadata != nil {
		d.Set("metadata", enCOSIntegrationFlattenMetadata(result.Metadata))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_integration_cos", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func enCOSIntegrationFlattenMetadata(metadataItem *en.IntegrationMetadata) []interface{} {

	metadataMap := make(map[string]interface{})
	if metadataItem.Endpoint != nil {
		metadataMap["endpoint"] = metadataItem.Endpoint
	}

	if metadataItem.CRN != nil {
		metadataMap["crn"] = metadataItem.CRN
	}

	if metadataItem.BucketName != nil {
		metadataMap["bucket_name"] = metadataItem.BucketName
	}
	return []interface{}{metadataMap}
}
