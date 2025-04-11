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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMEnSMTPConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSMTPConfigurationCreate,
		ReadContext:   resourceIBMEnSMTPConfigurationRead,
		UpdateContext: resourceIBMEnSMTPConfigurationUpdate,
		DeleteContext: resourceIBMEnSMTPConfigurationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_configuration", "instance_id"),
				Description:  "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_configuration", "name"),
				Description:  "SMTP name.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_configuration", "description"),
				Description:  "SMTP description.",
			},
			"domain": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_en_smtp_configuration", "domain"),
				Description:  "Domain Name.",
			},
			"verification_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SPF/DKIM.",
			},
			"config": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Payload describing a SMTP configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dkim": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DKIM attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"txt_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DMIM text name.",
									},
									"txt_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DMIM text value.",
									},
									"verification": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "dkim verification.",
									},
								},
							},
						},
						"en_authorization": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The en_authorization attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"verification": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "en_authorization verification.",
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
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created time.",
			},
			"en_smtp_configuration_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP ID.",
			},
		},
	}
}

func ResourceIBMEnSMTPConfigurationValidator() *validate.ResourceValidator {
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
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[a-zA-Z 0-9-_\/.?:'";,+=!#@$%^&*() ]*`,
			MinValueLength:             1,
			MaxValueLength:             250,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `[a-zA-Z 0-9-_\/.?:'";,+=!#@$%^&*() ]*`,
			MinValueLength:             1,
			MaxValueLength:             250,
		},
		validate.ValidateSchema{
			Identifier:                 "domain",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `.*`,
			MinValueLength:             1,
			MaxValueLength:             512,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_en_smtp_configuration", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMEnSMTPConfigurationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createSMTPConfigurationOptions := &en.CreateSMTPConfigurationOptions{}

	createSMTPConfigurationOptions.SetInstanceID(d.Get("instance_id").(string))
	createSMTPConfigurationOptions.SetName(d.Get("name").(string))
	createSMTPConfigurationOptions.SetDomain(d.Get("domain").(string))
	if _, ok := d.GetOk("description"); ok {
		createSMTPConfigurationOptions.SetDescription(d.Get("description").(string))
	}

	smtpCreateResponse, _, err := eventNotificationsClient.CreateSMTPConfigurationWithContext(context, createSMTPConfigurationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSMTPConfigurationWithContext failed: %s", err.Error()), "ibm_en_smtp_configuration", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createSMTPConfigurationOptions.InstanceID, *smtpCreateResponse.ID))

	return resourceIBMEnSMTPConfigurationRead(context, d, meta)
}

func resourceIBMEnSMTPConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSMTPConfigurationOptions := &en.GetSMTPConfigurationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "update")
		return tfErr.GetDiag()
	}

	getSMTPConfigurationOptions.SetInstanceID(parts[0])
	getSMTPConfigurationOptions.SetID(parts[1])

	smtpConfiguration, response, err := eventNotificationsClient.GetSMTPConfigurationWithContext(context, getSMTPConfigurationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSMTPConfigurationWithContext failed: %s", err.Error()), "ibm_en_smtp_configuration", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", smtpConfiguration.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(smtpConfiguration.Description) {
		if err = d.Set("description", smtpConfiguration.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if err = d.Set("domain", smtpConfiguration.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting domain: %s", err))
	}
	configMap, err := resourceIBMEnSMTPConfigurationSMTPConfigToMap(smtpConfiguration.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("config", []map[string]interface{}{configMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(smtpConfiguration.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("en_smtp_configuration_id", smtpConfiguration.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting en_smtp_configuration_id: %s", err))
	}

	return nil
}

func resourceIBMEnSMTPConfigurationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateSMTPConfigurationOptions := &en.UpdateSMTPConfigurationOptions{}
	verifySMTPConfiguration := &en.UpdateVerifySMTPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "update")
		return tfErr.GetDiag()
	}

	updateSMTPConfigurationOptions.SetInstanceID(parts[0])
	updateSMTPConfigurationOptions.SetID(parts[1])
	verifySMTPConfiguration.SetType(d.Get("verification_type").(string))
	verifySMTPConfiguration.SetInstanceID(parts[0])
	verifySMTPConfiguration.SetID(parts[1])

	hasChange := false
	hasChangeverification := false

	if d.HasChange("instance_id") {
		return diag.FromErr(err)
	}
	if d.HasChange("name") {
		updateSMTPConfigurationOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateSMTPConfigurationOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("verification_type") {
		verifySMTPConfiguration.SetType(d.Get("verification_type").(string))
		hasChangeverification = true
	}

	if hasChange {
		_, _, err = eventNotificationsClient.UpdateSMTPConfigurationWithContext(context, updateSMTPConfigurationOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSMTPConfigurationWithContext failed: %s", err.Error()), "ibm_en_smtp_configuration", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if hasChangeverification {
		_, _, err = eventNotificationsClient.UpdateVerifySMTPWithContext(context, verifySMTPConfiguration)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMEnSMTPConfigurationRead(context, d, meta)
}

func resourceIBMEnSMTPConfigurationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	eventNotificationsClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteSMTPConfigurationOptions := &en.DeleteSMTPConfigurationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_smtp_configuration", "delete")
		return tfErr.GetDiag()
	}

	deleteSMTPConfigurationOptions.SetInstanceID(parts[0])
	deleteSMTPConfigurationOptions.SetID(parts[1])

	_, err = eventNotificationsClient.DeleteSMTPConfigurationWithContext(context, deleteSMTPConfigurationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSMTPConfigurationWithContext: failed: %s", err.Error()), "ibm_en_smtp_configuration", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIBMEnSMTPConfigurationSMTPConfigToMap(model *en.SMTPConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Dkim != nil {
		dkimMap, err := resourceIBMEnSMTPConfigurationSmtpdkimAttributesToMap(model.Dkim)
		if err != nil {
			return modelMap, err
		}
		modelMap["dkim"] = []map[string]interface{}{dkimMap}
	}
	if model.EnAuthorization != nil {
		enAuthorizationMap, err := resourceIBMEnSMTPConfigurationEnAuthAttributesToMap(model.EnAuthorization)
		if err != nil {
			return modelMap, err
		}
		modelMap["en_authorization"] = []map[string]interface{}{enAuthorizationMap}
	}
	if model.Spf != nil {
		spfMap, err := resourceIBMEnSMTPConfigurationSpfAttributesToMap(model.Spf)
		if err != nil {
			return modelMap, err
		}
		modelMap["spf"] = []map[string]interface{}{spfMap}
	}
	return modelMap, nil
}

func resourceIBMEnSMTPConfigurationSmtpdkimAttributesToMap(model *en.SmtpdkimAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TxtName != nil {
		modelMap["txt_name"] = model.TxtName
	}
	if model.TxtValue != nil {
		modelMap["txt_value"] = model.TxtValue
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}

func resourceIBMEnSMTPConfigurationEnAuthAttributesToMap(model *en.EnAuthAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}

func resourceIBMEnSMTPConfigurationSpfAttributesToMap(model *en.SpfAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TxtName != nil {
		modelMap["txt_name"] = model.TxtName
	}
	if model.TxtValue != nil {
		modelMap["txt_value"] = model.TxtValue
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}
