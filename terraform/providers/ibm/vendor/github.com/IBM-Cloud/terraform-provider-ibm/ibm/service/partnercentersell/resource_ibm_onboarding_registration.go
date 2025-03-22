// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.99.1-daeb6e46-20250131-173156
 */

package partnercentersell

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/partnercentersellv1"
)

func ResourceIbmOnboardingRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingRegistrationCreate,
		ReadContext:   resourceIbmOnboardingRegistrationRead,
		UpdateContext: resourceIbmOnboardingRegistrationUpdate,
		DeleteContext: resourceIbmOnboardingRegistrationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_registration", "account_id"),
				Description:  "The ID of your account.",
			},
			"company_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_registration", "company_name"),
				Description:  "The name of your company that is displayed in the IBM Cloud catalog.",
			},
			"primary_contact": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The primary contact for your product.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the primary contact for your product.",
						},
						"email": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The email address of the primary contact for your product.",
						},
					},
				},
			},
			"default_private_catalog_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_registration", "default_private_catalog_id"),
				Description:  "The default private catalog in which products are created.",
			},
			"provider_access_group": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_registration", "provider_access_group"),
				Description:  "The onboarding access group for your team.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the registration was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the registration was updated.",
			},
		},
	}
}

func ResourceIbmOnboardingRegistrationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "account_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9]+$`,
			MinValueLength:             1,
			MaxValueLength:             32,
		},
		validate.ValidateSchema{
			Identifier:                 "company_name",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "default_private_catalog_id",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
		},
		validate.ValidateSchema{
			Identifier:                 "provider_access_group",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^AccessGroupId-[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             1,
			MaxValueLength:             50,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_registration", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingRegistrationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createRegistrationOptions := &partnercentersellv1.CreateRegistrationOptions{}

	createRegistrationOptions.SetAccountID(d.Get("account_id").(string))
	createRegistrationOptions.SetCompanyName(d.Get("company_name").(string))
	primaryContactModel, err := ResourceIbmOnboardingRegistrationMapToPrimaryContact(d.Get("primary_contact.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "create", "parse-primary_contact").GetDiag()
	}
	createRegistrationOptions.SetPrimaryContact(primaryContactModel)
	if _, ok := d.GetOk("default_private_catalog_id"); ok {
		createRegistrationOptions.SetDefaultPrivateCatalogID(d.Get("default_private_catalog_id").(string))
	}
	if _, ok := d.GetOk("provider_access_group"); ok {
		createRegistrationOptions.SetProviderAccessGroup(d.Get("provider_access_group").(string))
	}

	registration, _, err := partnerCenterSellClient.CreateRegistrationWithContext(context, createRegistrationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_registration", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*registration.ID)

	return resourceIbmOnboardingRegistrationRead(context, d, meta)
}

func resourceIbmOnboardingRegistrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRegistrationOptions := &partnercentersellv1.GetRegistrationOptions{}

	getRegistrationOptions.SetRegistrationID(d.Id())

	registration, response, err := partnerCenterSellClient.GetRegistrationWithContext(context, getRegistrationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_registration", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("account_id", registration.AccountID); err != nil {
		err = fmt.Errorf("Error setting account_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-account_id").GetDiag()
	}
	if err = d.Set("company_name", registration.CompanyName); err != nil {
		err = fmt.Errorf("Error setting company_name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-company_name").GetDiag()
	}
	primaryContactMap, err := ResourceIbmOnboardingRegistrationPrimaryContactToMap(registration.PrimaryContact)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "primary_contact-to-map").GetDiag()
	}
	if err = d.Set("primary_contact", []map[string]interface{}{primaryContactMap}); err != nil {
		err = fmt.Errorf("Error setting primary_contact: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-primary_contact").GetDiag()
	}
	if !core.IsNil(registration.DefaultPrivateCatalogID) {
		if err = d.Set("default_private_catalog_id", registration.DefaultPrivateCatalogID); err != nil {
			err = fmt.Errorf("Error setting default_private_catalog_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-default_private_catalog_id").GetDiag()
		}
	}
	if !core.IsNil(registration.ProviderAccessGroup) {
		if err = d.Set("provider_access_group", registration.ProviderAccessGroup); err != nil {
			err = fmt.Errorf("Error setting provider_access_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-provider_access_group").GetDiag()
		}
	}
	if !core.IsNil(registration.CreatedAt) {
		if err = d.Set("created_at", registration.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(registration.UpdatedAt) {
		if err = d.Set("updated_at", registration.UpdatedAt); err != nil {
			err = fmt.Errorf("Error setting updated_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "read", "set-updated_at").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingRegistrationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateRegistrationOptions := &partnercentersellv1.UpdateRegistrationOptions{}

	updateRegistrationOptions.SetRegistrationID(d.Id())

	hasChange := false

	patchVals := &partnercentersellv1.RegistrationPatch{}
	if d.HasChange("company_name") {
		newCompanyName := d.Get("company_name").(string)
		patchVals.CompanyName = &newCompanyName
		hasChange = true
	}
	if d.HasChange("primary_contact") {
		primaryContact, err := ResourceIbmOnboardingRegistrationMapToPrimaryContact(d.Get("primary_contact.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "update", "parse-primary_contact").GetDiag()
		}
		patchVals.PrimaryContact = primaryContact
		hasChange = true
	}
	if d.HasChange("default_private_catalog_id") {
		newDefaultPrivateCatalogID := d.Get("default_private_catalog_id").(string)
		patchVals.DefaultPrivateCatalogID = &newDefaultPrivateCatalogID
		hasChange = true
	}
	if d.HasChange("provider_access_group") {
		newProviderAccessGroup := d.Get("provider_access_group").(string)
		patchVals.ProviderAccessGroup = &newProviderAccessGroup
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateRegistrationOptions.RegistrationPatch = ResourceIbmOnboardingRegistrationRegistrationPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateRegistrationWithContext(context, updateRegistrationOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_registration", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingRegistrationRead(context, d, meta)
}

func resourceIbmOnboardingRegistrationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_registration", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteRegistrationOptions := &partnercentersellv1.DeleteRegistrationOptions{}

	deleteRegistrationOptions.SetRegistrationID(d.Id())

	_, err = partnerCenterSellClient.DeleteRegistrationWithContext(context, deleteRegistrationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteRegistrationWithContext failed: %s", err.Error()), "ibm_onboarding_registration", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingRegistrationMapToPrimaryContact(modelMap map[string]interface{}) (*partnercentersellv1.PrimaryContact, error) {
	model := &partnercentersellv1.PrimaryContact{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Email = core.StringPtr(modelMap["email"].(string))
	return model, nil
}

func ResourceIbmOnboardingRegistrationPrimaryContactToMap(model *partnercentersellv1.PrimaryContact) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["email"] = *model.Email
	return modelMap, nil
}

func ResourceIbmOnboardingRegistrationRegistrationPatchAsPatch(patchVals *partnercentersellv1.RegistrationPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "company_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["company_name"] = nil
	} else if !exists {
		delete(patch, "company_name")
	}
	path = "primary_contact"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["primary_contact"] = nil
	} else if !exists {
		delete(patch, "primary_contact")
	}
	path = "default_private_catalog_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["default_private_catalog_id"] = nil
	} else if !exists {
		delete(patch, "default_private_catalog_id")
	}
	path = "provider_access_group"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["provider_access_group"] = nil
	} else if !exists {
		delete(patch, "provider_access_group")
	}

	return patch
}
