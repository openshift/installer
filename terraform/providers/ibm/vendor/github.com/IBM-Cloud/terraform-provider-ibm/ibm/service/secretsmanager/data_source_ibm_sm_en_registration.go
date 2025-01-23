// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmEnRegistration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmEnRegistrationRead,

		Schema: map[string]*schema.Schema{
			"event_notifications_instance_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
		},
	}
}

func dataSourceIbmSmEnRegistrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", EnRegistrationResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getNotificationsRegistrationOptions := &secretsmanagerv2.GetNotificationsRegistrationOptions{}

	notificationsRegistration, response, err := secretsManagerClient.GetNotificationsRegistrationWithContext(context, getNotificationsRegistrationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetNotificationsRegistrationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetNotificationsRegistrationWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", EnRegistrationResourceName), "read")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", region, instanceId))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", EnRegistrationResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("event_notifications_instance_crn", notificationsRegistration.EventNotificationsInstanceCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting event_notifications_instance_crn"), fmt.Sprintf("(Data) %s", EnRegistrationResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
