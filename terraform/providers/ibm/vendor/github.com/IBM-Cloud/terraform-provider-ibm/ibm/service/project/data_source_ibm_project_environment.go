// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.92.1-44330004-20240620-143510
 */

package project

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/project-go-sdk/projectv1"
)

func DataSourceIbmProjectEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmProjectEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique project ID.",
			},
			"project_environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment ID.",
			},
			"project": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project that is referenced by this resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The definition of the project reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the project.",
									},
								},
							},
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An IBM Cloud resource name that uniquely identifies a resource.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"target_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target account ID derived from the authentication block values. The target account exists only if the environment currently has an authorization block.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The environment definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the environment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the environment. It's unique within the account across projects and regions.",
						},
						"authorizations": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trusted profile ID.",
									},
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The authorization method. You can authorize by using a trusted profile or an API key in Secrets Manager.",
									},
									"api_key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
										Description: "The IBM Cloud API Key. It can be either raw or pulled from the catalog via a `CRN` or `JSON` blob.",
									},
								},
							},
						},
						"inputs": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The input variables that are used for configuration definition and environment.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"compliance_profile": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile that is required for compliance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID for the compliance profile.",
									},
									"instance_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for the instance of a compliance profile.",
									},
									"instance_location": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The location of the compliance instance.",
									},
									"attachment_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for the attachment to a compliance profile.",
									},
									"profile_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the compliance profile.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmProjectEnvironmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_environment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectEnvironmentOptions := &projectv1.GetProjectEnvironmentOptions{}

	getProjectEnvironmentOptions.SetProjectID(d.Get("project_id").(string))
	getProjectEnvironmentOptions.SetID(d.Get("project_environment_id").(string))

	environment, _, err := projectClient.GetProjectEnvironmentWithContext(context, getProjectEnvironmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectEnvironmentWithContext failed: %s", err.Error()), "(Data) ibm_project_environment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getProjectEnvironmentOptions.ProjectID, *getProjectEnvironmentOptions.ID))

	project := []map[string]interface{}{}
	if environment.Project != nil {
		modelMap, err := DataSourceIbmProjectEnvironmentProjectReferenceToMap(environment.Project)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_project_environment", "read", "project-to-map").GetDiag()
		}
		project = append(project, modelMap)
	}
	if err = d.Set("project", project); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting project: %s", err), "(Data) ibm_project_environment", "read", "set-project").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(environment.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_project_environment", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("target_account", environment.TargetAccount); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting target_account: %s", err), "(Data) ibm_project_environment", "read", "set-target_account").GetDiag()
	}

	if err = d.Set("modified_at", flex.DateTimeToString(environment.ModifiedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting modified_at: %s", err), "(Data) ibm_project_environment", "read", "set-modified_at").GetDiag()
	}

	if err = d.Set("href", environment.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_project_environment", "read", "set-href").GetDiag()
	}

	definition := []map[string]interface{}{}
	if environment.Definition != nil {
		modelMap, err := DataSourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesResponseToMap(environment.Definition)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_project_environment", "read", "definition-to-map").GetDiag()
		}
		definition = append(definition, modelMap)
	}
	if err = d.Set("definition", definition); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting definition: %s", err), "(Data) ibm_project_environment", "read", "set-definition").GetDiag()
	}

	return nil
}

func DataSourceIbmProjectEnvironmentProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["href"] = *model.Href
	definitionMap, err := DataSourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = *model.Crn
	return modelMap, nil
}

func DataSourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesResponseToMap(model *projectv1.EnvironmentDefinitionRequiredPropertiesResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = *model.Description
	modelMap["name"] = *model.Name
	if model.Authorizations != nil {
		authorizationsMap, err := DataSourceIbmProjectEnvironmentProjectConfigAuthToMap(model.Authorizations)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
	}
	if model.Inputs != nil {
		inputs := make(map[string]interface{})
		for k, v := range model.Inputs {
			inputs[k] = flex.Stringify(v)
		}
		modelMap["inputs"] = inputs
	}
	if model.ComplianceProfile != nil {
		complianceProfileMap, err := DataSourceIbmProjectEnvironmentProjectComplianceProfileToMap(model.ComplianceProfile)
		if err != nil {
			return modelMap, err
		}
		modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
	}
	return modelMap, nil
}

func DataSourceIbmProjectEnvironmentProjectConfigAuthToMap(model *projectv1.ProjectConfigAuth) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = *model.TrustedProfileID
	}
	if model.Method != nil {
		modelMap["method"] = *model.Method
	}
	if model.ApiKey != nil {
		modelMap["api_key"] = *model.ApiKey
	}
	return modelMap, nil
}

func DataSourceIbmProjectEnvironmentProjectComplianceProfileToMap(model *projectv1.ProjectComplianceProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = *model.InstanceID
	}
	if model.InstanceLocation != nil {
		modelMap["instance_location"] = *model.InstanceLocation
	}
	if model.AttachmentID != nil {
		modelMap["attachment_id"] = *model.AttachmentID
	}
	if model.ProfileName != nil {
		modelMap["profile_name"] = *model.ProfileName
	}
	return modelMap, nil
}
