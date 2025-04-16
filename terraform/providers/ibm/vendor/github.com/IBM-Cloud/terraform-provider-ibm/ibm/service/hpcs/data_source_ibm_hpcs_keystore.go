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

func DataSourceIbmKeystore() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIbmKeystoreRead,

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
			"keystore_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the keystore.",
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the target keystore. It can be changed in the future.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Geographic location of the keystore, if available.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the keystore.",
			},
			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of groups that this keystore belongs to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
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
			"google_credentials": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the JSON key represented in the Base64 format.",
			},
			"google_location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location represents the geographical region where a Cloud KMS resource is stored and can be accessed. A key's location impacts the performance of applications using the key.",
			},
			"google_project_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project id associated with this keystore.",
			},
			"google_private_key_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private key id associated with this keystore.",
			},
			"google_key_ring": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of keys.",
			},
			"aws_region": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Region.",
			},
			"aws_access_key_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The access key id used for connecting to this instance of AWS KMS.",
			},
			"aws_secret_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret access key used for connecting to this instance of AWS KMS.",
			},
			"azure_service_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name of the key vault instance from the Azure portal.",
			},
			"azure_resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group in Azure.",
			},
			"azure_location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of the Azure Key Vault.",
			},
			"azure_service_principal_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure service principal client ID.",
			},
			"azure_service_principal_password": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Azure service principal password.",
			},
			"azure_tenant": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure tenant that the Key Vault is associated with,.",
			},
			"azure_subscription_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription ID in Azure.",
			},
			"azure_environment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure environment, usually 'Azure'.",
			},
			"ibm_api_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API endpoint of the IBM Cloud keystore.",
			},
			"ibm_iam_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint of the IAM service for this IBM Cloud keystore.",
			},
			"ibm_api_key": &schema.Schema{ // pragma: allowlist secret
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The IBM Cloud API key to be used for connecting to this IBM Cloud keystore.",
			},
			"ibm_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID of the IBM Cloud keystore.",
			},
			"ibm_variant": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Possible IBM Cloud KMS variants.",
			},
			"ibm_key_ring": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key ring of an IBM Cloud KMS Keystore.",
			},
		},
	}
}

func DataSourceIbmKeystoreRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getKeystoreOptions := &ukov4.GetKeystoreOptions{}

	region := d.Get("region").(string)
	instance_id := d.Get("instance_id").(string)
	vault_id := d.Get("uko_vault").(string)
	keystore_id := d.Get("keystore_id").(string)
	getKeystoreOptions.SetID(keystore_id)
	getKeystoreOptions.SetUKOVault(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	keystoreIntf, response, err := ukoClient.GetKeystoreWithContext(context, getKeystoreOptions)
	keystore := keystoreIntf.(*ukov4.Keystore)
	if err != nil {
		log.Printf("[DEBUG] GetKeystoreWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetKeystoreWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, vault_id, *getKeystoreOptions.ID))

	vault := []map[string]interface{}{}
	if keystore.Vault != nil {
		modelMap, err := DataSourceIbmKeystoreVaultReferenceToMap(keystore.Vault)
		if err != nil {
			return diag.FromErr(err)
		}
		vault = append(vault, modelMap)
	}
	if err = d.Set("vault", vault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault %s", err))
	}

	if err = d.Set("name", keystore.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("location", keystore.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}

	if err = d.Set("description", keystore.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
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

	if err = d.Set("google_credentials", keystore.GoogleCredentials); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting google_credentials: %s", err))
	}

	if err = d.Set("google_location", keystore.GoogleLocation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting google_location: %s", err))
	}

	if err = d.Set("google_project_id", keystore.GoogleProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting google_project_id: %s", err))
	}

	if err = d.Set("google_private_key_id", keystore.GooglePrivateKeyID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting google_private_key_id: %s", err))
	}

	if err = d.Set("google_key_ring", keystore.GoogleKeyRing); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting google_key_ring: %s", err))
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

	return nil
}

func DataSourceIbmKeystoreVaultReferenceToMap(model *ukov4.VaultReference) (map[string]interface{}, error) {
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
