// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package hpcs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/ibm-hpcs-uko-sdk/ukov4"
)

func DataSourceIbmManagedKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIbmManagedKeyRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the UKO instance this resource exists in.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region of the UKO instance this resource exists in.",
			},
			"key_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the key.",
			},
			"uko_vault": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the Vault in which the update is to take place.",
			},
			"vault": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Reference to a vault.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the referenced vault.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL that uniquely identifies your cloud resource.",
						},
					},
				},
			},
			"template": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Reference to a key template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the key template.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL that uniquely identifies your cloud resource.",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the managed key.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The label of the key.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the key.",
			},
			"size": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The size of the underlying cryptographic key or key pair. E.g. \"256\" for AES keys, or \"2048\" for RSA.",
			},
			"algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The algorithm of the key.",
			},
			"verification_patterns": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of verification patterns of the key (e.g. public key hash for RSA keys).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The method used for calculating the verification pattern.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The calculated value.",
						},
					},
				},
			},
			"activation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "First day when the key is active.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last day when the key is active.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Key-value pairs associated with the key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a tag.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value of a tag.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the key was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the key was last updated.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that created the key.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that last updated the key.",
			},
			"referenced_keystores": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "referenced keystores.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the target keystore.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of keystore.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL that uniquely identifies your cloud resource.",
						},
					},
				},
			},
			"instances": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "key instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"label_in_keystore": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The label of the key.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the key instance.",
						},
						"keystore": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Description of properties of a key within the context of keystores.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of keystore.",
									},
								},
							},
						},
						"google_key_protection_level": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"google_key_purpose": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"google_kms_algorithm": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
		},
	}
}

func DataSourceIbmManagedKeyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getManagedKeyOptions := &ukov4.GetManagedKeyOptions{}

	region := d.Get("region").(string)
	instance_id := d.Get("instance_id").(string)
	vault_id := d.Get("uko_vault").(string)
	key_id := d.Get("key_id").(string)
	getManagedKeyOptions.SetID(key_id)
	getManagedKeyOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	managedKey, response, err := ukoClient.GetManagedKeyWithContext(context, getManagedKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetManagedKeyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetManagedKeyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, vault_id, *getManagedKeyOptions.ID))

	vault := []map[string]interface{}{}
	if managedKey.Vault != nil {
		modelMap, err := DataSourceIbmManagedKeyVaultReferenceToMap(managedKey.Vault)
		if err != nil {
			return diag.FromErr(err)
		}
		vault = append(vault, modelMap)
	}
	if err = d.Set("vault", vault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault %s", err))
	}

	template := []map[string]interface{}{}
	if managedKey.Template != nil {
		modelMap, err := DataSourceIbmManagedKeyTemplateReferenceToMap(managedKey.Template)
		if err != nil {
			return diag.FromErr(err)
		}
		template = append(template, modelMap)
	}
	if err = d.Set("template", template); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting template %s", err))
	}

	if err = d.Set("description", managedKey.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("label", managedKey.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}

	if err = d.Set("state", managedKey.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	if err = d.Set("size", managedKey.Size); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}

	if err = d.Set("algorithm", managedKey.Algorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting algorithm: %s", err))
	}

	verificationPatterns := []map[string]interface{}{}
	if managedKey.VerificationPatterns != nil {
		for _, modelItem := range managedKey.VerificationPatterns {
			modelMap, err := DataSourceIbmManagedKeyKeyVerificationPatternToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			verificationPatterns = append(verificationPatterns, modelMap)
		}
	}
	if err = d.Set("verification_patterns", verificationPatterns); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting verification_patterns %s", err))
	}

	if err = d.Set("activation_date", flex.DateToString(managedKey.ActivationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting activation_date: %s", err))
	}

	if err = d.Set("expiration_date", flex.DateToString(managedKey.ExpirationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration_date: %s", err))
	}

	tags := []map[string]interface{}{}
	if managedKey.Tags != nil {
		for _, modelItem := range managedKey.Tags {
			modelMap, err := DataSourceIbmManagedKeyTagToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			tags = append(tags, modelMap)
		}
	}
	if err = d.Set("tags", tags); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tags %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(managedKey.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(managedKey.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("created_by", managedKey.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_by", managedKey.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	referencedKeystores := []map[string]interface{}{}
	if managedKey.ReferencedKeystores != nil {
		for _, modelItem := range managedKey.ReferencedKeystores {
			modelMap, err := DataSourceIbmManagedKeyTargetKeystoreReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			referencedKeystores = append(referencedKeystores, modelMap)
		}
	}
	if err = d.Set("referenced_keystores", referencedKeystores); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referenced_keystores %s", err))
	}

	instances := []map[string]interface{}{}
	if managedKey.Instances != nil {
		for _, modelItem := range managedKey.Instances {
			// TODO: I'm worried about this line
			modelMap, err := dataSourceIbmHpcsManagedKeyKeyInstanceToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			instances = append(instances, modelMap)
		}
	}
	if err = d.Set("instances", instances); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instances %s", err))
	}

	if err = d.Set("href", managedKey.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	return nil
}

func DataSourceIbmManagedKeyVaultReferenceToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func DataSourceIbmManagedKeyTemplateReferenceToMap(model *ukov4.TemplateReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func DataSourceIbmManagedKeyKeyVerificationPatternToMap(model *ukov4.KeyVerificationPattern) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Method != nil {
		modelMap["method"] = *model.Method
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmManagedKeyTagToMap(model *ukov4.Tag) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmManagedKeyTargetKeystoreReferenceToMap(model *ukov4.TargetKeystoreReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

// func DataSourceIbmManagedKeyKeyInstanceToMap(model *ukov4.KeyInstance) (map[string]interface{}, error) {
func dataSourceIbmHpcsManagedKeyKeyInstanceToMap(model ukov4.KeyInstanceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeyInstanceGoogleKms); ok {
		return dataSourceIbmHpcsManagedKeyKeyInstanceGoogleKmsToMap(model.(*ukov4.KeyInstanceGoogleKms))
	} else if _, ok := model.(*ukov4.KeyInstanceAwsKms); ok {
		return dataSourceIbmHpcsManagedKeyKeyInstanceAwsKmsToMap(model.(*ukov4.KeyInstanceAwsKms))
	} else if _, ok := model.(*ukov4.KeyInstanceIbmCloudKms); ok {
		return dataSourceIbmHpcsManagedKeyKeyInstanceIbmCloudKmsToMap(model.(*ukov4.KeyInstanceIbmCloudKms))
	} else if _, ok := model.(*ukov4.KeyInstanceAzure); ok {
		return dataSourceIbmHpcsManagedKeyKeyInstanceAzureToMap(model.(*ukov4.KeyInstanceAzure))
	} else if _, ok := model.(*ukov4.KeyInstance); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeyInstance)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.LabelInKeystore != nil {
			modelMap["label_in_keystore"] = *model.LabelInKeystore
		}
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Keystore != nil {
			keystoreMap, err := dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
			if err != nil {
				return modelMap, err
			}
			modelMap["keystore"] = []map[string]interface{}{keystoreMap}
		}
		if model.GoogleKeyProtectionLevel != nil {
			modelMap["google_key_protection_level"] = *model.GoogleKeyProtectionLevel
		}
		if model.GoogleKeyPurpose != nil {
			modelMap["google_key_purpose"] = *model.GoogleKeyPurpose
		}
		if model.GoogleKmsAlgorithm != nil {
			modelMap["google_kms_algorithm"] = *model.GoogleKmsAlgorithm
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ukov4.KeyInstanceIntf subtype encountered")
	}
}

func dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model *ukov4.InstanceInKeystore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = *model.Group
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func dataSourceIbmHpcsManagedKeyKeyInstanceGoogleKmsToMap(model *ukov4.KeyInstanceGoogleKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LabelInKeystore != nil {
		modelMap["label_in_keystore"] = *model.LabelInKeystore
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Keystore != nil {
		// TODO: I'm worried about this line
		keystoreMap, err := dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
		if err != nil {
			return modelMap, err
		}
		modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	}
	if model.GoogleKeyProtectionLevel != nil {
		modelMap["google_key_protection_level"] = *model.GoogleKeyProtectionLevel
	}
	if model.GoogleKeyPurpose != nil {
		modelMap["google_key_purpose"] = *model.GoogleKeyPurpose
	}
	if model.GoogleKmsAlgorithm != nil {
		modelMap["google_kms_algorithm"] = *model.GoogleKmsAlgorithm
	}
	return modelMap, nil
}

// func DataSourceIbmManagedKeyInstanceInKeystoreToMap(model *ukov4.InstanceInKeystore) (map[string]interface{}, error) {
func dataSourceIbmHpcsManagedKeyKeyInstanceAwsKmsToMap(model *ukov4.KeyInstanceAwsKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LabelInKeystore != nil {
		modelMap["label_in_keystore"] = *model.LabelInKeystore
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Keystore != nil {
		keystoreMap, err := dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
		if err != nil {
			return modelMap, err
		}
		modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	}
	return modelMap, nil
}

func dataSourceIbmHpcsManagedKeyKeyInstanceIbmCloudKmsToMap(model *ukov4.KeyInstanceIbmCloudKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LabelInKeystore != nil {
		modelMap["label_in_keystore"] = *model.LabelInKeystore
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Keystore != nil {
		keystoreMap, err := dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
		if err != nil {
			return modelMap, err
		}
		modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	}
	return modelMap, nil
}

func dataSourceIbmHpcsManagedKeyKeyInstanceAzureToMap(model *ukov4.KeyInstanceAzure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LabelInKeystore != nil {
		modelMap["label_in_keystore"] = *model.LabelInKeystore
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Keystore != nil {
		keystoreMap, err := dataSourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
		if err != nil {
			return modelMap, err
		}
		modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	}
	return modelMap, nil
}
