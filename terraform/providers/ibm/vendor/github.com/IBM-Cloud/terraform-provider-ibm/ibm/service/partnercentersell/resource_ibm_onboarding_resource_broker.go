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

func ResourceIbmOnboardingResourceBroker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmOnboardingResourceBrokerCreate,
		ReadContext:   resourceIbmOnboardingResourceBrokerRead,
		UpdateContext: resourceIbmOnboardingResourceBrokerUpdate,
		DeleteContext: resourceIbmOnboardingResourceBrokerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"env": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "env"),
				Description:  "The environment to fetch this object from.",
			},
			"auth_username": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "auth_username"),
				Description:  "The authentication username to reach the broker.",
			},
			"auth_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authentication password to reach the broker.",
			},
			"auth_scheme": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "auth_scheme"),
				Description:  "The supported authentication scheme for the broker.",
			},
			"resource_group_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cloud resource name of the resource group.",
			},
			"state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "state"),
				Description:  "The state of the broker.",
			},
			"broker_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL associated with the broker application.",
			},
			"allow_context_updates": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the resource controller will call the broker for any context changes to the instance. Currently, the only context related change is an instance name update.",
			},
			"catalog_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To enable the provisioning of your broker, set this parameter value to `service`.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "type"),
				Description:  "The type of the provisioning model.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_onboarding_resource_broker", "name"),
				Description:  "The name of the broker.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the pricing plan is available.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the account in which you manage the broker.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud resource name (CRN) of the broker.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the service broker was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the service broker was updated.",
			},
			"deleted_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the service broker was deleted.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of the user who created this broker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who dispatched this action.",
						},
						"user_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username of the user who dispatched this action.",
						},
					},
				},
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of the user who updated this broker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who dispatched this action.",
						},
						"user_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username of the user who dispatched this action.",
						},
					},
				},
			},
			"deleted_by": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of the user who deleted this broker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who dispatched this action.",
						},
						"user_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username of the user who dispatched this action.",
						},
					},
				},
			},
			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique identifier of the broker.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL associated with the broker.",
			},
		},
	}
}

func ResourceIbmOnboardingResourceBrokerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "env",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "auth_username",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "apikey",
		},
		validate.ValidateSchema{
			Identifier:                 "auth_scheme",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "bearer, bearer-crn",
		},
		validate.ValidateSchema{
			Identifier:                 "state",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "active, removed",
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "provision_behind, provision_through",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[ -~\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_onboarding_resource_broker", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmOnboardingResourceBrokerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createResourceBrokerOptions := &partnercentersellv1.CreateResourceBrokerOptions{}

	createResourceBrokerOptions.SetAuthScheme(d.Get("auth_scheme").(string))
	createResourceBrokerOptions.SetName(d.Get("name").(string))
	createResourceBrokerOptions.SetBrokerURL(d.Get("broker_url").(string))
	createResourceBrokerOptions.SetType(d.Get("type").(string))
	if _, ok := d.GetOk("auth_username"); ok {
		createResourceBrokerOptions.SetAuthUsername(d.Get("auth_username").(string))
	}
	if _, ok := d.GetOk("auth_password"); ok {
		createResourceBrokerOptions.SetAuthPassword(d.Get("auth_password").(string))
	}
	if _, ok := d.GetOk("resource_group_crn"); ok {
		createResourceBrokerOptions.SetResourceGroupCrn(d.Get("resource_group_crn").(string))
	}
	if _, ok := d.GetOk("state"); ok {
		createResourceBrokerOptions.SetState(d.Get("state").(string))
	}
	if _, ok := d.GetOk("allow_context_updates"); ok {
		createResourceBrokerOptions.SetAllowContextUpdates(d.Get("allow_context_updates").(bool))
	}
	if _, ok := d.GetOk("catalog_type"); ok {
		createResourceBrokerOptions.SetCatalogType(d.Get("catalog_type").(string))
	}
	if _, ok := d.GetOk("region"); ok {
		createResourceBrokerOptions.SetRegion(d.Get("region").(string))
	}
	if _, ok := d.GetOk("env"); ok {
		createResourceBrokerOptions.SetEnv(d.Get("env").(string))
	}

	broker, _, err := partnerCenterSellClient.CreateResourceBrokerWithContext(context, createResourceBrokerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateResourceBrokerWithContext failed: %s", err.Error()), "ibm_onboarding_resource_broker", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*broker.ID)

	return resourceIbmOnboardingResourceBrokerRead(context, d, meta)
}

func resourceIbmOnboardingResourceBrokerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getResourceBrokerOptions := &partnercentersellv1.GetResourceBrokerOptions{}

	getResourceBrokerOptions.SetBrokerID(d.Id())
	if _, ok := d.GetOk("env"); ok {
		getResourceBrokerOptions.SetEnv(d.Get("env").(string))
	}

	broker, response, err := partnerCenterSellClient.GetResourceBrokerWithContext(context, getResourceBrokerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetResourceBrokerWithContext failed: %s", err.Error()), "ibm_onboarding_resource_broker", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(broker.AuthUsername) {
		if err = d.Set("auth_username", broker.AuthUsername); err != nil {
			err = fmt.Errorf("Error setting auth_username: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-auth_username").GetDiag()
		}
	}
	if !core.IsNil(broker.AuthPassword) {
		if err = d.Set("auth_password", broker.AuthPassword); err != nil {
			err = fmt.Errorf("Error setting auth_password: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-auth_password").GetDiag()
		}
	}
	if err = d.Set("auth_scheme", broker.AuthScheme); err != nil {
		err = fmt.Errorf("Error setting auth_scheme: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-auth_scheme").GetDiag()
	}
	if !core.IsNil(broker.ResourceGroupCrn) {
		if err = d.Set("resource_group_crn", broker.ResourceGroupCrn); err != nil {
			err = fmt.Errorf("Error setting resource_group_crn: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-resource_group_crn").GetDiag()
		}
	}
	if !core.IsNil(broker.State) {
		if err = d.Set("state", broker.State); err != nil {
			err = fmt.Errorf("Error setting state: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-state").GetDiag()
		}
	}
	if err = d.Set("broker_url", broker.BrokerURL); err != nil {
		err = fmt.Errorf("Error setting broker_url: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-broker_url").GetDiag()
	}
	if !core.IsNil(broker.AllowContextUpdates) {
		if err = d.Set("allow_context_updates", broker.AllowContextUpdates); err != nil {
			err = fmt.Errorf("Error setting allow_context_updates: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-allow_context_updates").GetDiag()
		}
	}
	if !core.IsNil(broker.CatalogType) {
		if err = d.Set("catalog_type", broker.CatalogType); err != nil {
			err = fmt.Errorf("Error setting catalog_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-catalog_type").GetDiag()
		}
	}
	if err = d.Set("type", broker.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-type").GetDiag()
	}
	if err = d.Set("name", broker.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-name").GetDiag()
	}
	if !core.IsNil(broker.Region) {
		if err = d.Set("region", broker.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(broker.AccountID) {
		if err = d.Set("account_id", broker.AccountID); err != nil {
			err = fmt.Errorf("Error setting account_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-account_id").GetDiag()
		}
	}
	if !core.IsNil(broker.Crn) {
		if err = d.Set("crn", broker.Crn); err != nil {
			err = fmt.Errorf("Error setting crn: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-crn").GetDiag()
		}
	}
	if !core.IsNil(broker.CreatedAt) {
		if err = d.Set("created_at", flex.DateTimeToString(broker.CreatedAt)); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(broker.UpdatedAt) {
		if err = d.Set("updated_at", flex.DateTimeToString(broker.UpdatedAt)); err != nil {
			err = fmt.Errorf("Error setting updated_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-updated_at").GetDiag()
		}
	}
	if !core.IsNil(broker.DeletedAt) {
		if err = d.Set("deleted_at", flex.DateTimeToString(broker.DeletedAt)); err != nil {
			err = fmt.Errorf("Error setting deleted_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-deleted_at").GetDiag()
		}
	}
	if !core.IsNil(broker.CreatedBy) {
		createdByMap, err := ResourceIbmOnboardingResourceBrokerBrokerEventCreatedByUserToMap(broker.CreatedBy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "created_by-to-map").GetDiag()
		}
		if err = d.Set("created_by", []map[string]interface{}{createdByMap}); err != nil {
			err = fmt.Errorf("Error setting created_by: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-created_by").GetDiag()
		}
	}
	if !core.IsNil(broker.UpdatedBy) {
		updatedByMap, err := ResourceIbmOnboardingResourceBrokerBrokerEventUpdatedByUserToMap(broker.UpdatedBy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "updated_by-to-map").GetDiag()
		}
		if err = d.Set("updated_by", []map[string]interface{}{updatedByMap}); err != nil {
			err = fmt.Errorf("Error setting updated_by: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-updated_by").GetDiag()
		}
	}
	if !core.IsNil(broker.DeletedBy) {
		deletedByMap, err := ResourceIbmOnboardingResourceBrokerBrokerEventDeletedByUserToMap(broker.DeletedBy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "deleted_by-to-map").GetDiag()
		}
		if err = d.Set("deleted_by", []map[string]interface{}{deletedByMap}); err != nil {
			err = fmt.Errorf("Error setting deleted_by: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-deleted_by").GetDiag()
		}
	}
	if !core.IsNil(broker.Guid) {
		if err = d.Set("guid", broker.Guid); err != nil {
			err = fmt.Errorf("Error setting guid: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-guid").GetDiag()
		}
	}
	if !core.IsNil(broker.URL) {
		if err = d.Set("url", broker.URL); err != nil {
			err = fmt.Errorf("Error setting url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "read", "set-url").GetDiag()
		}
	}

	return nil
}

func resourceIbmOnboardingResourceBrokerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateResourceBrokerOptions := &partnercentersellv1.UpdateResourceBrokerOptions{}

	updateResourceBrokerOptions.SetBrokerID(d.Id())
	if _, ok := d.GetOk("env"); ok {
		updateResourceBrokerOptions.SetEnv(d.Get("env").(string))
	}

	hasChange := false

	patchVals := &partnercentersellv1.BrokerPatch{}
	if d.HasChange("auth_username") {
		newAuthUsername := d.Get("auth_username").(string)
		patchVals.AuthUsername = &newAuthUsername
		hasChange = true
	}
	if d.HasChange("auth_password") {
		newAuthPassword := d.Get("auth_password").(string)
		patchVals.AuthPassword = &newAuthPassword
		hasChange = true
	}
	if d.HasChange("auth_scheme") {
		newAuthScheme := d.Get("auth_scheme").(string)
		patchVals.AuthScheme = &newAuthScheme
		hasChange = true
	}
	if d.HasChange("resource_group_crn") {
		newResourceGroupCrn := d.Get("resource_group_crn").(string)
		patchVals.ResourceGroupCrn = &newResourceGroupCrn
		hasChange = true
	}
	if d.HasChange("state") {
		newState := d.Get("state").(string)
		patchVals.State = &newState
		hasChange = true
	}
	if d.HasChange("broker_url") {
		newBrokerURL := d.Get("broker_url").(string)
		patchVals.BrokerURL = &newBrokerURL
		hasChange = true
	}
	if d.HasChange("allow_context_updates") {
		newAllowContextUpdates := d.Get("allow_context_updates").(bool)
		patchVals.AllowContextUpdates = &newAllowContextUpdates
		hasChange = true
	}
	if d.HasChange("catalog_type") {
		newCatalogType := d.Get("catalog_type").(string)
		patchVals.CatalogType = &newCatalogType
		hasChange = true
	}
	if d.HasChange("type") {
		newType := d.Get("type").(string)
		patchVals.Type = &newType
		hasChange = true
	}
	if d.HasChange("region") {
		newRegion := d.Get("region").(string)
		patchVals.Region = &newRegion
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateResourceBrokerOptions.BrokerPatch = ResourceIbmOnboardingResourceBrokerBrokerPatchAsPatch(patchVals, d)

		_, _, err = partnerCenterSellClient.UpdateResourceBrokerWithContext(context, updateResourceBrokerOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateResourceBrokerWithContext failed: %s", err.Error()), "ibm_onboarding_resource_broker", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmOnboardingResourceBrokerRead(context, d, meta)
}

func resourceIbmOnboardingResourceBrokerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	partnerCenterSellClient, err := meta.(conns.ClientSession).PartnerCenterSellV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_onboarding_resource_broker", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteResourceBrokerOptions := &partnercentersellv1.DeleteResourceBrokerOptions{}

	deleteResourceBrokerOptions.SetBrokerID(d.Id())
	if _, ok := d.GetOk("env"); ok {
		deleteResourceBrokerOptions.SetEnv(d.Get("env").(string))
	}

	_, err = partnerCenterSellClient.DeleteResourceBrokerWithContext(context, deleteResourceBrokerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteResourceBrokerWithContext failed: %s", err.Error()), "ibm_onboarding_resource_broker", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmOnboardingResourceBrokerBrokerEventCreatedByUserToMap(model *partnercentersellv1.BrokerEventCreatedByUser) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UserID != nil {
		modelMap["user_id"] = *model.UserID
	}
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func ResourceIbmOnboardingResourceBrokerBrokerEventUpdatedByUserToMap(model *partnercentersellv1.BrokerEventUpdatedByUser) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UserID != nil {
		modelMap["user_id"] = *model.UserID
	}
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func ResourceIbmOnboardingResourceBrokerBrokerEventDeletedByUserToMap(model *partnercentersellv1.BrokerEventDeletedByUser) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UserID != nil {
		modelMap["user_id"] = *model.UserID
	}
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func ResourceIbmOnboardingResourceBrokerBrokerPatchAsPatch(patchVals *partnercentersellv1.BrokerPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "auth_username"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auth_username"] = nil
	} else if !exists {
		delete(patch, "auth_username")
	}
	path = "auth_password"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auth_password"] = nil
	} else if !exists {
		delete(patch, "auth_password")
	}
	path = "auth_scheme"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auth_scheme"] = nil
	} else if !exists {
		delete(patch, "auth_scheme")
	}
	path = "resource_group_crn"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["resource_group_crn"] = nil
	} else if !exists {
		delete(patch, "resource_group_crn")
	}
	path = "state"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["state"] = nil
	} else if !exists {
		delete(patch, "state")
	}
	path = "broker_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["broker_url"] = nil
	} else if !exists {
		delete(patch, "broker_url")
	}
	path = "allow_context_updates"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["allow_context_updates"] = nil
	} else if !exists {
		delete(patch, "allow_context_updates")
	}
	path = "catalog_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["catalog_type"] = nil
	} else if !exists {
		delete(patch, "catalog_type")
	}
	path = "type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	} else if !exists {
		delete(patch, "type")
	}
	path = "region"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["region"] = nil
	} else if !exists {
		delete(patch, "region")
	}

	return patch
}
