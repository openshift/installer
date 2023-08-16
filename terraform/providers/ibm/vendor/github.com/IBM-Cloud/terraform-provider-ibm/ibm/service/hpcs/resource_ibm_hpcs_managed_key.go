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

func ResourceIbmManagedKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIbmManagedKeyCreate,
		ReadContext:   ResourceIbmManagedKeyRead,
		UpdateContext: ResourceIbmManagedKeyUpdate,
		DeleteContext: ResourceIbmManagedKeyDelete,
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
			"key_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the key.",
			},
			"template_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_managed_key", "template_name"),
				Description:  "Name of the key template to use when creating a key.",
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
			"label": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_managed_key", "label"),
				Description:  "The label of the key.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Key-value pairs associated with the key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of a tag.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of a tag.",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_managed_key", "description"),
				Description:  "Description of the managed key.",
			},
			"template": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Reference to a key template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the key template.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A URL that uniquely identifies your cloud resource.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
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
							Required:    true,
							Description: "The method used for calculating the verification pattern.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the target keystore.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of keystore.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Required:    true,
							Description: "The v4 UUID used to uniquely identify the resource, as specified by RFC 4122.",
						},
						"label_in_keystore": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The label of the key.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Type of the key instance.",
						},
						"keystore": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Description of properties of a key within the context of keystores.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of keystore.",
									},
								},
							},
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmManagedKeyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "template_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Za-z][A-Za-z0-9-]+$`,
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
			Identifier:                 "label",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Za-z0-9._ \/-]+$`,
			MinValueLength:             1,
			MaxValueLength:             255,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.|\n)*`,
			// TODO: the old version was:
			// Regexp:                     `(.|\\n)*`,
			MinValueLength: 0,
			MaxValueLength: 200,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_hpcs_managed_key", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIbmManagedKeyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	createManagedKeyOptions := &ukov4.CreateManagedKeyOptions{}

	uko_vault := d.Get("uko_vault").(string)
	createManagedKeyOptions.SetUKOVault(uko_vault)
	createManagedKeyOptions.SetTemplateName(d.Get("template_name").(string))
	vaultModel, err := ResourceIbmManagedKeyMapToVaultReferenceInCreationRequest(d.Get("vault.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createManagedKeyOptions.SetVault(vaultModel)
	createManagedKeyOptions.SetLabel(d.Get("label").(string))
	if _, ok := d.GetOk("tags"); ok {
		var tags []ukov4.Tag
		for _, e := range d.Get("tags").([]interface{}) {
			value := e.(map[string]interface{})
			tagsItem, err := ResourceIbmManagedKeyMapToTag(value)
			if err != nil {
				return diag.FromErr(err)
			}
			tags = append(tags, *tagsItem)
		}
		createManagedKeyOptions.SetTags(tags)
	}
	if _, ok := d.GetOk("description"); ok {
		createManagedKeyOptions.SetDescription(d.Get("description").(string))
	}

	managedKey, response, err := ukoClient.CreateManagedKeyWithContext(context, createManagedKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateManagedKeyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateManagedKeyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, uko_vault, *managedKey.ID))

	return ResourceIbmManagedKeyRead(context, d, meta)
}

func ResourceIbmManagedKeyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getManagedKeyOptions := &ukov4.GetManagedKeyOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	key_id := id[3]
	getManagedKeyOptions.SetID(key_id)
	getManagedKeyOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	managedKey, response, err := ukoClient.GetManagedKeyWithContext(context, getManagedKeyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetManagedKeyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetManagedKeyWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("uko_vault", getManagedKeyOptions.UKOVault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uko_vault: %s", err))
	}
	if err = d.Set("key_id", key_id); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_id: %s", err))
	}
	// if err = d.Set("template_name", getManagedKeyOptions.TemplateName); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting template_name: %s", err))
	// }
	vaultMap, err := ResourceIbmManagedKeyVaultReferenceInCreationRequestToMap(managedKey.Vault)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vault", []map[string]interface{}{vaultMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault: %s", err))
	}
	if err = d.Set("label", managedKey.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}
	tags := []map[string]interface{}{}
	if managedKey.Tags != nil {
		for _, tagsItem := range managedKey.Tags {
			tagsItemMap, err := ResourceIbmManagedKeyTagToMap(&tagsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			tags = append(tags, tagsItemMap)
		}
	}
	if err = d.Set("tags", tags); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tags: %s", err))
	}
	if err = d.Set("description", managedKey.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if managedKey.Template != nil {
		templateMap, err := ResourceIbmManagedKeyTemplateReferenceToMap(managedKey.Template)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("template", []map[string]interface{}{templateMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting template: %s", err))
		}
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
		for _, verificationPatternsItem := range managedKey.VerificationPatterns {
			verificationPatternsItemMap, err := ResourceIbmManagedKeyKeyVerificationPatternToMap(&verificationPatternsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			verificationPatterns = append(verificationPatterns, verificationPatternsItemMap)
		}
	}
	if err = d.Set("verification_patterns", verificationPatterns); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting verification_patterns: %s", err))
	}
	if err = d.Set("activation_date", flex.DateToString(managedKey.ActivationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting activation_date: %s", err))
	}
	if err = d.Set("expiration_date", flex.DateToString(managedKey.ExpirationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration_date: %s", err))
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
	for _, referencedKeystoresItem := range managedKey.ReferencedKeystores {
		referencedKeystoresItemMap, err := ResourceIbmManagedKeyTargetKeystoreReferenceToMap(&referencedKeystoresItem)
		if err != nil {
			return diag.FromErr(err)
		}
		referencedKeystores = append(referencedKeystores, referencedKeystoresItemMap)
	}
	if err = d.Set("referenced_keystores", referencedKeystores); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referenced_keystores: %s", err))
	}
	instances := []map[string]interface{}{}
	for _, instancesItem := range managedKey.Instances {
		// TODO: Worried about this typing
		instancesItemMap, err := ResourceIbmManagedKeyKeyInstanceToMap(instancesItem)
		if err != nil {
			return diag.FromErr(err)
		}
		instances = append(instances, instancesItemMap)
	}
	if err = d.Set("instances", instances); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instances: %s", err))
	}
	if err = d.Set("href", managedKey.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func ResourceIbmManagedKeyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	updateManagedKeyOptions := &ukov4.UpdateManagedKeyOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	key_id := id[3]
	updateManagedKeyOptions.SetID(key_id)
	updateManagedKeyOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	hasChange := false

	if d.HasChange("uko_vault") || d.HasChange("template_name") || d.HasChange("vault") || d.HasChange("label") {
		updateManagedKeyOptions.SetUKOVault(d.Get("uko_vault").(string))
		// updateManagedKeyOptions.SetTemplateName(d.Get("template_name").(string))
		// vault, err := ResourceIbmManagedKeyMapToVaultReferenceInCreationRequest(d.Get("vault.0").(map[string]interface{}))
		// if err != nil {
		// 	return diag.FromErr(err)
		// }
		// updateManagedKeyOptions.SetUKOVault(vault)
		updateManagedKeyOptions.SetLabel(d.Get("label").(string))
		hasChange = true
	}
	if d.HasChange("tags") {
		// TODO: handle Tags of type TypeList -- not primitive, not model
		hasChange = true
	}
	if d.HasChange("description") {
		updateManagedKeyOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	etag := d.Get("etag").(string)
	updateManagedKeyOptions.SetIfMatch(etag)

	if hasChange {
		_, response, err := ukoClient.UpdateManagedKeyWithContext(context, updateManagedKeyOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateManagedKeyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateManagedKeyWithContext failed %s\n%s", err, response))
		}
		etag = response.Headers.Get("Etag")
	}

	// Support for changing state
	if d.HasChange("state") {
		prevIntf, newIntf := d.GetChange("state")
		prev := prevIntf.(string)
		new := newIntf.(string)
		if new == "deactivated" {
			if prev == "active" {
				deactivateManagedKeyOptions := &ukov4.DeactivateManagedKeyOptions{}

				deactivateManagedKeyOptions.SetIfMatch(etag)
				deactivateManagedKeyOptions.SetID(key_id)
				deactivateManagedKeyOptions.SetUKOVault(vault_id)

				_, response, err := ukoClient.DeactivateManagedKeyWithContext(context, deactivateManagedKeyOptions)
				if err != nil {
					log.Printf("[DEBUG] DeactivateManagedKeyWithContext failed %s\n%s", err, response)
					return diag.FromErr(fmt.Errorf("DeactivateManagedKeyWithContext failed %s\n%s", err, response))
				}
			} else {
				return diag.FromErr(fmt.Errorf("Deactivate managed key failed: Cannot deactivate key not in active state"))
			}
		} else if new == "destroyed" {
			if prev == "deactivated" || prev == "pre_activation" {
				destroyManagedKeyOptions := &ukov4.DestroyManagedKeyOptions{}

				destroyManagedKeyOptions.SetIfMatch(etag)
				destroyManagedKeyOptions.SetID(key_id)
				destroyManagedKeyOptions.SetUKOVault(vault_id)

				_, response, err := ukoClient.DestroyManagedKeyWithContext(context, destroyManagedKeyOptions)
				if err != nil {
					log.Printf("[DEBUG] DestroyManagedKeyWithContext failed %s\n%s", err, response)
					return diag.FromErr(fmt.Errorf("DestroyManagedKeyWithContext failed %s\n%s", err, response))
				}
			} else {
				return diag.FromErr(fmt.Errorf("Destroy managed key failed: Cannot destroy key not in deactivated state"))
			}
		} else if new == "active" {
			if prev == "deactivated" || prev == "pre_activation" {
				activateManagedKeyOptions := &ukov4.ActivateManagedKeyOptions{}

				activateManagedKeyOptions.SetIfMatch(etag)
				activateManagedKeyOptions.SetID(key_id)
				activateManagedKeyOptions.SetUKOVault(vault_id)

				_, response, err := ukoClient.ActivateManagedKeyWithContext(context, activateManagedKeyOptions)
				if err != nil {
					log.Printf("[DEBUG] ActivateManagedKeyWithContext failed %s\n%s", err, response)
					return diag.FromErr(fmt.Errorf("ActivateManagedKeyWithContext failed %s\n%s", err, response))
				}
			} else {
				return diag.FromErr(fmt.Errorf("Activate managed key failed: Cannot activate key not in deactivated or pre_activation state"))
			}
		}
	}

	return ResourceIbmManagedKeyRead(context, d, meta)
}

func ResourceIbmManagedKeyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: This is where we need to walk through for deleting keys
	// Don't allow deleting of non-destroyed keys
	if d.Get("state") != "destroyed" {
		return diag.FromErr(fmt.Errorf("Delete Key failed: Cannot delete key that is not destroyed"))
	}

	// Etag support
	etag := d.Get("etag").(string)

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	key_id := id[3]

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	deleteManagedKeyOptions := &ukov4.DeleteManagedKeyOptions{}

	deleteManagedKeyOptions.SetIfMatch(etag)
	deleteManagedKeyOptions.SetID(key_id)
	deleteManagedKeyOptions.SetUKOVault(vault_id)

	response, err := ukoClient.DeleteManagedKeyWithContext(context, deleteManagedKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteManagedKeyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteManagedKeyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIbmManagedKeyMapToVaultReferenceInCreationRequest(modelMap map[string]interface{}) (*ukov4.VaultReferenceInCreationRequest, error) {
	model := &ukov4.VaultReferenceInCreationRequest{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIbmManagedKeyMapToTag(modelMap map[string]interface{}) (*ukov4.Tag, error) {
	model := &ukov4.Tag{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIbmManagedKeyVaultReferenceInCreationRequestToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIbmManagedKeyTagToMap(model *ukov4.Tag) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	return modelMap, nil
}

func ResourceIbmManagedKeyTemplateReferenceToMap(model *ukov4.TemplateReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIbmManagedKeyKeyVerificationPatternToMap(model *ukov4.KeyVerificationPattern) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["method"] = model.Method
	modelMap["value"] = model.Value
	return modelMap, nil
}

func ResourceIbmManagedKeyTargetKeystoreReferenceToMap(model *ukov4.TargetKeystoreReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	modelMap["type"] = model.Type
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

// TODO: Worried about typing
func ResourceIbmManagedKeyKeyInstanceToMap(model ukov4.KeyInstanceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeyInstanceGoogleKms); ok {
		return resourceIbmHpcsManagedKeyKeyInstanceGoogleKmsToMap(model.(*ukov4.KeyInstanceGoogleKms))
	} else if _, ok := model.(*ukov4.KeyInstanceAwsKms); ok {
		return resourceIbmHpcsManagedKeyKeyInstanceAwsKmsToMap(model.(*ukov4.KeyInstanceAwsKms))
	} else if _, ok := model.(*ukov4.KeyInstanceIbmCloudKms); ok {
		return resourceIbmHpcsManagedKeyKeyInstanceIbmCloudKmsToMap(model.(*ukov4.KeyInstanceIbmCloudKms))
	} else if _, ok := model.(*ukov4.KeyInstanceAzure); ok {
		return resourceIbmHpcsManagedKeyKeyInstanceAzureToMap(model.(*ukov4.KeyInstanceAzure))
	} else if _, ok := model.(*ukov4.KeyInstance); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeyInstance)
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.LabelInKeystore != nil {
			modelMap["label_in_keystore"] = model.LabelInKeystore
		}
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.Keystore != nil {
			keystoreMap, err := resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
			if err != nil {
				return modelMap, err
			}
			modelMap["keystore"] = []map[string]interface{}{keystoreMap}
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
		return nil, fmt.Errorf("Unrecognized ukov4.KeyInstanceIntf subtype encountered")
	}
}

func resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model *ukov4.InstanceInKeystore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["group"] = model.Group
	modelMap["type"] = model.Type
	return modelMap, nil
}

func resourceIbmHpcsManagedKeyKeyInstanceGoogleKmsToMap(model *ukov4.KeyInstanceGoogleKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["label_in_keystore"] = model.LabelInKeystore
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	keystoreMap, err := resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
	if err != nil {
		return modelMap, err
	}
	modelMap["keystore"] = []map[string]interface{}{keystoreMap}
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

func resourceIbmHpcsManagedKeyKeyInstanceAwsKmsToMap(model *ukov4.KeyInstanceAwsKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["label_in_keystore"] = model.LabelInKeystore
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	keystoreMap, err := resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
	if err != nil {
		return modelMap, err
	}
	modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	return modelMap, nil
}

func resourceIbmHpcsManagedKeyKeyInstanceIbmCloudKmsToMap(model *ukov4.KeyInstanceIbmCloudKms) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["label_in_keystore"] = model.LabelInKeystore
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	keystoreMap, err := resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
	if err != nil {
		return modelMap, err
	}
	modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	return modelMap, nil
}

func resourceIbmHpcsManagedKeyKeyInstanceAzureToMap(model *ukov4.KeyInstanceAzure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["label_in_keystore"] = model.LabelInKeystore
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	keystoreMap, err := resourceIbmHpcsManagedKeyInstanceInKeystoreToMap(model.Keystore)
	if err != nil {
		return modelMap, err
	}
	modelMap["keystore"] = []map[string]interface{}{keystoreMap}
	return modelMap, nil
}
