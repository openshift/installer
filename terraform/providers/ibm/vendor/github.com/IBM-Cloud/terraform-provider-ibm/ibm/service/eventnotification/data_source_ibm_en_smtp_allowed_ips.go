// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnSMTPAllowedIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMTPAllowedIpsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"en_smtp_allowed_ips_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for SMTP.",
			},
			"subnets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The SMTP allowed Ips.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updated at.",
			},
		},
	}
}

func dataSourceIBMEnSMTPAllowedIpsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_smtp_allowed_ips", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSMTPAllowedIpsOptions := &eventnotificationsv1.GetSMTPAllowedIpsOptions{}

	getSMTPAllowedIpsOptions.SetInstanceID(d.Get("instance_id").(string))
	getSMTPAllowedIpsOptions.SetID(d.Get("en_smtp_allowed_ips_id").(string))

	smtpAllowedIPs, _, err := eventNotificationsClient.GetSMTPAllowedIpsWithContext(context, getSMTPAllowedIpsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSMTPAllowedIpsWithContext failed: %s", err.Error()), "(Data) ibm_en_smtp_allowed_ips", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMEnSMTPAllowedIpsID(d))

	if err = d.Set("updated_at", flex.DateTimeToString(smtpAllowedIPs.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_smtp_allowed_ips", "read")
		return tfErr.GetDiag()
	}

	if smtpAllowedIPs.Subnets != nil {
		err = d.Set("subnets", smtpAllowedIPs.Subnets)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnets %s", err))
		}
	}

	return nil
}

// dataSourceIBMEnSMTPAllowedIpsID returns a reasonable ID for the list.
func dataSourceIBMEnSMTPAllowedIpsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
