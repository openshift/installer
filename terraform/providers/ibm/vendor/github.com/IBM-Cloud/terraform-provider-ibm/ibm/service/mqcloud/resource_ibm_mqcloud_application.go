// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package mqcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudApplicationCreate,
		ReadContext:   resourceIbmMqcloudApplicationRead,
		DeleteContext: resourceIbmMqcloudApplicationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_application", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQ on Cloud service instance.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_application", "name"),
				Description:  "The name of the application - conforming to MQ rules.",
			},
			"create_api_key_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI to create a new apikey for the application.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this application.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the application which was allocated on creation, and can be used for delete calls.",
			},
		},
	}
}

func ResourceIbmMqcloudApplicationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z][-a-z0-9]*$`,
			MinValueLength:             1,
			MaxValueLength:             12,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_application", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudApplicationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Create Application failed %s", err))
	}

	createApplicationOptions := &mqcloudv1.CreateApplicationOptions{}

	createApplicationOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createApplicationOptions.SetName(d.Get("name").(string))

	applicationCreated, response, err := mqcloudClient.CreateApplicationWithContext(context, createApplicationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateApplicationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateApplicationWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createApplicationOptions.ServiceInstanceGuid, *applicationCreated.ID))

	return resourceIbmMqcloudApplicationRead(context, d, meta)
}

func resourceIbmMqcloudApplicationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getApplicationOptions := &mqcloudv1.GetApplicationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getApplicationOptions.SetServiceInstanceGuid(parts[0])
	getApplicationOptions.SetApplicationID(parts[1])

	applicationDetails, response, err := mqcloudClient.GetApplicationWithContext(context, getApplicationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetApplicationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetApplicationWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", applicationDetails.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service instance guid: %s", err))
	}
	if err = d.Set("create_api_key_uri", applicationDetails.CreateApiKeyURI); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting create_api_key_uri: %s", err))
	}
	if err = d.Set("href", applicationDetails.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("application_id", applicationDetails.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting application_id: %s", err))
	}

	return nil
}

func resourceIbmMqcloudApplicationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Delete Application failed %s", err))
	}

	deleteApplicationOptions := &mqcloudv1.DeleteApplicationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteApplicationOptions.SetServiceInstanceGuid(parts[0])
	deleteApplicationOptions.SetApplicationID(parts[1])

	response, err := mqcloudClient.DeleteApplicationWithContext(context, deleteApplicationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteApplicationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteApplicationWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
