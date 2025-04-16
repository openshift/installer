// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package hpcs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-hpcs-uko-sdk/ukov4"
)

func ResourceIbmKeyTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIbmKeyTemplateCreate,
		ReadContext:   ResourceIbmKeyTemplateRead,
		UpdateContext: ResourceIbmKeyTemplateUpdate,
		DeleteContext: ResourceIbmKeyTemplateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the UKO instance this resource exists in.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region of the UKO instance this resource exists in.",
			},
			"uko_vault": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the Vault in which the update is to take place.",
			},
			"vault": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "ID of the Vault where the entity is to be created in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_key_template", "name"),
				Description:  "Name of the template, it will be referenced when creating managed keys.",
			},
			"key": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Properties describing the properties of the managed key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The size of the underlying cryptographic key or key pair. E.g. \"256\" for AES keys, or \"2048\" for RSA.",
						},
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The algorithm of the key.",
						},
						"activation_date": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key activation date can be provided as a period definition (e.g. PY1 means 1 year).",
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key expiration date can be provided as a period definition (e.g. PY1 means 1 year).",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The state that the key will be in after generation.",
						},
					},
				},
			},
			"keystores": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "An array describing the type and group of target keystores the managed key is to be installed in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Which keystore group to distribute the key to.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of keystore.",
						},
						"google_key_protection_level": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"google_key_purpose": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"google_kms_algorithm": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_key_template", "description"),
				Description:  "Description of the key template.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the key template. Every time the key template is updated, the version will be updated automatically.",
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
		},
	}
}

func ResourceIbmKeyTemplateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Za-z][A-Za-z0-9-]*$`,
			MinValueLength:             1,
			MaxValueLength:             30,
		},
		validate.ValidateSchema{
			Identifier:                 "region",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "au-syd, in-che, jp-osa, jp-tok, kr-seo, eu-de, eu-gb, ca-tor, us-south, us-south-test, us-east, br-sao",
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.|\\n)*`,
			MinValueLength:             0,
			MaxValueLength:             200,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_hpcs_key_template", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIbmKeyTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	instance_id := d.Get("instance_id").(string)
	region := d.Get("region").(string)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	createKeyTemplateOptions := &ukov4.CreateKeyTemplateOptions{}

	createKeyTemplateOptions.SetUKOVault(d.Get("uko_vault").(string))
	vaultModel, err := ResourceIbmKeyTemplateMapToVaultReferenceInCreationRequest(d.Get("vault.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createKeyTemplateOptions.SetVault(vaultModel)
	createKeyTemplateOptions.SetName(d.Get("name").(string))
	keyModel, err := ResourceIbmKeyTemplateMapToKeyProperties(d.Get("key.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createKeyTemplateOptions.SetKey(keyModel)
	var keystores []ukov4.KeystoresPropertiesCreateIntf
	for _, e := range d.Get("keystores").([]interface{}) {
		value := e.(map[string]interface{})
		keystoresItem, err := resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreate(value)
		if err != nil {
			return diag.FromErr(err)
		}
		keystores = append(keystores, keystoresItem)
	}
	createKeyTemplateOptions.SetKeystores(keystores)
	if _, ok := d.GetOk("description"); ok {
		createKeyTemplateOptions.SetDescription(d.Get("description").(string))
	}

	template, response, err := ukoClient.CreateKeyTemplateWithContext(context, createKeyTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateKeyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateKeyTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, d.Get("uko_vault").(string), *template.ID))

	return ResourceIbmKeyTemplateRead(context, d, meta)
}

func ResourceIbmKeyTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getKeyTemplateOptions := &ukov4.GetKeyTemplateOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	template_id := id[3]
	getKeyTemplateOptions.SetID(template_id)
	getKeyTemplateOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	template, response, err := ukoClient.GetKeyTemplateWithContext(context, getKeyTemplateOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetKeyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetKeyTemplateWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("uko_vault", getKeyTemplateOptions.UKOVault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uko_vault: %s", err))
	}
	vaultMap, err := ResourceIbmKeyTemplateVaultReferenceInCreationRequestToMap(template.Vault)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vault", []map[string]interface{}{vaultMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault: %s", err))
	}
	if err = d.Set("name", template.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	keyMap, err := ResourceIbmKeyTemplateKeyPropertiesToMap(template.Key)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("key", []map[string]interface{}{keyMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key: %s", err))
	}
	keystores := []map[string]interface{}{}
	for _, keystoresItem := range template.Keystores {
		keystoresItemMap, err := resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateToMap(keystoresItem)
		if err != nil {
			return diag.FromErr(err)
		}
		keystores = append(keystores, keystoresItemMap)
	}
	if err = d.Set("keystores", keystores); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting keystores: %s", err))
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
	if err = d.Set("href", template.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	// TODO: I'm worried about this line
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}
	// if err = d.Set("version", flex.IntValue(template.Version)); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	// }

	return nil
}

func ResourceIbmKeyTemplateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	updateKeyTemplateOptions := &ukov4.UpdateKeyTemplateOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	template_id := id[3]
	updateKeyTemplateOptions.SetID(template_id)
	updateKeyTemplateOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	hasChange := false

	// TODO: Worried about this
	// if d.HasChange("key") || d.HasChange("keystores") {
	// 	keyprops, err := ResourceIbmKeyTemplateMapToKeyProperties(d.Get("key.0").(map[string]interface{}))
	if d.HasChange("uko_vault") || d.HasChange("vault") {
		updateKeyTemplateOptions.SetUKOVault(d.Get("uko_vault").(string))
		// vault, err := ResourceIbmKeyTemplateMapToVaultReferenceInCreationRequest(d.Get("vault.0").(map[string]interface{}))
		// if err != nil {
		// 	return diag.FromErr(err)
		// }
		// updateKeyTemplateOptions.SetUKOVault(vault)
		//
	}
	// if d.HasChange("name") {
	// 	updateKeyTemplateOptions.SetName(d.Get("name").(string))
	// }
	if d.HasChange("key") || d.HasChange("keystores") {
		keyprops, err := ResourceIbmKeyTemplateMapToKeyProperties(d.Get("key").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		var key *ukov4.KeyPropertiesUpdate
		key.Size = keyprops.Size
		key.ActivationDate = keyprops.ActivationDate
		key.ExpirationDate = keyprops.ExpirationDate
		key.State = keyprops.State
		updateKeyTemplateOptions.SetKey(key)
		// TODO: handle Keystores of type TypeList -- not primitive, not model
		hasChange = true
	}
	if d.HasChange("description") {
		updateKeyTemplateOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	// Etag support
	updateKeyTemplateOptions.SetIfMatch(d.Get("version").(string))

	if hasChange {
		_, response, err := ukoClient.UpdateKeyTemplateWithContext(context, updateKeyTemplateOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateKeyTemplateWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateKeyTemplateWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIbmKeyTemplateRead(context, d, meta)
}

func ResourceIbmKeyTemplateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteKeyTemplateOptions := &ukov4.DeleteKeyTemplateOptions{}

	// Etag support
	deleteKeyTemplateOptions.SetIfMatch(d.Get("version").(string))

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	template_id := id[3]

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	deleteKeyTemplateOptions.SetID(template_id)
	deleteKeyTemplateOptions.SetUKOVault(vault_id)

	response, err := ukoClient.DeleteKeyTemplateWithContext(context, deleteKeyTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteKeyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteKeyTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIbmKeyTemplateMapToVaultReferenceInCreationRequest(modelMap map[string]interface{}) (*ukov4.VaultReferenceInCreationRequest, error) {
	model := &ukov4.VaultReferenceInCreationRequest{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIbmKeyTemplateMapToKeyProperties(modelMap map[string]interface{}) (*ukov4.KeyProperties, error) {
	model := &ukov4.KeyProperties{}
	model.Size = core.StringPtr(modelMap["size"].(string))
	model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	model.ActivationDate = core.StringPtr(modelMap["activation_date"].(string))
	model.ExpirationDate = core.StringPtr(modelMap["expiration_date"].(string))
	model.State = core.StringPtr(modelMap["state"].(string))
	return model, nil
}

// TODO: worried about this
//
//	func ResourceIbmKeyTemplateMapToKeystoresProperties(modelMap map[string]interface{}) (*ukov4.KeystoresProperties, error) {
//		model := &ukov4.KeystoresProperties{}
//		model.Group = core.StringPtr(modelMap["group"].(string))
//		model.Type = core.StringPtr(modelMap["type"].(string))
func resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreate(modelMap map[string]interface{}) (ukov4.KeystoresPropertiesCreateIntf, error) {
	model := &ukov4.KeystoresPropertiesCreate{}
	if modelMap["group"] != nil && modelMap["group"].(string) != "" {
		model.Group = core.StringPtr(modelMap["group"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["google_key_protection_level"] != nil && modelMap["google_key_protection_level"].(string) != "" {
		model.GoogleKeyProtectionLevel = core.StringPtr(modelMap["google_key_protection_level"].(string))
	}
	if modelMap["google_key_purpose"] != nil && modelMap["google_key_purpose"].(string) != "" {
		model.GoogleKeyPurpose = core.StringPtr(modelMap["google_key_purpose"].(string))
	}
	if modelMap["google_kms_algorithm"] != nil && modelMap["google_kms_algorithm"].(string) != "" {
		model.GoogleKmsAlgorithm = core.StringPtr(modelMap["google_kms_algorithm"].(string))
	}
	return model, nil
}

func resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreateGoogleKms(modelMap map[string]interface{}) (*ukov4.KeystoresPropertiesCreateGoogleKms, error) {
	model := &ukov4.KeystoresPropertiesCreateGoogleKms{}
	if modelMap["group"] != nil && modelMap["group"].(string) != "" {
		model.Group = core.StringPtr(modelMap["group"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["google_key_protection_level"] != nil && modelMap["google_key_protection_level"].(string) != "" {
		model.GoogleKeyProtectionLevel = core.StringPtr(modelMap["google_key_protection_level"].(string))
	}
	if modelMap["google_key_purpose"] != nil && modelMap["google_key_purpose"].(string) != "" {
		model.GoogleKeyPurpose = core.StringPtr(modelMap["google_key_purpose"].(string))
	}
	if modelMap["google_kms_algorithm"] != nil && modelMap["google_kms_algorithm"].(string) != "" {
		model.GoogleKmsAlgorithm = core.StringPtr(modelMap["google_kms_algorithm"].(string))
	}
	return model, nil
}

func resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreateAwsKms(modelMap map[string]interface{}) (*ukov4.KeystoresPropertiesCreateAwsKms, error) {
	model := &ukov4.KeystoresPropertiesCreateAwsKms{}
	if modelMap["group"] != nil && modelMap["group"].(string) != "" {
		model.Group = core.StringPtr(modelMap["group"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreateIbmCloudKms(modelMap map[string]interface{}) (*ukov4.KeystoresPropertiesCreateIbmCloudKms, error) {
	model := &ukov4.KeystoresPropertiesCreateIbmCloudKms{}
	if modelMap["group"] != nil && modelMap["group"].(string) != "" {
		model.Group = core.StringPtr(modelMap["group"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIbmHpcsKeyTemplateMapToKeystoresPropertiesCreateAzure(modelMap map[string]interface{}) (*ukov4.KeystoresPropertiesCreateAzure, error) {
	model := &ukov4.KeystoresPropertiesCreateAzure{}
	if modelMap["group"] != nil && modelMap["group"].(string) != "" {
		model.Group = core.StringPtr(modelMap["group"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func ResourceIbmKeyTemplateVaultReferenceInCreationRequestToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIbmKeyTemplateKeyPropertiesToMap(model *ukov4.KeyProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["size"] = model.Size
	modelMap["algorithm"] = model.Algorithm
	modelMap["activation_date"] = model.ActivationDate
	modelMap["expiration_date"] = model.ExpirationDate
	modelMap["state"] = model.State
	return modelMap, nil
}

// TODO: Worried
// func ResourceIbmKeyTemplateKeystoresPropertiesToMap(model *ukov4.KeystoresProperties) (map[string]interface{}, error) {
func resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateToMap(model ukov4.KeystoresPropertiesCreateIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeystoresPropertiesCreateGoogleKms); ok {
		return resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateGoogleKmsToMap(model.(*ukov4.KeystoresPropertiesCreateGoogleKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateAwsKms); ok {
		return resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAwsKmsToMap(model.(*ukov4.KeystoresPropertiesCreateAwsKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateIbmCloudKms); ok {
		return resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateIbmCloudKmsToMap(model.(*ukov4.KeystoresPropertiesCreateIbmCloudKms))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreateAzure); ok {
		return resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAzureToMap(model.(*ukov4.KeystoresPropertiesCreateAzure))
	} else if _, ok := model.(*ukov4.KeystoresPropertiesCreate); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeystoresPropertiesCreate)
		if model.Group != nil {
			modelMap["group"] = model.Group
		}
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.GoogleKeyProtectionLevel != nil {
			modelMap["google_key_protection_level"] = model.GoogleKeyProtectionLevel
		}
		if model.GoogleKeyPurpose != nil {
			modelMap["google_key_purpose"] = model.GoogleKeyPurpose
		}
		if model.GoogleKmsAlgorithm != nil {
			modelMap["google_kms_algorithm"] = model.GoogleKmsAlgorithm
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ukov4.KeystoresPropertiesCreateIntf subtype encountered")
	}
}

func resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateGoogleKmsToMap(model *ukov4.KeystoresPropertiesCreateGoogleKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = model.Group
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.GoogleKeyProtectionLevel != nil {
		modelMap["google_key_protection_level"] = model.GoogleKeyProtectionLevel
	}
	if model.GoogleKeyPurpose != nil {
		modelMap["google_key_purpose"] = model.GoogleKeyPurpose
	}
	if model.GoogleKmsAlgorithm != nil {
		modelMap["google_kms_algorithm"] = model.GoogleKmsAlgorithm
	}
	return modelMap, nil
}

func resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAwsKmsToMap(model *ukov4.KeystoresPropertiesCreateAwsKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = model.Group
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateIbmCloudKmsToMap(model *ukov4.KeystoresPropertiesCreateIbmCloudKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = model.Group
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIbmHpcsKeyTemplateKeystoresPropertiesCreateAzureToMap(model *ukov4.KeystoresPropertiesCreateAzure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Group != nil {
		modelMap["group"] = model.Group
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}
