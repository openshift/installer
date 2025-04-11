// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmCodeEngineSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineSecretRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your secret.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Data container that allows to specify config parameters and their values as a key-value map. Each key field must consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value field can consists of any character and must not exceed a max length of 1048576 characters.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the secret instance, which is used to achieve optimistic locking.",
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specify the format of the secret.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new secret,  a URL is created identifying the location of the instance.",
			},
			"secret_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the secret.",
			},
			"service_access": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Properties for Service Access Secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_key": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The service credential associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the service credential associated with the secret.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the service credential associated with the secret.",
									},
								},
							},
						},
						"role": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A reference to the Role and Role CRN for service binding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN of the IAM Role for this service access secret.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role of the service credential.",
									},
								},
							},
						},
						"service_instance": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IBM Cloud service instance associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the IBM Cloud service instance associated with the secret.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of IBM Cloud service associated with the secret.",
									},
								},
							},
						},
						"serviceid": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A reference to a Service ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN value of a Service ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the Service ID.",
									},
								},
							},
						},
					},
				},
			},
			"service_operator": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Properties for the IBM Cloud Operator Secret.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apikey_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the apikey associated with the operator secret.",
						},
						"resource_group_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of resource groups (by ID) that the operator secret can bind services in.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"serviceid": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A reference to a Service ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN value of a Service ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the Service ID.",
									},
								},
							},
						},
						"user_managed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether the operator secret is user managed.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmCodeEngineSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_secret", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSecretOptions := &codeenginev2.GetSecretOptions{}

	getSecretOptions.SetProjectID(d.Get("project_id").(string))
	getSecretOptions.SetName(d.Get("name").(string))

	secret, _, err := codeEngineClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_secret", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSecretOptions.ProjectID, *getSecretOptions.Name))

	if !core.IsNil(secret.CreatedAt) {
		if err = d.Set("created_at", secret.CreatedAt); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_secret", "read", "set-created_at").GetDiag()
		}
	}

	if !core.IsNil(secret.Data) {
		if err = d.Set("data", secret.Data); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting data: %s", err), "(Data) ibm_code_engine_secret", "read", "set-data").GetDiag()
		}
	}

	if err = d.Set("entity_tag", secret.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_secret", "read", "set-entity_tag").GetDiag()
	}

	if !core.IsNil(secret.Format) {
		if err = d.Set("format", secret.Format); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting format: %s", err), "(Data) ibm_code_engine_secret", "read", "set-format").GetDiag()
		}
	}

	if !core.IsNil(secret.Href) {
		if err = d.Set("href", secret.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_secret", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(secret.ID) {
		if err = d.Set("secret_id", secret.ID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting secret_id: %s", err), "(Data) ibm_code_engine_secret", "read", "set-secret_id").GetDiag()
		}
	}

	if !core.IsNil(secret.Region) {
		if err = d.Set("region", secret.Region); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_secret", "read", "set-region").GetDiag()
		}
	}

	if !core.IsNil(secret.ResourceType) {
		if err = d.Set("resource_type", secret.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_secret", "read", "set-resource_type").GetDiag()
		}
	}

	if !core.IsNil(secret.ServiceAccess) {
		serviceAccess := []map[string]interface{}{}
		serviceAccessMap, err := DataSourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(secret.ServiceAccess)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_secret", "read", "service_access-to-map").GetDiag()
		}
		serviceAccess = append(serviceAccess, serviceAccessMap)
		if err = d.Set("service_access", serviceAccess); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting service_access: %s", err), "(Data) ibm_code_engine_secret", "read", "set-service_access").GetDiag()
		}
	}

	if !core.IsNil(secret.ServiceOperator) {
		serviceOperator := []map[string]interface{}{}
		serviceOperatorMap, err := DataSourceIbmCodeEngineSecretOperatorSecretPropsToMap(secret.ServiceOperator)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_secret", "read", "service_operator-to-map").GetDiag()
		}
		serviceOperator = append(serviceOperator, serviceOperatorMap)
		if err = d.Set("service_operator", serviceOperator); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting service_operator: %s", err), "(Data) ibm_code_engine_secret", "read", "set-service_operator").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(model *codeenginev2.ServiceAccessSecretProps) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	resourceKeyMap, err := DataSourceIbmCodeEngineSecretResourceKeyRefToMap(model.ResourceKey)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_key"] = []map[string]interface{}{resourceKeyMap}
	if model.Role != nil {
		roleMap, err := DataSourceIbmCodeEngineSecretRoleRefToMap(model.Role)
		if err != nil {
			return modelMap, err
		}
		modelMap["role"] = []map[string]interface{}{roleMap}
	}
	serviceInstanceMap, err := DataSourceIbmCodeEngineSecretServiceInstanceRefToMap(model.ServiceInstance)
	if err != nil {
		return modelMap, err
	}
	modelMap["service_instance"] = []map[string]interface{}{serviceInstanceMap}
	if model.Serviceid != nil {
		serviceidMap, err := DataSourceIbmCodeEngineSecretServiceIDRefToMap(model.Serviceid)
		if err != nil {
			return modelMap, err
		}
		modelMap["serviceid"] = []map[string]interface{}{serviceidMap}
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineSecretResourceKeyRefToMap(model *codeenginev2.ResourceKeyRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineSecretRoleRefToMap(model *codeenginev2.RoleRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineSecretServiceInstanceRefToMap(model *codeenginev2.ServiceInstanceRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineSecretServiceIDRefToMap(model *codeenginev2.ServiceIDRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func DataSourceIbmCodeEngineSecretOperatorSecretPropsToMap(model *codeenginev2.OperatorSecretProps) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["apikey_id"] = *model.ApikeyID
	modelMap["resource_group_ids"] = model.ResourceGroupIds
	serviceidMap, err := DataSourceIbmCodeEngineSecretServiceIDRefToMap(model.Serviceid)
	if err != nil {
		return modelMap, err
	}
	modelMap["serviceid"] = []map[string]interface{}{serviceidMap}
	modelMap["user_managed"] = *model.UserManaged
	return modelMap, nil
}
