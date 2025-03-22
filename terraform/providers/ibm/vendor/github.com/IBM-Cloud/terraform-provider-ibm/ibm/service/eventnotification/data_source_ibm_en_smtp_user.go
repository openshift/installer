// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnSMTPUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMTPUserRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"en_smtp_config_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for SMTP.",
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UserID.",
			},
			"smtp_config_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP confg Id.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP User description.",
			},
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain Name.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP user name.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updated time.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updated time.",
			},
		},
	}
}

func dataSourceIBMEnSMTPUserRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_smtp_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSMTPUserOptions := &en.GetSMTPUserOptions{}

	getSMTPUserOptions.SetInstanceID(d.Get("instance_id").(string))
	getSMTPUserOptions.SetID(d.Get("en_smtp_config_id").(string))
	getSMTPUserOptions.SetUserID(d.Get("user_id").(string))

	smtpUser, _, err := eventNotificationsClient.GetSMTPUserWithContext(context, getSMTPUserOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTemplateWithContext failed: %s", err.Error()), "(Data) ibm_en_smtp_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *getSMTPUserOptions.InstanceID, *getSMTPUserOptions.ID, *getSMTPUserOptions.UserID))

	if err = d.Set("smtp_config_id", smtpUser.SMTPConfigID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("description", smtpUser.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("domain", smtpUser.Domain); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("username", smtpUser.Username); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(smtpUser.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(smtpUser.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_smtp_user", "read")
		return tfErr.GetDiag()
	}

	return nil
}
