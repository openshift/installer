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

func DataSourceIbmKeyTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIbmKeyTemplateRead,

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
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the template.",
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
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of the key template. Every time the key template is updated, the version will be updated automatically.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the key template.",
			},
			"key": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Properties describing the properties of the managed key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"activation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key activation date can be provided as a period definition (e.g. PY1 means 1 year).",
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key expiration date can be provided as a period definition (e.g. PY1 means 1 year).",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state that the key will be in after generation.",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the key template.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the key template was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the key template was updated.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that created the key template.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that updated the key.",
			},
			"keystores": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Which keystore group to distribute the key to.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of keystore.",
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

func DataSourceIbmKeyTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getKeyTemplateOptions := &ukov4.GetKeyTemplateOptions{}

	region := d.Get("region").(string)
	instance_id := d.Get("instance_id").(string)
	vault_id := d.Get("uko_vault").(string)
	template_id := d.Get("template_id").(string)
	getKeyTemplateOptions.SetID(template_id)
	getKeyTemplateOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	template, response, err := ukoClient.GetKeyTemplateWithContext(context, getKeyTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] GetKeyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetKeyTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, d.Get("uko_vault").(string), *getKeyTemplateOptions.ID))

	vault := []map[string]interface{}{}
	if template.Vault != nil {
		modelMap, err := DataSourceIbmKeyTemplateVaultReferenceToMap(template.Vault)
		if err != nil {
			return diag.FromErr(err)
		}
		vault = append(vault, modelMap)
	}
	if err = d.Set("vault", vault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault %s", err))
	}

	if err = d.Set("version", flex.IntValue(template.Version)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	if err = d.Set("name", template.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	key := []map[string]interface{}{}
	if template.Key != nil {
		modelMap, err := DataSourceIbmKeyTemplateKeyPropertiesToMap(template.Key)
		if err != nil {
			return diag.FromErr(err)
		}
		key = append(key, modelMap)
	}
	if err = d.Set("key", key); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key %s", err))
	}

	if err = d.Set("description", template.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(template.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(template.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("created_by", template.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_by", template.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	keystores := []map[string]interface{}{}
	if template.Keystores != nil {
		for _, modelItem := range template.Keystores {
			// TODO: Worried about this line
			modelMap, err := dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			keystores = append(keystores, modelMap)
		}
	}
	if err = d.Set("keystores", keystores); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting keystores %s", err))
	}

	if err = d.Set("href", template.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	return nil
}

func DataSourceIbmKeyTemplateVaultReferenceToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
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

func DataSourceIbmKeyTemplateKeyPropertiesToMap(model *ukov4.KeyProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Size != nil {
		modelMap["size"] = *model.Size
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = *model.Algorithm
	}
	if model.ActivationDate != nil {
		modelMap["activation_date"] = *model.ActivationDate
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = *model.ExpirationDate
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	return modelMap, nil
}

// func DataSourceIbmKeyTemplateKeystoresPropertiesToMap(model *ukov4.KeystoresProperties) (map[string]interface{}, error) {
func dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateToMap(model ukov4.KeystoresPropertiesCreateIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeystoresPropertiesCreateGoogleKms); ok {
		return dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateGoogleKmsToMap(model.(*ukov4.KeystoresPropertiesCreateGoogleKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateAwsKms); ok {
		return dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAwsKmsToMap(model.(*ukov4.KeystoresPropertiesCreateAwsKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateIbmCloudKms); ok {
		return dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateIbmCloudKmsToMap(model.(*ukov4.KeystoresPropertiesCreateIbmCloudKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateAzure); ok {
		return dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAzureToMap(model.(*ukov4.KeystoresPropertiesCreateAzure))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreate); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeystoresPropertiesCreate)
		if model.Group != nil {
			modelMap["group"] = *model.Group
		}
		if model.Type != nil {
			modelMap["type"] = *model.Type
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
		return nil, fmt.Errorf("Unrecognized ukov4.KeystoresPropertiesCreateIntf subtype encountered")
	}
}

func dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateGoogleKmsToMap(model *ukov4.KeystoresPropertiesCreateGoogleKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = *model.Group
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
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

func dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAwsKmsToMap(model *ukov4.KeystoresPropertiesCreateAwsKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = *model.Group
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateIbmCloudKmsToMap(model *ukov4.KeystoresPropertiesCreateIbmCloudKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = *model.Group
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func dataSourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAzureToMap(model *ukov4.KeystoresPropertiesCreateAzure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = *model.Group
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}
