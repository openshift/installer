// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func DataSourceIBMCdToolchainToolSecuritycompliance() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMCdToolchainToolSecuritycomplianceRead,

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the toolchain.",
			},
			"tool_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the tool bound to the toolchain.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where tool can be found.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool CRN.",
			},
			"toolchain_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of toolchain which the tool is bound to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI representing the tool.",
			},
			"referent": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on URIs to access this resource through the UI or API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI representing the this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI representing the this resource through an API.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool name.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest tool update timestamp.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Parameters to be used to create the tool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Give this tool integration a name, for example: my-security-compliance.",
						},
						"evidence_repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "To collect and store evidence for all tasks performed, a Git repository is required as an evidence locker.",
						},
						"trigger_scan": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enabling trigger validation scans provides details for a pipeline task to trigger a scan.",
						},
						"location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"evidence_namespace": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kind of pipeline evidence to be displayed in Security and Compliance Center for this toolchain. The evidence locker will be searched for CD (Continuous Deployment) pipeline evidence, or for CC (Continuous Compliance) pipeline evidence.",
						},
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The IBM Cloud API key is used to access the Security and Compliance API. You can obtain your API key with 'ibmcloud iam api-key-create' or via the console at https://cloud.ibm.com/iam#/apikeys by clicking **Create API key** (Each API key only can be viewed once).",
						},
						"scope": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Select an existing scope name to narrow the focus of the validation scan. [Learn more.](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-scopes) ![](https://cloud.ibm.com/media/docs/images/icons/launch-glyph.svg).",
						},
						"profile": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Select an existing profile, where a profile is a collection of security controls. [Learn more.](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-profiles) ![](https://cloud.ibm.com/media/docs/images/icons/launch-glyph.svg).",
						},
						"trigger_info": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "trigger_info.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current configuration state of the tool.",
			},
		},
	}
}

func DataSourceIBMCdToolchainToolSecuritycomplianceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	getToolByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))
	getToolByIDOptions.SetToolID(d.Get("tool_id").(string))

	getToolByIDResponse, response, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if *getToolByIDResponse.ToolTypeID != "security_compliance" {
		return diag.FromErr(fmt.Errorf("Retrieved tool is not the correct type: %s", *getToolByIDResponse.ToolTypeID))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getToolByIDOptions.ToolchainID, *getToolByIDOptions.ToolID))

	if err = d.Set("resource_group_id", getToolByIDResponse.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}

	if err = d.Set("crn", getToolByIDResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("toolchain_crn", getToolByIDResponse.ToolchainCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}

	if err = d.Set("href", getToolByIDResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	referent := []map[string]interface{}{}
	if getToolByIDResponse.Referent != nil {
		modelMap, err := DataSourceIBMCdToolchainToolSecuritycomplianceToolReferentToMap(getToolByIDResponse.Referent)
		if err != nil {
			return diag.FromErr(err)
		}
		referent = append(referent, modelMap)
	}
	if err = d.Set("referent", referent); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent %s", err))
	}

	if err = d.Set("name", getToolByIDResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(getToolByIDResponse.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	parameters := []map[string]interface{}{}
	if getToolByIDResponse.Parameters != nil {
		remapFields := map[string]string{
			"api_key": "api-key",
		}
		modelMap := GetParametersFromRead(getToolByIDResponse.Parameters, DataSourceIBMCdToolchainToolSecuritycompliance(), remapFields)
		parameters = append(parameters, modelMap)
	}
	if err = d.Set("parameters", parameters); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters %s", err))
	}

	if err = d.Set("state", getToolByIDResponse.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	return nil
}

func DataSourceIBMCdToolchainToolSecuritycomplianceToolReferentToMap(model *cdtoolchainv2.ToolReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}
