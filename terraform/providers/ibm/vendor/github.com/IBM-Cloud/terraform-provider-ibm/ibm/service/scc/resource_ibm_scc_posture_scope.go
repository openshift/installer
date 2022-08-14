// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func ResourceIBMSccPostureScopes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccPostureScopesCreate,
		ReadContext:   resourceIBMSccPostureScopesRead,
		UpdateContext: resourceIBMSccPostureScopesUpdate,
		DeleteContext: resourceIBMSccPostureScopesDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_scope", "name"),
				Description:  "A unique name for your scope.",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_scope", "description"),
				Description:  "A detailed description of the scope.",
			},
			"collector_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The unique IDs of the collectors that are attached to the scope.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"credential_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_scope", "credential_id"),
				Description:  "The unique identifier of the credential.",
			},
			"credential_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_scope", "credential_type"),
				Description:  "The environment that the scope is targeted to.",
			},
		},
	}
}

func ResourceIBMSccPostureScopesValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\.,_\s]*$`,
			MinValueLength:             3,
			MaxValueLength:             50,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\.,_\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             255,
		},
		validate.ValidateSchema{
			Identifier:                 "credential_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\.,_\\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             50,
		},
		validate.ValidateSchema{
			Identifier:                 "credential_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "aws, azure, gcp, hosted, ibm, on_premise, openstack, services",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_posture_scope", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSccPostureScopesCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createScopeOptions := &posturemanagementv2.CreateScopeOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	createScopeOptions.SetAccountID(userDetails.UserAccount)

	createScopeOptions.SetName(d.Get("name").(string))
	createScopeOptions.SetDescription(d.Get("description").(string))
	collector_ids_int := d.Get("collector_ids").([]interface{})
	collector_ids := make([]string, len(collector_ids_int))
	for i, collector_id := range collector_ids_int {
		collector_ids[i] = collector_id.(string)
	}
	createScopeOptions.SetCollectorIds(collector_ids) //[]string{
	createScopeOptions.SetCredentialID(d.Get("credential_id").(string))
	createScopeOptions.SetCredentialType(d.Get("credential_type").(string))

	scope, response, err := postureManagementClient.CreateScopeWithContext(context, createScopeOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateScopeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateScopeWithContext failed %s\n%s", err, response))
	}

	d.SetId(*scope.ID)

	return resourceIBMSccPostureScopesRead(context, d, meta)
}

func resourceIBMSccPostureScopesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getScopesOptions := &posturemanagementv2.GetScopeDetailsOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	getScopesOptions.SetAccountID(accountID)
	getScopesOptions.SetID(d.Id())

	scope, response, err := postureManagementClient.GetScopeDetailsWithContext(context, getScopesOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetScopeDetailsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetScopeDetailsWithContext failed %s\n%s", err, response))
	}
	d.SetId(*scope.ID)

	return nil
}

func resourceIBMSccPostureScopesUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateScopeDetailsOptions := &posturemanagementv2.UpdateScopeDetailsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	updateScopeDetailsOptions.SetAccountID(userDetails.UserAccount)

	hasChange := false

	updateScopeDetailsOptions.SetID(d.Id())

	if d.HasChange("name") {
		updateScopeDetailsOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateScopeDetailsOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := postureManagementClient.UpdateScopeDetailsWithContext(context, updateScopeDetailsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateScopeDetailsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateScopeDetailsWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccPostureScopesRead(context, d, meta)
}

func resourceIBMSccPostureScopesDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteScopeOptions := &posturemanagementv2.DeleteScopeOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	deleteScopeOptions.SetAccountID(userDetails.UserAccount)

	deleteScopeOptions.SetID(d.Id())

	response, err := postureManagementClient.DeleteScopeWithContext(context, deleteScopeOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteScopeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteScopeWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
