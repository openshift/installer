// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMEnSMTPUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSMTPUserCreate,
		ReadContext:   resourceIBMEnSMTPUserRead,
		UpdateContext: resourceIBMEnSMTPUserUpdate,
		DeleteContext: resourceIBMEnSMTPUserDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_user", "instance_id"),
				Description:  "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_user", "description"),
				Description:  "SMTP User description.",
			},
			"smtp_config_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SMTP confg Id.",
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
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "SMTP user password.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updated time.",
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id.",
			},
		},
	}
}

func ResourceIBMEnSMTPUserValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]`,
			MinValueLength:             10,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `[a-zA-Z 0-9-_\/.?:'";,+=!#@$%^&*() ]*`,
			MinValueLength:             0,
			MaxValueLength:             250,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_en_smtp_user", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMEnSMTPUserCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createSMTPUserOptions := &en.CreateSMTPUserOptions{}

	createSMTPUserOptions.SetInstanceID(d.Get("instance_id").(string))
	createSMTPUserOptions.SetID(d.Get("smtp_config_id").(string))
	if _, ok := d.GetOk("description"); ok {
		createSMTPUserOptions.SetDescription(d.Get("description").(string))
	}

	smtpUserResponse, _, err := eventNotificationsClient.CreateSMTPUserWithContext(context, createSMTPUserOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("password", smtpUserResponse.Password); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting password: %s", err))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createSMTPUserOptions.InstanceID, *createSMTPUserOptions.ID, *smtpUserResponse.ID))

	return resourceIBMEnSMTPUserRead(context, d, meta)
}

func resourceIBMEnSMTPUserRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSMTPUserOptions := &en.GetSMTPUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getSMTPUserOptions.SetInstanceID(parts[0])
	getSMTPUserOptions.SetID(parts[1])
	getSMTPUserOptions.SetUserID(parts[2])

	smtpUser, response, err := eventNotificationsClient.GetSMTPUserWithContext(context, getSMTPUserOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if !core.IsNil(smtpUser.Description) {
		if err = d.Set("description", smtpUser.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if err = d.Set("smtp_config_id", smtpUser.SMTPConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting smtp_config_id: %s", err))
	}
	if err = d.Set("domain", smtpUser.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting domain: %s", err))
	}
	if err = d.Set("username", smtpUser.Username); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting username: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(smtpUser.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(smtpUser.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("user_id", smtpUser.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting user_id: %s", err))
	}

	return nil
}

func resourceIBMEnSMTPUserUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSMTPUserOptions := &en.UpdateSMTPUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateSMTPUserOptions.SetInstanceID(parts[0])
	updateSMTPUserOptions.SetID(parts[1])
	updateSMTPUserOptions.SetUserID(parts[2])

	hasChange := false

	if d.HasChange("instance_id") {
		return diag.FromErr(err)
	}
	if d.HasChange("description") {
		updateSMTPUserOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	if hasChange {
		_, _, err = eventNotificationsClient.UpdateSMTPUserWithContext(context, updateSMTPUserOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMEnSMTPUserRead(context, d, meta)
}

func resourceIBMEnSMTPUserDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSMTPUserOptions := &en.DeleteSMTPUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSMTPUserOptions.SetInstanceID(parts[0])
	deleteSMTPUserOptions.SetID(parts[1])
	deleteSMTPUserOptions.SetUserID(parts[2])

	_, err = eventNotificationsClient.DeleteSMTPUserWithContext(context, deleteSMTPUserOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
