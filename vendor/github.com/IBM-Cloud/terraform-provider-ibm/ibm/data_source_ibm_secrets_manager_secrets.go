// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMSecretsManagerSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSecretsManagerSecretsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secrets Manager instance GUID",
			},
			"secret_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeDataSourceValidator("ibm_secrets_manager_secrets", "secret_type"),
				Description:  "The secret type. Supported options include: arbitrary, iam_credentials, username_password.",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeDataSourceValidator("ibm_secrets_manager_secrets", "endpoint_type"),
				Description:  "Endpoint Type. 'public' or 'private'",
				Default:      "public",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The metadata that describes the resource array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resources in the resource array.",
						},
						"collection_total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of elements in the resource array.",
						},
					},
				},
			},
			"secrets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A collection of secret resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MIME type that represents the secret.",
						},
						"secret_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A human-readable alias to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.",
						},
						"secret_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret group to assign to this secret.If you omit this parameter, your secret is assigned to the `default` secret group.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels that you can use to filter for secrets in your instance.Up to 30 labels can be created. Labels can be between 2-30 characters, including spaces. Special characters not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
						"secret_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.",
						},
						"creation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date the secret was created. The date format follows RFC 3339.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the entity that created the secret.",
						},
						"last_update_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the actual secret is modified. The date format follows RFC 3339.",
						},
						"versions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An array that contains metadata for each secret version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the secret version.",
									},
									"creation_date": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date that the version of the secret was created.",
									},
									"created_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the entity that created the secret.",
									},
									"auto_rotated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the version of the secret was created by automatic rotation.",
									},
								},
							},
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date the secret material expires. The date format follows RFC 3339.You can set an expiration date on supported secret types at their creation. If you create a secret without specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the following secret types:- `arbitrary`- `username_password`.",
						},
						"payload": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The new secret data to assign to an `arbitrary` secret.",
						},
						"secret_data": &schema.Schema{
							Type:        schema.TypeMap,
							Sensitive:   true,
							Computed:    true,
							Description: "The secret data object",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The username to assign to this secret.",
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The password to assign to this secret.",
						},
						"next_rotation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that can be auto-rotated and have an existing rotation policy.",
						},
						"ttl": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time-to-live (TTL) or lease duration to assign to generated credentials.For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m` or `24h`.",
						},
						"access_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The access groups that define the capabilities of the service ID and API key that are generated for an`iam_credentials` secret.**Tip:** To find the ID of an access group, go to **Manage > Access (IAM) > Access groups** in the IBM Cloud console. Select the access group to inspect, and click **Details** to view its ID.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The API key that is generated for this secret.After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you want to continue to use the same API key for future read operations, see the `reuse_api_key` field.",
						},
						"service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service ID under which the API key (see the `api_key` field) is created. This service ID is added to the access groups that you assign for this secret.",
						},
						"reuse_api_key": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "(IAM credentials) Reuse the service ID and API key for future read operations.",
						},
					},
				},
			},
		},
	}
}

func datasourceIBMSecretsManagerSecretsValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	secretType := "arbitrary, iam_credentials, username_password"
	endpointType := "public, private"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "secret_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              secretType})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "endpoint_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              endpointType})

	ibmSecretsManagerSecretsdatasourceValidator := ResourceValidator{ResourceName: "ibm_secrets_manager_secrets", Schema: validateSchema}
	return &ibmSecretsManagerSecretsdatasourceValidator
}

func dataSourceIBMSecretsManagerSecretsRead(d *schema.ResourceData, meta interface{}) error {
	bluemixSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	region := bluemixSession.Config.Region

	secretsManagerClient, err := meta.(ClientSession).SecretsManagerV1()
	if err != nil {
		return err
	}
	rContollerClient, err := meta.(ClientSession).ResourceControllerAPIV2()
	if err != nil {
		return err
	}

	instanceID := d.Get("instance_id").(string)
	endpointType := d.Get("endpoint_type").(string)
	var smEndpointURL string

	rContollerAPI := rContollerClient.ResourceServiceInstanceV2()

	instanceData, err := rContollerAPI.GetInstance(instanceID)
	if err != nil {
		return err
	}
	instanceCRN := instanceData.Crn.String()

	crnData := strings.Split(instanceCRN, ":")

	if crnData[4] == "secrets-manager" {
		if endpointType == "private" {
			smEndpointURL = "https://" + instanceID + "private" + "." + region + ".secrets-manager.appdomain.cloud"
		} else {
			smEndpointURL = "https://" + instanceID + "." + region + ".secrets-manager.appdomain.cloud"
		}
		smUrl := envFallBack([]string{"IBMCLOUD_SECRETS_MANAGER_API_ENDPOINT"}, smEndpointURL)
		secretsManagerClient.Service.Options.URL = smUrl
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

	listAllSecretsOptions := &secretsmanagerv1.ListAllSecretsOptions{}

	listSecrets, response, err := secretsManagerClient.ListAllSecretsWithContext(context.TODO(), listAllSecretsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListAllSecretsWithContext failed %s\n%s", err, response)
		return err
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchResources []secretsmanagerv1.SecretResourceIntf
	var secretType string
	var suppliedFilter bool

	if v, ok := d.GetOk("secret_type"); ok {
		secretType = v.(string)
		suppliedFilter = true
		for _, data := range listSecrets.Resources {
			if rdata, ok := data.(*secretsmanagerv1.SecretResource); ok {
				if *rdata.SecretType == secretType {
					matchResources = append(matchResources, data)
				}
			}
		}
	} else {
		matchResources = listSecrets.Resources
	}
	listSecrets.Resources = matchResources

	if len(listSecrets.Resources) == 0 {
		return fmt.Errorf("no Resources found with secretType %s\nIf not specified, please specify more filters", secretType)
	}

	if suppliedFilter {
		d.SetId(secretType)
	} else {
		d.SetId(dataSourceIBMSecretsManagerSecretsID(d))
	}

	if listSecrets.Metadata != nil {
		err = d.Set("metadata", dataSourceListSecretsFlattenMetadata(*listSecrets.Metadata))
		if err != nil {
			return fmt.Errorf("Error setting metadata %s", err)
		}
	}

	if listSecrets.Resources != nil {
		err = d.Set("secrets", dataSourceListSecretsFlattenResources(listSecrets.Resources))
		if err != nil {
			return fmt.Errorf("Error setting resources %s", err)
		}
	}

	return nil
}

// dataSourceIBMSecretsManagerSecretsID returns a reasonable ID for the list.
func dataSourceIBMSecretsManagerSecretsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceListSecretsFlattenMetadata(result secretsmanagerv1.CollectionMetadata) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceListSecretsMetadataToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceListSecretsMetadataToMap(metadataItem secretsmanagerv1.CollectionMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.CollectionType != nil {
		metadataMap["collection_type"] = metadataItem.CollectionType
	}
	if metadataItem.CollectionTotal != nil {
		metadataMap["collection_total"] = metadataItem.CollectionTotal
	}

	return metadataMap
}

func dataSourceListSecretsFlattenResources(result []secretsmanagerv1.SecretResourceIntf) (resources []map[string]interface{}) {
	for _, resourcesItem := range result {
		if ritem, ok := resourcesItem.(*secretsmanagerv1.SecretResource); ok {
			resources = append(resources, dataSourceListSecretsResourcesToMap(*ritem))
		}
	}

	return resources
}

func dataSourceListSecretsResourcesToMap(resourcesItem secretsmanagerv1.SecretResource) (resourcesMap map[string]interface{}) {
	resourcesMap = map[string]interface{}{}

	if resourcesItem.Type != nil {
		resourcesMap["type"] = *resourcesItem.Type
	}
	if resourcesItem.ID != nil {
		resourcesMap["secret_id"] = *resourcesItem.ID
	}
	if resourcesItem.Name != nil {
		resourcesMap["name"] = *resourcesItem.Name
	}
	if resourcesItem.Description != nil {
		resourcesMap["description"] = *resourcesItem.Description
	}
	if resourcesItem.SecretGroupID != nil {
		resourcesMap["secret_group_id"] = *resourcesItem.SecretGroupID
	}
	if resourcesItem.Labels != nil {
		resourcesMap["labels"] = resourcesItem.Labels
	}
	if resourcesItem.State != nil {
		resourcesMap["state"] = *resourcesItem.State
	}
	if resourcesItem.StateDescription != nil {
		resourcesMap["state_description"] = *resourcesItem.StateDescription
	}
	if resourcesItem.SecretType != nil {
		resourcesMap["secret_type"] = *resourcesItem.SecretType
	}
	if resourcesItem.CRN != nil {
		resourcesMap["crn"] = *resourcesItem.CRN
	}
	if resourcesItem.CreationDate != nil {
		resourcesMap["creation_date"] = (*resourcesItem.CreationDate).String()
	}
	if resourcesItem.CreatedBy != nil {
		resourcesMap["created_by"] = *resourcesItem.CreatedBy
	}
	if resourcesItem.LastUpdateDate != nil {
		resourcesMap["last_update_date"] = (*resourcesItem.LastUpdateDate).String()
	}
	if resourcesItem.Versions != nil {
		versionsList := []map[string]interface{}{}
		for _, versionsItem := range resourcesItem.Versions {
			versionsList = append(versionsList, dataSourceListSecretsResourcesVersionsToMap(versionsItem))
		}
		resourcesMap["versions"] = versionsList
	}
	if resourcesItem.ExpirationDate != nil {
		resourcesMap["expiration_date"] = (*resourcesItem.ExpirationDate).String()
	}
	if resourcesItem.Payload != nil {
		resourcesMap["payload"] = *resourcesItem.Payload
	}
	if resourcesItem.SecretData != nil {
		secretData := resourcesItem.SecretData.(map[string]interface{})
		resourcesMap["secret_data"] = secretData
		if *resourcesItem.SecretType == "username_password" {
			resourcesMap["username"] = secretData["username"].(string)
			resourcesMap["password"] = secretData["password"].(string)
		} else if *resourcesItem.SecretType == "arbitrary" {
			resourcesMap["payload"] = secretData["payload"].(string)
		}
	}
	if resourcesItem.NextRotationDate != nil {
		resourcesMap["next_rotation_date"] = (*resourcesItem.NextRotationDate).String()
	}
	if resourcesItem.TTL != nil {
		resourcesMap["ttl"] = fmt.Sprintf("%v", resourcesItem.TTL)
	}
	if resourcesItem.AccessGroups != nil {
		resourcesMap["access_groups"] = resourcesItem.AccessGroups
	}
	if resourcesItem.APIKey != nil {
		resourcesMap["api_key"] = *resourcesItem.APIKey
	}
	if resourcesItem.ServiceID != nil {
		resourcesMap["service_id"] = *resourcesItem.ServiceID
	}
	if resourcesItem.ReuseAPIKey != nil {
		resourcesMap["reuse_api_key"] = *resourcesItem.ReuseAPIKey
	}

	return resourcesMap
}

func dataSourceListSecretsResourcesVersionsToMap(versionsItem secretsmanagerv1.SecretVersion) (versionsMap map[string]interface{}) {
	versionsMap = map[string]interface{}{}

	if versionsItem.ID != nil {
		versionsMap["id"] = *versionsItem.ID
	}
	if versionsItem.CreationDate != nil {
		versionsMap["creation_date"] = (*versionsItem.CreationDate).String()
	}
	if versionsItem.CreatedBy != nil {
		versionsMap["created_by"] = *versionsItem.CreatedBy
	}
	if versionsItem.AutoRotated != nil {
		versionsMap["auto_rotated"] = *versionsItem.AutoRotated
	}

	return versionsMap
}
