// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMEnCustomEmailDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnCustomEmailDestinationCreate,
		ReadContext:   resourceIBMEnCustomEmailDestinationRead,
		UpdateContext: resourceIBMEnCustomEmailDestinationUpdate,
		DeleteContext: resourceIBMEnCustomEmailDestinationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Destintion name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of Destination type smtp_custom.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Destination description.",
			},
			"collect_failed_events": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to collect the failed event in Cloud Object Storage bucket",
			},
			"verification_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_destination_custom_email", "verification_type"),
				Description:  "Verification Method spf/dkim.",
			},
			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Payload describing a destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Domain for the Custom Domain Email Destination",
									},
									"dkim": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The DKIM attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"public_key": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim public key.",
												},
												"selector": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim selector.",
												},
												"verification": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim verification.",
												},
											},
										},
									},
									"spf": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The SPF attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"txt_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf text name.",
												},
												"txt_value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf text value.",
												},
												"verification": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf verification.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"destination_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination ID",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
			"subscription_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of subscriptions.",
			},
			"subscription_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subscriptions.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ResourceIBMEnEmailDestinationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "verification_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "spf,dkim",
			MinValueLength:             1,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_en_destination_custom_email", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMEnCustomEmailDestinationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.CreateDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetName(d.Get("name").(string))

	options.SetType(d.Get("type").(string))
	options.SetCollectFailedEvents(d.Get("collect_failed_events").(bool))
	destinationtype := d.Get("type").(string)
	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("config"); ok {
		config := CustomEmaildestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}), destinationtype)
		options.SetConfig(&config)
	}

	result, _, err := enClient.CreateDestinationWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDestinationWithContext failed: %s", err.Error()), "ibm_en_destination_custom_email", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnServiceNowDestinationRead(context, d, meta)
}

func resourceIBMEnCustomEmailDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetDestinationWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDestinationWithContext failed: %s", err.Error()), "ibm_en_destination_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("destination_id", options.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if err = d.Set("collect_failed_events", result.CollectFailedEvents); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting CollectFailedEvents: %s", err))
	}

	if err = d.Set("description", result.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}

	if result.Config != nil {
		err = d.Set("config", enCustomEmailDestinationFlattenConfig(*result.Config))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting config %s", err))
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_count: %s", err))
	}

	if err = d.Set("subscription_names", result.SubscriptionNames); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_names: %s", err))
	}

	return nil
}

func resourceIBMEnCustomEmailDestinationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.UpdateDestinationOptions{}
	verifyCustomEmailDestinationConfiguration := &en.UpdateVerifyDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "update")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	verifyCustomEmailDestinationConfiguration.SetInstanceID(parts[0])
	verifyCustomEmailDestinationConfiguration.SetID(parts[1])
	hasChangeverification := false

	verifyCustomEmailDestinationConfiguration.SetType(d.Get("verification_type").(string))

	if ok := d.HasChanges("name", "description", "collect_failed_events", "config"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		if _, ok := d.GetOk("collect_failed_events"); ok {
			options.SetCollectFailedEvents(d.Get("collect_failed_events").(bool))
		}

		destinationtype := d.Get("type").(string)

		if d.HasChange("verification_type") {
			verifyCustomEmailDestinationConfiguration.SetType(d.Get("verification_type").(string))
			hasChangeverification = true
		}
		if _, ok := d.GetOk("config"); ok {
			config := CustomEmaildestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}), destinationtype)
			options.SetConfig(&config)
		}
		_, _, err := enClient.UpdateDestinationWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateDestinationWithContext failed: %s", err.Error()), "ibm_en_destination_custom_email", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		if hasChangeverification {
			_, _, err = enClient.UpdateVerifyDestinationWithContext(context, verifyCustomEmailDestinationConfiguration)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		return resourceIBMEnCustomEmailDestinationRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnCustomEmailDestinationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.DeleteDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_destination_custom_email", "delete")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteDestinationWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDestinationWithContext failed: %s", err.Error()), "ibm_en_destination_custom_email", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func CustomEmaildestinationConfigMapToDestinationConfig(configParams map[string]interface{}, destinationtype string) en.DestinationConfig {
	params := new(en.DestinationConfigOneOfCustomDomainEmailDestinationConfig)
	if configParams["domain"] != nil {
		params.Domain = core.StringPtr(configParams["domain"].(string))
	}

	if configParams["dkim"] != nil && len(configParams["dkim"].([]interface{})) > 0 {
		DkimModel, _ := resourceIBMEnDestinationMapToDkimAttributes(configParams["dkim"].([]interface{})[0].(map[string]interface{}))
		params.Dkim = &DkimModel
	}
	if configParams["spf"] != nil && len(configParams["spf"].([]interface{})) > 0 {
		SpfModel, _ := resourceIBMEnDestinationMapToSpfAttributes(configParams["spf"].([]interface{})[0].(map[string]interface{}))
		params.Spf = &SpfModel
	}

	destinationConfig := new(en.DestinationConfig)
	destinationConfig.Params = params
	return *destinationConfig
}

func resourceIBMEnDestinationMapToDkimAttributes(modelMap map[string]interface{}) (en.DkimAttributes, error) {
	model := new(en.DkimAttributes)
	if modelMap["public_key"] != nil && modelMap["public_key"].(string) != "" {
		model.PublicKey = core.StringPtr(modelMap["public_key"].(string))
	}
	if modelMap["selector"] != nil && modelMap["selector"].(string) != "" {
		model.Selector = core.StringPtr(modelMap["selector"].(string))
	}
	if modelMap["verification"] != nil && modelMap["verification"].(string) != "" {
		model.Verification = core.StringPtr(modelMap["verification"].(string))
	}
	return *model, nil
}

func resourceIBMEnDestinationMapToSpfAttributes(modelMap map[string]interface{}) (en.SpfAttributes, error) {
	model := new(en.SpfAttributes)
	if modelMap["txt_name"] != nil && modelMap["txt_name"].(string) != "" {
		model.TxtName = core.StringPtr(modelMap["txt_name"].(string))
	}
	if modelMap["txt_value"] != nil && modelMap["txt_value"].(string) != "" {
		model.TxtValue = core.StringPtr(modelMap["txt_value"].(string))
	}
	if modelMap["verification"] != nil && modelMap["verification"].(string) != "" {
		model.Verification = core.StringPtr(modelMap["verification"].(string))
	}
	return *model, nil
}
