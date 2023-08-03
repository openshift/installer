// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineSecretRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"data": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Data container that allows to specify config parameters and their values as a key-value map. Each key field must consist of alphanumeric characters, `-`, `_` or `.` and must not be exceed a max length of 253 characters. Each value field can consists of any character and must not be exceed a max length of 1048576 characters.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the secret instance, which is used to achieve optimistic locking.",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specify the format of the secret.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new secret,  a URL is created identifying the location of the instance.",
			},
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the secret.",
			},
			"service_access": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Properties for Service Access Secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The service credential associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the service credential associated with the secret.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the service credential associated with the secret.",
									},
								},
							},
						},
						"role": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A reference to the Role and Role CRN for service binding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN of the IAM Role for thise service access secret.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role of the service credential.",
									},
								},
							},
						},
						"service_instance": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IBM Cloud service instance associated with the secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the IBM Cloud service instance associated with the secret.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of IBM Cloud service associated with the secret.",
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

func dataSourceIbmCodeEngineSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getSecretOptions := &codeenginev2.GetSecretOptions{}

	getSecretOptions.SetProjectID(d.Get("project_id").(string))
	getSecretOptions.SetName(d.Get("name").(string))

	secret, response, err := codeEngineClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecretWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSecretOptions.ProjectID, *getSecretOptions.Name))

	if err = d.Set("created_at", secret.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if secret.Data != nil {
		if err = d.Set("data", secret.Data); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data %s", err))
		}
	}

	if err = d.Set("entity_tag", secret.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("format", secret.Format); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting format: %s", err))
	}

	if err = d.Set("href", secret.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("secret_id", secret.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_id: %s", err))
	}

	if err = d.Set("resource_type", secret.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	serviceAccess := []map[string]interface{}{}
	if secret.ServiceAccess != nil {
		modelMap, err := dataSourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(secret.ServiceAccess)
		if err != nil {
			return diag.FromErr(err)
		}
		serviceAccess = append(serviceAccess, modelMap)
	}
	if err = d.Set("service_access", serviceAccess); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_access %s", err))
	}

	return nil
}

func dataSourceIbmCodeEngineSecretServiceAccessSecretPropsToMap(model *codeenginev2.ServiceAccessSecretProps) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	resourceKeyMap, err := dataSourceIbmCodeEngineSecretResourceKeyRefToMap(model.ResourceKey)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_key"] = []map[string]interface{}{resourceKeyMap}
	if model.Role != nil {
		roleMap, err := dataSourceIbmCodeEngineSecretRoleRefToMap(model.Role)
		if err != nil {
			return modelMap, err
		}
		modelMap["role"] = []map[string]interface{}{roleMap}
	}
	serviceInstanceMap, err := dataSourceIbmCodeEngineSecretServiceInstanceRefToMap(model.ServiceInstance)
	if err != nil {
		return modelMap, err
	}
	modelMap["service_instance"] = []map[string]interface{}{serviceInstanceMap}
	return modelMap, nil
}

func dataSourceIbmCodeEngineSecretResourceKeyRefToMap(model *codeenginev2.ResourceKeyRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	return modelMap, nil
}

func dataSourceIbmCodeEngineSecretRoleRefToMap(model *codeenginev2.RoleRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = model.Crn
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	return modelMap, nil
}

func dataSourceIbmCodeEngineSecretServiceInstanceRefToMap(model *codeenginev2.ServiceInstanceRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}
