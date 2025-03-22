// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccInstanceSettings() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccInstanceSettingsRead,

		Schema: map[string]*schema.Schema{
			"event_notifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Event Notifications settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Event Notifications instance CRN.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the Event Notifications connection was updated.",
						},
						"source_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Security and Compliance Center instance CRN.",
						},
					},
				},
			},
			"object_storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Object Storage settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage instance CRN.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage bucket name.",
						},
						"bucket_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage bucket location.",
						},
						"bucket_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connected Cloud Object Storage bucket endpoint.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the bucket connection was updated.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccInstanceSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &securityandcompliancecenterapiv3.GetSettingsOptions{}
	getSettingsOptions.SetInstanceID(d.Get("instance_id").(string))

	settings, response, err := adminClient.GetSettingsWithContext(context, getSettingsOptions)

	service_url := adminClient.GetServiceURL()
	d.SetId(service_url)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(settings.EventNotifications) {
		eventNotificationsMap, err := dataSourceIbmSccInstanceSettingsEventNotificationsToMap(settings.EventNotifications)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("event_notifications", []map[string]interface{}{eventNotificationsMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting event_notifications: %s", err))
		}
	}
	if !core.IsNil(settings.ObjectStorage) {
		objectStorageMap, err := dataSourceIbmSccInstanceSettingsObjectStorageToMap(settings.ObjectStorage)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("object_storage", []map[string]interface{}{objectStorageMap}); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting object_storage: %s", err))
		}
	}
	return nil
}

func dataSourceIbmSccInstanceSettingsEventNotificationsToMap(model *securityandcompliancecenterapiv3.EventNotifications) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InstanceCRN != nil {
		modelMap["instance_crn"] = model.InstanceCRN
	}
	if model.UpdatedOn != nil {
		modelMap["updated_on"] = model.UpdatedOn.String()
	}
	if model.SourceID != nil {
		modelMap["source_id"] = model.SourceID
	}
	return modelMap, nil
}

func dataSourceIbmSccInstanceSettingsObjectStorageToMap(model *securityandcompliancecenterapiv3.ObjectStorage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InstanceCRN != nil {
		modelMap["instance_crn"] = model.InstanceCRN
	}
	if model.Bucket != nil {
		modelMap["bucket"] = model.Bucket
	}
	if model.BucketLocation != nil {
		modelMap["bucket_location"] = model.BucketLocation
	}
	if model.BucketEndpoint != nil {
		modelMap["bucket_endpoint"] = model.BucketEndpoint
	}
	if model.UpdatedOn != nil {
		modelMap["updated_on"] = model.UpdatedOn.String()
	}
	return modelMap, nil
}
