// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmEnRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmEnRegistrationCreate,
		ReadContext:   resourceIbmSmEnRegistrationRead,
		UpdateContext: resourceIbmSmEnRegistrationUpdate,
		DeleteContext: resourceIbmSmEnRegistrationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"event_notifications_instance_crn": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator(EnRegistrationResourceName, "event_notifications_instance_crn"),
				Description:  "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"event_notifications_source_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator(EnRegistrationResourceName, "event_notifications_source_name"),
				Description:  "The name that is displayed as a source that is in your Event Notifications instance.",
			},
			"event_notifications_source_description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator(EnRegistrationResourceName, "event_notifications_source_description"),
				Description:  "An optional description for the source  that is in your Event Notifications instance.",
			},
		},
	}
}

func ResourceIbmSmEnRegistrationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "event_notifications_instance_crn",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^crn:v[0-9](:([A-Za-z0-9-._~!$&'()*+,;=@\/]|%[0-9A-Z]{2})*){8}$`,
			MinValueLength:             9,
			MaxValueLength:             512,
		},
		validate.ValidateSchema{
			Identifier:                 "event_notifications_source_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             2,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "event_notifications_source_description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             0,
			MaxValueLength:             1024,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: EnRegistrationResourceName, Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSmEnRegistrationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", EnRegistrationResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createNotificationsRegistrationOptions := &secretsmanagerv2.CreateNotificationsRegistrationOptions{}

	createNotificationsRegistrationOptions.SetEventNotificationsInstanceCrn(d.Get("event_notifications_instance_crn").(string))
	createNotificationsRegistrationOptions.SetEventNotificationsSourceName(d.Get("event_notifications_source_name").(string))
	if _, ok := d.GetOk("event_notifications_source_description"); ok {
		createNotificationsRegistrationOptions.SetEventNotificationsSourceDescription(d.Get("event_notifications_source_description").(string))
	}

	_, response, err := secretsManagerClient.CreateNotificationsRegistrationWithContext(context, createNotificationsRegistrationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateNotificationsRegistrationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateNotificationsRegistrationWithContext failed %s\n%s", err, response), EnRegistrationResourceName, "create")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", region, instanceId))

	return resourceIbmSmEnRegistrationRead(context, d, meta)
}

func resourceIbmSmEnRegistrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 2 {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Wrong format of resource ID. To import event notification registration use the format `<region>/<instance_id>`"), EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getNotificationsRegistrationOptions := &secretsmanagerv2.GetNotificationsRegistrationOptions{}

	notificationsRegistration, response, err := secretsManagerClient.GetNotificationsRegistrationWithContext(context, getNotificationsRegistrationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetNotificationsRegistrationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetNotificationsRegistrationWithContext failed %s\n%s", err, response), EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("event_notifications_instance_crn", notificationsRegistration.EventNotificationsInstanceCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting event_notifications_instance_crn"), EnRegistrationResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmEnRegistrationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf(""), EnRegistrationResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createNotificationsRegistrationOptions := &secretsmanagerv2.CreateNotificationsRegistrationOptions{}

	hasChange := false

	if d.HasChange("event_notifications_instance_crn") || d.HasChange("event_notifications_source_name") {
		createNotificationsRegistrationOptions.SetEventNotificationsInstanceCrn(d.Get("event_notifications_instance_crn").(string))
		createNotificationsRegistrationOptions.SetEventNotificationsSourceName(d.Get("event_notifications_source_name").(string))
		hasChange = true
	}
	if d.HasChange("event_notifications_source_description") {
		createNotificationsRegistrationOptions.SetEventNotificationsSourceDescription(d.Get("event_notifications_source_description").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := secretsManagerClient.CreateNotificationsRegistrationWithContext(context, createNotificationsRegistrationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateNotificationsRegistrationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateNotificationsRegistrationWithContext failed %s\n%s", err, response), EnRegistrationResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmEnRegistrationRead(context, d, meta)
}

func resourceIbmSmEnRegistrationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", EnRegistrationResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteNotificationsRegistrationOptions := &secretsmanagerv2.DeleteNotificationsRegistrationOptions{}

	response, err := secretsManagerClient.DeleteNotificationsRegistrationWithContext(context, deleteNotificationsRegistrationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteNotificationsRegistrationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteNotificationsRegistrationWithContext failed %s\n%s", err, response), EnRegistrationResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
