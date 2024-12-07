// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmUsernamePasswordSecretMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmUsernamePasswordSecretMetadataRead,

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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"password_generation_policy": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Policy for auto-generated passwords.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The length of auto-generated passwords.",
						},
						"include_digits": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Include digits in auto-generated passwords.",
						},
						"include_symbols": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Include symbols in auto-generated passwords.",
						},
						"include_uppercase": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Include uppercase letters in auto-generated passwords.",
						},
					},
				},
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
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"next_rotation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
			},
		},
	}
}

func dataSourceIbmSmUsernamePasswordSecretMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretMetadataOptions := &secretsmanagerv2.GetSecretMetadataOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretMetadataOptions.SetID(secretId)

	usernamePasswordSecretMetadataIntf, response, err := secretsManagerClient.GetSecretMetadataWithContext(context, getSecretMetadataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretMetadataWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	usernamePasswordSecretMetadata := usernamePasswordSecretMetadataIntf.(*secretsmanagerv2.UsernamePasswordSecretMetadata)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", usernamePasswordSecretMetadata.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(usernamePasswordSecretMetadata.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", usernamePasswordSecretMetadata.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if usernamePasswordSecretMetadata.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(usernamePasswordSecretMetadata.CustomMetadata))
		for k, v := range usernamePasswordSecretMetadata.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", usernamePasswordSecretMetadata.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", usernamePasswordSecretMetadata.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(usernamePasswordSecretMetadata.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", usernamePasswordSecretMetadata.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", usernamePasswordSecretMetadata.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", usernamePasswordSecretMetadata.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(usernamePasswordSecretMetadata.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", usernamePasswordSecretMetadata.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(usernamePasswordSecretMetadata.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(usernamePasswordSecretMetadata.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	rotation := []map[string]interface{}{}
	if usernamePasswordSecretMetadata.Rotation != nil {
		modelMap, err := dataSourceIbmSmUsernamePasswordSecretMetadataRotationPolicyToMap(usernamePasswordSecretMetadata.Rotation)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		rotation = append(rotation, modelMap)
	}
	if err = d.Set("rotation", rotation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	passwordPolicy := []map[string]interface{}{}
	if usernamePasswordSecretMetadata.PasswordGenerationPolicy != nil {
		modelMap, err := passwordGenerationPolicyToMap(usernamePasswordSecretMetadata.PasswordGenerationPolicy)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		passwordPolicy = append(passwordPolicy, modelMap)
	}
	if err = d.Set("password_generation_policy", passwordPolicy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting password_generation_policy"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(usernamePasswordSecretMetadata.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("next_rotation_date", DateTimeToRFC3339(usernamePasswordSecretMetadata.NextRotationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting next_rotation_date"), fmt.Sprintf("(Data) %s_metadata", UsernamePasswordSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSmUsernamePasswordSecretMetadataRotationPolicyToMap(model secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.CommonRotationPolicy); ok {
		return dataSourceIbmSmUsernamePasswordSecretMetadataCommonRotationPolicyToMap(model.(*secretsmanagerv2.CommonRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.RotationPolicy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.RotationPolicy)
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
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.RotationPolicyIntf subtype encountered")
	}
}

func dataSourceIbmSmUsernamePasswordSecretMetadataCommonRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
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
