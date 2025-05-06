// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPublicCertificateConfigurationCALetsEncrypt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPublicCertificateConfigurationCALetsEncryptCreate,
		ReadContext:   resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead,
		UpdateContext: resourceIbmSmPublicCertificateConfigurationCALetsEncryptUpdate,
		DeleteContext: resourceIbmSmPublicCertificateConfigurationCALetsEncryptDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable unique name to assign to your configuration.To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.",
			},
			"lets_encrypt_environment": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The configuration of the Let's Encrypt CA environment.",
			},
			"lets_encrypt_private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The PEM encoded private key of your Lets Encrypt account.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if removeNewLineFromCertificate(oldValue) == removeNewLineFromCertificate(newValue) {
						return true
					}
					return false
				},
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
		},
	}
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmPublicCertificateConfigurationCALetsEncryptMapToConfigurationPrototype(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateConfigurationWithContext failed %s\n%s", err, response))
	}
	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncrypt)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		return diag.Errorf("Wrong format of resource ID. To import a CA configuration use the format `<region>/<instance_id>/<name>`")
	}
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(configName)

	configurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigurationWithContext failed %s\n%s", err, response))
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncrypt)

	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("config_type", configuration.ConfigType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config_type: %s", err))
	}
	if err = d.Set("name", configuration.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config name: %s", err))
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}
	if err = d.Set("lets_encrypt_environment", configuration.LetsEncryptEnvironment); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lets_encrypt_environment: %s", err))
	}
	if err = d.Set("lets_encrypt_private_key", configuration.LetsEncryptPrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lets_encrypt_private_key: %s", err))
	}

	return nil
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)
	updateConfigurationOptions.SetXSmAcceptConfigurationType("public_cert_configuration_ca_lets_encrypt")

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}

	if d.HasChange("lets_encrypt_private_key") {
		patchVals.LetsEncryptPrivateKey = core.StringPtr(d.Get("lets_encrypt_private_key").(string))
		hasChange = true
	}

	patchVals.LetsEncryptEnvironment = core.StringPtr(d.Get("lets_encrypt_environment").(string))
	if d.HasChange("lets_encrypt_environment") {
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateConfigurationWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteConfigurationOptions := &secretsmanagerv2.DeleteConfigurationOptions{}

	deleteConfigurationOptions.SetName(configName)

	response, err := secretsManagerClient.DeleteConfigurationWithContext(context, deleteConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteConfigurationWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.PublicCertificateConfigurationCALetsEncryptPrototype{}

	model.ConfigType = core.StringPtr("public_cert_configuration_ca_lets_encrypt")

	//if _, ok := d.GetOk("config_type"); ok {
	//	model.ConfigType = core.StringPtr(d.Get("config_type").(string))
	//}
	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("lets_encrypt_environment"); ok {
		model.LetsEncryptEnvironment = core.StringPtr(d.Get("lets_encrypt_environment").(string))
	}
	if _, ok := d.GetOk("lets_encrypt_private_key"); ok {
		model.LetsEncryptPrivateKey = core.StringPtr(formatCertificate(d.Get("lets_encrypt_private_key").(string)))
	}
	return model, nil
}
