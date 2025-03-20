// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

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
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccScope() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSccScopeRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_scope", "instance_id"),
				Description:  "The ID of the Security and Compliance Center instance.",
			},
			"scope_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_scope", "scope_id"),
				Description:  "The ID of the scope.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope description.",
			},
			"environment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scope environment. This value details what cloud provider the scope targets.",
			},
			// Manual Change: change name and value for scope_type and scope_id
			"properties": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A list of scopes/targets to exclude from a scope.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclusions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of scopes/targets to exclude from a scope.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A target to exclude from the ibm_scc_scope.",
						},
						"scope_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the target.",
						},
					},
				},
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the account associated with the scope.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the account or service ID who created the scope.",
			},
			"created_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the scope was created.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user or service ID who updated the scope.",
			},
			"updated_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the scope was updated.",
			},
			"attachment_count": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The number of attachments tied to the scope.",
			},
		},
	}
}

func DataSourceIbmSccScopeValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "scope_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_scope", Schema: validateSchema}
	return &resourceValidator
}

func dataSourceIbmSccScopeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getScopeOptions := &securityandcompliancecenterapiv3.GetScopeOptions{}

	getScopeOptions.SetInstanceID(d.Get("instance_id").(string))
	getScopeOptions.SetScopeID(d.Get("scope_id").(string))

	scope, response, err := securityAndComplianceCenterClient.GetScopeWithContext(context, getScopeOptions)
	if err != nil {
		log.Printf("[DEBUG] GetScopeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetScopeWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *scope.ID))

	if !core.IsNil(scope.Name) {
		if err = d.Set("name", scope.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(scope.Description) {
		if err = d.Set("description", scope.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(scope.Environment) {
		if err = d.Set("environment", scope.Environment); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting environment: %s", err))
		}
	}
	if !core.IsNil(scope.Properties) {
		// Manual Change
		if properties, err := resourceIBMSccScopeScopePropertyToMap(scope.Properties, d); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting properties or exclusions: %s\n%#v", err, properties))
		}
		// End Manual Change
	}
	if err = d.Set("account_id", scope.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("created_by", scope.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("created_on", flex.DateTimeToString(scope.CreatedOn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_on: %s", err))
	}
	if err = d.Set("updated_by", scope.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}
	if err = d.Set("updated_on", flex.DateTimeToString(scope.UpdatedOn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_on: %s", err))
	}
	if err = d.Set("attachment_count", scope.AttachmentCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting attachment_count: %s", err))
	}

	return nil
}
