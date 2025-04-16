// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureCreate,
		ReadContext:   resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead,
		UpdateContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureUpdate,
		DeleteContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable unique name to assign to your configuration.To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
			},
			"classic_infrastructure_username": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_public_certificate_configuration_dns_classic_infrastructure", "classic_infrastructure_username"),
				Description:  "The username that is associated with your classic infrastructure account.In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
			},
			"classic_infrastructure_password": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_public_certificate_configuration_dns_classic_infrastructure", "classic_infrastructure_password"),
				Description:  "Your classic infrastructure API key.For information about viewing and accessing your classic infrastructure API key, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
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
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
		},
	}
}

func ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "classic_infrastructure_username",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             2,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "classic_infrastructure_password",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             2,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_sm_public_certificate_configuration_dns_classic_infrastructure", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	bodyModelMap := map[string]interface{}{}
	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	if _, ok := d.GetOk("config_type"); ok {
		bodyModelMap["config_type"] = d.Get("config_type")
	}
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	if _, ok := d.GetOk("classic_infrastructure_username"); ok {
		bodyModelMap["classic_infrastructure_username"] = d.Get("classic_infrastructure_username")
	}
	if _, ok := d.GetOk("classic_infrastructure_password"); ok {
		bodyModelMap["classic_infrastructure_password"] = d.Get("classic_infrastructure_password")
	}
	convertedModel, err := resourceIbmSmConfigurationPublicCertificateClassicInfrastructureMapToPublicCertificateConfigurationDNSClassicInfrastructurePrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createConfigurationOptions.SetConfigurationPrototype(convertedModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateConfigurationWithContext failed %s\n%s", err, response))
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructure)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		return diag.Errorf("Wrong format of resource ID. To import a DNS configuration use the format `<region>/<instance_id>/<name>`")
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

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructure)
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if !core.IsNil(configuration.ConfigType) {
		if err = d.Set("config_type", configuration.ConfigType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting config_type: %s", err))
		}
	}
	if !core.IsNil(configuration.ClassicInfrastructureUsername) {
		if err = d.Set("classic_infrastructure_username", configuration.ClassicInfrastructureUsername); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting classic_infrastructure_username: %s", err))
		}
	}
	if !core.IsNil(configuration.ClassicInfrastructurePassword) {
		if err = d.Set("classic_infrastructure_password", configuration.ClassicInfrastructurePassword); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting classic_infrastructure_password: %s", err))
		}
	}
	if err = d.Set("name", configuration.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config name: %s", err))
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}
	if err = d.Set("created_by", configuration.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("created_at", DateTimeToRFC3339(configuration.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(configuration.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}
	if d.HasChange("classic_infrastructure_username") {
		newClassicInfrastructureUsername := d.Get("classic_infrastructure_username").(string)
		patchVals.ClassicInfrastructureUsername = &newClassicInfrastructureUsername
		hasChange = true
	}
	if d.HasChange("classic_infrastructure_password") {
		newClassicInfrastructurePassword := d.Get("classic_infrastructure_password").(string)
		patchVals.ClassicInfrastructurePassword = &newClassicInfrastructurePassword
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

	return resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmSmConfigurationPublicCertificateClassicInfrastructureMapToPublicCertificateConfigurationDNSClassicInfrastructurePrototype(modelMap map[string]interface{}) (*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructurePrototype, error) {
	model := &secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructurePrototype{}
	model.ConfigType = core.StringPtr("public_cert_configuration_dns_classic_infrastructure")
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.ClassicInfrastructureUsername = core.StringPtr(modelMap["classic_infrastructure_username"].(string))
	model.ClassicInfrastructurePassword = core.StringPtr(modelMap["classic_infrastructure_password"].(string))
	return model, nil
}
