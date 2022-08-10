// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
)

func DataSourceIBMSccNotificationSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSccAccountNotificationSettingsRead,

		Schema: map[string]*schema.Schema{
			"instance_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the Event Notifications instance that you want to connect.",
			},
		},
	}
}

func dataSourceIbmSccAccountNotificationSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceApiClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &adminserviceapiv1.GetSettingsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions.SetAccountID(userDetails.UserAccount)

	accountSettings, response, err := adminServiceApiClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}
	notificationsSettings := accountSettings.EventNotifications

	if err = d.Set("instance_crn", notificationsSettings.InstanceCrn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_crn: %s", err))
	}

	d.SetId("scc_admin_notification_settings")

	return nil
}

// dataSourceIbmSccAccountNotificationSettingsID returns a reasonable ID for the list.
func dataSourceIbmSccAccountNotificationSettingsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
