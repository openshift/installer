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

func ResourceIbmOnboardingProduct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingProductCreate,
		ReadContext:   resourceIbmOnboardingProductRead,
		UpdateContext: resourceIbmOnboardingProductUpdate,
		DeleteContext: resourceIbmOnboardingProductDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_product", "type"),
				Description:  "The type of the product.",
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
			"eccn_number": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Export Control Classification Number of your product.",
			},
			"ero_class": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ERO class of your product.",
			},
			"unspsc": &schema.Schema{
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The United Nations Standard Products and Services Code of your product.",
			},
			"tax_assessment": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tax assessment type of your product.",
			},
			"support": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The support information that is not displayed in the catalog, but available in ServiceNow.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"escalation_contacts": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of contacts in case of support escalations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the support escalation contact.",
									},
									"email": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The email address of the support escalation contact.",
									},
									"role": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The role of the support escalation contact.",
									},
								},
							},
						},
					},
				},
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IBM Cloud account ID of the provider.",
			},
			"private_catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private catalog that contains the product. Only applicable for software type products.",
			},
			"private_catalog_offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the linked private catalog product. Only applicable for software type products.",
			},
			"global_catalog_offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a global catalog object.",
			},
			"staging_global_catalog_offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of a global catalog object.",
			},
			"approver_resource_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the approval workflow of your product.",
			},
			"iam_registration_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM registration identifier.",
			},
		},
	}
}

func ResourceIbmOnboardingProductValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "professional_service, service, software",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_product", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingProductCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createOnboardingProductOptions := &partnercentersellv1.CreateOnboardingProductOptions{}

	createOnboardingProductOptions.SetType(d.Get("type").(string))
	primaryContactModel, err := ResourceIbmOnboardingProductMapToPrimaryContact(d.Get("primary_contact.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "create", "parse-primary_contact").GetDiag()
	}
	createOnboardingProductOptions.SetPrimaryContact(primaryContactModel)
	if _, ok := d.GetOk("eccn_number"); ok {
		createOnboardingProductOptions.SetEccnNumber(d.Get("eccn_number").(string))
	}
	if _, ok := d.GetOk("ero_class"); ok {
		createOnboardingProductOptions.SetEroClass(d.Get("ero_class").(string))
	}
	if _, ok := d.GetOk("unspsc"); ok {
		createOnboardingProductOptions.SetUnspsc(d.Get("unspsc").(float64))
	}
	if _, ok := d.GetOk("tax_assessment"); ok {
		createOnboardingProductOptions.SetTaxAssessment(d.Get("tax_assessment").(string))
	}
	if _, ok := d.GetOk("support"); ok {
		supportModel, err := ResourceIbmOnboardingProductMapToOnboardingProductSupport(d.Get("support.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "create", "parse-support").GetDiag()
		}
		createOnboardingProductOptions.SetSupport(supportModel)
	}

	onboardingProduct, _, err := partnerCenterSellClient.CreateOnboardingProductWithContext(context, createOnboardingProductOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateOnboardingProductWithContext failed: %s", err.Error()), "ibm_onboarding_product", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*onboardingProduct.ID)

	return resourceIbmOnboardingProductRead(context, d, meta)
}

func resourceIbmOnboardingProductRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOnboardingProductOptions := &partnercentersellv1.GetOnboardingProductOptions{}

	getOnboardingProductOptions.SetProductID(d.Id())

	onboardingProduct, response, err := partnerCenterSellClient.GetOnboardingProductWithContext(context, getOnboardingProductOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOnboardingProductWithContext failed: %s", err.Error()), "ibm_onboarding_product", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("type", onboardingProduct.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-type").GetDiag()
	}
	primaryContactMap, err := ResourceIbmOnboardingProductPrimaryContactToMap(onboardingProduct.PrimaryContact)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "primary_contact-to-map").GetDiag()
	}
	if err = d.Set("primary_contact", []map[string]interface{}{primaryContactMap}); err != nil {
		err = fmt.Errorf("Error setting primary_contact: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-primary_contact").GetDiag()
	}
	if !core.IsNil(onboardingProduct.EccnNumber) {
		if err = d.Set("eccn_number", onboardingProduct.EccnNumber); err != nil {
			err = fmt.Errorf("Error setting eccn_number: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-eccn_number").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.EroClass) {
		if err = d.Set("ero_class", onboardingProduct.EroClass); err != nil {
			err = fmt.Errorf("Error setting ero_class: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-ero_class").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.Unspsc) {
		if err = d.Set("unspsc", onboardingProduct.Unspsc); err != nil {
			err = fmt.Errorf("Error setting unspsc: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-unspsc").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.TaxAssessment) {
		if err = d.Set("tax_assessment", onboardingProduct.TaxAssessment); err != nil {
			err = fmt.Errorf("Error setting tax_assessment: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-tax_assessment").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.Support) {
		supportMap, err := ResourceIbmOnboardingProductOnboardingProductSupportToMap(onboardingProduct.Support)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "support-to-map").GetDiag()
		}
		if err = d.Set("support", []map[string]interface{}{supportMap}); err != nil {
			err = fmt.Errorf("Error setting support: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-support").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.AccountID) {
		if err = d.Set("account_id", onboardingProduct.AccountID); err != nil {
			err = fmt.Errorf("Error setting account_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-account_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.PrivateCatalogID) {
		if err = d.Set("private_catalog_id", onboardingProduct.PrivateCatalogID); err != nil {
			err = fmt.Errorf("Error setting private_catalog_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-private_catalog_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.PrivateCatalogOfferingID) {
		if err = d.Set("private_catalog_offering_id", onboardingProduct.PrivateCatalogOfferingID); err != nil {
			err = fmt.Errorf("Error setting private_catalog_offering_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-private_catalog_offering_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.GlobalCatalogOfferingID) {
		if err = d.Set("global_catalog_offering_id", onboardingProduct.GlobalCatalogOfferingID); err != nil {
			err = fmt.Errorf("Error setting global_catalog_offering_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-global_catalog_offering_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.StagingGlobalCatalogOfferingID) {
		if err = d.Set("staging_global_catalog_offering_id", onboardingProduct.StagingGlobalCatalogOfferingID); err != nil {
			err = fmt.Errorf("Error setting staging_global_catalog_offering_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-staging_global_catalog_offering_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.ApproverResourceID) {
		if err = d.Set("approver_resource_id", onboardingProduct.ApproverResourceID); err != nil {
			err = fmt.Errorf("Error setting approver_resource_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-approver_resource_id").GetDiag()
		}
	}
	if !core.IsNil(onboardingProduct.IamRegistrationID) {
		if err = d.Set("iam_registration_id", onboardingProduct.IamRegistrationID); err != nil {
			err = fmt.Errorf("Error setting iam_registration_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "read", "set-iam_registration_id").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingProductUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateOnboardingProductOptions := &partnercentersellv1.UpdateOnboardingProductOptions{}

	updateOnboardingProductOptions.SetProductID(d.Id())

	hasChange := false

	patchVals := &partnercentersellv1.OnboardingProductPatch{}
	if d.HasChange("primary_contact") {
		primaryContact, err := ResourceIbmOnboardingProductMapToPrimaryContact(d.Get("primary_contact.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "update", "parse-primary_contact").GetDiag()
		}
		patchVals.PrimaryContact = primaryContact
		hasChange = true
	}
	if d.HasChange("eccn_number") {
		newEccnNumber := d.Get("eccn_number").(string)
		patchVals.EccnNumber = &newEccnNumber
		hasChange = true
	}
	if d.HasChange("ero_class") {
		newEroClass := d.Get("ero_class").(string)
		patchVals.EroClass = &newEroClass
		hasChange = true
	}
	if d.HasChange("unspsc") {
		newUnspsc := d.Get("unspsc").(float64)
		patchVals.Unspsc = &newUnspsc
		hasChange = true
	}
	if d.HasChange("tax_assessment") {
		newTaxAssessment := d.Get("tax_assessment").(string)
		patchVals.TaxAssessment = &newTaxAssessment
		hasChange = true
	}
	if d.HasChange("support") {
		support, err := ResourceIbmOnboardingProductMapToOnboardingProductSupport(d.Get("support.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "update", "parse-support").GetDiag()
		}
		patchVals.Support = support
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateOnboardingProductOptions.OnboardingProductPatch = ResourceIbmOnboardingProductOnboardingProductPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateOnboardingProductWithContext(context, updateOnboardingProductOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateOnboardingProductWithContext failed: %s", err.Error()), "ibm_onboarding_product", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingProductRead(context, d, meta)
}

func resourceIbmOnboardingProductDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_product", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteOnboardingProductOptions := &partnercentersellv1.DeleteOnboardingProductOptions{}

	deleteOnboardingProductOptions.SetProductID(d.Id())

	_, err = partnerCenterSellClient.DeleteOnboardingProductWithContext(context, deleteOnboardingProductOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteOnboardingProductWithContext failed: %s", err.Error()), "ibm_onboarding_product", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingProductMapToPrimaryContact(modelMap map[string]interface{}) (*partnercentersellv1.PrimaryContact, error) {
	model := &partnercentersellv1.PrimaryContact{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Email = core.StringPtr(modelMap["email"].(string))
	return model, nil
}

func ResourceIbmOnboardingProductMapToOnboardingProductSupport(modelMap map[string]interface{}) (*partnercentersellv1.OnboardingProductSupport, error) {
	model := &partnercentersellv1.OnboardingProductSupport{}
	if modelMap["escalation_contacts"] != nil {
		escalationContacts := []partnercentersellv1.OnboardingProductSupportEscalationContactItems{}
		for _, escalationContactsItem := range modelMap["escalation_contacts"].([]interface{}) {
			escalationContactsItemModel, err := ResourceIbmOnboardingProductMapToOnboardingProductSupportEscalationContactItems(escalationContactsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			escalationContacts = append(escalationContacts, *escalationContactsItemModel)
		}
		model.EscalationContacts = escalationContacts
	}
	return model, nil
}

func ResourceIbmOnboardingProductMapToOnboardingProductSupportEscalationContactItems(modelMap map[string]interface{}) (*partnercentersellv1.OnboardingProductSupportEscalationContactItems, error) {
	model := &partnercentersellv1.OnboardingProductSupportEscalationContactItems{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["email"] != nil && modelMap["email"].(string) != "" {
		model.Email = core.StringPtr(modelMap["email"].(string))
	}
	if modelMap["role"] != nil && modelMap["role"].(string) != "" {
		model.Role = core.StringPtr(modelMap["role"].(string))
	}
	return model, nil
}

func ResourceIbmOnboardingProductPrimaryContactToMap(model *partnercentersellv1.PrimaryContact) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["email"] = *model.Email
	return modelMap, nil
}

func ResourceIbmOnboardingProductOnboardingProductSupportToMap(model *partnercentersellv1.OnboardingProductSupport) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EscalationContacts != nil {
		escalationContacts := []map[string]interface{}{}
		for _, escalationContactsItem := range model.EscalationContacts {
			escalationContactsItemMap, err := ResourceIbmOnboardingProductOnboardingProductSupportEscalationContactItemsToMap(&escalationContactsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			escalationContacts = append(escalationContacts, escalationContactsItemMap)
		}
		modelMap["escalation_contacts"] = escalationContacts
	}
	return modelMap, nil
}

func ResourceIbmOnboardingProductOnboardingProductSupportEscalationContactItemsToMap(model *partnercentersellv1.OnboardingProductSupportEscalationContactItems) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	if model.Role != nil {
		modelMap["role"] = *model.Role
	}
	return modelMap, nil
}

func ResourceIbmOnboardingProductOnboardingProductPatchAsPatch(patchVals *partnercentersellv1.OnboardingProductPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "primary_contact"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["primary_contact"] = nil
	} else if !exists {
		delete(patch, "primary_contact")
	}
	path = "eccn_number"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["eccn_number"] = nil
	} else if !exists {
		delete(patch, "eccn_number")
	}
	path = "ero_class"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["ero_class"] = nil
	} else if !exists {
		delete(patch, "ero_class")
	}
	path = "unspsc"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["unspsc"] = nil
	} else if !exists {
		delete(patch, "unspsc")
	}
	path = "tax_assessment"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tax_assessment"] = nil
	} else if !exists {
		delete(patch, "tax_assessment")
	}
	path = "support"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["support"] = nil
	} else if exists && patch["support"] != nil {
		ResourceIbmOnboardingProductOnboardingProductSupportAsPatch(patch["support"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "support")
	}

	return patch
}

func ResourceIbmOnboardingProductOnboardingProductSupportAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".escalation_contacts"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["escalation_contacts"] = nil
	} else if exists && patch["escalation_contacts"] != nil {
		escalation_contactsList := patch["escalation_contacts"].([]map[string]interface{})
		for i, escalation_contactsItem := range escalation_contactsList {
			ResourceIbmOnboardingProductOnboardingProductSupportEscalationContactItemsAsPatch(escalation_contactsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "escalation_contacts")
	}
}

func ResourceIbmOnboardingProductOnboardingProductSupportEscalationContactItemsAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = rootPath + ".email"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["email"] = nil
	} else if !exists {
		delete(patch, "email")
	}
	path = rootPath + ".role"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["role"] = nil
	} else if !exists {
		delete(patch, "role")
	}
}
