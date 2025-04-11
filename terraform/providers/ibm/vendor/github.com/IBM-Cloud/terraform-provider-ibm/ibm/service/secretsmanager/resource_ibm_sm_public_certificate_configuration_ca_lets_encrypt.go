// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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
			"lets_encrypt_preferred_chain": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefer the chain with an issuer matching this Subject Common Name.",
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
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigCALetsEncryptResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmPublicCertificateConfigurationCALetsEncryptMapToConfigurationPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigCALetsEncryptResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PublicCertConfigCALetsEncryptResourceName, "create")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncrypt)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a CA configuration use the format `<region>/<instance_id>/<name>`", PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncrypt)

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("config_type", configuration.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("lets_encrypt_environment", configuration.LetsEncryptEnvironment); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lets_encrypt_environment"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("lets_encrypt_preferred_chain", configuration.LetsEncryptPreferredChain); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lets_encrypt_preferred_chain"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("lets_encrypt_private_key", configuration.LetsEncryptPrivateKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting lets_encrypt_private_key"), PublicCertConfigCALetsEncryptResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigCALetsEncryptResourceName, "update")
		return tfErr.GetDiag()
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

	if d.HasChange("lets_encrypt_preferred_chain") {
		patchVals.LetsEncryptPreferredChain = core.StringPtr(d.Get("lets_encrypt_preferred_chain").(string))
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
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigCALetsEncryptResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigCALetsEncryptResourceName, "delete")
		return tfErr.GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigCALetsEncryptResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPublicCertificateConfigurationCALetsEncryptMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.PublicCertificateConfigurationCALetsEncryptPrototype{}

	model.ConfigType = core.StringPtr("public_cert_configuration_ca_lets_encrypt")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("lets_encrypt_environment"); ok {
		model.LetsEncryptEnvironment = core.StringPtr(d.Get("lets_encrypt_environment").(string))
	}
	if _, ok := d.GetOk("lets_encrypt_preferred_chain"); ok {
		model.LetsEncryptPreferredChain = core.StringPtr(d.Get("lets_encrypt_preferred_chain").(string))
	}
	if _, ok := d.GetOk("lets_encrypt_private_key"); ok {
		model.LetsEncryptPrivateKey = core.StringPtr(formatCertificate(d.Get("lets_encrypt_private_key").(string)))
	}
	return model, nil
}
