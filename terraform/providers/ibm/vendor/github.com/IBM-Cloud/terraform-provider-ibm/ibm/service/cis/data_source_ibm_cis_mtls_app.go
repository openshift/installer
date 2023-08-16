// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCISMtlsApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIBMCISMtlsAppRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_mtls_apps",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			"mtls_access_apps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for Access App Response.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application ID",
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application name",
						},
						"app_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Domain",
						},
						"app_aud": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Aud",
						},
						"allowed_idps": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of allowed idps.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"auto_redirect_to_identity": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Auto Redirect to Identity",
						},
						"session_duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Session Duration",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Type",
						},
						"app_uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application UID",
						},
						"app_created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Created At",
						},
						"app_updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Updated At",
						},
					},
				},
			},
			"mtls_access_app_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Access Policies Information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy name",
						},
						"policy_decision": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Decision",
						},
						"policy_precedence": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy Precedence",
						},
						"policy_uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy UID",
						},
						"policy_created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Created At",
						},
						"policy_updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Updated At",
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISMtlsAppValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISMTLSAppValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_mtls_apps",
		Schema:       validateSchema}
	return &iBMCISMTLSAppValidator
}

func dataIBMCISMtlsAppRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()

	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}

	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewListAccessApplicationsOptions(zoneID)
	result, resp, err := sess.ListAccessApplications(opt)
	if err != nil {
		log.Printf("[WARN] List all Applications failed: %v\n", resp)
		return diag.FromErr(err)
	}

	mtlsAppLists := make([]map[string]interface{}, 0)
	mtlsPolicyLists := make([]map[string]interface{}, 0)
	for _, appObj := range result.Result {
		mtlsAppList := map[string]interface{}{}
		mtlsAppList["app_id"] = *appObj.ID
		mtlsAppList["app_name"] = *appObj.Name
		mtlsAppList["app_domain"] = *appObj.Domain
		mtlsAppList["allowed_idps"] = appObj.AllowedIdps
		mtlsAppList["auto_redirect_to_identity"] = appObj.AutoRedirectToIdentity
		mtlsAppList["session_duration"] = appObj.SessionDuration
		mtlsAppList["app_type"] = appObj.Type
		mtlsAppList["app_uid"] = appObj.Uid
		mtlsAppList["app_created_at"] = appObj.CreatedAt
		mtlsAppList["app_updated_at"] = appObj.UpdatedAt

		PolicyOpt := sess.NewListAccessPoliciesOptions(zoneID, *appObj.ID)
		PolicyResult, PolicyResp, PolicyErr := sess.ListAccessPolicies(PolicyOpt)
		if PolicyErr != nil {
			log.Printf("[WARN] List all Policies failed: %v\n", PolicyResp)
			return diag.FromErr(PolicyErr)
		}

		for _, PolicyObj := range PolicyResult.Result {
			mtlsPolicyList := map[string]interface{}{}
			mtlsPolicyList["policy_id"] = *PolicyObj.ID
			mtlsPolicyList["policy_name"] = *PolicyObj.Name
			mtlsPolicyList["policy_decision"] = *PolicyObj.Decision
			mtlsPolicyList["policy_precedence"] = *PolicyObj.Precedence
			mtlsPolicyList["policy_uid"] = *PolicyObj.Uid
			mtlsPolicyList["policy_created_at"] = *PolicyObj.CreatedAt
			mtlsPolicyList["policy_updated_at"] = *PolicyObj.UpdatedAt

			// TODO Include, Exclude and Require of Interface type

			mtlsPolicyLists = append(mtlsPolicyLists, mtlsPolicyList)
		}
		mtlsAppLists = append(mtlsAppLists, mtlsAppList)

	}
	d.SetId(dataSourceCISMtlsAppCheckID(d))

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set("mtls_access_apps", mtlsAppLists)
	d.Set("mtls_access_app_policies", mtlsPolicyLists)

	return nil
}
func dataSourceCISMtlsAppCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
