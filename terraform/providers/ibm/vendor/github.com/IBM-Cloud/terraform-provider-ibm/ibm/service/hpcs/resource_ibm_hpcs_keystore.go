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

func ResourceIbmKeystore() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIbmKeystoreCreate,
		ReadContext:   ResourceIbmKeystoreRead,
		UpdateContext: ResourceIbmKeystoreUpdate,
		DeleteContext: ResourceIbmKeystoreDelete,
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
			"aws_region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Region.",
			},
			"aws_access_key_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The access key id used for connecting to this instance of AWS KMS.",
			},
			"aws_secret_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The secret access key used for connecting to this instance of AWS KMS.",
			},
			"azure_service_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service name of the key vault instance from the Azure portal.",
			},
			"azure_resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group in Azure.",
			},
			"azure_location": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the Azure Key Vault.",
			},
			"azure_service_principal_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure service principal client ID.",
			},
			"azure_service_principal_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure service principal password.",
			},
			"azure_tenant": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure tenant that the Key Vault is associated with,.",
			},
			"azure_subscription_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscription ID in Azure.",
			},
			"azure_environment": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure environment, usually 'Azure'.",
			},
			"ibm_variant": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Possible IBM Cloud KMS variants.",
			},
			"ibm_api_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API endpoint of the IBM Cloud keystore.",
			},
			"ibm_iam_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Endpoint of the IAM service for this IBM Cloud keystore.",
			},
			"ibm_api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.",
			},
			"ibm_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The instance ID of the IBM Cloud keystore.",
			},
			"ibm_key_ring": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key ring of an IBM Cloud KMS Keystore.",
			},
			"vault": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Reference to a vault.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the target keystore. It can be changed in the future.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the keystore.",
			},
			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of groups that this keystore belongs to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of keystore.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the target keystore was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the target keystore was last updated.",
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmKeystoreValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "region",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "au-syd, in-che, jp-osa, jp-tok, kr-seo, eu-de, eu-gb, ca-tor, us-south, us-south-test, us-east, br-sao",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_hpcs_keystore", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIbmKeystoreCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	createKeystoreOptions := &ukov4.CreateKeystoreOptions{}

	createKeystoreOptions.SetUKOVault(d.Get("uko_vault").(string))

	// Instead of setting keystore body this way, we need to get every value of d and create a keystore body with those parameters
	keystoreBodyModel, err := ResourceIbmKeystoreMapToKeystoreCreationRequest(DKeystoreToKeystoreBody(d).(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createKeystoreOptions.SetKeystoreBody(keystoreBodyModel)

	keystoreIntf, response, err := ukoClient.CreateKeystoreWithContext(context, createKeystoreOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateKeystoreWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateKeystoreWithContext failed %s\n%s", err, response))
	}

	keystore := keystoreIntf.(*ukov4.Keystore)
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, d.Get("uko_vault").(string), *keystore.ID))

	return ResourceIbmKeystoreRead(context, d, meta)
}

func ResourceIbmKeystoreRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getKeystoreOptions := &ukov4.GetKeystoreOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	keystore_id := id[3]
	getKeystoreOptions.SetID(keystore_id)
	getKeystoreOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	keystoreIntf, response, err := ukoClient.GetKeystoreWithContext(context, getKeystoreOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetKeystoreWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetKeystoreWithContext failed %s\n%s", err, response))
	}

	keystore := keystoreIntf.(*ukov4.Keystore)
	if err = d.Set("uko_vault", getKeystoreOptions.UKOVault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uko_vault: %s", err))
	}
	// TODO: handle argument of type KeystoreCreationRequest
	if keystore.Vault != nil {
		vaultMap, err := ResourceIbmKeystoreVaultReferenceToMap(keystore.Vault)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("vault", []map[string]interface{}{vaultMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vault: %s", err))
		}
	}
	if err = d.Set("name", keystore.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", keystore.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if keystore.Groups != nil {
		if err = d.Set("groups", keystore.Groups); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting groups: %s", err))
		}
	}
	if err = d.Set("type", keystore.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(keystore.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(keystore.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created_by", keystore.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_by", keystore.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}
	if err = d.Set("href", keystore.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("aws_region", keystore.AwsRegion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting aws_region: %s", err))
	}
	if err = d.Set("aws_access_key_id", keystore.AwsAccessKeyID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting aws_access_key_id: %s", err))
	}
	if err = d.Set("aws_secret_access_key", keystore.AwsSecretAccessKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting aws_secret_access_key: %s", err))
	}
	if err = d.Set("azure_service_name", keystore.AzureServiceName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_service_name: %s", err))
	}
	if err = d.Set("azure_resource_group", keystore.AzureResourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_resource_group: %s", err))
	}
	if err = d.Set("azure_location", keystore.AzureLocation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_location: %s", err))
	}
	if err = d.Set("azure_service_principal_client_id", keystore.AzureServicePrincipalClientID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_service_principal_client_id: %s", err))
	}
	if err = d.Set("azure_service_principal_password", keystore.AzureServicePrincipalPassword); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_service_principal_password: %s", err))
	}
	if err = d.Set("azure_tenant", keystore.AzureTenant); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_tenant: %s", err))
	}
	if err = d.Set("azure_subscription_id", keystore.AzureSubscriptionID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_subscription_id: %s", err))
	}
	if err = d.Set("azure_environment", keystore.AzureEnvironment); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting azure_environment: %s", err))
	}
	if err = d.Set("ibm_api_endpoint", keystore.IbmApiEndpoint); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_api_endpoint: %s", err))
	}
	if err = d.Set("ibm_iam_endpoint", keystore.IbmIamEndpoint); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_iam_endpoint: %s", err))
	}
	if err = d.Set("ibm_api_key", keystore.IbmApiKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_api_key: %s", err)) // pragma: allowlist secret
	}
	if err = d.Set("ibm_instance_id", keystore.IbmInstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_instance_id: %s", err))
	}
	if err = d.Set("ibm_variant", keystore.IbmVariant); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_variant: %s", err))
	}
	if err = d.Set("ibm_key_ring", keystore.IbmKeyRing); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ibm_key_ring: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func ResourceIbmKeystoreUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	updateKeystoreOptions := &ukov4.UpdateKeystoreOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	keystore_id := id[3]
	updateKeystoreOptions.SetID(keystore_id)
	updateKeystoreOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	hasChange := false

	// if d.HasChange() for every parameter, then again, create a keystore body out of the new parameters

	if d.HasChange("uko_vault") || DHasChanges(d) {
		updateKeystoreOptions.SetUKOVault(d.Get("uko_vault").(string))
		keystoreBody := ukov4.KeystoreUpdateRequestIntf(DKeystoreToKeystoreBodyUpdate(d))
		updateKeystoreOptions.SetKeystoreBody(keystoreBody)
		hasChange = true
	}

	updateKeystoreOptions.SetIfMatch(d.Get("version").(string))

	if hasChange {
		_, response, err := ukoClient.UpdateKeystoreWithContext(context, updateKeystoreOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateKeystoreWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateKeystoreWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIbmKeystoreRead(context, d, meta)
}

func ResourceIbmKeystoreDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteKeystoreOptions := &ukov4.DeleteKeystoreOptions{}

	// Etag support
	deleteKeystoreOptions.SetIfMatch(d.Get("version").(string))

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	keystore_id := id[3]

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	deleteKeystoreOptions.SetID(keystore_id)
	deleteKeystoreOptions.SetUKOVault(vault_id)

	response, err := ukoClient.DeleteKeystoreWithContext(context, deleteKeystoreOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteKeystoreWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteKeystoreWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequest(modelMap map[string]interface{}) (ukov4.KeystoreCreationRequestIntf, error) {
	discValue, ok := modelMap["type"]
	if ok {
		if discValue == "aws_kms" {
			return ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeAwsKmsCreate(modelMap)
		} else if discValue == "azure_key_vault" {
			return ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeAzureCreate(modelMap)
		} else if discValue == "ibm_cloud_kms" {
			return ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate(modelMap)
		} else {
			return nil, fmt.Errorf("unexpected value for discriminator property 'type' found in map: '%s'", discValue)
		}
	} else {
		return nil, fmt.Errorf("discriminator property 'type' not found in map")
	}
}

func ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap map[string]interface{}) (*ukov4.VaultReferenceInCreationRequest, error) {
	model := &ukov4.VaultReferenceInCreationRequest{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeAwsKmsCreate(modelMap map[string]interface{}) (*ukov4.KeystoreCreationRequestKeystoreTypeAwsKmsCreate, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeAwsKmsCreate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	model.AwsRegion = core.StringPtr(modelMap["aws_region"].(string))
	model.AwsAccessKeyID = core.StringPtr(modelMap["aws_access_key_id"].(string))
	model.AwsSecretAccessKey = core.StringPtr(modelMap["aws_secret_access_key"].(string))
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate(modelMap map[string]interface{}) (ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateIntf, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	model.IbmVariant = core.StringPtr(modelMap["ibm_variant"].(string))
	if modelMap["name"] != nil {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	if modelMap["ibm_api_endpoint"] != nil {
		model.IbmApiEndpoint = core.StringPtr(modelMap["ibm_api_endpoint"].(string))
	}
	if modelMap["ibm_iam_endpoint"] != nil {
		model.IbmIamEndpoint = core.StringPtr(modelMap["ibm_iam_endpoint"].(string))
	}
	if modelMap["ibm_api_key"] != nil {
		model.IbmApiKey = core.StringPtr(modelMap["ibm_api_key"].(string))
	}
	if modelMap["ibm_instance_id"] != nil {
		model.IbmInstanceID = core.StringPtr(modelMap["ibm_instance_id"].(string))
	}
	if modelMap["ibm_key_ring"] != nil {
		model.IbmKeyRing = core.StringPtr(modelMap["ibm_key_ring"].(string))
	}
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate(modelMap map[string]interface{}) (*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	model.IbmApiEndpoint = core.StringPtr(modelMap["ibm_api_endpoint"].(string))
	model.IbmIamEndpoint = core.StringPtr(modelMap["ibm_iam_endpoint"].(string))
	model.IbmApiKey = core.StringPtr(modelMap["ibm_api_key"].(string))
	model.IbmInstanceID = core.StringPtr(modelMap["ibm_instance_id"].(string))
	model.IbmVariant = core.StringPtr(modelMap["ibm_variant"].(string))
	if modelMap["ibm_key_ring"] != nil {
		model.IbmKeyRing = core.StringPtr(modelMap["ibm_key_ring"].(string))
	}
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate(modelMap map[string]interface{}) (ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateIntf, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	model.IbmVariant = core.StringPtr(modelMap["ibm_variant"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate(modelMap map[string]interface{}) (*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	model.IbmVariant = core.StringPtr(modelMap["ibm_variant"].(string))
	if modelMap["name"] != nil {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	return model, nil
}

func ResourceIbmKeystoreMapToKeystoreCreationRequestKeystoreTypeAzureCreate(modelMap map[string]interface{}) (*ukov4.KeystoreCreationRequestKeystoreTypeAzureCreate, error) {
	model := &ukov4.KeystoreCreationRequestKeystoreTypeAzureCreate{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	VaultModel, err := ResourceIbmKeystoreMapToVaultReferenceInCreationRequest(modelMap["vault"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Vault = VaultModel
	if modelMap["name"] != nil {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["groups"] != nil {
		groups := []string{}
		for _, groupsItem := range modelMap["groups"].([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		model.Groups = groups
	}
	model.AzureServiceName = core.StringPtr(modelMap["azure_service_name"].(string))
	model.AzureResourceGroup = core.StringPtr(modelMap["azure_resource_group"].(string))
	model.AzureLocation = core.StringPtr(modelMap["azure_location"].(string))
	model.AzureServicePrincipalClientID = core.StringPtr(modelMap["azure_service_principal_client_id"].(string))
	model.AzureServicePrincipalPassword = core.StringPtr(modelMap["azure_service_principal_password"].(string))
	model.AzureTenant = core.StringPtr(modelMap["azure_tenant"].(string))
	model.AzureSubscriptionID = core.StringPtr(modelMap["azure_subscription_id"].(string))
	model.AzureEnvironment = core.StringPtr(modelMap["azure_environment"].(string))
	return model, nil
}

func ResourceIbmKeystoreKeystoreCreationRequestToMap(model ukov4.KeystoreCreationRequestIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeAwsKmsCreate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeAwsKmsCreateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeAwsKmsCreate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeAzureCreate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeAzureCreateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeAzureCreate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequest); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeystoreCreationRequest)
		modelMap["type"] = model.Type
		vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
		if err != nil {
			return modelMap, err
		}
		modelMap["vault"] = []map[string]interface{}{vaultMap}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Groups != nil {
			modelMap["groups"] = model.Groups
		}
		if model.AwsRegion != nil {
			modelMap["aws_region"] = model.AwsRegion
		}
		if model.AwsAccessKeyID != nil {
			modelMap["aws_access_key_id"] = model.AwsAccessKeyID
		}
		if model.AwsSecretAccessKey != nil {
			modelMap["aws_secret_access_key"] = model.AwsSecretAccessKey
		}
		if model.AzureServiceName != nil {
			modelMap["azure_service_name"] = model.AzureServiceName
		}
		if model.AzureResourceGroup != nil {
			modelMap["azure_resource_group"] = model.AzureResourceGroup
		}
		if model.AzureLocation != nil {
			modelMap["azure_location"] = model.AzureLocation
		}
		if model.AzureServicePrincipalClientID != nil {
			modelMap["azure_service_principal_client_id"] = model.AzureServicePrincipalClientID
		}
		if model.AzureServicePrincipalPassword != nil {
			modelMap["azure_service_principal_password"] = model.AzureServicePrincipalPassword
		}
		if model.AzureTenant != nil {
			modelMap["azure_tenant"] = model.AzureTenant
		}
		if model.AzureSubscriptionID != nil {
			modelMap["azure_subscription_id"] = model.AzureSubscriptionID
		}
		if model.AzureEnvironment != nil {
			modelMap["azure_environment"] = model.AzureEnvironment
		}
		if model.IbmVariant != nil {
			modelMap["ibm_variant"] = model.IbmVariant
		}
		if model.IbmApiEndpoint != nil {
			modelMap["ibm_api_endpoint"] = model.IbmApiEndpoint
		}
		if model.IbmIamEndpoint != nil {
			modelMap["ibm_iam_endpoint"] = model.IbmIamEndpoint
		}
		if model.IbmApiKey != nil {
			modelMap["ibm_api_key"] = model.IbmApiKey
		}
		if model.IbmInstanceID != nil {
			modelMap["ibm_instance_id"] = model.IbmInstanceID
		}
		if model.IbmKeyRing != nil {
			modelMap["ibm_key_ring"] = model.IbmKeyRing
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ukov4.KeystoreCreationRequestIntf subtype encountered")
	}
}

func ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model *ukov4.VaultReferenceInCreationRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeAwsKmsCreateToMap(model *ukov4.KeystoreCreationRequestKeystoreTypeAwsKmsCreate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
	if err != nil {
		return modelMap, err
	}
	modelMap["vault"] = []map[string]interface{}{vaultMap}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Groups != nil {
		modelMap["groups"] = model.Groups
	}
	modelMap["aws_region"] = model.AwsRegion
	modelMap["aws_access_key_id"] = model.AwsAccessKeyID
	modelMap["aws_secret_access_key"] = model.AwsSecretAccessKey // pragma: allowlist secret
	return modelMap, nil
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateToMap(model ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreate)
		modelMap["type"] = model.Type
		vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
		if err != nil {
			return modelMap, err
		}
		modelMap["vault"] = []map[string]interface{}{vaultMap}
		modelMap["ibm_variant"] = model.IbmVariant
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Groups != nil {
			modelMap["groups"] = model.Groups
		}
		if model.IbmApiEndpoint != nil {
			modelMap["ibm_api_endpoint"] = model.IbmApiEndpoint
		}
		if model.IbmIamEndpoint != nil {
			modelMap["ibm_iam_endpoint"] = model.IbmIamEndpoint
		}
		if model.IbmApiKey != nil {
			modelMap["ibm_api_key"] = model.IbmApiKey
		}
		if model.IbmInstanceID != nil {
			modelMap["ibm_instance_id"] = model.IbmInstanceID
		}
		if model.IbmKeyRing != nil {
			modelMap["ibm_key_ring"] = model.IbmKeyRing
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateIntf subtype encountered")
	}
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreateToMap(model *ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsCreate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
	if err != nil {
		return modelMap, err
	}
	modelMap["vault"] = []map[string]interface{}{vaultMap}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Groups != nil {
		modelMap["groups"] = model.Groups
	}
	modelMap["ibm_api_endpoint"] = model.IbmApiEndpoint
	modelMap["ibm_iam_endpoint"] = model.IbmIamEndpoint
	modelMap["ibm_api_key"] = model.IbmApiKey // pragma: allowlist secret
	modelMap["ibm_instance_id"] = model.IbmInstanceID
	modelMap["ibm_variant"] = model.IbmVariant
	if model.IbmKeyRing != nil {
		modelMap["ibm_key_ring"] = model.IbmKeyRing
	}
	return modelMap, nil
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateToMap(model ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate); ok {
		return ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdateToMap(model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate))
	} else if _, ok := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreate)
		modelMap["type"] = model.Type
		vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
		if err != nil {
			return modelMap, err
		}
		modelMap["vault"] = []map[string]interface{}{vaultMap}
		modelMap["ibm_variant"] = model.IbmVariant
		modelMap["name"] = model.Name
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Groups != nil {
			modelMap["groups"] = model.Groups
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateIntf subtype encountered")
	}
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdateToMap(model *ukov4.KeystoreCreationRequestKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalExternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalCreateKeystoreTypeIbmCloudKmsInternalUpdateKeystoreTypeBaseUpdate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
	if err != nil {
		return modelMap, err
	}
	modelMap["vault"] = []map[string]interface{}{vaultMap}
	modelMap["ibm_variant"] = model.IbmVariant
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Groups != nil {
		modelMap["groups"] = model.Groups
	}
	return modelMap, nil
}

func ResourceIbmKeystoreKeystoreCreationRequestKeystoreTypeAzureCreateToMap(model *ukov4.KeystoreCreationRequestKeystoreTypeAzureCreate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	vaultMap, err := ResourceIbmKeystoreVaultReferenceInCreationRequestToMap(model.Vault)
	if err != nil {
		return modelMap, err
	}
	modelMap["vault"] = []map[string]interface{}{vaultMap}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Groups != nil {
		modelMap["groups"] = model.Groups
	}
	modelMap["azure_service_name"] = model.AzureServiceName
	modelMap["azure_resource_group"] = model.AzureResourceGroup
	modelMap["azure_location"] = model.AzureLocation
	modelMap["azure_service_principal_client_id"] = model.AzureServicePrincipalClientID
	modelMap["azure_service_principal_password"] = model.AzureServicePrincipalPassword // pragma: allowlist secret
	modelMap["azure_tenant"] = model.AzureTenant
	modelMap["azure_subscription_id"] = model.AzureSubscriptionID
	modelMap["azure_environment"] = model.AzureEnvironment
	return modelMap, nil
}

func ResourceIbmKeystoreVaultReferenceToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
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

func DKeystoreToKeystoreBody(d *schema.ResourceData) interface{} {
	keystoreBody := make(map[string]interface{})
	keystoreType := d.Get("type").(string)

	keystoreBody["type"] = keystoreType
	keystoreBody["name"] = d.Get("name").(string)
	keystoreBody["vault"] = d.Get("vault").([]interface{})
	keystoreBody["description"] = d.Get("description").(string)
	keystoreBody["groups"] = d.Get("groups").([]interface{})
	if keystoreType == "aws_kms" {
		keystoreBody["aws_region"] = d.Get("aws_region").(string)
		keystoreBody["aws_access_key_id"] = d.Get("aws_access_key_id").(string)
		keystoreBody["aws_secret_access_key"] = d.Get("aws_secret_access_key").(string)
	} else if keystoreType == "azure_key_vault" {
		keystoreBody["azure_service_name"] = d.Get("azure_service_name").(string)
		keystoreBody["azure_resource_group"] = d.Get("azure_resource_group").(string)
		keystoreBody["azure_location"] = d.Get("azure_location").(string)
		keystoreBody["azure_service_principal_client_id"] = d.Get("azure_service_principal_client_id").(string)
		keystoreBody["azure_service_principal_password"] = d.Get("azure_service_principal_password").(string)
		keystoreBody["azure_tenant"] = d.Get("azure_tenant").(string)
		keystoreBody["azure_subscription_id"] = d.Get("azure_subscription_id").(string)
		keystoreBody["azure_environment"] = d.Get("azure_environment").(string)
	} else if keystoreType == "ibm_cloud_kms" {
		ibm_variant := d.Get("ibm_variant").(string)
		keystoreBody["ibm_variant"] = ibm_variant
		if ibm_variant != "internal" {
			keystoreBody["ibm_api_endpoint"] = d.Get("ibm_api_endpoint").(string)
			keystoreBody["ibm_iam_endpoint"] = d.Get("ibm_iam_endpoint").(string)
			keystoreBody["ibm_api_key"] = d.Get("ibm_api_key").(string)
			keystoreBody["ibm_instance_id"] = d.Get("ibm_instance_id").(string)
			keystoreBody["ibm_key_ring"] = d.Get("ibm_key_ring").(string)
		}
	}

	return keystoreBody
}

func DKeystoreToKeystoreBodyUpdate(d *schema.ResourceData) *ukov4.KeystoreUpdateRequest {
	var keystoreBody ukov4.KeystoreUpdateRequest

	if d.Get("name") != nil && d.Get("name") != "" {
		keystoreBody.Name = core.StringPtr(d.Get("name").(string))
	}
	if d.Get("description") != nil && d.Get("description") != "" {
		keystoreBody.Description = core.StringPtr(d.Get("description").(string))
	}
	if d.Get("groups") != nil {
		groups := []string{}
		for _, groupsItem := range d.Get("groups").([]interface{}) {
			groups = append(groups, groupsItem.(string))
		}
		keystoreBody.Groups = groups
	}
	if d.Get("aws_region") != nil && d.Get("aws_region") != "" {
		keystoreBody.AwsRegion = core.StringPtr(d.Get("aws_region").(string))
	}
	if d.Get("aws_access_key_id") != nil && d.Get("aws_access_key_id") != "" {
		keystoreBody.AwsAccessKeyID = core.StringPtr(d.Get("aws_access_key_id").(string))
	}
	if d.Get("aws_secret_access_key") != nil && d.Get("aws_secret_access_key") != "" {
		keystoreBody.AwsSecretAccessKey = core.StringPtr(d.Get("aws_secret_access_key").(string))
	}
	if d.Get("azure_service_name") != nil && d.Get("azure_service_name") != "" {
		keystoreBody.AzureServiceName = core.StringPtr(d.Get("azure_service_name").(string))
	}
	if d.Get("azure_resource_group") != nil && d.Get("azure_resource_group") != "" {
		keystoreBody.AzureResourceGroup = core.StringPtr(d.Get("azure_resource_group").(string))
	}
	if d.Get("azure_location") != nil && d.Get("azure_location") != "" {
		keystoreBody.AzureLocation = core.StringPtr(d.Get("azure_location").(string))
	}
	if d.Get("azure_service_principal_client_id") != nil && d.Get("azure_service_principal_client_id") != "" {
		keystoreBody.AzureServicePrincipalClientID = core.StringPtr(d.Get("azure_service_principal_client_id").(string))
	}
	if d.Get("azure_service_principal_password") != nil && d.Get("azure_service_principal_password") != "" {
		keystoreBody.AzureServicePrincipalPassword = core.StringPtr(d.Get("azure_service_principal_password").(string))
	}
	if d.Get("azure_tenant") != nil && d.Get("azure_tenant") != "" {
		keystoreBody.AzureTenant = core.StringPtr(d.Get("azure_tenant").(string))
	}
	if d.Get("azure_subscription_id") != nil && d.Get("azure_subscription_id") != "" {
		keystoreBody.AzureSubscriptionID = core.StringPtr(d.Get("azure_subscription_id").(string))
	}
	if d.Get("azure_environment") != nil && d.Get("azure_environment") != "" {
		keystoreBody.AzureEnvironment = core.StringPtr(d.Get("azure_environment").(string))
	}
	if d.Get("ibm_api_endpoint") != nil && d.Get("ibm_api_endpoint") != "" {
		keystoreBody.IbmApiEndpoint = core.StringPtr(d.Get("ibm_api_endpoint").(string))
	}
	if d.Get("ibm_iam_endpoint") != nil && d.Get("ibm_iam_endpoint") != "" {
		keystoreBody.IbmIamEndpoint = core.StringPtr(d.Get("ibm_iam_endpoint").(string))
	}
	if d.Get("ibm_api_key") != nil && d.Get("ibm_api_key") != "" {
		keystoreBody.IbmApiKey = core.StringPtr(d.Get("ibm_api_key").(string))
	}
	if d.Get("ibm_instance_id") != nil && d.Get("ibm_instance_id") != "" {
		keystoreBody.IbmInstanceID = core.StringPtr(d.Get("ibm_instance_id").(string))
	}
	if d.Get("ibm_key_ring") != nil && d.Get("ibm_key_ring") != "" {
		keystoreBody.IbmKeyRing = core.StringPtr(d.Get("ibm_key_ring").(string))
	}

	return &keystoreBody
}

func DHasChanges(d *schema.ResourceData) bool {
	if d.HasChange("type") {
		return true
	}
	if d.HasChange("vault") {
		return true
	}
	if d.HasChange("description") {
		return true
	}
	if d.HasChange("groups") {
		return true
	}
	if d.HasChange("aws_region") {
		return true
	}
	if d.HasChange("aws_access_key_id") {
		return true
	}
	if d.HasChange("aws_secret_access_key") {
		return true
	}
	if d.HasChange("azure_service_name") {
		return true
	}
	if d.HasChange("azure_resource_group") {
		return true
	}
	if d.HasChange("azure_location") {
		return true
	}
	if d.HasChange("azure_service_principal_client_id") {
		return true
	}
	if d.HasChange("azure_service_principal_password") {
		return true
	}
	if d.HasChange("azure_tenant") {
		return true
	}
	if d.HasChange("azure_subscription_id") {
		return true
	}
	if d.HasChange("azure_environment") {
		return true
	}
	if d.HasChange("ibm_variant") {
		return true
	}
	if d.HasChange("ibm_api_endpoint") {
		return true
	}
	if d.HasChange("ibm_iam_endpoint") {
		return true
	}
	if d.HasChange("ibm_api_key") {
		return true
	}
	if d.HasChange("ibm_instance_id") {
		return true
	}
	if d.HasChange("ibm_key_ring") {
		return true
	}
	return false
}
