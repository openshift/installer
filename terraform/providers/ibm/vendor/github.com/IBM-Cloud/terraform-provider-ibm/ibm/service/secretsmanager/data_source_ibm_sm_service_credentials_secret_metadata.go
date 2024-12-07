// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func DataSourceIbmSmServiceCredentialsSecretMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmServiceCredentialsSecretMetadataRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the secret.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
			},
			"state_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A text representation of the secret state.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"versions_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of versions of the secret.",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time-to-live (TTL) or lease duration to assign to generated credentials.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"rotation": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Determines whether Secrets Manager rotates your secrets automatically.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotate": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether Secrets Manager rotates your secret automatically.Default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined interval.",
						},
						"interval": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The length of the secret rotation time interval.",
						},
						"unit": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The units for the secret rotation time interval.",
						},
					},
				},
			},
			"next_rotation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that the secret is scheduled for automatic rotation. The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
			},
			"source_service": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The properties required for creating the service credentials for the specified source service instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source service instance identifier.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A CRN that uniquely identifies a service credentials target.",
									},
								},
							},
						},
						"role": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The service-specific custom role object, CRN role is accepted. Refer to the service’s documentation for supported roles.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN role identifier for creating a service-id.",
									},
								},
							},
						},
						"iam": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source service IAM data is returned in case IAM credentials where created for this secret.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"apikey": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The IAM apikey metadata for the IAM credentials that were generated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM API key name for the generated service credentials.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM API key description for the generated service credentials.",
												},
											},
										},
									},
									"role": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The IAM role for the generate service credentials.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM role CRN assigned to the generated service credentials.",
												},
											},
										},
									},
									"serviceid": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The IAM serviceid for the generated service credentials.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IAM Service ID CRN.",
												},
											},
										},
									},
								},
							},
						},
						"resource_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source service resource key data of the generated service credentials.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource key CRN of the generated service credentials.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource key name of the generated service credentials.",
									},
								},
							},
						},
						"parameters": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The collection of parameters for the service credentials target.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSmServiceCredentialsSecretMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretMetadataOptions := &secretsmanagerv2.GetSecretMetadataOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretMetadataOptions.SetID(secretId)

	ServiceCredentialsSecretMetadataIntf, response, err := secretsManagerClient.GetSecretMetadataWithContext(context, getSecretMetadataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretMetadataWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	ServiceCredentialsSecretMetadata := ServiceCredentialsSecretMetadataIntf.(*secretsmanagerv2.ServiceCredentialsSecretMetadata)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", ServiceCredentialsSecretMetadata.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(ServiceCredentialsSecretMetadata.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", ServiceCredentialsSecretMetadata.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if ServiceCredentialsSecretMetadata.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(ServiceCredentialsSecretMetadata.CustomMetadata))
		for k, v := range ServiceCredentialsSecretMetadata.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", ServiceCredentialsSecretMetadata.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", ServiceCredentialsSecretMetadata.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if ServiceCredentialsSecretMetadata.Labels != nil {
		if err = d.Set("labels", ServiceCredentialsSecretMetadata.Labels); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting labels"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("locks_total", flex.IntValue(ServiceCredentialsSecretMetadata.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", ServiceCredentialsSecretMetadata.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", ServiceCredentialsSecretMetadata.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", ServiceCredentialsSecretMetadata.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(ServiceCredentialsSecretMetadata.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", ServiceCredentialsSecretMetadata.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(ServiceCredentialsSecretMetadata.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(ServiceCredentialsSecretMetadata.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("ttl", ServiceCredentialsSecretMetadata.TTL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ttl"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	rotation := []map[string]interface{}{}
	if ServiceCredentialsSecretMetadata.Rotation != nil {
		modelMap, err := dataSourceIbmSmServiceCredentialsSecretMetadataRotationPolicyToMap(ServiceCredentialsSecretMetadata.Rotation.(*secretsmanagerv2.RotationPolicy))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		rotation = append(rotation, modelMap)
	}
	if err = d.Set("rotation", rotation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("next_rotation_date", DateTimeToRFC3339(ServiceCredentialsSecretMetadata.NextRotationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting next_rotation_date"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	sourceServiceMap, err := dataSourceIbmSmServiceCredentialsSecretMetadataSourceServiceToMap(ServiceCredentialsSecretMetadata.SourceService)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	if len(sourceServiceMap) > 0 {
		if err = d.Set("source_service", []map[string]interface{}{sourceServiceMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting source_service"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if ServiceCredentialsSecretMetadata.ExpirationDate != nil {
		if err = d.Set("expiration_date", DateTimeToRFC3339(ServiceCredentialsSecretMetadata.ExpirationDate)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s_metadata", ServiceCredentialsSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func dataSourceIbmSmServiceCredentialsSecretMetadataRotationPolicyToMap(model *secretsmanagerv2.RotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	return modelMap, nil
}

func dataSourceIbmSmServiceCredentialsSecretMetadataSourceServiceToMap(sourceService *secretsmanagerv2.ServiceCredentialsSecretSourceServiceRO) (map[string]interface{}, error) {
	mainModelMap := make(map[string]interface{})
	if sourceService.Instance != nil {
		instanceMap := make(map[string]interface{})
		instanceModel := sourceService.Instance
		if instanceModel.Crn != nil {
			instanceMap["crn"] = instanceModel.Crn
		}
		mainModelMap["instance"] = []map[string]interface{}{instanceMap}
	}

	if sourceService.Role != nil {
		roleMap := make(map[string]interface{})
		roleModel := sourceService.Role
		if roleModel.Crn != nil {
			roleMap["crn"] = roleModel.Crn
		}
		mainModelMap["role"] = []map[string]interface{}{roleMap}
	}

	if sourceService.Iam != nil {
		iamMap := make(map[string]interface{})
		iamModel := sourceService.Iam

		// apikey
		if iamModel.Apikey != nil {
			iamApikeyMap := make(map[string]interface{})
			iamApikeyModel := iamModel.Apikey
			if iamApikeyModel.Name != nil {
				iamApikeyMap["name"] = iamApikeyModel.Name
			}
			if iamApikeyModel.Description != nil {
				iamApikeyMap["description"] = iamApikeyModel.Description
			}
			iamMap["apikey"] = []map[string]interface{}{iamApikeyMap}
		}

		// role
		if iamModel.Role != nil {
			iamRoleMap := make(map[string]interface{})
			iamRoleModel := iamModel.Role
			if iamRoleModel.Crn != nil {
				iamRoleMap["crn"] = iamRoleModel.Crn
			}
			iamMap["role"] = []map[string]interface{}{iamRoleMap}
		}

		// service id
		if iamModel.Serviceid != nil {
			iamServiceidMap := make(map[string]interface{})
			iamServiceidModel := iamModel.Serviceid
			if iamServiceidModel.Crn != nil {
				iamServiceidMap["crn"] = iamServiceidModel.Crn
			}
			iamMap["serviceid"] = []map[string]interface{}{iamServiceidMap}
		}

		mainModelMap["iam"] = []map[string]interface{}{iamMap}

	}

	if sourceService.ResourceKey != nil {
		resourceKeyMap := make(map[string]interface{})
		resourceKeyModel := sourceService.ResourceKey
		if resourceKeyModel.Crn != nil {
			resourceKeyMap["crn"] = resourceKeyModel.Crn
		}
		if resourceKeyModel.Name != nil {
			resourceKeyMap["name"] = resourceKeyModel.Name
		}
		mainModelMap["resource_key"] = []map[string]interface{}{resourceKeyMap}
	}

	if sourceService.Parameters != nil {
		parametersMap := sourceService.Parameters.GetProperties()
		for k, v := range parametersMap {
			parametersMap[k] = fmt.Sprint(v)
		}
		if sourceService.Parameters.ServiceidCrn != nil {
			if len(parametersMap) == 0 {
				parametersMap = make(map[string]interface{})
			}
			parametersMap["serviceid_crn"] = sourceService.Parameters.ServiceidCrn
		}
		mainModelMap["parameters"] = parametersMap
	}

	return mainModelMap, nil
}
